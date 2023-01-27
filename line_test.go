// Copyright (c) 2022 Stephan Lukits. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package lines

import (
	"testing"

	. "github.com/slukits/gounit"
	"github.com/slukits/lines/internal/api"
	"github.com/slukits/lines/internal/term"
)

type lineFX struct {
	api.Displayer
	Line
	w, h int
	gg   *Globals
}

func (x *lineFX) width() int {
	if x.w == 0 {
		x.w, _ = x.Size()
	}
	return x.w
}

func (x *lineFX) height() int {
	if x.h == 0 {
		_, x.h = x.Size()
	}
	return x.h
}

func (x *lineFX) redraw(tt *term.Fixture) api.CellsLine {
	rr, ss := x.display(x.width(), x.gg)
	for i, r := range rr {
		x.Display(i, 0, r, ss.of(i))
	}
	x.Redraw()
	return tt.CellsArea(0, 0, x.width(), 1)[0]
}

func (x *lineFX) highlighted(s Style) Style {
	if x.gg.Style(Highlight).AA() != 0 {
		if s.AA()&x.gg.Style(Highlight).AA() == 0 {
			s = s.WithAdded(x.gg.Style(Highlight).AA())
		} else {
			s = s.WithRemoved(x.gg.Style(Highlight).AA())
		}
	}
	if x.gg.Style(Highlight).FG() != DefaultColor {
		s = s.WithFG(x.gg.Style(Highlight).FG())
	}
	if x.gg.Style(Highlight).BG() != DefaultColor {
		s = s.WithBG(x.gg.Style(Highlight).BG())
	}
	return s
}

// NOTE since the point here is to determine what a line provides for
// the display it doesn't matter to what backend it goes as long as we
// can figure out what went to the display with what style attributes.
func newLineFX(t *T) (*term.Fixture, *lineFX) {
	ui, tt := term.LstFixture(t.GoT(), nil, 0)
	tt.PostResize(20, 1)
	return tt, &lineFX{Displayer: ui, gg: newGlobals(nil)}
}

func newLineFXDflt(t *T, s Style) (*term.Fixture, *lineFX) {
	tt, fx := newLineFX(t)
	fx.gg.SetStyle(Default, s)
	return tt, fx
}

type ALine struct{ Suite }

func (s *ALine) SetUp(t *T) { t.Parallel() }

func (s *ALine) Is_padded_with_spaces_if_zero(t *T) {
	tt, fx := newLineFX(t)
	scrLine := fx.redraw(tt)
	for _, c := range scrLine {
		t.True(c.Rune == ' ')
		t.Eq(c.Style, DefaultStyle)
	}
}

func (s *ALine) Uses_given_display_style_if_no_default_set(t *T) {
	exp := NewStyle(Blink, Yellow, Blue)
	tt, fx := newLineFXDflt(t, exp)
	scrLine := fx.redraw(tt)
	for _, c := range scrLine {
		t.True(c.Rune == ' ')
		t.Eq(c.Style, exp)
	}
}

func (s *ALine) Has_set_default_style_if_empty(t *T) {
	tt, fx := newLineFX(t)
	exp := NewStyle(Blink, Yellow, Blue)
	fx.setDefaultStyle(exp)
	scrLine := fx.redraw(tt)
	for _, c := range scrLine {
		t.True(c.Rune == ' ')
		t.Eq(c.Style, exp)
	}
}

func (s *ALine) Has_updated_default_style_attributes(t *T) {
	tt, fx := newLineFX(t)
	exp := DefaultStyle.WithAA(Dim)
	fx.withAA(Dim)
	scrLine := fx.redraw(tt)
	for _, c := range scrLine {
		t.True(c.Rune == ' ')
		t.Eq(c.Style, exp)
	}
}

