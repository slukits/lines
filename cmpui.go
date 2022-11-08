// Copyright (c) 2022 Stephan Lukits. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package lines

import (
	"bytes"

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

// Len returns the number of lines currently stored in a component.
// Note the number of component lines is independent of a component's
// available screen lines.
func (c *component) Len() int {
	return len(*c.ll)
}

// IsDirty is true if given component c is flagged dirty or one of its
// lines or gaps.
func (c *component) IsDirty() bool {
	ll, gg := c.ll.IsDirty(), c.gaps.isDirty()
	return ll || gg || c.dirty
}

// SetDirty flags a component as dirty having the effect that at the
// next syncing the component's screen area is cleared before it is
// written to.  Note usually you don't need this method since a
// component is automatically flagged dirty if its layout changed, if
// one of its global properties changed etc.
func (c *component) SetDirty() {
	c.dirty = true
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
		_, _, _, height := c.Dim().Area()
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
// content is written to the screen.
func (c *component) hardSync(rw runeWriter) {
	if !c.dirty {
		c.dirty = true
	}
	c.sync(rw)
}

// sync writes receiving components lines to the screen.
func (c *component) sync(rw runeWriter) {
	if c.dim.IsOffScreen() {
		if !c.dirty {
			c.dirty = true
		}
		return
	}
	sx, sy, sw, sh := c.dim.Area()
	if c.mod&Tailing == Tailing && c.Len() >= sh {
		c.setFirst(c.Len() - sh)
	}
	if c.dirty {
		c.syncCleared(rw)
		return
	}
	if c.gaps != nil && c.gaps.isDirty() {
		sx, sy, sw, sh = c.gaps.sync(sx, sy, sw, sh, rw, c.gg)
	}
	gg := c.gaps.isDirty()
	_ = gg
	if sw <= 0 || sh <= 0 {
		return
	}
	c.ll.ForDirty(c.first, func(i int, l *Line) (stop bool) {
		if i >= sh {
			return true
		}
		l.sync(sx, sy+i, sw, rw, c.gg)
		return false
	})
}

// clear fills the receiving component's printable area with spaces.
func (c *component) syncCleared(rw runeWriter) {
	sx, sy, sw, sh := c.dim.Rect()
	for y := sy; y < sy+sh; y++ {
		for x := sx; x < sx+sw; x++ {
			rw.Display(x, y, ' ', c.gg.Style(Default))
		}
	}
	c.dirty = false
	sx, sy, sw, sh = c.dim.Area()
	if c.gaps != nil {
		sx, sy, sw, sh = c.gaps.sync(sx, sy, sw, sh, rw, c.gg)
	}
	if sw <= 0 || sh <= 0 {
		return
	}
	c.ll.For(c.first, func(i int, l *Line) (stop bool) {
		if i >= sh {
			return true
		}
		l.sync(sx, sy+i, sw, rw, c.gg)
		return false
	})
}

// setFirst sets the first displayed line and in case it changes given
// component becomes also dirty (hence the indirection).
func (c *component) setFirst(f int) {
	if f < 0 || f == c.first || f >= c.Len() {
		return
	}

	c.first = f
	if c.dirty {
		return
	}
	c.dirty = true
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
