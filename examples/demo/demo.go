// Copyright (c) 2022 Stephan Lukits. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

/*
Package demo provides helper for the example packages in "examples".
*/
package demo

import (
	"fmt"

	"github.com/slukits/lines"
	"github.com/slukits/lines/internal/lyt"
)

// Demo implements features all demonstrations, i.e. menu, context,
// tool-tip and stacked, have in common.
type Demo struct {
	Titled
	dg                          dimGapper
	Next                        lines.Componenter
	DefaultWidth, DefaultHeight int
}

type dimGapper interface {
	Gaps(int) *lines.GapsWriter
	GapsLen() (_, _, _, _ int)
	Dim() *lyt.Dim
}

// Init sets up the embedding component's title and its default size.
func (d *Demo) InitDemo(dg dimGapper, e *lines.Env, title []rune) {
	if d.DefaultWidth == 0 {
		d.DefaultWidth = 25
	}
	if d.DefaultHeight == 0 {
		d.DefaultHeight = 8
	}
	d.dg = dg
	d.Title = title
	d.Default(dg, e)
	dg.Dim().SetWidth(d.width()).SetHeight(d.height())
}

func (d *Demo) WriteTip(s string) {
	rr := []rune(s)
	fmt.Fprint(d.dg.Gaps(1).TopLeft.FG(0xEEEEEE).BG(0x666666),
		string(rr[0]))
	fmt.Fprint(d.dg.Gaps(1).TopRight.FG(0xEEEEEE).BG(0x666666), "")
	fmt.Fprint(d.dg.Gaps(1).Top.FG(0xEEEEEE).BG(0x666666),
		string(rr[1:]))
}

type heighter interface {
	Height() int
}

func (d *Demo) height() int {
	if h, ok := d.dg.(heighter); ok {
		return h.Height()
	}
	return 8
}

type widther interface {
	Width() int
}

func (d *Demo) width() int {
	if w, ok := d.dg.(widther); ok {
		return w.Width()
	}
	return 25
}

// OnFocus switches to double-framing.
func (d *Demo) OnFocus(e *lines.Env) { d.Focused(d.dg, e) }

// OnFocusLost reverts to single-framing.
func (d *Demo) OnFocusLost(e *lines.Env) {
	d.Default(d.dg, e)
	if lt, _, _, _ := d.dg.GapsLen(); lt > 1 {
		fmt.Fprint(d.dg.Gaps(1).TopLeft, "")
		fmt.Fprint(d.dg.Gaps(1).TopRight, "")
		fmt.Fprint(d.dg.Gaps(1).Top, "")
	}
}

func (d *Demo) OnKey(e *lines.Env, k lines.Key, mm lines.ModifierMask) {
	if k != lines.TAB {
		return
	}
	e.Lines.Focus(d.Next)

	// otherwise each tab bubbles ub to the App-instance i.e. the
	// menu-demo is always focused.
	e.StopBubbling()
}
