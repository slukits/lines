// Copyright (c) 2022 Stephan Lukits. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package selects

import (
	"testing"

	. "github.com/slukits/gounit"
	"github.com/slukits/lines"
	"github.com/slukits/lines/cmp/fx"
)

type AZeroList struct{ Suite }

func (s *AZeroList) SetUp(t *T) { t.Parallel() }

func zeroFX(t *T) (*lines.Fixture, *List) {
	cmp := &List{}
	fx := fx.New(t, cmp)
	return fx, cmp
}

func (s *AZeroList) Has_no_items(t *T) {
	_, zl := zeroFX(t)
	t.True(zl.Items == nil)
}

func (s *AZeroList) Has_no_scrollable_liner(t *T) {
	_, zl := zeroFX(t)
	t.True(zl.SelectableLiner == nil)
}

func (s *AZeroList) Is_not_focusable(t *T) {
	fx, zl := zeroFX(t)
	t.FatalOn(fx.Lines.Update(zl, nil, func(e *lines.Env) {
		t.Not.True(zl.FF.Has(lines.Focusable))
	}))
}

func (s *AZeroList) Lines_are_not_focusable(t *T) {
	fx, zl := zeroFX(t)
	t.FatalOn(fx.Lines.Update(zl, nil, func(e *lines.Env) {
		t.Not.True(zl.FF.Has(lines.LinesFocusable))
	}))
}

func (s *AZeroList) Shows_the_zero_item(t *T) {
	fx, zl := zeroFX(t)
	t.Eq(fx.ScreenOf(zl).Trimmed(), NoItems)
}

func (s *AZeroList) Height_collapses_to_one(t *T) {
	fx, zl := zeroFX(t)
	fx.Lines.Update(fx.Root(), nil, func(e *lines.Env) {
		t.Eq(1, zl.Dim().Height())
	})
}

func (s *AZeroList) Width_collapses_to_zero_item_width(t *T) {
	fx, zl := zeroFX(t)
	fx.Lines.Update(fx.Root(), nil, func(e *lines.Env) {
		t.Eq(len(NoItems), zl.Dim().Width())
	})
}

type gapper interface {
	Gaps(int) *lines.GapsWriter
}

func (s *AZeroList) Without_item_has_height_one_plus_horizontal_gaps(t *T) {
	cmp := &fx.Wrap{Componenter: &List{},
		ONInit: func(cmp lines.Componenter, e *lines.Env) {
			lines.Print(cmp.(gapper).Gaps(0).Filling(), ' ')
		}}
	fx := fx.New(t, cmp)
	fx.Lines.Update(cmp, nil, func(e *lines.Env) {
		t.Eq(3, cmp.Componenter.(lines.Dimer).Dim().Height())
	})
}

func (s *AZeroList) Ignores_the_mouse(t *T) {
	fx, zl := zeroFX(t)
	cc, x, y := fx.Cells().Trimmed(), 0, 0
	fx.Lines.Update(zl, nil, func(e *lines.Env) {
		x, y, _, _ = zl.Dim().Printable()
	})

	fx.FireMove(x, y)
	t.True(cc.Equals(fx.Cells().Trimmed()))

	fx.FireClick(x, y)
	t.True(cc.Equals(fx.Cells().Trimmed()))
}

func TestAZeroList(t *testing.T) {
	t.Parallel()
	Run(&AZeroList{}, t)
}

type AList struct{ Suite }

func (s *AList) With_items_is_not_zero(t *T) {
	t.Not.True((&List{Items: []string{"item"}}).IsZero())
}

func (s *AList) With_scrollable_liner_is_not_zero(t *T) {
	t.Not.True((&List{SelectableLiner: &fx.SelectableLiner{}}).IsZero())
}

func (s *AList) Is_Focusable(t *T) {
	cmp := &List{Items: []string{"item"}}
	fx := fx.New(t, cmp)
	t.FatalOn(fx.Lines.Update(cmp, nil, func(e *lines.Env) {
		t.True(cmp.FF.Has(lines.Focusable))
	}))
}

