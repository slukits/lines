// Copyright (c) 2022 Stephan Lukits. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

import (
	"github.com/slukits/lines"
	"github.com/slukits/lines/cmp/fx"
	"github.com/slukits/lines/cmp/selects"
)

var (
	even   = lines.NewStyle(lines.ZeroStyle, lines.LightGreen, lines.DarkRed)
	hiEven = lines.NewStyle(lines.Bold, lines.DarkRed, lines.LightGreen)
	odd    = lines.NewStyle(lines.ZeroStyle, lines.Yellow, lines.DarkBlue)
	hiOdd  = lines.NewStyle(lines.Bold, lines.DarkBlue, lines.Yellow)
)

type menuBar struct {
	lines.Component
	lines.Stacking
}

// OnInit sets up the right hand menu bar consisting of two List
// instances and a blank filler component separating the two.  Also the
// min-width functionality for the bottom right quit-List is set.
func (c *menuBar) OnInit(e *lines.Env) {
	m, q := &menu{}, &quit{}
	q.minWidth = func(m *menu) func() int {
		return func() int {
			return m.width
		}
	}(m)
	c.CC = append(c.CC, m, &filler{}, q)
	e.Lines.Focus(c.CC[0])
}

func (c *menuBar) menu() *menu    { return c.CC[0].(*menu) }
func (c *menuBar) quitter() *quit { return c.CC[2].(*quit) }

// OnAfterLayout adapts the width of the right hand menu bar to its
// stacked components.
func (c *menuBar) OnAfterLayout(e *lines.Env, d lines.DD) (reflow bool) {
	if c.Dim().Width() != d(c.CC[0]).Width() {
		reflow = true
		c.Dim().SetWidth(d(c.CC[0]).Width())
	}
	return reflow
}

// menu embeds a selection List whose items are the available
// List-examples.
type menu struct {
	selects.List
	// width is set after OnLayout to be available for the quit-List
	// whose layout is calculated after the menu was calculated.
	width   int
	display lines.Componenter
}

// OnInit sets up the vertically gapped selection List of the right hand
// menu bar to choose from the available List-examples.
func (c *menu) OnInit(e *lines.Env) {
	c.Items = []string{
		lines.Filler + "empty List" + lines.Filler,
		lines.Filler + "simple  List" + lines.Filler,
		lines.Filler + "scrolling List" + lines.Filler,
		lines.Filler + "empty DropDown" + lines.Filler,
		lines.Filler + "drop-down List" + lines.Filler,
		lines.Filler + "drop-up List" + lines.Filler,
		lines.Filler + "unlabeled DropDown" + lines.Filler,
		lines.Filler + "style  picker" + lines.Filler,
	}
	c.Listener = c
	c.Dim().SetHeight(len(c.Items))
	lines.Print(c.Gaps(0).Vertical.At(0).Filling(), ' ')
	c.List.OnInit(e)
}

