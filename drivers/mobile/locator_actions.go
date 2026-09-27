package mobile

import (
	"encoding/json"
	"fmt"
	"time"
)

// Strategy names for WithStrategy / WithDefaultStrategy (the wire's
// LocatorStrategy).
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

// WithStrategy overrides the resolution strategy (StrategyAuto/Vision/
// Accessibility) for this locator call.
func WithStrategy(strategy string) CallOption {
	return func(c *callConfig) { c.strategy = strategy }
}

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

// resolveOverride returns callValue if the call set one, else driverDefault
// (which is "" if the driver has none either: a locator call omits a field
// entirely rather than forcing a client-side default onto it, unlike the
// vision calls' resolveEngine, so the server's own default applies).
func resolveOverride(callValue, driverDefault string) string {
	if callValue != "" {
		return callValue
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

// locatorCallConfig applies a locator call's options against
// defaultLocatorTimeout and resolves strategy/model/ocrEngine against the
// driver's own defaults, per resolveOverride. A timeout outside the
// protocol's range is an InvalidArgs error here rather than a rejected
// request on the device.
func (d *MobileDriver) locatorCallConfig(opts []CallOption) (callConfig, error) {
	cfg := applyCall(defaultLocatorTimeout, opts)
	if cfg.timeout < 0 || cfg.timeout > maxLocatorTimeout {
		return callConfig{}, &Error{
			Code:    CodeInvalidArgs,
			Message: fmt.Sprintf("locator timeout %s is outside [0, %s]", cfg.timeout, maxLocatorTimeout),
		}
	}
	cfg.strategy = resolveOverride(cfg.strategy, d.defaultStrategy)
	cfg.ocrEngine = resolveOverride(cfg.ocrEngine, d.defaultOCREngine)
	cfg.model = resolveOverride(cfg.model, d.defaultModel)
	return cfg, nil
}

// --- actions -----------------------------------------------------------

// Tap resolves the locator, waits until it is actionable (present and not
// moving), then taps its centre, all on the device in one round trip. Fails
// with ActionTimeout if the target never becomes actionable within the
// call's timeout.
func (l *Locator) Tap(opts ...CallOption) (LocatorResult, error) {
	cfg, err := l.driver.locatorCallConfig(opts)
	if err != nil {
		return LocatorResult{}, err
	}
	p := locatorTapParams{
		Locator:   l.toWire(),
		Strategy:  cfg.strategy,
		TimeoutMs: timeoutMsField(cfg.timeout),
		OcrEngine: cfg.ocrEngine,
		Model:     cfg.model,
	}
	raw, err := l.driver.locatorCall(methodLocatorTap, p, cfg.timeout)
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
	cfg, err := l.driver.locatorCallConfig(opts)
	if err != nil {
		return LocatorResult{}, err
	}
	p := locatorFillParams{
		Locator:   l.toWire(),
		Text:      text,
		Strategy:  cfg.strategy,
		TimeoutMs: timeoutMsField(cfg.timeout),
		OcrEngine: cfg.ocrEngine,
		Model:     cfg.model,
	}
	raw, err := l.driver.locatorCall(methodLocatorFill, p, cfg.timeout)
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
// instead, with no locator involved.
func (l *Locator) Press(key string, opts ...CallOption) (LocatorResult, error) {
	return l.driver.press(key, l.toWire(), opts)
}

// Press sends a named key (see the Key* constants) to whatever currently has
// focus. ResolvedBy/Bounds on the result are empty: nothing was located.
func (d *MobileDriver) Press(key string, opts ...CallOption) (LocatorResult, error) {
	return d.press(key, nil, opts)
}

func (d *MobileDriver) press(key string, loc *locatorWire, opts []CallOption) (LocatorResult, error) {
	cfg, err := d.locatorCallConfig(opts)
	if err != nil {
		return LocatorResult{}, err
	}
	p := locatorPressParams{
		Locator:   loc,
		Key:       key,
		Strategy:  cfg.strategy,
		TimeoutMs: timeoutMsField(cfg.timeout),
		OcrEngine: cfg.ocrEngine,
		Model:     cfg.model,
	}
	raw, err := d.locatorCall(methodLocatorPress, p, cfg.timeout)
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
	cfg, err := l.driver.locatorCallConfig(opts)
	if err != nil {
		return nil, err
	}
	p := locatorWaitForParams{
		Locator:   l.toWire(),
		State:     state,
		Strategy:  cfg.strategy,
		TimeoutMs: timeoutMsField(cfg.timeout),
		OcrEngine: cfg.ocrEngine,
		Model:     cfg.model,
	}
	raw, err := l.driver.locatorCall(methodLocatorWaitFor, p, cfg.timeout)
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
	cfg, err := l.driver.locatorCallConfig(opts)
	if err != nil {
		return nil, err
	}
	p := locatorQueryParams{
		Locator:   l.toWire(),
		Strategy:  cfg.strategy,
		TimeoutMs: timeoutMsField(cfg.timeout),
		OcrEngine: cfg.ocrEngine,
		Model:     cfg.model,
	}
	return l.driver.locatorCall(method, p, cfg.timeout)
}

// Count counts the targets matching the locator on the current screen, zero
// included. Never waits: use WaitFor first if the target may not be on
// screen yet.
func (l *Locator) Count(opts ...CallOption) (int, error) {
	cfg, err := l.driver.locatorCallConfig(opts)
	if err != nil {
		return 0, err
	}
	p := locatorCountParams{
		Locator:   l.toWire(),
		Strategy:  cfg.strategy,
		OcrEngine: cfg.ocrEngine,
		Model:     cfg.model,
	}
	// Count sends no device-side budget, so the in-flight-inference margin
	// that pads every waiting call does not apply to a timeout the caller
	// gave it: that is the whole deadline.
	var raw json.RawMessage
	if cfg.timeoutSet {
		raw, err = l.driver.call(methodLocatorCount, p, cfg.timeout)
	} else {
		raw, err = l.driver.locatorCall(methodLocatorCount, p, cfg.timeout)
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
