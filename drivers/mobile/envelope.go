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
