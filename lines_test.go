// Copyright (c) 2022 Stephan Lukits. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package lines

import (
	"fmt"
	"testing"

	. "github.com/slukits/gounit"
)

type _lines struct{ Suite }

func (s *_lines) SetUp(t *T) { t.Parallel() }

type initFX struct{ Component }

const expInit = "component-fixture initialized"

func (c *initFX) OnInit(e *Env) { fmt.Fprint(e, expInit) }

func (s *_lines) tt(t *T, cmp Componenter) *Fixture {
	return fx(t, cmp)
}

func (s *_lines) Initializes_initially_given_component(t *T) {
	init := false
	fx(t, &cmpFX{onInit: func(cf *cmpFX, e *Env) { init = true }})
	t.True(init)
}

func (s *_lines) Reports_quit_key_events_to_all_quitter(t *T) {
	quitReported := 0
	fx_ := fx(t, &cmpFX{})
	fx_.Lines.
		OnQuit(func() { quitReported++ }).
		OnQuit(func() { quitReported++ })
	fx_.FireRune('q')
	t.Eq(2, quitReported)
	fx_ = fx(t, &cmpFX{})
	fx_.Lines.
		OnQuit(func() { quitReported++ }).
		OnQuit(func() { quitReported++ })
	fx_.FireKey(CtrlC)
	t.Eq(4, quitReported)
	fx_ = fx(t, &cmpFX{})
	fx_.Lines.
		OnQuit(func() { quitReported++ }).
		OnQuit(func() { quitReported++ })
	fx_.FireKey(CtrlD)
	t.Eq(6, quitReported)
}

func (s *_lines) Reports_layout_after_initialization(t *T) {
	cmp := &cmpFX{}
	fx(t, cmp)
	t.True(cmp.T(onInit).Before(cmp.T(onLayout)))
}

func (s *_lines) Reports_update_to_provided_listener(t *T) {
	reported := false
	fx, cmp := fxCmp(t)
	t.FatalOn(fx.Lines.Update(cmp, nil, func(_ *Env) {
		reported = true
	}))
	t.True(reported)
}

func (s *_lines) Reports_update_without_listener_to_component(t *T) {
	reported := false
	cmp := &cmpFX{
		onUpdate: func(_ *cmpFX, _ *Env, _ interface{}) {
			reported = true
		},
	}
	fx := fx(t, cmp)
	t.FatalOn(fx.Lines.Update(cmp, nil, nil))
	t.True(reported)
}

/*
	+-App-----------------------------------------+
	| +-chainer---------------------------------+ |
	| | +-cmp1---------+ +-stacker------------+ | |
	| | |              | | +-cmp2-----------+ | | |
	| | |              | | |                | | | |
	| | |              | | |                | | | |
	| | |              | | +----------------+ | | |
	| | +--------------+ +--------------------+ | |
	| +-----------------------------------------+ |
	+---------------------------------------------+

Focusing cmp1 we expect that chainer and cmp1 get OnFocus reported.
Focusing then cmp2 we expect that stacker and cmp2 get OnFocus reported
but NOT chainer or App.
*/
func (s *_lines) Reports_focus_gain_to_all_parents_not_focused(t *T) {
	app, chn, stk := &stackingFX{}, &chainingFX{}, &stackingFX{}
	cmp1, cmp2 := &cmpFX{}, &cmpFX{}
	app.onInit = func(c *cmpFX, e *Env) { app.CC = append(app.CC, chn) }
	chn.onInit = func(c *cmpFX, e *Env) {
		chn.CC = append(chn.CC, cmp1, stk)
	}
	stk.onInit = func(c *cmpFX, e *Env) {
		stk.CC = append(stk.CC, cmp2)
	}
	fx := fx(t, app)
	t.True(app.N(onFocus) == 0 && chn.N(onFocus) == 0 &&
		cmp1.N(onFocus) == 0 && stk.N(onFocus) == 0 &&
		cmp2.N(onFocus) == 0)
	fx.Lines.Focus(cmp1)
	t.True(app.N(onFocus) == 0 && chn.N(onFocus) == 1 &&
		cmp1.N(onFocus) == 1 && stk.N(onFocus) == 0 &&
		cmp2.N(onFocus) == 0)
}

