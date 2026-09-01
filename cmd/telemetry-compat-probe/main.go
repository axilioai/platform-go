// Command telemetry-compat-probe validates AXI-1982 against a candidate Go
// SDK. Live checks are read-only and limited to dev/staging. Future frame
// kinds are tested only against a loopback capture/replay server.
package main

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"time"

	platformgo "github.com/axilioai/platform-go"
	"github.com/axilioai/platform-go/client"
	"github.com/axilioai/platform-go/drivers/telemetry"
	"github.com/axilioai/platform-go/option"
)

const apiKeyEnv = "AXILIO_API_KEY"

type fixture struct {
	ID        string `json:"id"`
	MinFrames int    `json:"min_frames,omitempty"`
}

type fixtureManifest struct {
	ManifestVersion int    `json:"manifest_version"`
	Environment     string `json:"environment"`
	SeedRevision    string `json:"seed_revision"`
	Fixtures        struct {
		NormalSession      fixture `json:"normal_session"`
		NormalEmptySession fixture `json:"normal_empty_session"`
		ExpiredSession     fixture `json:"expired_session"`
	} `json:"fixtures"`
}

type result struct {
	ID               string `json:"id"`
	Status           string `json:"status"`
	Classification   string `json:"classification"`
	StartedAt        string `json:"started_at"`
	DurationMillis   int64  `json:"duration_ms"`
	Attempts         int    `json:"attempts"`
	Expected         string `json:"expected"`
	ObservedRedacted string `json:"observed_redacted"`
	EvidenceFile     string `json:"evidence_file"`
}

type recorder struct {
	Results []result
}

func (r *recorder) check(id, expected string, fn func() (string, error)) {
	started := time.Now().UTC()
	observed, err := fn()
	status, classification := "PASS", "PASS"
	if err != nil {
		status, classification = "FAIL", "SDK_ARTIFACT_FAILURE"
		observed = redactedFailure(err)
	}
	r.Results = append(r.Results, result{
		ID: id, Status: status, Classification: classification,
		StartedAt: started.Format(time.RFC3339), DurationMillis: time.Since(started).Milliseconds(),
		Attempts: 1, Expected: expected, ObservedRedacted: observed,
	})
}

func redactedFailure(err error) string {
	fingerprint := sha256.Sum256([]byte(err.Error()))
	return fmt.Sprintf("error_type=%T fingerprint=%x", err, fingerprint)
}

func (r *recorder) passed() bool {
	for _, check := range r.Results {
		if check.Status != "PASS" {
			return false
		}
	}
	return true
}

type output struct {
	Environment           string   `json:"environment"`
	SDK                   string   `json:"sdk"`
	SDKRef                string   `json:"sdk_ref"`
	ArtifactSHA256        string   `json:"artifact_sha256"`
	SeedRevision          string   `json:"seed_revision"`
	FixtureManifestSHA256 string   `json:"fixture_manifest_sha256"`
	Verdict               string   `json:"verdict"`
	Results               []result `json:"results"`
}

func main() {
	var environment, baseURL, approvedDevOrigin, manifestPath, outputPath, sdkRef, artifactSHA string
	var replayOnly bool
	flag.StringVar(&environment, "env", "", "target environment: dev or staging")
	flag.StringVar(&baseURL, "base-url", "", "target API origin, without /api/v1")
	flag.StringVar(&approvedDevOrigin, "approved-dev-origin", "", "approved non-loopback dev API origin")
	flag.StringVar(&manifestPath, "fixture-manifest", "", "telemetry-compat fixture manifest")
	flag.StringVar(&outputPath, "output", "", "redacted JSON evidence output")
	flag.StringVar(&sdkRef, "sdk-ref", "", "exact candidate Go SDK git ref")
	flag.StringVar(&artifactSHA, "artifact-sha256", "", "SHA-256 of the tested candidate artifact")
	flag.BoolVar(&replayOnly, "replay-only", false, "skip deployed read-only checks")
	flag.Parse()

	if err := execute(context.Background(), environment, baseURL, approvedDevOrigin, manifestPath, outputPath, sdkRef, artifactSHA, replayOnly, os.Getenv(apiKeyEnv)); err != nil {
		fmt.Fprintf(os.Stderr, "telemetry-compat-probe: %v\n", err)
		os.Exit(1)
	}
}

