// Copyright (c) 2022 Stephan Lukits. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package lines

import (
	"fmt"
	"strings"
	"testing"
	"time"

	. "github.com/slukits/gounit"
)

type AComponent struct{ Suite }

func (s *AComponent) SetUp(t *T) { t.Parallel() }

func (s *AComponent) Access_panics_outside_event_processing(t *T) {
	_, cmp := fxCmp(t)
	t.Panics(func() { cmp.Dim().SetHeight(20) })
}

func cmpfx(t *T, d ...time.Duration) (*Fixture, *cmpFX) {
	cmp := &cmpFX{}
	var tt *Fixture
	if len(d) == 0 {
		tt = TermFixture(t.GoT(), 0, cmp)
	} else {
		tt = TermFixture(t.GoT(), d[0], cmp)
	}
	return tt, cmp
}

func xcmpfx(t *T, cmp Componenter, d ...time.Duration) *Fixture {
	if len(d) == 0 {
		return TermFixture(t.GoT(), 0, cmp)
	}
	return TermFixture(t.GoT(), d[0], cmp)
}

func (s *AComponent) Creates_needed_lines_on_write(t *T) {
	fx, cmp := fxCmp(t)
	t.FatalOn(fx.Lines.Update(cmp, nil, func(e *Env) {
		t.Eq(0, cmp.Len())
		fmt.Fprint(e, "first\nsecond\nthird")
		t.Eq(3, cmp.Len())
	}))
}

func (s *AComponent) Doesnt_change_line_count_on_line_overwrite(t *T) {
	fx, cmp := fxCmp(t)
	t.FatalOn(fx.Lines.Update(cmp, nil, func(e *Env) {
		cmp.LL.Mod(Overwriting)
		fmt.Fprint(e, "two\nlines")
		t.Eq(2, cmp.Len())
		fmt.Fprint(e, "one line")
		t.Eq(2, cmp.Len())
	}))

	// but second line is empty now
	t.Eq("one line", fx.Screen().Trimmed().String())
}

func (s *AComponent) Has_a_line_more_after_appending_an_line(t *T) {
	fx, cmp := fxCmp(t)
	t.FatalOn(fx.Lines.Update(cmp, nil, func(e *Env) {
		cmp.LL.Mod(Appending)
		fmt.Fprint(e, "two\nlines")
		fmt.Fprint(e, "one line")
		t.Eq(3, cmp.Len())
	}))
}

func (s *AComponent) Has_a_line_more_after_writing_to_tailing(t *T) {
	fx, cmp := fxCmp(t)
	t.FatalOn(fx.Lines.Update(cmp, nil, func(e *Env) {
		cmp.LL.Mod(Tailing)
		fmt.Fprint(e, "two\nlines")
		fmt.Fprint(e, "one line")
		t.Eq(3, cmp.Len())
	}))
}

func (s *AComponent) Shows_last_line_clipped_above_if_tailing(t *T) {
	fx, cmp := fxCmp(t)
	fx.FireResize(20, 2)
	t.FatalOn(fx.Lines.Update(cmp, nil, func(e *Env) {
		cmp.LL.Mod(Tailing)
		fmt.Fprint(e, "three\nlines\nat last")
	}))
	t.Eq("lines  \nat last", fx.Screen().Trimmed().String())
}

func (s *AComponent) Blanks_a_reset_line(t *T) {
	fx, cmp := fxCmp(t)
	fx.FireResize(20, 2)
	t.FatalOn(fx.Lines.Update(cmp, nil, func(e *Env) {
		fmt.Fprint(e, "first\nsecond")
	}))
	t.Eq("first \nsecond", fx.Screen().Trimmed().String())

	t.FatalOn(fx.Lines.Update(cmp, nil, func(e *Env) {
		cmp.Reset(-2) // no-op, coverage
		cmp.Reset(0)
	}))

	t.Eq("second", fx.Screen().Trimmed().String())
}

