// Copyright (c) 2022 Stephan Lukits. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package api

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
