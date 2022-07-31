// Copyright (c) 2022 Stephan Lukits. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package lines

import (
	"fmt"
	"testing"

	. "github.com/slukits/gounit"
)

type env struct{ Suite }

func (s *env) SetUp(t *T) { t.Parallel() }

type envCmpFX struct {
	Component
	env *Env
}

func (c *envCmpFX) OnInit(e *Env) { c.env = e }

func (s *env) Access_panics_outside_event_processing(t *T) {
	cmp := &envCmpFX{}
	ee, _ := Test(t.GoT(), cmp, -1)
	defer ee.QuitListening()
	t.Panics(func() { fmt.Fprint(cmp.env, "panics") })
}

func TestEnv(t *testing.T) {
	t.Parallel()
	Run(&_component{}, t)
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
