// Copyright (c) 2022 Stephan Lukits. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package lines

import (
	"fmt"
	"testing"

	. "github.com/slukits/gounit"
)

type cmpFX struct {
	Component
}

const expInit = "component-fixture initialized"

func (c *cmpFX) OnInit(e *Env) {
	fmt.Fprint(e, expInit)
}

type _component struct{ Suite }

func (s *_component) SetUp(t *T) { t.Parallel() }

func (s *_component) Access_panics_outside_event_processing(t *T) {
	cmp := &cmpFX{}
	ee, _ := Test(t.GoT(), cmp, -1)
	defer ee.QuitListening()
	t.Panics(func() { cmp.Dim().SetHeight(20) })
}

func TestComponent(t *testing.T) {
	t.Parallel()
	Run(&_component{}, t)
}
