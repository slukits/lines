// Copyright (c) 2022 Stephan Lukits. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package selects

import (
	"fmt"

	"github.com/slukits/lines"
)

type Value int

// SelectableLiner is a scrollable liner which can also provide the
// information about a preferred width for its items.
type SelectableLiner interface {
	lines.ScrollableLiner

	// MaxWidth should return the preferable width for the provided items.
	MaxWidth() int
}

// Styler provides style information for how to style elements of an
// selection list.
type Styler func(idx int) lines.Style

// Highlighter maps a given style to its highlighted version.
type Highlighter func(lines.Style) lines.Style

// List holds a list of given Items or uses a given ScrollableLiner as
// items-source.  These items are offers for selection whereas an
// selected item is reported to given Listener.
type List struct {
	component

	// Items are the items of a selection list.  Note items are
	// superseded by a scrollable liner.
	Items []string

	// Styler defines optionally the styles of given Items.  Note a
	// Styler is superseded by a selectable liner. An Item's default
	// style is the reversed component's default style.
	Styler Styler

	// Highlighter maps optionally a given style to its highlighted
	// version.  Note a highlighter provided by a selectable liner
	// supersedes this highlighter.  The highlighted default style is a
	// component's default style.
	Highlighter Highlighter

	// SelectableLiner is an alternative to define the list-elements
	// suited for big lists or lists with sophisticated styling of their
	// elements.
	SelectableLiner SelectableLiner

	// Listener which is either a [lines.Componenter] or a function
	// taking the item index as argument (func(int)) is informed about a
	// selected item.  In case of an Componenter the index is to
	// OnUpdate reported, i.e. given component must implement OnUpdate
	// to receive the selected item index.
	Listener interface{}

	// MaxWidth is the number of screen columns a list instance should
	// occupy at most without gaps
	MaxWidth int

	// MaxHeight is the number of screen lines a list instance should
	// occupy at most without gaps
	MaxHeight int

	// focus keeps track of the menu-item which is currently hovered by
	// the mouse pointer.
	focus int

	onLayout bool
}

func (l *List) IsZero() bool {
	return len(l.Items) == 0 && l.SelectableLiner == nil
}

// OnInit sets basic features and default styles.
func (l *List) OnInit(e *lines.Env) {
	l.focus = -1
	if l.IsZero() {
		return
	}
	l.FF.Set(lines.Focusable | lines.LinesSelectable |
		lines.HighlightEnabled)
	if l.SelectableLiner != nil {
		l.Src = &lines.ContentSource{Liner: l.SelectableLiner}
		if _, ok := l.SelectableLiner.(lines.Highlighter); !ok {
			if l.Highlighter != nil {
				l.Globals().SetHighlighter(l.Highlighter)
			}
		}
	}
	l.Globals().SetStyle(
		lines.Highlight, l.Globals().Style(lines.Default).Reverse())
	l.Globals().SetStyle(
		lines.Default, l.Globals().Style(lines.Default).Reverse())
	if l.SelectableLiner == nil && l.Highlighter != nil {
		l.Globals().SetHighlighter(l.Highlighter)
	}
}

func (l *List) OnLayout(e *lines.Env) (reflow bool) {

	if l.onLayout {
		l.onLayout = false
		return false
	}
	l.onLayout = true

	if l.IsZero() {
		l.zeroPrint(e)
		return
	}
	_, _, w, h := l.ContentArea()
	if l.maxHeight() > h {
		if !l.FF.Has(lines.Scrollable) {
			l.FF.Set(lines.Scrollable)
		}
	} else {
		if l.FF.Has(lines.Scrollable) {
			l.FF.Delete(lines.Scrollable)
		}
	}
	if l.maxHeight() < h {
		top, _, bottom, _ := l.GapsLen()
		l.Dim().SetHeight(l.maxHeight() + top + bottom)
		reflow = true
	}
	width := l.Width()
	if width < w {
		l.Dim().SetWidth(width)
		reflow = true
	}
	if l.SelectableLiner == nil {
		l.print(e)
	}
	return reflow
}

func (l *List) OnFocusLost(e *lines.Env) {
	if l.IsZero() {
		return
	}
	l.LL.Focus.Reset()
}

func (l *List) maxHeight() int {
	h := l.len()
	if l.MaxHeight > 0 && l.MaxHeight < h {
		return l.MaxHeight
	}
	return h
}

func (c *List) Width() int {
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

func (l *List) print(e *lines.Env) {
	for i, itm := range l.Items {
		if l.Styler != nil {
			fmt.Fprint(e.Sty(l.Styler(i)).LL(i), itm)
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

func (l *List) OnEnter(e *lines.Env, x, y int) {
	if l.IsZero() {
		return
	}
	l.LL.Focus.AtCoordinate(y)
}

// OnMove takes care of emphasizing the item hovered by the mouse
// cursor by moving the line-focus to it.
func (l *List) OnMove(e *lines.Env, x, y int) {
	if l.IsZero() || x < 0 || y < 0 {
		return
	}
	if l.LL.Focus.Screen() == y {
		return
	}
	l.LL.Focus.AtCoordinate(y)
}

// OnExit removes the line focus from a list.
func (l *List) OnExit(e *lines.Env) {
	if l.IsZero() {
		return
	}
	l.LL.Focus.Reset()
}

// OnClick reports the selected menu-item if any.
func (l *List) OnClick(e *lines.Env, x, y int) {
	if l.IsZero() {
		return
	}
	l.report(e, y)
}

func (l *List) OnLineSelection(e *lines.Env, _, sl int) {
	l.report(e, sl)
}

func (l *List) report(e *lines.Env, y int) {
	if l.Listener == nil {
		return
	}
	if y > l.len() || y < 0 {
		l.toListener(e, -1)
		return
	}
	l.toListener(e, l.First()+y)
}

func (l *List) toListener(e *lines.Env, idx int) {
	switch lst := l.Listener.(type) {
	case func(int):
		lst(idx)
	case lines.Componenter:
		e.Lines.Update(lst, Value(idx), nil)
	}
}

func (l *List) len() int {
	if l.SelectableLiner != nil {
		return l.SelectableLiner.Len()
	}
	return len(l.Items)
}

type ModalList struct {
	List

	// close a items-layer
	close func(*lines.Lines)

	// MaxWidth is the maximum width which defaults to the list element
	// with maximum width
	MaxWidth int

	// MaxHeight is the maximum height which defaults to the
	MaxHeight int

	pos *lines.LayerPos
}

// OnOutOfBoundClick reports a zero selection.
func (l *ModalList) OnOutOfBoundClick(e *lines.Env) bool {
	l.close(e.Lines)
	l.focus = -1
	l.report(e, -1)
	return false
}

// OnOutOfBoundMove lets us remove the emphasis from the last emphasized
// menu item.
func (l *ModalList) OnOutOfBoundMove(e *lines.Env) bool {
	if l.IsZero() {
		return false
	}
	l.LL.Focus.Reset()
	return false
}
