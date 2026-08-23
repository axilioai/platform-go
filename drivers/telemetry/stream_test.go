package telemetry

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"strings"
	"testing"
	"time"

	"github.com/coder/websocket"
)

// scriptConn is a scripted rawConn: recv drains the queued messages, then
// fails with err (the scripted force-close).
type scriptConn struct {
	messages [][]byte
	err      error
}

func (c *scriptConn) recv(_ context.Context) ([]byte, error) {
	if len(c.messages) > 0 {
		m := c.messages[0]
		c.messages = c.messages[1:]
		return m, nil
	}
	if c.err != nil {
		return nil, c.err
	}
	return nil, errors.New("scripted conn exhausted")
}

func (c *scriptConn) closeConn() error { return nil }

// scriptDialer hands out one scripted connection (or dial error) per dial,
// recording each attach URL so tests can assert the resume params.
type scriptDialer struct {
	script []func() (rawConn, error)
	urls   []string
}

func (d *scriptDialer) dial(_ context.Context, url string) (rawConn, error) {
	d.urls = append(d.urls, url)
	if len(d.script) == 0 {
		return nil, errors.New("dial script exhausted")
	}
	next := d.script[0]
	d.script = d.script[1:]
	return next()
}

// newTestStream builds a Stream over the dialer's first connection with an
// instant backoff.
func newTestStream(t *testing.T, d *scriptDialer) *Stream {
	t.Helper()
	s := &Stream{url: "wss://example.test/realtime/ws/telemetry?token=x", dial: d.dial, redialDelay: func(int) time.Duration { return time.Nanosecond }}
	conn, err := s.dial(context.Background(), s.attachURL())
	if err != nil {
		t.Fatalf("initial dial: %v", err)
	}
	s.conn = conn
	return s
}

const (
	sdkSpanJSON    = `{"kind":"span","span_type":"sdk_call","phase":"end","trace_id":"t1","span_id":"s1","name":"Screen.observe","start_time_unix_nano":100,"end_time_unix_nano":200,"status":{"code":"ok","message":""}}`
	logJSON        = `{"kind":"log","log_type":"output_log","severity":"INFO","body":"hello","time_unix_nano":150,"trace_id":"t1"}`
	endFrameJSON   = `{"kind":"span","span_type":"session","phase":"end","trace_id":"t1","span_id":"root","name":"session","start_time_unix_nano":1,"end_time_unix_nano":900,"status":{"code":"ok","message":""}}`
	legacyEndJSON  = `{"kind":"span","span_type":"phone_session","phase":"end","trace_id":"t1","span_id":"root","name":"session","start_time_unix_nano":1,"end_time_unix_nano":900,"status":{"code":"ok","message":""}}`
	unknownKindMsg = `{"kind":"metric","name":"cpu","value":42}`
)

func msgs(raw ...string) [][]byte {
	out := make([][]byte, len(raw))
	for i, r := range raw {
		out[i] = []byte(r)
	}
	return out
}

func mustNext(t *testing.T, s *Stream) *Frame {
	t.Helper()
	frame, err := s.Next(context.Background())
	if err != nil {
		t.Fatalf("Next: %v", err)
	}
	return frame
}

func TestStreamYieldsSingleAndArrayMessages(t *testing.T) {
	d := &scriptDialer{script: []func() (rawConn, error){func() (rawConn, error) {
		return &scriptConn{messages: msgs(
			sdkSpanJSON,
			"["+logJSON+","+sdkSpanJSON+"]",
			endFrameJSON,
		)}, nil
	}}}
	s := newTestStream(t, d)

	if f := mustNext(t, s); f.GetSpan() == nil || f.GetSpan().SpanID != "s1" {
		t.Fatalf("first frame = %+v, want sdk span s1", f)
	}
	if f := mustNext(t, s); f.GetLog() == nil || f.GetLog().Body != "hello" {
		t.Fatalf("second frame = %+v, want log from array batch", f)
	}
	if f := mustNext(t, s); f.GetSpan() == nil {
		t.Fatalf("third frame = %+v, want span from array batch", f)
	}
	if f := mustNext(t, s); !isSessionEndFrame(f) {
		t.Fatalf("fourth frame = %+v, want session end frame delivered", f)
	}
	if _, err := s.Next(context.Background()); err != io.EOF {
		t.Fatalf("after end frame: err = %v, want io.EOF", err)
	}
	if !s.Ended() {
		t.Fatal("Ended() = false after terminal end frame")
	}
	// EOF is sticky.
	if _, err := s.Next(context.Background()); err != io.EOF {
		t.Fatalf("second read after end: err = %v, want io.EOF", err)
	}
}

