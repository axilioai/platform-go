package telemetry

import (
	"context"
	"encoding/json"
	"io"
	"sort"
	"time"

	platformgo "github.com/axilioai/platform-go"
	"github.com/axilioai/platform-go/client"
	"github.com/axilioai/platform-go/option"
)

// Span types and attribute keys of the axilio.* vocabulary this package
// reads, mirroring the dashboard's trace viewer. Unknown values pass through
// untouched — the vocabulary is the contract's extension seam.
const (
	spanTypeSDKCall   = "sdk_call"
	spanTypeInference = "inference"

	attrDurationNs  = "axilio.duration_ns"
	attrInferenceID = "axilio.inference.id"
)

// _tracePageLimit is the archive page size: the op's maximum (1-1000), so a
// full trace costs the fewest round-trips.
const _tracePageLimit int64 = 1000

// _minGap is the smallest inter-call interval Summary counts as unobserved
// time; anything shorter is scheduling jitter, not a gap (same threshold as
// the dashboard timeline).
const _minGap = 5 * time.Millisecond

// framesLister is the one generated operation the helpers read. The runs
// subclient satisfies it; tests inject a fake.
type framesLister interface {
	SessionsListFrames(ctx context.Context, request *platformgo.SessionsListFramesRequest, opts ...option.RequestOption) (*platformgo.RunSessionFramesResponse, error)
}

// Session reads one session's telemetry: the durable trace archive through
// the API client, and — when constructed with the allocation's telemetry URL
// — the live leg.
type Session struct {
	id           string
	frames       framesLister
	telemetryURL string
	dial         dialFunc
}

// SessionOption configures a Session.
type SessionOption func(*Session)

// WithTelemetryURL enables the live leg (Tail, Logs with live=true) with the
// telemetry_url the allocation returned.
func WithTelemetryURL(telemetryURL string) SessionOption {
	return func(s *Session) { s.telemetryURL = telemetryURL }
}

