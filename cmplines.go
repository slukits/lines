// Copyright (c) 2022 Stephan Lukits. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// NOTE this is not the entry point to this package!  The central types
// are found in events.go, screen.go, component.go and env.go.  Also
// testing.go is a good place to start.  The lines type provides a
// sequence of lines which make up a component's content (not its screen
// representation).  The content of lines can be manipulated writing to
// an environment instance of type Env or to a writer provided by some
// of its methods.  Lines can not be accessed directly by a client since
// the lines do not need to represent a components logical lines due to
// line wrapping.

package lines

// LineFiller can be used in component content-lines indicating that a
// line l should fill up its whole width whereas its remaining empty
// space is spread equally over filler found in l.
const LineFiller = string(rune(29))

type lines []*line

// append given content lines to current content
func (ll *lines) append(
	lnFactory func() *line,
	ff LineFlags,
	sty Style,
	cc ...[]byte,
) {

	for _, c := range cc {
		l := lnFactory()
		l.content = string(c)
		l.ff = ff
		l.sty = sty
		*ll = append(*ll, l)
	}
}

// replaceAt replaces starting at given index the following lines with
// given content lines.  replaceAt is a no-op if idx < 0 or len(cc) == 0
func (ll *lines) replaceAt(
	lnFactory func() *line,
	idx, cell int,
	ff LineFlags,
	sty Style,
	cc ...[]byte,
) {
	if idx < 0 || len(cc) == 0 {
		return
	}
	for idx+len(cc) > len(*ll) {
		*ll = append(*ll, lnFactory())
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
