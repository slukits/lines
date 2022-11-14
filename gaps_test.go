// Copyright (c) 2022 Stephan Lukits. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package lines

import (
	"fmt"
	"strings"
	"testing"

	. "github.com/slukits/gounit"
	"github.com/slukits/lines/internal/api"
)

type _gaps struct{ Suite }

func (s *_gaps) fx(t *T, f ...func(*icmpFX, *Env)) (*Fixture, *icmpFX) {
	cmp := &icmpFX{}
	if len(f) > 0 {
		cmp.init = f[0]
	}
	return TermFixture(t.GoT(), 0, cmp), cmp
}

func (s *_gaps) Access_panics_outside_listener_callback(t *T) {
	_, cmp := s.fx(t)
	t.Panics(func() { cmp.Gaps(0) })
}

func (s *_gaps) Print_leaves_component_dirty(t *T) {
	tt, cmp := s.fx(t)
	exp := "written to top gap"

	tt.Lines.Update(tt.Root(), nil, func(e *Env) {
		fmt.Fprint(cmp.Gaps(0).Top, exp)
		t.True(cmp.IsDirty())
	})
}

func (s *_gaps) Synchronization_makes_component_clean(t *T) {
	tt, cmp := s.fx(t)
	exp := "written to top gap"

	tt.Lines.Update(tt.Root(), nil, func(e *Env) {
		fmt.Fprint(cmp.Gaps(0).Top, exp)
		t.True(cmp.IsDirty())
	})
	// we are promised everything is synced
	tt.Lines.Update(tt.Root(), nil, func(_ *Env) {
		t.Not.True(cmp.IsDirty())
	})
}

func (s *_gaps) Prints_to_top_gap_of_given_level(t *T) {
	tt, cmp := s.fx(t, func(c *icmpFX, e *Env) {
		c.Dim().SetHeight(3)
	})
	exp := "written to top gap"

	tt.Lines.Update(tt.Root(), nil, func(e *Env) {
		fmt.Fprint(cmp.Gaps(0).Top, exp)
	})
	t.Contains(tt.ScreenOf(cmp)[0], exp)
}

func (s *_gaps) Prints_to_bottom_gap_of_given_level(t *T) {
	tt, cmp := s.fx(t, func(c *icmpFX, e *Env) {
		c.Dim().SetHeight(3)
	})
	exp := "written to bottom gap"

	tt.Lines.Update(tt.Root(), nil, func(e *Env) {
		fmt.Fprint(cmp.Gaps(0).Bottom, exp)
	})
	t.Contains(tt.ScreenOf(cmp)[2], exp)
}

func (s *_gaps) Prints_to_left_gap_of_given_level(t *T) {
	tt, cmp := s.fx(t, func(c *icmpFX, e *Env) {
		c.Dim().SetWidth(3)
	})
	exp := "written to left gap"

	tt.Lines.Update(tt.Root(), nil, func(e *Env) {
		fmt.Fprint(cmp.Gaps(0).Left, exp)
	})
	t.Contains(tt.ScreenOf(cmp).Column(0), exp)
}

func (s *_gaps) Prints_to_right_gap_of_given_level(t *T) {
	tt, cmp := s.fx(t, func(c *icmpFX, e *Env) {
		c.Dim().SetWidth(3)
	})
	exp := "written to right gap"

	tt.Lines.Update(tt.Root(), nil, func(e *Env) {
		fmt.Fprint(cmp.Gaps(0).Right, exp)
	})
	t.Contains(tt.ScreenOf(cmp).Column(2), exp)
}

func (s *_gaps) Prints_to_horizontal_gaps_of_given_level(t *T) {
	tt, cmp := s.fx(t, func(c *icmpFX, e *Env) {
		c.Dim().SetHeight(3)
	})
	exp := "written to top/bottom gap"

	tt.Lines.Update(tt.Root(), nil, func(e *Env) {
		fmt.Fprint(cmp.Gaps(0).Horizontal, exp)
	})
	t.Contains(tt.ScreenOf(cmp)[0], exp)
	t.Contains(tt.ScreenOf(cmp)[2], exp)
}

