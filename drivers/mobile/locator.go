package mobile

// Locator is a lazy, immutable description of an on-screen target (the
// Locator tier: DCP's Locator.tap/fill/press/waitFor/boundingBox/text/count).
// Building one sends nothing; only an action or query method on it (Tap,
// Fill, Press, WaitFor, BoundingBox, Text, Count) reaches the device, and
// resolves the locator, auto-waits until actionable, and acts, all in one
// round trip. Nth, First, Within, Has and Filter each return a new Locator
// rather than mutating the receiver.
type Locator struct {
	driver *MobileDriver

	role, name, text, id, query, androidClassName string
	exact                                         bool
	states                                        []string
	nth                                           *int
	within, has                                   *Locator
}

// LocatorOption sets one predicate on a Locator at construction (via
// MobileDriver.Locator) or refinement (via Locator.Filter). Predicates
// combine as AND; a literal selector that matches nothing fails after the
// call's auto-wait.
type LocatorOption func(*Locator)

// Text matches visible text: substring and case-insensitive unless Exact is
// also given. Against the accessibility tree when the session has one, OCR
// otherwise (GetByText is this option alone).
func Text(text string) LocatorOption { return func(l *Locator) { l.text = text } }

// Exact requires an exact match on Text rather than a case-insensitive
// substring.
func Exact() LocatorOption { return func(l *Locator) { l.exact = true } }

// Role matches the accessibility role, e.g. "button", "textbox". Needs the
// accessibility tree: on a session without one this answers
// StrategyUnavailable regardless of strategy (today's phones have none).
func Role(role string) LocatorOption { return func(l *Locator) { l.role = role } }

// Name matches the accessible name. Needs the accessibility tree; see Role.
func Name(name string) LocatorOption { return func(l *Locator) { l.name = name } }

// ID matches a developer-assigned id, e.g. an Android resource id
// "com.example.app:id/save". Needs the accessibility tree; see Role.
func ID(id string) LocatorOption { return func(l *Locator) { l.id = id } }

// States requires the given accessibility states, e.g. "checked". Needs the
// accessibility tree; see Role.
func States(states ...string) LocatorOption {
	return func(l *Locator) { l.states = append([]string(nil), states...) }
}

// Query is a natural-language description of the target: ranked over the
// accessibility tree by a model when the session has one, read from the
// screen by the VLM otherwise.
func Query(query string) LocatorOption { return func(l *Locator) { l.query = query } }

// AndroidClassName matches the Android view class, e.g.
// "android.widget.Button". Native matching; not portable across device
// classes.
func AndroidClassName(className string) LocatorOption {
	return func(l *Locator) { l.androidClassName = className }
}

// Locator builds a general-purpose locator from options. GetByText, GetByRole
// and GetByID are shorthand for the common single-predicate cases.
func (d *MobileDriver) Locator(opts ...LocatorOption) *Locator {
	l := &Locator{driver: d}
	for _, o := range opts {
		o(l)
	}
	return l
}

// GetByText builds a locator matching visible text (Text(text) plus any
// further opts, e.g. Exact()). On today's phones (no accessibility tree) this
// resolves by OCR.
func (d *MobileDriver) GetByText(text string, opts ...LocatorOption) *Locator {
	return d.Locator(append([]LocatorOption{Text(text)}, opts...)...)
}

// GetByRole builds a locator matching an accessibility role. Today's phones
// have no accessibility tree, so acting on it answers StrategyUnavailable by
// design until a phone advertises one.
func (d *MobileDriver) GetByRole(role string, opts ...LocatorOption) *Locator {
	return d.Locator(append([]LocatorOption{Role(role)}, opts...)...)
}

// GetByID builds a locator matching a developer-assigned id. Today's phones
// have no accessibility tree, so acting on it answers StrategyUnavailable by
// design until a phone advertises one.
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
func (l *Locator) Within(other *Locator) *Locator {
	out := l.clone()
	out.within = other
	return out
}

// Has returns a new locator that only matches elements containing a match of
// other.
func (l *Locator) Has(other *Locator) *Locator {
	out := l.clone()
	out.has = other
	return out
}

// Filter returns a new locator with opts applied on top of the receiver's own
// predicates (AND'ed together).
func (l *Locator) Filter(opts ...LocatorOption) *Locator {
	out := l.clone()
	for _, o := range opts {
		o(out)
	}
	return out
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
