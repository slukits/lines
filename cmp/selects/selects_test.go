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

type ADropDown struct{ Suite }

func (s *ADropDown) SetUp(t *T) { t.Parallel() }

func (s *ADropDown) Has_a_height_of_one(t *T) {
	cmp := &DropDown{}
	fx := fx.New(t, cmp)
	fx.Lines.Update(cmp, nil, func(e *lines.Env) {
		t.Eq(1, cmp.Dim().Height())
	})
}

func (s *ADropDown) Displays_a_default_item_if_none_set(t *T) {
	fx := fx.New(t, &DropDown{})
	t.Contains(fx.ScreenOf(fx.Root()), NoItems)
}

const maxItem = "1234"

var lblFX = "label:"
var iiFX = []string{"12", maxItem, "123"}

func (s *ADropDown) Width_defaults_to_max_items_width_and_two(t *T) {
	// the three come from a blank after the max-width item-value and
	// the "Drop"-character.
	zeroWidth := len(NoItems) + 1 + len([]rune(Drop))
	fx_ := fx.New(t, &DropDown{})
	fx_.Lines.Update(fx_.Root(), nil, func(e *lines.Env) {
		t.Eq(zeroWidth, fx_.Root().(lines.Dimer).Dim().Width())
	})

	calcWidth := len(maxItem) + 1 + len([]rune(Drop))
	fx_ = fx.New(t, &DropDown{Items: iiFX})
	fx_.Lines.Update(fx_.Root(), nil, func(e *lines.Env) {
		t.Eq(calcWidth, fx_.Root().(lines.Dimer).Dim().Width())
	})
}

func (s *ADropDown) Shows_initially_set_default_item(t *T) {
	fx := fx.New(t, &DropDown{Items: iiFX, DefaultItem: 1})
	t.Contains(fx.ScreenOf(fx.Root()), maxItem)
}

func (s *ADropDown) Respects_set_maximum_width(t *T) {
	fx := fx.New(t, &DropDown{Items: iiFX, MaxWidth: 3})
	expWidth := 3 + 1 + len([]rune(Drop))
	fx.Lines.Update(fx.Root(), nil, func(e *lines.Env) {
		t.Eq(expWidth, fx.Root().(lines.Dimer).Dim().Width())
	})
}

func (s *ADropDown) Shortens_not_fitting_item_label(t *T) {
	fx := fx.New(t, &DropDown{Items: iiFX, MaxWidth: 3, DefaultItem: 1})
	t.Contains(fx.ScreenOf(fx.Root()), "12…")
}

func extractDropDownLayer(
	t *T, fx *lines.Fixture, cmp *DropDown,
) *ModalList {
	fx.FireComponentClick(cmp, 0, 0)
	lyr, ok := (*ModalList)(nil), true
	fx.Lines.Update(cmp, nil, func(e *lines.Env) {
		// NOTE we can not fail the test in an other go-routine
		lyr, ok = e.Focused().(*ModalList)
	})
	if !ok {
		t.Fatalf("drop down: extract layer: not found")
	}
	return lyr
}

func closeDropDownItems(fx *lines.Fixture, cmp *DropDown) {
	fx.FireComponentClick(cmp, len(maxItem)+1, 0)
}

func (s *ADropDown) Shows_max_height_many_items(t *T) {
	cmp := &DropDown{Items: iiFX, MaxHeight: 2}
	fx := fx.New(t, cmp)
	extractDropDownLayer(t, fx, cmp)
	t.SpaceMatched(fx.Screen(), "12", "1234")
	t.Not.SpaceMatched(fx.Screen(), "12", "1234", "123")
}

func (s *ADropDown) Items_styles_default_to_reversed_globals(t *T) {
	cmp := &DropDown{Items: iiFX}
	fx := fx.New(t, cmp)
	rvrDflt := fx.Lines.Globals.Style(lines.Highlight)
	lyr := extractDropDownLayer(t, fx, cmp)
	lyrScr := fx.CellsOf(lyr)
	for _, l := range lyrScr {
		t.True(l[0].Style == rvrDflt)
	}
}

func (s *ADropDown) Uses_given_items_styles(t *T) {
	sty := lines.DefaultStyle.WithBG(lines.DarkBlue).
		WithFG(lines.Silver).WithAA(lines.Bold)
	cmp := &DropDown{Items: iiFX}
	fx := fx.New(t, cmp)
	cmp.Styler = func(_ int) lines.Style { return sty }
	lyr := extractDropDownLayer(t, fx, cmp)
	lyrScr := fx.CellsOf(lyr)
	for _, l := range lyrScr {
		t.Eq(l[0].Style, sty)
	}
}

