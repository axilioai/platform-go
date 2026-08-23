package telemetry

import "errors"

// Code is a stable, machine-readable error classification for the telemetry
// leg. The taxonomy deliberately parallels drivers/mobile (the control leg):
// the two transports share a reconnect contract but not a package, so each
// leg's codes can evolve with its own wire.
type Code string

const (
	// CodeUnauthorized — the telemetry token was rejected (HTTP 401 on
	// attach). Terminal: the token allocate handed out is the only one there
	// is.
	CodeUnauthorized Code = "unauthorized"
	// CodeConnection — the telemetry WebSocket failed and could not be
	// re-established within the bounded reconnect budget. Retryable: the
	// allocation may still be live, so a fresh Tail can succeed.
	CodeConnection Code = "connection"
	// CodeSessionEnded — the session is over: the server refused the attach
	// with HTTP 403 (the allocation-scoped token dies at deallocation
	// regardless of its exp), or closed normally before the terminal
	// session end frame arrived. Terminal; never retried. The durable
	// archive (Session.Trace) still holds whatever was recorded.
	CodeSessionEnded Code = "session_ended"
	// CodeClosed — the stream was closed locally via Close while a read was
	// in flight or before it started.
	CodeClosed Code = "closed"
)

// Error is the single error type this package returns for telemetry-leg
// failures. Context cancellation is surfaced as the context's own error, and
// a cleanly ended stream as io.EOF — neither is wrapped in an Error.
type Error struct {
	Code      Code
	Message   string
	Retryable bool
}

func (e *Error) Error() string {
	if e.Message == "" {
		return string(e.Code)
	}
	return string(e.Code) + ": " + e.Message
}

func hasCode(err error, code Code) bool {
	var e *Error
	return errors.As(err, &e) && e.Code == code
}

// IsUnauthorized reports whether err says the telemetry token was rejected.
func IsUnauthorized(err error) bool { return hasCode(err, CodeUnauthorized) }

// IsConnection reports whether err is an exhausted-reconnect connection
// failure (retryable with a fresh Tail).
func IsConnection(err error) bool { return hasCode(err, CodeConnection) }

// IsSessionEnded reports whether err says the session is over.
func IsSessionEnded(err error) bool { return hasCode(err, CodeSessionEnded) }
