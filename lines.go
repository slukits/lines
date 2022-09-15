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
func (ll *lines) append(ff LineFlags, fmt *llFmt, cc ...[]byte) {
	for _, c := range cc {
		l := &line{
			content: string(c),
			dirty:   true,
			ff:      ff,
		}
		if fmt != nil {
			l.ss = lineStyles{Range{0, len(c)}: fmt.sty}
			l.fmt = fmt.mask
		}
		*ll = append(*ll, l)
	}
}

// replaceAt replaces starting at given index the following lines with
// given content lines.  replaceAt is a no-op if idx < 0 or len(cc) == 0
func (ll *lines) replaceAt(
	idx int, ff LineFlags, fmt *llFmt, cc ...[]byte,
) {
	if idx < 0 || len(cc) == 0 {
		return
	}
	for idx > len(*ll) {
		*ll = append(*ll, &line{dirty: true})
	}
	max, j := idx+len(cc), 0
	if max > len(*ll) {
		max = len(*ll)
	}
	for i := idx; i < max; i++ {
		(*ll)[i].set(string(cc[j]))
		(*ll)[i].ff = ff
		(*ll)[i].ss = nil
		(*ll)[i].fmt = 0
		if fmt != nil {
			(*ll)[i].ss = lineStyles{Range{0, len(cc[j])}: fmt.sty}
			(*ll)[i].fmt = fmt.mask
		}
		j++
	}
	ll.append(ff, fmt, cc[j:]...)
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

type runeWriter interface {
	SetContent(x, y int, r rune, combo []rune, s tcell.Style)
}

func (l *line) sync(x, y, width int, rw runeWriter, fmt llFmt) {
	l.dirty = false
	if l.fmt&filled == filled {
		sty := l.ss.of(0, fmt.sty, l.ff)
		for i := x; i < x+width; i++ {
			rw.SetContent(i, y, ' ', nil, sty)
		}
		l.stale = ""
	}
	l.toScreen(x, y, width, rw, fmt.sty)
	l.stale = ""
}

func (l *line) SwitchHighlighted() {
	if l.ff&Highlighted == Highlighted {
		l.ff &^= Highlighted
	} else {
		l.ff |= Highlighted
	}
	l.dirty = true
	l.fmt |= filled
}

func (l *line) toScreen(
	x, y, width int, rw runeWriter, sty tcell.Style,
) {

	forScr := l.forScreen(width)
	forScrLen := len(forScr)
	diffStale := len(l.stale) - forScrLen

	for i, r := range forScr {
		if i == width {
			break
		}
		rw.SetContent(x+i, y, r, nil, l.ss.of(i, sty, l.ff))
	}

	if forScrLen > width {
		return
	}

	for i := 0; i < diffStale; i++ {
		if forScrLen+i == width {
			break
		}
		rw.SetContent(
			x+forScrLen+i, y, ' ', nil, l.ss.of(forScrLen+i, sty, l.ff))
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

type lineStyles map[Range]tcell.Style

// of returns the style for given line-cell.  Note style ranges are
// defined relative to a lines origin; i.e. the first cell's style is
//
//	s.of(0, dflt)
func (s lineStyles) of(
	cell int, dflt tcell.Style, ff LineFlags,
) tcell.Style {

	if ff&Highlighted == Highlighted {
		dflt = dflt.Reverse(true)
	}
	if s == nil {
		return dflt
	}
	for r := range s {
		if !r.Contains(cell) {
			continue
		}
		if ff&Highlighted == Highlighted {
			return s.highlighted(s[r])
		}
		return s[r]
	}
	return dflt
}

func (s lineStyles) highlighted(sty tcell.Style) tcell.Style {
	fg, bg, _ := sty.Decompose()
	return sty.Foreground(bg).Background(fg)
}

type LineFlags uint

const (
	NotFocusable LineFlags = 1 << iota
	Highlighted
)
