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
//	    go func(ee *lines.Events) {
//	        time.Sleep(1*time.Second)
//	        ee.UpdateComponent(c, nil, func(e *lines.Env) {
//	             fmt.Fprint(e, "awoken") // will not panic
//	        })
//	    }(e.EE)
//	}
//
// An Env instance also informs about the triggering event see
// Evt-property.  NOTE the Evt-property can be nil, e.g. if a OnFocus
// event is reported.  Last but not least an Env instance provides
// features to communicate back to the reporting Events instance, e.g.
//
//	func (c *MyUIComponent) Runes(e *lines.Env, r rune) {
//	    fmt.Fprintf(e, "received rune: '%c'", r)
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
	write(lines []byte, at, cell int, ff LineFlags, sty Style) (int, error)
}

// Write writes to the screen area of the component having given
// environment.
func (e *Env) Write(bb []byte) (int, error) {
	return e.cmp.(cmpWriter).write(bb, -1, -1, 0,
		e.cmp.embedded().fmt.sty)
}

// Fmt sets the next write's formattings like centered.
// func (e *Env) Fmt(f FmtMask) *FmtWriter {
// 	return &FmtWriter{cmp: e.cmp.(cmpWriter), fmt: &llFmt{mask: f}}
// }

// Attr sets the next write's style attributes like bold.
func (e *Env) Attr(aa StyleAttribute) *FmtWriter {
	return &FmtWriter{cmp: e.cmp.(cmpWriter),
		sty: e.cmp.embedded().fmt.sty.WithAttrsAdded(aa)}
}

func (e *Env) AddStyleRange(idx int, sr SR, rr ...SR) {
	ll := *e.cmp.embedded().ll
	if idx < 0 || idx > len(ll) {
		return
	}
	ll[idx].addStyleRange(sr, rr...)
}

func (e *Env) SetLineFlags(idx int, ff LineFlags) {
	ll := *e.cmp.embedded().ll
	if idx < 0 || idx > len(ll) {
		return
	}
	ll[idx].setFlags(ff)
}

// FG sets the next write's foreground color.
func (e *Env) FG(color Color) *FmtWriter {
	return &FmtWriter{cmp: e.cmp.(cmpWriter),
		sty: e.cmp.embedded().fmt.sty.WithFG(color)}
}

// BG sets the next write's foreground color.
func (e *Env) BG(color Color) *FmtWriter {
	return &FmtWriter{cmp: e.cmp.(cmpWriter),
		sty: e.cmp.embedded().fmt.sty.WithBG(color)}
}

// LL returns a writer which writes to the line and its following lines
// at given index.
func (e *Env) LL(idx int, ff ...LineFlags) *locWriter {
	_ff := LineFlags(0)
	for _, f := range ff {
		_ff |= f
	}
	return &locWriter{line: idx, cell: -1, ff: _ff, cmp: e.cmp.(cmpWriter)}
}

// At returns a writer which writes to given line at given position
// adding given line flags to the line's flags.
func (e *Env) At(line, cell int, ff ...LineFlags) *locWriter {
	_ff := LineFlags(0)
	for _, f := range ff {
		_ff |= f
	}
	return &locWriter{line: line, cell: cell, ff: _ff, cmp: e.cmp.(cmpWriter)}
}

// Focused returns the currently focused component.  Please remember to
// ask your Events-instance (e.EE) for an update event of the focused
// component if you want it to be changed.
func (e *Env) Focused() Componenter {
	return e.Lines.scr.focus.userComponent()
}

// StopBubbling prevents any further reporting of an mouse or key event
// after the listener calling StopBubbling returns.
func (e *Env) StopBubbling() { e.flags |= envStopBubbling }

// ScreenSize provides the currently available screen size.  This might
// be useful during the OnInit event to do some layout
// calculations/settings before the layout manager layouts the
// components.  Or to investigate how a component's layout relates to
// the screen.
func (e *Env) ScreenSize() (width, height int) { return e.size() }

func (e *Env) reset() {
	e.Lines = nil
	e.Evt = nil
	e.cmp = nil
}
