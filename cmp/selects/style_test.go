// Copyright (c) 2022 Stephan Lukits. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package selects

import (
	"strings"
	"testing"
	"time"

	. "github.com/slukits/gounit"
	"github.com/slukits/lines"
	"github.com/slukits/lines/cmp/fx"
)

type property struct{ Suite }

func (s *property) SetUp(t *T) { t.Parallel() }

func (s *property) Has_properties_as_zero_label_if_monochrome(t *T) {
	fx := fx.New(t, &StyleProperty{Styles: &Styles{}})
	t.Contains(fx.ScreenOf(fx.Root()), PropertyNames[Properties])
}

func (s *property) Offers_background_foreground_reverse_if_mono(t *T) {
	fx := fx.New(t, &StyleProperty{Styles: &Styles{}})
	lry := extractDropDownLayer(t, fx, fx.Root().(*StyleProperty))
	t.Contains(fx.ScreenOf(lry), PropertyNames[ReverseFgBgProperty])
}

func (s *property) Offers_style_attributes_if_monochrome(t *T) {
	fx := fx.New(t, &StyleProperty{Styles: &Styles{}})
	lry := extractDropDownLayer(t, fx, fx.Root().(*StyleProperty))
	t.Contains(fx.ScreenOf(lry), PropertyNames[StyleAttributeProperty])
}

func (s *property) Offers_reset_properties_if_monochrome(t *T) {
	fx := fx.New(t, &StyleProperty{Styles: &Styles{}})
	lry := extractDropDownLayer(t, fx, fx.Root().(*StyleProperty))
	t.Contains(fx.ScreenOf(lry), PropertyNames[StyleAttributeProperty])
}

func (s *property) Of_monochrome_has_reverse_attrs_and_reset_only(t *T) {
	fx := fx.New(t, &StyleProperty{Styles: &Styles{}})
	lry := extractDropDownLayer(t, fx, fx.Root().(*StyleProperty))
	t.SpaceMatched(fx.ScreenOf(lry), PropertyNames[ReverseFgBgProperty],
		PropertyNames[StyleAttributeProperty],
		PropertyNames[ResetProperties])
	t.Not.Contains(fx.ScreenOf(lry), PropertyNames[ForegroundProperty])
	t.Not.Contains(fx.ScreenOf(lry), PropertyNames[BackgroundProperty])
}

func selectModalProperty(
	t *T, p string, fx *lines.Fixture, cmp *ModalList,
) {
	y := -1
	for i, l := range fx.CellsOf(cmp) {
		if !strings.Contains(l.String(), p) {
			continue
		}
		y = i
		break
	}
	if y == -1 {
		t.Fatalf("select property '%s': not found", p)
	}
	fx.FireComponentClick(cmp, 0, y)
}

// selectProperty opens given style-properties drop down component cmp's
// selection list and selects given property p in given terminal test
// fixture fx.  It ends test execution if p is not found.
func selectProperty(t *T, p string, fx *lines.Fixture, cmp *StyleProperty) {
	selectModalProperty(t, p, fx, extractDropDownLayer(t, fx, cmp))
}

type ppSty struct {
	lines.Component
	lines.Chaining
}

func fxPPSty(t *T, cr ColorRange, timeout ...time.Duration) (
	*lines.Fixture, *StyleProperty,
) {
	cmp, pp, ss := &ppSty{}, &StyleProperty{}, &Styles{Colors: cr}
	pp.Styles = ss
	cmp.CC = append(cmp.CC, pp, ss)
	fx := fx.New(t, cmp, timeout...)
	return fx, pp
}

func (c *ppSty) OnInit(e *lines.Env) { c.Dim().SetHeight(1) }

func (s *property) Mono_reverse_switches_fg_and_bg(t *T) {
	fx, pp := fxPPSty(t, Monochrome)
	stylesLabelStyle := fx.CellsOf(pp.Styles).Trimmed()[0][0].Style
	selectProperty(t, PropertyNames[ReverseFgBgProperty], fx, pp)
	reversedLabelStyle := fx.CellsOf(pp.Styles).Trimmed()[0][0].Style
	t.Eq(stylesLabelStyle.FG(), reversedLabelStyle.BG())
	t.Eq(stylesLabelStyle.BG(), reversedLabelStyle.FG())
	selectProperty(t, PropertyNames[ReverseFgBgProperty], fx, pp)
	reversedLabelStyle = fx.CellsOf(pp.Styles).Trimmed()[0][0].Style
	t.Eq(stylesLabelStyle.FG(), reversedLabelStyle.FG())
	t.Eq(stylesLabelStyle.BG(), reversedLabelStyle.BG())
}

