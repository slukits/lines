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
	c.CC = append(c.CC, blankDefault, &menuBar{})
	c.Dim().SetWidth(50).SetHeight(16)
}

// OnAfterInit can set it self as display to the menu
func (c *selectsDemo) OnAfterInit(e *lines.Env) {
	c.CC[1].(*menuBar).menu().display = c
	c.CC[1].(*menuBar).quitter().display = c
}

func (c *selectsDemo) OnUpdate(e *lines.Env, data interface{}) {
	c.CC[0] = data.(lines.Componenter)
	e.Lines.Redraw()
}

var blankDefault = &blank{}

type blank struct{ lines.Component }

var label = "placeholder for example"

func (c *blank) OnInit(e *lines.Env) {
	c.Dim().SetHeight(1)
	fmt.Fprint(e, lines.Filler+label+lines.Filler)
}
