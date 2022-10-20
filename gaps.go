// Copyright (c) 2022 Stephan Lukits. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package lines

type gaps struct {
	sty                                        Style
	top, right, bottom, left                   gap
	topLeft, topRight, bottomLeft, bottomRight corner
}

func newGaps(sty Style) *gaps {
	return &gaps{
		sty:         sty,
		top:         gap{gm: top},
		right:       gap{gm: right},
		bottom:      gap{gm: bottom},
		left:        gap{gm: left},
		topLeft:     corner{cm: topLeft},
		topRight:    corner{cm: topRight},
		bottomRight: corner{cm: bottomRight},
		bottomLeft:  corner{cm: bottomLeft},
	}
}

func (g *gaps) isDirty() bool {
	if g == nil {
		return false
	}
	return g.top.isDirty() || g.right.isDirty() ||
		g.bottom.isDirty() || g.left.isDirty() ||
		g.topLeft.isDirty() || g.topRight.isDirty() ||
		g.bottomRight.isDirty() || g.bottomLeft.isDirty()
}

func (g *gaps) forStyler(cb func(styler)) {
	for _, s := range []styler{&g.top, &g.right, &g.bottom, &g.left,
		&g.topLeft, &g.topRight, &g.bottomRight, &g.bottomLeft} {

		cb(s)
	}
}

// sync writes the gaps into given area and returns an area reduced by
// gaps.
func (g *gaps) sync(
	x, y, width, height int, w runeWriter,
) (_, _, _, _ int) {
	g.topLeft.sync(x, y, width, height, w, g.sty)
	g.topRight.sync(x, y, width, height, w, g.sty)
	g.bottomRight.sync(x, y, width, height, w, g.sty)
	g.bottomLeft.sync(x, y, width, height, w, g.sty)
	th := g.top.sync(x+1, y, width-2, height, w, g.sty)
	bh := g.bottom.sync(x+1, y+th, width-2, height-th, w, g.sty)
	lw := g.left.sync(x, y+1, width, height-2, w, g.sty)
	rw := g.right.sync(x+lw, y+1, width-lw, height-2, w, g.sty)
	return x + lw, y + th, width - (lw + rw), height - (th + bh)
}