func (s *ALine) Has_updated_default_foreground_color(t *T) {
	tt, fx := newLineFX(t)
	exp := DefaultStyle.WithFG(Green)
	fx.withFG(Green)
	scrLine := fx.redraw(tt)
	for _, c := range scrLine {
		t.True(c.Rune == ' ')
		t.Eq(c.Style, exp)
	}
}

func (s *ALine) Has_updated_default_background_color(t *T) {
	tt, fx := newLineFX(t)
	exp := DefaultStyle.WithBG(Red)
	fx.withBG(Red)
	scrLine := fx.redraw(tt)
	for _, c := range scrLine {
		t.True(c.Rune == ' ')
		t.Eq(c.Style, exp)
	}
}

func (s *ALine) Displays_set_content_space_padded_to_line_width(t *T) {
	_, fx := newLineFX(t)
	fx.set("0123456789")
	got, _ := fx.display(fx.width(), fx.gg)
	t.Eq("0123456789", string(got[:10]))
	t.Eq("          ", string(got[10:]))
}

func (s *ALine) Truncates_line_with_width_overflowing_content(t *T) {
	_, fx := newLineFX(t)
	fx.set("01234567890123456789012")
	got, _ := fx.display(fx.width(), fx.gg)
	t.Eq("01234567890123456789", string(got))
}

func (s *ALine) Displays_content_with_set_style(t *T) {
	exp := NewStyle(Blink, Yellow, Blue)
	tt, fx := newLineFX(t)
	fx.setStyled("0123456789", exp)
	scrLine := fx.redraw(tt)
	for _, c := range scrLine {
		switch c.Rune {
		case ' ':
			t.Eq(c.Style, fx.gg.Style(Default))
		default:
			t.Eq(c.Style, exp)
		}
	}
}

func (s *ALine) Has_content_set_at_zero_position(t *T) {
	_, fx := newLineFX(t)
	fx.setAt(0, []rune("0123456789"))
	got, _ := fx.display(fx.width(), fx.gg)
	t.Eq("0123456789", string(got[:10]))
	t.Eq("          ", string(got[10:]))
}

func (s *ALine) Has_content_set_at_given_position_space_padded(t *T) {
	_, fx := newLineFX(t)
	fx.setAt(8, []rune("0123456789"))
	got, _ := fx.display(fx.width(), fx.gg)
	t.Eq("        ", string(got[:8]))
	t.Eq("0123456789", string(got[8:18]))
	t.Eq("  ", string(got[18:]))
}

func (s *ALine) Styles_content_set_at_given_position(t *T) {
	exp := NewStyle(Blink, Yellow, Blue)
	tt, fx := newLineFX(t)
	fx.setStyledAt(5, []rune("0123456789"), exp)
	scrLine := fx.redraw(tt)
	for _, c := range scrLine {
		switch c.Rune {
		case ' ':
			t.Eq(c.Style, fx.gg.Style(Default))
		default:
			t.Eq(c.Style, exp)
		}
	}
}

func (s *ALine) Overwrites_content_after_given_position(t *T) {
	_, fx := newLineFX(t)
	fx.set("0123456789")
	got, _ := fx.display(10, fx.gg)
	t.Eq("0123456789", string(got))

	fx.setAt(2, []rune("42"))
	got, _ = fx.display(10, fx.gg)
	t.Eq("0142      ", string(got))
}

func (s *ALine) Fills_remaining_space_with_a_filling_rune(t *T) {
	_, fx := newLineFX(t)
	fx.setAtFilling(0, 'a')
	fx.setAt(1, []rune("0123456789"))
	got, _ := fx.display(fx.width(), fx.gg)
	t.Eq("aaaaaaaaaa0123456789", string(got))
	fx.set("01234")
	fx.setAtFilling(5, 'a')
	fx.setAt(6, []rune("56789"))
	got, _ = fx.display(fx.width(), fx.gg)
	t.Eq("01234aaaaaaaaaa56789", string(got))
	fx.setAt(0, []rune("0123456789")) // test filler truncation
	fx.setAtFilling(10, 'a')
	got, _ = fx.display(fx.width(), fx.gg)
	t.Eq("0123456789aaaaaaaaaa", string(got))
	fx.setAt(0, []rune("012"))
	fx.setAtFilling(3, 'a')
	fx.setAt(4, []rune("3456"))
	fx.setAtFilling(8, 'b')
	fx.setAt(9, []rune("789"))
	got, _ = fx.display(fx.width(), fx.gg)
	t.Eq("012aaaaa3456bbbbb789", string(got))
}

