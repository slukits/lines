// Copyright (c) 2022 Stephan Lukits. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"strings"

	"github.com/slukits/lines"
)

type example struct {
	lines.Component
	lines.Stacking
	explain       []string
	expNotFilling bool
	cmp           lines.Componenter
	dontFill      bool
}

func (c *example) OnInit(e *lines.Env) {
	c.CC = append(c.CC, &msg{
		txt: c.explain, notFilling: c.expNotFilling}, c.cmp)
	if !c.dontFill {
		c.CC = append(c.CC, &expFiller{})
	}
}

type msg struct {
	lines.Component
	notFilling bool
	txt        []string
}

func (c *msg) OnInit(e *lines.Env) {
	fmt.Fprint(e, strings.Join(c.txt, "\n"))
	if c.notFilling {
		c.Dim().SetHeight(len(c.txt) + 1)
	}
}

type expFiller struct{ lines.Component }
