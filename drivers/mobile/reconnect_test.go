package mobile

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

// matrixConn is a scripted rawConn for the reconnect matrix: each send is
// decoded and handed to respond, whose replies (responses and notification
// frames alike) are queued; once every queued frame is drained, recv fails
// with recvErr, the scripted force-close.
type matrixConn struct {
	respond   func(cmd dcpCommand) []dcpResponse
	recvErr   error
	preloaded [][]byte // frames delivered before any queued reply (stale replays)
	inbox     [][]byte
	sent      []dcpCommand
}

func (c *matrixConn) send(_ context.Context, data []byte) error {
	var cmd dcpCommand
	if err := json.Unmarshal(data, &cmd); err != nil {
		return err
	}
	c.sent = append(c.sent, cmd)
	if c.respond != nil {
		for _, r := range c.respond(cmd) {
			b, err := json.Marshal(r)
			if err != nil {
				return err
			}
			c.inbox = append(c.inbox, b)
		}
	}
	return nil
}

func (c *matrixConn) recv(_ context.Context) ([]byte, error) {
	if len(c.preloaded) > 0 {
		b := c.preloaded[0]
		c.preloaded = c.preloaded[1:]
		return b, nil
	}
	if len(c.inbox) > 0 {
		b := c.inbox[0]
		c.inbox = c.inbox[1:]
		return b, nil
	}
	if c.recvErr != nil {
		return nil, c.recvErr
	}
	return nil, errors.New("scripted conn exhausted")
}

func (c *matrixConn) closeConn() error { return nil }

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

func conns(cs ...rawConn) []func() (rawConn, error) {
	out := make([]func() (rawConn, error), len(cs))
	for i, c := range cs {
		c := c
		out[i] = func() (rawConn, error) { return c, nil }
	}
	return out
}

func testTransport(d *scriptDialer) *RemoteTransport {
	return &RemoteTransport{
		url:         "wss://connect.test/api/v1/realtime/ws/control?token=abc",
		openTimeout: time.Second,
		dial:        d.dial,
		redialDelay: func(int) time.Duration { return time.Millisecond },
	}
}

func okAll() func(cmd dcpCommand) []dcpResponse {
	return func(cmd dcpCommand) []dcpResponse {
		return []dcpResponse{{ID: cmd.ID, Result: json.RawMessage(`{}`)}}
	}
}

func notification(t *testing.T, method string, params any) []byte {
	t.Helper()
	raw, err := json.Marshal(params)
	if err != nil {
		t.Fatalf("marshal notification params: %v", err)
	}
	b, err := json.Marshal(dcpResponse{Method: method, Params: raw})
	if err != nil {
		t.Fatalf("marshal notification: %v", err)
	}
	return b
}

func keyOf(t *testing.T, cmd dcpCommand) string {
	t.Helper()
	var p struct {
		IdempotencyKey string `json:"idempotencyKey"`
	}
	if err := json.Unmarshal(cmd.Params, &p); err != nil {
		t.Fatalf("params not an object: %v", err)
	}
	return p.IdempotencyKey
}

