// Copyright (c) 2022 Stephan Lukits. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package lines

import (
	"github.com/slukits/lines/internal/api"
	"github.com/slukits/lines/internal/lyt"
)

func mouseCurry(
	l func(*Env, ButtonMask, int, int), bm ButtonMask, x, y int,
) func(*Env) {
	return func(e *Env) { l(e, bm, x, y) }
}

func continueReportOnModal(
	mc Modaler, cntx *rprContext, evt MouseEventer,
) bool {
	x, y := evt.Pos()
	continueReport := true
	if !cntx.scr.focus.Dim().Contains(x, y) {
		callback(cntx.scr.focus.userComponent(), cntx, func(e *Env) {
			continueReport = mc.OnOutOfBoundClick(e)
		})
	}
	return continueReport
}

func cancelOnModal(cntx *rprContext, evt MouseEventer) bool {
	x, y := evt.Pos()
	if _, ok := cntx.scr.focus.userComponent().(Modaler); ok {
		if !cntx.scr.focus.Dim().Contains(x, y) {
			return true
		}
	}
	return false
}

func cancelOnModalDrag(cntx *rprContext, evt *MouseDrag) bool {
	if _, ok := cntx.scr.focus.userComponent().(Modaler); !ok {
		return false
	}
	drg, ok := cntx.scr.focus.userComponent().(Drager)
	if !ok {
		return true
	}
	x, y := evt.Pos()
	callback(cntx.scr.focus.userComponent(), cntx,
		mouseCurry(drg.OnDrag, evt.Button(), x, y))
	return true
}

func cancelOnModalMove(cntx *rprContext, evt *MouseMove) bool {
	if _, ok := cntx.scr.focus.userComponent().(Modaler); !ok {
		return false
	}
	x, y := evt.Pos()
	if cntx.scr.focus.Dim().Contains(x, y) {
		return false
	}
	oob, ok := cntx.scr.focus.userComponent().(OutOfBoundMover)
	if !ok {
		return true
	}
	continueReport := false
	callback(cntx.scr.focus.userComponent(), cntx, func(e *Env) {
		continueReport = oob.OnOutOfBoundMove(e)
	})
	return !continueReport
}

func posCurry(l func(*Env, int, int), x, y int) func(*Env) {
	return func(e *Env) { l(e, x, y) }
}

func reportMouseMove(cntx *rprContext, evt *MouseMove) {

	if cancelOnModalMove(cntx, evt) {
		return
	}

	x, y := evt.Pos()
	path, err := cntx.scr.lyt.LocateAt(x, y)
	if err != nil || path == nil {
		return
	}
	reported := reportEnterExit(
		cntx, path[len(path)-1].(layoutComponenter), evt)
	if reported && len(path) == 1 {
		return
	}
	if reported {
		path = path[:len(path)-1]
	}

	reportBubbling(
		cntx, path, x, y, true,
		func(c Componenter) bool {
			_, ok := c.(Mover)
			return ok
		},
		func(c Componenter, x, y int) func(*Env) {
			return posCurry(c.(Mover).OnMove, x, y)
		},
	)
}

func reportEnterExit(
	cntx *rprContext, in layoutComponenter, evt *MouseMove,
) (reported bool) {
	x, y := evt.Pos()
	if cntx.scr.mouseIn == nil {
		if in.wrapped().InContentArea(x, y) {
			cntx.scr.mouseIn = in
			if e, ok := in.userComponent().(Enterer); ok {
				rx, ry := relative(in.wrapped().ContentArea, x, y)
				callback(in.userComponent(), cntx,
					posCurry(e.OnEnter, rx, ry))
				return true
			}
		}
		return false
	}

	// NOTE we need compare user components because different
	// layoutComponenter may wrap the same component if layered and
	// un-layered.
	if in.userComponent() == cntx.scr.mouseIn.userComponent() {
		if !in.wrapped().InContentArea(x, y) {
			cntx.scr.mouseIn = nil
			if e, ok := in.userComponent().(Exiter); ok {
				callback(in.userComponent(), cntx, e.OnExit)
				return true
			}
		}
		return false
	}

	ox, oy := evt.Origin()
	if cntx.scr.mouseIn.wrapped().InContentArea(ox, oy) {
		if e, ok := cntx.scr.mouseIn.userComponent().(Exiter); ok {
			callback(cntx.scr.mouseIn.userComponent(), cntx, e.OnExit)
			cntx.scr.mouseIn = nil
		}
	}

	if !in.wrapped().InContentArea(x, y) {
		cntx.scr.mouseIn = nil
		return
	}

	cntx.scr.mouseIn = in
	if e, ok := in.userComponent().(Enterer); ok {
		rx, ry := relative(in.wrapped().ContentArea, x, y)
		callback(in.userComponent(), cntx, posCurry(e.OnEnter, rx, ry))
		return true
	}
	return false
}

func relative(ar func() (x, y, w, h int), x, y int) (rx, ry int) {
	ax, ay, _, _ := ar()
	return x - ax, y - ay
}

func reportMouseClick(cntx *rprContext, evt *MouseClick) {
	if evt.Button()&(api.Primary|api.Secondary) == 0 {
		return
	}
	if mc, ok := cntx.scr.focus.userComponent().(Modaler); ok {
		if !continueReportOnModal(mc, cntx, evt) {
			return
		}
	}

	x, y := evt.Pos()
	path := focusedPath(cntx, evt, x, y)
	if path == nil {
		return
	}

	if evt.Button()&api.Primary == api.Primary {
		reportPrimary(cntx, evt, path, x, y)
		return
	}
	reportSecondary(cntx, evt, path, x, y)
}