func (s *AComponent) Truncates_lines_to_screen_area_on_reset_all(t *T) {
	fx, cmp := fxCmp(t)
	fx.FireResize(20, 2)
	t.FatalOn(fx.Lines.Update(cmp, nil, func(e *Env) {
		fmt.Fprint(e, "first\nsecond\nthird\nforth")
		t.Eq(4, cmp.Len())
	}))
	t.FatalOn(fx.Lines.Update(cmp, nil, func(e *Env) {
		cmp.Reset(All)
		t.Eq(2, cmp.Len())
	}))
}

func (s *AComponent) Scrolls_by_one_line_if_height_is_one(t *T) {
	fx, cmp := fxCmp(t)
	t.FatalOn(fx.Lines.Update(cmp, nil, func(e *Env) {
		cmp.Dim().SetHeight(1)
		fmt.Fprint(e, "first\nsecond")
	}))
	t.Eq("first", fx.Screen().Trimmed())

	t.FatalOn(fx.Lines.Update(cmp, nil, func(e *Env) {
		cmp.Scroll.Down()
	}))
	t.Eq("second", fx.Screen().Trimmed())
	t.FatalOn(fx.Lines.Update(cmp, nil, func(e *Env) {
		cmp.Scroll.Down()
	}))
	t.Eq("second", fx.Screen().Trimmed())

	t.FatalOn(fx.Lines.Update(cmp, nil, func(e *Env) {
		cmp.Scroll.Up()
	}))
	t.Eq("first", fx.Screen().Trimmed())
	t.FatalOn(fx.Lines.Update(cmp, nil, func(e *Env) {
		cmp.Scroll.Up()
	}))
	t.Eq("first", fx.Screen().Trimmed())
}

func (s *AComponent) Scrolls_to_last_line_if_last_displayed(t *T) {
	fx, cmp := fxCmp(t)
	t.FatalOn(fx.Lines.Update(cmp, nil, func(e *Env) {
		cmp.Dim().SetHeight(3)
		fmt.Fprint(e, "first\nsecond\nthird\nforth")
	}))
	t.Eq("first \nsecond\nthird ", fx.Screen().Trimmed())

	t.FatalOn(fx.Lines.Update(cmp, nil, func(e *Env) {
		cmp.Scroll.Down()
	}))
	t.Eq("second\nthird \nforth ", fx.Screen().Trimmed())
}

func (s *AComponent) Scrolls_to_first_line_if_first_displayed(t *T) {
	fx, cmp := fxCmp(t)
	t.FatalOn(fx.Lines.Update(cmp, nil, func(e *Env) {
		cmp.Dim().SetHeight(3)
		fmt.Fprint(e, "first\nsecond\nthird\nforth")
	}))
	t.Eq("first \nsecond\nthird ", fx.Screen().Trimmed().String())
	t.FatalOn(fx.Lines.Update(cmp, nil, func(e *Env) {
		cmp.Scroll.Down()
	}))
	t.Eq("second\nthird \nforth ", fx.Screen().Trimmed().String())

	t.FatalOn(fx.Lines.Update(cmp, nil, func(e *Env) {
		cmp.Scroll.Up()
	}))
	t.Eq("first \nsecond\nthird ", fx.Screen().Trimmed().String())
}

func (s *AComponent) Scrolls_down_by_90_percent_height(t *T) {
	fx, cmp := fxCmp(t)
	fx.FireResize(20, 30)
	t.FatalOn(fx.Lines.Update(cmp, nil, func(e *Env) {
		cmp.Dim().SetHeight(5)
		ll := make([]string, 60)
		for i := 0; i < 60; i++ {
			ll[i] = fmt.Sprintf("line %d", i+1)
		}
		fmt.Fprint(e, strings.Join(ll, "\n"))
		t.Eq(60, cmp.Len())
	}))

	exp := []string{}
	for i := 0; i < 5; i++ {
		exp = append(exp, fmt.Sprintf("line %d", i+5))
	}
	t.FatalOn(fx.Lines.Update(cmp, nil, func(e *Env) {
		cmp.Scroll.Down()
	}))
	t.Eq(strings.Join(exp, "\n"), fx.Screen().Trimmed().String())

	t.FatalOn(fx.Lines.Update(cmp, nil, func(e *Env) {
		cmp.Dim().SetHeight(15)
	}))

	exp = []string{}
	for i := 0; i < 15; i++ { // first is still at fifth line
		exp = append(exp, fmt.Sprintf("line %d", i+19))
	}
	t.FatalOn(fx.Lines.Update(cmp, nil, func(e *Env) {
		cmp.Scroll.Down()
	}))
	t.Eq(strings.Join(exp, "\n"), fx.Screen().Trimmed().String())

	t.FatalOn(fx.Lines.Update(cmp, nil, func(e *Env) {
		cmp.Scroll.ToTop()
		t.True(cmp.Scroll.IsAtTop())
		cmp.Dim().SetHeight(30)
	}))
	exp = []string{}
	for i := 0; i < 30; i++ { // first is still at fifth line
		exp = append(exp, fmt.Sprintf("line %d", i+28))
	}
	t.FatalOn(fx.Lines.Update(cmp, nil, func(e *Env) {
		cmp.Scroll.Down()
	}))
	t.Eq(strings.Join(exp, "\n"), fx.Screen().Trimmed().String())
}

