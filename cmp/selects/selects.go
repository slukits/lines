// Copyright (c) 2022 Stephan Lukits. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

/* selects.go provides APIs for the *items* see items.go */

package selects

import "github.com/slukits/lines"

type component = lines.Component
type chaining = lines.Chaining
type stacking = lines.Stacking

// NoLabel is a selection list's default label if none is set.
const NoLabel = "no-label:"

// NoItems is a selection list's default item if none are set.
const NoItems = "no items"

type Orientation string

const (
	Drop Orientation = "▼"
	Up   Orientation = "▲"
)

const NoDefault = -1

type DropDown struct {
	items

	fireUpdate func(lines.Componenter, interface{}, lines.Listener) error

	// Items of a selection list
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

	// DefaultItem is the index of the item which is selected if no item
	// is selected.  Note set DefaultItem to NoDefault if a zero input
	// is allowed.
	DefaultItem int

	// MaxWidth sets the maximum items-label width which defaults to the
	// width of the widest item.
	MaxWidth int

	MinWidth int

	// MaxHeight may be used to restrict the hight of dropped selection
	// list.
	MaxHeight int

	// Orientation indicates if DropDown actually drops down or if it
	// drops up.
	Orientation Orientation

	value int
}

func (c *DropDown) dfltItems() []string { return c.Items }
func (c *DropDown) dfltStyler() Styler  { return c.Styler }
func (c *DropDown) dfltHighlighter() Highlighter {
	return c.Highlighter
}
func (c *DropDown) dfltItem() int      { return c.DefaultItem }
func (c *DropDown) dfltMaxWidth() int  { return c.MaxWidth }
func (c *DropDown) dfltMinWidth() int  { return c.MinWidth }
func (c *DropDown) dfltMaxHeight() int { return c.MaxHeight }
func (c *DropDown) dfltOrientation() Orientation {
	return c.Orientation
}

func (c *DropDown) OnInit(e *lines.Env) {
	c.dd = c
	c.value = c.DefaultItem
	if c.Orientation == "" {
		c.Orientation = Drop
	}
	if c.MaxWidth < 0 {
		c.MaxWidth = 0
	}
	if c.Items == nil {
		c.Items = []string{NoItems}
	}
	c.listener = c
	c.Dim().SetHeight(1).SetWidth(c.width(true))
	c.fireUpdate = e.Lines.Update
	c.items.OnInit(e)
}

func (c *DropDown) newItems() (ii *items) {
	if c.MaxWidth < 0 {
		c.MaxWidth = 0
	}
	ii = &items{dd: c}
	if c.Items == nil {
		c.Items = []string{NoItems}
	}
	ii.listener = c
	return ii
}

func (c *DropDown) OnUpdate(e *lines.Env, data interface{}) {
	idx := int(data.(Value))
	if idx == -1 && c.hasDefault() {
		return
	}
	c.value = idx
	c.items.OnUpdate(e, idx)
}

func (c *DropDown) Value() int { return c.value }

// DropDownHrz is a horizontally labeled drop down selection list.
// While the zero-value is ready to use it's not very useful.  Usually
// you will set at least its Label and its Items.
type DropDownHrz struct {
	component
	chaining

	fireUpdate func(lines.Componenter, interface{}, lines.Listener) error

	// Label of a horizontally labeled selection list
	Label string

	// Items of a horizontally labeled selection list
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

	// DefaultItem is the index of the item which is selected if no item
	// is selected.  Note set DefaultItem to NoDefault if a zero input
	// is allowed.
	DefaultItem int

	// MaxWidth sets the maximum items-label width which defaults to the
	// width of the widest item.
	MaxWidth int

	MinWidth int

	// MaxHeight may be used to restrict the hight of dropped selection
	// list.
	MaxHeight int

	// Orientation indicates if DropDown actually drops down or if it
	// drops up.
	Orientation Orientation

	value int
}

func (c *DropDownHrz) dfltItems() []string { return c.Items }
func (c *DropDownHrz) dfltStyler() Styler  { return c.Styler }
func (c *DropDownHrz) dfltHighlighter() Highlighter {
	return c.Highlighter
}
func (c *DropDownHrz) dfltItem() int      { return c.DefaultItem }
func (c *DropDownHrz) dfltMaxWidth() int  { return c.MaxWidth }
func (c *DropDownHrz) dfltMinWidth() int  { return c.MinWidth }
func (c *DropDownHrz) dfltMaxHeight() int { return c.MaxHeight }
func (c *DropDownHrz) dfltOrientation() Orientation {
	return c.Orientation
}

