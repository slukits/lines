// Copyright (c) 2022 Stephan Lukits. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package lines

import (
	"errors"

	"github.com/gdamore/tcell/v2"
)

// ErrScreen is the error returned by the mocked screen creation.
var ErrScreen = errors.New("mock: screen: failing creation")

// ScreenFactory for the mocked screen creation/initialization.
type mockScreenFactory struct{ Fail, FailInit bool }

// NewScreen mock tcell's screen creation.
func (f *mockScreenFactory) NewScreen() (tcell.Screen, error) {
	if f.Fail {
		return nil, ErrScreen
	}
	return &screenMock{Fail: f.FailInit}, nil
}

// NewSimulationScreen mock tcell's simulation screen creation.
func (f *mockScreenFactory) NewSimulationScreen(string) tcell.SimulationScreen {
	return &screenMock{Fail: f.FailInit}
}

// ErrInit is the error returned by the mocked screen initialization.
var ErrInit = errors.New("mock: screen: failing initialization")

// ScreenMock mocks tcell's SimulationScreen implementation to mock-up its
// possible creation/initialization-errors.
type screenMock struct{ Fail bool }

// Init initializes the screen for use.
func (s *screenMock) Init() error {
	if s.Fail {
		return ErrInit
	}
	return nil
}

// InjectKeyBytes injects a stream of bytes corresponding to
// the native encoding (see charset).  It turns true if the entire
// set of bytes were processed and delivered as KeyEvents, false
// if any bytes were not fully understood.  Any bytes that are not
// fully converted are discarded.
func (s *screenMock) InjectKeyBytes(buf []byte) bool {
	panic("not implemented")
}

// InjectKey injects a key event.  The rune is a UTF-8 rune, post
// any translation.
func (s *screenMock) InjectKey(key tcell.Key, r rune, mod tcell.ModMask) {
	panic("not implemented")
}

// InjectMouse injects a mouse event.
func (s *screenMock) InjectMouse(x int, y int, buttons tcell.ButtonMask, mod tcell.ModMask) {
	panic("not implemented")
}

// SetSize resizes the underlying physical screen.  It also causes
// a resize event to be injected during the next Show() or Sync().
// A new physical contents array will be allocated (with data from
// the old copied), so any prior value obtained with GetContents
// won't be used anymore
func (s *screenMock) SetSize(width int, height int) {
	panic("not implemented")
}

// GetContents returns screen contents as an array of
// cells, along with the physical width & height.   Note that the
// physical contents will be used until the next time SetSize()
// is called.
func (s *screenMock) GetContents() (cells []tcell.SimCell, width int, height int) {
	panic("not implemented")
}

// GetCursor returns the cursor details.
func (s *screenMock) GetCursor() (x int, y int, visible bool) {
	panic("not implemented")
}

// Fini finalizes the screen also releasing resources.
func (s *screenMock) Fini() {
	panic("not implemented")
}

// Clear erases the screen.  The contents of any screen buffers
// will also be cleared.  This has the logical effect of
// filling the screen with spaces, using the global default style.
func (s *screenMock) Clear() {
	panic("not implemented")
}

// Fill fills the screen with the given character and style.
func (s *screenMock) Fill(_ rune, _ tcell.Style) {
	panic("not implemented")
}

// SetCell is an older API, and will be removed.  Please use
// SetContent instead; SetCell is implemented in terms of SetContent.
func (s *screenMock) SetCell(x int, y int, style tcell.Style, ch ...rune) {
	panic("not implemented")
}

// GetContent returns the contents at the given location.  If the
// coordinates are out of range, then the values will be 0, nil,
// StyleDefault.  Note that the contents returned are logical contents
// and may not actually be what is displayed, but rather are what will
// be displayed if Show() or Sync() is called.  The width is the width
// in screen cells; most often this will be 1, but some East Asian
// characters require two cells.
func (s *screenMock) GetContent(x int, y int) (mainc rune, combc []rune, style tcell.Style, width int) {
	panic("not implemented")
}

