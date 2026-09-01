package main

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestFailureEvidenceRedactsErrorDetails(t *testing.T) {
	got := redactedFailure(errors.New("GET /sessions/secret-id Authorization: secret-token"))
	if strings.Contains(got, "secret-id") || strings.Contains(got, "secret-token") {
		t.Fatalf("redacted failure leaked details: %s", got)
	}
	if !strings.Contains(got, "error_type=") || !strings.Contains(got, "fingerprint=") {
		t.Fatalf("redacted failure lacks safe diagnostics: %s", got)
	}
	if strings.Contains(got, "details") {
		t.Fatalf("redacted failure contains detail marker: %s", got)
	}
}

func TestReplayMatrixPasses(t *testing.T) {
	rec := new(recorder)
	runReplay(context.Background(), rec)
	if !rec.passed() {
		for _, result := range rec.Results {
			if result.Status != "PASS" {
				t.Errorf("%s: %s", result.ID, result.ObservedRedacted)
			}
		}
	}
	want := []string{
		"GO-REPLAY-01", "GO-REPLAY-02", "GO-REPLAY-03", "GO-REPLAY-04",
		"GO-REPLAY-05", "GO-REPLAY-06", "GO-REPLAY-07",
	}
	if len(rec.Results) != len(want) {
		t.Fatalf("checks = %d, want %d", len(rec.Results), len(want))
	}
	for i := range want {
		if rec.Results[i].ID != want[i] {
			t.Errorf("check[%d] = %s, want %s", i, rec.Results[i].ID, want[i])
		}
	}
}

