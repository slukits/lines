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

type LineHighlight struct {
	Suite
	Fixtures
}

func (s *LineHighlight) SetUp(t *T) {
	t.Parallel()
}

func (s *LineHighlight) TearDown(t *T) {
	ee, ok := s.Get(t).(*Events)
	if !ok {
		return
	}
	if !ee.IsListening() {
		return
	}
	ee.QuitListening()
}

func (s *LineHighlight) iFX(
	t *T, init func(*icmpFX, *Env), max ...int,
) (*Events, *Testing, *icmpFX) {
	fx := &icmpFX{init: init}
	ee, tt := Test(t.GoT(), fx, max...)
	s.Set(t, ee)
	ee.Listen()
	return ee, tt, fx
}

func (s *LineHighlight) Has_initially_no_line_highlighted(t *T) {
	s.iFX(t, func(c *icmpFX, e *Env) {
		t.Eq(-1, c.Highlight.Current())
		// no-ops for coverage
		c.Highlight.Next()
		c.Highlight.Reset()
	}, 1)
}

func (s *LineHighlight) Highlights_first_selectable_line(t *T) {
	ee, _, fx := s.iFX(t, func(c *icmpFX, e *Env) {
		fmt.Fprint(e, "first\nsecond")
		c.Highlight.Next()
		t.Eq(0, c.Highlight.Current())
	}, 2)

	ee.Update(fx, nil, func(e *Env) {
		fx.Highlight.Reset()
		fmt.Fprint(e.LL(0, NotSelectable), "first")
		fmt.Fprint(e.LL(1), "second")
		fx.Highlight.Next()
		t.Eq(1, fx.Highlight.Current())
	})
}

func (s *LineHighlight) Highlights_next_selectable_line(t *T) {
	s.iFX(t, func(c *icmpFX, e *Env) {
		fmt.Fprint(e.LL(0), "first")
		fmt.Fprint(e.LL(1, NotSelectable), "second")
		fmt.Fprint(e.LL(2), "third")
		c.Highlight.Next()
		t.Eq(0, c.Highlight.Current())
		c.Highlight.Next()
		t.Eq(2, c.Highlight.Current())
	}, 1)
}

func (s *LineHighlight) Resets_if_no_next_selectable(t *T) {
	s.iFX(t, func(c *icmpFX, e *Env) {
		fmt.Fprint(e.LL(0), "first")
		t.Eq(0, c.Highlight.Next())

		t.Eq(-1, c.Highlight.Next())

		fmt.Fprint(e.LL(1, NotSelectable), "second")
		fmt.Fprint(e.LL(2, NotSelectable), "third")

		t.Eq(0, c.Highlight.Next())
		t.Eq(-1, c.Highlight.Next())
	}, 1)
}

func (s *LineHighlight) Highlights_previous_selectable_line(t *T) {
	s.iFX(t, func(c *icmpFX, e *Env) {
		fmt.Fprint(e, "first\nsecond")
		t.Eq(-1, c.Highlight.Previous())
		t.Eq(0, c.Highlight.Next())
		t.Eq(-1, c.Highlight.Previous())

		fmt.Fprint(e.LL(0, NotSelectable), "first")
		fmt.Fprint(e.LL(1), "second")
		fmt.Fprint(e.LL(2), "third")
		t.Eq(1, c.Highlight.Next())
		t.Eq(-1, c.Highlight.Previous())

		t.Eq(1, c.Highlight.Next())
		t.Eq(2, c.Highlight.Next())
		t.Eq(1, c.Highlight.Previous())
	}, 1)
}

func (s *LineHighlight) Next_triggered_by_selectable_lines_feat(t *T) {
	ee, tt, fx := s.iFX(t, func(c *icmpFX, e *Env) {
		c.FF.Add(LinesSelectable)
		fmt.Fprint(e, "first\nsecond")
	}, 4)
	tt.FireKey(tcell.KeyDown, tcell.ModNone)
	tt.FireRune('j')
	ee.Update(fx, nil, func(e *Env) {
		t.Eq(1, fx.Highlight.Current())
	})
}

func (s *LineHighlight) Previous_triggered_by_selectable_lines_feat(t *T) {
	ee, tt, fx := s.iFX(t, func(c *icmpFX, e *Env) {
		c.FF.Add(LinesSelectable)
		fmt.Fprint(e, "first\nsecond\nthird")
		c.Highlight.Next()
		c.Highlight.Next()
		t.Eq(2, c.Highlight.Next())
	}, 4)
	tt.FireKey(tcell.KeyUp, tcell.ModNone)
	tt.FireRune('k')
	ee.Update(fx, nil, func(e *Env) {
		t.Eq(0, fx.Highlight.Current())
	})
}