func (s *property) Mono_reverse_switches_back_to_zero_label(t *T) {
	fx, pp := fxPPSty(t, Monochrome)
	selectProperty(t, PropertyNames[ReverseFgBgProperty], fx, pp)
	t.Contains(fx.ScreenOf(pp).Trimmed(), PropertyNames[Properties])
}

func (s *property) Mono_reset_sets_initial_style(t *T) {
	fx, pp := fxPPSty(t, Monochrome)
	initLabelStyle := fx.CellsOf(pp.Styles).Trimmed()[0][0].Style
	selectProperty(t, PropertyNames[ReverseFgBgProperty], fx, pp)
	selectProperty(t, PropertyNames[ResetProperties], fx, pp)
	resetLabelStyle := fx.CellsOf(pp.Styles).Trimmed()[0][0].Style
	t.Eq(initLabelStyle.FG(), resetLabelStyle.FG())
	t.Eq(initLabelStyle.BG(), resetLabelStyle.BG())
}

func (s *property) Mono_reset_switches_back_to_zero_label(t *T) {
	fx, pp := fxPPSty(t, Monochrome)
	selectProperty(t, PropertyNames[ReverseFgBgProperty], fx, pp)
	selectProperty(t, PropertyNames[ResetProperties], fx, pp)
	t.Contains(fx.ScreenOf(pp).Trimmed(), PropertyNames[Properties])
}

func (s *property) Mono_attributes_selection_becomes_label(t *T) {
	fx, pp := fxPPSty(t, Monochrome)
	selectProperty(t, PropertyNames[StyleAttributeProperty], fx, pp)
	t.Contains(fx.ScreenOf(pp).Trimmed(),
		PropertyNames[StyleAttributeProperty])
}

func (s *property) Mono_zero_selection_resets_attributes_selection(t *T) {
	fx, pp := fxPPSty(t, Monochrome)
	selectProperty(t, PropertyNames[StyleAttributeProperty], fx, pp)
	t.Contains(
		fx.ScreenOf(extractDropDownLayer(t, fx, pp.Styles)),
		lines.StyleAttributeNames[lines.Bold],
	)
	fx.FireComponentClick(pp.Styles, 0, 0)
	fx.FireComponentClick(pp, 0, 0) // open selection
	fx.FireComponentClick(pp, 0, 0) // select zero
	t.Not.Contains(
		fx.ScreenOf(extractDropDownLayer(t, fx, pp.Styles)),
		lines.StyleAttributeNames[lines.Bold],
	)
}

func (s *property) Has_no_zero_label_if_not_monochrome(t *T) {
	_, pp := fxPPSty(t, System8)
	t.Eq("", pp.ZeroLabel)
}

func (s *property) Label_defaults_to_foreground_prop_if_not_mono(t *T) {
	fx, pp := fxPPSty(t, System8)
	t.Contains(fx.ScreenOf(pp), PropertyNames[ForegroundProperty])
}

func (s *property) Offers_background_property_if_not_monochrome(t *T) {
	fx, pp := fxPPSty(t, System8)
	t.Contains(
		fx.ScreenOf(extractDropDownLayer(t, fx, pp)),
		PropertyNames[BackgroundProperty],
	)
}

func (s *property) Switches_styles_to_bg_color_selection(t *T) {
	fx, pp := fxPPSty(t, System8)
	t.True(pp.Styles.SelectingForeground())
	selectProperty(t, PropertyNames[BackgroundProperty], fx, pp)
	t.True(pp.Styles.SelectingBackground())
	t.Not.True(pp.Styles.SelectingForeground())
}

