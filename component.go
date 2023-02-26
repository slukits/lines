// Copyright (c) 2022 Stephan Lukits. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package lines

import (
	"github.com/slukits/lines/internal/api"
	"github.com/slukits/lines/internal/lyt"
)

// Component enables a user implemented UI-component to be processed by
// lines, i.e. all user UI-components which are provided to lines must
// embed this type.  NOTE accessing features of an embedded Component is
// only supported during an event reporting callback to embedding
// component c:
//
//	func (c *cmp) OnInit(_ *lines.Env) {
//	    go func() {
//	        time.Sleep(1*time.Second)
//	        c.Dim().SetHeight(5) // will panic
//	    }
//	    c.Dim().SetHeight(5) // will not panic
//	}
//
// Note that there are three rectangles on the screen associated with
// component:
//
//	c.Dim().Rect() // the component's "screen area" including margins
//	c.Dim().Area() // the "printable area", i.e. without margins
//	c.ContentArea() // the "content area", i.e. without margins and gaps
//
// Next to embedding the Component type a client component will usually
// also implement event listener interfaces to receive events like:
//   - [Initer] is informed once before a component becomes part of the layout
//   - [AfterIniter] is informed once after all components are initialize
//     before the layout is calculated
//   - [Layouter] is informed that a component's layout was calculated
//   - [Focuser]/[FocusLooser] is informed about focus gain/loss
//   - [Keyer] is informed about any user special key-press like 'enter' or 'tab'
//   - [Runer] is informed about user rune-key input
//   - [Mouser] is informed about any mouse event see also [Clicker]/[Contexter]
//   - [LineSelecter] is informed if a component's line was selected
//   - [LineFocuser] is informed if a component's line received the focus
type Component struct {

	// FF provides access and fine grained control over a components
	// default behavior.
	FF *Features

	// Register provides the api to register keys and runes listeners
	// for a component.  See Keyer and Runer for a more general way of
	// keyboard listening.
	Register *Listeners

	// Edit provides a component's API to control editing its content.
	Edit *Editor

	// bcknd to post Update and Focus events
	bcknd api.UIer

	// component provides properties/features of an Component.  A
	// Component can't do it directly if it should panic if it is used
	// outside an event reporting callback.  layoutCmp wraps the actual
	// component instance hence it is available for internal use and if
	// this Component is enabled its *components points to the
	// layoutCmp wrapped component.  If it is disabled *component is set
	// to nil.
	*component

	// layoutCmp wraps a *component-instance in a way that it can be
	// used by the layout-manager and that it provides user defined
	// components.  I.e. the problem here is that an embedded Component
	// instance might be enabled or not (*component is nil or not).
	// Thus we can't just pass on a user defined component to the layout
	// manager because its access to the components dimensions which are
	// provided by the wrapped *component might panic.  On the other
	// hand  we need to provide the user provided components
	// (ForStacked, ForChained) to the layout manager.  A
	// layoutComponenter combines a Component's wrapped (non nil)
	// *component with the ForStacked/ForChained method of the user
	// defined associated component.  Emitting other layoutComponenter
	// to the layout manager.
	layoutCmp layoutComponenter
}

type ComponentMode uint

const (

	// Overwriting a components content by an write operation.
	Overwriting ComponentMode = 1 << iota

	// Appending to a components content by an write operation.
	Appending

	// Tailing is appending and displaying the contents "tail"
	// especially if the display area cannot show all the content.
	Tailing
)

func (c *Component) initialize(
	userComponent Componenter, backend api.UIer, gg *Globals,
) layoutComponenter {

	if c.layoutCmp != nil { // already initialized
		return c.layoutCmp
	}

	c.bcknd = backend

	inner := &component{
		dim:     lyt.DimFilling(1, 1),
		ll:      &lines{},
		gg:      gg,
		userCmp: userComponent,
		mod:     Overwriting,
		dirty:   true,
	}
	c.FF = &Features{c: c}
	c.Register = &Listeners{c: c}
	inner.Scroll = &Scroller{c: c, bar: -1}
	inner.LL = newComponentLines(c)
	inner.gg.SetUpdateListener(cmpGlobalsClosure(inner))
	switch userComponent.(type) {
	case Stacker:
		c.layoutCmp = &stackingWrapper{component: inner}
	case Chainer:
		c.layoutCmp = &chainingWrapper{component: inner}
	default:
		c.layoutCmp = inner
	}
	return c.layoutCmp
}

