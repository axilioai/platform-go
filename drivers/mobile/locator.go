package mobile

import (
	"fmt"
	"strings"
)

// Locator is a lazy, immutable description of an on-screen target (the
// Locator tier: DCP's Locator.tap/fill/press/waitFor/boundingBox/text/count).
// Building one sends nothing; only an action or query method on it (Tap,
// Fill, Press, WaitFor, BoundingBox, Text, Count) reaches the device, and
// resolves the locator, auto-waits until actionable, and acts, all in one
// round trip. Nth, First, Within, Has and Filter each return a new Locator
// rather than mutating the receiver.
//
// How a locator resolves depends on the session. With accessibility mode on (the
// allocation default; see PhoneAllocateRequest.Accessibility and
// MobileDriver.Accessibility),
// the device matches the literal selectors (Role, Name, Text, ID, Value,
// States, ...) against the phone's accessibility tree, and a Query is
// ranked by a model over the tree's nodes. With it off, vision is the only
// resolver: a plain Text locator resolves by OCR, and a locator that also
// carries Query, Within, Has or Nth is resolved by one vision-model call,
// with a prompt composed from the whole locator. Under vision, Locator.Count
// needs a plain text locator and answers CodeInvalidArgs for one that also
// carries Query, Within or Has, since a vision-model call only resolves a
// single target per prompt.
//
// The tree-only selectors (Role, Name, ID, States, Value, WindowID, NodeID
// and the Android options) need the accessibility tree: on a session
// without one they answer CodeStrategyUnavailable whatever the strategy.
// They are never turned into a model prompt. Strategy picks the resolver
// explicitly (StrategyVision skips the tree even when it is on).
type Locator struct {
	driver *MobileDriver

	role, name, text, id, query, value, windowID, nodeID string
	androidClassName, androidPackageName                 string
	exact                                                bool
	states                                               []string
	nth                                                  *int
	within, has                                          *Locator

	// model, ocrEngine and strategy are this locator's own resolution
	// options (Model, OCREngine, Strategy). An action or query on the
	// locator resolves with these, falling back to the driver's
	// WithDefaultModel / WithDefaultOCREngine / WithDefaultStrategy, and
	// omitting the wire field entirely if neither is set. Within/Has only ever read the
	// selector fields off a nested locator: these three never travel with
	// it (one call resolves the whole locator, governed by the outer
	// locator's own options).
	model, ocrEngine string
	strategy         LocatorStrategy

	// buildErr is set by Within/Has when the locator passed in as the scope
	// carries its own Model/OCREngine/strategy (or already carries a
	// buildErr of its own): those options would otherwise be silently
	// dropped, so instead every action and query on this locator (and on
	// anything refined from it, since clone keeps the field) fails locally
	// with this error rather than sending a request that ignored what the
	// caller asked for.
	buildErr *Error
}

// LocatorOption configures a Locator at construction (via MobileDriver.Locator,
// GetByText, GetByRole or GetByID). Most are selector predicates (Text, Role,
// Name, ID, Exact, Query, ...), which combine as AND; a literal selector that
// matches nothing fails after the call's auto-wait. Model, OCREngine and
// Strategy are different: they are resolution options, setting how an
// action or query on this locator resolves it rather than narrowing what it
// matches.
type LocatorOption func(*Locator)

// Text matches visible text: substring and case-insensitive unless Exact is
// also given. Matched by OCR (GetByText is this option alone).
func Text(text string) LocatorOption { return func(l *Locator) { l.text = text } }

// Exact requires a case-sensitive exact match on Text, Name and Value rather
// than a case-insensitive substring.
func Exact() LocatorOption { return func(l *Locator) { l.exact = true } }

// Role matches the accessibility role, e.g. "button", "textbox". Needs the
// accessibility tree (see Locator).
func Role(role string) LocatorOption { return func(l *Locator) { l.role = role } }

// Name matches the accessible name: substring and case-insensitive unless
// Exact is also given. Needs the accessibility tree.
func Name(name string) LocatorOption { return func(l *Locator) { l.name = name } }

// ID matches a developer-assigned id, e.g. an Android resource id
// "com.example.app:id/save". Needs the accessibility tree.
func ID(id string) LocatorOption { return func(l *Locator) { l.id = id } }

// States requires the given accessibility states, e.g. "checked". Needs the
// accessibility tree.
func States(states ...string) LocatorOption {
	return func(l *Locator) { l.states = append([]string(nil), states...) }
}

