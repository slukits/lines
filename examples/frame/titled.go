// Copyright (c) 2022 Stephan Lukits. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package frame

import (
	"fmt"

	"github.com/slukits/lines"
)

type Titled struct{ Title []rune }

type gapper interface{ Gaps(int) *lines.GapsWriter }

func (f *Titled) Default(g gapper, e *lines.Env) {
	lines.Print(g.Gaps(0).Vertical.At(0).Filling(), '│')
	lines.Print(g.Gaps(0).Horizontal.At(0).Filling(), '─')
	fmt.Fprintf(g.Gaps(0).Corners, "┌┐┘└")
	lines.Print(g.Gaps(0).Top.At(1), f.Title)
	lines.Print(g.Gaps(0).Top.At(1+len(f.Title)).Filling(), '─')
}

func (f *Titled) Focused(g gapper, e *lines.Env) {
	lines.Print(g.Gaps(0).Vertical.At(0).Filling(), '║')
	lines.Print(g.Gaps(0).Horizontal.At(0).Filling(), '═')
	fmt.Fprintf(g.Gaps(0).Corners, "╔╗╝╚")
	lines.Print(g.Gaps(0).Top.At(1).AA(lines.Bold), f.Title)
	lines.Print(g.Gaps(0).Top.At(1+len(f.Title)).Filling(), '═')
}
