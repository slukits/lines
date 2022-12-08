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

type lineFocusFeat struct {
	Suite
}

func (s *lineFocusFeat) SetUp(t *T) { t.Parallel() }

func (s *lineFocusFeat) itt(
	t *T, init func(*icmpFX, *Env),
) (*Fixture, *icmpFX) {
	fx := &icmpFX{init: init}
	tt := TermFixture(t.GoT(), 0, fx)
	return tt, fx
}

func (s *lineFocusFeat) Has_initially_no_line_focused(t *T) {
	cmp := &cmpFX{onInit: func(c *cmpFX, e *Env) {
		t.Eq(-1, c.LL.Focus.Current())
		// no-ops for coverage
		c.LL.Focus.Next(false)
		c.LL.Focus.Reset()
		t.Eq(-1, c.LL.Focus.Current())
	}}
	fx(t, cmp)
}

type lfCmpFX struct {
	Component
	onInit func(*lfCmpFX, *Env)
	onLf   func(*lfCmpFX, *Env, int)
	onLfl  func(*lfCmpFX, *Env, int, int)
	onOf   func(*lfCmpFX, *Env, bool, bool)
	// the number of received line focus events
	lfN, lflN, ofN int
}

func (c *lfCmpFX) OnInit(e *Env) {
	c.FF.Add(LinesFocusable)
	if c.onInit == nil {
		return
	}
	c.onInit(c, e)
}

func (c *lfCmpFX) OnLineFocus(e *Env, cIdx, sIdx int) {
	c.lfN++
	if c.onLf == nil {
		return
	}
	c.onLf(c, e, cIdx)
}

func (c *lfCmpFX) OnLineFocusLost(e *Env, cIdx, sIdx int) {
	c.lflN++
	if c.onLfl == nil {
		return
	}
	c.onLfl(c, e, cIdx, sIdx)
}

