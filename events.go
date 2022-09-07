// Copyright (c) 2022 Stephan Lukits. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package lines

import (
	"fmt"
	"sync"
	"time"

	"github.com/gdamore/tcell/v2"
)

// Events reports user-input and programmatically posted events to the
// implemented Events-interfaces of the currently focused component and
// its ancestors.  Except for the Update-event which is created
// programmatically and allows to obtain concurrency save an event
// environment for any component at any time and reports to a listener
// callback if given or to the component for which the Update event was
// requested.
type Events struct {
	mutex *sync.Mutex

	// scr to report events to screen components.
	scr *screen

	// pollEvent to run the event loop.
	pollEvent func() tcell.Event

	// postEvent to queue events requested by screen components.
	postEvent func(tcell.Event) error

	// isListening is true if we are looping around the event queue.
	isListening bool

	// reported is called back after an event has been reported; useful
	// for testing and logging.
	reported func()

	// t if non-nil events is in testing-mode.
	t *Testing

	// synced sends a message at the end of a processed event.  The
	// client (usually a Testing instance) setting up this channel is
	// responsible for draining this channel.
	synced chan bool
}

// New returns a Events instance reporting events to given components.
// It panics if the screen can't be created or initialized.  The Events
// instance has the Quitable feature bound to 'q', ctrl-c and ctrl-d.
// The binding to 'q' may be removed.  The bindings to ctrl-c and ctrl-d
// may not be removed.  Use the Kiosk constructor for an Events instance
// without any quit bindings.
func New(cmp Componenter) *Events {
	ee := Events{
		scr:    newScreen(cmp),
		mutex:  &sync.Mutex{},
		synced: make(chan bool, 1),
	}
	ee.pollEvent = ee.scr.lib.PollEvent
	ee.postEvent = ee.scr.lib.PostEvent
	return &ee
}

// Componenter is the private interface a type must implement to be used
// as an terminal ui component by lines.  Embedding [lines.Component] in
// a type automatically fulfills this condition:
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
	initialize(Componenter) layoutComponenter
	// isInitialized returns true if embedded *component was wrapped
	// into a layout component.
	isInitialized() bool
	// embedded returns a reference to client-component's embedded
	// Component-instance.
	embedded() *Component
}

// Kiosk returns an Events instance without registered Quitable feature,
// i.e. the application can't be quit by the user.
func Kiosk(cmp Componenter) *Events {
	defaultFeatures = &features{
		keys: map[tcell.ModMask]map[tcell.Key]FeatureMask{},
		runes: map[rune]FeatureMask{
			0: NoFeature, // indicates the immutable default features
		},
		buttons: map[tcell.ModMask]map[tcell.ButtonMask]FeatureMask{},
	}
	return New(cmp)
}

// IsListening returns true if receiving Events instance is looping
// around the event queue.
func (ee *Events) IsListening() bool {
	ee.mutex.Lock()
	defer ee.mutex.Unlock()
	return ee.isListening
}

// Listen blocks and starts polling from the event loop reporting
// received events to the listener interface implementation of the
// currently focused component and its ancestors.  Listen returns if
// either a quit-event was received ('q', ctrl-c, ctrl-d input) or
// QuitListening was called.  NOTE in testing Listen is non-blocking,
// i.e. returns after it has reported all init events.
func (ee *Events) Listen() {
	if !ee.setListening() {
		return
	}
	if ee.t != nil {
		ee.t.t.Helper()
		ee.t.listen()
		return
	}
	ee.listen() // TODO: not covered
}

func (ee *Events) setListening() bool {
	ee.mutex.Lock()
	defer ee.mutex.Unlock()
	if ee.isListening {
		return false
	}
	ee.isListening = true
	return true
}

func (ee *Events) listen() {
	ee.scr.lib.EnableMouse()
	for {
		ev := ee.scr.lib.PollEvent()

		switch ev := ev.(type) {
		case nil: // event-loop ended
			return
		case *tcell.EventResize:
			width, height := ev.Size()
			ee.scr.setWidth(width).setHeight(height)
			reportInit(ee, ee.scr)
			ee.scr.hardSync(ee)
		default:
			if quit := report(ev, ee, ee.scr); quit {
				ee.quitListening()
				return
			}
			reportInit(ee, ee.scr)
			ee.scr.softSync(ee)
		}
		if ee.synced != nil {
			ee.synced <- true
		}
	}
}

// Reported calls back if an event was reported for logging and testing.
// Do not use if the Events-instance was obtained from Test.
func (ee *Events) Reported(listener func()) {
	ee.reported = listener
}

// QuitListening posts a quit event ending the event-loop, i.e.  after
// this event is processed IsListening will be false.  All components
// implementing the Quitter interface are notified before lines cleans
// up its resources.  In testing also a snapshot of the last screen
// content is made.
func (ee *Events) QuitListening() {
	if ee.isListening {
		var wait func()
		if ee.t != nil {
			wait = ee.t.registerEventSync(
				"test: quit listening: sync timed out")
		}
		ee.scr.lib.PostEvent(&quitEvent{when: time.Now()})
		if wait != nil {
			ee.t.t.Helper()
			wait()
		}
		return
	}
	ee.quitListening()
}

func (ee *Events) quitListening() {
	ee.mutex.Lock()
	defer ee.mutex.Unlock()
	if !ee.isListening {
		return
	}
	ee.isListening = false
	if ee.t != nil {
		ee.t.beforeFinalize()
	}
	ee.scr.lib.Fini()
	if ee.synced != nil {
		close(ee.synced)
	}
}

type quitEvent struct {
	when time.Time
}