func TestStreamTolerantReader(t *testing.T) {
	d := &scriptDialer{script: []func() (rawConn, error){func() (rawConn, error) {
		return &scriptConn{messages: msgs(
			unknownKindMsg,         // unknown kind: delivered as the unknown variant
			`{"foo":1}`,            // no kind, no type: dropped
			`{"type":"SOMETHING"}`, // unknown transport frame: dropped
			`not json`,             // garbage: dropped
			endFrameJSON,
		)}, nil
	}}}
	s := newTestStream(t, d)

	unknown := mustNext(t, s)
	if unknown.Kind != "metric" || unknown.GetSpan() != nil || unknown.GetLog() != nil {
		t.Fatalf("unknown-kind frame = %+v, want Kind=metric with nil variants", unknown)
	}
	raw, err := json.Marshal(unknown)
	if err != nil {
		t.Fatalf("marshal unknown frame: %v", err)
	}
	if string(raw) != unknownKindMsg {
		t.Fatalf("unknown frame raw JSON = %s, want the original bytes", raw)
	}
	// Everything dropped in between, straight to the end frame.
	if f := mustNext(t, s); !isSessionEndFrame(f) {
		t.Fatalf("next frame = %+v, want the end frame", f)
	}
	if _, err := s.Next(context.Background()); err != io.EOF {
		t.Fatalf("err = %v, want io.EOF", err)
	}
}

func TestStreamReconnectPresentsCursor(t *testing.T) {
	d := &scriptDialer{script: []func() (rawConn, error){
		func() (rawConn, error) {
			return &scriptConn{
				messages: msgs(`{"type":"CURSOR","cursor":"c-42"}`, sdkSpanJSON),
				err:      websocket.CloseError{Code: websocket.StatusGoingAway},
			}, nil
		},
		func() (rawConn, error) {
			return &scriptConn{messages: msgs(endFrameJSON)}, nil
		},
	}}
	s := newTestStream(t, d)

	if f := mustNext(t, s); f.GetSpan() == nil {
		t.Fatalf("frame before drop = %+v, want span", f)
	}
	if f := mustNext(t, s); !isSessionEndFrame(f) {
		t.Fatalf("frame after reattach = %+v, want end frame", f)
	}
	if len(d.urls) != 2 {
		t.Fatalf("dials = %d, want 2", len(d.urls))
	}
	if !strings.Contains(d.urls[0], "resume=1") || strings.Contains(d.urls[0], "cursor=") {
		t.Fatalf("first attach URL = %q, want resume=1 and no cursor", d.urls[0])
	}
	if !strings.Contains(d.urls[1], "cursor=c-42") || !strings.Contains(d.urls[1], "resume=1") {
		t.Fatalf("reattach URL = %q, want resume=1 and cursor=c-42", d.urls[1])
	}
}

func TestStreamResyncRequiredMarksGap(t *testing.T) {
	d := &scriptDialer{script: []func() (rawConn, error){
		func() (rawConn, error) {
			return &scriptConn{messages: msgs(`{"type":"CURSOR","cursor":"c-1"}`), err: websocket.CloseError{Code: websocket.StatusTryAgainLater}}, nil
		},
		func() (rawConn, error) {
			return &scriptConn{messages: msgs(`{"type":"RESYNC_REQUIRED"}`, endFrameJSON)}, nil
		},
	}}
	s := newTestStream(t, d)

	if f := mustNext(t, s); !isSessionEndFrame(f) {
		t.Fatalf("frame = %+v, want end frame (transport frames not yielded)", f)
	}
	if !s.Gapped() {
		t.Fatal("Gapped() = false after RESYNC_REQUIRED")
	}
}

