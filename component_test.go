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

type cmpFX struct{ Component }

type AComponent struct{ Suite }

func (s *AComponent) SetUp(t *T) { t.Parallel() }

func (s *AComponent) Access_panics_outside_event_processing(t *T) {
	cmp := &cmpFX{}
	TermFixture(t.GoT(), 0, cmp)
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
	tt, _ := cmpfx(t)
	t.FatalOn(tt.Lines.Update(tt.Root(), nil, func(e *Env) {
		fx := tt.Root().(*cmpFX)
		t.Eq(0, fx.Len())
		fmt.Fprint(e, "first\nsecond\nthird")
		t.Eq(3, fx.Len())
	}))
}

func (s *AComponent) Doesnt_change_line_count_on_line_overwrite(t *T) {
	tt, cmp := cmpfx(t)
	t.FatalOn(tt.Lines.Update(cmp, nil, func(e *Env) {
		cmp.LL.Mod(Overwriting)
		fmt.Fprint(e, "two\nlines")
		t.Eq(2, cmp.Len())
		fmt.Fprint(e, "one line")
		t.Eq(2, cmp.Len())
	}))

	// but second line is empty now
	t.Eq("one line", tt.Screen().Trimmed().String())
}

func (s *AComponent) Has_a_line_more_after_appending_an_line(t *T) {
	tt, cmp := cmpfx(t)
	t.FatalOn(tt.Lines.Update(cmp, nil, func(e *Env) {
		cmp.LL.Mod(Appending)
		fmt.Fprint(e, "two\nlines")
		fmt.Fprint(e, "one line")
		t.Eq(3, cmp.Len())
	}))
}

func (s *AComponent) Has_a_line_more_after_writing_to_tailing(t *T) {
	tt, cmp := cmpfx(t)
	t.FatalOn(tt.Lines.Update(cmp, nil, func(e *Env) {
		cmp.LL.Mod(Tailing)
		fmt.Fprint(e, "two\nlines")
		fmt.Fprint(e, "one line")
		t.Eq(3, cmp.Len())
	}))
}

func (s *AComponent) Shows_last_line_clipped_above_if_tailing(t *T) {
	tt, cmp := cmpfx(t)
	tt.FireResize(20, 2)
	t.FatalOn(tt.Lines.Update(cmp, nil, func(e *Env) {
		cmp.LL.Mod(Tailing)
		fmt.Fprint(e, "three\nlines\nat last")
	}))
	t.Eq("lines  \nat last", tt.Screen().Trimmed().String())
}

func (s *AComponent) Blanks_a_reset_line(t *T) {
	tt, cmp := cmpfx(t)
	tt.FireResize(20, 2)
	t.FatalOn(tt.Lines.Update(cmp, nil, func(e *Env) {
		fmt.Fprint(e, "first\nsecond")
	}))
	t.Eq("first \nsecond", tt.Screen().Trimmed().String())

	t.FatalOn(tt.Lines.Update(cmp, nil, func(e *Env) {
		cmp.Reset(-2) // no-op, coverage
		cmp.Reset(0)
	}))

	t.Eq("second", tt.Screen().Trimmed().String())
}

func (s *AComponent) fxCmp(t *T) (*Fixture, *cmpFX) {
	cmp := &cmpFX{}
	tt := TermFixture(t.GoT(), 0, cmp)
	return tt, cmp
}

func (s *AComponent) Truncates_lines_to_screen_area_on_reset_all(t *T) {
	tt, fx := s.fxCmp(t)
	tt.FireResize(20, 2)
	t.FatalOn(tt.Lines.Update(fx, nil, func(e *Env) {
		fmt.Fprint(e, "first\nsecond\nthird\nforth")
		t.Eq(4, fx.Len())
	}))
	t.FatalOn(tt.Lines.Update(fx, nil, func(e *Env) {
		fx.Reset(All)
		t.Eq(2, fx.Len())
	}))
}

func (s *AComponent) Scrolls_by_one_line_if_height_is_one(t *T) {
	tt, fx := s.fxCmp(t)
	t.FatalOn(tt.Lines.Update(fx, nil, func(e *Env) {
		fx.Dim().SetHeight(1)
		fmt.Fprint(e, "first\nsecond")
	}))
	t.Eq("first", tt.Screen().Trimmed())

	t.FatalOn(tt.Lines.Update(fx, nil, func(e *Env) {
		fx.Scroll.Down()
	}))
	t.Eq("second", tt.Screen().Trimmed())
	t.FatalOn(tt.Lines.Update(fx, nil, func(e *Env) {
		fx.Scroll.Down()
	}))
	t.Eq("second", tt.Screen().Trimmed())

	t.FatalOn(tt.Lines.Update(fx, nil, func(e *Env) {
		fx.Scroll.Up()
	}))
	t.Eq("first", tt.Screen().Trimmed())
	t.FatalOn(tt.Lines.Update(fx, nil, func(e *Env) {
		fx.Scroll.Up()
	}))
	t.Eq("first", tt.Screen().Trimmed())
}

