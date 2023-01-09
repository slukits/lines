// Copyright (c) 2022 Stephan Lukits. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

/* cmpfx.go defines the component fixtures for testing. */

package lines

import (
	"time"

	. "github.com/slukits/gounit"
)

type counter int

const (
	onInit counter = iota
	onLineFocus
	onLineFocusLost
	onLineSelection
	onLineOverflowing
	onCursor
	onRune
	onEdit
)

// TODO: refactor make this component replace all component fixtures
// which are not nesting, i.e. neither stacking nor chaining.
type cmpFX struct {
	Component
	onInit            func(*cmpFX, *Env)
	onLineFocus       func(_ *cmpFX, _ *Env, cIdx, sIdx int)
	onLineFocusLost   func(_ *cmpFX, _ *Env, cIdx, sIdx int)
	onLineSelection   func(_ *cmpFX, _ *Env, cIdx, sIdx int)
	onLineOverflowing func(_ *cmpFX, _ *Env, left, right bool)
	onCursor          func(_ *cmpFX, _ *Env, absOnly bool)
	onRune            func(*cmpFX, *Env, rune, ModifierMask)
	onEdit            func(*cmpFX, *Env, *Edit) bool
	cc                map[counter]int
}

func fx(t *T, cmp Componenter, timeout ...time.Duration) *Fixture {
	d := time.Duration(0)
	if len(timeout) > 0 {
		d = timeout[0]
	}
	if cmp == nil {
		cmp = &cmpFX{}
	}
	return TermFixture(t.GoT(), d, cmp)
}

func fxFF(
	t *T, ff FeatureMask, timeout ...time.Duration,
) (*Fixture, *cmpFX) {
	d := time.Duration(0)
	if len(timeout) > 0 {
		d = timeout[0]
	}
	cmp := &cmpFX{
		onInit: func(c *cmpFX, e *Env) {
			c.FF.Add(ff)
		},
	}
	return TermFixture(t.GoT(), d, cmp), cmp
}

func (c *cmpFX) increment(cn counter) {
	if c.cc == nil {
		c.cc = map[counter]int{}
	}
	c.cc[cn]++
}

func (c *cmpFX) N(cn counter) int { return c.cc[cn] }

func (c *cmpFX) OnInit(e *Env) {
	c.increment(onInit)
	if c.onInit == nil {
		return
	}
	c.onInit(c, e)
}

func (c *cmpFX) OnRune(e *Env, r rune, mm ModifierMask) {
	c.increment(onRune)
	if c.onRune == nil {
		return
	}
	c.onRune(c, e, r, mm)
}

func (c *cmpFX) OnLineFocus(e *Env, cIdx, sIdx int) {
	c.increment(onLineFocus)
	if c.onLineFocus == nil {
		return
	}
	c.onLineFocus(c, e, cIdx, sIdx)
}

func (c *cmpFX) OnLineFocusLost(e *Env, cIdx, sIdx int) {
	c.increment(onLineFocusLost)
	if c.onLineFocusLost == nil {
		return
	}
	c.onLineFocusLost(c, e, cIdx, sIdx)
}

func (c *cmpFX) OnLineSelection(e *Env, cIdx, sIdx int) {
	c.increment(onLineSelection)
	if c.onLineSelection == nil {
		return
	}
	c.onLineSelection(c, e, cIdx, sIdx)
}

func (c *cmpFX) OnLineOverflowing(e *Env, left, right bool) {
	c.increment(onLineOverflowing)
	if c.onLineOverflowing == nil {
		return
	}
	c.onLineOverflowing(c, e, left, right)
}

func (c *cmpFX) OnCursor(e *Env, absOnly bool) {
	c.increment(onCursor)
	if c.onCursor == nil {
		return
	}
	c.onCursor(c, e, absOnly)
}

func (c *cmpFX) OnEdit(e *Env, edt *Edit) bool {
	c.increment(onEdit)
	if c.onEdit == nil {
		return true
	}
	return c.onEdit(c, e, edt)
}
