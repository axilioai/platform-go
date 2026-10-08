package mobile

import (
	"fmt"
	"strings"
	"time"
)

// Accessibility is the phone's accessibility tree (the DCP Accessibility
// domain): read it as a snapshot, query it, and turn it on or off for the
// session. Get one with MobileDriver.Accessibility; it holds no state of its
// own, so it is cheap to call that each time.
//
// Locators do not need any of this: with the tree on, GetByRole, GetByID and
// the other literal selectors already resolve against it on the device.
// Accessibility is for reading the tree itself, e.g. to see what is on
// screen or to pin a node by its id (NodeID) for a later action.
//
// Every read fails with CodeStrategyUnavailable while the tree is off, and
// with CodeTreeUnavailable while a system dialog covers the app. On a
// device class that has no Accessibility domain at all, every call answers
// CodeUnknownOp.
type Accessibility struct {
	driver *MobileDriver
}

// Accessibility returns the accessibility tree accessor for this driver.
func (d *MobileDriver) Accessibility() *Accessibility { return &Accessibility{driver: d} }

// AccessibilityState is whether the tree is on for this session, and
// whether this session can turn it on and off.
type AccessibilityState struct {
	// Enabled reports the tree is on.
	Enabled bool `json:"enabled"`
	// Toggleable reports this session can call Enable and Disable.
	Toggleable bool `json:"toggleable"`
}

// AXTree is one snapshot of the accessibility tree: every node of the
// windows read, in document order, plus the windows themselves.
type AXTree struct {
	Nodes   []AXNode   `json:"nodes"`
	Windows []AXWindow `json:"windows"`
	// CapturedAt is when the device read the tree.
	CapturedAt time.Time `json:"captured_at"`
}

// Node returns the node with id, or nil if the snapshot has none.
func (t *AXTree) Node(id string) *AXNode {
	for i := range t.Nodes {
		if t.Nodes[i].NodeID == id {
			return &t.Nodes[i]
		}
	}
	return nil
}

// AXNode is one accessibility node: CDP's AXNode plus where it is on screen,
// which window it belongs to, what actions it supports, and its native
// attributes.
type AXNode struct {
	// NodeID is stable while the element lives. Pass it to NodeID (the
	// locator option), Partial or Children.
	NodeID string `json:"node_id"`
	// Ignored marks a layout-only node (CDP's ignored semantics).
	Ignored     bool         `json:"ignored"`
	Role        AXValue      `json:"role"`
	Name        *AXValue     `json:"name,omitempty"`
	Description *AXValue     `json:"description,omitempty"`
	Value       *AXValue     `json:"value,omitempty"`
	Properties  []AXProperty `json:"properties,omitempty"`
	ParentID    string       `json:"parent_id,omitempty"`
	ChildIDs    []string     `json:"child_ids"`
	// Bounds are frame pixels, the space the Touch methods address.
	Bounds   BBox   `json:"bounds"`
	WindowID string `json:"window_id"`
	// Actions are portable action names, e.g. "click", "longClick",
	// "scroll".
	Actions  []string        `json:"actions"`
	Platform *AXNodePlatform `json:"platform,omitempty"`
}

// AXValue is a CDP AXValue: a typed value whose Value is a string, bool or
// number.
type AXValue struct {
	// Type is the CDP AXValueType, e.g. "role", "computedString", "boolean".
	Type  string `json:"type"`
	Value any    `json:"value,omitempty"`
}

// String renders the value as text: "" for a nil AXValue or a missing
// value, so n.Name.String() is safe on a node with no name.
func (v *AXValue) String() string {
	if v == nil || v.Value == nil {
		return ""
	}
	if s, ok := v.Value.(string); ok {
		return s
	}
	return fmt.Sprint(v.Value)
}

// AXProperty is one CDP node property, e.g. "focusable", "disabled",
// "checked".
type AXProperty struct {
	Name  string  `json:"name"`
	Value AXValue `json:"value"`
}

// AXNodePlatform carries a node's native attributes, keyed by platform. Not
// portable across device classes.
type AXNodePlatform struct {
	Android *AXAndroidNode `json:"android,omitempty"`
}

