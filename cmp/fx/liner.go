// Copyright (c) 2022 Stephan Lukits. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package fx

import (
	"fmt"

	"github.com/slukits/lines"
)

type Liner struct {
	// II is the provided content of this type's instance which defaults
	// to 8 lines: 1st, 2nd, 3rd, 4th, ..., 8th
	II  []string
	sty lines.Style

	init bool
}

func (l *Liner) ensureInit() {
	if l.init {
		return
	}
	l.init = true
	l.sty = lines.DefaultStyle
}

func (l *Liner) InitLines(n int) *Liner {
	l.ensureInit()
	l.II = NStrings(n)
	return l
}

func (l *Liner) SetII(ii []string) *Liner {
	l.ensureInit()
	l.II = ii
	return l
}

func (l *Liner) SetSty(sty lines.Style) *Liner {
	l.ensureInit()
	l.sty = sty
	return l
}

func (l *Liner) Print(idx int, w *lines.EnvLineWriter) bool {
	l.ensureInit()
	if l.II == nil {
		l.InitLines(8)
	}
	if len(l.II) <= idx || idx < 0 {
		return false
	}
	if l.sty.IsDefault() {
		fmt.Fprintf(w, l.II[idx])
	} else {
		fmt.Fprintf(w.Sty(l.sty), l.II[idx])
	}
	return idx+1 < len(l.II)
}

type ScrollableLiner struct {
	Liner
}

func (l *ScrollableLiner) Len() int {
	l.ensureInit()
	if l.II == nil {
		l.InitLines(8)
	}
	return len(l.II)
}

type SelectableLiner struct {
	ScrollableLiner
}

func (l *SelectableLiner) MaxWidth() int {
	l.ensureInit()
	maxWdth := 0
	for _, i := range l.II {
		if maxWdth >= len(i) {
			continue
		}
		maxWdth = len(i)
	}
	return maxWdth
}

func (l *SelectableLiner) SetII(ii []string) *SelectableLiner {
	l.ScrollableLiner.SetII(ii)
	return l
}

func (l *SelectableLiner) SetSty(sty lines.Style) *SelectableLiner {
	l.ScrollableLiner.SetSty(sty)
	return l
}

func (l *SelectableLiner) IsFocusable(_ int) bool { return true }

type HighlightingLiner struct {
	SelectableLiner
	hi lines.Style
}

func (l *HighlightingLiner) SetII(ii []string) *HighlightingLiner {
	l.SelectableLiner.SetII(ii)
	return l
}

func (l *HighlightingLiner) SetSty(sty lines.Style) *HighlightingLiner {
	l.SelectableLiner.SetSty(sty)
	return l
}

func (l *HighlightingLiner) SetHi(sty lines.Style) *HighlightingLiner {
	l.hi = sty
	return l
}

func (l *HighlightingLiner) Highlight(s lines.Style) lines.Style {
	return l.hi
}
