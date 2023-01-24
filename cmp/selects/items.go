// Copyright (c) 2022 Stephan Lukits. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package selects

import (
	"fmt"

	"github.com/slukits/lines"
)

// items represents a selection's items-items layer.
type items struct {
	component

	// ii are the items-labels
	ii []string

	// styler returns for each item its corresponding styler iff set
	styler Styler

	// listener is updated with selected item
	listener func(int)

	maxWidth int

	defaultItem int
}

// OnInit prints the list items to the printable area.
func (ii *items) OnInit(e *lines.Env) {
	ii.Dim().SetWidth(ii.width(true))
	ii.resetItemsLabel(e)
}

func (ii *items) resetItemsLabel(e *lines.Env) {
	if !ii.hasDefault() {
		fmt.Fprint(e, lines.Filler+Drop)
		return
	}
	fmt.Fprint(e, ii.ii[ii.defaultItem]+lines.Filler+Drop)
}

func (ii *items) hasDefault() bool {
	return ii.defaultItem >= 0 && ii.defaultItem < len(ii.ii)
}

func (c *items) width(respectMax bool) int {
	maxWdth, decoration := 0, len([]rune(Drop))+1
	for _, i := range c.ii {
		if maxWdth >= len(i)+decoration {
			continue
		}
		maxWdth = len(i) + decoration
	}
	if !respectMax {
		return maxWdth
	}
	if c.maxWidth == 0 || c.maxWidth+decoration >= maxWdth {
		return maxWdth
	}
	return c.maxWidth + decoration
}

func (c *items) OnClick(e *lines.Env, x, y int) {
	l := &ModalList{
		List: List{
			Items:    c.ii,
			close:    c.close,
			Listener: c.listener,
			Styler:   c.styler,
		},
	}
	l.pos = lines.NewLayerPos(
		c.Dim().X(), c.Dim().Y()+1,
		c.width(false), len(c.ii),
	)
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
	if data.(int) == -1 {
		c.resetItemsLabel(e)
		return
	}
	fmt.Fprint(e, c.ii[data.(int)]+lines.Filler+"▼")
}
