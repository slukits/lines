// Copyright (c) 2022 Stephan Lukits. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package lines

import (
	"testing"
	"time"

	"github.com/slukits/lines/internal/api"
	"github.com/slukits/lines/internal/lyt"
	"github.com/slukits/lines/internal/term"
)

// StringScreen is the string representation of the screen lines at a
// particular point in time.
type StringScreen = api.StringScreen

// CellsScreen is a screen representation at a specific point in time of
// lines of cells which also provide information about their styling.
type CellsScreen = api.CellsScreen

// CellsLine represents a line of a [lines.CellsScreen].
type CellsLine = api.CellsLine

// backend prevents the backend testing instance from being a public
// property of an lines.Testing instance.
type backend = *term.Fixture

// Fixture augments lines.Lines instance created by a *Fixture
// constructor with useful features for testing like emulating user
// input or getting the current screen content.
//
// The [Lines.WaitForQuit] method provided by a Fixture instance is
// non-blocking.
//
// It is guaranteed that all methods of an Fixture/Lines-instances which
// trigger an event do not return before the event and all subsequently
// triggered events are processed and any writes to environments are
// printed to the screen.
type Fixture struct {
	backend
	Lines      *Lines
	terminated bool
	syncAdd    chan bool
	syncWait   chan (chan bool)
	t          *testing.T
}

// TermFixture provides a slightly differently behaving Events instance and an
// augmenting Testing instance adding features useful for testing.
//
// The here provided Events instance has a non-blocking Listen method
// and all its methods triggering events are guaranteed to return after
// the event and subsequently triggered events have been processed and
// the (simulation) screen is synchronized.  All event triggering
// methods start the event loop automatically if not started, i.e. a
// call to Listen can be skipped.
//
// The Testing instance provides an event countdown which ends the event
// loop once it is zero.  Provide as last argument 0 for an indefinitely
// running event loop.  The default is 1.  NOTE reported OnInit and
// OnLayout events are accumulated and each is counted as one reported
// event for the event countdown.
//
// Testing provides methods for firing user input events which start the
// event-loop if not started and do return after the event and
// subsequently triggered events have been processed and the screen has
// been synchronized.
func TermFixture(
	t *testing.T,
	timeout time.Duration,
	c Componenter,
) *Fixture {
	t.Helper()
	ll := &Lines{}
	ui, backend := term.NewFixture(t, timeout)
	ll.Globals = newGlobals(nil)
	ll.scr = newScreen(ui, c, ll.Globals)
	ll.Globals.propagation = globalsPropagationClosure(ll.scr)
	ll.backend = ui
	tt := &Fixture{
		backend: backend,
		Lines:   ll,
		t:       t,
	}
	backend.Listen(ll.listen)
	return tt
}

func (tt *Fixture) Root() Componenter {
	if tt.Lines.scr.lyt.Root == nil {
		tt.t.Fatal("testing: root: layout not initialized")
	}
	return tt.Lines.scr.root().userComponent()
}

// FireResize posts a resize event and returns after this event
// has been processed.  Is associated Events instance not listening
// it is started before the event is fired.  NOTE this event as such is
// not reported, i.e. the event countdown is not reduced through this
// event.  But subsequently triggered OnInit or OnLayout events are
// counting down if reported.
func (tt *Fixture) FireResize(width, height int) {
	tt.t.Helper()
	if width == 0 && height == 0 {
		return
	}
	tt.PostResize(width, height)
}

// FireRune posts given run-key-press event and returns after this event
// has been processed.  Note modifier keys are ignored for
// rune-triggered key-events.  Is associated Events instance not
// listening it is started before the event is fired.
func (tt *Fixture) FireRune(r rune, m ...Modifier) {
	tt.t.Helper()
	if len(m) == 0 {
		tt.PostRune(r, api.ZeroModifier)
	} else {
		tt.PostRune(r, m[0])
	}
}

// FireKey posts given special-key event and returns after this
// event has been processed.  Is associated Events instance not
// listening it is started before the event is fired.
func (tt *Fixture) FireKey(k api.Key, m ...Modifier) {
	tt.t.Helper()
	if len(m) == 0 {
		tt.PostKey(k, api.ZeroModifier)
	} else {
		tt.PostKey(k, m[0])
	}
}