// SetContent sets the contents of the given cell location.  If
// the coordinates are out of range, then the operation is ignored.
//
// The first rune is the primary non-zero width rune.  The array
// that follows is a possible list of combining characters to append,
// and will usually be nil (no combining characters.)
//
// The results are not displayd until Show() or Sync() is called.
//
// Note that wide (East Asian full width) runes occupy two cells,
// and attempts to place character at next cell to the right will have
// undefined effects.  Wide runes that are printed in the
// last column will be replaced with a single width space on output.
func (s *screenMock) SetContent(x int, y int, mainc rune, combc []rune, style tcell.Style) {
	panic("not implemented")
}

// SetStyle sets the default style to use when clearing the screen
// or when StyleDefault is specified.  If it is also StyleDefault,
// then whatever system/terminal default is relevant will be used.
func (s *screenMock) SetStyle(style tcell.Style) {
	panic("not implemented")
}

// ShowCursor is used to display the cursor at a given location.
// If the coordinates -1, -1 are given or are otherwise outside the
// dimensions of the screen, the cursor will be hidden.
func (s *screenMock) ShowCursor(x int, y int) {
	panic("not implemented")
}

// HideCursor is used to hide the cursor.  Its an alias for
// ShowCursor(-1, -1).sim
func (s *screenMock) HideCursor() {
	panic("not implemented")
}

// SetCursorStyle is used to set the cursor style.  If the style
// is not supported (or cursor styles are not supported at all),
// then this will have no effect.
func (s *screenMock) SetCursorStyle(_ tcell.CursorStyle) {
	panic("not implemented")
}

// Size returns the screen size as width, height.  This changes in
// response to a call to Clear or Flush.
func (s *screenMock) Size() (width int, height int) {
	panic("not implemented")
}

// ChannelEvents is an infinite loop that waits for an event and
// channels it into the user provided channel ch.  Closing the
// quit channel and calling the Fini method are cancellation
// signals.  When a cancellation signal is received the method
// returns after closing ch.
//
// This method should be used as a goroutine.
//
// NOTE: PollEvent should not be called while this method is running.
func (s *screenMock) ChannelEvents(ch chan<- tcell.Event, quit <-chan struct{}) {
	panic("not implemented")
}

// PollEvent waits for events to arrive.  Main application loops
// must spin on this to prevent the application from stalling.
// Furthermore, this will return nil if the Screen is finalized.
func (s *screenMock) PollEvent() tcell.Event {
	panic("not implemented")
}

// HasPendingEvent returns true if PollEvent would return an event
// without blocking.  If the screen is stopped and PollEvent would
// return nil, then the return value from this function is unspecified.
// The purpose of this function is to allow multiple events to be collected
// at once, to minimize screen redraws.
func (s *screenMock) HasPendingEvent() bool {
	panic("not implemented")
}

// PostEvent tries to post an event into the event stream.  This
// can fail if the event queue is full.  In that case, the event
// is dropped, and ErrEventQFull is returned.
func (s *screenMock) PostEvent(ev tcell.Event) error {
	panic("not implemented")
}

// Deprecated: PostEventWait is unsafe, and will be removed
// in the future.
//
// PostEventWait is like PostEvent, but if the queue is full, it
// blocks until there is space in the queue, making delivery
// reliable.  However, it is VERY important that this function
// never be called from within whatever event loop is polling
// with PollEvent(), otherwise a deadlock may arise.
//
// For this reason, when using this function, the use of a
// Goroutine is recommended to ensure no deadlock can occur.
func (s *screenMock) PostEventWait(ev tcell.Event) {
	panic("not implemented")
}

// EnableMouse enables the mouse.  (If your terminal supports it.)
// If no flags are specified, then all events are reported, if the
// terminal supports them.
func (s *screenMock) EnableMouse(_ ...tcell.MouseFlags) {
	panic("not implemented")
}

// DisableMouse disables the mouse.
func (s *screenMock) DisableMouse() {
	panic("not implemented")
}

// EnablePaste enables bracketed paste mode, if supported.
func (s *screenMock) EnablePaste() {
	panic("not implemented")
}

// DisablePaste disables bracketed paste mode.
func (s *screenMock) DisablePaste() {
	panic("not implemented")
}

// HasMouse returns true if the terminal (apparently) supports a
// mouse.  Note that the a return value of true doesn't guarantee that
// a mouse/pointing device is present; a false return definitely
// indicates no mouse support is available.
func (s *screenMock) HasMouse() bool {
	panic("not implemented")
}

