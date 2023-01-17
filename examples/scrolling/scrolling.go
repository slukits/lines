// Copyright (c) 2022 Stephan Lukits. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

/*
Package scrolling chains to components one of which gets its content
from a source with a liner implementation while the other one's content
is directly written to it.  Both have identical content and every line
with an even index is focusable.  Both components should behave exactly
the same.  These two scrolling components are displayed in two variants:
with gaps and without gaps.
*/
package main

import (
	"fmt"

	"github.com/slukits/ints"
	"github.com/slukits/lines"
)

func main() {
	lines.Term(&app{}).WaitForQuit()
}

type app struct {
	lines.Component
	lines.Stacking
}

func (c *app) OnInit(e *lines.Env) {
	// create nested stacked components
	c.CC = []lines.Componenter{
		&messageBar{}, &scrolling{}, &scrolling{gaps: true}}
}

func (c *app) OnAfterInit(e *lines.Env) {
	// calculate size to have all content horizontally and vertically
	// centered on the terminal screen.
	mb, sc, scg := c.CC[0].(*messageBar), c.CC[1].(*scrolling),
		c.CC[2].(*scrolling)
	c.Dim().SetHeight(mb.height() + sc.height() + scg.height())
	c.Dim().SetWidth(ints.Max(mb.width(), sc.width(), scg.width()))
	e.Lines.Focus(sc)
}

// messageBar displays how to use this example once it is executed.
type messageBar struct{ lines.Component }

var mm []string = []string{
	"PgUp/PgDn:" + lines.Filler + "scroll",
	"up/down:" + lines.Filler + "focus line",
	"click:" + lines.Filler + "focus component",
}

func (c *messageBar) OnInit(e *lines.Env) {
	for i := 0; i < 3; i++ { // print initial content
		fmt.Fprint(e.LL(i), mm[i])
	}
	c.Dim().SetHeight(c.height()) // set static height
}

func (c *messageBar) width() int {
	w := 0
	for _, m := range mm {
		if len(m) <= w {
			continue
		}
		w = len(m)
	}
	return w + 6
}

func (c *messageBar) height() int { return 3 }

// scrolling chains the components whose features should be demonstrated
// in two variants: scrolling/line-focusing with gaps and without gaps.
type scrolling struct {
	lines.Component
	lines.Chaining
	gaps bool
}

func (c *scrolling) OnInit(e *lines.Env) {
	c.FF.Set(lines.Focusable)
	c.Dim().SetHeight(c.height())
	if !c.gaps {
		frame(c, []rune("no gaps-trimmed highlight"), true)
	} else {
		frame(c, []rune("with gaps"), true)
	}
}

type cmpidx int

const (
	drcIdx cmpidx = iota
	scrIdx
)

func (c *scrolling) cmp(idx cmpidx) lines.Componenter {
	if len(c.CC) == 0 {
		c.CC = []lines.Componenter{
			&direct{gaps: c.gaps},
			&sourced{gaps: c.gaps},
		}
	}
	return c.CC[idx]
}

func (c *scrolling) height() int {
	if !c.gaps {
		return 7
	}
	return 9
}

func (c *scrolling) width() int {
	return c.cmp(drcIdx).(*direct).width() +
		c.cmp(scrIdx).(*sourced).width()
}

// direct is the scrolling component with selectable lines who has its
// content directly written to it (see direct.OnInit).
type direct struct {
	lines.Component
	gaps bool
}

var ll = []string{" 1st", " 2nd", " 3rd", " 4th", " 5th", " 6th",
	" 7th", " 8th", " 9th", "10th", "11th", "12th", "13th", "14th",
	"15th", "16th", "17th", "18th", "19th", "20th"}

func (c *direct) OnInit(e *lines.Env) {
	c.FF.Set(lines.LinesSelectable | lines.Scrollable)
	if c.gaps {
		frame(c, directTitle, false)
	}
	for i, l := range ll {
		fmt.Fprint(e.LL(i), lines.Filler+l+lines.Filler)
		if i%2 != 0 {
			continue
		}
		c.LL.By(i).Flag(lines.NotFocusable)
	}

	if !c.gaps {
		c.FF.Set(lines.TrimmedHighlightEnabled)
	} else {
		c.FF.Set(lines.HighlightEnabled)
	}
	c.Dim().SetWidth(c.width())
}

