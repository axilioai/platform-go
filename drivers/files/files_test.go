package files

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	platformgo "github.com/axilioai/platform-go"
	client "github.com/axilioai/platform-go/client"
	"github.com/axilioai/platform-go/option"
)

// These tests exist for the same reason platform-python's do: this package is
// hand-written against a generated client that regenerates underneath it, and
// while Go's compiler catches a renamed method, it cannot catch an option that
// silently stops being honoured. Push ignoring WithWait compiled fine and was
// wrong for a release.

// newTestClient points the generated client at a stub server.
func newTestClient(t *testing.T, h http.Handler) (*client.Client, func()) {
	t.Helper()
	srv := httptest.NewServer(h)
	return client.NewClient(
		option.WithBaseURL(srv.URL+"/api/v1"), // callers prepend the host; the generated default is the bare "/api/v1"
		option.WithAPIKey("axl_test"),
	), srv.Close
}

func writeJSON(t *testing.T, w http.ResponseWriter, v any) {
	t.Helper()
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(v); err != nil {
		t.Fatalf("encode response: %v", err)
	}
}

// deliveryJSON is the shape of a FileDeliverySummary on the wire.
func deliveryJSON(status string) map[string]any {
	return map[string]any{
		"id":         "del_1",
		"file_id":    "file_1",
		"phone_id":   "phn_1",
		"filename":   "demo.png",
		"mime_type":  "image/png",
		"size_bytes": 3,
		"status":     status,
		"created_at": "2026-07-26T00:00:00Z",
	}
}

func TestPushWithoutWaitReturnsAtDispatch(t *testing.T) {
	var gets int
	c, done := newTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			gets++
		}
		writeJSON(t, w, map[string]any{"delivery": deliveryJSON("dispatched")})
	}))
	defer done()

	d, err := Push(context.Background(), c, "phn_1", "file_1")
	if err != nil {
		t.Fatalf("Push: %v", err)
	}
	if d.Status != platformgo.FileDeliverySummaryStatusDispatched {
		t.Fatalf("status = %q, want dispatched", d.Status)
	}
	if gets != 0 {
		t.Fatalf("bare Push polled %d times; it must return at dispatch", gets)
	}
}

// TestPushWithWaitPollsToTerminal is the regression this file was added for:
// Push accepted WithWait and discarded it, so the one flow Push exists for —
// pushing a stored file to several phones — could never be waited on.
func TestPushWithWaitPollsToTerminal(t *testing.T) {
	var gets int
	c, done := newTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			writeJSON(t, w, map[string]any{"delivery": deliveryJSON("dispatched")})
			return
		}
		// The wait loop must read THIS delivery by id, not a page of the
		// newest ones — a busy phone would push it off the page.
		if got, want := r.URL.Path, "/api/v1/phones/phn_1/deliveries/del_1"; got != want {
			t.Errorf("poll path = %q, want %q", got, want)
		}
		gets++
		writeJSON(t, w, deliveryJSON("delivered"))
	}))
	defer done()

	d, err := Push(context.Background(), c, "phn_1", "file_1",
		WithWait(5*time.Second), WithPollInterval(10*time.Millisecond))
	if err != nil {
		t.Fatalf("Push: %v", err)
	}
	if d.Status != platformgo.FileDeliverySummaryStatusDelivered {
		t.Fatalf("status = %q, want delivered — WithWait was ignored", d.Status)
	}
	if gets == 0 {
		t.Fatal("wait loop never polled the delivery")
	}
}

// TestSendWaitMatchesPushWait pins the property that made Send's private copy
// of the wait loop removable: Send is Upload followed by Push, so wait
// behaves identically through either entry point.
func TestSendWaitMatchesPushWait(t *testing.T) {
	path := t.TempDir() + "/demo.png"
	if err := os.WriteFile(path, []byte("abc"), 0o600); err != nil {
		t.Fatalf("write fixture: %v", err)
	}

	var puts, polls int
	c, done := newTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodPut:
			puts++
			w.WriteHeader(http.StatusOK)
		case r.URL.Path == "/api/v1/uploads":
			writeJSON(t, w, map[string]any{
				"file":                      fileJSON("uploading"),
				"upload_url":                "http://" + r.Host + "/put-here",
				"upload_expires_in_seconds": 900,
			})
		case r.URL.Path == "/api/v1/uploads/file_1/complete":
			writeJSON(t, w, map[string]any{"file": fileJSON("ready")})
		case r.Method == http.MethodPost:
			writeJSON(t, w, map[string]any{"delivery": deliveryJSON("dispatched")})
		default:
			polls++
			writeJSON(t, w, deliveryJSON("delivered"))
		}
	}))
	defer done()

	d, err := Send(context.Background(), c, "phn_1", path,
		WithWait(5*time.Second), WithPollInterval(10*time.Millisecond))
	if err != nil {
		t.Fatalf("Send: %v", err)
	}
	if d.Status != platformgo.FileDeliverySummaryStatusDelivered {
		t.Fatalf("status = %q, want delivered", d.Status)
	}
	if puts != 1 {
		t.Fatalf("PUT count = %d, want exactly 1", puts)
	}
	if polls == 0 {
		t.Fatal("Send did not wait")
	}
}

func fileJSON(status string) map[string]any {
	return map[string]any{
		"id":         "file_1",
		"filename":   "demo.png",
		"mime_type":  "image/png",
		"size_bytes": 3,
		"status":     status,
		"created_at": "2026-07-26T00:00:00Z",
	}
}