func execute(ctx context.Context, environment, baseURL, approvedDevOrigin, manifestPath, outputPath, sdkRef, artifactSHA string, replayOnly bool, apiKey string) error {
	origin, err := validateTarget(environment, baseURL, approvedDevOrigin)
	if err != nil {
		return err
	}
	if outputPath == "" || sdkRef == "" || artifactSHA == "" {
		return errors.New("--output, --sdk-ref, and --artifact-sha256 are required")
	}
	if err := validateCandidateProvenance(sdkRef, artifactSHA); err != nil {
		return err
	}
	manifest, fixtureManifestSHA256, err := loadManifest(manifestPath, environment)
	if err != nil {
		return err
	}
	if !replayOnly && apiKey == "" {
		return errors.New("AXILIO_API_KEY is required for live validation")
	}

	rec := new(recorder)
	if !replayOnly {
		runLive(ctx, origin, apiKey, manifest, sdkRef, artifactSHA, rec)
	}
	runReplay(ctx, rec)

	verdict := "PASS"
	if !rec.passed() {
		verdict = "SDK_ARTIFACT_FAILURE"
	}
	for index := range rec.Results {
		rec.Results[index].EvidenceFile = filepath.Base(outputPath)
	}
	if err := writeOutput(outputPath, output{
		Environment: environment, SDK: "go", SDKRef: sdkRef,
		ArtifactSHA256: artifactSHA, SeedRevision: manifest.SeedRevision,
		FixtureManifestSHA256: fixtureManifestSHA256,
		Verdict:               verdict, Results: rec.Results,
	}); err != nil {
		return err
	}
	if verdict != "PASS" {
		return fmt.Errorf("validation verdict %s; see %s", verdict, outputPath)
	}
	fmt.Printf("AXI-1982 Go validation PASS; evidence %s\n", outputPath)
	return nil
}

func validateCandidateProvenance(sdkRef, artifactSHA string) error {
	if !isLowerHex(sdkRef, 40) {
		return errors.New("--sdk-ref must be a full lowercase 40-hex Git commit SHA")
	}
	if !isLowerHex(artifactSHA, 64) {
		return errors.New("--artifact-sha256 must be a full lowercase 64-hex SHA-256 digest")
	}
	return nil
}

func isLowerHex(value string, width int) bool {
	if len(value) != width {
		return false
	}
	for i := range value {
		if (value[i] < '0' || value[i] > '9') && (value[i] < 'a' || value[i] > 'f') {
			return false
		}
	}
	return true
}

func validateTarget(environment, rawBaseURL, approvedDevOrigin string) (string, error) {
	if environment != "dev" && environment != "staging" {
		return "", errors.New("environment must be dev or staging; production is refused")
	}
	rawBaseURL = strings.TrimRight(rawBaseURL, "/")
	parsed, err := url.Parse(rawBaseURL)
	if err != nil || parsed.Hostname() == "" {
		return "", errors.New("base URL must be an absolute origin")
	}
	hostname := strings.TrimRight(strings.ToLower(parsed.Hostname()), ".")
	if hostname == "api.axilio.ai" {
		return "", errors.New("production origin is refused")
	}
	if parsed.User != nil || parsed.RawQuery != "" || parsed.Fragment != "" || (parsed.Path != "" && parsed.Path != "/") {
		return "", errors.New("base URL must be an origin without credentials, path, query, or fragment")
	}
	if environment == "staging" && rawBaseURL != "https://staging-api.axilio.ai" {
		return "", errors.New("staging must use exactly https://staging-api.axilio.ai")
	}
	loopback := hostname == "127.0.0.1" || hostname == "localhost" || hostname == "::1"
	if parsed.Scheme != "https" && !(environment == "dev" && parsed.Scheme == "http" && loopback) {
		return "", errors.New("HTTPS is required; only loopback dev may use HTTP")
	}
	if environment == "dev" && !loopback {
		if approvedDevOrigin == "" {
			return "", errors.New("non-loopback dev requires --approved-dev-origin")
		}
		approvedDevOrigin = strings.TrimRight(approvedDevOrigin, "/")
		approved, approvedErr := url.Parse(approvedDevOrigin)
		if approvedErr != nil || approved == nil {
			return "", errors.New("approved dev origin must be an HTTPS origin")
		}
		approvedHostname := strings.TrimRight(strings.ToLower(approved.Hostname()), ".")
		if approved.Scheme != "https" || approvedHostname == "" || approved.User != nil || approved.RawQuery != "" || approved.Fragment != "" || (approved.Path != "" && approved.Path != "/") {
			return "", errors.New("approved dev origin must be an HTTPS origin")
		}
		if approvedHostname == "api.axilio.ai" {
			return "", errors.New("production origin is refused")
		}
		if rawBaseURL != approvedDevOrigin {
			return "", errors.New("dev base URL does not match approved dev origin")
		}
	}
	return rawBaseURL, nil
}