// Value matches the node's current value (a text field's contents, a
// slider's value): substring and case-insensitive unless Exact is also
// given. Needs the accessibility tree.
func Value(value string) LocatorOption { return func(l *Locator) { l.value = value } }

// WindowID only matches nodes in this window, as AXTree.Windows lists it.
// Needs the accessibility tree.
func WindowID(windowID string) LocatorOption {
	return func(l *Locator) { l.windowID = windowID }
}

// NodeID matches exactly the node with this id from an earlier
// Accessibility snapshot. Fails with CodeStaleNode once that node is gone;
// it never matches a different element. Needs the accessibility tree.
func NodeID(nodeID string) LocatorOption { return func(l *Locator) { l.nodeID = nodeID } }

// Query is a natural-language description of the target. With the
// accessibility tree on, a model ranks the tree's nodes; otherwise (or under
// StrategyVision) the VLM reads it off the screen.
func Query(query string) LocatorOption { return func(l *Locator) { l.query = query } }

// AndroidClassName matches the Android view class, e.g.
// "android.widget.Button". Native matching; not portable across device
// classes. Needs the accessibility tree.
func AndroidClassName(className string) LocatorOption {
	return func(l *Locator) { l.androidClassName = className }
}

// AndroidPackageName only matches nodes owned by this app, e.g.
// "com.example.app". Native matching; not portable across device classes.
// Needs the accessibility tree.
func AndroidPackageName(packageName string) LocatorOption {
	return func(l *Locator) { l.androidPackageName = packageName }
}

// --- resolution options -------------------------------------------------

// Model sets the vision model this locator resolves with: the VLM that
// reads the target off the screen when the locator resolves by a
// vision-model call (see Locator). Locator value wins over the driver's
// WithDefaultModel; with neither set, a locator call omits the field
// entirely and the server's own default applies.
func Model(model string) LocatorOption { return func(l *Locator) { l.model = model } }

// OCREngine sets the OCR engine this locator resolves with when it resolves
// as a plain-text OCR match (see Locator). Locator value wins over the
// driver's WithDefaultOCREngine; with neither set, a locator call omits the
// field entirely and the server's own default applies.
func OCREngine(engine string) LocatorOption { return func(l *Locator) { l.ocrEngine = engine } }

// Strategy overrides the resolution strategy (StrategyAuto, StrategyVision
// or StrategyAccessibility) for this locator. Locator value wins over the
// driver's WithDefaultStrategy; with neither set, a locator call omits the
// field entirely and the server applies StrategyAuto.
func Strategy(strategy LocatorStrategy) LocatorOption {
	return func(l *Locator) { l.strategy = strategy }
}

// Locator builds a general-purpose locator from options. GetByText,
// GetByRole and GetByID are shorthand for the common single-predicate
// cases.
func (d *MobileDriver) Locator(opts ...LocatorOption) *Locator {
	l := &Locator{driver: d}
	for _, o := range opts {
		o(l)
	}
	return l
}

// GetByText builds a locator matching visible text (Text(text) plus any
// further opts, e.g. Exact()). Resolves by OCR.
func (d *MobileDriver) GetByText(text string, opts ...LocatorOption) *Locator {
	return d.Locator(append([]LocatorOption{Text(text)}, opts...)...)
}

// GetByRole builds a locator matching an accessibility role (Role(role)
// plus any further opts, typically Name("Log in")), e.g.
// d.GetByRole("button", mobile.Name("Log in")). Resolves against the
// accessibility tree; on a session without one, acting on it answers
// CodeStrategyUnavailable.
func (d *MobileDriver) GetByRole(role string, opts ...LocatorOption) *Locator {
	return d.Locator(append([]LocatorOption{Role(role)}, opts...)...)
}

// GetByID builds a locator matching a developer-assigned id, e.g. an
// Android resource id "com.example.app:id/login". Resolves against the
// accessibility tree; on a session without one, acting on it answers
// CodeStrategyUnavailable.
func (d *MobileDriver) GetByID(id string, opts ...LocatorOption) *Locator {
	return d.Locator(append([]LocatorOption{ID(id)}, opts...)...)
}

// Nth returns a new locator that picks the nth match in reading order
// (zero-based).
func (l *Locator) Nth(n int) *Locator {
	out := l.clone()
	out.nth = &n
	return out
}

// First returns a new locator pinned to the first match: shorthand for
// Nth(0).
func (l *Locator) First() *Locator { return l.Nth(0) }

