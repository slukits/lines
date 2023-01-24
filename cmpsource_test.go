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

type ASourcedComponent struct{ Suite }

func (s *ASourcedComponent) Has_initially_a_dirty_source(t *T) {
	tt, cmp := fxCmp(t)
	tt.Lines.Update(cmp, nil, func(e *Env) {
		cmp.Src = &ContentSource{}
		t.True(cmp.Src.IsDirty())
	})
}

func (s *ASourcedComponent) Is_dirty_after_setting_its_source(t *T) {
	tt, cmp := fxCmp(t)
	tt.Lines.Update(cmp, nil, func(e *Env) {
		t.Not.True(cmp.IsDirty())
		cmp.Src = &ContentSource{}
		t.True(cmp.IsDirty())
	})
}

func (s *ASourcedComponent) Source_is_clean_after_first_sync(t *T) {
	tt, cmp := fxCmp(t)
	tt.Lines.Update(cmp, nil, func(e *Env) {
		cmp.Src = &ContentSource{}
	})
	tt.Lines.Update(cmp, nil, func(e *Env) {
		t.Not.True(cmp.Src.IsDirty())
	})
}

func (s *ASourcedComponent) Displays_the_first_n_source_lines(t *T) {
	tt, cmp := fxCmp(t)
	tt.FireResize(3, 2)
	tt.Lines.Update(cmp, nil, func(e *Env) {
		cmp.Src = &ContentSource{Liner: &linerFX{}}
	})

	t.Eq("1st\n2nd", tt.Screen())
}

func (s *ASourcedComponent) Has_no_more_lines_than_screen_lines(t *T) {
	cmp := &srcFX{liner: &focusableLinerFX{highlighted: true}}
	fx := fx(t, cmp, 20*time.Minute)
	fx.FireResize(3, 2)
	fx.Lines.Update(cmp, nil, func(e *Env) {
		t.Eq(len(*cmp.ll), cmp.ContentScreenLines())
		t.Eq(2, len(*cmp.ll))
	})
	fx.FireResize(3, 6)
	fx.Lines.Update(cmp, nil, func(e *Env) {
		t.Eq(len(*cmp.ll), cmp.ContentScreenLines())
		t.Eq(6, len(*cmp.ll))
	})
}

func (s *ASourcedComponent) Is_scrollable_if_source_liner_has_len(t *T) {
	tt, cmp := fxCmp(t)
	tt.FireResize(3, 2)
	tt.Lines.Update(cmp, nil, func(e *Env) {
		cmp.Src = &ContentSource{Liner: &scrollableLinerFX{}}
	})

	tt.Lines.Update(cmp, nil, func(e *Env) {
		t.True(cmp.FF.Has(Scrollable))
	})
}

func (s *ASourcedComponent) Has_dirty_source_on_scrolling(t *T) {
	tt, cmp := fxCmp(t)
	tt.FireResize(3, 1)
	t.FatalOn(tt.Lines.Update(cmp, nil, func(e *Env) {
		cmp.Src = &ContentSource{Liner: &scrollableLinerFX{}}
	}))

	t.FatalOn(tt.Lines.Update(cmp, nil, func(e *Env) {
		cmp.Scroll.Down()
		t.True(cmp.Src.IsDirty())
	}))
}

func (s *ASourcedComponent) Scrolls_by_one_line_if_height_is_one(t *T) {
	tt, cmp := fxCmp(t)
	tt.FireResize(3, 1)
	t.FatalOn(tt.Lines.Update(cmp, nil, func(e *Env) {
		cmp.Src = &ContentSource{Liner: &scrollableLinerFX{}}
	}))
	t.Eq("1st", tt.Screen())

	t.FatalOn(tt.Lines.Update(cmp, nil, func(e *Env) {
		cmp.Scroll.Down()
	}))
	t.Eq("2nd", tt.Screen())
	t.FatalOn(tt.Lines.Update(cmp, nil, func(e *Env) {
		cmp.Scroll.Down()
	}))
	t.Eq("3rd", tt.Screen())

	t.FatalOn(tt.Lines.Update(cmp, nil, func(e *Env) {
		cmp.Scroll.Up()
	}))
	t.Eq("2nd", tt.Screen())
	t.FatalOn(tt.Lines.Update(cmp, nil, func(e *Env) {
		cmp.Scroll.Up()
	}))
	t.Eq("1st", tt.Screen())
}

