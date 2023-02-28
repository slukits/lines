// Copyright (c) 2022 Stephan Lukits. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package lines

// Initer is implemented by components which want to be notified for
// initialization purposes before the first layout and before user input
// events are processed.  Note implement [AfterIniter] to be notified
// after all components were initialized; implement [Layouter] to be
// notified after each layout calculation of a component.
type Initer interface {

	// OnInit is executed before a component's layout is calculated the
	// first time.  Main use case of this event is to create nested
	// components and print initial content to a component.  NOTE if
	// several instances of a component type which implements OnInit are
	// added then OnInit is reported to each of this components.
	//
	// If specific dimensions should be set before the first layout
	// Env.ScreenSize provides the screen size to help with that.  If
	// these dimensions depend on subsequently created components
	// OnAfterInit is probably the right event.  If the calculated
	// layout should be adapted OnLayout may be the better event to do
	// so because it is called after the layout manager did its work.
	OnInit(*Env)
}

// AfterInit is implemented by components which want to be notified
// after all initial components have been created and their respective
// OnInit event was reported, see [Initer]; but before the first layout,
// see [Layouter].
type AfterIniter interface {

	// OnAfterInit is executed after all components have been
	// initialized and their respective OnInit event was reported.  I.e.
	// all components of an initial layout should be created and
	// available now.  Hence the main use case are initializations
	// depending on the existence of nested components. E.g. set the
	// focus or do size calculations depending on nested components.
	OnAfterInit(*Env)
}

// Focuser is implemented by components which want to be notified when
// they gain the focus.  Note if a component c has the focus and a
// component c' gets the focus then also all components which are
// parents of c' and are no parents of c get the focus, e.g.
//
//	+-App-----------------------------------------+
//	| +-chainer---------------------------------+ |
//	| | +-cmp1---------+ +-stacker------------+ | |
//	| | |              | | +-cmp2-----------+ | | |
//	| | |              | | |                | | | |
//	| | |              | | |                | | | |
//	| | |              | | +----------------+ | | |
//	| | +--------------+ +--------------------+ | |
//	| +-----------------------------------------+ |
//	+---------------------------------------------+
//
// if cmp1 has the focus and cmp2 gets the focus then stacker and cmp2
// get OnFocus reported (if implemented) while chainer and App which are
// parents of cmp1 and cmp2 don't.
type Focuser interface {

	// OnFocus is called back if an implementing component receives the
	// focus.
	OnFocus(*Env)
}

// FocusLooser is implemented by components which want to be informed
// when they loose the focus.  Note if a component c looses the focus to
// a component c' then all parents of c which are no parents of c' get
// also OnFocusLost reported, e.g.
//
//	+-App-----------------------------------------+
//	| +-chainer---------------------------------+ |
//	| | +-cmp1---------+ +-stacker------------+ | |
//	| | |              | | +-cmp2-----------+ | | |
//	| | |              | | |                | | | |
//	| | |              | | |                | | | |
//	| | |              | | +----------------+ | | |
//	| | +--------------+ +--------------------+ | |
//	| +-----------------------------------------+ |
//	+---------------------------------------------+
//
// if cmp2 looses the focus to cmp1 then cmp2 and stacker get
// OnFocusLost reported (if implemented) while chainer and App which are
// parents of cmp1 and cmp2 don't.
type FocusLooser interface {

	// OnFocusLost is called back if an implementing component looses
	// the focus.
	OnFocusLost(*Env)
}

// Updater is implemented by components which want to be informed about
// update events.  NOTE an Update event reaches only an Updater
// interface implementation if the [Lines.Update]-call was NOT provided
// with an event listener.
type Updater interface {

	// OnUpdate is reported for update requests without listener;
	// e.Evt.(*lines.UpdateEvent).Data provides the data optionally
	// provided to the update event registration.
	OnUpdate(e *Env, data interface{})
}

type rprContext struct {
	evt Eventer
	ll  *Lines
	scr *screen
}

// reportInit reports the init-event to the screen-components which
// haven't been initialized yet.
func reportInit(ll *Lines, scr *screen) {

	cntx := &rprContext{ll: ll, scr: scr}

	ii := []Componenter{}
	scr.forUninitialized(func(cmp Componenter) {
		if !cmp.hasLayoutWrapper() {
			return
		}

		if ic, ok := cmp.(Initer); ok {
			callback(cmp, cntx, ic.OnInit)
		}
		cmp.layoutComponent().wrapped().setInitialized()
		ii = append(ii, cmp.layoutComponent().userComponent())
	})
	for _, i := range ii {
		if ai, ok := i.(AfterIniter); ok {
			callback(i, cntx, ai.OnAfterInit)
		}
	}
}