func (s *property) Switches_styles_to_fg_color_selection(t *T) {
	fx, pp := fxPPSty(t, System8)
	t.True(pp.Styles.SelectingForeground())
	selectProperty(t, PropertyNames[BackgroundProperty], fx, pp)
	selectProperty(t, PropertyNames[ForegroundProperty], fx, pp)
	t.True(pp.Styles.SelectingForeground())
	t.Not.True(pp.Styles.SelectingBackground())
}

func (s *property) Switches_styles_to_attribute_selection(t *T) {
	fx, pp := fxPPSty(t, System8)
	t.True(pp.Styles.SelectingForeground())
	selectProperty(t, PropertyNames[StyleAttributeProperty], fx, pp)
	t.Not.True(pp.Styles.SelectingForeground())
	t.True(pp.Styles.SelectingStyleAttributes())
}

func (s *property) Colored_invert_switches_fg_and_bg(t *T) {
	fx, pp := fxPPSty(t, System8)
	stylesLabelStyle := fx.CellsOf(pp.Styles).Trimmed()[0][0].Style
	selectProperty(t, PropertyNames[ReverseFgBgProperty], fx, pp)
	invertedLabelStyle := fx.CellsOf(pp.Styles).Trimmed()[0][0].Style
	t.Eq(stylesLabelStyle.FG(), invertedLabelStyle.BG())
	t.Eq(stylesLabelStyle.BG(), invertedLabelStyle.FG())
	selectProperty(t, PropertyNames[ReverseFgBgProperty], fx, pp)
	invertedLabelStyle = fx.CellsOf(pp.Styles).Trimmed()[0][0].Style
	t.Eq(stylesLabelStyle.FG(), invertedLabelStyle.FG())
	t.Eq(stylesLabelStyle.BG(), invertedLabelStyle.BG())
}

func (s *property) Colored_reset_sets_initial_style(t *T) {
	fx, pp := fxPPSty(t, System8)
	initLabelStyle := fx.CellsOf(pp.Styles).Trimmed()[0][0].Style
	selectProperty(t, PropertyNames[ReverseFgBgProperty], fx, pp)
	selectProperty(t, PropertyNames[ResetProperties], fx, pp)
	resetLabelStyle := fx.CellsOf(pp.Styles).Trimmed()[0][0].Style
	t.Eq(initLabelStyle.FG(), resetLabelStyle.FG())
	t.Eq(initLabelStyle.BG(), resetLabelStyle.BG())
}

func TestPropertySelect(t *testing.T) {
	t.Parallel()
	Run(&property{}, t)
}

type styles struct{ Suite }

func (s *styles) Colors_default_to_monochrome(t *T) {
	t.Eq(Monochrome, (Styles{}).Colors)
}

func (s *styles) Monochrome_defaults_to_black_bg_and_white_fg(t *T) {
	mnc := lines.NewStyle(lines.ZeroStyle, lines.White, lines.Black)
	fx := fx.New(t, &Styles{})
	t.Eq(mnc, fx.Root().(*Styles).Value())
}

func (s *styles) Monochrome_defaults_to_zero_style_attributes(t *T) {
	mnc := lines.NewStyle(lines.ZeroStyle, lines.White, lines.Black)
	fx := fx.New(t, &Styles{})
	t.Eq(mnc, fx.Root().(*Styles).Value())
}

func (s *styles) Monochrome_label_defaults_to_color_range_name(t *T) {
	fx := fx.New(t, &Styles{})
	t.Contains(
		fx.ScreenOf(fx.Root()).Trimmed(),
		RangeNames[Monochrome],
	)
}

func (s *styles) Monochrome_prints_its_label_in_set_style(t *T) {
	mnc := lines.NewStyle(lines.ZeroStyle, lines.White, lines.Black)
	fx := fx.New(t, &Styles{})
	cc := fx.CellsOf(fx.Root()).Trimmed()[0]
	cc = cc[:len(cc)-2]
	for _, c := range cc {
		t.FatalIfNot(t.Eq(mnc, c.Style))
	}
}

func (s *styles) Attributes_use_selected_style_for_zero_label(t *T) {
	fx, pp := fxPPSty(t, Monochrome)
	selectProperty(t, PropertyNames[StyleAttributeProperty], fx, pp)
	got := fx.CellsOf(pp.Styles).Trimmed()[0][0].Style
	t.Eq(pp.Styles.Value(), got)
}

