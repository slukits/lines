// Copyright (c) 2022 Stephan Lukits. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package lines

import (
	"fmt"
	"strings"
	"testing"

	. "github.com/slukits/gounit"
)

type linerFX struct {
	cc []string
}

func (l *linerFX) Print(idx int, w *EnvLineWriter) bool {
	if len(l.cc) <= idx || idx < 0 {
		return false
	}
	fmt.Fprintf(w, l.cc[idx])
	return idx+1 < len(l.cc)
}

type ASourcedComponent struct{ Suite }

func (s *ASourcedComponent) Has_initially_a_dirty_source(t *T) {
	tt, cmp := cmpfx(t)
	tt.Lines.Update(cmp, nil, func(e *Env) {
		cmp.Src = &ContentSource{}
		t.True(cmp.Src.IsDirty())
	})
}

func (s *ASourcedComponent) Is_dirty_after_setting_its_source(t *T) {
	tt, cmp := cmpfx(t)
	tt.Lines.Update(cmp, nil, func(e *Env) {
		cmp.Src = &ContentSource{}
		t.True(cmp.IsDirty())
	})
}

func (s *ASourcedComponent) Source_is_clean_after_first_sync(t *T) {
	tt, cmp := cmpfx(t)
	tt.Lines.Update(cmp, nil, func(e *Env) {
		cmp.Src = &ContentSource{}
	})
	tt.Lines.Update(cmp, nil, func(e *Env) {
		t.Not.True(cmp.Src.IsDirty())
	})
}

func (s *ASourcedComponent) Displays_the_first_n_source_lines(t *T) {
	tt, cmp := cmpfx(t)
	tt.FireResize(3, 2)
	tt.Lines.Update(cmp, nil, func(e *Env) {
		cmp.Src = &ContentSource{
			Liner: &linerFX{cc: []string{"1st", "2nd", "3rd"}}}
	})

	t.Eq("1st\n2nd", tt.Screen())
}

type scrollLinerFX struct{ linerFX }

func (sl scrollLinerFX) Len() int { return len(sl.cc) }

func (s *ASourcedComponent) Is_scrollable_if_source_liner_has_len(t *T) {
	tt, cmp := cmpfx(t)
	tt.FireResize(3, 2)
	tt.Lines.Update(cmp, nil, func(e *Env) {
		cmp.Src = &ContentSource{Liner: &scrollLinerFX{
			linerFX: linerFX{cc: []string{"1st", "2nd", "3rd"}}}}
	})

	tt.Lines.Update(cmp, nil, func(e *Env) {
		t.True(cmp.FF.Has(Scrollable))
	})
}

func (s *ASourcedComponent) Has_dirty_source_on_scrolling(t *T) {
	tt, cmp := cmpfx(t)
	tt.FireResize(3, 1)
	t.FatalOn(tt.Lines.Update(cmp, nil, func(e *Env) {
		cmp.Src = &ContentSource{Liner: &scrollLinerFX{
			linerFX: linerFX{cc: []string{"1st", "2nd"}}}}
	}))

	t.FatalOn(tt.Lines.Update(cmp, nil, func(e *Env) {
		cmp.Scroll.Down()
		t.True(cmp.Src.IsDirty())
	}))
}

func (s *ASourcedComponent) Scrolls_by_one_line_if_height_is_one(t *T) {
	tt, cmp := cmpfx(t)
	tt.FireResize(3, 1)
	t.FatalOn(tt.Lines.Update(cmp, nil, func(e *Env) {
		cmp.Src = &ContentSource{Liner: &scrollLinerFX{
			linerFX: linerFX{cc: []string{"1st", "2nd"}}}}
	}))
	t.Eq("1st", tt.Screen())

	t.FatalOn(tt.Lines.Update(cmp, nil, func(e *Env) {
		cmp.Scroll.Down()
	}))
	t.Eq("2nd", tt.Screen().Trimmed())
	t.FatalOn(tt.Lines.Update(cmp, nil, func(e *Env) {
		cmp.Scroll.Down()
	}))
	t.Eq("2nd", tt.Screen())

	t.FatalOn(tt.Lines.Update(cmp, nil, func(e *Env) {
		cmp.Scroll.Up()
	}))
	t.Eq("1st", tt.Screen().Trimmed())
	t.FatalOn(tt.Lines.Update(cmp, nil, func(e *Env) {
		cmp.Scroll.Up()
	}))
	t.Eq("1st", tt.Screen())
}

