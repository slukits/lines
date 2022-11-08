// Copyright (c) 2022 Stephan Lukits. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package lines

func clickCurry(l func(*Env, int, int), x, y int) func(*Env) {
	return func(e *Env) { l(e, x, y) }
}

func mouseCurry(
	l func(*Env, ButtonMask, int, int), bm ButtonMask, x, y int,
) func(*Env) {
	return func(e *Env) { l(e, bm, x, y) }
}

// reportMouse makes sure that the smallest component containing
// the click coordinates has the focus; then the click events are
// reported bubbling.
func reportMouse(cntx *rprContext, evt MouseEventer) {

	x, y := evt.Pos()
	path, err := cntx.scr.lyt.LocateAt(x, y)
	if err != nil {
		return
	}
	if len(path) == 0 {
		return
	}
	lytCmp := path[len(path)-1].(layoutComponenter)

	if lytCmp != cntx.scr.focus {
		focusIfFocusable(cntx, lytCmp)
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

	if cntx.evt.(MouseEventer).Button()&Primary == ZeroButton {
		return
	}

	clk, ok := lc.userComponent().(Clicker)
	if !ok {
		return
	}

	env := callback(nil, cntx, clickCurry(clk.OnClick, x, y))
	return env&envStopBubbling == envStopBubbling
}

// reportContext reports a "right"-click if an according mouse button
// was received and the focused component implements corresponding
// listener.
func reportContext(
	cntx *rprContext, lc layoutComponenter, x, y int,
) (stopBubbling bool) {

	if cntx.evt.(MouseEventer).Button()&Secondary ==
		ZeroButton {
		return
	}

	clk, ok := lc.userComponent().(Contexter)
	if !ok {
		return
	}

	env := callback(nil, cntx, clickCurry(clk.OnContext, x, y))
	return env&envStopBubbling == envStopBubbling
}

func reportOnMouse(
	cntx *rprContext, lc layoutComponenter, x, y int,
) (stopBubbling bool) {

	msr, ok := lc.userComponent().(Mouser)
	if !ok {
		return false
	}
	env := callback(nil, cntx, mouseCurry(
		msr.OnMouse, cntx.evt.(MouseEventer).Button(), x, y))
	return env&envStopBubbling == envStopBubbling
}
