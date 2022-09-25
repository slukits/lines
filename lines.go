// Copyright (c) 2022 Stephan Lukits. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// NOTE this is not the entry point to this package!  The central types
// are found in events.go, screen.go, component.go and env.go.  Also
// testing.go is a good place to start.

package lines

import (
	"strings"

	"github.com/gdamore/tcell/v2"
)

// LineFiller can be used in component content-lines indicating that a
// line l should fill up its whole width whereas its remaining empty
// space is spread equally over filler found in l.
const LineFiller = string(rune(29))

type lines []*line

// append given content lines to current content
func (ll *lines) append(ff LineFlags, sty tcell.Style, cc ...[]byte) {
	for _, c := range cc {
		l := &line{
			content: string(c),
			dirty:   true,
			ff:      ff,
			sty:     sty,
		}
		*ll = append(*ll, l)
	}
}

// replaceAt replaces starting at given index the following lines with
// given content lines.  replaceAt is a no-op if idx < 0 or len(cc) == 0
func (ll *lines) replaceAt(
	idx, cell int, ff LineFlags, sty tcell.Style, cc ...[]byte,
) {
	if idx < 0 || len(cc) == 0 {
		return
	}
	for idx+len(cc) > len(*ll) {
		*ll = append(*ll, &line{dirty: true})
	}
	for i, j := idx, 0; i < idx+len(cc); i++ {
		(*ll)[i].replaceAt(cell, string(cc[j]), sty, ff)
		j++
	}
}

// IsDirty returns true if on of the lines is dirty.
func (ll lines) IsDirty() bool {
	for _, l := range ll {
		if !l.dirty {
			continue
		}
		return true
	}
	return false
}

// ForDirty calls back for every dirty line.
func (ll lines) ForDirty(cb func(int, *line)) {
	for i, l := range ll {
		if !l.dirty {
			return
		}
		cb(i, l)
	}
}

// For calls back for every line of given lines ll starting at given
// offset.
func (ll lines) For(offset int, cb func(int, *line) (stop bool)) {
	for i, l := range ll[offset:] {
		if cb(i, l) {
			return
		}
	}
}

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

// Contains returns true if given i is in the style range r
// [r.Start,r.End[.
func (r Range) Contains(i int) bool {
	return r.Start() <= i && i < r.End()
}

type lineStyles map[Range]tcell.Style

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

	sty tcell.Style

	ss lineStyles

	ff LineFlags
}

// Set updates the content of a line.
func (l *line) set(content string) *line {
	if content == l.content {
		return l
	}
	if !l.dirty {
		l.dirty = true
	}
	if l.stale == "" {
		l.stale = l.content
	}
	l.content = content
	return l
}

// replaceAt replaces current content with given content.  If cell is -1
// the whole content is replaced and given style is set for the whole
// line.  Is cell >= 0 the content is replaced starting at cell (padded
// with blanks if necessary) and given style is set for the range
// cell to len(content).  Is cell < -1 the call is ignored.
func (l *line) replaceAt(
	cell int, content string, s tcell.Style, ff ...LineFlags,
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
	_ff := LineFlags(0)
	for _, f := range ff {
		_ff |= f
	}
	l.ff = _ff
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
	SetContent(x, y int, r rune, combo []rune, s tcell.Style)
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
		rw.SetContent(i, y, ' ', nil, l.sty)
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

	forScr := l.forScreen(width)
	if l.ff&Highlighted > 0 {
		l.toScreenHighlighted(forScr, x, y, width, rw)
		return
	}
	forScrLen := len(forScr)

	for i, r := range forScr {
		if i == width {
			break
		}
		rw.SetContent(x+i, y, r, nil, l.ss.of(i, l.sty))
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

	_, _, aa := l.sty.Decompose()
	rvr := l.sty.Reverse(true)
	if aa&tcell.AttrReverse > 0 {
		rvr = l.sty.Reverse(false)
	}
	for i, r := range s {
		if i == width {
			break
		}
		switch {
		case i < trimL:
			rw.SetContent(x+i, y, ' ', nil, l.sty)
		case i+trimR >= len(s):
			rw.SetContent(x+i, y, ' ', nil, l.sty)
		default:
			rw.SetContent(x+i, y, r, nil, rvr)
		}
	}
	if len(s) < width {
		l.fill(x+len(s), y, width-len(s), rw)
	}
}

func (l *line) forScreen(width int) string {

	ff := strings.Split(l.content, LineFiller)
	if len(ff) == 1 {
		return l.content
	}

	blank := width - (len(l.content) - (len(ff) - 1)) // ignore filler
	if blank <= len(ff)-1 {
		return strings.Join(ff, " ")
	}

	dist := blank / (len(ff) - 1)
	if (len(ff)-1)*dist == blank {
		return strings.Join(ff, strings.Repeat(" ", dist))
	}

	b, fill := strings.Builder{}, strings.Repeat(" ", dist)
	rest := blank % (len(ff) - 1)
	for _, s := range ff {
		b.WriteString(s + fill)
		if rest > 0 {
			b.WriteRune(' ')
			rest--
		}
	}

	return strings.TrimSpace(b.String())
}

// of returns the style for given line-cell.  Note style ranges are
// defined relative to a lines origin; i.e. the first cell's style is
//
//	s.of(0, dflt)
func (s lineStyles) of(cell int, dflt tcell.Style) tcell.Style {

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

type LineFlags uint64

const (
	NotFocusable LineFlags = 1 << iota
	Highlighted
)
