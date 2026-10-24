package telemetry

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"
)

const (
	sdkSpanJSON    = `{"kind":"span","span_type":"sdk_call","phase":"end","trace_id":"t1","span_id":"s1","name":"Screen.observe","start_time_unix_nano":100,"end_time_unix_nano":200,"status":{"code":"ok","message":""}}`
	logJSON        = `{"kind":"log","log_type":"output_log","severity":"INFO","body":"hello","time_unix_nano":150,"trace_id":"t1"}`
	endFrameJSON   = `{"kind":"span","span_type":"session","phase":"end","trace_id":"t1","span_id":"root","name":"session","start_time_unix_nano":1,"end_time_unix_nano":900,"status":{"code":"ok","message":""}}`
	unknownKindMsg = `{"kind":"metric","name":"cpu","value":42}`
)

// sseServer answers the Nth request with the Nth response: a status and,
// for 200, a raw event-stream body written and flushed before the response
// ends. It records each request's Last-Event-ID and API key header.
type sseServer struct {
	mu        sync.Mutex
	responses []sseResponse
	cursors   []string
	keys      []string
}

type sseResponse struct {
	status int
	body   string
}

func (s *sseServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mu.Lock()
	s.cursors = append(s.cursors, r.Header.Get("Last-Event-ID"))
	s.keys = append(s.keys, r.Header.Get("X-Axilio-Api-Key"))
	resp := sseResponse{status: http.StatusServiceUnavailable}
	if n := len(s.cursors); n <= len(s.responses) {
		resp = s.responses[n-1]
	}
	s.mu.Unlock()
	if resp.status != http.StatusOK {
		w.WriteHeader(resp.status)
		return
	}
	w.Header().Set("Content-Type", "text/event-stream")
	_, _ = io.WriteString(w, resp.body)
	w.(http.Flusher).Flush()
}

func tailTest(t *testing.T, srv *sseServer) (*Stream, error) {
	t.Helper()
	ts := httptest.NewServer(srv)
	t.Cleanup(ts.Close)
	s := &Stream{
		url:         ts.URL,
		open:        httpOpener(tailConfig{client: ts.Client(), header: http.Header{"X-Axilio-Api-Key": {"axl_test"}}}),
		redialDelay: func(int) time.Duration { return time.Nanosecond },
	}
	return start(context.Background(), s)
}

func drain(t *testing.T, s *Stream) ([]*Frame, error) {
	t.Helper()
	var frames []*Frame
	for {
		f, err := s.Next(context.Background())
		if err != nil {
			return frames, err
		}
		frames = append(frames, f)
	}
}

// TestStreamReadsEventsAndResumes: frames arrive one per event (heartbeats
// and cursor-only events skipped, unknown kinds kept), a drop reopens with
// the last cursor as Last-Event-ID, and the session end frame ends the
// stream cleanly.
func TestStreamReadsEventsAndResumes(t *testing.T) {
	srv := &sseServer{responses: []sseResponse{
		{http.StatusOK, ": heartbeat\n\n" +
			"event: span\ndata: " + sdkSpanJSON + "\n\n" +
			"id: 1-1\nevent: metric\ndata: " + unknownKindMsg + "\n\n" +
			"id: 1-2\n\n"},
		{http.StatusOK, "id: 1-3\nevent: log\ndata: " + logJSON + "\n\n" +
			"id: 1-4\nevent: span\ndata: " + endFrameJSON + "\n\n"},
	}}
	s, err := tailTest(t, srv)
	if err != nil {
		t.Fatalf("Tail: %v", err)
	}
	frames, err := drain(t, s)
	if !errors.Is(err, io.EOF) || !s.Ended() {
		t.Fatalf("end = %v (ended %v), want io.EOF after the session end frame", err, s.Ended())
	}
	var kinds []string
	for _, f := range frames {
		kinds = append(kinds, f.Kind)
	}
	if want := []string{"span", "metric", "log", "span"}; !equal(kinds, want) {
		t.Fatalf("kinds = %v, want %v", kinds, want)
	}
	if want := []string{"", "1-2"}; !equal(srv.cursors, want) {
		t.Fatalf("Last-Event-ID per request = %v, want %v", srv.cursors, want)
	}
	if want := []string{"axl_test", "axl_test"}; !equal(srv.keys, want) {
		t.Fatalf("API key per request = %v, want %v", srv.keys, want)
	}
}

