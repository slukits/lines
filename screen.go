// Copyright (c) 2022 Stephan Lukits. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package lines

import (
	"github.com/slukits/lines/internal/api"
	"github.com/slukits/lines/internal/lyt"
)

type Dimensions struct {
	dim *lyt.Dim
}

// Screen returns the screen area of a component i.e. its area without
// clippings (including margins).
func (dd Dimensions) Screen() (x, y, width, height int) {
	return dd.dim.Screen()
}

// Printable returns the screen area of a component it can print to i.e.
// its area without margins and without clippings.
func (dd Dimensions) Printable() (x, y, width, height int) {
	return dd.dim.Printable()
}

// Width of component without margins.
func (dd Dimensions) Width() int { return dd.dim.Width() }

// Height of component without margins.
func (dd Dimensions) Height() int { return dd.dim.Height() }

// DD is a function which extracts layout information from a component
// without the need of being inside one of the components
// event-callbacks.
type DD func(Componenter) Dimensions

// Layouter is implemented by components which want to be notified if
// their layout has changed.
type Layouter interface {

	// OnLayout is called after the layout manager has changed the
	// screen area of a component.
	OnLayout(*Env) (reflow bool)
}

// AfterLayouter is implemented by components who want to be notified
// after all OnLayout-events have been reported and processed.
// Typically these are components nesting other components, i.e. Stacker
// and Chainer, who want to adjust their layout after their nested
// components have finished their layout settings.
type AfterLayouter interface {

	// OnAfterLayout implementations are called after all OnLayout
	// events have been reported and processed.  Since the typical
	// use-case is to adjust the layout of stacker and chainer according
	// to the layout of their nested components a "dimensions function"
	// DD is passed along which allows to extract layout dimensions for
	// an arbitrary component of the layout.
	OnAfterLayout(*Env, DD) (reflow bool)
}

// Stacker is implemented by components which want to provide nested
// components which are vertically stacked in the layout.
type Stacker interface {

	// ForStacked calls back for each component of this Stacker
	// respectively until the callback asks to stop.
	ForStacked(func(Componenter) (stop bool))
}

// Chainer is implemented by components which want to provided nested
// components which are horizontally chained in the layout.
type Chainer interface {

	// ForChained calls back for each component of this Chainer
	// respectively until the callback asks to stop.
	ForChained(func(Componenter) (stop bool))
}

type cursor [3]int

// X returns set x coordinate of given cursor c's position which is -1
// if the cursor is not shown.
func (c cursor) X() int { return c[0] }

// Y returns set y coordinate of given cursor c's position which is -1
// if the cursor is not shown.
func (c cursor) Y() int { return c[1] }

// Style returns given cursor c's style which is the ZeroCursor if c is
// not shown.
func (c *cursor) Style() CursorStyle { return CursorStyle((*c)[2]) }

// Resets sets given cursor c's x and y coordinate to -1 and its style
// to ZeroCursor.
func (c *cursor) Reset() *cursor {
	(*c)[0], (*c)[1], (*c)[2] = -1, -1, int(ZeroCursor)
	return c
}

// Removed returns true if given cursor c's x coordinate is set to -1
func (c cursor) Removed() bool { return c[0] == -1 }

// Coordinates returns set x and y coordinates which are -1 if given
// cursor c is not shown.
func (c *cursor) Coordinates() (x, y int) {
	x, y = (*c)[0], (*c)[1]
	return
}

// set sets given cursor c's position given coordinates x and y as well
// as its cursor style given cursor style cs if later is not the
// ZeroCursor.  Note if x or y is negative c is reset otherwise x and y
// are set while cs is only set if not the ZeroCursor.
func (c *cursor) set(x, y int, cs CursorStyle) *cursor {
	if x < 0 || y < 0 {
		c.Reset()
		return c
	}
	(*c)[0], (*c)[1] = x, y
	if cs != ZeroCursor {
		(*c)[2] = int(cs)
	}
	return c
}

