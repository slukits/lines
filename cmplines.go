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

const filler = rune(29)

// A Filler in a string printed to a component's environment e indicates
// that a line l should fill up its whole width whereas its remaining
// empty space is distributed equally over filler found in l.
//
//	fmt.Fprint(e.LL(0), lines.Filler+"centered first line"+lines.Filler)
//
// See [EnvAtWriter.Filling] for a more sophisticated filling mechanism.
const Filler = string(filler)

// ComponentLines provides the API to manipulate ui-aspects of a
// component's lines like which line has the focus.  A component's lines
// are accessed through its LL-property.  To manipulate their content
// print to an Env(ironment) instance provided to an event listener
// implementation.
type ComponentLines struct {
	c     *Component
	Focus *LineFocus
}

// Mod sets how given component lines cll are maintained.
func (cll *ComponentLines) Mod(cm ComponentMode) {
	switch cm {
	case Appending:
		cll.c.mod &^= Overwriting | Tailing
		cll.c.mod |= Appending
	case Overwriting:
		cll.c.mod &^= Appending | Tailing
		cll.c.mod |= Overwriting
	case Tailing:
		cll.c.mod &^= Appending | Overwriting
		cll.c.mod |= Tailing
	}
}

// AA set the style attributes of all content lines of associated
// component to given style attributes aa.
func (cll *ComponentLines) AA(aa StyleAttributeMask) {
	for _, l := range *cll.c.ll {
		l.withAA(aa)
	}
}

// FG set the foreground color of all content lines of associated
// component to given color c.
func (cll *ComponentLines) FG(c Color) {
	for _, l := range *cll.c.ll {
		l.withFG(c)
	}
}

// BG set the background color of all content lines of associated
// component to given color c.
func (cll *ComponentLines) BG(c Color) {
	for _, l := range *cll.c.ll {
		l.withBG(c)
	}
}

// Len returns the number of component lines, which is independent from
// the number of screen lines.
func (cll *ComponentLines) Len() int { return cll.c.Len() }

// By returns the component line with given non negative index idx.  By
// panics if idx is negative.  Is idx < [ComponentLines.Len] lines are
// padded accordingly.
func (cll *ComponentLines) By(idx int) *Line {
	if idx < 0 {
		panic("lines: component lines: negative line index given")
	}
	return cll.c.ll.padded(idx)
}

func newComponentLines(c *Component) *ComponentLines {
	return &ComponentLines{
		c:     c,
		Focus: &LineFocus{c: c, current: -1, hlType: Highlighted},
	}
}

type lines []*Line

// append given content lines to current content
func (ll *lines) append(sty *Style, cc ...[]byte) {

	for _, c := range cc {
		l := Line{rr: []rune(string(c)), ff: dirty}
		if sty != nil {
			l.setDefaultStyle(*sty)
		}
		*ll = append(*ll, &l)
	}
}

// replaceAt replaces starting at given index the following lines with
// given content lines.  replaceAt is a no-op if idx < 0 or len(cc) == 0
func (ll *lines) replaceAt(
	idx, cell int,
	sty *Style,
	cc ...[]byte,
) {
	if idx < 0 || len(cc) == 0 {
		return
	}
	for idx+len(cc) > len(*ll) {
		l := Line{ff: dirty}
		*ll = append(*ll, &l)
	}
	l := (*ll)[idx]
	if sty == nil {
		l.setAt(cell, []rune(string(cc[0])))
	} else {
		l.setStyledAt(cell, []rune(string(cc[0])), *sty)
	}
	if len(cc) < 2 {
		return
	}
	for i := 1; i < len(cc); i++ {
		l := (*ll)[idx+i]
		l.set(string(cc[i]))
		if sty != nil {
			l.setDefaultStyle(*sty)
		}
	}
}

func (ll *lines) padded(idx int) *Line {
	if idx < len(*ll) {
		return (*ll)[idx]
	}
	for idx >= len(*ll) || idx == 0 && len(*ll) == 0 {
		l := Line{ff: dirty}
		*ll = append(*ll, &l)
	}
	return (*ll)[idx]
}

// IsDirty returns true if on of the lines is dirty.
func (ll lines) IsDirty() bool {
	if len(ll) == 0 {
		return false
	}
	for _, l := range ll {
		if !l.isDirty() {
			continue
		}
		return true
	}
	return false
}

// ForDirty calls back for every dirty line.
func (ll lines) ForDirty(offset int, cb func(int, *Line) (stop bool)) {
	for i, l := range ll[offset:] {
		if !l.isDirty() {
			continue
		}
		if cb(i, l) {
			return
		}
	}
}

// For calls back for every line of given lines ll starting at given
// offset.
func (ll lines) For(offset int, cb func(int, *Line) (stop bool)) {
	for i, l := range ll[offset:] {
		if cb(i, l) {
			return
		}
	}
}
