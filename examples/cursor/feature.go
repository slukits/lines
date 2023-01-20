// Copyright (c) 2022 Stephan Lukits. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"

	"github.com/slukits/lines"
	"github.com/slukits/lines/examples/demo"
)

type featureDemo struct {
	lines.Component
	demo.Demo
}

const hellip = '…'

var featureTitle []rune = []rune("highlighted sourced line|cell-focus-feature")

func enum(i int) string {
	switch i {
	case 1:
		return "st"
	case 2:
		return "nd"
	case 3:
		return "rd"
	default:
		return "th"
	}
}

func (c *featureDemo) OnInit(e *lines.Env) {
	c.InitDemo(c, e, featureTitle)
	c.Dim().SetWidth(64).SetHeight(10)
	c.FF.Set(lines.CellFocusable | lines.HighlightEnabled)
	fmt.Fprint(c.Gaps(1).Top, "")   // create gap for focused tip
	fmt.Fprint(c.Gaps(1).Left, "")  // create gap for left overflow
	fmt.Fprint(c.Gaps(1).Right, "") // create gap for right overflow
	c.Src = &lines.ContentSource{Liner: newLiner()}
}

func (c *featureDemo) OnFocus(e *lines.Env) {
	c.Demo.OnFocus(e)
	c.WriteTip("every second line is focusable with arrow keys")
}

func (c *featureDemo) OnFocusLost(e *lines.Env) {
	c.Demo.OnFocusLost(e)
	e.Lines.RemoveCursor()
}

func (c *featureDemo) OnLineOverflowing(_ *lines.Env, left, right bool) {
	lines.Print(c.Gaps(1).Left.At(c.LL.Focus.Screen()), ' ')
	lines.Print(c.Gaps(1).Right.At(c.LL.Focus.Screen()), ' ')
	if left {
		lines.Print(c.Gaps(1).Left.At(c.LL.Focus.Screen()), hellip)
	}
	if right {
		lines.Print(c.Gaps(1).Right.At(c.LL.Focus.Screen()), hellip)
	}
}

func (c *featureDemo) OnLineFocusLost(_ *lines.Env, _, sIdx int) {
	lines.Print(c.Gaps(1).Left.At(sIdx), ' ')
	lines.Print(c.Gaps(1).Right.At(sIdx), ' ')
}

type focusableLiner struct {
	cc           []string
	notFocusable map[int]bool
}

func newLiner() *focusableLiner {
	lr := focusableLiner{notFocusable: map[int]bool{}}
	for i := 0; i < 20; i++ {
		if i%2 == 0 {
			lr.notFocusable[i] = true
		}
		if i != 0 && i%3 == 0 {
			lr.cc = append(lr.cc,
				fmt.Sprintf(
					"%d%s line of focusable lines with arrow-keys "+
						"movable with some line overflowing",
					i+1, enum(i+1),
				))
			continue
		}
		lr.cc = append(lr.cc, fmt.Sprintf(
			"%d%s line of focusable lines with arrow-keys movable cursor",
			i+1, enum(i+1),
		))
	}
	return &lr
}

// Print prints the line with given index idx to given line writer w and
// returns true if there are lines with a greater index left to write.
func (l *focusableLiner) Print(idx int, w *lines.EnvLineWriter) bool {
	if len(l.cc) <= idx || idx < 0 {
		return false
	}
	fmt.Fprintf(w, l.cc[idx])
	return idx+1 < len(l.cc)
}

// Len returns the total number of content lines a liner implementation
// can provide to its associated component.
func (l *focusableLiner) Len() int { return len(l.cc) }

// IsFocusable returns true iff the line with given index idx is
// focusable.
func (l *focusableLiner) IsFocusable(idx int) bool {
	return !l.notFocusable[idx]
}

// Highlighted indicates if focusable lines are highlighted if focused.
// And in case they are highlighted if they should be trimmed
// highlighted.
func (l *focusableLiner) Highlighted() (bool, bool) {
	return true, false
}
