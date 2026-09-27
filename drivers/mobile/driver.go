package mobile

import (
	"context"
	"encoding/json"
	"time"
)

// Default deadlines. Input/screenshot calls use defaultCallTimeout; the
// Observe vision call defaults to visionTimeout; locator calls have their own
// budget (see defaultLocatorTimeout in locator_actions.go). All take a
// per-call WithTimeout override.
const (
	defaultCallTimeout = 30 * time.Second
	visionTimeout      = 10 * time.Second
	defaultOpenTimeout = 10 * time.Second
)

// MobileDriver drives a paired phone through a transport (the DCP control
// WebSocket). It is the Go twin of platform-python's MobileDriver: an
// ergonomic observe/locator/tap/type API over literal CDP method frames.
//
// The DefaultOCREngine / DefaultModel / DefaultStrategy session defaults feed
// the vision and locator calls: any call that takes a per-call WithOCREngine /
// WithModel / WithStrategy uses the driver default when the call omits one,
// so a script sets the premium engine (or a specific VLM, or a resolution
// strategy) once instead of on every call. A per-call option always wins.
// Observe (the vision call) always sends an explicit OCR engine, falling
// back to "free" when neither is set; a locator call, with neither set,
// omits the field entirely so the server's own default applies.
type MobileDriver struct {
	tp transport

	defaultOCREngine string
	defaultModel     string
	defaultStrategy  string
	openTimeout      time.Duration
}

// Option configures a MobileDriver at construction.
type Option func(*MobileDriver)

// WithDefaultOCREngine sets the session-wide OCR engine for vision and
// locator calls.
func WithDefaultOCREngine(engine string) Option {
	return func(d *MobileDriver) { d.defaultOCREngine = engine }
}

// WithDefaultModel sets the session-wide VLM for locator calls that fall back
// to the vision model.
func WithDefaultModel(model string) Option {
	return func(d *MobileDriver) { d.defaultModel = model }
}

// WithDefaultStrategy sets the session-wide locator resolution strategy
// (StrategyAuto/Vision/Accessibility).
func WithDefaultStrategy(strategy string) Option {
	return func(d *MobileDriver) { d.defaultStrategy = strategy }
}

// WithOpenTimeout sets how long the first call waits to open the control socket.
func WithOpenTimeout(d time.Duration) Option {
	return func(m *MobileDriver) { m.openTimeout = d }
}

// ConnectRemote builds a driver for a remotely-allocated phone over its DCP
// control URL. controlURL is the value returned by the allocate call — a wss://
// URL with the scoped, allocation-bound control token already embedded. The
// socket opens lazily on the first call.
func ConnectRemote(controlURL string, opts ...Option) *MobileDriver {
	d := &MobileDriver{openTimeout: defaultOpenTimeout}
	for _, o := range opts {
		o(d)
	}
	d.tp = newRemoteTransport(controlURL, d.openTimeout)
	return d
}

// newDriver wraps a transport directly (test seam).
func newDriver(tp transport, opts ...Option) *MobileDriver {
	d := &MobileDriver{openTimeout: defaultOpenTimeout}
	for _, o := range opts {
		o(d)
	}
	d.tp = tp
	return d
}

// Close releases the underlying transport. The next call reconnects.
func (d *MobileDriver) Close() error { return d.tp.close() }

// --- per-call options -------------------------------------------------------

type callConfig struct {
	ocrEngine string
	model     string
	timeout   time.Duration
	strategy  string
	// timeoutSet records that the caller passed WithTimeout, so a call with
	// no device-side budget (Locator.count) can treat it as the whole
	// deadline.
	timeoutSet bool
}

// CallOption tunes a single vision or locator call.
type CallOption func(*callConfig)

// WithOCREngine overrides the OCR engine for this call.
func WithOCREngine(engine string) CallOption {
	return func(c *callConfig) { c.ocrEngine = engine }
}

// WithModel overrides the model for this call: the VLM under vision
// resolution, or the model a locator's natural-language Query is ranked with
// over the accessibility tree.
func WithModel(model string) CallOption {
	return func(c *callConfig) { c.model = model }
}

// WithTimeout overrides the deadline for this call.
func WithTimeout(d time.Duration) CallOption {
	return func(c *callConfig) {
		c.timeout = d
		c.timeoutSet = d > 0
	}
}

func (d *MobileDriver) resolveEngine(c callConfig) string {
	if c.ocrEngine != "" {
		return c.ocrEngine
	}
	if d.defaultOCREngine != "" {
		return d.defaultOCREngine
	}
	return "free"
}

func (d *MobileDriver) call(method string, params any, timeout time.Duration) (json.RawMessage, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	return d.tp.call(ctx, method, params)
}

// --- vision -----------------------------------------------------------------

// Observe captures the current frame and returns a typed Screen.
func (d *MobileDriver) Observe(opts ...CallOption) (*Screen, error) {
	cfg := applyCall(visionTimeout, opts)
	raw, err := d.call(methodScreenObserve, observeParams{OcrEngine: d.resolveEngine(cfg)}, cfg.timeout)
	if err != nil {
		return nil, err
	}
	return screenFromWire(raw)
}

// --- input ------------------------------------------------------------------

