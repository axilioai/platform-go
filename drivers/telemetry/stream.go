// Package telemetry is the read-only client for a session's telemetry: the
// live event stream (telemetry_url, Server-Sent Events on the realtime host)
// and helpers that reconstruct the dashboard's trace view from the durable
// frame archive (sessions_list_frames). Live and archive speak the same
// unified frame envelope and differ only in cardinality: the live stream
// carries start and end phases of a span, the archive one completed frame per
// span.
package telemetry

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"io"
	"math/rand/v2"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

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

// Reopen budget: bounded exponential backoff with full jitter, same schedule
// as the control leg. The budget counts consecutive failures, so a long tail
// survives any number of isolated drops.
const (
	_maxRedials = 6
	_redialBase = 250 * time.Millisecond
	_redialCap  = 8 * time.Second
)

// _maxLineBytes caps one line of the stream. Client protection only: frames
// are small, and this keeps a broken server from growing the buffer forever.
const _maxLineBytes = 16 << 20

// _idleTimeout ends a stream that has gone silent. The server sends a
// heartbeat every 15 s, so silence this long means a dead connection, and
// the stream reopens from its cursor.
const _idleTimeout = 45 * time.Second

// event is one server-sent event. hasID is set when the event carried an id
// field, which moves the cursor even when there is no data.
type event struct {
	id, name, data string
	hasID          bool
}

// eventSource reads the events of one stream response.
type eventSource interface {
	next() (event, error)
	close() error
}

// openFunc opens the stream, resuming after cursor when it is set.
type openFunc func(ctx context.Context, url, cursor string) (eventSource, error)

// Stream is a live tail of one session's telemetry. Next yields frames in
// order; the stream reopens transparently on a drop (sending its cursor as
// Last-Event-ID, so delivery resumes with no gap) and ends cleanly with
// io.EOF after the session end frame, or when the server reports the session
// over with nothing left to send. The server backfills the archive on a fresh
// stream and deduplicates, so each frame arrives once.
//
// Next and Close may be called from different goroutines; Next itself is not
// safe for concurrent use.
type Stream struct {
	url         string
	open        openFunc
	redialDelay func(attempt int) time.Duration

	// mu guards src and closed, the only state Close shares with Next. Next
	// never holds it across a blocking read, so Close can always interrupt one.
	mu     sync.Mutex
	src    eventSource
	closed bool

	// The rest belongs to the goroutine calling Next.
	cursor string
	ended  bool // the session is over: end frame seen, or 204 from the server
	eof    bool
	err    error // sticky terminal error
}

// TailOption configures Tail.
type TailOption func(*tailConfig)

type tailConfig struct {
	header http.Header
	client *http.Client
}

// WithHeader sends header on every request of the stream, for example the
// X-Axilio-Api-Key credential when the URL carries no token.
func WithHeader(header http.Header) TailOption {
	return func(c *tailConfig) { c.header = header.Clone() }
}

// WithHTTPClient sets the HTTP client. It must not set a total Timeout: the
// stream is long-lived.
func WithHTTPClient(client *http.Client) TailOption {
	return func(c *tailConfig) { c.client = client }
}

// StreamURL is the event stream of a session on the API host (apiBaseURL
// ends in /api/v1), for callers that hold a session id and an API key but no
// telemetry_url. Send the key with WithHeader.
func StreamURL(apiBaseURL, sessionID string) string {
	return strings.TrimRight(apiBaseURL, "/") + "/realtime/sessions/" + url.PathEscape(sessionID) + "/events"
}

// Tail opens a session's telemetry stream. telemetryURL is the session's
// telemetry_url (its read-only token embedded) or StreamURL with an API key
// header. It fails fast with CodeUnauthorized, CodeNotFound or
// CodeTelemetryOff when the server refuses the stream.
func Tail(ctx context.Context, telemetryURL string, opts ...TailOption) (*Stream, error) {
	cfg := tailConfig{client: &http.Client{}}
	for _, opt := range opts {
		opt(&cfg)
	}
	return start(ctx, &Stream{url: telemetryURL, open: httpOpener(cfg), redialDelay: redialDelay})
}

// start opens s for the first time.
func start(ctx context.Context, s *Stream) (*Stream, error) {
	src, err := s.open(ctx, s.url, "")
	if isSessionOver(err) {
		s.ended, s.eof = true, true
		return s, nil
	}
	if err != nil {
		return nil, classifyOpenErr(err)
	}
	s.src = src
	return s, nil
}