// The force-close matrix (the design doc's original SDK kill criterion):
// each retryable close class (1001 going away, 1013 try again later, 1011
// internal error, abrupt TCP loss) must be survived by a redial that
// re-sends the interrupted command under a fresh id with the SAME
// idempotency key, so the caller's Tap returns success exactly once. Each
// terminal class (1000 session ended / superseded, 4409 control held)
// must surface its code without a single redial.
func TestReconnect_ForceCloseMatrix(t *testing.T) {
	retryable := []struct {
		name string
		err  error
	}{
		{"1001 going away", websocket.CloseError{Code: websocket.StatusGoingAway, Reason: "Server shutting down"}},
		{"1013 try again later", websocket.CloseError{Code: websocket.StatusTryAgainLater, Reason: "Server shutting down"}},
		{"1011 internal error", websocket.CloseError{Code: websocket.StatusInternalError, Reason: "welcome send failed"}},
		{"abrupt loss", io.ErrUnexpectedEOF},
	}
	for _, c := range retryable {
		t.Run(c.name, func(t *testing.T) {
			conn1 := &matrixConn{recvErr: c.err} // swallows the command, then dies
			conn2 := &matrixConn{respond: okAll()}
			dialer := &scriptDialer{script: conns(conn1, conn2)}
			d := newDriver(testTransport(dialer))

			if err := d.Tap(Coords{X: 120, Y: 640}); err != nil {
				t.Fatalf("Tap after %s = %v, want success via redial", c.name, err)
			}
			if len(dialer.urls) != 2 {
				t.Fatalf("dials = %d, want 2 (initial + one redial)", len(dialer.urls))
			}
			if len(conn1.sent) != 1 || len(conn2.sent) != 1 {
				t.Fatalf("sends = %d/%d, want the one command on each connection", len(conn1.sent), len(conn2.sent))
			}
			orig, resend := conn1.sent[0], conn2.sent[0]
			if resend.Method != methodTouchTap {
				t.Fatalf("re-sent method = %q, want %q", resend.Method, methodTouchTap)
			}
			if key := keyOf(t, orig); key == "" || key != keyOf(t, resend) {
				t.Fatalf("idempotency keys: original %q, re-send %q; must be one non-empty key reused verbatim", key, keyOf(t, resend))
			}
			if resend.ID <= orig.ID {
				t.Fatalf("re-send id %d not greater than original %d: ids must stay monotonic across the redial (AXI-1293)", resend.ID, orig.ID)
			}
		})
	}

	terminal := []struct {
		name     string
		err      error
		wantCode Code
	}{
		{"1000 session ended", websocket.CloseError{Code: websocket.StatusNormalClosure, Reason: "session ended"}, CodeSessionEnded},
		{"1000 superseded", websocket.CloseError{Code: websocket.StatusNormalClosure, Reason: "Superseded"}, CodeSessionEnded},
		{"4409 control held", websocket.CloseError{Code: _closeControlHeld, Reason: "control_held:editor"}, CodeControlHeld},
	}
	for _, c := range terminal {
		t.Run(c.name, func(t *testing.T) {
			conn1 := &matrixConn{recvErr: c.err}
			dialer := &scriptDialer{script: conns(conn1, &matrixConn{respond: okAll()})}
			d := newDriver(testTransport(dialer))

			err := d.Tap(Coords{X: 1, Y: 2})
			e := asError(err)
			if e == nil || e.Code != c.wantCode {
				t.Fatalf("Tap after %s = %v, want code %s", c.name, err, c.wantCode)
			}
			if e.Retryable {
				t.Fatalf("%s marked retryable; it is terminal by contract", c.name)
			}
			if len(dialer.urls) != 1 {
				t.Fatalf("dials = %d after a terminal close, want 1 (never auto-retried)", len(dialer.urls))
			}
		})
	}
}

// Reattach HTTP status is the out-of-band liveness signal: 403 (allocation
// no longer active) ends the retry loop as CodeSessionEnded; 401 (bad
// token) as CodeUnauthorized.
func TestReconnect_ReattachStatusTerminal(t *testing.T) {
	cases := []struct {
		name     string
		status   int
		wantCode Code
	}{
		{"403 session over", 403, CodeSessionEnded},
		{"401 bad token", 401, CodeUnauthorized},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			conn1 := &matrixConn{recvErr: websocket.CloseError{Code: websocket.StatusGoingAway}}
			dialer := &scriptDialer{script: []func() (rawConn, error){
				func() (rawConn, error) { return conn1, nil },
				func() (rawConn, error) {
					return nil, &dialHTTPError{status: c.status, err: errors.New("upgrade refused")}
				},
			}}
			d := newDriver(testTransport(dialer))

			err := d.Tap(Coords{X: 1, Y: 2})
			e := asError(err)
			if e == nil || e.Code != c.wantCode || e.Retryable {
				t.Fatalf("Tap = %v, want terminal code %s", err, c.wantCode)
			}
			if len(dialer.urls) != 2 {
				t.Fatalf("dials = %d, want 2 (the refused reattach ends the loop)", len(dialer.urls))
			}
		})
	}
}

