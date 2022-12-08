// Copyright (c) 2022 Stephan Lukits. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package lyt

type update struct{ width, height int }

type dim struct{ x, y, width, hight int }

// Dim holds the information which is needed and set during the layout
// process.
type Dim struct {
	x, y, width, height int

	// updates holds requests for width/height updates before the
	// updating reflow and the fillsWidth/fillsHeight during the
	// updating reflow and is set to zero after the updating reflow.
	update update

	// isDirty is set if the dimensions of a Dim-instance change to find
	// easily the components which needs redrawing after a Reflow of the
	// layout
	isDirty bool

	// filler indicate if components may grow or shrink in the
	// respective dimension.  A zero-value means no growing or
	// shrinking; a positive value means can grow or shrink but cannot
	// shrink further than its value.
	fillsWidth, fillsHeight int

	// clipper indicate that a component has only a part of its needed
	// area visible (width-clipWidth/height-clipHeight) on the screen.
	// A component that can not be layouted on the screen has
	// clipWidth/clipHeight set to its width/height indicating that it
	// is not on the screen (see IsOffScreen).
	clipWidth, clipHeight int

	// margins are used to deal with left-over space which may happens
	// if an inside component and all its components to the right or
	// bottom have a fixed widths or heights and don't need the
	// available space.
	mrgTop, mrgRight, mrgBottom, mrgLeft int
}

func needsToBeNonZero(w, h int) {
	if w > 0 && h > 0 {
		return
	}
	panic("dim: constructor: dimensions/self must not be zero")
}

// DimFilling creates a dimension which can extend its width and hight
// arbitrarily but not shrink below given sizes.  The later must be
// positive.
func DimFilling(fw, fh int) *Dim {
	needsToBeNonZero(fw, fh)
	return &Dim{fillsWidth: fw, fillsHeight: fh, isDirty: true}
}

// DimFillingWidth creates a dimension which can extend its width
// arbitrarily but not shrink below given width while it has a fixed
// height.  Both arguments must be positive.
func DimFillingWidth(fw, h int) *Dim {
	needsToBeNonZero(fw, h)
	return &Dim{fillsWidth: fw, height: h, isDirty: true}
}

// DimFillingHeight creates a dimension which can extend its height
// arbitrarily but not shrink below given height while it has a fixed
// width.  Both arguments must be positive.
func DimFillingHeight(w, fh int) *Dim {
	needsToBeNonZero(w, fh)
	return &Dim{width: w, fillsHeight: fh, isDirty: true}
}

// DimFixed creates a dimension with a fixed width and height.  Both
// arguments must be positive.
func DimFixed(w, h int) *Dim {
	needsToBeNonZero(w, h)
	return &Dim{width: w, height: h, isDirty: true}
}

// X returns the x-component of a Dimer's origin.
func (d *Dim) X() int { return d.x }

// Y returns the y-component of a Dimer's origin.
func (d *Dim) Y() int { return d.y }

// Width returns a Dimer's width.
func (d *Dim) Width() int { return d.width }

// Height returns a Dimer's height.
func (d *Dim) Height() int { return d.height }

// UpdateWidth ensures that at the next Reflow of the layout the
// associated Dimer's width + delta is layouted as requested.
// Subsequent UpdateWidth calls overwrite each other if no Resize was
// requested from the associated layout Manager m in between.   The main
// use case of this feature is that it provides especially for width
// fillers the means to control their filling width.  E.g. an initial
// layout of a Chainer's Dimers might be:
//
//	+------------+------------+
//	|     d      |      d'    |
//	+------------+------------+
//
// whereas d and d' are width-fillers having the available width evenly
// distributed.  Calling now d.Dim().UpdateWidth(2) followed by
// m.Reflow() will result in
//
//	+--------------+----------+
//	|       d      |    d'    |
//	+--------------+----------+
//
// UpdateWidth calls resulting in a none positive width are ignored.
func (d *Dim) UpdateWidth(delta int) *Dim {
	if delta == 0 {
		return d // TODO: coverage
	}
	u := d.width + delta
	if u <= 0 {
		return d
	}
	d.update.width = u
	return d
}

// UpdateHeight ensures that at the next Reflow of the layout the
// associated Dimer's height + delta is layouted as requested.
// subsequent UpdateHeight calls overwrite each other if no Reflow was
// requested from the associated layout Manager m in between.   The main
// use case of this feature is that it provides especially for hight
// fillers the means to control their filling height.  E.g. an initial
// layout of a Stacker's Dimers might be the left layout:
//
//	+------------+                           +------------+
//	|            |                           |            |
//	|     d      |                           |     d      |
//	|            |                           |            |
//	|            |                           |            |
//	+------------+                           |            |
//	|            |                           |            |
//	|     d'     |                           +------------+
//	|            |                           |            |
//	|            |                           |     d'     |
//	+------------+                           +------------+
//
// whereas d and d' are height-fillers having the available height
// evenly distributed.  Calling now d.Dim().UpdateHeight(2) followed by
// m.Reflow() will result in the right layout.  UpdateHeight calls
// resulting in a none positive height are ignored.
func (d *Dim) UpdateHeight(delta int) *Dim {
	if delta == 0 {
		return d // TODO: coverage
	}
	h := d.height + delta
	if h <= 0 {
		return d
	}
	d.update.height = h
	return d
}

