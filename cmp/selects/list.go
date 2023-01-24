// Copyright (c) 2022 Stephan Lukits. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package selects

import (
	"fmt"

	"github.com/slukits/lines"
)

// SelectableLiner is a scrollable liner which can also provide the
// information about a preferred width for its items.
type SelectableLiner interface {
	lines.ScrollableLiner

	// MaxWidth should return the preferable width for the provided items.
	MaxWidth() int
}

// List holds a list of given Items or uses a given ScrollableLiner as
// items-source.  These items are offers for selection whereas an
// selected item is reported to given Listener.
type List struct {
	component

	// Items are the items of a selection list.  Note items are
	// superseded by a scrollable liner.
	Items []string

	// Styler defines optionally the (highlighted) styles of given
	// Items.  Note a Styler is superseded by a scrollable liner
	// and defaults to Reversed back- and foreground as default style
	// and to default back- and foreground as highlight-style.
	Styler Styler

	// SelectableLiner is an alternative to define the list-elements
	// suited for big lists or lists with sophisticated styling of their
	// elements.
	SelectableLiner SelectableLiner

	// Listener is informed about a selected item
	Listener func(int)

	// close a items-layer
	close func(*lines.Lines)

	// focus keeps track of the menu-item which is currently hovered by
	// the mouse pointer.
	focus int
}

func (l *List) IsZero() bool {
	return len(l.Items) == 0 && l.SelectableLiner == nil
}

// OnInit prints the list items to the printable area.
func (l *List) OnInit(e *lines.Env) {
	l.focus = -1
	if l.IsZero() {
		return
	}
	l.FF.Set(lines.Focusable | lines.LinesFocusable)
	if l.SelectableLiner != nil {
		l.Src = &lines.ContentSource{Liner: l.SelectableLiner}
	}
	l.Globals().SetStyle(
		lines.Default, l.Globals().Style(lines.Default).Reverse())
}

func (l *List) OnLayout(e *lines.Env) (reflow bool) {
	if l.IsZero() {
		l.zeroPrint(e)
		return
	}
	_, _, w, h := l.ContentArea()
	if l.len() > h {
		l.FF.Set(lines.Scrollable)
	}
	if l.len() < h {
		top, _, bottom, _ := l.GapsLen()
		l.Dim().SetHeight(l.len() + top + bottom)
		reflow = true
	}
	width := l.maxWidth()
	if width < w {
		l.Dim().SetWidth(width)
	}
	if l.SelectableLiner == nil {
		l.print(0, e)
	}
	return reflow
}

func (c *List) maxWidth() int {
	_, right, _, left := c.GapsLen()
	if c.SelectableLiner != nil {
		return left + c.SelectableLiner.MaxWidth() + right
	}
	maxWdth := 0
	for _, i := range c.Items {
		if maxWdth >= len(i) {
			continue
		}
		maxWdth = len(i)
	}
	return left + maxWdth + right
}

func (l *List) len() int {
	if l.SelectableLiner != nil {
		return l.SelectableLiner.Len()
	}
	return len(l.Items)
}

func (l *List) styler_() func(int, bool) lines.Style {
	if l.Styler != nil {
		return l.Styler
	}
	return styler(l)
}

func (l *List) print(start int, e *lines.Env) {
	_, _, _, max := l.Dim().Printable()
	for i, itm := range l.Items[start:] {
		if i == max {
			break
		}
		if l.Styler != nil {
			fmt.Fprint(e.Sty(l.Styler(i+start, false)).LL(i), itm)
			continue
		}
		fmt.Fprint(e.LL(i), itm)
	}
}

func (l *List) zeroPrint(e *lines.Env) {
	top, right, bottom, left := l.GapsLen()
	l.Dim().SetWidth(right + len([]rune(NoItems)) + left).
		SetHeight(top + 1 + bottom)
	fmt.Fprint(e, NoItems)
}

// styler returns default items styler closure styler.  NOTE a styler
// call will panic outside a listener callback.
func styler(l *List) (styler func(int, bool) lines.Style) {
	var sty, hi lines.Style
	hi = l.Globals().Style(lines.Default)
	switch hi.AA() & lines.Reverse {
	case lines.Reverse:
		sty = hi.WithAA(hi.AA() &^ lines.Reverse)
	default:
		sty = hi.WithAA(hi.AA() | lines.Reverse)
	}
	styler = func(_ int, highlight bool) lines.Style {
		if highlight {
			return hi
		}
		return sty
	}
	return styler
}

// OnMove takes care of emphasizing the item hovered by the mouse
// cursor.
func (l *List) OnMove(e *lines.Env, x, y int) {
	if l.IsZero() {
		return
	}
	l.setFocus(e, y)
}

func (l *List) OnExit(e *lines.Env) {
	if l.IsZero() {
		return
	}
	l.setFocus(e, -1)
}

func (l *List) setFocus(e *lines.Env, idx int) {
	if l.focus >= 0 && l.focus != idx {
		l.resetFocus(e)
	}
	if idx > len(l.Items) || idx < 0 || l.focus == idx {
		return
	}
	fmt.Fprint(e.Sty(l.styler_()(idx, true)).LL(idx), l.Items[idx])
	l.focus = idx
}

func (l *List) resetFocus(e *lines.Env) {
	itm := l.Items[l.focus]
	fmt.Fprint(e.LL(l.focus).Sty(l.styler_()(l.focus, false)), itm)
	l.focus = -1
}

// OnClick reports the selected menu-item if any.
func (l *List) OnClick(e *lines.Env, x, y int) {
	if l.IsZero() {
		return
	}
	l.report(e, y)
}

func (l *List) report(e *lines.Env, idx int) {
	l.close(e.Lines)
	l.focus = -1
	if idx > len(l.Items) || idx < 0 {
		l.Listener(-1)
		return
	}
	l.Listener(idx)
}

type ModalList struct {
	List

	// MaxWidth is the maximum width which defaults to the list element
	// with maximum width
	MaxWidth int

	// MaxHeight is the maximum height which defaults to the
	MaxHeight int

	pos *lines.LayerPos
}

// OnOutOfBoundClick closes the menu-items
func (l *ModalList) OnOutOfBoundClick(e *lines.Env) bool {
	l.report(e, -1)
	return false
}

// OnOutOfBoundMove lets us remove the emphasis from the last emphasized
// menu item.
func (l *ModalList) OnOutOfBoundMove(e *lines.Env) bool {
	l.setFocus(e, -1)
	return false
}
