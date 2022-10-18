// Copyright (c) 2022 Stephan Lukits. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package lines

import (
	"fmt"
	"strings"
	"testing"

	. "github.com/slukits/gounit"
)

type env struct{ Suite }

// func (s *env) SetUp(t *T) { t.Parallel() }

func (s *env) tt(t *T, cmp Componenter) *Testing {
	return TermFixture(t.GoT(), 0, cmp)
}

type envCmpFX struct {
	Component
	env  *Env
	test func(e *Env)
}

func (c *envCmpFX) OnInit(e *Env) {
	c.env = e
	if c.test == nil {
		return
	}
	c.test(e)
}

func (c *envCmpFX) OnUpdate(e *Env) {
	f, ok := e.Evt.(*UpdateEvent).Data.(func(e *Env))
	if !ok {
		return
	}
	f(e)
}

func (s *env) Access_panics_outside_event_processing(t *T) {
	cmp := &envCmpFX{}
	s.tt(t, cmp)
	t.Panics(func() { fmt.Fprint(cmp.env, "panics") })
}

func (s *env) Provides_the_display_size(t *T) {
	tt := s.tt(t, &envCmpFX{})
	width, height := tt.Size()
	tt.Lines.Update(tt.Root(), nil, func(e *Env) {
		envWidth, envHeight := e.DisplaySize()
		t.Eq(width, envWidth)
		t.Eq(height, envHeight)
	})
}

func (s *env) Provides_writer_for_the_nth_line(t *T) {
	tt := s.tt(t, &envCmpFX{test: func(e *Env) {
		fmt.Fprint(e.LL(0), "first line")
		fmt.Fprint(e.LL(7), "eighth line")
	}})

	sl := strings.Split(tt.Screen().Trimmed().String(), "\n")
	t.FatalIfNot(t.Eq(8, len(sl)))
	t.Eq("first line", strings.TrimSpace(sl[0]))
	t.Eq("eighth line", strings.TrimSpace(sl[7]))
}

func (s *env) Overwrites_given_line_and_following(t *T) {
	fxCmp := &envCmpFX{test: func(e *Env) {
		fmt.Fprint(e.LL(0), "first line")
		fmt.Fprint(e.LL(7), "eighth line")
	}}
	tt := s.tt(t, fxCmp)
	tt.Lines.Update(fxCmp, nil, func(e *Env) {
		fmt.Fprint(e.LL(6), "seventh line\n"+
			"short 8th\nninth line")
	})

	sl := tt.Screen().Trimmed()
	t.FatalIfNot(t.Eq(9, len(sl)))
	t.Eq(strings.TrimSpace(sl[0]), "first line")
	t.Eq(strings.TrimSpace(sl[6]), "seventh line")
	t.Eq(strings.TrimSpace(sl[7]), "short 8th")
	t.Eq(strings.TrimSpace(sl[8]), "ninth line")
}

func (s *env) Changes_fore_and_background_for_line_s_content(t *T) {
	tt := s.tt(t, &icmpFX{init: func(_ *icmpFX, e *Env) {
		fmt.Fprint(e.BG(Red).FG(White),
			"text with read back- and white foreground")
	}})
	l := tt.Cells()[0]
	str := strings.TrimSpace(l.String())
	for i := range str {
		t.True(l.HasBG(i, Red))
		t.True(l.HasFG(i, White))
		t.Eq(l[i].Rune, int32(str[i]))
	}
}

func (s *env) Changes_fore_and_background_for_whole_line(t *T) {
	tt := s.tt(t, &icmpFX{init: func(_ *icmpFX, e *Env) {
		fmt.Fprint(e.BG(Red).FG(White).LL(0), "line with space")
	}})
	l0 := tt.Cells()[0]
	for i := range l0 {
		t.True(l0.HasBG(i, Red))
		t.True(l0.HasFG(i, White))
	}
}

func (s *env) Changes_fore_and_background_for_partial_line(t *T) {
	tt := s.tt(t, &icmpFX{init: func(_ *icmpFX, e *Env) {
		fmt.Fprint(e.BG(Red).FG(White).At(0, 0), "un-filled with space")
		fmt.Fprint(e.BG(Red).FG(White).LL(1), "filled with space")
	}})

	l0 := tt.Cells()[0]
	for i := range l0 {
		if l0[i].Rune != ' ' || len(l0) > i+1 && l0[i+1].Rune != ' ' {
			t.True(l0.HasBG(i, Red))
			t.True(l0.HasFG(i, White))
			continue
		}
		t.Not.True(l0.HasBG(i, Red))
		t.Not.True(l0.HasFG(i, White))
	}

	l1 := tt.Cells()[1]
	for i := range l1 {
		t.True(l1.HasBG(i, Red))
		t.True(l1.HasFG(i, White))
	}
}

func (s *env) Changes_line_style_for_a_range_of_runes(t *T) {
	tt := s.tt(t, &icmpFX{init: func(_ *icmpFX, e *Env) {
		fmt.Fprint(e, "\t")
		fmt.Fprint(e.FG(White).BG(Red).At(0, 1), "red")
		fmt.Fprint(e.At(0, 1+len("red")), LineFiller+"right")
	}})
	l0, exp := tt.Cells()[0], Range{4, 7}
	str := l0.String()
	for i := range str {
		if exp.Contains(i) {
			t.True(l0.HasBG(i, Red))
			continue
		}
		t.Not.True(l0.HasBG(i, Red))
	}
}

func TestEnv(t *testing.T) {
	t.Parallel()
	Run(&env{}, t)
}
