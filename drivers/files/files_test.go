package files

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
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
