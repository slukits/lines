// Copyright (c) 2022 Stephan Lukits. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package api

import "time"

// StyleAttributeMask defines the looks of a style, i.e. the looks of a print
// to the screen/window.
type StyleAttributeMask int32

const (
	Bold StyleAttributeMask = 1 << iota
	Blink
	Reverse
	Underline
	Dim
	Italic
	StrikeThrough
	Invalid
	ZeroStyle StyleAttributeMask = 0
)

type StyleSemantic uint8

// Style represents what a print to the screen should look like.  A zero
// Style instance has fore- and background color set to black.  Use its
// With* methods to create a style with desired properties:
//
// myStyle := (lines.Style{}).WithFG(lines.White)
type Style struct {
	// AA is the style attribute mask providing set style attributes
	aa StyleAttributeMask

	// FG provides a style's foreground color
	fg Color

	// BG provides a style's background color
	bg Color

	sm StyleSemantic
}

// DefaultStyle has no attributes and "default" colors.  The semantics
// of the later is decided by the backend implementation.  Use the With*
// methods to create new styles from the default style.
var DefaultStyle = Style{fg: DefaultColor, bg: DefaultColor}

// NewStyle creates a new style with given style attributes and given fore-
// and background color.
func NewStyle(aa StyleAttributeMask, fg, bg Color) Style {
	return Style{aa: aa, fg: fg, bg: bg}
}

func (s Style) AA() StyleAttributeMask { return s.aa }
func (s Style) FG() Color              { return s.fg }
func (s Style) BG() Color              { return s.bg }
func (s Style) SMN() StyleSemantic     { return s.sm }

// IsDefault returns true if given Style s is the DefaultStyle.
func (s Style) IsDefault() bool { return s == DefaultStyle }

// Reverse returns given Style s reversed, i.e. with the Reverse-bit set
// if s has it not set or with the Reverse-bit unset otherwise.
func (s Style) Reverse() Style {
	if s.aa&Reverse == Reverse {
		return Style{fg: s.fg, bg: s.bg, aa: s.aa &^ Reverse, sm: s.sm}
	}
	return Style{fg: s.fg, bg: s.bg, aa: s.aa | Reverse, sm: s.sm}
}

// Switch removes given attribute a if it is set or adds it otherwise.
func (s Style) Switch(a StyleAttributeMask) Style {
	if s.aa&a == a {
		return Style{fg: s.fg, bg: s.bg, aa: s.aa &^ a, sm: s.sm}
	}
	return Style{fg: s.fg, bg: s.bg, aa: s.aa | a, sm: s.sm}
}

// Invert returns a style having given Style s foreground color as
// background color and its background color as foreground color.
func (s Style) Invert() Style {
	return Style{fg: s.bg, bg: s.fg, aa: s.aa, sm: s.sm}
}

// WithAdded returns given style with given attribute mask added.
func (s Style) WithAdded(aa StyleAttributeMask) Style {
	return Style{fg: s.fg, bg: s.bg, aa: s.aa | aa, sm: s.sm}
}

// WithRemoved returns given style without given attribute mask.
func (s Style) WithRemoved(aa StyleAttributeMask) Style {
	return Style{fg: s.fg, bg: s.bg, aa: s.aa &^ aa, sm: s.sm}
}

// WithAA returns given style with its attributes set to given attribute
// mask.
func (s Style) WithAA(aa StyleAttributeMask) Style {
	return Style{fg: s.fg, bg: s.bg, aa: aa, sm: s.sm}
}

// WithFG returns given style with its foreground color set to given
// color.
func (s Style) WithFG(c Color) Style {
	return Style{fg: c, bg: s.bg, aa: s.aa, sm: s.sm}
}

// WithBG returns given style with its background color set to given
// color.
func (s Style) WithBG(c Color) Style {
	return Style{fg: s.fg, bg: c, aa: s.aa, sm: s.sm}
}

