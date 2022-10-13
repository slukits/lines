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
	ui, tt := LstFixture(t.GoT(), nil, 0)
	t.FatalOn(tt.PostResize(16, 5))

	tt.Display(fx, api.Style{})
	ui.Redraw()
	t.Eq(fx, tt.Screen().String())
}

func (s *_testing) Reports_trimmed_string_representation_of_screen(t *T) {
	ui, tt := LstFixture(t.GoT(), nil, 0)

	tt.Display(fx, api.Style{})
	ui.Redraw()
	t.Eq(exp, tt.Screen().Trimmed().String())
}

func (s *_testing) Reports_string_of_given_screen_area(t *T) {
	ui, tt := LstFixture(t.GoT(), nil, 0)

	tt.Display(fx, api.Style{})
	ui.Redraw()
	t.Eq(exp, tt.ScreenArea(3, 1, 13, 3).String())
}

func (s *_testing) Reports_cells_of_screen(t *T) {
	ui, tt := LstFixture(t.GoT(), nil, 0)
	tt.PostResize(16, 5)

	tt.Display(fx, api.Style{})
	ui.Redraw()
	t.Eq(fx, tt.Cells().String())
}

func (s *_testing) Reports_trimmed_cells_of_screen(t *T) {
	ui, tt := LstFixture(t.GoT(), nil, 0)

	tt.Display(fx, api.Style{})
	ui.Redraw()
	t.Eq(exp, tt.Cells().Trimmed().String())
}

func (s *_testing) Reports_cells_of_screen_area(t *T) {
	ui, tt := LstFixture(t.GoT(), nil, 0)

	tt.Display(fx, api.Style{})
	ui.Redraw()
	str := tt.CellsArea(3, 1, 13, 3).String()
	t.Eq(exp, str)
}

func (s *_testing) Returns_from_evt_post_after_sub_posts_processed(t *T) {
	var ui *UI
	runeReported := false
	ui, tt := LstFixture(t.GoT(), func(evt api.Eventer) {
		switch evt.(type) {
		case api.KeyEventer:
			ui.Post(newRuneEvent('r', api.ZeroModifier))
		case api.RuneEventer:
			runeReported = true
		}
	}, 0)

	tt.PostKey(api.Enter, api.ZeroModifier)

	t.True(runeReported)
}

func TestTesting(t *testing.T) {
	t.Parallel()
	Run(&_testing{}, t)
}
