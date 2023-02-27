// Copyright (c) 2022 Stephan Lukits. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package demo

import (
	"fmt"

	"github.com/slukits/lines"
)

type Gapper interface{ Gaps(int) *lines.GapsWriter }

type Titled struct {
	Title []rune
	Gapper
}

func (tt *Titled) Single(e *lines.Env) *Titled {
	lines.Print(tt.Gaps(0).Vertical.At(0).Filling(), '│')
	lines.Print(tt.Gaps(0).Horizontal.At(0).Filling(), '─')
	fmt.Fprintf(tt.Gaps(0).Corners, "┌┐┘└")
	lines.Print(tt.Gaps(0).Top.At(1), tt.Title)
	lines.Print(tt.Gaps(0).Top.At(1+len(tt.Title)).Filling(), '─')
	return tt
}

func (tt *Titled) Double(e *lines.Env) *Titled {
	lines.Print(tt.Gaps(0).Vertical.At(0).Filling(), '║')
	lines.Print(tt.Gaps(0).Horizontal.At(0).Filling(), '═')
	fmt.Fprintf(tt.Gaps(0).Corners, "╔╗╝╚")
	lines.Print(tt.Gaps(0).Top.At(1).AA(lines.Bold), tt.Title)
	lines.Print(tt.Gaps(0).Top.At(1+len(tt.Title)).Filling(), '═')
	return tt
}

func (tt *Titled) Styled(e *lines.Env, sty lines.Style) *Titled {
	lines.Print(tt.Gaps(0).Sty(sty).Vertical.At(0).Filling(), '│')
	lines.Print(tt.Gaps(0).Horizontal.At(0).Filling(), '─')
	fmt.Fprintf(tt.Gaps(0).Corners, "┌┐┘└")
	lines.Print(tt.Gaps(0).Top.At(1), tt.Title)
	lines.Print(tt.Gaps(0).Top.At(1+len(tt.Title)).Filling(), '─')
	return tt
}

func (tt *Titled) Default(e *lines.Env) {
	lines.Print(tt.Gaps(0).Vertical.At(0).Filling(), '│')
	lines.Print(tt.Gaps(0).Horizontal.At(0).Filling(), '─')
	fmt.Fprintf(tt.Gaps(0).Corners, "┌┐┘└")
	lines.Print(tt.Gaps(0).Top.At(1), tt.Title)
	lines.Print(tt.Gaps(0).Top.At(1+len(tt.Title)).Filling(), '─')
}

func (f *Titled) Focused(g Gapper, e *lines.Env) {
	lines.Print(g.Gaps(0).Vertical.At(0).Filling(), '║')
	lines.Print(g.Gaps(0).Horizontal.At(0).Filling(), '═')
	fmt.Fprintf(g.Gaps(0).Corners, "╔╗╝╚")
	lines.Print(g.Gaps(0).Top.At(1).AA(lines.Bold), f.Title)
	lines.Print(g.Gaps(0).Top.At(1+len(f.Title)).Filling(), '═')
}