func (s *AComponent) Scrolls_to_last_line_if_last_displayed(t *T) {
	tt, fx := s.fxCmp(t)
	t.FatalOn(tt.Lines.Update(fx, nil, func(e *Env) {
		fx.Dim().SetHeight(3)
		fmt.Fprint(e, "first\nsecond\nthird\nforth")
	}))
	t.Eq("first \nsecond\nthird ", tt.Screen().Trimmed())

	t.FatalOn(tt.Lines.Update(fx, nil, func(e *Env) {
		fx.Scroll.Down()
	}))
	t.Eq("second\nthird \nforth ", tt.Screen().Trimmed())
}

func (s *AComponent) Scrolls_to_first_line_if_first_displayed(t *T) {
	tt, fx := s.fxCmp(t)
	t.FatalOn(tt.Lines.Update(fx, nil, func(e *Env) {
		fx.Dim().SetHeight(3)
		fmt.Fprint(e, "first\nsecond\nthird\nforth")
	}))
	t.Eq("first \nsecond\nthird ", tt.Screen().Trimmed().String())
	t.FatalOn(tt.Lines.Update(fx, nil, func(e *Env) {
		fx.Scroll.Down()
	}))
	t.Eq("second\nthird \nforth ", tt.Screen().Trimmed().String())

	t.FatalOn(tt.Lines.Update(fx, nil, func(e *Env) {
		fx.Scroll.Up()
	}))
	t.Eq("first \nsecond\nthird ", tt.Screen().Trimmed().String())
}

func (s *AComponent) Scrolls_down_by_90_percent_height(t *T) {
	tt, fx := s.fxCmp(t)
	tt.FireResize(20, 30)
	t.FatalOn(tt.Lines.Update(fx, nil, func(e *Env) {
		fx.Dim().SetHeight(5)
		ll := make([]string, 60)
		for i := 0; i < 60; i++ {
			ll[i] = fmt.Sprintf("line %d", i+1)
		}
		fmt.Fprint(e, strings.Join(ll, "\n"))
		t.Eq(60, fx.Len())
	}))

	exp := []string{}
	for i := 0; i < 5; i++ {
		exp = append(exp, fmt.Sprintf("line %d", i+5))
	}
	t.FatalOn(tt.Lines.Update(fx, nil, func(e *Env) {
		fx.Scroll.Down()
	}))
	t.Eq(strings.Join(exp, "\n"), tt.Screen().Trimmed().String())

	t.FatalOn(tt.Lines.Update(fx, nil, func(e *Env) {
		fx.Dim().SetHeight(15)
	}))

	exp = []string{}
	for i := 0; i < 15; i++ { // first is still at fifth line
		exp = append(exp, fmt.Sprintf("line %d", i+19))
	}
	t.FatalOn(tt.Lines.Update(fx, nil, func(e *Env) {
		fx.Scroll.Down()
	}))
	t.Eq(strings.Join(exp, "\n"), tt.Screen().Trimmed().String())

	t.FatalOn(tt.Lines.Update(fx, nil, func(e *Env) {
		fx.Scroll.ToTop()
		t.True(fx.Scroll.IsAtTop())
		fx.Dim().SetHeight(30)
	}))
	exp = []string{}
	for i := 0; i < 30; i++ { // first is still at fifth line
		exp = append(exp, fmt.Sprintf("line %d", i+28))
	}
	t.FatalOn(tt.Lines.Update(fx, nil, func(e *Env) {
		fx.Scroll.Down()
	}))
	t.Eq(strings.Join(exp, "\n"), tt.Screen().Trimmed().String())
}

