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

type mouseEvent struct{ evt *tcell.EventMouse }

func newMouseEvent(x, y int, b api.ButtonMask, m api.ModifierMask) api.MouseEventer {
	return &mouseEvent{evt: tcell.NewEventMouse(
		x, y,
		apiToTcellButtons[b],
		apiToTcellMods[m],
	)}
}

func (e *mouseEvent) Pos() (int, int) { return e.evt.Position() }

func (e *mouseEvent) Button() api.ButtonMask {
	return tcellToApiButtons[e.evt.Buttons()]
}

func (e *mouseEvent) Mod() api.ModifierMask {
	return tcellToApiMods[e.evt.Modifiers()]
}

func (e *mouseEvent) When() time.Time { return e.evt.When() }

func (e *mouseEvent) Source() interface{} { return e.evt }