var directTitle = []rune("direct")

// frame given component g with given title centered in the top.
func frame(g Gapper, title []rune, fillTitle bool) {
	lines.Print(g.Gaps(0).Bottom.At(0).Filling(), '─')
	lines.Print(g.Gaps(0).Vertical.At(0).Filling(), '│')
	fmt.Fprintf(g.Gaps(0).Corners, "┌┐┘└")
	if !fillTitle {
		lines.Print(g.Gaps(0).Top.At(0), title)
		return
	}
	lines.Print(g.Gaps(0).Top.At(0).Filling(), '─')
	lines.Print(g.Gaps(0).Top.At(1), title)
	lines.Print(g.Gaps(0).Top.At(len(title)+1).Filling(), '─')
}

type Gapper interface{ Gaps(int) *lines.GapsWriter }

func (c *direct) OnFocusLost(_ *lines.Env) {
	if c.gaps {
		c.Gaps(0).AA(lines.ZeroStyle)
		c.LL.Focus.Reset()
		return
	}
	c.LL.AA(c.Globals().AA(lines.Default))
	c.LL.Focus.Reset()
}

func (c *direct) OnFocus(_ *lines.Env) {
	if c.gaps {
		c.Gaps(0).AA(lines.Bold)
		return
	}
	c.LL.AA(lines.Bold)
}

func (c *direct) width() int {
	if c.gaps {
		return len(directTitle) + 2
	}
	return len(directTitle)
}

// sourced is the scrolling component with selectable lines whose
// content comes from a source with a Liner implementation (see
// sourced.OnInit).
type sourced struct {
	lines.Component
	gaps bool
}

func (c *sourced) OnInit(e *lines.Env) {
	if c.gaps {
		frame(c, sourcedTitle, false)
		c.Src = &lines.ContentSource{Liner: newLiner()}
	} else {
		c.Src = &lines.ContentSource{Liner: newLinerTrimmed()}
	}
	c.Dim().SetWidth(c.width())
}

var sourcedTitle = []rune("sourced")

func (c *sourced) OnFocusLost(_ *lines.Env) {
	if c.gaps {
		c.Gaps(0).AA(lines.ZeroStyle)
		return
	}
	c.LL.AA(c.Globals().AA(lines.Default))
}

func (c *sourced) OnFocus(_ *lines.Env) {
	if c.gaps {
		c.Gaps(0).AA(lines.Bold)
		return
	}
	c.LL.AA(lines.Bold)
}

func (c *sourced) width() int {
	if c.gaps {
		return len(sourcedTitle) + 2
	}
	return len(sourcedTitle)
}

type liner struct {
	cc           []string
	notFocusable map[int]bool
}

func newLiner() *liner {
	lr := liner{notFocusable: map[int]bool{}}
	for i, l := range ll {
		lr.cc = append(lr.cc, lines.Filler+l+lines.Filler)
		if i%2 != 0 {
			continue
		}
		lr.notFocusable[i] = true
	}
	return &lr
}

// Print prints the line with given index idx to given line writer w and
// returns true if there are lines with a greater index left to write.
func (l *liner) Print(idx int, w *lines.EnvLineWriter) bool {
	if len(l.cc) <= idx || idx < 0 {
		return false
	}
	fmt.Fprintf(w, l.cc[idx])
	return idx+1 < len(l.cc)
}

// Len returns the total number of content lines a liner implementation
// can provide to its associated component.
func (sl liner) Len() int { return len(sl.cc) }

// IsFocusable returns true iff the line with given index idx is
// focusable.
func (l *liner) IsFocusable(idx int) bool {
	return !l.notFocusable[idx]
}

// Highlighted indicates if focusable lines are highlighted if focused.
// And in case they are highlighted if they should be trimmed
// highlighted.
func (l *liner) Highlighted() (bool, bool) { return true, false }

type linerTrimmed struct {
	*liner
}

func newLinerTrimmed() *linerTrimmed {
	return &linerTrimmed{liner: newLiner()}
}

func (l *linerTrimmed) Highlighted() (bool, bool) { return true, true }
