// Copyright (c) 2022 Stephan Lukits. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package lines

import (
	"strings"
)

type line struct {
	// dirty is set to true if the a line's content changes.
	dirty bool
	// stale if not zero contains the content at the last screen update.
	stale string
	// content is alines current content if stale is not zero the
	// content has not been written to the screen yet.
	content string

	// fmt can hold formattings for a line's content respectively
	// sub-strings of it.  It defaults to the component's formattings.
	fmt FmtMask

	// settings coming from parent component's lines factory.
	global *global

	sty Style

	ss lineStyles

	ff LineFlags
}

// Set updates the content of a line.
func (l *line) reset(sty Style, ff LineFlags) *line {
	l.sty = sty
	l.ss = nil
	l.ff = ff
	if !l.dirty {
		l.dirty = true
	}
	if l.stale == "" {
		l.stale = l.content
	}
	l.content = ""
	return l
}

func (l *line) addStyleRange(sr SR, rr ...SR) {
	if l.ss == nil {
		l.ss = lineStyles{}
	}
	l.ss[sr.Range] = sr.Style
	for _, r := range rr {
		l.ss[r.Range] = r.Style
	}
	if !l.dirty {
		l.dirty = true
	}
}

func (l *line) setFlags(ff LineFlags) {
	l.ff = ff
	if !l.dirty {
		l.dirty = true
	}
}

// replaceAt replaces current content with given content.  If cell is -1
// the whole content is replaced and given style is set for the whole
// line.  Is cell >= 0 the content is replaced starting at cell (padded
// with blanks if necessary) and given style is set for the range
// cell to len(content).  Is cell < -1 the call is ignored.
func (l *line) replaceAt(
	cell int, content string, s Style, ff LineFlags,
) {
	if cell < -1 {
		return
	}
	if !l.dirty {
		l.dirty = true
	}
	if l.stale == "" {
		l.stale = l.content
	}
	l.ff = ff
	if cell == -1 {
		l.sty = s
		l.content = content
		return
	}
	if l.ss == nil {
		l.ss = lineStyles{}
	}
	l.ss[Range{cell, cell + len(content)}] = s
	if len(l.content) < cell {
		l.content = l.content +
			strings.Repeat(" ", cell-len(l.content)) + content
		return
	}
	l.content = l.content[:cell] + content
}

type runeWriter interface {
	Display(x, y int, r rune, s Style)
}

func (l *line) sync(x, y, width int, rw runeWriter) {
	l.dirty = false
	if l.fmt&onetimeFilled > 0 {
		l.fill(x, y, width, rw)
		l.fmt &^= onetimeFilled
	}
	l.toScreen(x, y, width, rw)
	l.stale = ""
}

func (l *line) fill(x, y, width int, rw runeWriter) {
	if width < 0 || x < 0 {
		return
	}
	for i := x; i < x+width; i++ {
		rw.Display(i, y, ' ', l.sty)
	}
	l.stale = ""
}

func (l *line) SwitchHighlighted() {
	if l.ff&Highlighted == Highlighted {
		l.ff &^= Highlighted
		l.fmt |= onetimeFilled
	} else {
		l.ff |= Highlighted
	}
	l.dirty = true
}

func (l *line) IsHighlighted() bool {
	return l.ff&Highlighted == Highlighted
}

func (l *line) IsFocusable() bool {
	return l.ff&NotFocusable == 0
}

func (l *line) toScreen(x, y, width int, rw runeWriter) {

	forScr, ss := l.forScreen(width)
	if l.ff&Highlighted > 0 {
		l.toScreenHighlighted(forScr, x, y, width, rw)
		return
	}
	forScrLen := len(forScr)

	for i, r := range forScr {
		if i == width {
			break
		}
		rw.Display(x+i, y, r, ss.of(i, l.sty))
	}

	if forScrLen >= width {
		return
	}

	l.fill(x+forScrLen, y, width-forScrLen, rw)
}

func (l *line) toScreenHighlighted(
	s string, x, y, width int, rw runeWriter,
) {

	trimL, trimR := 0, 0
	for _, r := range s {
		if r != ' ' {
			break
		}
		trimL++
	}
	for i := len(s) - 1; i >= 0; i-- {
		if s[i] != ' ' {
			break
		}
		trimR++
	}

	if trimL == len(s) {
		l.fill(x, y, width, rw)
	}

	var rvr Style
	if l.sty.AA()&Reverse == 0 {
		rvr = l.sty.WithAdded(Reverse)
	} else {
		rvr = l.sty.WithRemoved(Reverse)
	}
	for i, r := range s {
		if i == width {
			break
		}
		switch {
		case i < trimL:
			rw.Display(x+i, y, ' ', l.sty)
		case i+trimR >= len(s):
			rw.Display(x+i, y, ' ', l.sty)
		default:
			rw.Display(x+i, y, r, rvr)
		}
	}
	if len(s) < width {
		l.fill(x+len(s), y, width-len(s), rw)
	}
}

