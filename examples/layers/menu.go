// Copyright (c) 2022 Stephan Lukits. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"strings"

	"github.com/slukits/lines"
	"github.com/slukits/lines/examples/demo"
)

// menuDemo separates the menu demo from other demos and stacks a
// menu bar as well as a printing area menuTxt.
type menuDemo struct {
	lines.Component
	demo.Demo
	lines.Stacking
}

var menuTitle []rune = []rune("menu-demo")

// OnInit sets up the component structure of the menu-demo.
func (c *menuDemo) OnInit(e *lines.Env) {
	c.InitDemo(c, e, menuTitle)
	menuBar, menuTxt := &menuBar{}, &menuTxt{}
	menuBar.reporter = menuTxt
	c.CC = append(c.CC, menuBar, menuTxt)
}

// A menuBar instance chains the menu-bar's menus.
type menuBar struct {
	lines.Component

	// Chaining chains the menu-bar's menus.
	lines.Chaining

	// reporter provides access to the menu-demo's text area.
	reporter lines.Componenter
}

const (
	file = "file"
	help = "help"
)

// OnInit sets up the menus.  To achive the effect that the last menu is
// flushed right a empty component is used.  Component embedder are by
// default filling, i.e. use up all remaining space.
func (c *menuBar) OnInit(e *lines.Env) {
	c.Dim().SetHeight(1)
	c.AA(lines.Reverse)
	c.CC = []lines.Componenter{
		&menu{label: file, reporter: c.reporter},
		&filler{}, // flush help menu to the right
		&menu{label: help, reporter: c.reporter},
	}
}

// OnLayout keeps an open menu sticked to its menu-bar in case of layout
// changes.
func (mb *menuBar) OnLayout(e *lines.Env) bool {
	for i := 0; i < 3; i++ {
		switch i {
		case 0, 2:
			if m := mb.CC[i].(*menu); m.pos != nil {
				e.Lines.Update(m, nil, func(e *lines.Env) {
					m.pos.MoveTo(m.Dim().X(), m.Dim().Y()+1)
				})
			}
		}
	}
	return false
}

// menu is a menu-bar's menu component and may be associated with the
// layer of its menu items.
type menu struct {
	lines.Component

	// reporter can print a selected menu item to the menu-demo text
	// area
	reporter lines.Componenter

	// label is the name of the menu shown in the menu bar
	label string

	// pos if not nil is the position of the menu's items layer.
	// Keeping the position here allows a menu to move its items layer
	// in case of an layout change
	pos *lines.LayerPos
}

const separator = lines.Filler

var mnItems = map[string]string{
	file: "open\nclose\n" + separator + "\nquit",
	help: "browse\nabout",
}

// OnInit sets the given menu m's width to its label width to not
// consume more space in the menu bar as necessary and prints its label
// into the menu bar.
func (m *menu) OnInit(e *lines.Env) {
	m.Dim().SetWidth(len([]rune(m.label)))
	fmt.Fprint(e, m.label)
}

// OnClick creates the modal layer of menu itemes of given menu m.
func (m *menu) OnClick(e *lines.Env, x, y int) {
	ii := &items{
		ii:       strings.Split(mnItems[m.label], "\n"),
		close:    m.close,
		reporter: m.reporter,
	}
	m.pos = lines.NewLayerPos(
		m.Dim().X(), m.Dim().Y()+1,
		ii.Width(), len(ii.ii),
	)
	m.Layered(e, ii, m.pos)
}

// close closes the menu m's modal layer of menu items.
func (m *menu) close(ll *lines.Lines) {
	ll.Update(m, nil, func(e *lines.Env) {
		m.RemoveLayer(e)
		m.pos = nil
	})
}

// filler is a component taking up all remaining space of its parent
// component.
type filler struct{ lines.Component }

// items represents a menu's items-layer
type items struct {
	lines.Component

	// close a items-layer
	close func(*lines.Lines)

	// ii are the items-labels
	ii []string

	// reporter is updated with selcted items
	reporter lines.Componenter

	// focus keeps track of the menu-item which is currently hovered by
	// the mouse pointer.
	focus int
}

// OnInit prints the menu items to the menu-items printable area.
func (ii *items) OnInit(e *lines.Env) {
	for i, itm := range ii.ii {
		if itm == separator {
			lines.Print(e.LL(i).At(0).Filling(), '─')
			continue
		}
		fmt.Fprint(e.LL(i), itm)
	}
	ii.focus = -1
}

// OnMove takes care of emphasizing the menu-item hoverd by the mouse
// cursor.
func (ii *items) OnMove(e *lines.Env, x, y int) {
	if ii.focus >= 0 && ii.focus != y {
		ii.resetFocus(e)
	}
	if y > len(ii.ii) || y < 0 || ii.ii[y] == separator || ii.focus == y {
		return
	}
	fmt.Fprint(e.LL(y).BG(lines.Salmon).AA(lines.Reverse|lines.Bold),
		ii.ii[y])
	ii.focus = y
}

// OnOutOfBoundMove lets us remove the emphasis from the last emphasized
// menu item.
func (ii *items) OnOutOfBoundMove(e *lines.Env) bool {
	if ii.focus >= 0 {
		ii.resetFocus(e)
	}
	return false
}

func (ii *items) resetFocus(e *lines.Env) {
	fmt.Fprint(
		e.LL(ii.focus).Sty(ii.Globals().Style(lines.Default)),
		ii.ii[ii.focus],
	)
	ii.focus = -1
}

// OnClick reports the selected menu-item if any.
func (ii *items) OnClick(e *lines.Env, x, y int) {
	if y > len(ii.ii) || y < 0 || ii.ii[y] == separator {
		return
	}
	ii.close(e.Lines)
	ii.focus = -1
	e.Lines.Update(
		ii.reporter, fmt.Sprintf("selected: '%s'", ii.ii[y]), nil)
}

// OnOutOfBoundClick closes the menu-items
func (ii *items) OnOutOfBoundClick(e *lines.Env) bool {
	ii.close(e.Lines)
	return false
}

// Width returns the number of runes of the longest menu-item label.
func (ii items) Width() (n int) {
	for _, i := range ii.ii {
		if len(i) <= n {
			continue
		}
		n = len(i)
	}
	return n
}

// menuTxt represents the menu-demo's text-area.
type menuTxt struct{ lines.Component }

// OnUpdate prints to the menu-demo's text-area.
func (m *menuTxt) OnUpdate(e *lines.Env, selectedItem interface{}) {
	fmt.Fprint(e.LL(1), lines.Filler+selectedItem.(string)+lines.Filler)
}