func (s *ASourcedComponent) Scrolls_to_end_if_last_line_displayed(t *T) {
	tt, cmp := fxCmp(t)
	tt.FireResize(3, 3)
	t.FatalOn(tt.Lines.Update(cmp, nil, func(e *Env) {
		cmp.Src = &ContentSource{Liner: &scrollableLinerFX{}}
		cmp.Src.Liner.(*scrollableLinerFX).initLines(4)
	}))
	t.Eq("1st\n2nd\n3rd", tt.Screen())

	t.FatalOn(tt.Lines.Update(cmp, nil, func(e *Env) {
		cmp.Scroll.Down()
	}))
	t.Eq("2nd\n3rd\n4th", tt.Screen())
}

func (s *ASourcedComponent) Scrolls_to_first_if_first_line_displayed(
	t *T,
) {
	tt, cmp := fxCmp(t)
	tt.FireResize(3, 3)
	t.FatalOn(tt.Lines.Update(cmp, nil, func(e *Env) {
		cmp.Src = &ContentSource{Liner: &scrollableLinerFX{}}
		cmp.Src.Liner.(*scrollableLinerFX).initLines(4)
	}))
	t.Eq("1st\n2nd\n3rd", tt.Screen())

	t.FatalOn(tt.Lines.Update(cmp, nil, func(e *Env) {
		cmp.Scroll.Down()
	}))
	t.Eq("2nd\n3rd\n4th", tt.Screen())

	t.FatalOn(tt.Lines.Update(cmp, nil, func(e *Env) {
		cmp.Scroll.Up()
	}))
	t.Eq("1st\n2nd\n3rd", tt.Screen())
}

func (s *ASourcedComponent) Scrolls_to_top_and_bottom(t *T) {
	tt, cmp := fxCmp(t)
	tt.FireResize(3, 2)
	t.FatalOn(tt.Lines.Update(cmp, nil, func(e *Env) {
		cmp.Src = &ContentSource{Liner: &scrollableLinerFX{}}
		cmp.Src.Liner.(*scrollableLinerFX).initLines(5)
	}))
	t.Eq("1st\n2nd", tt.Screen())

	t.FatalOn(tt.Lines.Update(cmp, nil, func(e *Env) {
		cmp.Scroll.ToBottom()
	}))
	t.Eq("4th\n5th", tt.Screen())

	t.FatalOn(tt.Lines.Update(cmp, nil, func(e *Env) {
		cmp.Scroll.ToTop()
	}))
	t.Eq("1st\n2nd", tt.Screen())
}

func (s *ASourcedComponent) Scrolls_down_by_90_percent_height(t *T) {
	fx, cmp := fxCmp(t)
	fx.FireResize(20, 30)
	t.FatalOn(fx.Lines.Update(cmp, nil, func(e *Env) {
		cmp.Dim().SetHeight(5)
		cmp.Src = &ContentSource{Liner: &scrollableLinerFX{}}
		cmp.Src.Liner.(*scrollableLinerFX).initLines(60)
		t.Eq(60, cmp.Src.Liner.(ScrollableLiner).Len())
	}))

	exp := []string{}
	for i := 0; i < 5; i++ {
		exp = append(exp, fmt.Sprintf("%dth", i+5))
	}
	t.FatalOn(fx.Lines.Update(cmp, nil, func(e *Env) {
		cmp.Scroll.Down()
	}))
	t.Eq(strings.Join(exp, "\n"), fx.Screen().Trimmed())

	t.FatalOn(fx.Lines.Update(cmp, nil, func(e *Env) {
		cmp.Dim().SetHeight(15)
	}))

	exp = []string{}
	for i := 0; i < 15; i++ { // first is still at fifth line
		exp = append(exp, fmt.Sprintf("%dth", i+19))
	}
	t.FatalOn(fx.Lines.Update(cmp, nil, func(e *Env) {
		cmp.Scroll.Down()
	}))
	t.Eq(strings.Join(exp, "\n"), fx.Screen().Trimmed())

	t.FatalOn(fx.Lines.Update(cmp, nil, func(e *Env) {
		cmp.Scroll.ToTop()
		t.True(cmp.Scroll.IsAtTop())
		cmp.Dim().SetHeight(30)
	}))

	exp = []string{}
	for i := 0; i < 30; i++ { // first is still at fifth line
		exp = append(exp, fmt.Sprintf("%dth", i+28))
	}
	t.FatalOn(fx.Lines.Update(cmp, nil, func(e *Env) {
		cmp.Scroll.Down()
	}))
	t.Eq(strings.Join(exp, "\n"), fx.Screen().Trimmed().String())
}

