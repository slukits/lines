// Copyright (c) 2022 Stephan Lukits. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package lines

import (
	"fmt"
	"time"

	"github.com/slukits/lines/internal/api"
	"github.com/slukits/lines/internal/lyt"
	"github.com/slukits/lines/internal/term"
)

// Eventer is the interface which all reported events implement.
type Eventer = api.Eventer

// QuitEventer is reported when a Lines-instance is quit.
type QuitEventer = api.QuitEventer

// ResizeEventer is reported when the Lines-display was resized.
type ResizeEventer = api.ResizeEventer

// Dimer provides dimensions of a component in the layout.  Note each
// type embedding [lines.Component] type implements the Dimer interface.
type Dimer = lyt.Dimer

type Lines struct {
	// scr to report resize events to screen components.
	scr *screen

	// backend is needed to post events.
	backend api.EventProcessor
}

// Term returns a Lines instance with a terminal backend displaying and
// reporting events to given componenter and its nested components.
// Given componenter has the Quitable feature set to 'q', ctrl-c and
// ctrl-d.  The binding to 'q' may be removed.  The bindings to ctrl-c
// and ctrl-d may not be removed.  Use the TermKiosk constructor for an
// setup without any quit bindings.  Term panics if the terminal screen
// can't be obtained.
func Term(cmp Componenter) *Lines {
	ll := Lines{}
	ll.backend = term.New(ll.listen)
	ll.scr = newScreen(ll.backend.(api.UIer), cmp)
	return &ll
}

// Componenter is the private interface a type must implement to be used
// as an lines ui component.  Embedding [lines.Component] in a type
// automatically fulfills this condition:
//
//	type MyTUIComponent struct { lines.Component }
//	lines.New(&MyTUIComponent{}).Listen()
type Componenter interface {

	// enable makes the embedded component usable for the client, i.e.
	// accessing its properties and methods won't panic.
	enable()

	// disable makes the embedded component unusable for the client,
	// i.e. accessing its properties and methods is likely to panic.
	disable()

	// hasLayoutWrapper is true if a component is part of the layout and
	// its layout has been calculated by the layout manager.
	hasLayoutWrapper() bool

	// layoutComponent is a wrapper around a client-component and its
	// embedded component independent of being enabled/disabled.  It
	// combines the client-components stacking or chaining aspects and
	// the internally calculated dimensional aspects of a component.
	layoutComponent() layoutComponenter

	// initialize sets up the embedded *component instance and wraps it
	// together with the client-instance in a layoutComponenter which is
	// returned.
	initialize(Componenter, api.UIer) layoutComponenter

	// isInitialized returns true if embedded *component was wrapped
	// into a layout component.
	isInitialized() bool

	// embedded returns a reference to client-component's embedded
	// Component-instance.
	embedded() *Component

	// backend to post Update and Focus events on a user Componenter
	// implementation.
	backend() api.UIer
}

// TermKiosk returns a Lines instance without registered Quitable feature,
// i.e. the application can't be quit by the user by default.
func TermKiosk(cmp Componenter) *Lines {
	defaultFeatures = &features{
		keys: map[Modifier]map[Key]FeatureMask{},
		runes: map[Modifier]map[rune]FeatureMask{ZeroModifier: {
			0: NoFeature, // indicates the immutable default features
		}},
		buttons: map[Modifier]map[Button]FeatureMask{},
	}
	return Term(cmp)
}

// Quit quits given lines instance's backend and unblocks WaitForQuit.
func (ee *Lines) Quit() { ee.backend.Quit() }

// OnQuit registers given function to be called on quitting
// event-polling and -reporting.
func (ll *Lines) OnQuit(listener func()) { ll.backend.OnQuit(listener) }

// WaitForQuit blocks until given Lines-instance is quit.  (Except a
// Lines instance provided by a [Fixture] in which case WaitForQuit is
// not blocking.)
func (ee *Lines) WaitForQuit() { ee.backend.WaitForQuit() }

// Update posts an update event into the event queue which is reported
// either to given listener if not nil or to given componenter if given
// listener is nil.  Given data will be provided by the reported Update
// event.  Update is a no-op if componenter and listener are nil.
func (ll *Lines) Update(
	cmp Componenter, data interface{}, l Listener,
) error {
	if cmp == nil && l == nil {
		return nil
	}
	return ll.backend.Post(&UpdateEvent{
		when: time.Now(),
		cmp:  cmp,
		lst:  l,
		Data: data,
	})
}

// UpdateEvent is created by an Update call on Lines.  Its Data field
// provides the data which was passed to that Update call.
type UpdateEvent struct {
	when time.Time
	// NOTE we can not extract the componenter from the component
	// without risking a race condition hence we leave it to the
	// reporter to do so.
	cmp Componenter
	lst Listener

	// Data provided to an update event listener
	Data interface{}
}

// When of an update event is set to time.Now()
func (u *UpdateEvent) When() time.Time { return u.when }

func (u *UpdateEvent) Source() interface{} { return u }

func (l *Lines) listen(evt api.Eventer) {
	switch evt := evt.(type) {
	case ResizeEventer:
		width, height := evt.Size()
		l.scr.setWidth(width).setHeight(height)
		reportInit(l, l.scr)
		l.scr.hardSync(l)
	default:
		if quit := report(evt, l, l.scr); quit {
			l.backend.Quit()
			return
		}
		reportInit(l, l.scr)
		l.scr.softSync(l)
	}
}

// MoveFocus posts a new MoveFocus event into the event loop which once
// it is polled calls the currently focused component's OnFocusLost
// implementation while given component's OnFocus implementation is
// executed.  Finally the focus is set to given component.  MoveFocus
// fails if the event-loop is full returned error will wrap tcell's
// *PostEvent* error.  MoveFocus is an no-op if Componenter is nil.
func (c *Lines) Focus(cmp Componenter) error {
	return c.backend.Post(&moveFocusEvent{
		when: time.Now(),
		cmp:  cmp,
	})
}

// moveFocusEvent is posted by calling MoveFocus for a programmatically
// change of focus.  This event-instance is not provided to the user.
type moveFocusEvent struct {
	when time.Time
	// NOTE we can not extract the componenter from the component
	// without risking a race condition hence we leave it to the
	// reporter to do so.
	cmp Componenter
}

func (e *moveFocusEvent) When() time.Time { return e.when }

func (e *moveFocusEvent) Source() interface{} { return e }

// AtWriter is for printing runes at specific screen cells commonly used
// to define differentiated stylings.
type AtWriter interface {
	WriteAt(rr []rune)
}

// Print to an AtWriter.  The most common AtWriter of lines are provided
// by Env and Gaps instances.
func Print(w AtWriter, rr interface{}) {
	if rr == nil {
		return
	}
	switch rr := rr.(type) {
	case rune:
		w.WriteAt([]rune{rr})
	case []rune:
		w.WriteAt(rr)
	default:
		panic(fmt.Sprintf(
			"lines: print: expected rune/rune-slice; got %T", rr))
	}
}
