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

// backend prevents the backend testing instance from being a public
// property of an lines.Testing instance.
type backend = *term.Testing

// Testing augments lines.Events instance created by *Test* with useful
// features for testing like firing an event or getting the current
// screen content as string.
//
// An Events/Testing-instances may not be used concurrently.
//
// An Events.Listen-method becomes non-blocking and starts the
// event-loop in its own go-routine.
//
// All event triggering methods start event-listening if it is not
// already started.
//
// It is guaranteed that all methods of an Events/Testing-instances
// which trigger an event do not return before the event is processed
// and any writes to environments are printed to the screen.  This holds
// also true if an event triggering method is called within a listener
// callback.
type Testing struct {
	backend
	ll         *Lines
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
) (*Lines, *Testing) {
	t.Helper()
	ll := &Lines{}
	ui, backend := term.Fixture(t, timeout)
	ll.scr = newScreen(ui, c)
	ll.backend = ui
	tt := &Testing{
		backend: backend,
		ll:      ll,
		t:       t,
	}
	backend.Listen(ll.listen)
	return ll, tt
}

func (tt *Testing) Root() Componenter {
	if tt.ll.scr.lyt.Root == nil {
		tt.t.Fatal("testing: root: layout not initialized")
	}
	return tt.ll.scr.root().userComponent()
}

// FireResize posts a resize event and returns after this event
// has been processed.  Is associated Events instance not listening
// it is started before the event is fired.  NOTE this event as such is
// not reported, i.e. the event countdown is not reduced through this
// event.  But subsequently triggered OnInit or OnLayout events are
// counting down if reported.
func (tt *Testing) FireResize(width, height int) {
	tt.t.Helper()
	if width == 0 && height == 0 {
		return
	}
	if err := tt.PostResize(width, height); err != nil {
		tt.t.Fatal(err)
	}
}

// FireRune posts given run-key-press event and returns after this event
// has been processed.  Note modifier keys are ignored for
// rune-triggered key-events.  Is associated Events instance not
// listening it is started before the event is fired.
func (tt *Testing) FireRune(r rune, m ...Modifier) {
	tt.t.Helper()
	if err := tt.PostRune(r, api.ZeroModifier); err != nil {
		tt.t.Fatal(err)
	}
}

// FireKey posts given special-key event and returns after this
// event has been processed.  Is associated Events instance not
// listening it is started before the event is fired.
func (tt *Testing) FireKey(k api.Key, m ...Modifier) {
	tt.t.Helper()
	var err error
	if len(m) == 0 {
		err = tt.PostKey(k, api.ZeroModifier)
	} else {
		err = tt.PostKey(k, m[0])
	}
	if err != nil {
		tt.t.Fatal(err)
	}
}

// FireClick posts a first button click at given coordinates and returns
// after this event has been processed.  Is associated Events instance
// not listening it is started before the event is fired.  Are given
// coordinates outside the available screen area the call is ignored.
func (tt *Testing) FireClick(x, y int) {
	tt.t.Helper()
	width, height := tt.ll.scr.backend.Size()
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
func (tt *Testing) FireContext(x, y int) {
	tt.t.Helper()
	width, height := tt.ll.scr.backend.Size()
	if x < 0 || y < 0 || x >= width || y >= height {
		return
	}
	tt.PostMouse(x, y, api.Secondary, api.ZeroModifier)
}

// FireMouse posts a mouse event with provided arguments and returns
// after this event has been processed.  Is associated Events instance
// not listening it is started before the event is fired.  Are given
// coordinates outside the available screen area the call is ignored.
func (tt *Testing) FireMouse(
	x, y int, bm api.Button, mm api.Modifier,
) {
	tt.t.Helper()
	width, height := tt.ll.scr.backend.Size()
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
func (tt *Testing) FireComponentClick(c Componenter, x, y int) {
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
func (tt *Testing) FireComponentContext(c Componenter, x, y int) {
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
// screen-portion, i.e.  including margins and without clippings.  The
// returned StringScreen is nil if given componenter is not part of the
// layout or off-screen.  NOTE do not use this method inside an
// Update-event listener.
func (tt *Testing) ScreenOf(c Componenter) api.StringScreen {
	if !c.hasLayoutWrapper() {
		return nil
	}
	dim := c.layoutComponent().wrapped().Dim()
	if dim.IsOffScreen() {
		return nil
	}
	return tt.ScreenArea(dim.Rect())
}

// CellsOf provides a lines of cells representation of given component's
// screen-portion, i.e.  including margins and without clippings.  A
// CellsScreen provides next to a string representation also style
// information for each screen coordinate.  The returned CellsScreen is
// nil if given componenter is not part of the layout or off-screen.
// NOTE do not use this method inside an Update-event listener.
func (tt *Testing) CellsOf(c Componenter) api.CellsScreen {
	if !c.hasLayoutWrapper() {
		return nil
	}
	dim := c.layoutComponent().wrapped().Dim()
	if dim.IsOffScreen() {
		return nil
	}
	return tt.CellsArea(dim.Rect())
}
