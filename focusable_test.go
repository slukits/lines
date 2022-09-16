// Copyright (c) 2022 Stephan Lukits. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package lines

import (
	"fmt"
	"testing"

	"github.com/gdamore/tcell/v2"
	. "github.com/slukits/gounit"
)

type lineFocus struct {
	Suite
	Fixtures
}

func (s *lineFocus) SetUp(t *T) {
	t.Parallel()
}

func (s *lineFocus) TearDown(t *T) {
	ee, ok := s.Get(t).(*Events)
	if !ok {
		return
	}
	if !ee.IsListening() {
		return
	}
	ee.QuitListening()
}

func (s *lineFocus) iFX(
	t *T, init func(*icmpFX, *Env), max ...int,
) (*Events, *Testing, *icmpFX) {
	fx := &icmpFX{init: init}
	ee, tt := Test(t.GoT(), fx, max...)
	s.Set(t, ee)
	ee.Listen()
	return ee, tt, fx
}

func (s *lineFocus) Has_initially_no_line_focused(t *T) {
	s.iFX(t, func(c *icmpFX, e *Env) {
		t.Eq(-1, c.Focus.Current())
		// no-ops for coverage
		c.Focus.Next(false)
		c.Focus.Reset(false)
		t.Eq(-1, c.Focus.Current())
	}, 1)
}

type lfCmpFX struct {
	Component
	onInit func(*lfCmpFX, *Env)
	onLf   func(*lfCmpFX, *Env, int)
	lfN    int
}

func (c *lfCmpFX) OnInit(e *Env) {
	c.FF.Add(LinesFocusable)
	if c.onInit == nil {
		return
	}
	c.onInit(c, e)
}

func (c *lfCmpFX) OnLineFocus(e *Env, idx int) {
	c.lfN++
	if c.onLf == nil {
		return
	}
	c.onLf(c, e, idx)
}

func (s *lineFocus) lfFX(t *T,
	init func(*lfCmpFX, *Env),
	onLf func(*lfCmpFX, *Env, int),
	max ...int,
) (*Events, *Testing, *lfCmpFX) {
	fx := &lfCmpFX{
		onInit: init,
		onLf:   onLf,
	}
	ee, tt := Test(t.GoT(), fx, max...)
	s.Set(t, ee)
	ee.Listen()
	return ee, tt, fx
}

func (s *lineFocus) Focuses_first_focusable_line(t *T) {
	ee, tt, fx := s.lfFX(t,
		func(c *lfCmpFX, e *Env) { // OnInit
			c.FF.Add(LinesFocusable)
			fmt.Fprint(e, "first\nsecond")
		},
		func(c *lfCmpFX, e *Env, i int) { // OnLineFocus
			switch c.lfN {
			case 1:
				t.Eq(0, c.Focus.Current())
			case 2:
				t.Eq(1, c.Focus.Current())
			}
		}, 4)

	tt.FireRune('j')

	ee.Update(fx, nil, func(e *Env) {
		fx.Focus.Reset(false)
		fmt.Fprint(e.LL(0, NotFocusable), "first")
		fmt.Fprint(e.LL(1), "second")
	})

	tt.FireKey(tcell.KeyDown)

	t.Eq(2, fx.lfN)
}

func (s *lineFocus) Focuses_next_focusable_line(t *T) {
	_, tt, fx := s.lfFX(t,
		func(c *lfCmpFX, e *Env) {
			fmt.Fprint(e.LL(0), "first")
			fmt.Fprint(e.LL(1, NotFocusable), "second")
			fmt.Fprint(e.LL(2), "third")
		},
		func(c *lfCmpFX, e *Env, idx int) {
			switch c.lfN {
			case 1:
				t.Eq(0, c.Focus.Current())
			case 2:
				t.Eq(2, c.Focus.Current())
			}
		}, 3)

	tt.FireRune('j')
	tt.FireKey(tcell.KeyDown)

	t.Eq(2, fx.lfN)
}

