// Copyright (c) 2022 Stephan Lukits. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"

	"github.com/slukits/lines"
	"github.com/slukits/lines/examples/demo"
	"github.com/slukits/lines/internal/api"
)

type clickDemo struct {
	lines.Component
	demo.Demo
}

var clickTitle []rune = []rune("left-click")

func (c *clickDemo) OnInit(e *lines.Env) {
	c.InitDemo(c, e, clickTitle)
	c.Dim().SetWidth(32)
}

func (c *clickDemo) OnFocus(e *lines.Env) {
	c.Demo.OnFocus(e)
	c.WriteTip("click to position cursor")
}

func (c *clickDemo) OnFocusLost(e *lines.Env) {
	c.Demo.OnFocusLost(e)
	e.Lines.RemoveCursor()
	c.Reset(lines.All)
}

func (c *clickDemo) OnClick(e *lines.Env, x, y int) {
	_, _, cw, ch := c.ContentArea()
	lTop, _, _, lLeft := c.GapsLen()
	if x-lLeft < 0 || x-lLeft >= cw || y-lTop < 0 || y-lTop > ch {
		return
	}
	c.SetCursor(y-lTop, x-lLeft, api.BlockCursorSteady)
	c.Reset(lines.All)
	c.labelCursor(e, x-lLeft, y-lTop, cw, ch)
}

func (c *clickDemo) labelCursor(e *lines.Env, cx, cy, cw, ch int) {
	rel := fmt.Sprintf("r(l,c):(%d,%d)", cy, cx)
	lTop, _, _, lLeft := c.GapsLen()
	abs := fmt.Sprintf("a(x,y):(%d,%d)", c.Dim().X()+lLeft+cx,
		c.Dim().Y()+lTop+cy)
	switch cy {
	case 0:
		fmt.Fprintf(e.LL(cy+1), "%s%s,%s%[1]s", lines.Filler, abs, rel)
	case ch - 1:
		fmt.Fprintf(e.LL(ch-2), "%s%s,%s%[1]s", lines.Filler, abs, rel)
	default:
		if cx+len(abs) < cw {
			lines.Print(e.LL(cy-1).At(cx+1), []rune(abs))
		} else {
			lines.Print(e.LL(cy-1).At(cw-len(abs)), []rune(abs))
		}
		if cx+len(rel) < cw {
			lines.Print(e.LL(cy+1).At(cx+1), []rune(rel))
		} else {
			lines.Print(e.LL(cy+1).At(cw-len(rel)), []rune(rel))
		}
	}
}
