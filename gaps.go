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
		topLeft:     corner{cm: topLeft, dirty: true},
		topRight:    corner{cm: topRight, dirty: true},
		bottomRight: corner{cm: bottomRight, dirty: true},
		bottomLeft:  corner{cm: bottomLeft, dirty: true},
	}
}

func (gg *gaps) isDirty() bool {
	if gg == nil {
		return false
	}
	return gg.top.isDirty() || gg.right.isDirty() ||
		gg.bottom.isDirty() || gg.left.isDirty() ||
		gg.topLeft.isDirty() || gg.topRight.isDirty() ||
		gg.bottomRight.isDirty() || gg.bottomLeft.isDirty()
}

func (gg *gaps) Len() (top, right, bottom, left int) {
	if gg == nil {
		return 0, 0, 0, 0
	}
	return len(gg.top.ll), len(gg.right.ll), len(gg.bottom.ll),
		len(gg.left.ll)
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
	x, y, width, height int, w runeWriter, ggl *Globals,
) {
	gglSty := ggl.Style(Default)
	gg.topLeft.sync(x, y, width, height, w, gglSty)
	gg.topRight.sync(x, y, width, height, w, gglSty)
	gg.bottomRight.sync(x, y, width, height, w, gglSty)
	gg.bottomLeft.sync(x, y, width, height, w, gglSty)
	th := gg.syncTop(x, y, width, height, w, ggl)
	gg.syncBottom(x, y+th, width, height-th, w, ggl)
	lw := gg.syncLeft(x, y, width, height, w, ggl)
	gg.syncRight(x+lw, y, width-lw, height, w, ggl)
}

func (gg *gaps) syncTop(
	x, y, width, height int, rw runeWriter, ggl *Globals,
) int {
	for i, l := range gg.top.ll {
		if width <= 0 || i == height {
			return i
		}
		if len(gg.left.ll) >= i+1 {
			width--
			x++
		}
		if len(gg.right.ll) >= i+1 {
			width--
		}
		l.sync(x, y+i, width, rw, ggl)
	}
	return len(gg.top.ll)
}

func (gg *gaps) syncBottom(
	x, y, width, height int, rw runeWriter, ggl *Globals,
) int {
	for i, l := range gg.bottom.ll {
		if width <= 0 || i == height {
			return i
		}
		if len(gg.left.ll) >= i+1 {
			width--
			x++
		}
		if len(gg.right.ll) >= i+1 {
			width--
		}
		l.sync(x, y+height-(i+1), width, rw, ggl)
	}
	return len(gg.bottom.ll)
}

func (gg *gaps) syncLeft(
	x, y, width, height int, rw runeWriter, ggl *Globals,
) int {
	for i, l := range gg.left.ll {
		if height <= 0 || i == width {
			return i
		}
		if len(gg.top.ll) >= i+1 {
			height--
			y++
		}
		if len(gg.bottom.ll) >= i+1 {
			height--
		}
		l.vsync(x+i, y, height, rw, ggl)
	}
	return len(gg.left.ll)
}

func (gg *gaps) syncRight(
	x, y, width, height int, rw runeWriter, ggl *Globals,
) int {
	for i, l := range gg.right.ll {
		if height <= 0 || i == width {
			return i
		}
		if len(gg.top.ll) >= i+1 {
			height--
			y++
		}
		if len(gg.bottom.ll) >= i+1 {
			height--
		}
		l.vsync(x+width-(i+1), y, height, rw, ggl)
	}
	return len(gg.right.ll)
}
