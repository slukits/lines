// Copyright (c) 2022 Stephan Lukits. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package fx

import (
	"time"

	"github.com/slukits/lines"
)

type Counter int

const (
	OnInit Counter = iota
	OnLayout
	OnUpdate
	OnFocus
	OnFocusLost
	OnLineFocus
	OnLineFocusLost
	OnLineSelection
	OnLineOverflowing
	OnCursor
	OnRune
	OnEdit
	OnOutOfBoundClick
	OnOutOfBoundMove
	OnMouse
	OnKey
)

// TODO: refactor make this component replace all component fixtures
// which are not nesting in any lines test.
type Cmp struct {
	lines.Component
	OnInit            func(*Cmp, *lines.Env)
	OnLayout          func(*Cmp, *lines.Env)
	OnUpdate          func(*Cmp, *lines.Env, interface{})
	OnFocus           func(*Cmp, *lines.Env)
	OnFocusLost       func(*Cmp, *lines.Env)
	OnLineFocus       func(_ *Cmp, _ *lines.Env, cIdx, sIdx int)
	OnLineFocusLost   func(_ *Cmp, _ *lines.Env, cIdx, sIdx int)
	OnLineSelection   func(_ *Cmp, _ *lines.Env, cIdx, sIdx int)
	OnLineOverflowing func(_ *Cmp, _ *lines.Env, left, right bool)
	OnCursor          func(_ *Cmp, _ *lines.Env, absOnly bool)
	OnRune            func(*Cmp, *lines.Env, rune, lines.ModifierMask)
	OnEdit            func(*Cmp, *lines.Env, *lines.Edit) bool
	OnMouse           func(*Cmp, *lines.Env, lines.ButtonMask, int, int)
	OnKey             func(*Cmp, *lines.Env, lines.Key, lines.ModifierMask)
	cc                map[Counter]int
	tt                map[Counter]time.Time
	gaps              bool
}
