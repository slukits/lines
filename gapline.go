// Copyright (c) 2022 Stephan Lukits. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package lines

import (
	"sort"
	"strings"
)

type gapLine struct {
	rr      []rune
	ss      styleRanges
	filling bool
	fillAt  []int
}

func (l *gapLine) setDefaultStyle(s Style) {
	l.ss = newStyleRanges(s)
}

func (l *gapLine) ensureStyleRanges() styleRanges {
	if l.ss == nil {
		l.ss = newStyleRanges(DefaultStyle)
	}
	return l.ss
}

func (l *gapLine) withAA(aa StyleAttributeMask) {
	l.ensureStyleRanges().withAA(aa)
}

func (l *gapLine) withFG(c Color) {
	l.ensureStyleRanges().withFG(c)
}

func (l *gapLine) withBG(c Color) {
	l.ensureStyleRanges().withBG(c)
}

func (l *gapLine) set(s string) {
	l.rr = []rune(s)
}

func (l *gapLine) setStyled(s string, sty Style) {
	l.set(s)
	l.ss = styleRanges{Range{0, len(l.rr)}: sty}
}

func (l *gapLine) setAt(at int, rr []rune) {
	l.padTo(at)
	l.rr = append(l.rr[:at], rr...)
}

func (l *gapLine) setStyledAt(at int, rr []rune, sty Style) {
	l.setAt(at, rr)
	if l.ss == nil {
		l.ss = styleRanges{}
	}
	l.ss[Range{at, at + len(rr)}] = sty
}

func (l *gapLine) setFilling(at int, r rune) {
	l.padTo(at)
	l.rr = append(l.rr[:at], r)
	l.fillAt = append(l.fillAt, at)
}

func (l *gapLine) setAtFilling(at int, r rune) {
	l.padTo(at)
	l.rr = append(l.rr[:at], r)
	l.fillAt = append(l.fillAt, at)
}

func (l *gapLine) setStyledAtFilling(at int, r rune, sty Style) {
	l.setAtFilling(at, r)
	if l.ss == nil {
		l.ss = styleRanges{}
	}
	l.ss[Range{at, at + 1}] = sty
}

func (l *gapLine) display(width int, dflt Style) ([]rune, styleRanges) {
	// TODO: adapt style range
	ss := l.ss.copyWithDefault(dflt)
	if len(l.rr) == 0 {
		return []rune(strings.Repeat(" ", width)), ss
	}
	if len(l.fillAt) > 0 {
		return l.displayFilledAt(width, ss)
	}
	if l.filling {
		return l.displayFilling(width, ss)
	}
	if len(l.rr) < width {
		return append(
			l.rr, []rune(strings.Repeat(" ", width-len(l.rr)))...), ss
	}
	return l.rr, l.ss
}

func (l *gapLine) displayFilling(
	width int, ss styleRanges,
) ([]rune, styleRanges) {

	c := width / len(l.rr)
	if c == 0 {
		if len(l.rr) < width {
			rr := append([]rune{}, l.rr...)
			return append(rr, l.rr[:width-len(rr)]...), ss
		}
		return l.rr, ss
	}
	rr := []rune{}
	for i := 0; i < c; i++ {
		rr = append(rr, l.rr...)
	}
	if len(rr) < width {
		return append(rr, l.rr[:width-len(rr)]...), ss
	}
	return rr, ss
}

func (l *gapLine) displayFilledAt(
	width int, ss styleRanges,
) ([]rune, styleRanges) {
	if len(l.rr) >= width {
		return l.rr, ss
	}
	f := (width - (len(l.rr) - len(l.fillAt))) / len(l.fillAt)
	mf := (width - (len(l.rr) - len(l.fillAt))) % len(l.fillAt)
	ff := map[int][]rune{}
	for i := 0; i < len(l.fillAt)-mf; i++ {
		ff[l.fillAt[i]] = []rune(
			strings.Repeat(string(l.rr[l.fillAt[i]]), f))
	}
	for i := len(l.fillAt) - mf; i < len(l.fillAt); i++ {
		ff[l.fillAt[i]] = []rune(
			strings.Repeat(string(l.rr[l.fillAt[i]]), f+1))
	}
	rr, last := []rune{}, 0
	for _, at := range l.fillAt {
		rr = append(rr, l.rr[last:at]...)
		rr = append(rr, ff[at]...)
		last = at + 1
	}
	return rr, l.adjustAtStyles(ff, ss)
}

func (l *gapLine) adjustAtStyles(
	filler map[int][]rune, ss styleRanges,
) styleRanges {
	ff := make([]int, 0, len(filler))
	for k := range filler {
		ff = append(ff, k)
	}
	sort.Ints(ff)
	for _, f := range ff {
		ss.expand(f, len(filler[f])-1)
	}
	return ss
}

func (l *gapLine) padTo(at int) {
	if len(l.rr) >= at {
		return
	}
	l.rr = append(l.rr, []rune(strings.Repeat(" ", at-len(l.rr)))...)
}
