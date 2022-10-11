// Copyright (c) 2022 Stephan Lukits. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package api

import (
	"time"

	"github.com/gdamore/tcell/v2"
)

// Color represents an rgb color.  Predefined colors are expressed in
// the typical hex-format 0xRRGGBB whereas R, G and B are hex-digits,
// i.e. red is 0xFF0000.
type Color int32

const sty = tcell.AttrBlink

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

// ScreenWriter implementation provides the screen/a window as a set of
// lines and cells to which a rune at a given position with a given
// style can be written.
type ScreenWriter interface {

	// Display "writes" given run with given style at given coordinates
	// to the screen.
	Display(rune, int, int, Style)

	// Update updates the screen.
	Update()

	// Redraw redraws the screen.
	Redraw()

	// Size reports the available screen/window size where as the width
	// is the number of single width runes fitting in a line and the
	// height is the number of lines fitting on the screen.
	Size() (int, int)
}

// EventProcessor provides user input events and programmatically posted
// events.
type EventProcessor interface {
	// Poll provides user input events and programmatically posted
	// events around which lines is looping.
	Poll() Eventer
	// Post posts given event to the event loop.
	Post(Eventer)
}

// An UIer implementation provides the functionality lines needs to
// provide its features.
type UIer interface {
	ScreenWriter
	EventProcessor
}

type TestCell struct {
	Rune rune
	Sty  Style
}

// TestLine represents a line of a [lines.TestScreen].
type TestLine []TestCell

type TestScreen []TestLine

// Tester implementation augments an UIer implementation with additional
// functionality for testing.
type Tester interface {

	// String returns a string representation of the screen/window
	// content.
	String() string

	// StringArea returns a string representation of given screen/window
	// area.
	StringArea(x, y, width, height int) string

	// TrimString reduces given string to its minimum number of
	// non-empty lines whereas the lines are trimmed to contain all non
	// white space runes:
	//
	// 	+--------------------+
	// 	|                    |       +------------+
	// 	|   upper left       |       |upper left  |
	// 	|                    |  =>   |            |
	// 	|          right     |       |       right|
	// 	|      bottom        |       |   bottom   |
	// 	|                    |       +------------+
	// 	+--------------------+
	//
	TrimString(string) string

	// Screen returns the content of a test screen, i.e. in addition to
	// a string representation also the style information is provided.
	Screen() TestScreen

	// ScreenArea returns the test screen of given screen area, i.e. in
	// addition to a string representation also the style information is
	// provided.
	ScreenArea(x, y, width, height int) TestScreen

	// TrimScreen reduces like TrimString a TestScreen to its minimal
	// size still providing all cells with content other than white
	// space.
	TrimScreen(TestScreen) TestScreen
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

	// Button returns the button number of the mouse event.
	Button() Button

	// Mod reports the pressed modifier key like shift, alt, ...
	Mod() Modifier

	// Pos returns the x- and y-screen-coordinates of a mouse event.
	Pos() (int, int)
}

// ResizeEventer implementation is reported on a screen/window-size
// change.
type ResizeEventer interface {

	// Size reports the width, i.e. the number of runes fitting in a
	// screen/window line, and the height, i.e. the number of lines fitting on
	// the screen/window, of the resize event.
	Size() (int, int)
}

// QuitEventer implementation is reported on to quit the application.
type QuitEventer interface {
	// Quit is only there for a discriminating type switch.
	Quit()
}
