package mobile

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"sync"
	"time"

	"github.com/coder/websocket"
)

// transport is the seam every driver call goes through: one round-trip per call.
// method is a DCP method name ("Domain.method"); params marshals to the CDP
// params object. Returns the raw result payload (or a *Error).
type transport interface {
	call(ctx context.Context, method string, params any) (json.RawMessage, error)
	close() error
}

// rawConn is the minimal WebSocket surface the transport needs, so tests can
// inject a scripted connection without pulling in a real socket.
type rawConn interface {
	send(ctx context.Context, data []byte) error
	recv(ctx context.Context) ([]byte, error)
	closeConn() error
}

// _readLimit matches the hub's control-WS read limit (realtime_router.go): the
// default 32 KiB is far too small for base64 screenshot frames.
const _readLimit = 16 << 20

// RemoteTransport speaks the DCP control WebSocket — literal CDP frames. The
// driver already emits DCP method names, so this does no translation: each call
// goes out as {"id","method","params"} and the matching {"id","result"|"error"}
// comes back. One WebSocket per allocation; the allocation lease outlives the
// socket, so a drop is recovered rather than surfaced (the reconnect contract):
//
//   - Every attach opts in to cursor checkpoints (resume=1); the transport
//     tracks the latest Axilio.cursor and presents it on reattach, so the
//     server resumes delivery where this client left off.
//   - A retryable close (1001/1013/1011, abrupt TCP loss) triggers a bounded
//     backoff redial against the same URL (the control token deliberately
//     outlives the session cap). A terminal close (1000 session ended /
//     superseded, 4409 control held) or a 403/401 on reattach surfaces
//     immediately and is never auto-retried.
//   - The one in-flight command is re-sent after a successful reattach under
//     a fresh request id; mutating input commands carry a transport-generated
//     idempotencyKey (reused verbatim on the re-send), so the executor dedups
//     and the command executes exactly once.
//   - Request ids stay monotonic across redials — a resumed socket can still
//     deliver a pre-drop response, and a reused id would mismatch it to the
//     wrong call (the AXI-1293 wedge).
//   - If Handshake succeeded on a previous connection, it is replayed
//     internally after each reattach before work resumes: capability state is
//     per-connection.
type RemoteTransport struct {
	url         string
	openTimeout time.Duration
	// dial is injectable for tests; production opens a real WebSocket.
	dial func(ctx context.Context, url string) (rawConn, error)
	// redialDelay is the backoff schedule seam; tests shrink it to keep the
	// force-close matrix fast.
	redialDelay func(attempt int) time.Duration

	mu     sync.Mutex
	conn   rawConn
	nextID int64
	// cursor is the latest Axilio.cursor checkpoint — the opaque resume
	// token presented on reattach. Empty until the first checkpoint (or
	// after a resync, whose window the server could not replay).
	cursor string
	// handshakeParams holds the params of the last successful
	// Protocol.handshake, for the internal replay after a reattach. Nil
	// until the caller performs one.
	handshakeParams json.RawMessage
}

func newRemoteTransport(url string, openTimeout time.Duration) *RemoteTransport {
	return &RemoteTransport{
		url:         url,
		openTimeout: openTimeout,
		dial:        dialWS,
		redialDelay: redialDelay,
	}
}

func (t *RemoteTransport) call(ctx context.Context, method string, params any) (json.RawMessage, error) {
	t.mu.Lock()
	defer t.mu.Unlock()

	raw, err := marshalParams(params)
	if err != nil {
		return nil, err
	}
	// The idempotency key is minted once per logical command and reused
	// verbatim on every re-send of it, which is exactly what lets the
	// executor's ledger answer a duplicate without re-executing.
	if isMutatingMethod(method) {
		raw, err = withIdempotencyKey(raw)
		if err != nil {
			return nil, err
		}
	}

	for redials := 0; ; {
		result, err := t.attempt(ctx, method, raw)
		if err == nil {
			if method == methodProtocolHandshake {
				t.handshakeParams = raw
			}
			return result, nil
		}
		e := asError(err)
		if e == nil || e.Code != CodeConnection || !e.Retryable {
			return nil, err
		}
		redials++
		if redials > _maxRedials {
			return nil, err
		}
		if sleepCtx(ctx, t.redialDelay(redials-1)) != nil {
			return nil, &Error{Code: CodeTimeout, Message: "call timed out while reconnecting", Retryable: true}
		}
	}
}

func (t *RemoteTransport) close() error {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.closeLocked()
}

// attempt performs one connect-if-needed + send + await round of the call.
// A retryable connection failure comes back as CodeConnection/Retryable and
// the caller's redial loop decides; everything else is final for the call.
func (t *RemoteTransport) attempt(ctx context.Context, method string, params json.RawMessage) (json.RawMessage, error) {
	conn, err := t.ensureConnected(ctx)
	if err != nil {
		return nil, err
	}
	t.nextID++
	id := t.nextID
	frame, err := buildRawFrame(id, method, params)
	if err != nil {
		return nil, err
	}
	if err := conn.send(ctx, frame); err != nil {
		t.closeLocked()
		return nil, t.sendRecvErr(ctx, err)
	}
	return t.awaitReply(ctx, conn, id)
}

