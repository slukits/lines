// Copyright (c) 2022 Stephan Lukits. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"strings"

	"github.com/slukits/lines"
	"github.com/slukits/lines/examples/demo"
)

type titleFrame struct{ title []rune }

type gapper interface{ Gaps(int) *lines.GapsWriter }

func (f *titleFrame) frame(g gapper, e *lines.Env) {
	lines.Print(g.Gaps(0).Vertical.At(0).Filling(), '│')
	lines.Print(g.Gaps(0).Horizontal.At(0).Filling(), '─')
	fmt.Fprintf(g.Gaps(0).Corners, "┌┐┘└")
	lines.Print(g.Gaps(0).Top.At(1), f.title)
	lines.Print(g.Gaps(0).Top.At(1+len(f.title)).Filling(), '─')
}

type app struct {
	lines.Component
	lines.Stacking
	demo.Titled
}

var appTitle []rune = []rune("gaps demo")

func (c *app) OnInit(e *lines.Env) {
	c.Title = appTitle
	c.Default(c, e)
	r1 := &row{}
	r1.CC = []lines.Componenter{&simpleGap{}, &ttlTopCentered{},
		&ttlTopLeft{}, &ttlTopRight{}}
	r2 := &row{}
	r2.CC = []lines.Componenter{&ttlOthers{}, &nested{}, &gutterGap{}}
	c.CC = []lines.Componenter{r1, r2}
}

type row struct {
	lines.Component
	lines.Chaining
}

type simpleGap struct {
	lines.Component
	lines.Stacking
}

func (c *simpleGap) OnInit(e *lines.Env) {
	c.Gaps(0).AA(lines.Reverse)
	c.Gaps(0).Corners.AA(lines.Reverse)
	c.CC = []lines.Componenter{&simple{}}
	c.Dim().SetWidth(18).SetHeight(7)
}

type simple struct{ lines.Component }

func (c *simple) OnInit(e *lines.Env) {
	c.Dim().SetWidth(len("simple")).SetHeight(2)
	fmt.Fprint(e, "simple\nframe")
}

type ttlTopCentered struct {
	lines.Component
	lines.Stacking
	demo.Titled
}

func (c *ttlTopCentered) OnInit(e *lines.Env) {
	c.Title = []rune("centered")
	c.Default(c, e)
	c.CC = []lines.Componenter{&titled{label: "top-title"}}
	c.Dim().SetWidth(18).SetHeight(7)
}

type ttlTopLeft struct {
	lines.Component
	lines.Stacking
}

func (c *ttlTopLeft) OnInit(e *lines.Env) {
	c.CC = []lines.Componenter{&titled{label: "top-title"}}
	lines.Print(c.Gaps(0).Vertical.At(0).Filling(), '│')
	lines.Print(c.Gaps(0).Bottom.At(0).Filling(), '─')
	fmt.Fprint(c.Gaps(0).Corners, "┌┐┘└")
	fmt.Fprint(c.Gaps(0).Top, "left")
	lines.Print(c.Gaps(0).Top.At(len("left")).Filling(), '─')
	c.Dim().SetWidth(18).SetHeight(7)
}

type ttlTopRight struct {
	lines.Component
	lines.Stacking
}

func (c *ttlTopRight) OnInit(e *lines.Env) {
	c.CC = []lines.Componenter{&titled{label: "top-title"}}
	lines.Print(c.Gaps(0).Vertical.At(0).Filling(), '│')
	lines.Print(c.Gaps(0).Bottom.At(0).Filling(), '─')
	fmt.Fprint(c.Gaps(0).Corners, "┌┐┘└")
	lines.Print(c.Gaps(0).Top.At(0).Filling(), '─')
	lines.Print(c.Gaps(0).Top.At(1), []rune("right"))
	c.Dim().SetWidth(18).SetHeight(7)
}

type ttlOthers struct {
	lines.Component
	lines.Stacking
}

func (c *ttlOthers) OnInit(e *lines.Env) {
	c.CC = []lines.Componenter{&titled{label: "others-titled"}}
	c.Dim().SetWidth(18).SetHeight(7)
	fmt.Fprint(c.Gaps(0).Corners, "┌┐┘└")
	lines.Print(c.Gaps(0).Top.At(0).Filling(), '─')
	lines.Print(c.Gaps(0).Left.At(0).Filling(), '│')
	lines.Print(c.Gaps(0).Left.At(1).AA(lines.Reverse), []rune("mid"))
	lines.Print(c.Gaps(0).Left.At(4).Filling(), '│')
	lines.Print(c.Gaps(0).Bottom.At(0).Filling(), '─')
	lines.Print(c.Gaps(0).Bottom.At(1).AA(lines.Reverse), []rune("right"))
	lines.Print(c.Gaps(0).Right.At(0).AA(lines.Reverse), []rune("top"))
	lines.Print(c.Gaps(0).Right.At(3).Filling(), '│')
}

type titled struct {
	lines.Component
	label string
}

func (c *titled) OnInit(e *lines.Env) {
	c.Dim().SetWidth(len(c.label)).SetHeight(1)
	fmt.Fprint(e, c.label)
}

type nested struct {
	lines.Component
	lines.Stacking
	idx int
}

var nestTitles [][]rune = [][]rune{
	[]rune("ne"), []rune("st"), []rune("ed")}

func (c *nested) OnInit(e *lines.Env) {
	c.Dim().SetWidth(18 - 2*c.idx).SetHeight(7 - 2*c.idx)
	fmt.Fprint(c.Gaps(0).Corners, "┌┐┘└")
	lines.Print(c.Gaps(0).Vertical.At(0).Filling(), '│')
	lines.Print(c.Gaps(0).Bottom.At(0).Filling(), '─')
	var indent []rune
	if c.idx > 0 {
		indent = []rune(strings.Repeat("─", c.idx))
		lines.Print(c.Gaps(0).Top.At(0), indent)
	}
	lines.Print(c.Gaps(0).Top.At(len(indent)), nestTitles[c.idx])
	lines.Print(c.Gaps(0).Top.At(len(indent)+len(nestTitles[c.idx])).
		Filling(), '─')
	if c.idx+1 < 3 { // nest
		c.CC = append(c.CC, &nested{idx: c.idx + 1})
	}
}

type gutterGap struct {
	lines.Component
	lines.Stacking
}

func (c *gutterGap) OnInit(e *lines.Env) {
	for _, i := range []int{0, 1} {
		c.Gaps(i).Left.AA(lines.Reverse)
		fmt.Fprint(c.Gaps(i).TopLeft.AA(lines.Reverse), " ")
		fmt.Fprint(c.Gaps(i).BottomLeft.AA(lines.Reverse), " ")
		if i == 0 {
			lines.Print(c.Gaps(0).Top.At(0).AA(lines.Reverse), ' ')
			lines.Print(c.Gaps(0).Bottom.At(0).AA(lines.Reverse), ' ')
		}
	}
	c.CC = []lines.Componenter{&gutter{}}
	c.Dim().SetWidth(18).SetHeight(7)
	lines.Print(c.Gaps(0).Left.At(0), '●')
	lines.Print(c.Gaps(0).Left.At(2), '☺')
	lines.Print(c.Gaps(0).Left.At(4), '')
}

type gutter struct{ lines.Component }

func (c *gutter) OnInit(e *lines.Env) {
	c.Dim().SetWidth(len("left-hand")).SetHeight(2)
	fmt.Fprint(e, "left-hand\ngutter")
}

func main() {
	lines.Term(&app{}).WaitForQuit()
}
