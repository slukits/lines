// Copyright (c) 2022 Stephan Lukits. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package lines

import (
	"fmt"
	"testing"

	. "github.com/slukits/gounit"
)

type AScroller struct{ Suite }

func (s *AScroller) SetUp(t *T) { t.Parallel() }

func (s *AScroller) Is_initially_at_top(t *T) {
	cmp := &cmpFX{
		onInit: func(c *cmpFX, e *Env) {
			c.FF.Set(Scrollable)
			c.Dim().SetWidth(4).SetHeight(2)
			fmt.Fprint(e, "12\n3456\n789")
		}}
	fx := fx(t, cmp)
	fx.Lines.Update(cmp, nil, func(e *Env) {
		t.True(cmp.Scroll.IsAtTop())
	})
	t.Eq("12  \n3456", fx.ScreenOf(cmp))
}

func (s *AScroller) Scrolls_down_on_page_down_key(t *T) {
	cmp := &cmpFX{
		onInit: func(c *cmpFX, e *Env) {
			c.FF.Set(Scrollable)
			c.Dim().SetWidth(4).SetHeight(2)
			fmt.Fprint(e, "12\n3456\n789")
		}}
	fx := fx(t, cmp).FireKey(PgDn)
	t.Eq("3456\n789 ", fx.ScreenOf(cmp))
}

func (s *AScroller) Is_at_top_on_scrolling_to_top(t *T) {
	cmp := &cmpFX{
		onInit: func(c *cmpFX, e *Env) {
			c.FF.Set(Scrollable)
			c.Dim().SetWidth(4).SetHeight(2)
			fmt.Fprint(e, "12\n3456\n789")
		}}
	fx := fx(t, cmp).FireKey(PgDn)
	t.Eq("3456\n789 ", fx.ScreenOf(cmp))

	fx.Lines.Update(cmp, nil, func(e *Env) {
		cmp.Scroll.ToTop()
	})
	t.Eq("12  \n3456", fx.ScreenOf(cmp))
}

func TestAScroller(t *testing.T) {
	t.Parallel()
	Run(&AScroller{}, t)
}

type AScrollBar struct{ Suite }

func (s *AScrollBar) SetUp(t *T) { t.Parallel() }

func (s *AScrollBar) Position_not_shown_if_higher_than_content(t *T) {
	cmp := &cmpFX{
		onInit: func(c *cmpFX, e *Env) {
			c.Scroll.Bar = true
			fmt.Fprint(e, "12\n3456\n789")
		}}
	fx, sdb := fx(t, cmp).FireResize(5, 3), ScrollBarDef{}
	t.FatalOn(fx.Lines.Update(cmp, nil, func(e *Env) {
		t.True(cmp.Len() <= cmp.Dim().Height())
		sdb = cmp.Globals().ScrollBarDef()
	}))
	t.Eq("12   \n3456 \n789  ", fx.Screen())
	for _, c := range fx.Cells().Column(4) {
		t.Eq(c.Style, sdb.Style)
	}
}

func (s *AScrollBar) Is_at_zero_position_if_at_top(t *T) {
	cmp := &cmpFX{
		onInit: func(c *cmpFX, e *Env) {
			c.Scroll.Bar = true
			fmt.Fprint(e, "12\n3456\n789")
		}}
	fx := fx(t, cmp).FireResize(5, 2)
	t.FatalIfNot(t.Eq("12   \n3456 ", fx.Screen()))
	sdb := fx.Scroll.BarDef(cmp)
	for i, c := range fx.Cells().Column(4) {
		switch i {
		case 0:
			t.FatalIfNot(t.Eq(c.Style, sdb.Position))
		default:
			t.FatalIfNot(t.Eq(c.Style, sdb.Style))
		}
	}
}

func (s *AScrollBar) Is_at_max_position_if_at_bottom(t *T) {
	cmp := &cmpFX{
		onInit: func(c *cmpFX, e *Env) {
			c.Scroll.Bar = true
			fmt.Fprint(e, "12\n3456\n789")
		}}
	fx := fx(t, cmp).FireResize(5, 2)
	t.FatalIfNot(t.Eq("12   \n3456 ", fx.Screen()))
	fx.Scroll.ToBottom(cmp)
	sdb := fx.Scroll.BarDef(cmp)
	max := len(fx.Cells().Column(4)) - 1
	for i, c := range fx.Cells().Column(4) {
		switch i {
		case max:
			t.FatalIfNot(t.Eq(c.Style, sdb.Position))
		default:
			t.FatalIfNot(t.Eq(c.Style, sdb.Style))
		}
	}
}

