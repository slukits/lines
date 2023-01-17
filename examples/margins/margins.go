// Copyright (c) 2022 Stephan Lukits. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"

	"github.com/slukits/lines"
)

type cmp struct{ lines.Component }

func (c *cmp) OnInit(e *lines.Env) {
	c.Dim().SetWidth(20).SetHeight(1)
}

func (c *cmp) OnLayout(e *lines.Env) (reflow bool) {
	_, mr, _, ml := c.Dim().Margin()
	x, y, w, h := c.Dim().Printable()
	fmt.Fprintf(e, "%d:%d (%d,%d) %dx%d", ml, mr, x, y, w, h)
	return false
}

type app struct {
	lines.Component
	lines.Chaining
}

func (c *app) OnInit(e *lines.Env) {
	c.Dim().SetHeight(1)
	c.CC = append(c.CC, &cmp{}, &cmp{}, &cmp{})
}

func main() { lines.Term(&app{}).WaitForQuit() }