// awaitReply reads until the frame that echoes id. Id-less frames are
// notifications: the Axilio.* transport notifications are intercepted
// (cursor tracking, resync) before the skip, everything else is skipped —
// the DCP request/response path carries no other notifications today.
// Stale responses (a pre-drop reply redelivered after a resume) have older
// ids and are skipped by the same match.
func (t *RemoteTransport) awaitReply(ctx context.Context, conn rawConn, id int64) (json.RawMessage, error) {
	for {
		data, err := conn.recv(ctx)
		if err != nil {
			t.closeLocked()
			return nil, t.sendRecvErr(ctx, err)
		}
		resp, err := decodeResponse(data)
		if err != nil {
			t.closeLocked()
			return nil, err
		}
		if resp.ID == 0 {
			t.observeNotification(resp)
			continue
		}
		if resp.ID != id {
			continue
		}
		if resp.Error != nil {
			return nil, fromDCPError(resp.Error)
		}
		return resp.Result, nil
	}
}

// observeNotification handles one id-less transport frame.
func (t *RemoteTransport) observeNotification(resp *dcpResponse) {
	switch resp.Method {
	case methodAxilioCursor:
		var p struct {
			Cursor string `json:"cursor"`
		}
		if json.Unmarshal(resp.Params, &p) == nil && p.Cursor != "" {
			t.cursor = p.Cursor
		}
	case methodAxilioResyncRequired:
		// The retained window expired: the server could not replay the gap
		// and continued live. Nothing is lost on this request/response
		// path — the transport never relies on replayed responses (the
		// in-flight command is always re-sent, and the executor dedups) —
		// but the held cursor predates the window, so drop it rather than
		// re-present a known-stale token on the next reattach.
		t.cursor = ""
	}
}

func (t *RemoteTransport) ensureConnected(ctx context.Context) (rawConn, error) {
	if t.conn != nil {
		return t.conn, nil
	}
	dctx, cancel := context.WithTimeout(ctx, t.openTimeout)
	defer cancel()
	conn, err := t.dial(dctx, t.attachURL())
	if err != nil {
		return nil, classifyDialErr(err)
	}
	t.conn = conn
	// Capability state is per-connection: replay the caller's handshake
	// before any command resumes on the new socket.
	if t.handshakeParams != nil {
		if err := t.replayHandshake(ctx, conn); err != nil {
			t.closeLocked()
			return nil, err
		}
	}
	return conn, nil
}

// replayHandshake re-runs the remembered Protocol.handshake on a fresh
// connection. Failures surface like any connection/call failure: a dropped
// socket mid-replay is retryable, a DCP error is final for the call.
func (t *RemoteTransport) replayHandshake(ctx context.Context, conn rawConn) error {
	t.nextID++
	id := t.nextID
	frame, err := buildRawFrame(id, methodProtocolHandshake, t.handshakeParams)
	if err != nil {
		return err
	}
	if err := conn.send(ctx, frame); err != nil {
		t.closeLocked()
		return t.sendRecvErr(ctx, err)
	}
	_, err = t.awaitReply(ctx, conn, id)
	return err
}

func (t *RemoteTransport) closeLocked() error {
	if t.conn == nil {
		return nil
	}
	err := t.conn.closeConn()
	t.conn = nil
	return err
}

// sendRecvErr classifies an I/O error: a per-call deadline becomes a
// timeout, a terminal close (session over / control held) surfaces as its
// own non-retryable code, and everything else — 1001 going away, 1013 try
// again later, 1011 internal error, abrupt TCP loss — is a retryable
// connection failure the redial loop recovers. Always surfaced after the
// socket has been dropped so a late reply can't be misread as the next
// call's.
func (t *RemoteTransport) sendRecvErr(ctx context.Context, err error) error {
	if ctx.Err() == context.DeadlineExceeded {
		return &Error{Code: CodeTimeout, Message: "call timed out", Retryable: true}
	}
	switch websocket.CloseStatus(err) {
	case websocket.StatusNormalClosure:
		// 1000: "session ended", or this connection was superseded by a
		// newer one. Either way this transport is done.
		return &Error{Code: CodeSessionEnded, Message: "control websocket closed: " + err.Error()}
	case _closeControlHeld:
		return &Error{Code: CodeControlHeld, Message: "another controller holds the session's control lease: " + err.Error()}
	default:
		return &Error{Code: CodeConnection, Message: "control websocket I/O failed: " + err.Error(), Retryable: true}
	}
}

// --- real WebSocket ---------------------------------------------------------

// dialHTTPError carries the HTTP status of a refused WebSocket upgrade so
// the reattach classification can tell a dead session (403) from a bad
// token (401) from a transient failure.
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
// liveness signal: 403 means the allocation is no longer active (terminal),
// 401 means the token is bad (terminal); anything else is transient.
func classifyDialErr(err error) error {
	var de *dialHTTPError
	if errors.As(err, &de) {
		switch de.status {
		case http.StatusForbidden:
			return &Error{Code: CodeSessionEnded, Message: "session is no longer active: " + err.Error()}
		case http.StatusUnauthorized:
			return &Error{Code: CodeUnauthorized, Message: "control token rejected: " + err.Error()}
		}
	}
	return &Error{Code: CodeConnection, Message: "cannot connect to control websocket: " + err.Error(), Retryable: true}
}

type wsRawConn struct{ c *websocket.Conn }

func (w *wsRawConn) send(ctx context.Context, data []byte) error {
	return w.c.Write(ctx, websocket.MessageText, data)
}

func (w *wsRawConn) recv(ctx context.Context) ([]byte, error) {
	_, data, err := w.c.Read(ctx)
	return data, err
}

func (w *wsRawConn) closeConn() error {
	return w.c.Close(websocket.StatusNormalClosure, "")
}
