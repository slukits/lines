// Copyright (c) 2022 Stephan Lukits. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package lines

import (
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/slukits/lines/internal/lyt"
)

// Testing augments lines.Events instance created by *Test* with useful
// features for testing like firing an event or getting the current
// screen content as string.
//
// An Events/Testing-instances may not be used concurrently.
//
// An Events.Listen-method becomes non-blocking and starts the
// event-loop in its own go-routine.
//
// All event triggering methods start event-listening if it is not
// already started.
//
// It is guaranteed that all methods of an Events/Testing-instances
// which trigger an event do not return before the event is processed
// and any writes to environments are printed to the screen.  This holds
// also true if an event triggering method is called within a listener
// callback.
type Testing struct {
	ee                        *Events
	lib                       tcell.SimulationScreen
	autoTerminate, terminated bool
	mutex                     *sync.Mutex
	syncAdd                   chan bool
	t                         *testing.T

	// Max is the number of reported events after which the
	// event-loop of a register-fixture is terminated.  Max is
	// decremented after each reported event.  I.e. events for which no
	// listener is registered are not counted.
	Max int

	// LastScreen provides the screen content right before quitting
	// listening.  NOTE it is guaranteed that that this snapshot is
	// taken *after* all lines-updates have made it to the screen.
	LastScreen TestScreen

	// Timeout defines how long an event-triggering method waits for the
	// event to be processed.  It defaults to 100ms.
	Timeout time.Duration
}

// Test provides a slightly differently behaving Events instance and an
// augmenting Testing instance adding features useful for testing.
//
// The here provided Events instance has a non-blocking Listen method
// and all its methods triggering events are guaranteed to return after
// the event and subsequently triggered events have been processed and
// the (simulation) screen is synchronized.  All event triggering
// methods start the event loop automatically if not started, i.e. a
// call to Listen can be skipped.
//
// The Testing instance provides an event countdown which ends the event
// loop once it is zero.  Provide as last argument 0 for an indefinitely
// running event loop.  The default is 1.  NOTE reported OnInit and
// OnLayout events are accumulated and each is counted as one reported
// event for the event countdown.
//
// Testing provides methods for firing user input events which start the
// event-loop if not started and do return after the event and
// subsequently triggered events have been processed and the screen has
// been synchronized.
func Test(t *testing.T, c Componenter, max ...int) (*Events, *Testing) {
	t.Helper()
	scr, err := newSim(c)
	if err != nil {
		t.Fatalf("test: init sim: %v", err)
	}
	ee := &Events{
		mutex:     &sync.Mutex{},
		scr:       scr,
		pollEvent: scr.lib.PollEvent,
		postEvent: scr.lib.PostEvent,
		synced:    make(chan bool, 1),
	}
	ee.t = &Testing{
		ee:      ee,
		lib:     scr.lib.(tcell.SimulationScreen),
		t:       t,
		Timeout: 200 * time.Millisecond,
		mutex:   &sync.Mutex{},
	}
	switch len(max) {
	case 0:
		ee.t.SetMax(1)
	default:
		ee.t.SetMax(max[0])
	}
	return ee, ee.t
}

// SetMax defines the maximum number of reported events before the event
// loop is terminated automatically.  If m is 0 (or lower) listening
// doesn't stop automatically.
func (tt *Testing) SetMax(m int) *Events {
	switch {
	case m <= 0:
		if tt.ee.reported != nil {
			tt.ee.reported = nil
		}
		tt.autoTerminate = false
	default:
		tt.ee.Reported(decrement(tt))
		tt.autoTerminate = true
	}
	tt.Max = m
	return tt.ee
}

func decrement(tt *Testing) func() {
	return func() {
		tt.Max--
	}
}

// listen posts the initial resize event and starts listening for events
// in a new go-routine.  listen returns after the initial resize has
// completed.
func (tt *Testing) listen() *Events {
	tt.t.Helper()
	if tt.terminated {
		panic("listening has already been terminated.")
	}
	tt.ee.setListening()
	go tt.ee.listen()
	// fmt.Println("register: listen")
	wait := tt.registerEventSync("test: listen: sync timed out")
	err := tt.lib.PostEvent(tcell.NewEventResize(tt.lib.Size()))
	if err != nil { // TODO: coverage
		tt.t.Fatalf("test: listen: post resize: %v", err)
	}
	if wait != nil {
		wait()
		tt.checkTermination()
	}
	return tt.ee
}

