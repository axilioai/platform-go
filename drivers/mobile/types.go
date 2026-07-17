// Package mobile is the hand-written phone driver for the Axilio Go SDK: the Go
// twin of platform-python's axilio.drivers.mobile. It drives a paired phone
// through a Transport (the DCP control WebSocket), turning an ergonomic
// observe/find/tap/type API into literal CDP method frames. It lives alongside
// the Fern-generated REST client and is preserved across regen by the rsync
// exclude list in scripts/regen.sh — .fernignore does not work here (see that
// script's header).
package mobile

import (
	"regexp"
	"strings"
	"time"
)

// Coords is a point in frame-space pixels.
type Coords struct {
	X int `json:"x"`
	Y int `json:"y"`
}

// BBox is an axis-aligned box in frame-space pixels (top-left origin).
type BBox struct {
	X      int `json:"x"`
	Y      int `json:"y"`
	Width  int `json:"width"`
	Height int `json:"height"`
}

// Center returns the midpoint of the box.
func (b BBox) Center() Coords {
	return Coords{X: b.X + b.Width/2, Y: b.Y + b.Height/2}
}

// Source records how an element was located.
type Source string

const (
	// SourceOCR is an element found by OCR (carries text).
	SourceOCR Source = "ocr"
	// SourceVLM is an element found by the vision model (no text).
	SourceVLM Source = "vlm"
)

// Element is one located element — the universal selector return type; actions
// chain off it (Tap, LongPress, TypeInto, SwipeTo). The unexported driver
// back-reference is what those chained actions drive; it is never serialised.
type Element struct {
	BBox       BBox    `json:"bbox"`
	Center     Coords  `json:"center"`
	Confidence float64 `json:"confidence"`
	// Text is the OCR text, empty for a VLM-sourced element.
	Text   string `json:"text,omitempty"`
	Source Source `json:"source"`

	driver *MobileDriver
}

// Tap taps at the element's center.
func (e Element) Tap() error { return e.driver.tapXY(e.Center.X, e.Center.Y) }

// LongPress presses and holds at the element's center for durationMs.
func (e Element) LongPress(durationMs int) error {
	return e.driver.longPressXY(e.Center.X, e.Center.Y, durationMs)
}

// TypeInto taps the element, then types text into it.
func (e Element) TypeInto(text string) error {
	if err := e.Tap(); err != nil {
		return err
	}
	return e.driver.typeText(text)
}

// SwipeTo swipes from this element's center to other's center over durationMs.
func (e Element) SwipeTo(other Element, durationMs int) error {
	return e.driver.swipeXY(e.Center.X, e.Center.Y, other.Center.X, other.Center.Y, durationMs)
}

// IconBox is one YOLO-detected icon (rectangle-only; the icon model has a single class).
type IconBox struct {
	BBox       BBox    `json:"bbox"`
	Center     Coords  `json:"center"`
	Confidence float64 `json:"confidence"`
}

// Screen is an immutable snapshot of one observed frame.
type Screen struct {
	Texts      []Element `json:"texts"`
	Icons      []IconBox `json:"icons"`
	Hash       string    `json:"hash"`
	Width      int       `json:"width"`
	Height     int       `json:"height"`
	CapturedAt time.Time `json:"captured_at"`
}

// FindText returns the first OCR element whose text matches, or nil. When exact
// is false the match is case-insensitive substring.
func (s Screen) FindText(text string, exact bool) *Element {
	for i := range s.Texts {
		if matches(s.Texts[i].Text, text, exact) {
			return &s.Texts[i]
		}
	}
	return nil
}

// FindAllText returns every OCR element matching the criteria. contains and
// pattern are mutually exclusive; passing both is an error. With neither, every
// text element is returned.
func (s Screen) FindAllText(contains, pattern string) ([]Element, error) {
	if contains != "" && pattern != "" {
		return nil, &Error{Code: CodeInvalidArgs, Message: "FindAllText: pass at most one of contains / pattern"}
	}
	switch {
	case contains != "":
		needle := strings.ToLower(contains)
		out := make([]Element, 0, len(s.Texts))
		for _, el := range s.Texts {
			if el.Text != "" && strings.Contains(strings.ToLower(el.Text), needle) {
				out = append(out, el)
			}
		}
		return out, nil
	case pattern != "":
		re, err := regexp.Compile(pattern)
		if err != nil {
			return nil, &Error{Code: CodeInvalidArgs, Message: "FindAllText: bad pattern: " + err.Error()}
		}
		out := make([]Element, 0, len(s.Texts))
		for _, el := range s.Texts {
			if el.Text != "" && re.MatchString(el.Text) {
				out = append(out, el)
			}
		}
		return out, nil
	default:
		return append([]Element(nil), s.Texts...), nil
	}
}

func matches(text, needle string, exact bool) bool {
	if text == "" {
		return false
	}
	if exact {
		return text == needle
	}
	return strings.Contains(strings.ToLower(text), strings.ToLower(needle))
}