func (s *styles) Attributes_have_color_range_name_as_zero_label(t *T) {
	fx, pp := fxPPSty(t, Monochrome)
	selectProperty(t, PropertyNames[StyleAttributeProperty], fx, pp)
	t.Eq(pp.Styles.ZeroLabel, RangeNames[Monochrome])
}

// selectStyleAspects opens given styles drop down component cmp's
// selection list and selects given style aspect a in given terminal
// test fixture fx.  It ends test execution if a is not found.
func selectStyleAspect(t *T, a string, fx *lines.Fixture, cmp *Styles) {
	selectModalProperty(t, a, fx, extractDropDownLayer(t, fx, cmp))
}

func (s *styles) Value_has_style_attr_added_on_attr_selection(t *T) {
	fx, pp := fxPPSty(t, Monochrome)
	t.Not.True(pp.Styles.Value().AA()&lines.Bold == lines.Bold)
	selectProperty(t, PropertyNames[StyleAttributeProperty], fx, pp)
	selectStyleAspect(
		t, lines.StyleAttributeNames[lines.Bold], fx, pp.Styles)
	t.True(pp.Styles.Value().AA()&lines.Bold == lines.Bold)
}

func (s *styles) Attr_selection_is_a_noop_if_no_attr_selected(t *T) {
	fx, pp := fxPPSty(t, Monochrome)
	selectProperty(t, PropertyNames[StyleAttributeProperty], fx, pp)
	lyr := extractDropDownLayer(t, fx, pp.Styles)
	t.Contains(fx.Screen(), lines.StyleAttributeNames[lines.Bold])
	exp := pp.Styles.Value().AA()
	x, y, _, _ := fx.Dim(lyr).Printable()
	fx.FireClick(x, y-1)
	t.Not.Contains(fx.Screen(), lines.StyleAttributeNames[lines.Bold])
	t.Eq(exp, pp.Styles.Value().AA())
}

func (s *styles) Zero_attribute_label_has_an_attr_change_applied(t *T) {
	fx, pp := fxPPSty(t, Monochrome)
	t.Not.True(pp.Styles.Value().AA()&lines.Bold == lines.Bold)
	selectProperty(t, PropertyNames[StyleAttributeProperty], fx, pp)
	aaBefore := fx.CellsOf(pp.Styles)[0][0].Style.AA()
	t.Eq(pp.Styles.Value().AA(), aaBefore)
	selectStyleAspect(
		t, lines.StyleAttributeNames[lines.Bold], fx, pp.Styles)
	aaAfter := fx.CellsOf(pp.Styles)[0][0].Style.AA()
	t.Eq(pp.Styles.Value().AA(), aaAfter)
	t.Not.Eq(aaBefore, aaAfter)
}

func fxAttributes(t *T, timeout ...time.Duration) (
	*lines.Fixture, *StyleProperty, *ModalList,
) {
	fx, pp := fxPPSty(t, Monochrome)
	selectProperty(t, PropertyNames[StyleAttributeProperty], fx, pp)
	selectStyleAspect(
		t, lines.StyleAttributeNames[lines.Bold], fx, pp.Styles)
	return fx, pp, extractDropDownLayer(t, fx, pp.Styles)
}

func (s *styles) Attributes_list_items_have_no_style_attributes(t *T) {
	fx, _, lyr := fxAttributes(t)
	t.Eq(lines.ZeroStyle, fx.CellsOf(lyr)[0][0].Style.AA())
}

func (s *styles) Selected_attributes_are_marked_in_item_list(t *T) {
	fx, _, lyr := fxAttributes(t)
	t.True(strings.HasSuffix(
		string(fx.ScreenOf(lyr)[0]),
		SelectedMark,
	))
}

func (s *styles) Deselected_attributes_are_unmarked_in_item_list(t *T) {
	fx, pp, lyr := fxAttributes(t)
	t.True(strings.HasSuffix(
		string(fx.ScreenOf(lyr)[0]),
		SelectedMark,
	))
	fx.FireComponentClick(pp.Styles, 0, 0)
	selectStyleAspect(
		t, lines.StyleAttributeNames[lines.Bold], fx, pp.Styles)
	lyr = extractDropDownLayer(t, fx, pp.Styles)
	t.Not.True(strings.HasSuffix(
		string(fx.ScreenOf(lyr)[0]),
		SelectedMark,
	))
}