type screen struct {
	lyt     *lyt.Manager
	backend api.Displayer
	focus   layoutComponenter
	mouseIn layoutComponenter
	cursor  *cursor
}

func newScreen(backend api.UIer, cmp Componenter, gg *Globals) *screen {
	scr := &screen{backend: backend}
	if cmp == nil {
		cmp = &Component{}
	}
	gg.scr = scr
	lc := cmp.initialize(cmp, backend, gg.clone())
	lc.wrapped().ensureFeatures()
	scr.lyt = &lyt.Manager{Root: lc}
	scr.focus = lc
	scr.cursor = (&cursor{}).Reset()
	return scr
}

func (s *screen) setRoot(c Componenter, gg *Globals) {
	if c == nil {
		return
	}
	lc := c.initialize(c, s.backend.(api.UIer), gg.clone())
	lc.wrapped().ensureFeatures()
	s.lyt.SetRoot(lc)
	s.focus = lc
}

func (s *screen) root() *component {
	return s.lyt.Root.(layoutComponenter).wrapped()
}

func (s *screen) setSize(w, h int, ll *Lines) func() {
	s.lyt.Width, s.lyt.Height = w, h
	x, y := s.cursor.Coordinates()
	if x == -1 || y == -1 {
		return nil
	}
	lc := s.cursorComponent()
	s.setCursor(-1, -1, ZeroCursor)
	if lc == nil {
		return nil
	}
	if !lyt.NewRect(lc.wrapped().ContentArea()).Has(x, y) {
		usr := lc.userComponent()
		if c, ok := usr.(Cursorer); ok {
			usr.enable()
			callback(lc.userComponent(), &rprContext{ll: ll, scr: s}, func(c Cursorer) func(e *Env) {
				return func(e *Env) { c.OnCursor(e, false) }
			}(c))
		}
		return nil
	}
	cx, cy, _, _ := lc.wrapped().ContentArea()
	return func(ll *Lines, lc layoutComponenter, x, y int) func() {
		return func() {
			if !s.lyt.Has(lc, nil) || lc.Dim().IsOffScreen() {
				return
			}
			// check if the content area got smaller, i.e. doesn't
			// include the cursor any more.
			_, _, cw, ch := lc.wrapped().ContentArea()
			absOnly := true
			if x >= cw {
				x = cw - 1
				absOnly = false
			}
			if y >= ch {
				y = ch - 1
				absOnly = false
			}
			lc.wrapped().setCursor(y, x)
			usr := lc.userComponent()
			if c, ok := usr.(Cursorer); ok {
				usr.enable()
				callback(
					lc.userComponent(),
					&rprContext{ll: ll, scr: s},
					func(c Cursorer) func(e *Env) {
						return func(e *Env) { c.OnCursor(e, absOnly) }
					}(c))
			}
			s.softSync(ll)
		}
	}(ll, lc, x-cx, y-cy)
}

// setCursor calls the backend to set the cursor at given coordinates
// (x,y) having optionally given style cs[0].  Note screen keeps a local
// copy of the current cursor state to determine the component having
// the cursor set (which must not necessarily be the focused component).
func (s *screen) setCursor(x, y int, cs ...CursorStyle) {
	s.cursor.set(s.backend.SetCursor(x, y, cs...))
}

func (s *screen) cursorComponent() layoutComponenter {
	if s.cursor.Removed() {
		return nil
	}
	path, err := s.lyt.LocateAt(s.cursor.Coordinates())
	if err != nil || len(path) == 0 {
		return nil
	}
	return path[len(path)-1].(layoutComponenter)
}

// forBaseComponents calls back for the components of the base layer.
func (s *screen) forBaseComponents(cb func(layoutComponenter)) {
	s.lyt.ForDimer(nil, func(d lyt.Dimer) (stop bool) {
		cb(d.(layoutComponenter))
		return false
	})
}

type dimerIterator interface {
	ForDimer(lyt.Dimer, func(lyt.Dimer) bool) *lyt.Layers
}