func (s *AList) Has_focusable_lines(t *T) {
	cmp := &List{Items: []string{"item"}}
	fx := fx.New(t, cmp)
	t.FatalOn(fx.Lines.Update(cmp, nil, func(e *lines.Env) {
		t.True(cmp.FF.Has(lines.LinesFocusable))
	}))
}

func (s *AList) Is_scrollable_if_items_exceed_component_height(t *T) {
	cmp := &List{Items: fx.NStrings(50)}
	fx := fx.New(t, cmp)
	fx.Lines.Update(cmp, nil, func(e *lines.Env) {
		t.True(cmp.FF.Has(lines.Scrollable))
	})
}

func (s *AList) Height_collapses_to_items_count_plus_gaps(t *T) {
	cmp, before := &List{Items: fx.NStrings(10)}, 0
	wrp := &fx.Wrap{Componenter: cmp,
		ONInit: func(c lines.Componenter, e *lines.Env) {
			lines.Print(cmp.Gaps(0).Filling(), ' ')
			lines.Print(cmp.Gaps(1).Filling(), ' ')
		},
		ONLayout: func(c lines.Componenter, e *lines.Env) (reflow bool) {
			before = cmp.Dim().Height()
			return false
		}}
	fx := fx.New(t, wrp)
	fx.Lines.Update(cmp, nil, func(e *lines.Env) {
		t.Not.True(cmp.FF.Has(lines.Scrollable))
		t.Eq(14, cmp.Dim().Height())
		t.True(before > 14)
	})
	t.True(before > 0)
}

func (s *AList) Height_collapses_to_lines_len_plus_gaps(t *T) {
	cmp, before := &List{SelectableLiner: &fx.SelectableLiner{}}, 0
	cmp.SelectableLiner.(*fx.SelectableLiner).InitLines(10)
	wrp := &fx.Wrap{Componenter: cmp,
		ONInit: func(c lines.Componenter, e *lines.Env) {
			lines.Print(cmp.Gaps(0).Filling(), ' ')
			lines.Print(cmp.Gaps(1).Filling(), ' ')
		},
		ONLayout: func(c lines.Componenter, e *lines.Env) (reflow bool) {
			before = cmp.Dim().Height()
			return false
		}}
	fx := fx.New(t, wrp)
	fx.Lines.Update(cmp, nil, func(e *lines.Env) {
		t.Not.True(cmp.FF.Has(lines.Scrollable))
		t.Eq(14, cmp.Dim().Height())
		t.True(before > 14)
	})
	t.True(before > 0)
}

func (s *AList) Width_collapses_to_max_items_width_and_gaps(t *T) {
	cmp, before := &List{Items: []string{"12", "1234", "123"}}, 0
	wrp := &fx.Wrap{Componenter: cmp,
		ONInit: func(c lines.Componenter, e *lines.Env) {
			lines.Print(cmp.Gaps(0).Filling(), ' ')
			lines.Print(cmp.Gaps(1).Filling(), ' ')
		},
		ONLayout: func(c lines.Componenter, e *lines.Env) (reflow bool) {
			before = cmp.Dim().Width()
			return false
		}}
	fx := fx.New(t, wrp)
	fx.Lines.Update(cmp, nil, func(e *lines.Env) {
		t.Not.True(cmp.FF.Has(lines.Scrollable))
		t.Eq(8, cmp.Dim().Width())
		t.True(before > 8)
	})
	t.True(before > 0)
}

