package mobile

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"
)

// fakeConn is a scripted rawConn: each sent command is decoded, handed to a
// responder, and the reply queued for the next recv (mirrors the Python
// FakeWS test seam).
type fakeConn struct {
	responder func(dcpCommand) dcpResponse
	sent      []dcpCommand
	inbox     [][]byte
}

func (f *fakeConn) send(_ context.Context, data []byte) error {
	var cmd dcpCommand
	if err := json.Unmarshal(data, &cmd); err != nil {
		return err
	}
	f.sent = append(f.sent, cmd)
	b, _ := json.Marshal(f.responder(cmd))
	f.inbox = append(f.inbox, b)
	return nil
}

func (f *fakeConn) recv(_ context.Context) ([]byte, error) {
	if len(f.inbox) == 0 {
		return nil, errors.New("no more frames")
	}
	b := f.inbox[0]
	f.inbox = f.inbox[1:]
	return b, nil
}

func (f *fakeConn) closeConn() error { return nil }

func driverWith(fc *fakeConn) *MobileDriver {
	rt := &RemoteTransport{
		url:         "wss://connect.test/api/v1/realtime/ws/control?token=abc",
		openTimeout: time.Second,
		dial:        func(context.Context, string) (rawConn, error) { return fc, nil },
	}
	return newDriver(rt)
}

func okResp(cmd dcpCommand, result any) dcpResponse {
	raw, _ := json.Marshal(result)
	return dcpResponse{ID: cmd.ID, Result: raw}
}

func TestTapSendsInputTap(t *testing.T) {
	fc := &fakeConn{responder: func(cmd dcpCommand) dcpResponse { return dcpResponse{ID: cmd.ID} }}
	d := driverWith(fc)

	if err := d.Tap(Coords{X: 5, Y: 6}); err != nil {
		t.Fatalf("Tap: %v", err)
	}
	if len(fc.sent) != 1 || fc.sent[0].Method != methodInputTap {
		t.Fatalf("want one Input.tap, got %+v", fc.sent)
	}
	var p tapParams
	if err := json.Unmarshal(fc.sent[0].Params, &p); err != nil {
		t.Fatalf("params: %v", err)
	}
	if p.X != 5 || p.Y != 6 {
		t.Fatalf("want {5,6}, got %+v", p)
	}
	if fc.sent[0].ID != 1 {
		t.Fatalf("want id 1, got %d", fc.sent[0].ID)
	}
}

func TestObserveParsesScreen(t *testing.T) {
	fc := &fakeConn{responder: func(cmd dcpCommand) dcpResponse {
		return okResp(cmd, map[string]any{
			"texts": []map[string]any{
				{"text": "Search", "confidence": 0.9, "bbox": map[string]int{"x": 10, "y": 20, "width": 100, "height": 40}},
			},
			"icons":       []map[string]any{{"confidence": 0.7, "bbox": map[string]int{"x": 0, "y": 0, "width": 24, "height": 24}}},
			"hash":        "abc",
			"width":       1080,
			"height":      2400,
			"captured_at": 1_700_000_000_000,
		})
	}}
	d := driverWith(fc)

	screen, err := d.Observe()
	if err != nil {
		t.Fatalf("Observe: %v", err)
	}
	if fc.sent[0].Method != methodScreenObserve {
		t.Fatalf("want Screen.observe, got %q", fc.sent[0].Method)
	}
	if len(screen.Texts) != 1 || screen.Texts[0].Text != "Search" {
		t.Fatalf("bad texts: %+v", screen.Texts)
	}
	// center = (x + w/2, y + h/2) = (60, 40)
	if got := screen.Texts[0].Center; got.X != 60 || got.Y != 40 {
		t.Fatalf("bad center: %+v", got)
	}
	if len(screen.Icons) != 1 || screen.Width != 1080 || screen.Hash != "abc" {
		t.Fatalf("bad screen: %+v", screen)
	}
	if screen.CapturedAt.UnixMilli() != 1_700_000_000_000 {
		t.Fatalf("bad captured_at: %v", screen.CapturedAt)
	}

	// default OCR engine flows through as "free"
	var op observeParams
	_ = json.Unmarshal(fc.sent[0].Params, &op)
	if op.OcrEngine != "free" {
		t.Fatalf("want free engine, got %q", op.OcrEngine)
	}
}

