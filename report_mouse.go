// Copyright (c) 2022 Stephan Lukits. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package lines

import "github.com/gdamore/tcell/v2"

// Clicker is implemented by components which want to be informed about
// a "left"-mouse click event in their printable area.  If the clicked
// component, i.e. the component with the smallest layout area
// containing the event coordinates, does not have the focus an OnFocus
// event is reported first if and only if the clicked component has the
// Focusable feature.
type Clicker interface {

	// OnClick gets "left click"-events reported.  x and y provide the
	// click coordinates translated into the layouted area of the
	// receiving component.  E.g. y == 3 means that the component's
	// third line was clicked.  Env.Evt provides the reported
	// tcell.EventMouse event.
	OnClick(_ *Env, x, y int)
}

// Contexter is implemented by components which want to be informed
// about a mouse "right click"-event in their printable area.  If the
// clicked, i.e. the component with the smallest layout area containing
// the event coordinates,component does not have the focus an OnFocus
// event is reported first if and only if the clicked component has the
// Focusable feature.
//
// TODO: implement: see if event can also be reported for a potential
// context-menu key press (having x/y set to -1 then?).
type Contexter interface {

	// OnContext gets "right click"-events reported.  x and y provide
	// the click coordinates translated into the layouted area of the
	// receiving component.  E.g. y == 3 means that the component's
	// third line was clicked.  Env.Evt provides the reported
	// tcell.EventMouse event.
	OnContext(_ *Env, x, y int)
}

// Mouser is implemented by components who want to be informed about all
// mouse event in their printable area as they are reported by tcell.
// If the clicked component, i.e. the component with the smallest layout
// area containing the event coordinates, does not have the focus an
// OnFocus event is reported first if and only if the clicked component
// has the Focusable feature.  Mouse events are reported bubbling.
type Mouser interface {

	// OnMouse gets any mouse event reported as it is reported by tcell.
	// x and y provide the event coordinates translated into the
	// layouted area of the receiving component.  E.g. y == 3 means that
	// the component's third line was clicked.  Env.Evt provides the
	// reported tcell.EventMouse event.
	OnMouse(_ *Env, x, y int)
}

func mouseCurry(l func(*Env, int, int), x, y int) func(*Env) {
	return func(e *Env) { l(e, x, y) }
}

// reportMouse makes sure that the smallest component containing
// the click coordinates has the focus; then the click events are
// reported bubbling.
func reportMouse(cntx *rprContext) {

	x, y := cntx.evt.(*tcell.EventMouse).Position()
	path, err := cntx.scr.lyt.LocateAt(x, y)
	if err != nil {
		return
	}
	if len(path) == 0 {
		return
	}
	lytCmp := path[len(path)-1].(layoutComponenter)

	if lytCmp != cntx.scr.focus {
		focusIfFocusable(lytCmp, cntx)
	}

	for i := len(path) - 1; i >= 0; i-- {
		lc := path[i].(layoutComponenter)
		rx := x - lc.Dim().X()
		ry := y - lc.Dim().Y()

		if sb := reportClick(cntx, lc, rx, ry); sb {
			break
		}
		if sb := reportContext(cntx, lc, rx, ry); sb {
			break
		}
		if sb := reportOnMouse(cntx, lc, rx, ry); sb {
			break
		}
	}
}

// reportClick reports a "left"-click if an according mouse button
// was received and the focused component implements corresponding
// listener.
func reportClick(
	cntx *rprContext, lc layoutComponenter, x, y int,
) (stopBubbling bool) {

	if cntx.evt.(*tcell.EventMouse).Buttons()&tcell.ButtonPrimary ==
		tcell.ButtonNone {
		return
	}

	clk, ok := lc.userComponent().(Clicker)
	if !ok {
		return
	}

	env := callback(nil, cntx, mouseCurry(clk.OnClick, x, y))
	return env&envStopBubbling == envStopBubbling
}

// reportContext reports a "right"-click if an according mouse button
// was received and the focused component implements corresponding
// listener.
func reportContext(
	cntx *rprContext, lc layoutComponenter, x, y int,
) (stopBubbling bool) {

	if cntx.evt.(*tcell.EventMouse).Buttons()&tcell.ButtonSecondary ==
		tcell.ButtonNone {
		return
	}

	clk, ok := lc.userComponent().(Contexter)
	if !ok {
		return
	}

	env := callback(nil, cntx, mouseCurry(clk.OnContext, x, y))
	return env&envStopBubbling == envStopBubbling
}

func reportOnMouse(
	cntx *rprContext, lc layoutComponenter, x, y int,
) (stopBubbling bool) {

	msr, ok := lc.userComponent().(Mouser)
	if !ok {
		return false
	}
	env := callback(nil, cntx, mouseCurry(msr.OnMouse, x, y))
	return env&envStopBubbling == envStopBubbling
}
