package mobile

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"
)

// driverWithOpts is driverWith plus construction-time Options, for the
// driver-level default tests below.
func driverWithOpts(fc *fakeConn, opts ...Option) *MobileDriver {
	rt := &RemoteTransport{
		url:         "wss://connect.test/api/v1/realtime/ws/control?token=abc",
		openTimeout: time.Second,
		dial:        func(context.Context, string) (rawConn, error) { return fc, nil },
	}
	return newDriver(rt, opts...)
}

func TestLocatorIsLazy(t *testing.T) {
	fc := &fakeConn{responder: func(cmd dcpCommand) dcpResponse { return dcpResponse{ID: cmd.ID} }}
	d := driverWith(fc)

	loc := d.GetByText("Continue").Within(d.GetByID("panel")).Nth(2).Filter("the enabled one")
	if len(fc.sent) != 0 {
		t.Fatalf("building a locator must send nothing, sent=%d", len(fc.sent))
	}
	_ = loc
}

func TestLocatorTapWireShapeAndDefaults(t *testing.T) {
	fc := &fakeConn{responder: func(cmd dcpCommand) dcpResponse {
		return okResp(cmd, map[string]any{
			"resolvedBy": "ocr",
			"bounds":     map[string]int{"x": 1, "y": 2, "width": 3, "height": 4},
			"tookMs":     42,
		})
	}}
	d := driverWith(fc)

	res, err := d.GetByText("Continue", Exact()).Tap()
	if err != nil {
		t.Fatalf("Tap: %v", err)
	}
	if res.ResolvedBy != "ocr" || res.TookMs != 42 || res.Bounds != (BBox{X: 1, Y: 2, Width: 3, Height: 4}) {
		t.Fatalf("bad result: %+v", res)
	}
	if fc.sent[0].Method != methodLocatorTap {
		t.Fatalf("wire method = %q, want %q", fc.sent[0].Method, methodLocatorTap)
	}

	var p locatorTapParams
	if err := json.Unmarshal(fc.sent[0].Params, &p); err != nil {
		t.Fatalf("params: %v", err)
	}
	if p.Locator == nil || p.Locator.Text != "Continue" || !p.Locator.Exact {
		t.Fatalf("bad locator on wire: %+v", p.Locator)
	}
	// Nothing set at call or driver level: strategy/model/ocrEngine are
	// omitted entirely (the server applies its own default) rather than a
	// client-side fallback forced onto the wire.
	if p.Strategy != "" || p.Model != "" || p.OcrEngine != "" {
		t.Fatalf("want strategy/model/ocrEngine omitted, got %+v", p)
	}
	// The device-side timeout budget is always sent explicitly, defaulting
	// to 5000ms.
	if p.TimeoutMs != 5000 {
		t.Fatalf("want default timeoutMs 5000, got %d", p.TimeoutMs)
	}

	// A raw JSON check that unset fields are actually absent, not just
	// zero-valued after decode.
	var raw map[string]any
	_ = json.Unmarshal(fc.sent[0].Params, &raw)
	for _, k := range []string{"strategy", "model", "ocrEngine"} {
		if _, present := raw[k]; present {
			t.Fatalf("wire params carry %q, want it omitted: %s", k, fc.sent[0].Params)
		}
	}
}

func TestLocatorOptionPrecedenceOverDriverDefault(t *testing.T) {
	fc := &fakeConn{responder: func(cmd dcpCommand) dcpResponse { return okResp(cmd, map[string]any{}) }}
	d := driverWithOpts(fc,
		WithDefaultStrategy(StrategyAccessibility), WithDefaultModel("driver-model"), WithDefaultOCREngine("driver-engine"))

	// Driver defaults flow through when the call sets nothing.
	if _, err := d.Locator(Text("x")).Tap(); err != nil {
		t.Fatalf("Tap: %v", err)
	}
	var p locatorTapParams
	_ = json.Unmarshal(fc.sent[0].Params, &p)
	if p.Strategy != StrategyAccessibility || p.Model != "driver-model" || p.OcrEngine != "driver-engine" {
		t.Fatalf("want driver defaults, got %+v", p)
	}

	// A per-call option wins over the driver default.
	if _, err := d.Locator(Text("x")).Tap(WithStrategy(StrategyVision), WithModel("call-model"), WithOCREngine("call-engine")); err != nil {
		t.Fatalf("Tap: %v", err)
	}
	_ = json.Unmarshal(fc.sent[1].Params, &p)
	if p.Strategy != StrategyVision || p.Model != "call-model" || p.OcrEngine != "call-engine" {
		t.Fatalf("want call overrides, got %+v", p)
	}
}

