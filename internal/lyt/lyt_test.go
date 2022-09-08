// Copyright (c) 2022 Stephan Lukits. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package lyt

import (
	"testing"

	. "github.com/slukits/gounit"
)

type dimerFixture struct {
	dim *Dim
}

func (df *dimerFixture) Dim() *Dim { return df.dim }

type dimerFactory struct{}

// df produces test-fixtures implementing the Dimer interface which must
// be implemented by any layouted component.
var df = &dimerFactory{}

// New creates a Dimer with zero dimensions.
func (df *dimerFactory) New() Dimer {
	return &dimerFixture{&Dim{}}
}

// Of creates a Dimer with given dimensions.
func (df *dimerFactory) Of(d *Dim) Dimer {
	return &dimerFixture{d}
}

// Filling provides filling Dimer with default filling width and height
// (20, 3).
func (df *dimerFactory) Filling() Dimer {
	return &dimerFixture{DimFilling(10, 3)}
}

// FillingOne provides filling Dimer with filling width 1 and filling
// height 1.
func (df *dimerFactory) FillingOne() Dimer {
	return &dimerFixture{DimFilling(1, 1)}
}

// FillingH provides filling Dimer with default filling width and given
// filling fHeight (20, height).
func (df *dimerFactory) FillingH(fHeight int) Dimer {
	return &dimerFixture{DimFilling(20, fHeight)}
}

// FillingW provides filling Dimer with given filling fWidth and default
// filling height (fWidth, 8).
func (df *dimerFactory) FillingW(fWidth int) Dimer {
	return &dimerFixture{DimFilling(fWidth, 8)}
}

// FillingWH provides filling Dimer with given filling fWidth and
// fHeight.
func (df *dimerFactory) FillingWH(fWidth, fHeight int) Dimer {
	return &dimerFixture{DimFilling(fWidth, fHeight)}
}

// Fixed provides fixed Dimer with default width and height (20, 8).
func (df *dimerFactory) Fixed() Dimer {
	return &dimerFixture{DimFixed(20, 8)}
}

// FixedH provides fixed Dimer with default width and given height (20,
// height).
func (df *dimerFactory) FixedH(height int) Dimer {
	return &dimerFixture{DimFixed(20, height)}
}

// FixedW provides fixed Dimer with given width and default height
// (width, 8).
func (df *dimerFactory) FixedW(width int) Dimer {
	return &dimerFixture{DimFixed(width, 8)}
}

// FixedWH provides fixed Dimer with given width and height.

func (df *dimerFactory) FixedWH(width, height int) Dimer {
	return &dimerFixture{DimFixed(width, height)}
}

// FillingFixed provides Dimer with given filling fWidth and fixed
// height.
func (df *dimerFactory) FillingFixed(fWidth, height int) Dimer {
	return &dimerFixture{DimFillingWidth(fWidth, height)}
}

// FixedFilling provides Dimer with given fixed width and filling
// fHeight.
func (df *dimerFactory) FixedFilling(width, fHeight int) Dimer {
	return &dimerFixture{DimFillingHeight(width, fHeight)}
}

// Screen provides fixed Dimer with default terminal screen size (80,
// 25).
func (df *dimerFactory) Screen() Dimer {
	return &dimerFixture{&Dim{width: 80, height: 25}}
}

// ScreenW provides fixed Dimer with default height and given width
// (width, 25).
func (df *dimerFactory) ScreenW(width int) Dimer {
	return &dimerFixture{&Dim{width: width, height: 25}}
}

// ScreenH provides fixed Dimer with default width and given height
// (80, height).
func (df *dimerFactory) ScreenH(height int) Dimer {
	return &dimerFixture{&Dim{width: 80, height: height}}
}

// managerFactory creates Manger-fixtures
type managerFactory struct{ width, height int }

var mf = &managerFactory{}

// New provides a Manager with 80x25-Root.
func (mf *managerFactory) New() *Manager {
	return &Manager{Root: df.Screen()}
}

// Of provides a Manager with given Dimer as root.
func (mf *managerFactory) Of(d Dimer) *Manager {
	return &Manager{Root: d}
}

func (mf *managerFactory) ScreenOf(d Dimer) *Manager {
	return &Manager{Width: 80, Height: 25, Root: d}
}

