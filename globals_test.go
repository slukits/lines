// Copyright (c) 2022 Stephan Lukits. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package lines

import (
	"testing"

	. "github.com/slukits/gounit"
)

// _globals is an white box test of implementation details hence we
// might replace it by time with according black box tests.
type _globals struct{ Suite }

func SetUp(t *T) { t.Parallel() }

func (*_globals) Tab_width_defaults_to_four(t *T) {
	gg := newGlobals(nil)
	t.Eq(4, gg.TabWidth())
}

func (*_globals) Ignore_non_positive_tab_width_update(t *T) {
	gg := newGlobals(nil).SetTabWidth(0)
	t.Eq(4, gg.TabWidth())
}

type globalsFX struct{ gg *Globals }

func (gg *globalsFX) globals() *Globals { return gg.gg }

func (*_globals) Propagate_updated_tab_width(t *T) {
	var g1, g2 globaler
	gg := newGlobals(func(f func(globaler)) {
		f(g1)
		f(g2)
	})
	g1, g2 = &globalsFX{gg: gg.clone()}, &globalsFX{gg: gg.clone()}
	gg.SetTabWidth(3)
	t.Eq(3, g1.globals().TabWidth())
	t.Eq(3, g2.globals().TabWidth())
}

func (*_globals) Ignore_propagated_tab_width_which_has_been_set(t *T) {
	var g1, g2 globaler
	gg := newGlobals(func(f func(globaler)) {
		f(g1)
		f(g2)
	})
	g1, g2 = &globalsFX{gg: gg.clone()}, &globalsFX{gg: gg.clone()}
	g2.globals().SetTabWidth(6)
	gg.SetTabWidth(3)
	t.Eq(3, g1.globals().TabWidth())
	t.Eq(6, g2.globals().TabWidth())
}

func (*_globals) Report_a_tab_width_update(t *T) {
	exp, got := globalTabWidth, globalsUpdates(0)
	gg := newGlobals(nil).SetUpdateListener(
		func(gu globalsUpdates, st StyleType, gsu globalStyleUpdates) {
			got = gu
		})
	t.Not.Eq(exp, got)
	gg.SetTabWidth(8)
	t.Eq(exp, got)
}

func (*_globals) Report_a_tab_width_update_through_propagation(t *T) {
	var g1 globaler
	gg := newGlobals(func(f func(globaler)) { f(g1) })
	g1 = &globalsFX{gg: gg.clone()}
	exp, got := globalTabWidth, globalsUpdates(0)
	g1.globals().SetUpdateListener(
		func(gu globalsUpdates, st StyleType, gsu globalStyleUpdates) {
			got = gu
		})
	t.Not.Eq(exp, got)
	gg.SetTabWidth(8)
	t.Eq(exp, got)
}

// Default_Styles_default_to_default_styles the later is here
// the *DefaultStyle*-variable.
func (*_globals) Default_styles_defaults_to_default_style(t *T) {
	gg := newGlobals(nil)
	t.Eq(DefaultStyle, gg.Style(Default))
	gg = &Globals{}
	t.Eq(DefaultStyle.AA(), gg.AA(Default))
	t.Eq(DefaultStyle.FG(), gg.FG(Default))
	t.Eq(DefaultStyle.BG(), gg.BG(Default))
	gg.ss = map[StyleType]Style{}
	t.Eq(DefaultStyle.AA(), gg.AA(Default))
	t.Eq(DefaultStyle.FG(), gg.FG(Default))
	t.Eq(DefaultStyle.BG(), gg.BG(Default))
}

func (*_globals) Propagate_updated_style(t *T) {
	var g1, g2 globaler
	gg := newGlobals(func(f func(globaler)) {
		f(g1)
		f(g2)
	})
	g1, g2 = &globalsFX{gg: gg.clone()}, &globalsFX{gg: gg.clone()}
	exp := DefaultStyle.WithAA(Dim).WithFG(Yellow).WithBG(Blue)
	gg.SetStyle(Default, exp)
	t.Eq(g1.globals().Style(Default), exp)
	t.Eq(g2.globals().Style(Default), exp)
}

func (*_globals) Ignore_propagated_style_which_has_been_set(t *T) {
	var g1, g2 globaler
	gg := newGlobals(func(f func(globaler)) {
		f(g1)
		f(g2)
	})
	g1, g2 = &globalsFX{gg: gg.clone()}, &globalsFX{gg: gg.clone()}
	exp := DefaultStyle.WithAA(Dim).WithFG(Yellow).WithBG(Blue)
	other := DefaultStyle.WithAA(Blink).WithFG(Green).WithBG(Red)
	g2.globals().SetStyle(Default, other)
	gg.SetStyle(Default, exp)
	t.Eq(g1.globals().Style(Default), exp)
	t.Not.Eq(g2.globals().Style(Default), exp)
	t.Eq(g2.globals().Style(Default), other)
}

