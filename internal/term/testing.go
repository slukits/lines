// Copyright (c) 2022 Stephan Lukits. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package term

import (
	"strings"
	"testing"

	"github.com/gdamore/tcell/v2"
	"github.com/slukits/lines/internal/api"
)

const tstWidth, tstHeight = 80, 25

type Testing struct {
	t      *testing.T
	ui     *UI
	Width  int
	Height int
}

// Fixture instantiates a new UI with an simulation screen and an
// testing instance.
func Fixture(t *testing.T) (*UI, *Testing) {
	t.Helper()
	ui := initUI(tcell.NewSimulationScreen("UTF-8"))
	ui.lib.PostEvent(tcell.NewEventResize(tstWidth, tstHeight))
	t.Cleanup(func() { ui.Quit() })
	return ui, &Testing{
		t: t, ui: ui, Width: tstWidth, Height: tstHeight}
}

func (tt *Testing) Display(s string, sty api.Style) {
	if len(s) == 0 {
		return
	}
	for y, s := range strings.Split(s, "\n") {
		for x, r := range s {
			tt.ui.Display(x, y, r, sty)
		}
	}
}

func (tt *Testing) PostKey(k api.Key, m api.Modifier) error {
	return tt.ui.Post(newKeyEvent(k, m))
}

func (tt *Testing) PostRune(r rune, m api.Modifier) error {
	return tt.ui.Post(newRuneEvent(r, m))
}

func (tt *Testing) PostMouse(
	x, y int, b api.Button, m api.Modifier,
) error {
	return tt.ui.Post(newMouseEvent(x, y, b, m))
}

func (tt *Testing) PostResize(width, height int) error {
	tt.ui.lib.(tcell.SimulationScreen).SetSize(width, height)
	return tt.ui.Post(newResize(width, height))
}

func (tt *Testing) Screen() api.StringScreen {
	bld, screen := &strings.Builder{}, api.StringScreen{}
	b, w, _ := tt.ui.lib.(tcell.SimulationScreen).GetContents()
	for i, c := range b {
		bld.WriteRune(c.Runes[0])
		if (i+1)%w == 0 {
			screen = append(screen, bld.String())
			bld.Reset()
		}
	}
	return screen
}

func (tt *Testing) ScreenArea(x, y, width, height int) api.StringScreen {
	bld, screen := &strings.Builder{}, api.StringScreen{}
	tt.screenArea(x, y, width, height, func(line []tcell.SimCell) {
		for _, c := range line {
			bld.WriteRune(c.Runes[0])
		}
		screen = append(screen, bld.String())
		bld.Reset()
	})
	return screen
}

func (tt *Testing) Cells() api.CellsScreen {
	b, w, _ := tt.ui.lib.(tcell.SimulationScreen).GetContents()
	if w == 0 {
		return api.CellsScreen{}
	}
	cs, line := api.CellsScreen{api.CellsLine{}}, 0
	styler := tcellToApiStyleClosure()
	for i, c := range b {
		cs[line] = append(cs[line], api.TestCell{
			Rune: c.Runes[0], Sty: styler(c.Style),
		})
		if (i+1)%w == 0 && i+1 < len(b) {
			line++
			cs = append(cs, api.CellsLine{})
		}
	}
	return cs
}

func (tt *Testing) CellsArea(x, y, width, height int) api.CellsScreen {
	cs, line := api.CellsScreen{}, -1
	styler := tcellToApiStyleClosure()
	tt.screenArea(x, y, width, height, func(l []tcell.SimCell) {
		line++
		cs = append(cs, api.CellsLine{})
		for _, c := range l {
			cs[line] = append(cs[line], api.TestCell{
				Rune: c.Runes[0], Sty: styler(c.Style),
			})
		}
	})
	return cs
}

func (tt *Testing) screenArea(
	x, y, width, height int,
	cb func(line []tcell.SimCell),
) {
	if width == 0 || height == 0 {
		return
	}
	b, w, _ := tt.ui.lib.(tcell.SimulationScreen).GetContents()
	lineCount := len(b) / w
	if width+x > w || y+height > lineCount {
		return
	}
	for i := 0; i < lineCount; i++ {
		if i < y {
			continue
		}
		if i-y >= height {
			break
		}
		cb(b[i*w : i*w+w][x : width+x])
	}
}
