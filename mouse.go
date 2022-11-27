// Copyright (c) 2022 Stephan Lukits. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package lines

import "github.com/slukits/lines/internal/api"

// A MouseEventer is implemented by a reported mouse event.  Mouse
// events may be received from a component by implementing the [Mouser]
// interface.  MouseEventer.Source() provides the backend event
// triggering the mouse event.
type MouseEventer = api.MouseEventer

// A ButtonMask mask is reported by a mouse event to a [Mouser]
// implementing component.
type ButtonMask = api.ButtonMask

const (
	Button1    ButtonMask = api.Button1 // Usually the left (primary) mouse button.
	Button2    ButtonMask = api.Button2 // Usually the right (secondary) mouse button.
	Button3    ButtonMask = api.Button3 // Usually the middle mouse button.
	Button4    ButtonMask = api.Button4 // Often a side button (thumb/next).
	Button5    ButtonMask = api.Button5 // Often a side button (thumb/prev).
	Button6    ButtonMask = api.Button6
	Button7    ButtonMask = api.Button7
	Button8    ButtonMask = api.Button8
	WheelUp    ButtonMask = api.WheelUp    // Wheel motion up/away from user.
	WheelDown  ButtonMask = api.WheelDown  // Wheel motion down/towards user.
	WheelLeft  ButtonMask = api.WheelLeft  // Wheel motion to left.
	WheelRight ButtonMask = api.WheelRight // Wheel motion to right.
	ZeroButton ButtonMask = api.ZeroButton // No button or wheel events.

	Primary   ButtonMask = Button1
	Secondary ButtonMask = Button2
	Middle    ButtonMask = Button3
)

// Mover is implemented by components which want to be informed about
// a mouse movement events in(to) their printable area.
type Mover interface {

	// OnMove implementation of a component c gets mouse movements
	// reported whereas provided coordinates x and y of the movement are
	// translated into the "screen area" of c.
	OnMove(e *Env, x, y int)
}

// Clicker is implemented by components which want to be informed about
// a "left"-mouse click event in their "screen area".  If the clicked
// component, i.e. the component with the smallest layout area
// containing the event coordinates, does not have the focus an OnFocus
// event is reported first if and only if the clicked component has the
// Focusable feature.  See [Mouser] event interface for a more general
// mouse event handling.
type Clicker interface {

	// OnClick implementation of a component c gets "left click"-events
	// reported which is an aggregated mouse event.  I.e. Mouser
	// implementer will not receive a "left click" event.  x and y
	// provide the click coordinates translated into the "screen area"
	// (c.Dim().Rect()) of c.  This event bubbles; use e.StopBubbling()
	// to suppress further bubbling.  Note e.Evt.(*lines.MouseClick).Mod
	// provides also the modifiers information.
	OnClick(e *Env, x, y int)
}

// Contexter is implemented by components which want to be informed
// about a mouse "right click" event in their screen area.  If the
// clicked component, i.e. the component with the smallest layout area
// containing the event coordinates, does not have the focus an OnFocus
// event is reported first if and only if the clicked component has the
// Focusable feature.  See [Mouser] event interface for a more general
// mouse event handling.
//
// TODO: implement: see if event can also be reported for a potential
// context-menu key press (having x/y set to -1 then?).
type Contexter interface {

	// OnContext implementation of a component c gets "right click"
	// events reported which is an aggregated mouse event.  I.e. Mouser
	// implementer will not receive a "right click" event.  x and y
	// provide the click coordinates translated into the "screen area"
	// (c.Dim().Rect()) of c.  This event bubbles; use e.StopBubbling()
	// to suppress further bubbling.  Note e.Evt.(*lines.MouseClick).Mod
	// provides also the modifiers information.
	OnContext(e *Env, x, y int)
}

// Drager is implemented by a component which wants to be informed about
// mouse movement while a button is pressed.  NOTE reported coordinates
// are absolute coordinates.
type Drager interface {
	// OnDrag implementation of a component c gets move movement
	// reported while a button is pressed.  Given environment provides
	// the origin of a drag e.Evt(*lines.MouseDrag).Origin().  Reported
	// coordinates are absolute coordinates.
	OnDrag(e *Env, _ ButtonMask, x, y int)
}

// Dropper is implemented by a component which wants to be informed when
// a (sequence of) drags ends.
type Dropper interface {
	// OnDrop implementation of a component c gets reported the absolute
	// position where a sequence of mouse movements ended while given
	// button was pressed.
	OnDrop(_ *Env, _ ButtonMask, x, y int)
}

// Mouser is implemented by components who want to be informed about all
// mouse event in their "screen area".  Note no implicit focusing
// happens no matter what the event.
type Mouser interface {

	// OnMouse implementation of a component c gets any mouse event
	// reported.  x and y provide the click coordinates translated into
	// the "screen area" (c.Dim().Rect()) of c.  Mouse events bubble; use
	// e.StopBubbling() to suppress further bubbling.  Note
	// e.Evt.(MouseEventer).Mod() provides also the modifiers information and
	// e.Evt.Source() provides the event object reported by the backend.
	OnMouse(e *Env, _ ButtonMask, x, y int)
}

// Enterer is implemented by components which want to be informed if the
// mouse pointer enters their printable area.
type Enterer interface {

	// OnEnter implementation of a component c gets the first mouse move
	// inside c's printable area reported.
	OnEnter(e *Env)
}

// Exiter is implemented by components which want to be informed if the
// mouse pointer exits their printable area.
type Exiter interface {

	// OnExit implementation of a component c gets the first mouse move
	// outside c's printable area reported.
	OnExit(e *Env)
}

// MouseClick is an MouseEventer implementation reporting the
// aggregation of mouse events as a mouse click of arbitrary button.
// Implementing the Clicker interface one gets clicks of the primary
// button reported.  Implementing the Contexter interface one gets
// clicks of the secondary button reported.
type MouseClick = api.MouseClick

type MouseMove = api.MouseMove

type MouseDrag = api.MouseDrag

type MouseDrop = api.MouseDrop
