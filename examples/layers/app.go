// Copyright (c) 2022 Stephan Lukits. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

import (
	"github.com/slukits/lines"
	"github.com/slukits/lines/examples/demo"
)

// App is the root component of the layers example which is passed to
// the lines.Term constructor in main.go.  By embedding lines.Component
// App satisfies the Componenter interface.  By embedding lines.Stacking
// further components may be nested into App in a stacking fashion.
// Stacking adds a CC-slice of type Componenter to the App component for
// the stacked components.  Finally demo.Demo is a helper type
// defined in examples/demo to frame components and implement movement
// between components using the tab-key.
type App struct {
	lines.Component
	lines.Stacking
	demo.Demo
}

// appTitle is of type runes because the title needs to be printed at a
// specific position in the frame.  If we want to write at a specific
// position in lines we always need to provide runes and use
// line.Printer.  If we just want to print at a specific line we can use
// fmt.Fprint and provide a string.  But a line internally is always
// represented as slice of runes.
var appTitle []rune = []rune("layers demo")

// OnInit sets up the components structure of the layers-demo.
func (c *App) OnInit(e *lines.Env) {

	c.Init(c, e, appTitle)

	// set up the nested components and how they relate to each other.
	c.CC = []lines.Componenter{&row{}, &row{}}
	menu, context, toolTip, stacked := &menuDemo{}, &context{}, &toolTip{},
		&stacked{}
	c.CC[0].(*row).CC = []lines.Componenter{menu, context}
	c.CC[1].(*row).CC = []lines.Componenter{toolTip, stacked}

	// make demos focusable with the tab-key
	c.Next, menu.Next, context.Next, toolTip.Next, stacked.Next =
		menu, context, toolTip, stacked, menu

	// have the layers demo horizontally and vertically centered on
	// bigger screens.
	c.Dim().SetWidth(72).SetHeight(24)

	// make demos focusable by mouse-click
	c.FF.Set(lines.Focusable)
}

// row is a structuring component for App which stacks rows which in
// turn chain components, i.e. we can build in this way a n x m layout.
type row struct {
	lines.Component
	lines.Chaining
}