func (mf *managerFactory) WHOf(width, height int, d Dimer) *Manager {
	return &Manager{Width: width, Height: height, Root: d}
}

func (mf *managerFactory) WHFilling(width, height int) *Manager {
	return &Manager{Width: width, Height: height, Root: df.FillingOne()}
}

type manager struct{ Suite }

func (s *manager) SetUp(t *T) { t.Parallel() }

func (s *manager) Operations_fail_if_root_unset(t *T) {
	m := mf.Of(nil)
	t.ErrIs(m.Reflow(nil), ErrLyt)
	_, err := m.Locate(nil)
	t.ErrIs(err, ErrLyt)
}

func (s *manager) Operations_fail_if_root_s_height_not_positive(t *T) {
	m := mf.Of(df.ScreenH(0))
	t.ErrIs(m.Reflow(nil), ErrLyt)
}

func (s *manager) Operations_fail_if_root_s_width_not_positive(t *T) {
	m := mf.Of(df.ScreenW(0))
	t.ErrIs(m.Reflow(nil), ErrLyt)
}

func (s *manager) Centers_vertically_fixed_height_root(t *T) {
	m := mf.ScreenOf(df.FillingFixed(1, 15))
	top, _, bottom, _ := m.Root.Dim().Margin()
	t.True(top == 0 && top == bottom)

	t.FatalOn(m.Reflow(nil))

	top, _, bottom, _ = m.Root.Dim().Margin()
	t.Eq(5, top)
	t.Eq(5, bottom)
}

func (s *manager) Centers_horizontally_fixed_width_root(t *T) {
	m := mf.ScreenOf(df.FixedFilling(40, 1))
	_, right, _, left := m.Root.Dim().Margin()
	t.True(right == 0 && left == right)

	t.FatalOn(m.Reflow(nil))

	_, right, _, left = m.Root.Dim().Margin()
	t.Eq(20, right)
	t.Eq(20, left)
}

func (s *manager) Makes_root_without_size_filling(t *T) {
	m := mf.ScreenOf(&dimerFixture{dim: &Dim{}})
	t.False(m.Root.Dim().IsFillingHeight())
	t.False(m.Root.Dim().IsFillingWidth())

	t.FatalOn(m.Reflow(nil))

	t.True(m.Root.Dim().IsFillingHeight())
	t.True(m.Root.Dim().IsFillingWidth())
}

func (s *manager) Assigns_its_width_and_height_to_root(t *T) {
	m := mf.ScreenOf(&dimerFixture{dim: &Dim{}})
	t.True(m.Root.Dim().Width() == 0 &&
		m.Root.Dim().Width() == m.Root.Dim().Height())

	t.FatalOn(m.Reflow(nil))

	t.Eq(m.Width, m.Root.Dim().Width())
	t.Eq(m.Height, m.Root.Dim().Height())
}

func (s *manager) Accounts_for_clipping_checking_consistency(t *T) {
	test := func(sc Dimer, d Dimer) {
		d.Dim().width += 2
		d.Dim().height += 2
		fx := mf.Of(sc)
		t.False(fx.HasConsistentLayout())
		d.Dim().clipHeight = 2
		d.Dim().clipWidth = 2
		t.True(fx.HasConsistentLayout())
	}
	sd, cd := df.Screen(), df.Screen()
	st, cn := sf.New(sd), cf.New(cd)
	test(st, sd)
	test(cn, cd)
}

func (s *manager) Accounts_for_margins_checking_consistency(t *T) {
	test := func(sc Dimer, d Dimer) {
		d.Dim().width -= 2
		d.Dim().height -= 2
		fx := mf.Of(sc)
		t.False(fx.HasConsistentLayout())
		d.Dim().mrgTop = 1
		d.Dim().mrgRight = 1
		d.Dim().mrgBottom = 1
		d.Dim().mrgLeft = 1
		t.True(fx.HasConsistentLayout())
	}
	sd, cd := df.Screen(), df.Screen()
	st, cn := sf.New(sd), cf.New(cd)
	test(st, sd)
	test(cn, cd)
}

