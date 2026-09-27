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
// Today's phones have no accessibility tree, so vision (OCR/VLM) is the only
// resolver: a plain Text locator resolves by OCR; a locator that also
// carries Query, Within, Has or Nth is instead resolved by one vision-model
// call, with a prompt composed from the whole locator (so Nth on a
// Query-based locator works). Locator.Count is the exception: it needs a
// plain text locator, and answers CodeInvalidArgs for one that also carries
// Query, Within or Has, since counting needs every independent match and a
// vision-model call only resolves a single target per prompt. Role/id
// selectors and a choice of resolution strategy arrive with accessibility
// support in a later release.
type Locator struct {
	driver *MobileDriver

	role, name, text, id, query, androidClassName string
	exact                                         bool
	states                                        []string
	nth                                           *int
	within, has                                   *Locator

	// model, ocrEngine and strategy are this locator's own resolution
	// options (Model/OCREngine, plus the unexported strategy held back until
	// accessibility ships). An action or query on the locator resolves with
	// these, falling back to the driver's WithDefaultModel /
	// WithDefaultOCREngine / withDefaultStrategy, and omitting the wire
	// field entirely if neither is set. Within/Has only ever read the
	// selector fields off a nested locator: these three never travel with
	// it (one call resolves the whole locator, governed by the outer
	// locator's own options).
	model, ocrEngine string
	strategy         locatorStrategy

	// buildErr is set by Within/Has when the locator passed in as the scope
	// carries its own Model/OCREngine/strategy (or already carries a
	// buildErr of its own): those options would otherwise be silently
	// dropped, so instead every action and query on this locator (and on
	// anything refined from it, since clone keeps the field) fails locally
	// with this error rather than sending a request that ignored what the
	// caller asked for.
	buildErr *Error
}

// LocatorOption configures a Locator at construction (via MobileDriver.Locator)
// or refinement (via Locator.Filter). Most are selector predicates (Text,
// Exact, Query), which combine as AND; a literal selector that matches
// nothing fails after the call's auto-wait. Model and OCREngine are
// different: they are resolution options, setting how an action or query on
// this locator resolves it rather than narrowing what it matches.
type LocatorOption func(*Locator)

// Text matches visible text: substring and case-insensitive unless Exact is
// also given. Matched by OCR (GetByText is this option alone).
func Text(text string) LocatorOption { return func(l *Locator) { l.text = text } }

// Exact requires an exact match on Text rather than a case-insensitive
// substring.
func Exact() LocatorOption { return func(l *Locator) { l.exact = true } }

// withRole matches the accessibility role, e.g. "button", "textbox". Needs
// the accessibility tree: on a session without one this answers
// StrategyUnavailable regardless of strategy (today's phones have none).
// Held back until accessibility support ships.
func withRole(role string) LocatorOption { return func(l *Locator) { l.role = role } }

// withName matches the accessible name. Needs the accessibility tree; see
// withRole. Held back until accessibility support ships.
func withName(name string) LocatorOption { return func(l *Locator) { l.name = name } }

// withID matches a developer-assigned id, e.g. an Android resource id
// "com.example.app:id/save". Needs the accessibility tree; see withRole.
// Held back until accessibility support ships.
func withID(id string) LocatorOption { return func(l *Locator) { l.id = id } }

// withStates requires the given accessibility states, e.g. "checked". Needs
// the accessibility tree; see withRole. Held back until accessibility
// support ships.
func withStates(states ...string) LocatorOption {
	return func(l *Locator) { l.states = append([]string(nil), states...) }
}

// Query is a natural-language description of the target, read from the
// screen by the VLM.
func Query(query string) LocatorOption { return func(l *Locator) { l.query = query } }

// withAndroidClassName matches the Android view class, e.g.
// "android.widget.Button". Native matching; not portable across device
// classes. Held back until accessibility support ships.
func withAndroidClassName(className string) LocatorOption {
	return func(l *Locator) { l.androidClassName = className }
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

// withStrategy overrides the resolution strategy (strategyAuto,
// strategyVision or strategyAccessibility) for this locator. Locator value
// wins over the driver's withDefaultStrategy; with neither set, a locator
// call omits the field entirely and the server applies strategyAuto. Held
// back until accessibility support ships.
func withStrategy(strategy locatorStrategy) LocatorOption {
	return func(l *Locator) { l.strategy = strategy }
}

// Locator builds a general-purpose locator from options. GetByText is
// shorthand for the common single-predicate case; getByRole and getByID
// (unexported, held back until accessibility support ships) are the
// equivalents for role/id selectors.
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

// getByRole builds a locator matching an accessibility role. Today's phones
// have no accessibility tree, so acting on it answers StrategyUnavailable by
// design until a phone advertises one. Held back until accessibility
// support ships.
func (d *MobileDriver) getByRole(role string, opts ...LocatorOption) *Locator {
	return d.Locator(append([]LocatorOption{withRole(role)}, opts...)...)
}

// getByID builds a locator matching a developer-assigned id. Today's phones
// have no accessibility tree, so acting on it answers StrategyUnavailable by
// design until a phone advertises one. Held back until accessibility
// support ships.
func (d *MobileDriver) getByID(id string, opts ...LocatorOption) *Locator {
	return d.Locator(append([]LocatorOption{withID(id)}, opts...)...)
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
	var set []string
	if other.model != "" {
		set = append(set, "Model")
	}
	if other.ocrEngine != "" {
		set = append(set, "OCREngine")
	}
	if other.strategy != "" {
		set = append(set, "Strategy")
	}
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
		Role:   l.role,
		Name:   l.name,
		Text:   l.text,
		Exact:  l.exact,
		Id:     l.id,
		States: l.states,
		Query:  l.query,
		Nth:    l.nth,
		Within: l.within.toWire(),
		Has:    l.has.toWire(),
	}
	if l.androidClassName != "" {
		w.Platform = &locatorPlatformWire{Android: &androidLocatorWire{ClassName: l.androidClassName}}
	}
	return w
}