func loadManifest(path, environment string) (fixtureManifest, string, error) {
	var manifest fixtureManifest
	if path == "" {
		return manifest, "", errors.New("--fixture-manifest is required")
	}
	body, err := os.ReadFile(path)
	if err != nil {
		return manifest, "", fmt.Errorf("read fixture manifest: %w", err)
	}
	if err := json.Unmarshal(body, &manifest); err != nil {
		return manifest, "", fmt.Errorf("parse fixture manifest: %w", err)
	}
	if manifest.ManifestVersion != 1 || manifest.Environment != environment {
		return manifest, "", errors.New("fixture manifest version/environment mismatch")
	}
	if manifest.SeedRevision == "" || manifest.SeedRevision == "unknown" {
		return manifest, "", errors.New("fixture manifest provenance missing")
	}
	if manifest.Fixtures.NormalSession.ID == "" || manifest.Fixtures.NormalEmptySession.ID == "" || manifest.Fixtures.ExpiredSession.ID == "" {
		return manifest, "", errors.New("fixture manifest is missing a required session")
	}
	digest := sha256.Sum256(body)
	return manifest, fmt.Sprintf("%x", digest), nil
}

func newClient(baseURL, apiKey string) *client.Client {
	httpClient := &http.Client{
		Timeout:       10 * time.Second,
		CheckRedirect: func(_ *http.Request, _ []*http.Request) error { return http.ErrUseLastResponse },
	}
	return client.NewClient(
		option.WithBaseURL(strings.TrimRight(baseURL, "/")+"/api/v1"),
		option.WithAPIKey(apiKey),
		option.WithHTTPClient(httpClient),
		option.WithoutRetries(),
	)
}

func getPage(ctx context.Context, c *client.Client, sessionID string, limit, offset int64) (*platformgo.RunSessionFramesResponse, error) {
	request := new(platformgo.SessionsListFramesRequest)
	request.SetSessionID(sessionID)
	request.SetLimit(platformgo.Int64(limit))
	request.SetOffset(platformgo.Int64(offset))
	return c.Runs.SessionsListFrames(ctx, request)
}

func validateNormalEmptyPage(ctx context.Context, c *client.Client, sessionID string) (string, error) {
	page, err := getPage(ctx, c, sessionID, 7, 3)
	if err != nil {
		return "", err
	}
	if page.RetentionExpired || page.Frames == nil || len(page.Frames) != 0 || page.Total != 0 || page.Limit != 7 || page.Offset != 3 || page.SdkCallCosts == nil || len(page.SdkCallCosts) != 0 || page.InferenceCosts == nil || len(page.InferenceCosts) != 0 {
		return "", errors.New("normal-empty page shape mismatch")
	}
	return "retention=false frames=[] total=0 maps=object-empty pagination=7/3", nil
}