func (s *AComponent) Scrolls_up_by_90_percent_height(t *T) {
	tt, fx := s.fxCmp(t)
	tt.FireResize(20, 30)
	t.FatalOn(tt.Lines.Update(fx, nil, func(e *Env) {
		fx.Dim().SetHeight(5)
		ll := make([]string, 60)
		for i := 0; i < 60; i++ {
			if i+1 < 10 {
				ll[i] = fmt.Sprintf("line 0%d", i+1)
				continue
			}
			ll[i] = fmt.Sprintf("line %d", i+1)
		}
		fmt.Fprint(e, strings.Join(ll, "\n"))
		t.Eq(60, fx.Len())
		fx.Scroll.ToBottom()
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
	t.FatalOn(tt.Lines.Update(fx, nil, func(e *Env) {
		fx.Scroll.Up()
	}))
	t.Eq(strings.Join(exp, "\n"), tt.Screen().Trimmed().String())

	t.FatalOn(tt.Lines.Update(fx, nil, func(e *Env) {
		fx.Dim().SetHeight(15)
	}))

	exp = []string{}
	for i := 0; i < 15; i++ {
		exp = append(exp, fmt.Sprintf("line %d", i+32))
	}
	t.FatalOn(tt.Lines.Update(fx, nil, func(e *Env) {
		fx.Scroll.Up()
	}))
	t.Eq(strings.Join(exp, "\n"), tt.Screen().Trimmed().String())

	t.FatalOn(tt.Lines.Update(fx, nil, func(e *Env) {
		fx.Dim().SetHeight(30)
		fx.Scroll.ToBottom()
		t.True(fx.Scroll.IsAtBottom())
	}))
	exp = []string{}
	for i := 0; i < 30; i++ {
		if i+4 < 10 {
			exp = append(exp, fmt.Sprintf("line 0%d", i+4))
			continue
		}
		exp = append(exp, fmt.Sprintf("line %d", i+4))
	}
	t.FatalOn(tt.Lines.Update(fx, nil, func(e *Env) {
		fx.Scroll.Up()
	}))
	t.Eq(strings.Join(exp, "\n"), tt.Screen().Trimmed().String())
}

func (s *AComponent) Scrolls_to_top_on_reset_all(t *T) {
	tt, fx := s.fxCmp(t)
	tt.FireResize(20, 2)
	t.FatalOn(tt.Lines.Update(fx, nil, func(e *Env) {
		fmt.Fprint(e, "first\nsecond\nthird\nforth")
		fx.Scroll.ToBottom()
		t.Not.True(fx.Scroll.IsAtTop())
	}))
	t.FatalOn(tt.Lines.Update(fx, nil, func(e *Env) {
		fx.Reset(All)
		t.True(fx.Scroll.IsAtTop())
	}))
}

type uiCmpFX struct {
	Component
	init func(c *uiCmpFX, e *Env)
}

func (c *uiCmpFX) OnInit(e *Env) {
	if c.init == nil {
		return
	}
	c.init(c, e)
}

func (c *uiCmpFX) OnUpdate(e *Env) {
	data := e.Evt.(*UpdateEvent).Data.(map[int]string)
	for idx, content := range data {
		fmt.Fprint(e.LL(idx), content)
	}
	for i := 0; i < c.Len(); i++ {
		if _, ok := data[i]; ok {
			continue
		}
		c.Reset(i)
	}
}

type fillerFX struct{ Component }

type rplStackFX struct {
	Component
	Stacking
	long  string
	short string
}

func (c *rplStackFX) OnInit(_ *Env) {
	c.CC = []Componenter{&fillerFX{}, &icmpFX{
		init: func(ic *icmpFX, e *Env) {
			ic.Dim().SetHeight(1)
			fmt.Fprint(e, c.long)
		}}}
}

func (c *rplStackFX) OnUpdate(e *Env, data interface{}) {
	c.CC[1] = data.(Componenter)
}

func (s *AComponent) Is_replaceable(t *T) {
	fx := &rplStackFX{
		long:  "a rather long long long line",
		short: "a short line",
	}
	tt := xcmpfx(t, fx)
	t.Eq(fx.long, tt.ScreenOf(fx).Trimmed().String())

	t.FatalOn(tt.Lines.Update(fx, &icmpFX{
		init: func(ic *icmpFX, e *Env) {
			ic.Dim().SetHeight(1)
			fmt.Fprint(e, fx.short)
		}}, nil))

	str := tt.ScreenOf(fx).Trimmed().String()
	t.Eq(fx.short, str)
}

func (s *AComponent) Updates_tab_expansions_on_tab_width_change(t *T) {
	tt, cmp := cmpfx(t)
	tt.FireResize(11, 2)

	tb, exp, expTB := "", 8, strings.Repeat(" ", 8)
	tt.Lines.Update(cmp, nil, func(e *Env) {
		tb = strings.Repeat(" ", cmp.Globals().TabWidth())
		fmt.Fprint(e, "\t1st\n\t2nd")
	})
	t.Not.Eq(tb, expTB)
	t.True(strings.HasPrefix(tt.Screen()[0], tb+"1st"))
	t.True(strings.HasPrefix(tt.Screen()[1], tb+"2nd"))

	tt.Lines.Update(cmp, nil, func(e *Env) {
		cmp.Globals().SetTabWidth(exp)
	})

	t.True(strings.HasPrefix(tt.Screen()[0], expTB+"1st"))
	t.True(strings.HasPrefix(tt.Screen()[1], expTB+"2nd"))
}

func TestComponent(t *testing.T) {
	t.Parallel()
	Run(&AComponent{}, t)
}