// NewSession builds the telemetry view of one session.
func NewSession(c *client.Client, sessionID string, opts ...SessionOption) *Session {
	s := &Session{id: sessionID, frames: c.Runs, dial: dialWS}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

// Tail attaches to the session's live telemetry leg. Requires the session to
// have been constructed with WithTelemetryURL.
func (s *Session) Tail(ctx context.Context) (*Stream, error) {
	if s.telemetryURL == "" {
		return nil, &Error{Code: CodeUnauthorized, Message: "no telemetry URL: construct the session with WithTelemetryURL(allocation.TelemetryURL)"}
	}
	stream := &Stream{url: s.telemetryURL, dial: s.dial, redialDelay: redialDelay}
	conn, err := stream.dial(ctx, stream.attachURL())
	if err != nil {
		return nil, classifyDialErr(err)
	}
	stream.conn = conn
	return stream, nil
}

// TraceSpan is one completed span with its billed cost joined in. Billed
// cost is a read-time billing join, never a frame attribute: sdk_call spans
// bill by span id, inference spans carry the per-inference detail keyed by
// axilio.inference.id.
type TraceSpan struct {
	Span *platformgo.RunSpanFrame
	// BilledMicrodollars is the post-markup billed cost (what the invoice
	// charges), 0 for unbilled spans.
	BilledMicrodollars int64
}

// Duration is the span's measured duration: the producer's monotonic
// axilio.duration_ns when stamped, wall-clock end−start otherwise (which can
// cross a clock-skew edge).
func (t *TraceSpan) Duration() time.Duration {
	return spanDuration(t.Span)
}

// Trace is the ordered, cost-joined trace of one session — the dashboard
// trace viewer's data, reconstructed from the archive.
type Trace struct {
	SessionID string
	// Spans in start order, costs joined.
	Spans []*TraceSpan
	// Logs in time order.
	Logs []*platformgo.RunLogFrame
	// Unknown carries frames of a kind this SDK version doesn't know,
	// verbatim (tolerant reader: render generically, never drop).
	Unknown []*Frame
	// RetentionExpired is true when the trace is past the org's retention
	// window: frames are withheld and the lists above are empty.
	RetentionExpired bool
}

// Trace fetches the session's full durable trace, paging through the
// archive, and joins the response-level billed-cost maps onto the spans.
func (s *Session) Trace(ctx context.Context) (*Trace, error) {
	trace := &Trace{SessionID: s.id}
	sdkCallCosts := map[string]int64{}
	inferenceCosts := map[string]int64{}

	var frames []*Frame
	for offset := int64(0); ; {
		request := &platformgo.SessionsListFramesRequest{}
		request.SetSessionID(s.id)
		request.SetLimit(platformgo.Int64(_tracePageLimit))
		request.SetOffset(platformgo.Int64(offset))
		page, err := s.frames.SessionsListFrames(ctx, request)
		if err != nil {
			return nil, err
		}
		if page.RetentionExpired {
			trace.RetentionExpired = true
			return trace, nil
		}
		frames = append(frames, page.Frames...)
		for id, cost := range page.SdkCallCosts {
			sdkCallCosts[id] = cost
		}
		for id, cost := range page.InferenceCosts {
			inferenceCosts[id] = cost
		}
		offset += int64(len(page.Frames))
		if len(page.Frames) == 0 || offset >= page.Total {
			break
		}
	}

	for _, frame := range frames {
		switch {
		case frame.GetSpan() != nil:
			trace.Spans = append(trace.Spans, &TraceSpan{
				Span:               frame.GetSpan(),
				BilledMicrodollars: billedCost(frame.GetSpan(), sdkCallCosts, inferenceCosts),
			})
		case frame.GetLog() != nil:
			trace.Logs = append(trace.Logs, frame.GetLog())
		default:
			trace.Unknown = append(trace.Unknown, frame)
		}
	}
	sort.SliceStable(trace.Spans, func(i, j int) bool {
		return trace.Spans[i].Span.StartTimeUnixNano < trace.Spans[j].Span.StartTimeUnixNano
	})
	sort.SliceStable(trace.Logs, func(i, j int) bool {
		return trace.Logs[i].TimeUnixNano < trace.Logs[j].TimeUnixNano
	})
	return trace, nil
}

// Summary is the trace viewer's header rollup.
type Summary struct {
	SessionID string
	// Total is the session's wall-clock span: the session-root span's
	// lifetime when present, earliest-to-latest telemetry otherwise.
	Total time.Duration
	// SDK is the summed duration of the sdk_call spans.
	SDK time.Duration
	// Unobserved is the summed time inside Total not covered by any
	// sdk_call — a sleep, an OCR loop, an untraced network call. Gaps under
	// 5ms are dropped as scheduling jitter.
	Unobserved time.Duration
	// BillableCalls counts the sdk_call spans with a nonzero billed cost.
	BillableCalls int
	// CostMicrodollars is the summed billed cost of those calls.
	CostMicrodollars int64
}

// Summary fetches the trace and reduces it to the header rollup.
func (s *Session) Summary(ctx context.Context) (*Summary, error) {
	trace, err := s.Trace(ctx)
	if err != nil {
		return nil, err
	}
	return summarize(s.id, trace), nil
}

// summarize mirrors the dashboard's timeline math: anchor on the session
// root (fallback: earliest telemetry), end at the latest span end or log,
// sweep the sdk_call intervals for gaps.
func summarize(sessionID string, trace *Trace) *Summary {
	summary := &Summary{SessionID: sessionID}

	var anchor int64
	var end int64
	anchored := false
	for _, ts := range trace.Spans {
		span := ts.Span
		if span.SpanType == spanTypeSession || span.SpanType == spanTypeSessionLegacy {
			anchor = span.StartTimeUnixNano
			anchored = true
		}
	}
	if !anchored {
		// A phone session need not be a run: without a root span, anchor on
		// the earliest telemetry of any shape.
		for _, ts := range trace.Spans {
			if start := ts.Span.StartTimeUnixNano; start > 0 && (!anchored || start < anchor) {
				anchor = start
				anchored = true
			}
		}
		for _, log := range trace.Logs {
			if at := log.TimeUnixNano; at > 0 && (!anchored || at < anchor) {
				anchor = at
				anchored = true
			}
		}
	}
	if !anchored {
		return summary
	}
	for _, ts := range trace.Spans {
		spanEnd := spanEndNano(ts.Span)
		if spanEnd == 0 {
			spanEnd = ts.Span.StartTimeUnixNano
		}
		if spanEnd > end {
			end = spanEnd
		}
	}
	for _, log := range trace.Logs {
		if log.TimeUnixNano > end {
			end = log.TimeUnixNano
		}
	}
	summary.Total = time.Duration(max(end-anchor, 0))

	// One interval per sdk_call, clamped to the timeline; Spans is already
	// start-ordered, which is what the gap sweep needs.
	var cursor time.Duration
	for _, ts := range trace.Spans {
		if ts.Span.SpanType != spanTypeSDKCall {
			continue
		}
		duration := spanDuration(ts.Span)
		summary.SDK += duration
		if ts.BilledMicrodollars > 0 {
			summary.BillableCalls++
			summary.CostMicrodollars += ts.BilledMicrodollars
		}
		start := time.Duration(max(ts.Span.StartTimeUnixNano-anchor, 0))
		if gap := start - cursor; gap >= _minGap {
			summary.Unobserved += gap
		}
		if callEnd := start + duration; callEnd > cursor {
			cursor = callEnd
		}
	}
	if tail := summary.Total - cursor; tail >= _minGap {
		summary.Unobserved += tail
	}
	return summary
}

// Logs returns the session's log frames. With live=false the archive's logs
// stream in time order and the stream ends immediately; with live=true the
// live leg is tailed and log frames stream as they happen until the session
// ends (io.EOF, like Stream).
func (s *Session) Logs(ctx context.Context, live bool) (*LogStream, error) {
	if !live {
		trace, err := s.Trace(ctx)
		if err != nil {
			return nil, err
		}
		return &LogStream{queue: trace.Logs}, nil
	}
	stream, err := s.Tail(ctx)
	if err != nil {
		return nil, err
	}
	return &LogStream{stream: stream}, nil
}

// LogStream yields log frames; Next returns io.EOF at the end of the logs.
type LogStream struct {
	queue  []*platformgo.RunLogFrame
	stream *Stream // nil in archive mode
}

// Next returns the next log frame, skipping non-log frames in live mode.
func (l *LogStream) Next(ctx context.Context) (*platformgo.RunLogFrame, error) {
	if l.stream == nil {
		if len(l.queue) == 0 {
			return nil, io.EOF
		}
		log := l.queue[0]
		l.queue = l.queue[1:]
		return log, nil
	}
	for {
		frame, err := l.stream.Next(ctx)
		if err != nil {
			return nil, err
		}
		if log := frame.GetLog(); log != nil {
			return log, nil
		}
	}
}

// Close tears down the live leg; a no-op in archive mode.
func (l *LogStream) Close() error {
	if l.stream == nil {
		return nil
	}
	return l.stream.Close()
}

// billedCost joins a span to the response-level billing maps: sdk_call spans
// bill by span id; inference spans surface the per-inference detail behind
// the parent call's cost.
func billedCost(span *platformgo.RunSpanFrame, sdkCallCosts, inferenceCosts map[string]int64) int64 {
	switch span.SpanType {
	case spanTypeSDKCall:
		return sdkCallCosts[span.SpanID]
	case spanTypeInference:
		id, ok := span.Attributes[attrInferenceID].(string)
		if !ok {
			return 0
		}
		return inferenceCosts[id]
	default:
		return 0
	}
}

// spanDuration prefers the producer's monotonic axilio.duration_ns over
// wall-clock end−start, which can cross a clock-skew edge.
func spanDuration(span *platformgo.RunSpanFrame) time.Duration {
	if ns, ok := attrNumber(span.Attributes, attrDurationNs); ok && ns > 0 {
		return time.Duration(ns)
	}
	if end := spanEndNano(span); end > span.StartTimeUnixNano {
		return time.Duration(end - span.StartTimeUnixNano)
	}
	return 0
}

// spanEndNano reads a span's end time, mapping the omitted (in-flight) case
// to 0 — the sentinel this package keys on. Since spec 0.82.0 the field is
// optional on the wire: live start-phase frames carry no end time.
func spanEndNano(span *platformgo.RunSpanFrame) int64 {
	if span == nil || span.EndTimeUnixNano == nil {
		return 0
	}
	return *span.EndTimeUnixNano
}

// attrNumber reads a numeric attribute, tolerating the JSON number types an
// attributes map can carry.
func attrNumber(attributes map[string]any, key string) (int64, bool) {
	switch v := attributes[key].(type) {
	case float64:
		return int64(v), true
	case int64:
		return v, true
	case int:
		return int64(v), true
	case json.Number:
		n, err := v.Int64()
		if err != nil {
			return 0, false
		}
		return n, true
	default:
		return 0, false
	}
}