func (s *_gaps) Prints_to_vertical_gaps_of_given_level(t *T) {
	tt, cmp := s.fx(t, func(c *icmpFX, e *Env) {
		c.Dim().SetWidth(3)
	})
	exp := "to left/right gap"

	tt.Lines.Update(tt.Root(), nil, func(e *Env) {
		fmt.Fprint(cmp.Gaps(0).Vertical, exp)
	})
	t.Contains(tt.ScreenOf(cmp).Column(0), exp)
	t.Contains(tt.ScreenOf(cmp).Column(2), exp)
}

func (s *_gaps) Prints_to_all_gaps_of_given_level(t *T) {
	tt, cmp := s.fx(t, func(c *icmpFX, e *Env) {
		c.Dim().SetWidth(3).SetHeight(3)
	})
	exp := " • \n• •\n • "

	tt.Lines.Update(tt.Root(), nil, func(e *Env) {
		Print(cmp.Gaps(0).Filling(), '•')
	})
	t.Eq(exp, tt.ScreenOf(cmp))
}

func (s *_gaps) Prints_to_corners(t *T) {
	tt, cmp := s.fx(t, func(c *icmpFX, e *Env) {
		c.Dim().SetWidth(3).SetHeight(3)
	})
	exp := "+ x\n   \ny z"

	tt.Lines.Update(tt.Root(), nil, func(e *Env) {
		fmt.Fprint(cmp.Gaps(0).Corners, "+", "x", "z", "y")
	})
	t.Contains(tt.ScreenOf(cmp), exp)

	exp = "z y\n   \nx +"
	tt.Lines.Update(tt.Root(), nil, func(e *Env) {
		cmp.SetDirty()
		fmt.Fprint(cmp.Gaps(0).TopLeft, "z")
		fmt.Fprint(cmp.Gaps(0).TopRight, "y")
		fmt.Fprint(cmp.Gaps(0).BottomRight, "+")
		fmt.Fprint(cmp.Gaps(0).BottomLeft, "x")
	})
	t.Contains(tt.ScreenOf(cmp), exp)
}

func (s *_gaps) Prints_filling(t *T) {
	tt, cmp := s.fx(t, func(c *icmpFX, e *Env) {
		c.Dim().SetWidth(6).SetHeight(4)
	})
	exp := " ---- \n|    |\n|    |\n ---- "

	tt.Lines.Update(tt.Root(), nil, func(e *Env) {
		Print(cmp.Gaps(0).Horizontal.At(0).Filling(), '-')
		Print(cmp.Gaps(0).Vertical.At(0).Filling(), '|')
	})

	t.Contains(tt.ScreenOf(cmp), exp)

	exp = " tttt \nl    r\nl    r\n bbbb "
	tt.Lines.Update(tt.Root(), nil, func(e *Env) {
		cmp.SetDirty()
		Print(cmp.Gaps(0).Top.At(0).Filling(), 't')
		Print(cmp.Gaps(0).Right.At(0).Filling(), 'r')
		Print(cmp.Gaps(0).Bottom.At(0).Filling(), 'b')
		Print(cmp.Gaps(0).Left.At(0).Filling(), 'l')
	})

	t.Contains(tt.ScreenOf(cmp), exp)
}

func (s *_gaps) Prints_reset_gap_filling(t *T) {
	tt, cmp := s.fx(t, func(c *icmpFX, e *Env) {
		c.Dim().SetWidth(4).SetHeight(1)
	})
	exp := " -- "

	tt.Lines.Update(tt.Root(), nil, func(e *Env) {
		Print(cmp.Gaps(0).Top.At(0).Filling(), '-')
	})

	t.Eq(exp, tt.ScreenOf(cmp))

	exp = " -  "
	tt.Lines.Update(tt.Root(), nil, func(e *Env) {
		fmt.Fprint(cmp.Gaps(0).Top, "-")
	})
	t.Eq(exp, tt.ScreenOf(cmp))
}

