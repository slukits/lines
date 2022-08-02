# Overview

lines provides an unopinionated, well tested and documented,
terminal UI library which does the heavy lifting for you when it
comes to

* concurrency safety
* event handling
* layout handling
* content/format handling
* feature handling
* testing

The motivation is to provide an UI-library with a small API and few
powerful features that lets its users quickly implement an terminal
ui exactly as needed.

```go
    package main

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
```

New provides an Events-instance reporting user input and
programmatically posted events to listener implementations of
provided components.  While listener implementations print to an
provided environment which is associated with the portion of the
screen of their component.

lines wraps the package [tcell](https://github.com/gdamore/tcell)
which does the heavy lifting on the terminal side.  I didn't see a
point in making something well done worse hence I didn't wrap the
constants and types which are defined by tcell and are used for
event-handling and styling.  I.e. you will have to make yourself
acquainted with tcell's Key, ModeMap, ButtonMask and AttrMask
constants as well as its Style and Color type as needed.  I also
tried to take care that lines doesn't "remove" features that tcell
provides.

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
        go func(ee *lines.Events) {
            time.Sleep(1*time.Second)
            ee.Update(c, nil, func(e *lines.Env) {
                 fmt.Fprint(e, "awoken") // will not panic
            })
        }(e.EE)
    }
```

Also using functionality or properties provided by embedded Component
instance after a listener has returned doesn't work.

```go
    func (c *Cmp) OnInit(e *lines.Env) {
        go func() {
            time.Sleep(1*time.Second)
            c.FF.Add(Scrollable) // panic or race condition
            c.Dim().SetWidth(42) // panic or race condition
        }()
    }
```

It is only save to pass (the initially created) events instance e.EE
on to a go routine where at the end provided update mechanisms of
said Events-instance are used to report back to a component.

# Event handling

The majority of lines' interfaces are for event handling.  Is such an
interface implemented in a component, corresponding events are
reported to that component.  E.g. OnKey, OnFocus, OnLayout are
methods of such interfaces.

# Layout handling

lines comes with a layout manager which does most of the work for
you.  If fine grained control is needed the embedded Component's Dim
method informs about positioning and size and also provides features
to change the later.  One can also control there if a component is
filling, i.e. uses up unused space, or if its size is fixed.
Components can be arbitrarily nested by embedding either the Stacking
or Chaining type in a component or by implementing the Stacker or
Chainer interface.  The layout manager is not smart enough to handle
a component which is both stacking and chaining other components
hence it silently ignores the chained components.

# Content and format handling

The Env(ironment) instance passed to a event listener is associated with
the screen portion of the component the event is reported to.  Writing
to the environment prints provided content to its screen portion.  Env's
methods Fmt, BG, FG, LL, Pos give fine grained control of what is
printed where and how.  Fmt stands for formatting like bold or centered.
FG, BG lets you set fore- and background color.  LL lets you address a
specific line, Pos a line and a column.  Each of these methods
return the Env instance, i.e. we can do this

```go
fmt.Fprintln(e.Fmt(lines.Centered).LL(5), "a centered line")
```

The above prints "a centered line" centered into the component's fifth
line.  While e.Fmt binds the formatting to the next printed text there
is a similar API on component level provided by the embedded Component
instance: Component.Fmt, .BG, .FG sets formatting directives for each
printed content of a component.  There is also the property Component.GG
which makes optional gaps around a component accessible.  And the method
Component.Mod controls if a component's content is overwritten or
appended, or if it is shown tailed.

# Feature handling

Features of a component are accessed and controlled through the FF
property of embedded Component-type.  Features are features for the
end user of a terminal application, e.g. Scrollable.  Lets assume we
have implement the components App, MessageBar, Statusbar, Workspace
and Panel.  Lets further assume component App stacks the components
MessageBar, Workspace and Statusbar while a Workspace  chains two
panel instances p1 and p2.

    APP--------------------------+
      |           mb             |
      WS-----------+-------------+
      |            |             |
      |    p1      |      p2     |
      |            |             |
      +------------+-------------+
      |           sb             |
      +--------------------------+

Finally we have some event interfaces implemented for p1 and p2.  Now
we want to test implemented event handler and fire events for them
(provided by the Testing type).  Nothing will happen because
initially App will have the focus and that's not changing unless
lines is told to do so

```go
    func (a *App) OnInit(e *Env) {
        // ...
        e.EE.MoveFocus(p1)
    }
```

Now p1 gets its events.  We have our App started and click into p2
where we have an OnClick implementation.  (Which might tries to move
the focus to itself :)  But the click is never reported because p2
has not the feature Focusable set.  While it might seem obvious to us
that p2 should receive the focus if clicked, it is not obvious to
lines.  For lines our Statusbar and our Panel are Componenter without
any further semantics.  Now the question is: will we want also that
the statusbar gets the focus if the user clicks on it?  Maybe, maybe
not.  Hence lines doesn't try to be smart about such things and
implements the features concept instead

```go
    func (ws *Workspace) OnInit(e *Env) {
        ws.Features.AddRecursively(Focusable | Selectable)
    }
```

This one line has  the following consequences: If the user clicks on
p1 or p2 the respective component gets the focus and events about
focus gain and loss are reported, the mouse click is reported.  If
ws, p1 or p2 has the focus and the user presses the the Tab key the
"next" panel gains the focus, focus gain and loss are reported, the
Tab key is reported.  If the user presses shift-tab the "previous"
panel gains the focus ...  While there is probably enough going on to
justify one line of code more importantly this line of code
represents full control over what is going on.

# Testing

lines comes with testing facilities:

```go
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
        ee, tt := lines.Test(t, fx)
        ee.Listen()
        if tt.LastScreen != fx.exp {
            t.Errorf("expected: '%s'; got '%s'", fx.exp, tt.LastScreen)
        }
    }
```

lines can be asked for a slightly modified Events instance augmented
with a lines.Testing instance which provides some convenience for
testing.  A testing Events instance's Listen method is not blocking
and all Methods which post an event are guaranteed to return after
that event and potentially subsequently triggered events were
processed and the screen was synchronized.

The main features of an Testing-instance are:

* an event-countdown which automatically terminates the event loop,
* providing methods for firing user input events
* providing the simulated terminal screen's content.

All methods posting (user) events are guaranteed to return after the
event and potentially subsequently triggered events were processed
and the screen was synchronized.

# TODO

All examples and scenarios mentioned in this overview are implemented
the API is frozen and the package is used in a production
environment.  But there are still some features lacking which I'd
like lines to have in order to be in some sense complete.  You can
learn about these features by parsing the code base for "// TODO:
implement".

Enjoy!