// Copyright (c) 2022 Stephan Lukits. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package lines

import "github.com/slukits/lines/internal/api"

// A MouseEventer is implemented by a reported mouse event.
type MouseEventer = api.MouseEventer

// A ButtonMask mask is reported by a mouse event.
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
