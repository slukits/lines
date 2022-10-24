// Copyright (c) 2022 Stephan Lukits. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package lines

type gapMask uint

const (
	top gapMask = 1 << iota
	right
	bottom
	left
	filling
)

type gap struct {
	gm    gapMask
	ll    []line
	dirty bool
}

func (g *gap) ensureLevel(l int) *line {
	if len(g.ll) > l {
		return &g.ll[l]
	}
	for i := len(g.ll); i <= l; i++ {
		g.ll = append(g.ll, line{})
	}
	return &g.ll[l]
}

func (g *gap) setDefaultStyle(level int, s Style) {
	if !g.dirty {
		g.dirty = true
	}
	g.ensureLevel(level).setDefaultStyle(s)
}

func (g *gap) withAA(level int, aa StyleAttributeMask) {
	if !g.dirty {
		g.dirty = true
	}
	g.ensureLevel(level).withAA(aa)
}

func (g *gap) withFG(level int, c Color) {
	if !g.dirty {
		g.dirty = true
	}
	g.ensureLevel(level).withFG(c)
}

func (g *gap) withBG(level int, c Color) {
	if !g.dirty {
		g.dirty = true
	}
	g.ensureLevel(level).withBG(c)
}

func (g *gap) isDirty() bool { return g.dirty }

func (g *gap) set(level int, s string) {
	if !g.dirty {
		g.dirty = true
	}
	g.ensureLevel(level).set(s)
}

func (g *gap) setAt(level, at int, rr []rune) {
	if !g.dirty {
		g.dirty = true
	}
	g.ensureLevel(level).setAt(at, rr)
}

func (g *gap) setStyledAt(level, at int, rr []rune, sty *Style) {
	if !g.dirty {
		g.dirty = true
	}
	g.ensureLevel(level).setStyledAt(at, rr, *sty)
}

func (g *gap) setAtFilling(level, at int, r rune) {
	if !g.dirty {
		g.dirty = true
	}
	g.ensureLevel(level).setAtFilling(at, r)
}

func (g *gap) setStyledAtFilling(level, at int, r rune, sty *Style) {
	if !g.dirty {
		g.dirty = true
	}
	g.ensureLevel(level).setStyledAtFilling(at, r, *sty)
}

func (g *gap) sync(
	x, y, width, height int, rw runeWriter, gg *globals,
) int {

	if g.dirty {
		g.dirty = false
	}

	switch g.gm & (top | right | bottom | left) {
	case top:
		return g.syncTop(x, y, width, height, rw, gg)
	case bottom:
		return g.syncBottom(x, y, width, height, rw, gg)
	case left:
		return g.syncLeft(x, y, width, height, rw, gg)
	case right:
		return g.syncRight(x, y, width, height, rw, gg)
	}

	return 0
}

func (g *gap) syncTop(
	x, y, width, height int, rw runeWriter, gg *globals,
) int {

	for i, l := range g.ll {
		if width <= 0 || i == height {
			return i
		}
		rr, ss := l.display(width, gg)
		for j, r := range rr {
			if j == width {
				break
			}
			rw.Display(x+j, y+i, r, ss.of(j))
		}
		x++
		width -= 2
	}

	return len(g.ll)
}

func (g *gap) syncBottom(
	x, y, width, height int, rw runeWriter, gg *globals,
) int {
	for i, l := range g.ll {
		if width <= 0 || i == height {
			return i
		}
		rr, ss := l.display(width, gg)
		for j, r := range rr {
			if j == width {
				break
			}
			rw.Display(x+j, y+height-(i+1), r, ss.of(j))
		}
		x++
		width -= 2
	}

	return len(g.ll)
}

func (g *gap) syncLeft(
	x, y, width, height int, rw runeWriter, gg *globals,
) int {
	for i, l := range g.ll {
		if height <= 0 || i == width {
			return i
		}
		rr, ss := l.display(height, gg)
		for j, r := range rr {
			if j == height {
				break
			}
			rw.Display(x+i, y+j, r, ss.of(j))
		}
		y++
		height -= 2
	}
	return len(g.ll)
}

func (g *gap) syncRight(
	x, y, width, height int, rw runeWriter, gg *globals,
) int {
	for i, l := range g.ll {
		if height <= 0 || i == width {
			return i
		}
		rr, ss := l.display(height, gg)
		for j, r := range rr {
			if j == height {
				break
			}
			rw.Display(x+width-(i+1), y+j, r, ss.of(j))
		}
		y++
		height -= 2
	}
	return len(g.ll)
}
