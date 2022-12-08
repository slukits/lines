// Copyright (c) 2022 Stephan Lukits. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package lines

import (
	"fmt"
	"testing"

	. "github.com/slukits/gounit"
)

type _Fixture struct{ Suite }

func (s *_Fixture) SetUp(t *T) { t.Parallel() }

type clickFX struct {
	Component
	clicked, context    bool
	x, y, width, height int
	rx, ry              int
}

func (c *clickFX) OnLayout(_ *Env) bool {
	c.x, c.y, c.width, c.height = c.Dim().Printable()
	return false
}

func (c *clickFX) OnClick(_ *Env, rx, ry int) {
	c.clicked = true
	c.rx, c.ry = rx, ry
}

func (c *clickFX) OnContext(_ *Env, rx, ry int) {
	c.context = true
	c.rx, c.ry = rx, ry
}

func (s *_Fixture) tt(t *T, cmp Componenter) *Fixture {
	return TermFixture(t.GoT(), 0, cmp)
}

func (s *_Fixture) Component_click_is_reported_to_component(t *T) {
	tt := s.tt(t, &clickFX{})
	tt.FireComponentClick(tt.Root(), 0, 0)
	t.True(tt.Root().(*clickFX).clicked)
}

func (s *_Fixture) Component_coordinates_are_reported_on_click(t *T) {
	fx := &clickFX{}
	tt := s.tt(t, fx)
	x, y := fx.width/2, fx.height/2
	tt.FireComponentClick(fx, x, y)
	t.Eq(x, fx.rx)
	t.Eq(y, fx.ry)
}

func (s *_Fixture) Ignores_component_click_if_coordinates_outside(t *T) {
	fx := &clickFX{}
	tt := s.tt(t, fx)

	tt.FireComponentClick(fx, -1, 0)

	t.Not.True(fx.clicked)

	tt.FireComponentClick(fx, 0, -1)

	t.Not.True(fx.clicked)

	tt.FireComponentClick(fx, fx.x+fx.width+1, 0)

	t.Not.True(fx.clicked)

	tt.FireComponentClick(fx, 0, fx.y+fx.height+1)

	t.Not.True(fx.clicked)
}

func (s *_Fixture) Context_is_reported_to_component(t *T) {
	fx := &clickFX{}
	tt := s.tt(t, fx)

	tt.FireComponentContext(fx, fx.width-1, fx.height-1)

	t.True(fx.context)
}

func (s *_Fixture) Component_coordinates_are_reported_on_context(t *T) {
	fx := &clickFX{}
	tt := s.tt(t, fx)

	x, y := fx.width/2, fx.height/2
	tt.FireComponentContext(fx, x, y)

	t.Eq(x, fx.rx)
	t.Eq(y, fx.ry)
}

func (s *_Fixture) Ignores_component_context_if_coordinates_outside(
	t *T,
) {
	fx := &clickFX{}
	tt := s.tt(t, fx)

	tt.FireComponentContext(fx, -1, 0)

	t.Not.True(fx.context)

	tt.FireComponentContext(fx, 0, -1)

	t.Not.True(fx.context)

	tt.FireComponentContext(fx, fx.x+fx.width+1, 0)

	t.Not.True(fx.context)

	tt.FireComponentContext(fx, 0, fx.y+fx.height+1)

	t.Not.True(fx.context)
}

func (s *_Fixture) Report_screen_portion_of_component(t *T) {
	// TODO: add at least a second component :)))
	exp := "123456\n223456\n323456\n423456\n523456\n623456"
	fx := &icmpFX{init: func(c *icmpFX, e *Env) {
		c.Dim().SetHeight(6).SetWidth(6)
		fmt.Fprint(e, exp)
	}}
	tt := s.tt(t, fx)

	ts := tt.ScreenOf(fx).Trimmed()
	t.Eq(exp, ts.String())
	t.Not.Eq(exp, tt.Screen().String())
	t.Eq("123456", ts[0])
}

// func (s *_Testing) Provides_line_s_cell_styles(t *T) {
// 	t.TODO()
// }

func TestFixture(t *testing.T) {
	t.Parallel()
	Run(&_Fixture{}, t)
}
