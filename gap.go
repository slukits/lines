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
	gm gapMask
	ll []*Line
}

func (g *gap) ensureLevel(l int) *Line {
	if len(g.ll) > l {
		return g.ll[l]
	}
	for i := len(g.ll); i <= l; i++ {
		g.ll = append(g.ll, &Line{ff: dirty})
	}
	return g.ll[l]
}

func (g *gap) reset(l int, sty *Style) {
	g.ensureLevel(l).reset(0, sty)
}

func (g *gap) setDefaultStyle(level int, s Style) {
	g.ensureLevel(level).setDefaultStyle(s)
}

func (g *gap) withAA(level int, aa StyleAttributeMask) {
	g.ensureLevel(level).withAA(aa)
}

func (g *gap) withFG(level int, c Color) {
	g.ensureLevel(level).withFG(c)
}

func (g *gap) withBG(level int, c Color) {
	g.ensureLevel(level).withBG(c)
}

func (g *gap) isDirty() bool {
	for _, l := range g.ll {
		if !l.isDirty() {
			continue
		}
		return true
	}
	return false
}

func (g *gap) set(level int, s string) {
	g.ensureLevel(level).set(s)
}

func (g *gap) setStyled(level int, s string, sty *Style) {
	g.ensureLevel(level).setStyled(s, *sty)
}

func (g *gap) setAt(level, at int, rr []rune) {
	g.ensureLevel(level).setAt(at, rr)
}

func (g *gap) setStyledAt(level, at int, rr []rune, sty *Style) {
	g.ensureLevel(level).setStyledAt(at, rr, *sty)
}

func (g *gap) setAtFilling(level, at int, r rune) {
	g.ensureLevel(level).setAtFilling(at, r)
}

func (g *gap) setStyledAtFilling(level, at int, r rune, sty *Style) {
	g.ensureLevel(level).setStyledAtFilling(at, r, *sty)
}
