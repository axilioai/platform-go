package mobile

import (
	"context"
	"encoding/json"
	"time"
)

// Default deadlines. Input/screenshot calls use defaultCallTimeout; the vision
// calls (observe/find) default to visionTimeout and take a per-call override.
const (
	defaultCallTimeout = 30 * time.Second
	visionTimeout      = 10 * time.Second
	defaultOpenTimeout = 10 * time.Second
	defaultPollEvery   = 300 * time.Millisecond
)

// MobileDriver drives a paired phone through a transport (the DCP control
// WebSocket). It is the Go twin of platform-python's MobileDriver: an ergonomic
// observe/find/tap/type API over literal CDP method frames.
//
// The DefaultOCREngine / DefaultModel session defaults feed the vision calls:
// any call that takes a per-call WithOCREngine / WithModel uses the driver
// default when the call omits one, so a script sets the premium engine (or a
// specific VLM) once instead of on every call. A per-call option always wins;
// with neither, the engine falls back to "free" and the model to the
// server-side default.
type MobileDriver struct {
	tp transport

	defaultOCREngine string
	defaultModel     string
	openTimeout      time.Duration
}

// Option configures a MobileDriver at construction.
type Option func(*MobileDriver)

// WithDefaultOCREngine sets the session-wide OCR engine for vision calls.
func WithDefaultOCREngine(engine string) Option {
	return func(d *MobileDriver) { d.defaultOCREngine = engine }
}

// WithDefaultModel sets the session-wide VLM for Find.
func WithDefaultModel(model string) Option {
	return func(d *MobileDriver) { d.defaultModel = model }
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
}

// CallOption tunes a single vision call.
type CallOption func(*callConfig)

// WithOCREngine overrides the OCR engine for this call.
func WithOCREngine(engine string) CallOption {
	return func(c *callConfig) { c.ocrEngine = engine }
}

// WithModel overrides the VLM for this Find call.
func WithModel(model string) CallOption {
	return func(c *callConfig) { c.model = model }
}

// WithTimeout overrides the deadline for this call.
func WithTimeout(d time.Duration) CallOption {
	return func(c *callConfig) { c.timeout = d }
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
	return d.screenFromWire(raw)
}

// FindText returns the first OCR element matching text (one Observe per call),
// or nil if none. With exact=false the match is case-insensitive substring.
func (d *MobileDriver) FindText(text string, exact bool, opts ...CallOption) (*Element, error) {
	screen, err := d.Observe(opts...)
	if err != nil {
		return nil, err
	}
	return screen.FindText(text, exact), nil
}

// FindAllText returns every OCR element matching the criteria (one Observe per
// call). contains and pattern are mutually exclusive.
func (d *MobileDriver) FindAllText(contains, pattern string, opts ...CallOption) ([]Element, error) {
	screen, err := d.Observe(opts...)
	if err != nil {
		return nil, err
	}
	return screen.FindAllText(contains, pattern)
}

// Find locates an element semantically via the vision model (through the
// on-device agent). Returns an *Error with CodeElementNotFound if nothing
// matches.
func (d *MobileDriver) Find(query string, opts ...CallOption) (*Element, error) {
	cfg := applyCall(visionTimeout, opts)
	p := findParams{Query: query, OcrEngine: d.resolveEngine(cfg)}
	if cfg.model != "" {
		p.Model = cfg.model
	} else if d.defaultModel != "" {
		p.Model = d.defaultModel
	}
	raw, err := d.call(methodScreenFind, p, cfg.timeout)
	if err != nil {
		return nil, err
	}
	var wire wireFind
	if err := unmarshalResult(raw, &wire); err != nil {
		return nil, err
	}
	if wire.Found == nil {
		return nil, &Error{Code: CodeElementNotFound, Message: "no element on screen matched query: " + query}
	}
	return d.elementFromFound(wire.Found), nil
}

// --- waits ------------------------------------------------------------------

// WaitForText polls FindText until the target appears or timeout elapses.
func (d *MobileDriver) WaitForText(text string, timeout time.Duration, exact bool, opts ...CallOption) (*Element, error) {
	deadline := time.Now().Add(timeout)
	for {
		el, err := d.FindText(text, exact, opts...)
		if err != nil {
			return nil, err
		}
		if el != nil {
			return el, nil
		}
		if time.Now().After(deadline) {
			return nil, &Error{Code: CodeTimeout, Message: "text not found within deadline: " + text, Retryable: true}
		}
		time.Sleep(defaultPollEvery)
	}
}