func cmpGlobalsClosure(
	c *component,
) func(globalsUpdates, StyleType, globalStyleUpdates) {
	return func(gu globalsUpdates, st StyleType, gsu globalStyleUpdates) {
		if !c.dirty {
			c.dirty = true
		}
	}
}

func (c *Component) backend() api.UIer { return c.bcknd }

// enable component for client usage.
func (c *Component) enable() {
	if c.component != nil {
		return
	}
	c.component = c.layoutCmp.wrapped()
}

// disable component for client usage.
func (c *Component) disable() {
	if c.component == nil {
		return
	}
	c.component = nil
}

// isEnabled returns true if given Component c's internal component is
// set; false otherwise.
func (c *Component) isEnabled() bool {
	return c.component != nil
}

// isInitialized returns true if embedded component-instance is wrapped
// in a layout component and has been initialized.
func (c *Component) isInitialized() bool {
	if c.layoutCmp == nil {
		return false
	}
	return c.layoutCmp.wrapped().initialized
}

func (c *Component) hasLayoutWrapper() bool {
	return c.layoutCmp != nil
}

func (c *Component) layoutComponent() layoutComponenter {
	return c.layoutCmp
}

// isNesting returns true if the component is stacking or chaining other
// components.
func (c *Component) isNesting() bool {
	if !c.isInitialized() {
		return false
	}
	switch c.layoutCmp.(type) {
	case *stackingWrapper:
		return true
	case *chainingWrapper:
		return true
	}
	return false
}

func (c *Component) embedded() *Component { return c }

// Gaps returns a gaps writer at given leven allowing to do framing,
// padding or guttering around a component's content:
//
//	fmt.Fprint(c.Gaps(0).AA(Reverse).Filling(), "•")
//	fmt.Fprint(c.Gaps(0).Corners.AA(Reverse), "•")
//	c.Gaps(1).AA(Reverse)
//	c.Gaps(1).Corners.AA(Reverse)
func (c *Component) Gaps(level int) *GapsWriter {
	if c.gaps == nil {
		c.gaps = newGaps(c.gg.Style(Default))
	}
	return newGapsWriter(level, c.gaps)
}

// Globals provides access to the API for manipulating component c
// specific globally inherited properties like tab-width.  Note to
// change such a property globally use the [Lines]-instance ll which
// layouts c.  ll's Globals-property provides the same Api but
// propagates manipulations to all components of the layout.
func (c *Component) Globals() *Globals { return c.gg }

// SetCursor of given component c to given line and column with
// optionally given cursor style within c's content area.  I.e.
// line=0, column=0 refers to the first column in the first content-line
// while gaps and margins are ignored.  Note setting the cursor is only
// effective after the first layout.  SetCursor is a no-op for stacker
// and chainer.  [Lines.SetCursor] allows for absolute cursor
// positioning whereas there the cursor is set to (x,y) coordinates,
// i.e. it has switched argument order compared to (line,column).
func (c *Component) SetCursor(
	line, column int, cs ...CursorStyle,
) *Component {
	c.setCursor(line, column, cs...)
	return c
}

// CursorPosition returns relative to given component c's content origin
// the line and column index of the cursor in the content area and true
// if c has the cursor set; otherwise -1, -1 and false is return.
func (c *Component) CursorPosition() (line, column int, _ bool) {
	return c.cursorPosition()
}

// component is the actual implementation of a lines-Component.
type component struct {
	userCmp     Componenter
	gg          *Globals
	dim         *lyt.Dim
	mod         ComponentMode
	initialized bool
	ll          *lines

	// LL provides an API for ui-aspects of a component's lines like is
	// a line focusable, are they tailing maintained etc.  To manipulate
	// the content of component lines use an Env(ironment) instance of a
	// reported event.
	LL *ComponentLines

	// Scroll provides a component's API for scrolling.
	Scroll *Scroller

	lst                *listeners
	ff                 *features
	gaps               *gaps
	dirty, cursorMoved bool

	Src *ContentSource

	// _first holds the content line index of the _first displayed line
	_first int

	// slctd hold the index of the currently selected line
	slctd int
}

// component gets the component out of a layoutComponenter without using
// a type-switch.
func (c *component) wrapped() *component { return c }