func (s *AComponent) Scrolls_up_by_90_percent_height(t *T) {
	fx, cmp := fxCmp(t)
	fx.FireResize(20, 30)
	t.FatalOn(fx.Lines.Update(cmp, nil, func(e *Env) {
		cmp.Dim().SetHeight(5)
		ll := make([]string, 60)
		for i := 0; i < 60; i++ {
			if i+1 < 10 {
				ll[i] = fmt.Sprintf("line 0%d", i+1)
				continue
			}
			ll[i] = fmt.Sprintf("line %d", i+1)
		}
		fmt.Fprint(e, strings.Join(ll, "\n"))
		t.Eq(60, cmp.Len())
		cmp.Scroll.ToBottom()
	}))

	exp := []string{}
	for i := 0; i < 5; i++ {
		exp = append(exp, fmt.Sprintf("line %d", i+56))
	}
	t.Eq(strings.Join(exp, "\n"), fx.Screen().Trimmed().String())
	exp = []string{}
	for i := 0; i < 5; i++ {
		exp = append(exp, fmt.Sprintf("line %d", i+52))
	}
	t.FatalOn(fx.Lines.Update(cmp, nil, func(e *Env) {
		cmp.Scroll.Up()
	}))
	t.Eq(strings.Join(exp, "\n"), fx.Screen().Trimmed().String())

	t.FatalOn(fx.Lines.Update(cmp, nil, func(e *Env) {
		cmp.Dim().SetHeight(15)
	}))

	exp = []string{}
	for i := 0; i < 15; i++ {
		exp = append(exp, fmt.Sprintf("line %d", i+32))
	}
	t.FatalOn(fx.Lines.Update(cmp, nil, func(e *Env) {
		cmp.Scroll.Up()
	}))
	t.Eq(strings.Join(exp, "\n"), fx.Screen().Trimmed().String())

	t.FatalOn(fx.Lines.Update(cmp, nil, func(e *Env) {
		cmp.Dim().SetHeight(30)
		cmp.Scroll.ToBottom()
		t.True(cmp.Scroll.IsAtBottom())
	}))
	exp = []string{}
	for i := 0; i < 30; i++ {
		if i+4 < 10 {
			exp = append(exp, fmt.Sprintf("line 0%d", i+4))
			continue
		}
		exp = append(exp, fmt.Sprintf("line %d", i+4))
	}
	t.FatalOn(fx.Lines.Update(cmp, nil, func(e *Env) {
		cmp.Scroll.Up()
	}))
	t.Eq(strings.Join(exp, "\n"), fx.Screen().Trimmed().String())
}

func (s *AComponent) Scrolls_to_top_on_reset_all(t *T) {
	fx, cmp := fxCmp(t)
	fx.FireResize(20, 2)
	t.FatalOn(fx.Lines.Update(cmp, nil, func(e *Env) {
		fmt.Fprint(e, "first\nsecond\nthird\nforth")
		cmp.Scroll.ToBottom()
		t.Not.True(cmp.Scroll.IsAtTop())
	}))
	t.FatalOn(fx.Lines.Update(cmp, nil, func(e *Env) {
		cmp.Reset(All)
		t.True(cmp.Scroll.IsAtTop())
	}))
}

