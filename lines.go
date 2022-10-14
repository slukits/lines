// Copyright (c) 2022 Stephan Lukits. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package lines

import (
	"github.com/slukits/lines/internal/api"
	"github.com/slukits/lines/internal/term"
)

// Eventer is the interface which all reported events implement.
type Eventer = api.Eventer

// QuitEventer is reported when a Lines-instance is quit.
type QuitEventer = api.QuitEventer

// ResizeEventer is reported when the Lines-display was resized.
type ResizeEventer = api.ResizeEventer

type Lines struct {
	// scr to report resize events to screen components.
	scr *screen

	// backend is needed to post events.
	backend api.EventProcessor
}

// Term returns a Lines up with a terminal backend displaying and reporting
// events to given componenter.  Given componenter has the Quitable
// feature set to 'q', ctrl-c and ctrl-d.  The binding to 'q' may be
// removed.  The bindings to ctrl-c and ctrl-d may not be removed.  Use
// the TermKiosk constructor for an setup without any quit bindings.
// Term panics if the terminal screen can't be obtained.
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
	initialize(
		Componenter, interface{ Post(Eventer) error },
	) layoutComponenter

	// isInitialized returns true if embedded *component was wrapped
	// into a layout component.
	isInitialized() bool

	// embedded returns a reference to client-component's embedded
	// Component-instance.
	embedded() *Component

	// backend to post Update and Focus events on a user Componenter
	// implementation.
	backend() interface{ Post(Eventer) error }
}

// TermKiosk returns an Events instance without registered Quitable feature,
// i.e. the application can't be quit by the user.
func TermKiosk(cmp Componenter) {
	defaultFeatures = &features{
		keys: map[Modifier]map[Key]FeatureMask{},
		runes: map[Modifier]map[rune]FeatureMask{ZeroModifier: {
			0: NoFeature, // indicates the immutable default features
		}},
		buttons: map[Modifier]map[Button]FeatureMask{},
	}
	Term(cmp)
}

// Quit quits given lines instance's backend and unblocks WaitForQuit.
func (ee *Lines) Quit() { ee.backend.Quit() }

// WaitForQuit blocks until given Lines-instance is quit.
func (ee *Lines) WaitForQuit() { ee.backend.WaitForQuit() }

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
