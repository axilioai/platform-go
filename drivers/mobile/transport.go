package mobile

import (
	"context"
	"encoding/json"
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
// comes back. One WebSocket per allocation; it reconnects lazily on the next
// call after a drop (the allocation lease outlives the socket).
type RemoteTransport struct {
	url         string
	openTimeout time.Duration
	// dial is injectable for tests; production opens a real WebSocket.
	dial func(ctx context.Context, url string) (rawConn, error)

	mu     sync.Mutex
	conn   rawConn
	nextID int64
}

func newRemoteTransport(url string, openTimeout time.Duration) *RemoteTransport {
	return &RemoteTransport{url: url, openTimeout: openTimeout, dial: dialWS}
}

func (t *RemoteTransport) call(ctx context.Context, method string, params any) (json.RawMessage, error) {
	t.mu.Lock()
	defer t.mu.Unlock()

	conn, err := t.ensureConnected()
	if err != nil {
		return nil, err
	}
	t.nextID++
	id := t.nextID
	frame, err := buildFrame(id, method, params)
	if err != nil {
		return nil, err
	}
	if err := conn.send(ctx, frame); err != nil {
		t.closeLocked()
		return nil, sendRecvErr(ctx, err)
	}
	return t.awaitReply(ctx, conn, id)
}

func (t *RemoteTransport) close() error {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.closeLocked()
}

// awaitReply reads until the frame that echoes id. Notifications (events, no id)
// and stale frames are skipped: the DCP request/response path carries no
// notifications and the telemetry up-channel routes them separately.
func (t *RemoteTransport) awaitReply(ctx context.Context, conn rawConn, id int64) (json.RawMessage, error) {
	for {
		data, err := conn.recv(ctx)
		if err != nil {
			t.closeLocked()
			return nil, sendRecvErr(ctx, err)
		}
		resp, err := decodeResponse(data)
		if err != nil {
			t.closeLocked()
			return nil, err
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

func (t *RemoteTransport) ensureConnected() (rawConn, error) {
	if t.conn != nil {
		return t.conn, nil
	}
	dctx, cancel := context.WithTimeout(context.Background(), t.openTimeout)
	defer cancel()
	conn, err := t.dial(dctx, t.url)
	if err != nil {
		return nil, &Error{Code: CodeConnection, Message: "cannot connect to control websocket: " + err.Error()}
	}
	t.conn = conn
	return conn, nil
}

func (t *RemoteTransport) closeLocked() error {
	if t.conn == nil {
		return nil
	}
	err := t.conn.closeConn()
	t.conn = nil
	return err
}

// sendRecvErr classifies an I/O error: a per-call deadline becomes a timeout,
// everything else a connection error. Both are surfaced after the socket has
// been dropped so a late reply can't be misread as the next call's.
func sendRecvErr(ctx context.Context, err error) error {
	if ctx.Err() == context.DeadlineExceeded {
		return &Error{Code: CodeTimeout, Message: "call timed out", Retryable: true}
	}
	return &Error{Code: CodeConnection, Message: "control websocket I/O failed: " + err.Error()}
}

// --- real WebSocket ---------------------------------------------------------

func dialWS(ctx context.Context, url string) (rawConn, error) {
	c, _, err := websocket.Dial(ctx, url, nil)
	if err != nil {
		return nil, err
	}
	c.SetReadLimit(_readLimit)
	return &wsRawConn{c: c}, nil
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
