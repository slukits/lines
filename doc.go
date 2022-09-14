// Copyright (c) 2022 Stephan Lukits. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

/*
Package lines provides an unopinionated, well tested and documented,
terminal UI library.  It does the heavy lifting for you when it comes to

  - concurrency safety
  - event handling
  - layout handling
  - content/format handling
  - feature handling
  - testing

The motivation is to provide an UI-library with a small API and few
powerful features that lets its users quickly implement an terminal ui
exactly as needed.

	import (
	    fmt

	    "github.com/slukits/lines"
	)

type Cmp struct { lines.Component }

	func (c *Cmp) OnInit(e *lines.Env) { // Env: component environment
	    c.FF.Add(lines.Scrollable) // FF: component features
	    fmt.Fprintf(e, "%s %s", "hello", "world")
	}

	func main() { lines.New(&Cmp{}).Listen() } // blocking

New provides an Events-instance reporting user input and
programmatically posted events to listener implementations of client
provided components.  While client listener implementations print to an
provided environment which is associated with the component's portion of
the screen.

lines wraps the package https://github.com/gdamore/tcell which does the
heavy lifting on the terminal side.  I didn't see a point in making
something well done worse hence I didn't wrap the constants and types
which are defined by tcell and used for event-handling and styling.
I.e. you will have to make yourself acquainted with tcell's Key,
ModeMap, ButtonMask and AttrMask constants as well as its Style and
Color type as needed.  I also tried to take care that lines doesn't
“remove” features that tcell provides.

# Concurrency safety

What doesn't work

	func (c *Cmp) OnInit(e *lines.Env) {
	    go func() {
	        time.Sleep(1*time.Second)
	        fmt.Fprint(e, "awoken") // will panic
	    }()
	}

what does work

	func (c *Cmp) OnInit(e *lines.Env) {
	    go func(ee *lines.Events) {
	        time.Sleep(1*time.Second)
	        ee.Update(c, nil, func(e *lines.Env) {
	             fmt.Fprint(e, "awoken") // will not panic
	        })
	    }(e.EE)
	}

Also using functionality or properties provided by embedded Component
instance after a listener has returned doesn't work.

	func (c *Cmp) OnInit(e *lines.Env) {
	    go func() {
	        time.Sleep(1*time.Second)
	        c.FF.Set(Scrollable) // panics or creates a race condition
	        c.Dim().SetWidth(42) // panics or creates a race condition
	    }()
	}

It is only save to pass (the initially created) events instance e.EE on
to a go routine where at the end provided update mechanisms of said
Events-instance are used to report back to a component.

# Event handling

The majority of lines' interfaces are for event handling.  Is such an
interface implemented in a component, corresponding events are reported
to that component.  E.g. OnKey, OnFocus, OnLayout are methods of such
interfaces.  Keyboard and mouse events are bubbling up from the
focused/clicked component through all enclosing ancestors.  The
environment instance of such a reported bubbling event may be used to
suppress bubbling: e.StopBubbling().

# Layout handling

lines comes with a layout manager which does most of the work for you.
If fine grained control is needed the embedded Component's Dim method
informs about positioning and size and also provides features to change
the later.  One can also control there if a component is filling, i.e.
uses up unused space, or if its size is fixed.  Components can be
arbitrarily nested by embedding either the Stacking or Chaining type in
a component or by implementing the Stacker or Chainer interface.  The
layout manager is not smart enough to handle a component which is both
stacking and chaining other components hence it silently ignores the
chained components.

# Content and format handling

The Env(ironment) instance passed to a event listener is associated with
the screen portion of the component the event is reported to.  Writing
to the environment prints provided content to its screen portion.  Env's
methods Fmt, Sty, BG, FG, LL, Pos give fine grained control of what is
printed where and how.  Fmt stands for formatting like centered.  Sty
refers to tcell's style attributes of type tcell.AttrMask.  Methods  FG,
BG lets you set fore- and background color.  LL lets you address a
specific line, Pos a line and a column.  Each of these methods return a
writer implementation, i.e. we can do this

	fmt.Fprint(
	    e.Fmt(lines.Centered).Sty(tcell.AttrBold).LL(5),
	    "a centered bold line",
	)

The above prints "a centered bold line" centered in bold letters into the
component's fifth line.  While the e-methods bind formatting and styles
to the next printed text there is a similar API on component level
provided by the embedded Component instance: Component.Fmt, .BG, .FG
sets formatting directives for each printed content of a component.
There is also the property Component.GG which makes optional gaps around
a component accessible.  And the method Component.Mod controls if a
component's content is overwritten or appended, or if it is shown
tailed.

# Feature handling

The feature concept answers the question after the default behavior of
an ui-component.  While we probably expect that we can scroll a
component whose content doesn't fit in its screen area, do we also want
a component whose content is shown tailed to be able to scroll up and
down? Maybe, maybe not.

Lets assume we have implemented the components App, MessageBar,
Statusbar, Workspace and Panel.  Lets further assume component App
stacks the components MessageBar, Workspace and Statusbar while a
Workspace  chains two panel instances p1 and p2.

	APP--------------------------+
	  |           mb             |
	  WS-----------+-------------+
	  |            |             |
	  |    p1      |      p2     |
	  |            |             |
	  +------------+-------------+
	  |           sb             |
	  +--------------------------+

If we now click into p2 we probably expect p2 to receive the focus.  But
if we click into the statusbar will we also want sb to receive the
focus? Maybe, maybe not.

For lines all UI-components have the same semantics, i.e. for lines
there is no difference between a working space's panel and a statusbar.
Further more you were promised an unopinionated terminal ui-library.
Hence the feature concept was implemented providing convenient and fine
grained control over component behavior:

	type MyComponent { lines.Component }

	func (c *MyComponent) OnInit(_ *Env) {
	    c.FF.Add(lines.Scrollable)
	}

Now our component will react on page up/down key-presses if its content
doesn't fit (vertically) into c's screen area.  But we don't need to set
each feature for each component separately:

	func (ws *Workspace) OnInit(e *Env) {
	    ws.FF.AddRecursively(lines.Scrollable | lines.Focusable)
	}

The above line makes all descendants of the workspace scrolling and
focusable, i.e. they scroll on page-up/-down key-presses and they
receive the focus (event) if clicked.

# Testing

lines comes with testing facilities:

	import (
	    "testing"

	    "github.com/slukits/lines"
	)


	type CmpFixture struct {
	    lines.Component
	    exp string
	}


	func (c *CmpFixture) OnInit(e *lines.Env) {
	    fmt.Fprint(e, c.exp)
	}

	func TestComponentInitialization(t *T) {
	    fx := &CmpFixture{ exp: "init-reported" }
	    ee, tt := lines.Test(t, fx, 1)
	    ee.Listen()
	    if fx.exp != tt.LastScreen.String() {
	        t.Errorf("expected: '%s'; got '%s'",
	            fx.exp, tt.LastScreen.String())
	    }
	}

lines can be asked for a slightly modified Events instance augmented
with a lines.Testing instance which provides some convenience for
testing.  A testing Events instance's Listen method is not blocking and
all Methods which post an event are guaranteed to return after that
event and potentially subsequently triggered events were processed and
the screen was synchronized.

The main features of an Testing-instance are:

  - an event-countdown which automatically terminates the event loop,
  - providing methods for firing user input events
  - providing the simulated terminal screen's content and its styles.

All methods posting (user) events are guaranteed to return not before
the event and potentially subsequently triggered events were processed
and the screen was synchronized.

# TODO

All examples and scenarios mentioned in this overview are implemented
the and the package is used in a production environment.  No breaking
changes will be made to reach version 1.  But there are still some
features lacking which I'd like lines to have in order to be in some
sense complete.  You can learn about these features by parsing the code
base for “// TODO: implement”.

Enjoy!
*/
package lines
