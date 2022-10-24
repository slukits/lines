// Copyright (c) 2022 Stephan Lukits. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package lines

import (
	"strings"
)

type cmpLine struct {
	// dirty is set to true if the a line's content changes.
	dirty bool
	// stale if not zero contains the content at the last screen update.
	stale string
	// content is a line's current content if stale is not zero the
	// content has not been written to the screen yet.
	content string

	// fmt can hold formattings for a line's content respectively
	// sub-strings of it.  It defaults to the component's formattings.
	fmt FmtMask

	// settings coming from parent component's lines factory.
	global *globals

	sty Style

	ss styleRanges

	ff LineFlagsZZZ
}

// Set updates the content of a line.
func (l *cmpLine) reset(sty Style, ff LineFlagsZZZ) *cmpLine {
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

func (l *cmpLine) addStyleRange(sr SR, rr ...SR) {
	if l.ss == nil {
		l.ss = styleRanges{}
	}
	l.ss[sr.Range] = sr.Style
	for _, r := range rr {
		l.ss[r.Range] = r.Style
	}
	if !l.dirty {
		l.dirty = true
	}
}

func (l *cmpLine) setFlags(ff LineFlagsZZZ) {
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
func (l *cmpLine) replaceAt(
	cell int, content string, s Style, ff LineFlagsZZZ,
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
		l.ss = styleRanges{}
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

func (l *cmpLine) sync(x, y, width int, rw runeWriter) {
	l.dirty = false
	if l.fmt&onetimeFilled > 0 {
		l.fill(x, y, width, rw)
		l.fmt &^= onetimeFilled
	}
	l.toScreen(x, y, width, rw)
	l.stale = ""
}

func (l *cmpLine) fill(x, y, width int, rw runeWriter) {
	if width < 0 || x < 0 {
		return
	}
	for i := x; i < x+width; i++ {
		rw.Display(i, y, ' ', l.sty)
	}
	l.stale = ""
}

func (l *cmpLine) SwitchHighlighted() {
	if l.ff&HighlightedZZZ == HighlightedZZZ {
		l.ff &^= HighlightedZZZ
		l.fmt |= onetimeFilled
	} else {
		l.ff |= HighlightedZZZ
	}
	l.dirty = true
}

func (l *cmpLine) IsHighlighted() bool {
	return l.ff&HighlightedZZZ == HighlightedZZZ
}

func (l *cmpLine) IsFocusable() bool {
	return l.ff&NotFocusableZZZ == 0
}

func (l *cmpLine) toScreen(x, y, width int, rw runeWriter) {

	forScr, ss := l.forScreen(width, l.sty)
	if l.ff&HighlightedZZZ > 0 {
		l.toScreenHighlighted(forScr, x, y, width, rw)
		return
	}
	forScrLen := len(forScr)

	for i, r := range forScr {
		if i == width {
			break
		}
		rw.Display(x+i, y, r, ss.of(i))
	}

	if forScrLen >= width {
		return
	}

	l.fill(x+forScrLen, y, width-forScrLen, rw)
}

func (l *cmpLine) toScreenHighlighted(
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
func (l *cmpLine) forScreen(width int, dflt Style) (string, styleRanges) {
	if len(l.content) == 0 {
		return "", l.ss.copyWithDefault(dflt)
	}

	content := l.content
	ss := l.ss.copyWithDefault(dflt)

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

type LineFlagsZZZ uint64

const (
	NotFocusableZZZ LineFlagsZZZ = 1 << iota
	HighlightedZZZ
)
