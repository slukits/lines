// Copyright (c) 2022 Stephan Lukits. All rights reserved.

// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package lines

type envMask uint64

const (
	envStopBubbling envMask = 1 << iota
)

// Env is an environment provided to event listeners when their event is
// reported.  An Env instance implements the io.Writer interface and is
// associated with a portion of the screen it writes to.  Writing to an
// Env instance after the event listener has returned will panic, e.g.
//
//	func (c *MyUIComponent) myListener(e *lines.Env) {
//	    go func() {
//	        time.Sleep(1*time.Second)
//	        fmt.Fprint(e, "awoken") // will panic
//	    }()
//	}
//
//	func (c *MyUIComponent) myListener(e *lines.Env) {
//	    go func(ll *lines.Lines) {
//	        time.Sleep(1*time.Second)
//	        ll.UpdateComponent(c, nil, func(e *lines.Env) {
//	             fmt.Fprint(e, "awoken") // will not panic
//	        })
//	    }(e.Lines)
//	}
//
// An Env e instance also informs about the triggering backend event see
// Evt-property which may be nil.  Last but not least an Env instance
// provides features to communicate back to the reporting [Lines]
// instance, e.g. [Env.StopBubbling]:
//
//	func (c *MyUIComponent) Runes(e *lines.Env, r rune) {
//	    fmt.Fprintf(
//	        e, "received rune: '%c' of event %T",
//	        r, e.Env.Source(),
//	    )
//	    // the event stops bubbling through enclosing components
//	    e.StopBubbling()
//	}
type Env struct {
	cmp Componenter

	size func() (int, int)

	// Lines is the Lines instance providing given environment
	// instance.  Use it to post Update or Focus events.
	Lines *Lines

	// Evt is the event which triggered the creation of the environment
	// instance.  NOTE with Evt.Source() a backend event may be accessed.
	Evt Eventer

	flags envMask
}

type cmpWriter interface {
	globals() *globals
	write(lines []byte, at, cell int, sty *Style) (int, error)
	writeAt(rr []rune, at, cell int, sty *Style)
	writeAtFilling(r rune, at, cell int, sty *Style)
}

// Write given bytes bb to the screen area of the component whose event
// handler was called with given environment e.  For that purpose given
// bytes a broken into screen lines at new-lines and optionally set
// style-attributes as well as fore- and background colors are passed as
// default style to these lines.  NOTE all previous content of the
// component is removed.
func (e *Env) Write(bb []byte) (int, error) {
	return e.cmp.(cmpWriter).write(bb, -1, -1, nil)
}

// Sty sets the next lines write's style, i.e. its style attributes and
// fore- and background color.
func (e *Env) Sty(s Style) *EnvWriter {
	return &EnvWriter{
		cmp: e.cmp.(cmpWriter),
		sty: &s,
	}
}

// AA sets the next write's style attributes like [Bold].
func (e *Env) AA(aa StyleAttributeMask) *EnvWriter {
	sty := e.cmp.embedded().gg.Style(Default).WithAdded(aa)
	return &EnvWriter{
		cmp: e.cmp.(cmpWriter),
		sty: &sty,
	}
}

// FG sets the next write's foreground color.
func (e *Env) FG(color Color) *EnvWriter {
	sty := e.cmp.embedded().gg.Style(Default).WithFG(color)
	return &EnvWriter{
		cmp: e.cmp.(cmpWriter),
		sty: &sty,
	}
}

// BG sets the next write's background color.
func (e *Env) BG(color Color) *EnvWriter {
	sty := e.cmp.embedded().gg.Style(Default).WithBG(color)
	return &EnvWriter{
		cmp: e.cmp.(cmpWriter),
		sty: &sty,
	}
}

// LL returns a writer which writes to the line and its following lines
// at given index.
func (e *Env) LL(idx int) *EnvLineWriter {
	return &EnvLineWriter{
		line: idx,
		cmp:  e.cmp.(cmpWriter),
	}
}

// Focused returns the currently focused component.  Please remember to
// ask your [Lines]-instance (e.Lines) for an update event of the focused
// component if you want it to be changed.
func (e *Env) Focused() Componenter {
	return e.Lines.scr.focus.userComponent()
}

// StopBubbling prevents any further reporting of an mouse or key event
// after the listener calling StopBubbling returns.
func (e *Env) StopBubbling() { e.flags |= envStopBubbling }

// ScreenSize provides the currently total screen size.  Use
// [Component.Dim] to get layout information about the environment
// receiving component.
func (e *Env) ScreenSize() (width, height int) { return e.size() }

func (e *Env) reset() {
	e.Lines = nil
	e.Evt = nil
	e.cmp = nil
}