// FireResize posts a resize event and returns after this event
// has been processed.  Is associated Events instance not listening
// it is started before the event is fired.  NOTE this event as such is
// not reported, i.e. the event countdown is not reduced through this
// event.  But subsequently triggered OnInit or OnLayout events are
// counting down if reported.
func (tt *Testing) FireResize(width, height int) *Events {
	tt.t.Helper()
	if !tt.ee.IsListening() {
		tt.listen()
	}
	if width == 0 && height == 0 {
		return tt.ee
	}
	w, h := tt.lib.Size()
	if w == width && h == height {
		return tt.ee
	}
	if width == 0 {
		width = w
	}
	if height == 0 {
		height = h
	}
	tt.lib.SetSize(width, height)
	// fmt.Println("register: resize")
	wait := tt.registerEventSync("test: set number of lines: sync timed out")
	err := tt.lib.PostEvent(tcell.NewEventResize(width, height))
	if err != nil {
		tt.t.Fatal(err) // TODO: not covered
	}
	if wait != nil {
		wait()
		tt.checkTermination()
	}
	return tt.ee
}

// FireRune posts given run-key-press event and returns after this event
// has been processed.  Note modifier keys are ignored for
// rune-triggered key-events.  Is associated Events instance not
// listening it is started before the event is fired.
func (tt *Testing) FireRune(r rune) *Events {
	if !tt.ee.IsListening() {
		tt.listen()
	}
	// fmt.Println("register: fire rune")
	wait := tt.registerEventSync("test: fire rune: sync timed out")
	tt.lib.InjectKey(tcell.KeyRune, r, tcell.ModNone)
	if wait != nil {
		wait()
		tt.checkTermination()
	}
	return tt.ee
}

// FireKey posts given special-key event and returns after this
// event has been processed.  Is associated Events instance not
// listening it is started before the event is fired.
func (tt *Testing) FireKey(k tcell.Key, m ...tcell.ModMask) *Events {
	tt.t.Helper()
	if !tt.ee.IsListening() {
		tt.listen()
	}
	// fmt.Println("register: fire key")
	wait := tt.registerEventSync("test: fire key: sync timed out")
	if len(m) == 0 {
		tt.lib.InjectKey(k, 0, tcell.ModNone)
	} else {
		tt.lib.InjectKey(k, 0, m[0])
	}
	if wait != nil {
		wait()
		tt.checkTermination()
	}
	return tt.ee
}

// FireClick posts a first button click at given coordinates and returns
// after this event has been processed.  Is associated Events instance
// not listening it is started before the event is fired.  Are given
// coordinates outside the available screen area the call is ignored.
func (tt *Testing) FireClick(x, y int) *Events {
	tt.t.Helper()
	if !tt.ee.IsListening() {
		tt.listen()
	}
	width, height := tt.lib.Size()
	if x < 0 || y < 0 || x >= width || y >= height {
		return tt.ee
	}
	// fmt.Println("register: fire click")
	wait := tt.registerEventSync("test: fire click: sync timed out")
	tt.lib.InjectMouse(x, y, tcell.ButtonPrimary, tcell.ModNone)
	if wait != nil {
		// fmt.Println("wait for mouse")
		wait()
		tt.checkTermination()
	}
	return tt.ee
}

// FireContext posts a second button click at given coordinates and
// returns after this event has been processed.  Is associated Events
// instance not listening it is started before the event is fired.  Are
// given coordinates outside the available screen area the call is
// ignored.
func (tt *Testing) FireContext(x, y int) *Events {
	tt.t.Helper()
	if !tt.ee.IsListening() {
		tt.listen()
	}
	width, height := tt.lib.Size()
	if x < 0 || y < 0 || x >= width || y >= height {
		return tt.ee
	}
	// fmt.Println("register: fire context")
	wait := tt.registerEventSync("test: fire click: sync timed out")
	tt.lib.InjectMouse(x, y, tcell.ButtonSecondary, tcell.ModNone)
	if wait != nil {
		wait()
		tt.checkTermination()
	}
	return tt.ee
}