// reportPrimary reports a "left"-click if an according mouse button
// was received and the focused component implements corresponding
// listener.
func reportPrimary(
	cntx *rprContext, evt *MouseClick, path []lyt.Dimer, x, y int,
) {
	if scroll(path[len(path)-1], true, x, y) {
		return
	}
	reportBubbling(
		cntx, path, x, y, true,
		func(c Componenter) bool {
			_, ok := c.(Clicker)
			return ok
		},
		func(c Componenter, x, y int) func(*Env) {
			return posCurry(c.(Clicker).OnClick, x, y)
		},
	)
}

func reportSecondary(
	cntx *rprContext, evt *MouseClick, path []lyt.Dimer, x, y int,
) {
	if scroll(path[len(path)-1], false, x, y) {
		return
	}
	reportBubbling(
		cntx, path, x, y, true,
		func(c Componenter) bool {
			_, ok := c.(Contexter)
			return ok
		},
		func(c Componenter, x, y int) func(*Env) {
			return posCurry(c.(Contexter).OnContext, x, y)
		},
	)
}

func scroll(d lyt.Dimer, down bool, x, y int) bool {
	cmp := d.(layoutComponenter).wrapped()
	cmp.userCmp.embedded().enable()
	defer cmp.userCmp.embedded().disable()
	if cmp.Scroll.BarContains(x, y) {
		if down {
			cmp.Scroll.Down()
		} else {
			cmp.Scroll.Up()
		}
		return true
	}
	return false
}

func focusedPath(
	cntx *rprContext, evt *MouseClick, x, y int,
) []lyt.Dimer {

	path, err := cntx.scr.lyt.LocateAt(x, y)
	if err != nil || len(path) == 0 {
		return nil
	}

	// find the deepest nested focusable component in the path ...
	var focus layoutComponenter
	for i := len(path) - 1; i >= 0; i-- {
		ff := path[i].(layoutComponenter).wrapped().ff
		if ff == nil {
			continue
		}
		f := ff.buttonFeature(evt.Button(), evt.Mod())
		if f&Focusable == NoFeature {
			continue
		}
		focus = path[i].(layoutComponenter)
		break
	}
	if focus == nil {
		return path
	}
	// ... and focus it
	moveFocus(focus.userComponent(), cntx)
	return path
}

func reportMouseDrag(cntx *rprContext, evt *MouseDrag) {
	if cancelOnModalDrag(cntx, evt) {
		return
	}

	x, y := evt.Pos()
	path, err := cntx.scr.lyt.LocateAt(x, y)
	if err != nil || path == nil {
		return
	}

	reportBubbling(
		cntx, path, x, y, false,
		func(c Componenter) bool {
			_, ok := c.(Drager)
			return ok
		},
		func(c Componenter, x, y int) func(*Env) {
			return mouseCurry(c.(Drager).OnDrag, evt.Button(), x, y)
		},
	)
}

func reportMouseDrop(cntx *rprContext, evt *MouseDrop) {
	if cancelOnModal(cntx, evt) {
		return
	}

	x, y := evt.Pos()
	path, err := cntx.scr.lyt.LocateAt(x, y)
	if err != nil || path == nil {
		return
	}

	reportBubbling(
		cntx, path, x, y, false,
		func(c Componenter) bool {
			_, ok := c.(Dropper)
			return ok
		},
		func(c Componenter, x, y int) func(*Env) {
			return mouseCurry(c.(Dropper).OnDrop, evt.Button(), x, y)
		},
	)
}

func reportMouse(cntx *rprContext, evt MouseEventer) {
	if cancelOnModal(cntx, evt) {
		return
	}

	path, err := cntx.scr.lyt.LocateAt(evt.Pos())
	if err != nil || len(path) == 0 {
		return
	}

	x, y := evt.Pos()
	for i := len(path) - 1; i >= 0; i-- {

		lc := path[i].(layoutComponenter)
		clk, ok := lc.userComponent().(Mouser)
		if !ok {
			continue
		}

		rx := x - lc.Dim().X()
		ry := y - lc.Dim().Y()
		env := callback(lc.userComponent(), cntx, mouseCurry(
			clk.OnMouse, evt.Button(), rx, ry))

		if env&envStopBubbling == envStopBubbling {
			break
		}
	}
}

func reportBubbling(
	cntx *rprContext, path []lyt.Dimer, x, y int, relative bool,
	implements func(Componenter) bool,
	curry func(_ Componenter, x, y int) func(*Env),
) {
	for i := len(path) - 1; i >= 0; i-- {

		usrCmp := path[i].(layoutComponenter).userComponent()
		wrapped := usrCmp.layoutComponent().wrapped()
		if !implements(usrCmp) || (relative && !wrapped.InContentArea(x, y)) {
			continue
		}

		ax, ay, _, _ := wrapped.ContentArea()
		rx, ry := x, y
		if relative {
			rx, ry = rx-ax, ry-ay
		}
		env := callback(usrCmp, cntx, curry(usrCmp, rx, ry))

		if env&envStopBubbling == envStopBubbling {
			break
		}
	}
}
