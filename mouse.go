// Copyright (c) 2022 Stephan Lukits. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package lines

import "github.com/slukits/lines/internal/api"

// A MouseEventer is implemented by a reported mouse event.
type MouseEventer = api.MouseEventer

// A Button mask is reported by a mouse event.
type Button = api.Button

const (
	Button1    Button = api.Button1 // Usually the left (primary) mouse button.
	Button2    Button = api.Button2 // Usually the right (secondary) mouse button.
	Button3    Button = api.Button3 // Usually the middle mouse button.
	Button4    Button = api.Button4 // Often a side button (thumb/next).
	Button5    Button = api.Button5 // Often a side button (thumb/prev).
	Button6    Button = api.Button6
	Button7    Button = api.Button7
	Button8    Button = api.Button8
	WheelUp    Button = api.WheelUp    // Wheel motion up/away from user.
	WheelDown  Button = api.WheelDown  // Wheel motion down/towards user.
	WheelLeft  Button = api.WheelLeft  // Wheel motion to left.
	WheelRight Button = api.WheelRight // Wheel motion to right.
	ZeroButton Button = api.ZeroButton // No button or wheel events.

	Primary   Button = Button1
	Secondary Button = Button2
	Middle    Button = Button3
)
