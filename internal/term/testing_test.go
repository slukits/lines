// Copyright (c) 2022 Stephan Lukits. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package term

import (
	"testing"

	. "github.com/slukits/gounit"
	"github.com/slukits/lines/internal/api"
)

type _testing struct{ Suite }

const fx = `                
   upper left   
           width
     bottom     
                `
const exp = `upper left   
        width
  bottom     `

func (s *_testing) Reports_string_representation_of_screen(t *T) {
	ui, tt := Fixture(t.GoT())
	c := make(chan api.Eventer)
	go func() { ui.Poll(); c <- ui.Poll() }()
	tt.PostResize(16, 5)
	select {
	case <-c:
	case <-t.Timeout(0):
		t.Fatal("initial resize-event timed out")
	}
	tt.Display(fx, api.Style{})
	ui.Redraw()
	t.Eq(fx, tt.Screen().String())
}

func (s *_testing) Reports_trimmed_string_representation_of_screen(t *T) {
	ui, tt := Fixture(t.GoT())
	c := make(chan api.Eventer)
	go func() { c <- ui.Poll() }()
	select {
	case <-c:
	case <-t.Timeout(0):
		t.Fatal("initial resize-event timed out")
	}

	tt.Display(fx, api.Style{})
	ui.Redraw()
	t.Eq(exp, tt.Screen().Trimmed().String())
}

func (s *_testing) Reports_string_of_given_screen_area(t *T) {
	ui, tt := Fixture(t.GoT())
	c := make(chan api.Eventer)
	go func() { c <- ui.Poll() }()
	select {
	case <-c:
	case <-t.Timeout(0):
		t.Fatal("initial resize-event timed out")
	}

	tt.Display(fx, api.Style{})
	ui.Redraw()
	t.Eq(exp, tt.ScreenArea(3, 1, 13, 3).String())
}

func (s *_testing) Reports_cells_of_screen(t *T) {
	ui, tt := Fixture(t.GoT())
	c := make(chan api.Eventer)
	go func() { ui.Poll(); c <- ui.Poll() }()
	tt.PostResize(16, 5)
	select {
	case <-c:
	case <-t.Timeout(0):
		t.Fatal("initial resize-event timed out")
	}
	tt.Display(fx, api.Style{})
	ui.Redraw()
	t.Eq(fx, tt.Cells().String())
}

func (s *_testing) Reports_trimmed_cells_of_screen(t *T) {
	ui, tt := Fixture(t.GoT())
	c := make(chan api.Eventer)
	go func() { c <- ui.Poll() }()
	select {
	case <-c:
	case <-t.Timeout(0):
		t.Fatal("initial resize-event timed out")
	}

	tt.Display(fx, api.Style{})
	ui.Redraw()
	t.Eq(exp, tt.Cells().Trimmed().String())
}

func (s *_testing) Reports_cells_of_screen_area(t *T) {
	ui, tt := Fixture(t.GoT())
	c := make(chan api.Eventer)
	go func() { c <- ui.Poll() }()
	select {
	case <-c:
	case <-t.Timeout(0):
		t.Fatal("initial resize-event timed out")
	}

	tt.Display(fx, api.Style{})
	ui.Redraw()
	str := tt.CellsArea(3, 1, 13, 3).String()
	t.Eq(exp, str)
}

func TestTesting(t *testing.T) {
	t.Parallel()
	Run(&_testing{}, t)
}