func (s *AScrollBar) Goes_according_position_down_on_scrolling(t *T) {
	cmp := &cmpFX{
		onInit: func(c *cmpFX, e *Env) {
			c.FF.Set(Scrollable)
			c.Scroll.Bar = true
			fmt.Fprint(e, "1st\n2nd\n3rd\n4th\n5th\n6th")
		}}
	fx := fx(t, cmp).FireResize(4, 3)
	t.FatalIfNot(t.Eq("1st \n2nd \n3rd ", fx.Screen()))
	fx.FireKeys(PgDn)
	t.FatalIfNot(t.Eq("3rd \n4th \n5th ", fx.Screen()))
	sdb := fx.Scroll.BarDef(cmp)
	for i, c := range fx.Cells().Column(3) {
		switch i {
		case 1:
			t.FatalIfNot(t.Eq(c.Style, sdb.Position))
		default:
			t.FatalIfNot(t.Eq(c.Style, sdb.Style))
		}
	}
}

func (s *AScrollBar) Goes_according_position_up_on_scrolling(t *T) {
	cmp := &cmpFX{
		onInit: func(c *cmpFX, e *Env) {
			c.FF.Set(Scrollable)
			c.Scroll.Bar = true
			fmt.Fprint(e, "1st\n2nd\n3rd\n4th\n5th\n6th")
		}}
	fx := fx(t, cmp).FireResize(4, 3)
	fx.Scroll.ToBottom(cmp)
	t.FatalIfNot(t.Eq("4th \n5th \n6th ", fx.Screen()))
	fx.FireKeys(PgUp)
	t.FatalIfNot(t.Eq("2nd \n3rd \n4th ", fx.Screen()))
	sdb := fx.Scroll.BarDef(cmp)
	for i, c := range fx.Cells().Column(3) {
		switch i {
		case 1:
			t.FatalIfNot(t.Eq(c.Style, sdb.Position))
		default:
			t.FatalIfNot(t.Eq(c.Style, sdb.Style))
		}
	}
}

func (s *AScrollBar) Knows_its_containing_coordinates(t *T) {
	cmp := &cmpFX{
		onInit: func(c *cmpFX, e *Env) {
			c.FF.Set(Scrollable)
			c.Scroll.Bar = true
			Print(c.Gaps(0).Top.At(0).Filling(), ' ')
			fmt.Fprint(e, "1st\n2nd\n3rd\n4th\n5th\n6th")
		}}
	fx := fx(t, cmp).FireResize(4, 4)
	fx.Lines.Update(cmp, nil, func(e *Env) {
		for i := 1; i < 4; i++ {
			t.True(cmp.Scroll.BarContains(3, i))
		}
	})
}

func (s *AScrollBar) Scrolls_down_on_a_left_click(t *T) {
	cmp := &cmpFX{
		onInit: func(c *cmpFX, e *Env) {
			c.FF.Set(Scrollable)
			c.Scroll.Bar = true
			Print(c.Gaps(0).Top.At(0).Filling(), ' ')
			fmt.Fprint(e, "1st\n2nd\n3rd\n4th\n5th\n6th")
		}}
	fx := fx(t, cmp).FireResize(4, 4)
	t.Eq("    \n1st \n2nd \n3rd ", fx.Screen())
	fx.FireClick(3, 1)
	t.Eq("    \n3rd \n4th \n5th ", fx.Screen())
}

func (s *AScrollBar) Scrolls_up_on_a_right_click(t *T) {
	cmp := &cmpFX{
		onInit: func(c *cmpFX, e *Env) {
			c.FF.Set(Scrollable)
			c.Scroll.Bar = true
			Print(c.Gaps(0).Top.At(0).Filling(), ' ')
			fmt.Fprint(e, "1st\n2nd\n3rd\n4th\n5th\n6th")
		}}
	fx := fx(t, cmp).FireResize(4, 4)
	fx.Scroll.ToBottom(cmp)
	t.Eq("    \n4th \n5th \n6th ", fx.Screen())
	fx.FireContext(3, 1)
	t.Eq("    \n2nd \n3rd \n4th ", fx.Screen())
}

func TestAScrollBar(t *testing.T) {
	t.Parallel()
	Run(&AScrollBar{}, t)
}
