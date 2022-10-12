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

func (s *AnUI) Has_initially_testing_s_width_and_height(t *T) {
	ui, tt := Fixture(t.GoT())
	width, height := ui.Size()
	t.Eq(tt.Width, width)
	t.Eq(tt.Height, height)
}

func (s *AnUI) Reports_a_resize_event_as_first_event(t *T) {
	ui, _ := Fixture(t.GoT())
	c := make(chan api.Eventer)

	go func() { c <- ui.Poll() }()
	select {
	case evt := <-c:
		_, ok := evt.(api.ResizeEventer)
		t.True(ok)
	case <-t.Timeout(0):
		t.Fatal("polling resize timed out")
	}
}

type evtFX struct{}

func (s *evtFX) When() time.Time     { return time.Now() }
func (s *evtFX) Source() interface{} { return s }

func (s *AnUI) Reports_a_posted_event(t *T) {
	ui, _ := Fixture(t.GoT())
	c := make(chan api.Eventer)
	go func() { ui.Poll(); c <- ui.Poll() }()

	t.FatalOn(ui.Post(&evtFX{}))

	select {
	case evt := <-c:
		_, ok := evt.(*evtFX)
		t.True(ok)
		_, ok = evt.Source().(*evtFX)
		t.True(ok)
	case <-t.Timeout(0):
		t.Fatal("polling posted event timed out")
	}
}

func (s *AnUI) Reports_a_key_eventer(t *T) {
	ui, tt := Fixture(t.GoT())
	c := make(chan api.Eventer)
	go func() { ui.Poll(); c <- ui.Poll() }()

	t.FatalOn(tt.PostKey(api.Enter, api.ZeroModifier))

	select {
	case evt := <-c:
		e, ok := evt.(api.KeyEventer)
		t.True(ok)
		t.Eq(api.Enter, e.Key())
		_, ok = e.Source().(*tcell.EventKey)
		t.True(ok)
	case <-t.Timeout(0):
		t.Fatal("polling key-event timed out")
	}
}

func (s *AnUI) Reports_a_rune_eventer(t *T) {
	ui, tt := Fixture(t.GoT())
	c := make(chan api.Eventer)
	go func() { ui.Poll(); c <- ui.Poll() }()

	t.FatalOn(tt.PostRune('x', api.ZeroModifier))

	select {
	case evt := <-c:
		e, ok := evt.(api.RuneEventer)
		t.True(ok)
		t.Eq('x', e.Rune())
		_, ok = e.Source().(*tcell.EventKey)
		t.True(ok)
	case <-t.Timeout(0):
		t.Fatal("polling rune-event timed out")
	}
}

func (s *AnUI) Reports_a_mouse_eventer(t *T) {
	ui, tt := Fixture(t.GoT())
	c := make(chan api.Eventer)
	go func() { ui.Poll(); c <- ui.Poll() }()

	t.FatalOn(tt.PostMouse(4, 2, api.Primary, api.ZeroModifier))

	select {
	case evt := <-c:
		e, ok := evt.(api.MouseEventer)
		t.True(ok)
		t.Eq(api.Primary, e.Button())
		x, y := e.Pos()
		t.Eq(4, x)
		t.Eq(2, y)
		_, ok = e.Source().(*tcell.EventMouse)
		t.True(ok)
	case <-t.Timeout(0):
		t.Fatal("polling mouse-event timed out")
	}
}

func (s *AnUI) Reports_a_resize_eventer(t *T) {
	ui, tt := Fixture(t.GoT())
	c := make(chan api.Eventer)
	go func() { ui.Poll(); c <- ui.Poll() }()

	t.FatalOn(tt.PostResize(42, 22))

	select {
	case evt := <-c:
		e, ok := evt.(api.ResizeEventer)
		t.True(ok)
		width, height := e.Size()
		t.Eq(42, width)
		t.Eq(22, height)
		uiWidth, uiHeight := ui.Size()
		t.Eq(width, uiWidth)
		t.Eq(height, uiHeight)
		_, ok = e.Source().(*tcell.EventResize)
		t.True(ok)
	case <-t.Timeout(0):
		t.Fatal("polling resize-event timed out")
	}
}

func (s *AnUI) Reports_the_nil_event_on_quitting(t *T) {
	ui, _ := Fixture(t.GoT())
	c := make(chan api.Eventer)
	go func() { ui.Poll(); c <- ui.Poll() }()

	ui.Quit()

	select {
	case evt := <-c:
		if evt != nil {
			go func() { c <- ui.Poll() }()
			select { // there is some-times an other resize ???
			case evt := <-c:
				t.True(evt == nil)
			case <-t.Timeout(0):
				t.Fatal("polling nil-event timed out")
			}
		}
	case <-t.Timeout(0):
		t.Fatal("polling nil-event timed out")
	}
}

func (s *AnUI) Displays_given_rune_at_given_position(t *T) {
	ui, tt := Fixture(t.GoT())
	c := make(chan api.Eventer)
	go func() { ui.Poll(); c <- ui.Poll() }()
	tt.PostResize(3, 3)
	select {
	case <-c:
	case <-t.Timeout(0):
		t.Fatal("resize-event timed out")
	}

	ui.Display(1, 1, 'x', api.Style{})
	ui.Redraw()
	t.Eq("x", tt.Screen().Trimmed().String())
	t.Eq("x", tt.ScreenArea(1, 1, 1, 1).String())
}

func TestAnUI(t *testing.T) {
	t.Parallel()
	Run(&AnUI{}, t)
}
