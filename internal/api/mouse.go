// Copyright (c) 2022 Stephan Lukits. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package api

import "time"

// ButtonMask mask is used to report/post user mouse input events.
type ButtonMask int32

const (
	Button1 ButtonMask = 1 << iota // Usually the left (primary) mouse button.
	Button2                        // Usually the right (secondary) mouse button.
	Button3                        // Usually the middle mouse button.
	Button4                        // Often a side button (thumb/next).
	Button5                        // Often a side button (thumb/prev).
	Button6
	Button7
	Button8
	WheelUp                   // Wheel motion up/away from user.
	WheelDown                 // Wheel motion down/towards user.
	WheelLeft                 // Wheel motion to left.
	WheelRight                // Wheel motion to right.
	ZeroButton ButtonMask = 0 // No button or wheel events.

	Primary   = Button1
	Secondary = Button2
	Middle    = Button3
)

// aggregatedMouseEvent combines several mouse events to a single event
// for convenience, e.g. an event with button b != 0 followed by two
// events  0-button events at the same position are reported as a click.
type aggregatedMouseEvent struct {
	x, y int
	b    ButtonMask
	m    ModifierMask
	when time.Time
}

func (e *aggregatedMouseEvent) Pos() (int, int) { return e.x, e.y }

func (e *aggregatedMouseEvent) Button() ButtonMask { return e.b }

func (e *aggregatedMouseEvent) Mod() ModifierMask { return e.m }

func (e *aggregatedMouseEvent) When() time.Time { return e.when }

func (e *aggregatedMouseEvent) Source() interface{} { return nil }

// MouseClick is an MouseEventer implementation reporting the
// aggregation of mouse events as a mouse click of arbitrary button.
type MouseClick struct {
	aggregatedMouseEvent
}

func NewMouseClick(to MouseEventer) *MouseClick {
	x, y := to.Pos()
	return &MouseClick{aggregatedMouseEvent: aggregatedMouseEvent{
		x: x, y: y, b: to.Button(), m: to.Mod(), when: to.When()},
	}
}

type MouseMove struct {
	aggregatedMouseEvent
	ox, oy int
}

func NewMouseMove(ox, oy int, to MouseEventer) *MouseMove {
	x, y := to.Pos()
	return &MouseMove{
		aggregatedMouseEvent: aggregatedMouseEvent{
			x: x, y: y, b: to.Button(), m: to.Mod(), when: to.When()},
		ox: ox, oy: oy,
	}
}

// Origin reports coordinates from where the mouse move started.  NOTE
// these coordinates are absolute coordinates since the movement might
// started outside the component.
func (d *MouseMove) Origin() (x, y int) { return d.ox, d.oy }

type MouseDrag struct {
	aggregatedMouseEvent
	ox, oy int
}

func NewMouseDrag(ox, oy int, to MouseEventer) *MouseDrag {
	x, y := to.Pos()
	return &MouseDrag{
		aggregatedMouseEvent: aggregatedMouseEvent{
			x: x, y: y, b: to.Button(), m: to.Mod(), when: to.When()},
		ox: ox, oy: oy,
	}
}

// Origin returns the absolute coordinates where the drag started.
func (d *MouseDrag) Origin() (x, y int) { return d.ox, d.oy }

type MouseDrop struct {
	aggregatedMouseEvent
}

func NewMouseDrop(to MouseEventer) *MouseDrop {
	x, y := to.Pos()
	return &MouseDrop{aggregatedMouseEvent: aggregatedMouseEvent{
		x: x, y: y, b: to.Button(), m: to.Mod(), when: to.When()},
	}
}