// AXAndroidNode is an Android node's native attributes.
type AXAndroidNode struct {
	ClassName          string `json:"class_name,omitempty"`
	ViewIDResourceName string `json:"view_id_resource_name,omitempty"`
	PackageName        string `json:"package_name,omitempty"`
	Text               string `json:"text,omitempty"`
	ContentDescription string `json:"content_description,omitempty"`
	HintText           string `json:"hint_text,omitempty"`
}

// AXWindowType is the kind of window an AXWindow is.
type AXWindowType string

// Window types (the wire's AXWindow.type).
const (
	AXWindowApplication AXWindowType = "application"
	AXWindowInputMethod AXWindowType = "inputMethod"
	AXWindowSystem      AXWindowType = "system"
	AXWindowOverlay     AXWindowType = "overlay"
	AXWindowOther       AXWindowType = "other"
)

// AXWindow is one on-screen window in a snapshot.
type AXWindow struct {
	WindowID string       `json:"window_id"`
	Title    string       `json:"title,omitempty"`
	Type     AXWindowType `json:"type"`
	// App is the package that owns the window, when known.
	App     string `json:"app,omitempty"`
	Focused bool   `json:"focused"`
	Bounds  BBox   `json:"bounds"`
	// RootID is the id of the window's root node, when the snapshot read it.
	RootID string `json:"root_id,omitempty"`
}

// AXQuery selects nodes for Accessibility.Query. Every set field must
// match (AND); an empty AXQuery matches every visible node.
type AXQuery struct {
	// Role matches the accessibility role, e.g. "button".
	Role string
	// Name matches the exact computed accessible name.
	Name string
	// Selector is a locator matched the same way the Locator methods match
	// it (Query is not accepted here). Its Model, OCREngine and Strategy
	// would have no effect, so setting any of them fails locally with
	// CodeInvalidArgs.
	Selector *Locator
}

// --- snapshot options -------------------------------------------------------

type snapshotConfig struct {
	timeout         time.Duration
	depth           *int
	windowID        string
	interestingOnly *bool
}

// SnapshotOption tunes one Accessibility.Snapshot call. WithTimeout is one;
// WithDepth, WithWindow and WithInterestingOnly are the others.
type SnapshotOption interface {
	applySnapshot(*snapshotConfig)
}

type snapshotOptionFunc func(*snapshotConfig)

func (f snapshotOptionFunc) applySnapshot(c *snapshotConfig) { f(c) }

func (o TimeoutOption) applySnapshot(c *snapshotConfig) { c.timeout = time.Duration(o) }

// WithDepth limits a snapshot to depth levels below each window root; 0 is
// the roots only. Without it the snapshot is the whole tree.
func WithDepth(depth int) SnapshotOption {
	return snapshotOptionFunc(func(c *snapshotConfig) { c.depth = &depth })
}

// WithWindow limits a snapshot to one window (an AXWindow.WindowID).
// Without it the snapshot covers every window.
func WithWindow(windowID string) SnapshotOption {
	return snapshotOptionFunc(func(c *snapshotConfig) { c.windowID = windowID })
}

// WithInterestingOnly sets whether a snapshot drops layout-only nodes. The
// device drops them by default; pass false for every node.
func WithInterestingOnly(interestingOnly bool) SnapshotOption {
	return snapshotOptionFunc(func(c *snapshotConfig) { c.interestingOnly = &interestingOnly })
}

// --- calls -----------------------------------------------------------------

// State reports whether the tree is on and whether this session can toggle
// it.
func (a *Accessibility) State(opts ...ActionOption) (AccessibilityState, error) {
	cfg := newActionConfig(defaultCallTimeout, opts)
	raw, err := a.driver.call(methodAccessibilityGetState, nil, cfg.timeout)
	if err != nil {
		return AccessibilityState{}, err
	}
	var w wireAccessibilityState
	if err := unmarshalResult(raw, &w); err != nil {
		return AccessibilityState{}, err
	}
	return AccessibilityState(w), nil
}

// Enable turns the tree on for this session and returns once the device
// confirms it is on. Only a session whose State is Toggleable can call it;
// elsewhere it answers CodeUnknownOp. While on, the accessibility service is
// visible to apps on the phone.
func (a *Accessibility) Enable(opts ...ActionOption) error {
	return a.toggle(methodAccessibilityEnable, opts)
}

