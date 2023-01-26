// Copyright (c) 2022 Stephan Lukits. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package lines

import (
	"fmt"
	"sync"
	"time"

	"github.com/slukits/lines/internal/api"
	"github.com/slukits/lines/internal/lyt"
	"github.com/slukits/lines/internal/term"
)

/*
Eventer is the interface which all reported events implement.  Note each
Env instance has an Env.Evt property of type Eventer whereas
Env.Evt.Source() provides the backend event if there is any.  The
following event interfaces with their reported event are defined:
  - [Initer]: OnInit(*Env): once before any other
  - [AfterIniter]: OnAfterInit(*Env): once after OnInit
  - [Focuser]: OnFocus(*Env): see [Lines.Focus], [Env.Focused]
  - [FocusLooser]: OnFocusLost(*Env) see [Lines.Focus], [Env.Focused]
  - [Updater]: OnUpdate(*Env, interface{}): see [Lines.Update]
  - [Layouter]: OnLayout(*Env) bool: after layout change
  - [Keyer]: OnKey(*Env, Key, ModifierMask): special key like Esc
  - [Runer]: OnRune(*Env, rune, ModifierMask)
  - [Enterer]: OnEnter(*Env): mouse-pointer entered component
  - [Exiter]: OnExit(*Env): mouse-pointer leaves component
  - [Mouser]: OnMouse(*Env, ButtonMask, int, int): any mouse-event
  - [Clicker]: OnClick(_ *Env, x, y int): primary button click
  - [Contexter]: OnContext(_ *Env, x, y int): secondary button click
  - [Drager]: OnDrag(*Env, ButtonMask, int, int): mouse-move with
    pressed button
  - [Dropper]: OnDrop(*Env, ButtonMask, int, int): button release after
    mouse-move with pressed button
  - [Modaler]: OnOutOfBoundClick(*Env) bool: for modal layers
  - [OutOfBoundMover]: OnOutOfBoundMove(*Env) bool: for modal layers
  - [LineSelecter]: OnLineSelection(*Env, int): [LineSelectable]
*/
type Eventer = api.Eventer

// resizeEventer is reported when the Lines-display was resized.
type resizeEventer = api.ResizeEventer

// Dimer provides dimensions of a component in the layout.  Note each
// type embedding [lines.Component] implements the Dimer interface.
type Dimer = lyt.Dimer

// quitting controls the quitting keys and runes of a lines instance.
type quitting struct {
	*sync.Mutex
	kk map[Key]bool
	rr map[rune]bool
}

// AddKey adds given Key k to the quitting keys of parent Lines
// instance.
func (q *quitting) AddKey(k Key) *quitting {
	if q == nil {
		return q
	}
	q.Lock()
	defer q.Unlock()
	if q.kk == nil {
		q.kk = map[Key]bool{}
	}
	q.kk[k] = true
	return q
}

// DelKey removes given Key k from the quitting keys of parent Lines
// instance.
func (q *quitting) DelKey(k Key) *quitting {
	if q == nil {
		return q
	}
	q.Lock()
	defer q.Unlock()
	delete(q.kk, k)
	return q
}

// AddRune adds given rune r to the quitting runes.
func (q *quitting) AddRune(r rune) *quitting {
	if q == nil {
		return q
	}
	q.Lock()
	defer q.Unlock()
	if q.rr == nil {
		q.rr = map[rune]bool{}
	}
	q.rr[r] = true
	return q
}

// DelRune removes given rune r from the quitting runes.
func (q *quitting) DelRune(r rune) *quitting {
	if q == nil {
		return q
	}
	q.Lock()
	defer q.Unlock()
	delete(q.rr, r)
	return q
}

// Rune returns true if given rune r quits the event-loop of parent
// Lines instance.
func (q *quitting) Rune(r rune) bool {
	if q == nil || q.rr == nil {
		return false
	}
	return q.rr[r]
}

// Key returns true if given Key k quits the event-loop of parent Lines
// instance.  Note CtrlC and CtrlD returns always true if parent Lines
// instance wasn't created by a Kiosk-constructor.
func (q *quitting) Key(k Key) bool {
	if q == nil {
		return false
	}
	if k == CtrlC || k == CtrlD {
		return true
	}
	if q.kk == nil {
		return false
	}
	return q.kk[k]
}

