// Copyright (c) 2022 Stephan Lukits. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package lines

import (
	"testing"

	. "github.com/slukits/gounit"
)

type _Testing struct{ Suite }

func (s *_Testing) SetUp(t *T) { t.Parallel() }

type clickFX struct {
	Component
	clicked, context    bool
	x, y, width, height int
	rx, ry              int
}

func (c *clickFX) OnLayout(_ *Env) {
	c.x, c.y, c.width, c.height = c.Dim().Area()
}

func (c *clickFX) OnClick(_ *Env, rx, ry int) {
	c.clicked = true
	c.rx, c.ry = rx, ry
}

func (c *clickFX) OnContext(_ *Env, rx, ry int) {
	c.context = true
	c.rx, c.ry = rx, ry
}

func (s *_Testing) tt(t *T, cmp Componenter) (*Lines, *Testing) {
	return TermFixture(t.GoT(), 0, cmp)
}

func (s *_Testing) Component_click_is_reported_to_component(t *T) {
	_, tt := s.tt(t, &clickFX{})
	tt.FireComponentClick(tt.Root(), 0, 0)
	t.True(tt.Root().(*clickFX).clicked)
}

func (s *_Testing) Component_coordinates_are_reported_on_click(t *T) {
	fx := &clickFX{}
	_, tt := s.tt(t, fx)
	x, y := fx.width/2, fx.height/2
	tt.FireComponentClick(fx, x, y)
	t.Eq(x, fx.rx)
	t.Eq(y, fx.ry)
}

// func (s *_Testing) Ignores_component_click_if_coordinates_outside(t *T) {
// 	fx := &clickFX{}
// 	ee, tt := Test(t.GoT(), fx, 2)
// 	defer ee.QuitListening()
// 	tt.FireComponentClick(fx, -1, 0)
// 	t.Not.True(fx.clicked)
// 	tt.FireComponentClick(fx, 0, -1)
// 	t.Not.True(fx.clicked)
// 	tt.FireComponentClick(fx, fx.x+fx.width+1, 0)
// 	t.Not.True(fx.clicked)
// 	tt.FireComponentClick(fx, 0, fx.y+fx.height+1)
// 	t.Not.True(fx.clicked)
// 	t.True(ee.IsListening())
// }
//
// func (s *_Testing) Starts_listening_on_a_fired_component_context(t *T) {
// 	fx := &clickFX{}
// 	ee, tt := Test(t.GoT(), fx)
// 	defer ee.QuitListening()
// 	t.Not.True(ee.IsListening())
// 	tt.FireComponentContext(fx, 0, 0)
// 	t.True(ee.IsListening())
// }
//
// func (s *_Testing) Counts_down_two_on_reported_component_context(t *T) {
// 	fx := &clickFX{}
// 	ee, tt := Test(t.GoT(), fx, 2)
// 	t.Not.True(ee.IsListening())
// 	t.Not.True(tt.FireComponentContext(fx, 0, 0).IsListening())
// }
//
// func (s *_Testing) Context_is_reported_to_component(t *T) {
// 	fx := &clickFX{}
// 	ee, tt := Test(t.GoT(), fx, 2)
// 	ee.Listen() // layouts fx
// 	tt.FireComponentContext(fx, fx.width-1, fx.height-1)
// 	t.True(fx.context)
// }
//
// func (s *_Testing) Component_coordinates_are_reported_on_context(t *T) {
// 	fx := &clickFX{}
// 	ee, tt := Test(t.GoT(), fx, 2)
// 	ee.Listen()
// 	x, y := fx.width/2, fx.height/2
// 	tt.FireComponentContext(fx, x, y)
// 	t.Eq(x, fx.rx)
// 	t.Eq(y, fx.ry)
// }
//
// func (s *_Testing) Ignores_component_context_if_coordinates_outside(
// 	t *T,
// ) {
// 	fx := &clickFX{}
// 	ee, tt := Test(t.GoT(), fx, 2)
// 	defer ee.QuitListening()
// 	tt.FireComponentContext(fx, -1, 0)
// 	t.Not.True(fx.context)
// 	tt.FireComponentContext(fx, 0, -1)
// 	t.Not.True(fx.context)
// 	tt.FireComponentContext(fx, fx.x+fx.width+1, 0)
// 	t.Not.True(fx.context)
// 	tt.FireComponentContext(fx, 0, fx.y+fx.height+1)
// 	t.Not.True(fx.context)
// 	t.True(ee.IsListening())
// }
//
// func (s *_Testing) Provides_trimmed_screen(t *T) {
// 	ee, tt := Test(t.GoT(), &icmpFX{init: func(c *icmpFX, e *Env) {
// 		fmt.Fprint(e,
// 			"                    \n"+
// 				"   upper left       \n"+
// 				"                    \n"+
// 				"          right     \n"+
// 				"      bottom        \n"+
// 				"                    ",
// 		)
// 	}}, 0)
// 	tt.FireResize(20, 6)
// 	defer ee.QuitListening()
// 	exp := "upper left  \n" +
// 		"            \n" +
// 		"       right\n" +
// 		"   bottom   "
//
// 	t.Eq(exp, tt.ScreenZZZ().String())
// }
//
// func (s *_Testing) Provides_string_with_screen_content(t *T) {
// 	ee, tt := Test(t.GoT(), &icmpFX{init: func(c *icmpFX, e *Env) {
// 		for i := 0; i < 20; i++ {
// 			if i < 10 {
// 				fmt.Fprintf(e.LL(i), "line 0%d", i)
// 				continue
// 			}
// 			fmt.Fprintf(e.LL(i), "line %d", i)
// 		}
// 	}}, 0)
// 	tt.FireResize(7, 20)
// 	defer ee.QuitListening()
//
// 	exp := []string{}
// 	for i := 0; i < 20; i++ {
// 		if i < 10 {
// 			exp = append(exp, fmt.Sprintf("line 0%d", i))
// 			continue
// 		}
// 		exp = append(exp, fmt.Sprintf("line %d", i))
// 	}
// 	t.Eq(strings.Join(exp, "\n"), tt.StringZZZ())
// }
//
// func (s *_Testing) Report_screen_portion_of_component(t *T) {
// 	exp := "123456\n223456\n323456\n423456\n523456\n623456"
// 	fx := &icmpFX{init: func(c *icmpFX, e *Env) {
// 		c.Dim().SetHeight(6).SetWidth(6)
// 		fmt.Fprint(e, exp)
// 	}}
// 	ee, tt := Test(t.GoT(), fx)
// 	ee.Listen()
// 	defer ee.QuitListening()
//
// 	ts := tt.Trim(tt.ScreenOf(fx))
// 	t.Eq(exp, ts.String())
// 	t.Not.Eq(exp, tt.FullScreen().String())
// 	t.Eq("123456", ts[0].String())
// }
//
// func (s *_Testing) Provides_line_s_cell_styles(t *T) {
// }

func TestTesting(t *testing.T) {
	t.Parallel()
	Run(&_Testing{}, t)
}