func (s *styles) Zero_label_is_initially_color_range_name(t *T) {
	_, pp := fxPPSty(t, System8)
	t.Eq(RangeNames[System8], pp.Styles.ZeroLabel)
}

func (s *styles) Zero_label_has_initially_initial_style(t *T) {
	fx, pp := fxPPSty(t, System8)
	got := fx.CellsOf(pp.Styles)[0][0].Style
	t.Eq(pp.Styles.Value(), got)
}

func (s *styles) Fg_selection_offers_all_colors_except_for_bg(t *T) {
	fx, pp := fxPPSty(t, System8)
	lyr := extractDropDownLayer(t, fx, pp.Styles)
	bg := pp.Styles.Value().BG()
	lyrScreen := fx.ScreenOf(lyr).Trimmed().String()
	t.Not.Contains(lyrScreen, lines.ColorNames[bg])
	for _, c := range system8Colors {
		if c == bg {
			continue
		}
		t.Contains(lyrScreen, lines.ColorNames[c])
	}
}

func (s *styles) Fg_selection_has_currently_set_color_Marked(t *T) {
	fx, pp := fxPPSty(t, System8)
	lyrScreen := fx.ScreenOf(extractDropDownLayer(t, fx, pp.Styles))
	fgName := lines.ColorNames[pp.Styles.Value().FG()]
	for _, l := range lyrScreen {
		if !strings.Contains(l, fgName) {
			continue
		}
		t.True(strings.HasSuffix(l, SelectedMark))
	}
}

func (s *styles) Fg_selection_items_are_colored_with_selectable_colors(
	t *T,
) {
	fx, pp := fxPPSty(t, System8)
	lyrCells := fx.CellsOf(extractDropDownLayer(t, fx, pp.Styles))
	for _, l := range lyrCells {
		t.FatalIfNot(t.Not.Eq("", lines.ColorNames[l[0].Style.FG()]))
		t.FatalIfNot(t.Contains(l, lines.ColorNames[l[0].Style.FG()]))
	}
}

func (s *styles) Fg_selection_items_have_constant_bg_color(t *T) {
	fx, pp := fxPPSty(t, System8)
	expBG := pp.Styles.Value().BG()
	lyrCells := fx.CellsOf(extractDropDownLayer(t, fx, pp.Styles))
	for _, l := range lyrCells {
		t.FatalIfNot(t.Eq(expBG, l[0].Style.BG()))
	}
}

func (s *styles) Fg_selection_is_set_as_value_s_fg_color(t *T) {
	fx, pp := fxPPSty(t, System8)
	lyr := extractDropDownLayer(t, fx, pp.Styles)
	expSty := fx.CellsOf(lyr)[0][0].Style
	t.Not.Eq(expSty.FG(), pp.Styles.Value().FG())
	fx.FireComponentClick(lyr, 0, 0)
	t.Eq(expSty.FG(), pp.Styles.Value().FG())
}

func (s *styles) Fg_selection_updates_selection_color_mark(t *T) {
	fx, pp := fxPPSty(t, System8)
	lyr := extractDropDownLayer(t, fx, pp.Styles)
	initMark := func() int {
		for idx, l := range fx.ScreenOf(lyr) {
			if !strings.HasSuffix(l, SelectedMark) {
				continue
			}
			return idx
		}
		t.Fatal("no initially marked foreground color")
		return -1
	}()
	expSty := fx.CellsOf(lyr)[0][0].Style
	t.Not.Eq(expSty.FG(), pp.Styles.Value().FG())
	fx.FireComponentClick(lyr, 0, 0)
	lyr = extractDropDownLayer(t, fx, pp.Styles)
	t.Not.True(strings.HasSuffix(
		fx.ScreenOf(lyr)[initMark], SelectedMark))
	t.True(strings.HasSuffix(
		fx.ScreenOf(lyr)[0], SelectedMark))
}

