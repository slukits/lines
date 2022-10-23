// Copyright (c) 2022 Stephan Lukits. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package lines

import (
	"sort"
	"strings"
)

// A line structure stores the content, style-information and line
// filler of a line and provides operations to calculate certain display
// representations of a line's content.
type line struct {
	rr     []rune
	ss     styleRanges
	fillAt []int
}

// pad returns a rune slice whose len is given width by padding given
// rune slice rr with spaces accordingly.  pad will panic if len(rr) >
// width.
func (l *line) pad(rr []rune, width int) []rune {
	c := width - len(rr)
	return append(rr, []rune(strings.Repeat(" ", c))...)
}

// setDefaultStyle sets given line l's zeroRange style of its style
// ranges.
func (l *line) setDefaultStyle(s Style) {
	l.ss = newStyleRanges(s)
}

// ensureStyleRanges returns given line l's style ranges and initializes
// them first if necessary.
func (l *line) ensureStyleRanges() styleRanges {
	if l.ss == nil {
		l.ss = newStyleRanges(DefaultStyle)
	}
	return l.ss
}

// withAA sets of given line l its default style's style attributes.
// NOTE if no default style has been set (see line.setDefaultStyle) the
// DefaultStyle is used and modified with given attributes.
func (l *line) withAA(aa StyleAttributeMask) {
	l.ensureStyleRanges().withAA(aa)
}

// withFG sets of given line l its default style's foreground color.
// NOTE if no default style has been set (see line.setDefaultStyle) the
// DefaultStyle is used and modified with given foreground color.
func (l *line) withFG(c Color) {
	l.ensureStyleRanges().withFG(c)
}

// withBG sets of given line l its default style's background color.
// NOTE if no default style has been set (see line.setDefaultStyle) the
// DefaultStyle is used and modified with given background color.
func (l *line) withBG(c Color) {
	l.ensureStyleRanges().withBG(c)
}

// set sets given line l's content converting given string s to a rune
// slice.
func (l *line) set(s string) {
	l.rr = []rune(s)
	if len(l.fillAt) > 0 {
		l.fillAt = []int{}
	}
}

// setStyled sets given line l's content to the conversion of given
// string s to a rune slice and adds a style range for this slice with
// given style.  Note the main difference to set is that if the line has
// a default style d which differs from provided style and the display
// is wider than the l's content then the remaining blanks have the
// style d.
func (l *line) setStyled(s string, sty Style) {
	l.set(s)
	l.ss = styleRanges{Range{0, len(l.rr)}: sty}
}

// setAt sets given rune slice rr at given position p in given line l's
// content overwriting what has been at and after this position.  If
// needed rr is padded with spaces to p.
func (l *line) setAt(p int, rr []rune) {
	l.padTo(p)
	l.rr = append(l.rr[:p], rr...)
	l.truncateAt(p)
}

// truncateAt truncates fillers and styles at and after given position
// p.
func (l *line) truncateAt(p int) {
	truncate := -1
	for i, f := range l.fillAt {
		if f < p {
			continue
		}
		truncate = i
		break
	}
	if truncate >= 0 {
		l.fillAt = l.fillAt[:truncate]
	}
	for r := range l.ss {
		if r.Start() < p || r == zeroRange {
			continue
		}
		delete(l.ss, r)
	}
}

// padTo pads a line's l content rr to given position p with spaces.
func (l *line) padTo(p int) {
	if len(l.rr) >= p {
		return
	}
	l.rr = append(l.rr, []rune(strings.Repeat(" ", p-len(l.rr)))...)
}

// setStyledAt sets given rune slice rr at given position p in given
// line l's content overwriting what has been at and after this position
// and adds a corresponding style range.  If needed l's content is
// padded with spaces until p.
func (l *line) setStyledAt(at int, rr []rune, sty Style) {
	l.setAt(at, rr)
	if l.ss == nil {
		l.ss = styleRanges{}
	}
	l.ss[Range{at, at + len(rr)}] = sty
}

// setAtFilling sets given rune r at given position p of given line l's
// content truncating all possibly following content.
func (l *line) setAtFilling(p int, r rune) {
	l.padTo(p)
	l.truncateAt(p)
	l.rr = append(l.rr[:p], r)
	l.fillAt = append(l.fillAt, p)
}

// setStyledAtFilling sets given rune r at given position p of given
// line l's content truncating all possibly following content and adds
// a corresponding style range.
func (l *line) setStyledAtFilling(at int, r rune, sty Style) {
	l.setAtFilling(at, r)
	if l.ss == nil {
		l.ss = styleRanges{}
	}
	l.ss[Range{at, at + 1}] = sty
}

// display returns a line's calculated content depending on given width
// and set filler as well as corresponding style ranges ready to print
// to the screen.
func (l *line) display(width int, dflt Style) ([]rune, styleRanges) {
	ss := l.ss.copyWithDefault(dflt)
	if len(l.rr) == 0 {
		return []rune(strings.Repeat(" ", width)), ss
	}
	rr := append([]rune{}, l.rr...)
	if len(rr) >= width {
		return rr[:width], ss
	}
	if len(l.fillAt) > 0 {
		return l.expandFillerAt(rr, width, ss)
	}
	if len(rr) < width {
		rr = l.pad(rr, width)
	}
	return rr, ss
}

// expandFillerAt expands runes marked as fillers by an equal amount in
// given runes rr to fit given width exactly adjusting given style
// ranges accordingly.
func (l *line) expandFillerAt(
	rr []rune, width int, ss styleRanges,
) ([]rune, styleRanges) {
	f := (width - (len(rr) - len(l.fillAt))) / len(l.fillAt)
	mf := (width - (len(rr) - len(l.fillAt))) % len(l.fillAt)
	ff := map[int][]rune{}
	for i := 0; i < len(l.fillAt)-mf; i++ {
		ff[l.fillAt[i]] = []rune(
			strings.Repeat(string(rr[l.fillAt[i]]), f))
	}
	for i := len(l.fillAt) - mf; i < len(l.fillAt); i++ {
		ff[l.fillAt[i]] = []rune(
			strings.Repeat(string(rr[l.fillAt[i]]), f+1))
	}
	_rr, last := []rune{}, 0
	for _, at := range l.fillAt {
		_rr = append(_rr, rr[last:at]...)
		_rr = append(_rr, ff[at]...)
		last = at + 1
	}
	if last < len(rr) {
		_rr = append(_rr, rr[last:]...)
	}
	return _rr, l.adjustFillerExpansionStyles(ff, ss)
}

func (l *line) adjustFillerExpansionStyles(
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
