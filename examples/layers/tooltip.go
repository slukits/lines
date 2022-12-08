// Copyright (c) 2022 Stephan Lukits. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"

	"github.com/slukits/lines"
	"github.com/slukits/lines/examples/demo"
)

type toolTip struct {
	lines.Component
	demo.Demo

	// currentLine keeps track of the line having its tool-tip shown
	currentLine int
}

var toolTipTitle []rune = []rune("tool-tip-demo")

// OnInit prints the content lines to the tool-tip demo component.
func (c *toolTip) OnInit(e *lines.Env) {
	c.Init(c, e, toolTipTitle)
	line := func(i int) string {
		i++
		switch i {
		case 1, 2, 3:
			return []string{"", "1st", "2nd", "3rd"}[i]
		default:
			return fmt.Sprintf("%dth", i)
		}
	}
	for i := 0; i < c.Height()-2; i++ {
		fmt.Fprintf(e.LL(i), "hover %s line for tip", line(i))
	}
	c.currentLine = 0
}

// OnMove keeps track of tool tip lines hovered by the mouse cursor and
// sets according tool-tip layers.
func (c *toolTip) OnMove(e *lines.Env, x, y int) {
	_, _, aw, ah := c.Dim().Printable()
	if c.currentLine == y-1 && x > 0 && x < aw-1 {
		return
	}
	if y < 1 || y > ah-2 || x < 1 || x > c.Dim().Width()-2 {
		c.reset(e)
		return
	}

	c.currentLine = y - 1
	cx, cy, _, _ := c.ContentArea()
	layer := &tip{line: y}
	c.Layered(e, layer,
		lines.NewLayerPos(cx+x, cy+y, layer.width(), 1).SetZ(1000))
}

func (c *toolTip) reset(e *lines.Env) {
	if c.currentLine == -1 {
		return
	}
	c.RemoveLayer(e)
	c.currentLine = -1
}

// uncomment if tool tip should be removed if hovered by the mouse
// func (c *toolTip) OnExit(e *lines.Env) {
// 	c.reset(e)
// }

type tip struct {
	lines.Component
	line int
}

const content = "component overlapping tool-tip for line %d"

func (c *tip) OnInit(e *lines.Env) {
	c.AA(lines.Reverse)
	fmt.Fprintf(e, content, c.line)
}

func (c *tip) width() int { return len(fmt.Sprintf(content, c.line)) }
