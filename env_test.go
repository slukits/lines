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
	}})
	ee.Listen()

	sl := strings.Split(tt.LastScreen, "\n")
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

	sl := strings.Split(tt.LastScreen, "\n")
	t.FatalIfNot(t.Eq(9, len(sl)))
	t.Eq(strings.TrimSpace(sl[0]), "first line")
	t.Eq(strings.TrimSpace(sl[6]), "seventh line")
	t.Eq(strings.TrimSpace(sl[7]), "short 8th")
	t.Eq(strings.TrimSpace(sl[8]), "ninth line")
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
