// Package telemetry is the read-only client for a session's telemetry: the
// live WebSocket leg allocate returns in telemetry_url, and helpers that
// reconstruct the dashboard's trace view from the durable frame archive
// (sessions_list_frames). Live and archive speak the same unified frame
// envelope and differ only in cardinality — the live stream carries start and
// end phases of a span, the archive one completed frame per span.
package telemetry

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"math/rand/v2"
	"net/http"
	"net/url"
	"time"

	"github.com/coder/websocket"

	platformgo "github.com/axilioai/platform-go"
)

// Frame is one telemetry frame in the unified envelope: a span or a log,
// discriminated on kind, with unknown kinds surfaced as the explicit unknown
// variant (Kind set, both pointers nil, raw JSON retained) per the tolerant
// reader contract.
type Frame = platformgo.RunSessionFramesResponseFramesItem

// Span types of the session-root span. The retired phone_session value is
// still emitted for traces recorded before the 2026-08-21 vocabulary cutover;
// a tolerant reader treats both as the root.
const (
	spanTypeSession       = "session"
	spanTypeSessionLegacy = "phone_session"
	spanPhaseEnd          = "end"
)

// Telemetry-leg transport frames (AXI-1740): type-tagged JSON below the frame
// envelope, disjoint from every frame kind so a non-upgraded viewer drops
// them as unknown. CURSOR checkpoints delivery (presented back on reattach);
// RESYNC_REQUIRED says the presented cursor fell off the live window, so the
// archive is the only complete record of what was missed.
const (
	transportFrameCursor         = "CURSOR"
	transportFrameResyncRequired = "RESYNC_REQUIRED"
)

// Redial budget: bounded exponential backoff with full jitter, same schedule
// as the control leg. The budget counts consecutive failures — a successful
// reattach resets it, so a long tail survives any number of isolated drops.
const (
	_maxRedials = 6
	_redialBase = 250 * time.Millisecond
	_redialCap  = 8 * time.Second
)

// _readLimit matches the hub's write side: frames are small, but the limit
// stays aligned with the control leg's so a burst of array-batched frames
// never trips it.
const _readLimit = 16 << 20

// rawConn is the minimal WebSocket surface the stream needs, so tests can
// inject a scripted connection without a real socket. The telemetry leg is
// read-only (the hub swallows anything a viewer sends), so there is no send.
type rawConn interface {
	recv(ctx context.Context) ([]byte, error)
	closeConn() error
}

type dialFunc func(ctx context.Context, url string) (rawConn, error)

// Stream is a live tail of one session's telemetry. Next yields frames in
// arrival order; the stream reconnects transparently on transient drops
// (presenting its cursor so delivery resumes where it left off) and ends
// cleanly — io.EOF — after the terminal session end frame, the platform's
// only session-end signal (peer_disconnected is never emitted).
//
// Next and Close may be called from different goroutines; Next itself is not
// safe for concurrent use.
type Stream struct {
	url         string
	dial        dialFunc
	redialDelay func(attempt int) time.Duration

	conn    rawConn
	cursor  string
	pending []*Frame
	ended   bool // terminal session end frame observed
	eof     bool
	gapped  bool
	closed  bool
	err     error // sticky terminal error
}

// Tail attaches to a session's live telemetry leg. telemetryURL is the value
// allocate returned (wss …/realtime/ws/telemetry, credential embedded); the
// dial fails fast with CodeUnauthorized / CodeSessionEnded when the server
// refuses the attach.
func Tail(ctx context.Context, telemetryURL string) (*Stream, error) {
	s := &Stream{
		url:         telemetryURL,
		dial:        dialWS,
		redialDelay: redialDelay,
	}
	conn, err := s.dial(ctx, s.attachURL())
	if err != nil {
		return nil, classifyDialErr(err)
	}
	s.conn = conn
	return s, nil
}

