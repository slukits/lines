// Copyright (c) 2022 Stephan Lukits. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package api_test

import (
	"testing"

	. "github.com/slukits/gounit"
	"github.com/slukits/lines"
	"github.com/slukits/lines/internal/api"
	"github.com/slukits/lines/internal/term"
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
	ui, tt := term.LstFixture(t.GoT(), nil, 0)
	tt.PostResize(16, 5)

	tt.Display(fx, ui.NewStyle())
	ui.Redraw()
	t.Eq(fx, tt.Screen().String())
}

func (s *_testing) fx(t *T) (*term.UI, *term.Fixture) {
	ui, tt := term.LstFixture(t.GoT(), nil, 0)
	tt.PostResize(16, 5)
	return ui, tt
}

func (s *_testing) Reports_trimmed_string_representation_of_screen(t *T) {
	ui, tt := s.fx(t)

	tt.Display(fx, ui.NewStyle())
	ui.Redraw()
	t.Eq(exp, tt.Screen().Trimmed().String())
}

func (s *_testing) Reports_string_of_given_screen_area(t *T) {
	ui, tt := s.fx(t)

	tt.Display(fx, ui.NewStyle())
	ui.Redraw()
	t.Eq(exp, tt.ScreenArea(3, 1, 13, 3).String())
}

func (s *_testing) Reports_cells_of_screen(t *T) {
	ui, tt := s.fx(t)

	tt.Display(fx, ui.NewStyle())
	ui.Redraw()
	t.Eq(fx, tt.Cells().String())
}

func (s *_testing) Reports_trimmed_cells_of_screen(t *T) {
	ui, tt := s.fx(t)

	tt.Display(fx, ui.NewStyle())
	ui.Redraw()
	t.Eq(exp, tt.Cells().Trimmed().String())
}

func (s *_testing) Reports_cells_of_screen_area(t *T) {
	ui, tt := s.fx(t)

	tt.Display(fx, ui.NewStyle())
	ui.Redraw()
	str := tt.CellsArea(3, 1, 13, 3).String()
	t.Eq(exp, str)
}

func (s *_testing) Reports_style_information_falsy_if_out_of_bound(t *T) {
	ui, tt := s.fx(t)
	sty := api.NewStyle(api.Italic, lines.Green, lines.Silver)
	tt.Display(fx, sty)
	ui.Redraw()
	scr := tt.Cells().Trimmed()

	t.Not.True(
		scr.HasFG(0, -1, lines.Green) ||
			scr.HasFG(0, len(scr), lines.Green) ||
			scr.HasFG(-1, 0, lines.Green) ||
			scr.HasFG(len(scr[0]), 0, lines.Green) ||
			scr.HasBG(0, -1, lines.Silver) ||
			scr.HasBG(-1, 0, lines.Silver) ||
			scr.HasAA(0, -1, api.Italic) ||
			scr.HasAA(-1, 0, api.Italic),
	)
}

func (s *_testing) Reports_style_information_of_screen_cells(t *T) {
	ui, tt := s.fx(t)
	sty := api.NewStyle(api.Italic, lines.Green, lines.Silver)
	tt.Display(fx, sty)
	ui.Redraw()
	scr := tt.Cells().Trimmed()

	t.True(
		scr.HasFG(0, 0, lines.Green) &&
			scr.HasBG(0, 0, lines.Silver) &&
			scr.HasAA(0, 0, api.Italic),
	)
}

func TestTesting(t *testing.T) {
	t.Parallel()
	Run(&_testing{}, t)
}
