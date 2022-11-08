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

// Clicker is implemented by components which want to be informed about
// a "left"-mouse click event in their printable area.  If the clicked
// component, i.e. the component with the smallest layout area
// containing the event coordinates, does not have the focus an OnFocus
// event is reported first if and only if the clicked component has the
// Focusable feature.  See [Mouser] event interface for a more general
// mouse event handling.
type Clicker interface {

	// OnClick gets "left click"-events reported.  x and y provide the
	// click coordinates translated into the layouted area of the
	// receiving component.  E.g. y == 3 means that the component's
	// third line was clicked.  This event bubbles use e.StopBubbling()
	// to suppress further bubbling.  Note e.Evt.Source() provides the
	// event object reported by the backend.
	OnClick(e *Env, x, y int)
}

// Contexter is implemented by components which want to be informed
// about a mouse "right click"-event in their printable area.  If the
// clicked component, i.e. the component with the smallest layout area
// containing the event coordinates, does not have the focus an OnFocus
// event is reported first if and only if the clicked component has the
// Focusable feature.  See [Mouser] event interface for a more general
// mouse event handling.
//
// TODO: implement: see if event can also be reported for a potential
// context-menu key press (having x/y set to -1 then?).
type Contexter interface {

	// OnContext gets "right click"-events reported whereas provided x
	// and y mouse-coordinates are translated into the printable area of
	// the receiving component.  E.g. y == 3 means that the component's
	// third line was clicked.  This event bubbles use e.StopBubbling()
	// to suppress further bubbling.  Note e.Evt.Source() provides the
	// event object reported by the backend.
	OnContext(e *Env, x, y int)
}

// Mouser is implemented by components who want to be informed about all
// mouse event in their printable area.  If the clicked component, i.e.
// the component with the smallest layout area containing the event
// coordinates, does not have the focus an OnFocus event is reported
// first if and only if the clicked component has the Focusable feature.
type Mouser interface {

	// OnMouse gets any mouse event reported whereas provided x and y
	// mouse-coordinates are translated into the printable area of the
	// receiving component.  E.g. y == 3 means that the component's
	// third line was clicked.  This event bubbles use e.StopBubbling()
	// to suppress further bubbling.  Note e.Evt.Source() provides the
	// event object reported by the backend.
	OnMouse(e *Env, _ ButtonMask, x, y int)
}
