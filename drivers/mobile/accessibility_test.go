package mobile

import (
	"encoding/json"
	"testing"
	"time"
)

// axNodeWire is one node in the camelCase shape the device sends.
func axNodeWire(id, role, name string, children ...string) map[string]any {
	return map[string]any{
		"nodeId":   id,
		"ignored":  false,
		"role":     map[string]any{"type": "role", "value": role},
		"name":     map[string]any{"type": "computedString", "value": name},
		"childIds": children,
		"bounds":   map[string]int{"x": 10, "y": 20, "width": 100, "height": 40},
		"windowId": "w1",
		"actions":  []string{"click"},
	}
}

func TestSnapshotDecodesTree(t *testing.T) {
	root := axNodeWire("n1", "generic", "", "n2")
	button := axNodeWire("n2", "button", "Log in")
	button["parentId"] = "n1"
	button["properties"] = []map[string]any{{"name": "focusable", "value": map[string]any{"type": "boolean", "value": true}}}
	button["platform"] = map[string]any{"android": map[string]any{
		"className": "android.widget.Button", "viewIdResourceName": "com.example.app:id/login", "packageName": "com.example.app",
	}}
	fc := &fakeConn{responder: func(cmd dcpCommand) dcpResponse {
		return okResp(cmd, map[string]any{
			"nodes": []any{root, button},
			"windows": []any{map[string]any{
				"windowId": "w1", "type": "application", "app": "com.example.app", "focused": true,
				"bounds": map[string]int{"x": 0, "y": 0, "width": 1080, "height": 2400}, "rootId": "n1",
			}},
			"capturedAt": 1_700_000_000_000,
		})
	}}
	d := driverWith(fc)

	tree, err := d.Accessibility().Snapshot()
	if err != nil {
		t.Fatalf("Snapshot: %v", err)
	}
	if fc.sent[0].Method != methodAccessibilityGetFullAXTree {
		t.Fatalf("wire method = %q, want %q", fc.sent[0].Method, methodAccessibilityGetFullAXTree)
	}
	// No options: every field is omitted so the device defaults apply (the
	// whole tree, every window, interesting nodes only).
	var raw map[string]any
	_ = json.Unmarshal(fc.sent[0].Params, &raw)
	if len(raw) != 0 {
		t.Fatalf("want empty params, got %s", fc.sent[0].Params)
	}

	if len(tree.Nodes) != 2 || len(tree.Windows) != 1 {
		t.Fatalf("bad tree: %+v", tree)
	}
	if !tree.CapturedAt.Equal(time.UnixMilli(1_700_000_000_000)) {
		t.Fatalf("bad capturedAt: %v", tree.CapturedAt)
	}
	n := tree.Node("n2")
	if n == nil {
		t.Fatal("Node(n2) = nil")
	}
	if n.Role.String() != "button" || n.Name.String() != "Log in" || n.ParentID != "n1" || n.WindowID != "w1" {
		t.Fatalf("bad node: %+v", n)
	}
	if n.Bounds != (BBox{X: 10, Y: 20, Width: 100, Height: 40}) || n.Bounds.Center() != (Coords{X: 60, Y: 40}) {
		t.Fatalf("bad bounds: %+v", n.Bounds)
	}
	if len(n.Properties) != 1 || n.Properties[0].Name != "focusable" || n.Properties[0].Value.String() != "true" {
		t.Fatalf("bad properties: %+v", n.Properties)
	}
	if n.Platform == nil || n.Platform.Android == nil || n.Platform.Android.ViewIDResourceName != "com.example.app:id/login" {
		t.Fatalf("bad platform: %+v", n.Platform)
	}
	if n.Description.String() != "" || n.Value != nil {
		t.Fatalf("want absent description/value, got %+v / %+v", n.Description, n.Value)
	}
	w := tree.Windows[0]
	if w.Type != AXWindowApplication || w.RootID != "n1" || !w.Focused || w.App != "com.example.app" {
		t.Fatalf("bad window: %+v", w)
	}
	if tree.Node("missing") != nil {
		t.Fatal("Node(missing) should be nil")
	}
}

func TestSnapshotOptionsOnTheWire(t *testing.T) {
	fc := &fakeConn{responder: func(cmd dcpCommand) dcpResponse {
		return okResp(cmd, map[string]any{"nodes": []any{}, "windows": []any{}, "capturedAt": 0})
	}}
	d := driverWith(fc)

	// depth=0 and interestingOnly=false are meaningful values that differ
	// from omitting the field, so both must reach the wire explicitly.
	if _, err := d.Accessibility().Snapshot(WithDepth(0), WithWindow("w9"), WithInterestingOnly(false), WithTimeout(time.Second)); err != nil {
		t.Fatalf("Snapshot: %v", err)
	}
	var raw map[string]any
	_ = json.Unmarshal(fc.sent[0].Params, &raw)
	if raw["depth"] != float64(0) || raw["windowId"] != "w9" || raw["interestingOnly"] != false {
		t.Fatalf("bad params: %s", fc.sent[0].Params)
	}

	if _, err := d.Accessibility().Snapshot(WithDepth(-1)); !hasCode(err, CodeInvalidArgs) {
		t.Fatalf("negative depth: want InvalidArgs, got %v", err)
	}
	if len(fc.sent) != 1 {
		t.Fatalf("a locally rejected snapshot must send nothing, sent=%d", len(fc.sent))
	}
}