func (s *ASourcedComponent) Scrolls_to_end_if_last_line_displayed(t *T) {
	tt, cmp := cmpfx(t)
	tt.FireResize(3, 3)
	t.FatalOn(tt.Lines.Update(cmp, nil, func(e *Env) {
		cmp.Src = &ContentSource{Liner: &scrollLinerFX{
			linerFX: linerFX{cc: []string{"1st", "2nd", "3rd", "4th"}}}}
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
	tt, cmp := cmpfx(t)
	tt.FireResize(3, 3)
	t.FatalOn(tt.Lines.Update(cmp, nil, func(e *Env) {
		cmp.Src = &ContentSource{Liner: &scrollLinerFX{
			linerFX: linerFX{cc: []string{"1st", "2nd", "3rd", "4th"}}}}
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
	tt, cmp := cmpfx(t)
	tt.FireResize(3, 2)
	t.FatalOn(tt.Lines.Update(cmp, nil, func(e *Env) {
		cmp.Src = &ContentSource{Liner: &scrollLinerFX{
			linerFX: linerFX{cc: []string{"1st", "2nd", "3rd", "4th", "5th"}}}}
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
	tt, cmp := cmpfx(t)
	tt.FireResize(20, 30)
	t.FatalOn(tt.Lines.Update(cmp, nil, func(e *Env) {
		cmp.Dim().SetHeight(5)
		ll := make([]string, 60)
		for i := 0; i < 60; i++ {
			ll[i] = fmt.Sprintf("line %d", i+1)
		}
		cmp.Src = &ContentSource{Liner: &scrollLinerFX{linerFX: linerFX{cc: ll}}}
		t.Eq(60, cmp.Src.Liner.(ScrollableLiner).Len())
	}))

	exp := []string{}
	for i := 0; i < 5; i++ {
		exp = append(exp, fmt.Sprintf("line %d", i+5))
	}
	t.FatalOn(tt.Lines.Update(cmp, nil, func(e *Env) {
		cmp.Scroll.Down()
	}))
	t.Eq(strings.Join(exp, "\n"), tt.Screen().Trimmed())

	t.FatalOn(tt.Lines.Update(cmp, nil, func(e *Env) {
		cmp.Dim().SetHeight(15)
	}))

	exp = []string{}
	for i := 0; i < 15; i++ { // first is still at fifth line
		exp = append(exp, fmt.Sprintf("line %d", i+19))
	}
	t.FatalOn(tt.Lines.Update(cmp, nil, func(e *Env) {
		cmp.Scroll.Down()
	}))
	t.Eq(strings.Join(exp, "\n"), tt.Screen().Trimmed())

	t.FatalOn(tt.Lines.Update(cmp, nil, func(e *Env) {
		cmp.Scroll.ToTop()
		t.True(cmp.Scroll.IsAtTop())
		cmp.Dim().SetHeight(30)
	}))

	exp = []string{}
	for i := 0; i < 30; i++ { // first is still at fifth line
		exp = append(exp, fmt.Sprintf("line %d", i+28))
	}
	t.FatalOn(tt.Lines.Update(cmp, nil, func(e *Env) {
		cmp.Scroll.Down()
	}))
	t.Eq(strings.Join(exp, "\n"), tt.Screen().Trimmed().String())
}

func (s *ASourcedComponent) Scrolls_up_by_90_percent_height(t *T) {
	tt, cmp := cmpfx(t)
	tt.FireResize(20, 30)
	t.FatalOn(tt.Lines.Update(cmp, nil, func(e *Env) {
		cmp.Dim().SetHeight(5)
		ll := make([]string, 60)
		for i := 0; i < 60; i++ {
			if i+1 < 10 {
				ll[i] = fmt.Sprintf("line 0%d", i+1)
				continue
			}
			ll[i] = fmt.Sprintf("line %d", i+1)
		}
		cmp.Src = &ContentSource{Liner: &scrollLinerFX{linerFX: linerFX{cc: ll}}}
		t.Eq(60, cmp.Src.Liner.(ScrollableLiner).Len())
		cmp.Scroll.ToBottom()
	}))

	exp := []string{}
	for i := 0; i < 5; i++ {
		exp = append(exp, fmt.Sprintf("line %d", i+56))
	}
	t.Eq(strings.Join(exp, "\n"), tt.Screen().Trimmed().String())
	exp = []string{}
	for i := 0; i < 5; i++ {
		exp = append(exp, fmt.Sprintf("line %d", i+52))
	}
	t.FatalOn(tt.Lines.Update(cmp, nil, func(e *Env) {
		cmp.Scroll.Up()
	}))
	t.Eq(strings.Join(exp, "\n"), tt.Screen().Trimmed().String())

	t.FatalOn(tt.Lines.Update(cmp, nil, func(e *Env) {
		cmp.Dim().SetHeight(15) // first changes to 46
	}))

	exp = []string{}
	for i := 0; i < 15; i++ {
		exp = append(exp, fmt.Sprintf("line %d", i+32))
	}
	t.FatalOn(tt.Lines.Update(cmp, nil, func(e *Env) {
		cmp.Scroll.Up()
	}))
	t.Eq(strings.Join(exp, "\n"), tt.Screen().Trimmed().String())

	t.FatalOn(tt.Lines.Update(cmp, nil, func(e *Env) {
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
	t.FatalOn(tt.Lines.Update(cmp, nil, func(e *Env) {
		cmp.Scroll.Up()
	}))
	t.Eq(strings.Join(exp, "\n"), tt.Screen().Trimmed().String())
}

func (s *ASourcedComponent) Scrolls_to_top_on_reset_all(t *T) {
	tt, cmp := cmpfx(t)
	tt.FireResize(20, 2)
	t.FatalOn(tt.Lines.Update(cmp, nil, func(e *Env) {
		cmp.Src = &ContentSource{Liner: &scrollLinerFX{
			linerFX: linerFX{cc: []string{"1st", "2nd", "3rd", "4th"}}}}
		cmp.Scroll.ToBottom()
		t.Not.True(cmp.Scroll.IsAtTop())
	}))
	t.FatalOn(tt.Lines.Update(cmp, nil, func(e *Env) {
		cmp.Reset(All)
		t.True(cmp.Scroll.IsAtTop())
	}))
}

type srcFcsFX struct {
	Component
	onLf func(*srcFcsFX, *Env, int)
	// the number of received line focus events
	lfN int
}

func (c *srcFcsFX) OnInit(e *Env) { c.FF.Add(LinesFocusable) }

func (c *srcFcsFX) OnLineFocus(e *Env, idx int) {
	c.lfN++
	if c.onLf == nil {
		return
	}
	c.onLf(c, e, idx)
}

// focusLinerFX implements FocusableLiner.
type focusLinerFX struct {
	scrollLinerFX
	notFocusable map[int]bool
}

func (l *focusLinerFX) IsFocusable(idx int) bool {
	return !l.notFocusable[idx]
}
func (l *focusLinerFX) Highlighted() (bool, bool) { return true, false }

func (s *ASourcedComponent) Is_focusable_on_focusable_src_liner(t *T) {
	tt, cmp := cmpfx(t)
	t.FatalOn(tt.Lines.Update(cmp, nil, func(e *Env) {
		liner := &focusLinerFX{}
		liner.cc = []string{"1st", "2nd", "3rd", "4th"}
		liner.notFocusable = map[int]bool{1: true, 2: true, 3: true}
		cmp.Src = &ContentSource{Liner: liner}
	}))

	t.FatalOn(tt.Lines.Update(cmp, nil, func(e *Env) {
		cmp.FF.Has(LinesFocusable)
	}))
}

func (s *ASourcedComponent) Focuses_first_focusable_line(t *T) {
	cmp := &srcFcsFX{
		onLf: func(c *srcFcsFX, _ *Env, _ int) {
			t.Eq(3, c.LL.Focus.Current())
		}}
	tt := xcmpfx(t, cmp)
	tt.FireResize(3, 2)
	t.FatalOn(tt.Lines.Update(cmp, nil, func(e *Env) {
		liner := &focusLinerFX{}
		liner.cc = []string{"1st", "2nd", "3rd", "4th"}
		liner.notFocusable = map[int]bool{0: true, 1: true, 2: true}
		cmp.Src = &ContentSource{Liner: liner}
	}))

	tt.FireRune('j')

	t.Eq(1, cmp.lfN)
}

func (s *ASourcedComponent) Focuses_next_focusable_line(t *T) {
	cmp := &srcFcsFX{
		onLf: func(c *srcFcsFX, _ *Env, _ int) {
			switch c.lfN {
			case 1:
				t.Eq(3, c.LL.Focus.Current())
			case 2:
				t.Eq(7, c.LL.Focus.Current())
			}
		}}
	tt := xcmpfx(t, cmp)
	tt.FireResize(3, 2)
	t.FatalOn(tt.Lines.Update(cmp, nil, func(e *Env) {
		liner := &focusLinerFX{}
		liner.cc = []string{
			"1st", "2nd", "3rd", "4th", "5th", "6th", "7th", "8th"}
		liner.notFocusable = map[int]bool{
			0: true, 1: true, 2: true, 4: true, 5: true, 6: true}
		cmp.Src = &ContentSource{Liner: liner}
	}))

	tt.FireRune('j')
	tt.FireKey(Down)

	t.Eq(2, cmp.lfN)
}

func (s *ASourcedComponent) Resets_focus_if_next_not_focusable(t *T) {
	cmp := &srcFcsFX{
		onLf: func(c *srcFcsFX, _ *Env, _ int) {
			switch c.lfN % 2 {
			case 0:
				t.Eq(-1, c.LL.Focus.Current())
			case 1:
				t.Eq(0, c.LL.Focus.Current())
			}
		}}
	tt := xcmpfx(t, cmp)
	tt.FireResize(3, 2)
	t.FatalOn(tt.Lines.Update(cmp, nil, func(e *Env) {
		liner := &focusLinerFX{}
		liner.cc = []string{"1st"}
		cmp.Src = &ContentSource{Liner: liner}
	}))

	tt.FireRune('j')
	tt.FireRune('j')

	t.FatalOn(tt.Lines.Update(cmp, nil, func(e *Env) {
		t.Eq(-1, cmp.LL.Focus.Current())
		cmp.Src.Liner.(*focusLinerFX).cc = []string{"1st", "2nd", "3rd"}
		cmp.Src.Liner.(*focusLinerFX).notFocusable = map[int]bool{
			1: true, 2: true}
	}))

	tt.FireKey(Down)
	tt.FireKey(Down)

	t.Eq(4, cmp.lfN)
}

func (s *ASourcedComponent) Focuses_previous_focusable_line(t *T) {
	cmp := &srcFcsFX{
		onLf: func(c *srcFcsFX, _ *Env, _ int) {
			switch c.lfN {
			case 1:
				t.Eq(3, c.LL.Focus.Current())
			case 2:
				t.Eq(7, c.LL.Focus.Current())
			case 3:
				t.Eq(3, c.LL.Focus.Current())
			case 4:
				t.Eq(-1, c.LL.Focus.Current())
			}
		}}
	tt := xcmpfx(t, cmp)
	tt.FireResize(3, 2)
	t.FatalOn(tt.Lines.Update(cmp, nil, func(e *Env) {
		liner := &focusLinerFX{}
		liner.cc = []string{
			"1st", "2nd", "3rd", "4th", "5th", "6th", "7th", "8th"}
		liner.notFocusable = map[int]bool{
			0: true, 1: true, 2: true, 4: true, 5: true, 6: true}
		cmp.Src = &ContentSource{Liner: liner}
	}))

	tt.FireRune('j')
	tt.FireKey(Down)

	tt.FireRune('k')
	tt.FireKey(Up)

	t.Eq(4, cmp.lfN)
}

func (s *ASourcedComponent) Triggers_reset_on_unfocusable_feature(t *T) {
	cmp := &srcFcsFX{
		onLf: func(c *srcFcsFX, _ *Env, _ int) {
			switch c.lfN {
			case 1:
				t.Eq(0, c.LL.Focus.Current())
			case 2:
				t.Eq(-1, c.LL.Focus.Current())
			}
		}}
	tt := xcmpfx(t, cmp)
	tt.FireResize(3, 2)
	t.FatalOn(tt.Lines.Update(cmp, nil, func(e *Env) {
		liner := &focusLinerFX{}
		liner.cc = []string{"1st", "2nd"}
		cmp.Src = &ContentSource{Liner: liner}
	}))

	tt.FireRune('j')
	tt.FireKey(Esc)

	t.Eq(2, cmp.lfN)
}

func (s *ASourcedComponent) Scrolls_inside_its_gaps(t *T) {
	cmp := &srcLsFX{gaps: true}
	tt := xcmpfx(t, cmp)
	tt.FireResize(5, 3)
	exp := "•••••\n•1st•\n•••••"
	t.Eq(exp, tt.Screen())

	scrollTest := func(scroll func(), exp string) {
		tt.Lines.Update(cmp, nil, func(e *Env) {
			scroll()
		})
		t.Eq(exp, tt.Screen())
	}
	scrollTest(cmp.Scroll.Down, "•••••\n•2nd•\n•••••")
	scrollTest(cmp.Scroll.Down, "•••••\n•3rd•\n•••••")
	scrollTest(cmp.Scroll.Up, "•••••\n•2nd•\n•••••")
	scrollTest(cmp.Scroll.ToTop, "•••••\n•1st•\n•••••")
	scrollTest(cmp.Scroll.ToBottom, "•••••\n•8th•\n•••••")

	tt.FireResize(5, 5)
	exp = "•••••\n•6th•\n•7th•\n•8th•\n•••••"
	t.Eq(exp, tt.Screen())

	scrollTest(cmp.Scroll.Up, "•••••\n•4th•\n•5th•\n•6th•\n•••••")
	scrollTest(cmp.Scroll.Up, "•••••\n•2nd•\n•3rd•\n•4th•\n•••••")
	scrollTest(cmp.Scroll.ToBottom, "•••••\n•6th•\n•7th•\n•8th•\n•••••")
	scrollTest(cmp.Scroll.ToTop, "•••••\n•1st•\n•2nd•\n•3rd•\n•••••")
}

func (s *ASourcedComponent) Scrolls_to_next_focusable_within_gaps(
	t *T,
) {
	cmp := &srcLsFX{gaps: true}
	tt := xcmpfx(t, cmp)
	tt.FireResize(5, 5)
	t.Eq("•••••\n•1st•\n•2nd•\n•3rd•\n•••••", tt.Screen())

	tt.FireKey(Down)
	t.Eq("•••••\n•3rd•\n•4th•\n•5th•\n•••••", tt.Screen())
	tt.Lines.Update(cmp, nil, func(e *Env) {
		t.True(cmp.LL.By(1).IsFlagged(Highlighted))
	})

	tt.FireKey(Down)
	t.Eq("•••••\n•6th•\n•7th•\n•8th•\n•••••", tt.Screen())
	tt.Lines.Update(cmp, nil, func(e *Env) {
		t.Not.True(cmp.LL.By(1).IsFlagged(Highlighted))
		t.True(cmp.LL.By(2).IsFlagged(Highlighted))
	})
}

func (s *ASourcedComponent) Scrolls_to_previous_focusable_within_gaps(
	t *T,
) {
	cmp := &srcLsFX{gaps: true}
	tt := xcmpfx(t, cmp)
	tt.FireResize(5, 5)
	t.Eq("•••••\n•1st•\n•2nd•\n•3rd•\n•••••", tt.Screen())

	tt.FireKey(Up)
	t.Eq("•••••\n•6th•\n•7th•\n•8th•\n•••••", tt.Screen())
	tt.Lines.Update(cmp, nil, func(e *Env) {
		t.True(cmp.LL.By(2).IsFlagged(Highlighted))
	})

	tt.FireKey(Up)
	t.Eq("•••••\n•4th•\n•5th•\n•6th•\n•••••", tt.Screen())
	tt.Lines.Update(cmp, nil, func(e *Env) {
		t.Not.True(cmp.LL.By(2).IsFlagged(Highlighted))
		t.True(cmp.LL.By(0).IsFlagged(Highlighted))
	})
}

func (s *ASourcedComponent) Remembers_highlighted_line_on_scrolling(
	t *T,
) {
	cmp := &srcLsFX{gaps: true}
	tt := xcmpfx(t, cmp)
	tt.FireResize(5, 5)
	t.Eq("•••••\n•1st•\n•2nd•\n•3rd•\n•••••", tt.Screen())

	tt.FireKey(Down)
	t.Eq("•••••\n•3rd•\n•4th•\n•5th•\n•••••", tt.Screen())
	tt.Lines.Update(cmp, nil, func(e *Env) {
		t.True(cmp.LL.By(1).IsFlagged(Highlighted))
	})

	tt.Lines.Update(cmp, nil, func(e *Env) {
		cmp.Scroll.Up()
	})
	t.Eq("•••••\n•1st•\n•2nd•\n•3rd•\n•••••", tt.Screen())
	tt.Lines.Update(cmp, nil, func(e *Env) {
		for i := 0; i < 3; i++ {
			t.Not.True(cmp.LL.By(i).IsFlagged(Highlighted))
		}
	})

	tt.Lines.Update(cmp, nil, func(e *Env) {
		cmp.Scroll.Down()
	})
	t.Eq("•••••\n•3rd•\n•4th•\n•5th•\n•••••", tt.Screen())
	tt.Lines.Update(cmp, nil, func(e *Env) {
		t.True(cmp.LL.By(1).IsFlagged(Highlighted))
	})

	tt.Lines.Update(cmp, nil, func(e *Env) {
		cmp.Scroll.Down()
	})
	t.Eq("•••••\n•5th•\n•6th•\n•7th•\n•••••", tt.Screen())
	tt.Lines.Update(cmp, nil, func(e *Env) {
		for i := 0; i < 3; i++ {
			t.Not.True(cmp.LL.By(i).IsFlagged(Highlighted))
		}
	})

	tt.Lines.Update(cmp, nil, func(e *Env) {
		cmp.Scroll.Up()
	})
	t.Eq("•••••\n•3rd•\n•4th•\n•5th•\n•••••", tt.Screen())
	tt.Lines.Update(cmp, nil, func(e *Env) {
		t.True(cmp.LL.By(1).IsFlagged(Highlighted))
	})
}

func (s *ASourcedComponent) Inverts_bg_fg_of_focused_if_highlighted(
	t *T,
) {
	cmp := &srcFcsFX{
		onLf: func(c *srcFcsFX, _ *Env, _ int) {
			switch c.lfN {
			case 1:
				t.Eq(1, c.LL.Focus.Current())
			}
		}}
	tt := xcmpfx(t, cmp)
	tt.FireResize(3, 2)
	t.FatalOn(tt.Lines.Update(cmp, nil, func(e *Env) {
		cmp.FF.Add(HighlightedFocusable)
		liner := &focusLinerFX{}
		liner.cc = []string{"1st", "2nd"}
		liner.notFocusable = map[int]bool{0: true}
		cmp.Src = &ContentSource{Liner: liner}
	}))

	l2 := tt.CellsOf(cmp).Trimmed()[1]
	for x := range l2 {
		t.Not.True(l2.HasAA(x, Reverse))
	}

	tt.FireKey(Down)
	t.Eq(1, cmp.lfN)

	l2 = tt.CellsOf(cmp).Trimmed()[1]
	for x := range l2 {
		if x < len("line 2") {
			t.True(l2.HasAA(x, Reverse))
			continue
		}
		t.Not.True(l2.HasAA(x, Reverse))
	}
}

type srcLsFX struct {
	srcFcsFX
	onLS func(*srcLsFX, *Env, int)
	lsN  int
	gaps bool
}

// OnInit initializes a source with eight lines of which the 4th and the
// 8th are focusable.
func (c *srcLsFX) OnInit(e *Env) {
	c.FF.Add(LinesSelectable)
	if c.gaps {
		Print(c.Gaps(0).Filling(), '•')
		fmt.Fprint(c.Gaps(0).Corners, "•")
	}
	liner := &focusLinerFX{}
	liner.cc = []string{
		"1st", "2nd", "3rd", "4th", "5th", "6th", "7th", "8th"}
	liner.notFocusable = map[int]bool{
		0: true, 1: true, 2: true, 4: true, 5: true, 6: true}
	c.Src = &ContentSource{Liner: liner}
}

func (c *srcLsFX) OnLineSelection(e *Env, idx int) {
	c.lsN++
	if c.onLS == nil {
		return
	}
	c.onLS(c, e, idx)
}

func (s *ASourcedComponent) Gets_its_selected_lines_reported(t *T) {
	cmp := &srcLsFX{}
	cmp.onLf = func(c *srcFcsFX, _ *Env, _ int) {
		switch c.lfN {
		case 1:
			t.Eq(3, c.LL.Focus.Current())
		case 2:
			t.Eq(7, c.LL.Focus.Current())
		}
	}
	cmp.onLS = func(c *srcLsFX, _ *Env, i int) { t.Eq(7, i) }
	tt := xcmpfx(t, cmp)
	tt.FireResize(3, 2)

	tt.FireRune('j')
	tt.FireRune('j')
	tt.FireKey(Enter)

	t.Eq(2, cmp.lfN)
	t.Eq(1, cmp.lsN)
}

func TestALiner(t *testing.T) {
	t.Parallel()
	Run(&ASourcedComponent{}, t)
}