func TestLocatorNestedFields(t *testing.T) {
	fc := &fakeConn{responder: func(cmd dcpCommand) dcpResponse { return okResp(cmd, map[string]any{}) }}
	d := driverWith(fc)

	within := d.GetByID("panel")
	has := d.GetByText("badge")
	loc := d.Locator(Role("button"), Name("Save"), States("checked"), AndroidClassName("android.widget.Button")).
		Within(within).Has(has).Nth(0)

	if _, err := loc.Tap(); err != nil {
		t.Fatalf("Tap: %v", err)
	}
	var p locatorTapParams
	_ = json.Unmarshal(fc.sent[0].Params, &p)
	w := p.Locator
	if w.Role != "button" || w.Name != "Save" || len(w.States) != 1 || w.States[0] != "checked" {
		t.Fatalf("bad locator: %+v", w)
	}
	if w.Within == nil || w.Within.Id != "panel" {
		t.Fatalf("bad within: %+v", w.Within)
	}
	if w.Has == nil || w.Has.Text != "badge" {
		t.Fatalf("bad has: %+v", w.Has)
	}
	if w.Platform == nil || w.Platform.Android == nil || w.Platform.Android.ClassName != "android.widget.Button" {
		t.Fatalf("bad platform: %+v", w.Platform)
	}
	// Nth(0) must round-trip as an explicit 0, not an omitted field: a nil
	// pointer and a pointer-to-zero are different states on the wire.
	if w.Nth == nil || *w.Nth != 0 {
		t.Fatalf("want nth=0 present, got %v", w.Nth)
	}

	var raw map[string]any
	_ = json.Unmarshal(fc.sent[0].Params, &raw)
	locRaw, _ := raw["locator"].(map[string]any)
	if _, present := locRaw["nth"]; !present {
		t.Fatalf("wire locator omits nth even though it was set to 0: %s", fc.sent[0].Params)
	}
}

func TestLocatorFillPressWaitForBoundingBoxTextCount(t *testing.T) {
	cases := []struct {
		name       string
		call       func(l *Locator) error
		wantMethod string
		result     map[string]any
	}{
		{
			"Fill",
			func(l *Locator) error { _, err := l.Fill("hello"); return err },
			methodLocatorFill,
			map[string]any{"resolvedBy": "a11y", "bounds": map[string]int{"x": 0, "y": 0, "width": 1, "height": 1}, "tookMs": 1},
		},
		{
			"Press",
			func(l *Locator) error { _, err := l.Press(KeyEnter); return err },
			methodLocatorPress,
			map[string]any{"resolvedBy": "a11y", "bounds": map[string]int{"x": 0, "y": 0, "width": 1, "height": 1}, "tookMs": 1},
		},
		{
			"BoundingBox",
			func(l *Locator) error { _, err := l.BoundingBox(); return err },
			methodLocatorBoundingBox,
			map[string]any{"resolvedBy": "vlm", "bounds": map[string]int{"x": 0, "y": 0, "width": 1, "height": 1}, "tookMs": 1, "modelName": "m"},
		},
		{
			"Text",
			func(l *Locator) error { _, err := l.Text(); return err },
			methodLocatorText,
			map[string]any{"text": "hi", "resolvedBy": "ocr", "bounds": map[string]int{"x": 0, "y": 0, "width": 1, "height": 1}, "tookMs": 1},
		},
		{
			"Count",
			func(l *Locator) error { _, err := l.Count(); return err },
			methodLocatorCount,
			map[string]any{"count": 3, "resolvedBy": "ocr", "tookMs": 1},
		},
		{
			"WaitForVisible",
			func(l *Locator) error { _, err := l.WaitFor(StateVisible); return err },
			methodLocatorWaitFor,
			map[string]any{"resolvedBy": "ocr", "bounds": map[string]int{"x": 0, "y": 0, "width": 1, "height": 1}, "tookMs": 1},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			fc := &fakeConn{responder: func(cmd dcpCommand) dcpResponse { return okResp(cmd, c.result) }}
			d := driverWith(fc)
			if err := c.call(d.GetByText("x")); err != nil {
				t.Fatalf("%s: %v", c.name, err)
			}
			if fc.sent[0].Method != c.wantMethod {
				t.Fatalf("%s: wire method = %q, want %q", c.name, fc.sent[0].Method, c.wantMethod)
			}
		})
	}
}

func TestLocatorWaitForHiddenReturnsNilResult(t *testing.T) {
	fc := &fakeConn{responder: func(cmd dcpCommand) dcpResponse {
		return okResp(cmd, map[string]any{"tookMs": 5})
	}}
	d := driverWith(fc)
	res, err := d.GetByText("Spinner").WaitFor(StateHidden)
	if err != nil {
		t.Fatalf("WaitFor: %v", err)
	}
	if res != nil {
		t.Fatalf("want nil result for hidden, got %+v", res)
	}
}

