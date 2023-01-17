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

type env struct{ Suite }

func (s *env) SetUp(t *T) { t.Parallel() }

func (s *env) Access_panics_outside_event_processing(t *T) {
	outsideEnv := (*Env)(nil)
	fx(t, &cmpFX{
		onInit: func(cf *cmpFX, e *Env) { outsideEnv = e },
	})
	t.Panics(func() { fmt.Fprint(outsideEnv, "panics") })
}

func (s *env) Provides_the_display_size(t *T) {
	fx := fx(t, &cmpFX{})
	width, height := fx.Size()
	fx.Lines.Update(fx.Root(), nil, func(e *Env) {
		envWidth, envHeight := e.ScreenSize()
		t.Eq(width, envWidth)
		t.Eq(height, envHeight)
	})
}

func (s *env) Prints_to_component_starting_at_top_left_corner(t *T) {
	stacking := &stackingFX{}
	stacking.CC = append(stacking.CC,
		&cmpFX{onInit: func(_ *cmpFX, e *Env) {
			fmt.Fprint(e, "1st")
		}},
		&cmpFX{onInit: func(_ *cmpFX, e *Env) {
			fmt.Fprint(e, "3rd")
		}},
	)
	fx := fx(t, stacking)
	fx.FireResize(4, 4)
	t.Eq("1st \n    \n3rd \n    ", fx.Screen())
}

func (s *env) Print_breaks_string_at_line_breaks_to_screen_lines(t *T) {
	fx := fx(t, &cmpFX{onInit: func(_ *cmpFX, e *Env) {
		fmt.Fprint(e, "1st\n2nd\n3rd\n4th\n5th")
	}})
	fx.FireResize(3, 4)
	t.Eq("1st\n2nd\n3rd\n4th", fx.Screen())
}

func (s *env) Defaults_printed_lines_colors_to_component_globals(t *T) {
	var exp Style
	fx := fx(t, &cmpFX{onInit: func(c *cmpFX, e *Env) {
		fmt.Fprint(e, "1st\n2nd\n3rd")
		exp = c.Globals().Style(Default)
	}})
	fx.FireResize(3, 3)
	for _, l := range fx.Cells() {
		for _, c := range l {
			c.Style.Equals(exp)
		}
	}
}

func (s *env) Defaults_printed_lines_styles_to_component_globals(t *T) {
	var exp Style
	fx := fx(t, &cmpFX{onInit: func(c *cmpFX, e *Env) {
		fmt.Fprint(e, "1st\n2nd\n3rd")
		exp = c.Globals().Style(Default)
	}})
	fx.FireResize(3, 3)
	for _, l := range fx.Cells() {
		for _, c := range l {
			c.Style.Equals(exp)
		}
	}
}

func (s *env) Sets_colors_for_printed_screen_lines(t *T) {
	var exp, ggSty Style
	fx := fx(t, &cmpFX{onInit: func(c *cmpFX, e *Env) {
		fmt.Fprint(e.FG(Green).BG(Red), "1st\n2nd")
		exp = c.Globals().Style(Default).WithFG(Green).WithBG(Red)
		ggSty = c.Globals().Style(Default)
	}})
	fx.FireResize(4, 3)
	for i, l := range fx.Cells() {
		sty := exp
		if i == 2 {
			sty = ggSty
		}
		for _, c := range l {
			t.True(c.Style.Equals(sty))
		}
	}
}

func (s *env) Sets_styles_for_printed_screen_lines(t *T) {
	var exp, ggSty Style
	fx := fx(t, &cmpFX{onInit: func(c *cmpFX, e *Env) {
		fmt.Fprint(e.AA(Blink), "1st\n2nd")
		exp = c.Globals().Style(Default).WithAA(Blink)
		ggSty = c.Globals().Style(Default)
	}})
	fx.FireResize(3, 3)
	for i, l := range fx.Cells() {
		sty := exp
		if i == 2 {
			sty = ggSty
		}
		for _, c := range l {
			t.True(c.Style.Equals(sty))
		}
	}
}

func (s *env) Provides_writer_for_the_nth_line(t *T) {
	fx := fx(t, &cmpFX{onInit: func(_ *cmpFX, e *Env) {
		fmt.Fprint(e.LL(0), "first line")
		fmt.Fprint(e.LL(7), "eighth line")
	}})

	sl := strings.Split(fx.Screen().Trimmed().String(), "\n")
	t.FatalIfNot(t.Eq(8, len(sl)))
	t.Eq("first line", strings.TrimSpace(sl[0]))
	t.Eq("eighth line", strings.TrimSpace(sl[7]))
}

func (s *env) Overwrites_given_line_and_following(t *T) {
	cmp := &cmpFX{onInit: func(_ *cmpFX, e *Env) {
		fmt.Fprint(e.LL(0), "first line")
		fmt.Fprint(e.LL(7), "eighth line")
	}}
	fx := fx(t, cmp)
	fx.Lines.Update(cmp, nil, func(e *Env) {
		fmt.Fprint(e.LL(6), "seventh line\n"+
			"short 8th\nninth line")
	})

	sl := fx.Screen().Trimmed()
	t.FatalIfNot(t.Eq(9, len(sl)))
	t.Eq(strings.TrimSpace(sl[0]), "first line")
	t.Eq(strings.TrimSpace(sl[6]), "seventh line")
	t.Eq(strings.TrimSpace(sl[7]), "short 8th")
	t.Eq(strings.TrimSpace(sl[8]), "ninth line")
}