func TestAccessibilityStateEnableDisable(t *testing.T) {
	fc := &fakeConn{responder: func(cmd dcpCommand) dcpResponse {
		if cmd.Method == methodAccessibilityGetState {
			return okResp(cmd, map[string]any{"enabled": true, "toggleable": true})
		}
		return okResp(cmd, map[string]any{})
	}}
	d := driverWith(fc)
	a := d.Accessibility()

	st, err := a.State()
	if err != nil {
		t.Fatalf("State: %v", err)
	}
	if st != (AccessibilityState{Enabled: true, Toggleable: true}) {
		t.Fatalf("bad state: %+v", st)
	}
	if err := a.Disable(); err != nil {
		t.Fatalf("Disable: %v", err)
	}
	if err := a.Enable(); err != nil {
		t.Fatalf("Enable: %v", err)
	}

	wantMethods := []string{methodAccessibilityGetState, methodAccessibilityDisable, methodAccessibilityEnable}
	for i, m := range wantMethods {
		if fc.sent[i].Method != m {
			t.Fatalf("call %d: method = %q, want %q", i, fc.sent[i].Method, m)
		}
	}
	// Enable/Disable mutate the device, so they carry an idempotency key;
	// getState is a read and stays keyless.
	for i, wantKey := range []bool{false, true, true} {
		var raw map[string]any
		_ = json.Unmarshal(fc.sent[i].Params, &raw)
		if _, has := raw["idempotencyKey"]; has != wantKey {
			t.Fatalf("%s: idempotencyKey present=%v, want %v", fc.sent[i].Method, has, wantKey)
		}
	}
}

func TestAccessibilityQueryPartialChildren(t *testing.T) {
	fc := &fakeConn{responder: func(cmd dcpCommand) dcpResponse {
		return okResp(cmd, map[string]any{"nodes": []any{axNodeWire("n2", "button", "Log in")}})
	}}
	d := driverWith(fc)
	a := d.Accessibility()

	nodes, err := a.Query(AXQuery{Role: "button", Name: "Log in", Selector: d.Locator(AndroidPackageName("com.example.app"))})
	if err != nil {
		t.Fatalf("Query: %v", err)
	}
	if len(nodes) != 1 || nodes[0].NodeID != "n2" {
		t.Fatalf("bad nodes: %+v", nodes)
	}
	var q queryAXTreeParams
	_ = json.Unmarshal(fc.sent[0].Params, &q)
	if fc.sent[0].Method != methodAccessibilityQueryAXTree || q.Role != "button" || q.AccessibleName != "Log in" ||
		q.Selector == nil || q.Selector.Platform.Android.PackageName != "com.example.app" {
		t.Fatalf("bad query params: %s", fc.sent[0].Params)
	}

	if _, err := a.Partial("n2", false); err != nil {
		t.Fatalf("Partial: %v", err)
	}
	var raw map[string]any
	_ = json.Unmarshal(fc.sent[1].Params, &raw)
	// fetchRelatives defaults to true on the device, so false must be sent.
	if fc.sent[1].Method != methodAccessibilityGetPartialAXTree || raw["nodeId"] != "n2" || raw["fetchRelatives"] != false {
		t.Fatalf("bad partial params: %s", fc.sent[1].Params)
	}

	if _, err := a.Children("n1"); err != nil {
		t.Fatalf("Children: %v", err)
	}
	_ = json.Unmarshal(fc.sent[2].Params, &raw)
	if fc.sent[2].Method != methodAccessibilityGetChildAXNodes || raw["id"] != "n1" {
		t.Fatalf("bad children params: %s", fc.sent[2].Params)
	}
}

func TestAccessibilityQueryRejectsSelectorResolutionOptions(t *testing.T) {
	fc := &fakeConn{responder: func(cmd dcpCommand) dcpResponse { return okResp(cmd, map[string]any{}) }}
	d := driverWith(fc)

	_, err := d.Accessibility().Query(AXQuery{Selector: d.GetByRole("button", Strategy(StrategyVision))})
	if !hasCode(err, CodeInvalidArgs) {
		t.Fatalf("want InvalidArgs, got %v", err)
	}
	if len(fc.sent) != 0 {
		t.Fatalf("a locally rejected query must send nothing, sent=%d", len(fc.sent))
	}
}

func TestAccessibilityErrorsMap(t *testing.T) {
	cases := []struct {
		kind string
		code int
		is   func(error) bool
	}{
		{kindStrategyUnavailable, -32010, IsStrategyUnavailable},
		{kindTreeUnavailable, -32011, IsTreeUnavailable},
		{kindStaleNode, -32012, IsStaleNode},
	}
	for _, c := range cases {
		t.Run(c.kind, func(t *testing.T) {
			fc := &fakeConn{responder: func(cmd dcpCommand) dcpResponse {
				return dcpResponse{ID: cmd.ID, Error: &dcpError{
					Code: c.code, Message: "boom", Data: &dcpErrorData{Kind: c.kind},
				}}
			}}
			_, err := driverWith(fc).Accessibility().Children("n1")
			if !c.is(err) {
				t.Fatalf("want %s, got %v", c.kind, err)
			}
		})
	}
}