func (s *manager) Leafs_printable_hight_after_update(t *T) {
	fx := mf.Of(sf.Of(df.FixedWH(10, 10), df.FillingOne()))
	fx.Reflow(nil)
	_, _, _, h := fx.Root.(*stackerFX).dd[0].Dim().Area()
	t.Eq(10, h)
	fx.Root.(*stackerFX).dd[0].Dim().UpdateHeight(1)
	fx.Reflow(nil)
	t.False(fx.Root.(*stackerFX).dd[0].Dim().IsDirty())
	_, _, _, h = fx.Root.(*stackerFX).dd[0].Dim().Area()
	t.Eq(10, h)
}

func (s *manager) Leafs_printable_width_after_update(t *T) {
	fx := mf.Of(sf.Of(df.FixedWH(10, 10), df.FillingOne()))
	fx.Reflow(nil)
	_, _, w, _ := fx.Root.(*stackerFX).dd[0].Dim().Area()
	t.Eq(10, w)
	fx.Root.(*stackerFX).dd[0].Dim().UpdateWidth(1)
	fx.Reflow(nil)
	t.False(fx.Root.(*stackerFX).dd[0].Dim().IsDirty())
	_, _, w, _ = fx.Root.(*stackerFX).dd[0].Dim().Area()
	t.Eq(10, w)
}

func (s *manager) Updates_printable_height_after_update(t *T) {
	fx := mf.Of(sf.Of(
		df.FixedWH(10, 10), df.FillingOne(), df.FillingOne()))
	fx.Reflow(nil)
	_, _, _, h0 := fx.Root.(*stackerFX).dd[0].Dim().Area()
	_, _, _, h1 := fx.Root.(*stackerFX).dd[1].Dim().Area()
	t.Eq(5, h0)
	t.Eq(5, h1)
	fx.Root.(*stackerFX).dd[0].Dim().UpdateHeight(2)
	fx.Reflow(nil)
	_, _, _, h0 = fx.Root.(*stackerFX).dd[0].Dim().Area()
	_, _, _, h1 = fx.Root.(*stackerFX).dd[1].Dim().Area()
	t.True(fx.Root.(*stackerFX).dd[0].Dim().IsFillingHeight())
	t.True(fx.Root.(*stackerFX).dd[1].Dim().IsFillingHeight())
	t.Eq(7, h0)
	t.Eq(3, h1)
}

func (s *manager) Locates_a_dimer_in_layout(t *T) {
	exp := df.Fixed()
	fx := sf.New(sf.Filling(df.Fixed(), sf.Filling(
		df.Fixed(), exp, df.Fixed()), df.Filling()))
	got, err := mf.Of(fx).Locate(exp)
	t.FatalOn(err)
	t.FatalIfNot(t.Eq(3, len(got)))
	has := false
	got[len(got)-1].(Stacker).ForStacked(func(d Dimer) (stop bool) {
		if d == exp {
			has = true
		}
		return false
	})
	t.True(has)
}

func (s *manager) Provides_a_narrowing_path_of_dimer_around_coordinates(
	t *T,
) {
	exp3 := df.Fixed()
	exp2 := cf.Filling(df.Filling(), exp3, df.Filling())
	exp1 := sf.Of(df.Screen(), df.Filling(), exp2, df.Filling())
	fx := mf.Of(exp1)
	t.FatalOn(fx.Reflow(nil))
	path, err := fx.LocateAt(exp3.Dim().x+1, exp3.Dim().y+1)
	t.FatalOn(err)
	t.True(exp1 == path[0])
	t.True(exp2 == path[1])
	t.True(exp3 == path[2])
}

func (s *manager) Provides_nil_path_if_dimer_not_locatable(t *T) {
	exp := df.Fixed()
	fx := sf.New(sf.Filling(df.Fixed(), sf.Filling(
		df.Fixed(), df.Fixed()), df.Filling()))
	got, err := mf.Of(fx).Locate(exp)
	t.FatalOn(err)
	t.Eq(0, len(got))
}

func (s *manager) Recursively_layouts_stacked_dimer(t *T) {
	fx := sf.New(df.Fixed(), sf.Filling(df.Fixed(), df.Filling()))
	t.False(fx.HasConsistentLayout())
	t.False(fx.dd[1].(*stackerFX).HasConsistentLayout())
	t.FatalOn((&Manager{Root: fx}).Reflow(nil))
	t.True(fx.HasConsistentLayout())
	t.True(fx.dd[1].(*stackerFX).HasConsistentLayout())
}

func TestManager(t *testing.T) {
	t.Parallel()
	Run(&manager{}, t)
}