func (s *ASourcedComponent) Scrolls_up_by_90_percent_height(t *T) {
	fx, cmp := fxCmp(t)
	fx.FireResize(20, 30)
	t.FatalOn(fx.Lines.Update(cmp, nil, func(e *Env) {
		cmp.Dim().SetHeight(5)
		cmp.Src = &ContentSource{Liner: &scrollableLinerFX{}}
		cmp.Src.Liner.(*scrollableLinerFX).initLines(60)
		t.Eq(60, cmp.Src.Liner.(ScrollableLiner).Len())
		cmp.Scroll.ToBottom()
	}))

	exp := []string{}
	for i := 0; i < 5; i++ {
		exp = append(exp, fmt.Sprintf("%dth", i+56))
	}
	t.Eq(strings.Join(exp, "\n"), fx.Screen().Trimmed())
	exp = []string{}
	for i := 0; i < 5; i++ {
		exp = append(exp, fmt.Sprintf("%dth", i+52))
	}
	t.FatalOn(fx.Lines.Update(cmp, nil, func(e *Env) {
		cmp.Scroll.Up()
	}))
	t.Eq(strings.Join(exp, "\n"), fx.Screen().Trimmed())

	t.FatalOn(fx.Lines.Update(cmp, nil, func(e *Env) {
		cmp.Dim().SetHeight(15) // first changes to 46
	}))

	exp = []string{}
	for i := 0; i < 15; i++ {
		exp = append(exp, fmt.Sprintf("%dth", i+32))
	}
	t.FatalOn(fx.Lines.Update(cmp, nil, func(e *Env) {
		cmp.Scroll.Up()
	}))
	t.Eq(strings.Join(exp, "\n"), fx.Screen().Trimmed())

	t.FatalOn(fx.Lines.Update(cmp, nil, func(e *Env) {
		cmp.Dim().SetHeight(30)
		cmp.Scroll.ToBottom()
		t.True(cmp.Scroll.IsAtBottom())
	}))
	exp = []string{}
	for i := 0; i < 30; i++ {
		if i+4 < 10 {
			exp = append(exp, fmt.Sprintf("%dth ", i+4))
			continue
		}
		exp = append(exp, fmt.Sprintf("%dth", i+4))
	}
	t.FatalOn(fx.Lines.Update(cmp, nil, func(e *Env) {
		cmp.Scroll.Up()
	}))
	t.Eq(strings.Join(exp, "\n"), fx.Screen().Trimmed())
}

func (s *ASourcedComponent) Scrolls_to_top_on_reset_all(t *T) {
	tt, cmp := fxCmp(t)
	tt.FireResize(20, 2)
	t.FatalOn(tt.Lines.Update(cmp, nil, func(e *Env) {
		cmp.Src = &ContentSource{Liner: &scrollableLinerFX{}}
		cmp.Scroll.ToBottom()
		t.Not.True(cmp.Scroll.IsAtTop())
	}))
	t.FatalOn(tt.Lines.Update(cmp, nil, func(e *Env) {
		cmp.Reset(All)
		t.True(cmp.Scroll.IsAtTop())
	}))
}

