// Copyright (c) 2022 Stephan Lukits. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package lines

import (
	"github.com/gdamore/tcell/v2"
	"github.com/slukits/lines/internal/lyt"
)

// Layouter is implemented by components which want to be notified if
// their layout has changed.
type Layouter interface {

	// OnLayout is called after the layout manager has changed the
	// available screen area of a component.
	OnLayout(*Env)
}

// Stacker is implemented by components which want to provide nested
// components in a vertical manner.
type Stacker interface {

	// ForStacked calls back for each component of this Stacker
	// respectively until the callback asks to stop.
	ForStacked(func(Componenter) (stop bool))
}

// Chainer is implemented by components which want to provided nested
// components in a horizontal manner.
type Chainer interface {

	// ForChained calls back for each component of this Chainer
	// respectively until the callback asks to stop.
	ForChained(func(Componenter) (stop bool))
}

type screen struct {
	lyt   *lyt.Manager
	lib   tcell.Screen
	focus layoutComponenter
}

// zeroComponent is set a Component's component after a listener of its
// componenter returns
var zeroComponent *component

// newSim returns a new Screen instance wrapping tcell's simulation
// screen for testing purposes.
func newSim(cmp Componenter) (*screen, error) {
	lib := screenFactory.NewSimulationScreen("")
	if err := lib.Init(); err != nil {
		return nil, err
	}
	scr := &screen{lib: lib}
	if cmp != nil {
		lc := cmp.initialize(cmp)
		lc.wrapped().ensureFeatures()
		scr.lyt = &lyt.Manager{Root: lc}
		scr.focus = lc
	} else {
		cmp := &Component{}
		lc := cmp.initialize(cmp)
		lc.wrapped().ensureFeatures()
		scr.lyt = &lyt.Manager{Root: lc}
		scr.focus = lc
	}
	return scr, nil
}

func (s *screen) root() *component {
	return s.lyt.Root.(layoutComponenter).wrapped()
}

func (s *screen) setWidth(w int) *screen {
	s.lyt.Root.Dim().SetWidth(w)
	return s
}

func (s *screen) setHeight(h int) *screen {
	s.lyt.Root.Dim().SetHeight(h)
	return s
}

func (s *screen) forComponent(cb func(Componenter)) {
	s.lyt.ForDimer(nil, func(d lyt.Dimer) (stop bool) {
		cb(d.(layoutComponenter).userComponent())
		return false
	})
}

func (s *screen) forUninitialized(cb func(Componenter)) {
	s.forComponent(func(cmp Componenter) {
		if cmp.isInitialized() {
			return
		}
		cb(cmp)
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

func (s *screen) hardSync(ee *Events) {
	s.syncReLayouted(ee)
	s.syncStillDirty()
	s.lib.Sync()
}

func (s *screen) softSync(ee *Events) {
	s.syncReLayouted(ee)
	s.syncStillDirty()
	s.lib.Show()
}

func (s *screen) syncReLayouted(ee *Events) {
	if s.lyt.IsDirty() {
		reported := false
		s.lyt.Reflow(func(d lyt.Dimer) {
			cmp := d.(layoutComponenter).userComponent()
			if lyt, ok := cmp.(Layouter); ok {
				cmp.enable()
				env := &Env{cmp: cmp, EE: ee, Evt: nil}
				lyt.OnLayout(env)
				env.reset()
				cmp.disable()
				if !reported {
					reported = true
				}
			}
			cmp.sync(s.lib)
		})
		if reported {
			reportReported(ee)
		}
	}
}

func (s *screen) syncStillDirty() {
	s.lyt.ForDimer(nil, func(d lyt.Dimer) (stop bool) {
		cmp := d.(layoutComponenter).userComponent()
		if cmp.isDirty() {
			cmp.sync(s.lib)
		}
		return false
	})
}

// screenFactory is used to create new tcell-screens for production or
// for simulation.  export_test.go makes it possible to replace this
// screen factory with a screen-factory mocking up tcell's screen
// creation errors so they can be tested.
var screenFactory screenFactoryer = &defaultFactory{}

type defaultFactory struct{}

func (f *defaultFactory) NewScreen() (tcell.Screen, error) {
	return tcell.NewScreen()
}

func (f *defaultFactory) NewSimulationScreen(
	s string,
) tcell.SimulationScreen {
	return tcell.NewSimulationScreen(s)
}

type screenFactoryer interface {
	NewScreen() (tcell.Screen, error)
	NewSimulationScreen(string) tcell.SimulationScreen
}

// Stacking embedded in a component makes the component implement the
// Stacker interface.  Typically the Componenter slice is filled in a
// component's OnInit-listener.
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
// Chainer interface.  Typically the Componenter slice is filled in a
// component's OnInit-listener.
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
