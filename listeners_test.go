// Copyright (c) 2022 Stephan Lukits. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package lines

import (
	"testing"

	. "github.com/slukits/gounit"
)

type Register struct{ Suite }

func (s *Register) tt(t *T, cmp Componenter) *Testing {
	return TermFixture(t.GoT(), 0, cmp)
}

func (s *Register) Reports_to_listener_on_registered_key(t *T) {
	fx, keyReported := &cmpFX{}, false
	tt := s.tt(t, fx)
	fx.Register.Key(Enter, ZeroModifier, func(e *Env) {
		keyReported = true
	})
	tt.FireKey(Enter, ZeroModifier)
	t.True(keyReported)
}

func (s *Register) Reports_to_listener_on_registered_rune(t *T) {
	fx, runeReported := &cmpFX{}, false
	tt := s.tt(t, fx)
	fx.Register.Rune('x', ZeroModifier, func(e *Env) {
		runeReported = true
	})
	tt.FireRune('x', ZeroModifier)
	t.True(runeReported)
}

func (s *Register) Updates_a_key_listener(t *T) {
	fx, keyReported, otherReported := &cmpFX{}, false, false
	tt := s.tt(t, fx)
	fx.Register.Key(Enter, ZeroModifier, func(e *Env) {
		keyReported = true
	})
	tt.FireKey(Enter, ZeroModifier)
	t.True(keyReported)
	fx.Register.Key(Enter, ZeroModifier, func(e *Env) {
		otherReported = true
	})
	tt.FireKey(Enter, ZeroModifier)
	t.True(otherReported)
}

func (s *Register) Updates_a_rune_listener(t *T) {
	fx, runeReported, otherReported := &cmpFX{}, false, false
	tt := s.tt(t, fx)
	fx.Register.Rune('x', ZeroModifier, func(e *Env) {
		runeReported = true
	})
	tt.FireRune('x', ZeroModifier)
	t.True(runeReported)
	fx.Register.Rune('x', ZeroModifier, func(e *Env) {
		otherReported = true
	})
	tt.FireRune('x', ZeroModifier)
	t.True(otherReported)
}

func (s *Register) Removes_a_key_listener(t *T) {
	fx, keyReported := &cmpFX{}, 0
	tt := s.tt(t, fx)
	fx.Register.Key(Enter, ZeroModifier, func(e *Env) {
		keyReported++
	})
	tt.FireKey(Enter, ZeroModifier)
	t.Eq(1, keyReported)
	fx.Register.Key(Enter, ZeroModifier, nil)
	tt.FireKey(Enter, ZeroModifier)
	t.Eq(1, keyReported)
}

func (s *Register) Removes_a_rune_listener(t *T) {
	fx, runeReported := &cmpFX{}, 0
	tt := s.tt(t, fx)
	fx.Register.Rune('x', ZeroModifier, func(e *Env) {
		runeReported++
	})
	tt.FireRune('x', ZeroModifier)
	t.Eq(1, runeReported)
	fx.Register.Rune('x', ZeroModifier, nil)
	tt.FireRune('x', ZeroModifier)
	t.Eq(1, runeReported)
}

func TestRegister(t *testing.T) {
	t.Parallel()
	Run(&Register{}, t)
}