func (s *ASourcedComponent) Is_focusable_on_focusable_src_liner(t *T) {
	tt, cmp := fxCmp(t)
	t.FatalOn(tt.Lines.Update(cmp, nil, func(e *Env) {
		cmp.Src = &ContentSource{Liner: &focusableLinerFX{}}
	}))
	t.FatalOn(tt.Lines.Update(cmp, nil, func(e *Env) {
		cmp.FF.Has(LinesFocusable)
	}))
}

func (s *ASourcedComponent) Focuses_first_focusable_line(t *T) {
	cmp := &srcFX{}
	cmp.onLineFocus = func(c *cmpFX, _ *Env, _, _ int) {
		t.Eq(3, c.LL.Focus.Current())
	}
	fx := fx(t, cmp)
	fx.FireResize(3, 2)
	t.FatalOn(fx.Lines.Update(cmp, nil, func(e *Env) {
		liner := (&focusableLinerFX{}).initLines(4)
		liner.focusable = func(idx int) bool { return idx == 3 }
		cmp.Src = &ContentSource{Liner: liner}
	}))

	fx.FireKey(Down)
	t.Eq(1, cmp.N(onLineFocus))
}

func (s *ASourcedComponent) Focuses_next_focusable_line(t *T) {
	cmp := &srcFX{}
	cmp.onLineFocus = func(c *cmpFX, _ *Env, _, _ int) {
		switch c.N(onFocusLost) {
		case 1:
			t.Eq(3, c.LL.Focus.Current())
		case 2:
			t.Eq(7, c.LL.Focus.Current())
		}
	}
	fx := fx(t, cmp)
	fx.FireResize(3, 2)
	t.FatalOn(fx.Lines.Update(cmp, nil, func(e *Env) {
		liner := &focusableLinerFX{}
		liner.focusable = func(idx int) bool {
			return idx == 3 || idx == 7
		}
		cmp.Src = &ContentSource{Liner: liner}
	}))

	fx.FireKey(Down)
	fx.FireKey(Down)
	t.Eq(2, cmp.N(onLineFocus))
}

func (s *ASourcedComponent) Resets_focus_if_next_not_focusable(t *T) {
	cmp := &srcFX{}
	cmp.onLineFocus = func(c *cmpFX, _ *Env, _, _ int) {
		t.Eq(0, c.LL.Focus.Current())
	}
	fx := fx(t, cmp)
	fx.FireResize(3, 2)
	t.FatalOn(fx.Lines.Update(cmp, nil, func(e *Env) {
		liner := (&focusableLinerFX{}).initLines(1)
		cmp.Src = &ContentSource{Liner: liner}
	}))

	fx.FireKey(Down)
	fx.FireKey(Down)
	t.FatalOn(fx.Lines.Update(cmp, nil, func(e *Env) {
		t.Eq(-1, cmp.LL.Focus.Current())
		cmp.Src.Liner.(*focusableLinerFX).cc = []string{
			"1st", "2nd", "3rd"}
		cmp.Src.Liner.(*focusableLinerFX).focusable = func(idx int) bool {
			return idx == 0
		}
	}))

	fx.FireKey(Down)
	fx.FireKey(Down)
	t.FatalOn(fx.Lines.Update(cmp, nil, func(e *Env) {
		t.Eq(-1, cmp.LL.Focus.Current())
	}))
	t.Eq(2, cmp.N(onLineFocus))
}