// WaitUntilGone polls until text disappears or timeout elapses.
func (d *MobileDriver) WaitUntilGone(text string, timeout time.Duration, exact bool, opts ...CallOption) error {
	deadline := time.Now().Add(timeout)
	for {
		el, err := d.FindText(text, exact, opts...)
		if err != nil {
			return err
		}
		if el == nil {
			return nil
		}
		if time.Now().After(deadline) {
			return &Error{Code: CodeTimeout, Message: "text still present after deadline: " + text, Retryable: true}
		}
		time.Sleep(defaultPollEvery)
	}
}

// WaitFor polls Observe until predicate(screen) is true or timeout elapses.
func (d *MobileDriver) WaitFor(predicate func(Screen) bool, timeout time.Duration) (*Screen, error) {
	deadline := time.Now().Add(timeout)
	for {
		screen, err := d.Observe()
		if err != nil {
			return nil, err
		}
		if predicate(*screen) {
			return screen, nil
		}
		if time.Now().After(deadline) {
			return nil, &Error{Code: CodeTimeout, Message: "predicate not satisfied within deadline", Retryable: true}
		}
		time.Sleep(defaultPollEvery)
	}
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
	_, err := d.call(methodInputKeyPress, keyPressParams{Key: key}, defaultCallTimeout)
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
	_, err := d.call(methodInputTap, tapParams{X: x, Y: y}, defaultCallTimeout)
	return err
}

func (d *MobileDriver) longPressXY(x, y, durationMs int) error {
	_, err := d.call(methodInputLongPress, longPressParams{X: x, Y: y, DurationMs: durationMs}, defaultCallTimeout)
	return err
}

func (d *MobileDriver) swipeXY(x1, y1, x2, y2, durationMs int) error {
	_, err := d.call(methodInputSwipe, swipeParams{X1: x1, Y1: y1, X2: x2, Y2: y2, DurationMs: durationMs}, defaultCallTimeout)
	return err
}

func (d *MobileDriver) typeText(text string) error {
	_, err := d.call(methodInputTypeText, typeTextParams{Text: text}, defaultCallTimeout)
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

// --- wire params/results (snake_case, mirrors sandbox/internal/sdkserver) ----

type tapParams struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type longPressParams struct {
	X          int `json:"x"`
	Y          int `json:"y"`
	DurationMs int `json:"duration_ms"`
}

type swipeParams struct {
	X1         int `json:"x1"`
	Y1         int `json:"y1"`
	X2         int `json:"x2"`
	Y2         int `json:"y2"`
	DurationMs int `json:"duration_ms"`
}

type typeTextParams struct {
	Text string `json:"text"`
}

type keyPressParams struct {
	Key string `json:"key"`
}

type observeParams struct {
	OcrEngine string `json:"ocr_engine"`
}

type findParams struct {
	Query     string `json:"query"`
	Model     string `json:"model,omitempty"`
	OcrEngine string `json:"ocr_engine"`
}

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

type wireFound struct {
	BBox       wireBBox `json:"bbox"`
	Text       string   `json:"text"`
	Confidence float64  `json:"confidence"`
}

type wireFind struct {
	Found *wireFound `json:"found"`
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

func (d *MobileDriver) elementFromText(w wireText) Element {
	b := bboxFromWire(w.BBox)
	return Element{BBox: b, Center: b.Center(), Confidence: w.Confidence, Text: w.Text, Source: SourceOCR, driver: d}
}

func (d *MobileDriver) elementFromFound(w *wireFound) *Element {
	b := bboxFromWire(w.BBox)
	source := SourceVLM
	if w.Text != "" {
		source = SourceOCR
	}
	el := Element{BBox: b, Center: b.Center(), Confidence: w.Confidence, Text: w.Text, Source: source, driver: d}
	return &el
}

func (d *MobileDriver) iconFromWire(w wireIcon) IconBox {
	b := bboxFromWire(w.BBox)
	return IconBox{BBox: b, Center: b.Center(), Confidence: w.Confidence}
}

func (d *MobileDriver) screenFromWire(raw json.RawMessage) (*Screen, error) {
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
		s.Texts = append(s.Texts, d.elementFromText(t))
	}
	for _, ic := range w.Icons {
		s.Icons = append(s.Icons, d.iconFromWire(ic))
	}
	return s, nil
}
