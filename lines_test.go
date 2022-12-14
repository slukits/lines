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
	return fx(t, cmp)
}

func (s *_lines) Initializes_initially_given_component(t *T) {
	fx := fx(t, &initFX{})
	t.Eq(expInit, fx.Screen().Trimmed().String())
}

func (s *_lines) Reports_quit_key_events_to_all_quitter(t *T) {
	for _, k := range quitableFeatures.keysOf(Quitable) {
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

func (c *lytCmpFX) OnLayout(*Env) bool {
	c.lyt = time.Now()
	return false
}

func (s *_lines) Reports_layout_after_initialization(t *T) {
	cmp := &lytCmpFX{}
	fx(t, cmp)
	t.True(cmp.init.Before(cmp.lyt))
}

type updLstCmpFX struct {
	Component
	reported bool
}

func (s *_lines) Reports_update_to_provided_listener(t *T) {
	cmp := &updLstCmpFX{}
	fx := fx(t, cmp)
	t.FatalOn(fx.Lines.Update(cmp, nil, func(_ *Env) {
		cmp.reported = true
	}))
	t.True(cmp.reported)
}

type updCmpFX struct {
	Component
	reported bool
}

func (c *updCmpFX) OnUpdate(e *Env, _ interface{}) { c.reported = true }

func (s *_lines) Reports_update_without_listener_to_component(t *T) {
	cmp := &updCmpFX{}
	fx := fx(t, cmp)
	t.FatalOn(fx.Lines.Update(cmp, nil, nil))
	t.True(cmp.reported)
}

type fcsCmp struct {
	Component
	onFocus, onFocusLost int
}

func (c *fcsCmp) OnFocus(_ *Env)     { c.onFocus++ }
func (c *fcsCmp) OnFocusLost(_ *Env) { c.onFocusLost++ }

type fcsApp struct {
	fcsCmp
	Stacking
}

func (c *fcsApp) OnInit(_ *Env) { c.CC = append(c.CC, &fcsChn{}) }

func (c *fcsApp) chainer() *fcsChn { return c.CC[0].(*fcsChn) }
func (c *fcsApp) cmp1() *fcsCmp    { return c.chainer().CC[0].(*fcsCmp) }
func (c *fcsApp) stacker() *fcsStk { return c.chainer().CC[1].(*fcsStk) }
func (c *fcsApp) cmp2() *fcsCmp    { return c.stacker().CC[0].(*fcsCmp) }

type fcsChn struct {
	fcsCmp
	Chaining
}

func (c *fcsChn) OnInit(_ *Env) {
	c.CC = append(c.CC, &fcsCmp{}, &fcsStk{})
}

type fcsStk struct {
	fcsCmp
	Stacking
}

func (c *fcsStk) OnInit(_ *Env) { c.CC = append(c.CC, &fcsCmp{}) }

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
	app := &fcsApp{}
	fx := fx(t, app)
	t.True(app.onFocus == 0 && app.chainer().onFocus == 0 &&
		app.cmp1().onFocus == 0 && app.stacker().onFocus == 0 &&
		app.cmp2().onFocus == 0)

	fx.Lines.Focus(app.cmp1())
	t.True(app.onFocus == 0 && app.chainer().onFocus == 1 &&
		app.cmp1().onFocus == 1 && app.stacker().onFocus == 0 &&
		app.cmp2().onFocus == 0)

	fx.Lines.Focus(app.cmp2())
	t.True(app.onFocus == 0 && app.chainer().onFocus == 1 &&
		app.cmp1().onFocus == 1 && app.stacker().onFocus == 1 &&
		app.cmp2().onFocus == 1)
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
	app := &fcsApp{}
	fx := fx(t, app)
	fx.Lines.Focus(app.cmp2())
	t.True(app.onFocusLost == 0 && app.chainer().onFocusLost == 0 &&
		app.cmp1().onFocusLost == 0 && app.stacker().onFocusLost == 0 &&
		app.cmp2().onFocusLost == 0)

	fx.Lines.Focus(app.cmp1())
	t.True(app.onFocusLost == 0 && app.chainer().onFocusLost == 0 &&
		app.cmp1().onFocusLost == 0 && app.stacker().onFocusLost == 1 &&
		app.cmp2().onFocusLost == 1)
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
	cmp := newStacking(&fcsCmpFX{})
	fx := fx(t, cmp)
	t.FatalOn(fx.Lines.Focus(cmp))
	t.Not.True(cmp.lostFocusReported)
	t.Not.True(cmp.CC[0].(*fcsCmpFX).hasFocus)
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
	cmp := newStacking(&fcsCmpFX{})
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
	fx := fx(t, &layeredFX{layer: &mdlLayerFX{}})
	fx.FireResize(3, 3)
	fx.Root().(*layeredFX).addLayer(fx)
	fx.Lines.Update(fx.Root(), nil, func(e *Env) {
		t.Eq(fx.Root().(*layeredFX).layer, e.Focused())
	})

	fx.Lines.Focus(fx.Root())
	fx.Lines.Update(fx.Root(), nil, func(e *Env) {
		t.Eq(fx.Root().(*layeredFX).layer, e.Focused())
	})
}

func (s *_lines) Has_no_cursor_component_if_cursor_removed(t *T) {
	stacking := newStacking(&cmpFX{}, &cmpFX{})
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
	stacking := newStacking(&cmpFX{}, &cmpFX{})
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
	stk, chn := &stackedCmpFX{}, &chnFX{}
	chn.CC = append(chn.CC, &cmpFX{}, &cmpFX{onInit: func(cf *cmpFX, e *Env) {
		cf.Dim().SetWidth(32).SetHeight(8)
	}})
	stk.CC = append(stk.CC, &cmpFX{}, chn)
	tt := fx(t, stk)
	cmp := chn.CC[1].(*cmpFX)
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
	stk, chn := &stackedCmpFX{}, &chnFX{}
	chn.CC = append(chn.CC, &cmpFX{}, &cmpFX{onInit: func(cf *cmpFX, e *Env) {
		cf.Dim().SetWidth(32).SetHeight(8)
	}})
	stk.CC = append(stk.CC, &cmpFX{}, chn)
	tt := fx(t, stk)
	cmp, cx, cy := chn.CC[1].(*cmpFX), 0, 0
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
	fx := fx(t, cmp, 20*time.Minute)
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