func (s *ADropDown) Item_highlight_defaults_to_global_style(t *T) {
	cmp := &DropDown{Items: iiFX}
	fx := fx.New(t, cmp)
	lyr, x, y := extractDropDownLayer(t, fx, cmp), 0, 0
	fx.Lines.Update(lyr, nil, func(e *lines.Env) {
		x, y, _, _ = lyr.Dim().Printable()
	})

	fx.FireMove(x+1, y+1) // move to/highlight the second layer item
	cell := fx.CellsOf(lyr)[1][0]
	t.Eq(cell.Style, fx.Lines.Globals.Style(lines.Default))
}

func (s *ADropDown) Item_highlight_uses_given_style(t *T) {
	sty := lines.DefaultStyle.WithBG(lines.DarkBlue).
		WithFG(lines.Silver).WithAA(lines.Bold)
	cmp := &DropDown{
		Items: iiFX,
		Highlighter: func(_ lines.Style) lines.Style {
			return sty
		}}
	fx := fx.New(t, cmp)
	lyr, x, y := extractDropDownLayer(t, fx, cmp), 0, 0
	fx.Lines.Update(lyr, nil, func(e *lines.Env) {
		x, y, _, _ = lyr.Dim().Printable()
	})

	fx.FireMove(x+1, y+1) // move to/highlight the second layer item
	cell := fx.CellsOf(lyr)[1][0]
	t.Eq(cell.Style, sty)
}

func (s *ADropDown) Removes_highlight_from_focus_loosing_item(t *T) {
	cmp := &DropDown{Items: iiFX}
	fx := fx.New(t, cmp)
	lyr, x, y := extractDropDownLayer(t, fx, cmp), 0, 0
	fx.Lines.Update(lyr, nil, func(e *lines.Env) {
		x, y, _, _ = lyr.Dim().Printable()
	})

	fx.FireMove(x+1, y+1) // move to/highlight the second layer item
	cell := fx.CellsOf(lyr)[1][0]
	t.Eq(cell.Style, fx.Lines.Globals.Style(lines.Default))

	fx.FireMove(x+1, y) // move out of item
	cell = fx.CellsOf(lyr)[1][0]
	t.Eq(cell.Style, fx.Lines.Globals.Style(lines.Default).
		WithAA(lines.Reverse))
}

func (s *ADropDown) Removes_highlight_on_out_of_bounds_move(t *T) {
	cmp := &DropDown{Items: iiFX}
	fx := fx.New(t, cmp)
	lyr, x, y := extractDropDownLayer(t, fx, cmp), 0, 0
	fx.Lines.Update(lyr, nil, func(e *lines.Env) {
		x, y, _, _ = lyr.Dim().Printable()
	})

	fx.FireMove(x+1, y+1) // move to/highlight the second layer item
	cell := fx.CellsOf(lyr)[1][0]
	t.Eq(cell.Style, fx.Lines.Globals.Style(lines.Default))

	fx.FireMove(x+1, y-1) // move out of layer
	cell = fx.CellsOf(lyr)[1][0]
	t.Eq(cell.Style, fx.Lines.Globals.Style(lines.Default).
		WithAA(lines.Reverse))
}

func (s *ADropDown) Selects_clicked_item(t *T) {
	cmp := &DropDown{Items: iiFX}
	fx := fx.New(t, cmp)
	lyr, x, y := extractDropDownLayer(t, fx, cmp), 0, 0
	fx.Lines.Update(lyr, nil, func(e *lines.Env) {
		x, y, _, _ = lyr.Dim().Printable()
	})

	fx.FireClick(x+1, y+1) // select the second layer item
	t.Eq(1, cmp.Value())
	t.Contains(fx.ScreenOf(cmp), maxItem)
}

func (s *ADropDown) Closes_on_item_selection(t *T) {
	cmp := &DropDown{Items: iiFX}
	fx := fx.New(t, cmp)
	lyr, x, y := extractDropDownLayer(t, fx, cmp), 0, 0
	fx.Lines.Update(lyr, nil, func(e *lines.Env) {
		x, y, _, _ = lyr.Dim().Printable()
	})

	fx.FireClick(x+1, y+1) // select the second layer item
	t.Eq(1, cmp.Value())
	t.Contains(fx.ScreenOf(cmp), maxItem)
	t.Not.SpaceMatched(fx.Screen(), "12", "\n", "1234", "\n", "123")
}

