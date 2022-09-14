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

type cmpFX struct{ Component }

type _component struct{ Suite }

func (s *_component) SetUp(t *T) { t.Parallel() }

func (s *_component) Access_panics_outside_event_processing(t *T) {
	cmp := &cmpFX{}
	ee, _ := Test(t.GoT(), cmp, 0)
	defer ee.QuitListening()
	t.Panics(func() { cmp.Dim().SetHeight(20) })
}

func (s *_component) Has_same_line_count_if_one_line_overwrite(t *T) {
	cmp := &cmpFX{}
	ee, tt := Test(t.GoT(), cmp)
	ee.Update(cmp, nil, func(e *Env) {
		cmp.Mod(Overwriting)
		fmt.Fprint(e, "two\nlines")
		t.Eq(2, cmp.Len())
		fmt.Fprint(e, "one line")
		t.Eq(2, cmp.Len())
	})
	t.False(ee.IsListening())
	// but second line is empty now
	t.Eq("one line", tt.LastScreen.String())
}

func (s *_component) Has_a_line_more_after_appending_an_line(t *T) {
	cmp := &cmpFX{}
	ee, _ := Test(t.GoT(), cmp)
	ee.Update(cmp, nil, func(e *Env) {
		cmp.Mod(Appending)
		fmt.Fprint(e, "two\nlines")
		fmt.Fprint(e, "one line")
		t.Eq(3, cmp.Len())
	})
}

func (s *_component) Has_a_line_more_after_writing_to_tailing(t *T) {
	cmp := &cmpFX{}
	ee, _ := Test(t.GoT(), cmp)
	ee.Update(cmp, nil, func(e *Env) {
		cmp.Mod(Tailing)
		fmt.Fprint(e, "two\nlines")
		fmt.Fprint(e, "one line")
		t.Eq(3, cmp.Len())
	})
}

func (s *_component) Shows_last_line_clips_above_if_tailing(t *T) {
	cmp := &cmpFX{}
	ee, tt := Test(t.GoT(), cmp)
	tt.FireResize(20, 2)
	ee.Update(cmp, nil, func(e *Env) {
		cmp.Mod(Tailing)
		fmt.Fprint(e, "three\nlines\nat last")
	})
	t.Eq("lines  \nat last", tt.LastScreen.String())
}

func (s *_component) Blanks_a_reset_line(t *T) {
	cmp := &cmpFX{}
	ee, tt := Test(t.GoT(), cmp, 2)
	tt.FireResize(20, 2)
	ee.Update(cmp, nil, func(e *Env) {
		fmt.Fprint(e, "first\nsecond")
	})
	t.Eq("first \nsecond", tt.Screen().String())

	ee.Update(cmp, nil, func(e *Env) {
		cmp.Reset(-2) // no-op, coverage
		cmp.Reset(0)
	})

	t.Eq("second", tt.LastScreen.String())
}

func (s *_component) fxCmp(
	t *T, countdown ...int,
) (*Events, *Testing, *cmpFX) {
	cmp := &cmpFX{}
	ee, tt := Test(t.GoT(), cmp, countdown...)
	return ee, tt, cmp
}

func (s *_component) Scrolls_by_one_line_if_height_is_one(t *T) {
	ee, tt, fx := s.fxCmp(t, 5)
	ee.Update(fx, nil, func(e *Env) {
		fx.Dim().SetHeight(1)
		fmt.Fprint(e, "first\nsecond")
	})
	t.Eq("first", tt.Screen().String())

	ee.Update(fx, nil, func(e *Env) { fx.Scroll.Down() })
	t.Eq("second", tt.Screen().String())
	ee.Update(fx, nil, func(e *Env) { fx.Scroll.Down() })
	t.Eq("second", tt.Screen().String())

	ee.Update(fx, nil, func(e *Env) { fx.Scroll.Up() })
	t.Eq("first", tt.Screen().String())
	ee.Update(fx, nil, func(e *Env) { fx.Scroll.Up() })
	t.Eq("first", tt.LastScreen.String())
}

