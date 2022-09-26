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

type Line struct{ Suite }

func (s *Line) Is_filled_at_line_fillers(t *T) {
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

func (s *line) Uses_added_style_range_on_next_sync(t *T) {
	var r Range
	fxSR := SR{Range: *r.SetStart(2).SetEnd(5),
		Style: tcell.StyleDefault.Background(tcell.ColorRed)}
	fx := &icmpFX{init: func(c *icmpFX, e *Env) {
		c.Dim().SetHeight(1).SetWidth(8)
		fmt.Fprintf(e, "12345678")
		e.AddStyleRange(0, fxSR)
	}}
	ee, tt := Test(t.GoT(), fx, 1)
	ee.Listen()

	got := tt.LastScreen[0].Styles()
	for i := range "12345678" {
		if fxSR.Contains(i) {
			t.True(got.Of(i).HasBG(tcell.ColorRed))
			continue
		}
		t.Not.True(got.Of(i).HasBG(tcell.ColorRed))
	}
}

func (s *Line) Adjusts_styles_on_centered(t *T) {
	fxCnt := LineFiller + "ab" + LineFiller
	fxSR, expRng := SR{Range: Range{1, 3}}, Range{2, 4}
	fxSR.Style = tcell.StyleDefault.Background(tcell.ColorRed)
	fx := &icmpFX{init: func(c *icmpFX, e *Env) {
		fmt.Fprint(e, fxCnt)
		e.AddStyleRange(0, fxSR)
	}}
	ee, tt := Test(t.GoT(), fx)
	tt.FireResize(6, 1)
	defer ee.QuitListening()

	gotCnt := tt.FullScreen()[0]
	t.Eq("  ab  ", gotCnt.String())
	got := gotCnt.Styles()
	for i := range gotCnt.String() {
		if expRng.Contains(i) {
			t.True(got.Of(i).HasBG(tcell.ColorRed))
			continue
		}
		t.Not.True(got.Of(i).HasBG(tcell.ColorRed))
	}
}

func (s *Line) fxSR(x1, x2 int, c tcell.Color) SR {
	return SR{
		Range: Range{x1, x2},
		Style: tcell.StyleDefault.Background(c),
	}
}

var (
	red    = tcell.ColorRed
	green  = tcell.ColorGreen
	blue   = tcell.ColorBlue
	yellow = tcell.ColorYellow
)

func (s *Line) Adjusts_styles_on_evenly_distributed_line_filler(t *T) {
	fxCnt := "a" + LineFiller + "bc" + LineFiller + "d"
	fxR, fxG, fxB := s.fxSR(1, 2, red), s.fxSR(2, 4, green),
		s.fxSR(5, 6, blue)
	expR, expG, expB := Range{1, 3}, Range{3, 5}, Range{7, 8}
	fx := &icmpFX{init: func(c *icmpFX, e *Env) {
		fmt.Fprint(e, fxCnt)
		e.AddStyleRange(0, fxR, fxG, fxB)
	}}
	ee, tt := Test(t.GoT(), fx)
	tt.FireResize(8, 1)
	defer ee.QuitListening()

	gotCnt := tt.FullScreen()[0]
	t.Eq("a  bc  d", gotCnt.String())
	got := gotCnt.Styles()
	for i := range gotCnt.String() {
		if expR.Contains(i) {
			t.True(got.Of(i).HasBG(red))
			continue
		}
		if expG.Contains(i) {
			t.True(got.Of(i).HasBG(green))
			continue
		}
		if expB.Contains(i) {
			t.True(got.Of(i).HasBG(blue))
			continue
		}
		t.True(got.Of(i).HasBG(tcell.ColorDefault))
	}
}

func (s *Line) Adjusts_styles_on_unevenly_distributed_line_filler(t *T) {
	fxCnt := "a" + LineFiller + "bc" + LineFiller + "d" + LineFiller + "ef"
	fxR, fxG, fxB, fxY := s.fxSR(1, 2, red), s.fxSR(2, 4, green),
		s.fxSR(5, 6, blue), s.fxSR(7, 9, yellow)
	expR, expG, expB, expY := Range{1, 3}, Range{3, 5}, Range{7, 8},
		Range{9, 11}
	fx := &icmpFX{init: func(c *icmpFX, e *Env) {
		fmt.Fprint(e, fxCnt)
		e.AddStyleRange(0, fxR, fxG, fxB, fxY)
	}}
	ee, tt := Test(t.GoT(), fx)
	tt.FireResize(11, 1)
	defer ee.QuitListening()

	gotCnt := tt.FullScreen()[0]
	t.Eq("a  bc  d ef", gotCnt.String())
	got := gotCnt.Styles()
	for i := range gotCnt.String() {
		if expR.Contains(i) {
			t.True(got.Of(i).HasBG(red))
			continue
		}
		if expG.Contains(i) {
			t.True(got.Of(i).HasBG(green))
			continue
		}
		if expB.Contains(i) {
			t.True(got.Of(i).HasBG(blue))
			continue
		}
		if expY.Contains(i) {
			t.True(got.Of(i).HasBG(yellow))
			continue
		}
		t.True(got.Of(i).HasBG(tcell.ColorDefault))
	}
}

func (s *Line) Adjusts_styles_on_tab_expansion(t *T) {
	fxCnt := "\t\tred"
	fxR := s.fxSR(2, 5, red)
	expR := Range{8, 11}
	fx := &icmpFX{init: func(c *icmpFX, e *Env) {
		fmt.Fprint(e, fxCnt)
		e.AddStyleRange(0, fxR)
	}}
	ee, tt := Test(t.GoT(), fx)
	tt.FireResize(11, 1)
	defer ee.QuitListening()

	gotCnt := tt.FullScreen()[0]
	t.Eq("        red", gotCnt.String())
	got := gotCnt.Styles()
	for i := range gotCnt.String() {
		if expR.Contains(i) {
			t.True(got.Of(i).HasBG(red))
			continue
		}
		t.True(got.Of(i).HasBG(tcell.ColorDefault))
	}
}

func (s *Line) Adjusts_styles_on_tab_and_line_filler_expansion(t *T) {
	fxCnt := "\tred" + LineFiller + "g" + LineFiller + "b"
	fxR, fxG, fxB := s.fxSR(1, 4, red), s.fxSR(5, 6, green),
		s.fxSR(7, 8, blue)
	expR, expG, expB := Range{4, 7}, Range{9, 10}, Range{11, 12}
	fx := &icmpFX{init: func(c *icmpFX, e *Env) {
		fmt.Fprint(e, fxCnt)
		e.AddStyleRange(0, fxR, fxG, fxB)
	}}
	ee, tt := Test(t.GoT(), fx)
	tt.FireResize(12, 1)
	defer ee.QuitListening()

	gotCnt := tt.FullScreen()[0]
	t.Eq("    red  g b", gotCnt.String())
	got := gotCnt.Styles()
	for i := range gotCnt.String() {
		if expR.Contains(i) {
			t.True(got.Of(i).HasBG(red))
			continue
		}
		if expG.Contains(i) {
			t.True(got.Of(i).HasBG(green))
			continue
		}
		if expB.Contains(i) {
			t.True(got.Of(i).HasBG(blue))
			continue
		}
		t.True(got.Of(i).HasBG(tcell.ColorDefault))
	}
}

func TestLineRun(t *testing.T) {
	Run(&Line{}, t)
}