func (s *AComponent) Is_replaceable(t *T) {
	cmp, long := &stackingFX{}, "a rather long long long line"
	cmp.CC = append(cmp.CC, &cmpFX{}, &cmpFX{
		onInit: func(c *cmpFX, e *Env) {
			c.Dim().SetHeight(1)
			fmt.Fprint(e, long)
		}},
	)
	cmp.onUpdate = func(_ *cmpFX, _ *Env, data interface{}) {
		cmp.CC[1] = data.(Componenter)
	}
	tt := xcmpfx(t, cmp)
	t.Eq(long, tt.ScreenOf(cmp).Trimmed().String())

	t.FatalOn(tt.Lines.Update(cmp, &cmpFX{
		onInit: func(c *cmpFX, e *Env) {
			c.Dim().SetHeight(1)
			fmt.Fprint(e, "short line")
		}}, nil))
	str := tt.ScreenOf(cmp).Trimmed().String()
	t.Eq("short line", str)
}

func (s *AComponent) Updates_tab_expansions_on_tab_width_change(t *T) {
	fx, cmp := fxCmp(t)
	fx.FireResize(11, 2)

	tb, exp, expTB := "", 8, strings.Repeat(" ", 8)
	fx.Lines.Update(cmp, nil, func(e *Env) {
		tb = strings.Repeat(" ", cmp.Globals().TabWidth())
		fmt.Fprint(e, "\t1st\n\t2nd")
	})
	t.Not.Eq(tb, expTB)
	t.True(strings.HasPrefix(fx.Screen()[0], tb+"1st"))
	t.True(strings.HasPrefix(fx.Screen()[1], tb+"2nd"))

	fx.Lines.Update(cmp, nil, func(e *Env) {
		cmp.Globals().SetTabWidth(exp)
	})

	t.True(strings.HasPrefix(fx.Screen()[0], expTB+"1st"))
	t.True(strings.HasPrefix(fx.Screen()[1], expTB+"2nd"))
}

func (s *AComponent) Is_lines_cursor_component_on_set_cursor(t *T) {
	stacking := &stackingFX{}
	stacking.CC = append(stacking.CC, &cmpFX{}, &cmpFX{})
	fx := fx(t, stacking)
	cmp := stacking.CC[1].(*cmpFX)
	fx.Lines.Update(cmp, nil, func(e *Env) {
		cmp.SetCursor(0, 0, BlockCursorBlinking)
	})
	t.Eq(cmp, fx.Lines.CursorComponent())
}

func (s *AComponent) Ignores_setting_cursor_outside_content(t *T) {
	stacking := &stackingFX{}
	stacking.CC = append(stacking.CC, &cmpFX{}, &cmpFX{})
	fx := fx(t, stacking)
	cmp := stacking.CC[1].(*cmpFX)
	fx.Lines.Update(cmp, nil, func(e *Env) {
		_, _, w, _ := cmp.ContentArea()
		cmp.SetCursor(0, w, BlockCursorBlinking)
	})
	t.Eq(Componenter(nil), fx.Lines.CursorComponent())
	fx.Lines.Update(cmp, nil, func(e *Env) {
		_, _, _, h := cmp.ContentArea()
		cmp.SetCursor(h, 0, BlockCursorBlinking)
	})
	t.Eq(Componenter(nil), fx.Lines.CursorComponent())
}

func (s *AComponent) Stacking_or_chaining_ignores_cursor_sets(t *T) {
	stacking := &stackingFX{}
	stacking.CC = append(stacking.CC, &cmpFX{}, &cmpFX{})
	fx := fx(t, stacking)
	fx.Lines.Update(stacking, nil, func(e *Env) {
		stacking.SetCursor(0, 0)
	})
	t.Eq(Componenter(nil), fx.Lines.CursorComponent())
	chaining := &chainingFX{}
	chaining.CC = append(chaining.CC, &cmpFX{})
	fx.Lines.Update(chaining, nil, func(e *Env) {
		chaining.SetCursor(0, 0)
	})
	t.Eq(Componenter(nil), fx.Lines.CursorComponent())
}

func TestComponent(t *testing.T) {
	t.Parallel()
	Run(&AComponent{}, t)
}