// Next returns the next frame. It returns io.EOF once the session is over
// (check Ended to tell that from a Close), the context's own error when ctx
// ends, and a terminal *Error otherwise. After a non-nil error every
// subsequent call returns the same result.
func (s *Stream) Next(ctx context.Context) (*Frame, error) {
	for {
		if s.err != nil {
			return nil, s.err
		}
		if s.eof {
			return nil, io.EOF
		}
		if s.ended {
			s.eof = true
			_ = s.closeSrc()
			return nil, io.EOF
		}
		src, live := s.liveSrc()
		if !live {
			s.err = errClosed()
			return nil, s.err
		}

		ev, err := readEvent(ctx, src)
		if err != nil {
			if fail := s.readErr(ctx, err); fail != nil {
				return nil, fail
			}
			continue // reopened, or the session turned out to be over
		}
		if ev.hasID {
			s.cursor = ev.id
		}
		if ev.data == "" {
			continue // a cursor-only event
		}
		frame := new(Frame)
		if err := json.Unmarshal([]byte(ev.data), frame); err != nil {
			// Not a frame at all (no kind). Drop it rather than fail a live
			// stream over one message; an unknown kind is NOT this case: it
			// unmarshals into the explicit unknown variant and is delivered.
			continue
		}
		if isSessionEndFrame(frame) {
			s.ended = true
		}
		return frame, nil
	}
}

// readEvent reads one event, unblocking when ctx ends.
func readEvent(ctx context.Context, src eventSource) (event, error) {
	stop := context.AfterFunc(ctx, func() { _ = src.close() })
	defer stop()
	return src.next()
}

// Ended reports whether the session is over: the session end frame was
// delivered, or the server answered a reopen with nothing left to send.
// False after an io.EOF means the stream was closed locally.
func (s *Stream) Ended() bool { return s.ended }

// Close tears the stream down. A blocked Next returns with CodeClosed.
func (s *Stream) Close() error {
	s.markClosed()
	return s.closeSrc()
}

// closeSrc detaches the response, then closes it with no lock held.
func (s *Stream) closeSrc() error {
	src := s.takeSrc()
	if src == nil {
		return nil
	}
	return src.close()
}

// The state helpers below are the only code that touches mu. Each locks with
// defer and does one state read or write, so no path can leave mu held and no
// network I/O runs under it.

func (s *Stream) markClosed() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.closed = true
}

func (s *Stream) isClosed() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.closed
}

func (s *Stream) liveSrc() (eventSource, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.src, !s.closed
}

func (s *Stream) takeSrc() eventSource {
	s.mu.Lock()
	defer s.mu.Unlock()
	src := s.src
	s.src = nil
	return src
}

// attach installs a reopened response and reports false, leaving it to the
// caller, if Close has already run.
func (s *Stream) attach(src eventSource) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.closed {
		return false
	}
	s.src = src
	return true
}

func errClosed() *Error { return &Error{Code: CodeClosed, Message: "stream closed"} }

// readErr handles one failed read: a context error or a terminal condition
// comes back as the error to surface (sticky); a drop runs the bounded reopen
// loop and returns nil once reopened or once the server says the session is
// over.
func (s *Stream) readErr(ctx context.Context, err error) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}
	if s.isClosed() {
		s.err = errClosed()
		return s.err
	}
	_ = s.closeSrc()
	for attempt := 1; ; attempt++ {
		if s.isClosed() {
			s.err = errClosed()
			return s.err
		}
		if attempt > _maxRedials {
			s.err = &Error{Code: CodeConnection, Message: "telemetry stream lost and could not be re-established", Retryable: true}
			return s.err
		}
		if sleepCtx(ctx, s.redialDelay(attempt-1)) != nil {
			return ctx.Err()
		}
		src, openErr := s.open(ctx, s.url, s.cursor)
		if isSessionOver(openErr) {
			s.ended = true
			return nil
		}
		if openErr == nil {
			if !s.attach(src) {
				_ = src.close()
				s.err = errClosed()
				return s.err
			}
			return nil
		}
		classified := classifyOpenErr(openErr)
		var e *Error
		if errors.As(classified, &e) && !e.Retryable {
			s.err = classified
			return s.err
		}
	}
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
	return span.Phase == spanPhaseEnd || spanEndNano(span) > 0
}

