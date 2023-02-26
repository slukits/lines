// Copyright (c) 2022 Stephan Lukits. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

import (
	"github.com/slukits/lines"
	"github.com/slukits/lines/examples/demo"
)

var coloringDemoString = []string{
	"// TestType helper for color scheme definitions.",
	"type TestType struct {",
	"	hello string",
	"}",
	"",
	"/* main go-function called on execution. */",
	"func main() {",
	" 	tt := TestType{hello: \"hello yourself\"}",
	" 	fmt.Printf(\"hello %T: %s\", tt, tt.hello)",
	"}",
}

func main() {
	lines.Term(&app{}).WaitForQuit()
}

type app struct {
	lines.Component
	lines.Stacking
	demo.Demo
}

func (c *app) OnInit(e *lines.Env) {
	c.CC = append(c.CC, &header{}, &settings{}, &save{})
	c.Init(c, e, []rune("text color picker demo"))
}

func (c *app) Width() int  { return 70 }
func (c *app) Height() int { return 22 }

type header struct {
	lines.Component
	lines.Chaining
}

func (c *header) OnInit(e *lines.Env) {
	sty := lines.DefaultStyle.WithAA(lines.Reverse)
	c.CC = append(c.CC,
		NewVSelection("color-sets: ", 0, 0,
			item{label: "monochrome", style: sty},
			item{label: "eight-colors", style: sty},
			item{label: "linux-tty", style: sty},
			item{label: "sixteen-colors", style: sty},
			item{label: "256-colors", style: sty},
			item{label: "true-color", style: sty},
		),
		NewVSelection("type: ", 0, 0,
			item{label: "dark", style: sty},
			item{label: "light", style: sty},
		),
	)
	c.Dim().SetHeight(3)
	lines.Print(c.Gaps(0).Horizontal.At(0).Filling(), ' ')
}

type save struct{ lines.Component }

func (c *save) OnInit(e *lines.Env) {
	c.Dim().SetHeight(2)
}
