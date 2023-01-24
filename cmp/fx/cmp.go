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
	NInit Counter = iota
	NLayout
	NUpdate
	NFocus
	NFocusLost
	NLineFocus
	NLineFocusLost
	NLineSelection
	NLineOverflowing
	NCursor
	NRune
	NEdit
	NOutOfBoundClick
	NOutOfBoundMove
	NMouse
	NKey
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

// Wraps a comptonenter enabling the user to inject code on certain
// events.
type Wrap struct {
	lines.Componenter
	ONInit   func(lines.Componenter, *lines.Env)
	ONLayout func(lines.Componenter, *lines.Env) (reflow bool)
	cc       map[Counter]int
}

func (w *Wrap) N(c Counter) int { return w.cc[c] }

func (w *Wrap) OnInit(e *lines.Env) {
	w.cc = map[Counter]int{}
	w.cc[NInit]++
	if w.ONInit != nil {
		w.ONInit(w.Componenter, e)
	}
	i, ok := w.Componenter.(lines.Initer)
	if !ok {
		return
	}
	i.OnInit(e)
}

func (w *Wrap) OnLayout(e *lines.Env) (reflow bool) {
	w.cc[NLayout]++
	if w.ONLayout != nil {
		w.ONLayout(w.Componenter, e)
	}
	l, ok := w.Componenter.(lines.Layouter)
	if !ok {
		return
	}
	return l.OnLayout(e)
}

type Chaining struct {
	lines.Component
	lines.Chaining
}

func (c *Chaining) Set(
	cmp lines.Componenter, cc ...lines.Componenter,
) *Chaining {
	c.CC = append(c.CC, cmp)
	c.CC = append(c.CC, cc...)
	return c
}
