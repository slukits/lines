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
	line
	w, h int
	dflt Style
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
	rr, ss := x.display(x.width(), x.dflt)
	for i, r := range rr {
		x.Display(i, 0, r, ss.of(i))
	}
	x.Redraw()
	return tt.CellsArea(0, 0, x.width(), 1)[0]
}

// NOTE since the point here is to determine what a line provides for
// the display it doesn't matter to what backend it goes as long as we
// can figure out what went to the display with what style attributes.
func fx(t *T) (*term.Fixture, *lineFX) {
	ui, tt := term.LstFixture(t.GoT(), nil, 0)
	tt.PostResize(20, 1)
	return tt, &lineFX{Displayer: ui, dflt: DefaultStyle}
}

func fxDflt(t *T, s Style) (*term.Fixture, *lineFX) {
	tt, fx := fx(t)
	fx.dflt = s
	return tt, fx
}

type ALine struct{ Suite }

func (s *ALine) SetUp(t *T) { t.Parallel() }

func (s *ALine) Is_padded_with_spaces_if_zero(t *T) {
	tt, fx := fx(t)
	scrLine := fx.redraw(tt)
	for _, c := range scrLine {
		t.True(c.Rune == ' ')
		t.True(c.Style.Equals(DefaultStyle))
	}
}

func (s *ALine) Uses_given_display_style_if_no_default_set(t *T) {
	exp := NewStyle(Blink, Yellow, Blue)
	tt, fx := fxDflt(t, exp)
	scrLine := fx.redraw(tt)
	for _, c := range scrLine {
		t.True(c.Rune == ' ')
		t.True(c.Style.Equals(exp))
	}
}

func (s *ALine) Has_set_default_style_if_empty(t *T) {
	tt, fx := fx(t)
	exp := NewStyle(Blink, Yellow, Blue)
	fx.setDefaultStyle(exp)
	scrLine := fx.redraw(tt)
	for _, c := range scrLine {
		t.True(c.Rune == ' ')
		t.True(c.Style.Equals(exp))
	}
}

func (s *ALine) Has_updated_default_style_attributes(t *T) {
	tt, fx := fx(t)
	exp := DefaultStyle.WithAA(Dim)
	fx.withAA(Dim)
	scrLine := fx.redraw(tt)
	for _, c := range scrLine {
		t.True(c.Rune == ' ')
		t.True(c.Style.Equals(exp))
	}
}

func (s *ALine) Has_updated_default_foreground_color(t *T) {
	tt, fx := fx(t)
	exp := DefaultStyle.WithFG(Green)
	fx.withFG(Green)
	scrLine := fx.redraw(tt)
	for _, c := range scrLine {
		t.True(c.Rune == ' ')
		t.True(c.Style.Equals(exp))
	}
}

func (s *ALine) Has_updated_default_background_color(t *T) {
	tt, fx := fx(t)
	exp := DefaultStyle.WithBG(Red)
	fx.withBG(Red)
	scrLine := fx.redraw(tt)
	for _, c := range scrLine {
		t.True(c.Rune == ' ')
		t.True(c.Style.Equals(exp))
	}
}

func (s *ALine) Displays_set_content_space_padded_to_line_width(t *T) {
	_, fx := fx(t)
	fx.set("0123456789")
	got, _ := fx.display(fx.width(), fx.dflt)
	t.Eq("0123456789", string(got[:10]))
	t.Eq("          ", string(got[10:]))
}

func (s *ALine) Truncates_line_with_width_overflowing_content(t *T) {
	_, fx := fx(t)
	fx.set("01234567890123456789012")
	got, _ := fx.display(fx.width(), fx.dflt)
	t.Eq("01234567890123456789", string(got))
}

func (s *ALine) Displays_content_with_set_style(t *T) {
	exp := NewStyle(Blink, Yellow, Blue)
	tt, fx := fx(t)
	fx.setStyled("0123456789", exp)
	scrLine := fx.redraw(tt)
	for _, c := range scrLine {
		switch c.Rune {
		case ' ':
			t.True(c.Style.Equals(fx.dflt))
		default:
			t.True(c.Style.Equals(exp))
		}
	}
}