/*
	+-App-----------------------------------------+
	| +-chainer---------------------------------+ |
	| | +-cmp1---------+ +-stacker------------+ | |
	| | |              | | +-cmp2-----------+ | | |
	| | |              | | |                | | | |
	| | |              | | |                | | | |
	| | |              | | +----------------+ | | |
	| | +--------------+ +--------------------+ | |
	| +-----------------------------------------+ |
	+---------------------------------------------+

Focusing cmp2 we expect no component to loos its focus.
Focusing then cmp1 we expect cmp2 and stacker to loos its focus but NOT
chainer or App.
*/
func (s *_lines) Reports_focus_loss_to_all_parents_not_focused_to(t *T) {
	app, chn, stk := &stackingFX{}, &chainingFX{}, &stackingFX{}
	cmp1, cmp2 := &cmpFX{}, &cmpFX{}
	app.onInit = func(c *cmpFX, e *Env) { app.CC = append(app.CC, chn) }
	chn.onInit = func(c *cmpFX, e *Env) {
		chn.CC = append(chn.CC, cmp1, stk)
	}
	stk.onInit = func(c *cmpFX, e *Env) {
		stk.CC = append(stk.CC, cmp2)
	}
	fx := fx(t, app)
	fx.Lines.Focus(cmp2)
	t.True(
		app.N(onFocusLost) == 0 && chn.N(onFocusLost) == 0 &&
			cmp1.N(onFocusLost) == 0 && stk.N(onFocusLost) == 0 &&
			cmp2.N(onFocusLost) == 0)

	fx.Lines.Focus(cmp1)
	t.True(app.N(onFocusLost) == 0 && chn.N(onFocusLost) == 0 &&
		cmp1.N(onFocusLost) == 0 && stk.N(onFocusLost) == 1 &&
		cmp2.N(onFocusLost) == 1)
}

type stackedCmpFX struct {
	Component
	Stacking
	lostFocusReported bool
}

func newStacking(cc ...Componenter) *stackedCmpFX {
	return &stackedCmpFX{Stacking: Stacking{CC: cc}}
}

func (c *stackedCmpFX) stacked(idx int) *cmpFX {
	return c.CC[idx].(*cmpFX)
}

func (c *stackedCmpFX) OnFocusLost(*Env) { c.lostFocusReported = true }

type fcsCmpFX struct {
	Component
	hasFocus bool
}

func (c *fcsCmpFX) OnFocus(*Env)     { c.hasFocus = true }
func (c *fcsCmpFX) OnFocusLost(*Env) { c.hasFocus = false }

func (s *_lines) Ignores_focus_on_focused_component(t *T) {
	cmp := &stackingFX{}
	cmp.CC = append(cmp.CC, &cmpFX{})
	fx := fx(t, cmp)
	fx.Lines.Update(cmp, nil, func(e *Env) {
		t.Eq(cmp, e.Focused())
	})
	t.FatalOn(fx.Lines.Focus(cmp))
	t.Eq(0, cmp.N(onFocusLost))
	t.Eq(0, cmp.N(onFocus))
}

func (s *_lines) Moves_focus_to_a_newly_set_root(t *T) {
	fx := fx(t, &cmpFX{})
	t.FatalOn(fx.Lines.Update(fx.Root(), nil, func(e *Env) {
		t.Eq(fx.Root(), e.Focused())
	}))

	cmp := &cmpFX{}
	fx.Lines.SetRoot(cmp)
	t.Eq(fx.Root(), cmp)
	t.FatalOn(fx.Lines.Update(fx.Root(), nil, func(e *Env) {
		t.Eq(fx.Root(), cmp)
	}))
}

func (s *_lines) Focuses_root_if_focused_component_removed(t *T) {
	cmp := &stackingFX{}
	cmp.CC = append(cmp.CC, &cmpFX{})
	fx := fx(t, cmp)
	t.FatalOn(fx.Lines.Focus(cmp.CC[0]))
	fx.Lines.Update(fx.Root(), nil, func(e *Env) {
		t.Eq(cmp.CC[0], e.Focused())
	})

	fx.Lines.Update(fx.Root(), nil, func(e *Env) {
		cmp.CC = []Componenter{}
	})
	fx.Lines.Update(fx.Root(), nil, func(e *Env) {
		t.Eq(fx.Root(), e.Focused())
	})
}

func (s *_lines) Ignores_focus_moving_if_modal_layer_is_focused(t *T) {
	cmp, modalLayer := &cmpFX{}, &modalLayerFX{}
	fx := fx(t, cmp)
	fx.FireResize(3, 3)
	fx.Lines.Update(cmp, nil, func(e *Env) {
		cmp.Layered(e, modalLayer, nil)
	})
	fx.Lines.Update(fx.Root(), nil, func(e *Env) {
		t.Eq(modalLayer, e.Focused())
	})

	fx.Lines.Focus(cmp)
	fx.Lines.Update(cmp, nil, func(e *Env) {
		t.Eq(modalLayer, e.Focused())
	})
}

func (s *_lines) Has_no_cursor_component_if_cursor_removed(t *T) {
	stacking := &stackingFX{}
	stacking.CC = append(stacking.CC, &cmpFX{}, &cmpFX{})
	tt := fx(t, stacking)
	cmp := stacking.CC[1].(*cmpFX)
	tt.Lines.Update(cmp, nil, func(e *Env) {
		cmp.SetCursor(0, 0, BlockCursorBlinking)
	})
	t.Eq(cmp, tt.Lines.CursorComponent())
	tt.Lines.Update(cmp, nil, func(e *Env) {
		cmp.SetCursor(0, 0, ZeroCursor)
	})
	t.Eq(Componenter(nil), tt.Lines.CursorComponent())
	tt.Lines.Update(cmp, nil, func(e *Env) {
		cmp.SetCursor(0, 0, BlockCursorBlinking)
	})
	t.Eq(cmp, tt.Lines.CursorComponent())
	tt.Lines.RemoveCursor()
	t.Eq(Componenter(nil), tt.Lines.CursorComponent())
}

