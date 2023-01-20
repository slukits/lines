// Copyright (c) 2022 Stephan Lukits. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package lines

import (
	"bytes"

	"github.com/slukits/ints"
	"github.com/slukits/lines/internal/lyt"
)

// AA replaces a component's style attributes like bold or dimmed.  Note
// changing the style attributes of the default style of the Lines'
// globals will have no effect on components whose style attributes has
// been set.
func (c *component) AA(aa StyleAttributeMask) *component {
	c.gg.SetAA(Default, aa) // globals update listener sets c dirty
	return c
}

// FG replaces a component's foreground color.  Note changing the
// foreground color of the default style of the Lines' globals will have
// no effect on components whose foreground color has been set.
func (c *component) FG(color Color) *component {
	c.gg.SetFG(Default, color) // globals update listener sets c dirty
	return c
}

// BG replaces a component's background color.  Note changing the
// background color of the default style of the Lines' globals will have
// no effect on components whose background color has been set.
func (c *component) BG(color Color) *component {
	c.gg.SetBG(Default, color) // globals update listener sets c dirty
	return c
}

// Sty replaces a component's style, i.e. its style attributes and its
// fore- and background color.  Note changing the default style of the
// Lines' globals will have no effect on components whose style has been
// set.
func (c *component) Sty(s Style) *component {
	c.gg.SetStyle(Default, s) // globals update listener sets c dirty
	return c
}

// Len returns the number of lines currently maintained by a component.
// Note the number of component lines is independent of a component's
// available screen lines.
func (c *component) Len() int {
	if c.Src != nil {
		if sl, ok := c.Src.Liner.(ScrollableLiner); ok {
			return sl.Len()
		}
	}
	return len(*c.ll)
}

// IsDirty is true if given component c is flagged dirty or one of its
// lines or gaps.
func (c *component) IsDirty() bool {
	return c.dirty || c.cursorMoved || c.ll.IsDirty() ||
		c.gaps.isDirty() || c.Src.IsDirty() || c.haveDirtyNested()
}

func (c *component) isDirtyButCursorMoved() bool {
	return c.dirty || c.ll.IsDirty() || c.gaps.isDirty() ||
		c.Src.IsDirty() || c.haveDirtyNested()
}

func (c *component) haveDirtyNested() (dirty bool) {
	c.forNested(func(n *component) (stop bool) {
		if !n.IsDirty() {
			return false
		}
		dirty = true
		return true
	})
	return dirty
}

// SetDirty flags a component as dirty having the effect that at the
// next syncing the component's screen area is cleared before it is
// written to.  Note usually you don't need this method since a
// component is automatically flagged dirty if its layout changed, if
// one of its global properties changed etc.
func (c *component) SetDirty() {
	if c.dirty {
		return
	}
	c.dirty = true
	c.forNested(func(n *component) (stop bool) {
		n.SetDirty()
		return false
	})
}

// Dim provides a components layout dimensions and features to adapt
// them.
func (c *component) Dim() *lyt.Dim { return c.dim }

// All indicates for an operation with a line-index that the operation
// should be executed for all lines, e.g. [Component.Reset] on a component.
const All = -1

// Reset blanks out the content of the line (or all lines) with given
// index the next time it is printed to the screen.  Provide line flags
// if for example a Reset line should not be focusable:
//
//	c.Reset(lines.All, lines.NotFocusable)
//
// If provided lines index is -1, see [All]-constant, Reset scrolls to
// the top, truncates its lines to the available screen-lines and resets
// the remaining lines.
func (c *component) Reset(idx int, ff ...LineFlags) {
	if idx < -1 || idx >= c.Len() {
		return
	}
	_ff := LineFlags(0)
	for _, f := range ff {
		_ff |= f
	}

	if idx == -1 {
		c.setFirst(0)
		height := c.ContentScreenLines()
		if len(*c.ll) > height {
			ll := (*c.ll)[:height]
			c.ll = &ll
		}
		for _, l := range *c.ll {
			l.reset(_ff, nil)
		}
		return
	}

	(*c.ll)[idx].reset(_ff, nil)
}

// hardSync clears the screen area of receiving component before its
// content is (re)written to the screen.
func (c *component) hardSync(rw runeWriter) {
	if !c.dirty {
		c.SetDirty()
	}
	c.sync(rw)
}

// sync writes receiving components lines to the screen.
func (c *component) sync(rw runeWriter) {
	if c.cursorMoved {
		c.cursorMoved = false
		if !c.isDirtyButCursorMoved() {
			return
		}
	}
	if c.dim.IsOffScreen() {
		if !c.dirty {
			c.dirty = true
		}
		return
	}
	cll := c.ContentScreenLines()
	if c.mod&Tailing == Tailing && c.Len() >= cll {
		c.setFirst(c.Len() - cll)
	}
	if c.dirty {
		c.syncCleared(rw)
	} else {
		c.syncContent(rw)
	}
	c.forNested(func(n *component) (stop bool) {
		n.sync(rw)
		return false
	})
}

