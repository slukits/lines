// Copyright (c) 2022 Stephan Lukits. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package lines

type gapsWriter struct {
	gg          *gaps
	gm          gapMask
	level       int
	sty         *Style
	Top         *gapWriter
	Bottom      *gapWriter
	Left        *gapWriter
	Right       *gapWriter
	Horizontal  *gapWriter
	Vertical    *gapWriter
	TopLeft     *cornerWriter
	TopRight    *cornerWriter
	BottomRight *cornerWriter
	BottomLeft  *cornerWriter
	Corners     *cornerWriter
}

func newGapsWriter(level int, gg *gaps) *gapsWriter {
	ggw := &gapsWriter{
		gg:          gg,
		gm:          top | right | bottom | left,
		level:       level,
		TopLeft:     &cornerWriter{gg: gg, cm: topLeft, level: level},
		TopRight:    &cornerWriter{gg: gg, cm: topRight, level: level},
		BottomRight: &cornerWriter{gg: gg, cm: bottomRight, level: level},
		BottomLeft:  &cornerWriter{gg: gg, cm: bottomLeft, level: level},
		Corners:     &cornerWriter{gg: gg, cm: allCorners, level: level},
	}
	ggw.Top = &gapWriter{ggw: ggw, gm: top, level: level}
	ggw.Bottom = &gapWriter{ggw: ggw, gm: bottom, level: level}
	ggw.Left = &gapWriter{ggw: ggw, gm: left, level: level}
	ggw.Right = &gapWriter{ggw: ggw, gm: right, level: level}
	ggw.Horizontal = &gapWriter{ggw: ggw, gm: top | bottom, level: level}
	ggw.Vertical = &gapWriter{ggw: ggw, gm: left | right, level: level}
	return ggw
}

type styler interface {
	setDefaultStyle(int, Style)
	withAA(int, StyleAttributeMask)
	withFG(int, Color)
	withBG(int, Color)
}

func (ggw *gapsWriter) initStyle(sty Style) {
	ggw.gg.forStyler(func(s styler) {
		s.setDefaultStyle(ggw.level, sty)
	})
	ggw.sty = &sty
}

// AA stets given style attributes of selected gap-level.
func (ggw *gapsWriter) AA(aa StyleAttributeMask) *gapsWriter {
	if ggw.sty == nil {
		ggw.initStyle(ggw.gg.sty.WithAA(aa))
		return ggw
	}
	ggw.gg.forStyler(func(s styler) {
		s.withAA(ggw.level, aa)
	})
	return ggw
}

func (ggw *gapsWriter) FG(c Color) *gapsWriter {
	if ggw.sty == nil {
		ggw.initStyle(ggw.gg.sty.WithFG(c))
		return ggw
	}
	ggw.gg.forStyler(func(s styler) {
		s.withFG(ggw.level, c)
	})
	return ggw
}

func (ggw *gapsWriter) BG(c Color) *gapsWriter {
	if ggw.sty == nil {
		ggw.initStyle(ggw.gg.sty.WithBG(c))
		return ggw
	}
	ggw.gg.forStyler(func(s styler) {
		s.withBG(ggw.level, c)
	})
	return ggw
}

func (ggw *gapsWriter) Filling() *allGapsWriter {
	return &allGapsWriter{ggw: ggw}
}

type allGapsWriter struct{ ggw *gapsWriter }

func (agg *allGapsWriter) Write(bb []byte) (int, error) {
	for _, g := range selectGaps(agg.ggw.gg, top|right|bottom|left) {
		g.set(agg.ggw.level, string(bb))
		g.filling(agg.ggw.level, true)
	}
	return len(bb), nil
}

type gapWriter struct {
	gm    gapMask
	ggw   *gapsWriter
	level int
	sty   *Style
}

func (w *gapWriter) Write(bb []byte) (int, error) {
	if len(bb) == 0 {
		return 0, nil
	}
	write := func(g *gap) {
		g.set(w.level, string(bb))
		g.filling(w.level, w.gm&filling != 0)
	}
	for _, g := range selectGaps(w.ggw.gg, w.gm) {
		write(g)
	}
	return len(bb), nil
}

