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

	text, query string
	exact       bool
	nth         *int
	within, has *Locator

	// model and ocrEngine are this locator's own resolution options
	// (Model/OCREngine). An action or query on the locator resolves with
	// these, falling back to the driver's WithDefaultModel /
	// WithDefaultOCREngine, and omitting the wire field entirely if neither
	// is set. Within/Has only ever read the selector fields off a nested
	// locator: these two never travel with it (one call resolves the whole
	// locator, governed by the outer locator's own options).
	model, ocrEngine string

	// buildErr is set by Within/Has when the locator passed in as the scope
	// carries its own Model/OCREngine (or already carries a buildErr of its
	// own): those options would otherwise be silently dropped, so instead
	// every action and query on this locator (and on anything refined from
	// it, since clone keeps the field) fails locally with this error rather
	// than sending a request that ignored what the caller asked for.
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

// Query is a natural-language description of the target, read from the
// screen by the VLM.
func Query(query string) LocatorOption { return func(l *Locator) { l.query = query } }

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

// Locator builds a general-purpose locator from options. GetByText is
// shorthand for the common single-predicate case.
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
// Model/OCREngine govern how the whole thing resolves. If other itself
// carries one of those options, set on itself rather than inherited from a
// driver default, that would silently ignore what the caller asked for, so
// this records a build error on the returned locator instead: every action
// or query on it (and on anything further refined from it) fails locally
// with a CodeInvalidArgs *Error, without sending anything.
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
// when it also carries its own Model/OCREngine.
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
	if len(set) == 0 {
		return nil
	}
	return &Error{
		Code: CodeInvalidArgs,
		Message: fmt.Sprintf(
			"%s: the inner locator sets %s; set Model and OCREngine on the outer locator instead",
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

// toWire converts a possibly-nil Locator into the generated wire shape (nil
// in, nil out; used directly for the optional locator on Locator.press and
// for the recursive within/has). The wire shape still carries Role, Name,
// Id, States and Platform (the vendored contract; the device still rejects
// them clearly), but the public API has nothing that sets them, so toWire
// never populates them.
func (l *Locator) toWire() *locatorWire {
	if l == nil {
		return nil
	}
	return &locatorWire{
		Text:   l.text,
		Exact:  l.exact,
		Query:  l.query,
		Nth:    l.nth,
		Within: l.within.toWire(),
		Has:    l.has.toWire(),
	}
}
