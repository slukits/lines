// Copyright (c) 2022 Stephan Lukits. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

/*
ui implements an UIer wrapping tcell for lines terminal ui.
*/

package term

import (
	"fmt"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/slukits/lines/internal/api"
)

type UI struct {

	// lib the tcell terminal screen which is the simulation screen in
	// case of testing
	lib tcell.Screen

	// styler avoids unnecessary api-style conversions if a sequence of
	// runes is requested to be displayed with the same style
	styler func(api.Style) tcell.Style

	// hasQuit indicates that UI.Quit was already called on an ui
	// instance to avoid two calls of u.lib.Fini() which panics at its
	// second call.
	hasQuit bool
}

func New() *UI {
	lib, err := tcell.NewScreen()
	if err != nil {
		panic(fmt.Sprintf(
			"lines: term: new: can't obtain screen: %v", err))
	}
	return initUI(lib)
}

func initUI(lib tcell.Screen) *UI {
	if err := lib.Init(); err != nil {
		panic(fmt.Sprintf(
			"lines: term: new: can't obtain screen: %v", err))
	}
	return &UI{lib: lib, styler: apiToTcellStyleClosure()}
}

// Size returns the ui's screen size.
func (u *UI) Size() (int, int) { return u.lib.Size() }

func (u *UI) Quit() {
	if u.hasQuit {
		return
	}
	u.hasQuit = true
	u.lib.Fini()
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

func (u *UI) Poll() api.Eventer {
	evt := u.lib.PollEvent()
	switch evt := evt.(type) {
	case *tcell.EventResize:
		return &resize{evt: evt}
	case nil:
		return nil
	default:
		return evt.(api.Eventer)
	}
}

func (u *UI) Post(evt api.Eventer) error {
	if err := u.lib.PostEvent(evt); err != nil {
		return err
	}
	return nil
}

func (u *UI) Display(x, y int, r rune, s api.Style) {
	u.lib.SetContent(x, y, r, nil, u.styler(s))
}

func (u *UI) Redraw() { u.lib.Sync() }

func (u *UI) Update() { u.lib.Show() }