// SetHeight sets given hight.  A minimal filling height is removed.  A
// none negative height is ignored.
func (d *Dim) SetHeight(h int) *Dim {
	if h <= 0 || d.height == h {
		return d // TODO: coverage
	}
	d.height = h
	if d.fillsHeight != 0 {
		d.fillsHeight = 0
	}
	d.isDirty = true
	return d
}

// SetWidth sets given width.  A minimal filling width is removed.  A
// none negative width is ignored.
func (d *Dim) SetWidth(w int) *Dim {
	if w <= 0 || d.width == w {
		return d // TODO: coverage
	}
	d.width = w
	if d.fillsWidth != 0 {
		d.fillsWidth = 0
	}
	d.isDirty = true
	return d
}

// SetFilling makes a Dimer filling respectively sets a filling Dimer's
// minimum width and height.  A none positive width or height is ignored.
func (d *Dim) SetFilling(width, height int) {
	if width > 0 {
		d.fillsWidth = width
	}
	if height > 0 {
		d.fillsHeight = height
	}
}

// IsUpdated returns true after UpdateWidth or UpdateHeight was used and
// the layout has not been recalculated yet.
func (d *Dim) IsUpdated() bool {
	return d.update.height != 0 || d.update.width != 0
}

// IsDirty returns true if a layout reflow changed a Dimer's area in the
// layout.
func (d *Dim) IsDirty() bool { return d.isDirty }

// setClean flags a dirty dimer as clean.
func (d *Dim) setClean() { d.isDirty = false }

func (d *Dim) cleanUp() {
	if d.IsDirty() {
		d.setClean()
	}
	if d.IsUpdated() {
		d.clearUpdate()
	}
}

// IsFillingWidth returns true if a Dimer's width can arbitrarily grow.
func (d *Dim) IsFillingWidth() bool { return d.fillsWidth > 0 }

// IsFillingHeight returns true if a Dimer's height can arbitrarily
// grow.
func (d *Dim) IsFillingHeight() bool { return d.fillsHeight > 0 }

func (d *Dim) prepareLayout() *dim {
	x, y, width, height := d.Printable()
	before := &dim{x: x, y: y, width: width, hight: height}
	if d.IsUpdated() {
		d.fixateUpdate()
	}
	return before
}

func (d *Dim) finalizeLayout(before *dim) {
	x, y, width, height := d.Printable()
	if before.x != x || before.y != y ||
		before.width != width || before.hight != height {
		d.isDirty = true
	}
}

func (d *Dim) fixateUpdate() {
	if d.update.height != 0 {
		d.height = d.update.height
		if d.fillsHeight > 0 {
			d.update.height = d.fillsHeight
			d.fillsHeight = 0
		}
	}
	if d.update.width != 0 {
		d.width = d.update.width
		if d.fillsWidth > 0 {
			d.update.width = d.fillsWidth
			d.fillsWidth = 0
		}
	}
}

func (d *Dim) clearUpdate() {
	if d.update.height > 0 {
		if d.update.height > d.height {
			d.fillsHeight = d.height // TODO: coverage
		} else {
			d.fillsHeight = d.update.height
		}
	}
	if d.update.width > 0 {
		if d.update.width > d.width {
			d.fillsWidth = d.width // TODO: coverage
		} else {
			d.fillsWidth = d.update.width
		}
	}
}

// Printable provides a Dimer's printable area in a layout, i.e. with
// clippings and without margins.
func (d *Dim) Printable() (x, y, width, height int) {
	if d.IsOffScreen() {
		return 0, 0, 0, 0
	}
	width, height = d.width, d.height
	if d.clipWidth > 0 {
		width -= d.clipWidth
	}
	if d.clipHeight > 0 {
		height -= d.clipHeight
	}
	x, y = d.x, d.y
	if d.mrgLeft != 0 {
		x += d.mrgLeft
	}
	if d.mrgTop != 0 {
		y += d.mrgTop
	}
	return x, y, width, height
}

// Screen provides a Dimer's screen area in the layout, i.e. the
// [Dim.Printable] rectangular with margins.
func (d *Dim) Screen() (x, y, width, height int) {
	if d.IsOffScreen() {
		return 0, 0, 0, 0
	}
	mt, mr, mb, ml := d.Margin()
	if mt == 0 && mr == 0 && mb == 0 && ml == 0 {
		return d.Printable()
	}
	width = d.width
	if mr == 0 && ml == 0 {
		if d.clipWidth > 0 {
			width -= d.clipWidth
		}
	} else {
		width += ml + mr
	}
	height = d.height
	if mt == 0 && mb == 0 {
		if d.clipHeight > 0 {
			height -= d.clipHeight
		}
	} else {
		height += mt + mb
	}
	return d.x, d.y, width, height
}