func (*_globals) Report_a_style_update(t *T) {
	exp, got := Default, StyleType(0)
	gg := newGlobals(nil).SetUpdateListener(
		func(gu globalsUpdates, st StyleType, gsu globalStyleUpdates) {
			got = st
		})
	sty := DefaultStyle.WithAA(Dim).WithFG(Yellow).WithBG(Blue)
	gg.SetStyle(Default, sty)
	t.Eq(exp, got)
}

func (*_globals) Report_a_style_update_through_propagation(t *T) {
	var g1 globaler
	gg := newGlobals(func(f func(globaler)) { f(g1) })
	g1 = &globalsFX{gg: gg.clone()}
	exp, got := Default, StyleType(0)
	g1.globals().SetUpdateListener(
		func(gu globalsUpdates, st StyleType, gsu globalStyleUpdates) {
			got = st
		})
	sty := DefaultStyle.WithAA(Dim).WithFG(Yellow).WithBG(Blue)
	gg.SetStyle(Default, sty)
	t.Eq(exp, got)
}

func (*_globals) Propagate_updated_style_aspect(t *T) {
	var g1, g2 globaler
	gg := newGlobals(func(f func(globaler)) {
		f(g1)
		f(g2)
	})
	g1, g2 = &globalsFX{gg: gg.clone()}, &globalsFX{gg: gg.clone()}
	exp := DefaultStyle.WithAA(Dim).WithFG(Yellow).WithBG(Blue)
	gg.SetAA(Default, exp.AA())
	t.Eq(exp.AA(), g1.globals().Style(Default).AA())
	t.Eq(exp.AA(), g2.globals().Style(Default).AA())
	gg.SetFG(Default, exp.FG())
	t.Eq(exp.FG(), g1.globals().Style(Default).FG())
	t.Eq(exp.FG(), g2.globals().Style(Default).FG())
	gg.SetBG(Default, exp.BG())
	t.Eq(exp.BG(), g1.globals().Style(Default).BG())
	t.Eq(exp.BG(), g2.globals().Style(Default).BG())
}

func (*_globals) Ignore_propagated_style_aspect_which_has_been_set(t *T) {
	var g1, g2 globaler
	gg := newGlobals(func(f func(globaler)) {
		f(g1)
		f(g2)
	})
	g1, g2 = &globalsFX{gg: gg.clone()}, &globalsFX{gg: gg.clone()}
	exp := DefaultStyle.WithAA(Dim).WithFG(Yellow).WithBG(Blue)
	other := DefaultStyle.WithAA(Blink).WithFG(Green).WithBG(Red)
	g2.globals().SetAA(Highlight, other.AA())
	gg.SetAA(Highlight, exp.AA())
	t.Eq(exp.AA(), g1.globals().AA(Highlight))
	t.Not.Eq(exp.AA(), g2.globals().Style(Highlight).AA())
	t.Eq(other.AA(), g2.globals().AA(Highlight))
	g2.globals().SetFG(Highlight, other.FG())
	gg.SetFG(Highlight, exp.FG())
	t.Eq(exp.FG(), g1.globals().Style(Highlight).FG())
	t.Not.Eq(exp.FG(), g2.globals().FG(Highlight))
	t.Eq(other.FG(), g2.globals().FG(Highlight))
	g2.globals().SetBG(Highlight, other.BG())
	gg.SetBG(Highlight, exp.BG())
	t.Eq(exp.BG(), g1.globals().BG(Highlight))
	t.Not.Eq(exp.BG(), g2.globals().Style(Highlight).BG())
	t.Eq(other.BG(), g2.globals().BG(Highlight))
}

func (*_globals) Highlight_maps_default_to_highlight_style(t *T) {
	gg := newGlobals(nil)
	t.Eq(gg.Highlight(gg.Style(Default)), gg.Style(Highlight))
}

