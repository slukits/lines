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
	OnLayout(*Env) (reflow bool)
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
	mouseIn layoutComponenter
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

func (s *screen) setRoot(c Componenter, gg *globals) {
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

func (s *screen) setWidth(w int) *screen {
	s.lyt.Width = w
	return s
}

func (s *screen) setHeight(h int) *screen {
	s.lyt.Height = h
	return s
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
	s.syncReflowLayout(ll, nil)
	s.lyt.Root.(layoutComponenter).wrapped().hardSync(s.backend)
	s.lyt.Layers.For(func(l *lyt.Layer) (stop bool) {
		l.Root.(layoutComponenter).wrapped().hardSync(s.backend)
		return false
	})
	s.backend.Redraw()
	s.ensureFocus(ll)
}

// softSync reflows the layout and  hard-syncs every component whose
// layout changed and soft-syncs all remaining dirty components.
// NOTE reflowing the layout is always necessary because we don't know
// if the user added/removed any components.
func (s *screen) softSync(ll *Lines) {
	s.syncReflowLayout(ll, func(c Componenter) {
		c.layoutComponent().wrapped().SetDirty()
	})
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
	s.backend.Update()
	s.ensureFocus(ll)
}

type components []*component

func (cc components) has(c *component) bool {
	for _, _c := range cc {
		if _c == c {
			return true
		}
	}
	return false
}

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
func (s *screen) syncReflowLayout(lines *Lines, cb func(Componenter)) {
	cntx := &rprContext{ll: lines, scr: s}
	for reflow := s.lyt.IsDirty(); reflow; {
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