var framed = `
*xxxxxxxx*
x+------+x
x|      |x
x| size |x
x| 10x8 |x
x|      |x
x+------+x
*xxxxxxxx*
`

func (s *_gaps) Frame_its_component_s_content(t *T) {
	tt, cmp := s.fx(t, func(c *icmpFX, e *Env) {
		c.Dim().SetWidth(10).SetHeight(8)
	})
	exp := strings.TrimSpace(framed)

	tt.Lines.Update(cmp, nil, func(e *Env) {
		Print(cmp.Gaps(0).Horizontal.At(0).Filling(), 'x')
		Print(cmp.Gaps(0).Vertical.At(0).Filling(), 'x')
		fmt.Fprint(cmp.Gaps(0).Corners, "*")
		Print(cmp.Gaps(1).Horizontal.At(0).Filling(), '-')
		Print(cmp.Gaps(1).Vertical.At(0).Filling(), '|')
		fmt.Fprint(cmp.Gaps(1).Corners, "+")
		Print(cmp.Gaps(2).Horizontal.At(0).Filling(), ' ')
		Print(cmp.Gaps(2).Vertical.At(0).Filling(), ' ')
		fmt.Fprint(cmp.Gaps(2).Corners, " ")
		fmt.Fprint(e, "size\n10x8")
	})

	t.Eq(exp, tt.ScreenOf(cmp))
}

var positioned = `
+    top   +
            
            
           r
l          i
e          g
f          h
t          t
            
            
+  bottom  +
`

func (s *_gaps) Print_at_given_position(t *T) {
	tt, cmp := s.fx(t, func(c *icmpFX, e *Env) {
		c.Dim().SetWidth(12).SetHeight(11)
	})
	exp := strings.TrimSpace(positioned)

	tt.Lines.Update(cmp, nil, func(e *Env) {
		Print(cmp.Gaps(0).Top.At(4), []rune("top"))
		Print(cmp.Gaps(0).Right.At(2), []rune("right"))
		Print(cmp.Gaps(0).Bottom.At(2), []rune("bottom"))
		Print(cmp.Gaps(0).Left.At(3), []rune("left"))
		fmt.Fprint(cmp.Gaps(0).Corners, "+")
	})
	t.Eq(exp, tt.ScreenOf(cmp))
}

func (s *_gaps) At_prints_reset_gap_filling(t *T) {
	tt, cmp := s.fx(t, func(c *icmpFX, e *Env) {
		c.Dim().SetWidth(4).SetHeight(1)
	})
	exp := " -- "

	tt.Lines.Update(tt.Root(), nil, func(e *Env) {
		Print(cmp.Gaps(0).Top.At(0).Filling(), '-')
	})

	t.Eq(exp, tt.ScreenOf(cmp))

	exp = " -  "
	tt.Lines.Update(tt.Root(), nil, func(e *Env) {
		Print(cmp.Gaps(0).Top.At(0), '-')
	})
	t.Eq(exp, tt.ScreenOf(cmp))
}

func (s *_gaps) Filles_from_given_position_on(t *T) {
	tt, cmp := s.fx(t, func(c *icmpFX, e *Env) {
		c.Dim().SetWidth(9).SetHeight(1)
	})
	exp, fx := " --top-- ", []rune("top")

	tt.Lines.Update(cmp, nil, func(e *Env) {
		Print(cmp.Gaps(0).Top.At(0).Filling(), '-')
		Print(cmp.Gaps(0).Top.At(1), fx)
		Print(cmp.Gaps(0).Top.At(len(fx)+1).Filling(), '-')
	})

	t.Eq(exp, tt.ScreenOf(cmp))
}