func (s *styles) Bg_selection_items_are_colored_with_selectable_colors(
	t *T,
) {
	fx, pp := fxPPSty(t, System8)
	selectProperty(t, PropertyNames[BackgroundProperty], fx, pp)
	lyrCells := fx.CellsOf(extractDropDownLayer(t, fx, pp.Styles))
	for _, l := range lyrCells {
		t.FatalIfNot(t.Not.Eq("", lines.ColorNames[l[0].Style.BG()]))
		t.FatalIfNot(t.Contains(l, lines.ColorNames[l[0].Style.BG()]))
	}
}

func (s *styles) Bg_selection_has_currently_set_color_Marked(t *T) {
	fx, pp := fxPPSty(t, System8)
	selectProperty(t, PropertyNames[BackgroundProperty], fx, pp)
	lyrScreen := fx.ScreenOf(extractDropDownLayer(t, fx, pp.Styles))
	bgName := lines.ColorNames[pp.Styles.Value().BG()]
	for _, l := range lyrScreen {
		if !strings.Contains(l, bgName) {
			continue
		}
		t.True(strings.HasSuffix(l, SelectedMark))
	}
}

func (s *styles) Bg_selection_items_have_constant_bg_color(t *T) {
	fx, pp := fxPPSty(t, System8)
	selectProperty(t, PropertyNames[BackgroundProperty], fx, pp)
	expFG := pp.Styles.Value().FG()
	lyrCells := fx.CellsOf(extractDropDownLayer(t, fx, pp.Styles))
	for _, l := range lyrCells {
		t.FatalIfNot(t.Eq(expFG, l[0].Style.FG()))
	}
}

func (s *styles) Bg_selection_is_set_as_value_s_bg_color(t *T) {
	fx, pp := fxPPSty(t, System8)
	selectProperty(t, PropertyNames[BackgroundProperty], fx, pp)
	lyr := extractDropDownLayer(t, fx, pp.Styles)
	expSty := fx.CellsOf(lyr)[1][0].Style
	t.Not.Eq(expSty.BG(), pp.Styles.Value().BG())
	fx.FireComponentClick(lyr, 0, 1)
	t.Eq(expSty.BG(), pp.Styles.Value().BG())
}

func (s *styles) Bg_selection_updates_selection_color_mark(t *T) {
	fx, pp := fxPPSty(t, System8)
	selectProperty(t, PropertyNames[BackgroundProperty], fx, pp)
	lyr := extractDropDownLayer(t, fx, pp.Styles)
	initMark := func() int {
		for idx, l := range fx.ScreenOf(lyr) {
			if !strings.HasSuffix(l, SelectedMark) {
				continue
			}
			return idx
		}
		t.Fatal("no initially marked foreground color")
		return -1
	}()
	expSty := fx.CellsOf(lyr)[1][0].Style
	t.Not.Eq(expSty.BG(), pp.Styles.Value().BG())
	fx.FireComponentClick(lyr, 0, 1)
	lyr = extractDropDownLayer(t, fx, pp.Styles)
	t.Not.True(strings.HasSuffix(
		fx.ScreenOf(lyr)[initMark], SelectedMark))
	t.True(strings.HasSuffix(
		fx.ScreenOf(lyr)[1], SelectedMark))
}

func (s *styles) Attr_selection_items_have_style_attrs_labels(t *T) {
	fx, pp := fxPPSty(t, System8)
	selectProperty(t, PropertyNames[StyleAttributeProperty], fx, pp)
	lyrScr := fx.ScreenOf(extractDropDownLayer(t, fx, pp.Styles)).String()
	for _, a := range aa {
		t.Contains(lyrScr, lines.StyleAttributeNames[a])
	}
}

func (s *styles) Attr_selection_has_inverse_value_colors_as_style(t *T) {
	fx, pp := fxPPSty(t, System8)
	selectProperty(t, PropertyNames[StyleAttributeProperty], fx, pp)
	lyrCll := fx.CellsOf(extractDropDownLayer(t, fx, pp.Styles))
	expSty := pp.Styles.Value().Invert().WithAA(lines.ZeroStyle)
	for _, l := range lyrCll {
		t.Eq(expSty, l[0].Style)
	}
}

func TestStyles(t *testing.T) {
	t.Parallel()
	Run(&styles{}, t)
}
