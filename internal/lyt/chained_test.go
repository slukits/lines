// Copyright (c) 2022 Stephan Lukits. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package lyt

import (
	"testing"

	. "github.com/slukits/gounit"
)

type chainerFX struct {
	Dimer
	dd []Dimer
}

func (cf *chainerFX) ForChained(cb func(Dimer) (stop bool)) {
	for _, d := range cf.dd {
		if cb(d) {
			return
		}
	}
}

// Width returns the consumed width in a layout of a chainer fixture.
func (cf *chainerFX) Width() int {
	return cf.Dim().layoutWidth()
}

// SumLayoutWidths returns the totally consumed layouted width of a
// Chainer fixture's chained Dimers, i.e. widths reduced by clippings
// respectively increased by margins.
func (cf *chainerFX) SumLayoutWidths() int {
	w := 0
	cf.ForChained(func(d Dimer) (stop bool) {
		w += d.Dim().layoutWidth()
		return false
	})
	return w
}

// HasConsistentLayout is true iff every Chainer's Dimer's layouted
// height is the layouted height of the Stacker and the sum of all
// Dimers' layouted widths equals the layouted width of the Stacker.
func (cf *chainerFX) HasConsistentLayout() (ok bool) {
	ok = true
	cf.ForChained(func(d Dimer) (stop bool) {
		if d.Dim().layoutHeight() != cf.Dim().layoutHeight() {
			ok = false
			return true
		}
		return false
	})
	return cf.Dim().layoutWidth() == cf.SumLayoutWidths() && ok
}

type chainerFactory struct{}

var cf = &chainerFactory{}

// New produces a new Chainer-implementation instance with default
// screen dimensions which provides given dimers.
func (cf *chainerFactory) New(dd ...Dimer) *chainerFX {
	return &chainerFX{
		Dimer: df.Screen(),
		dd:    dd,
	}
}

// Of produces a new Stacker-implementation instance wrapping given Dimer
// and providing the remaining dimers.
func (sf *chainerFactory) Of(d Dimer, dd ...Dimer) *stackerFX {
	return &stackerFX{
		Dimer: d,
		dd:    dd,
	}
}

func (cf *chainerFactory) Filling(dd ...Dimer) *chainerFX {
	return &chainerFX{
		Dimer: df.FillingWH(1, 1),
		dd:    dd,
	}
}

type chained struct{ Suite }

func (s *chained) SetUp(t *T) { t.Parallel() }

func (s *chained) Fails_if_it_has_fixed_dimer_with_zero_width(t *T) {
	fx := &Manager{cf.New(df.Of(&Dim{width: 80, height: 25}))}
	t.FatalOn(fx.Reflow(nil))
	fx = &Manager{cf.New(df.Of(&Dim{width: 0, height: 25}))}
	t.ErrIs(fx.Reflow(nil), ErrDim)
}

// fxWF1 chainer with sole width filler.
var fxWF1 = func() *chainerFX { return cf.New(df.Filling()) }

// fxWF2 stacker with width filler at the beginning.
var fxWF2 = func() *chainerFX {
	return cf.New(df.Filling(), df.Fixed())
}

// fxWF3 stacker with width filler in between.
var fxWF3 = func() *chainerFX {
	return cf.New(df.Fixed(), df.Filling(), df.Fixed())
}

// fxWF4 stacker with width filler at the end.
var fxWF4 = func() *chainerFX {
	return cf.New(df.Fixed(), df.Filling())
}

func (s *chained) With_width_filler_consume_all_width(t *T) {
	for _, fx := range []*chainerFX{fxWF1(), fxWF2(), fxWF3(), fxWF4()} {
		// TODO: fails if -race and -count=10
		t.True(fx.Width() > fx.SumLayoutWidths())
		t.FatalOn((&Manager{Root: fx}).Reflow(nil))
		t.Eq(fx.Width(), fx.SumLayoutWidths())
	}
}

func (s *chained) With_width_filler_have_no_margins_at_fixed(t *T) {
	fx := cf.New(df.Filling(), df.Fixed(), df.Filling(), df.Fixed(),
		df.Filling())
	for _, d := range fx.dd {
		d.Dim().mrgRight, d.Dim().mrgLeft = 1, 1
	}
	t.FatalOn((&Manager{Root: fx}).Reflow(nil))
	for _, d := range fx.dd {
		t.Eq(0, d.Dim().mrgRight)
		t.Eq(0, d.Dim().mrgLeft)
	}
}

func (s *chained) With_underflowing_width_have_tb_margins(t *T) {
	// single margin (fx1), distributed margins (fx2), exact fit (fx3)
	fx1, fx2 := cf.New(df.Fixed()), cf.New(df.Fixed(), df.Fixed())
	fx3 := cf.New(df.Fixed(), df.Fixed(), df.FixedW(40))
	// []fixtures[]dimerInFixture[]marginTopBottom
	exp := [][][2]int{{{30, 30}}, {{13, 6}, {7, 14}},
		{{0, 0}, {0, 0}, {0, 0}}}
	for i, fx := range []*chainerFX{fx1, fx2, fx3} {
		t.FatalOn((&Manager{Root: fx}).Reflow(nil))
		for j, d := range fx.dd {
			_, mr, _, ml := d.Dim().Margin()
			t.Eq(exp[i][j][0], ml)
			t.Eq(exp[i][j][1], mr)
		}
	}
}

