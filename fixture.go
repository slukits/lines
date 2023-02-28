// Copyright (c) 2022 Stephan Lukits. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package lines

import (
	"sync"
	"testing"
	"time"

	"github.com/slukits/lines/internal/api"
	"github.com/slukits/lines/internal/lyt"
	"github.com/slukits/lines/internal/term"
)

// StringScreen is the string representation of the screen lines at a
// particular point in time of [Fixture]'s [Lines] instance.  E.g. see
// [Fixture.ScreenOf] or Fixture.Screen.  NOTE use StringScreen's
// Trimmed-method to minimize the reported screen area.
type StringScreen = api.StringScreen

// CellsScreen is a screen representation at a specific point in time of
// of a [Fixtures]'s [Lines] instances.  E.g. see [Fixture.CellsOf] or
// Fixture.Cells.  NOTE use CellsScreen's Trimmed-method to minimize the
// reported screen area.
type CellsScreen = api.CellsScreen

// CellsLine represents a line of a [CellsScreen] providing of each cell
// in the line its displayed rune and style information for
// test-evaluations.
type CellsLine = api.CellsLine

// backend prevents the backend testing instance from being a public
// property of an lines.Testing instance.
type backend = *term.Fixture

// Fixture augments the [Lines] instance created by a *Fixture
// constructor like [TermFixture] with useful features for testing like
// emulating user input or getting the current screen content.
//
// Note The [Lines.WaitForQuit] method provided by a Fixture instance is
// non-blocking.
//
// It is guaranteed that all methods of an Fixture's Lines-instances
// which trigger an event do not return before the event and all
// subsequently triggered events are processed and any writes to
// environments are printed to the screen.
type Fixture struct {
	backend
	terminated bool
	syncAdd    chan bool
	syncWait   chan (chan bool)
	t          *testing.T

	// Lines instance created by the fixture constructor reporting
	// events to Componenter of the layout.
	Lines *Lines

	Scroll Scroll
}

// TermFixture returns a Fixture instance with a slightly differently
// behaving [Lines] instance which has given component c as root and
// features useful for testing.  Potentially occurring errors during the
// usage of a Fixture fatales given testing instance t.
//
// The here created [Lines] instance has a non-blocking
// [Lines.WaitForQuit] method and all its methods triggering events are
// guaranteed to return after the event and subsequently triggered
// events have been processed and the (simulation) screen is
// synchronized.  Or an event triggering method fatales t if given
// duration timeout has passed before all events have been processed.
//
// Testing provides methods for firing user input events like
// [Fixture.FireRune] and retrieving the content of the screen and its
// stylings.  Also user input emulating events do not return before they
// were processed along with subsequently triggered events and all
// prints to the screen have been synchronized within timeout.
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
	ll.Quitting = &quitting{Mutex: &sync.Mutex{}}
	ll.Quitting.AddRune('q')
	tt := &Fixture{
		backend: backend,
		Lines:   ll,
		t:       t,
	}
	tt.Scroll = Scroll{fx: tt}
	backend.Listen(ll.listen)
	return tt
}

// Root returns the initially to the fixture constructor given
// component.  It fatales the test if root is nil.
func (fx *Fixture) Root() Componenter {
	if fx.Lines.scr.lyt.Root == nil {
		fx.t.Fatal("testing: root: layout not initialized")
	}
	return fx.Lines.scr.root().userComponent()
}

func (fx *Fixture) Dim(c Componenter) (d *lyt.Dim) {
	fx.Lines.Update(c, nil, func(e *Env) {
		d = c.layoutComponent().Dim()
	})
	return d
}

// FireResize posts a resize event and returns after this event has been
// processed.  NOTE this event as such is not reported but it triggers
// OnInit and OnLayout events of components which are not initialized or
// whose layout dimensions have changed.
func (fx *Fixture) FireResize(width, height int) *Fixture {
	fx.t.Helper()
	if width == 0 && height == 0 {
		return fx
	}
	fx.PostResize(width, height)
	return fx
}

