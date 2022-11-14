# Overview

lines provides an unopinionated simple UI library which does the heavy
lifting for you when it comes to

* concurrency safety
* event handling
* layout handling
* content/format handling
* feature handling
* testing

Following example times out in the go playground since it is blocking:

```go
package main

import (
    fmt

    "github.com/slukits/lines"
)

type Cmp struct { lines.Component }

func (c *Cmp) OnInit(e *lines.Env) {
    c.Dim().SetWidth(len("hello world")).SetHeight(1)
    fmt.Fprint(e, "hello world")
}

func main() {
    lines.Term(&Cmp{}).WaitForQuit()
}
```

Term provides an Lines-instance with a terminal backend.  It reports
user input and programmatically posted events to listener
implementations of client provided components.  While client listener
implementations print to an provided environment which is associated
with the component's portion of the screen.  lines is designed to 
add further backends like "shiny" or "fyne" for graphical displays.  As
of now lines has only a terminal backend which wraps the package
[tcell](https://github.com/gdamore/tcell).

Above "hello world"-program takes over a terminal screen printing
horizontally and vertically centered "hello world" to it.  "hello world"
stays centered in case the screen is a terminal window which changes its
size.  Ctrl-c, Ctrl-d or q will quit the application.  Note SetWidth in
above example works as expected because "hello world" consists of ASCII
characters only.  Is that not guaranteed you will want to count runes
instead of bytes.  Setting width and height is not necessary.  Left out
in above example "hello world" is printed to the screen starting in the
upper left corner.


# Concurrency safety

What doesn't work

```go
func (c *Cmp) OnInit(e *lines.Env) {
    go func() {
        time.Sleep(1*time.Second)
        fmt.Fprint(e, "awoken") // will panic
    }()
}
```

what does work

```go
func (c *Cmp) OnInit(e *lines.Env) {
    go func(ll *lines.Lines) {
        time.Sleep(1*time.Second)
        ll.Update(c, nil, func(e *lines.Env) {
             fmt.Fprint(e, "awoken") // will not panic
        })
    }(e.Lines)
}
```

Also using functionality or properties provided by embedded Component
instance in a function that doesn't return in the executing listener
won't work.

```go
func (c *Cmp) OnInit(e *lines.Env) {
    go func() {
        time.Sleep(1*time.Second)
        c.FF.Add(Scrollable) // panic or race condition
        c.Dim().SetWidth(42) // panic or race condition
    }()
}
```

It is only save to pass (the initially created) Lines instance 
on to a go routine where at the end provided update mechanisms of
said Lines-instance are used to report back to a component.

# Event handling

The majority of lines' interfaces are for event handling.  Is such an
interface implemented in a component, corresponding events are reported
to that component.  E.g. OnKey, OnFocus, OnLayout are methods of such
interfaces.  Keyboard and mouse events are bubbling up from the
focused/clicked component through all enclosing ancestors.  The
environment instance e of such a reported bubbling event may be used to
suppress bubbling: e.StopBubbling().

# Layout handling

lines comes with a layout manager which does most of the work for you.
If fine grained control is needed the embedded Component's Dim method
informs about positioning and size and also provides features to change
the later.  One can also control there if a component is filling, i.e.
uses up unused space, or if its size is fixed.  Components can be
arbitrarily nested by embedding either the Stacking or Chaining type in
a component or by implementing either the Stacker or Chainer interface.

# Content and format handling

The Env(ironment) instance passed to a event listener is associated with
the screen portion of the component the event is reported to.  Printing
to the environment prints provided content to its screen portion.  Env
comes with a few methods which give fine grained control over what is
printer where and how. e.g.:

```go
fmt.Fprint(e.LL(5).AA(lines.Bold),
    lines.Filler + "a centered bold line" + lines.Filler),
)
```

The above prints "a centered bold line" centered in bold letters into
the component's fifth line.  Note the line will stay centered if the
component's size changes.  If there should be many lines associated with
a component cmp without storing them in the component a component's
source can be set:

```go
cmp.Src = &lines.ComponentSource{Liner: MyLinerImplementation}
```

Whereas the liner implementation prints the lines as requested by cmp.
A similar printing API is provided by embedded Component's
Gaps(index)-method

```go
c.Gaps(0).AA(lines.Reverse)
c.Gaps(0).Corners.AA(lines.Revers)
```

above is as of now the simplest way to frame a component.  Gaps allow to
do all sorts of framing, padding and guttering of a component.  See
examples/gaps for an introduction of what can be done with gaps.  Remember
that Src and Gaps will panic if accessed outside a listener callback.

# Feature handling

The feature concept answers the question after the default behavior of
an ui-component.  Lets assume we have implemented the components App,
MessageBar, Statusbar, Workspace and Panel.  Lets further assume
component App stacks the components MessageBar, Workspace and Statusbar
while a Workspace  chains two panel instances p1 and p2.

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
focus? Maybe, maybe not.  To avoid opinionated default behavior lines
implements the "Feature" concept whose API may be accessed by the
FF-property of the embedded Component instance.

```go
type MyComponent { lines.Component }

func (c *MyComponent) OnInit(_ *lines.Env) {
    c.FF.Add(lines.Scrollable)
}
```

Now our component will react on page up/down key-presses if its content
doesn't fit (vertically) into c's screen area.  But we don't need to set
each feature for each component separately:

```go
func (ws *Workspace) OnInit(e *lines.Env) {
    ws.FF.AddRecursively(lines.Focusable)
}
```

The above line makes all descendants of the workspace focusable, i.e.
they receive the focus if clicked with the mouse on them.

An other way to acquire features without fuss is to set a components
content source which needs to have a Liner-implementation.  According to
the features of such a Liner-implementation the component's features are
set.  E.g. is a Liner a ScrollableLiner then the component gets
automatically the scrollable feature set.  examples/scrolling has
examples for features usage.

# Testing

lines comes with testing facilities:

```go
import (
    "testing"

    "github.com/slukits/lines"
)

type CmpFixture struct { lines.Component }

func TestUpdateListenerIsCalled(t *T) {
    tt := lines.TermFixture(t, 0, &CmpFixture{})
    exp :=  "update listener called"
    tt.Lines.Update(tt.Root(), nil, func(e *lines.Env) {
        fmt.Fprint(e, exp)
    })
    if exp != tt.Screen().Trimmed().String() {
        t.Errorf("expected: '%s'; got '%s'", exp,
            tt.Screen().Trimmed().String())
    }
}
```

lines can be asked for a test fixture.  The main features of an
Fixture-instance are:

* setting up and tearing down a corresponding Lines-instance whose
  WaitForQuit method is not blocking.

* providing methods for emulating user input events

* guaranteeing that no event-triggering method returns before its event
  and all subsequently triggered events have been processed and all
  writes made it to the screen.  (The second argument of TermFixture is
  a timeout within all such events need to be processed.)

* providing the screen's content and its styles.

Enjoy!
