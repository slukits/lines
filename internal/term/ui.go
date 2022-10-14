// Copyright (c) 2022 Stephan Lukits. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

/*
ui implements an UIer wrapping tcell for lines terminal ui.
*/

package term

import (
	"fmt"
	"sync/atomic"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/slukits/lines/internal/api"
)

type UI struct {

	// listener is informed about new events.
	listener func(api.Eventer)

	// lib the tcell terminal screen which is the simulation screen in
	// case of testing
	lib tcell.Screen

	// styler avoids unnecessary api-style conversions if a sequence of
	// runes is requested to be displayed with the same style
	styler func(api.Style) tcell.Style

	// hasQuit indicates that UI.Quit was already called on an ui
	// instance to avoid two calls of u.lib.Fini() which panics at its
	// second call.
	hasQuit atomic.Bool

	waitForQuit chan struct{}

	transactional atomic.Value
}

func New(listener func(api.Eventer)) *UI {
	lib, err := tcell.NewScreen()
	if err != nil {
		panic(fmt.Sprintf(
			"lines: term: new: can't obtain screen: %v", err))
	}

	return initUI(lib, listener)
}

func (u *UI) Lib() interface{} { return u.lib }

func initUI(lib tcell.Screen, l func(api.Eventer)) *UI {
	if err := lib.Init(); err != nil {
		panic(fmt.Sprintf(
			"lines: term: new: can't obtain screen: %v", err))
	}
	lib.EnableMouse()
	lib.EnablePaste()
	ui := &UI{
		lib:         lib,
		styler:      apiToTcellStyleClosure(),
		waitForQuit: make(chan struct{}),
		listener:    l,
	}
	go ui.poll()
	return ui
}

// WaitForQuit returns a channel which is closed if the event-loop is
// quit.
func (u *UI) WaitForQuit() {
	<-u.waitForQuit
}

// EnableEventTransactions guarantees that an event-post p not returns
// before all other posted events during the processing of p have been
// processed.
func (u *UI) EnableTransactionalEventPosts(timeout time.Duration) {
	u.transactional.CompareAndSwap(nil, &transactional{
		ui: u, timeout: timeout, waiting: make(chan bool)})
}

// Size returns the ui's screen size.
func (u *UI) Size() (int, int) { return u.lib.Size() }

func (u *UI) Quit() {
	if !u.hasQuit.CompareAndSwap(false, true) {
		return
	}
	if u.listener != nil {
		u.listener(&QuitEvent{})
	}
	u.lib.Fini()
	close(u.waitForQuit)
}

type QuitEvent struct{}

func (q *QuitEvent) Quitting()           {}
func (q *QuitEvent) Source() interface{} { return q }
func (q *QuitEvent) When() time.Time     { return time.Now() }

func (u *UI) poll() {
	for {
		evt := u.lib.PollEvent()
		if evt == nil {
			return
		}
		if u.listener == nil {
			if u.transactional.Load() != nil {
				u.transactional.Load().(*transactional).polled()
			}
			continue
		}
		switch evt := evt.(type) {
		case api.Eventer:
			u.listener(evt)
		case *tcell.EventResize:
			u.listener(&resize{evt: evt})
		case *tcell.EventKey:
			if evt.Key() == tcell.KeyRune {
				u.listener(&runeEvent{evt: evt})
				break
			}
			u.listener(&keyEvent{evt: evt})
		case *tcell.EventMouse:
			u.listener(&mouseEvent{evt: evt})
		case *tcell.EventPaste:
			u.listener(&bracketPaste{evt: evt})
		default:
			panic(fmt.Sprintf("unknown event type: %T", evt))
		}
		if u.transactional.Load() != nil {
			u.transactional.Load().(*transactional).polled()
		}
	}
}

type resize struct{ evt *tcell.EventResize }

func newResize(width, height int) api.ResizeEventer {
	return &resize{evt: tcell.NewEventResize(
		width, height,
	)}
}

func (r *resize) Source() interface{} { return r.evt }
func (r *resize) When() time.Time     { return r.evt.When() }
func (r *resize) Size() (int, int)    { return r.evt.Size() }

type bracketPaste struct{ evt *tcell.EventPaste }

func newBracketPaste(start bool) *bracketPaste {
	return &bracketPaste{evt: tcell.NewEventPaste(start)}
}

func (r *bracketPaste) Source() interface{} { return r.evt }
func (r *bracketPaste) When() time.Time     { return r.evt.When() }
func (r *bracketPaste) Start() bool         { return r.evt.Start() }
func (r *bracketPaste) End() bool           { return r.evt.End() }

func (u *UI) Post(evt api.Eventer) error {
	if u.transactional.Load() == nil {
		return u.lib.PostEvent(evt)
	}
	return u.transactional.Load().(*transactional).Post(evt)
}

func (u *UI) Display(x, y int, r rune, s api.Style) {
	u.lib.SetContent(x, y, r, nil, u.styler(s))
}

func (u *UI) Redraw() { u.lib.Sync() }

func (u *UI) Update() { u.lib.Show() }