func (c *menu) OnUpdate(e *lines.Env, data interface{}) {
	exp := &example{}
	switch int(data.(selects.Value)) {
	case 0:
		exp.explain = []string{
			"The zero-list 'List{}' is",
			"usable holding the default",
			"zero-element 'no items'",
		}
		exp.cmp = &selects.List{}
		e.Lines.Update(c.display, exp, nil)
	case 1:
		exp.explain = []string{
			"A simple list, i.e. one hav-",
			"ing space for its items,",
			"reduces its size to its items",
			"which may be selected by",
			"keyboard or mouse.",
		}
		exp.cmp = &simple{}
		exp.dontFill = true
		e.Lines.Update(c.display, exp, nil)
	case 2:
		exp.explain = []string{
			"A list having not enough",
			"space for its items, shows",
			"by default a scrollbar and",
			"is scrollable. Left-click on",
			"the scrollbar scrolls down",
			"right-click up.",
		}
		exp.cmp = &scrolling{}
		exp.dontFill = true
		e.Lines.Update(c.display, exp, nil)
	case 3:
		exp.explain = []string{
			"An empty drop-down List",
			"without label gets the de-",
			"fault label, has the default",
			"item and nothing to drop.",
		}
		exp.cmp = &selects.DropDownHrz{}
		e.Lines.Update(c.display, exp, nil)
	case 4:
		exp.explain = []string{
			"A drop-down List shows its",
			"values on click.",
		}
		exp.cmp = &selects.DropDownVrt{
			Items:       fx.NStrings(20),
			Label:       "DropDown",
			MaxHeight:   5,
			DefaultItem: selects.NoDefault,
		}
		exp.expNotFilling = true
		e.Lines.Update(c.display, exp, nil)
	case 5:
		exp.explain = []string{
			"A drop List at the bottom",
			"may show its items upwards.",
		}
		sty := func(idx int) lines.Style {
			if idx%2 == 0 {
				return even
			}
			return odd
		}
		hi := func(sty lines.Style) lines.Style {
			if sty == even {
				return hiEven
			}
			return hiOdd
		}
		exp.cmp = &selects.DropDownHrz{
			Items:       fx.NStrings(20),
			Label:       "Drop-Up",
			MaxHeight:   5,
			Orientation: selects.Up,
			Styler:      sty,
			Highlighter: hi,
		}
		exp.dontFill = true
		e.Lines.Update(c.display, exp, nil)
	case 6:
		exp.explain = []string{
			"A drop-down list without label whose",
			"item-labels greater 10 don't fit into",
			"the available width.",
		}
		exp.cmp = &selects.DropDown{
			Items:     fx.NStrings(20),
			MaxHeight: 8,
			MaxWidth:  3,
		}
		e.Lines.Update(c.display, exp, nil)
		exp.expNotFilling = true
	case 7:
		exp.explain = []string{
			"Two combined drop-downs letting a user",
			"select a style.",
		}
		exp.cmp = &styleSelection{}
		e.Lines.Update(c.display, exp, nil)
		exp.expNotFilling = true
	default:
		e.Lines.Update(c.display, BLANK, nil)
	}
}

// OnLayout sets the width property for the subsequently layout
// calculations of the "quit"-List.
func (c *menu) OnLayout(e *lines.Env) bool {
	reflow := c.List.OnLayout(e)
	c.width = c.Dim().Width()
	return reflow
}

// filler fills the space between the examples List and the quit-"List"
type filler struct{ lines.Component }

// OnInit sets the default style and highlight style of the given filler
// c to the default styles of the List-instances.
func (c *filler) OnInit(e *lines.Env) {
	dflt := c.Globals().Style(lines.Default)
	c.Globals().SetStyle(
		lines.Default, c.Globals().Style(lines.Highlight))
	c.Globals().SetStyle(lines.Highlight, dflt)
}

type styleSelection struct {
	lines.Component
	lines.Chaining
}

func (c *styleSelection) OnInit(e *lines.Env) {
	ss := &selects.Styles{}
	c.CC = append(c.CC, &selects.StyleProperty{Styles: ss}, ss)
}

// quit provides the bottom list of our right hand menu-list which
// contains only the quit-item
type quit struct {
	selects.List
	minWidth func() int
	display  lines.Componenter
}

// OnInit sets content of the embedded List instance to a horizontally
// centered "quit"-item and registers a Listener with embedded List to
// execute the quitting of the example.
func (c *quit) OnInit(e *lines.Env) {
	c.Items = []string{
		lines.Filler + "revoke" + lines.Filler,
		lines.Filler + "redraw" + lines.Filler,
		lines.Filler + "quit" + lines.Filler,
	}
	c.Listener = c.handler(e.Lines)
	c.List.OnInit(e)
}

// OnLayout adjusts the width of embedded List which by default tightens
// width to its widest item which in this case is the rather narrow
// "quit"-item.  Set minWidth-property should provide the width of the
// "examples"-List of the stacking right hand menu-bar so the quit list
// has the same width as the examples-list.
func (c *quit) OnLayout(e *lines.Env) bool {
	if !c.List.OnLayout(e) {
		return false
	}

	if c.minWidth() > c.Dim().Width() {
		c.Dim().SetWidth(c.minWidth())
	}
	return true
}

// handler provides a closure to quit the application once the user
// clicks the quit-item.
func (c *quit) handler(ll *lines.Lines) func(int) {
	return func(i int) {
		switch i {
		case 0:
			ll.Update(c.display, BLANK, nil)
			// we need this redraw since the example place holder
			// is already initialized and has a layout.  An alternative
			// would be to create it each time new, e.g.
			// ll.Update(c.display, &blank{}, nil)
			ll.Redraw()
		case 1:
			ll.Redraw()
		case 2:
			ll.Quit()
		default:
			ll.Update(c.display, BLANK, nil)
		}
	}
}
