// Copyright (c) 2022 Stephan Lukits. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"

	"github.com/slukits/lines"
)

type app struct {
	lines.Component
	lines.Stacking
}

func (c *app) OnInit(_ *lines.Env) {
	c.Dim().SetHeight(5).SetWidth(len(msgWaiting))
	msg := &msg{}
	c.CC = []lines.Componenter{msg, &lineFocus{msg: msg}}
}

const (
	msgWaiting  = "waiting for selecting key"
	msgNothing  = "no line is focused"
	msgFocused  = "on line focused: %d"
	msgSelected = "on line selected: %d"
)

// msg has the role of a message bar reporting what's currently selected
// in the selecting component.
type msg struct{ lines.Component }

func (c *msg) OnInit(e *lines.Env) {
	c.Dim().SetHeight(2)
	fmt.Fprint(
		e.BG(lines.Green).FG(lines.Blue),
		msgWaiting,
	)
}

func (c *msg) OnUpdate(e *lines.Env) {
	fmt.Fprint(e, e.Evt.(*lines.UpdateEvent).Data.(string))
}

// lineFocus component has a list of lines of which a few can be
// focused and subsequently selected.
type lineFocus struct {
	lines.Component
	msg *msg
}

const (
	fs  = "line %d: %s"
	yes = "focusable"
	no  = "not focusable"
)

func (c *lineFocus) OnInit(e *lines.Env) {
	c.FF.Add(lines.LinesSelectable)
	c.Dim().SetHeight(3)
	if err := e.Lines.Focus(c); err != nil {
		panic(err)
	}
	fmt.Fprintf(e.LL(0), fs, 0, yes)
	fmt.Fprintf(e.LL(1), fs, 1, no)
	fmt.Fprintf(e.LL(2), fs, 2, no)
	fmt.Fprintf(e.LL(3), fs, 3, no)
	fmt.Fprintf(e.LL(4), fs, 4, yes)
	fmt.Fprintf(e.LL(5), fs, 5, no)
	fmt.Fprintf(e.LL(6), fs, 6, yes)
	for _, i := range []int{1, 2, 3, 5} {
		c.LL.By(i).Flag(lines.NotFocusable)
	}
}

func (c *lineFocus) OnLineFocus(e *lines.Env, idx int) {
	if c.LL.Focus.Current() < 0 {
		if err := e.Lines.Update(c.msg, msgNothing, nil); err != nil {
			panic(err)
		}
		return
	}
	err := e.Lines.Update(
		c.msg, fmt.Sprintf(msgFocused, c.LL.Focus.Current()), nil)
	if err != nil {
		panic(err)
	}
}

func (c *lineFocus) OnLineSelection(e *lines.Env, idx int) {
	if c.LL.Focus.Current() < 0 {
		if err := e.Lines.Update(c.msg, msgNothing, nil); err != nil {
			panic(err)
		}
		return
	}
	err := e.Lines.Update(
		c.msg, fmt.Sprintf(msgSelected, c.LL.Focus.Current()), nil)
	if err != nil {
		panic(err)
	}
}

func main() {
	lines.Term(&app{}).WaitForQuit()
}