func TestDetectMIME(t *testing.T) {
	// Every extension the API's upload whitelist accepts must resolve to the
	// exact type the server expects, on any host. HEIC is the one that
	// motivated the table: mime.TypeByExtension does not reliably know it.
	for ext, want := range map[string]string{
		".heic": "image/heic",
		".jpg":  "image/jpeg",
		".JPEG": "image/jpeg", // case-insensitive
		".png":  "image/png",
		".webp": "image/webp",
		".gif":  "image/gif",
		".mp4":  "video/mp4",
		".webm": "video/webm",
		".mov":  "video/quicktime",
		".3gp":  "video/3gpp",
		".mkv":  "video/x-matroska",
	} {
		if got := detectMIME("demo" + ext); got != want {
			t.Errorf("detectMIME(demo%s) = %q, want %q", ext, got, want)
		}
	}

	if got := detectMIME("demo.unknown-ext"); got != defaultMIME {
		t.Errorf("unknown extension = %q, want %q", got, defaultMIME)
	}
	// Whatever the host database says, the API compares the bare type, so
	// parameters must never reach the registration.
	if got := detectMIME("notes.txt"); strings.Contains(got, ";") {
		t.Errorf("detectMIME returned parameters: %q", got)
	}
}

// TestUploadStreamsAndPinsHeaders is the regression guard for the 1 GB limit:
// Upload used to os.ReadFile the whole payload into memory. It also pins the
// two headers the presigned PUT signature covers — a mismatch there is
// rejected by storage with an opaque error.
func TestUploadStreamsAndPinsHeaders(t *testing.T) {
	path := t.TempDir() + "/demo.heic"
	payload := bytes.Repeat([]byte("axilio"), 4096) // 24 KiB
	if err := os.WriteFile(path, payload, 0o600); err != nil {
		t.Fatalf("write fixture: %v", err)
	}

	var (
		gotBody          []byte
		gotType          string
		gotLength        int64
		gotChunked       bool
		registeredType   string
		registeredLength float64
	)
	c, done := newTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodPut:
			gotType = r.Header.Get("Content-Type")
			gotLength = r.ContentLength
			gotChunked = len(r.TransferEncoding) > 0
			gotBody, _ = io.ReadAll(r.Body)
			w.WriteHeader(http.StatusOK)
		case r.URL.Path == "/api/v1/uploads":
			var body map[string]any
			_ = json.NewDecoder(r.Body).Decode(&body)
			registeredType, _ = body["mime_type"].(string)
			registeredLength, _ = body["size_bytes"].(float64)
			writeJSON(t, w, map[string]any{
				"file":                      fileJSON("uploading"),
				"upload_url":                "http://" + r.Host + "/put-here",
				"upload_expires_in_seconds": 900,
			})
		default:
			writeJSON(t, w, map[string]any{"file": fileJSON("ready")})
		}
	}))
	defer done()

	if _, err := Upload(context.Background(), c, path); err != nil {
		t.Fatalf("Upload: %v", err)
	}

	if !bytes.Equal(gotBody, payload) {
		t.Errorf("storage received %d bytes, want %d", len(gotBody), len(payload))
	}
	if gotLength != int64(len(payload)) {
		t.Errorf("Content-Length = %d, want %d", gotLength, len(payload))
	}
	if gotChunked {
		t.Error("request used chunked encoding; the presigned signature pins an exact length")
	}
	// The whole point of the extension table: this must not be octet-stream.
	if gotType != "image/heic" || registeredType != "image/heic" {
		t.Errorf("content type sent=%q registered=%q, want image/heic both", gotType, registeredType)
	}
	if int64(registeredLength) != int64(len(payload)) {
		t.Errorf("registered size_bytes = %v, want %d", registeredLength, len(payload))
	}
}

func TestUploadUsesInjectedHTTPClient(t *testing.T) {
	path := t.TempDir() + "/demo.png"
	if err := os.WriteFile(path, []byte("abc"), 0o600); err != nil {
		t.Fatalf("write fixture: %v", err)
	}
	c, done := newTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPut {
			w.WriteHeader(http.StatusOK)
			return
		}
		if r.URL.Path == "/api/v1/uploads" {
			writeJSON(t, w, map[string]any{
				"file":                      fileJSON("uploading"),
				"upload_url":                "http://" + r.Host + "/put-here",
				"upload_expires_in_seconds": 900,
			})
			return
		}
		writeJSON(t, w, map[string]any{"file": fileJSON("ready")})
	}))
	defer done()

	used := false
	custom := &http.Client{Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
		used = true
		return http.DefaultTransport.RoundTrip(r)
	})}

	if _, err := Upload(context.Background(), c, path, WithHTTPClient(custom)); err != nil {
		t.Fatalf("Upload: %v", err)
	}
	if !used {
		t.Error("WithHTTPClient was ignored; the storage PUT did not go through the injected client")
	}

	// A nil client must not disable uploads.
	if _, err := Upload(context.Background(), c, path, WithHTTPClient(nil)); err != nil {
		t.Fatalf("Upload with nil client: %v", err)
	}
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }

func TestUploadRejectsMissingFileBeforeRegistering(t *testing.T) {
	// Opening first means a bad path costs no library row and no quota.
	var calls int
	c, done := newTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		calls++
		w.WriteHeader(http.StatusOK)
	}))
	defer done()

	_, err := Upload(context.Background(), c, t.TempDir()+"/does-not-exist.png")
	if err == nil {
		t.Fatal("expected an error for a missing file")
	}
	if !errors.Is(err, os.ErrNotExist) {
		t.Errorf("error %v does not wrap os.ErrNotExist; %%w was dropped somewhere", err)
	}
	if calls != 0 {
		t.Errorf("made %d API calls before reading the file; expected none", calls)
	}
}