// Every attach opts in to checkpoints; a reattach presents the latest
// Axilio.cursor so the server resumes delivery where this client left off.
func TestReconnect_CursorTrackedAndPresented(t *testing.T) {
	calls := 0
	conn1 := &matrixConn{
		recvErr: websocket.CloseError{Code: websocket.StatusGoingAway},
		respond: func(cmd dcpCommand) []dcpResponse {
			calls++
			if calls == 1 {
				// Checkpoint interleaved before the response, as the hub
				// emits it per delivered batch.
				return []dcpResponse{
					{Method: methodAxilioCursor, Params: json.RawMessage(`{"cursor":"1755861234567-0"}`)},
					{ID: cmd.ID, Result: json.RawMessage(`{}`)},
				}
			}
			return nil // second command: swallowed, then recvErr fires
		},
	}
	conn2 := &matrixConn{respond: okAll()}
	dialer := &scriptDialer{script: conns(conn1, conn2)}
	d := newDriver(testTransport(dialer))

	if err := d.Tap(Coords{X: 1, Y: 2}); err != nil {
		t.Fatalf("first Tap: %v", err)
	}
	if err := d.Tap(Coords{X: 3, Y: 4}); err != nil {
		t.Fatalf("second Tap (across the drop): %v", err)
	}
	if len(dialer.urls) != 2 {
		t.Fatalf("dials = %d, want 2", len(dialer.urls))
	}
	first, second := dialer.urls[0], dialer.urls[1]
	if !strings.Contains(first, "resume=1") || strings.Contains(first, "cursor=") {
		t.Fatalf("first attach %q: want resume=1 and no cursor", first)
	}
	if !strings.Contains(second, "resume=1") || !strings.Contains(second, "cursor=1755861234567-0") {
		t.Fatalf("reattach %q: want resume=1 and the tracked cursor", second)
	}
}

// Axilio.resyncRequired means the held cursor predates the retained window;
// the transport must drop it rather than re-present a known-stale token.
func TestReconnect_ResyncClearsCursor(t *testing.T) {
	calls := 0
	conn1 := &matrixConn{
		recvErr: websocket.CloseError{Code: websocket.StatusGoingAway},
		preloaded: [][]byte{
			notification(t, methodAxilioCursor, map[string]string{"cursor": "5-0"}),
			notification(t, methodAxilioResyncRequired, map[string]string{"requested": "5-0", "oldest": "9-0"}),
		},
		respond: func(cmd dcpCommand) []dcpResponse {
			calls++
			if calls == 1 {
				return []dcpResponse{{ID: cmd.ID, Result: json.RawMessage(`{}`)}}
			}
			return nil // second command: swallowed, then recvErr fires
		},
	}
	conn2 := &matrixConn{respond: okAll()}
	dialer := &scriptDialer{script: conns(conn1, conn2)}
	d := newDriver(testTransport(dialer))

	if err := d.Tap(Coords{X: 1, Y: 2}); err != nil {
		t.Fatalf("first Tap: %v", err)
	}
	// conn1's replies are drained, so the next call hits recvErr and redials.
	if err := d.Tap(Coords{X: 3, Y: 4}); err != nil {
		t.Fatalf("second Tap: %v", err)
	}
	if got := dialer.urls[1]; strings.Contains(got, "cursor=") {
		t.Fatalf("reattach %q presented a cursor after a resync; it must be dropped", got)
	}
}

// AXI-1293 regression, resumed-connection flavor: a pre-drop response
// redelivered on the resumed socket must never be matched to the re-sent
// command: ids stay monotonic and the transport matches strictly by the
// fresh id.
func TestReconnect_StaleReplayedResponseSkipped(t *testing.T) {
	conn1 := &matrixConn{recvErr: websocket.CloseError{Code: websocket.StatusGoingAway}}
	var conn2 *matrixConn
	conn2 = &matrixConn{
		respond: func(cmd dcpCommand) []dcpResponse {
			return []dcpResponse{{ID: cmd.ID, Result: json.RawMessage(`{"source":"fresh"}`)}}
		},
	}
	dialer := &scriptDialer{script: conns(conn1, conn2)}
	tp := testTransport(dialer)

	// The stale replay: the pre-drop response for the ORIGINAL id arrives
	// first on the resumed socket (at-least-once redelivery).
	conn2.preloaded = [][]byte{[]byte(`{"id":1,"result":{"source":"stale-replay"}}`)}

	raw, err := tp.call(context.Background(), methodScreenObserve, observeParams{})
	if err != nil {
		t.Fatalf("call across drop: %v", err)
	}
	if string(raw) != `{"source":"fresh"}` {
		t.Fatalf("result = %s, want the fresh response, never the stale replay", raw)
	}
	if got := conn2.sent[0].ID; got <= conn1.sent[0].ID {
		t.Fatalf("re-send id %d not monotonic past original %d", got, conn1.sent[0].ID)
	}
}

