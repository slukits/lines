// Copyright (c) 2022 Stephan Lukits. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package lines

import "github.com/slukits/lines/internal/api"

type cornerWriter struct {
	gg    *gaps
	sty   *Style
	cm    cornerMask
	level int
}

func (g *cornerWriter) Write(bb []byte) (int, error) {
	var rr []rune
	if len(bb) == 0 {
		rr = []rune(" ")
	} else {
		rr = []rune(string(bb))
	}
	switch g.cm & allCorners {
	case topLeft:
		g.gg.topLeft.set(g.level, rr[0], g.sty)
	case topRight:
		g.gg.topRight.set(g.level, rr[0], g.sty)
	case bottomRight:
		g.gg.bottomRight.set(g.level, rr[0], g.sty)
	case bottomLeft:
		g.gg.bottomLeft.set(g.level, rr[0], g.sty)
	case allCorners:
		tl, tr, br, bl := cornerRunes(rr)
		g.gg.topLeft.set(g.level, tl, g.sty)
		g.gg.topRight.set(g.level, tr, g.sty)
		g.gg.bottomRight.set(g.level, br, g.sty)
		g.gg.bottomLeft.set(g.level, bl, g.sty)
	}
	return len(bb), nil
}

func (cr *cornerWriter) currentStyle() *Style {
	if cr.sty != nil {
		return cr.sty
	}
	return &cr.gg.sty
}

// AA stets given style attributes for the printed corner-rune to
// this corner-writer of selected corners.
func (cw *cornerWriter) AA(aa StyleAttributeMask) *cornerWriter {
	sty := cw.currentStyle().WithAA(aa)
	cw.sty = &sty
	return cw
}

// FG stets given foreground color for the printed corner-rune to
// this corner-writer of selected corners.
func (cw *cornerWriter) FG(c Color) *cornerWriter {
	sty := cw.currentStyle().WithFG(c)
	cw.sty = &sty
	return cw
}

// BG stets given background color for the printed corner-rune to
// this corner-writer of selected corners.
func (cw *cornerWriter) BG(c Color) *cornerWriter {
	sty := cw.currentStyle().WithBG(c)
	cw.sty = &sty
	return cw
}

func cornerRunes(rr []rune) (tl, tr, br, bl rune) {
	switch len(rr) {
	case 1:
		tl, tr, br, bl = rr[0], rr[0], rr[0], rr[0]
	case 4:
		tl, tr, br, bl = rr[0], rr[1], rr[2], rr[3]
	}
	return tl, tr, br, bl
}

type cornerMask uint

const (
	topLeft cornerMask = 1 << iota
	topRight
	bottomRight
	bottomLeft

	allCorners cornerMask = topLeft | topRight | bottomRight | bottomLeft
)

type corner struct {
	cm    cornerMask
	rr    []cornerRune
	dirty bool
}

func (c *corner) ensureLevel(l int) *cornerRune {
	if len(c.rr) > l {
		return &c.rr[l]
	}
	for i := len(c.rr); i <= l; i++ {
		c.rr = append(c.rr, cornerRune{})
	}
	return &c.rr[l]
}

func (c *corner) isDirty() bool { return c.dirty }

func (c *corner) setDefaultStyle(level int, s Style) {
	if !c.dirty {
		c.dirty = true
	}
	c.ensureLevel(level).setDefaultStyle(s)
}

func (c *corner) withAA(level int, aa StyleAttributeMask) {
	if !c.dirty {
		c.dirty = true
	}
	c.ensureLevel(level).withAA(aa)
}

func (cr *corner) withFG(level int, c Color) {
	if !cr.dirty {
		cr.dirty = true
	}
	cr.ensureLevel(level).withFG(c)
}

func (cr *corner) withBG(level int, c Color) {
	if !cr.dirty {
		cr.dirty = true
	}
	cr.ensureLevel(level).withBG(c)
}

func (c *corner) set(level int, r rune, s *Style) {
	if !c.dirty {
		c.dirty = true
	}
	cr := c.ensureLevel(level)
	cr.r = r
	cr.sty = s
	// if s != nil {
	// 	cr.sty = s
	// }
}

func (c *corner) sync(x, y, width, height int, w runeWriter, sty Style) {

	if c.dirty {
		c.dirty = false
	}

	switch c.cm & allCorners {
	case topLeft:
		c.syncTL(x, y, w, sty)
	case topRight:
		c.syncTR(x+width-1, y, w, sty)
	case bottomRight:
		c.syncBR(x+width-1, y+height-1, w, sty)
	case bottomLeft:
		c.syncBL(x, y+height-1, w, sty)
	}
}

func (c *corner) syncTL(x, y int, w runeWriter, dflt Style) {
	for i, cr := range c.rr {
		r, sty := cr.display(&dflt)
		w.Display(x+i, y+i, r, *sty)
	}
}

func (c *corner) syncTR(x, y int, w runeWriter, dflt Style) {
	for i, cr := range c.rr {
		r, sty := cr.display(&dflt)
		w.Display(x-i, y+i, r, *sty)
	}
}

func (c *corner) syncBR(x, y int, w runeWriter, dflt Style) {
	for i, cr := range c.rr {
		r, sty := cr.display(&dflt)
		w.Display(x-i, y-i, r, *sty)
	}
}

func (c *corner) syncBL(x, y int, w runeWriter, dflt Style) {
	for i, cr := range c.rr {
		r, sty := cr.display(&dflt)
		w.Display(x+i, y-i, r, *sty)
	}
}

type cornerRune struct {
	sty *Style
	r   rune
}

func (r *cornerRune) setDefaultStyle(s Style) { r.sty = &s }

func (r *cornerRune) ensureStyle() *Style {
	if r.sty != nil {
		return r.sty
	}
	sty := api.DefaultStyle
	r.sty = &sty
	return r.sty
}

func (r *cornerRune) withAA(aa StyleAttributeMask) {
	sty := r.ensureStyle().WithAA(aa)
	r.sty = &sty
}

func (r *cornerRune) withFG(c Color) {
	sty := r.ensureStyle().WithFG(c)
	r.sty = &sty
}

func (r *cornerRune) withBG(c Color) {
	sty := r.ensureStyle().WithBG(c)
	r.sty = &sty
}

func (rn *cornerRune) display(dflt *Style) (rune, *Style) {
	r := rn.r
	if r == 0 {
		r = ' '
	}
	if rn.sty == nil {
		return r, dflt
	}
	return r, rn.sty
}
