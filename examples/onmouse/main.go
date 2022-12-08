// Copyright (c) 2022 Stephan Lukits. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"time"

	"github.com/slukits/lines"
	"github.com/slukits/lines/examples/frame"
)

func main() {
	lines.Term(&app{}).WaitForQuit()
}

type app struct {
	lines.Component
	lines.Chaining
	frame.Titled
}

func (c *app) OnInit(e *lines.Env) {
	c.Dim().SetWidth(64).SetHeight(24)
	c.Title = []rune("mouse-observer")
	c.Default(c, e)
	c.CC = append(c.CC, &eventsAggregates{}, &report{})
}

func (c *app) OnAfterInit(_ *lines.Env) {
	c.CC[0].(*eventsAggregates).CC[0].(*events).reporter = c.CC[1]
	c.CC[0].(*eventsAggregates).CC[0].(*events).aggregates =
		c.CC[0].(*eventsAggregates).CC[1]
}

type eventsAggregates struct {
	lines.Component
	lines.Stacking
}

func (c *eventsAggregates) OnInit(e *lines.Env) {
	c.CC = append(c.CC, &events{}, &aggregates{})
}

type events struct {
	lines.Component
	frame.Titled
	reporter   lines.Componenter
	last       time.Time
	aggregates lines.Componenter
}

func (c *events) OnInit(e *lines.Env) {
	c.Title = []rune(" create events ")
	c.Default(c, e)
}

func (c *events) OnMouse(e *lines.Env, b lines.ButtonMask, x, y int) {
	m := e.Evt.(lines.MouseEventer).Mod()
	if c.last.IsZero() {
		e.Lines.Update(
			c.reporter,
			fmt.Sprintf("%d(%d): (%d,%d), 0ms", b, m, x, y),
			nil,
		)
		c.last = e.Evt.When()
		return
	}
	e.Lines.Update(
		c.reporter,
		fmt.Sprintf("%d(%d): (%d,%d), %v",
			b, m, x, y, e.Evt.When().Sub(c.last)),
		nil,
	)
	c.last = e.Evt.When()
}

func (c *events) OnMove(e *lines.Env, x, y int) {
	ox, oy := e.Evt.(*lines.MouseMove).Origin()
	data := fmt.Sprintf("move: (%d,%d) ~> (%d,%d)", ox, oy, x, y)
	e.Lines.Update(c.aggregates, data, nil)
}

func (c *events) OnClick(e *lines.Env, x, y int) {
	data := fmt.Sprintf("click: (%d,%d)", x, y)
	e.Lines.Update(c.aggregates, data, nil)
}

func (c *events) OnContext(e *lines.Env, x, y int) {
	data := fmt.Sprintf("context: (%d,%d)", x, y)
	e.Lines.Update(c.aggregates, data, nil)
}

func (c *events) OnDrag(e *lines.Env, b lines.ButtonMask, x, y int) {
	ox, oy := e.Evt.(*lines.MouseDrag).Origin()
	data := fmt.Sprintf("drag(%d): (%d,%d) ~> (%d,%d)", b, ox, oy, x, y)
	e.Lines.Update(c.aggregates, data, nil)
}

func (c *events) OnDrop(e *lines.Env, b lines.ButtonMask, x, y int) {
	data := fmt.Sprintf("drop(%d): (%d,%d)", b, x, y)
	e.Lines.Update(c.aggregates, data, nil)
}

type aggregates struct {
	lines.Component
	frame.Titled
	ee []string
}

func (c *aggregates) OnInit(e *lines.Env) {
	c.Title = []rune(" aggregates ")
	c.Default(c, e)
}

func (c *aggregates) OnUpdate(e *lines.Env, event interface{}) {
	_, _, _, h := c.ContentArea()
	c.ee = append([]string{event.(string)}, c.ee...)
	if len(c.ee) >= h {
		c.ee = c.ee[:h]
	}
	for i, evt := range c.ee {
		fmt.Fprint(e.LL(i), evt)
	}
}

type report struct {
	lines.Component
	frame.Titled
	ee []string
}

func (c *report) OnInit(e *lines.Env) {
	c.Title = []rune(" report events ")
	c.Default(c, e)
}

func (c *report) OnUpdate(e *lines.Env, event interface{}) {
	_, _, _, h := c.ContentArea()
	c.ee = append([]string{event.(string)}, c.ee...)
	if len(c.ee) >= h {
		c.ee = c.ee[:h]
	}
	for i, evt := range c.ee {
		fmt.Fprint(e.LL(i), evt)
	}
}