// TestStreamSkipsFramesAReplayRepeats: a reopen whose cursor the live window
// has trimmed past gets the archive replayed (no ids), repeating frames the
// caller already has; each is returned once.
func TestStreamSkipsFramesAReplayRepeats(t *testing.T) {
	otherSpanJSON := strings.Replace(sdkSpanJSON, `"span_id":"s1"`, `"span_id":"s2"`, 1)
	srv := &sseServer{responses: []sseResponse{
		{http.StatusOK, "id: 1-1\nevent: span\ndata: " + sdkSpanJSON + "\n\n" +
			"id: 1-2\nevent: log\ndata: " + logJSON + "\n\n"},
		{http.StatusOK, "event: span\ndata: " + sdkSpanJSON + "\n\n" +
			"event: log\ndata: " + logJSON + "\n\n" +
			"event: span\ndata: " + otherSpanJSON + "\n\n" +
			"id: 2-1\nevent: span\ndata: " + endFrameJSON + "\n\n"},
	}}
	s, err := tailTest(t, srv)
	if err != nil {
		t.Fatalf("Tail: %v", err)
	}
	frames, err := drain(t, s)
	if !errors.Is(err, io.EOF) || !s.Ended() {
		t.Fatalf("end = %v (ended %v), want io.EOF after the session end frame", err, s.Ended())
	}
	var got []string
	for _, f := range frames {
		if span := f.GetSpan(); span != nil {
			got = append(got, span.SpanID)
		} else {
			got = append(got, f.Kind)
		}
	}
	if want := []string{"s1", "log", "s2", "root"}; !equal(got, want) {
		t.Fatalf("frames = %v, want %v", got, want)
	}
}

// TestStreamStatuses: how each server answer ends a stream, on the first
// open and on a reopen.
func TestStreamStatuses(t *testing.T) {
	cases := []struct {
		name      string
		responses []sseResponse
		wantErr   func(error) bool
		wantEnded bool
	}{
		{"204 on open: the session is over", []sseResponse{{status: http.StatusNoContent}}, isEOF, true},
		{"204 on reopen: nothing left after the cursor", []sseResponse{{http.StatusOK, "id: 1-1\nevent: span\ndata: " + sdkSpanJSON + "\n\n"}, {status: http.StatusNoContent}}, isEOF, true},
		{"401: credential rejected", []sseResponse{{status: http.StatusUnauthorized}}, IsUnauthorized, false},
		{"404: no such session", []sseResponse{{status: http.StatusNotFound}}, IsNotFound, false},
		{"409: telemetry off", []sseResponse{{status: http.StatusConflict}}, IsTelemetryOff, false},
		{"drops and 5xx until the budget runs out", []sseResponse{{http.StatusOK, ""}}, IsConnection, false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			s, err := tailTest(t, &sseServer{responses: tc.responses})
			if err == nil {
				_, err = drain(t, s)
			}
			if !tc.wantErr(err) {
				t.Fatalf("err = %v", err)
			}
			if s != nil && s.Ended() != tc.wantEnded {
				t.Fatalf("Ended = %v, want %v", s.Ended(), tc.wantEnded)
			}
		})
	}
}

// TestStreamCloseUnblocksNext: Close from another goroutine ends a Next that
// is blocked on a silent stream.
func TestStreamCloseUnblocksNext(t *testing.T) {
	release := make(chan struct{})
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		w.(http.Flusher).Flush()
		<-release
	}))
	t.Cleanup(ts.Close)
	t.Cleanup(func() { close(release) })
	s, err := Tail(context.Background(), ts.URL)
	if err != nil {
		t.Fatalf("Tail: %v", err)
	}
	done := make(chan error, 1)
	go func() {
		_, err := s.Next(context.Background())
		done <- err
	}()
	time.Sleep(50 * time.Millisecond)
	_ = s.Close()
	select {
	case err := <-done:
		if !hasCode(err, CodeClosed) {
			t.Fatalf("Next after Close = %v, want CodeClosed", err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("Next stayed blocked after Close")
	}
}

func isEOF(err error) bool { return errors.Is(err, io.EOF) }

func equal(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
