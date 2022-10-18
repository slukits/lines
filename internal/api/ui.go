// Copyright (c) 2022 Stephan Lukits. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package api

import "time"

// Color represents an rgb color.  Predefined colors are expressed in
// the typical hex-format 0xRRGGBB whereas R, G and B are hex-digits,
// i.e. red is 0xFF0000.
type Color int32

// StyleAttribute define the looks of a style, i.e. the looks of a print
// to the screen.
type StyleAttribute int32

const (
	Bold StyleAttribute = 1 << iota
	Blink
	Reverse
	Underline
	Dim
	Italic
	StrikeThrough
	Invalid
	ZeroStyle StyleAttribute = 0
)

// Style represents what a print to the screen should look like.
type Style struct {
	// AA is the style attribute mask providing set style attributes
	AA StyleAttribute

	// FG provides a style's foreground color
	FG Color

	// BG provides a style's background color
	BG Color
}

func (s Style) Equals(other Style) bool {
	return s.AA == other.AA && s.FG == other.FG && s.BG == other.BG
}

func (s Style) WithAttrsAdded(aa StyleAttribute) Style {
	return Style{FG: s.FG, BG: s.BG, AA: s.AA | aa}
}

func (s Style) WithAttrsRemoved(aa StyleAttribute) Style {
	return Style{FG: s.FG, BG: s.BG, AA: s.AA &^ aa}
}

func (s Style) WithAttrs(aa StyleAttribute) Style {
	return Style{FG: s.FG, BG: s.BG, AA: aa}
}

func (s Style) WithFG(c Color) Style {
	return Style{FG: c, BG: s.BG, AA: s.AA}
}

func (s Style) WithBG(c Color) Style {
	return Style{FG: s.FG, BG: c, AA: s.AA}
}

// Displayer implementation provides the screen/a window as a set of
// lines and cells to which a rune at a given position with a given
// style can be written.
type Displayer interface {

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
	Mod() Modifier
}

// RuneEventer implementation is reported on a user rune input event.
type RuneEventer interface {
	Eventer

	// Rune reports the pressed rune.
	Rune() rune

	// Mod reports the pressed modifier key like shift, alt, ...
	Mod() Modifier
}

// MouseEventer implementation is reported on a user-input mouse event.
type MouseEventer interface {
	Eventer

	// Button returns the buttons mask of the mouse event.
	Button() Button

	// Mod reports the pressed modifier key like shift, alt, ...
	Mod() Modifier

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

type QuitEventer interface {
	Eventer
	// Quitting is only defined to discriminate from Eventer.
	Quitting()
}
