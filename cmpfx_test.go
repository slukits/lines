// Copyright (c) 2022 Stephan Lukits. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

/* cmpfx.go defines the component fixtures for testing. */

package lines

import (
	"fmt"
	"time"

	"github.com/slukits/gounit"
)

type counter int

const (
	onInit counter = iota
	onLayout
	onUpdate
	onFocus
	onFocusLost
	onLineFocus
	onLineFocusLost
	onLineSelection
	onLineOverflowing
	onCursor
	onRune
	onEdit
	onOutOfBoundClick
	onOutOfBoundMove
	onMouseN
	onKeyN
)

// TODO: refactor make this component replace all component fixtures
// which are not nesting, i.e. neither stacking nor chaining.
type cmpFX struct {
	Component
	onInit            func(*cmpFX, *Env)
	onLayout          func(*cmpFX, *Env)
	onUpdate          func(*cmpFX, *Env, interface{})
	onFocus           func(*cmpFX, *Env)
	onFocusLost       func(*cmpFX, *Env)
	onLineFocus       func(_ *cmpFX, _ *Env, cIdx, sIdx int)
	onLineFocusLost   func(_ *cmpFX, _ *Env, cIdx, sIdx int)
	onLineSelection   func(_ *cmpFX, _ *Env, cIdx, sIdx int)
	onLineOverflowing func(_ *cmpFX, _ *Env, left, right bool)
	onCursor          func(_ *cmpFX, _ *Env, absOnly bool)
	onRune            func(*cmpFX, *Env, rune, ModifierMask)
	onEdit            func(*cmpFX, *Env, *Edit) bool
	onMouse           func(*cmpFX, *Env, ButtonMask, int, int)
	onKey             func(*cmpFX, *Env, Key, ModifierMask)
	cc                map[counter]int
	tt                map[counter]time.Time
	gaps              bool
}

func fx(
	t *gounit.T, cmp Componenter, timeout ...time.Duration,
) *Fixture {
	d := time.Duration(0)
	if len(timeout) > 0 {
		d = timeout[0]
	}
	if cmp == nil {
		cmp = &cmpFX{}
	}
	return TermFixture(t.GoT(), d, cmp)
}

func fxCmp(
	t *gounit.T, timeout ...time.Duration,
) (*Fixture, *cmpFX) {
	d := time.Duration(0)
	if len(timeout) > 0 {
		d = timeout[0]
	}
	cmp := &cmpFX{}
	return TermFixture(t.GoT(), d, cmp), cmp
}

func fxFF(
	t *gounit.T, ff FeatureMask, timeout ...time.Duration,
) (*Fixture, *cmpFX) {
	d := time.Duration(0)
	if len(timeout) > 0 {
		d = timeout[0]
	}
	cmp := &cmpFX{
		onInit: func(c *cmpFX, e *Env) {
			c.FF.Set(ff)
		},
	}
	return TermFixture(t.GoT(), d, cmp), cmp
}

func (c *cmpFX) increment(cn counter) {
	if c.cc == nil {
		c.cc = map[counter]int{}
	}
	if c.tt == nil {
		c.tt = map[counter]time.Time{}
	}
	c.cc[cn]++
	c.tt[cn] = time.Now()
}

func (c *cmpFX) N(cn counter) int { return c.cc[cn] }

func (c *cmpFX) T(cn counter) time.Time { return c.tt[cn] }

func (c *cmpFX) OnInit(e *Env) {
	c.increment(onInit)
	if c.gaps {
		Print(c.Gaps(0).Filling(), '•')
		fmt.Fprint(c.Gaps(0).Corners, "•")
	}
	if c.onInit == nil {
		return
	}
	c.onInit(c, e)
}

func (c *cmpFX) OnLayout(e *Env) bool {
	c.increment(onLayout)
	if c.onLayout == nil {
		return false
	}
	c.onLayout(c, e)
	return false
}

func (c *cmpFX) OnUpdate(e *Env, data interface{}) {
	c.increment(onUpdate)
	if c.onUpdate == nil {
		return
	}
	c.onUpdate(c, e, data)
}

