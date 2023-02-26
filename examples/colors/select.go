// Copyright (c) 2022 Stephan Lukits. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"

	"github.com/slukits/lines"
)

type VSelection struct {
	lines.Component
	lines.Chaining
	Value      int
	fireUpdate func(lines.Componenter, interface{}, lines.Listener) error
}

func NewVSelection(
	label string,
	maxWidth int,
	default_ int,
	ii ...item,
) *VSelection {
	vs := &VSelection{Value: default_}
	vs.CC = append(vs.CC,
		&label_{lbl: label},
		&select_{
			ii:                ii,
			maxWidth:          maxWidth,
			default_:          default_,
			selectionListener: vs.setValue,
		},
	)
	return vs
}

func (c *VSelection) OnInit(e *lines.Env) {
	c.Dim().SetWidth(
		c.CC[0].(*label_).width() + c.CC[1].(*select_).width())
	c.Dim().SetHeight(3)
	c.fireUpdate = e.Lines.Update
}

func (c *VSelection) setValue(idx int) {
	if idx == -1 && c.CC[1].(*select_).hasDefault() {
		return
	}
	c.Value = idx
	c.fireUpdate(c.CC[1].(*select_), idx, nil)
}

type label_ struct {
	lines.Component
	lbl string
}

func (c *label_) OnInit(e *lines.Env) { fmt.Fprint(e, c.lbl) }

func (c *label_) width() int { return len(c.lbl) }

type select_ struct {
	lines.Component
	ii                []item
	maxWidth          int
	default_          int
	pos               *lines.LayerPos
	selectionListener func(int)
}

func (c *select_) reset(e *lines.Env) {
	if !c.hasDefault() {
		fmt.Fprint(e, lines.Filler+"▼")
		return
	}
	fmt.Fprint(e, c.ii[c.default_].label+lines.Filler+"▼")
}

func (c *select_) hasDefault() bool {
	return c.default_ >= 0 && c.default_ < len(c.ii)
}

func (c *select_) OnInit(e *lines.Env) {
	c.Dim().SetWidth(c.width())
	c.reset(e)
}

func (c *select_) OnClick(e *lines.Env, x, y int) {
	l := &list{
		ii:       c.ii,
		close:    c.close,
		listener: c.selectionListener,
	}
	c.pos = lines.NewLayerPos(
		c.Dim().X(), c.Dim().Y()+1,
		c.width(), len(c.ii),
	)
	if !c.hasDefault() {
		c.reset(e)
	}
	c.Layered(e, l, c.pos)
}

func (c *select_) OnUpdate(e *lines.Env, data interface{}) {
	if data.(int) == -1 {
		c.reset(e)
		return
	}
	label := c.ii[data.(int)].label
	fmt.Fprint(e, label+lines.Filler+"▼")
}

// close closes the menu m's modal layer of menu items.
func (c *select_) close(ll *lines.Lines) {
	ll.Update(c, nil, func(e *lines.Env) {
		c.RemoveLayer(e)
		c.pos = nil
	})
}

func (c *select_) width() int {
	maxWdth := 0
	for _, i := range c.ii {
		if maxWdth >= len(i.label)+2 {
			continue
		}
		maxWdth = len(i.label) + 2
	}
	if c.maxWidth == 0 || c.maxWidth >= maxWdth {
		return maxWdth
	}
	return c.maxWidth
}

type item struct {
	label string
	style lines.Style
}

// list represents a selection's items-list layer.
type list struct {
	lines.Component

	// close a items-layer
	close func(*lines.Lines)

	// ii are the items-labels
	ii []item

	// listener is updated with selected item
	listener func(int)

	// focus keeps track of the menu-item which is currently hovered by
	// the mouse pointer.
	focus int
}

// OnInit prints the list items to the printable area.
func (l *list) OnInit(e *lines.Env) {
	for i, itm := range l.ii {
		fmt.Fprint(e.Sty(itm.style).LL(i), itm.label)
	}
	l.focus = -1
}

// OnMove takes care of emphasizing the menu-item hoverd by the mouse
// cursor.
func (l *list) OnMove(e *lines.Env, x, y int) {
	if l.focus >= 0 && l.focus != y {
		l.resetFocus(e)
	}
	if y > len(l.ii) || y < 0 || l.focus == y {
		return
	}
	itm := l.ii[y]
	sty := itm.style.WithAA(itm.style.AA() &^ lines.Reverse)
	fmt.Fprint(e.Sty(sty).LL(y), l.ii[y].label)
	l.focus = y
}

// OnOutOfBoundMove lets us remove the emphasis from the last emphasized
// menu item.
func (l *list) OnOutOfBoundMove(e *lines.Env) bool {
	if l.focus >= 0 {
		l.resetFocus(e)
	}
	return false
}

func (l *list) resetFocus(e *lines.Env) {
	itm := l.ii[l.focus]
	fmt.Fprint(e.LL(l.focus).Sty(itm.style), itm.label)
	l.focus = -1
}

// OnClick reports the selected menu-item if any.
func (l *list) OnClick(e *lines.Env, x, y int) {
	if y > len(l.ii) || y < 0 {
		return
	}
	l.close(e.Lines)
	l.focus = -1
	l.listener(y)
}

// OnOutOfBoundClick closes the menu-items
func (l *list) OnOutOfBoundClick(e *lines.Env) bool {
	l.close(e.Lines)
	l.focus = -1
	l.listener(-1)
	return false
}