func (s *_lines) Provides_component_to_accordingly_set_cursor(t *T) {
	stacking := &stackingFX{}
	stacking.CC = append(stacking.CC, &cmpFX{}, &cmpFX{})
	tt := fx(t, stacking)
	cmp := stacking.CC[1].(*cmpFX)
	tt.Lines.Update(cmp, nil, func(e *Env) {
		tt.Lines.SetCursor(cmp.Dim().X(), cmp.Dim().Y())
	})
	t.Eq(cmp, tt.Lines.CursorComponent())
}

func (s *_lines) Ignores_setting_cursor_outside_the_screen(t *T) {
	tt, width, height := fx(t, &cmpFX{}), 0, 0
	tt.Lines.Update(tt.Root(), nil, func(e *Env) {
		width, height = e.ScreenSize()
	})
	tt.Lines.SetCursor(-1, 0)
	t.Eq(Componenter(nil), tt.Lines.CursorComponent())
	tt.Lines.SetCursor(0, -1)
	t.Eq(Componenter(nil), tt.Lines.CursorComponent())
	tt.Lines.SetCursor(width, 0)
	t.Eq(Componenter(nil), tt.Lines.CursorComponent())
	tt.Lines.SetCursor(0, height)
	t.Eq(Componenter(nil), tt.Lines.CursorComponent())
}

func (s *_lines) Removes_cursor_not_in_content_area_on_resize(t *T) {
	stacking, chaining := &stackingFX{}, &chainingFX{}
	chaining.CC = append(chaining.CC, &cmpFX{}, &cmpFX{
		onInit: func(cf *cmpFX, e *Env) {
			cf.Dim().SetWidth(32).SetHeight(8)
		}},
	)
	stacking.CC = append(stacking.CC, &cmpFX{}, chaining)
	tt := fx(t, stacking)
	cmp := chaining.CC[1].(*cmpFX)
	tt.Lines.Update(cmp, nil, func(e *Env) {
		tt.Lines.SetCursor(cmp.Dim().X(), cmp.Dim().Y())
		_, _, haveCursor := cmp.CursorPosition()
		t.Not.True(haveCursor)
	})
	_, _, haveCursor := tt.Lines.CursorPosition()
	t.True(haveCursor)
	tt.FireResize(70, 20)
	_, _, haveCursor = tt.Lines.CursorPosition()
	t.Not.True(haveCursor)
}

func (s *_lines) Adjusts_set_content_area_cursor_on_resize(t *T) {
	stacking, chaining := &stackingFX{}, &chainingFX{}
	chaining.CC = append(chaining.CC, &cmpFX{}, &cmpFX{
		onInit: func(cf *cmpFX, e *Env) {
			cf.Dim().SetWidth(32).SetHeight(8)
		}},
	)
	stacking.CC = append(stacking.CC, &cmpFX{}, chaining)
	tt := fx(t, stacking)
	cmp, cx, cy := chaining.CC[1].(*cmpFX), 0, 0
	tt.Lines.Update(cmp, nil, func(e *Env) {
		cx, cy, _, _ = cmp.ContentArea()
		tt.Lines.SetCursor(cx+1, cy+1)
	})
	testInvariant := func() {
		x, y, _ := tt.Lines.CursorPosition()
		t.Eq(cx+1, x)
		t.Eq(cy+1, y)
	}
	testInvariant()
	tt.FireResize(70, 20)
	tt.Lines.Update(cmp, nil, func(e *Env) {
		cx, cy, _, _ = cmp.ContentArea()
	})
	testInvariant()
}

func (s *_lines) Reports_cursor_change_on_resize(t *T) {
	cmp := &cmpFX{
		onInit: func(c *cmpFX, e *Env) {
			c.FF.Add(CellFocusable)
			c.Dim().SetWidth(10).SetHeight(2)
			fmt.Fprint(e, "1st\n2nd\n3rd")
			e.Lines.SetCursor(0, 0)
		},
		onCursor: func(c *cmpFX, e *Env, absOnly bool) {
			switch c.N(onCursor) {
			case 1:
				x, y, cursorSet := e.Lines.CursorPosition()
				t.True(!absOnly && !cursorSet && x == y && x == -1)
			case 2:
				x, y, cursorSet := c.CursorPosition()
				t.True(absOnly && cursorSet && x == y && x == 1)
			}
		},
	}
	fx := fx(t, cmp)
	x, y, cursorSet := fx.Lines.CursorPosition()
	t.True(cursorSet && x == y && x == 0)
	fx.FireResize(40, 10)

	t.FatalIfNot(t.Eq(1, cmp.N(onCursor)))

	fx.Lines.Update(cmp, nil, func(e *Env) {
		cmp.SetCursor(1, 1)
		ln, cl, cursorSet := cmp.CursorPosition()
		t.True(cursorSet && ln == cl && ln == 1)
	})
	fx.FireResize(3, 2)

	t.FatalIfNot(t.Eq(2, cmp.N(onCursor)))
}

func TestLines(t *testing.T) {
	t.Parallel()
	Run(&_lines{}, t)
}