// redialDelay returns the wait before reopen attempt (0-based): uniformly
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

// --- HTTP -------------------------------------------------------------------

// statusError carries the HTTP status of a refused stream.
type statusError struct{ status int }

func (e *statusError) Error() string { return "telemetry stream answered " + http.StatusText(e.status) }

// isSessionOver reports a 204: the session has ended and nothing is left
// after the cursor.
func isSessionOver(err error) bool {
	var se *statusError
	return errors.As(err, &se) && se.status == http.StatusNoContent
}

// httpOpener opens the stream with a GET. The response outlives the ctx it
// was opened with (Next takes its own), so the request runs on a detached
// context that ctx can cancel only until the headers arrive.
func httpOpener(cfg tailConfig) openFunc {
	return func(ctx context.Context, url, cursor string) (eventSource, error) {
		reqCtx, cancel := context.WithCancel(context.WithoutCancel(ctx))
		stop := context.AfterFunc(ctx, cancel)
		defer stop()
		req, err := http.NewRequestWithContext(reqCtx, http.MethodGet, url, nil)
		if err != nil {
			cancel()
			return nil, err
		}
		for k, v := range cfg.header {
			req.Header[k] = v
		}
		req.Header.Set("Accept", "text/event-stream")
		if cursor != "" {
			req.Header.Set("Last-Event-ID", cursor)
		}
		resp, err := cfg.client.Do(req)
		if err != nil {
			cancel()
			return nil, err
		}
		if resp.StatusCode != http.StatusOK {
			_ = resp.Body.Close()
			cancel()
			return nil, &statusError{status: resp.StatusCode}
		}
		return newSSEReader(resp.Body, cancel), nil
	}
}

// classifyOpenErr maps a refused (re)open onto the error taxonomy: 401 is a
// bad credential, 404 an unknown session (or another org's), 409 a session
// with telemetry off; all terminal. Anything else is transient.
func classifyOpenErr(err error) error {
	var se *statusError
	if errors.As(err, &se) {
		switch se.status {
		case http.StatusUnauthorized, http.StatusForbidden:
			return &Error{Code: CodeUnauthorized, Message: "telemetry credential rejected: " + err.Error()}
		case http.StatusNotFound:
			return &Error{Code: CodeNotFound, Message: "session not found: " + err.Error()}
		case http.StatusConflict:
			return &Error{Code: CodeTelemetryOff, Message: "telemetry is off for this session"}
		}
	}
	return &Error{Code: CodeConnection, Message: "cannot open telemetry stream: " + err.Error(), Retryable: true}
}

// sseReader parses a text/event-stream body. A stream that stays silent past
// _idleTimeout is closed, which surfaces as a read error and a reopen.
type sseReader struct {
	body   io.ReadCloser
	cancel context.CancelFunc
	sc     *bufio.Scanner
	idle   *time.Timer
}

func newSSEReader(body io.ReadCloser, cancel context.CancelFunc) *sseReader {
	sc := bufio.NewScanner(body)
	sc.Buffer(make([]byte, 0, 64<<10), _maxLineBytes)
	r := &sseReader{body: body, cancel: cancel, sc: sc}
	r.idle = time.AfterFunc(_idleTimeout, func() { _ = r.close() })
	return r
}

// next reads one event per the SSE spec: field lines up to a blank line,
// comments (":") ignored, multi-line data joined with "\n".
func (r *sseReader) next() (event, error) {
	var (
		ev   event
		data []string
		seen bool
	)
	for r.sc.Scan() {
		r.idle.Reset(_idleTimeout)
		line := r.sc.Text()
		if line == "" {
			if !seen {
				continue
			}
			ev.data = strings.Join(data, "\n")
			return ev, nil
		}
		if strings.HasPrefix(line, ":") {
			continue
		}
		field, value, _ := strings.Cut(line, ":")
		value = strings.TrimPrefix(value, " ")
		seen = true
		switch field {
		case "id":
			ev.id, ev.hasID = value, true
		case "event":
			ev.name = value
		case "data":
			data = append(data, value)
		}
	}
	if err := r.sc.Err(); err != nil {
		return event{}, err
	}
	return event{}, io.EOF
}

func (r *sseReader) close() error {
	r.idle.Stop()
	r.cancel()
	return r.body.Close()
}