// Next returns the next frame. It returns io.EOF after the terminal session
// end frame has been delivered (check Ended to distinguish that clean end
// from a Close), the context's own error when ctx ends, and a terminal
// *Error otherwise. After a non-nil error every subsequent call returns the
// same result.
func (s *Stream) Next(ctx context.Context) (*Frame, error) {
	for {
		if len(s.pending) > 0 {
			frame := s.pending[0]
			s.pending = s.pending[1:]
			if isSessionEndFrame(frame) {
				s.ended = true
			}
			return frame, nil
		}
		if s.err != nil {
			return nil, s.err
		}
		if s.eof {
			return nil, io.EOF
		}
		if s.ended {
			// The end frame was delivered and nothing is buffered behind
			// it: the session is over, so the socket is done too.
			s.eof = true
			_ = s.closeConn()
			return nil, io.EOF
		}
		if s.closed {
			s.err = &Error{Code: CodeClosed, Message: "stream closed"}
			return nil, s.err
		}

		data, err := s.conn.recv(ctx)
		if err != nil {
			if fail := s.recvErr(ctx, err); fail != nil {
				return nil, fail
			}
			continue // reattached; read again
		}
		s.ingest(data)
	}
}

// Ended reports whether the terminal session end frame was observed. False
// after an io.EOF means the stream ended without it (Close, or a normal
// server close that raced the frame) — the durable archive is then the
// authoritative record.
func (s *Stream) Ended() bool { return s.ended }

// Gapped reports whether a reattach fell off the live retention window
// (RESYNC_REQUIRED): frames were missed and the live stream alone is no
// longer the full trace. Refetch the archive (Session.Trace) to true up.
func (s *Stream) Gapped() bool { return s.gapped }

// Close tears the stream down. A blocked Next returns with CodeClosed.
func (s *Stream) Close() error {
	s.closed = true
	return s.closeConn()
}

func (s *Stream) closeConn() error {
	if s.conn == nil {
		return nil
	}
	err := s.conn.closeConn()
	s.conn = nil
	return err
}

// recvErr handles one failed read: context errors and terminal conditions
// come back as the error to surface (sticky), a transient drop runs the
// bounded redial loop and returns nil once reattached.
func (s *Stream) recvErr(ctx context.Context, err error) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}
	if s.closed {
		s.err = &Error{Code: CodeClosed, Message: "stream closed"}
		return s.err
	}
	if websocket.CloseStatus(err) == websocket.StatusNormalClosure {
		// A normal close without the end frame: the session is over, but
		// the terminal frame is best-effort and this one lost the race (or
		// was dropped edge-side). The archive holds whatever was recorded.
		s.err = &Error{Code: CodeSessionEnded, Message: "telemetry websocket closed before the session end frame: " + err.Error()}
		return s.err
	}
	// Everything else — 1001 going away, 1013 try again later, 1011
	// internal error, abrupt TCP loss — is a transient drop the redial
	// loop recovers.
	_ = s.closeConn()
	for attempt := 1; ; attempt++ {
		if attempt > _maxRedials {
			s.err = &Error{Code: CodeConnection, Message: "telemetry websocket lost and could not be re-established", Retryable: true}
			return s.err
		}
		if sleepCtx(ctx, s.redialDelay(attempt-1)) != nil {
			return ctx.Err()
		}
		conn, dialErr := s.dial(ctx, s.attachURL())
		if dialErr == nil {
			s.conn = conn
			return nil
		}
		classified := classifyDialErr(dialErr)
		var e *Error
		if errors.As(classified, &e) && !e.Retryable {
			s.err = classified
			return s.err
		}
	}
}

// ingest parses one wire message: a single frame object, a JSON array of
// frame objects, or a transport frame. Malformed or unrecognized content is
// dropped, never an error — the tolerant reader contract.
func (s *Stream) ingest(data []byte) {
	trimmed := bytes.TrimSpace(data)
	if len(trimmed) == 0 {
		return
	}
	if trimmed[0] == '[' {
		var elements []json.RawMessage
		if err := json.Unmarshal(trimmed, &elements); err != nil {
			return
		}
		for _, element := range elements {
			s.ingestObject(element)
		}
		return
	}
	s.ingestObject(trimmed)
}

