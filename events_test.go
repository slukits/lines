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

type events struct{ Suite }

func (s *events) SetUp(t *T) {
	t.Parallel()
}

type initFX struct{ Component }

const expInit = "component-fixture initialized"

func (c *initFX) OnInit(e *Env) { fmt.Fprint(e, expInit) }

func (s *events) Initializes_initially_given_component(t *T) {
	ee, tt := Test(t.GoT(), &initFX{})
	ee.Listen()
	t.Eq(expInit, tt.LastScreen.String())
}

type quitCmpFX struct {
	Component
	quitReported bool
}

func (c *quitCmpFX) OnQuit() { c.quitReported = true }

type twoQuittersFX struct {
	Component
	q1, q2 *quitCmpFX
}

func (x *twoQuittersFX) ForStacked(cb func(Componenter) (stop bool)) {
	cb(x.q1)
	cb(x.q2)
}

func (s *events) Reports_quit_key_events_to_all_quitter(t *T) {
	for _, k := range defaultFeatures.keysOf(Quitable) {
		fx := &twoQuittersFX{q1: &quitCmpFX{}, q2: &quitCmpFX{}}
		ee, tt := Test(t.GoT(), fx, -1)
		ee.Listen()
		tt.FireKey(k.Key, k.Mod)
		t.True(fx.q1.quitReported)
		t.True(fx.q2.quitReported)
		t.False(ee.IsListening())
	}
}

type lytCmpFX struct {
	Component
	init, lyt time.Time
}

func (c *lytCmpFX) OnInit(*Env) { c.init = time.Now() }

func (c *lytCmpFX) OnLayout(*Env) { c.lyt = time.Now() }

func (s *events) Reports_layout_after_initialization(t *T) {
	fx := &lytCmpFX{}
	ee, _ := Test(t.GoT(), fx, 2)
	ee.Listen()
	t.True(fx.init.Before(fx.lyt))
}

type updLstCmpFX struct {
	Component
	reported bool
}

func (s *events) Reports_update_to_provided_listener(t *T) {
	fx := &updLstCmpFX{}
	ee, _ := Test(t.GoT(), fx, 1)
	ee.Listen()
	ee.Update(fx, nil, func(_ *Env) { fx.reported = true })
	t.True(fx.reported)
	t.False(ee.IsListening())
}

type updCmpFX struct {
	Component
	reported bool
}

func (c *updCmpFX) OnUpdate(e *Env) { c.reported = true }

func (s *events) Reports_update_without_listener_to_component(t *T) {
	fx := &updCmpFX{}
	ee, _ := Test(t.GoT(), fx, 1)
	ee.Listen()
	ee.Update(fx, nil, nil)
	t.True(fx.reported)
	t.False(ee.IsListening())
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

func (s *events) Reports_moved_focus_gaining_and_loosing(t *T) {
	fx := &stackedCmpFX{cc: []Componenter{&fcsCmpFX{}}}
	ee, _ := Test(t.GoT(), fx, 2)
	ee.MoveFocus(fx.cc[0])
	t.True(fx.lostFocus)
	t.True(fx.cc[0].(*fcsCmpFX).gainedFocus)
	t.False(ee.IsListening())
}

func TestEvents(t *testing.T) {
	t.Parallel()
	Run(&events{}, t)
}
