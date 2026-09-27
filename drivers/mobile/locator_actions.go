package mobile

import (
	"encoding/json"
	"fmt"
	"time"
)

// Strategy names for the Strategy locator option / WithDefaultStrategy (the
// wire's LocatorStrategy).
const (
	// StrategyAuto uses the accessibility tree when the session has one and
	// the screen otherwise. The server-side default when strategy is omitted.
	StrategyAuto = "auto"
	// StrategyVision always resolves from the screen (OCR/VLM), even on a
	// session that has an accessibility tree.
	StrategyVision = "vision"
	// StrategyAccessibility always resolves from the accessibility tree;
	// answers StrategyUnavailable on a session without one.
	StrategyAccessibility = "accessibility"
)

// States for Locator.WaitFor (the wire's LocatorWaitForParams.state).
const (
	StateVisible = "visible"
	StateHidden  = "hidden"
	// StateEnabled needs the accessibility tree: vision cannot tell enabled
	// from disabled.
	StateEnabled = "enabled"
)

// Locator call budgets. defaultLocatorTimeout is the device-side auto-wait
// budget (the wire's timeoutMs) when the call doesn't override it, matching
// the contract's own default. locatorCallMargin pads the transport-level
// deadline past that budget: an inference already running when timeoutMs
// ends is allowed to finish server-side, so the round trip's own deadline
// must outlive it rather than racing it.
const (
	defaultLocatorTimeout = 5 * time.Second
	locatorCallMargin     = 15 * time.Second
)

// LocatorResult is the outcome of a locator action or query: where the
// target was when the device acted on it, how it was resolved, and how long
// resolving took (auto-wait included).
type LocatorResult struct {
	// ResolvedBy is how the target was found: "a11y", "ocr", or "vlm". Empty
	// when the call carried no locator (MobileDriver.Press) or the result is
	// from a WaitFor("hidden") (nothing to report).
	ResolvedBy string `json:"resolved_by,omitempty"`
	// Bounds is where the target was when the device acted, in the pixel
	// space the Touch methods address. Zero when ResolvedBy is empty.
	Bounds BBox `json:"bounds"`
	// TookMs is the end-to-end time on the device, auto-wait included.
	TookMs int64 `json:"took_ms"`
	// ModelName is the model that resolved the target; set only when
	// ResolvedBy is "vlm".
	ModelName string `json:"model_name,omitempty"`
}

// wireLocatorResult decodes any of LocatorResult / LocatorPressResult /
// LocatorWaitForResult / LocatorQueryParams's LocatorResult: they share this
// shape, differing only in which fields the contract requires versus leaves
// optional, which JSON decoding doesn't need to distinguish.
type wireLocatorResult struct {
	ResolvedBy string   `json:"resolvedBy,omitempty"`
	Bounds     wireBBox `json:"bounds"`
	TookMs     int64    `json:"tookMs"`
	ModelName  string   `json:"modelName,omitempty"`
}

type wireLocatorTextResult struct {
	Text       string   `json:"text"`
	ResolvedBy string   `json:"resolvedBy"`
	Bounds     wireBBox `json:"bounds"`
	TookMs     int64    `json:"tookMs"`
	ModelName  string   `json:"modelName,omitempty"`
}

type wireLocatorCountResult struct {
	Count      int    `json:"count"`
	ResolvedBy string `json:"resolvedBy"`
	TookMs     int64  `json:"tookMs"`
	ModelName  string `json:"modelName,omitempty"`
}

func locatorResultFromWire(w wireLocatorResult) LocatorResult {
	return LocatorResult{
		ResolvedBy: w.ResolvedBy,
		Bounds:     bboxFromWire(w.Bounds),
		TookMs:     w.TookMs,
		ModelName:  w.ModelName,
	}
}

// resolveOverride returns locatorValue if the locator set one, else
// driverDefault (which is "" if the driver has none either: a locator call
// omits a field entirely rather than forcing a client-side default onto it,
// unlike the vision call's resolveEngine, so the server's own default
// applies).
func resolveOverride(locatorValue, driverDefault string) string {
	if locatorValue != "" {
		return locatorValue
	}
	return driverDefault
}

// timeoutMsField converts a call's resolved timeout to the wire's
// milliseconds integer.
func timeoutMsField(d time.Duration) int { return int(d / time.Millisecond) }

// locatorCall performs one Locator.* round trip. deviceTimeout is the budget
// already sent as the wire's timeoutMs (or, for Locator.count, the caller's
// WithTimeout with no wire effect); the transport deadline is padded past it
// by locatorCallMargin so the call's own deadline never races the device-side
// wait it is bounding.
func (d *MobileDriver) locatorCall(method string, params any, deviceTimeout time.Duration) (json.RawMessage, error) {
	return d.call(method, params, deviceTimeout+locatorCallMargin)
}

// maxLocatorTimeout is the protocol's ceiling on a locator call's device-side
// budget (the wire's timeoutMs maximum).
const maxLocatorTimeout = 60 * time.Second