func (s *ASourcedComponent) Focuses_previous_focusable_line(t *T) {
	cmp := &srcFX{}
	cmp.onLineFocus = func(c *cmpFX, _ *Env, _, _ int) {
		switch c.N(onLineFocus) {
		case 1:
			t.Eq(3, c.LL.Focus.Current())
		case 2:
			t.Eq(7, c.LL.Focus.Current())
		case 3:
			t.Eq(3, c.LL.Focus.Current())
		}
	}
	fx := fx(t, cmp)
	fx.FireResize(3, 2)
	t.FatalOn(fx.Lines.Update(cmp, nil, func(e *Env) {
		liner := &focusableLinerFX{}
		liner.focusable = func(idx int) bool {
			return idx == 3 || idx == 7
		}
		cmp.Src = &ContentSource{Liner: liner}
	}))

	fx.FireKey(Down)
	fx.FireKey(Down)
	fx.FireKey(Up)
	fx.FireKey(Up)
	t.FatalOn(fx.Lines.Update(cmp, nil, func(e *Env) {
		t.Eq(-1, cmp.LL.Focus.Current())
	}))
	t.Eq(3, cmp.N(onLineFocus))
}

func (s *ASourcedComponent) Triggers_reset_on_unfocusable_feature(t *T) {
	cmp := &srcFX{}
	cmp.onLineFocus = func(c *cmpFX, _ *Env, _, _ int) {
		switch c.N(onLineFocus) {
		case 1:
			t.Eq(0, c.LL.Focus.Current())
		case 2:
			t.Eq(-1, c.LL.Focus.Current())
		}
	}
	fx := fx(t, cmp)
	fx.FireResize(3, 2)
	t.FatalOn(fx.Lines.Update(cmp, nil, func(e *Env) {
		liner := (&focusableLinerFX{}).initLines(2)
		cmp.Src = &ContentSource{Liner: liner}
	}))

	fx.FireKey(Down)
	fx.FireKey(Esc)

	t.Eq(1, cmp.N(onLineFocus))
}

func (s *ASourcedComponent) Scrolls_to_next_focusable_line(t *T) {
	cmp := &srcFX{liner: &focusableLinerFX{highlighted: true}}
	cmp.onInit = func(c *cmpFX, e *Env) {
		c.Dim().SetWidth(3).SetHeight(2)
	}
	fx := fx(t, cmp)
	fx.FireResize(3, 2)
	fx.FireKeys(Down, PgDn, Down)
	fx.Lines.Update(cmp, nil, func(e *Env) {
		t.Eq(1, cmp.LL.Focus.Current())
	})
	t.Eq("2nd\n3rd", fx.Screen())
	fx.FireKeys(Down, Down, Down, Down, Down, PgUp, PgUp, Up)
	t.Eq("5th\n6th", fx.Screen())
	fx.Lines.Update(cmp, nil, func(e *Env) {
		t.Eq(5, cmp.LL.Focus.Current())
	})
}

func (s *ASourcedComponent) Scrolls_inside_its_gaps(t *T) {
	cmp := &srcFX{cmpFX: cmpFX{gaps: true, onInit: func(c *cmpFX, e *Env) {
		c.Src = &ContentSource{Liner: &focusableLinerFX{}}
	}}}
	fx := fx(t, cmp)
	fx.FireResize(5, 3)
	exp := "•••••\n•1st•\n•••••"
	t.Eq(exp, fx.Screen())

	scrollTest := func(scroll func(), exp string) {
		fx.Lines.Update(cmp, nil, func(e *Env) {
			scroll()
		})
		t.Eq(exp, fx.Screen())
	}
	scrollTest(cmp.Scroll.Down, "•••••\n•2nd•\n•••••")
	scrollTest(cmp.Scroll.Down, "•••••\n•3rd•\n•••••")
	scrollTest(cmp.Scroll.Up, "•••••\n•2nd•\n•••••")
	scrollTest(cmp.Scroll.ToTop, "•••••\n•1st•\n•••••")
	scrollTest(cmp.Scroll.ToBottom, "•••••\n•8th•\n•••••")

	fx.FireResize(5, 5)
	exp = "•••••\n•6th•\n•7th•\n•8th•\n•••••"
	t.Eq(exp, fx.Screen())

	scrollTest(cmp.Scroll.Up, "•••••\n•4th•\n•5th•\n•6th•\n•••••")
	scrollTest(cmp.Scroll.Up, "•••••\n•2nd•\n•3rd•\n•4th•\n•••••")
	scrollTest(cmp.Scroll.ToBottom, "•••••\n•6th•\n•7th•\n•8th•\n•••••")
	scrollTest(cmp.Scroll.ToTop, "•••••\n•1st•\n•2nd•\n•3rd•\n•••••")
}