// FireMouse posts a mouse event with provided arguments and returns
// after this event has been processed.  Is associated Events instance
// not listening it is started before the event is fired.  Are given
// coordinates outside the available screen area the call is ignored.
func (tt *Testing) FireMouse(
	x, y int, bm tcell.ButtonMask, mm tcell.ModMask,
) *Events {
	tt.t.Helper()
	if !tt.ee.IsListening() {
		tt.listen()
	}
	width, height := tt.lib.Size()
	if x < 0 || y < 0 || x >= width || y >= height {
		return tt.ee
	}
	// fmt.Println("register: fire mouse")
	wait := tt.registerEventSync("test: fire mouse: sync timed out")
	tt.lib.InjectMouse(x, y, bm, mm)
	if wait != nil {
		wait()
		tt.checkTermination()
	}
	return tt.ee
}

// FireComponentClick posts an click on given relative coordinate in
// given componenter.  Is associated Events instance not listening it is
// started before the event is fired.  Note if x or y are outside the
// component's screen area or the component is not part of the layout no
// click will be fired.
func (tt *Testing) FireComponentClick(c Componenter, x, y int) *Events {
	tt.t.Helper()
	if !tt.ee.IsListening() { // to calculate the layout
		tt.listen()
	}
	if !c.hasLayoutWrapper() {
		return tt.ee
	}
	ox, oy, ok := isInside(c.layoutComponent().wrapped().dim, x, y)
	if !ok {
		return tt.ee
	}
	return tt.FireClick(ox+x, oy+y)
}

// FireComponentContext posts an "right"-click on given relative
// coordinate in given componenter.  Is associated Events instance not
// listening it is started before the event is fired.  Note if x or y
// are outside the component's screen area or the component is not part
// of the layout no click will be fired.
func (tt *Testing) FireComponentContext(c Componenter, x, y int) *Events {
	tt.t.Helper()
	if !tt.ee.IsListening() { // to calculate the layout
		tt.listen()
	}
	if !c.hasLayoutWrapper() {
		return tt.ee
	}
	ox, oy, ok := isInside(c.layoutComponent().wrapped().dim, x, y)
	if !ok {
		return tt.ee
	}
	return tt.FireContext(ox+x, oy+y)
}

func isInside(dim *lyt.Dim, x, y int) (ox, oy int, ok bool) {
	if x < 0 || y < 0 || dim.IsOffScreen() {
		return 0, 0, false
	}
	_, _, width, height := dim.Area()
	if x >= width {
		return 0, 0, false
	}
	if y >= height {
		return 0, 0, false
	}
	return dim.X(), dim.Y(), true
}

func (tt *Testing) registerEventSync(err string) (wait func()) {
	tt.mutex.Lock()
	defer tt.mutex.Unlock()
	if tt.syncAdd != nil {
		tt.syncAdd <- true
		return nil
	}
	sw := newSyncWait(tt, err)
	tt.syncAdd = make(chan bool)
	go syncGroup(sw.c, tt.syncAdd, tt.ee.synced)
	return sw.f
}

func (tt *Testing) checkTermination() {
	if !tt.autoTerminate {
		return
	}
	if tt.Max <= 0 {
		// the last reported event might was a quit event,
		if tt.ee.IsListening() { // i.e. we stopped already listening
			tt.ee.QuitListening()
		}
	}
}

// syncGroup is send of in a go routing which waits for so many sync
// events until a counter fed by add is zero.
func syncGroup(waite chan struct{}, add chan bool, less chan bool) {
	n := 0
	for {
		select {
		case add := <-add:
			if !add {
				close(waite)
				// fmt.Println("dbg: add: closed")
				return
			}
			n++
			// fmt.Printf("dbg: added: %d\n", n)
		case less := <-less:
			if !less || n == 0 {
				close(waite)
				// fmt.Println("dbg: less: closed")
				return
			}
			n--
			// fmt.Printf("dbg: lessened: %d\n", n)
		}
	}
}