// Tap taps once at coords.
func (d *MobileDriver) Tap(c Coords) error { return d.tapXY(c.X, c.Y) }

// LongPress presses and holds at coords for durationMs.
func (d *MobileDriver) LongPress(c Coords, durationMs int) error {
	return d.longPressXY(c.X, c.Y, durationMs)
}

// Swipe swipes from start to end over durationMs.
func (d *MobileDriver) Swipe(start, end Coords, durationMs int) error {
	return d.swipeXY(start.X, start.Y, end.X, end.Y, durationMs)
}

// TypeText types a string of US-layout-typable text.
func (d *MobileDriver) TypeText(text string) error { return d.typeText(text) }

// KeyPress presses a named key (see the Key* constants), e.g. KeyEnter.
func (d *MobileDriver) KeyPress(key string) error {
	_, err := d.call(methodKeyboardKeyPress, keyPressParams{Key: key}, defaultCallTimeout)
	return err
}

// Screenshot captures the current frame as PNG-encoded bytes.
func (d *MobileDriver) Screenshot() ([]byte, error) {
	raw, err := d.call(methodScreenScreenshot, nil, defaultCallTimeout)
	if err != nil {
		return nil, err
	}
	var wire wireScreenshot
	if err := unmarshalResult(raw, &wire); err != nil {
		return nil, err
	}
	if len(wire.PngBase64) == 0 {
		return nil, &Error{Code: CodeInternal, Message: "screenshot returned no image"}
	}
	return wire.PngBase64, nil
}

func (d *MobileDriver) tapXY(x, y int) error {
	_, err := d.call(methodTouchTap, tapParams{X: x, Y: y}, defaultCallTimeout)
	return err
}

func (d *MobileDriver) longPressXY(x, y, durationMs int) error {
	_, err := d.call(methodTouchLongPress, longPressParams{X: x, Y: y, DurationMs: durationMs}, defaultCallTimeout)
	return err
}

func (d *MobileDriver) swipeXY(x1, y1, x2, y2, durationMs int) error {
	_, err := d.call(methodTouchSwipe, swipeParams{X1: x1, Y1: y1, X2: x2, Y2: y2, DurationMs: durationMs}, defaultCallTimeout)
	return err
}

func (d *MobileDriver) typeText(text string) error {
	_, err := d.call(methodKeyboardTypeText, typeTextParams{Text: text}, defaultCallTimeout)
	return err
}

func applyCall(defaultTimeout time.Duration, opts []CallOption) callConfig {
	cfg := callConfig{timeout: defaultTimeout}
	for _, o := range opts {
		o(&cfg)
	}
	if cfg.timeout <= 0 {
		cfg.timeout = defaultTimeout
	}
	return cfg
}

// --- wire result frames (snake_case). The input-param frames live in the
// generated wire_gen.go; these result shapes feed the hand-written converters
// below into the ergonomic Screen/Element types, so they stay hand-written. ---

type wireBBox struct {
	X      int `json:"x"`
	Y      int `json:"y"`
	Width  int `json:"width"`
	Height int `json:"height"`
}

type wireText struct {
	Text       string   `json:"text"`
	BBox       wireBBox `json:"bbox"`
	Confidence float64  `json:"confidence"`
}

type wireIcon struct {
	BBox       wireBBox `json:"bbox"`
	Confidence float64  `json:"confidence"`
}

type wireObserve struct {
	Texts      []wireText `json:"texts"`
	Icons      []wireIcon `json:"icons"`
	Hash       string     `json:"hash"`
	Width      int        `json:"width"`
	Height     int        `json:"height"`
	CapturedAt int64      `json:"captured_at"` // epoch-millis
}

type wireScreenshot struct {
	// Go's json decodes a base64 string straight into []byte.
	PngBase64 []byte `json:"png_base64"`
}

func unmarshalResult(raw json.RawMessage, into any) error {
	if len(raw) == 0 {
		return nil
	}
	if err := json.Unmarshal(raw, into); err != nil {
		return &Error{Code: CodeInternal, Message: "malformed result: " + err.Error()}
	}
	return nil
}

func bboxFromWire(w wireBBox) BBox {
	return BBox{X: w.X, Y: w.Y, Width: w.Width, Height: w.Height}
}

func elementFromText(w wireText) Element {
	b := bboxFromWire(w.BBox)
	return Element{BBox: b, Center: b.Center(), Confidence: w.Confidence, Text: w.Text, Source: SourceOCR}
}

func iconFromWire(w wireIcon) IconBox {
	b := bboxFromWire(w.BBox)
	return IconBox{BBox: b, Center: b.Center(), Confidence: w.Confidence}
}

func screenFromWire(raw json.RawMessage) (*Screen, error) {
	var w wireObserve
	if err := unmarshalResult(raw, &w); err != nil {
		return nil, err
	}
	s := &Screen{
		Hash:       w.Hash,
		Width:      w.Width,
		Height:     w.Height,
		CapturedAt: time.UnixMilli(w.CapturedAt).UTC(),
	}
	for _, t := range w.Texts {
		s.Texts = append(s.Texts, elementFromText(t))
	}
	for _, ic := range w.Icons {
		s.Icons = append(s.Icons, iconFromWire(ic))
	}
	return s, nil
}
