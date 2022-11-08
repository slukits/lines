// Copyright (c) 2022 Stephan Lukits. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package lines

// GapsWriter allows to get more specific gap writer to create gaps
// between contents of different components.  Use [Component.Gaps]
// method of a component c to obtain a gaps writer for c.
type GapsWriter struct {
	gg    *gaps
	gm    gapMask
	level int
	sty   *Style

	// Top writes to the top gap, i.e. first line of a component's
	// screen area (plus provided level).
	Top *gapWriter

	// Bottom writes to the bottom gap, i.e. last line of a component's
	// screen area (minus provided level).
	Bottom *gapWriter

	// Left writes to the left gap, i.e. first column of a component's
	// screen area (plus provided level).
	Left *gapWriter

	// Right writes to the right gap, i.e. last column of a component's
	// screen area (minus provided level).
	Right *gapWriter

	// Horizontal writes to bottom and top gap.
	Horizontal *gapWriter

	// Vertical writes to left and right gap.
	Vertical *gapWriter

	// TopLeft writes to the top left corner at selected level.
	TopLeft *cornerWriter

	// TopRight writes to the top right corner at selected level.
	TopRight *cornerWriter

	// BottomRight writes to the bottom right corner at selected level.
	BottomRight *cornerWriter

	// BottomLeft writes to the bottom left corner at selected level.
	BottomLeft *cornerWriter

	// Corners writes to all corners at selected level.
	Corners *cornerWriter
}

func newGapsWriter(level int, gg *gaps) *GapsWriter {
	ggw := &GapsWriter{
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

func (ggw *GapsWriter) initStyle(sty Style) {
	ggw.gg.forStyler(func(s styler) {
		s.setDefaultStyle(ggw.level, sty)
	})
	ggw.sty = &sty
}

// AA stets given style attributes aa for selected gap-level.
func (ggw *GapsWriter) AA(aa StyleAttributeMask) *GapsWriter {
	if ggw.sty == nil {
		ggw.initStyle(ggw.gg.sty.WithAA(aa))
		return ggw
	}
	ggw.gg.forStyler(func(s styler) {
		s.withAA(ggw.level, aa)
	})
	return ggw
}

// FG sets given color c as foreground color for selected gap-level.
func (ggw *GapsWriter) FG(c Color) *GapsWriter {
	if ggw.sty == nil {
		ggw.initStyle(ggw.gg.sty.WithFG(c))
		return ggw
	}
	ggw.gg.forStyler(func(s styler) {
		s.withFG(ggw.level, c)
	})
	return ggw
}

// BG sets given color c as background color for selected gap-level.
func (ggw *GapsWriter) BG(c Color) *GapsWriter {
	if ggw.sty == nil {
		ggw.initStyle(ggw.gg.sty.WithBG(c))
		return ggw
	}
	ggw.gg.forStyler(func(s styler) {
		s.withBG(ggw.level, c)
	})
	return ggw
}

// Sty stets given style s as style for selected gap-level, i.e. sets
// style attributes and colors.
func (ggw *GapsWriter) Sty(s Style) *GapsWriter {
	ggw.initStyle(s)
	return ggw
}

// Filling returns a filling writer filling all gaps of selected level
// with what's printed to it.  Note [lines.Print] is needed to print to
// a filling writer.
func (ggw *GapsWriter) Filling() *allGapsFiller {
	return &allGapsFiller{ggw: ggw}
}

type allGapsFiller struct{ ggw *GapsWriter }

func (af *allGapsFiller) WriteAt(rr []rune) {
	for _, g := range selectGaps(af.ggw.gg, top|right|bottom|left) {
		g.setAtFilling(af.ggw.level, 0, rr[0])
	}
}

type allGapsWriter struct{ ggw *GapsWriter }

func (agg *allGapsWriter) Write(bb []byte) (int, error) {
	for _, g := range selectGaps(agg.ggw.gg, top|right|bottom|left) {
		g.set(agg.ggw.level, string(bb))
	}
	return len(bb), nil
}

type gapWriter struct {
	gm    gapMask
	ggw   *GapsWriter
	level int
	sty   *Style
}

func (w *gapWriter) Write(bb []byte) (int, error) {
	if len(bb) == 0 {
		return 0, nil
	}
	write := func(g *gap) {
		g.set(w.level, string(bb))
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

func (g *gapWriter) initStyle(sty Style) {
	for _, gp := range selectGaps(g.ggw.gg, g.gm) {
		gp.setDefaultStyle(g.level, sty)
	}
	g.sty = &sty
}

// AA stets given style attributes aa for selected gap's next write at
// selected level.
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

// FG stets given color c as foreground color for selected gap's next
// write at selected level.
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

// BG stets given color c as background color for selected gap's next
// write at selected level.
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

// Sty stets given style s for selected gap's next write at selected
// level, i.e. sets style attributes and colors.
func (g *gapWriter) Sty(s Style) *gapWriter {
	g.initStyle(s)
	return g
}

type gapAtWriter struct {
	ggw   *GapsWriter
	sty   *Style
	gm    gapMask
	level int
	at    int
}

// Filling indicates that the first rune of the next write is flagged as
// filling at the select position of selected gap at selected level.
func (aw *gapAtWriter) Filling() *gapAtWriter {
	aw.gm |= filling
	return aw
}

// WriteAt writes given runes rr at the select position of selected gap
// at selected level.  Applying optionally set style information over
// the range of printed runes.
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

// BG stets given style for the printed sequence of runes to this
// at-writer of selected gap(s) at selected gap-level.
func (aw *gapAtWriter) Sty(s Style) *gapAtWriter {
	aw.sty = &s
	return aw
}