// forScreen calculates from the line's content the string which should
// be written to the screen as well as the styles by expanding leading
// tabs and line filler.
func (l *line) forScreen(width int) (string, lineStyles) {
	if len(l.content) == 0 {
		return "", l.ss
	}

	content := l.content
	ss := l.ss.copy()

	tt := 0 // leading tab expansion
	for _, r := range content {
		if r != '\t' {
			break
		}
		tt++
	}
	for i := 0; i < tt; i++ {
		ss.expand(i*l.global.tabWidth, l.global.tabWidth-1)
		content = content[:i*l.global.tabWidth] +
			strings.Repeat(" ", l.global.tabWidth) +
			content[i*l.global.tabWidth+1:]
	}

	ff := strings.Split(content, LineFiller)
	if len(ff) == 1 {
		return content, ss
	}

	blank := width - (len(content) - (len(ff) - 1)) // ignore filler
	if blank <= len(ff)-1 {                         // overflow width
		if string(content[0]) == LineFiller { // remove leading filler
			ss.contract(0, 1)
			return strings.TrimSpace(strings.Join(ff, " ")), ss
		}
		// don't trim we might have expanded tabs
		return strings.Join(ff, " "), ss
	}

	ffPoints := []int{len(ff[0])}
	for i, s := range ff[1 : len(ff)-1] {
		ffPoints = append(ffPoints, ffPoints[i]+len(s)+1)
	}

	dist := blank / (len(ff) - 1)
	if (len(ff)-1)*dist == blank {
		for i, p := range ffPoints {
			expDist := dist - 1
			if expDist == 0 { // replace filler with single space
				expDist = 1
			}
			ss.expand(p+i*expDist, expDist)
		}
		return strings.Join(ff, strings.Repeat(" ", dist)), ss
	}

	b, fill := strings.Builder{}, strings.Repeat(" ", dist)
	rest := blank % (len(ff) - 1)
	plusOne := 0
	for i, s := range ff {
		b.WriteString(s + fill)
		if rest > 0 {
			ss.expand(ffPoints[i]+(dist)*plusOne, dist)
			b.WriteRune(' ')
			rest--
			plusOne++
		} else {
			if i >= len(ffPoints) {
				continue
			}
			expDist := dist - 1
			if expDist == 0 { // replace filler with single space
				expDist = 1
			}
			ss.expand(
				ffPoints[i]+(dist)*plusOne+expDist*(i-(plusOne-1)), expDist)
		}
	}

	return b.String(), ss
}

type LineFlags uint64

const (
	NotFocusable LineFlags = 1 << iota
	Highlighted
)

// Range is a two component array of which the first component should be
// smaller than the second, i.e. r.Start() <= r.End() if r is a
// Range-instance.
type Range [2]int

// Start index of a [lines.TestLine] style range.  Not the start index
// is inclusive.
func (r Range) Start() int { return r[0] }

// End index of a [lines.TestLine] style range.  Note the end index is
// exclusive.
func (r Range) End() int { return r[1] }

// SetStart sets given range's start index.
func (r *Range) SetStart(s int) *Range {
	(*r)[0] = s
	return r
}

// copy returns a copy of given range.
func (r Range) copy() Range {
	return Range{r[0], r[1]}
}

// IncrementStart increments a ranges start value by one.
func (r *Range) IncrementStart() { (*r)[0]++ }

// IncrementEnd increments a ranges end value by one.
func (r *Range) IncrementEnd() { (*r)[1]++ }

// SetEnd sets given range's end index.
func (r *Range) SetEnd(e int) *Range {
	(*r)[1] = e
	return r
}

// Shift increases start and and index by given s.
func (r Range) Shift(s int) Range {
	return Range{r[0] + s, r[1] + s}
}

// Contains returns true if given i is in the style range r
// [r.Start,r.End[.
func (r Range) Contains(i int) bool {
	return r.Start() <= i && i < r.End()
}

// SR represents a ranged style which may be set for a line see
// [Env.AddStyleRange].
type SR struct {
	Range
	Style
}

type lineStyles map[Range]Style

func (s lineStyles) copy() lineStyles {
	if s == nil {
		return nil
	}
	cp := lineStyles{}
	for r, s := range s {
		cp[r.copy()] = s
	}
	return cp
}

func (s lineStyles) expand(at, by int) {
	if s == nil {
		return
	}
	update := map[Range]Range{}
	for r := range s {
		if r.End() <= at {
			continue
		}
		switch {
		case r.Contains(at):
			update[r] = Range{r.Start(), r.End() + by}
		default:
			update[r] = r.Shift(by)
		}
	}
	for k, u := range update {
		if u.Start() < u.End() {
			s[u] = s[k]
		}
		delete(s, k)
	}
}

func (s lineStyles) contract(at, by int) {
	s.expand(at, -by)
}

// of returns the style for given line-cell.  Note style ranges are
// defined relative to a lines origin; i.e. the first cell's style is
//
//	s.of(0, dflt)
func (s lineStyles) of(cell int, dflt Style) Style {

	if s == nil {
		return dflt
	}
	for r := range s {
		if !r.Contains(cell) {
			continue
		}
		return s[r]
	}
	return dflt
}