// FireClick posts a first button click at given coordinates and returns
// after this event has been processed.  Is associated Events instance
// not listening it is started before the event is fired.  Are given
// coordinates outside the available screen area the call is ignored.
func (tt *Fixture) FireClick(x, y int) {
	tt.t.Helper()
	width, height := tt.Lines.scr.backend.Size()
	if x < 0 || y < 0 || x >= width || y >= height {
		return
	}
	tt.PostMouse(x, y, api.Primary, api.ZeroModifier)
}

// FireContext posts a second button click at given coordinates and
// returns after this event has been processed.  Is associated Events
// instance not listening it is started before the event is fired.  Are
// given coordinates outside the available screen area the call is
// ignored.
func (tt *Fixture) FireContext(x, y int) {
	tt.t.Helper()
	width, height := tt.Lines.scr.backend.Size()
	if x < 0 || y < 0 || x >= width || y >= height {
		return
	}
	tt.PostMouse(x, y, api.Secondary, api.ZeroModifier)
}

// FireMouse posts a mouse event with provided arguments and returns
// after this event has been processed.  Is associated Events instance
// not listening it is started before the event is fired.  Are given
// coordinates outside the available screen area the call is ignored.
func (tt *Fixture) FireMouse(
	x, y int, bm api.Button, mm api.Modifier,
) {
	tt.t.Helper()
	width, height := tt.Lines.scr.backend.Size()
	if x < 0 || y < 0 || x >= width || y >= height {
		return
	}
	tt.PostMouse(x, y, bm, mm)
}

// FireComponentClick posts an click on given relative coordinate in
// given componenter.  Is associated Events instance not listening it is
// started before the event is fired.  Note if x or y are outside the
// component's screen area or the component is not part of the layout no
// click will be fired.
func (tt *Fixture) FireComponentClick(c Componenter, x, y int) {
	tt.t.Helper()
	if !c.hasLayoutWrapper() {
		return
	}
	ox, oy, ok := isInside(c.layoutComponent().wrapped().dim, x, y)
	if !ok {
		return
	}
	tt.FireClick(ox+x, oy+y)
}

// FireComponentContext posts an "right"-click on given relative
// coordinate in given componenter.  Is associated Events instance not
// listening it is started before the event is fired.  Note if x or y
// are outside the component's screen area or the component is not part
// of the layout no click will be fired.
func (tt *Fixture) FireComponentContext(c Componenter, x, y int) {
	tt.t.Helper()
	if !c.hasLayoutWrapper() {
		return
	}
	ox, oy, ok := isInside(c.layoutComponent().wrapped().dim, x, y)
	if !ok {
		return
	}
	tt.FireContext(ox+x, oy+y)
}

func isInside(dim *lyt.Dim, x, y int) (ox, oy int, ok bool) {
	if x < 0 || y < 0 || dim.IsOffScreen() {
		return 0, 0, false
	}
	_, _, width, height := dim.Area()
	if x >= width {
		return 0, 0, false
	}
	if y >= height {
		return 0, 0, false
	}
	return dim.X(), dim.Y(), true
}

// ScreenOf provides a string representation of given component's
// screen-area, i.e. without margins and without clippings.  The
// returned StringScreen is nil if given componenter is not part of the
// layout or off-screen.  Note call ScreenArea(c.Dim().Rect()) inside an
// Update-event callback to get the ScreenArea of a component including
// layout margins.
func (tt *Fixture) ScreenOf(c Componenter) api.StringScreen {
	if !c.hasLayoutWrapper() {
		return nil
	}
	dim := c.layoutComponent().wrapped().Dim()
	if dim.IsOffScreen() {
		return nil
	}
	return tt.ScreenArea(dim.Area())
}

// CellsOf provides a lines of cells representation of given component's
// screen-portion, i.e.  including margins and without clippings.  A
// CellsScreen provides next to a string representation also style
// information for each screen coordinate.  The returned CellsScreen is
// nil if given componenter is not part of the layout or off-screen.
// Note call CellsArea(c.Dim().Rect()) inside an Update-event callback
// to get the ScreenArea of a component including layout margins.
func (tt *Fixture) CellsOf(c Componenter) api.CellsScreen {
	if !c.hasLayoutWrapper() {
		return nil
	}
	dim := c.layoutComponent().wrapped().Dim()
	if dim.IsOffScreen() {
		return nil
	}
	return tt.CellsArea(dim.Area())
}
