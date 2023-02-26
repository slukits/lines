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

type Fixture struct {
	t      *testing.T
	ui     *UI
	Width  int
	Height int
}

// TimeOut provides the duration given fixture waits for a polled event
// and all its subordinately triggered events to be processed.
func (fx *Fixture) TimeOut() time.Duration {
	if fx.ui.transactional.Load() == nil {
		return time.Duration(0)
	}
	return fx.ui.transactional.Load().(*transactional).timeout
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
) (*UI, *Fixture) {

	t.Helper()

	if timeout == 0 {
		timeout = 10 * time.Second
	}
	ui := initUI(tcell.NewSimulationScreen("UTF-8"), listener, false)
	t.Cleanup(func() { ui.Quit() })
	ui.EnableTransactionalEventPosts(timeout)
	close(ui.waitForQuit)

	tt := &Fixture{
		t: t, ui: ui, Width: tstWidth, Height: tstHeight}
	tt.PostResize(tstWidth, tstHeight)

	return ui, tt
}

// NewFixture instantiates a new UI with an simulation screen and an
// testing instance.  In order to receive events [Testing.Listen] must
// be called additionally.  The ui is setup for transactional events,
// i.e. an event-post p returns not before all event-posts have been
// processed which were posted during p's processing.  Processing of p
// times out after given timeout.  A zero-timeout defaults to 100ms.
func NewFixture(t *testing.T, timeout time.Duration) (*UI, *Fixture) {
	t.Helper()

	if timeout == 0 {
		timeout = 10 * time.Second
	}
	ui := initUI(tcell.NewSimulationScreen("UTF-8"), nil, false)
	t.Cleanup(func() { ui.Quit() })
	ui.EnableTransactionalEventPosts(timeout)
	close(ui.waitForQuit)

	tt := &Fixture{
		t: t, ui: ui, Width: tstWidth, Height: tstHeight}

	return ui, tt
}

// Listen sets the listener and posts initial resize event.  Listen is an
// no-op if already a listener is set.
func (tt *Fixture) Listen(l func(api.Eventer)) {
	tt.t.Helper()
	if tt.ui.listener != nil {
		return
	}
	tt.ui.listener = l
	tt.PostResize(tstWidth, tstHeight)
}

func (tt *Fixture) Display(s string, sty api.Style) {
	tt.t.Helper()
	if len(s) == 0 {
		return
	}
	for y, s := range strings.Split(s, "\n") {
		for x, r := range s {
			tt.ui.Display(x, y, r, sty)
		}
	}
}

func (tt *Fixture) PostKey(k api.Key, m api.ModifierMask) {
	tt.t.Helper()
	// NOTE UI.Post is used instead of tcell.PostEvent to make
	// transactional event processing work, i.e. an event post doesn't
	// return before all triggered sub-events are processed and all
	// content updates made it to the screen.
	if err := tt.ui.Post(newKeyEvent(k, m)); err != nil {
		tt.t.Fatalf("post: key: %v", err)
	}
}

func (tt *Fixture) PostRune(r rune, m api.ModifierMask) {
	tt.t.Helper()
	// NOTE UI.Post is used instead of tcell.PostEvent to make
	// transactional event processing work, i.e. an event post doesn't
	// return before all triggered sub-events are processed and all
	// content updates made it to the screen.
	if err := tt.ui.Post(newRuneEvent(r, m)); err != nil {
		tt.t.Fatalf("post: rune: %v", err)
	}
}

func newMouseEvent(
	x, y int, b api.ButtonMask, m api.ModifierMask,
) api.MouseEventer {
	return &mouseEvent{evt: tcell.NewEventMouse(
		x, y,
		apiButtonsToTcell(b),
		apiModifiersToTcell(m),
	)}
}

func (tt *Fixture) PostMouse(
	x, y int, b api.ButtonMask, m api.ModifierMask,
) {
	tt.t.Helper()
	// NOTE UI.Post is used instead of tcell.PostEvent to make
	// transactional event processing work, i.e. an event post doesn't
	// return before all triggered sub-events are processed and all
	// content updates made it to the screen.
	if err := tt.ui.Post(newMouseEvent(x, y, b, m)); err != nil {
		tt.t.Fatalf("post: mouse: %v", err)
	}
}

type frameEvent struct {
	exec func()
	when time.Time
}

func (f *frameEvent) Source() interface{} { return f }
func (f *frameEvent) When() time.Time     { return f.when }

func (f *frameEvent) Exec() {
	f.exec()
}

func (tt *Fixture) PostClick(
	x, y int, b api.ButtonMask, m api.ModifierMask,
) {
	tt.t.Helper()
	if b == 0 {
		return
	}
	err := tt.ui.Post(&frameEvent{exec: func() {
		if err := tt.ui.Post(newMouseEvent(x, y, b, m)); err != nil {
			tt.t.Fatalf("post click: mouse down: %v", err)
		}
		if err := tt.ui.Post(newMouseEvent(x, y, 0, 0)); err != nil {
			tt.t.Fatalf("post click: mouse up: %v", err)
		}
	}})
	if err != nil {
		tt.t.Fatalf("post click: %v", err)
	}
}