// FireRune posts given run-key-press event and returns after this event
// has been processed.
func (fx *Fixture) FireRune(r rune, m ...ModifierMask) *Fixture {
	fx.t.Helper()
	if len(m) == 0 {
		fx.PostRune(r, api.ZeroModifier)
	} else {
		fx.PostRune(r, m[0])
	}
	return fx
}

// FireKey posts given special-key event and returns after this
// event has been processed.
func (fx *Fixture) FireKey(k api.Key, m ...ModifierMask) *Fixture {
	fx.t.Helper()
	if len(m) == 0 {
		fx.PostKey(k, api.ZeroModifier)
	} else {
		fx.PostKey(k, m[0])
	}
	return fx
}

// FireKeys for given keys k_0,...,k_n in given fixture fx is a shortcut
// for
//
//	fx.FireKey(k_0, line.ZeroModifier)
//	// ...
//	fx.FireKey(k_n, line.ZeroModifier)
func (fx *Fixture) FireKeys(kk ...api.Key) *Fixture {
	fx.t.Helper()
	for _, k := range kk {
		fx.FireKey(k)
	}
	return fx
}

// FireMove posts a mouse move to given coordinates; an other two given
// ints will be reported as the origin of the mouse move.  Are any given
// coordinates outside the screen area the call is ignored.
func (fx *Fixture) FireMove(x, y int, xy ...int) *Fixture {
	fx.t.Helper()
	if !fx.validCoordinates(x, y, xy...) {
		return fx
	}
	fx.PostMove(x, y, xy...)
	return fx
}

// FireClick posts a first (left) button click at given coordinates and
// returns after this event has been processed.  Are given coordinates
// outside the printable screen area the call is ignored.
func (fx *Fixture) FireClick(x, y int) *Fixture {
	fx.t.Helper()
	if !fx.validCoordinates(x, y) {
		return fx
	}
	fx.PostClick(x, y, api.Primary, api.ZeroModifier)
	return fx
}

// FireContext posts a secondary (right) button click at given coordinates
// and returns after this event has been processed.  Are given
// coordinates outside the screen area the call is ignored.
func (fx *Fixture) FireContext(x, y int) *Fixture {
	fx.t.Helper()
	if !fx.validCoordinates(x, y) {
		return fx
	}
	fx.PostClick(x, y, api.Secondary, api.ZeroModifier)
	return fx
}

func (fx *Fixture) FireDragNDrop(
	x, y int, b ButtonMask, mm ModifierMask, xy ...int,
) *Fixture {
	fx.t.Helper()
	if !fx.validCoordinates(x, y, xy...) {
		return fx
	}
	var dx, dy int
	if len(xy) >= 2 {
		dx, dy = xy[0], xy[1]
	}
	fx.PostDrag(dx, dy, b, mm)(x, y)
	return fx
}

// FireMouse posts a mouse event with provided arguments and returns
// after this event has been processed.  Are given coordinates outside
// the printable screen area the call is ignored.
func (fx *Fixture) FireMouse(
	x, y int, bm api.ButtonMask, mm api.ModifierMask,
) *Fixture {
	fx.t.Helper()
	if !fx.validCoordinates(x, y) {
		return fx
	}

	fx.PostMouse(x, y, bm, mm)
	return fx
}

func (fx *Fixture) validCoordinates(x, y int, xy ...int) bool {
	width, height := fx.Lines.scr.backend.Size()
	if x < 0 || y < 0 || x >= width || y >= height {
		return false
	}
	if len(xy) >= 2 {
		if xy[0] < 0 || xy[1] < 0 || xy[0] >= width || xy[1] >= height {
			return false
		}
	}
	return true
}