// Lines listens to a backend implementation's reporting of events and
// controls the event reporting to client components (see [Component])
// and their layout accordingly.  Use one of the constructors [Term],
// [TermKiosk] or [TermFixture] to obtain a Lines-instance.
type Lines struct {

	// scr to report resize events to screen components.
	scr *screen

	// Quitting provides an instance to control key/rune-bindings for
	// terminating the event-loop and unblocking [Lines.WaitForQuit].
	// Note if Quitting is set the keys CtrlC/CtrlD always quit.
	Quitting *quitting

	// backend is needed to post events.
	backend api.EventProcessor

	// Globals are properties whose changing is propagated to all its
	// clones in components who update iff the updated property is still
	// in sync with the origin.
	Globals *globals
}

func newTerm(cmp Componenter) *Lines {
	ll := Lines{}
	ll.backend = term.New(ll.listen)
	ll.Globals = newGlobals(nil)
	ll.scr = newScreen(ll.backend.(api.UIer), cmp, ll.Globals)
	ll.Globals.propagation = globalsPropagationClosure(ll.scr)
	return &ll
}

// Term returns a Lines l instance with a terminal backend displaying
// and reporting events to given component cmp and its nested
// components.  By default 'q', ctrl-c and ctrl-d terminate the
// event-loop and unblock  l.WaitForQuit().  The binding to 'q' may be
// removed.  The bindings to ctrl-c and ctrl-d may not be removed.  Use
// the [TermKiosk] constructor for an setup without any quit bindings.
// Term panics if the terminal screen can't be obtained.  NOTE to create
// a Componenter-instance define a type which embeds the [Component]
// type:
//
//	type myComponent struct { lines.Component }
//	lines.Term(&myComponent{}).WaitForQuit()
//
// Leverage [Lines.OnQuit] registration if you want to be informed about
// the quit event which is triggered by user-input that is associated
// with the quitable feature or by calling [Lines.Quit].  After your
// application is initialized you typically will want to wait while
// processing user input until the quit event occurs using
// [Lines.WaitForQuit].
func Term(cmp Componenter) *Lines {
	ll := newTerm(cmp)
	ll.Quitting = &quitting{Mutex: &sync.Mutex{}}
	ll.Quitting.AddRune('q')
	return ll
}

// Componenter is the private interface a type must implement to be used
// as an lines ui component.  Embedding [lines.Component] in a type
// automatically fulfills this condition:
//
//	type MyCmp struct { lines.Component }
//	lines.Term(&MyCmp{}).WaitForQuit()
//
// a Componenter implementation is informed about application events if
// it also implements event listener interfaces like
//   - [Initer] is informed once before a component becomes part of the layout
//   - [Layouter] is informed that a component's layout was calculated
//   - [Focuser]/[FocusLooser] is informed about focus gain/loss
//   - [Keyer] is informed about any user special key-press like 'enter' or 'tab'
//   - [Runer] is informed about user rune-key input
//   - [Mouser] is informed about any mouse event see also [Clicker]/[Contexter]
//   - [LineSelecter] is informed if a component's line was selected
//   - [LineFocuser] is informed if a component's line received the focus
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
	initialize(Componenter, api.UIer, *globals) layoutComponenter

	// isInitialized returns true if embedded *component was wrapped
	// into a layout component.
	isInitialized() bool

	// embedded returns a reference to client-component's embedded
	// Component-instance.
	embedded() *Component

	// backend to post Update and Focus events on a user Componenter
	// implementation.
	backend() api.UIer

	// globals returns a components global properties.
	globals() *globals
}

// TermKiosk returns a Lines instance like [Term] but without registered
// Quitable feature, i.e. the application can't be quit by the user by
// default.
func TermKiosk(cmp Componenter) *Lines {
	return newTerm(cmp)
}

// SetRoot replaces currently used root component by given component.
func (ll *Lines) SetRoot(c Componenter) error {
	if c == nil {
		return nil
	}
	return ll.backend.Post(&rootEvent{
		when:    time.Now(),
		newRoot: c,
	})
}

// moveFocusEvent is posted by calling MoveFocus for a programmatically
// change of focus.  This event-instance is not provided to the user.
type rootEvent struct {
	when    time.Time
	newRoot Componenter
}

func (e *rootEvent) When() time.Time { return e.when }

func (e *rootEvent) Source() interface{} { return e }

// Quit posts a quit event which consequently closes given Lines
// instance's backend and unblocks WaitForQuit.
func (ll *Lines) Quit() { ll.backend.Quit() }

// OnQuit registers given function to be called on quitting
// event-polling and -reporting.
func (ll *Lines) OnQuit(listener func()) *Lines {
	ll.backend.OnQuit(listener)
	return ll
}

// WaitForQuit blocks until given Lines-instance is quit.  (Except a
// Lines instance provided by a [Fixture] in which case WaitForQuit is
// not blocking.)
func (ll *Lines) WaitForQuit() { ll.backend.WaitForQuit() }

