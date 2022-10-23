// Copyright (c) 2022 Stephan Lukits. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package lines

import (
	"fmt"
	"testing"

	. "github.com/slukits/gounit"
	"github.com/slukits/lines/internal/api"
)

type CmpLine struct{ Suite }

func (s *CmpLine) tt(t *T, c Componenter) *Fixture {
	return TermFixture(t.GoT(), 0, c)
}

func (s *CmpLine) Is_filled_at_line_fillers(t *T) {
	fx := &icmpFX{init: func(c *icmpFX, e *Env) {
		c.Dim().SetHeight(1).SetWidth(8)
		fmt.Fprintf(e, "a%sb", LineFiller)
	}}
	tt := s.tt(t, fx)

	t.Eq("a      b", tt.ScreenOf(fx).Trimmed().String())

	t.FatalOn(tt.Lines.Update(fx, nil, func(e *Env) {
		fmt.Fprintf(e, "a%sb%[1]sc", LineFiller)
	}))

	t.Eq("a   b  c", tt.ScreenOf(fx).Trimmed().String())

	t.FatalOn(tt.Lines.Update(fx, nil, func(e *Env) {
		fmt.Fprintf(e, "ab%scd%[1]sef%[1]sgh", LineFiller)
	}))

	t.Eq("ab cd ef", tt.ScreenOf(fx).Trimmed().String())
}

func (s *CmpLine) Uses_added_style_range_on_next_sync(t *T) {
	fxSR := SR{
		Range: Range{2, 5},
		Style: api.DefaultStyle.WithBG(Red),
	}
	fx := &icmpFX{init: func(c *icmpFX, e *Env) {
		c.Dim().SetHeight(1).SetWidth(8)
		fmt.Fprintf(e, "12345678")
		e.AddStyleRange(0, fxSR)
	}}
	tt := s.tt(t, fx)

	l0 := tt.CellsOf(fx)[0]
	for i := range l0 {
		if fxSR.contains(i) {
			t.True(l0.HasBG(i, Red))
			continue
		}
		t.Not.True(l0.HasBG(i, Red))
	}
}

func (s *CmpLine) Adjusts_styles_on_centered(t *T) {
	fxCnt := LineFiller + "ab" + LineFiller
	expRng := Range{2, 4}
	fx := &icmpFX{init: func(c *icmpFX, e *Env) {
		fmt.Fprint(e, fxCnt)
		e.AddStyleRange(0, SR{
			Range: Range{1, 3},
			Style: e.NewStyle().WithBG(Red),
		})
	}}
	tt := s.tt(t, fx)
	tt.FireResize(6, 1)

	gotCnt := tt.Cells()[0]
	t.Eq("  ab  ", gotCnt.String())
	for i := range gotCnt {
		if expRng.contains(i) {
			t.True(gotCnt.HasBG(i, Red))
			continue
		}
		t.Not.True(gotCnt.HasBG(i, Red))
	}
}

func fxSR(x1, x2 int, c Color) SR {
	return SR{
		Range: Range{x1, x2},
		Style: api.DefaultStyle.WithBG(c),
	}
}

func (s *CmpLine) Adjusts_styles_on_evenly_distributed_line_filler(t *T) {
	fxCnt := "a" + LineFiller + "bc" + LineFiller + "d"
	fxR, fxG, fxB := fxSR(1, 2, Red), fxSR(2, 4, Green),
		fxSR(5, 6, Blue)
	dflt := api.DefaultStyle
	expR, expG, expB := Range{1, 3}, Range{3, 5}, Range{7, 8}
	fx := &icmpFX{init: func(c *icmpFX, e *Env) {
		fmt.Fprint(e, fxCnt)
		e.AddStyleRange(0, fxR, fxG, fxB)
	}}
	tt := s.tt(t, fx)
	tt.FireResize(8, 1)

	gotCnt := tt.Cells()[0]
	t.Eq("a  bc  d", gotCnt.String())
	for i := range gotCnt {
		if expR.contains(i) {
			t.True(gotCnt.HasBG(i, Red))
			continue
		}
		if expG.contains(i) {
			t.True(gotCnt.HasBG(i, Green))
			continue
		}
		if expB.contains(i) {
			t.True(gotCnt.HasBG(i, Blue))
			continue
		}
		t.True(gotCnt.HasBG(i, dflt.BG()))
	}
}