func (s *ASourcedComponent) Scrolls_to_next_focusable_within_gaps(
	t *T,
) {
	cmp := &srcFX{cmpFX: cmpFX{gaps: true, onInit: func(c *cmpFX, e *Env) {
		c.Src = &ContentSource{Liner: &focusableLinerFX{highlighted: true}}
		c.Src.Liner.(*focusableLinerFX).focusable = func(idx int) bool {
			return idx == 3 || idx == 7
		}
	}}}
	fx := fx(t, cmp)
	fx.FireResize(5, 5)
	t.Eq("•••••\n•1st•\n•2nd•\n•3rd•\n•••••", fx.Screen())

	fx.FireKey(Down)
	t.Eq("•••••\n•3rd•\n•4th•\n•5th•\n•••••", fx.Screen())
	fx.Lines.Update(cmp, nil, func(e *Env) {
		t.True(cmp.LL.By(1).IsFlagged(Highlighted))
	})

	fx.FireKey(Down)
	t.Eq("•••••\n•6th•\n•7th•\n•8th•\n•••••", fx.Screen())
	fx.Lines.Update(cmp, nil, func(e *Env) {
		t.Not.True(cmp.LL.By(1).IsFlagged(Highlighted))
		t.True(cmp.LL.By(2).IsFlagged(Highlighted))
	})
}

func (s *ASourcedComponent) Scrolls_to_previous_focusable_within_gaps(
	t *T,
) {
	cmp := &srcFX{cmpFX: cmpFX{gaps: true, onInit: func(c *cmpFX, e *Env) {
		c.Src = &ContentSource{Liner: &focusableLinerFX{highlighted: true}}
		c.Src.Liner.(*focusableLinerFX).focusable = func(idx int) bool {
			return idx == 3 || idx == 7
		}
	}}}
	fx := fx(t, cmp)
	fx.FireResize(5, 5)
	t.Eq("•••••\n•1st•\n•2nd•\n•3rd•\n•••••", fx.Screen())

	fx.FireKey(Up)
	t.Eq("•••••\n•6th•\n•7th•\n•8th•\n•••••", fx.Screen())
	fx.Lines.Update(cmp, nil, func(e *Env) {
		t.True(cmp.LL.By(2).IsFlagged(Highlighted))
	})

	fx.FireKey(Up)
	t.Eq("•••••\n•4th•\n•5th•\n•6th•\n•••••", fx.Screen())
	fx.Lines.Update(cmp, nil, func(e *Env) {
		t.Not.True(cmp.LL.By(2).IsFlagged(Highlighted))
		t.True(cmp.LL.By(0).IsFlagged(Highlighted))
	})
}