func (c *DropDownHrz) OnInit(e *lines.Env) {
	c.value = c.DefaultItem
	if c.Orientation == "" {
		c.Orientation = Drop
	}
	c.CC = append(c.CC, c.newLabel(), c.newItems())
	c.Dim().SetHeight(1).SetWidth(
		c.CC[0].(*label).width() + c.CC[1].(*items).width(true))
	c.fireUpdate = e.Lines.Update
}

func (c *DropDownHrz) newLabel() (l *label) {
	l = &label{lbl: c.Label}
	if l.lbl == "" {
		l.lbl = NoLabel
	}
	return l
}

func (c *DropDownHrz) newItems() (ii *items) {
	if c.MaxWidth < 0 {
		c.MaxWidth = 0
	}
	ii = &items{dd: c}
	if c.Items == nil {
		c.Items = []string{NoItems}
	}
	ii.listener = c
	return ii
}

func (c *DropDownHrz) OnUpdate(e *lines.Env, data interface{}) {
	idx := int(data.(Value))
	if idx == -1 && c.CC[1].(*items).hasDefault() {
		return
	}
	c.value = idx
	c.fireUpdate(c.CC[1].(*items), idx, nil)
}

func (c *DropDownHrz) Value() int { return c.value }

// DropDownVrt is a vertically labeled drop down selection list.
// While the zero-value is ready to use it's not very useful.  Usually
// you will set at least its Label and its Items.
type DropDownVrt struct {
	component
	stacking

	fireUpdate func(lines.Componenter, interface{}, lines.Listener) error

	// Label of a horizontally labeled selection list
	Label string

	// Items of a horizontally labeled selection list
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

	// DefaultItem is the index of the item which is selected if no item
	// is selected.  Note set DefaultItem to NoDefault if a zero input
	// is allowed.
	DefaultItem int

	// MaxWidth sets the maximum items-label width which defaults to the
	// width of the widest item.
	MaxWidth int

	MinWidth int

	// MaxHeight may be used to restrict the hight of dropped selection
	// list.
	MaxHeight int

	// Orientation indicates if DropDown actually drops down or if it
	// drops up.
	Orientation Orientation

	value int
}

func (c *DropDownVrt) dfltItems() []string { return c.Items }
func (c *DropDownVrt) dfltStyler() Styler  { return c.Styler }
func (c *DropDownVrt) dfltHighlighter() Highlighter {
	return c.Highlighter
}
func (c *DropDownVrt) dfltItem() int      { return c.DefaultItem }
func (c *DropDownVrt) dfltMaxWidth() int  { return c.MaxWidth }
func (c *DropDownVrt) dfltMinWidth() int  { return c.MinWidth }
func (c *DropDownVrt) dfltMaxHeight() int { return c.MaxHeight }
func (c *DropDownVrt) dfltOrientation() Orientation {
	return c.Orientation
}

func (c *DropDownVrt) OnInit(e *lines.Env) {
	c.value = c.DefaultItem
	if c.Orientation == "" {
		c.Orientation = Drop
	}
	label, items := c.newLabel(), c.newItems()
	if c.MinWidth == 0 {
		c.MinWidth = label.width()
	}
	width := items.width(true)
	if c.MinWidth > width {
		width = c.MinWidth
	}
	c.Dim().SetHeight(2).SetWidth(width)
	c.CC = append(c.CC, label, items)
	c.fireUpdate = e.Lines.Update
}

func (c *DropDownVrt) newLabel() (l *label) {
	l = &label{lbl: c.Label}
	if l.lbl == "" {
		l.lbl = NoLabel
	}
	return l
}

func (c *DropDownVrt) newItems() (ii *items) {
	if c.MaxWidth < 0 {
		c.MaxWidth = 0
	}
	ii = &items{dd: c}
	if c.Items == nil {
		c.Items = []string{NoItems}

	}
	ii.listener = c
	return ii
}

func (c *DropDownVrt) OnUpdate(e *lines.Env, data interface{}) {
	idx := int(data.(Value))
	if idx == -1 && c.CC[1].(*items).hasDefault() {
		return
	}
	c.value = idx
	c.fireUpdate(c.CC[1].(*items), idx, nil)
}

func (c *DropDownVrt) Value() int { return c.value }