func TestStreamRedialBudgetExhausted(t *testing.T) {
	script := []func() (rawConn, error){func() (rawConn, error) {
		return &scriptConn{err: websocket.CloseError{Code: websocket.StatusInternalError}}, nil
	}}
	for i := 0; i < _maxRedials; i++ {
		script = append(script, func() (rawConn, error) { return nil, errors.New("dial refused") })
	}
	d := &scriptDialer{script: script}
	s := newTestStream(t, d)

	_, err := s.Next(context.Background())
	if !IsConnection(err) {
		t.Fatalf("err = %v, want CodeConnection", err)
	}
	if len(d.urls) != 1+_maxRedials {
		t.Fatalf("dials = %d, want initial + %d redials", len(d.urls), _maxRedials)
	}
	// Terminal errors are sticky.
	if _, err2 := s.Next(context.Background()); !IsConnection(err2) {
		t.Fatalf("second err = %v, want the sticky CodeConnection", err2)
	}
}

func TestStreamReattachClassification(t *testing.T) {
	cases := []struct {
		name   string
		status int
		check  func(error) bool
		code   Code
	}{
		{name: "403 means the session ended", status: 403, check: IsSessionEnded, code: CodeSessionEnded},
		{name: "401 means the token is bad", status: 401, check: IsUnauthorized, code: CodeUnauthorized},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			d := &scriptDialer{script: []func() (rawConn, error){
				func() (rawConn, error) {
					return &scriptConn{err: websocket.CloseError{Code: websocket.StatusGoingAway}}, nil
				},
				func() (rawConn, error) {
					return nil, &dialHTTPError{status: tc.status, err: errors.New("upgrade refused")}
				},
			}}
			s := newTestStream(t, d)
			_, err := s.Next(context.Background())
			if !tc.check(err) {
				t.Fatalf("err = %v, want %s", err, tc.code)
			}
		})
	}
}

func TestStreamNormalCloseWithoutEndFrame(t *testing.T) {
	d := &scriptDialer{script: []func() (rawConn, error){func() (rawConn, error) {
		return &scriptConn{messages: msgs(sdkSpanJSON), err: websocket.CloseError{Code: websocket.StatusNormalClosure}}, nil
	}}}
	s := newTestStream(t, d)

	mustNext(t, s)
	_, err := s.Next(context.Background())
	if !IsSessionEnded(err) {
		t.Fatalf("err = %v, want CodeSessionEnded", err)
	}
	if s.Ended() {
		t.Fatal("Ended() = true, but the terminal end frame was never observed")
	}
}

func TestStreamLegacySessionRootEndsCleanly(t *testing.T) {
	d := &scriptDialer{script: []func() (rawConn, error){func() (rawConn, error) {
		return &scriptConn{messages: msgs(legacyEndJSON)}, nil
	}}}
	s := newTestStream(t, d)

	if f := mustNext(t, s); !isSessionEndFrame(f) {
		t.Fatalf("frame = %+v, want legacy phone_session end recognized as terminal", f)
	}
	if _, err := s.Next(context.Background()); err != io.EOF {
		t.Fatalf("err = %v, want io.EOF", err)
	}
}

func TestTailClassifiesInitialDial(t *testing.T) {
	// Tail against a dead session fails fast: exercised through the package
	// seam by classifying the same dialHTTPError dialWS produces.
	err := classifyDialErr(&dialHTTPError{status: 403, err: errors.New("refused")})
	if !IsSessionEnded(err) {
		t.Fatalf("classifyDialErr(403) = %v, want CodeSessionEnded", err)
	}
	if e := classifyDialErr(errors.New("connection refused")); !IsConnection(e) {
		t.Fatalf("classifyDialErr(net) = %v, want retryable CodeConnection", e)
	}
}
