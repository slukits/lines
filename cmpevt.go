// Copyright (c) 2022 Stephan Lukits. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package lines

import (
	"time"
)

// Update posts an update event into the event loop and once it is
// polled it is reported to given listener with the environment for
// given component.  Is given listener is nil the event is to the
// OnUpdate implementation of given component reported.  Update returns
// fails if given component is not initialized or the backend fails.
func (c *Component) Update(data interface{}, l Listener) error {
	return c.bcknd.Post(&UpdateEvent{
		when: time.Now(),
		cmp:  c,
		lst:  l,
		Data: data,
	})
}

// UpdateEvent is created by an Update call on Events.  Its Data field
// provides the data which was passed to that Update call.
type UpdateEvent struct {
	when time.Time
	// NOTE we can not extract the componenter from the component
	// without risking a race condition hence we leave it to the
	// reporter to do so.
	cmp *Component
	lst Listener

	// Data provided to an update event listener
	Data interface{}
}

// When of an update event is set to time.Now()
func (u *UpdateEvent) When() time.Time { return u.when }

func (u *UpdateEvent) Source() interface{} { return u }

// MoveFocus posts a new MoveFocus event into the event loop which once
// it is polled calls the currently focused component's OnFocusLost
// implementation while given component's OnFocus implementation is
// executed.  Finally the focus is set to given component.  MoveFocus
// fails if the event-loop is full returned error will wrap tcell's
// *PostEvent* error.  MoveFocus is an no-op if Componenter is nil.
func (c *Component) Focus() error {
	return c.bcknd.Post(&moveFocusEvent{
		when: time.Now(),
		cmp:  c,
	})
}

// moveFocusEvent is posted by calling MoveFocus for a programmatically
// change of focus.  This event-instance is not provided to the user.
type moveFocusEvent struct {
	when time.Time
	// NOTE we can not extract the componenter from the component
	// without risking a race condition hence we leave it to the
	// reporter to do so.
	cmp *Component
}

func (e *moveFocusEvent) When() time.Time { return e.when }

func (e *moveFocusEvent) Source() interface{} { return e }
