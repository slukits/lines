// Copyright (c) 2022 Stephan Lukits. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

/*
ui implements an UIer wrapping tcell for lines terminal ui.
*/

package term

import (
	"fmt"
	"sync"
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

	defaultStyle api.Style

	// styler avoids unnecessary api-style conversions if a sequence of
	// runes is requested to be displayed with the same style
	styler func(api.Style) tcell.Style

	// hasQuit indicates that UI.Quit was already called on an ui
	// instance to avoid two calls of u.lib.Fini() which panics at its
	// second call.
	hasQuit atomic.Bool

	waitForQuit chan struct{}

	transactional atomic.Value

	*sync.Mutex

	quitter []func()
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
		lib:          lib,
		Mutex:        &sync.Mutex{},
		defaultStyle: api.DefaultStyle,
		styler:       apiToTcellStyleClosure(),
		waitForQuit:  make(chan struct{}),
		listener:     l,
	}
	go ui.poll()
	return ui
}

// WaitForQuit returns a channel which is closed if the event-loop is
// quit.
func (u *UI) WaitForQuit() {
	<-u.waitForQuit
}

// OnQuit registers given function to be called on quitting.
func (u *UI) OnQuit(listener func()) {
	u.Lock()
	defer u.Unlock()
	u.quitter = append(u.quitter, listener)
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

// Quit polling and reporting events, inform all listeners about about
// it and reset the terminal screen.
func (u *UI) Quit() {
	if !u.hasQuit.CompareAndSwap(false, true) {
		return
	}
	u.Lock()
	defer u.Unlock()
	for _, l := range u.quitter {
		l()
	}
	u.lib.Fini()
	select {
	case <-u.waitForQuit:
		// non-blocking i.e. must be closed already
	default:
		close(u.waitForQuit)
	}
}

func (u *UI) poll() {
	for {
		evt := u.lib.PollEvent()
		if evt == nil {
			return
		}
		if evt, ok := evt.(*screenEvent); ok {
			u.handleScreenEvent(evt)
			continue
		}
		if u.listener == nil {
			u.handleTransactional()
			continue
		}
		switch evt := evt.(type) {
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
			e, ok := evt.(api.Eventer)
			if !ok {
				panic(fmt.Sprintf("unknown event type: %T", evt))
			}
			u.listener(e)
		}
		u.handleTransactional()
	}
}

func (u *UI) handleTransactional() {
	if u.transactional.Load() != nil {
		u.transactional.Load().(*transactional).polled()
	}
}

func (u *UI) handleScreenEvent(evt *screenEvent) {
	evt.grab()
	u.handleTransactional()
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

// NewStyle returns a new style corresponding to tcell's default style.
func (u *UI) NewStyle() api.Style {
	return u.defaultStyle
}
