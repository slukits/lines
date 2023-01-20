// Copyright (c) 2022 Stephan Lukits. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

import (
	"github.com/slukits/lines"
	"github.com/slukits/lines/examples/demo"
)

type stacked struct {
	lines.Component
	lines.Chaining
	demo.Demo
	pos   *lines.LayerPos
	first bool
}

var stackedTitle []rune = []rune("stacked-demo")

// OnInit creates the dummy off-screen components which are then
// layered with the stacking layers.
func (c *stacked) OnInit(e *lines.Env) {
	c.InitDemo(c, e, stackedTitle)
	c.CC = append(c.CC, &layered{}, &layered{}, &layered{}, &layered{})
	c.first = true
}

// OnLayout initializes the layers in dependance of the calculated
// layout of the stacked-demo after the first layout.  On further
// layouts it is only made sure that the stacked layers stay relative to
// the stacked-demo at the same palce in case the over-all layout
// changes.
func (c *stacked) OnLayout(e *lines.Env) bool {
	x, y, _, _ := c.ContentArea()
	if c.first {
		for i := 0; i < 4; i++ {
			e.Lines.Update(c.CC[i], nil, func(i int) func(*lines.Env) {
				return func(e *lines.Env) { c.addLayer(e, i, x, y) }
			}(i))
		}
		c.first = false
		return false
	}
	c.reposition(e)
	return false
}

func (c *stacked) addLayer(e *lines.Env, idx int, x, y int) {
	var pos *lines.LayerPos
	var color lines.Color
	switch idx {
	case 0:
		pos, color = lines.NewLayerPos(x+7, y, 9, 3), lines.Red
	case 1:
		pos, color = lines.NewLayerPos(x+1, y+2, 9, 3), lines.Blue
	case 2:
		pos, color = lines.NewLayerPos(x+13, y+2, 9, 3), lines.Yellow
	case 3:
		pos, color = lines.NewLayerPos(x+7, y+4, 9, 3), lines.Green
	}
	c.CC[idx].(*layered).pos = pos
	c.CC[idx].(*layered).Layered(
		e, &lyr{color: color, focus: c.focus(idx)}, pos)
}

// reposition keeps the layers relative to the stacked-demo component
// at the same place.
func (c *stacked) reposition(e *lines.Env) {
	x, y, _, _ := c.ContentArea()
	for i := 0; i < 4; i++ {
		lrd := c.CC[i].(*layered)
		switch i {
		case 0:
			e.Lines.Update(lrd, nil, func(i int, lrd *layered) func(*lines.Env) {
				return func(_ *lines.Env) { lrd.pos.MoveTo(x+7, y) }
			}(i, lrd))
		case 1:
			e.Lines.Update(lrd, nil, func(i int, lrd *layered) func(*lines.Env) {
				return func(_ *lines.Env) { lrd.pos.MoveTo(x+1, y+2) }
			}(i, lrd))
		case 2:
			e.Lines.Update(lrd, nil, func(i int, lrd *layered) func(*lines.Env) {
				return func(_ *lines.Env) { lrd.pos.MoveTo(x+13, y+2) }
			}(i, lrd))
		case 3:
			e.Lines.Update(lrd, nil, func(i int, lrd *layered) func(*lines.Env) {
				return func(_ *lines.Env) { lrd.pos.MoveTo(x+7, y+4) }
			}(i, lrd))
		}
	}
}

// focus moves layer with given dummy-index to the top by adjusting its
// z-axis-level.
func (c *stacked) focus(idx int) func() {
	return func() {
		for i := 0; i < 4; i++ {
			lrd := c.CC[i].(*layered)
			if i == idx {
				lrd.pos.SetZ(10)
				continue
			}
			lrd.pos.SetZ(i)
		}
	}
}

type layered struct {
	lines.Component
	pos *lines.LayerPos
}

// OnInit sets the layerd-dummy off-screen.
func (c *layered) OnInit(e *lines.Env) {
	c.Dim().SetWidth(0)
}

type lyr struct {
	lines.Component
	color lines.Color
	focus func()
}

// OnInit set the background color.
func (c *lyr) OnInit(e *lines.Env) { c.BG(c.color) }

// OnClick moves the clicked layer to the top.
func (c *lyr) OnClick(e *lines.Env, x, y int) { c.focus() }
