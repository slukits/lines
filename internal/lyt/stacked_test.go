// Copyright (c) 2022 Stephan Lukits. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package lyt

import (
	"testing"

	. "github.com/slukits/gounit"
)

type stackerFX struct {
	Dimer
	dd []Dimer
}

func (sf *stackerFX) ForStacked(cb func(Dimer) (stop bool)) {
	for _, d := range sf.dd {
		if cb(d) {
			return
		}
	}
}

// Height returns the consumed hight in a layout of a stacker fixture.
func (sf *stackerFX) Height() int {
	return sf.Dim().layoutHeight()
}

// SumLayoutHeights returns the totally consumed layouted hight of a
// Stacker fixture's stacked Dimers, i.e. heights reduced by clippings
// respectively increased by margins.
func (sf *stackerFX) SumLayoutHeights() int {
	h := 0
	sf.ForStacked(func(d Dimer) (stop bool) {
		h += d.Dim().layoutHeight()
		return false
	})
	return h
}

// HasConsistentLayout is true iff every Stacker's Dimer's layouted
// width is the width of the Stacker and the sum of all Dimer's heights
// equals the height of the Stacker.
func (sf *stackerFX) HasConsistentLayout() (ok bool) {
	lHeightSum, ok := 0, true
	sf.ForStacked(func(d Dimer) (stop bool) {
		if d.Dim().layoutWidth() != sf.Dim().layoutWidth() {
			ok = false
			return true
		}
		lHeightSum += d.Dim().layoutHeight()
		return false
	})
	return sf.Dim().layoutHeight() == lHeightSum && ok
}

type stackerFactory struct{}

var sf = &stackerFactory{}

// New produces a new Stacker-implementation instance with default
// screen dimensions which provides given dimers.
func (sf *stackerFactory) New(dd ...Dimer) *stackerFX {
	return &stackerFX{
		Dimer: df.Screen(),
		dd:    dd,
	}
}

// Of produces a new Stacker-implementation instance wrapping given Dimer
// and providing the remaining dimers.
func (sf *stackerFactory) Of(d Dimer, dd ...Dimer) *stackerFX {
	return &stackerFX{
		Dimer: d,
		dd:    dd,
	}
}

// Filling produces a new Stacker-implementation instance which provides
// given dimers and has its fillsWidth/fillsHeight set to 1.
func (sf *stackerFactory) Filling(dd ...Dimer) *stackerFX {
	return &stackerFX{
		Dimer: df.FillingWH(1, 1),
		dd:    dd,
	}
}

// fxHF1 stacker with sole height filler.
var fxHF1 = func() *stackerFX { return sf.New(df.Filling()) }

// fxHF2 stacker with height filler at the beginning.
var fxHF2 = func() *stackerFX {
	return sf.New(df.Filling(), df.Fixed())
}

// fxHF3 stacker with height filler in between.
var fxHF3 = func() *stackerFX {
	return sf.New(df.Fixed(), df.Filling(), df.Fixed())
}

// fxHF4 stacker with height filler at the end.
var fxHF4 = func() *stackerFX {
	return sf.New(df.Fixed(), df.Filling())
}

type stacked struct{ Suite }

func (s *stacked) SetUp(t *T) { t.Parallel() }

func (s *stacked) Fails_if_it_has_fixed_dimer_with_zero_height(t *T) {
	fx := &Manager{Root: sf.New(df.Of(&Dim{width: 80, height: 25}))}
	t.FatalOn(fx.Reflow(nil))
	fx = &Manager{Root: sf.New(df.Of(&Dim{width: 80, height: 0}))}
	t.ErrIs(fx.Reflow(nil), ErrDim)
}

func (s *stacked) With_height_filler_consume_all_hight(t *T) {
	for _, fx := range []*stackerFX{fxHF1(), fxHF2(), fxHF3(), fxHF4()} {
		t.True(fx.Height() > fx.SumLayoutHeights())
		t.FatalOn((&Manager{Root: fx}).Reflow(nil))
		t.Eq(fx.Height(), fx.SumLayoutHeights())
	}
}

func (s *stacked) With_height_filler_have_no_margins_at_fixed(t *T) {
	fx := sf.New(df.Filling(), df.Fixed(), df.Filling(), df.Fixed(),
		df.Filling())
	for _, d := range fx.dd {
		d.Dim().mrgTop, d.Dim().mrgBottom = 1, 1
	}
	t.FatalOn((&Manager{Root: fx}).Reflow(nil))
	for _, d := range fx.dd {
		t.Eq(0, d.Dim().mrgTop)
		t.Eq(0, d.Dim().mrgBottom)
	}
}

func (s *stacked) With_underflowing_height_have_tb_margins(t *T) {
	// single margin (fx1), distributed margins (fx2), exact fit (fx3)
	fx1, fx2 := sf.New(df.Fixed()), sf.New(df.Fixed(), df.Fixed())
	fx3 := sf.New(df.Fixed(), df.Fixed(), df.FixedH(9))
	// []fixtures[]dimerInFixture[]marginTopBottom
	exp := [][][2]int{{{8, 9}}, {{3, 1}, {2, 3}},
		{{0, 0}, {0, 0}, {0, 0}}}
	for i, fx := range []*stackerFX{fx1, fx2, fx3} {
		t.FatalOn((&Manager{Root: fx}).Reflow(nil))
		for j, d := range fx.dd {
			mt, _, mb, _ := d.Dim().Margin()
			t.Eq(exp[i][j][0], mt)
			t.Eq(exp[i][j][1], mb)
		}
	}
}

