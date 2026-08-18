package mobile

import "testing"

func TestHandshakeParsesResult(t *testing.T) {
	fc := &fakeConn{responder: func(cmd dcpCommand) dcpResponse {
		if cmd.Method != methodProtocolHandshake {
			t.Errorf("handshake sent method %q, want %q", cmd.Method, methodProtocolHandshake)
		}
		return okResp(cmd, map[string]any{
			"protocol_version": 1,
			"device": map[string]any{
				"device_id": "p1", "platform": "android", "form_factor": "phone",
				"input_modalities": []string{"touch", "keyboard"},
				"screen_width": 1080, "screen_height": 2400,
			},
			"domains":      []string{"Device", "Keyboard", "Protocol", "Screen", "Touch"},
			"capabilities": []string{"Touch.tap", "Screen.observe"},
		})
	}}
	d := driverWith(fc)

	hs, err := d.Handshake()
	if err != nil {
		t.Fatalf("handshake: %v", err)
	}
	if hs.ProtocolVersion != 1 {
		t.Fatalf("protocol_version = %d, want 1", hs.ProtocolVersion)
	}
	if hs.Device.Platform != "android" || hs.Device.FormFactor != "phone" {
		t.Fatalf("device descriptor not parsed: %+v", hs.Device)
	}
	if !hs.Supports("Touch.tap") || hs.Supports("Pointer.click") {
		t.Fatalf("Supports() wrong: caps=%v", hs.Capabilities)
	}
	if !hs.HasDomain("Touch") || hs.HasDomain("Pointer") {
		t.Fatalf("HasDomain() wrong: domains=%v", hs.Domains)
	}
}

func TestHandshakePropagatesError(t *testing.T) {
	// Every executor implements the handshake, so it is not skew-tolerant: an
	// error (including UnknownOp) propagates rather than being swallowed (AXI-1753).
	fc := &fakeConn{responder: func(cmd dcpCommand) dcpResponse {
		return dcpResponse{ID: cmd.ID, Error: &dcpError{
			Code: -32601, Message: "unknown method",
			Data: &dcpErrorData{Kind: kindUnknownOp},
		}}
	}}
	d := driverWith(fc)

	if _, err := d.Handshake(); err == nil {
		t.Fatal("handshake must return the error, not swallow it")
	} else if !hasCode(err, CodeUnknownOp) {
		t.Fatalf("want unknown_op error, got %v", err)
	}
}

func TestDeviceInfoParses(t *testing.T) {
	fc := &fakeConn{responder: func(cmd dcpCommand) dcpResponse {
		if cmd.Method != methodDeviceInfo {
			t.Errorf("sent method %q, want %q", cmd.Method, methodDeviceInfo)
		}
		return okResp(cmd, map[string]any{
			"device_id": "p1", "platform": "android", "form_factor": "phone",
			"input_modalities": []string{"touch", "keyboard"},
			"model": "Pixel 8", "os_version": "16",
			"screen_width": 1080, "screen_height": 2400,
		})
	}}
	d := driverWith(fc)

	info, err := d.DeviceInfo()
	if err != nil {
		t.Fatalf("device info: %v", err)
	}
	if info.Model != "Pixel 8" || len(info.InputModalities) != 2 {
		t.Fatalf("descriptor not parsed: %+v", info)
	}
}
