package telemetry

import "errors"

// Code is a stable, machine-readable error classification for the telemetry
// stream. The taxonomy deliberately parallels drivers/mobile (the control
// leg): the two transports share a reconnect contract but not a package, so
// each can evolve with its own wire.
//
// The end of a session is not an error: the stream returns io.EOF and Ended
// reports true.
type Code string

const (
	// CodeUnauthorized: the credential was rejected (HTTP 401 or 403). The
	// token in a telemetry_url expires; fetch the session again for a fresh
	// one.
	CodeUnauthorized Code = "unauthorized"
	// CodeNotFound: no such session in the caller's organization (HTTP 404).
	CodeNotFound Code = "not_found"
	// CodeTelemetryOff: the session was started with telemetry off, so
	// nothing is streamed or stored (HTTP 409).
	CodeTelemetryOff Code = "telemetry_off"
	// CodeConnection: the stream failed and could not be reopened within the
	// bounded budget. Retryable: a fresh Tail can succeed.
	CodeConnection Code = "connection"
	// CodeClosed: the stream was closed locally via Close while a read was
	// in flight or before it started.
	CodeClosed Code = "closed"
)

// Error is the single error type this package returns for telemetry-stream
// failures. Context cancellation is surfaced as the context's own error, and
// a cleanly ended stream as io.EOF; neither is wrapped in an Error.
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

// IsUnauthorized reports whether err says the credential was rejected.
func IsUnauthorized(err error) bool { return hasCode(err, CodeUnauthorized) }

// IsNotFound reports whether err says the session does not exist for the caller.
func IsNotFound(err error) bool { return hasCode(err, CodeNotFound) }

// IsTelemetryOff reports whether err says the session has telemetry off.
func IsTelemetryOff(err error) bool { return hasCode(err, CodeTelemetryOff) }

// IsConnection reports whether err is an exhausted-reconnect connection
// failure (retryable with a fresh Tail).
func IsConnection(err error) bool { return hasCode(err, CodeConnection) }