func (s *AList) Width_collapses_to_liner_width_and_gaps(t *T) {
	cmp, before := &List{SelectableLiner: &fx.SelectableLiner{}}, 0
	cmp.SelectableLiner.(*fx.SelectableLiner).II =
		[]string{"12", "1234", "123"}
	wrp := &fx.Wrap{Componenter: cmp,
		ONInit: func(c lines.Componenter, e *lines.Env) {
			lines.Print(cmp.Gaps(0).Filling(), ' ')
			lines.Print(cmp.Gaps(1).Filling(), ' ')
		},
		ONLayout: func(c lines.Componenter, e *lines.Env) (reflow bool) {
			before = cmp.Dim().Width()
			return false
		}}
	fx := fx.New(t, wrp)
	fx.Lines.Update(cmp, nil, func(e *lines.Env) {
		t.Not.True(cmp.FF.Has(lines.Scrollable))
		t.Eq(8, cmp.Dim().Width())
		t.True(before > 8)
	})
	t.True(before > 0)
}

func (s *AList) Displays_its_items(t *T) {
	cmpII := &List{Items: []string{"12", "1234", "123"}}
	cmpLN := &List{SelectableLiner: &fx.SelectableLiner{}}
	cmpLN.SelectableLiner.(*fx.SelectableLiner).II =
		[]string{"12", "1234", "123"}
	cmp := (&fx.Chaining{}).Set(cmpII, cmpLN)
	fx := fx.New(t, cmp)
	t.FatalIfNot(t.Eq("12  \n1234\n123 ", fx.ScreenOf(cmpII)))
	t.FatalIfNot(t.Eq(fx.ScreenOf(cmpII), fx.ScreenOf(cmpLN)))
	t.True(fx.CellsOf(cmpII).Equals(fx.CellsOf(cmpLN)))
}

func (s *AList) Style_defaults_to_the_reversed_global(t *T) {
	cmpII := &List{Items: []string{"12", "1234", "123"}}
	cmpLN := &List{SelectableLiner: &fx.SelectableLiner{}}
	cmpLN.SelectableLiner.(*fx.SelectableLiner).II =
		[]string{"12", "1234", "123"}
	cmp := (&fx.Chaining{}).Set(cmpII, cmpLN)
	fx := fx.New(t, cmp)
	t.FatalIfNot(t.Eq("12  \n1234\n123 ", fx.ScreenOf(cmpII)))
	glbRvr := fx.Lines.Globals.Style(lines.Default).Reverse()
	for _, l := range fx.CellsOf(cmpII) {
		t.Eq(glbRvr, l[0].Style)
	}
	for _, l := range fx.CellsOf(cmpLN) {
		t.Eq(glbRvr, l[0].Style)
	}
}

func (s *AList) Applies_provided_styler(t *T) {
	sty := lines.DefaultStyle.WithBG(lines.DarkBlue).
		WithFG(lines.Silver).WithAA(lines.Bold)
	cmp := &List{
		Items: []string{"12", "1234", "123"},
		Styler: func(_ int, highlight bool) lines.Style {
			if !highlight {
				return sty
			}
			return lines.DefaultStyle
		},
	}
	fx := fx.New(t, cmp)
	t.FatalIfNot(t.Eq("12  \n1234\n123 ", fx.ScreenOf(cmp)))
	for _, l := range fx.CellsOf(cmp) {
		t.Eq(sty, l[0].Style)
	}
}

func (s *AList) Liner_supersedes_items_and_styler(t *T) {
	sty := lines.DefaultStyle.WithBG(lines.DarkBlue).
		WithFG(lines.Silver).WithAA(lines.Bold)
	styLN := lines.DefaultStyle.WithBG(lines.Yellow).
		WithFG(lines.Red).WithAA(lines.Italic)
	cmp := &List{
		SelectableLiner: (&fx.SelectableLiner{}).
			SetII([]string{"12", "1234", "123"}).SetSty(styLN),
		Items:  []string{"blub", "42", "hurz"},
		Styler: func(_ int, _ bool) lines.Style { return sty },
	}
	fx := fx.New(t, cmp)
	t.FatalIfNot(t.Eq("12  \n1234\n123 ", fx.ScreenOf(cmp)))
	for _, l := range fx.CellsOf(cmp) {
		t.Eq(styLN, l[0].Style)
	}
}

func TestAList(t *testing.T) {
	t.Parallel()
	Run(&AList{}, t)
}