func (s *ALine) Has_content_set_at_zero_position(t *T) {
	_, fx := fx(t)
	fx.setAt(0, []rune("0123456789"))
	got, _ := fx.display(fx.width(), fx.dflt)
	t.Eq("0123456789", string(got[:10]))
	t.Eq("          ", string(got[10:]))
}

func (s *ALine) Has_content_set_at_given_position_space_padded(t *T) {
	_, fx := fx(t)
	fx.setAt(10, []rune("0123456789"))
	got, _ := fx.display(fx.width(), fx.dflt)
	t.Eq("          ", string(got[:10]))
	t.Eq("0123456789", string(got[10:]))
}

func (s *ALine) Styles_content_set_at_given_position(t *T) {
	exp := NewStyle(Blink, Yellow, Blue)
	tt, fx := fx(t)
	fx.setStyledAt(10, []rune("0123456789"), exp)
	scrLine := fx.redraw(tt)
	for _, c := range scrLine {
		switch c.Rune {
		case ' ':
			t.True(c.Style.Equals(fx.dflt))
		default:
			t.True(c.Style.Equals(exp))
		}
	}
}

func (s *ALine) Overwrites_content_after_given_position(t *T) {
	_, fx := fx(t)
	fx.set("0123456789")
	got, _ := fx.display(10, fx.dflt)
	t.Eq("0123456789", string(got))

	fx.setAt(2, []rune("42"))
	got, _ = fx.display(10, fx.dflt)
	t.Eq("0142      ", string(got))
}

func (s *ALine) Fills_remaining_space_with_a_filling_rune(t *T) {
	_, fx := fx(t)
	fx.setAtFilling(0, 'a')
	fx.setAt(1, []rune("0123456789"))
	got, _ := fx.display(fx.width(), fx.dflt)
	t.Eq("aaaaaaaaaa0123456789", string(got))
	fx.set("01234")
	fx.setAtFilling(5, 'a')
	fx.setAt(6, []rune("56789"))
	got, _ = fx.display(fx.width(), fx.dflt)
	t.Eq("01234aaaaaaaaaa56789", string(got))
	fx.setAt(0, []rune("0123456789")) // test filler truncation
	fx.setAtFilling(10, 'a')
	got, _ = fx.display(fx.width(), fx.dflt)
	t.Eq("0123456789aaaaaaaaaa", string(got))
}

func (s *ALine) Styles_filling_runes_with_set_style(t *T) {
	s1, s2 := NewStyle(Blink, Yellow, Blue), NewStyle(Dim, Green, Red)
	_, fx := fx(t)
	fx.setStyledAtFilling(0, 'a', s1)
	fx.setAt(1, []rune("0123456789"))
	_, ss := fx.display(fx.width(), fx.dflt)
	for r, s := range ss {
		switch r.End() {
		case 10:
			t.Eq(0, r.Start())
			s.Equals(s1)
		case 20:
			t.Eq(10, r.Start())
			s.Equals(fx.dflt)
		}
	}
	fx.set("01234")
	fx.setStyledAtFilling(5, 'a', s1)
	fx.setStyledAt(6, []rune("56789"), s2)
	_, ss = fx.display(fx.width(), fx.dflt)
	for r, s := range ss {
		switch r.End() {
		case 5:
			t.Eq(0, r.Start())
			s.Equals(fx.dflt)
		case 15:
			t.Eq(5, r.Start())
			s.Equals(s1)
		case 20:
			t.Eq(15, r.Start())
			s.Equals(s2)
		}
	}
	fx.setAt(0, []rune("0123456789"))
	fx.setStyledAtFilling(10, 'a', s1)
	_, ss = fx.display(fx.width(), fx.dflt)
	for r, s := range ss {
		switch r.End() {
		case 10:
			t.Eq(0, r.Start())
			s.Equals(fx.dflt)
		case 20:
			t.Eq(10, r.Start())
			s.Equals(s1)
		}
	}
}

func TestALine(t *testing.T) {
	t.Parallel()
	Run(&ALine{}, t)
}