func (s *lineFocusFeat) lfFX(
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

func (s *lineFocusFeat) Reports_and_focuses_first_focusable_line(t *T) {
	cmp := &cmpFX{
		onInit: func(c *cmpFX, e *Env) {
			c.FF.Add(LinesFocusable)
			fmt.Fprint(e, "first\nsecond")
		},
		onLineFocus: func(c *cmpFX, e *Env, i int) {
			switch c.N(onLineFocus) {
			case 1:
				t.Eq(0, c.LL.Focus.Current())
			case 2:
				t.Eq(1, c.LL.Focus.Current())
			}
		},
	}
	fx := fx(t, cmp)

	fx.FireKey(Down)

	t.FatalOn(fx.Lines.Update(cmp, nil, func(e *Env) {
		cmp.LL.Focus.Reset()
		cmp.LL.By(0).Flag(NotFocusable)
	}))

	fx.FireKey(Down)
	t.Eq(2, cmp.N(onLineFocus))
}

func (s *lineFocusFeat) Reports_line_focused_loss(t *T) {
	cmp := &cmpFX{
		onInit: func(c *cmpFX, e *Env) { // OnInit
			c.FF.Add(LinesFocusable)
			fmt.Fprint(e, "first\nsecond")
		},
		onLineFocusLost: func(lcf *cmpFX, e *Env, cIdx, sIdx int) {
			t.Eq(0, cIdx)
			t.Eq(sIdx, cIdx)
		},
	}
	tt := fx(t, cmp)
	tt.FireKey(Down)
	tt.FireKey(Down)
	t.Eq(1, cmp.N(onLineFocusLost))
}

func (s *lineFocusFeat) Reports_overflowing_line_on_focus(t *T) {
	cmp := &cmpFX{
		onInit: func(c *cmpFX, e *Env) { // OnInit
			c.FF.Add(LinesFocusable)
			fmt.Fprint(e, "an overflowing line")
		},
		onLineOverflowing: func(_ *cmpFX, _ *Env, left, right bool) {
			t.True(right)
			t.Not.True(left)
		},
	}
	fx := fx(t, cmp)
	fx.FireResize(len("an overflowing"), 1)
	fx.FireKey(Down)

	t.Eq(1, cmp.N(onLineOverflowing))
}

func (s *lineFocusFeat) Reports_overflowing_line_on_cursor_move(t *T) {
	cmp := &cmpFX{
		onInit: func(c *cmpFX, e *Env) { // OnInit
			c.FF.Add(CellFocusable)
			fmt.Fprint(e, "1234")
		},
		onLineOverflowing: func(c *cmpFX, _ *Env, left, right bool) {
			switch c.N(onLineOverflowing) {
			case 1:
				t.Not.True(left)
				t.True(right)
			case 2:
				t.True(left)
				t.True(right)
			case 3:
				t.True(left)
				t.Not.True(right)
			case 4:
				t.True(left)
				t.True(right)
			case 5:
				t.Not.True(left)
				t.True(right)
			case 6:
				t.True(left)
				t.Not.True(right)
			}
		},
	}
	fx := fx(t, cmp)
	fx.FireResize(2, 1)

	fx.FireKeys(Down, Right, Right)
	t.Eq(2, cmp.N(onLineOverflowing))

	fx.FireKey(Right)
	t.Eq(3, cmp.N(onLineOverflowing))

	fx.FireKeys(Left, Left)
	t.Eq(4, cmp.N(onLineOverflowing))

	fx.FireKey(Home)
	t.Eq(5, cmp.N(onLineOverflowing))

	fx.FireKey(End)
	t.Eq(6, cmp.N(onLineOverflowing))
}

func (s *lineFocusFeat) Reports_cursor_position_changes(t *T) {
	cmp := &cmpFX{
		onInit: func(c *cmpFX, e *Env) {
			c.FF.Add(CellFocusable)
			fmt.Fprint(e, "123")
		},
	}
	fx := fx(t, cmp)
	fx.FireResize(2, 1)

	fx.FireKeys(Down, Right)
	t.Eq(2, cmp.N(onCursor))
	t.Eq(fx.Screen(), "12")

	fx.FireKey(Right)
	t.Eq(2, cmp.N(onCursor))
	t.Eq(fx.Screen(), "23")

	fx.FireKeys(Home, Home, Left)
	t.Eq(3, cmp.N(onCursor))
	t.Eq(fx.Screen(), "12")

	fx.FireKeys(End, End, Right)
	t.Eq(4, cmp.N(onCursor))
	t.Eq(fx.Screen(), "23")
}

func (s *lineFocusFeat) Focuses_next_focusable_line(t *T) {
	cmp := &cmpFX{
		onInit: func(c *cmpFX, e *Env) {
			c.FF.Add(LinesFocusable)
			fmt.Fprint(e, "first\nsecond\nthird")
			c.LL.By(1).Flag(NotFocusable)
		},
		onLineFocus: func(c *cmpFX, e *Env, idx int) {
			switch c.N(onLineFocus) {
			case 1:
				t.Eq(0, c.LL.Focus.Current())
			case 2:
				t.Eq(2, c.LL.Focus.Current())
			}
		},
	}
	tt := fx(t, cmp)

	tt.FireKey(Down)
	tt.FireKey(Down)

	t.Eq(2, cmp.N(onLineFocus))
}

func (s *lineFocusFeat) Resets_if_no_next_focusable(t *T) {
	cmp := &lfCmpFX{
		onInit: func(c *lfCmpFX, e *Env) { // OnInit
			fmt.Fprint(e.LL(0), "first")
		},
		onLf: func(c *lfCmpFX, e *Env, i int) { // OnLineFocus
			t.Eq(0, c.LL.Focus.Current())
		},
	}
	fx := fx(t, cmp)

	fx.FireKey(Down)
	fx.FireKey(Down)

	t.FatalOn(fx.Lines.Update(cmp, nil, func(e *Env) {
		t.Eq(-1, cmp.LL.Focus.Current())
		fmt.Fprint(e.LL(1), "second")
		cmp.LL.By(1).Flag(NotFocusable)
		fmt.Fprint(e.LL(2), "third")
		cmp.LL.By(2).Flag(NotFocusable)
	}))

	fx.FireKey(Down)
	fx.FireKey(Down)

	t.FatalOn(fx.Lines.Update(cmp, nil, func(e *Env) {
		t.Eq(-1, cmp.LL.Focus.Current())
	}))
	t.Eq(2, cmp.lfN)
}

func (s *lineFocusFeat) Focuses_previous_focusable_line(t *T) {
	cmp := &lfCmpFX{
		onInit: func(lcf *lfCmpFX, e *Env) { // OnInit
			fmt.Fprint(e, "first\nsecond")
		},
		onLf: func(c *lfCmpFX, e *Env, i int) { // OnLineFocused
			switch c.lfN {
			case 1:
				t.Eq(0, c.LL.Focus.Current())
			case 3, 5, 7:
				t.Eq(1, c.LL.Focus.Current())
			case 6:
				t.Eq(2, c.LL.Focus.Current())
			}
		},
	}
	tt := fx(t, cmp)

	tt.FireKey(Down) // case 1
	tt.FireKey(Up)

	t.FatalOn(tt.Lines.Update(cmp, nil, func(e *Env) {
		t.Eq(-1, cmp.LL.Focus.Current())
		fmt.Fprint(e.LL(0), "first")
		cmp.LL.By(0).Flag(NotFocusable)
		fmt.Fprint(e.LL(1), "second")
		fmt.Fprint(e.LL(2), "third")
	}))

	tt.FireKey(Down) // case 3
	tt.FireKey(Up)
	t.FatalOn(tt.Lines.Update(cmp, nil, func(e *Env) {
		t.Eq(-1, cmp.LL.Focus.Current())
	}))
	tt.FireKey(Down) // case 5
	tt.FireKey(Down) // case 6
	tt.FireKey(Up)   // case 7

	t.Eq(5, cmp.lfN)
}

func (s *lineFocusFeat) Reset_triggered_by_unfocusable_feature(t *T) {
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

	tt.FireKey(Down)
	tt.FireKey(Esc)

	t.Eq(1, fx.lfN)
}

func (s *lineFocusFeat) Scrolls_to_next_highlighted_line(t *T) {
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

	t.Eq("line 5\nline 6", tt.ScreenOf(fx).Trimmed())

	t.Eq(1, fx.lfN)
}

func (s *lineFocusFeat) Scrolls_to_previous_highlighted_line(t *T) {
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

func (s *lineFocusFeat) Inverts_bg_fg_of_focused_if_highlighted(t *T) {
	cmp := &lfCmpFX{
		onInit: func(c *lfCmpFX, e *Env) {
			c.FF.Add(LinesHighlightedFocusable)
			fmt.Fprint(e.LL(0), "line 1")
		},
		onLf: func(c *lfCmpFX, e *Env, i int) {
			t.Eq(0, c.LL.Focus.Current())
		},
	}
	tt := fx(t, cmp)
	tt.FireResize(len("line n"), 1)

	l1 := tt.Cells()[0]
	for x := range l1 {
		t.Not.True(l1.HasAA(x, Reverse))
	}

	tt.FireKey(Down)
	t.Eq(1, cmp.lfN)

	l1 = tt.Cells()[0]
	for x := range l1 {
		t.True(l1.HasAA(x, Reverse))
	}
}

func (s *lineFocusFeat) Moves_highlight_to_next_focused_line(t *T) {
	cmp := &lfCmpFX{
		onInit: func(c *lfCmpFX, e *Env) {
			c.FF.Add(LinesHighlightedFocusable)
			fmt.Fprint(e.LL(0), "line 1")
			fmt.Fprint(e.LL(1), "line 2")
			c.LL.By(1).Flag(NotFocusable)
			fmt.Fprint(e.LL(2), "line 3")
		},
		onLf: func(c *lfCmpFX, e *Env, i int) {
			switch c.lfN {
			case 1:
				t.Eq(0, c.LL.Focus.Current())
			case 2:
				t.Eq(2, c.LL.Focus.Current())
			}
		},
	}
	tt := fx(t, cmp, 20*time.Minute)
	tt.FireResize(len("line n"), 2)
	tt.FireKey(Down)
	tt.FireKey(Down)

	t.Eq(2, cmp.lfN)
	t.Eq("line 2\nline 3", tt.Screen())

	l1 := tt.Cells()[0]
	for x := range l1 {
		t.FatalIfNot(t.Not.True(l1.HasAA(x, Reverse)))
	}
	l3 := tt.Cells()[1]
	for x := range l3 {
		t.FatalIfNot(t.True(l3.HasAA(x, Reverse)))
	}
}

func (s *lineFocusFeat) Trims_highlight_to_non_blanks(t *T) {
	var dflSty, hiSty Style
	cmp := &cmpFX{
		onInit: func(c *cmpFX, e *Env) {
			c.FF.Add(LinesHighlightedFocusable)
			c.LL.Focus.Trimmed()
			fmt.Fprint(e.LL(0), " no_blanks ")
			dflSty = c.Globals().Style(Default)
			hiSty = dflSty.WithAA(Reverse)
		},
	}
	fx := fx(t, cmp)
	fx.FireResize(len(" no_blanks "), 1)
	fx.FireKey(Down)
	l1 := fx.Cells()[0]
	for i, x := range l1 {
		switch i {
		case 0, 10:
			x.Style.Equals(dflSty)
		default:
			x.Style.Equals(hiSty)
		}
	}
	fx.FireKey(Down)
	fx.Lines.Update(cmp, nil, func(e *Env) {
		cmp.LL.Focus.Trimmed()
	})
	fx.FireKey(Down)
	l1 = fx.Cells()[0]
	for _, x := range l1 {
		x.Style.Equals(hiSty)
	}
}

type lsCmpFX struct {
	Component
	onIN func(*lsCmpFX, *Env)
	// onLF := OnLineFocus
	onLF func(*lsCmpFX, *Env, int)
	// onLFL := OnLineFocusLost
	onLFL func(_ *lsCmpFX, _ *Env, cIdx, sIdx int)
	// onLS := OnLineSelection
	onLS     func(*lsCmpFX, *Env, int)
	lfN, lsN int
}

func (c *lsCmpFX) OnInit(e *Env) {
	if c.onIN == nil {
		return
	}
	c.onIN(c, e)
}

func (c *lsCmpFX) OnLineFocus(e *Env, cIdx, sIdx int) {
	c.lfN++
	if c.onLF == nil {
		return
	}
	c.onLF(c, e, cIdx)
}

func (c *lsCmpFX) OnLineSelection(e *Env, cIdx, sIdx int) {
	c.lsN++
	if c.onLS == nil {
		return
	}
	c.onLS(c, e, cIdx)
}

func (s *lineFocusFeat) slFX(t *T,
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

func (s *lineFocusFeat) Reports_focused_line_on_line_selection(t *T) {
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

	tt.FireKey(Down)
	tt.FireKey(Down)
	tt.FireKey(Enter)

	t.Eq(2, fx.lfN)
	t.Eq(1, fx.lsN)
}

func (s *lineFocusFeat) Has_cursor_on_focused_line_if_cell_focusable(t *T) {
	cmp := &cmpFX{onInit: func(cf *cmpFX, e *Env) {
		cf.FF.Add(CellFocusable)
		fmt.Fprint(e, "first line")
	}}
	tt := fx(t, cmp)
	testCursor := func() {
		tt.Lines.Update(cmp, nil, func(e *Env) {
			ln, cl, haveCursor := cmp.CursorPosition()
			t.True(haveCursor)
			t.Eq(0, ln)
			t.Eq(0, cl)
		})
	}
	tt.FireKey(Down)
	t.Eq(cmp, tt.Lines.CursorComponent())
	testCursor()
	tt.FireKey(Down)
	tt.Lines.Update(cmp, nil, func(e *Env) {
		_, _, haveCursor := cmp.CursorPosition()
		t.Eq(false, haveCursor)
	})
	tt.FireKey(Up)
	t.Eq(cmp, tt.Lines.CursorComponent())
	testCursor()
}

func (s *lineFocusFeat) Moves_cursor_on_next_cell_focus_feature(t *T) {
	cmp := &cmpFX{onInit: func(cf *cmpFX, e *Env) {
		cf.FF.Add(CellFocusable)
		fmt.Fprint(e, "first line")
	}}
	tt := fx(t, cmp)
	tt.FireKey(Down)
	tt.FireKey(Right)
	tt.Lines.Update(cmp, nil, func(e *Env) {
		ln, cl, haveCursor := cmp.CursorPosition()
		t.True(haveCursor)
		t.Eq(0, ln)
		t.Eq(1, cl)
	})
}

func (s *lineFocusFeat) Moves_content_if_cursor_at_border(t *T) {
	cmp := &cmpFX{onInit: func(cf *cmpFX, e *Env) {
		cf.FF.Add(CellFocusable)
		fmt.Fprint(e, "123456")
		cf.Dim().SetWidth(4).SetHeight(1)
	}}
	tt := fx(t, cmp)
	t.Eq("1234", tt.ScreenOf(cmp))
	tt.FireKey(Down)
	for i := 0; i < 3; i++ {
		tt.FireKey(Right)
	}
	t.Eq("1234", tt.ScreenOf(cmp))
	tt.Lines.Update(cmp, nil, func(e *Env) {
		_, cl, haveCursor := cmp.CursorPosition()
		t.True(haveCursor)
		t.Eq(3, cl)
	})
	tt.FireKey(Right)
	t.Eq("2345", tt.ScreenOf(cmp))

	for i := 0; i < 4; i++ {
		tt.FireKey(Left)
	}
	t.Eq("1234", tt.ScreenOf(cmp))
}

func (s *lineFocusFeat) Moves_content_and_cursor_on_end_cell_feature(
	t *T,
) {
	cmp := &cmpFX{onInit: func(cf *cmpFX, e *Env) {
		cf.FF.Add(CellFocusable)
		fmt.Fprint(e, "12345")
		cf.Dim().SetWidth(6).SetHeight(1)
	}}
	tt := fx(t, cmp)
	t.Eq("12345 ", tt.ScreenOf(cmp))
	tt.FireKeys(Down, End)
	tt.Lines.Update(cmp, nil, func(e *Env) {
		_, cl, haveCursor := cmp.CursorPosition()
		t.True(haveCursor)
		t.Eq(4, cl)
	})
	tt.FireKey(Home)
	tt.Lines.Update(cmp, nil, func(e *Env) {
		cmp.Dim().SetWidth(4)
		_, cl, _ := cmp.CursorPosition()
		t.Eq(0, cl)
	})
	t.Eq("1234", tt.ScreenOf(cmp))
	tt.FireKey(End)
	tt.Lines.Update(cmp, nil, func(e *Env) {
		_, cl, haveCursor := cmp.CursorPosition()
		t.True(haveCursor)
		t.Eq(3, cl)
	})
	t.Eq("2345", tt.ScreenOf(cmp))
}

func (s *lineFocusFeat) Moves_content_and_cursor_on_home_cell_feature(
	t *T,
) {
	cmp := &cmpFX{onInit: func(cf *cmpFX, e *Env) {
		cf.FF.Add(CellFocusable)
		fmt.Fprint(e, "12345")
		cf.Dim().SetWidth(6).SetHeight(1)
	}}
	tt := fx(t, cmp)
	t.Eq("12345 ", tt.ScreenOf(cmp))
	tt.FireKey(Down)
	tt.FireKey(End)
	tt.FireKey(Home)
	tt.Lines.Update(cmp, nil, func(e *Env) {
		_, cl, haveCursor := cmp.CursorPosition()
		t.True(haveCursor)
		t.Eq(0, cl)
		cmp.Dim().SetWidth(4)
	})
	t.Eq("1234", tt.ScreenOf(cmp))
	tt.FireKey(End)
	t.Eq("2345", tt.ScreenOf(cmp))
	tt.FireKey(Home)
	tt.Lines.Update(cmp, nil, func(e *Env) {
		_, cl, haveCursor := cmp.CursorPosition()
		t.True(haveCursor)
		t.Eq(0, cl)
	})
	t.Eq("1234", tt.ScreenOf(cmp))
}

func (s *lineFocusFeat) Moves_cursor_on_line_focusable_features(t *T) {
	cmp := &cmpFX{onInit: func(cf *cmpFX, e *Env) {
		cf.FF.Add(CellFocusable)
		fmt.Fprint(e, "123456\n1234\n123456")
		cf.Dim().SetWidth(6).SetHeight(3)
	}}
	tt := fx(t, cmp)
	testCursor := func(ln, cl int) {
		tt.Lines.Update(cmp, nil, func(e *Env) {
			_ln, _cl, haveCursor := cmp.CursorPosition()
			t.True(haveCursor)
			t.Eq(ln, _ln)
			t.Eq(cl, _cl)
		})
	}
	tt.FireKey(Down)
	tt.FireKey(End)
	testCursor(0, 5)
	tt.FireKey(Down)
	testCursor(1, 3)
	tt.FireKey(Down)
	testCursor(2, 3)
	tt.FireKey(End)
	testCursor(2, 5)
	tt.FireKey(Up)
	testCursor(1, 3)
	tt.FireKey(Up)
	testCursor(0, 3)
}

func (s *lineFocusFeat) No_ops_on_cell_focusable_features(t *T) {
	cmp := &cmpFX{onInit: func(cf *cmpFX, e *Env) {
		cf.FF.Add(CellFocusable)
		fmt.Fprint(e, "12345")
		cf.Dim().SetWidth(3).SetHeight(1)
	}}
	tt := fx(t, cmp)
	testCursor := func(ln, cl int, have bool) {
		tt.Lines.Update(cmp, nil, func(e *Env) {
			_ln, _cl, _have := cmp.CursorPosition()
			t.Eq(ln, _ln)
			t.Eq(cl, _cl)
			t.Eq(have, _have)
		})
	}
	ln, cl, have := 0, 0, false
	tt.Lines.Update(cmp, nil, func(e *Env) {
		ln, cl, have = cmp.CursorPosition()
	})
	tt.FireKey(Right)
	testCursor(ln, cl, have)
	tt.FireKey(Left)
	testCursor(ln, cl, have)
	tt.FireKey(End)
	testCursor(ln, cl, have)
	tt.FireKey(Home)
	testCursor(ln, cl, have)
	tt.FireKey(Down)
	tt.Lines.Update(cmp, nil, func(e *Env) {
		ln, cl, have = cmp.CursorPosition()
	})
	tt.FireKey(Left)
	testCursor(ln, cl, have)
	tt.FireKey(Home)
	testCursor(ln, cl, have)
	tt.FireKey(End)
	tt.Lines.Update(cmp, nil, func(e *Env) {
		ln, cl, have = cmp.CursorPosition()
	})
	tt.FireKey(Right)
	t.Eq("345", tt.ScreenOf(cmp)) // catched an incrementStart bug
	testCursor(ln, cl, have)
	tt.FireKey(End)
	testCursor(ln, cl, have)
}

func (s lineFocusFeat) Removes_cursor_if_line_focus_resets(t *T) {
	cmp := &cmpFX{onInit: func(cf *cmpFX, e *Env) {
		cf.FF.Add(CellFocusable)
		fmt.Fprint(e, "12345\n12345")
		cf.Dim().SetWidth(6).SetHeight(2)
	}}
	tt := fx(t, cmp)
	haveCursor := func(b bool, n int) {
		tt.Lines.Update(cmp, nil, func(e *Env) {
			_, _, have := cmp.CursorPosition()
			t.Eq(b, have)
		})
		t.Eq(n, cmp.N(onCursor))
	}
	haveCursor(false, 0)
	tt.FireKey(Down)
	haveCursor(true, 1)
	tt.FireKey(Up)
	haveCursor(false, 2)
	for i := 0; i < 3; i++ {
		tt.FireKey(Down)
	}
	haveCursor(false, 5)
}

func (s lineFocusFeat) Resets_line_start_on_line_focus_change(t *T) {
	cmp := &cmpFX{onInit: func(cf *cmpFX, e *Env) {
		cf.FF.Add(CellFocusable)
		fmt.Fprint(e, "12345\n12345")
		cf.Dim().SetWidth(3).SetHeight(2)
	}}
	tt := fx(t, cmp)
	tt.FireKey(Down)
	tt.FireKey(End)
	t.Eq("345\n123", tt.ScreenOf(cmp))
	tt.FireKey(Down)
	t.Eq("123\n123", tt.ScreenOf(cmp))
	tt.FireKey(End)
	t.Eq("123\n345", tt.ScreenOf(cmp))
	tt.FireKey(Up)
	t.Eq("123\n123", tt.ScreenOf(cmp))
}

func TestLineFocus(t *testing.T) {
	t.Parallel()
	Run(&lineFocusFeat{}, t)
}