// FireComponentClick posts an first (left) button click on given
// relative coordinate in the printable area of given componenter.  Note
// if x or y are outside the printable area of c or c is not part of the
// layout no click will be fired.
func (fx *Fixture) FireComponentClick(c Componenter, x, y int) *Fixture {
	fx.t.Helper()
	if !c.hasLayoutWrapper() {
		return fx
	}
	ox, oy, ok := isInside(c.layoutComponent().wrapped().dim, x, y)
	if !ok {
		return fx
	}
	fx.FireClick(ox+x, oy+y)
	return fx
}

// FireComponentContext posts an second (right) button click on given
// relative coordinate in to printable area of given componenter c.
// Note if x or y are outside the printable area of c or c is not part
// of the layout no click will be fired.
func (fx *Fixture) FireComponentContext(
	c Componenter, x, y int,
) *Fixture {
	fx.t.Helper()
	if !c.hasLayoutWrapper() {
		return fx
	}
	ox, oy, ok := isInside(c.layoutComponent().wrapped().dim, x, y)
	if !ok {
		return fx
	}
	fx.FireContext(ox+x, oy+y)
	return fx
}

func isInside(dim *lyt.Dim, x, y int) (ox, oy int, ok bool) {
	if x < 0 || y < 0 || dim.IsOffScreen() {
		return 0, 0, false
	}
	ox, oy, width, height := dim.Printable()
	if x >= width {
		return 0, 0, false
	}
	if y >= height {
		return 0, 0, false
	}
	return ox, oy, true
}

// ScreenOf provides a string representation of given component's
// screen-area, i.e. without margins and without clippings.  The
// returned StringScreen is nil if given componenter is not part of the
// layout or off-screen.  Note call ScreenArea(c.Dim().Rect()) to get
// the ScreenArea of a component including layout margins.
func (fx *Fixture) ScreenOf(c Componenter) api.StringScreen {
	if !c.hasLayoutWrapper() {
		return nil
	}
	dim := c.layoutComponent().wrapped().Dim()
	if dim.IsOffScreen() {
		return nil
	}
	return fx.ScreenArea(dim.Printable())
}

// ContentOf provides a string representation of given component's
// printable area without gaps.  The returned StringScreen is nil if
// given componenter is not part of the layout or off-screen.
func (fx *Fixture) ContentOf(c Componenter) api.StringScreen {
	if !c.hasLayoutWrapper() {
		return nil
	}
	cmp := c.layoutComponent().wrapped()
	if cmp.Dim().IsOffScreen() {
		return nil
	}
	return fx.ScreenArea(cmp.ContentArea())
}

// CellsOf provides a lines of cells representation of given component's
// printable screen-portion, i.e.  without margins and clippings.  A
// CellsScreen provides next to a string representation also style
// information for each screen coordinate.  The returned CellsScreen is
// nil if given componenter is not part of the layout or off-screen.
// Note call CellsArea(c.Dim().Screen()) to get the screen area of a
// component including margins.
func (fx *Fixture) CellsOf(c Componenter) api.CellsScreen {
	if !c.hasLayoutWrapper() {
		return nil
	}
	dim := c.layoutComponent().wrapped().Dim()
	if dim.IsOffScreen() {
		return nil
	}
	return fx.CellsArea(dim.Printable())
}

type Scroll struct{ fx *Fixture }

// BarDef retrieves given Componenter c's scroll bar definition sparing
// fx.Lines.Update...
func (s Scroll) BarDef(c Componenter) (sbd ScrollBarDef) {
	err := s.fx.Lines.Update(c, nil, func(e *Env) {
		sbd = c.embedded().gg.scrollBarDef
	})
	if err != nil {
		s.fx.t.Fatal(err)
	}
	return sbd
}

// ToBottom scrolls given Componenter c to the bottom sparing
// fx.Lines.Update...
func (s Scroll) ToBottom(c Componenter) {
	err := s.fx.Lines.Update(c, nil, func(e *Env) {
		c.embedded().Scroll.ToBottom()
	})
	if err != nil {
		s.fx.t.Fatal(err)
	}
}