// Disable turns the tree off for this session and returns once the device
// confirms it is off. See Enable.
func (a *Accessibility) Disable(opts ...ActionOption) error {
	return a.toggle(methodAccessibilityDisable, opts)
}

func (a *Accessibility) toggle(method string, opts []ActionOption) error {
	cfg := newActionConfig(defaultCallTimeout, opts)
	_, err := a.driver.call(method, accessibilityToggleParams{}, cfg.timeout)
	return err
}

// Snapshot reads the tree: every window (or one, WithWindow), in document
// order, without layout-only nodes unless WithInterestingOnly(false).
func (a *Accessibility) Snapshot(opts ...SnapshotOption) (*AXTree, error) {
	cfg := snapshotConfig{timeout: defaultCallTimeout}
	for _, o := range opts {
		o.applySnapshot(&cfg)
	}
	if cfg.timeout <= 0 {
		cfg.timeout = defaultCallTimeout
	}
	if cfg.depth != nil && *cfg.depth < 0 {
		return nil, &Error{Code: CodeInvalidArgs, Message: fmt.Sprintf("snapshot depth %d is negative", *cfg.depth)}
	}
	params := getFullAXTreeParams{
		Depth:           cfg.depth,
		WindowId:        cfg.windowID,
		InterestingOnly: cfg.interestingOnly,
	}
	raw, err := a.driver.call(methodAccessibilityGetFullAXTree, params, cfg.timeout)
	if err != nil {
		return nil, err
	}
	var w wireAXTree
	if err := unmarshalResult(raw, &w); err != nil {
		return nil, err
	}
	return axTreeFromWire(w), nil
}

// Query returns every visible node matching q, in reading order. Never
// waits: use a Locator's WaitFor first if the target may not be on screen
// yet.
func (a *Accessibility) Query(q AXQuery, opts ...ActionOption) ([]AXNode, error) {
	if err := axQuerySelectorErr(q.Selector); err != nil {
		return nil, err
	}
	cfg := newActionConfig(defaultCallTimeout, opts)
	params := queryAXTreeParams{
		AccessibleName: q.Name,
		Role:           q.Role,
		Selector:       q.Selector.toWire(),
	}
	return a.nodes(methodAccessibilityQueryAXTree, params, cfg.timeout)
}

// axQuerySelectorErr rejects a Query selector that carries its own build
// error, or resolution options a tree query would silently ignore.
func axQuerySelectorErr(sel *Locator) error {
	if sel == nil {
		return nil
	}
	if sel.buildErr != nil {
		return sel.buildErr
	}
	if set := sel.resolutionOptionsSet(); len(set) > 0 {
		return &Error{
			Code: CodeInvalidArgs,
			Message: fmt.Sprintf("Query: the selector sets %s, which a tree query does not use; drop them",
				strings.Join(set, ", ")),
		}
	}
	return nil
}

// Partial returns one node by id and, when fetchRelatives is true, its
// ancestors, siblings and children too. Fails with CodeStaleNode once the
// node is gone.
func (a *Accessibility) Partial(nodeID string, fetchRelatives bool, opts ...ActionOption) ([]AXNode, error) {
	cfg := newActionConfig(defaultCallTimeout, opts)
	params := getPartialAXTreeParams{NodeId: nodeID, FetchRelatives: &fetchRelatives}
	return a.nodes(methodAccessibilityGetPartialAXTree, params, cfg.timeout)
}

// Children returns the direct children of the node with id nodeID. Fails
// with CodeStaleNode once the node is gone.
func (a *Accessibility) Children(nodeID string, opts ...ActionOption) ([]AXNode, error) {
	cfg := newActionConfig(defaultCallTimeout, opts)
	return a.nodes(methodAccessibilityGetChildAXNodes, getChildAXNodesParams{Id: nodeID}, cfg.timeout)
}

func (a *Accessibility) nodes(method string, params any, timeout time.Duration) ([]AXNode, error) {
	raw, err := a.driver.call(method, params, timeout)
	if err != nil {
		return nil, err
	}
	var w wireAXNodes
	if err := unmarshalResult(raw, &w); err != nil {
		return nil, err
	}
	return axNodesFromWire(w.Nodes), nil
}

// --- wire result frames (camelCase). Results stay hand-written, like the
// other domains': they convert into the exported snake_case types above. ---

