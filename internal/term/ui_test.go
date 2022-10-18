// Copyright (c) 2022 Stephan Lukits. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package term

import (
	"testing"
	"time"

	"github.com/gdamore/tcell/v2"
	. "github.com/slukits/gounit"
	"github.com/slukits/lines/internal/api"
)

type AnUI struct{ Suite }

func (s *AnUI) SetUp(t *T) { t.Parallel() }

func (s *AnUI) Has_initially_testing_s_width_and_height(t *T) {
	ui, tt := LstFixture(t.GoT(), nil, 0)
	width, height := ui.Size()
	t.Eq(tt.Width, width)
	t.Eq(tt.Height, height)
}

func (s *AnUI) Reports_a_resize_event_as_first_event(t *T) {
	eventReceived := false

	LstFixture(t.GoT(), func(evt api.Eventer) {
		if _, ok := evt.(api.ResizeEventer); ok {
			eventReceived = true
		} else {
			t.True(eventReceived)
		}
	}, 0)

	t.True(eventReceived)
}

func (s *AnUI) Reports_a_quit_event_when_quitting(t *T) {
	quitted := false
	ui, _ := LstFixture(t.GoT(), nil, 0)
	ui.OnQuit(func() { quitted = true })

	ui.Quit()
	t.True(quitted)
}

type evtFX struct{}

func (s *evtFX) When() time.Time     { return time.Now() }
func (s *evtFX) Source() interface{} { return s }

func (s *AnUI) Reports_a_posted_event(t *T) {
	eventReceived := false
	ui, _ := LstFixture(t.GoT(), func(e api.Eventer) {
		_, eventReceived = e.(*evtFX)
	}, 0)

	t.FatalOn(ui.Post(&evtFX{}))

	t.True(eventReceived)
}

func (s *AnUI) Reports_a_key_eventer(t *T) {
	eventReceived := false
	_, tt := LstFixture(t.GoT(), func(evt api.Eventer) {
		if e, ok := evt.(api.KeyEventer); ok {
			t.Eq(api.Enter, e.Key())
			_, ok = e.Source().(*tcell.EventKey)
			t.True(ok)
			eventReceived = true
		}
	}, 0)

	tt.PostKey(api.Enter, api.ZeroModifier)

	t.True(eventReceived)
}

func (s *AnUI) Reports_a_rune_eventer(t *T) {
	eventReceived := false
	_, tt := LstFixture(t.GoT(), func(evt api.Eventer) {
		if e, ok := evt.(api.RuneEventer); ok {
			t.Eq('x', e.Rune())
			_, ok = e.Source().(*tcell.EventKey)
			t.True(ok)
			eventReceived = true
		}
	}, 0)

	tt.PostRune('x', api.ZeroModifier)

	t.True(eventReceived)
}

func (s *AnUI) Reports_a_mouse_eventer(t *T) {
	eventReceived := false
	_, tt := LstFixture(t.GoT(), func(evt api.Eventer) {
		if e, ok := evt.(api.MouseEventer); ok {
			t.Eq(api.Primary, e.Button())
			x, y := e.Pos()
			t.Eq(4, x)
			t.Eq(2, y)
			_, ok = e.Source().(*tcell.EventMouse)
			t.True(ok)
			eventReceived = true
		}
	}, 0)

	tt.PostMouse(4, 2, api.Primary, api.ZeroModifier)

	t.True(eventReceived)
}

func (s *AnUI) Reports_a_resize_eventer(t *T) {
	reported := 0
	ui, tt := LstFixture(t.GoT(), func(evt api.Eventer) {
		if e, ok := evt.(api.ResizeEventer); ok {
			reported++
			if reported != 2 {
				return
			}
			width, height := e.Size()
			t.Eq(42, width)
			t.Eq(22, height)
			_, ok = e.Source().(*tcell.EventResize)
			t.True(ok)
		}
	}, 0)

	tt.PostResize(42, 22)

	t.FatalIfNot(t.Eq(2, reported))
	uiWidth, uiHeight := ui.Size()
	t.Eq(42, uiWidth)
	t.Eq(22, uiHeight)
}

func (s *AnUI) Displays_given_rune_at_given_position(t *T) {
	ui, tt := LstFixture(t.GoT(), nil, 10*time.Minute)
	tt.PostResize(3, 3)

	ui.Display(1, 1, 'x', api.Style{})
	ui.Redraw()
	t.Eq("x", tt.Screen().Trimmed().String())
	t.Eq("x", tt.ScreenArea(1, 1, 1, 1).String())
}

func TestAnUI(t *testing.T) {
	t.Parallel()
	Run(&AnUI{}, t)
}