func (s *lineFocus) Resets_if_no_next_focusable(t *T) {
	ee, tt, fx := s.lfFX(t,
		func(c *lfCmpFX, e *Env) { // OnInit
			fmt.Fprint(e.LL(0), "first")
		},
		func(c *lfCmpFX, e *Env, i int) { // OnLineFocus
			switch c.lfN % 2 {
			case 0:
				t.Eq(-1, c.Focus.Current())
			case 1:
				t.Eq(0, c.Focus.Current())
			}
		})

	tt.FireRune('j')
	tt.FireRune('j')

	ee.Update(fx, nil, func(e *Env) {
		t.Eq(-1, fx.Focus.Current())
		fmt.Fprint(e.LL(1, NotFocusable), "second")
		fmt.Fprint(e.LL(2, NotFocusable), "third")
	})

	tt.FireKey(tcell.KeyDown)
	tt.FireKey(tcell.KeyDown)

	t.Eq(4, fx.lfN)
}

func (s *lineFocus) Focuses_previous_focusable_line(t *T) {
	ee, tt, fx := s.lfFX(t,
		func(lcf *lfCmpFX, e *Env) { // OnInit
			fmt.Fprint(e, "first\nsecond")
		},
		func(c *lfCmpFX, e *Env, i int) { // OnLineFocused
			switch c.lfN {
			case 1:
				t.Eq(0, c.Focus.Current())
			case 2, 4:
				t.Eq(-1, c.Focus.Current())
			case 3, 5, 7:
				t.Eq(1, c.Focus.Current())
			case 6:
				t.Eq(2, c.Focus.Current())
			}
		})

	tt.FireRune('j') // case 1
	tt.FireRune('k') // case 2

	ee.Update(fx, nil, func(e *Env) {
		fmt.Fprint(e.LL(0, NotFocusable), "first")
		fmt.Fprint(e.LL(1), "second")
		fmt.Fprint(e.LL(2), "third")
	})

	tt.FireRune('j') // case 3
	tt.FireRune('k') // case 4
	tt.FireRune('j') // case 5
	tt.FireRune('j') // case 6
	tt.FireRune('k') // case 7

	t.Eq(7, fx.lfN)
}

func (s *lineFocus) Reset_triggered_by_unfocusable_feature(t *T) {
	_, tt, fx := s.lfFX(t,
		func(c *lfCmpFX, e *Env) {
			fmt.Fprint(e, "first\nsecond")
		},
		func(c *lfCmpFX, e *Env, i int) {
			switch c.lfN {
			case 1:
				t.Eq(0, c.Focus.Current())
			case 2:
				t.Eq(-1, c.Focus.Current())
			}
		})

	tt.FireRune('j')
	tt.FireKey(tcell.KeyEsc, tcell.ModNone)

	t.Eq(2, fx.lfN)
}

func (s *lineFocus) Scrolls_to_next_highlighted_line(t *T) {
	ee, tt, fx := s.lfFX(t,
		func(c *lfCmpFX, e *Env) {
			c.dim.SetHeight(2)
			fmt.Fprint(e.LL(0, NotFocusable), "line 1")
			fmt.Fprint(e.LL(1, NotFocusable), "line 2")
			fmt.Fprint(e.LL(2, NotFocusable), "line 3")
			fmt.Fprint(e.LL(3, NotFocusable), "line 4")
			fmt.Fprint(e.LL(4, NotFocusable), "line 5")
			fmt.Fprint(e.LL(5), "line 6")
			fmt.Fprint(e.LL(6), "line 7")
		},
		func(c *lfCmpFX, e *Env, i int) {
			t.Eq(5, c.Focus.Current())
		})

	tt.FireKey(tcell.KeyDown)

	ee.Update(fx, nil, func(e *Env) {
		t.Eq("line 5\nline 6", tt.Screen().String())
	})

	t.Eq(1, fx.lfN)
}