func (s *_component) Scrolls_to_last_line_is_last_displayed(t *T) {
	ee, tt, fx := s.fxCmp(t, 2)
	ee.Update(fx, nil, func(e *Env) {
		fx.Dim().SetHeight(3)
		fmt.Fprint(e, "first\nsecond\nthird\nforth")
	})
	t.Eq("first \nsecond\nthird ", tt.Screen().String())

	ee.Update(fx, nil, func(e *Env) { fx.Scroll.Down() })
	t.Eq("second\nthird \nforth ", tt.LastScreen.String())
}

func (s *_component) Scrolls_to_first_line_is_first_displayed(t *T) {
	ee, tt, fx := s.fxCmp(t, 3)
	ee.Update(fx, nil, func(e *Env) {
		fx.Dim().SetHeight(3)
		fmt.Fprint(e, "first\nsecond\nthird\nforth")
	})
	t.Eq("first \nsecond\nthird ", tt.Screen().String())
	ee.Update(fx, nil, func(e *Env) { fx.Scroll.Down() })
	t.Eq("second\nthird \nforth ", tt.Screen().String())

	ee.Update(fx, nil, func(e *Env) { fx.Scroll.Up() })
	t.Eq("first \nsecond\nthird ", tt.LastScreen.String())
}

func (s *_component) Scrolls_down_by_90_percent_height(t *T) {
	ee, tt, fx := s.fxCmp(t, 6)
	tt.FireResize(20, 30)
	ee.Update(fx, nil, func(e *Env) {
		fx.Dim().SetHeight(5)
		ll := make([]string, 60)
		for i := 0; i < 60; i++ {
			ll[i] = fmt.Sprintf("line %d", i+1)
		}
		fmt.Fprint(e, strings.Join(ll, "\n"))
		t.Eq(60, fx.Len())
	})

	exp := []string{}
	for i := 0; i < 5; i++ {
		exp = append(exp, fmt.Sprintf("line %d", i+5))
	}
	ee.Update(fx, nil, func(e *Env) { fx.Scroll.Down() })
	t.Eq(strings.Join(exp, "\n"), tt.Screen().String())

	ee.Update(fx, nil, func(e *Env) { fx.Dim().SetHeight(15) })

	exp = []string{}
	for i := 0; i < 15; i++ { // first is still at fifth line
		exp = append(exp, fmt.Sprintf("line %d", i+19))
	}
	ee.Update(fx, nil, func(e *Env) { fx.Scroll.Down() })
	t.Eq(strings.Join(exp, "\n"), tt.Screen().String())

	ee.Update(fx, nil, func(e *Env) {
		fx.Scroll.ToTop()
		t.True(fx.Scroll.IsAtTop())
		fx.Dim().SetHeight(30)
	})
	exp = []string{}
	for i := 0; i < 30; i++ { // first is still at fifth line
		exp = append(exp, fmt.Sprintf("line %d", i+28))
	}
	ee.Update(fx, nil, func(e *Env) { fx.Scroll.Down() })
	t.Eq(strings.Join(exp, "\n"), tt.LastScreen.String())
}

