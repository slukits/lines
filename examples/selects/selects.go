// Copyright (c) 2022 Stephan Lukits. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"

	"github.com/slukits/lines"
	"github.com/slukits/lines/examples/demo"
)

func main() {
	lines.Term(&selectsDemo{}).WaitForQuit()
}

type selectsDemo struct {
	lines.Component
	lines.Chaining
	titled *demo.Titled
}

var title = []rune(" selects demo ")

func (c *selectsDemo) OnInit(e *lines.Env) {
	c.titled = (&demo.Titled{Gapper: c, Title: title}).Single(e)
	c.CC = append(c.CC, &display{}, &menuBar{})
	c.Dim().SetWidth(50).SetHeight(16)
}

// OnAfterInit can set it self as display to the menu
func (c *selectsDemo) OnAfterInit(e *lines.Env) {
	c.CC[1].(*menuBar).menu().display = c.CC[0]
	c.CC[1].(*menuBar).quitter().display = c.CC[0]
}

type display struct {
	lines.Component
	lines.Stacking
}

func (c *display) OnInit(e *lines.Env) {
	c.CC = append(c.CC, BLANK)
}

func (c *display) OnUpdate(e *lines.Env, data interface{}) {
	c.CC[0] = data.(lines.Componenter)
}

var BLANK = &blank{}

type blank struct{ lines.Component }

var label = "placeholder for example"

func (c *blank) OnInit(e *lines.Env) {
	c.Dim().SetWidth(len(label)).SetHeight(1)
	fmt.Fprint(e, label)
}