func (s *lineFocus) Scrolls_to_previous_highlighted_line(t *T) {
	ee, tt, fx := s.lfFX(t,
		func(c *lfCmpFX, e *Env) {
			c.dim.SetHeight(2)
			fmt.Fprint(e.LL(0, NotFocusable), "line 1")
			fmt.Fprint(e.LL(1), "line 2")
			fmt.Fprint(e.LL(2, NotFocusable), "line 3")
			fmt.Fprint(e.LL(3, NotFocusable), "line 4")
			fmt.Fprint(e.LL(4, NotFocusable), "line 5")
			fmt.Fprint(e.LL(5), "line 6")
			fmt.Fprint(e.LL(6), "line 7")
		},
		func(c *lfCmpFX, e *Env, i int) {
			switch c.lfN {
			case 1, 3:
				t.Eq(1, c.Focus.Current())
			case 2:
				t.Eq(5, c.Focus.Current())
			}
		})

	tt.FireKey(tcell.KeyDown)
	tt.FireKey(tcell.KeyDown)

	ee.Update(fx, nil, func(e *Env) {
		t.Eq("line 5\nline 6", tt.Screen().String())
	})

	tt.FireKey(tcell.KeyUp)

	ee.Update(fx, nil, func(e *Env) {
		t.Eq("line 2\nline 3", tt.Screen().String())
	})

	t.Eq(3, fx.lfN)
}

func (s *lineFocus) Inverts_bg_fg_of_focused_if_highlighted(t *T) {
	_, tt, fx := s.lfFX(t,
		func(c *lfCmpFX, e *Env) {
			c.dim.SetHeight(2)
			fmt.Fprint(e.LL(0, NotFocusable), "line 1")
			fmt.Fprint(e.LL(1), "line 2")
		},
		func(c *lfCmpFX, e *Env, i int) {
			switch c.lfN {
			case 1:
				t.Eq(1, c.Focus.Current())
			}
		})

	scr := tt.Screen()
	l1, l2 := scr[0], scr[1]
	t.Eq(l1.Styles().Of(0).BG(), l2.Styles().Of(0).BG())
	t.Eq(l1.Styles().Of(0).FG(), l2.Styles().Of(0).FG())

	tt.FireKey(tcell.KeyDown)
	t.Eq(1, fx.lfN)

	t.Eq(l1.Styles().Of(0).BG(), l2.Styles().Of(0).FG())
	t.Eq(l1.Styles().Of(0).FG(), l2.Styles().Of(0).BG())

}

type lsCmpFX struct {
	Component
	onIN     func(*lsCmpFX, *Env)
	onLF     func(*lsCmpFX, *Env, int)
	onLS     func(*lsCmpFX, *Env, int)
	lfN, lsN int
}

func (c *lsCmpFX) OnInit(e *Env) {
	if c.onIN == nil {
		return
	}
	c.onIN(c, e)
}

func (c *lsCmpFX) OnLineFocus(e *Env, idx int) {
	c.lfN++
	if c.onLF == nil {
		return
	}
	c.onLF(c, e, idx)
}

func (c *lsCmpFX) OnLineSelection(e *Env, idx int) {
	c.lsN++
	if c.onLS == nil {
		return
	}
	c.onLS(c, e, idx)
}

func (s *lineFocus) slFX(t *T,
	onIN func(*lsCmpFX, *Env),
	onLF func(*lsCmpFX, *Env, int),
	onLS func(*lsCmpFX, *Env, int),
	max ...int,
) (*Events, *Testing, *lsCmpFX) {
	fx := &lsCmpFX{
		onIN: onIN,
		onLF: onLF,
		onLS: onLS,
	}
	ee, tt := Test(t.GoT(), fx, max...)
	s.Set(t, ee)
	ee.Listen()
	return ee, tt, fx
}

func (s *lineFocus) Reports_focused_line_on_line_selection(t *T) {
	_, tt, fx := s.slFX(t,
		func(c *lsCmpFX, e *Env) { // OnInit
			c.FF.Add(LinesSelectable)
			fmt.Fprint(e, "first\nsecond")
		},
		func(c *lsCmpFX, e *Env, i int) { // OnLineFocus
			switch c.lfN {
			case 1:
				t.Eq(0, c.Focus.Current())
			case 2:
				t.Eq(1, c.Focus.Current())
			}
		},
		func(c *lsCmpFX, e *Env, i int) { // OnLineSelection
			t.Eq(1, i)
		},
	)

	tt.FireRune('j')
	tt.FireRune('j')
	tt.FireKey(tcell.KeyEnter)

	t.Eq(2, fx.lfN)
	t.Eq(1, fx.lsN)
}

func TestLineFocus(t *testing.T) {
	t.Parallel()
	Run(&lineFocus{}, t)
}