func (tt *Fixture) PostMove(x, y int, xy ...int) {
	tt.t.Helper()
	err := tt.ui.Post(&frameEvent{exec: func() {
		if len(xy) >= 2 {
			err := tt.ui.Post(newMouseEvent(xy[0], xy[1], 0, 0))
			if err != nil {
				tt.t.Fatalf("post move origin: %v", err)
			}
			err = tt.ui.Post(newMouseEvent(xy[0], xy[1], 0, 0))
			if err != nil {
				tt.t.Fatalf("post move origin: stop: %v", err)
			}
		}
		if err := tt.ui.Post(newMouseEvent(x, y, 0, 0)); err != nil {
			tt.t.Fatalf("post move: %v", err)
		}
		if err := tt.ui.Post(newMouseEvent(x, y, 0, 0)); err != nil {
			tt.t.Fatalf("post move: stop: %v", err)
		}
	}})
	if err != nil {
		tt.t.Fatalf("post move: %v", err)
	}
}

func (tt *Fixture) PostDrag(
	x, y int, b api.ButtonMask, m api.ModifierMask,
) (drop func(x, y int)) {
	return func(dx, dy int) {
		err := tt.ui.Post(&frameEvent{exec: func() {
			err := tt.ui.Post(newMouseEvent(x, y, b, m))
			if err != nil {
				tt.t.Fatalf("post drag origin: %v", err)
			}
			err = tt.ui.Post(newMouseEvent(dx, dy, b, m))
			if err != nil {
				tt.t.Fatalf("post drag: %v", err)
			}
			err = tt.ui.Post(newMouseEvent(dx, dy, 0, 0))
			if err != nil {
				tt.t.Fatalf("post drag drop: %v", err)
			}
		}})
		if err != nil {
			tt.t.Fatalf("post drag 'n' drop: %v", err)
		}
	}
}

func (tt *Fixture) PostBracketPaste(paste string) {
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

func (tt *Fixture) PostResize(width, height int) {
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

func (tt *Fixture) Screen() api.StringScreen {
	bld, screen := &strings.Builder{}, api.StringScreen{}
	err := tt.ui.Post(&screenEvent{when: time.Now(), grab: func() {
		b, w, _ := tt.ui.lib.(tcell.SimulationScreen).GetContents()
		for i, c := range b {
			bld.WriteRune(c.Runes[0])
			if (i+1)%w == 0 {
				screen = append(screen, bld.String())
				bld.Reset()
			}
		}
	}})
	if err != nil {
		tt.t.Fatalf("testing: cells-are: screen-event: %v", err)
	}
	return screen
}

func (tt *Fixture) ScreenArea(x, y, width, height int) api.StringScreen {
	bld, screen := &strings.Builder{}, api.StringScreen{}
	err := tt.ui.Post(&screenEvent{when: time.Now(), grab: func() {
		tt.screenArea(x, y, width, height, func(line []tcell.SimCell) {
			for _, c := range line {
				bld.WriteRune(c.Runes[0])
			}
			screen = append(screen, bld.String())
			bld.Reset()
		})
	}})
	if err != nil {
		tt.t.Fatalf("testing: cells-are: screen-event: %v", err)
	}
	return screen
}

func (tt *Fixture) Cells() api.CellsScreen {
	cs, line := api.CellsScreen{api.CellsLine{}}, 0
	err := tt.ui.Post(&screenEvent{when: time.Now(), grab: func() {
		b, w, _ := tt.ui.lib.(tcell.SimulationScreen).GetContents()
		if w == 0 {
			cs = api.CellsScreen{}
			return
		}
		styler := tcellToApiStyleClosure()
		for i, c := range b {
			cs[line] = append(cs[line], api.TestCell{
				Rune: c.Runes[0], Style: styler(c.Style),
			})
			if (i+1)%w == 0 && i+1 < len(b) {
				line++
				cs = append(cs, api.CellsLine{})
			}
		}
	}})
	if err != nil {
		tt.t.Fatalf("testing: cells-are: screen-event: %v", err)
	}
	return cs
}

func (tt *Fixture) CellsArea(x, y, width, height int) api.CellsScreen {
	cs, line := api.CellsScreen{}, -1
	styler := tcellToApiStyleClosure()
	err := tt.ui.Post(&screenEvent{when: time.Now(), grab: func() {
		tt.screenArea(x, y, width, height, func(l []tcell.SimCell) {
			line++
			cs = append(cs, api.CellsLine{})
			for _, c := range l {
				cs[line] = append(cs[line], api.TestCell{
					Rune: c.Runes[0], Style: styler(c.Style),
				})
			}
		})
	}})
	if err != nil {
		tt.t.Fatalf("testing: cells-are: screen-event: %v", err)
	}
	return cs
}

// screenEvent is used to read the content of a simulation screen while
// making sure no other event is writing to it.
type screenEvent struct {
	when time.Time
	grab func()
}

func (e *screenEvent) When() time.Time     { return e.when }
func (e *screenEvent) Source() interface{} { return e }

// Size returns the number of available lines (height) and the number of
// runes per line (width) off the terminal screen/display.
func (tt *Fixture) Size() (width, height int) {
	return tt.ui.lib.Size()
}

func (tt *Fixture) screenArea(
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