// Within returns a new locator that only matches inside a match of other.
// Refinements only ever narrow: calling Within again scopes the new ancestor
// inside the earlier one (target in other, other in the previous scope)
// rather than dropping it.
//
// other contributes only its selector fields: the outer locator's own
// Model/OCREngine/Strategy govern how the whole thing resolves. If other
// itself carries one of those options, set on itself rather than inherited
// from a driver default, that would silently ignore what the caller asked
// for, so this records a build error on the returned locator instead: every
// action or query on it (and on anything further refined from it) fails
// locally with a CodeInvalidArgs *Error, without sending anything.
func (l *Locator) Within(other *Locator) *Locator {
	out := l.clone()
	if out.buildErr == nil {
		out.buildErr = scopeOptionsErr("Within", other)
	}
	out.within = chainScope(other, l.within, func(c *Locator) **Locator { return &c.within })
	return out
}

// Has returns a new locator that only matches elements containing a match of
// other. Calling Has again chains the new descendant onto the earlier one
// rather than dropping it, so the locator only ever narrows.
//
// other contributes only its selector fields; see Within for what happens
// when it also carries its own Model/OCREngine/Strategy.
func (l *Locator) Has(other *Locator) *Locator {
	out := l.clone()
	if out.buildErr == nil {
		out.buildErr = scopeOptionsErr("Has", other)
	}
	out.has = chainScope(other, l.has, func(c *Locator) **Locator { return &c.has })
	return out
}

// scopeOptionsErr reports the build error other's own resolution options (or
// its own pre-existing build error) put on a locator that adopts it as a
// Within/Has scope, or nil if other is a plain selector. who is "Within" or
// "Has", named in the resulting message.
func scopeOptionsErr(who string, other *Locator) *Error {
	if other == nil {
		return nil
	}
	if other.buildErr != nil {
		return other.buildErr
	}
	set := other.resolutionOptionsSet()
	if len(set) == 0 {
		return nil
	}
	return &Error{
		Code: CodeInvalidArgs,
		Message: fmt.Sprintf(
			"%s: the inner locator sets %s; set Model, OCREngine and Strategy on the outer locator instead",
			who, strings.Join(set, ", "),
		),
	}
}

// resolutionOptionsSet names the resolution options set on l itself (not
// inherited from a driver default), in a fixed order.
func (l *Locator) resolutionOptionsSet() []string {
	var set []string
	if l.model != "" {
		set = append(set, "Model")
	}
	if l.ocrEngine != "" {
		set = append(set, "OCREngine")
	}
	if l.strategy != "" {
		set = append(set, "Strategy")
	}
	return set
}

// Filter returns a new locator that keeps the receiver's literal selectors
// (which narrow the candidates) and adds a natural-language query the model
// ranks the survivors by. A second Filter appends to the first query rather
// than replacing it.
func (l *Locator) Filter(query string) *Locator {
	out := l.clone()
	if out.query == "" {
		out.query = query
	} else {
		out.query += ", " + query
	}
	return out
}

// chainScope returns a copy of scope with prev attached at the innermost end
// of its chain (followed through field), so an earlier Within/Has scope is
// kept, not replaced. The receivers are never mutated.
func chainScope(scope, prev *Locator, field func(*Locator) **Locator) *Locator {
	if prev == nil {
		return scope
	}
	if scope == nil {
		return prev
	}
	c := scope.clone()
	next := field(c)
	*next = chainScope(*next, prev, field)
	return c
}

func (l *Locator) clone() *Locator {
	c := *l
	return &c
}

// toWire converts a possibly-nil Locator into the generated wire shape
// (nil in, nil out; used directly for the optional locator on Locator.press
// and for the recursive within/has).
func (l *Locator) toWire() *locatorWire {
	if l == nil {
		return nil
	}
	w := &locatorWire{
		Role:     l.role,
		Name:     l.name,
		Text:     l.text,
		Exact:    l.exact,
		Id:       l.id,
		States:   l.states,
		Query:    l.query,
		Value:    l.value,
		WindowId: l.windowID,
		NodeId:   l.nodeID,
		Nth:      l.nth,
		Within:   l.within.toWire(),
		Has:      l.has.toWire(),
	}
	if l.androidClassName != "" || l.androidPackageName != "" {
		w.Platform = &locatorPlatformWire{Android: &androidLocatorWire{
			ClassName:   l.androidClassName,
			PackageName: l.androidPackageName,
		}}
	}
	return w
}