// forComponent calls back for all components including components of
// layers whereas the layers are processed in the order they where
// provided to the layout manager.
func (s *screen) forComponent(cb func(layoutComponenter)) {
	var recurseOverComponents func(dimerIterator, func(layoutComponenter))
	recurseOverComponents = func(
		m dimerIterator, cb func(layoutComponenter),
	) {
		ll := m.ForDimer(nil, func(d lyt.Dimer) (stop bool) {
			cb(d.(layoutComponenter))
			return false
		})
		ll.For(func(l *lyt.Layer) (stop bool) {
			recurseOverComponents(l, cb)
			return
		})
	}
	recurseOverComponents(s.lyt, cb)
}

func (s *screen) forUninitialized(cb func(Componenter)) {
	s.forComponent(func(cmp layoutComponenter) {
		if cmp.userComponent().isInitialized() {
			return
		}
		cb(cmp.userComponent())
	})
}

// forFocused calls back for the focused components and all its
// ancestors in the layout.
func (s *screen) forFocused(cb func(layoutComponenter) (stop bool)) {
	s.forBubbling(s.focus, cb)
}

// forBubbling calls back for given layouter and all its ancestors in
// the layout.
func (s *screen) forBubbling(
	lc layoutComponenter, cb func(layoutComponenter) (stop bool),
) {
	path, err := s.lyt.Locate(lc)
	if err != nil {
		return
	}
	if cb(lc) {
		return
	}
	for i := len(path) - 1; i >= 0; i-- {
		if cb(path[i].(layoutComponenter)) {
			return
		}
	}
}

// hardSync reflows the layout and hard-syncs every component.
func (s *screen) hardSync(ll *Lines) {
	s.syncReflowLayout(ll, true, nil)
	s.syncAfterLayout(ll)
	s.lyt.Root.(layoutComponenter).wrapped().hardSync(s.backend)
	s.lyt.Layers.For(func(l *lyt.Layer) (stop bool) {
		l.Root.(layoutComponenter).wrapped().hardSync(s.backend)
		return false
	})
	s.backend.Redraw()
	s.ensureFocus(ll)
}

func dimensionsOf(cmp Componenter) Dimensions {
	return Dimensions{
		dim: cmp.embedded().layoutCmp.wrapped().dim}
}

// softSync reflows the layout and  hard-syncs every component whose
// layout changed and soft-syncs all remaining dirty components.
// NOTE reflowing the layout is always necessary because we don't know
// if the user added/removed any components.
func (s *screen) softSync(ll *Lines) {
	s.syncReflowLayout(ll, false, func(c Componenter) {
		c.layoutComponent().wrapped().SetDirty()
	})
	update := s.backend.Update
	s.syncAfterLayout(ll)
	// if s.syncAfterLayout(ll) {
	// 	s.lyt.Root.(layoutComponenter).wrapped().hardSync(s.backend)
	// 	s.lyt.Layers.For(func(l *lyt.Layer) (stop bool) {
	// 		l.Root.(layoutComponenter).wrapped().hardSync(s.backend)
	// 		return false
	// 	})
	// 	update = s.backend.Redraw
	// }
	if !s.lyt.Root.(layoutComponenter).wrapped().IsDirty() &&
		s.lyt.Layers == nil {
		s.ensureFocus(ll)
		return
	}
	if s.lyt.Root.(layoutComponenter).wrapped().IsDirty() {
		s.lyt.Root.(layoutComponenter).wrapped().sync(s.backend)
	}
	s.lyt.Layers.For(func(l *lyt.Layer) (stop bool) {
		l.Root.(layoutComponenter).wrapped().hardSync(s.backend)
		return false
	})
	update()
	s.ensureFocus(ll)
}

func (s *screen) syncAfterLayout(ll *Lines) bool {
	cntx, reflow := &rprContext{ll: ll, scr: s}, false
	s.lyt.ForDimer(nil, func(d lyt.Dimer) (stop bool) {
		cmp := d.(layoutComponenter).userComponent()
		if al, ok := cmp.(AfterLayouter); ok {
			callback(cmp, cntx, func(e *Env) {
				if reflow {
					al.OnAfterLayout(e, dimensionsOf)
					return
				}
				reflow = al.OnAfterLayout(e, dimensionsOf)
			})
		}
		return
	})
	if reflow {
		s.lyt.Reflow(nil)
		return true
	}
	return false
}

