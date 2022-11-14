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

func (gg *gaps) isDirty() bool {
	if gg == nil {
		return false
	}
	top := gg.top.isDirty()
	return top || gg.right.isDirty() ||
		gg.bottom.isDirty() || gg.left.isDirty() ||
		gg.topLeft.isDirty() || gg.topRight.isDirty() ||
		gg.bottomRight.isDirty() || gg.bottomLeft.isDirty()
}

func (gg *gaps) forStyler(cb func(styler)) {
	for _, s := range []styler{&gg.top, &gg.right, &gg.bottom, &gg.left,
		&gg.topLeft, &gg.topRight, &gg.bottomRight, &gg.bottomLeft} {

		cb(s)
	}
}

// sync writes the gaps into given area and returns an area reduced by
// gaps.
func (gg *gaps) sync(
	x, y, width, height int, w runeWriter, ggl *globals,
) {
	gglSty := ggl.Style(Default)
	gg.topLeft.sync(x, y, width, height, w, gglSty)
	gg.topRight.sync(x, y, width, height, w, gglSty)
	gg.bottomRight.sync(x, y, width, height, w, gglSty)
	gg.bottomLeft.sync(x, y, width, height, w, gglSty)
	th := gg.top.sync(x+1, y, width-2, height, w, ggl)
	gg.bottom.sync(x+1, y+th, width-2, height-th, w, ggl)
	lw := gg.left.sync(x, y+1, width, height-2, w, ggl)
	gg.right.sync(x+lw, y+1, width-lw, height-2, w, ggl)
}