func (s *ADropDown) Closes_on_out_of_bounds_click(t *T) {
	cmp := &DropDown{Items: iiFX}
	fx := fx.New(t, cmp)
	lyr, x, y := extractDropDownLayer(t, fx, cmp), 0, 0
	fx.Lines.Update(lyr, nil, func(e *lines.Env) {
		x, y, _, _ = lyr.Dim().Printable()
	})
	t.Contains(fx.Screen(), maxItem)

	fx.FireClick(x+1, y-1) // click out of bounds
	t.Not.Contains(fx.Screen(), maxItem)
}

func (s *ADropDown) Zeros_on_zero_select_iff_no_default(t *T) {
	cmp := &DropDown{Items: iiFX, DefaultItem: NoDefault, MaxWidth: -5}
	fx := fx.New(t, cmp)
	lyr, x, y := extractDropDownLayer(t, fx, cmp), 0, 0
	fx.Lines.Update(lyr, nil, func(e *lines.Env) {
		x, y, _, _ = lyr.Dim().Printable()
	})
	fx.FireClick(x, y) // select the second layer item
	t.Eq(0, cmp.Value())
	t.Contains(fx.ScreenOf(cmp), "12")

	lyr, x, y = extractDropDownLayer(t, fx, cmp), 0, 0
	fx.Lines.Update(lyr, nil, func(e *lines.Env) {
		x, y, _, _ = lyr.Dim().Printable()
	})
	fx.FireClick(x, y-1)
	t.Eq(fmt.Sprintf("     %s", Drop), fx.ScreenOf(cmp))
}

func TestADropDown(t *testing.T) {
	t.Parallel()
	Run(&ADropDown{}, t)
}

type AHrzDrop struct{ Suite }

func (s *AHrzDrop) SetUp(t *T) { t.Parallel() }

func (s *AHrzDrop) Has_a_height_of_one(t *T) {
	cmp := &DropDownHrz{}
	fx := fx.New(t, cmp)
	fx.Lines.Update(cmp, nil, func(e *lines.Env) {
		t.Eq(1, cmp.Dim().Height())
	})
}

func (s *AHrzDrop) Displays_a_default_label_if_none_set(t *T) {
	fx := fx.New(t, &DropDownHrz{})
	t.Contains(fx.ScreenOf(fx.Root()), NoLabel)
}

func (s *AHrzDrop) Displays_a_default_item_if_none_set(t *T) {
	fx := fx.New(t, &DropDownHrz{})
	t.Contains(fx.ScreenOf(fx.Root()), NoItems)
}

func (s *AHrzDrop) Width_defaults_to_label_plus_items_width_and_three(
	t *T,
) {
	// the three come from a blank after the label and after the
	// max-width item-value and the "Drop"-character.
	zeroWidth := len(NoLabel) + 1 + len(NoItems) + 1 + len([]rune(Drop))
	fx_ := fx.New(t, &DropDownHrz{})
	fx_.Lines.Update(fx_.Root(), nil, func(e *lines.Env) {
		t.Eq(zeroWidth, fx_.Root().(lines.Dimer).Dim().Width())
	})

	calcWidth := len(maxItem) + 1 + len(lblFX) + 1 + len([]rune(Drop))
	fx_ = fx.New(t, &DropDownHrz{Label: lblFX, Items: iiFX})
	fx_.Lines.Update(fx_.Root(), nil, func(e *lines.Env) {
		t.Eq(calcWidth, fx_.Root().(lines.Dimer).Dim().Width())
	})
}

func (s *AHrzDrop) Respects_set_maximum_width(t *T) {
	// TODO: stumbled over this: looks wrong!
	fx := fx.New(t, &DropDownHrz{Label: lblFX, Items: iiFX, MaxWidth: 3})
	calcWidth := len(maxItem) + len(lblFX) + 1 + len([]rune(Drop))
	fx.Lines.Update(fx.Root(), nil, func(e *lines.Env) {
		t.Eq(calcWidth, fx.Root().(lines.Dimer).Dim().Width())
	})
}

func (s *AHrzDrop) Without_default_has_a_blank_items_label(t *T) {
	fx := fx.New(t, &DropDownHrz{
		Label: lblFX, Items: iiFX, DefaultItem: NoDefault})
	t.Eq(
		fmt.Sprintf("%s      %s", lblFX, Drop),
		fx.ScreenOf(fx.Root()),
	)
}

func extractII(cmp *DropDownHrz) *items { return cmp.CC[1].(*items) }

type horizontalFX struct {
	component
	lines.Stacking
}

type filler struct{ component }

func hrzFX(t *T, hs *DropDownHrz, tt ...time.Duration) *lines.Fixture {
	hsFX := &horizontalFX{}
	hsFX.CC = append(hsFX.CC, hs, &filler{})
	fx := fx.New(t, hsFX, tt...)
	fx.FireResize(len(maxItem)+1+len(lblFX)+1+len([]rune(Drop)), 4)
	return fx
}