func (s *chained) With_underflowing_width_consume_all_width(t *T) {
	// single margin (fx1), distributed margins (fx2), exact fit (fx3)
	fx1, fx2 := cf.New(df.Fixed()), cf.New(df.Fixed(), df.Fixed())
	fx3 := cf.New(df.Fixed(), df.Fixed(), df.FixedW(40))
	for _, fx := range []*chainerFX{fx1, fx2, fx3} {
		t.FatalOn((&Manager{Root: fx}).Reflow(nil))
		t.Eq(fx.Width(), fx.SumLayoutWidths())
	}
}

func (s *chained) With_overflowing_width_are_clipped(t *T) {
	// single overflow (fx1), last overflows (fx2), off-screen (fx3)
	fx1, fx2 := cf.New(df.FixedW(81)), cf.New(df.Fixed(), df.FixedW(61))
	fx3 := cf.New(df.Fixed(), df.FixedW(61), df.Fixed())
	fx4 := cf.New(
		df.FixedW(22), df.FillingW(20), df.Fixed(), df.FillingW(20))
	t.FatalOn((&Manager{Root: fx1}).Reflow(nil))
	clippedWidth, _ := fx1.dd[0].Dim().Clip()
	t.Eq(1, clippedWidth)
	t.FatalOn((&Manager{Root: fx2}).Reflow(nil))
	clippedWidth, _ = fx2.dd[1].Dim().Clip()
	t.Eq(1, clippedWidth)
	t.FatalOn((&Manager{Root: fx3}).Reflow(nil))
	clippedWidth, _ = fx3.dd[1].Dim().Clip()
	t.Eq(1, clippedWidth)
	t.True(fx3.dd[2].Dim().IsOffScreen())
	t.FatalOn((&Manager{Root: fx4}).Reflow(nil))
	clippedWidth, _ = fx4.dd[3].Dim().Clip()
	t.Eq(2, clippedWidth)
}

func (s *chained) With_fitting_layout_has_no_clipping(t *T) {
	fx := cf.New(df.Fixed(), df.Filling(), df.Fixed())
	for _, d := range fx.dd {
		d.Dim().clipWidth = 1
	}
	t.FatalOn((&Manager{Root: fx}).Reflow(nil))
	for _, d := range fx.dd {
		t.Eq(0, d.Dim().clipWidth)
	}
}

func (s *chained) Width_filler_have_width_evenly_distributed(t *T) {
	// i.e. min-height 10 available 80 => 50/3: 17,17,16 => 27,27,26
	fx := cf.New(df.Filling(), df.Filling(), df.Filling())
	t.FatalOn((&Manager{Root: fx}).Reflow(nil))
	for i, d := range fx.dd {
		_, _, got, _ := d.Dim().Area()
		switch i {
		case 2:
			t.Eq(26, got)
		default:
			t.Eq(27, got)
		}
	}
}

func (s *chained) Height_fillers_have_no_top_or_bottom_margins(t *T) {
	fx := cf.New(df.Filling(), df.Filling(), df.Filling())
	for _, d := range fx.dd {
		d.Dim().mrgTop, d.Dim().mrgBottom = 1, 1
	}
	t.FatalOn((&Manager{Root: fx}).Reflow(nil))
	for _, d := range fx.dd {
		mt, _, mb, _ := d.Dim().Margin()
		t.True(mt == mb && mb == 0)
	}
}

func (s *chained) Height_fillers_consume_all_height(t *T) {
	fx := cf.New(df.Filling(), df.Filling(), df.Filling())
	_, _, _, exp := fx.Dim().Area()
	t.FatalOn((&Manager{Root: fx}).Reflow(nil))
	for _, d := range fx.dd {
		_, _, _, got := d.Dim().Area()
		t.Eq(exp, got)
	}
}

func (s *chained) With_underflowing_height_have_tb_margins(t *T) {
	fx := cf.New(df.FixedH(9))
	t.FatalOn((&Manager{Root: fx}).Reflow(nil))
	gotTM, _, gotBM, _ := fx.dd[0].Dim().Margin()
	t.True(gotTM == gotBM && gotBM == 8)
}

func (s *chained) With_underflowing_height_consume_all_height(t *T) {
	fx := cf.New(df.FixedH(9))
	t.FatalOn((&Manager{Root: fx}).Reflow(nil))
	t.Eq(fx.Dim().height, fx.dd[0].Dim().layoutHeight())
}

func (s *chained) With_overflowing_height_are_clipped_to_height(t *T) {
	fx := cf.New(df.FixedH(30), df.FillingH(30))
	t.FatalOn((&Manager{Root: fx}).Reflow(nil))
	_, gotFixed := fx.dd[0].Dim().Clip()
	t.Eq(5, gotFixed)
	_, gotFilled := fx.dd[1].Dim().Clip()
	t.Eq(5, gotFilled)
	t.Eq(fx.Dim().height, fx.dd[0].Dim().layoutHeight())
	t.Eq(fx.Dim().height, fx.dd[1].Dim().layoutHeight())
}

func TestChained(t *testing.T) {
	t.Parallel()
	Run(&chained{}, t)
}
