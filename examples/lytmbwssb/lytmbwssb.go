// Copyright (c) 2022 Stephan Lukits. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"

	"github.com/slukits/lines"
)

type App struct {
	lines.Component
	lines.Stacking
}

func (c *App) OnInit(e *lines.Env) {
	c.CC = append(c.CC, &MessageBar{}, &WS{}, &Statusbar{})
}

type MessageBar struct{ lines.Component }

func (c *MessageBar) OnInit(e *lines.Env) {
	c.Dim().SetHeight(3)
}

func (c *MessageBar) OnLayout(e *lines.Env) {
	fmt.Fprintf(e, "message bar: %d,%d: %dx%d",
		c.Dim().X(), c.Dim().Y(), c.Dim().Width(), c.Dim().Height())
}

type Statusbar struct{ lines.Component }

func (c *Statusbar) OnInit(e *lines.Env) {
	c.Dim().SetHeight(3)
}

func (c *Statusbar) OnLayout(e *lines.Env) {
	fmt.Fprintf(e, "statusbar: %d,%d: %dx%d",
		c.Dim().X(), c.Dim().Y(), c.Dim().Width(), c.Dim().Height())
}

type WS struct {
	lines.Component
	lines.Chaining
}

func (c *WS) OnInit(e *lines.Env) {
	c.CC = append(c.CC, &Panel{}, &Panel{})
	c.FF.AddRecursive(lines.Focusable)
	e.EE.MoveFocus(c.CC[0])
}

type Panel struct{ lines.Component }

func (c *Panel) OnInit(e *lines.Env) {
	c.Mod(lines.Tailing)
}

func (c *Panel) OnLayout(e *lines.Env) {
	fmt.Fprintf(e, "panel: %d,%d: %dx%d",
		c.Dim().X(), c.Dim().Y(), c.Dim().Width(), c.Dim().Height())
}

func (c *Panel) OnFocusLost(e *lines.Env) {
	fmt.Fprint(e, "lost focus")
}

func (c *Panel) OnFocus(e *lines.Env) {
	fmt.Fprint(e, "got focus")
}

func (c *Panel) OnClick(e *lines.Env, x int, y int) {
	fmt.Fprintf(e, "clicked (%d,%d)", x, y)
}

func main() { lines.New(&App{}).Listen() }