func (s *LineHighlight) Reset_triggered_by_selectable_lines_feat(t *T) {
	ee, tt, fx := s.iFX(t, func(c *icmpFX, e *Env) {
		c.FF.Add(LinesSelectable)
		fmt.Fprint(e, "first\nsecond")
		t.Eq(0, c.Highlight.Next())
	}, 3)
	tt.FireKey(tcell.KeyEsc, tcell.ModNone)
	ee.Update(fx, nil, func(e *Env) {
		t.Eq(-1, fx.Highlight.Current())
	})
}

type lsCmpFx struct {
	*icmpFX
	ls func(*lsCmpFx, *Env, int)
}

func (s *lsCmpFx) OnLineSelection(e *Env, idx int) {
	if s.ls == nil {
		return
	}
	s.ls(s, e, idx)
}

func (s *LineHighlight) slFX(
	t *T, init func(*icmpFX, *Env), ls func(*lsCmpFx, *Env, int), max ...int,
) (*Events, *Testing, *lsCmpFx) {
	fx := &lsCmpFx{icmpFX: &icmpFX{init: init}, ls: ls}
	ee, tt := Test(t.GoT(), fx, max...)
	s.Set(t, ee)
	ee.Listen()
	return ee, tt, fx
}

func (s *LineHighlight) Reports_highlighted_line_by_sl_feature(t *T) {
	reported := false
	_, tt, _ := s.slFX(t,
		func(c *icmpFX, e *Env) {
			c.FF.Add(LinesSelectable)
			fmt.Fprint(e, "first\nsecond")
			c.Highlight.Next()
			t.Eq(1, c.Highlight.Next())
		},
		func(lcf *lsCmpFx, e *Env, i int) {
			t.Eq(1, i)
			reported = true
		},
		2,
	)
	tt.FireKey(tcell.KeyEnter, tcell.ModNone)
	t.True(reported)
}

func (s *LineHighlight) Scrolls_to_next_highlighted_line(t *T) {
	ee, tt, fx := s.iFX(t, func(c *icmpFX, e *Env) {
		c.FF.Add(LinesSelectable)
		c.dim.SetHeight(2)
		fmt.Fprint(e.LL(0, NotSelectable), "line 1")
		fmt.Fprint(e.LL(1, NotSelectable), "line 2")
		fmt.Fprint(e.LL(2, NotSelectable), "line 3")
		fmt.Fprint(e.LL(3, NotSelectable), "line 4")
		fmt.Fprint(e.LL(4, NotSelectable), "line 5")
		fmt.Fprint(e.LL(5), "line 6")
		fmt.Fprint(e.LL(6), "line 7")
	}, 3)
	tt.FireKey(tcell.KeyDown, tcell.ModNone)
	ee.Update(fx, nil, func(e *Env) {
		t.Eq("line 5\nline 6", tt.Screen().String())
	})
}

func (s *LineHighlight) Scrolls_to_previous_highlighted_line(t *T) {
	ee, tt, fx := s.iFX(t, func(c *icmpFX, e *Env) {
		c.FF.Add(LinesSelectable)
		c.dim.SetHeight(2)
		fmt.Fprint(e.LL(0, NotSelectable), "line 1")
		fmt.Fprint(e.LL(1), "line 2")
		fmt.Fprint(e.LL(2, NotSelectable), "line 3")
		fmt.Fprint(e.LL(3, NotSelectable), "line 4")
		fmt.Fprint(e.LL(4, NotSelectable), "line 5")
		fmt.Fprint(e.LL(5), "line 6")
		fmt.Fprint(e.LL(6), "line 7")
	}, 5)
	tt.FireKey(tcell.KeyDown, tcell.ModNone)
	tt.FireKey(tcell.KeyDown, tcell.ModNone)
	ee.Update(fx, nil, func(e *Env) {
		t.Eq("line 5\nline 6", tt.Screen().String())
	})

	tt.FireKey(tcell.KeyUp, tcell.ModNone)
	ee.Update(fx, nil, func(e *Env) {
		t.Eq("line 2\nline 3", tt.Screen().String())
	})
}

func TestLineHighlight(t *testing.T) {
	t.Parallel()
	Run(&LineHighlight{}, t)
}
