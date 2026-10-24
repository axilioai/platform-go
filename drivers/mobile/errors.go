package mobile

import "errors"

// Code is a stable, machine-readable error classification. Each DCP error kind
// maps 1:1 onto one of these (see fromDCPError); the driver also raises a few
// locally (CodeTimeout from a call deadline, CodeConnection from a failed dial).
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
	// CodeTimeout — a call or a wait loop exceeded its deadline (retryable).
	CodeTimeout Code = "timeout"
	// CodeActionTimeout: a locator action's device-side auto-wait exceeded
	// its timeoutMs budget without the target becoming actionable. Not
	// retryable: an inference already running when the budget ended was
	// allowed to finish, so a bare retry races the same clock again.
	CodeActionTimeout Code = "action_timeout"
	// CodeStrategyUnavailable: the call needs a resolver the session doesn't
	// have. A locator with a tree-only selector (role, name, id, states,
	// value, windowId, nodeId, platform) or StrategyAccessibility, or an
	// Accessibility tree read, on a session whose accessibility tree is
	// off. Not retryable as is: turn the tree on (Accessibility().Enable,
	// where State().Toggleable) or change the selector.
	CodeStrategyUnavailable Code = "strategy_unavailable"
	// CodeTreeUnavailable: the accessibility tree is on but has no app
	// window to read, typically because a system dialog such as a runtime
	// permission prompt is covering the app. Not retryable as is: dismiss
	// the dialog (vision still sees it) and try again.
	CodeTreeUnavailable Code = "tree_unavailable"
	// CodeStaleNode: a node id from an earlier snapshot (NodeID, or
	// Accessibility().Partial / Children) no longer exists. Never matched
	// against a different element. Take a fresh snapshot.
	CodeStaleNode Code = "stale_node"
)

// Error is the single error type this package returns. Classify it with the
// helpers (IsTimeout, IsActionTimeout, ...) or by comparing .Code.
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

// IsDeviceOffline reports whether err is a transient device-offline (retryable).
func IsDeviceOffline(err error) bool { return hasCode(err, CodeDeviceOffline) }

// IsActionTimeout reports whether err is a locator action whose device-side
// auto-wait exceeded its timeoutMs budget.
func IsActionTimeout(err error) bool { return hasCode(err, CodeActionTimeout) }

// IsStrategyUnavailable reports whether err is a call that needs a resolver
// the session doesn't have (e.g. role/id need the accessibility tree, and
// it is off for this session).
func IsStrategyUnavailable(err error) bool { return hasCode(err, CodeStrategyUnavailable) }

// IsTreeUnavailable reports whether err is the accessibility tree having no
// app window to read (a system dialog is covering the app).
func IsTreeUnavailable(err error) bool { return hasCode(err, CodeTreeUnavailable) }

// IsStaleNode reports whether err is a node id that no longer exists.
func IsStaleNode(err error) bool { return hasCode(err, CodeStaleNode) }

// IsRetryable reports whether err carries the retryable flag.
func IsRetryable(err error) bool {
	e := asError(err)
	return e != nil && e.Retryable
}

// _kindToCode maps a DCP error frame's data.kind (PascalCase) onto a Code.
var _kindToCode = map[string]Code{
	kindUnknownOp:           CodeUnknownOp,
	kindInvalidArgs:         CodeInvalidArgs,
	kindNoAllocation:        CodeNoAllocation,
	kindNotConnected:        CodeNotConnected,
	kindDeviceOffline:       CodeDeviceOffline,
	kindTimeout:             CodeTimeout,
	kindUnauthorized:        CodeUnauthorized,
	kindInternal:            CodeInternal,
	kindCanceled:            CodeCanceled,
	kindActionTimeout:       CodeActionTimeout,
	kindStrategyUnavailable: CodeStrategyUnavailable,
	kindTreeUnavailable:     CodeTreeUnavailable,
	kindStaleNode:           CodeStaleNode,
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
