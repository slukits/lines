// Copyright (c) 2022 Stephan Lukits. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package lines

import (
	"fmt"
	"testing"

	. "github.com/slukits/gounit"
)

type lineFocus struct{ Suite }

func (s *lineFocus) SetUp(t *T) { t.Parallel() }

func (s *lineFocus) Has_initially_no_line_focused(t *T) {
	cmp := &cmpFX{
		onInit: func(c *cmpFX, e *Env) {
			t.Eq(-1, c.LL.Focus.Current())
			// NoOps for coverage
			c.LL.Focus.Next()
			c.LL.Focus.Reset()
			t.Eq(-1, c.LL.Focus.Current())
		}}
	fx(t, cmp)
}

func (s *lineFocus) Focuses_first_focusable_line_on_down_key(t *T) {
	cmp := &cmpFX{
		onInit: func(c *cmpFX, e *Env) {
			c.FF.Set(LinesFocusable)
			fmt.Fprint(e, "first\nsecond")
		},
		onLineFocus: func(c *cmpFX, e *Env, cIdx, sIdx int) {
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

func (s *lineFocus) Reports_line_focus_loss(t *T) {
	cmp := &cmpFX{
		onInit: func(c *cmpFX, e *Env) { // OnInit
			c.FF.Set(LinesFocusable)
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

func (s *lineFocus) Reports_overflowing_line_on_focus(t *T) {
	cmp := &cmpFX{
		onInit: func(c *cmpFX, e *Env) { // OnInit
			c.FF.Set(LinesFocusable)
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

func (s *lineFocus) Reports_overflowing_line_on_cursor_move(t *T) {
	cmp := &cmpFX{
		onInit: func(c *cmpFX, e *Env) { // OnInit
			c.FF.Set(CellFocusable)
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

func (s *lineFocus) Reports_cursor_position_changes(t *T) {
	cmp := &cmpFX{
		onInit: func(c *cmpFX, e *Env) {
			c.FF.Set(CellFocusable)
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

func (s *lineFocus) Focuses_next_focusable_line_on_down_key(t *T) {
	cmp := &cmpFX{
		onInit: func(c *cmpFX, e *Env) {
			c.FF.Set(LinesFocusable)
			fmt.Fprint(e, "first\nsecond\nthird")
			c.LL.By(1).Flag(NotFocusable)
		},
		onLineFocus: func(c *cmpFX, e *Env, cIdx, sIdx int) {
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

func (s *lineFocus) Resets_if_no_next_focusable(t *T) {
	cmp := &cmpFX{
		onInit: func(c *cmpFX, e *Env) {
			c.FF.Set(LinesFocusable)
			fmt.Fprint(e.LL(0), "first")
		},
		onLineFocus: func(c *cmpFX, e *Env, cIdx, sIdx int) {
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
	t.Eq(2, cmp.cc[onLineFocus])
	t.Eq(2, cmp.cc[onLineFocusLost])
}

func (s *lineFocus) Focuses_previous_focusable_line_on_up_key(t *T) {
	cmp := &cmpFX{
		onInit: func(c *cmpFX, e *Env) {
			c.FF.Set(LinesFocusable)
			fmt.Fprint(e, "first\nsecond")
		},
		onLineFocus: func(c *cmpFX, e *Env, cIdx, sIdx int) {
			switch c.cc[onLineFocus] {
			case 1:
				t.Eq(0, c.LL.Focus.Current())
			case 2, 3, 5:
				t.Eq(1, c.LL.Focus.Current())
			case 4:
				t.Eq(2, c.LL.Focus.Current())
			}
		},
		onLineFocusLost: func(c *cmpFX, e *Env, cIdx, sIdx int) {
			switch c.cc[onLineFocusLost] {
			case 1:
				t.Eq(0, cIdx)
				t.Eq(-1, c.LL.Focus.Current())
			case 2:
				t.Eq(1, cIdx)
				t.Eq(-1, c.LL.Focus.Current())
			case 3:
				t.Eq(1, cIdx)
				t.Eq(2, c.LL.Focus.Current())
			case 4:
				t.Eq(2, cIdx)
				t.Eq(1, c.LL.Focus.Current())
			}
		},
	}
	tt := fx(t, cmp)

	tt.FireKey(Down)
	tt.FireKey(Up)
	t.True(cmp.cc[onLineFocus] == 1 && cmp.cc[onLineFocusLost] == 1)

	t.FatalOn(tt.Lines.Update(cmp, nil, func(e *Env) {
		t.Eq(-1, cmp.LL.Focus.Current())
		fmt.Fprint(e.LL(0), "first")
		cmp.LL.By(0).Flag(NotFocusable)
		fmt.Fprint(e.LL(1), "second")
		fmt.Fprint(e.LL(2), "third")
	}))

	tt.FireKey(Down)
	tt.FireKey(Up)
	t.True(cmp.cc[onLineFocus] == 2 && cmp.cc[onLineFocusLost] == 2)

	tt.FireKey(Down)
	t.True(cmp.cc[onLineFocus] == 3 && cmp.cc[onLineFocusLost] == 2)

	tt.FireKey(Down)
	t.True(cmp.cc[onLineFocus] == 4 && cmp.cc[onLineFocusLost] == 3)

	tt.FireKey(Up)
	t.True(cmp.cc[onLineFocus] == 5 && cmp.cc[onLineFocusLost] == 4)
}

func (s *lineFocus) Reset_triggered_by_unfocusable_feature(t *T) {
	cmp := &cmpFX{
		onInit: func(c *cmpFX, e *Env) {
			c.FF.Set(LinesFocusable)
			fmt.Fprint(e, "first\nsecond")
		},
		onLineFocus: func(c *cmpFX, e *Env, cIdx, sIdx int) {
			t.Eq(0, c.LL.Focus.Current())
		},
		onLineFocusLost: func(c *cmpFX, e *Env, cIdx, sIdx int) {
			t.Eq(-1, c.LL.Focus.Current())
		},
	}
	tt := fx(t, cmp)

	tt.FireKey(Down)
	t.Eq(1, cmp.cc[onLineFocus])
	tt.FireKey(Esc)
	t.Eq(1, cmp.cc[onLineFocus])
	t.Eq(1, cmp.cc[onLineFocusLost])
}

func (s *lineFocus) Scrolls_to_next_focusable_line(t *T) {
	cmp := &cmpFX{
		onInit: func(c *cmpFX, e *Env) {
			c.FF.Set(LinesFocusable)
			c.dim.SetHeight(2)
			for i := 0; i < 7; i++ {
				fmt.Fprintf(e.LL(i), "line %d", i+1)
			}
			for i := 0; i < 5; i++ {
				c.LL.By(i).Flag(NotFocusable)
			}
		},
		onLineFocus: func(c *cmpFX, e *Env, _, _ int) {
			t.Eq(5, c.LL.Focus.Current())
		},
	}
	fx := fx(t, cmp)
	fx.FireKey(Down)

	t.Eq("line 5\nline 6", fx.ScreenOf(cmp).Trimmed())

	t.Eq(1, cmp.cc[onLineFocus])
}

func (s *lineFocus) Scrolls_to_previous_focusable_line(t *T) {
	cmp := &cmpFX{
		onInit: func(c *cmpFX, e *Env) {
			c.FF.Set(LinesFocusable)
			c.dim.SetHeight(2)
			for i := 0; i < 7; i++ {
				fmt.Fprintf(e.LL(i), "line %d", i+1)
			}
			for _, idx := range []int{0, 2, 3, 4} {
				c.LL.By(idx).Flag(NotFocusable)
			}
		},
		onLineFocus: func(c *cmpFX, e *Env, _, _ int) {
			switch c.cc[onLineFocus] {
			case 1, 3:
				t.Eq(1, c.LL.Focus.Current())
			case 2:
				t.Eq(5, c.LL.Focus.Current())
			}
		},
	}
	fx := fx(t, cmp)
	fx.FireKey(Down)
	fx.FireKey(Down)

	t.Eq("line 5\nline 6", fx.ScreenOf(cmp).Trimmed().String())

	fx.FireKey(Up)

	t.Eq("line 2\nline 3", fx.ScreenOf(cmp).Trimmed().String())
	t.Eq(3, cmp.cc[onLineFocus])
}

func (s *lineFocus) Inverts_bg_fg_of_focused_if_highlighted(t *T) {
	cmp := &cmpFX{
		onInit: func(c *cmpFX, e *Env) {
			c.FF.Set(LinesFocusable | HighlightEnabled)
			fmt.Fprint(e.LL(0), "line 1")
		},
		onLineFocus: func(c *cmpFX, e *Env, _, _ int) {
			t.Eq(0, c.LL.Focus.Current())
		},
	}
	fx := fx(t, cmp)
	fx.FireResize(len("line n"), 1)
	l1 := fx.Cells()[0]
	for x := range l1 {
		t.Not.True(l1.HasAA(x, Reverse))
	}

	fx.FireKey(Down)
	t.Eq(1, cmp.cc[onLineFocus])
	l1 = fx.Cells()[0]
	for x := range l1 {
		t.True(l1.HasAA(x, Reverse))
	}
}

func (s *lineFocus) Moves_highlight_to_next_focused_line(t *T) {
	cmp := &cmpFX{
		onInit: func(c *cmpFX, e *Env) {
			c.FF.Set(LinesFocusable | HighlightEnabled)
			fmt.Fprint(e.LL(0), "line 1")
			fmt.Fprint(e.LL(1), "line 2")
			c.LL.By(1).Flag(NotFocusable)
			fmt.Fprint(e.LL(2), "line 3")
		},
		onLineFocus: func(c *cmpFX, e *Env, _, _ int) {
			switch c.cc[onLineFocus] {
			case 1:
				t.Eq(0, c.LL.Focus.Current())
			case 2:
				t.Eq(2, c.LL.Focus.Current())
			}
		},
	}
	tt := fx(t, cmp)
	tt.FireResize(len("line n"), 2)
	tt.FireKeys(Down, Down)

	t.Eq(2, cmp.cc[onLineFocus])
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

func (s *lineFocus) Removes_highlight_if_unfocused(t *T) {
	cmp := &cmpFX{
		onInit: func(c *cmpFX, e *Env) {
			c.FF.Set(LinesFocusable | TrimmedHighlightEnabled)
			fmt.Fprint(e.LL(0), " no_blanks ")
		},
	}
	fx := fx(t, cmp)
	fx.FireResize(len(" no_blanks "), 1)
	fx.Lines.Update(cmp, nil, func(e *Env) {
		t.Eq(-1, cmp.LL.Focus.Current())
		t.Not.True(cmp.LL.By(0).IsFlagged(TrimmedHighlighted))
	})
	fx.FireKey(Down)
	fx.Lines.Update(cmp, nil, func(e *Env) {
		t.Eq(0, cmp.LL.Focus.Current())
		t.True(cmp.LL.By(0).IsFlagged(TrimmedHighlighted))
	})
	fx.FireKey(Down)
	fx.Lines.Update(cmp, nil, func(e *Env) {
		t.Eq(-1, cmp.LL.Focus.Current())
		t.Not.True(cmp.LL.By(0).IsFlagged(TrimmedHighlighted))
		cmp.FF.Set(HighlightEnabled)
	})
	fx.FireKey(Down)
	fx.Lines.Update(cmp, nil, func(e *Env) {
		t.Eq(0, cmp.LL.Focus.Current())
		t.True(cmp.LL.By(0).IsFlagged(Highlighted))
	})
}

func (s *lineFocus) Trims_highlight_to_non_blanks(t *T) {
	var dflSty, hiSty Style
	cmp := &cmpFX{
		onInit: func(c *cmpFX, e *Env) {
			c.FF.Set(LinesFocusable | TrimmedHighlightEnabled)
			fmt.Fprint(e.LL(0), " no_blanks ")
			dflSty = c.Globals().Style(Default)
			hiSty = dflSty.WithAA(Reverse)
		},
	}
	fx := fx(t, cmp)
	fx.FireResize(len(" no_blanks "), 1)
	fx.Lines.Update(cmp, nil, func(e *Env) {
		t.Not.True(cmp.LL.By(0).IsFlagged(TrimmedHighlighted))
	})
	fx.FireKey(Down)
	fx.Lines.Update(cmp, nil, func(e *Env) {
		t.True(cmp.LL.By(0).IsFlagged(TrimmedHighlighted))
	})
	l1 := fx.Cells()[0]
	for i, x := range l1 {
		switch i {
		case 0, 10:
			t.Eq(x.Style, dflSty)
		default:
			t.Eq(x.Style, hiSty)
		}
	}
	fx.FireKey(Down)
	fx.Lines.Update(cmp, nil, func(e *Env) {
		cmp.FF.Set(HighlightEnabled)
	})
	fx.FireKey(Down)
	l1 = fx.Cells()[0]
	for _, x := range l1 {
		t.Eq(x.Style, hiSty)
	}
}

func (s *lineFocus) Reports_focused_line_on_line_selection(t *T) {
	cmp := &cmpFX{
		onInit: func(c *cmpFX, e *Env) {
			c.FF.Set(LinesSelectable)
			fmt.Fprint(e, "first\nsecond")
		},
		onLineFocus: func(c *cmpFX, e *Env, cIdx, sIdx int) {
			switch c.cc[onLineFocus] {
			case 1:
				t.Eq(0, c.LL.Focus.Current())
			case 2:
				t.Eq(1, c.LL.Focus.Current())
			}
		},
		onLineSelection: func(c *cmpFX, e *Env, cIdx, _ int) {
			t.Eq(1, cIdx)
		},
	}
	tt := fx(t, cmp)

	tt.FireKey(Down)
	tt.FireKey(Down)
	tt.FireKey(Enter)

	t.Eq(2, cmp.cc[onLineFocus])
	t.Eq(1, cmp.cc[onLineSelection])
}

func (s *lineFocus) Resets_line_start_on_line_focus_change(t *T) {
	cmp := &cmpFX{onInit: func(cf *cmpFX, e *Env) {
		cf.FF.Set(CellFocusable)
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
	Run(&lineFocus{}, t)
}