// WithSemantics returns given style s with given semantics (ID) set.
// Note the lines package doesn't define any semantics (yet).  The sole
// purpose of the semantics property is to give the user an option to
// discriminate two semantics which might accidentally have the same
// colors and attributes but have two different meanings.
func (s Style) WithSemantics(sm StyleSemantic) Style {
	return Style{fg: s.fg, bg: s.bg, aa: s.aa, sm: sm}
}

// Displayer implementation provides the screen/a window as a set of
// lines and cells to which a rune at a given position with a given
// style can be written.
type Displayer interface {

	// NewStyle must be used to obtain a new style with backend specific
	// defaults.
	NewStyle() Style

	// Display "writes" given run with given style at given coordinates
	// to the screen.
	Display(int, int, rune, Style)

	// Update updates the screen.
	Update()

	// Redraw redraws the screen.
	Redraw()

	// Size reports the available screen/window size whereas the width
	// is the number of single width runes fitting in a line and the
	// height is the number of lines fitting on the screen.
	Size() (int, int)

	// SetCursor positions the screen/window cursor at given coordinates
	// x and y having optionally given cursor style cs and returns the
	// actually set cursor position with actually set cursor style:
	//  - if x and y are outside the screen or they are inside and given
	//    cursor style is the ZeroCursor the -1, -1, ZeroCursor is
	//    returned.
	//  - if x and y are inside the screen and no cursor style is given
	//    x and y and ZeroCursor is return indicating that the cursor
	//    style has not changed.
	//  - if x and y are inside the screen and a non-zero cursor style
	//    is given the arguments are returned as received.
	SetCursor(x, y int, cs ...CursorStyle) (int, int, CursorStyle)

	// Colors provide the number of available (ANSI) colors.  In case of
	// a monochrome screen 0 is returned.
	Colors() int
}

// EventProcessor provides user input events and programmatically posted
// events.
type EventProcessor interface {

	// Post posts given event to the event loop.
	Post(Eventer) error

	// Quit event polling.
	Quit()

	// WaitForQuit blocks until the backend was quit.
	WaitForQuit()

	// OnQuit registers given function to be called on quitting.
	OnQuit(listener func())
}

// An UIer implementation provides the functionality lines needs to
// provide its features.
type UIer interface {
	Displayer
	EventProcessor

	// Lib provides access to the encapsulated backend library.
	Lib() interface{}
}

// Eventer is the abstract interface which must be implemented by all
// reported/posted events.
type Eventer interface {
	// When returns the creation time of an event.
	When() time.Time
	// Source returns the wrapped event of the backend.
	Source() interface{}
}

// KeyEventer implementation is reported on a user special-key input
// event like "enter" or "backspace".
type KeyEventer interface {
	Eventer

	// Key reports the pressed key.
	Key() Key

	// Mod reports the pressed modifier key like shift, alt, ...
	Mod() ModifierMask
}

// RuneEventer implementation is reported on a user rune input event.
type RuneEventer interface {
	Eventer

	// Rune reports the pressed rune.
	Rune() rune

	// Mod reports the pressed modifier key like shift, alt, ...
	Mod() ModifierMask
}

// MouseEventer implementation is reported on a user-input mouse event.
type MouseEventer interface {
	Eventer

	// Button returns the buttons mask of the mouse event.
	Button() ButtonMask

	// Mod reports the pressed modifier key like shift, alt, ...
	Mod() ModifierMask

	// Pos returns the x- and y-screen-coordinates of a mouse event.
	Pos() (int, int)
}

// ResizeEventer implementation is reported on a screen/window-size
// change.
type ResizeEventer interface {
	Eventer

	// Size reports the width, i.e. the number of runes fitting in a
	// screen/window line, and the height, i.e. the number of lines fitting on
	// the screen/window, of the resize event.
	Size() (int, int)
}

// Gaps is internally used for layout calculations of nested components
// of either stacking or chaining gaped components.  The layout wrapper
// of stacking/chaining components created during a component's
// initialization process provides a Gaps instance to the layout
// manager.
type Gaps struct {
	Top, Right, Bottom, Left int
}
