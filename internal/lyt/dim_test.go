// Copyright (c) 2022 Stephan Lukits. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package lyt

import (
	"testing"

	. "github.com/slukits/gounit"
)

type _Dim struct{ Suite }

func SetUp(t *T) { t.Parallel() }

func (s *_Dim) Panics_if_constructor_receives_zero_dimension(t *T) {
	vv := [][]int{{0, 1}, {1, 0}, {0, 0}}
	ff := []func(int, int) *Dim{
		DimFilling, DimFillingWidth, DimFillingHeight, DimFixed}
	for _, f := range ff {
		for _, v := range vv {
			t.Panics(func() { f(v[0], v[1]) })
		}
	}
}

func (s *_Dim) Ignores_width_update_resulting_in_non_positive_width(
	t *T,
) {
	fx := DimFixed(10, 10)
	t.Not.True(fx.IsUpdated())
	fx.UpdateWidth(-10)
	t.Not.True(fx.IsUpdated())
}

func (s *_Dim) Is_updated_after_valid_width_update(t *T) {
	fx := DimFixed(10, 10)
	t.Not.True(fx.IsUpdated())
	fx.UpdateWidth(-5)
	t.True(fx.IsUpdated())
}

func (s *_Dim) Ignores_height_update_resulting_in_non_positive_height(
	t *T,
) {
	fx := DimFixed(10, 10)
	t.Not.True(fx.IsUpdated())
	fx.UpdateHeight(-10)
	t.Not.True(fx.IsUpdated())
}

func (s *_Dim) Is_updated_after_valid_height_update(t *T) {
	fx := DimFixed(10, 10)
	t.Not.True(fx.IsUpdated())
	fx.UpdateHeight(-5)
	t.True(fx.IsUpdated())
}

func (s *_Dim) Stops_filling_width_if_width_set(t *T) {
	fx := DimFillingWidth(1, 10)
	t.True(fx.IsFillingWidth())
	fx.SetWidth(5)
	t.Not.True(fx.IsFillingWidth())
}

func (s *_Dim) Stops_filling_height_if_height_set(t *T) {
	fx := DimFillingHeight(10, 1)
	t.True(fx.IsFillingHeight())
	fx.SetHeight(5)
	t.Not.True(fx.IsFillingHeight())
}

func (s *_Dim) Is_not_dirty_after_being_cleaned(t *T) {
	fx := DimFixed(10, 10)
	d := fx.prepareLayout()
	fx.SetWidth(5)
	fx.finalizeLayout(d)
	t.True(fx.IsDirty())
	fx.setClean()
	t.Not.True(fx.IsDirty())
}

func TestDim(t *testing.T) {
	t.Parallel()
	Run(&_Dim{}, t)
}
