// Copyright (c) 2022 Stephan Lukits. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package lyt

import (
	"fmt"

	"github.com/slukits/lines/internal/api"
)

type gapper interface {
	Gaps() api.Gaps
}

func layoutChainer(c Chainer) (err error) {
	minWidth, filler, n, err := minChainWidth(c)
	if err != nil {
		return err
	}
	_, _, chainerWidth, chainerHeight := area(c)
	if chainerWidth <= 0 || chainerHeight <= 0 {
		layoutChainedOffScreen(c)
		return nil
	}
	if minWidth > chainerWidth {
		layoutChainerOverflowing(c)
		return nil
	}
	if filler > 0 {
		layoutFilledChainer(c, minWidth, filler)
		return nil
	}
	layoutFixedChainerUnderflowing(c, minWidth, n)
	return err
}

func layoutChainedOffScreen(c Chainer) {
	c.ForChained(func(d Dimer) (stop bool) {
		d.Dim().setOffScreen()
		return false
	})
}

func layoutFilledChainer(c Chainer, minWidth, filler int) {
	x, y, chainerWidth, chainerHeight := area(c)
	distribute := (chainerWidth - minWidth) / filler
	distributeModulo := (chainerWidth - minWidth) % filler
	shiftX := 0
	c.ForChained(func(d Dimer) (stop bool) {
		d.Dim().setOrigin(x+shiftX, y)
		d.Dim().setLayoutedHeight(chainerHeight, 0)
		if d.Dim().fillsWidth == 0 { // fixed width Dimer
			d.Dim().setLayoutedWidth(d.Dim().width, 0)
			shiftX += d.Dim().layoutWidth()
			return false
		}
		fillerWidth := distribute + d.Dim().fillsWidth
		if distributeModulo > 0 {
			distributeModulo--
			fillerWidth++
		}
		d.Dim().setLayoutedWidth(fillerWidth, 0)
		shiftX += d.Dim().layoutWidth()
		return false
	})
}

func layoutFixedChainerUnderflowing(c Chainer, minWidth, n int) {
	x, y, chainerWidth, chainerHeight := area(c)
	mm := calculateMargins(chainerWidth-minWidth, n)
	shiftX, i := 0, 0
	c.ForChained(func(d Dimer) (stop bool) {
		d.Dim().setOrigin(x+shiftX, y)
		d.Dim().setLayoutedHeight(chainerHeight, 0)
		d.Dim().setLayoutedWidth(d.Dim().width+mm.sum(i), mm.right(i))
		shiftX += d.Dim().layoutWidth()
		i++
		return false
	})
}

func layoutChainerOverflowing(c Chainer) {
	x, y, chainerWidth, chainerHeight := area(c)
	shiftX := 0
	c.ForChained(func(d Dimer) (stop bool) {
		if shiftX >= chainerWidth {
			d.Dim().setOffScreen()
			return false
		}
		d.Dim().setOrigin(x+shiftX, y)
		d.Dim().setLayoutedHeight(chainerHeight, 0)
		if d.Dim().fillsHeight == 0 {
			if chainerWidth-shiftX-d.Dim().width < 0 { // overflow?
				d.Dim().setLayoutedWidth(chainerWidth-shiftX, 0)
				shiftX += d.Dim().layoutWidth()
				return false
			}
			d.Dim().setLayoutedWidth(d.Dim().width, 0)
			shiftX += d.Dim().layoutWidth()
			return false
		}
		if chainerWidth-shiftX-d.Dim().fillsWidth < 0 { // overflow?
			d.Dim().setLayoutedWidth(chainerWidth-shiftX, 0)
			shiftX += d.Dim().layoutWidth()
			return false
		}
		d.Dim().setLayoutedWidth(d.Dim().fillsWidth, 0)
		shiftX += d.Dim().layoutWidth()
		return false
	})
}

func minChainWidth(c Chainer) (minWidth, filler, n int, err error) {
	c.ForChained(func(d Dimer) (stop bool) {
		n++
		if d.Dim().fillsWidth == 0 {
			if d.Dim().width == 0 {
				err = fmt.Errorf("%w%s", ErrDim,
					"chain-layout: Dimer must be filling or have width")
				return true
			}
			minWidth += d.Dim().width
			return false
		}
		minWidth += d.Dim().fillsWidth
		filler++
		return false
	})
	if err != nil {
		return 0, 0, 0, err
	}
	return minWidth, filler, n, nil
}

func area(d Dimer) (x, y, w, h int) {
	x, y, w, h = d.Dim().Area()
	if ggr, ok := d.(gapper); ok {
		gg := ggr.Gaps()
		x += gg.Left
		y += gg.Top
		w -= (gg.Left + gg.Right)
		h -= (gg.Top + gg.Bottom)
	}
	return x, y, w, h
}
