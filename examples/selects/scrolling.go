// Copyright (c) 2022 Stephan Lukits. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"

	"github.com/slukits/lines"
	"github.com/slukits/lines/cmp/selects"
)

type scrolling struct {
	lines.Component
	lines.Stacking
}

func (c *scrolling) OnInit(e *lines.Env) {
	c.CC = append(c.CC, &scrollingSplits{})
}

type scrollingSplits struct {
	lines.Component
	lines.Chaining
}

func (c *scrollingSplits) OnInit(e *lines.Env) {
	c.Dim().SetWidth(20).SetHeight(5)
	c.CC = append(c.CC, &scrollingList{}, &simpleMsg{})
	c.CC[0].(*scrollingList).Listener = c.CC[1]
}

func (c *scrollingSplits) OnAfterInit(e *lines.Env) {
	e.Lines.Focus(c.CC[0])
}

type scrollingList struct{ selects.List }

func (l *scrollingList) OnInit(e *lines.Env) {
	for i := 0; i < 20; i++ {
		l.Items = append(l.Items, fmt.Sprintf("item %d", i+1))
	}
	l.List.OnInit(e)
}