func (s *env) Fills_line_with_blanks_at_line_filler(t *T) {
	cmp := &cmpFX{onInit: func(_ *cmpFX, e *Env) {
		fmt.Fprint(e,
			"first line"+Filler+"filled\nsecond line"+Filler+"filled")
	}}
	fx := fx(t, cmp)
	fx.FireResize(20, 2)

	expFst := "first line    filled"
	t.Eq(expFst, fx.Screen()[0])
	expSnd := "second line   filled"
	t.Eq(expSnd, fx.Screen()[1])

	fx.Lines.Update(cmp, nil, func(e *Env) {
		cmp.Reset(All)
		fmt.Fprint(e.LL(0),
			"first line"+Filler+"filled\n"+
				"second line"+Filler+"filled",
		)
	})

	t.Eq(expFst, fx.Screen()[0])
	t.Eq(expSnd, fx.Screen()[1])
}

func (s *env) Changes_fore_and_background_for_partial_line(t *T) {
	fx := fx(t, &cmpFX{onInit: func(_ *cmpFX, e *Env) {
		Print(
			e.LL(0).At(0).BG(Red).FG(White),
			[]rune("un-filled with space"),
		)
		fmt.Fprint(e.LL(1).BG(Red).FG(White), "filled with space")
	}})

	l0 := fx.Cells()[0]
	for i := range l0 {
		if l0[i].Rune != ' ' || len(l0) > i+1 && l0[i+1].Rune != ' ' {
			t.True(l0.HasBG(i, Red))
			t.True(l0.HasFG(i, White))
			continue
		}
		t.Not.True(l0.HasBG(i, Red))
		t.Not.True(l0.HasFG(i, White))
	}

	l1 := fx.Cells()[1]
	for i := range l1 {
		t.True(l1.HasBG(i, Red))
		t.True(l1.HasFG(i, White))
	}
}

func (s *env) Changes_line_style_for_a_range_of_runes(t *T) {
	fx := fx(t, &cmpFX{onInit: func(_ *cmpFX, e *Env) {
		fmt.Fprint(e, "\t")
		Print(e.LL(0).At(1).FG(White).BG(Red), []rune("red"))
		Print(e.LL(0).At(1+len("red")), []rune(Filler+"right"))
	}})

	l0, exp := fx.Cells()[0], Range{4, 7}
	str := l0.String()
	for i := range str {
		if exp.contains(i) {
			t.True(l0.HasBG(i, Red))
			continue
		}
		t.Not.True(l0.HasBG(i, Red))
	}
}

func (s *env) Prints_filling_to_component_line(t *T) {
	red, blue, green := []rune("red"), []rune("blue"), []rune("green")
	fx := fx(t, &cmpFX{onInit: func(_ *cmpFX, e *Env) {
		fmt.Fprint(e, "\t\t")
		Print(e.LL(0).At(2).FG(White).BG(Red), red)
		Print(e.LL(0).At(2+len(red)).Filling(), '_')
		Print(e.LL(0).At(2+len(red)+1).FG(Yellow).BG(Blue), blue)
		Print(e.LL(0).At(2+len(red)+1+len(blue)).BG(Brown).FG(Salmon).
			Filling(), '_')
		Print(e.LL(0).At(2+len(red)+1+len(blue)+1).FG(Black).BG(Green),
			green)
	}})
	fx.PostResize(40, 1)

	exp := "        red__________blue__________green"
	t.Eq(exp, fx.Screen())
	l := fx.Cells()[0]
	for i, c := range l {
		switch i {
		case 0, 1, 2, 3, 4, 5, 6, 7:
			t.True(c.Style.Equals(fx.Lines.Globals.Style(Default)))
		case 8, 9, 10:
			t.True(c.Style.Equals(fx.Lines.Globals.Style(Default).
				WithFG(White).WithBG(Red)))
		case 11, 12, 13, 14, 15, 16, 17, 18, 19, 20:
			t.True(c.Style.Equals(fx.Lines.Globals.Style(Default)))
		case 21, 22, 23, 24:
			t.True(c.Style.Equals(fx.Lines.Globals.Style(Default).
				WithFG(Yellow).WithBG(Blue)))
		case 25, 26, 27, 28, 29, 30, 31, 32, 33, 34:
			t.True(c.Style.Equals(fx.Lines.Globals.Style(Default).
				WithFG(Salmon).WithBG(Brown)))
		case 35, 36, 37, 38, 39:
			t.True(c.Style.Equals(fx.Lines.Globals.Style(Default).
				WithFG(Black).WithBG(Green)))
		}
	}
}

func TestEnv(t *testing.T) {
	t.Parallel()
	Run(&env{}, t)
}