func runLive(ctx context.Context, baseURL, apiKey string, manifest fixtureManifest, sdkRef, artifactSHA string, rec *recorder) {
	c := newClient(baseURL, apiKey)
	rec.check("GO-01", "generated client parses the normal page as today's known frames", func() (string, error) {
		page, err := getPage(ctx, c, manifest.Fixtures.NormalSession.ID, 100, 0)
		if err != nil {
			return "", err
		}
		if page.RetentionExpired || len(page.Frames) < max(1, manifest.Fixtures.NormalSession.MinFrames) || page.SdkCallCosts == nil || page.InferenceCosts == nil {
			return "", errors.New("normal page shape mismatch")
		}
		for _, frame := range page.Frames {
			if frame.GetSpan() == nil && frame.GetLog() == nil {
				return "", fmt.Errorf("today's backend returned unknown kind %q", frame.GetKind())
			}
		}
		return fmt.Sprintf("frames=%d maps=object", len(page.Frames)), nil
	})
	rec.check("GO-02", "generated client accepts expired null cost maps as nil", func() (string, error) {
		page, err := getPage(ctx, c, manifest.Fixtures.ExpiredSession.ID, 7, 3)
		if err != nil {
			return "", err
		}
		if !page.RetentionExpired || len(page.Frames) != 0 || page.Total != 0 || page.SdkCallCosts != nil || page.InferenceCosts != nil || page.Limit != 7 || page.Offset != 3 {
			return "", errors.New("expired page shape mismatch")
		}
		return "retention=true frames=[] maps=nil pagination=7/3", nil
	})
	rec.check("GO-03", "high-level expired trace is empty without error", func() (string, error) {
		trace, err := telemetry.NewSession(c, manifest.Fixtures.ExpiredSession.ID).Trace(ctx)
		if err != nil {
			return "", err
		}
		if !trace.RetentionExpired || len(trace.Spans)+len(trace.Logs)+len(trace.Unknown) != 0 {
			return "", errors.New("expired high-level trace mismatch")
		}
		return "retention=true spans=0 logs=0 unknown=0", nil
	})
	rec.check("GO-04", "high-level normal trace keeps every known frame", func() (string, error) {
		page, err := getPage(ctx, c, manifest.Fixtures.NormalSession.ID, 1000, 0)
		if err != nil {
			return "", err
		}
		if page.Total != int64(len(page.Frames)) {
			return "", errors.New("normal fixture no longer fits one raw page")
		}
		rawSpans, rawLogs := 0, 0
		for _, frame := range page.Frames {
			if frame.GetSpan() != nil {
				rawSpans++
			}
			if frame.GetLog() != nil {
				rawLogs++
			}
		}
		trace, err := telemetry.NewSession(c, manifest.Fixtures.NormalSession.ID).Trace(ctx)
		if err != nil {
			return "", err
		}
		if trace.RetentionExpired || len(trace.Spans) != rawSpans || len(trace.Logs) != rawLogs || len(trace.Unknown) != 0 {
			return "", errors.New("normal high-level trace mismatch")
		}
		return fmt.Sprintf("spans=%d logs=%d unknown=0", len(trace.Spans), len(trace.Logs)), nil
	})
	rec.check("GO-05", "candidate source/artifact provenance is canonical lowercase full-length hex", func() (string, error) {
		if err := validateCandidateProvenance(sdkRef, artifactSHA); err != nil {
			return "", err
		}
		return fmt.Sprintf("sdk_ref=%s artifact_sha256=%s", sdkRef, artifactSHA), nil
	})
	rec.check("GO-06", "normal empty session is non-expired with empty arrays and object cost maps", func() (string, error) {
		return validateNormalEmptyPage(ctx, c, manifest.Fixtures.NormalEmptySession.ID)
	})
}

var metricFrame = map[string]any{
	"kind": "metric", "name": "axi.e2e.synthetic", "value": 0.72,
	"nested": map[string]any{"unit": "ratio"}, "tags": []any{"a", "b"},
}

func spanFrame() map[string]any {
	return map[string]any{
		"kind": "span", "phase": "end", "span_type": "sdk_call",
		"trace_id": strings.Repeat("0", 32), "span_id": strings.Repeat("1", 16),
		"name": "Screen.observe", "start_time_unix_nano": 1, "end_time_unix_nano": 2,
		"status": map[string]any{"code": "ok", "message": ""}, "attributes": map[string]any{},
	}
}

func logFrame(index int) map[string]any {
	return map[string]any{
		"kind": "log", "log_type": "output_log", "trace_id": strings.Repeat("0", 32),
		"span_id": strings.Repeat("1", 16), "time_unix_nano": 10 + index,
		"severity": "INFO", "body": "synthetic", "attributes": map[string]any{},
	}
}

func replayResponse(frames []any, total, limit, offset int, retention, nullMaps bool) map[string]any {
	var costs any = map[string]int64{}
	if nullMaps {
		costs = nil
	}
	return map[string]any{
		"frames": frames, "total": total, "limit": limit, "offset": offset,
		"retention_expired": retention, "sdk_call_costs": costs, "inference_costs": costs,
	}
}