// Contains returns true if the dimers screen area (see [Dim.Rect])
// contains given point with the coordinates x and y.
func (d *Dim) Contains(x, y int) bool {
	return d.contains(d.Screen, x, y)
}

// PrintableContains returns true if the dimers printable area (see
// [Dim.Area]) contains given point with the coordinates x and y.
func (d *Dim) PrintableContains(x, y int) bool {
	return d.contains(d.Printable, x, y)
}

func (d *Dim) contains(a func() (x, y, w, h int), x, y int) bool {
	ax, ay, aw, ah := a()
	if ay <= y && ax <= x && aw+ax > x && ah+ay > y {
		return true
	}
	return false
}

// Clip how much of the components (minimal) width and hight is clipped.
func (d *Dim) Clip() (width, height int) {
	return d.clipWidth, d.clipHeight
}

// Margin is used to fill a layout's empty space.  I.e. width + right +
// left is the actual consumed width of a Dimer in a calculated layout.
func (d *Dim) Margin() (top, right, bottom, left int) {
	return d.mrgTop, d.mrgRight, d.mrgBottom, d.mrgLeft
}

func (d *Dim) resetClippingsAndMargins() {
	d.clipHeight, d.clipWidth = 0, 0
	d.mrgTop, d.mrgRight, d.mrgBottom, d.mrgLeft = 0, 0, 0, 0
}

// IsOffScreen returns true if a Dimer's width and hight is fully
// clipped.
func (d *Dim) IsOffScreen() bool {
	return d.clipWidth == d.width && d.clipHeight == d.height
}

func (d *Dim) setLayedOutWidth(w, mrgRight int) {
	if d.width >= w || d.fillsWidth <= w {
		if d.mrgLeft > 0 || d.mrgRight > 0 {
			d.mrgLeft, d.mrgRight = 0, 0
		}
	}
	if d.width <= w {
		if d.clipWidth > 0 {
			d.clipWidth = 0
		}
	}
	if d.width == w {
		return
	}
	if d.fillsWidth > 0 {
		if d.fillsWidth > w {
			d.clipWidth = d.fillsWidth - w
			d.width = d.fillsWidth
		} else {
			d.width = w
		}
		return
	}
	if d.width > w {
		d.clipWidth = d.width - w
		return
	}
	// hence d.width < w
	if mrgRight > 0 {
		d.mrgRight = mrgRight
		d.mrgLeft = w - d.width - mrgRight
		return
	}
	d.mrgLeft = (w - d.width) / 2
	d.mrgRight = w - d.width - d.mrgLeft
}

func (d *Dim) setLayedOutHeight(h int, mrgBottom int) {
	if d.height >= h || d.fillsHeight <= h {
		if d.mrgTop > 0 || d.mrgBottom > 0 {
			d.mrgTop, d.mrgBottom = 0, 0
		}
	}
	if d.height <= h {
		if d.clipHeight > 0 {
			d.clipHeight = 0
		}
	}
	if d.height == h {
		return
	}
	if d.fillsHeight > 0 {
		if d.fillsHeight > h {
			d.clipHeight = d.fillsHeight - h
			d.height = d.fillsHeight
		} else {
			d.height = h
		}
		return
	}
	if d.height > h {
		d.clipHeight = d.height - h
		return
	}
	// hence d.height < h
	if mrgBottom > 0 {
		d.mrgBottom = mrgBottom
		d.mrgTop = h - d.height - mrgBottom
		return
	}
	d.mrgTop = (h - d.height) / 2
	d.mrgBottom = h - d.height - d.mrgTop
}

func (d *Dim) setOrigin(x, y int) {
	if d.x == x && d.y == y {
		return
	}
	d.x = x
	d.y = y
	d.isDirty = true
}

// layoutedHeight is the height with without its clipping or with its
// top/bottom margins added if there is/are any clipping or margins, i.e.
// the consumed height in the layout.
func (d *Dim) layoutHeight() int {
	if d.clipHeight > 0 {
		return d.height - d.clipHeight
	}
	if d.mrgTop > 0 || d.mrgBottom > 0 {
		return d.height + d.mrgTop + d.mrgBottom
	}
	return d.height
}

// layoutedWidth is the width with without its clipping or with its
// left/right margins added if there is/are any clipping or margins, i.e.
// the consumed width in the layout.
func (d *Dim) layoutWidth() int {
	if d.clipWidth > 0 {
		return d.width - d.clipWidth
	}
	if d.mrgRight > 0 || d.mrgLeft > 0 {
		return d.width + d.mrgRight + d.mrgLeft
	}
	return d.width
}

func (d *Dim) setOffScreen() {
	if d.height == 0 {
		d.height = d.fillsHeight // TODO: coverage
	}
	if d.width == 0 {
		d.width = d.fillsWidth // TODO: coverage
	}
	d.clipHeight = d.height
	d.clipWidth = d.width
}
