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
	ee := &editors{reporter: &evtReporter{}}
	c.CC = append(c.CC, &doc{}, ee, ee.reporter)
}

func (c *app) Width() int { return 61 }

func (c *app) Height() int { return 22 }

type doc struct {
	lines.Component
}

const (
	left  = "left is an editor reporting directly to the component"
	right = "right is an editor reporting to the component's source"
	below = "below reported events of an editor input are displayed"
	exp   = "press <insert>-key to activate and <esc> to deactivate"
)

func (c *doc) OnInit(e *lines.Env) {
	fmt.Fprint(e.FG(0xEEEEEE).BG(0x666666),
		left+"\n"+right+"\n"+below+"\n"+exp)
	fmt.Fprint(c.Gaps(0).Bottom, "")
	c.Dim().SetWidth(59).SetHeight(5)
}

type editors struct {
	lines.Component
	lines.Chaining
	reporter lines.Componenter
}

func (c *editors) OnInit(e *lines.Env) {
	c.CC = append(c.CC, &edt{reporter: c.reporter},
		&edtSrc{reporter: c.reporter})
}

type edt struct {
	lines.Component
	demo.Demo
	reporter lines.Componenter
}

func (c *edt) OnInit(e *lines.Env) {
	c.InitDemo(c, e, []rune("direct-inactive"))
	c.FF.Set(lines.Editable)
}

func (c *edt) OnEdit(e *lines.Env, edt *lines.Edit) (suppress bool) {
	switch edt.Type {
	case lines.Resume:
		c.Title = []rune("direct-active")
		c.Double(e)
	case lines.Suspend:
		c.Title = []rune("direct-inactive")
		c.Double(e)
	}
	e.Lines.Update(c.reporter, &edit{source: "direct", Edit: edt}, nil)
	return
}

type edtSrc struct {
	lines.Component
	demo.Demo
	reporter lines.Componenter
}

func (c *edtSrc) OnInit(e *lines.Env) {
	c.InitDemo(c, e, []rune("sourced-inactive"))
	c.Src = &lines.ContentSource{Liner: &editLiner{
		cc:       []string{""},
		ll:       e.Lines,
		reporter: c.reporter,
	}}
}

func (c *edtSrc) OnEdit(e *lines.Env, edt *lines.Edit) (suppress bool) {
	switch edt.Type {
	case lines.Resume:
		c.Title = []rune("sourced-active")
		c.Double(e)
	case lines.Suspend:
		c.Title = []rune("sourced-inactive")
		c.Double(e)
	}
	return
}

type editLiner struct {
	cc       []string
	ll       *lines.Lines
	reporter lines.Componenter
}

func (l *editLiner) Print(idx int, w *lines.EnvLineWriter) bool {
	if len(l.cc) <= idx || idx < 0 {
		return false
	}
	fmt.Fprintf(w, l.cc[idx])
	return idx+1 < len(l.cc)
}

func (l *editLiner) OnEdit(edt *lines.Edit) (suppress bool) {
	l.ll.Update(l.reporter, &edit{source: "sourced", Edit: edt}, nil)
	return false
}

func (l *editLiner) IsFocusable(idx int) bool { return true }

func (l *editLiner) Len() int { return len(l.cc) }

type edit struct {
	*lines.Edit
	source string
}

func (edt *edit) String() string {
	evt := ""
	switch edt.Type {
	case lines.Resume:
		evt = "resume"
	case lines.Suspend:
		evt = "suspend"
	case lines.Ins:
		evt = "insert"
	default:
		evt = "unknown event"
	}
	return fmt.Sprintf("%s: %s at line: %d and column: %d rune: '%c'",
		edt.source, evt, edt.Line, edt.Cell, edt.Rune)
}

type evtReporter struct {
	lines.Component
	demo.Demo
	reported []string
}

func (c *evtReporter) OnInit(e *lines.Env) {
	c.InitDemo(c, e, []rune("edit-events-LIFO"))
}

func (c *evtReporter) OnUpdate(e *lines.Env, data interface{}) {
	scrLines := c.ContentScreenLines()
	if len(c.reported) > scrLines {
		c.reported = append([]string{data.(*edit).String()},
			c.reported[:scrLines-1]...)
	} else {
		c.reported = append([]string{data.(*edit).String()},
			c.reported...)
	}
	for i, s := range c.reported {
		fmt.Fprint(e.LL(i), s)
	}
}

func (c *evtReporter) Width() int { return 57 }