func (c *cmpFX) OnFocus(e *Env) {
	c.increment(onFocus)
	if c.onFocus == nil {
		return
	}
	c.onFocus(c, e)
}

func (c *cmpFX) OnFocusLost(e *Env) {
	c.increment(onFocusLost)
	if c.onFocusLost == nil {
		return
	}
	c.onFocusLost(c, e)
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

func (c *cmpFX) OnMouse(e *Env, bm ButtonMask, x, y int) {
	c.increment(onMouseN)
	if c.onMouse == nil {
		return
	}
	c.onMouse(c, e, bm, x, y)
}

func (c *cmpFX) OnKey(e *Env, k Key, mm ModifierMask) {
	c.increment(onKeyN)
	if c.onKey == nil {
		return
	}
	c.onKey(c, e, k, mm)
}

type stackingFX struct {
	cmpFX
	Stacking
}

type chainingFX struct {
	cmpFX
	Chaining
}

type modalLayerFX struct {
	cmpFX
	onOutOfBoundClick func(*modalLayerFX, *Env)
	onOutOfBoundMove  func(*modalLayerFX, *Env)
}

func (l *modalLayerFX) OnOutOfBoundClick(e *Env) bool {
	l.increment(onOutOfBoundClick)
	if l.onOutOfBoundClick != nil {
		l.onOutOfBoundClick(l, e)
	}
	return false
}

func (l *modalLayerFX) OnOutOfBoundMove(e *Env) bool {
	l.increment(onOutOfBoundMove)
	if l.onOutOfBoundMove != nil {
		l.onOutOfBoundMove(l, e)
	}
	return false
}

type framingFX struct {
	cmpFX
	filler rune
}

func (c *framingFX) OnInit(e *Env) {
	Print(c.Gaps(0).Filling(), c.filler)
	fmt.Fprint(c.Gaps(0).Corners, string(c.filler))
	c.cmpFX.OnInit(e)
}

type linerFX struct {
	// cc is the provided content of this type's instance which defaults
	// to 8 lines: 1st, 2nd, 3rd, 4th, ..., 8th
	cc []string
}

func lineString(no int) string {
	switch no {
	case 1:
		return "1st"
	case 2:
		return "2nd"
	case 3:
		return "3rd"
	}
	return fmt.Sprintf("%dth", no)
}

func (l *linerFX) initLines(n int) *linerFX {
	if n < 0 {
		panic("new liner fixture: negative number of lines")
	}
	cc := []string{}
	for i := 0; i < n; i++ {
		cc = append(cc, lineString(i+1))
	}
	l.cc = cc
	return l
}

func (l *linerFX) Print(idx int, w *EnvLineWriter) bool {
	if len(l.cc) == 0 {
		l.initLines(8)
	}
	if len(l.cc) <= idx || idx < 0 {
		return false
	}
	fmt.Fprintf(w, l.cc[idx])
	return idx+1 < len(l.cc)
}

type scrollableLinerFX struct {
	linerFX
}

func (l *scrollableLinerFX) Len() int {
	if len(l.cc) == 0 {
		l.initLines(8)
	}
	return len(l.cc)
}

type focusableLinerFX struct {
	scrollableLinerFX
	// returns to a given line index if the line is focusable or not.
	// focusable defaults to func(_ int) bool { return true }
	focusable func(idx int) bool
}

func (l *focusableLinerFX) initLines(n int) *focusableLinerFX {
	l.scrollableLinerFX.initLines(n)
	return l
}

func (l *focusableLinerFX) IsFocusable(idx int) bool {
	if l.focusable == nil {
		return true
	}
	return l.focusable(idx)
}

type srcFX struct {
	cmpFX
	liner Liner
}

func (c *srcFX) OnInit(e *Env) {
	if c.liner == nil {
		c.liner = &linerFX{}
	}
	c.Src = &ContentSource{Liner: c.liner}
	c.cmpFX.OnInit(e)
}
