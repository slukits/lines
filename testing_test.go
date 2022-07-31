// Copyright (c) 2022 Stephan Lukits. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package lines

import (
	"testing"

	"github.com/gdamore/tcell/v2"
	. "github.com/slukits/gounit"
)

type _Testing struct{ Suite }

func (s *_Testing) SetUp(t *T) { t.Parallel() }

func (s *_Testing) Starts_non_blocking_listening_with_listen_call(t *T) {
	ee, _ := Test(t.GoT(), nil)
	t.False(ee.IsListening())
	ee.Listen()
	t.True(ee.IsListening())
	ee.QuitListening()
	t.False(ee.IsListening())
}

func (s *_Testing) Starts_listening_if_a_resize_is_fired(t *T) {
	ee, tt := Test(t.GoT(), nil)
	defer ee.QuitListening()
	t.False(ee.IsListening())
	t.True(tt.FireResize(22, 42).IsListening())
}

func (s *_Testing) Starts_listening_if_a_key_is_fired(t *T) {
	ee, tt := Test(t.GoT(), nil, -1) // listen for ever
	defer ee.QuitListening()
	t.False(ee.IsListening())
	t.True(tt.FireKey(tcell.KeyBS, 0).IsListening())
}

func (s *_Testing) Starts_listening_a_rune_is_fired(t *T) {
	ee, tt := Test(t.GoT(), nil, -1)
	defer ee.QuitListening()
	t.False(ee.IsListening())
	t.True(tt.FireRune('r').IsListening())
}

func (s *_Testing) Starts_listening_with_update_request(t *T) {
	fx := &cmpFX{}
	ee, _ := Test(t.GoT(), &cmpFX{}, 3) // TODO: clarify what's reported
	defer ee.QuitListening()            // during this test
	t.False(ee.IsListening())
	t.FatalOn(ee.Update(fx, nil, nil))
	t.True(ee.IsListening())
}

func TestTesting(t *testing.T) {
	t.Parallel()
	Run(&_Testing{}, t)
}
