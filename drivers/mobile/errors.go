package mobile

import "errors"

// Code is a stable, machine-readable error classification. Each DCP error kind
// maps 1:1 onto one of these (see fromDCPError); the driver also raises a few
// locally (CodeTimeout from wait loops, CodeElementNotFound from find,
// CodeConnection from a failed dial).
type Code string

const (
	// CodeUnknownOp — the executor doesn't recognise the method (version skew).
	CodeUnknownOp Code = "unknown_op"
	// CodeInvalidArgs — params failed validation.
	CodeInvalidArgs Code = "invalid_args"
	// CodeNoAllocation — the daemon has no active allocation.
	CodeNoAllocation Code = "no_allocation"
	// CodeNotConnected — the executor couldn't reach the on-device agent.
	CodeNotConnected Code = "not_connected"
	// CodeDeviceOffline — the device is transiently unavailable (retryable).
	CodeDeviceOffline Code = "device_offline"
	// CodeUnauthorized — the session token was rejected.
	CodeUnauthorized Code = "unauthorized"
	// CodeInternal — an unclassified executor-side failure.
	CodeInternal Code = "internal"
	// CodeCanceled — the operation was canceled (deadline or context).
	CodeCanceled Code = "canceled"
	// CodeConnection means the control WebSocket failed and could not be
	// re-established within the transport's bounded reconnect budget.
	// Retryable: the allocation may still be live, so a later call (or a
	// caller-level retry) can succeed.
	CodeConnection Code = "connection"
	// CodeSessionEnded means the session is over: the server closed with 1000
	// ("session ended", or this connection was superseded by a newer one)
	// or refused the reattach with HTTP 403 (allocation no longer active).
	// Terminal; never retried.
	CodeSessionEnded Code = "session_ended"
	// CodeControlHeld means another controller holds the session's control
	// lease (close code 4409). Terminal for this transport; surfaced, never
	// auto-retried (a retry loop against a held lease is a thundering herd
	// aimed at the one-controller guardrail).
	CodeControlHeld Code = "control_held"
	// CodeElementNotFound — a selector found nothing.
	CodeElementNotFound Code = "element_not_found"
	// CodeTimeout — a call or a wait loop exceeded its deadline (retryable).
	CodeTimeout Code = "timeout"
)

// Error is the single error type this package returns. Classify it with the
// helpers (IsTimeout, IsElementNotFound, ...) or by comparing .Code.
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

// asError unwraps err to a *Error, or nil.
func asError(err error) *Error {
	var e *Error
	if errors.As(err, &e) {
		return e
	}
	return nil
}

func hasCode(err error, code Code) bool {
	e := asError(err)
	return e != nil && e.Code == code
}

// IsTimeout reports whether err is a timeout (a call or wait loop deadline).
func IsTimeout(err error) bool { return hasCode(err, CodeTimeout) }

// IsElementNotFound reports whether err is a find that matched nothing.
func IsElementNotFound(err error) bool { return hasCode(err, CodeElementNotFound) }

// IsDeviceOffline reports whether err is a transient device-offline (retryable).
func IsDeviceOffline(err error) bool { return hasCode(err, CodeDeviceOffline) }

// IsRetryable reports whether err carries the retryable flag.
func IsRetryable(err error) bool {
	e := asError(err)
	return e != nil && e.Retryable
}

// _kindToCode maps a DCP error frame's data.kind (PascalCase) onto a Code.
var _kindToCode = map[string]Code{
	kindUnknownOp:       CodeUnknownOp,
	kindInvalidArgs:     CodeInvalidArgs,
	kindNoAllocation:    CodeNoAllocation,
	kindNotConnected:    CodeNotConnected,
	kindDeviceOffline:   CodeDeviceOffline,
	kindElementNotFound: CodeElementNotFound,
	kindTimeout:         CodeTimeout,
	kindUnauthorized:    CodeUnauthorized,
	kindInternal:        CodeInternal,
	kindCanceled:        CodeCanceled,
}

// fromDCPError maps a DCP error frame's error object onto an *Error. An unknown
// kind degrades to CodeInternal.
func fromDCPError(de *dcpError) *Error {
	kind := kindInternal
	retryable := false
	if de.Data != nil {
		if de.Data.Kind != "" {
			kind = de.Data.Kind
		}
		retryable = de.Data.Retryable
	}
	code, ok := _kindToCode[kind]
	if !ok {
		code = CodeInternal
	}
	return &Error{Code: code, Message: de.Message, Retryable: retryable}
}