func report(
	evt Eventer, ee *Lines, scr *screen,
) {

	cntx := &rprContext{evt: evt, ll: ee, scr: scr}

	switch evt := evt.(type) {
	case *UpdateEvent:
		reportUpdate(cntx, evt)
	case *moveFocusEvent:
		reportMoveFocus(cntx, evt)
	case RuneEventer:
		reportRune(cntx, evt)
	case KeyEventer:
		reportKey(cntx, evt)
	case *MouseMove:
		reportMouseMove(cntx, evt)
	case *MouseClick:
		reportMouseClick(cntx, evt)
	case *MouseDrag:
		reportMouseDrag(cntx, evt)
	case *MouseDrop:
		reportMouseDrop(cntx, evt)
	case MouseEventer:
		reportMouse(cntx, evt)
	}
}

func updateCurry(
	cb func(*Env, interface{}), data interface{},
) func(*Env) {
	return func(e *Env) { cb(e, data) }
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
	callback(evt.cmp, cntx, updateCurry(upd.OnUpdate, evt.Data))
}

func reportMoveFocus(cntx *rprContext, evt *moveFocusEvent) {
	if !evt.cmp.isInitialized() ||
		evt.cmp.layoutComponent().wrapped().dim.IsOffScreen() {
		return
	}
	if _, ok := cntx.scr.focus.userComponent().(Modaler); ok {
		return
	}
	moveFocus(evt.cmp, cntx)
}

func moveFocus(to Componenter, cntx *rprContext) {
	if to == cntx.scr.focus.userComponent() {
		return
	}
	if nestedInFocused(to, cntx) {
		setFocus(to, cntx)
		return
	}
	reportLostFocus(to.layoutComponent(), cntx)
	setFocus(to, cntx)
}

func nestedInFocused(to Componenter, cntx *rprContext) bool {
	path, err := cntx.scr.lyt.Locate(to.layoutComponent())
	if err != nil {
		return false
	}
	for _, d := range path {
		if d != cntx.scr.focus {
			continue
		}
		return true
	}
	return false
}

// setFocus reports OnFocus to all parents of the component to focus
// which are not parents of the currently focused component and set the
// screen focus.
func setFocus(to Componenter, cntx *rprContext) {
	fPath, e1 := cntx.scr.lyt.Locate(cntx.scr.focus)
	fPath = append(fPath, cntx.scr.focus)
	tPath, e2 := cntx.scr.lyt.Locate(to.layoutComponent())
	if e1 != nil || e2 != nil {
		return
	}
	tPath = append(tPath, to.layoutComponent())
	tPathIdx := 0
	for i, d := range fPath {
		if i >= len(tPath) {
			break
		}
		if tPath[i] != d {
			break
		}
		tPathIdx++
	}
	for _, d := range tPath[tPathIdx:] {
		usrCmp := d.(layoutComponenter).userComponent()
		fcs, ok := usrCmp.(Focuser)
		if ok {
			callback(usrCmp, cntx, fcs.OnFocus)
		}
	}
	cntx.scr.focus = to.layoutComponent()
}

// reportLostFocus reports focus lost to which are no parents of the
// given component the focus is set to.  E.g.
//
// +-App--------------------------------------------------+
// | +-container1------------+  +-container2------------+ |
// | | +--------+ +--------+ |  | +--------+ +--------+ | |
// | | | panel1 | | panel2 | |  | | panel3 | | panel4 | | |
// | | +--------+ +--------+ |  | +--------+ +--------+ | |
// | +-----------------------+  +-----------------------+ |
// +------------------------------------------------------+
//
// if the focus is moved from panel1 to panel2 then container1 doesn't
// get focus lost reported.  If on the other hand the focus is moved
// from panel2 to panel4 then container1 gets focus lost reported.
func reportLostFocus(to layoutComponenter, cntx *rprContext) {
	path, err := cntx.scr.lyt.Locate(cntx.scr.focus)
	if err != nil {
		return
	}
	path = append(path, cntx.scr.focus)
	lostPathIdx := 0
	if fPath, err := cntx.scr.lyt.Locate(to); err == nil {
		for i, d := range fPath {
			if len(path) <= i || d != path[i] {
				break
			}
			lostPathIdx++
		}
	}
	for _, d := range path[lostPathIdx:] {
		usrCmp := d.(layoutComponenter).userComponent()
		lsr, ok := usrCmp.(FocusLooser)
		if ok {
			callback(usrCmp, cntx, lsr.OnFocusLost)
		}
	}
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

	hasBeenEnabled := cmp.isEnabled()
	if !hasBeenEnabled {
		cmp.enable()
	}
	cb(env)
	if !hasBeenEnabled {
		cmp.disable()
	}
	env.reset()

	return env.flags
}
