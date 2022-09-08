// Copyright (c) 2022 Stephan Lukits. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package lines

import "github.com/gdamore/tcell/v2"

// Initer is implemented by components which want to be notified for
// initialization purposes before the first layout and before user input
// events are processed.
type Initer interface {

	// OnInit provides to its implementation an environment before it is
	// layouted the first time.  NOTE if several instances of a
	// component type which implements OnInit are added then OnInit is
	// reported to each of this components.
	//
	// If specific dimensions should be set before the first layout
	// Env.ScreenSize provides the screen size to help with that.  If
	// the calculated layout should be adapted OnLayout may be the
	// better event to do so because it is called after the layout
	// manager did its work.
	OnInit(*Env)
}

// Focuser is implemented by components which want to be notified when
// they gain the focus.
type Focuser interface {

	// OnFocus is called back if an implementing component receives the
	// focus.
	OnFocus(*Env)
}

// FocusLooser is implemented by components which want to be informed
// when they loose the focus.
type FocusLooser interface {

	// OnFocusLost get focus loss reported.
	OnFocusLost(*Env)
}

// Quitter is implemented by components which want to be informed when
// the event-loop ends.
type Quitter interface {

	// OnQuit is reported to all components implementing this interface not
	// only to the one having the focus before the event loop ends.
	OnQuit()
}

// Updater is implemented by components which want to be informed about
// update events.  NOTE an Update event reaches only an Updater
// interface if the Update-call on Events was NOT provided with an
// event listener.
type Updater interface {

	// OnUpdate is reported for update requests without listener;
	// *Env.Evt.(*lines.UpdateEvent).Data provides the data optionally
	// provided to the update event registration.
	OnUpdate(*Env)
}

type rprContext struct {
	evt tcell.Event
	ee  *Events
	scr *screen
}

// reportInit reports the init-event to the screen-components which
// haven't been initialized yet.
func reportInit(ee *Events, scr *screen) {

	var reportedInit bool
	cntx := &rprContext{ee: ee, scr: scr}

	scr.forUninitialized(func(cmp Componenter) {
		if !cmp.hasLayoutWrapper() {
			return
		}

		if ic, ok := cmp.(Initer); ok {
			callback(cmp, cntx, ic.OnInit)
			if !reportedInit {
				reportedInit = true
			}
		}
		registerKeys(cmp, cntx)
		registerRunes(cmp, cntx)
		cmp.layoutComponent().wrapped().setInitialized()
	})
	if reportedInit {
		reportReported(ee)
	}
}

func report(
	evt tcell.Event, ee *Events, scr *screen,
) (quit bool) {

	cntx := &rprContext{evt: evt, ee: ee, scr: scr}

	switch evt := evt.(type) {
	case *UpdateEvent:
		reportUpdate(cntx)
	case *moveFocusEvent:
		reportMoveFocus(cntx)
	case *updateKeysEvent:
		registerKeys(evt.cmp, cntx)
	case *updateRunesEvent:
		registerRunes(evt.cmp, cntx)
	case *tcell.EventKey:
		switch evt.Key() {
		case tcell.KeyRune:
			return reportRune(cntx)
		default:
			return reportKey(cntx)
		}
	case *tcell.EventMouse:
		reportMouse(cntx)
	case *quitEvent:
		reportQuit(cntx)
		return true
	}
	return false
}

func reportUpdate(cntx *rprContext) {
	evt := cntx.evt.(*UpdateEvent)
	if evt.lst != nil {
		callback(evt.cmp, cntx, evt.lst)
		return
	}
	upd, ok := evt.cmp.(Updater)
	if !ok {
		return
	}
	callback(evt.cmp, cntx, upd.OnUpdate)
}

func reportMoveFocus(cntx *rprContext) {
	moveFocus(cntx.evt.(*moveFocusEvent).cmp, cntx)
}

func moveFocus(cmp Componenter, cntx *rprContext) {
	fls, ok := cntx.scr.focus.userComponent().(FocusLooser)
	if ok {
		callback(cntx.scr.focus.userComponent(), cntx, fls.OnFocusLost)
	}
	fcs, ok := cmp.(Focuser)
	if ok {
		callback(cmp, cntx, fcs.OnFocus)
	}
	cntx.scr.focus = cmp.layoutComponent()
}

func mouseFocusable(
	evt *tcell.EventMouse,
) func(ff *features, recursive bool) bool {

	return func(ff *features, recursive bool) bool {

		if ff == nil {
			return false
		}
		f := ff.buttonFeature(evt.Buttons(), evt.Modifiers())

		if !recursive && f&Focusable != NoFeature {
			return true
		}

		if recursive && f&(Focusable|_recursive) != NoFeature {
			return true
		}

		return false
	}
}

func focusIfFocusable(cmp layoutComponenter, cntx *rprContext) bool {

	focusAble := true // we abuse this as indicator for recursive ...
	var isFocusable func(*features, bool) bool
	switch evt := cntx.evt.(type) {
	case *tcell.EventMouse:
		isFocusable = mouseFocusable(evt)
	default:
		return false
	}

	cntx.scr.forBubbling(cmp, func(lc layoutComponenter) (stop bool) {
		if focusAble {
			if isFocusable(lc.wrapped().ff, !focusAble) {
				return true
			}
			focusAble = false
			return
		}
		if isFocusable(lc.wrapped().ff, !focusAble) {
			focusAble = true
			return true
		}
		return
	})

	if !focusAble {
		return false
	}

	moveFocus(cmp.userComponent(), cntx)
	return true
}

func sizeClosure(scr tcell.Screen) func() (int, int) {
	return func() (int, int) { return scr.Size() }
}

func callback(
	cmp Componenter, cntx *rprContext, cb func(*Env),
) (flags envMask) {

	if cmp == nil {
		cmp = cntx.scr.focus.userComponent()
	}
	if !cmp.hasLayoutWrapper() {
		return
	}
	env := &Env{cmp: cmp, EE: cntx.ee, Evt: cntx.evt,
		size: sizeClosure(cntx.scr.lib)}

	cmp.enable()
	cb(env)
	cmp.disable()
	env.reset()

	if cntx.evt != nil {
		reportReported(cntx.ee)
	}
	return env.flags
}

func reportQuit(cntx *rprContext) {
	reported := false
	cntx.scr.forComponent(func(c Componenter) {
		if qtt, ok := c.(Quitter); ok {
			qtt.OnQuit()
			if !reported {
				reported = true
			}
		}
	})
	if !reported {
		return
	}
	reportReported(cntx.ee)
}

func reportReported(ee *Events) {
	if ee.reported == nil {
		return
	}
	ee.reported()
}
