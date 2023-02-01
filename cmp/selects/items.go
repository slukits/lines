// Copyright (c) 2022 Stephan Lukits. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package selects

import (
	"fmt"

	"github.com/slukits/lines"
)

type defaulter interface {
	dfltItems() []string
	dfltStyler() Styler
	dfltHighlighter() Highlighter
	dfltItem() int
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
	if !ii.hasDefault() {
		fmt.Fprint(e, lines.Filler+string(ii.dd.dfltOrientation()))
		return
	}
	fmt.Fprint(e, ii.dd.dfltItems()[ii.dd.dfltItem()]+lines.Filler+
		string(ii.dd.dfltOrientation()))
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
	if c.dd.dfltOrientation() == Drop {
		l.pos = lines.NewLayerPos(
			c.Dim().X(), c.Dim().Y()+1,
			c.width(false), maxHeight,
		)
	}
	if c.dd.dfltOrientation() == Up {
		l.pos = lines.NewLayerPos(
			c.Dim().X(), c.Dim().Y()-maxHeight,
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

func (c *items) OnUpdate(e *lines.Env, data interface{}) {
	c.RemoveLayer(e)
	if data.(int) == -1 {
		c.resetItemsLabel(e)
		return
	}
	fmt.Fprint(e,
		c.dd.dfltItems()[data.(int)]+lines.Filler+
			string(c.dd.dfltOrientation()))
}