func TestFindFoundAndNotFound(t *testing.T) {
	// found
	fc := &fakeConn{responder: func(cmd dcpCommand) dcpResponse {
		return okResp(cmd, map[string]any{
			"found": map[string]any{"text": "Continue", "confidence": 0.8, "bbox": map[string]int{"x": 0, "y": 0, "width": 200, "height": 80}},
		})
	}}
	d := driverWith(fc)
	el, err := d.Find("the continue button", WithOCREngine("premium"), WithModel("openai/gpt-5"))
	if err != nil {
		t.Fatalf("Find: %v", err)
	}
	if el.Source != SourceOCR || el.Center.X != 100 || el.Center.Y != 40 {
		t.Fatalf("bad element: %+v", el)
	}
	var fp findParams
	_ = json.Unmarshal(fc.sent[0].Params, &fp)
	if fp.Query != "the continue button" || fp.OcrEngine != "premium" || fp.Model != "openai/gpt-5" {
		t.Fatalf("bad find params: %+v", fp)
	}

	// not found -> CodeElementNotFound
	fc2 := &fakeConn{responder: func(cmd dcpCommand) dcpResponse {
		return okResp(cmd, map[string]any{"found": nil})
	}}
	d2 := driverWith(fc2)
	if _, err := d2.Find("nope"); !IsElementNotFound(err) {
		t.Fatalf("want element-not-found, got %v", err)
	}
}

func TestErrorFrameMapsToCode(t *testing.T) {
	fc := &fakeConn{responder: func(cmd dcpCommand) dcpResponse {
		return dcpResponse{ID: cmd.ID, Error: &dcpError{
			Code: -32004, Message: "device is rebooting",
			Data: &dcpErrorData{Kind: kindDeviceOffline, Retryable: true},
		}}
	}}
	d := driverWith(fc)

	err := d.Tap(Coords{X: 1, Y: 2})
	if !IsDeviceOffline(err) || !IsRetryable(err) {
		t.Fatalf("want retryable device_offline, got %v", err)
	}
}

func TestElementChainTaps(t *testing.T) {
	fc := &fakeConn{responder: func(cmd dcpCommand) dcpResponse {
		if cmd.Method == methodScreenObserve {
			return okResp(cmd, map[string]any{
				"texts": []map[string]any{
					{"text": "Settings", "confidence": 0.9, "bbox": map[string]int{"x": 0, "y": 0, "width": 100, "height": 20}},
				},
			})
		}
		return dcpResponse{ID: cmd.ID}
	}}
	d := driverWith(fc)

	el, err := d.FindText("settings", false)
	if err != nil || el == nil {
		t.Fatalf("FindText: %v el=%v", err, el)
	}
	if err := el.Tap(); err != nil {
		t.Fatalf("chained Tap: %v", err)
	}
	last := fc.sent[len(fc.sent)-1]
	if last.Method != methodInputTap {
		t.Fatalf("want chained Input.tap, got %q", last.Method)
	}
	var p tapParams
	_ = json.Unmarshal(last.Params, &p)
	if p.X != 50 || p.Y != 10 { // center of {0,0,100,20}
		t.Fatalf("want tap at center {50,10}, got %+v", p)
	}
}

func TestScreenshotDecodesBytes(t *testing.T) {
	fc := &fakeConn{responder: func(cmd dcpCommand) dcpResponse {
		// Go marshals []byte as base64, matching the wire.
		return okResp(cmd, map[string]any{"png_base64": []byte("PNGDATA")})
	}}
	d := driverWith(fc)
	png, err := d.Screenshot()
	if err != nil {
		t.Fatalf("Screenshot: %v", err)
	}
	if string(png) != "PNGDATA" {
		t.Fatalf("bad png bytes: %q", png)
	}
}
