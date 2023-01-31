// Copyright (c) 2022 Stephan Lukits. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package selects

import (
	"fmt"
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

func (s *AList) Has_highlight_able_lines(t *T) {
	cmp := &List{Items: []string{"item"}}
	fx := fx.New(t, cmp)
	t.FatalOn(fx.Lines.Update(cmp, nil, func(e *lines.Env) {
		t.True(cmp.FF.Has(lines.HighlightEnabled))
	}))
}

func (s *AList) Has_selectable_lines(t *T) {
	cmp := &List{Items: []string{"item"}}
	fx := fx.New(t, cmp)
	t.FatalOn(fx.Lines.Update(cmp, nil, func(e *lines.Env) {
		t.True(cmp.FF.Has(lines.LineSelectable))
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
		Styler: func(_ int) lines.Style {
			return sty
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
		Styler: func(_ int) lines.Style { return sty },
	}
	fx := fx.New(t, cmp)
	t.FatalIfNot(t.Eq("12  \n1234\n123 ", fx.ScreenOf(cmp)))
	for _, l := range fx.CellsOf(cmp) {
		t.Eq(styLN, l[0].Style)
	}
}

func (s *AList) Highlights_its_first_item_on_down_key(t *T) {
	cmp := &List{Items: []string{"12", "1234", "123"}}
	fx_ := fx.New(t, cmp)
	var hi lines.Style
	fx_.Lines.Update(cmp, nil, func(e *lines.Env) {
		hi = cmp.Globals().Style(lines.Highlight)
	})
	fx_.FireKey(lines.Down)
	t.Eq(hi, fx_.CellsOf(cmp)[0][1].Style)

	cmp = &List{SelectableLiner: &fx.SelectableLiner{}}
	cmp.SelectableLiner.(*fx.SelectableLiner).II =
		[]string{"12", "1234", "123"}
	fx_ = fx.New(t, cmp)
	fx_.FireKey(lines.Down)
	t.Eq(hi, fx_.CellsOf(cmp)[0][0].Style)
}

func (s *AList) Highlights_a_mouse_hovered_list_item(t *T) {
	cmp := &List{Items: []string{"12", "1234", "123"}}
	fx := fx.New(t, cmp)
	x, y := 0, 0
	fx.Lines.Update(cmp, nil, func(e *lines.Env) {
		x, y, _, _ = cmp.Dim().Printable()
	})
	dflt := fx.Lines.Globals.Style(lines.Highlight)

	fx.FireMove(x, y) // move to/highlight the first item
	t.Eq(fx.CellsOf(cmp)[0][0].Style, fx.Lines.Globals.Style(lines.Default))
	t.Eq(fx.CellsOf(cmp)[1][0].Style, dflt)
	fx.FireMove(x+1, y+2) // move to/highlight the third item
	t.Eq(fx.CellsOf(cmp)[2][0].Style, fx.Lines.Globals.Style(lines.Default))
	t.Eq(fx.CellsOf(cmp)[0][0].Style, dflt)
}

func (s *AList) Uses_liner_highlighter(t *T) {
	hi := lines.DefaultStyle.WithBG(lines.DarkBlue).
		WithFG(lines.Silver).WithAA(lines.Bold)
	cmp := &List{SelectableLiner: (&fx.HighlightingLiner{}).SetHi(hi)}
	cmp.SelectableLiner.(*fx.HighlightingLiner).II =
		[]string{"12", "1234", "123"}
	fx := fx.New(t, cmp)
	fx.FireKey(lines.Down) // focus first line
	t.Eq(fx.CellsOf(cmp)[0][0].Style, hi)

	x, y := 0, 0
	fx.Lines.Update(cmp, nil, func(e *lines.Env) {
		x, y, _, _ = cmp.Dim().Printable()
	})
	fx.FireMove(x, y+2) // focus third line
	t.Eq(fx.CellsOf(cmp)[2][0].Style, hi)
}

func (s *AList) Removes_line_focus_on_exit(t *T) {
	lst := &List{Items: []string{"12", "1234", "123"}}
	cmp := (&fx.Stacking{}).Set(lst, &fx.Cmp{})
	fx := fx.New(t, cmp)
	x, y, height := 0, 0, 0
	fx.Lines.Update(lst, nil, func(e *lines.Env) {
		x, y, _, height = lst.Dim().Printable()
	})
	dflt := fx.Lines.Globals.Style(lines.Highlight)

	fx.FireMove(x, y) // move to/highlight the first item
	t.Eq(fx.CellsOf(lst)[0][0].Style, fx.Lines.Globals.Style(lines.Default))
	fx.FireMove(x, y+height)
	t.Eq(fx.CellsOf(lst)[0][0].Style, dflt)
}

func (s *AList) Removes_line_focus_on_focus_lost(t *T) {
	lst, c := &List{Items: []string{"12", "1234", "123"}}, &fx.Cmp{}
	cmp := (&fx.Stacking{}).Set(lst, c)
	fx := fx.New(t, cmp)
	fx.Lines.Focus(lst)
	x, y := 0, 0
	fx.Lines.Update(lst, nil, func(e *lines.Env) {
		x, y, _, _ = lst.Dim().Printable()
	})
	dflt := fx.Lines.Globals.Style(lines.Highlight)

	fx.FireMove(x, y) // move to/highlight the first item
	t.Eq(fx.CellsOf(lst)[0][0].Style, fx.Lines.Globals.Style(lines.Default))
	fx.Lines.Focus(c)
	t.Eq(fx.CellsOf(lst)[0][0].Style, dflt)
}

func (s *AList) Reports_clicked_list_item(t *T) {
	reported := -1
	cmp := &List{
		Items: []string{"12", "1234", "123"},
		Listener: &fx.Cmp{
			ONUpdate: func(c *fx.Cmp, e *lines.Env, i interface{}) {
				reported = int(i.(Value))
			}},
	}
	fx_ := fx.New(t, (&fx.Chaining{}).Set(
		cmp, cmp.Listener.(lines.Componenter)))
	x, y := 0, 0
	fx_.Lines.Update(cmp, nil, func(e *lines.Env) {
		x, y, _, _ = cmp.Dim().Printable()
	})

	fx_.FireClick(x, y) // select the first item
	t.Eq(1, cmp.Listener.(*fx.Cmp).N(fx.NUpdate))
	t.Eq(0, reported)
}

func (s *AList) Reports_clicked_liner_item(t *T) {
	reported := -1
	cmp := &List{
		SelectableLiner: (&fx.SelectableLiner{}),
		Listener:        func(i int) { reported = i },
	}
	cmp.SelectableLiner.(*fx.SelectableLiner).II =
		[]string{"12", "1234", "123"}
	fx := fx.New(t, cmp)
	x, y := 0, 0
	fx.Lines.Update(cmp, nil, func(e *lines.Env) {
		x, y, _, _ = cmp.Dim().Printable()
	})

	fx.FireClick(x, y+1) // select the second item
	t.Eq(1, reported)
}

func (s *AList) Reports_selected_list_item(t *T) {
	reported := -1
	cmp := &List{
		Items: []string{"12", "1234", "123"},
		Listener: &fx.Cmp{
			ONUpdate: func(c *fx.Cmp, e *lines.Env, i interface{}) {
				reported = int(i.(Value))
			}},
	}
	fx_ := fx.New(t, (&fx.Chaining{}).Set(
		cmp, cmp.Listener.(lines.Componenter)))
	fx_.Lines.Focus(fx_.Root().(*fx.Chaining).CC[0])
	fx_.FireKeys(lines.Down, lines.Enter)
	t.Eq(0, reported)
}

func (s *AList) Reports_selected_liner_item(t *T) {
	reported := -1
	cmp := &List{
		SelectableLiner: (&fx.SelectableLiner{}),
		Listener:        func(i int) { reported = i },
	}
	cmp.SelectableLiner.(*fx.SelectableLiner).II =
		[]string{"12", "1234", "123"}
	fx := fx.New(t, cmp)

	fx.FireKeys(lines.Down, lines.Down, lines.Down, lines.Enter)
	t.Eq(2, reported)
}

func (s *AList) Scrolls_to_off_screen_list_items(t *T) {
	cmp := &List{Items: []string{"12", "3456", "789"}}
	fx := fx.Sized(t, 10, 2, cmp)
	t.Contains(fx.Screen(), "3456")
	t.Not.Contains(fx.Screen(), "789")

	fx.FireKeys(lines.Down, lines.Down, lines.Down)
	t.Contains(fx.Screen(), "789")
}

func (s *AList) Reports_scrolled_to_list_item(t *T) {
	reported := 1
	cmp := &List{
		Items:    []string{"12", "3456", "789"},
		Listener: func(i int) { reported = i },
	}
	fx := fx.Sized(t, 10, 2, cmp)
	t.Contains(fx.Screen(), "3456")
	t.Not.Contains(fx.Screen(), "789")

	fx.FireKeys(lines.Down, lines.Down, lines.Down)
	t.Contains(fx.Screen(), "789")
	fx.FireKey(lines.Enter)
	t.Eq(2, reported)
}

func (s *AList) Scrolls_to_off_screen_source_liner_item(t *T) {
	cmp := &List{SelectableLiner: (&fx.SelectableLiner{})}
	cmp.SelectableLiner.(*fx.SelectableLiner).II =
		[]string{"12", "3456", "789"}
	fx := fx.Sized(t, 10, 2, cmp)
	t.Contains(fx.Screen(), "3456")
	t.Not.Contains(fx.Screen(), "789")

	fx.FireKeys(lines.Down, lines.Down, lines.Down)
	t.Contains(fx.Screen(), "789")
}

func (s *AList) Reports_scrolled_to_source_liner_item(t *T) {
	reported := -1
	cmp := &List{
		SelectableLiner: (&fx.SelectableLiner{}),
		Listener:        func(i int) { reported = i },
	}
	cmp.SelectableLiner.(*fx.SelectableLiner).II =
		[]string{"12", "3456", "789"}
	fx := fx.Sized(t, 10, 2, cmp)

	fx.FireKeys(lines.Down, lines.Down, lines.Down, lines.Enter)
	t.Eq(2, reported)
}

type bugList struct {
	List
}

func (l *bugList) OnInit(e *lines.Env) {
	l.Items = fx.NStrings(20)
	l.List.OnInit(e)
}

// Scrolling_bug originated in the misconception that an as filling
// flagged rune fills until the next rune not flagged as filling while
// the implementation splits remaining space equally over the filler.
// Anyway its a good test checking if scrolling works.
func (s *AList) Scrolling_bug(t *T) {
	cmp := &bugList{}
	fx := fx.Sized(t, 5, 5, cmp)
	t.Eq(
		fx.Lines.Globals.Style(lines.Default),
		fx.Cells()[0][4].Style,
	)
	sl, dflt, hi := 0, lines.DefaultStyle, lines.DefaultStyle
	fx.Lines.Update(cmp, nil, func(e *lines.Env) {
		sl = cmp.ContentScreenLines()
		dflt = cmp.Globals().Style(lines.Default)
		hi = cmp.Globals().Style(lines.Highlight)
	})
	t.Eq(hi, fx.Lines.Globals.Style(lines.Default))
	t.Eq(dflt, fx.Lines.Globals.Style(lines.Highlight))
	t.Eq(
		fx.Lines.Globals.Style(lines.Default),
		fx.Cells()[0][4].Style,
	)
	for i := 0; i < sl; i++ {
		cc := fx.Cells().Column(4)
		for j, c := range cc {
			if j == i && !t.Eq(c.Style, hi) {
				fmt.Printf("hi: %d: %d\n", i, j)
			}
			if j != i && !t.Eq(c.Style, dflt) {
				fmt.Printf("dflt: %d: %d\n", i, j)
			}
		}
		fx.Lines.Update(cmp, nil, func(e *lines.Env) {
			t.Eq(i, cmp.Scroll.BarPosition())
			cmp.Scroll.Down()
		})
	}
}

func TestAList(t *testing.T) {
	t.Parallel()
	Run(&AList{}, t)
}