func (s *ASourcedComponent) Remembers_highlighted_line_on_scrolling(
	t *T,
) {
	cmp := &srcFX{cmpFX: cmpFX{gaps: true, onInit: func(c *cmpFX, e *Env) {
		c.Src = &ContentSource{Liner: &focusableLinerFX{highlighted: true}}
		c.Src.Liner.(*focusableLinerFX).focusable = func(idx int) bool {
			return idx == 3 || idx == 7
		}
	}}}
	fx := fx(t, cmp)
	fx.FireResize(5, 5)
	t.Eq("•••••\n•1st•\n•2nd•\n•3rd•\n•••••", fx.Screen())

	fx.FireKey(Down)
	t.Eq("•••••\n•3rd•\n•4th•\n•5th•\n•••••", fx.Screen())
	fx.Lines.Update(cmp, nil, func(e *Env) {
		t.True(cmp.LL.By(1).IsFlagged(Highlighted))
	})

	fx.Lines.Update(cmp, nil, func(e *Env) {
		cmp.Scroll.Up()
	})
	t.Eq("•••••\n•1st•\n•2nd•\n•3rd•\n•••••", fx.Screen())
	fx.Lines.Update(cmp, nil, func(e *Env) {
		for i := 0; i < 3; i++ {
			t.Not.True(cmp.LL.By(i).IsFlagged(Highlighted))
		}
	})

	fx.Lines.Update(cmp, nil, func(e *Env) {
		cmp.Scroll.Down()
	})
	t.Eq("•••••\n•3rd•\n•4th•\n•5th•\n•••••", fx.Screen())
	fx.Lines.Update(cmp, nil, func(e *Env) {
		t.True(cmp.LL.By(1).IsFlagged(Highlighted))
	})

	fx.Lines.Update(cmp, nil, func(e *Env) {
		cmp.Scroll.Down()
	})
	t.Eq("•••••\n•5th•\n•6th•\n•7th•\n•••••", fx.Screen())
	fx.Lines.Update(cmp, nil, func(e *Env) {
		for i := 0; i < 3; i++ {
			t.Not.True(cmp.LL.By(i).IsFlagged(Highlighted))
		}
	})

	fx.Lines.Update(cmp, nil, func(e *Env) {
		cmp.Scroll.Up()
	})
	t.Eq("•••••\n•3rd•\n•4th•\n•5th•\n•••••", fx.Screen())
	fx.Lines.Update(cmp, nil, func(e *Env) {
		t.True(cmp.LL.By(1).IsFlagged(Highlighted))
	})
}

func (s *ASourcedComponent) Inverts_bg_fg_of_focused_if_highlighted(
	t *T,
) {
	cmp := &srcFX{cmpFX: cmpFX{onInit: func(c *cmpFX, e *Env) {
		liner := (&focusableLinerFX{highlighted: true}).initLines(2)
		liner.focusable = func(idx int) bool { return idx == 1 }
		c.Src = &ContentSource{Liner: liner}
		c.dim.SetWidth(3).SetHeight(2)
	},
		onLineFocus: func(c *cmpFX, _ *Env, _, _ int) {
			switch c.N(onLineFocus) {
			case 1:
				t.Eq(1, c.LL.Focus.Current())
			}
		},
	}}
	fx := fx(t, cmp)
	fx.FireResize(5, 2)

	l2 := fx.Cells()[1]
	for x := range l2 {
		t.Not.True(l2.HasAA(x, Reverse))
	}

	fx.FireKey(Down)
	t.Eq(1, cmp.N(onLineFocus))
	l2 = fx.Cells()[1]
	for x := range l2 {
		if x != 0 && x != 4 {
			t.True(l2.HasAA(x, Reverse))
			continue
		}
		t.Not.True(l2.HasAA(x, Reverse))
	}
}

func (s *ASourcedComponent) Gets_its_selected_lines_reported(t *T) {
	cmp := &srcFX{cmpFX: cmpFX{gaps: true,
		onInit: func(c *cmpFX, e *Env) {
			c.FF.Set(LineSelectable)
			liner := &focusableLinerFX{highlighted: true}
			liner.focusable = func(idx int) bool {
				return idx == 3 || idx == 7
			}
			c.Src = &ContentSource{Liner: liner}
		},
		onLineFocus: func(c *cmpFX, _ *Env, _, _ int) {
			switch c.N(onLineFocus) {
			case 1:
				t.Eq(3, c.LL.Focus.Current())
			case 2:
				t.Eq(7, c.LL.Focus.Current())
			}
		},
		onLineSelection: func(_ *cmpFX, _ *Env, cIdx, _ int) {
			t.Eq(7, cIdx)
		},
	}}
	fx := fx(t, cmp)
	fx.FireResize(3, 2)

	fx.FireKey(Down)
	fx.FireKey(Down)
	fx.FireKey(Enter)

	t.Eq(2, cmp.N(onLineFocus))
	t.Eq(1, cmp.N(onLineSelection))
}

func TestASourcedComponent(t *testing.T) {
	t.Parallel()
	Run(&ASourcedComponent{}, t)
}