// Colors returns the number of colors.  All colors are assumed to
// use the ANSI color map.  If a terminal is monochrome, it will
// return 0.
func (s *screenMock) Colors() int {
	panic("not implemented")
}

// Show makes all the content changes made using SetContent() visible
// on the display.
//
// It does so in the most efficient and least visually disruptive
// manner possible.
func (s *screenMock) Show() {
	panic("not implemented")
}

// Sync works like Show(), but it updates every visible cell on the
// physical display, assuming that it is not synchronized with any
// internal model.  This may be both expensive and visually jarring,
// so it should only be used when believed to actually be necessary.
//
// Typically this is called as a result of a user-requested redraw
// (e.g. to clear up on screen corruption caused by some other program),
// or during a resize event.
func (s *screenMock) Sync() {
	panic("not implemented")
}

// CharacterSet returns information about the character set.
// This isn't the full locale, but it does give us the input/output
// character set.  Note that this is just for diagnostic purposes,
// we normally translate input/output to/from UTF-8, regardless of
// what the user's environment is.
func (s *screenMock) CharacterSet() string {
	panic("not implemented")
}

// RegisterRuneFallback adds a fallback for runes that are not
// part of the character set -- for example one could register
// o as a fallback for ø.  This should be done cautiously for
// characters that might be displayed ordinarily in language
// specific text -- characters that could change the meaning of
// of written text would be dangerous.  The intention here is to
// facilitate fallback characters in pseudo-graphical applications.
//
// If the terminal has fallbacks already in place via an alternate
// character set, those are used in preference.  Also, standard
// fallbacks for graphical characters in the ACSC terminfo string
// are registered implicitly.
//
// The display string should be the same width as original rune.
// This makes it possible to register two character replacements
// for full width East Asian characters, for example.
//
// It is recommended that replacement strings consist only of
// 7-bit ASCII, since other characters may not display everywhere.
func (s *screenMock) RegisterRuneFallback(r rune, subst string) {
	panic("not implemented")
}

// UnregisterRuneFallback unmaps a replacement.  It will unmap
// the implicit ASCII replacements for alternate characters as well.
// When an unmapped char needs to be displayed, but no suitable
// glyph is available, '?' is emitted instead.  It is not possible
// to "disable" the use of alternate characters that are supported
// by your terminal except by changing the terminal database.
func (s *screenMock) UnregisterRuneFallback(r rune) {
	panic("not implemented")
}

// CanDisplay returns true if the given rune can be displayed on
// this screen.  Note that this is a best guess effort -- whether
// your fonts support the character or not may be questionable.
// Mostly this is for folks who work outside of Unicode.
//
// If checkFallbacks is true, then if any (possibly imperfect)
// fallbacks are registered, this will return true.  This will
// also return true if the terminal can replace the glyph with
// one that is visually indistinguishable from the one requested.
func (s *screenMock) CanDisplay(r rune, checkFallbacks bool) bool {
	panic("not implemented")
}

// Resize does nothing, since its generally not possible to
// ask a screen to resize, but it allows the Screen to implement
// the View interface.
func (s *screenMock) Resize(_ int, _ int, _ int, _ int) {
	panic("not implemented")
}

// HasKey returns true if the keyboard is believed to have the
// key.  In some cases a keyboard may have keys with this name
// but no support for them, while in others a key may be reported
// as supported but not actually be usable (such as some emulators
// that hijack certain keys).  Its best not to depend to strictly
// on this function, but it can be used for hinting when building
// menus, displayed hot-keys, etc.  Note that KeyRune (literal
// runes) is always true.
func (s *screenMock) HasKey(_ tcell.Key) bool {
	panic("not implemented")
}

// Suspend pauses input and output processing.  It also restores the
// terminal settings to what they were when the application started.
// This can be used to, for example, run a sub-shell.
func (s *screenMock) Suspend() error {
	panic("not implemented")
}

// Resume resumes after Suspend().
func (s *screenMock) Resume() error {
	panic("not implemented")
}

// Beep attempts to sound an OS-dependent audible alert and returns an error
// when unsuccessful.
func (s *screenMock) Beep() error {
	panic("not implemented")
}
