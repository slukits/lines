// Copyright (c) 2022 Stephan Lukits. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

/*
Package edit chains two editable components one of which gets its content
from a source with a liner implementation while the other one's content
is directly written to it.  Both start out empty and can be
filled/edited by the user.
*/
package main

import (
	"fmt"

	"github.com/slukits/lines"
	"github.com/slukits/lines/examples/demo"
)

func main() {
	lines.Term(&app{}).WaitForQuit()
}

type app struct {
	lines.Component
	lines.Stacking
	demo.Demo
}

func (c *app) OnInit(e *lines.Env) {
	c.InitDemo(c, e, []rune("editor demo"))
	c.CC = append(c.CC, &doc{}, &editors{}, &evtReporter{})
}

func (c *app) Width() int { return 61 }

func (c *app) Height() int { return 20 }

type doc struct {
	lines.Component
}

const (
	left  = "left is an editor reporting directly to the component"
	right = "right is an editor reporting to the component's source"
	below = "below reported events of an editor input are displayed"
)

func (c *doc) OnInit(e *lines.Env) {
	fmt.Fprint(e.FG(0xEEEEEE).BG(0x666666), left+"\n"+right+"\n"+below)
	c.Dim().SetWidth(59).SetHeight(3)
}

type editors struct {
	lines.Component
	lines.Chaining
}

func (c *editors) OnInit(e *lines.Env) {
	c.CC = append(c.CC, &edt{}, &edtSrc{})
}

type edt struct {
	lines.Component
	demo.Demo
}

func (c *edt) OnInit(e *lines.Env) {
	c.InitDemo(c, e, []rune("direct-inactive"))
	c.FF.Set(lines.Editable)
}

type edtSrc struct {
	lines.Component
	demo.Demo
}

func (c *edtSrc) OnInit(e *lines.Env) {
	c.InitDemo(c, e, []rune("sourced-inactive"))
	c.Src = &lines.ContentSource{Liner: &editLiner{}}
}

type editLiner struct {
	cc []string
}

func (l *editLiner) Print(idx int, w *lines.EnvLineWriter) bool {
	if len(l.cc) <= idx || idx < 0 {
		return false
	}
	fmt.Fprintf(w, l.cc[idx])
	return idx+1 < len(l.cc)
}

type evtReporter struct {
	lines.Component
	demo.Demo
}

func (c *evtReporter) OnInit(e *lines.Env) {
	c.InitDemo(c, e, []rune("events"))
}

func (c *evtReporter) Width() int { return 57 }