func hrzCmpFX(
	t *T, hi func(lines.Style) lines.Style, tt ...time.Duration,
) (*lines.Fixture, *DropDownHrz) {
	hs := &DropDownHrz{Label: lblFX, Items: iiFX, Highlighter: hi}
	return hrzFX(t, hs, tt...), hs
}

// extractLayer clicks on given selects.Horizontal instance cmp's items
// component and retrieves the consequently displayed layer with items.
func extractLayer(t *T, fx *lines.Fixture, cmp *DropDownHrz) *ModalList {
	fx.FireComponentClick(extractII(cmp), len(maxItem)+1, 0)
	lyr, ok := (*ModalList)(nil), true
	fx.Lines.Update(cmp, nil, func(e *lines.Env) {
		// NOTE we can not fail the test in an other go-routine
		lyr, ok = e.Focused().(*ModalList)
	})
	if !ok {
		t.Fatalf("drop down horizontal: extract layer: not found")
	}
	return lyr
}

func closeItemList(fx *lines.Fixture, cmp *DropDownHrz) {
	fx.FireComponentClick(extractII(cmp), len(maxItem)+1, 0)
}

func (s *AHrzDrop) Items_styles_default_to_reversed_globals(t *T) {
	fx, cmp := hrzCmpFX(t, nil)
	rvrDflt := fx.Lines.Globals.Style(lines.Highlight)
	lyr := extractLayer(t, fx, cmp)
	lyrScr := fx.CellsOf(lyr)
	for _, l := range lyrScr {
		t.True(l[0].Style == rvrDflt)
	}
}

func (s *AHrzDrop) Uses_given_items_styles(t *T) {
	sty := lines.DefaultStyle.WithBG(lines.DarkBlue).
		WithFG(lines.Silver).WithAA(lines.Bold)
	fx, cmp := hrzCmpFX(t, nil)
	cmp.Styler = func(_ int) lines.Style { return sty }
	lyr := extractLayer(t, fx, cmp)
	lyrScr := fx.CellsOf(lyr)
	for _, l := range lyrScr {
		t.Eq(l[0].Style, sty)
	}
}

func (s *AHrzDrop) Item_highlight_defaults_to_global_style(t *T) {
	fx, cmp := hrzCmpFX(t, nil)
	lyr, x, y := extractLayer(t, fx, cmp), 0, 0
	fx.Lines.Update(lyr, nil, func(e *lines.Env) {
		x, y, _, _ = lyr.Dim().Printable()
	})

	fx.FireMove(x+1, y+1) // move to/highlight the second layer item
	cell := fx.CellsOf(lyr)[1][0]
	t.Eq(cell.Style, fx.Lines.Globals.Style(lines.Default))
}

func (s *AHrzDrop) Item_highlight_uses_given_style(t *T) {
	sty := lines.DefaultStyle.WithBG(lines.DarkBlue).
		WithFG(lines.Silver).WithAA(lines.Bold)
	fx, cmp := hrzCmpFX(t, func(_ lines.Style) lines.Style {
		return sty
	})
	lyr, x, y := extractLayer(t, fx, cmp), 0, 0
	fx.Lines.Update(lyr, nil, func(e *lines.Env) {
		x, y, _, _ = lyr.Dim().Printable()
	})

	fx.FireMove(x+1, y+1) // move to/highlight the second layer item
	cell := fx.CellsOf(lyr)[1][0]
	t.Eq(cell.Style, sty)
}

func (s *AHrzDrop) Removes_highlight_from_focus_loosing_item(t *T) {
	fx, cmp := hrzCmpFX(t, nil)
	lyr, x, y := extractLayer(t, fx, cmp), 0, 0
	fx.Lines.Update(lyr, nil, func(e *lines.Env) {
		x, y, _, _ = lyr.Dim().Printable()
	})

	fx.FireMove(x+1, y+1) // move to/highlight the second layer item
	cell := fx.CellsOf(lyr)[1][0]
	t.Eq(cell.Style, fx.Lines.Globals.Style(lines.Default))

	fx.FireMove(x+1, y) // move out of item
	cell = fx.CellsOf(lyr)[1][0]
	t.Eq(cell.Style, fx.Lines.Globals.Style(lines.Default).
		WithAA(lines.Reverse))
}

