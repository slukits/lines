// Copyright (c) 2022 Stephan Lukits. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// NOTE this is not the entry point to this package!  The central types
// are found in events.go, screen.go, component.go and env.go.  Also
// testing.go is a good place to start.

package lines

import "github.com/gdamore/tcell/v2"

type lines []*line

func (ll *lines) replace(cc ...[]byte) {
	new := []*line{}
	for _, c := range cc {
		new = append(new, &line{
			content: string(c),
			dirty:   true,
		})
	}
	*ll = new
}

func (ll *lines) append(cc ...[]byte) {
	for _, c := range cc {
		*ll = append(*ll, &line{
			content: string(c),
			dirty:   true,
		})
	}
}

func (ll lines) IsDirty() bool {
	for _, l := range ll {
		if !l.dirty {
			continue
		}
		return true
	}
	return false
}

func (ll lines) ForDirty(cb func(int, *line)) {
	for i, l := range ll {
		if !l.dirty {
			return
		}
		cb(i, l)
	}
}

func (ll lines) For(cb func(int, *line) (stop bool)) {
	for i, l := range ll {
		if cb(i, l) {
			return
		}
	}
}

type line struct {
	dirty   bool
	stale   string
	content string
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

func (l *line) sync(x, y, width int, rw runeWriter, sty tcell.Style) {
	l.dirty = false
	if len(l.content) >= len(l.stale) {
		l.setLonger(x, y, width, rw, sty)
	} else {
		l.setShorter(x, y, width, rw, sty)
	}
	l.stale = ""
}

func (l *line) setShorter(
	x, y, width int, rw runeWriter, sty tcell.Style,
) {

	base, add := len(l.content), len(l.stale)-len(l.content)
	l.setLonger(x, y, width, rw, sty)
	for i := 0; i < add; i++ {
		if i == width {
			break
		}
		rw.SetContent(x+base+i, y, ' ', nil, sty)
	}
}

func (l *line) setLonger(
	x, y, width int, rw runeWriter, sty tcell.Style,
) {

	for i, r := range l.content {
		if i == width {
			break
		}
		rw.SetContent(x+i, y, r, nil, sty)
	}
}