// A handshake performed by the caller is replayed internally after every
// reattach, before the interrupted command resumes: capability state is
// per-connection.
func TestReconnect_HandshakeReplayedOnReattach(t *testing.T) {
	handshakeResult := json.RawMessage(`{"protocol_version":1,"device":{"device_id":"d1","platform":"android","form_factor":"phone","input_modalities":["touch"],"screen_width":1080,"screen_height":2400},"domains":["Touch"],"capabilities":["Touch.tap"]}`)
	respondAll := func(cmd dcpCommand) []dcpResponse {
		if cmd.Method == methodProtocolHandshake {
			return []dcpResponse{{ID: cmd.ID, Result: handshakeResult}}
		}
		return []dcpResponse{{ID: cmd.ID, Result: json.RawMessage(`{}`)}}
	}
	calls := 0
	conn1 := &matrixConn{
		recvErr: websocket.CloseError{Code: websocket.StatusGoingAway},
		respond: func(cmd dcpCommand) []dcpResponse {
			calls++
			if calls == 1 {
				return respondAll(cmd)
			}
			return nil // the tap that dies with the connection
		},
	}
	conn2 := &matrixConn{respond: respondAll}
	dialer := &scriptDialer{script: conns(conn1, conn2)}
	d := newDriver(testTransport(dialer))

	if _, err := d.Handshake(); err != nil {
		t.Fatalf("Handshake: %v", err)
	}
	if err := d.Tap(Coords{X: 1, Y: 2}); err != nil {
		t.Fatalf("Tap across drop: %v", err)
	}
	if len(conn2.sent) != 2 {
		t.Fatalf("resumed connection got %d frames, want handshake + tap", len(conn2.sent))
	}
	if conn2.sent[0].Method != methodProtocolHandshake || conn2.sent[1].Method != methodTouchTap {
		t.Fatalf("resumed order = %q, %q; want the handshake replayed before the re-sent command",
			conn2.sent[0].Method, conn2.sent[1].Method)
	}
}

// The redial budget is bounded: persistent dial failure surfaces a
// retryable connection error after initial + _maxRedials dials, never an
// unbounded loop.
func TestReconnect_BudgetBounded(t *testing.T) {
	dials := 0
	tp := testTransport(&scriptDialer{})
	tp.dial = func(context.Context, string) (rawConn, error) {
		dials++
		return nil, errors.New("connect unreachable")
	}
	d := newDriver(tp)

	err := d.Tap(Coords{X: 1, Y: 2})
	e := asError(err)
	if e == nil || e.Code != CodeConnection || !e.Retryable {
		t.Fatalf("Tap = %v, want a retryable connection error after the budget", err)
	}
	if dials != _maxRedials+1 {
		t.Fatalf("dials = %d, want %d (initial + %d redials)", dials, _maxRedials+1, _maxRedials)
	}
}

// Mutating input carries a transport-minted idempotency key; reads stay
// keyless by contract (they are naturally safe and keyless reads keep the
// executor's ledger small).
func TestKeysOnlyOnMutatingMethods(t *testing.T) {
	conn := &matrixConn{respond: okAll()}
	d := newDriver(testTransport(&scriptDialer{script: conns(conn)}))

	if err := d.Tap(Coords{X: 1, Y: 2}); err != nil {
		t.Fatalf("Tap: %v", err)
	}
	if _, err := d.Observe(); err != nil {
		t.Fatalf("Observe: %v", err)
	}
	if key := keyOf(t, conn.sent[0]); key == "" {
		t.Fatal("Tap params carry no idempotencyKey")
	}
	if key := keyOf(t, conn.sent[1]); key != "" {
		t.Fatalf("Observe params carry idempotencyKey %q; reads are keyless by contract", key)
	}
}
