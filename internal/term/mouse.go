// Copyright (c) 2022 Stephan Lukits. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package term

import (
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/slukits/lines/internal/api"
)

var apiToTcellButtons = map[api.ButtonMask]tcell.ButtonMask{
	api.Button1:    tcell.Button1,
	api.Button2:    tcell.Button2,
	api.Button3:    tcell.Button3,
	api.Button4:    tcell.Button4,
	api.Button5:    tcell.Button5,
	api.Button6:    tcell.Button6,
	api.Button7:    tcell.Button7,
	api.Button8:    tcell.Button8,
	api.WheelUp:    tcell.WheelUp,
	api.WheelDown:  tcell.WheelDown,
	api.WheelLeft:  tcell.WheelLeft,
	api.WheelRight: tcell.WheelRight,
	api.ZeroButton: tcell.ButtonNone,
}

var tcellToApiButtons = map[tcell.ButtonMask]api.ButtonMask{
	tcell.Button1:    api.Button1,
	tcell.Button2:    api.Button2,
	tcell.Button3:    api.Button3,
	tcell.Button4:    api.Button4,
	tcell.Button5:    api.Button5,
	tcell.Button6:    api.Button6,
	tcell.Button7:    api.Button7,
	tcell.Button8:    api.Button8,
	tcell.WheelUp:    api.WheelUp,
	tcell.WheelDown:  api.WheelDown,
	tcell.WheelLeft:  api.WheelLeft,
	tcell.WheelRight: api.WheelRight,
	tcell.ButtonNone: api.ZeroButton,
}

// mouseEvent wraps a tcell mouse event to adapt it to the
// api.MouseEventer interface.
type mouseEvent struct{ evt *tcell.EventMouse }

func (e *mouseEvent) Pos() (int, int) { return e.evt.Position() }

func (e *mouseEvent) Button() api.ButtonMask {
	return tcellToApiButtons[e.evt.Buttons()]
}

func (e *mouseEvent) Mod() api.ModifierMask {
	return tcellToApiMods[e.evt.Modifiers()]
}

func (e *mouseEvent) When() time.Time { return e.evt.When() }

func (e *mouseEvent) Source() interface{} { return e.evt }

func mouseAggregator() func(e *tcell.EventMouse) api.MouseEventer {

	var last *tcell.EventMouse
	inDrag, ox, oy := false, 0, 0

	var clear = func(e *tcell.EventMouse) {
		last = nil
		ox, oy = e.Position()
		if inDrag {
			inDrag = false
		}
	}

	var eqBB = func(exp tcell.ButtonMask, ee ...*tcell.EventMouse) bool {
		for _, e := range ee {
			if exp == e.Buttons() {
				continue
			}
			return false
		}
		return true
	}

	var eqPos = func(
		e, other *tcell.EventMouse, ee ...*tcell.EventMouse,
	) bool {
		x, y := e.Position()
		ox, oy := other.Position()
		for _, _e := range ee {
			_x, _y := _e.Position()
			if x == _x && y == _y {
				continue
			}
			return false
		}
		return x == ox && y == oy
	}

	var zeroEvt = func(e, other *tcell.EventMouse) bool {
		return e.Buttons() == other.Buttons() && eqPos(e, other)
	}

	var zeroBtt = func(e *tcell.EventMouse) bool {
		return e.Buttons() == tcell.ButtonNone
	}

	return func(e *tcell.EventMouse) (evt api.MouseEventer) {
		switch last {
		case nil:
			if e.Buttons() == tcell.ButtonNone {
				// ignore zero-button without movement
				x, y := e.Position()
				if ox == x && oy == y {
					return
				}
				evt = api.NewMouseMove(ox, oy, &mouseEvent{evt: e})
				clear(e)
				return evt
			}
			last = e
			ox, oy = e.Position()
			return nil
		default:
			if zeroEvt(last, e) {
				return nil
			}
			if eqBB(last.Buttons(), e) {
				if !inDrag {
					inDrag = true
				}
				last = e
				return api.NewMouseDrag(ox, oy, &mouseEvent{evt: e})
			}
			if inDrag {
				inDrag = false
				evt = api.NewMouseDrop(&mouseEvent{evt: last})
			} else {
				evt = api.NewMouseClick(&mouseEvent{evt: last})
			}
			if zeroBtt(e) {
				clear(e)
				return evt
			}
			last = e
			return evt
		}
	}
}
