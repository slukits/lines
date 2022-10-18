// Copyright (c) 2022 Stephan Lukits. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package lines

import (
	"bytes"

	"github.com/slukits/lines/internal/lyt"
)

// Mod sets how given components content is maintained.
func (c *component) Mod(cm ComponentMode) {
	switch cm {
	case Appending:
		c.mod &^= Overwriting | Tailing
		c.mod |= Appending
	case Overwriting:
		c.mod &^= Appending | Tailing
		c.mod |= Overwriting
	case Tailing:
		c.mod &^= Appending | Overwriting
		c.mod |= Tailing
	}
}

// Sty replaces a component's style attributes like bold or dimmed.
func (c *component) Sty(attr StyleAttribute) {
	c.fmt.sty = c.fmt.sty.WithAttrs(attr)
}

// FG replaces a component's foreground color.
func (c *component) FG(color Color) {
	c.dirty = true
	c.fmt.sty = c.fmt.sty.WithFG(color)
}

// BG replaces a component's background color.
func (c *component) BG(color Color) {
	c.dirty = true
	c.fmt.sty = c.fmt.sty.WithBG(color)
}

// Len returns the number of lines currently stored in a component.
// Note the line number is independent of a component's associated
// screen area.
func (c *component) Len() int {
	return len(*c.ll)
}

// IsDirty is true if this component is flagged dirty or one of its
// lines.
func (c *component) IsDirty() bool {
	if c.Len() == 0 {
		return c.dirty
	}
	return c.ll.IsDirty() || c.dirty
}

// SetDirty flags a component as dirty having the effect that at the
// next syncing the component's screen area is cleared before it is
// written to.
func (c *component) SetDirty() {
	c.dirty = true
}

// Dim provides a components layout dimensions and features to adapt
// them.
func (c *component) Dim() *lyt.Dim { return c.dim }

// All indicates for an operation with a line-index that the operation
// should be executed for all lines, e.g. Reset on a component.
const All = -1

// Reset blanks out the content of the line or all lines with given
// index the next time it is printed to the screen.  Provide line flags
// if for example a reset line should not be focusable.  If provided
// lines index is -1 (see All-constant) Rest scrolls to the top,
// truncates its lines to the screen-area-height and resets the
// remaining lines.
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
			l.reset(c.fmt.sty, _ff)
		}
		return
	}

	(*c.ll)[idx].reset(c.fmt.sty, _ff)
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
	sx, sy, sw, sh := c.dim.Area()
	if c.mod&Tailing == Tailing && c.Len() >= sh {
		c.setFirst(c.Len() - sh)
	}
	if c.dirty {
		c.clear(rw)
		c.dirty = false
		c.ll.For(c.first, func(i int, l *line) (stop bool) {
			if i >= sh {
				return true
			}
			l.sync(sx, sy+i, sw, rw)
			return false
		})
	}
	c.ll.ForDirty(c.first, func(i int, l *line) (stop bool) {
		if i >= sh {
			return true
		}
		l.sync(sx, sy+i, sw, rw)
		return false
	})
}

// clear fills the receiving component's printable area with spaces.
func (c *component) clear(rw runeWriter) {
	sx, sy, sw, sh := c.dim.Rect()
	for y := sy; y < sy+sh; y++ {
		for x := sx; x < sx+sw; x++ {
			rw.Display(x, y, ' ', c.fmt.sty)
		}
	}
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
	bb []byte, line, cell int, ff LineFlags, sty Style,
) (int, error) {
	switch {
	case c.mod&(Appending|Tailing) != 0:
		c.ll.append(
			c.lineFactory, ff, sty, bytes.Split(bb, []byte("\n"))...)
	default:
		if line == -1 {
			c.Reset(line)
			line = 0
		}
		c.ll.replaceAt(
			c.lineFactory, line, cell, ff, sty,
			bytes.Split(bb, []byte("\n"))...)
	}
	return len(bb), nil
}

func (c *component) lineFactory() *line {
	return &line{
		sty:    c.fmt.sty,
		dirty:  true,
		global: c.global,
	}
}