type syncWait struct {
	c chan struct{}
	f func()
}

func newSyncWait(tt *Testing, err string) *syncWait {
	c := make(chan struct{})
	return &syncWait{c: c, f: syncClosure(tt, err, c)}
}

func syncClosure(tt *Testing, err string, c chan struct{}) func() {
	return func() {
		tt.t.Helper()
		select {
		case <-time.After(tt.Timeout):
			tt.t.Fatal(err)
		case <-c:
			// fmt.Println("dbg: stopped waiting")
			tt.mutex.Lock()
			defer tt.mutex.Unlock()
			tt.syncAdd = nil
		}
	}
}

func (tt *Testing) beforeFinalize() {
	tt.mutex.Lock()
	defer tt.mutex.Unlock()
	tt.terminated = true
	tt.LastScreen = tt.Screen()
	if tt.syncAdd != nil {
		close(tt.syncAdd)
	}
}

func (tt *Testing) String() {
	b, w, _ := tt.lib.GetContents()
	bld, j := strings.Builder{}, 1
	for i, c := range b {
		r := ' '
		if len(c.Runes) > 0 {
			r = c.Runes[0]
		}
		bld.WriteRune(r)
		if (i+j)%w == 0 {
			j++
			if j == w {
				j = 0
			}
			bld.WriteString("\n")
		}
	}
}

// Screen returns a trimmed cells matrix which may be stringified or
// investigated for expected styling.  The screen content is trimmed to
// the smallest possible rectangle containing all non blank cells:
//
//	+--------------------+
//	|                    |       +------------+
//	|   upper left       |       |upper left  |
//	|                    |  =>   |            |
//	|          right     |       |       right|
//	|      bottom        |       |   bottom   |
//	|                    |       +------------+
//	+--------------------+
//
// A cell is considered blank if its rune is ' ', '\t' or '\r'
func (tt *Testing) Screen() TestScreen {
	return tt.FullScreen().trimVertical().trimHorizontal()
}

// FullScreen returns a matrix of test cells holding a copy of each
// tcell's sim-cell rune and style information.
func (tt *Testing) FullScreen() TestScreen {
	ts, start := TestScreen{}, 0
	b, w, h := tt.lib.GetContents()
	for i := 0; i < h; i++ {
		l := TestLine{}
		for _, c := range b[start : start+w] {
			if len(c.Runes) == 0 {
				l = append(l, testCell{r: ' ', sty: c.Style})
				continue
			}
			l = append(l, testCell{r: c.Runes[0], sty: c.Style})
		}
		ts = append(ts, l)
		start += w
	}
	return ts
}

// TestScreen is a trimmed convenience representation of a tcell
// simulation screen to evaluate the simulation screen's state
// against an expected state.
type TestScreen []TestLine

// String returns a string representation of a [lines.TestScreen]
func (s TestScreen) String() string {
	if len(s) == 0 {
		return ""
	}
	b := strings.Builder{}
	for _, l := range s {
		for _, c := range l {
			b.WriteRune(c.r)
		}
		b.WriteRune('\n')
	}
	return b.String()[:b.Len()-1]
}

func (s TestScreen) trimVertical() TestScreen {
	if len(s) == 0 {
		return s
	}

	blankAtBeginning := 0
	for _, l := range s {
		if !l.isBlank() {
			break
		}
		blankAtBeginning++
	}
	blankAtEnd := 0
	for i := len(s) - 1; i >= 0; i-- {
		if !s[i].isBlank() {
			break
		}
		blankAtEnd++
	}

	if len(s) == blankAtBeginning {
		return TestScreen{}
	}
	return s[blankAtBeginning : len(s)-blankAtEnd]
}

