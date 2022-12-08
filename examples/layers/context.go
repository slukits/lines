// Copyright (c) 2022 Stephan Lukits. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"

	"github.com/slukits/lines"
	"github.com/slukits/lines/examples/demo"
	"github.com/slukits/lines/examples/frame"
)

// context associates the context-demo content-area with a context menu
// layer which is drag and dropable.
type context struct {
	lines.Component
	lines.Stacking
	demo.Demo
}

var contextTitle []rune = []rune("context-demo")

// OnInit sets up the context-area component.
func (c *context) OnInit(e *lines.Env) {
	c.Init(c, e, contextTitle)
	c.CC = append(c.CC, &contextArea{})
}

type contextArea struct{ lines.Component }

// OnContext associates the context area with the context menu layer.
func (c *contextArea) OnContext(e *lines.Env, x, y int) {
	pos := lines.NewLayerPos(
		c.Dim().X()+x, c.Dim().Y()+y,
		len(cntxItems[0])+2, 6,
	)
	c.Layered(
		e, &contextMenu{close: c.close, reporter: c, pos: pos}, pos)
}

func (c *contextArea) close(ll *lines.Lines) {
	ll.Update(c, nil, func(e *lines.Env) {
		c.RemoveLayer(e)
	})
}

// OnUpdate reports a selected context menu item.
func (m *contextArea) OnUpdate(e *lines.Env, selectedItem interface{}) {
	fmt.Fprint(e.LL(2), lines.Filler+selectedItem.(string)+lines.Filler)
}

// contextMenu represents the layer showing a modal context menu.
type contextMenu struct {
	lines.Component

	// Titled frames the context menu items
	frame.Titled

	// close since a context menu is modal it is the only component which
	// can trigger its closeing.  The functionality is provided by
	// associated component.
	close func(ll *lines.Lines)

	// focus keeps track of the contet menu item currently hovered by
	// the mouse cursor.
	focus int

	// reporter updates are triggered for reporting selected context
	// menu item.
	reporter lines.Componenter

	// pos is the layer position of a conetxt menu allowing to move it
	// around.
	pos *lines.LayerPos

	// dragDelta is the distance of drag-starting x-coordinate to the
	// context-menu's origing x-coordinate.
	dragDelta int
}

var cntxItems = []string{
	"context-func-1",
	"context-func-2",
	"context-func-3",
	"context-func-4",
}

// OnInit prints the menu items to the context menu layer's lines.
func (c *contextMenu) OnInit(e *lines.Env) {
	c.AA(lines.Reverse)
	c.Title = []rune("context-menu")
	c.Default(c, e)
	for i, mi := range cntxItems {
		fmt.Fprint(e.LL(i), mi)
	}
	c.focus, c.dragDelta = -1, -1
}

// OnOutOfBoundClick we close the context menu.
func (c *contextMenu) OnOutOfBoundClick(e *lines.Env) bool {
	c.close(e.Lines)
	return false
}

// OnOutOfBoundMove is used to reset the last emphasized menu-item.
func (c *contextMenu) OnOutOfBoundMove(e *lines.Env) bool {
	if c.focus >= 0 {
		c.resetFocus(e)
	}
	return false
}

func (c *contextMenu) resetFocus(e *lines.Env) {
	fmt.Fprint(
		e.LL(c.focus).Sty(c.Globals().Style(lines.Default)),
		cntxItems[c.focus],
	)
	c.focus = -1
}

// OnDrag starting in context menu's title-line moves the context menu
// to draged coordinates.  Note we also increase the context-menu
// layer's z-level to overlay the layers of the stacked-demo.
func (c *contextMenu) OnDrag(e *lines.Env, b lines.ButtonMask, x, y int) {
	ox, oy := e.Evt.(*lines.MouseDrag).Origin()
	px, py, _, _ := c.Dim().Printable()
	if !c.inDrag() && (b != lines.Primary || py != oy) {
		return
	}
	if !c.inDrag() {
		c.dragDelta = ox - px
		c.pos.SetZ(100)
	}
	c.pos.MoveTo(x-c.dragDelta, y)
}

func (c *contextMenu) inDrag() bool { return c.dragDelta >= 0 }

// OnDrop resets the drag-delta and context-menu layer's the z-level of
// the last drag.
func (c *contextMenu) OnDrop(e *lines.Env, b lines.ButtonMask, x, y int) {
	c.dragDelta = -1
	c.pos.SetZ(0)
}

// OnMove is leveraged to emphasize context-menu items which are hoverd
// by the mouse.
func (c *contextMenu) OnMove(e *lines.Env, x, y int) {
	if c.focus >= 0 && c.focus != y-1 {
		c.resetFocus(e)
	}
	// account for gaps
	if y > len(cntxItems) || y < 1 || c.focus == y-1 {
		return
	}
	c.focus = y - 1
	fmt.Fprint(e.LL(c.focus).BG(lines.Cornsilk).
		AA(lines.Reverse|lines.Bold), cntxItems[c.focus])
}

// OnClick reports clicked context menu item.
func (c *contextMenu) OnClick(e *lines.Env, x, y int) {
	if y > len(cntxItems)+1 || y < 1 { // account for gaps
		return
	}
	c.close(e.Lines)
	c.focus = -1
	e.Lines.Update(
		c.reporter, fmt.Sprintf("'%s'", cntxItems[y-1]), nil)
}
