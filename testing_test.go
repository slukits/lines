// Copyright (c) 2022 Stephan Lukits. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package lines

import (
	"testing"

	"github.com/gdamore/tcell/v2"
	. "github.com/slukits/gounit"
)

type _Testing struct{ Suite }

func (s *_Testing) SetUp(t *T) { t.Parallel() }

func (s *_Testing) Starts_non_blocking_listening_with_listen_call(t *T) {
	ee, _ := Test(t.GoT(), nil)
	t.False(ee.IsListening())
	ee.Listen()
	t.True(ee.IsListening())
	ee.QuitListening()
	t.False(ee.IsListening())
}

func (s *_Testing) Starts_listening_if_a_resize_is_fired(t *T) {
	ee, tt := Test(t.GoT(), nil)
	defer ee.QuitListening()
	t.False(ee.IsListening())
	t.True(tt.FireResize(22, 42).IsListening())
}

func (s *_Testing) Starts_listening_if_a_key_is_fired(t *T) {
	ee, tt := Test(t.GoT(), nil, -1) // listen for ever
	defer ee.QuitListening()
	t.False(ee.IsListening())
	t.True(tt.FireKey(tcell.KeyBS, 0).IsListening())
}

func (s *_Testing) Starts_listening_if_a_rune_is_fired(t *T) {
	ee, tt := Test(t.GoT(), nil, -1)
	defer ee.QuitListening()
	t.False(ee.IsListening())
	t.True(tt.FireRune('r').IsListening())
}

func (s *_Testing) Starts_listening_with_update_request(t *T) {
	fx := &cmpFX{}
	ee, _ := Test(t.GoT(), &cmpFX{}, 3) // TODO: clarify what's reported
	defer ee.QuitListening()            // during this test
	t.False(ee.IsListening())
	t.FatalOn(ee.Update(fx, nil, nil))
	t.True(ee.IsListening())
}

func (s *_Testing) Starts_listening_on_a_fired_mouse_event(t *T) {
	ee, tt := Test(t.GoT(), nil, -1) // listen for ever
	t.False(ee.IsListening())
	t.True(tt.FireClick(0, 0).IsListening())
	ee.QuitListening()

	ee, tt = Test(t.GoT(), nil, -1)
	t.False(ee.IsListening())
	t.True(tt.FireContext(0, 0).IsListening())
	ee.QuitListening()

	ee, tt = Test(t.GoT(), nil, -1)
	t.False(ee.IsListening())
	t.True(tt.FireMouse(
		0, 0, tcell.Button1, tcell.ModNone).IsListening())
	ee.QuitListening()
}

type clickFX struct {
	Component
	clicked, context bool
}

func (c *clickFX) OnClick(_ *Env, _, _ int) {
	c.clicked = true
}

func (c *clickFX) OnContext(_ *Env, _, _ int) {
	c.context = true
}

func (s *_Testing) Starts_listening_on_a_fired_component_click(t *T) {
	fx := &clickFX{}
	ee, tt := Test(t.GoT(), fx, 0) // listen for ever
	defer ee.QuitListening()
	t.False(ee.IsListening())
	t.True(tt.FireComponentClick(fx).IsListening())
}

func (s *_Testing) Counts_down_two_on_reported_component_click(t *T) {
	fx := &clickFX{}
	ee, tt := Test(t.GoT(), fx, 2)
	t.False(ee.IsListening())
	t.False(tt.FireComponentClick(fx).IsListening())
}

func (s *_Testing) Component_click_is_reported_to_component(t *T) {
	fx := &clickFX{}
	_, tt := Test(t.GoT(), fx, 2)
	tt.FireComponentClick(fx)
	t.True(fx.clicked)
}

func (s *_Testing) Starts_listening_on_a_fired_context(t *T) {
	fx := &clickFX{}
	ee, tt := Test(t.GoT(), fx, 0) // listen for ever
	defer ee.QuitListening()
	t.False(ee.IsListening())
	t.True(tt.FireComponentContext(fx).IsListening())
}

func (s *_Testing) Counts_down_two_on_reported_context(t *T) {
	fx := &clickFX{}
	ee, tt := Test(t.GoT(), fx, 2)
	t.False(ee.IsListening())
	t.False(tt.FireComponentContext(fx).IsListening())
}

func (s *_Testing) Context_is_reported_to_component(t *T) {
	fx := &clickFX{}
	_, tt := Test(t.GoT(), fx, 2)
	tt.FireComponentContext(fx)
	t.True(fx.context)
}

func TestTesting(t *testing.T) {
	t.Parallel()
	Run(&_Testing{}, t)
}