// locatorParams is what one locator round trip resolves before it builds its
// wire params: the strategy/model/ocrEngine to send (precedence: the
// locator's own Model/OCREngine/Strategy option, else the driver's
// WithDefault*, else omitted so the server's own default applies) and the
// call's timeout budget.
type locatorParams struct {
	strategy, model, ocrEngine string
	timeout                    time.Duration
	// timeoutSet records that the call passed WithTimeout, so Locator.Count
	// (which sends no device-side budget) can treat it as the whole
	// deadline instead of padding it with locatorCallMargin.
	timeoutSet bool
}

// validateLocatorTimeout rejects a timeout outside the protocol's range
// locally, as an InvalidArgs error, rather than sending a request the device
// would reject anyway.
func validateLocatorTimeout(cfg callConfig) error {
	if cfg.timeout < 0 || cfg.timeout > maxLocatorTimeout {
		return &Error{
			Code:    CodeInvalidArgs,
			Message: fmt.Sprintf("locator timeout %s is outside [0, %s]", cfg.timeout, maxLocatorTimeout),
		}
	}
	return nil
}

// resolveParams applies opts (a locator action or query's only meaningful
// CallOption is WithTimeout: strategy/model/ocrEngine are set on the locator
// itself, not per call) against defaultLocatorTimeout, then resolves the
// locator's own strategy/model/ocrEngine against the driver's defaults, per
// resolveOverride.
func (l *Locator) resolveParams(opts []CallOption) (locatorParams, error) {
	cfg := applyCall(defaultLocatorTimeout, opts)
	if err := validateLocatorTimeout(cfg); err != nil {
		return locatorParams{}, err
	}
	d := l.driver
	return locatorParams{
		strategy:   resolveOverride(l.strategy, d.defaultStrategy),
		model:      resolveOverride(l.model, d.defaultModel),
		ocrEngine:  resolveOverride(l.ocrEngine, d.defaultOCREngine),
		timeout:    cfg.timeout,
		timeoutSet: cfg.timeoutSet,
	}, nil
}

// locatorCallTimeout resolves just the timeout budget for a locator call with
// no locator to resolve against (MobileDriver.Press): it takes no resolution
// options, so its wire strategy/model/ocrEngine stay empty.
func locatorCallTimeout(opts []CallOption) (locatorParams, error) {
	cfg := applyCall(defaultLocatorTimeout, opts)
	if err := validateLocatorTimeout(cfg); err != nil {
		return locatorParams{}, err
	}
	return locatorParams{timeout: cfg.timeout, timeoutSet: cfg.timeoutSet}, nil
}

// --- actions -----------------------------------------------------------

// Tap resolves the locator, waits until it is actionable (present and not
// moving), then taps its centre, all on the device in one round trip. Fails
// with ActionTimeout if the target never becomes actionable within the
// call's timeout.
func (l *Locator) Tap(opts ...CallOption) (LocatorResult, error) {
	p, err := l.resolveParams(opts)
	if err != nil {
		return LocatorResult{}, err
	}
	wp := locatorTapParams{
		Locator:   l.toWire(),
		Strategy:  p.strategy,
		TimeoutMs: timeoutMsField(p.timeout),
		OcrEngine: p.ocrEngine,
		Model:     p.model,
	}
	raw, err := l.driver.locatorCall(methodLocatorTap, wp, p.timeout)
	if err != nil {
		return LocatorResult{}, err
	}
	var w wireLocatorResult
	if err := unmarshalResult(raw, &w); err != nil {
		return LocatorResult{}, err
	}
	return locatorResultFromWire(w), nil
}

// Fill resolves and waits as Tap, focuses the target by tapping it, then
// types text into it.
func (l *Locator) Fill(text string, opts ...CallOption) (LocatorResult, error) {
	p, err := l.resolveParams(opts)
	if err != nil {
		return LocatorResult{}, err
	}
	wp := locatorFillParams{
		Locator:   l.toWire(),
		Text:      text,
		Strategy:  p.strategy,
		TimeoutMs: timeoutMsField(p.timeout),
		OcrEngine: p.ocrEngine,
		Model:     p.model,
	}
	raw, err := l.driver.locatorCall(methodLocatorFill, wp, p.timeout)
	if err != nil {
		return LocatorResult{}, err
	}
	var w wireLocatorResult
	if err := unmarshalResult(raw, &w); err != nil {
		return LocatorResult{}, err
	}
	return locatorResultFromWire(w), nil
}

// Press resolves and waits as Tap, focuses the target, then presses key.
// Use MobileDriver.Press to send the key to whatever currently has focus
// instead, with no locator involved (and so no resolution options).
func (l *Locator) Press(key string, opts ...CallOption) (LocatorResult, error) {
	p, err := l.resolveParams(opts)
	if err != nil {
		return LocatorResult{}, err
	}
	return l.driver.press(key, l.toWire(), p)
}

