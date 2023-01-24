// Copyright (c) 2022 Stephan Lukits. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package selects

import "github.com/slukits/lines"

type component = lines.Component
type chaining = lines.Chaining

// NoLabel is a selection list's default label if none is set.
const NoLabel = "no-label:"

// NoItems is a selection list's default item if none are set.
const NoItems = "no items"

const Drop = "▼"

const NoDefault = -1

// Styler provides style information for how to style elements of an
// selection list.
type Styler func(idx int, highlight bool) lines.Style

// Horizontal is a horizontally labeled drop down selection list.  While
// the zero-value is ready to use it's not very useful.  Usually you
// will set at least its Label and its Items.
type Horizontal struct {
	component
	chaining

	fireUpdate func(lines.Componenter, interface{}, lines.Listener) error

	// Label of a horizontally labeled selection list
	Label string

	// Items of a horizontally labeled selection list
	Items []string

	// Styler returns for each item-index its corresponding style
	// respectively its highlighted style.
	Styler Styler

	// DefaultItem is the index of the item which is selected if no item
	// is selected.  Note set DefaultItem to NoDefault if a zero input
	// is allowed.
	DefaultItem int

	// MaxWidth sets the maximum items-label width which defaults to the
	// width of the widest item.
	MaxWidth int

	value int
}

func (c *Horizontal) OnInit(e *lines.Env) {
	c.value = c.DefaultItem
	c.CC = append(c.CC, c.newLabel(), c.newItems())
	c.Dim().SetHeight(1).SetWidth(
		c.CC[0].(*label).width() + c.CC[1].(*items).width(true))
	c.fireUpdate = e.Lines.Update
}

func (c *Horizontal) newLabel() (l *label) {
	l = &label{lbl: c.Label}
	if l.lbl == "" {
		l.lbl = NoLabel
	}
	return l
}

func (c *Horizontal) newItems() (ii *items) {
	if c.MaxWidth < 0 {
		c.MaxWidth = 0
	}
	ii = &items{
		ii:          c.Items,
		styler:      c.Styler,
		defaultItem: c.DefaultItem,
		maxWidth:    c.MaxWidth,
	}
	if ii.ii == nil {
		ii.ii = []string{NoItems}
	}
	ii.listener = c.setValue
	return ii
}

func (c *Horizontal) setValue(idx int) {
	if idx == -1 && c.CC[1].(*items).hasDefault() {
		return
	}
	c.value = idx
	c.fireUpdate(c.CC[1].(*items), idx, nil)
}

func (c *Horizontal) Value() int { return c.value }
