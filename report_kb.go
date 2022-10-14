// Copyright (c) 2022 Stephan Lukits. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package lines

import (
	"github.com/slukits/lines/internal/api"
)

// Keyer is implemented by components who want to take over the user's
// key-input if they are focused.
type Keyer interface {

	// OnKey is provided with every key-press and reported modifiers which
	// were pressed at the same time.
	OnKey(*Env, Key, Modifier)
}

// Runer is implemented by components who want to take over the user's
// rune-input if they are focused.
type Runer interface {

	// OnRune is provided with every rune-input.
	OnRune(*Env, rune, Modifier)
}

func reportKey(cntx *rprContext, evt api.KeyEventer) (quit bool) {
	sb := false
	stopBubbling := func() bool {
		sb = true
		return true
	}
	cntx.scr.forFocused(func(c layoutComponenter) (stop bool) {
		if sb := reportKeyListener(c, evt, cntx); sb {
			return stopBubbling()
		}
		if sb := reportOnKey(c, evt, cntx); sb {
			return stopBubbling()
		}
		return false
	})
	if sb {
		return false
	}
	execKeyFeature(cntx, evt)
	return cntx.scr.root().ff.keyQuits(evt.Key())
}

func reportKeyListener(
	c layoutComponenter, evt api.KeyEventer, cntx *rprContext,
) (stopBubbling bool) {
	if c.wrapped().lst == nil {
		return false
	}
	l, ok := c.wrapped().lst.keyListenerOf(evt.Key(), evt.Mod())
	if !ok {
		return false
	}
	env := callback(c.userComponent(), cntx, l)
	return env&envStopBubbling == envStopBubbling
}

func keyCurry(
	evt api.KeyEventer, cb func(*Env, api.Key, api.Modifier),
) func(*Env) {
	return func(e *Env) {
		cb(e, evt.Key(), evt.Mod())
	}
}

func execKeyFeature(cntx *rprContext, evt api.KeyEventer) {
	usr := cntx.scr.focus.userComponent()
	f := usr.layoutComponent().wrapped().ff.keyFeature(
		evt.Key(), evt.Mod())
	if f == NoFeature {
		return
	}
	usr.enable()
	defer usr.disable()
	execute(cntx, usr, f)
}

func reportOnKey(
	c layoutComponenter, evt api.KeyEventer, cntx *rprContext,
) (stopBubbling bool) {
	kyr, ok := c.userComponent().(Keyer)
	if !ok {
		return false
	}
	env := callback(c.userComponent(), cntx, keyCurry(evt, kyr.OnKey))
	return env&envStopBubbling == envStopBubbling
}

func reportRune(cntx *rprContext, evt RuneEventer) (quit bool) {
	sb := false
	stopBubbling := func() bool {
		sb = true
		return true
	}
	cntx.scr.forFocused(func(c layoutComponenter) (stop bool) {
		if sb := reportRuneListener(c, evt, cntx); sb {
			return stopBubbling()
		}
		if sb := reportOnRune(c, evt, cntx); sb {
			return stopBubbling()
		}
		return false
	})
	if sb {
		return false
	}
	execRuneFeature(cntx, evt)
	return cntx.scr.root().ff.runeQuits(evt.Rune())
}

func reportRuneListener(
	c layoutComponenter, evt RuneEventer, cntx *rprContext,
) (stopBubbling bool) {
	if c.wrapped().lst == nil {
		return false
	}
	l, ok := c.wrapped().lst.runeListenerOf(evt.Rune(), evt.Mod())
	if !ok {
		return false
	}
	env := callback(c.userComponent(), cntx, l)
	return env&envStopBubbling == envStopBubbling
}

func runeCurry(
	evt RuneEventer, cb func(*Env, rune, Modifier),
) func(*Env) {
	return func(e *Env) { cb(e, evt.Rune(), evt.Mod()) }
}

func reportOnRune(
	c layoutComponenter, evt RuneEventer, cntx *rprContext,
) (stopBubbling bool) {
	rnr, ok := c.userComponent().(Runer)
	if !ok {
		return false
	}
	env := callback(c.userComponent(), cntx, runeCurry(evt, rnr.OnRune))
	return env&envStopBubbling == envStopBubbling
}

func execRuneFeature(cntx *rprContext, evt api.RuneEventer) {
	usr := cntx.scr.focus.userComponent()
	f := usr.layoutComponent().wrapped().ff.runeFeature(
		evt.Rune(), evt.Mod())
	if f == NoFeature {
		return
	}
	usr.enable()
	defer usr.disable()
	execute(cntx, usr, f)
}
