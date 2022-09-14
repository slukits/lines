// Copyright (c) 2022 Stephan Lukits. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package lines

import (
	"fmt"
	"strings"
	"testing"

	"github.com/gdamore/tcell/v2"
	. "github.com/slukits/gounit"
)

type _Testing struct{ Suite }

func (s *_Testing) SetUp(t *T) { t.Parallel() }

func (s *_Testing) Starts_non_blocking_listening_with_listen_call(t *T) {
	ee, _ := Test(t.GoT(), nil, 1)
	t.False(ee.IsListening())
	ee.Listen()
	t.True(ee.IsListening())
	ee.QuitListening()
	t.False(ee.IsListening())
}

func (s *_Testing) Starts_listening_if_a_resize_is_fired(t *T) {
	ee, tt := Test(t.GoT(), nil, 1)
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

func (s *_Testing) Starts_listening_on_a_fired_component_click(t *T) {
	fx := &clickFX{}
	ee, tt := Test(t.GoT(), fx) // listen for ever
	defer ee.QuitListening()
	t.False(ee.IsListening())
	t.True(tt.FireComponentClick(fx, 0, 0).IsListening())
}

func (s *_Testing) Counts_down_two_on_reported_component_click(t *T) {
	fx := &clickFX{}
	ee, tt := Test(t.GoT(), fx, 2) // plus OnLayout
	t.False(ee.IsListening())
	tt.FireComponentClick(fx, 0, 0)
	t.False(ee.IsListening())
}

func (s *_Testing) Component_click_is_reported_to_component(t *T) {
	fx := &clickFX{}
	_, tt := Test(t.GoT(), fx, 2)
	tt.FireComponentClick(fx, 0, 0)
	t.True(fx.clicked)
}

func (s *_Testing) Component_coordinates_are_reported_on_click(t *T) {
	fx := &clickFX{}
	ee, tt := Test(t.GoT(), fx, 2)
	ee.Listen()
	x, y := fx.width/2, fx.height/2
	tt.FireComponentClick(fx, x, y)
	t.Eq(x, fx.rx)
	t.Eq(y, fx.ry)
}

func (s *_Testing) Ignores_component_click_if_coordinates_outside(t *T) {
	fx := &clickFX{}
	ee, tt := Test(t.GoT(), fx, 2)
	defer ee.QuitListening()
	tt.FireComponentClick(fx, -1, 0)
	t.False(fx.clicked)
	tt.FireComponentClick(fx, 0, -1)
	t.False(fx.clicked)
	tt.FireComponentClick(fx, fx.x+fx.width+1, 0)
	t.False(fx.clicked)
	tt.FireComponentClick(fx, 0, fx.y+fx.height+1)
	t.False(fx.clicked)
	t.True(ee.IsListening())
}

func (s *_Testing) Starts_listening_on_a_fired_component_context(t *T) {
	fx := &clickFX{}
	ee, tt := Test(t.GoT(), fx)
	defer ee.QuitListening()
	t.False(ee.IsListening())
	tt.FireComponentContext(fx, 0, 0)
	t.True(ee.IsListening())
}

func (s *_Testing) Counts_down_two_on_reported_component_context(t *T) {
	fx := &clickFX{}
	ee, tt := Test(t.GoT(), fx, 2)
	t.False(ee.IsListening())
	t.False(tt.FireComponentContext(fx, 0, 0).IsListening())
}

func (s *_Testing) Context_is_reported_to_component(t *T) {
	fx := &clickFX{}
	ee, tt := Test(t.GoT(), fx, 2)
	ee.Listen() // layouts fx
	tt.FireComponentContext(fx, fx.width-1, fx.height-1)
	t.True(fx.context)
}

func (s *_Testing) Component_coordinates_are_reported_on_context(t *T) {
	fx := &clickFX{}
	ee, tt := Test(t.GoT(), fx, 2)
	ee.Listen()
	x, y := fx.width/2, fx.height/2
	tt.FireComponentContext(fx, x, y)
	t.Eq(x, fx.rx)
	t.Eq(y, fx.ry)
}

func (s *_Testing) Ignores_component_context_if_coordinates_outside(
	t *T,
) {
	fx := &clickFX{}
	ee, tt := Test(t.GoT(), fx, 2)
	defer ee.QuitListening()
	tt.FireComponentContext(fx, -1, 0)
	t.False(fx.context)
	tt.FireComponentContext(fx, 0, -1)
	t.False(fx.context)
	tt.FireComponentContext(fx, fx.x+fx.width+1, 0)
	t.False(fx.context)
	tt.FireComponentContext(fx, 0, fx.y+fx.height+1)
	t.False(fx.context)
	t.True(ee.IsListening())
}

func (s *_Testing) Provides_trimmed_screen(t *T) {
	ee, tt := Test(t.GoT(), &icmpFX{init: func(c *icmpFX, e *Env) {
		fmt.Fprint(e,
			"                    \n"+
				"   upper left       \n"+
				"                    \n"+
				"          right     \n"+
				"      bottom        \n"+
				"                    ",
		)
	}}, 0)
	tt.FireResize(20, 6)
	defer ee.QuitListening()
	exp := "upper left  \n" +
		"            \n" +
		"       right\n" +
		"   bottom   "

	t.Eq(exp, tt.Screen().String())
}

func (s *_Testing) Provides_string_with_screen_content(t *T) {
	ee, tt := Test(t.GoT(), &icmpFX{init: func(c *icmpFX, e *Env) {
		for i := 0; i < 20; i++ {
			if i < 10 {
				fmt.Fprintf(e.LL(i), "line 0%d", i)
				continue
			}
			fmt.Fprintf(e.LL(i), "line %d", i)
		}
	}}, 0)
	tt.FireResize(7, 20)
	defer ee.QuitListening()

	exp := []string{}
	for i := 0; i < 20; i++ {
		if i < 10 {
			exp = append(exp, fmt.Sprintf("line 0%d", i))
			continue
		}
		exp = append(exp, fmt.Sprintf("line %d", i))
	}
	t.Eq(strings.Join(exp, "\n"), tt.String())
}

func (s *_Testing) Report_screen_portion_of_component(t *T) {
	exp := "123456\n223456\n323456\n423456\n523456\n623456"
	fx := &icmpFX{init: func(c *icmpFX, e *Env) {
		c.Dim().SetHeight(6).SetWidth(6)
		fmt.Fprint(e, exp)
	}}
	ee, tt := Test(t.GoT(), fx)
	ee.Listen()
	defer ee.QuitListening()

	ts := tt.Trim(tt.ScreenOf(fx))
	t.Eq(exp, ts.String())
	t.Neq(exp, tt.FullScreen().String())
	t.Eq("123456", ts[0].String())
}

func (s *_Testing) Provides_line_s_cell_styles(t *T) {
}

func TestTesting(t *testing.T) {
	t.Parallel()
	Run(&_Testing{}, t)
}
