// Copyright (c) 2022 Stephan Lukits. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"

	"github.com/slukits/lines"
	"github.com/slukits/lines/examples/demo"
)

type arrowsDemo struct {
	lines.Component
	demo.Demo
}

var arrowsTitle []rune = []rune("cursor-keys")

func (c *arrowsDemo) OnInit(e *lines.Env) {
	c.Init(c, e, arrowsTitle)
	c.Dim().SetWidth(32)
	fmt.Fprint(e.LL(1),
		"first line of content\n",
		"second line of content\n",
		"third line of content\n",
		"fourth line of content\n",
		"fifth line of content\n",
	)
}

func (c *arrowsDemo) OnFocus(e *lines.Env) {
	c.Demo.OnFocus(e)
	c.WriteTip("arrow-keys move cursor")
	c.SetCursor(0, 0, lines.BlockCursorBlinking)
}

func (c *arrowsDemo) OnFocusLost(e *lines.Env) {
	c.Demo.OnFocusLost(e)
	e.Lines.RemoveCursor()
}

func (c *arrowsDemo) OnKey(
	e *lines.Env, k lines.Key, m lines.ModifierMask,
) {
	line, column, cursorSet := c.CursorPosition()
	if !cursorSet {
		c.SetCursor(0, 0, lines.BlockCursorBlinking)
		return
	}
	_, _, w, h := c.ContentArea()
	switch k {
	case lines.Up:
		if line == 0 {
			return
		}
		c.SetCursor(line-1, column)
	case lines.Right:
		if column+1 == w {
			return
		}
		c.SetCursor(line, column+1)
	case lines.Down:
		if line+1 == h {
			return
		}
		c.SetCursor(line+1, column)
	case lines.Left:
		if column == 0 {
			return
		}
		c.SetCursor(line, column-1)
	default:
		c.Demo.OnKey(e, k, m)
	}
}