func (s *_component) Scrolls_up_by_90_percent_height(t *T) {
	ee, tt, fx := s.fxCmp(t, 6)
	tt.FireResize(20, 30)
	ee.Update(fx, nil, func(e *Env) {
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
	})

	exp := []string{}
	for i := 0; i < 5; i++ {
		exp = append(exp, fmt.Sprintf("line %d", i+56))
	}
	t.Eq(strings.Join(exp, "\n"), tt.Screen().String())
	exp = []string{}
	for i := 0; i < 5; i++ {
		exp = append(exp, fmt.Sprintf("line %d", i+52))
	}
	ee.Update(fx, nil, func(e *Env) { fx.Scroll.Up() })
	t.Eq(strings.Join(exp, "\n"), tt.Screen().String())

	ee.Update(fx, nil, func(e *Env) { fx.Dim().SetHeight(15) })

	exp = []string{}
	for i := 0; i < 15; i++ { // first is still at 52nd line
		exp = append(exp, fmt.Sprintf("line %d", i+38))
	}
	ee.Update(fx, nil, func(e *Env) { fx.Scroll.Up() })
	t.Eq(strings.Join(exp, "\n"), tt.Screen().String())

	ee.Update(fx, nil, func(e *Env) {
		fx.Dim().SetHeight(30)
		fx.Scroll.ToBottom()
		t.True(fx.Scroll.IsAtBottom())
	})
	exp = []string{}
	for i := 0; i < 30; i++ {
		if i+4 < 10 {
			exp = append(exp, fmt.Sprintf("line 0%d", i+4))
			continue
		}
		exp = append(exp, fmt.Sprintf("line %d", i+4))
	}
	ee.Update(fx, nil, func(e *Env) { fx.Scroll.Up() })
	t.Eq(strings.Join(exp, "\n"), tt.LastScreen.String())
}

func (s *_component) Updates(t *T) {
	cmp := &uiCmpFX{init: func(c *uiCmpFX, e *Env) {
		fmt.Fprint(e, "initial value")
	}}
	ee, tt := Test(t.GoT(), cmp, 0)
	ee.Listen()
	defer ee.QuitListening()
	tt.FireResize(13, 7)
	str := strings.TrimSpace(tt.String())
	t.Eq("initial value", str)
	linesUpdate := map[int]string{
		0: "line 00",
		1: "line 01",
		2: "line 02",
		3: "line 03",
		4: "line 04",
	}
	if err := ee.Update(cmp, linesUpdate, nil); err != nil {
		t.Fatalf("gounit: view: update: lines: %v", err)
	}
	str = strings.TrimSpace(tt.String())
	exp := make([]string, 5)
	for i, v := range linesUpdate {
		exp[i] = fmt.Sprintf("%s      ", v)
	}
	t.Eq(strings.TrimSpace(strings.Join(exp, "\n")), str)
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

func (c *rplStackFX) OnUpdate(e *Env) {
	cmp := e.Evt.(*UpdateEvent).Data.(Componenter)
	c.CC[1] = cmp
}

func (s *_component) Is_replaceable(t *T) {
	fx := &rplStackFX{
		long:  "a rather long long long line",
		short: "a short line",
	}
	ee, tt := Test(t.GoT(), fx, 0)
	ee.Listen()
	defer ee.QuitListening()
	t.Eq(fx.long, tt.Trim(tt.ScreenOf(fx)).String())

	ee.Update(fx, &icmpFX{
		init: func(ic *icmpFX, e *Env) {
			ic.Dim().SetHeight(1)
			fmt.Fprint(e, fx.short)
		}}, nil)

	str := tt.Trim(tt.ScreenOf(fx)).String()
	t.Eq(fx.short, str)
}

func (s *_component) Fills_line_at_unit_separators(t *T) {
	fx := &icmpFX{init: func(c *icmpFX, e *Env) {
		c.Dim().SetHeight(1).SetWidth(8)
		fmt.Fprintf(e, "a%sb", LineFiller)
	}}
	ee, tt := Test(t.GoT(), fx, 3)
	ee.Listen()

	t.Eq("a      b", tt.Screen().String())

	ee.Update(fx, nil, func(e *Env) {
		fmt.Fprintf(e, "a%sb%[1]sc", LineFiller)
	})

	t.Eq("a   b  c", tt.Screen().String())

	ee.Update(fx, nil, func(e *Env) {
		fmt.Fprintf(e, "ab%scd%[1]sef%[1]sgh", LineFiller)
	})

	t.Eq("ab cd ef", tt.LastScreen.String())
}

func TestComponent(t *testing.T) {
	t.Parallel()
	Run(&_component{}, t)
}