func (s *AHrzDrop) Removes_highlight_on_out_of_bounds_move(t *T) {
	fx, cmp := hrzCmpFX(t, nil)
	lyr, x, y := extractLayer(t, fx, cmp), 0, 0
	fx.Lines.Update(lyr, nil, func(e *lines.Env) {
		x, y, _, _ = lyr.Dim().Printable()
	})

	fx.FireMove(x+1, y+1) // move to/highlight the second layer item
	cell := fx.CellsOf(lyr)[1][0]
	t.Eq(cell.Style, fx.Lines.Globals.Style(lines.Default))

	fx.FireMove(x+1, y-1) // move out of layer
	cell = fx.CellsOf(lyr)[1][0]
	t.Eq(cell.Style, fx.Lines.Globals.Style(lines.Default).
		WithAA(lines.Reverse))
}

func (s *AHrzDrop) Selects_clicked_item(t *T) {
	fx, cmp := hrzCmpFX(t, nil)
	lyr, x, y := extractLayer(t, fx, cmp), 0, 0
	fx.Lines.Update(lyr, nil, func(e *lines.Env) {
		x, y, _, _ = lyr.Dim().Printable()
	})

	fx.FireClick(x+1, y+1) // select the second layer item
	t.Eq(1, cmp.Value())
	t.Contains(fx.ScreenOf(cmp), maxItem)
}

func (s *AHrzDrop) Closes_on_item_selection(t *T) {
	fx, cmp := hrzCmpFX(t, nil)
	lyr, x, y := extractLayer(t, fx, cmp), 0, 0
	fx.Lines.Update(lyr, nil, func(e *lines.Env) {
		x, y, _, _ = lyr.Dim().Printable()
	})

	fx.FireClick(x+1, y+1) // select the second layer item
	t.Eq(1, cmp.Value())
	t.Contains(fx.ScreenOf(cmp), maxItem)
	t.Not.SpaceMatched(fx.Screen(), "12", "\n", "1234", "\n", "123")
}

func (s *AHrzDrop) Zeros_on_zero_select_iff_no_default(t *T) {
	cmp := &DropDownHrz{
		Label: lblFX, Items: iiFX, DefaultItem: NoDefault, MaxWidth: -5}
	fx := hrzFX(t, cmp)
	lyr, x, y := extractLayer(t, fx, cmp), 0, 0
	fx.Lines.Update(lyr, nil, func(e *lines.Env) {
		x, y, _, _ = lyr.Dim().Printable()
	})
	fx.FireClick(x, y) // select the second layer item
	t.Eq(0, cmp.Value())
	t.Contains(fx.ScreenOf(cmp), "12")

	lyr, x, y = extractLayer(t, fx, cmp), 0, 0
	fx.Lines.Update(lyr, nil, func(e *lines.Env) {
		x, y, _, _ = lyr.Dim().Printable()
	})
	fx.FireClick(x, y-1)
	t.Eq(fmt.Sprintf("%s      %s", lblFX, Drop), fx.ScreenOf(cmp))
}

func TestAHrzDrop(t *testing.T) {
	t.Parallel()
	Run(&AHrzDrop{}, t)
}

type AVrtDrop struct{ Suite }

func (s *AVrtDrop) SetUp(t *T) { t.Parallel() }

func (s *AVrtDrop) Has_a_height_of_two(t *T) {
	cmp := &DropDownVrt{}
	fx := fx.New(t, cmp)
	fx.Lines.Update(cmp, nil, func(e *lines.Env) {
		t.Eq(2, cmp.Dim().Height())
	})
}

func (s *AVrtDrop) Displays_a_default_label_if_none_set(t *T) {
	fx := fx.New(t, &DropDownHrz{})
	t.Contains(fx.ScreenOf(fx.Root()), NoLabel)
}

func (s *AVrtDrop) Displays_a_default_item_if_none_set(t *T) {
	fx := fx.New(t, &DropDownHrz{})
	t.Contains(fx.ScreenOf(fx.Root()), NoItems)
}

func (s *AVrtDrop) Without_default_has_a_blank_items_label(t *T) {
	fx := fx.New(t, &DropDownHrz{
		Label: lblFX, Items: iiFX, DefaultItem: NoDefault})
	t.Eq(
		fmt.Sprintf("%s      %s", lblFX, Drop),
		fx.ScreenOf(fx.Root()),
	)
}

func (s *AVrtDrop) Selects_clicked_item(t *T) {
	cmp := &DropDownVrt{Label: lblFX, Items: iiFX}
	fx := fx.Sized(t, 7, 7, cmp)
	fx.FireClick(4, 3)
	fx.FireClick(4, 5)
	t.Eq(1, cmp.Value())
	t.Contains(fx.ScreenOf(cmp), maxItem)
}

func TestAVrtDrop(t *testing.T) {
	t.Parallel()
	Run(&AVrtDrop{}, t)
}