func TestValidateCandidateProvenanceRequiresCanonicalDigests(t *testing.T) {
	validRef := strings.Repeat("a", 40)
	validArtifact := strings.Repeat("b", 64)
	tests := []struct {
		name, sdkRef, artifactSHA string
		wantErr                   bool
	}{
		{"canonical", validRef, validArtifact, false},
		{"placeholder ref", "candidate-sha", validArtifact, true},
		{"short ref", strings.Repeat("a", 39), validArtifact, true},
		{"uppercase ref", strings.Repeat("A", 40), validArtifact, true},
		{"non-hex ref", strings.Repeat("g", 40), validArtifact, true},
		{"placeholder artifact", validRef, "artifact-sha", true},
		{"short artifact", validRef, strings.Repeat("b", 63), true},
		{"uppercase artifact", validRef, strings.Repeat("B", 64), true},
		{"non-hex artifact", validRef, strings.Repeat("z", 64), true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateCandidateProvenance(tt.sdkRef, tt.artifactSHA)
			if (err != nil) != tt.wantErr {
				t.Fatalf("validateCandidateProvenance() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestExecuteRejectsPlaceholderProvenanceBeforeWritingEvidence(t *testing.T) {
	outputPath := filepath.Join(t.TempDir(), "go.json")
	err := execute(context.Background(), "dev", "http://127.0.0.1", "", "", outputPath, "candidate-sha", "artifact-sha", true, "")
	if err == nil || !strings.Contains(err.Error(), "--sdk-ref") {
		t.Fatalf("execute() error = %v, want strict sdk-ref error", err)
	}
	if _, statErr := os.Stat(outputPath); !os.IsNotExist(statErr) {
		t.Fatalf("evidence file was created for invalid provenance: %v", statErr)
	}
}

func TestValidateNormalEmptyPageRequiresExactShape(t *testing.T) {
	tests := []struct {
		name, body string
		wantErr    bool
	}{
		{
			"exact empty page",
			`{"frames":[],"total":0,"limit":7,"offset":3,"retention_expired":false,"sdk_call_costs":{},"inference_costs":{}}`,
			false,
		},
		{
			"null frames",
			`{"frames":null,"total":0,"limit":7,"offset":3,"retention_expired":false,"sdk_call_costs":{},"inference_costs":{}}`,
			true,
		},
		{
			"null maps",
			`{"frames":[],"total":0,"limit":7,"offset":3,"retention_expired":false,"sdk_call_costs":null,"inference_costs":null}`,
			true,
		},
		{
			"nonempty maps",
			`{"frames":[],"total":0,"limit":7,"offset":3,"retention_expired":false,"sdk_call_costs":{"span":1},"inference_costs":{"inference":1}}`,
			true,
		},
		{
			"retention expired",
			`{"frames":[],"total":0,"limit":7,"offset":3,"retention_expired":true,"sdk_call_costs":{},"inference_costs":{}}`,
			true,
		},
		{
			"wrong pagination",
			`{"frames":[],"total":0,"limit":100,"offset":0,"retention_expired":false,"sdk_call_costs":{},"inference_costs":{}}`,
			true,
		},
		{
			"nonempty frames",
			`{"frames":[{"kind":"metric","value":1}],"total":1,"limit":7,"offset":3,"retention_expired":false,"sdk_call_costs":{},"inference_costs":{}}`,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path != "/api/v1/phones/sessions/normal-empty/frames" || r.URL.Query().Get("limit") != "7" || r.URL.Query().Get("offset") != "3" {
					http.Error(w, `{"error":"unexpected request"}`, http.StatusBadRequest)
					return
				}
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write([]byte(tt.body))
			}))
			defer server.Close()

			_, err := validateNormalEmptyPage(context.Background(), newClient(server.URL, "axl_test"), "normal-empty")
			if (err != nil) != tt.wantErr {
				t.Fatalf("validateNormalEmptyPage() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRegenPreservesTelemetryCompatibilityProbe(t *testing.T) {
	repoRoot := filepath.Clean(filepath.Join("..", ".."))
	regen, err := os.ReadFile(filepath.Join(repoRoot, "scripts", "regen.sh"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(regen), "--exclude='cmd/telemetry-compat-probe'") {
		t.Fatal("scripts/regen.sh does not preserve cmd/telemetry-compat-probe")
	}
	workflow, err := os.ReadFile(filepath.Join(repoRoot, ".github", "workflows", "regen.yml"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(workflow), "cmd/telemetry-compat-probe/**") {
		t.Fatal("regen workflow add-paths does not include cmd/telemetry-compat-probe")
	}
}

func TestValidateTarget(t *testing.T) {
	tests := []struct {
		name, environment, baseURL string
		wantErr                    bool
	}{
		{"staging exact", "staging", "https://staging-api.axilio.ai", false},
		{"dev loopback", "dev", "http://127.0.0.1:8080", false},
		{"production origin", "staging", "https://api.axilio.ai", true},
		{"production environment", "production", "https://api.axilio.ai", true},
		{"arbitrary staging", "staging", "https://example.invalid", true},
		{"remote cleartext dev", "dev", "http://example.invalid", true},
		{"path", "dev", "https://dev.example.invalid/api/v1", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := validateTarget(tt.environment, tt.baseURL, "")
			if (err != nil) != tt.wantErr {
				t.Fatalf("err = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateTargetRefusesTrailingDotProductionAndArbitraryHTTPSDev(t *testing.T) {
	for _, target := range []string{"https://api.axilio.ai.", "https://example.invalid"} {
		if _, err := validateTarget("dev", target, ""); err == nil {
			t.Fatalf("validateTarget(dev, %q) succeeded", target)
		}
	}
}

func TestValidateTargetAllowsOnlyApprovedRemoteDevOrigin(t *testing.T) {
	approved := "https://dev-api.example.invalid"
	if got, err := validateTarget("dev", approved, approved); err != nil || got != approved {
		t.Fatalf("validateTarget approved dev = %q, %v", got, err)
	}
	if _, err := validateTarget("dev", approved, "https://other-dev.example.invalid"); err == nil {
		t.Fatal("validateTarget accepted mismatched dev origin")
	}
}

func TestOutputIncludesFixtureAndFileProvenance(t *testing.T) {
	tempDir := t.TempDir()
	manifestPath := filepath.Join(tempDir, "fixtures.json")
	manifestBody := []byte(`{"manifest_version":1,"environment":"dev","seed_revision":"seed-sha","fixtures":{"normal_session":{"id":"normal"},"normal_empty_session":{"id":"empty"},"expired_session":{"id":"expired"}}}`)
	if err := os.WriteFile(manifestPath, manifestBody, 0o600); err != nil {
		t.Fatal(err)
	}
	outputPath := filepath.Join(tempDir, "go.json")
	if err := execute(context.Background(), "dev", "http://127.0.0.1", "", manifestPath, outputPath, strings.Repeat("a", 40), strings.Repeat("b", 64), true, ""); err != nil {
		t.Fatal(err)
	}
	body, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatal(err)
	}
	var got output
	if err := json.Unmarshal(body, &got); err != nil {
		t.Fatal(err)
	}
	digest := sha256.Sum256(manifestBody)
	if got.SeedRevision != "seed-sha" || got.FixtureManifestSHA256 != fmt.Sprintf("%x", digest) {
		t.Fatalf("fixture provenance = seed %q hash %q", got.SeedRevision, got.FixtureManifestSHA256)
	}
	for _, result := range got.Results {
		if result.EvidenceFile != filepath.Base(outputPath) {
			t.Fatalf("%s evidence_file = %q", result.ID, result.EvidenceFile)
		}
	}
}