func (s *Stream) ingestObject(data []byte) {
	// Transport frames are type-tagged and kind-less, so probing type first
	// cleanly separates them from the frame envelope.
	var probe struct {
		Type   string `json:"type"`
		Cursor string `json:"cursor"`
	}
	if err := json.Unmarshal(data, &probe); err != nil {
		return
	}
	switch probe.Type {
	case transportFrameCursor:
		if probe.Cursor != "" {
			s.cursor = probe.Cursor
		}
		return
	case transportFrameResyncRequired:
		s.gapped = true
		return
	}
	if probe.Type != "" {
		return // unknown transport frame: drop, per the tolerant reader
	}
	frame := new(Frame)
	if err := json.Unmarshal(data, frame); err != nil {
		// No kind discriminant at all — not a frame. Drop rather than
		// erroring a live stream over one unrecognized message. A frame
		// with an unknown kind is NOT this case: it unmarshals into the
		// explicit unknown variant and is delivered.
		return
	}
	s.pending = append(s.pending, frame)
}

// isSessionEndFrame recognizes the terminal frame: the session-root span
// (neutral or legacy name) closing.
func isSessionEndFrame(f *Frame) bool {
	span := f.GetSpan()
	if span == nil {
		return false
	}
	if span.SpanType != spanTypeSession && span.SpanType != spanTypeSessionLegacy {
		return false
	}
	return span.Phase == spanPhaseEnd || span.EndTimeUnixNano > 0
}

// attachURL is the telemetry URL plus the resume params: every attach opts
// in to cursor checkpoints (resume=1), and a reattach that holds a cursor
// presents it so delivery continues where this client left off.
func (s *Stream) attachURL() string {
	u, err := url.Parse(s.url)
	if err != nil {
		// An unparsable URL fails at dial with a real error; don't mask it.
		return s.url
	}
	q := u.Query()
	q.Set("resume", "1")
	if s.cursor != "" {
		q.Set("cursor", s.cursor)
	}
	u.RawQuery = q.Encode()
	return u.String()
}

// redialDelay returns the wait before redial attempt (0-based): uniformly
// random in (0, min(cap, base·2^attempt)]. Full jitter so a fleet of viewers
// that lost the same pod doesn't reconnect in lockstep.
func redialDelay(attempt int) time.Duration {
	limit := _redialBase
	for i := 0; i < attempt && limit < _redialCap; i++ {
		limit <<= 1
	}
	if limit > _redialCap {
		limit = _redialCap
	}
	return time.Duration(rand.Int64N(int64(limit))) + 1
}

// sleepCtx waits d or until ctx ends, returning ctx's error in that case.
func sleepCtx(ctx context.Context, d time.Duration) error {
	timer := time.NewTimer(d)
	defer timer.Stop()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return nil
	}
}

// --- real WebSocket ---------------------------------------------------------

// dialHTTPError carries the HTTP status of a refused upgrade so the attach
// classification can tell a dead session (403) from a bad token (401) from a
// transient failure.
type dialHTTPError struct {
	status int
	err    error
}

func (e *dialHTTPError) Error() string { return e.err.Error() }
func (e *dialHTTPError) Unwrap() error { return e.err }

func dialWS(ctx context.Context, url string) (rawConn, error) {
	c, resp, err := websocket.Dial(ctx, url, nil) //nolint:bodyclose // coder/websocket owns resp.Body.
	if err != nil {
		if resp != nil {
			return nil, &dialHTTPError{status: resp.StatusCode, err: err}
		}
		return nil, err
	}
	c.SetReadLimit(_readLimit)
	return &wsRawConn{c: c}, nil
}

// classifyDialErr maps a failed (re)attach onto the contract's out-of-band
// liveness signal: 403 means the allocation is no longer active — the
// allocation-scoped token dies at deallocation regardless of its exp — and
// 401 means the token is bad; both terminal. Anything else is transient.
func classifyDialErr(err error) error {
	var de *dialHTTPError
	if errors.As(err, &de) {
		switch de.status {
		case http.StatusForbidden:
			return &Error{Code: CodeSessionEnded, Message: "session is no longer active: " + err.Error()}
		case http.StatusUnauthorized:
			return &Error{Code: CodeUnauthorized, Message: "telemetry token rejected: " + err.Error()}
		}
	}
	return &Error{Code: CodeConnection, Message: "cannot connect to telemetry websocket: " + err.Error(), Retryable: true}
}

type wsRawConn struct{ c *websocket.Conn }

func (w *wsRawConn) recv(ctx context.Context) ([]byte, error) {
	_, data, err := w.c.Read(ctx)
	return data, err
}

func (w *wsRawConn) closeConn() error {
	return w.c.Close(websocket.StatusNormalClosure, "")
}
