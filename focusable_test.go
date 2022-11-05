// Copyright (c) 2022 Stephan Lukits. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package lines

import (
	"fmt"
	"testing"

	. "github.com/slukits/gounit"
)

type lineFocus struct {
	Suite
}

func (s *lineFocus) SetUp(t *T) { t.Parallel() }

func (s *lineFocus) itt(
	t *T, init func(*icmpFX, *Env),
) (*Fixture, *icmpFX) {
	fx := &icmpFX{init: init}
	tt := TermFixture(t.GoT(), 0, fx)
	return tt, fx
}

func (s *lineFocus) Has_initially_no_line_focused(t *T) {
	s.itt(t, func(c *icmpFX, e *Env) {
		t.Eq(-1, c.LL.Focus.Current())
		// no-ops for coverage
		c.LL.Focus.Next(false)
		c.LL.Focus.Reset(false)
		t.Eq(-1, c.LL.Focus.Current())
	})
}

type lfCmpFX struct {
	Component
	onInit func(*lfCmpFX, *Env)
	onLf   func(*lfCmpFX, *Env, int)
	// the number of received line focus events
	lfN int
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

func (s *lineFocus) lfFX(
	t *T,
	init func(*lfCmpFX, *Env),
	onLf func(*lfCmpFX, *Env, int),
) (*Fixture, *lfCmpFX) {
	fx := &lfCmpFX{
		onInit: init,
		onLf:   onLf,
	}
	tt := TermFixture(t.GoT(), 0, fx)
	return tt, fx
}

func (s *lineFocus) Focuses_first_focusable_line(t *T) {
	tt, fx := s.lfFX(t,
		func(c *lfCmpFX, e *Env) { // OnInit
			c.FF.Add(LinesFocusable)
			fmt.Fprint(e, "first\nsecond")
		},
		func(c *lfCmpFX, e *Env, i int) { // OnLineFocus
			switch c.lfN {
			case 1:
				t.Eq(0, c.LL.Focus.Current())
			case 2:
				t.Eq(1, c.LL.Focus.Current())
			}
		})

	tt.FireRune('j')

	t.FatalOn(tt.Lines.Update(fx, nil, func(e *Env) {
		fx.LL.Focus.Reset(false)
		fmt.Fprint(e.LL(0), "first")
		fx.LL.By(0).Flag(NotFocusable)
		fmt.Fprint(e.LL(1), "second")
	}))

	tt.FireKey(Down)
	t.Eq(2, fx.lfN)
}

func (s *lineFocus) Focuses_next_focusable_line(t *T) {
	tt, fx := s.lfFX(t,
		func(c *lfCmpFX, e *Env) {
			fmt.Fprint(e.LL(0), "first")
			fmt.Fprint(e.LL(1), "second")
			c.LL.By(1).Flag(NotFocusable)
			fmt.Fprint(e.LL(2), "third")
		},
		func(c *lfCmpFX, e *Env, idx int) {
			switch c.lfN {
			case 1:
				t.Eq(0, c.LL.Focus.Current())
			case 2:
				t.Eq(2, c.LL.Focus.Current())
			}
		})

	tt.FireRune('j')
	tt.FireKey(Down)

	t.Eq(2, fx.lfN)
}

func (s *lineFocus) Resets_if_no_next_focusable(t *T) {
	tt, fx := s.lfFX(t,
		func(c *lfCmpFX, e *Env) { // OnInit
			fmt.Fprint(e.LL(0), "first")
		},
		func(c *lfCmpFX, e *Env, i int) { // OnLineFocus
			switch c.lfN % 2 {
			case 0:
				t.Eq(-1, c.LL.Focus.Current())
			case 1:
				t.Eq(0, c.LL.Focus.Current())
			}
		})

	tt.FireRune('j')
	tt.FireRune('j')

	t.FatalOn(tt.Lines.Update(fx, nil, func(e *Env) {
		t.Eq(-1, fx.LL.Focus.Current())
		fmt.Fprint(e.LL(1), "second")
		fx.LL.By(1).Flag(NotFocusable)
		fmt.Fprint(e.LL(2), "third")
		fx.LL.By(2).Flag(NotFocusable)
	}))

	tt.FireKey(Down)
	tt.FireKey(Down)

	t.Eq(4, fx.lfN)
}

func (s *lineFocus) Focuses_previous_focusable_line(t *T) {
	tt, fx := s.lfFX(t,
		func(lcf *lfCmpFX, e *Env) { // OnInit
			fmt.Fprint(e, "first\nsecond")
		},
		func(c *lfCmpFX, e *Env, i int) { // OnLineFocused
			switch c.lfN {
			case 1:
				t.Eq(0, c.LL.Focus.Current())
			case 2, 4:
				t.Eq(-1, c.LL.Focus.Current())
			case 3, 5, 7:
				t.Eq(1, c.LL.Focus.Current())
			case 6:
				t.Eq(2, c.LL.Focus.Current())
			}
		})

	tt.FireRune('j') // case 1
	tt.FireRune('k') // case 2

	t.FatalOn(tt.Lines.Update(fx, nil, func(e *Env) {
		fmt.Fprint(e.LL(0), "first")
		fx.LL.By(0).Flag(NotFocusable)
		fmt.Fprint(e.LL(1), "second")
		fmt.Fprint(e.LL(2), "third")
	}))

	tt.FireRune('j') // case 3
	tt.FireRune('k') // case 4
	tt.FireRune('j') // case 5
	tt.FireRune('j') // case 6
	tt.FireRune('k') // case 7

	t.Eq(7, fx.lfN)
}

func (s *lineFocus) Reset_triggered_by_unfocusable_feature(t *T) {
	tt, fx := s.lfFX(t,
		func(c *lfCmpFX, e *Env) {
			fmt.Fprint(e, "first\nsecond")
		},
		func(c *lfCmpFX, e *Env, i int) {
			switch c.lfN {
			case 1:
				t.Eq(0, c.LL.Focus.Current())
			case 2:
				t.Eq(-1, c.LL.Focus.Current())
			}
		})

	tt.FireRune('j')
	tt.FireKey(Esc)

	t.Eq(2, fx.lfN)
}

func (s *lineFocus) Scrolls_to_next_highlighted_line(t *T) {
	tt, fx := s.lfFX(t,
		func(c *lfCmpFX, e *Env) {
			c.dim.SetHeight(2)
			for i := 0; i < 7; i++ {
				fmt.Fprintf(e.LL(i), "line %d", i+1)
			}
			for i := 0; i < 5; i++ {
				c.LL.By(i).Flag(NotFocusable)
			}
		},
		func(c *lfCmpFX, e *Env, i int) {
			t.Eq(5, c.LL.Focus.Current())
		})

	tt.FireKey(Down)

	t.Eq("line 5\nline 6", tt.ScreenOf(fx).Trimmed().String())

	t.Eq(1, fx.lfN)
}

func (s *lineFocus) Scrolls_to_previous_highlighted_line(t *T) {
	tt, fx := s.lfFX(t,
		func(c *lfCmpFX, e *Env) {
			c.dim.SetHeight(2)
			for i := 0; i < 7; i++ {
				fmt.Fprintf(e.LL(i), "line %d", i+1)
			}
			for _, idx := range []int{0, 2, 3, 4} {
				c.LL.By(idx).Flag(NotFocusable)
			}
		},
		func(c *lfCmpFX, e *Env, i int) {
			switch c.lfN {
			case 1, 3:
				t.Eq(1, c.LL.Focus.Current())
			case 2:
				t.Eq(5, c.LL.Focus.Current())
			}
		})

	tt.FireKey(Down)
	tt.FireKey(Down)

	t.Eq("line 5\nline 6", tt.ScreenOf(fx).Trimmed().String())

	tt.FireKey(Up)

	t.Eq("line 2\nline 3", tt.ScreenOf(fx).Trimmed().String())
	t.Eq(3, fx.lfN)
}

func (s *lineFocus) Inverts_bg_fg_of_focused_if_highlighted(t *T) {
	tt, fx := s.lfFX(t,
		func(c *lfCmpFX, e *Env) {
			c.FF.Add(HighlightedFocusable)
			c.dim.SetHeight(2)
			fmt.Fprint(e.LL(0), "line 1")
			c.LL.By(0).Flag(NotFocusable)
			fmt.Fprint(e.LL(1), "line 2")
		},
		func(c *lfCmpFX, e *Env, i int) {
			switch c.lfN {
			case 1:
				t.Eq(1, c.LL.Focus.Current())
			}
		})

	l2 := tt.CellsOf(fx).Trimmed()[1]
	for x := range l2 {
		t.Not.True(l2.HasAttr(x, Reverse))
	}

	tt.FireKey(Down)
	t.Eq(1, fx.lfN)

	l2 = tt.CellsOf(fx).Trimmed()[1]
	for x := range l2 {
		if x < len("line 2") {
			t.True(l2.HasAttr(x, Reverse))
			continue
		}
		t.Not.True(l2.HasAttr(x, Reverse))
	}
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
) (*Fixture, *lsCmpFX) {
	fx := &lsCmpFX{
		onIN: onIN,
		onLF: onLF,
		onLS: onLS,
	}
	tt := TermFixture(t.GoT(), 0, fx)
	return tt, fx
}

func (s *lineFocus) Reports_focused_line_on_line_selection(t *T) {
	tt, fx := s.slFX(t,
		func(c *lsCmpFX, e *Env) { // OnInit
			c.FF.Add(LinesSelectable)
			fmt.Fprint(e, "first\nsecond")
		},
		func(c *lsCmpFX, e *Env, i int) { // OnLineFocus
			switch c.lfN {
			case 1:
				t.Eq(0, c.LL.Focus.Current())
			case 2:
				t.Eq(1, c.LL.Focus.Current())
			}
		},
		func(c *lsCmpFX, e *Env, i int) { // OnLineSelection
			t.Eq(1, i)
		},
	)

	tt.FireRune('j')
	tt.FireRune('j')
	tt.FireKey(Enter)

	t.Eq(2, fx.lfN)
	t.Eq(1, fx.lsN)
}

func TestLineFocus(t *testing.T) {
	t.Parallel()
	Run(&lineFocus{}, t)
}
