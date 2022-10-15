// Copyright (c) 2022 Stephan Lukits. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package lines

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
	evt Eventer
	ll  *Lines
	scr *screen
}

// reportInit reports the init-event to the screen-components which
// haven't been initialized yet.
func reportInit(ll *Lines, scr *screen) {

	var reportedInit bool
	cntx := &rprContext{ll: ll, scr: scr}

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
		cmp.layoutComponent().wrapped().setInitialized()
	})
}

func report(
	evt Eventer, ee *Lines, scr *screen,
) (quit bool) {

	cntx := &rprContext{evt: evt, ll: ee, scr: scr}

	switch evt := evt.(type) {
	case QuitEventer:
		reportQuit(cntx)
	case *UpdateEvent:
		reportUpdate(cntx, evt)
	case *moveFocusEvent:
		reportMoveFocus(cntx, evt)
	case RuneEventer:
		return reportRune(cntx, evt)
	case KeyEventer:
		return reportKey(cntx, evt)
	case MouseEventer:
		reportMouse(cntx, evt)
	}
	return false
}

func reportUpdate(cntx *rprContext, evt *UpdateEvent) {
	if evt.lst != nil {
		callback(evt.cmp, cntx, evt.lst)
		return
	}
	if !evt.cmp.isInitialized() {
		return
	}
	upd, ok := evt.cmp.layoutComponent().userComponent().(Updater)
	if !ok {
		return
	}
	callback(evt.cmp, cntx, upd.OnUpdate)
}

func reportMoveFocus(cntx *rprContext, evt *moveFocusEvent) {
	if !evt.cmp.isInitialized() {
		return
	}
	moveFocus(evt.cmp, cntx)
}

func moveFocus(cmp Componenter, cntx *rprContext) {
	usrCmp := cntx.scr.focus.userComponent()
	if cmp == usrCmp {
		return
	}
	fls, ok := usrCmp.(FocusLooser)
	if ok {
		callback(usrCmp, cntx, fls.OnFocusLost)
	}
	fcs, ok := cmp.(Focuser)
	if ok {
		callback(cmp, cntx, fcs.OnFocus)
	}
	cntx.scr.focus = cmp.layoutComponent()
}

func mouseFocusable(
	evt MouseEventer,
) func(ff *features, recursive bool) bool {

	return func(ff *features, recursive bool) bool {

		if ff == nil {
			return false
		}
		f := ff.buttonFeature(evt.Button(), evt.Mod())

		if !recursive && f&Focusable != NoFeature {
			return true
		}

		if recursive && f&(Focusable|_recursive) != NoFeature {
			return true
		}

		return false
	}
}

func focusIfFocusable(cntx *rprContext, cmp layoutComponenter) bool {

	focusAble := true // we abuse this as indicator for recursive ...
	var isFocusable func(*features, bool) bool
	switch evt := cntx.evt.(type) {
	case MouseEventer:
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

func cbEnv(cntx *rprContext, cmp Componenter) *Env {
	return &Env{
		cmp:   cmp,
		Lines: cntx.ll,
		Evt:   cntx.evt,
		size:  cntx.scr.backend.Size,
	}
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
	env := cbEnv(cntx, cmp)

	cmp.enable()
	cb(env)
	cmp.disable()
	env.reset()

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
}
