package telemetry

import (
	"context"
	"encoding/json"
	"io"
	"testing"
	"time"

	platformgo "github.com/axilioai/platform-go"
	"github.com/axilioai/platform-go/option"
)

// fakeLister serves scripted archive pages and records the requested offsets.
type fakeLister struct {
	pages   []string // raw RunSessionFramesResponse JSON, one per call
	offsets []int64
}

func (f *fakeLister) SessionsListFrames(_ context.Context, request *platformgo.SessionsListFramesRequest, _ ...option.RequestOption) (*platformgo.RunSessionFramesResponse, error) {
	if request.Offset != nil {
		f.offsets = append(f.offsets, *request.Offset)
	}
	if len(f.pages) == 0 {
		return nil, io.ErrUnexpectedEOF
	}
	raw := f.pages[0]
	f.pages = f.pages[1:]
	response := new(platformgo.RunSessionFramesResponse)
	if err := json.Unmarshal([]byte(raw), response); err != nil {
		return nil, err
	}
	return response, nil
}

// The scripted trace: a 10s session with two sdk calls (one billed, one with
// an inference child), a log, and an unknown-kind frame, split over two pages
// to exercise pagination and cost-map merging.
const (
	tracePage1 = `{
		"frames": [
			{"kind":"span","span_type":"session","phase":"end","trace_id":"t1","span_id":"root","name":"session","start_time_unix_nano":0,"end_time_unix_nano":10000000000,"status":{"code":"ok","message":""}},
			{"kind":"span","span_type":"sdk_call","phase":"end","trace_id":"t1","span_id":"call-1","name":"Screen.observe","start_time_unix_nano":1000000000,"end_time_unix_nano":2000000000,"status":{"code":"ok","message":""}},
			{"kind":"span","span_type":"inference","phase":"end","trace_id":"t1","span_id":"inf-span","parent_span_id":"call-1","name":"inference","start_time_unix_nano":1200000000,"end_time_unix_nano":1800000000,"attributes":{"axilio.inference.id":"inf-1"},"status":{"code":"ok","message":""}}
		],
		"sdk_call_costs": {"call-1": 700},
		"inference_costs": {},
		"limit": 1000, "offset": 0, "total": 5, "retention_expired": false
	}`
	tracePage2 = `{
		"frames": [
			{"kind":"span","span_type":"sdk_call","phase":"end","trace_id":"t1","span_id":"call-2","name":"Screen.tap","start_time_unix_nano":4000000000,"end_time_unix_nano":5900000000,"attributes":{"axilio.duration_ns":2000000000},"status":{"code":"ok","message":""}},
			{"kind":"log","log_type":"output_log","severity":"INFO","body":"done","time_unix_nano":9000000000,"trace_id":"t1"},
			{"kind":"metric","name":"cpu","value":1}
		],
		"sdk_call_costs": {},
		"inference_costs": {"inf-1": 300},
		"limit": 1000, "offset": 3, "total": 6, "retention_expired": false
	}`
)

func newTestSession(lister framesLister) *Session {
	return &Session{id: "sess-1", frames: lister}
}

func TestTracePaginatesAndJoinsCosts(t *testing.T) {
	lister := &fakeLister{pages: []string{tracePage1, tracePage2}}
	trace, err := newTestSession(lister).Trace(context.Background())
	if err != nil {
		t.Fatalf("Trace: %v", err)
	}
	if len(lister.offsets) != 2 || lister.offsets[0] != 0 || lister.offsets[1] != 3 {
		t.Fatalf("offsets = %v, want [0 3]", lister.offsets)
	}
	if len(trace.Spans) != 4 {
		t.Fatalf("spans = %d, want 4", len(trace.Spans))
	}
	// Start-ordered regardless of page boundaries.
	for i, want := range []string{"root", "call-1", "inf-span", "call-2"} {
		if got := trace.Spans[i].Span.SpanID; got != want {
			t.Fatalf("span[%d] = %s, want %s", i, got, want)
		}
	}
	costs := map[string]int64{}
	for _, ts := range trace.Spans {
		costs[ts.Span.SpanID] = ts.BilledMicrodollars
	}
	// sdk_call joins by span id, inference by axilio.inference.id; the rest
	// (session root, the unbilled call) stay 0.
	if costs["call-1"] != 700 || costs["inf-span"] != 300 || costs["root"] != 0 || costs["call-2"] != 0 {
		t.Fatalf("cost join = %v, want call-1:700 inf-span:300 others:0", costs)
	}
	if len(trace.Logs) != 1 || trace.Logs[0].Body != "done" {
		t.Fatalf("logs = %+v, want the one output_log", trace.Logs)
	}
	if len(trace.Unknown) != 1 || trace.Unknown[0].Kind != "metric" {
		t.Fatalf("unknown = %+v, want the metric frame kept verbatim", trace.Unknown)
	}
}