func (s *CmpLine) Adjusts_styles_on_unevenly_distributed_line_filler(t *T) {
	fxCnt := "a" + LineFiller + "bc" + LineFiller + "d" + LineFiller + "ef"
	fxR, fxG, fxB, fxY := fxSR(1, 2, Red), fxSR(2, 4, Green),
		fxSR(5, 6, Blue), fxSR(7, 9, Yellow)
	expR, expG, expB, expY := Range{1, 3}, Range{3, 5}, Range{7, 8},
		Range{9, 11}
	dflt := api.DefaultStyle
	fx := &icmpFX{init: func(c *icmpFX, e *Env) {
		fmt.Fprint(e, fxCnt)
		e.AddStyleRange(0, fxR, fxG, fxB, fxY)
	}}
	tt := s.tt(t, fx)
	tt.FireResize(11, 1)

	gotCnt := tt.Cells()[0]
	t.Eq("a  bc  d ef", gotCnt.String())
	for i := range gotCnt {
		if expR.contains(i) {
			t.True(gotCnt.HasBG(i, Red))
			continue
		}
		if expG.contains(i) {
			t.True(gotCnt.HasBG(i, Green))
			continue
		}
		if expB.contains(i) {
			t.True(gotCnt.HasBG(i, Blue))
			continue
		}
		if expY.contains(i) {
			t.True(gotCnt.HasBG(i, Yellow))
			continue
		}
		t.True(gotCnt.HasBG(i, dflt.BG()))
	}
}

func (s *CmpLine) Adjusts_styles_on_tab_expansion(t *T) {
	fxCnt := "\t\tred"
	fxR := fxSR(2, 5, Red)
	expR, dflt := Range{8, 11}, api.DefaultStyle
	fx := &icmpFX{init: func(c *icmpFX, e *Env) {
		fmt.Fprint(e, fxCnt)
		e.AddStyleRange(0, fxR)
	}}
	tt := s.tt(t, fx)
	tt.FireResize(11, 1)

	gotCnt := tt.Cells()[0]
	t.Eq("        red", gotCnt.String())
	for i := range gotCnt {
		if expR.contains(i) {
			t.True(gotCnt.HasBG(i, Red))
			continue
		}
		t.True(gotCnt.HasBG(i, dflt.BG()))
	}
}

func (s *CmpLine) Adjusts_styles_on_tab_and_line_filler_expansion(t *T) {
	fxCnt := "\tred" + LineFiller + "g" + LineFiller + "b"
	fxR, fxG, fxB := fxSR(1, 4, Red), fxSR(5, 6, Green),
		fxSR(7, 8, Blue)
	expR, expG, expB := Range{4, 7}, Range{9, 10}, Range{11, 12}
	dflt := api.DefaultStyle
	fx := &icmpFX{init: func(c *icmpFX, e *Env) {
		fmt.Fprint(e, fxCnt)
		e.AddStyleRange(0, fxR, fxG, fxB)
	}}
	tt := s.tt(t, fx)
	tt.FireResize(12, 1)

	gotCnt := tt.Cells()[0]
	t.Eq("    red  g b", gotCnt.String())
	for i := range gotCnt {
		if expR.contains(i) {
			t.True(gotCnt.HasBG(i, Red))
			continue
		}
		if expG.contains(i) {
			t.True(gotCnt.HasBG(i, Green))
			continue
		}
		if expB.contains(i) {
			t.True(gotCnt.HasBG(i, Blue))
			continue
		}
		t.True(gotCnt.HasBG(i, dflt.BG()))
	}
}

func TestLineRun(t *testing.T) {
	Run(&CmpLine{}, t)
}
