// Copyright (c) 2022 Stephan Lukits. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package lines

import (
	"fmt"
	"testing"

	. "github.com/slukits/gounit"
)

type cmpFX struct{ Component }

type _component struct{ Suite }

func (s *_component) SetUp(t *T) { t.Parallel() }

func (s *_component) Access_panics_outside_event_processing(t *T) {
	cmp := &cmpFX{}
	ee, _ := Test(t.GoT(), cmp, -1)
	defer ee.QuitListening()
	t.Panics(func() { cmp.Dim().SetHeight(20) })
}

func (s *_component) Has_one_line_after_over_writing_one_line(t *T) {
	cmp := &cmpFX{}
	ee, _ := Test(t.GoT(), cmp)
	ee.Update(cmp, nil, func(e *Env) {
		cmp.Mod(Overwriting)
		fmt.Fprint(e, "two\nlines")
		t.Eq(2, cmp.Len())
		fmt.Fprint(e, "one line")
		t.Eq(1, cmp.Len())
	})
	t.False(ee.IsListening())
}

func (s *_component) Has_a_line_more_after_appending_an_line(t *T) {
	cmp := &cmpFX{}
	ee, _ := Test(t.GoT(), cmp)
	ee.Update(cmp, nil, func(e *Env) {
		cmp.Mod(Appending)
		fmt.Fprint(e, "two\nlines")
		fmt.Fprint(e, "one line")
		t.Eq(3, cmp.Len())
	})
}

func (s *_component) Has_a_line_more_after_writing_to_tailing(t *T) {
	cmp := &cmpFX{}
	ee, _ := Test(t.GoT(), cmp)
	ee.Update(cmp, nil, func(e *Env) {
		cmp.Mod(Tailing)
		fmt.Fprint(e, "two\nlines")
		fmt.Fprint(e, "one line")
		t.Eq(3, cmp.Len())
	})
}

func (s *_component) Shows_last_line_clips_above_if_tailing(t *T) {
	cmp := &cmpFX{}
	ee, tt := Test(t.GoT(), cmp)
	tt.FireResize(20, 2)
	ee.Update(cmp, nil, func(e *Env) {
		cmp.Mod(Tailing)
		fmt.Fprint(e, "three\nlines\nat last")
	})
	t.Eq("lines\nat last", tt.LastScreen)
}

func (s *_component) Blanks_a_reset_line(t *T) {
	cmp := &cmpFX{}
	ee, tt := Test(t.GoT(), cmp, 2)
	tt.FireResize(20, 2)
	ee.Update(cmp, nil, func(e *Env) {
		fmt.Fprint(e, "first\nsecond")
	})
	t.Eq(tt.String(), "first\nsecond")

	ee.Update(cmp, nil, func(e *Env) {
		cmp.Reset(-1) // no-op, coverage
		cmp.Reset(0)
	})

	t.Eq("second", tt.LastScreen)
}

func TestComponent(t *testing.T) {
	t.Parallel()
	Run(&_component{}, t)
}