func (u *quitEvent) When() time.Time { return u.when }

// Update posts a new event into the event loop which calls once it is
// polled given components update listener.  Is Listener nil the given
// components Updater implementation is informed about the event.  Given
// data will be provided by the Env instance of the receiving listener:
//
//	func(c *Cmp) OnUpdate(e *lines.Env) {
//	    data := e.Evt.(*lines.UpdateEvent).Data.(*MyType)
//	}
//
// Update fails if the event-loop is full; returned error will wrap
// tcell's *PostEvent* error.  Update is an no-op if Componenter is nil.
// NOTE in testing Update starts the event loop if not running and
// returns after the event was fully processed.
func (ee *Events) Update(
	cmp Componenter, data interface{}, l Listener,
) error {
	if cmp == nil {
		return nil
	}
	if ee.t != nil && !ee.isListening {
		ee.t.listen()
	}
	evt := &UpdateEvent{
		when: time.Now(),
		cmp:  cmp,
		lst:  l,
		Data: data,
	}
	var wait func()
	if ee.t != nil {
		wait = ee.t.registerEventSync("test: update: sync timed out")
	}
	if err := ee.scr.lib.PostEvent(evt); err != nil {
		return fmt.Errorf(errEventFmt, err)
	}
	if wait != nil {
		ee.t.t.Helper()
		wait()
	}
	return nil
}

// errEventFmt is the error message for a failing update-event post.
var errEventFmt = "can't post event: %w"

// UpdateEvent is created by an Update call on Events.  Its Data field
// provides the data which was passed to that Update call.
type UpdateEvent struct {
	when time.Time
	cmp  Componenter
	lst  Listener

	// Data provided to an update event listener
	Data interface{}
}

// When of an update event is set to time.Now()
func (u *UpdateEvent) When() time.Time { return u.when }

// MoveFocus posts a new MoveFocus event into the event loop which once
// it is polled calls the currently focused component's OnFocusLost
// implementation while given component's OnFocus implementation is
// executed.  Finally the focus is set to given component.  MoveFocus
// fails if the event-loop is full returned error will wrap tcell's
// *PostEvent* error.  MoveFocus is an no-op if Componenter is nil.
// NOTE in testing MoveFocus starts the event loop if not running and
// returns after the event was fully processed.
func (ee *Events) MoveFocus(cmp Componenter) error {
	if cmp == nil {
		return nil
	}
	if ee.t != nil && !ee.isListening {
		ee.t.listen()
	}
	evt := &moveFocusEvent{
		when: time.Now(),
		cmp:  cmp,
	}
	var wait func()
	if ee.t != nil {
		wait = ee.t.registerEventSync("test: move-focus: sync timed out")
	}
	if err := ee.scr.lib.PostEvent(evt); err != nil {
		return fmt.Errorf(errEventFmt, err)
	}
	if wait != nil {
		ee.t.t.Helper()
		wait()
	}
	return nil
}

// moveFocusEvent is posted by calling MoveFocus for a programmatically
// change of focus.  This event-instance is not provided to the user.
type moveFocusEvent struct {
	when time.Time
	cmp  Componenter
}

func (u *moveFocusEvent) When() time.Time { return u.when }

// UpdateRunes posts a new UpdateRunes event into the event loop which
// once it is polled calls the given component's Runes implementation
// and registers/overwrites/removes provided runes-registrations.
// UpdateRunes fails if the event-loop is full; returned error will wrap
// tcell's *PostEvent* error.  UpdateRunes is an no-op if Componenter is
// nil.  NOTE in testing UpdateRunes starts the event loop if not
// started and returns after the event was fully processed.
func (ee *Events) UpdateRunes(cmp Componenter) error {
	if cmp == nil {
		return nil
	}
	if ee.t != nil && !ee.isListening {
		ee.t.listen()
	}
	evt := &updateRunesEvent{
		when: time.Now(),
		cmp:  cmp,
	}
	var wait func()
	if ee.t != nil {
		wait = ee.t.registerEventSync("test: update-keys: sync timed out")
	}
	if err := ee.scr.lib.PostEvent(evt); err != nil {
		return fmt.Errorf(errEventFmt, err)
	}
	if wait != nil {
		ee.t.t.Helper()
		wait()
	}
	return nil
}

type updateRunesEvent struct {
	when time.Time
	cmp  Componenter
}

func (u *updateRunesEvent) When() time.Time { return u.when }

// UpdateKeys posts a new UpdateKeys event into the event loop which
// once it is polled calls the given component's Keys implementation and
// registers/overwrites/removes provided keys-registrations.  UpdateKeys
// fails if the event-loop is full; returned error will wrap tcell's
// *PostEvent* error.  UpdateKeys is an no-op if Componenter is nil.
// NOTE in testing UpdateKeys starts the event loop if not running and
// returns after the event was fully processed.
func (ee *Events) UpdateKeys(cmp Componenter) error {
	if cmp == nil {
		return nil
	}
	if ee.t != nil && !ee.isListening {
		ee.t.listen()
	}
	evt := &updateKeysEvent{
		when: time.Now(),
		cmp:  cmp,
	}
	var wait func()
	if ee.t != nil {
		wait = ee.t.registerEventSync("test: update-keys: sync timed out")
	}
	if err := ee.scr.lib.PostEvent(evt); err != nil {
		return fmt.Errorf(errEventFmt, err)
	}
	if wait != nil {
		ee.t.t.Helper()
		wait()
	}
	return nil
}

type updateKeysEvent struct {
	when time.Time
	cmp  Componenter
}

func (u *updateKeysEvent) When() time.Time { return u.when }
