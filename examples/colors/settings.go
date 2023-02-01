// Copyright (c) 2022 Stephan Lukits. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"

	"github.com/slukits/lines"
)

type settings struct {
	lines.Component
	lines.Stacking
}

func (c *settings) OnInit(e *lines.Env) {
	c.CC = append(c.CC,
		&setting{label: "background: ", height: 2},
		&setting{label: "foreground: ", height: 2},
		&setting{label: "keyword: ", height: 2},
		&setting{label: "operator: ", height: 2},
		&setting{label: "comment: ", height: 2},
		&setting{label: "identifier: ", height: 2},
		&setting{label: "type: ", height: 2},
		&setting{label: "error: ", height: 1},
	)
}

type setting struct {
	lines.Component
	label  string
	height int
}

func (c *setting) OnInit(e *lines.Env) {
	fmt.Fprint(e, c.label)
	c.Dim().SetHeight(c.height)
}
