// Copyright (c) 2022 Stephan Lukits. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"strconv"

	"github.com/slukits/lines"
	"github.com/slukits/lines/cmp/selects"
)

type simple struct {
	lines.Component
	lines.Stacking
}

func (c *simple) OnInit(e *lines.Env) {
	c.CC = append(c.CC, &simpleSplits{})
}

type simpleSplits struct {
	lines.Component
	lines.Chaining
}

func (c *simpleSplits) OnInit(e *lines.Env) {
	c.Dim().SetWidth(20).SetHeight(4)
	c.CC = append(c.CC, &simpleList{}, &simpleMsg{})
	c.CC[0].(*simpleList).Listener = c.CC[1]
}

func (c *simpleSplits) OnAfterInit(e *lines.Env) {
	e.Lines.Focus(c.CC[0])
}

type simpleList struct{ selects.List }

func (l *simpleList) OnInit(e *lines.Env) {
	for i := 0; i < 4; i++ {
		l.Items = append(l.Items, fmt.Sprintf("item %d", i+1))
	}
	l.List.OnInit(e)
}

type simpleMsg struct{ lines.Component }

func (l *simpleMsg) OnInit(e *lines.Env) {
	fmt.Fprint(e.LL(0), lines.Filler+"no  item"+lines.Filler)
	fmt.Fprint(e.LL(1), lines.Filler+"selected"+lines.Filler)
	l.Dim().SetHeight(2)
}

func (l *simpleMsg) OnUpdate(e *lines.Env, data interface{}) {
	v := strconv.Itoa(int(data.(selects.Value)) + 1)
	fmt.Fprint(e.LL(0), lines.Filler+"item "+v+lines.Filler)
	fmt.Fprint(e.LL(1), lines.Filler+"selected"+lines.Filler)
}