// Press sends a named key (see the Key* constants) to whatever currently has
// focus. ResolvedBy/Bounds on the result are empty: nothing was located, and
// with no locator this call takes no resolution options, only WithTimeout.
func (d *MobileDriver) Press(key string, opts ...CallOption) (LocatorResult, error) {
	p, err := locatorCallTimeout(opts)
	if err != nil {
		return LocatorResult{}, err
	}
	return d.press(key, nil, p)
}

func (d *MobileDriver) press(key string, loc *locatorWire, p locatorParams) (LocatorResult, error) {
	wp := locatorPressParams{
		Locator:   loc,
		Key:       key,
		Strategy:  p.strategy,
		TimeoutMs: timeoutMsField(p.timeout),
		OcrEngine: p.ocrEngine,
		Model:     p.model,
	}
	raw, err := d.locatorCall(methodLocatorPress, wp, p.timeout)
	if err != nil {
		return LocatorResult{}, err
	}
	var w wireLocatorResult
	if err := unmarshalResult(raw, &w); err != nil {
		return LocatorResult{}, err
	}
	return locatorResultFromWire(w), nil
}

// WaitFor blocks on the device until the locator reaches state (StateVisible
// by default), re-resolving against each fresh frame; one call replaces a
// client-side polling loop. Returns nil for StateHidden (nothing was
// located, so there is nothing to report); for the other states it returns
// where the target was found.
func (l *Locator) WaitFor(state string, opts ...CallOption) (*LocatorResult, error) {
	p, err := l.resolveParams(opts)
	if err != nil {
		return nil, err
	}
	wp := locatorWaitForParams{
		Locator:   l.toWire(),
		State:     state,
		Strategy:  p.strategy,
		TimeoutMs: timeoutMsField(p.timeout),
		OcrEngine: p.ocrEngine,
		Model:     p.model,
	}
	raw, err := l.driver.locatorCall(methodLocatorWaitFor, wp, p.timeout)
	if err != nil {
		return nil, err
	}
	if state == StateHidden {
		return nil, nil
	}
	var w wireLocatorResult
	if err := unmarshalResult(raw, &w); err != nil {
		return nil, err
	}
	res := locatorResultFromWire(w)
	return &res, nil
}

// BoundingBox waits until the locator resolves, then returns the target's
// bounds.
func (l *Locator) BoundingBox(opts ...CallOption) (LocatorResult, error) {
	raw, err := l.runQuery(methodLocatorBoundingBox, opts)
	if err != nil {
		return LocatorResult{}, err
	}
	var w wireLocatorResult
	if err := unmarshalResult(raw, &w); err != nil {
		return LocatorResult{}, err
	}
	return locatorResultFromWire(w), nil
}

// Text waits until the locator resolves, then returns the target's text.
func (l *Locator) Text(opts ...CallOption) (string, error) {
	raw, err := l.runQuery(methodLocatorText, opts)
	if err != nil {
		return "", err
	}
	var w wireLocatorTextResult
	if err := unmarshalResult(raw, &w); err != nil {
		return "", err
	}
	return w.Text, nil
}

// runQuery runs the shared LocatorQueryParams shape (BoundingBox and Text both
// use it).
func (l *Locator) runQuery(method string, opts []CallOption) (json.RawMessage, error) {
	p, err := l.resolveParams(opts)
	if err != nil {
		return nil, err
	}
	wp := locatorQueryParams{
		Locator:   l.toWire(),
		Strategy:  p.strategy,
		TimeoutMs: timeoutMsField(p.timeout),
		OcrEngine: p.ocrEngine,
		Model:     p.model,
	}
	return l.driver.locatorCall(method, wp, p.timeout)
}

// Count counts the targets matching the locator on the current screen, zero
// included. Never waits: use WaitFor first if the target may not be on
// screen yet. Under vision resolution, Count needs a plain text locator: one
// that also carries Query, Within or Has answers CodeInvalidArgs, since
// counting needs every independent match and a vision-model call only
// resolves a single target per prompt.
func (l *Locator) Count(opts ...CallOption) (int, error) {
	p, err := l.resolveParams(opts)
	if err != nil {
		return 0, err
	}
	wp := locatorCountParams{
		Locator:   l.toWire(),
		Strategy:  p.strategy,
		OcrEngine: p.ocrEngine,
		Model:     p.model,
	}
	// Count sends no device-side budget, so the in-flight-inference margin
	// that pads every waiting call does not apply to a timeout the caller
	// gave it: that is the whole deadline.
	var raw json.RawMessage
	if p.timeoutSet {
		raw, err = l.driver.call(methodLocatorCount, wp, p.timeout)
	} else {
		raw, err = l.driver.locatorCall(methodLocatorCount, wp, p.timeout)
	}
	if err != nil {
		return 0, err
	}
	var w wireLocatorCountResult
	if err := unmarshalResult(raw, &w); err != nil {
		return 0, err
	}
	return w.Count, nil
}