func (c *component) globals() *Globals {
	if c == nil {
		return nil
	}
	return c.gg
}

// cursorPosition returns relative to given component c's content origin
// the line and column index of the cursor in the content area and true
// if c has the cursor set; otherwise -1, -1 and false is return.
func (c *component) cursorPosition() (line, column int, _ bool) {
	if c.gg == nil {
		return -1, -1, false
	}
	if c.gg.scr.cursor.Removed() {
		return -1, -1, false
	}
	x, y := c.gg.scr.cursor.Coordinates()
	cx, cy, cw, ch := c.ContentArea()
	x -= cx
	y -= cy
	if x < 0 || x >= cw || y < 0 || y >= ch {
		return -1, -1, false
	}
	return y, x, true
}

func (c *component) setCursor(
	line, column int, cs ...CursorStyle,
) {
	if _, ok := c.userCmp.layoutComponent().(lyt.Stacker); ok {
		return
	}
	if _, ok := c.userCmp.layoutComponent().(lyt.Chainer); ok {
		return
	}
	if line < 0 || column < 0 {
		c.gg.setCursor(line, column)
		if !c.cursorMoved {
			c.cursorMoved = true
		}
		return
	}

	x, y, w, h := c.ContentArea()
	if line >= h || column >= w {
		return
	}
	c.gg.setCursor(y+line, x+column, cs...)
	if !c.cursorMoved {
		c.cursorMoved = true
	}
}

func (c *component) userComponent() Componenter {
	if c == nil {
		return nil
	}
	return c.userCmp
}

func (c *component) setInitialized() {
	c.initialized = true
}

func (c *component) ensureFeatures() *features {
	if c.ff != nil {
		return c.ff
	}
	c.ff = &features{}
	return c.ff
}

func (c *component) ensureListeners() {
	if c.lst != nil {
		return
	}
	c.lst = &listeners{}
}

// layoutComponenter combines the user-provided component with its
// internally created component to have all information to layout the
// component together: Dimer-implementation is provide by the internally
// created component while the user-component potentially provides nested
// user-components.
type layoutComponenter interface {
	lyt.Dimer
	wrapped() *component
	userComponent() Componenter
}

// stackingWrapper wraps a stacking user-component for the layout
// manager.  Avoiding panics on Gaps- or Dim-access through the layout
// manager
type stackingWrapper struct{ *component }

func (sw *stackingWrapper) Gaps() api.Gaps {
	if sw.gaps == nil {
		return api.Gaps{}
	}
	return api.Gaps{
		Top:    len(sw.gaps.top.ll),
		Right:  len(sw.gaps.right.ll),
		Bottom: len(sw.gaps.bottom.ll),
		Left:   len(sw.gaps.left.ll),
	}
}

func (sw *stackingWrapper) ForStacked(cb func(lyt.Dimer) bool) {
	sw.userCmp.(Stacker).ForStacked(func(cmp Componenter) bool {
		if !cmp.hasLayoutWrapper() {
			cmp.initialize(
				cmp,
				sw.userCmp.backend(),
				sw.globals().clone(),
			)
			if sw.ff.all() != NoFeature {
				cmp.embedded().layoutCmp.wrapped().ff = sw.ff.copy()
			}
		}
		return cb(cmp.layoutComponent())
	})
}

// chainingWrapper wraps a chaining user-component for the layout
// manager.  Avoiding panics on Gaps- or Dim-access through the layout
// manager
type chainingWrapper struct{ *component }

func (sw *chainingWrapper) Gaps() api.Gaps {
	if sw.gaps == nil {
		return api.Gaps{}
	}
	return api.Gaps{
		Top:    len(sw.gaps.top.ll),
		Right:  len(sw.gaps.right.ll),
		Bottom: len(sw.gaps.bottom.ll),
		Left:   len(sw.gaps.left.ll),
	}
}

func (cw *chainingWrapper) ForChained(cb func(lyt.Dimer) bool) {
	cw.userCmp.(Chainer).ForChained(func(cmp Componenter) bool {
		if !cmp.hasLayoutWrapper() {
			cmp.initialize(
				cmp,
				cw.userCmp.backend(),
				cw.globals().clone(),
			)
			if cw.ff.all() != NoFeature {
				cmp.embedded().layoutCmp.wrapped().ff = cw.ff.copy()
			}
		}
		return cb(cmp.layoutComponent())
	})
}
