// Copyright (c) 2022 Stephan Lukits. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package lines

import (
	"fmt"
	"testing"

	. "github.com/slukits/gounit"
)

type cellFocus struct{ Suite }

func (s *cellFocus) SetUp(t *T) { t.Parallel() }

func (s *cellFocus) Has_cursor_on_focused_line_if_cell_focusable(t *T) {
	cmp := &cmpFX{onInit: func(cf *cmpFX, e *Env) {
		cf.FF.Set(CellFocusable)
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
		t.Not.True(haveCursor)
	})
	tt.FireKey(Up)
	t.Eq(cmp, tt.Lines.CursorComponent())
	testCursor()
}

func (s *cellFocus) Moves_cursor_on_next_cell_focus_feature(t *T) {
	cmp := &cmpFX{onInit: func(cf *cmpFX, e *Env) {
		cf.FF.Set(CellFocusable)
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

func (s *cellFocus) Moves_content_if_cursor_at_border(t *T) {
	cmp := &cmpFX{onInit: func(cf *cmpFX, e *Env) {
		cf.FF.Set(CellFocusable)
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

func (s *cellFocus) Moves_content_and_cursor_on_end_cell_feature(
	t *T,
) {
	cmp := &cmpFX{onInit: func(cf *cmpFX, e *Env) {
		cf.FF.Set(CellFocusable)
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

func (s *cellFocus) Moves_content_and_cursor_on_home_cell_feature(
	t *T,
) {
	cmp := &cmpFX{onInit: func(cf *cmpFX, e *Env) {
		cf.FF.Set(CellFocusable)
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

func (s *cellFocus) Moves_cursor_on_line_focusable_features(t *T) {
	cmp := &cmpFX{onInit: func(cf *cmpFX, e *Env) {
		cf.FF.Set(CellFocusable)
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

func (s *cellFocus) No_ops_on_cell_focusable_features(t *T) {
	cmp := &cmpFX{onInit: func(cf *cmpFX, e *Env) {
		cf.FF.Set(CellFocusable)
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

func (s *cellFocus) Removes_cursor_if_line_focus_resets(t *T) {
	cmp := &cmpFX{onInit: func(cf *cmpFX, e *Env) {
		cf.FF.Set(CellFocusable)
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

func (s *cellFocus) test_cursor_pos(
	fx *Fixture, t *T, cmp *cmpFX, slIdx, scIdx int,
) {
	fx.Lines.Update(cmp, nil, func(e *Env) {
		slIdx, scIdx, hasCursor := cmp.CursorPosition()
		t.FatalIfNot(t.True((hasCursor)))
		t.Eq(slIdx, slIdx)
		t.Eq(scIdx, scIdx)
	})
}

func (s *cellFocus) Sets_cursor_after_last_rune(t *T) {
	cmp := &cmpFX{onInit: func(cf *cmpFX, e *Env) {
		cf.FF.Set(CellFocusable)
		cf.LL.Focus.EolAfterLastRune()
		fmt.Fprint(e, "12345\n6789\n12345")
		cf.Dim().SetWidth(3).SetHeight(2)
	}}

	fx := fx(t, cmp)
	fx.FireKeys(Down, End, End, Right)
	t.Eq(fx.ScreenOf(cmp)[0], "45 ")
	s.test_cursor_pos(fx, t, cmp, 0, 2)

	fx.FireKey(Down)
	s.test_cursor_pos(fx, t, cmp, 1, 2)
	fx.FireKey(Right)
	s.test_cursor_pos(fx, t, cmp, 1, 2)
	t.Eq(fx.ScreenOf(cmp)[1], "789")
	fx.FireKey(Right)
	s.test_cursor_pos(fx, t, cmp, 1, 2)
	t.Eq(fx.ScreenOf(cmp)[1], "89 ")
	fx.FireKey(Right)
	s.test_cursor_pos(fx, t, cmp, 1, 2)
	t.Eq(fx.ScreenOf(cmp)[1], "89 ")
}

func TestCellFocus(t *testing.T) {
	t.Parallel()
	Run(&cellFocus{}, t)
}
