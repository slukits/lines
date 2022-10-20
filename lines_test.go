// Copyright (c) 2022 Stephan Lukits. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package lines

import (
	"fmt"
	"testing"
	"time"

	. "github.com/slukits/gounit"
)

type _lines struct{ Suite }

func (s *_lines) SetUp(t *T) { t.Parallel() }

type initFX struct{ Component }

const expInit = "component-fixture initialized"

func (c *initFX) OnInit(e *Env) { fmt.Fprint(e, expInit) }

func (s *_lines) tt(t *T, cmp Componenter) *Fixture {
	return TermFixture(t.GoT(), 0, cmp)
}

func (s *_lines) Initializes_initially_given_component(t *T) {
	tt := s.tt(t, &initFX{})
	t.Eq(expInit, tt.Screen().Trimmed().String())
}

func (s *_lines) Reports_quit_key_events_to_all_quitter(t *T) {
	for _, k := range defaultFeatures.keysOf(Quitable) {
		q1, q2 := false, false
		tt := s.tt(t, &cmpFX{})
		tt.Lines.OnQuit(func() { q1 = true })
		tt.Lines.OnQuit(func() { q2 = true })

		tt.FireKey(k.Key, k.Mod)

		t.True(q1)
		t.True(q2)
	}
}

type lytCmpFX struct {
	Component
	init, lyt time.Time
}

func (c *lytCmpFX) OnInit(*Env) { c.init = time.Now() }

func (c *lytCmpFX) OnLayout(*Env) { c.lyt = time.Now() }

func (s *_lines) Reports_layout_after_initialization(t *T) {
	fx := &lytCmpFX{}
	s.tt(t, fx)
	t.True(fx.init.Before(fx.lyt))
}

type updLstCmpFX struct {
	Component
	reported bool
}

func (s *_lines) Reports_update_to_provided_listener(t *T) {
	fx := &updLstCmpFX{}
	tt := s.tt(t, fx)
	t.FatalOn(tt.Lines.Update(fx, nil, func(_ *Env) {
		fx.reported = true
	}))
	t.True(fx.reported)
}

type updCmpFX struct {
	Component
	reported bool
}

func (c *updCmpFX) OnUpdate(e *Env) { c.reported = true }

func (s *_lines) Reports_update_without_listener_to_component(t *T) {
	fx := &updCmpFX{}
	tt := s.tt(t, fx)
	t.FatalOn(tt.Lines.Update(fx, nil, nil))
	t.True(fx.reported)
}

type stackedCmpFX struct {
	Component
	lostFocus bool
	cc        []Componenter
}

func (c *stackedCmpFX) ForStacked(cb func(Componenter) (stop bool)) {
	for _, cmp := range c.cc {
		if !cb(cmp) {
			return
		}
	}
}

func (c *stackedCmpFX) OnFocusLost(*Env) { c.lostFocus = true }

type fcsCmpFX struct {
	Component
	gainedFocus bool
}

func (c *fcsCmpFX) OnFocus(*Env) { c.gainedFocus = true }

func (s *_lines) Reports_moved_focus_gaining_and_loosing(t *T) {
	fx := &stackedCmpFX{cc: []Componenter{&fcsCmpFX{}}}
	tt := s.tt(t, fx)
	t.FatalOn(tt.Lines.Focus(fx.cc[0]))
	t.True(fx.lostFocus)
	t.True(fx.cc[0].(*fcsCmpFX).gainedFocus)
}

func (s *_lines) Ignores_focus_on_focused_component(t *T) {
	fx := &stackedCmpFX{cc: []Componenter{&fcsCmpFX{}}}
	tt := s.tt(t, fx)
	t.FatalOn(tt.Lines.Focus(fx))
	t.Not.True(fx.lostFocus)
	t.Not.True(fx.cc[0].(*fcsCmpFX).gainedFocus)
}

func TestLines(t *testing.T) {
	t.Parallel()
	Run(&_lines{}, t)
}