func (s *stacked) With_underflowing_height_consume_all_height(t *T) {
	// single margin (fx1), distributed margins (fx2), exact fit (fx3)
	fx1, fx2 := sf.New(df.Fixed()), sf.New(df.Fixed(), df.Fixed())
	fx3 := sf.New(df.Fixed(), df.Fixed(), df.FixedH(9))
	for _, fx := range []*stackerFX{fx1, fx2, fx3} {
		t.FatalOn((&Manager{Root: fx}).Reflow(nil))
		t.Eq(fx.Height(), fx.SumLayoutHeights())
	}
}

func (s *stacked) With_overflowing_height_are_clipped(t *T) {
	// single overflow (fx1), last overflows (fx2), off-screen (fx3)
	fx1, fx2 := sf.New(df.FixedH(26)), sf.New(df.Fixed(), df.FixedH(18))
	fx3 := sf.New(df.Fixed(), df.FixedH(18), df.Fixed())
	fx4 := sf.New(df.FixedH(13), df.Filling(), df.Fixed(), df.Filling())
	t.FatalOn((&Manager{Root: fx1}).Reflow(nil))
	_, clippedHeight := fx1.dd[0].Dim().Clip()
	t.Eq(1, clippedHeight)
	t.FatalOn((&Manager{Root: fx2}).Reflow(nil))
	_, clippedHeight = fx2.dd[1].Dim().Clip()
	t.Eq(1, clippedHeight)
	t.FatalOn((&Manager{Root: fx3}).Reflow(nil))
	_, clippedHeight = fx3.dd[1].Dim().Clip()
	t.Eq(1, clippedHeight)
	t.True(fx3.dd[2].Dim().IsOffScreen())
	t.FatalOn((&Manager{Root: fx4}).Reflow(nil))
	_, clippedHeight = fx4.dd[3].Dim().Clip()
	t.Eq(2, clippedHeight)
}

func (s *stacked) With_fitting_layout_has_no_clipping(t *T) {
	fx := sf.New(df.Fixed(), df.Filling(), df.Fixed())
	for _, d := range fx.dd {
		d.Dim().clipHeight = 1
	}
	t.FatalOn((&Manager{Root: fx}).Reflow(nil))
	for _, d := range fx.dd {
		t.Eq(0, d.Dim().clipHeight)
	}
}

func (s *stacked) Height_filler_have_height_evenly_distributed(t *T) {
	// i.e. min-height 9 available 25 => 16/3: 6, 5, 5 => 9, 8, 8
	fx := sf.New(df.Filling(), df.Filling(), df.Filling())
	t.FatalOn((&Manager{Root: fx}).Reflow(nil))
	for i, d := range fx.dd {
		_, _, _, got := d.Dim().Area()
		switch i {
		case 0:
			t.Eq(9, got)
		default:
			t.Eq(8, got)
		}
	}
}

func (s *stacked) Width_fillers_have_no_left_or_right_margins(t *T) {
	fx := sf.New(df.Filling(), df.Filling(), df.Filling())
	for _, d := range fx.dd {
		d.Dim().mrgLeft, d.Dim().mrgRight = 1, 1
	}
	t.FatalOn((&Manager{Root: fx}).Reflow(nil))
	for _, d := range fx.dd {
		_, mr, _, ml := d.Dim().Margin()
		t.True(mr == ml && mr == 0)
	}
}

func (s *stacked) Width_fillers_consume_all_width(t *T) {
	fx := sf.New(df.Filling(), df.Filling(), df.Filling())
	_, _, exp, _ := fx.Dim().Area()
	t.FatalOn((&Manager{Root: fx}).Reflow(nil))
	for _, d := range fx.dd {
		_, _, got, _ := d.Dim().Area()
		t.Eq(exp, got)
	}
}

func (s *stacked) With_underflowing_width_have_lr_margins(t *T) {
	fx := sf.New(df.Fixed())
	t.FatalOn((&Manager{Root: fx}).Reflow(nil))
	_, gotRM, _, gotLM := fx.dd[0].Dim().Margin()
	t.True(gotRM == gotLM && gotLM == 30)
}

func (s *stacked) With_underflowing_width_consume_all_width(t *T) {
	fx := sf.New(df.Fixed())
	t.FatalOn((&Manager{Root: fx}).Reflow(nil))
	t.Eq(fx.Dim().width, fx.dd[0].Dim().layoutWidth())
}

func (s *stacked) With_overflowing_width_are_clipped_to_width(t *T) {
	fx := sf.New(df.FixedW(90), df.FillingW(90))
	t.FatalOn((&Manager{Root: fx}).Reflow(nil))
	gotFixed, _ := fx.dd[0].Dim().Clip()
	t.Eq(10, gotFixed)
	gotFilled, _ := fx.dd[1].Dim().Clip()
	t.Eq(10, gotFilled)
	t.Eq(fx.Dim().width, fx.dd[0].Dim().layoutWidth())
	t.Eq(fx.Dim().width, fx.dd[1].Dim().layoutWidth())
}

func TestStacked(t *testing.T) {
	t.Parallel()
	Run(&stacked{}, t)
}