func TestSummaryMirrorsTimelineMath(t *testing.T) {
	lister := &fakeLister{pages: []string{tracePage1, tracePage2}}
	summary, err := newTestSession(lister).Summary(context.Background())
	if err != nil {
		t.Fatalf("Summary: %v", err)
	}
	if summary.Total != 10*time.Second {
		t.Fatalf("Total = %v, want 10s (session root lifetime)", summary.Total)
	}
	// call-1 is 1s wall clock; call-2 measures 2s via axilio.duration_ns
	// (preferred over its 1.9s wall clock).
	if summary.SDK != 3*time.Second {
		t.Fatalf("SDK = %v, want 3s", summary.SDK)
	}
	// Gaps on the sdk_call sweep: [0,1s), [2s,4s), [6s,10s] = 7s.
	if summary.Unobserved != 7*time.Second {
		t.Fatalf("Unobserved = %v, want 7s", summary.Unobserved)
	}
	if summary.BillableCalls != 1 || summary.CostMicrodollars != 700 {
		t.Fatalf("billing rollup = %d calls / %d µ$, want 1 / 700 (inference detail stays on its own span)", summary.BillableCalls, summary.CostMicrodollars)
	}
}

func TestTraceRetentionExpired(t *testing.T) {
	lister := &fakeLister{pages: []string{`{"frames":[],"sdk_call_costs":{},"inference_costs":{},"limit":1000,"offset":0,"total":0,"retention_expired":true}`}}
	trace, err := newTestSession(lister).Trace(context.Background())
	if err != nil {
		t.Fatalf("Trace: %v", err)
	}
	if !trace.RetentionExpired || len(trace.Spans) != 0 {
		t.Fatalf("trace = %+v, want empty with RetentionExpired", trace)
	}
}

func TestLogsArchiveMode(t *testing.T) {
	lister := &fakeLister{pages: []string{tracePage1, tracePage2}}
	logs, err := newTestSession(lister).Logs(context.Background(), false)
	if err != nil {
		t.Fatalf("Logs: %v", err)
	}
	first, err := logs.Next(context.Background())
	if err != nil || first.Body != "done" {
		t.Fatalf("first log = %+v (%v), want the archived output_log", first, err)
	}
	if _, err := logs.Next(context.Background()); err != io.EOF {
		t.Fatalf("err = %v, want io.EOF", err)
	}
}

func TestLogsLiveModeFiltersToLogs(t *testing.T) {
	d := &scriptDialer{script: []func() (rawConn, error){func() (rawConn, error) {
		return &scriptConn{messages: msgs(sdkSpanJSON, logJSON, endFrameJSON)}, nil
	}}}
	session := &Session{id: "sess-1", telemetryURL: "wss://example.test/ws", dial: d.dial}
	logs, err := session.Logs(context.Background(), true)
	if err != nil {
		t.Fatalf("Logs(live): %v", err)
	}
	first, err := logs.Next(context.Background())
	if err != nil || first.Body != "hello" {
		t.Fatalf("first live log = %+v (%v), want the log frame, spans skipped", first, err)
	}
	if _, err := logs.Next(context.Background()); err != io.EOF {
		t.Fatalf("err = %v, want io.EOF after the session end frame", err)
	}
}

func TestTailRequiresTelemetryURL(t *testing.T) {
	session := &Session{id: "sess-1"}
	if _, err := session.Tail(context.Background()); err == nil {
		t.Fatal("Tail without a telemetry URL succeeded, want an error")
	}
}