func (s *ALine) Expands_filler_style_preserving(t *T) {
	s1, s2 := NewStyle(Blink, Yellow, Blue), NewStyle(Dim, Green, Red)
	_, fx := newLineFX(t)
	fx.setStyledAtFilling(0, 'a', s1)
	fx.setAt(1, []rune("0123456789"))
	_, ss := fx.display(fx.width(), fx.gg)
	for r, s := range ss {
		switch r.End() {
		case 10:
			t.Eq(0, r.Start())
			t.Eq(s, s1)
		case 20:
			t.Eq(10, r.Start())
			t.Eq(s, fx.gg.Style(Default))
		}
	}
	fx.set("01234")
	fx.setStyledAtFilling(5, 'a', s1)
	fx.setStyledAt(6, []rune("56789"), s2)
	_, ss = fx.display(fx.width(), fx.gg)
	for r, s := range ss {
		switch r.End() {
		case 5:
			t.Eq(0, r.Start())
			t.Eq(s, fx.gg.Style(Default))
		case 15:
			t.Eq(5, r.Start())
			t.Eq(s, (s1))
		case 20:
			t.Eq(15, r.Start())
			t.Eq(s, s2)
		}
	}
	fx.setAt(0, []rune("0123456789"))
	fx.setStyledAtFilling(10, 'a', s1)
	_, ss = fx.display(fx.width(), fx.gg)
	for r, s := range ss {
		switch r.End() {
		case 10:
			t.Eq(0, r.Start())
			t.Eq(s, fx.gg.Style(Default))
		case 20:
			t.Eq(10, r.Start())
			t.Eq(s, s1)
		}
	}
}

func (s *ALine) Expands_leading_tabs_style_preserving(t *T) {
	s1, s2 := NewStyle(Blink, Yellow, Blue), NewStyle(Dim, Green, Red)
	tt, fx := newLineFX(t)
	fx.gg.tabWidth = 5
	fx.setStyledAt(0, []rune{'\t', '\t'}, s1)
	fx.setStyledAt(2, []rune("0123456789"), s2)

	l := fx.redraw(tt)

	t.Eq("          0123456789", l.String())
	for _, c := range l {
		switch c.Rune {
		case ' ':
			t.Eq(c.Style, s1)
		default:
			t.Eq(c.Style, s2)
		}
	}
}

func (s *ALine) Is_highlighted_if_highlight_flag_set(t *T) {
	tt, fx := newLineFX(t)
	fx.Switch(Highlighted)
	fx.gg.SetStyle(Default, fx.gg.Style(Default).WithAA(Dim))
	fx.gg.SetStyle(Highlight,
		NewStyle(Dim, RebeccaPurple, DarkGoldenrod))
	l, hStyle := fx.redraw(tt), fx.highlighted(fx.gg.Style(Default))
	t.Not.Eq(fx.gg.Style(Default), hStyle)
	for _, c := range l {
		t.Eq(c.Style, hStyle)
	}
	fx.Switch(Highlighted | TrimmedHighlighted)
	l = fx.redraw(tt)
	for _, c := range l {
		t.Eq(c.Style, hStyle)
	}
}