func (c *component) syncContent(rw runeWriter) {
	if c.gaps != nil && c.gaps.isDirty() {
		gx, gy, gw, gh := c.Dim().Printable()
		c.gaps.sync(gx, gy, gw, gh, rw, c.gg)
	}
	cx, cy, cw, ch := c.ContentArea()
	if cw <= 0 || ch <= 0 {
		return
	}
	if c.Src.IsDirty() {
		c.Src.cleanup(c)
	}
	c.ll.ForDirty(c._first, func(i int, l *Line) (stop bool) {
		if i >= ch {
			return true
		}
		l.sync(cx, cy+i, cw, rw, c.gg)
		return false
	})
}

// clear fills the receiving component's printable area with spaces.
func (c *component) syncCleared(rw runeWriter) {
	if c.first()+c.ContentScreenLines() > c.Len() {
		c.setFirst(ints.Max(0, c.Len()-c.ContentScreenLines()))
	}
	cx, cy, cw, ch := c.dim.Screen()
	for y := cy; y < cy+ch; y++ {
		for x := cx; x < cx+cw; x++ {
			rw.Display(x, y, ' ', c.gg.Style(Default))
		}
	}
	c.dirty = false
	if c.gaps != nil {
		gx, gy, gw, gh := c.Dim().Printable()
		c.gaps.sync(gx, gy, gw, gh, rw, c.gg)
	}
	cx, cy, cw, ch = c.ContentArea()
	if cw <= 0 || ch <= 0 {
		return
	}
	if c.Src.IsDirty() {
		c.Src.cleanup(c)
	} else if c.Src != nil {
		if len(*c.ll) > ch { // keep no of lines and screen hight in sync
			*c.ll = (*c.ll)[:ch]
		}
		c.Src.sync(ch, c)
	}
	c.ll.For(c._first, func(i int, l *Line) (stop bool) {
		if i >= ch {
			return true
		}
		l.sync(cx, cy+i, cw, rw, c.gg)
		return false
	})
}

func (c *component) forNested(cb func(n *component) (stop bool)) {
	lc := c.userCmp.layoutComponent()
	var nested func(func(Dimer) (stop bool))
	if s, ok := lc.(lyt.Stacker); ok {
		nested = s.ForStacked
	} else { // ignore chainer if stacker and chainer are provided
		if c, ok := lc.(lyt.Chainer); ok {
			nested = c.ForChained
		}
	}
	if nested == nil {
		return
	}
	nested(func(d Dimer) (stop bool) {
		n := d.(layoutComponenter).wrapped()
		return cb(n)
	})
}

func (c *component) ContentArea() (x, y, w, h int) {
	x, y, w, h = c.dim.Printable()
	if c.gaps == nil {
		return
	}
	return x + len(c.gaps.left.ll),
		y + len(c.gaps.top.ll),
		w - len(c.gaps.left.ll) - len(c.gaps.right.ll),
		h - len(c.gaps.top.ll) - len(c.gaps.bottom.ll)
}

func (c *component) ContentScreenLines() int {
	_, _, _, sh := c.dim.Printable()
	if c.gaps == nil {
		return sh
	}
	return sh - (len(c.gaps.top.ll) + len(c.gaps.bottom.ll))
}

// setFirst sets the first displayed line and in case it changes given
// component it becomes also dirty (hence the indirection).  setFirst
// takes also an optionally set component source into account.
func (c *component) setFirst(f int) {
	if c.Src != nil {
		c.Src.setFirst(f)
		return
	}

	if f < 0 || f == c._first || f >= c.Len() {
		return
	}

	c._first = f
	if c.dirty {
		return
	}
	c.dirty = true
}

func (c *component) ensureAsManyLineInstancesAsScreenLines() {
	bb := [][]byte{}
	for i := 0; i < c.ContentScreenLines()-len(*c.ll); i++ {
		bb = append(bb, []byte(""))
	}
	if len(bb) == 0 {
		return
	}
	sty := c.globals().Style(Default)
	c.ll.replaceAt(len(*c.ll), -1, &sty, bb...)
}

// first returns the index of the first line displayed taking an
// optionally set component source into account.
func (c *component) first() int {
	if c.Src != nil {
		return c.Src.first
	}
	return c._first
}

func (c *component) write(
	bb []byte, line, cell int, sty *Style,
) (int, error) {
	switch {
	case c.mod&(Appending|Tailing) != 0:
		c.ll.append(sty, bytes.Split(bb, []byte("\n"))...)
	default:
		if line == -1 {
			c.Reset(line)
			line = 0
		}
		c.ll.replaceAt(
			line, cell, sty,
			bytes.Split(bb, []byte("\n"))...)
	}
	return len(bb), nil
}

func (c *component) writeAt(
	rr []rune, line, cell int, sty *Style,
) {
	if line < 0 || cell < 0 || len(rr) == 0 {
		return
	}
	l := c.ll.padded(line)
	if sty == nil {
		l.setAt(cell, rr)
	} else {
		l.setStyledAt(cell, rr, *sty)
	}
}

func (c *component) writeAtFilling(
	r rune, line, cell int, sty *Style,
) {
	if line < 0 || cell < 0 {
		return
	}
	l := c.ll.padded(line)
	if sty == nil {
		l.setAtFilling(cell, r)
	} else {
		l.setStyledAtFilling(cell, r, *sty)
	}
}
