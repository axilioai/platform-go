package mobile

import "encoding/json"

// DCP (Device Control Protocol) wire constants + codec. DCP is literal CDP: a
// command rides the wire as {"id","method","params"}, the reply echoes the id
// with exactly one of "result"/"error". This mirrors
// commons/go/types/backend/realtime/dcp.go in the monorepo — redeclared here so
// the public SDK stays free of the internal commons module (the same way
// platform-python's drivers/mobile/_envelope.py mirrors it).

// DCP method names ("Domain.method"). The driver's helpers translate their
// ergonomic API to these verbatim, the way Playwright translates to CDP.
const (
	methodInputTap         = "Input.tap"
	methodInputLongPress   = "Input.longPress"
	methodInputSwipe       = "Input.swipe"
	methodInputTypeText    = "Input.typeText"
	methodInputKeyPress    = "Input.keyPress"
	methodScreenScreenshot = "Screen.screenshot"
	methodScreenObserve    = "Screen.observe"
	methodScreenFind       = "Screen.find"
)

// DCP error kinds (the data.kind on a CDP error frame). PascalCase to mirror the
// Go side; mapped to Code in fromDCPError.
const (
	kindUnknownOp       = "UnknownOp"
	kindInvalidArgs     = "InvalidArgs"
	kindNoAllocation    = "NoAllocation"
	kindNotConnected    = "NotConnected"
	kindDeviceOffline   = "DeviceOffline"
	kindElementNotFound = "ElementNotFound"
	kindTimeout         = "Timeout"
	kindUnauthorized    = "Unauthorized"
	kindInternal        = "Internal"
	kindCanceled        = "Canceled"
)

// dcpCommand is a client->device request frame.
type dcpCommand struct {
	ID     int64           `json:"id"`
	Method string          `json:"method"`
	Params json.RawMessage `json:"params,omitempty"`
}

// dcpResponse is a device->client reply frame (exactly one of Result/Error).
type dcpResponse struct {
	ID     int64           `json:"id"`
	Result json.RawMessage `json:"result,omitempty"`
	Error  *dcpError       `json:"error,omitempty"`
}

// dcpError is the error object on a reply frame.
type dcpError struct {
	Code    int           `json:"code"`
	Message string        `json:"message"`
	Data    *dcpErrorData `json:"data,omitempty"`
}

type dcpErrorData struct {
	Kind      string `json:"kind"`
	Retryable bool   `json:"retryable"`
}

// buildFrame marshals a command frame. params may be nil (no params).
func buildFrame(id int64, method string, params any) ([]byte, error) {
	cmd := dcpCommand{ID: id, Method: method}
	if params != nil {
		raw, err := json.Marshal(params)
		if err != nil {
			return nil, &Error{Code: CodeInternal, Message: "marshal params: " + err.Error()}
		}
		cmd.Params = raw
	}
	return json.Marshal(cmd)
}

// decodeResponse parses a reply frame.
func decodeResponse(data []byte) (*dcpResponse, error) {
	var resp dcpResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, &Error{Code: CodeInternal, Message: "malformed JSON frame: " + err.Error()}
	}
	return &resp, nil
}
