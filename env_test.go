// Copyright (c) 2022 Stephan Lukits. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package lines

import (
	"fmt"
	"strings"
	"testing"

	"github.com/gdamore/tcell/v2"
	. "github.com/slukits/gounit"
)

type env struct{ Suite }

func (s *env) SetUp(t *T) { t.Parallel() }

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
	ee, _ := Test(t.GoT(), cmp, 0)
	defer ee.QuitListening()
	t.Panics(func() { fmt.Fprint(cmp.env, "panics") })
}

func (s *env) Provides_writer_for_the_nth_line(t *T) {
	ee, tt := Test(t.GoT(), &envCmpFX{test: func(e *Env) {
		fmt.Fprint(e.LL(0), "first line")
		fmt.Fprint(e.LL(7), "eighth line")
	}}, 1)
	ee.Listen()

	sl := strings.Split(tt.LastScreen.String(), "\n")
	t.FatalIfNot(t.Eq(8, len(sl)))
	t.Eq(strings.TrimSpace(sl[0]), "first line")
	t.Eq(strings.TrimSpace(sl[7]), "eighth line")
}

func (s *env) Overwrites_given_line_and_following(t *T) {
	fxCmp := &envCmpFX{test: func(e *Env) {
		fmt.Fprint(e.LL(0), "first line")
		fmt.Fprint(e.LL(7), "eighth line")
	}}
	ee, tt := Test(t.GoT(), fxCmp, 2)
	ee.Listen()
	ee.Update(fxCmp, func(e *Env) {
		fmt.Fprint(e.LL(6), "seventh line\n"+
			"short 8th\nninth line")
	}, nil)

	sl := strings.Split(tt.LastScreen.String(), "\n")
	t.FatalIfNot(t.Eq(9, len(sl)))
	t.Eq(strings.TrimSpace(sl[0]), "first line")
	t.Eq(strings.TrimSpace(sl[6]), "seventh line")
	t.Eq(strings.TrimSpace(sl[7]), "short 8th")
	t.Eq(strings.TrimSpace(sl[8]), "ninth line")
}

func (s *env) Changes_fore_and_background_for_line_s_content(t *T) {
	ee, tt := Test(t.GoT(), &icmpFX{init: func(_ *icmpFX, e *Env) {
		fmt.Fprint(
			e.BG(tcell.ColorRed).FG(tcell.ColorWhite),
			"text with read back- and white foreground",
		)
	}}, 1)
	ee.Listen()
	ss := tt.LastScreen[0].Styles()
	str := tt.LastScreen.String()
	l := tt.LastScreen[0]
	for i := range l {
		t.True(ss.Of(i).HasBG(tcell.ColorRed))
		t.True(ss.Of(i).HasFG(tcell.ColorWhite))
		t.Eq(l[i].r, int32(str[i]))
	}
}

func (s *env) Changes_fore_and_background_for_whole_line(t *T) {
	ee, tt := Test(t.GoT(), &icmpFX{init: func(_ *icmpFX, e *Env) {
		fmt.Fprint(e, "define width for last screen")
		fmt.Fprint(
			e.BG(tcell.ColorRed).Filled().FG(tcell.ColorWhite).LL(1),
			"line with space",
		)
	}}, 1)
	ee.Listen()
	ss := tt.LastScreen[1].Styles()
	l := tt.LastScreen[1]
	for i := range l {
		t.True(ss.Of(i).HasBG(tcell.ColorRed))
		t.True(ss.Of(i).HasFG(tcell.ColorWhite))
	}
}

func (s *env) Changes_fore_and_background_for_partial_line(t *T) {
	ee, tt := Test(t.GoT(), &icmpFX{init: func(_ *icmpFX, e *Env) {
		fmt.Fprint(e, "define width for last screen")
		fmt.Fprint(
			e.BG(tcell.ColorRed).FG(tcell.ColorWhite).LL(1),
			"un-filled with space",
		)
		fmt.Fprint(
			e.BG(tcell.ColorRed).Filled().FG(tcell.ColorWhite).LL(2),
			"filled with space",
		)
	}}, 1)
	ee.Listen()

	ss := tt.LastScreen[1].Styles()
	l := tt.LastScreen[1]
	for i := range l {
		if l[i].r != ' ' || len(l) > i+1 && l[i+1].r != ' ' {
			t.True(ss.Of(i).HasBG(tcell.ColorRed))
			t.True(ss.Of(i).HasFG(tcell.ColorWhite))
			continue
		}
		t.False(ss.Of(i).HasBG(tcell.ColorRed))
		t.False(ss.Of(i).HasFG(tcell.ColorWhite))
	}

	ss = tt.LastScreen[2].Styles()
	l = tt.LastScreen[2]
	for i := range l {
		t.True(ss.Of(i).HasBG(tcell.ColorRed))
		t.True(ss.Of(i).HasFG(tcell.ColorWhite))
	}
}

func TestEnv(t *testing.T) {
	t.Parallel()
	Run(&env{}, t)
}

// type A struct{}
//
// type B struct{ a *A }
//
// func (b *B) A() *A { return b.a }
//
// func (b *B) init() *B {
// 	b.a = &A{}
// 	return b
// }
//
// func (b *B) reset(_b *B) {
// 	*b = *_b
// }
//
// type C struct{ B }
//
// func (s *events) Test(t *T) {
// 	c := &C{}
// 	bInitialized := *c.init()
// 	tBeforeReset := fmt.Sprintf("%T::%[1]v", bInitialized.A())
// 	t.Log(tBeforeReset)
// 	bZero := &B{}
// 	c.reset(bZero)
// 	tInit := fmt.Sprintf("%T::%[1]v", bInitialized.A())
// 	tZero := fmt.Sprintf("%T::%[1]v", bZero.A())
// 	t.Log(tInit, tZero)
// 	cZero := fmt.Sprintf("%T::%[1]v", c.A())
// 	c.reset(&bInitialized)
// 	cInit := fmt.Sprintf("%T::%[1]v", c.A())
// 	t.Log(cZero, cInit)
// }