type wireAccessibilityState struct {
	Enabled    bool `json:"enabled"`
	Toggleable bool `json:"toggleable"`
}

type wireAXValue struct {
	Type  string `json:"type"`
	Value any    `json:"value,omitempty"`
}

type wireAXProperty struct {
	Name  string      `json:"name"`
	Value wireAXValue `json:"value"`
}

type wireAXAndroidNode struct {
	ClassName          string `json:"className,omitempty"`
	ViewIDResourceName string `json:"viewIdResourceName,omitempty"`
	PackageName        string `json:"packageName,omitempty"`
	Text               string `json:"text,omitempty"`
	ContentDescription string `json:"contentDescription,omitempty"`
	HintText           string `json:"hintText,omitempty"`
}

type wireAXNode struct {
	NodeID      string           `json:"nodeId"`
	Ignored     bool             `json:"ignored"`
	Role        wireAXValue      `json:"role"`
	Name        *wireAXValue     `json:"name,omitempty"`
	Description *wireAXValue     `json:"description,omitempty"`
	Value       *wireAXValue     `json:"value,omitempty"`
	Properties  []wireAXProperty `json:"properties,omitempty"`
	ParentID    string           `json:"parentId,omitempty"`
	ChildIDs    []string         `json:"childIds"`
	Bounds      wireBBox         `json:"bounds"`
	WindowID    string           `json:"windowId"`
	Actions     []string         `json:"actions"`
	Platform    *struct {
		Android *wireAXAndroidNode `json:"android,omitempty"`
	} `json:"platform,omitempty"`
}

type wireAXWindow struct {
	WindowID string   `json:"windowId"`
	Title    string   `json:"title,omitempty"`
	Type     string   `json:"type"`
	App      string   `json:"app,omitempty"`
	Focused  bool     `json:"focused"`
	Bounds   wireBBox `json:"bounds"`
	RootID   string   `json:"rootId,omitempty"`
}

type wireAXTree struct {
	Nodes      []wireAXNode   `json:"nodes"`
	Windows    []wireAXWindow `json:"windows"`
	CapturedAt int64          `json:"capturedAt"` // epoch-millis
}

type wireAXNodes struct {
	Nodes []wireAXNode `json:"nodes"`
}

func axValueFromWire(w *wireAXValue) *AXValue {
	if w == nil {
		return nil
	}
	return &AXValue{Type: w.Type, Value: w.Value}
}

func axNodeFromWire(w wireAXNode) AXNode {
	n := AXNode{
		NodeID:      w.NodeID,
		Ignored:     w.Ignored,
		Role:        AXValue{Type: w.Role.Type, Value: w.Role.Value},
		Name:        axValueFromWire(w.Name),
		Description: axValueFromWire(w.Description),
		Value:       axValueFromWire(w.Value),
		ParentID:    w.ParentID,
		ChildIDs:    append([]string{}, w.ChildIDs...),
		Bounds:      bboxFromWire(w.Bounds),
		WindowID:    w.WindowID,
		Actions:     append([]string{}, w.Actions...),
	}
	for _, p := range w.Properties {
		n.Properties = append(n.Properties, AXProperty{Name: p.Name, Value: AXValue{Type: p.Value.Type, Value: p.Value.Value}})
	}
	if w.Platform != nil && w.Platform.Android != nil {
		an := AXAndroidNode(*w.Platform.Android)
		n.Platform = &AXNodePlatform{Android: &an}
	}
	return n
}

func axNodesFromWire(ws []wireAXNode) []AXNode {
	out := make([]AXNode, 0, len(ws))
	for _, w := range ws {
		out = append(out, axNodeFromWire(w))
	}
	return out
}

func axTreeFromWire(w wireAXTree) *AXTree {
	t := &AXTree{
		Nodes:      axNodesFromWire(w.Nodes),
		Windows:    make([]AXWindow, 0, len(w.Windows)),
		CapturedAt: time.UnixMilli(w.CapturedAt).UTC(),
	}
	for _, win := range w.Windows {
		t.Windows = append(t.Windows, AXWindow{
			WindowID: win.WindowID,
			Title:    win.Title,
			Type:     AXWindowType(win.Type),
			App:      win.App,
			Focused:  win.Focused,
			Bounds:   bboxFromWire(win.Bounds),
			RootID:   win.RootID,
		})
	}
	return t
}