func selectGaps(gg *gaps, gm gapMask) []*gap {
	switch gm & (top | right | bottom | left) {
	case top:
		return []*gap{&gg.top}
	case bottom:
		return []*gap{&gg.bottom}
	case left:
		return []*gap{&gg.left}
	case right:
		return []*gap{&gg.right}
	case top | bottom:
		return []*gap{&gg.top, &gg.bottom}
	case left | right:
		return []*gap{&gg.left, &gg.right}
	case top | right | bottom | left:
		return []*gap{&gg.top, &gg.right, &gg.bottom, &gg.left}
	}
	return nil
}

func (g *gapWriter) At(idx int) *gapAtWriter {
	return &gapAtWriter{
		ggw:   g.ggw,
		gm:    g.gm & (top | right | bottom | left),
		level: g.level,
		at:    idx,
	}
}

func (g *gapWriter) Filling() *gapWriter {
	g.gm |= filling
	return g
}

func (g *gapWriter) initStyle(sty Style) {
	for _, gp := range selectGaps(g.ggw.gg, g.gm) {
		gp.setDefaultStyle(g.level, sty)
	}
	g.sty = &sty
}

// AA stets given style attributes of selected gap-level.
func (g *gapWriter) AA(aa StyleAttributeMask) *gapWriter {
	if g.sty == nil {
		g.initStyle(g.ggw.gg.sty.WithAA(aa))
		return g
	}
	for _, gp := range selectGaps(g.ggw.gg, g.gm) {
		gp.withAA(g.level, aa)
	}
	return g
}

func (g *gapWriter) FG(c Color) *gapWriter {
	if g.sty == nil {
		g.initStyle(g.ggw.gg.sty.WithFG(c))
		return g
	}
	for _, gp := range selectGaps(g.ggw.gg, g.gm) {
		gp.withFG(g.level, c)
	}
	return g
}

func (g *gapWriter) BG(c Color) *gapWriter {
	if g.sty == nil {
		g.initStyle(g.ggw.gg.sty.WithBG(c))
		return g
	}
	for _, gp := range selectGaps(g.ggw.gg, g.gm) {
		gp.withBG(g.level, c)
	}
	return g
}

type gapAtWriter struct {
	ggw   *gapsWriter
	sty   *Style
	gm    gapMask
	level int
	at    int
}

func (aw *gapAtWriter) Filling() *gapAtWriter {
	aw.gm |= filling
	return aw
}

func (aw *gapAtWriter) WriteAt(rr []rune) {
	if len(rr) == 0 {
		return
	}
	write := func(g *gap) {
		if aw.gm&filling == 0 {
			if aw.sty == nil {
				g.setAt(aw.level, aw.at, rr)
			} else {
				g.setStyledAt(aw.level, aw.at, rr, aw.sty)
			}
		} else {
			if aw.sty == nil {
				g.setAtFilling(aw.level, aw.at, rr[0])
			} else {
				g.setStyledAtFilling(
					aw.level, aw.at, rr[0], aw.sty)
			}
		}
		g.filling(aw.level, false)
	}
	for _, g := range selectGaps(aw.ggw.gg, aw.gm) {
		write(g)
	}
}

func (aw *gapAtWriter) currentStyle() *Style {
	if aw.sty != nil {
		return aw.sty
	}
	return &aw.ggw.gg.sty
}

// AA stets given style attributes for the printed sequence of runes to
// this at-writer of selected gaps at selected gap-level.
func (aw *gapAtWriter) AA(aa StyleAttributeMask) *gapAtWriter {
	sty := aw.currentStyle().WithAA(aa)
	aw.sty = &sty
	return aw
}

// FG stets given foreground color for the printed sequence of runes to
// this at-writer of selected gaps at selected gap-level.
func (aw *gapAtWriter) FG(c Color) *gapAtWriter {
	sty := aw.currentStyle().WithFG(c)
	aw.sty = &sty
	return aw
}

// BG stets given background color for the printed sequence of runes to
// this at-writer of selected gaps at selected gap-level.
func (aw *gapAtWriter) BG(c Color) *gapAtWriter {
	sty := aw.currentStyle().WithBG(c)
	aw.sty = &sty
	return aw
}