func TestDriverPressWithNoLocator(t *testing.T) {
	fc := &fakeConn{responder: func(cmd dcpCommand) dcpResponse {
		return okResp(cmd, map[string]any{"tookMs": 3})
	}}
	d := driverWith(fc)
	res, err := d.Press(KeyEnter)
	if err != nil {
		t.Fatalf("Press: %v", err)
	}
	if res.ResolvedBy != "" || res.Bounds != (BBox{}) {
		t.Fatalf("want empty resolvedBy/bounds with no locator, got %+v", res)
	}
	var p locatorPressParams
	_ = json.Unmarshal(fc.sent[0].Params, &p)
	if p.Locator != nil {
		t.Fatalf("want nil locator on the wire, got %+v", p.Locator)
	}
	if p.Key != KeyEnter {
		t.Fatalf("want key %q, got %q", KeyEnter, p.Key)
	}
}

func TestLocatorMutatingCallsCarryIdempotencyKey(t *testing.T) {
	cases := []struct {
		name   string
		call   func(l *Locator) error
		mutate bool
	}{
		{"Tap", func(l *Locator) error { _, err := l.Tap(); return err }, true},
		{"Fill", func(l *Locator) error { _, err := l.Fill("x"); return err }, true},
		{"Press", func(l *Locator) error { _, err := l.Press(KeyEnter); return err }, true},
		{"WaitFor", func(l *Locator) error { _, err := l.WaitFor(StateVisible); return err }, false},
		{"BoundingBox", func(l *Locator) error { _, err := l.BoundingBox(); return err }, false},
		{"Text", func(l *Locator) error { _, err := l.Text(); return err }, false},
		{"Count", func(l *Locator) error { _, err := l.Count(); return err }, false},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			fc := &fakeConn{responder: func(cmd dcpCommand) dcpResponse {
				return okResp(cmd, map[string]any{"count": 0})
			}}
			d := driverWith(fc)
			if err := c.call(d.GetByText("x")); err != nil {
				t.Fatalf("%s: %v", c.name, err)
			}
			var raw map[string]any
			_ = json.Unmarshal(fc.sent[0].Params, &raw)
			_, hasKey := raw["idempotencyKey"]
			if hasKey != c.mutate {
				t.Fatalf("%s: idempotencyKey present=%v, want %v (params=%s)", c.name, hasKey, c.mutate, fc.sent[0].Params)
			}
		})
	}
}

func TestLocatorErrorMapping(t *testing.T) {
	cases := []struct {
		name string
		kind string
		code int
		is   func(error) bool
	}{
		{"ActionTimeout", kindActionTimeout, -32009, IsActionTimeout},
		{"StrategyUnavailable", kindStrategyUnavailable, -32010, IsStrategyUnavailable},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			fc := &fakeConn{responder: func(cmd dcpCommand) dcpResponse {
				return dcpResponse{ID: cmd.ID, Error: &dcpError{
					Code: c.code, Message: "boom",
					Data: &dcpErrorData{Kind: c.kind, Retryable: false},
				}}
			}}
			d := driverWith(fc)
			_, err := d.GetByRole("button").Tap()
			if !c.is(err) {
				t.Fatalf("%s: want matching error, got %v", c.name, err)
			}
			if IsRetryable(err) {
				t.Fatalf("%s: want not retryable, got retryable", c.name)
			}
		})
	}
}

func TestRefinementsOnlyNarrow(t *testing.T) {
	// A second Within keeps the first scope (chained, not replaced); a
	// second Filter appends its query; the receivers are unchanged.
	d := &MobileDriver{}
	base := d.GetByText("Save")
	loc := base.Within(d.GetByText("Dialog")).Within(d.GetByText("Card")).Filter("the primary one").Filter("enabled")
	w := loc.toWire()
	if w.Within == nil || w.Within.Text != "Card" || w.Within.Within == nil || w.Within.Within.Text != "Dialog" {
		t.Fatalf("within chain = %+v, want Card inside Dialog", w.Within)
	}
	if w.Query != "the primary one, enabled" {
		t.Fatalf("query = %q", w.Query)
	}
	if base.toWire().Within != nil || base.toWire().Query != "" {
		t.Fatal("refinement mutated its receiver")
	}
}

func TestLocatorTimeoutOutsideProtocolRangeFailsLocally(t *testing.T) {
	d := &MobileDriver{}
	_, err := d.GetByText("Save").Tap(WithTimeout(61 * time.Second))
	var me *Error
	if !errors.As(err, &me) || me.Code != CodeInvalidArgs {
		t.Fatalf("err = %v, want a local CodeInvalidArgs", err)
	}
}