func (s *_gaps) Style_defaults_to_component_s_style(t *T) {
	tt, cmp := s.fx(t, func(c *icmpFX, e *Env) {
		c.Dim().SetWidth(3).SetHeight(3)
	})
	tt.Lines.Update(cmp, nil, func(e *Env) {
		cmp.AA(Reverse).FG(Blue).BG(Yellow)
		Print(cmp.Gaps(0).Filling(), ' ')
		fmt.Fprint(cmp.Gaps(0).Corners, " ")
	})

	testStyle := func(c api.TestCell) {
		t.True(c.Style.AA() == Reverse && c.Style.FG() == Blue &&
			c.Style.BG() == Yellow)
	}

	cc := tt.CellsOf(cmp)
	for _, c := range cc[0] {
		testStyle(c)
	}
	for _, c := range cc[2] {
		testStyle(c)
	}
	testStyle(cc[1][0])
	testStyle(cc[1][2])
}

func (s *_gaps) Have_set_style(t *T) {
	tt, cmp := s.fx(t, func(c *icmpFX, e *Env) {
		c.Dim().SetWidth(3).SetHeight(3)
	})
	tt.Lines.Update(cmp, nil, func(e *Env) {
		cmp.Gaps(0).AA(Reverse).FG(Blue).BG(Yellow)
	})

	testStyle := func(c api.TestCell) {
		t.True(c.Style.AA() == Reverse && c.Style.FG() == Blue &&
			c.Style.BG() == Yellow)
	}

	cc := tt.CellsOf(cmp)
	for _, c := range cc[0] {
		testStyle(c)
	}
	for _, c := range cc[2] {
		testStyle(c)
	}
	testStyle(cc[1][0])
	testStyle(cc[1][2])
}

func (s *_gaps) Have_set_vertical_style(t *T) {
	tt, cmp := s.fx(t, func(c *icmpFX, e *Env) {
		c.Dim().SetWidth(3).SetHeight(3)
	})
	tt.Lines.Update(cmp, nil, func(e *Env) {
		cmp.Gaps(0).Vertical.AA(Reverse).FG(Blue).BG(Yellow)
	})

	testStyle := func(c api.TestCell) {
		t.True(c.Style.AA() == Reverse && c.Style.FG() == Blue &&
			c.Style.BG() == Yellow)
	}

	cc := tt.CellsOf(cmp)
	testStyle(cc[1][0])
	testStyle(cc[1][2])
}

func (s *_gaps) Have_set_horizontal_style(t *T) {
	tt, cmp := s.fx(t, func(c *icmpFX, e *Env) {
		c.Dim().SetWidth(3).SetHeight(3)
	})
	tt.Lines.Update(cmp, nil, func(e *Env) {
		cmp.Gaps(0).Horizontal.AA(Reverse).FG(Blue).BG(Yellow)
	})

	testStyle := func(c api.TestCell) {
		t.True(c.Style.AA() == Reverse && c.Style.FG() == Blue &&
			c.Style.BG() == Yellow)
	}

	cc := tt.CellsOf(cmp)
	testStyle(cc[0][1])
	testStyle(cc[2][1])
}

func (s *_gaps) Have_set_at_writer_style(t *T) {
	tt, cmp := s.fx(t, func(c *icmpFX, e *Env) {
		c.Dim().SetWidth(9).SetHeight(1)
	})
	tt.Lines.Update(cmp, nil, func(e *Env) {
		Print(cmp.Gaps(0).Top.At(0).Filling().
			AA(Bold).FG(Yellow).BG(Blue), '-')
		Print(cmp.Gaps(0).Top.At(1).AA(Blink).FG(Blue).BG(Yellow),
			[]rune("top"))
		Print(cmp.Gaps(0).Top.At(4).Filling().
			AA(Bold).FG(Yellow).BG(Blue), '-')
	})

	t.Eq(" --top-- ", tt.CellsOf(cmp)[0])
	frame := (Style{}).WithAA(Bold).WithFG(Yellow).WithBG(Blue)
	title := NewStyle(Blink, Blue, Yellow)
	for _, c := range tt.CellsOf(cmp)[0] {
		switch c.Rune {
		case '-':
			c.Style.Equals(frame)
		case ' ':
		default:
			c.Style.Equals(title)
		}
	}
}

