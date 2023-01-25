// Copyright (c) 2022 Stephan Lukits. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package selects

import (
	"fmt"
	"testing"
	"time"

	. "github.com/slukits/gounit"
	"github.com/slukits/lines"
	"github.com/slukits/lines/cmp/fx"
)

type AHselection struct{ Suite }

func (s *AHselection) Has_a_height_of_one(t *T) {
	cmp := &Horizontal{}
	fx := fx.New(t, cmp)
	fx.Lines.Update(cmp, nil, func(e *lines.Env) {
		t.Eq(1, cmp.Dim().Height())
	})
}

func (s *AHselection) Displays_a_default_label_if_none_set(t *T) {
	fx := fx.New(t, &Horizontal{})
	t.Contains(fx.ScreenOf(fx.Root()), NoLabel)
}

func (s *AHselection) Displays_a_default_item_if_none_set(t *T) {
	fx := fx.New(t, &Horizontal{})
	t.Contains(fx.ScreenOf(fx.Root()), NoItems)
}

const maxItem = "1234"

var lblFX = "label:"
var iiFX = []string{"12", maxItem, "123"}

func (s *AHselection) Width_defaults_to_label_plus_items_width_and_three(
	t *T,
) {
	// the three come from a blank after the label and after the
	// max-width item-value and the "Drop"-character.
	zeroWidth := len(NoLabel) + 1 + len(NoItems) + 1 + len([]rune(Drop))
	fx_ := fx.New(t, &Horizontal{})
	fx_.Lines.Update(fx_.Root(), nil, func(e *lines.Env) {
		t.Eq(zeroWidth, fx_.Root().(lines.Dimer).Dim().Width())
	})

	calcWidth := len(maxItem) + 1 + len(lblFX) + 1 + len([]rune(Drop))
	fx_ = fx.New(t, &Horizontal{Label: lblFX, Items: iiFX})
	fx_.Lines.Update(fx_.Root(), nil, func(e *lines.Env) {
		t.Eq(calcWidth, fx_.Root().(lines.Dimer).Dim().Width())
	})
}

func (s *AHselection) Respects_set_maximum_width(t *T) {
	fx := fx.New(t, &Horizontal{Label: lblFX, Items: iiFX, MaxWidth: 3})
	calcWidth := len(maxItem) + len(lblFX) + 1 + len([]rune(Drop))
	fx.Lines.Update(fx.Root(), nil, func(e *lines.Env) {
		t.Eq(calcWidth, fx.Root().(lines.Dimer).Dim().Width())
	})
}

func (s *AHselection) Without_default_has_a_blank_items_label(t *T) {
	fx := fx.New(t, &Horizontal{
		Label: lblFX, Items: iiFX, DefaultItem: NoDefault})
	t.Eq(
		fmt.Sprintf("%s      %s", lblFX, Drop),
		fx.ScreenOf(fx.Root()),
	)
}

func extractII(cmp *Horizontal) *items { return cmp.CC[1].(*items) }

type horizontalFX struct {
	component
	lines.Stacking
}

type filler struct{ component }

func hrzFX(t *T, hs *Horizontal, tt ...time.Duration) *lines.Fixture {
	hsFX := &horizontalFX{}
	hsFX.CC = append(hsFX.CC, hs, &filler{})
	fx := fx.New(t, hsFX, tt...)
	fx.FireResize(len(maxItem)+1+len(lblFX)+1+len([]rune(Drop)), 4)
	return fx
}

func hrzCmpFX(
	t *T, styler func(int) lines.Style, tt ...time.Duration,
) (*lines.Fixture, *Horizontal) {
	hs := &Horizontal{Label: lblFX, Items: iiFX, Styler: styler}
	return hrzFX(t, hs, tt...), hs
}

// extractLayer clicks on given selects.Horizontal instance cmp's items
// component and retrieves the consequently displayed layer with items.
func extractLayer(t *T, fx *lines.Fixture, cmp *Horizontal) *ModalList {
	fx.FireComponentClick(extractII(cmp), len(maxItem)+1, 0)
	lyr := (*ModalList)(nil)
	fx.Lines.Update(cmp, nil, func(e *lines.Env) {
		l, ok := e.Focused().(*ModalList)
		t.FatalIfNot(t.True(ok))
		lyr = l
	})
	return lyr
}

func closeItemList(fx *lines.Fixture, cmp *Horizontal) {
	fx.FireComponentClick(extractII(cmp), len(maxItem)+1, 0)
}

func (s *AHselection) Item_styles_default_to_reversed_globals(t *T) {
	fx, cmp := hrzCmpFX(t, nil)
	rvrDflt := fx.Lines.Globals.Style(lines.Default).WithAA(
		fx.Lines.Globals.Style(lines.Default).AA() | lines.Reverse)
	lyr := extractLayer(t, fx, cmp)
	lyrScr := fx.CellsOf(lyr)
	for _, l := range lyrScr {
		t.True(l[0].Style == rvrDflt)
	}

	closeItemList(fx, cmp)
	fx.Lines.Globals.SetAA(
		lines.Default,
		fx.Lines.Globals.AA(lines.Default)|lines.Reverse,
	)

	rvrDflt = fx.Lines.Globals.Style(lines.Default).WithAA(
		fx.Lines.Globals.Style(lines.Default).AA() &^ lines.Reverse)
	lyr = extractLayer(t, fx, cmp)
	lyrScr = fx.CellsOf(lyr)
	for _, l := range lyrScr {
		t.True(l[0].Style == rvrDflt)
	}
}

