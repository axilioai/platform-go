package mobile

import "encoding/json"

// DCP (Device Control Protocol) frame codec. DCP is literal CDP: a command
// rides the wire as {"id","method","params"}, the reply echoes the id with
// exactly one of "result"/"error". The method-name constants, error-kind
// constants, and input-param frames are generated from the vendored contract
// into wire_gen.go (see scripts/gen_dcp_wire.go); this file is the hand-written
// transport codec that stays put across regens.

// dcpCommand is a client->device request frame.
type dcpCommand struct {
	ID     int64           `json:"id"`
	Method string          `json:"method"`
	Params json.RawMessage `json:"params,omitempty"`
}

// dcpResponse is a device->client frame. A reply echoes the command id with
// exactly one of Result/Error; an id-less frame is a notification and
// carries Method/Params instead (the Axilio.* transport notifications are
// the only ones emitted today).
type dcpResponse struct {
	ID     int64           `json:"id"`
	Result json.RawMessage `json:"result,omitempty"`
	Error  *dcpError       `json:"error,omitempty"`
	Method string          `json:"method,omitempty"`
	Params json.RawMessage `json:"params,omitempty"`
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

// marshalParams marshals a call's params object; nil means no params.
func marshalParams(params any) (json.RawMessage, error) {
	if params == nil {
		return nil, nil
	}
	raw, err := json.Marshal(params)
	if err != nil {
		return nil, &Error{Code: CodeInternal, Message: "marshal params: " + err.Error()}
	}
	return raw, nil
}

// buildRawFrame marshals a command frame from already-marshaled params,
// the shape the transport works in, so a re-send after a reconnect reuses
// the exact same params bytes (idempotency key included).
func buildRawFrame(id int64, method string, params json.RawMessage) ([]byte, error) {
	return json.Marshal(dcpCommand{ID: id, Method: method, Params: params})
}

// decodeResponse parses a reply frame.
func decodeResponse(data []byte) (*dcpResponse, error) {
	var resp dcpResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, &Error{Code: CodeInternal, Message: "malformed JSON frame: " + err.Error()}
	}
	return &resp, nil
}
