package mobile

import (
	"context"
	"encoding/json"
	"math/rand/v2"
	"net/url"
	"strings"
	"time"

	"github.com/coder/websocket"
	"github.com/google/uuid"
)

// Reconnect machinery for the DCP control transport: the client half of the
// raw-WebSocket reconnect contract. The server persists every frame to a
// per-session stream and interleaves opt-in cursor checkpoints; this file
// owns presenting the cursor back, pacing the redial, and making the
// re-send of the interrupted command idempotent.

// Transport notification methods — the vendor Axilio domain the server uses
// for resume machinery below the CDP frame. Strictly opt-in server-side
// (only sent because every attach carries resume=1), and structurally
// id-less so they can never be mistaken for a response.
const (
	methodAxilioCursor         = "Axilio.cursor"
	methodAxilioResyncRequired = "Axilio.resyncRequired"
)

// _closeControlHeld is the application close code for "another controller
// holds the session's control lease" (4409). Terminal by contract: a retry
// loop against a held lease is the one-controller model's failure mode.
const _closeControlHeld websocket.StatusCode = 4409

// Redial budget: bounded exponential backoff with full jitter. Base/cap are
// sized for an interactive SDK call — the worst case (~15s of accumulated
// backoff across the budget) stays inside a default call deadline while
// covering a connect pod replacement.
const (
	_maxRedials = 6
	_redialBase = 250 * time.Millisecond
	_redialCap  = 8 * time.Second
)

// redialDelay returns the wait before redial attempt (0-based): uniformly
// random in (0, min(cap, base·2^attempt)]. Full jitter so a fleet of
// clients that lost the same pod doesn't reconnect in lockstep.
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

// attachURL is the control URL plus the resume params: every attach opts in
// to checkpoints (resume=1), and a reattach that holds a cursor presents it
// so delivery continues where this client left off. The same URL stays
// valid across redials by design — the control token outlives the session
// cap.
func (t *RemoteTransport) attachURL() string {
	u, err := url.Parse(t.url)
	if err != nil {
		// An unparsable control URL fails at dial with a real error; don't
		// mask it here.
		return t.url
	}
	q := u.Query()
	q.Set("resume", "1")
	if t.cursor != "" {
		q.Set("cursor", t.cursor)
	}
	u.RawQuery = q.Encode()
	return u.String()
}

// isMutatingMethod reports whether a DCP method mutates device state — the
// interaction domains (Touch.* / Keyboard.*). Only these carry idempotency
// keys: reads are naturally safe to re-send, and keyless reads are what
// keep the executor's ledger small.
func isMutatingMethod(method string) bool {
	domain, _, ok := strings.Cut(method, ".")
	if !ok {
		return false
	}
	return domain == "Touch" || domain == "Keyboard"
}

// withIdempotencyKey stamps a fresh client-generated key into a mutating
// command's marshaled params. Minted once per logical command (the caller
// re-sends the same bytes), so the executor's per-allocation ledger can
// answer a post-reconnect re-send with the recorded response instead of
// executing twice.
func withIdempotencyKey(params json.RawMessage) (json.RawMessage, error) {
	fields := map[string]json.RawMessage{}
	if len(params) > 0 {
		if err := json.Unmarshal(params, &fields); err != nil {
			return nil, &Error{Code: CodeInternal, Message: "params not an object: " + err.Error()}
		}
	}
	key, err := json.Marshal(uuid.NewString())
	if err != nil {
		return nil, &Error{Code: CodeInternal, Message: "marshal idempotency key: " + err.Error()}
	}
	fields["idempotencyKey"] = key
	out, err := json.Marshal(fields)
	if err != nil {
		return nil, &Error{Code: CodeInternal, Message: "marshal keyed params: " + err.Error()}
	}
	return out, nil
}