func replayHandler(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(parts) != 6 || strings.Join(parts[:4], "/") != "api/v1/phones/sessions" || parts[5] != "frames" {
		http.NotFound(w, r)
		return
	}
	sessionID := parts[4]
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit == 0 {
		limit = 100
	}
	var payload map[string]any
	switch {
	case sessionID == "mixed":
		payload = replayResponse([]any{spanFrame(), metricFrame, logFrame(0)}, 3, limit, offset, false, false)
	case sessionID == "extra-known":
		span := spanFrame()
		span["span_type"], span["future_field"] = "future_role", map[string]any{"x": 1}
		log := logFrame(0)
		log["log_type"] = "future_log"
		payload = replayResponse([]any{span, log}, 2, limit, offset, false, false)
	case sessionID == "malformed-span":
		payload = replayResponse([]any{map[string]any{"kind": "span", "phase": "end"}}, 1, limit, offset, false, false)
	case sessionID == "malformed-log":
		payload = replayResponse([]any{map[string]any{"kind": "log", "body": "synthetic"}}, 1, limit, offset, false, false)
	case strings.HasPrefix(sessionID, "bad-kind-"):
		name := strings.TrimPrefix(sessionID, "bad-kind-")
		badKinds := map[string]any{"empty": "", "null": nil, "number": 7, "object": map[string]any{"x": 1}, "array": []any{"metric"}}
		frame := map[string]any{}
		if name != "missing" {
			frame["kind"] = badKinds[name]
		}
		payload = replayResponse([]any{frame}, 1, limit, offset, false, false)
	case sessionID == "expired":
		payload = replayResponse([]any{}, 0, limit, offset, true, true)
	case sessionID == "paged" && offset == 0:
		frames := make([]any, 0, 1000)
		span := spanFrame()
		frames = append(frames, span)
		for i := 1; i < 1000; i++ {
			frames = append(frames, logFrame(i))
		}
		payload = replayResponse(frames, 1001, 1000, 0, false, false)
		payload["sdk_call_costs"] = map[string]int64{strings.Repeat("1", 16): 12}
	case sessionID == "paged" && offset == 1000:
		payload = replayResponse([]any{metricFrame}, 1001, 1000, 1000, false, true)
	default:
		http.Error(w, `{"error":"unexpected replay request"}`, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(payload)
}

func runReplay(ctx context.Context, rec *recorder) {
	server := httptest.NewServer(http.HandlerFunc(replayHandler))
	defer server.Close()
	c := newClient(server.URL, "axl_loopback")

	rec.check("GO-REPLAY-01", "future metric survives between typed known siblings", func() (string, error) {
		page, err := getPage(ctx, c, "mixed", 100, 0)
		if err != nil {
			return "", err
		}
		if len(page.Frames) != 3 || page.Frames[0].GetSpan() == nil || page.Frames[1].GetKind() != "metric" || page.Frames[1].GetSpan() != nil || page.Frames[1].GetLog() != nil || page.Frames[2].GetLog() == nil {
			return "", errors.New("mixed generated classification mismatch")
		}
		trace, err := telemetry.NewSession(c, "mixed").Trace(ctx)
		if err != nil || len(trace.Spans) != 1 || len(trace.Logs) != 1 || len(trace.Unknown) != 1 {
			return "", errors.New("mixed high-level classification mismatch")
		}
		unknownBody, err := json.Marshal(trace.Unknown[0])
		if err != nil {
			return "", err
		}
		var unknownPayload any
		if err := json.Unmarshal(unknownBody, &unknownPayload); err != nil {
			return "", errors.New("high-level unknown is not JSON")
		}
		metricBody, _ := json.Marshal(metricFrame)
		var metricPayload any
		_ = json.Unmarshal(metricBody, &metricPayload)
		if !reflect.DeepEqual(unknownPayload, metricPayload) {
			return "", errors.New("high-level unknown payload changed")
		}
		return "generated=span,unknown,log high=1/1/1", nil
	})
	rec.check("GO-REPLAY-02", "new fields and role strings remain additive inside known kinds", func() (string, error) {
		page, err := getPage(ctx, c, "extra-known", 100, 0)
		if err != nil {
			return "", err
		}
		if len(page.Frames) != 2 || page.Frames[0].GetSpan() == nil || page.Frames[1].GetLog() == nil || page.Frames[0].GetSpan().GetSpanType() != "future_role" || page.Frames[1].GetLog().GetLogType() != "future_log" {
			return "", errors.New("future known values were not retained")
		}
		if _, ok := page.Frames[0].GetSpan().GetExtraProperties()["future_field"]; !ok {
			return "", errors.New("future known field was dropped")
		}
		return "known models retained extra field and future role strings", nil
	})
	rec.check("GO-REPLAY-03", "malformed known kind remains known; Go does not promise required-field validation", func() (string, error) {
		spanPage, err := getPage(ctx, c, "malformed-span", 100, 0)
		if err != nil {
			return "", err
		}
		logPage, err := getPage(ctx, c, "malformed-log", 100, 0)
		if err != nil {
			return "", err
		}
		if len(spanPage.Frames) != 1 || spanPage.Frames[0].GetSpan() == nil || spanPage.Frames[0].GetKind() != "span" {
			return "", errors.New("malformed known span fell through to unknown")
		}
		if len(logPage.Frames) != 1 || logPage.Frames[0].GetLog() == nil || logPage.Frames[0].GetKind() != "log" {
			return "", errors.New("malformed known log fell through to unknown")
		}
		return "accepted as known span/log with zero-value missing fields (characterized, not guaranteed validation)", nil
	})
	rec.check("GO-REPLAY-04", "only non-empty unknown strings enter the fallback", func() (string, error) {
		for _, name := range []string{"missing", "empty", "null", "number", "object", "array"} {
			if _, err := getPage(ctx, c, "bad-kind-"+name, 100, 0); err == nil {
				return "", fmt.Errorf("bad kind %s parsed successfully", name)
			}
		}
		return "missing,empty,null,number,object,array rejected", nil
	})
	rec.check("GO-REPLAY-05", "expired body preserves generated/high-level null semantics", func() (string, error) {
		page, err := getPage(ctx, c, "expired", 7, 3)
		if err != nil {
			return "", err
		}
		trace, err := telemetry.NewSession(c, "expired").Trace(ctx)
		if err != nil {
			return "", err
		}
		if !page.RetentionExpired || len(page.Frames) != 0 || page.Total != 0 || page.Limit != 7 || page.Offset != 3 || page.SdkCallCosts != nil || page.InferenceCosts != nil {
			return "", errors.New("expired generated replay semantics mismatch")
		}
		if !trace.RetentionExpired || len(trace.Spans) != 0 || len(trace.Logs) != 0 || len(trace.Unknown) != 0 {
			return "", errors.New("expired replay semantics mismatch")
		}
		return "generated retention/empty/0/7/3/maps=nil high empty/expired", nil
	})
	rec.check("GO-REPLAY-06", "two-page aggregation keeps 1001 items and page-one cost", func() (string, error) {
		trace, err := telemetry.NewSession(c, "paged").Trace(ctx)
		if err != nil {
			return "", err
		}
		if len(trace.Spans) != 1 || len(trace.Logs) != 999 || len(trace.Unknown) != 1 || trace.Spans[0].BilledMicrodollars != 12 {
			return "", fmt.Errorf("paged aggregation mismatch spans=%d logs=%d unknown=%d", len(trace.Spans), len(trace.Logs), len(trace.Unknown))
		}
		return "items=1001 pages=2 known=1000 unknown=1 cost=12", nil
	})
	rec.check("GO-REPLAY-07", "unknown decode and re-marshal preserves semantic JSON", func() (string, error) {
		page, err := getPage(ctx, c, "mixed", 100, 0)
		if err != nil {
			return "", err
		}
		body, err := json.Marshal(page.Frames[1])
		if err != nil {
			return "", err
		}
		var got, want any
		if json.Unmarshal(body, &got) != nil {
			return "", errors.New("re-marshaled unknown is not JSON")
		}
		wantBody, _ := json.Marshal(metricFrame)
		_ = json.Unmarshal(wantBody, &want)
		if !reflect.DeepEqual(got, want) {
			return "", errors.New("unknown semantic JSON changed")
		}
		return "semantic-roundtrip=true", nil
	})
}

func writeOutput(path string, value output) error {
	if dir := filepath.Dir(path); dir != "." {
		if err := os.MkdirAll(dir, 0o700); err != nil {
			return fmt.Errorf("create evidence directory: %w", err)
		}
	}
	body, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return err
	}
	body = append(body, '\n')
	if err := os.WriteFile(path, body, 0o600); err != nil {
		return fmt.Errorf("write evidence: %w", err)
	}
	return nil
}