type components []*component

type layereder interface {
	Layered() layoutComponenter
}

func (s *screen) ensureFocus(ll *Lines) {
	if modal := s.haveModal(); modal != nil {
		if modal == s.focus {
			return
		}
		moveFocus(modal.userComponent(), &rprContext{ll: ll, scr: s})
		return
	}
	if s.lyt.Has(s.focus, nil) || s.lyt.Layers.Have(s.focus) {
		return
	}
	moveFocus(
		s.lyt.Root.(layoutComponenter).userComponent(),
		&rprContext{ll: ll, scr: s},
	)
}

func (s *screen) haveModal() (lc layoutComponenter) {
	s.lyt.Layers.ForReversed(func(l *lyt.Layer) (stop bool) {
		_, ok := l.Root.(layoutComponenter).userComponent().(Modaler)
		if !ok {
			return false
		}
		lc = l.Root.(layoutComponenter)
		return true
	})
	return lc
}

// syncReflowLayout reflows the layout and reports to every component
// with changed layout implementing Layouter.  It also calls back for
// every component with changed layout if callback not nil.
func (s *screen) syncReflowLayout(
	lines *Lines, hard bool, cb func(Componenter),
) {
	cntx, count := &rprContext{ll: lines, scr: s}, 0
	reflow := s.lyt.IsDirty()
	if hard {
		reflow = true
	}
	for reflow && count < 10 {
		reflow = false
		ll := s.lyt.Layers
		s.lyt.Reflow(func(d lyt.Dimer) {
			cmp := d.(layoutComponenter).userComponent()
			if lyt, ok := cmp.(Layouter); ok {
				callback(cmp, cntx, func(e *Env) {
					reflow = lyt.OnLayout(e)
				})
			}
			if cb != nil {
				cb(cmp)
			}
		})
		if cb != nil {
			s.lyt.Layers.BaseDimerLayeredByReMovedLayers(
				s.lyt, ll, func(d lyt.Dimer) {
					cb(d.(layoutComponenter).userComponent())
				})
		}
		if reflow {
			count++
			reportInit(lines, s)
		}
	}
}

// Stacking embedded in a component makes the component implement the
// Stacker interface.  Typically the Componenter slice is filled in a
// component's [Initer]-listener:
//
//	type stackedCmp struct { lines.Component }
//
//	type myCmp struct{
//		lines.Component
//		lines.Stacking
//	}
//
//	func (c *myCmp) OnInit(_ *lines.Env) {
//		for i := 0; i < 3; i++ {
//			c.CC = append(c.CC, &stackedCmp{})
//		}
//	}
type Stacking struct{ CC []Componenter }

// ForStacked calls back for each component of this Stacker respectively
// until the callback asks to stop.
func (s Stacking) ForStacked(cb func(Componenter) (stop bool)) {
	for _, c := range s.CC {
		if cb(c) {
			return
		}
	}
}

// Chaining embedded in a component makes the component implement the
// Chainer interface.  Typically the Componenter slice CC is filled in a
// component's OnInit-listener:
//
//	type chainedCmp struct { lines.Component }
//
//	type myCmp struct{
//		lines.Component
//		lines.Chaining
//	}
//
//	func (c *myCmp) OnInit(_ *lines.Env) {
//		for i := 0; i < 3; i++ {
//			c.CC = append(c.CC, &chainedCmp{})
//		}
//	}
type Chaining struct {

	// CC holds the chained components
	CC []Componenter
}

// ForChained calls back for each component of this Chainer respectively
// until the callback asks to stop.
func (cg Chaining) ForChained(cb func(Componenter) (stop bool)) {
	for _, c := range cg.CC {
		if cb(c) {
			return
		}
	}
}