func (s TestScreen) trimHorizontal() TestScreen {
	if len(s) == 0 {
		return s
	}
	leftTrim, rightTrim := len(s[0]), len(s[0])
	for _, l := range s {
		if leftTrim > 0 && leftTrim > l.blankPrefix() {
			leftTrim = l.blankPrefix()
		}
		if rightTrim > 0 && rightTrim > l.blankSuffix() {
			rightTrim = l.blankSuffix()
		}
		if leftTrim == 0 && rightTrim == 0 {
			break
		}
	}
	for i, l := range s {
		s[i] = l[leftTrim : len(l)-rightTrim]
	}
	return s
}

// TestLine represents a line of a [lines.TestScreen].
type TestLine []testCell

type testCell struct {
	r   rune
	sty tcell.Style
}

const blanks = " \r\t"

// Styles returns a test screen's test line's styles to validate a test
// line's cell's styles like foreground color, background color or style
// attributes.
func (l TestLine) Styles() LineTestStyles {
	if len(l) == 0 {
		return nil
	}
	cfg, cbg, caa := l[0].sty.Decompose()
	ss, cr := LineTestStyles{}, Range{0}
	for i, c := range l {
		fg, bg, aa := c.sty.Decompose()
		if cbg == bg && cfg == fg && caa == aa {
			continue
		}
		cr[1] = i
		ss[cr] = TestStyle{bg: cbg, fg: cfg, aa: caa}
		cr = Range{i}
		cfg, cbg, caa = fg, bg, aa
	}
	cr[1] = len(l)
	ss[cr] = TestStyle{bg: cbg, fg: cfg, aa: caa}
	return ss
}

func (l TestLine) isBlank() bool {
	for _, c := range l {
		if strings.ContainsRune(blanks, c.r) {
			continue
		}
		return false
	}
	return true
}

func (l TestLine) blankPrefix() int {
	n := 0
	for _, c := range l {
		if strings.ContainsRune(blanks, c.r) {
			n++
			continue
		}
		break
	}
	return n
}

func (l TestLine) blankSuffix() int {
	n := 0
	for i := len(l) - 1; i >= 0; i-- {
		if strings.ContainsRune(blanks, l[i].r) {
			n++
			continue
		}
		break
	}
	return n
}

// Range is a two component array of which the first component should be
// smaller than the second, i.e. r.Start() <= r.End() if r is a
// Range-instance.
type Range [2]int

// Start index of a [lines.TestLine] style range.  Not the start index
// is inclusive.
func (r Range) Start() int { return r[0] }

// End index of a [lines.TestLine] style range.  Note the end index is
// exclusive.
func (r Range) End() int { return r[1] }

// Contains returns true if given i is in the style range r
// [r.Start,r.End[.
func (r Range) Contains(i int) bool {
	return r.Start() <= i && i < r.End()
}

// LineTestStyles are provided by a [lines.TestLine] of a [lines.TestScreen]
// to validate a line cell's styles like foreground color, background
// color or style attributes (see [tcell.AttrMask]).
type LineTestStyles map[Range]TestStyle

var defaultTestStyle = func() TestStyle {
	fg, bg, aa := tcell.StyleDefault.Decompose()
	return TestStyle{bg: bg, fg: fg, aa: aa}
}()

// Of returns an [lines.TestStyle] instance for given line styles' cell.
func (s LineTestStyles) Of(cell int) TestStyle {
	for r := range s {
		if !r.Contains(cell) {
			continue
		}
		return s[r]
	}
	return defaultTestStyle
}

// TestStyle returned by a test screen's line for one of its cells to
// verify its style like background color, foreground color or style
// attribute.
type TestStyle struct {
	fg, bg tcell.Color
	aa     tcell.AttrMask
}

// Has returns true if test style s has given style attribute attr;
// false otherwise.
func (s TestStyle) Has(attr tcell.AttrMask) bool {
	return s.aa&attr == attr
}

// HasBG returns true if test style s has given color as background
// color; false otherwise.
func (s TestStyle) HasBG(color tcell.Color) bool {
	return s.bg == color
}

// HasFG returns true if test style s has given color as foreground
// color; false otherwise.
func (s TestStyle) HasFG(color tcell.Color) bool {
	return s.fg == color
}
