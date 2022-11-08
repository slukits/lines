// Copyright (c) 2022 Stephan Lukits. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package lines

import (
	"github.com/slukits/lines/internal/api"
	"github.com/slukits/lines/internal/lyt"
)

// Layouter is implemented by components which want to be notified if
// their layout has changed.
type Layouter interface {

	// OnLayout is called after the layout manager has changed the
	// screen area of a component.
	OnLayout(*Env)
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

type screen struct {
	lyt     *lyt.Manager
	backend api.Displayer
	focus   layoutComponenter
}

func newScreen(backend api.UIer, cmp Componenter, gg *globals) *screen {
	scr := &screen{backend: backend}
	if cmp == nil {
		cmp = &Component{}
	}
	lc := cmp.initialize(cmp, backend, gg.clone())
	lc.wrapped().ensureFeatures()
	scr.lyt = &lyt.Manager{Root: lc}
	scr.focus = lc
	return scr
}

func (s *screen) root() *component {
	return s.lyt.Root.(layoutComponenter).wrapped()
}

func (s *screen) setWidth(w int) *screen {
	s.lyt.Width = w
	return s
}

func (s *screen) setHeight(h int) *screen {
	s.lyt.Height = h
	return s
}

func (s *screen) forComponent(cb func(layoutComponenter)) {
	s.lyt.ForDimer(nil, func(d lyt.Dimer) (stop bool) {
		cb(d.(layoutComponenter))
		return false
	})
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
	s.syncReLayout(ll, nil)
	s.lyt.ForDimer(nil, func(d lyt.Dimer) (stop bool) {
		wrapped := d.(layoutComponenter).wrapped()
		wrapped.hardSync(s.backend)
		return false
	})
	s.backend.Redraw()
}

// softSync reflows the layout and  hard-syncs every component whose
// layout changed and soft-syncs all remaining dirty components.
// NOTE reflowing the layout is always necessary because we don't know
// if the user added any new components.
func (s *screen) softSync(ll *Lines) {
	s.syncReLayout(ll, func(cmp Componenter) {
		cmp.layoutComponent().wrapped().hardSync(s.backend)
	})
	s.syncDirty()
	s.backend.Update()
}

// syncReLayout reflows the layout and reports to every component with
// changed layout implementing Layouter.  It also calls back for every
// component with changed layout if callback not nil.
func (s *screen) syncReLayout(ll *Lines, cb func(Componenter)) {
	if s.lyt.IsDirty() {
		reported := false
		cntx := &rprContext{ll: ll, scr: s}
		s.lyt.Reflow(func(d lyt.Dimer) {
			cmp := d.(layoutComponenter).userComponent()
			if lyt, ok := cmp.(Layouter); ok {
				callback(cmp, cntx, lyt.OnLayout)
				if !reported {
					reported = true
				}
			}
			if cb != nil {
				cb(cmp)
			}
		})
	}
}

// syncDirty syncs component with updated content.
func (s *screen) syncDirty() {
	s.lyt.ForDimer(nil, func(d lyt.Dimer) (stop bool) {
		// cmp := d.(layoutComponenter).userComponent()
		wrapped := d.(layoutComponenter).wrapped()
		if wrapped.IsDirty() {
			wrapped.sync(s.backend)
		}
		return false
	})
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
type Chaining struct{ CC []Componenter }

// ForChained calls back for each component of this Chainer respectively
// until the callback asks to stop.
func (cg Chaining) ForChained(cb func(Componenter) (stop bool)) {
	for _, c := range cg.CC {
		if cb(c) {
			return
		}
	}
}