var lytInit = `
•••top•••
•       •
•       •
•••••••••
`

var lytChange = `
••••top••••
•         •
•         •
•         •
•         •
•••••••••••
`

func (s *_gaps) Adapt_fillers_on_layout_change(t *T) {
	tt, cmp := s.fx(t, func(c *icmpFX, e *Env) {
		c.Dim().SetWidth(9).SetHeight(4)
	})
	tt.Lines.Update(cmp, nil, func(e *Env) {
		fmt.Fprint(cmp.Gaps(0).Corners.AA(Reverse), "•")
		Print(cmp.Gaps(0).Top.At(0).AA(Bold).Filling(), '•')
		Print(cmp.Gaps(0).Top.At(1).AA(Blink), []rune("top"))
		Print(cmp.Gaps(0).Top.At(4).AA(Bold).Filling(), '•')
		Print(cmp.Gaps(0).Vertical.AA(Dim).At(0).Filling(), '•')
		Print(cmp.Gaps(0).Bottom.AA(Dim).At(0).Filling(), '•')
	})

	t.Eq(strings.TrimSpace(lytInit), tt.ScreenOf(cmp))

	tt.Lines.Update(cmp, nil, func(e *Env) {
		cmp.Dim().SetWidth(11).SetHeight(6)
	})

	t.Eq(strings.TrimSpace(lytChange), tt.ScreenOf(cmp))
	cc := tt.CellsOf(cmp)
	t.True(cc[0].HasAA(0, Reverse))
	t.True(cc[0].HasAA(10, Reverse))
	t.True(cc[0].HasAA(3, Bold))
	t.True(cc[0].HasAA(4, Blink))
	t.True(cc[0].HasAA(6, Blink))
	t.True(cc[0].HasAA(7, Bold))
	t.True(cc[4].HasAA(0, Dim))
	t.True(cc[5].HasAA(1, Dim))
	t.True(cc[5].HasAA(0, Reverse))
}

var fillAll = `
••••••••
•      •
•      •
••••••••
`

func (s *_gaps) Filler_fills_whole_level(t *T) {
	tt, cmp := s.fx(t, func(c *icmpFX, e *Env) {
		c.Dim().SetWidth(8).SetHeight(4)
	})
	tt.Lines.Update(cmp, nil, func(e *Env) {
		Print(cmp.Gaps(0).Filling(), '•')
		fmt.Fprint(cmp.Gaps(0).Corners, "•")
	})
	t.Eq(strings.TrimSpace(fillAll), tt.ScreenOf(cmp))
}

var nested = `
•stacker•
••inner••
••     ••
••     ••
•••••••••
•••••••••
`

type gapCmp struct {
	Component
	top string
}

func (c *gapCmp) OnInit(e *Env) {
	Print(c.Gaps(0).Filling(), '•')
	fmt.Fprint(c.Gaps(0).Top, c.top)
	fmt.Fprint(c.Gaps(0).Corners, "•")
}

type gapStackerCmp struct {
	gapCmp
	Stacking
}

func (c *gapStackerCmp) OnInit(e *Env) {
	c.CC = append(c.CC, &gapCmp{top: "inner"})
	c.top = "stacker"
	c.gapCmp.OnInit(e)
}

func (s *_gaps) Of_nested_components_are_displayed(t *T) {
	tt := TermFixture(t.GoT(), 0, &gapStackerCmp{})
	tt.FireResize(9, 6)
	t.Eq(strings.TrimSpace(nested), tt.Screen())
}

func TestGaps(t *testing.T) {
	t.Parallel()
	Run(&_gaps{}, t)
}