func (s *AHselection) Uses_given_items_styles(t *T) {
	sty := lines.DefaultStyle.WithBG(lines.DarkBlue).
		WithFG(lines.Silver).WithAA(lines.Bold)
	fx, cmp := hrzCmpFX(t, func(i int) lines.Style {
		return sty
	})

	lyr := extractLayer(t, fx, cmp)
	lyrScr := fx.CellsOf(lyr)
	for _, l := range lyrScr {
		t.True(l[0].Style == sty)
	}
}

// func (s *AHselection) Item_highlight_defaults_to_global_style(t *T) {
// 	fx, cmp := hrzCmpFX(t, nil)
// 	lyr, x, y := extractLayer(t, fx, cmp), 0, 0
// 	fx.Lines.Update(lyr, nil, func(e *lines.Env) {
// 		x, y, _, _ = lyr.Dim().Printable()
// 	})
//
// 	fx.FireMove(x+1, y+1) // move to/highlight the second layer item
// 	cell := fx.CellsOf(lyr)[1][0]
// 	t.True(cell.Style.Equals(fx.Lines.Globals.Style(lines.Default)))
// }
//
// func (s *AHselection) Item_highlight_uses_given_style(t *T) {
// 	sty := lines.DefaultStyle.WithBG(lines.DarkBlue).
// 		WithFG(lines.Silver).WithAA(lines.Bold)
// 	fx, cmp := hrzCmpFX(t, func(i int, highlight bool) lines.Style {
// 		if !highlight {
// 			return lines.DefaultStyle.WithAA(lines.Reverse)
// 		}
// 		return sty
// 	})
// 	lyr, x, y := extractLayer(t, fx, cmp), 0, 0
// 	fx.Lines.Update(lyr, nil, func(e *lines.Env) {
// 		x, y, _, _ = lyr.Dim().Printable()
// 	})
//
// 	fx.FireMove(x+1, y+1) // move to/highlight the second layer item
// 	cell := fx.CellsOf(lyr)[1][0]
// 	t.True(cell.Style.Equals(sty))
// }
//
// func (s *AHselection) Removes_highlight_from_focus_loosing_item(t *T) {
// 	fx, cmp := hrzCmpFX(t, nil)
// 	lyr, x, y := extractLayer(t, fx, cmp), 0, 0
// 	fx.Lines.Update(lyr, nil, func(e *lines.Env) {
// 		x, y, _, _ = lyr.Dim().Printable()
// 	})
//
// 	fx.FireMove(x+1, y+1) // move to/highlight the second layer item
// 	cell := fx.CellsOf(lyr)[1][0]
// 	t.True(cell.Style.Equals(fx.Lines.Globals.Style(lines.Default)))
//
// 	fx.FireMove(x+1, y) // move out of item
// 	cell = fx.CellsOf(lyr)[1][0]
// 	t.True(cell.Style.Equals(fx.Lines.Globals.Style(lines.Default).
// 		WithAA(lines.Reverse)))
// }
//
// func (s *AHselection) Removes_highlight_on_out_of_bounds_move(t *T) {
// 	fx, cmp := hrzCmpFX(t, nil)
// 	lyr, x, y := extractLayer(t, fx, cmp), 0, 0
// 	fx.Lines.Update(lyr, nil, func(e *lines.Env) {
// 		x, y, _, _ = lyr.Dim().Printable()
// 	})
//
// 	fx.FireMove(x+1, y+1) // move to/highlight the second layer item
// 	cell := fx.CellsOf(lyr)[1][0]
// 	t.True(cell.Style.Equals(fx.Lines.Globals.Style(lines.Default)))
//
// 	fx.FireMove(x+1, y-1) // move out of layer
// 	cell = fx.CellsOf(lyr)[1][0]
// 	t.True(cell.Style.Equals(fx.Lines.Globals.Style(lines.Default).
// 		WithAA(lines.Reverse)))
// }

func (s *AHselection) Selects_clicked_item(t *T) {
	fx, cmp := hrzCmpFX(t, nil)
	lyr, x, y := extractLayer(t, fx, cmp), 0, 0
	fx.Lines.Update(lyr, nil, func(e *lines.Env) {
		x, y, _, _ = lyr.Dim().Printable()
	})

	fx.FireClick(x+1, y+1) // select the second layer item
	t.Eq(1, cmp.Value())
	t.Contains(fx.ScreenOf(cmp), maxItem)
}

// func (s *AHselection) Zeros_on_zero_select_iff_no_default(t *T) {
// 	cmp := &Horizontal{
// 		Label: lblFX, Items: iiFX, DefaultItem: NoDefault, MaxWidth: -5}
// 	fx := hrzFX(t, cmp)
// 	lyr, x, y := extractLayer(t, fx, cmp), 0, 0
// 	fx.Lines.Update(lyr, nil, func(e *lines.Env) {
// 		x, y, _, _ = lyr.Dim().Printable()
// 	})
// 	fx.FireClick(x, y) // select the second layer item
// 	t.Eq(0, cmp.Value())
// 	t.Contains(fx.ScreenOf(cmp), "12")
//
// 	lyr, x, y = extractLayer(t, fx, cmp), 0, 0
// 	fx.Lines.Update(lyr, nil, func(e *lines.Env) {
// 		x, y, _, _ = lyr.Dim().Printable()
// 	})
// 	fx.FireClick(x, y-1)
// 	t.Eq(fmt.Sprintf("%s      %s", lblFX, Drop), fx.ScreenOf(cmp))
// }

func TestAHselection(t *testing.T) {
	t.Parallel()
	Run(&AHselection{}, t)
}
