// Copyright (c) 2022 Stephan Lukits. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

/* cmpfx.go defines the component fixtures for testing. */

package lines

type counter int

const (
	onInit counter = iota
	onLineFocus
	onLineFocusLost
	onLineSelection
	onLineOverflowing
	onCursor
)

// TODO: refactor make this component replace all fixture component
// which are not nesting, i.e. neither stacking nor chaining.
type cmpFX struct {
	Component
	onInit            func(*cmpFX, *Env)
	onLineFocus       func(*cmpFX, *Env, int)
	onLineFocusLost   func(_ *cmpFX, _ *Env, cIdx, sIdx int)
	onLineSelection   func(*cmpFX, *Env, int)
	onLineOverflowing func(_ *cmpFX, _ *Env, left, right bool)
	onCursor          func(_ *cmpFX, _ *Env, absOnly bool)
	cc                map[counter]int
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

func (c *cmpFX) OnLineFocus(e *Env, cIdx, sIdx int) {
	c.increment(onLineFocus)
	if c.onLineFocus == nil {
		return
	}
	c.onLineFocus(c, e, cIdx)
}

func (c *cmpFX) OnLineFocusLost(e *Env, cIdx, sIdx int) {
	c.increment(onLineFocusLost)
	if c.onLineFocusLost == nil {
		return
	}
	c.onLineFocusLost(c, e, cIdx, sIdx)
}

func (c *cmpFX) OnLineSelection(e *Env, i int) {
	c.increment(onLineSelection)
	if c.onLineSelection == nil {
		return
	}
	c.onLineSelection(c, e, i)
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
