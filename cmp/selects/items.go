// Copyright (c) 2022 Stephan Lukits. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package selects

import (
	"fmt"

	"github.com/slukits/lines"
)

// LabelStyle is used to query from a Styler the style for a items'
// component label.
const LabelStyle = -1

type defaulter interface {
	dfltItems() []string
	dfltStyler() Styler
	dfltHighlighter() Highlighter
	dfltItem() int
	dfltZeroLabel() string
	dfltMaxWidth() int
	dfltMinWidth() int
	dfltMaxHeight() int
	dfltOrientation() Orientation
}

// items represents a selection's items-items layer.
type items struct {
	component

	dd defaulter

	// listener is updated with selected item
	listener lines.Componenter
}

// OnInit prints the list items to the printable area.
func (ii *items) OnInit(e *lines.Env) {
	ii.Dim().SetWidth(ii.width(true))
	ii.resetItemsLabel(e)
}

func (ii *items) resetItemsLabel(e *lines.Env) {
	if ii.dd.dfltStyler() != nil {
		ii.resetStyledItemsLabel(e)
		return
	}
	if !ii.hasDefault() {
		fmt.Fprint(e, ii.dd.dfltZeroLabel()+
			lines.Filler+string(ii.dd.dfltOrientation()))
		return
	}
	lbl := ii.calculateLabel([]rune(ii.dd.dfltItems()[ii.dd.dfltItem()]))
	fmt.Fprint(e, lbl+lines.Filler+string(ii.dd.dfltOrientation()))
}

func (ii *items) resetStyledItemsLabel(e *lines.Env) {
	sty := ii.dd.dfltStyler()(LabelStyle)
	if !ii.hasDefault() {
		fmt.Fprint(e.Sty(sty), ii.dd.dfltZeroLabel()+
			lines.Filler+string(ii.dd.dfltOrientation()))
		return
	}
	lbl := ii.calculateLabel([]rune(ii.dd.dfltItems()[ii.dd.dfltItem()]))
	fmt.Fprint(e.Sty(sty),
		lbl+lines.Filler+string(ii.dd.dfltOrientation()))
}

func (ii *items) hasDefault() bool {
	i := ii.dd.dfltItem()
	return i >= 0 && i < len(ii.dd.dfltItems())
}

func (c *items) width(respectMax bool) int {
	maxWdth, decoration := 0, len([]rune(Drop))+1
	for _, i := range c.dd.dfltItems() {
		if maxWdth >= len(i)+decoration {
			continue
		}
		maxWdth = len(i) + decoration
	}
	if len([]rune(c.dd.dfltZeroLabel()))+decoration > maxWdth {
		maxWdth = len([]rune(c.dd.dfltZeroLabel())) + decoration
	}
	if !respectMax {
		if c.dd.dfltMinWidth() > maxWdth {
			return c.dd.dfltMinWidth()
		}
		return maxWdth
	}
	dfltMaxWidth := c.dd.dfltMaxWidth()
	if dfltMaxWidth == 0 || dfltMaxWidth+decoration >= maxWdth {
		if c.dd.dfltMinWidth() > maxWdth {
			return c.dd.dfltMinWidth()
		}
		return maxWdth
	}
	if c.dd.dfltMinWidth() > dfltMaxWidth+decoration {
		return c.dd.dfltMinWidth()
	}
	return dfltMaxWidth + decoration
}

func (c *items) OnClick(e *lines.Env, x, y int) {
	if len(c.dd.dfltItems()) == 1 && c.dd.dfltItems()[0] == NoItems {
		return
	}
	l := &ModalList{
		List: List{
			Items:       c.dd.dfltItems(),
			Listener:    c.listener,
			Styler:      c.dd.dfltStyler(),
			Highlighter: c.dd.dfltHighlighter(),
		},
		close: c.close,
	}
	maxHeight := c.dd.dfltMaxHeight()
	if maxHeight == 0 || maxHeight > len(c.dd.dfltItems()) {
		maxHeight = len(c.dd.dfltItems())
	}
	px, py, _, _ := c.Dim().Printable()
	if c.dd.dfltOrientation() == Drop {
		l.pos = lines.NewLayerPos(
			px, py+1,
			c.width(false), maxHeight,
		)
	}
	if c.dd.dfltOrientation() == Up {
		l.pos = lines.NewLayerPos(
			px, py-maxHeight,
			c.width(false), maxHeight,
		)
	}
	if !c.hasDefault() {
		c.resetItemsLabel(e)
	}
	c.Layered(e, l, l.pos)
}

// close closes the menu m's modal layer of menu items.
func (c *items) close(ll *lines.Lines) {
	ll.Update(c, nil, func(e *lines.Env) {
		c.RemoveLayer(e)
	})
}

func (ii *items) OnUpdate(e *lines.Env, data interface{}) {
	ii.RemoveLayer(e)
	idx := int(data.(Value))
	if idx == -1 {
		ii.resetItemsLabel(e)
		return
	}
	lbl := ii.calculateLabel([]rune(ii.dd.dfltItems()[idx]))
	if ii.dd.dfltStyler() != nil {
		fmt.Fprint(
			e.Sty(ii.dd.dfltStyler()(idx)),
			lbl+lines.Filler+string(ii.dd.dfltOrientation()),
		)
		return
	}
	fmt.Fprint(e, lbl+lines.Filler+string(ii.dd.dfltOrientation()))
}

// calculateLabel returns given label possibly shortened if it doesn't
// fit together with the orientation into the printable width
func (ii *items) calculateLabel(lbl []rune) string {
	orientationWidth := len([]rune(ii.dd.dfltOrientation()))
	lblWidth := len(lbl)
	_, _, pw, _ := ii.Dim().Printable()
	maxLblWidth := pw - 1 - orientationWidth
	if maxLblWidth < lblWidth {
		return string(append(lbl[:maxLblWidth-1], '…'))
	}
	return string(lbl)
}
