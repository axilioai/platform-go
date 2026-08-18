package mobile

// Protocol.handshake / Device.info — the DCP capability-negotiation surface
// (AXI-1753). Every executor implements the handshake, so it is not treated as
// skew-tolerant: a failure (including UnknownOp) propagates as an error rather
// than being swallowed. Call it once on connect to gate behavior on the
// advertised capabilities.

// DeviceInfo is the static device descriptor (the Device.info result and the
// handshake's `device`). Fields mirror the DCP DeviceInfo schema; a client sizes
// behavior from InputModalities / FormFactor, never from a hardcoded platform.
type DeviceInfo struct {
	DeviceID        string   `json:"device_id"`
	Platform        string   `json:"platform"`
	FormFactor      string   `json:"form_factor"`
	InputModalities []string `json:"input_modalities"`
	Model           string   `json:"model,omitempty"`
	OSVersion       string   `json:"os_version,omitempty"`
	ScreenWidth     int      `json:"screen_width"`
	ScreenHeight    int      `json:"screen_height"`
}

// HandshakeResult is the Protocol.handshake reply. Capabilities is the exact
// method-name set the executor answers; Domains is the capability-profile set it
// advertises (the domain prefixes).
type HandshakeResult struct {
	ProtocolVersion int        `json:"protocol_version"`
	Device          DeviceInfo `json:"device"`
	Domains         []string   `json:"domains"`
	Capabilities    []string   `json:"capabilities"`
}

// Supports reports whether the executor advertised the given DCP method name
// (e.g. "Touch.tap").
func (h *HandshakeResult) Supports(method string) bool {
	for _, m := range h.Capabilities {
		if m == method {
			return true
		}
	}
	return false
}

// HasDomain reports whether the executor advertised the given capability domain
// (e.g. "Touch").
func (h *HandshakeResult) HasDomain(domain string) bool {
	for _, d := range h.Domains {
		if d == domain {
			return true
		}
	}
	return false
}

// Handshake performs the DCP Protocol.handshake and returns the executor's
// advertised protocol version, device descriptor, domains, and capabilities.
// Call it once on connect to gate behavior on the returned capabilities instead
// of assuming the phone surface. Any error (including UnknownOp from a server
// that somehow lacks the method) is returned as-is.
func (d *MobileDriver) Handshake(opts ...CallOption) (*HandshakeResult, error) {
	cfg := applyCall(defaultCallTimeout, opts)
	raw, err := d.call(methodProtocolHandshake, handshakeParams{}, cfg.timeout)
	if err != nil {
		return nil, err
	}
	out := &HandshakeResult{}
	if err := unmarshalResult(raw, out); err != nil {
		return nil, err
	}
	return out, nil
}

// DeviceInfo performs Device.info and returns the static device descriptor.
func (d *MobileDriver) DeviceInfo(opts ...CallOption) (*DeviceInfo, error) {
	cfg := applyCall(defaultCallTimeout, opts)
	raw, err := d.call(methodDeviceInfo, struct{}{}, cfg.timeout)
	if err != nil {
		return nil, err
	}
	out := &DeviceInfo{}
	if err := unmarshalResult(raw, out); err != nil {
		return nil, err
	}
	return out, nil
}