func (*_globals) Propagate_set_highlighter(t *T) {
	var g1, g2 globaler
	gg := newGlobals(func(f func(globaler)) {
		f(g1)
		f(g2)
	})
	g1, g2 = &globalsFX{gg: gg.clone()}, &globalsFX{gg: gg.clone()}
	h1 := g1.globals().Highlight(gg.Style(Default))
	h2 := g2.globals().Highlight(gg.Style(Default))
	t.Eq(h1, h2)

	sty := DefaultStyle.WithBG(Green).WithFG(Salmon1).WithAA(Blink)
	gg.SetHighlighter(func(_ Style) Style { return sty })
	h1 = g1.globals().Highlight(gg.Style(Default))
	h2 = g2.globals().Highlight(gg.Style(Default))
	t.Eq(h1, h2)
	t.Eq(h1, sty)
}

func (*_globals) Default_highlighter_uses_propagated_highlight_style(
	t *T,
) {
	var g globaler
	gg := newGlobals(func(f func(globaler)) { f(g) })
	g = &globalsFX{gg: gg.clone()}
	sty := DefaultStyle.WithBG(Green).WithFG(Salmon1).WithAA(Blink)
	gg.SetStyle(Highlight, sty)
	t.Eq(g.globals().Highlight(gg.Style(Default)), sty)
}

func (*_globals) Propagate_default_if_zero_highlighter_set(t *T) {
	var g globaler
	gg := newGlobals(func(f func(globaler)) { f(g) })
	g = &globalsFX{gg: gg.clone()}
	hiSty := DefaultStyle.WithBG(LightSlateBlue).WithFG(Black).WithAA(Dim)
	g.globals().SetStyle(Highlight, hiSty)
	h := g.globals().Highlight(gg.Style(Default))
	t.Eq(h, hiSty)

	sty := DefaultStyle.WithBG(Green).WithFG(Salmon1).WithAA(Blink)
	gg.SetHighlighter(func(_ Style) Style { return sty })
	t.Not.Eq(h, g.globals().Highlight(gg.Style(Default)))
	t.Eq(sty, g.globals().Highlight(gg.Style(Default)))

	gg.SetHighlighter(nil)
	t.Eq(h, g.globals().Highlight(gg.Style(Default)))
}

func (*_globals) Ignore_propagation_if_highlighter_has_been_set(t *T) {
	var g globaler
	gg := newGlobals(func(f func(globaler)) { f(g) })
	g = &globalsFX{gg: gg.clone()}
	hiSty := DefaultStyle.WithBG(LightSteelBlue).WithFG(Black).WithAA(Dim)
	g.globals().SetHighlighter(func(s Style) Style { return hiSty })
	h := g.globals().Highlight(gg.Style(Default))

	gg.SetHighlighter(nil)
	t.Eq(h, g.globals().Highlight(gg.Style(Default)))
}

func (*_globals) Report_highlighter_update(t *T) {
	hlReported := 0
	var g globaler
	gg := newGlobals(func(f func(globaler)) { f(g) })
	g = &globalsFX{gg: gg.clone()}
	gg.SetUpdateListener(
		func(gu globalsUpdates, st StyleType, gsu globalStyleUpdates) {
			if gu == globalHighlighter {
				hlReported++
			}
		},
	)
	g.globals().SetUpdateListener(
		func(gu globalsUpdates, st StyleType, gsu globalStyleUpdates) {
			if gu == globalHighlighter {
				hlReported++
			}
		},
	)

	hiSty := DefaultStyle.WithBG(LightSteelBlue).WithFG(Black).WithAA(Dim)
	gg.SetHighlighter(func(s Style) Style { return hiSty })
	t.Eq(2, hlReported)
	gg.SetHighlighter(nil)
	t.Eq(4, hlReported)
}

// func (*_globals) Report_a_tab_width_update(t *T) {
// 	exp, got := globalTabWidth, globalsUpdates(0)
// 	gg := newGlobals(nil).SetUpdateListener(
// 		func(gu globalsUpdates, st StyleType, gsu globalStyleUpdates) {
// 			got = gu
// 		})
// 	t.Not.Eq(exp, got)
// 	gg.SetTabWidth(8)
// 	t.Eq(exp, got)
// }
//
// func (*_globals) Report_a_tab_width_update_through_propagation(t *T) {
// 	var g1 globaler
// 	gg := newGlobals(func(f func(globaler)) { f(g1) })
// 	g1 = &globalsFX{gg: gg.clone()}
// 	exp, got := globalTabWidth, globalsUpdates(0)
// 	g1.globals().SetUpdateListener(
// 		func(gu globalsUpdates, st StyleType, gsu globalStyleUpdates) {
// 			got = gu
// 		})
// 	t.Not.Eq(exp, got)
// 	gg.SetTabWidth(8)
// 	t.Eq(exp, got)
// }

func TestGlobals(t *testing.T) {
	t.Parallel()
	Run(&_globals{}, t)
}