func (s *ALine) Is_highlighted_trimmed_if_corresponding_flag_set(t *T) {
	tt, fx := newLineFX(t)
	fx.Switch(TrimmedHighlighted)
	fx.setAt(4, []rune("0123456789"))
	l, hStyle := fx.redraw(tt), fx.highlighted(fx.gg.Style(Default))
	t.Not.Eq(fx.gg.Style(Default), hStyle)
	for _, c := range l {
		switch c.Rune {
		case ' ':
			t.Eq(c.Style, fx.gg.Style(Default))
		default:
			t.Eq(c.Style, hStyle)
		}
	}
}

func (s *ALine) Adapts_styles_overlapping_trimmed_highlighted_range(t *T) {
	s1, s2 := NewStyle(Blink, Yellow, Blue), NewStyle(Dim, Green, Red)
	tt, fx := newLineFX(t)
	fx.Switch(TrimmedHighlighted)
	fx.setStyledAt(2, []rune("  012"), s1)
	fx.setAt(7, []rune("3456"))
	fx.setStyledAt(11, []rune("789  "), s2)
	hs, hs1, hs2 := fx.highlighted(fx.gg.Style(Default)),
		fx.highlighted(s1), fx.highlighted(s2)
	l := fx.redraw(tt)
	for i, c := range l {
		switch i {
		case 0, 1:
			t.Eq(c.Style, fx.gg.Style(Default))
		case 2, 3:
			t.Eq(c.Style, s1)
		case 4, 5, 6:
			t.Eq(c.Style, hs1)
		case 7, 8, 9, 10:
			t.Eq(c.Style, hs)
		case 11, 12, 13:
			t.Eq(c.Style, hs2)
		case 14, 15:
			t.Eq(c.Style, s2)
		default:
			t.Eq(c.Style, fx.gg.Style(Default))
		}
	}
}

func (s *ALine) Adapts_enclosed_styles_in_trimmed_highlighted(t *T) {
	s1, s2 := NewStyle(Blink, Yellow, Blue), NewStyle(Dim, Green, Red)
	tt, fx := newLineFX(t)
	fx.Switch(TrimmedHighlighted)
	fx.setAt(4, []rune("01"))
	fx.setStyledAt(6, []rune("23"), s1)
	fx.setAt(8, []rune("456"))
	fx.setStyledAt(11, []rune("78"), s2)
	fx.setAt(13, []rune("9"))
	hs, hs1, hs2 := fx.highlighted(fx.gg.Style(Default)), fx.highlighted(s1),
		fx.highlighted(s2)
	l := fx.redraw(tt)
	for i, c := range l {
		switch i {
		case 0, 1, 2, 3:
			t.Eq(c.Style, fx.gg.Style(Default))
		case 4, 5:
			t.Eq(c.Style, hs)
		case 6, 7:
			t.Eq(c.Style, hs1)
		case 8, 9, 10:
			t.Eq(c.Style, hs)
		case 11, 12:
			t.Eq(c.Style, hs2)
		case 13:
			t.Eq(c.Style, hs)
		default:
			t.Eq(c.Style, fx.gg.Style(Default))
		}
	}
}

func (s *ALine) Adapts_enclosing_style_of_trimmed_highlighted(t *T) {
	s1 := NewStyle(Blink, Yellow, Blue)
	tt, fx := newLineFX(t)
	fx.Switch(TrimmedHighlighted)
	fx.setStyledAt(3, []rune(" 0123456789 "), s1)
	hs1 := fx.highlighted(s1)
	l := fx.redraw(tt)
	for i, c := range l {
		switch i {
		case 0, 1, 2:
			t.Eq(c.Style, fx.gg.Style(Default))
		case 3:
			t.Eq(c.Style, s1)
		case 4, 5, 6, 7, 8, 9, 10, 11, 12, 13:
			t.Eq(c.Style, hs1)
		case 14:
			t.Eq(c.Style, s1)
		default:
			t.Eq(c.Style, fx.gg.Style(Default))
		}
	}
}

func TestALine(t *testing.T) {
	t.Parallel()
	Run(&ALine{}, t)
}