// Update posts an update event into the event queue which is reported
// either to given listener if not nil or to given componenter if given
// listener is nil.  Given data will be provided by the reported Update
// event.  Update is a no-op if componenter and listener are nil.
//
// Note if Update is used in a test to retrieve component properties
// which may be only retrieved in a listener callback then you can not
// cancel the test-run inside the listener.  "FatalNow" and friends
// must be only called in the go-routine of the test but the
// Update-listener callback during testing is done from the go-routine
// which is running the event-loop and which is different from the
// go-routine running the test.
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

// UpdateEvent is created by an [Lines.Update] call.  Its Data field
// provides the data which was passed to that Update call.  To get
// notified of an update event a component must implement [Updater]:
//
//	func (c *myCmp) OnUpdate(e *lines.Env) {
//	    d := e.Evt.(*lines.UpdateEvent).Data
//	}
//
// Note was an explicit listener passed to [Lines.Update] the event is not
// reported to the component.
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

func (ll *Lines) listen(evt api.Eventer) {
	if evt, ok := evt.(RuneEventer); ok {
		if ll.Quitting.Rune(evt.Rune()) {
			ll.backend.Quit()
			return
		}
	}
	if evt, ok := evt.(KeyEventer); ok {
		if ll.Quitting.Key(evt.Key()) {
			ll.backend.Quit()
			return
		}
	}
	switch evt := evt.(type) {
	case *rootEvent:
		ll.scr.setRoot(evt.newRoot, ll.Globals)
		ll.scr.hardSync(ll)
	case resizeEventer:
		width, height := evt.Size()
		postSync := ll.scr.setSize(width, height, ll)
		reportInit(ll, ll.scr)
		ll.scr.hardSync(ll)
		if postSync != nil {
			postSync()
		}
	default:
		report(evt, ll, ll.scr)
		reportInit(ll, ll.scr)
		ll.scr.softSync(ll)
	}
}

// MoveFocus posts a new MoveFocus event into the event loop which once
// it is polled calls the currently focused component's OnFocusLost
// implementation while given component's OnFocus implementation is
// executed.  Finally the focus is set to given component.  MoveFocus
// fails if the event-loop is full returned error will wrap tcell's
// *PostEvent* error.  MoveFocus is an no-op if Componenter is nil, is
// not part of the layout, is off-screen or is already focused.
func (ll *Lines) Focus(cmp Componenter) error {
	return ll.backend.Post(&moveFocusEvent{
		when: time.Now(),
		cmp:  cmp,
	})
}

// CursorComponent returns the component where currently the cursor is
// set or nil if the cursor is not set.
func (ll *Lines) CursorComponent() Componenter {
	cc := ll.scr.cursorComponent()
	if cc == nil {
		return nil
	}
	return cc.userComponent()
}

// CursorPosition returns given lines ll instance's cursor screen
// position (x,y) and true in case the cursor is set; otherwise -1, -1
// and false is returned.
func (ll *Lines) CursorPosition() (x, y int, _ bool) {
	if ll.scr.cursor.Removed() {
		return -1, -1, false
	}
	x, y = ll.scr.cursor.Coordinates()
	return x, y, true
}

// RemoveCursor removes the cursor from the screen.
func (ll *Lines) RemoveCursor() {
	if lc := ll.scr.cursorComponent(); lc != nil {
		lc.wrapped().setCursor(-1, -1, ZeroCursor)
		if ll.scr.cursor.Removed() {
			return
		}
	}
	ll.scr.setCursor(-1, -1, ZeroCursor)
}

// SetCursor sets the cursor to given coordinates (x,y) on the screen
// having optionally given cursor style cs[0].  The call is ignored
// respectively a currently set cursor is removed iff x and y are
// outside the screen area or if they are inside the screen area but
// cs[0] is the zero cursor style.
func (ll *Lines) SetCursor(x, y int, cs ...CursorStyle) {
	ll.scr.setCursor(x, y, cs...)
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
// to define differentiated stylings printing to a component's
// environment e:
//
//	fmt.Fprint(e.LL(0), "An ")
//	lines.Print(e.LL(0).At(3).AA(Italic), "italic")
//	lines.Print(e.LL(0).At(9), " word")
type AtWriter interface {
	WriteAt(rr []rune)
}

// Print to an AtWriter a rune or slice of runes.  The common AtWriter
// of lines are provided by a component listener environment and a
// components gaps:
//
//	lines.Print(e.LL(0).At(0), []rune("print to first line's first cell"))
//	lines.Print(cmp.Gaps(0).Left.At(5), []rune("print to left gutter"))
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
