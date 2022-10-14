// Copyright (c) 2022 Stephan Lukits. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package term

import (
	"strings"
	"testing"
	"time"

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

// LstFixture instantiates a new UI with an simulation screen and an
// testing instance reporting all events to provided listener.  The ui
// is setup for transactional events, i.e. an event-post p returns not
// before all event-posts have been processed which were posted during
// p's processing.  Processing of p times out after given timeout.
// LstFixture does not returns before initial resize event wasn't
// processed.  Listener may be nil and a zero-timeout defaults to 100ms.
func LstFixture(
	t *testing.T, listener func(api.Eventer), timeout time.Duration,
) (*UI, *Testing) {

	t.Helper()

	if timeout == 0 {
		timeout = 100 * time.Millisecond
	}
	ui := initUI(tcell.NewSimulationScreen("UTF-8"), listener)
	t.Cleanup(func() { ui.Quit() })
	ui.EnableTransactionalEventPosts(timeout)

	tt := &Testing{
		t: t, ui: ui, Width: tstWidth, Height: tstHeight}
	tt.PostResize(tstWidth, tstHeight)

	return ui, tt
}

// Fixture instantiates a new UI with an simulation screen and an
// testing instance.  In order to receive events [Testing.Listen] must
// be called additionally.  The ui is setup for transactional events,
// i.e. an event-post p returns not before all event-posts have been
// processed which were posted during p's processing.  Processing of p
// times out after given timeout.  A zero-timeout defaults to 100ms.
func Fixture(t *testing.T, timeout time.Duration) (*UI, *Testing) {
	t.Helper()

	if timeout == 0 {
		timeout = 100 * time.Millisecond
	}
	ui := initUI(tcell.NewSimulationScreen("UTF-8"), nil)
	t.Cleanup(func() { ui.Quit() })
	ui.EnableTransactionalEventPosts(timeout)

	tt := &Testing{
		t: t, ui: ui, Width: tstWidth, Height: tstHeight}

	return ui, tt
}

// Listen sets the listener and posts initial resize event.  Listen is an
// no-op if already a listener is set.
func (tt *Testing) Listen(l func(api.Eventer)) {
	if tt.ui.listener != nil {
		return
	}
	tt.ui.listener = l
	tt.PostResize(tstWidth, tstHeight)
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

func (tt *Testing) PostKey(k api.Key, m api.Modifier) {
	tt.t.Helper()
	if err := tt.ui.Post(newKeyEvent(k, m)); err != nil {
		tt.t.Fatalf("post: key: %v", err)
	}
}

func (tt *Testing) PostRune(r rune, m api.Modifier) {
	tt.t.Helper()
	if err := tt.ui.Post(newRuneEvent(r, m)); err != nil {
		tt.t.Fatalf("post: rune: %v", err)
	}
}

func (tt *Testing) PostMouse(
	x, y int, b api.Button, m api.Modifier,
) {
	tt.t.Helper()
	if err := tt.ui.Post(newMouseEvent(x, y, b, m)); err != nil {
		tt.t.Fatalf("post: mouse: %v", err)
	}
}

func (tt *Testing) PostBracketPaste(paste string) {
	tt.t.Helper()
	if len(paste) == 0 {
		return
	}
	if err := tt.ui.Post(newBracketPaste(true)); err != nil {
		tt.t.Fatalf("post: bracket paste: start: %v", err)
	}
	for _, r := range paste {
		err := tt.ui.Post(newRuneEvent(r, api.ZeroModifier))
		if err != nil {
			tt.t.Fatalf("post: bracket paste: rune: %v", err)
		}
	}
	if err := tt.ui.Post(newBracketPaste(false)); err != nil {
		tt.t.Fatalf("post: bracket paste: end: %v", err)
	}
}

func (tt *Testing) PostResize(width, height int) {
	tt.t.Helper()
	if width == 0 && height == 0 {
		return
	}
	w, h := tt.ui.Size()
	if width == 0 {
		width = w
	}
	if height == 0 {
		height = h
	}
	tt.ui.lib.(tcell.SimulationScreen).SetSize(width, height)
	if err := tt.ui.Post(newResize(width, height)); err != nil {
		tt.t.Fatalf("post: resize: %v", err)
	}
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
