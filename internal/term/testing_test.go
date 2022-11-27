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
	tt.PostResize(16, 5)

	tt.Display(fx, ui.NewStyle())
	ui.Redraw()
	t.Eq(fx, tt.Screen().String())
}

func (s *_testing) Reports_trimmed_string_representation_of_screen(t *T) {
	ui, tt := LstFixture(t.GoT(), nil, 0)

	tt.Display(fx, ui.NewStyle())
	ui.Redraw()
	t.Eq(exp, tt.Screen().Trimmed().String())
}

func (s *_testing) Reports_string_of_given_screen_area(t *T) {
	ui, tt := LstFixture(t.GoT(), nil, 0)

	tt.Display(fx, ui.NewStyle())
	ui.Redraw()
	t.Eq(exp, tt.ScreenArea(3, 1, 13, 3).String())
}

func (s *_testing) Reports_cells_of_screen(t *T) {
	ui, tt := LstFixture(t.GoT(), nil, 0)
	tt.PostResize(16, 5)

	tt.Display(fx, ui.NewStyle())
	ui.Redraw()
	t.Eq(fx, tt.Cells().String())
}

func (s *_testing) Reports_trimmed_cells_of_screen(t *T) {
	ui, tt := LstFixture(t.GoT(), nil, 0)

	tt.Display(fx, ui.NewStyle())
	ui.Redraw()
	t.Eq(exp, tt.Cells().Trimmed().String())
}

func (s *_testing) Reports_cells_of_screen_area(t *T) {
	ui, tt := LstFixture(t.GoT(), nil, 0)

	tt.Display(fx, ui.NewStyle())
	ui.Redraw()
	str := tt.CellsArea(3, 1, 13, 3).String()
	t.Eq(exp, str)
}

func (s *_testing) Reports_style_of_cell(t *T) {
	ui, tt := LstFixture(t.GoT(), nil, 0)

	tt.Display(fx, ui.NewStyle().
		WithFG(api.Yellow).WithBG(api.Blue).WithAA(api.Italic))
	ui.Redraw()
	for _, l := range tt.CellsArea(3, 1, 13, 3) {
		for i := range l {
			t.True(l.HasFG(i, api.Yellow))
			t.True(l.HasBG(i, api.Blue))
			t.True(l.HasAA(i, api.Italic))
		}
	}
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

func (s *_testing) Reports_mouse_click(t *T) {
	var clickReported api.ButtonMask
	_, tt := LstFixture(t.GoT(), func(evt api.Eventer) {
		switch evt := evt.(type) {
		case *api.MouseClick:
			clickReported = evt.Button()
		}
	}, 0)

	tt.PostClick(0, 0, api.Button3, 0)
	t.Eq(api.Button3, clickReported)
}

func (s *_testing) Reports_mouse_move(t *T) {
	expX, expY, expOX, expOY := 10, 5, 0, 0
	gotX, gotY, gotOX, gotOY := 0, 0, 0, 0
	_, tt := LstFixture(t.GoT(), func(evt api.Eventer) {
		switch evt := evt.(type) {
		case *api.MouseMove:
			gotX, gotY = evt.Pos()
			gotOX, gotOY = evt.Origin()
		}
	}, 0)

	tt.PostMove(expX, expY)
	t.True(expX == gotX && expY == gotY)
	t.True(expOX == gotOX && expOY == gotOY)
}

func (s *_testing) Reports_drag_and_drop(t *T) {
	expX, expY, expDX, expDY := 5, 2, 10, 5
	gotX, gotY, gotDX, gotDY := 0, 0, 0, 0
	_, tt := LstFixture(t.GoT(), func(evt api.Eventer) {
		switch evt := evt.(type) {
		case *api.MouseDrag:
			gotX, gotY = evt.Origin()
		case *api.MouseDrop:
			gotDX, gotDY = evt.Pos()
		}
	}, 0)

	drop := tt.PostDrag(expX, expY, api.Primary, api.ZeroModifier)
	drop(expDX, expDY)
	t.True(expX == gotX && expY == gotY)
	t.True(expDX == gotDX && expDY == gotDY)
}

func TestTesting(t *testing.T) {
	t.Parallel()
	Run(&_testing{}, t)
}
