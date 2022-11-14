// Copyright (c) 2022 Stephan Lukits. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package lyt

import "fmt"

var ErrDim = fmt.Errorf("%w"+"dimensions: ", ErrLyt)

func layoutStacker(s Stacker) (err error) {
	minHeight, filler, n, err := minStackHeight(s)
	if err != nil {
		return err
	}
	_, _, stackerWidth, stackerHeight := area(s)
	if stackerWidth <= 0 || stackerHeight <= 0 {
		layoutStackedOffScreen(s)
		return nil
	}
	if minHeight > stackerHeight {
		layoutStackerOverflowing(s)
		return nil
	}
	if filler > 0 {
		layoutFilledStacker(s, minHeight, filler)
		return nil
	}
	layoutFixedStackerUnderflowing(s, minHeight, n)
	return err
}

func layoutStackedOffScreen(s Stacker) {
	s.ForStacked(func(d Dimer) (stop bool) {
		d.Dim().setOffScreen()
		return false
	})
}

// layoutFilledStacker not overflowing distributes remaining height
// equally equally amongst fillers.
func layoutFilledStacker(s Stacker, minHeight, filler int) {
	x, y, stackerWidth, stackerHeight := area(s)
	distribute := (stackerHeight - minHeight) / filler
	distributeModulo := (stackerHeight - minHeight) % filler
	shiftY := 0
	s.ForStacked(func(d Dimer) (stop bool) {
		d.Dim().setOrigin(x, y+shiftY)
		d.Dim().setLayoutedWidth(stackerWidth, 0)
		if d.Dim().fillsHeight == 0 { // fixed height Dimer
			d.Dim().setLayoutedHeight(d.Dim().height, 0)
			shiftY += d.Dim().layoutHeight()
			return false
		}
		fillerHeight := distribute + d.Dim().fillsHeight
		if distributeModulo > 0 {
			distributeModulo--
			fillerHeight++
		}
		d.Dim().setLayoutedHeight(fillerHeight, 0)
		shiftY += d.Dim().layoutHeight()
		return false
	})
}

// layoutFixedStackerUnderflowing distributes remaining hight equally as
// margins over the fixed height Dimers.
func layoutFixedStackerUnderflowing(s Stacker, minHeight, n int) {
	x, y, stackerWidth, stackerHeight := area(s)
	mm := calculateMargins(stackerHeight-minHeight, n)
	shiftY, i := 0, 0
	s.ForStacked(func(d Dimer) (stop bool) {
		d.Dim().setOrigin(x, y+shiftY)
		d.Dim().setLayoutedWidth(stackerWidth, 0)
		d.Dim().setLayoutedHeight(d.Dim().height+mm.sum(i), mm.bottom(i))
		shiftY += d.Dim().layoutHeight()
		i++
		return false
	})
}

type margins [][2]int

func (m margins) sum(idx int) int { return m[idx][0] + m[idx][1] }

func (m margins) right(idx int) int  { return m[idx][1] }
func (m margins) bottom(idx int) int { return m[idx][1] }

// calculateMargins for a Stacker/Chainer with only fixed components in
// a way that the space for margins is evenly distributed around and
// between the components, e.g. 3 components with 20 units available
// space would result in the margins (5,2), (3,2), (3,5), i.e. we have
// around and between the components 5 units margin.
func calculateMargins(l, n int) margins {
	if n == 1 {
		return [][2]int{{l / 2, l - l/2}}
	}
	distribute := l / (n + 1) // one margin more than components
	modulo := l % (n + 1)
	mm := make([][2]int, n)
	for i := 0; i < n; i++ {
		if i == 0 {
			mm[i][0] = distribute
			mm[i][1] = distribute / 2
			continue
		}
		if i < n-1 {
			mm[i][0] = distribute - distribute/2
			mm[i][1] = distribute / 2
		}
		if i+1 == n {
			mm[i][0] = distribute - distribute/2
			mm[i][1] = distribute
		}
		if i+modulo == n {
			mm[i][1] += 1
			modulo--
		}
	}
	return mm
}

// layoutStackerOverflowing set fillers to their minimum height and
// clips/puts off-screen overflowing Dimers.
func layoutStackerOverflowing(s Stacker) {
	x, y, stackerWidth, stackerHeight := area(s)
	shiftY := 0
	s.ForStacked(func(d Dimer) (stop bool) {
		if shiftY >= stackerHeight {
			d.Dim().setOffScreen()
			return false
		}
		d.Dim().setOrigin(x, y+shiftY)
		d.Dim().setLayoutedWidth(stackerWidth, 0)
		if d.Dim().fillsHeight == 0 {
			if stackerHeight-shiftY-d.Dim().height < 0 { // overflow?
				d.Dim().setLayoutedHeight(stackerHeight-shiftY, 0)
				shiftY += d.Dim().layoutHeight()
				return false
			}
			d.Dim().setLayoutedHeight(d.Dim().height, 0)
			shiftY += d.Dim().layoutHeight()
			return false
		}
		if stackerHeight-shiftY-d.Dim().fillsHeight < 0 { // overflow?
			d.Dim().setLayoutedHeight(stackerHeight-shiftY, 0)
			shiftY += d.Dim().layoutHeight()
			return false
		}
		d.Dim().setLayoutedHeight(d.Dim().fillsHeight, 0)
		shiftY += d.Dim().layoutHeight()
		return false
	})
}

// minSackHeight calculates the minimal stack height to layout a stack
// as side effect it also counts the filler and Dimer in a stack and
// checks for consistency of the provided dimensions, i.e. height and
// fillsHeight can not both be zero.
func minStackHeight(s Stacker) (minHeight, filler, n int, err error) {
	s.ForStacked(func(d Dimer) (stop bool) {
		n++
		if d.Dim().fillsHeight == 0 {
			if d.Dim().height == 0 {
				err = fmt.Errorf("%w%s", ErrDim,
					"stack-layout: Dimer must be filling or have height")
				return true
			}
			minHeight += d.Dim().height
			return false
		}
		minHeight += d.Dim().fillsHeight
		filler++
		return false
	})
	if err != nil {
		return 0, 0, 0, err
	}
	return minHeight, filler, n, nil
}
