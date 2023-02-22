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

func (s *property) Has_properties_as_zero_label(t *T) {
	fx := fx.New(t, &StyleProperty{Styles: &Styles{}})
	t.Contains(fx.ScreenOf(fx.Root()), PropertyNames[Properties])
}

func (s *property) Offers_background_property(t *T) {
	fx := fx.New(t, &StyleProperty{})
	lry := extractDropDownLayer(t, fx, fx.Root().(*StyleProperty))
	t.Contains(fx.ScreenOf(lry), PropertyNames[BackgroundProperty])
}

func (s *property) Offers_background_foreground_reverse(t *T) {
	fx := fx.New(t, &StyleProperty{})
	lry := extractDropDownLayer(t, fx, fx.Root().(*StyleProperty))
	t.Contains(fx.ScreenOf(lry), PropertyNames[ReverseFgBgProperty])
}

func (s *property) Offers_style_attributes(t *T) {
	fx := fx.New(t, &StyleProperty{})
	lry := extractDropDownLayer(t, fx, fx.Root().(*StyleProperty))
	t.Contains(fx.ScreenOf(lry), PropertyNames[StyleAttributeProperty])
}

func (s *property) Offers_reset_properties(t *T) {
	fx := fx.New(t, &StyleProperty{})
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

func newPPSty(t *T, timeout ...time.Duration) (
	*lines.Fixture, *StyleProperty,
) {
	cmp, pp, ss := &ppSty{}, &StyleProperty{}, &Styles{}
	pp.Styles = ss
	cmp.CC = append(cmp.CC, pp, ss)
	fx := fx.New(t, cmp, timeout...)
	return fx, pp
}

func (c *ppSty) OnInit(e *lines.Env) { c.Dim().SetHeight(1) }

func (s *property) Mono_reverse_switches_fg_and_bg(t *T) {
	fx, pp := newPPSty(t)
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
	fx, pp := newPPSty(t)
	selectProperty(t, PropertyNames[ReverseFgBgProperty], fx, pp)
	t.Contains(fx.ScreenOf(pp).Trimmed(), PropertyNames[Properties])
}

func (s *property) Mono_reset_sets_initial_style(t *T) {
	fx, pp := newPPSty(t)
	initLabelStyle := fx.CellsOf(pp.Styles).Trimmed()[0][0].Style
	selectProperty(t, PropertyNames[ReverseFgBgProperty], fx, pp)
	selectProperty(t, PropertyNames[ResetProperties], fx, pp)
	resetLabelStyle := fx.CellsOf(pp.Styles).Trimmed()[0][0].Style
	t.Eq(initLabelStyle.FG(), resetLabelStyle.FG())
	t.Eq(initLabelStyle.BG(), resetLabelStyle.BG())
}

func (s *property) Mono_reset_switches_back_to_zero_label(t *T) {
	fx, pp := newPPSty(t)
	selectProperty(t, PropertyNames[ReverseFgBgProperty], fx, pp)
	selectProperty(t, PropertyNames[ResetProperties], fx, pp)
	t.Contains(fx.ScreenOf(pp).Trimmed(), PropertyNames[Properties])
}

func (s *property) Mono_attributes_selection_becomes_label(t *T) {
	fx, pp := newPPSty(t)
	selectProperty(t, PropertyNames[StyleAttributeProperty], fx, pp)
	t.Contains(fx.ScreenOf(pp).Trimmed(),
		PropertyNames[StyleAttributeProperty])
}

func (s *property) Mono_zero_selection_resets_attributes_selection(t *T) {
	fx, pp := newPPSty(t)
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

func TestPropertySelect(t *testing.T) {
	t.Parallel()
	Run(&property{}, t)
}

type styles struct{ Suite }

func (s *styles) Colors_default_to_monochrome(t *T) {
	t.Eq(Monochrome, (Styles{}).Colors)
}

func (s *styles) Monochrome_default_to_black_bg_and_white_fg(t *T) {
	mnc := lines.NewStyle(lines.ZeroStyle, lines.White, lines.Black)
	fx := fx.New(t, &Styles{})
	t.Eq(mnc, fx.Root().(*Styles).Value())
}

func (s *styles) Monochrome_defaults_to_zero_style_attributes(t *T) {
	mnc := lines.NewStyle(lines.ZeroStyle, lines.White, lines.Black)
	fx := fx.New(t, &Styles{})
	t.Eq(mnc, fx.Root().(*Styles).Value())
}

func (s *styles) Monochrome_label_defaults_to_initial_fg_color_name(t *T) {
	mnc := lines.NewStyle(lines.ZeroStyle, lines.White, lines.Black)
	fx := fx.New(t, &Styles{})
	t.Contains(
		fx.ScreenOf(fx.Root()).Trimmed(),
		lines.ColorNames[mnc.FG()],
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
	fx, pp := newPPSty(t)
	selectProperty(t, PropertyNames[StyleAttributeProperty], fx, pp)
	got := fx.CellsOf(pp.Styles).Trimmed()[0][0].Style
	t.Eq(pp.Styles.Value(), got)
	t.Contains(fx.ScreenOf(pp.Styles), lines.ColorNames[got.FG()])
	t.Eq(pp.Styles.ZeroLabel, lines.ColorNames[got.FG()])
}

func (s *styles) Attributes_have_fg_color_name_as_zero_label(t *T) {
	fx, pp := newPPSty(t)
	selectProperty(t, PropertyNames[StyleAttributeProperty], fx, pp)
	t.Eq(pp.Styles.ZeroLabel, lines.ColorNames[pp.Styles.value.FG()])
}

// selectStyleAspects opens given styles drop down component cmp's
// selection list and selects given style aspect a in given terminal
// test fixture fx.  It ends test execution if a is not found.
func selectStyleAspect(t *T, a string, fx *lines.Fixture, cmp *Styles) {
	selectModalProperty(t, a, fx, extractDropDownLayer(t, fx, cmp))
}

func (s *styles) Value_has_style_attr_added_on_attr_selection(t *T) {
	fx, pp := newPPSty(t)
	t.Not.True(pp.Styles.Value().AA()&lines.Bold == lines.Bold)
	selectProperty(t, PropertyNames[StyleAttributeProperty], fx, pp)
	selectStyleAspect(
		t, lines.StyleAttributeNames[lines.Bold], fx, pp.Styles)
	t.True(pp.Styles.Value().AA()&lines.Bold == lines.Bold)
}

func (s *styles) Attr_selection_is_a_noop_if_no_attr_selected(t *T) {
	fx, pp := newPPSty(t)
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
	fx, pp := newPPSty(t)
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
	fx, pp := newPPSty(t)
	selectProperty(t, PropertyNames[StyleAttributeProperty], fx, pp)
	selectStyleAspect(
		t, lines.StyleAttributeNames[lines.Bold], fx, pp.Styles)
	return fx, pp, extractDropDownLayer(t, fx, pp.Styles)
}

func (s *styles) Attributes_list_items_have_no_style_attributes(t *T) {
	fx, _, lyr := fxAttributes(t)
	t.Eq(lines.ZeroStyle, fx.CellsOf(lyr)[0][0].Style.AA())
}

func (s *styles) Selected_attributes_are_flagged_in_item_list(t *T) {
	fx, _, lyr := fxAttributes(t)
	t.True(strings.HasSuffix(
		string(fx.ScreenOf(lyr)[0]),
		SelectedAttr,
	))
}

func (s *styles) Deselected_attributes_are_unflagged_in_item_list(t *T) {
	fx, pp, lyr := fxAttributes(t)
	t.True(strings.HasSuffix(
		string(fx.ScreenOf(lyr)[0]),
		SelectedAttr,
	))
	fx.FireComponentClick(pp.Styles, 0, 0)
	selectStyleAspect(
		t, lines.StyleAttributeNames[lines.Bold], fx, pp.Styles)
	lyr = extractDropDownLayer(t, fx, pp.Styles)
	t.Not.True(strings.HasSuffix(
		string(fx.ScreenOf(lyr)[0]),
		SelectedAttr,
	))
}

func TestStyles(t *testing.T) {
	t.Parallel()
	Run(&styles{}, t)
}
