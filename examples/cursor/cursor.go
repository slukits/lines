// Copyright (c) 2022 Stephan Lukits. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"

	"github.com/slukits/lines"
	"github.com/slukits/lines/examples/demo"
)

func main() { lines.Term(&App{}).WaitForQuit() }

// App is the root component of the cursor example which is passed to
// the lines.Term constructor above.  By embedding lines.Component
// App satisfies the Componenter interface.  By embedding lines.Stacking
// further components may be nested into App in a stacking fashion.
// Stacking adds a CC-slice of type Componenter to the App component for
// the stacked components.  Finally demo.Demo is a helper type
// defined in examples/demo to frame components and implement movement
// between components using the tab-key.
type App struct {
	// lines.Component makes App satisfy the lines.Componenter interface
	lines.Component
	// lines.Stacking makes  App satisfy the lines.Stacker interface and
	// enables App to nest components (in a stacking fashion)
	lines.Stacking
	demo.Demo
}

// appTitle is of type runes because the title needs to be printed at a
// specific position in the frame.  If we want to write at a specific
// position in lines we always need to provide runes and use
// line.Printer.  If we just want to print at a specific line we can use
// fmt.Fprint and provide a string.  But a line internally is always
// represented as slice of runes.
var appTitle []rune = []rune("cursor demo")

// OnInit sets up the components structure of the cursor-demo.
func (c *App) OnInit(e *lines.Env) {

	c.InitDemo(c, e, appTitle)

	// set up the nested components and how they relate to each other.
	c.CC = []lines.Componenter{&row{}, &row{}}
	arrows, clicks, feature := &arrowsDemo{}, &clickDemo{}, &featureDemo{}
	c.CC[0].(*row).CC = []lines.Componenter{arrows, clicks}
	c.CC[1].(*row).CC = []lines.Componenter{feature}

	// write some instructions in the second gap line (to not overwrite
	// the frame)
	fmt.Fprint(c.Gaps(1).Top, "use the mouse or tab-key to move to a demo")

	// make the tab-key select the "first" demo
	c.Next, arrows.Next, clicks.Next, feature.Next = arrows, clicks,
		feature, c

	// make demos focusable by mouse-click
	c.FF.Set(lines.Focusable)

	// have the cursor demo horizontally and vertically centered on
	// bigger screens.
	c.Dim().SetWidth(72).SetHeight(24)
}

func (c *App) OnFocusLost(e *lines.Env) {}

// row is a structuring component for App which stacks rows which in
// turn chain components, i.e. we can build in this way a n x m layout.
type row struct {
	lines.Component
	lines.Chaining
}
