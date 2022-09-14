// Copyright (c) 2022 Stephan Lukits. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package lines

import "github.com/gdamore/tcell/v2"

// KeysRegisterer implementations are ask at initialization time if they
// want to register listeners for particular keys.  To update the
// registered keys of a component use Events.UpdateKeys.
type KeysRegisterer interface {

	// Keys implementation is provided with a callback function of which
	// each call maps a key and its modifier to an event listener.
	// Invalid calls are ignored.  A provided nil listener removes the
	// rune binding.
	Keys(KeyRegistration)
}

// RunesRegisterer implementations are ask at initialization time if they
// want to register listeners for particular runes.  To update the
// registered runes of a component use Events.UpdateRunes.
type RunesRegisterer interface {

	// Runes implementation is provided with a callback function of
	// which each call maps a rune to an event listener.  Invalid calls
	// are ignored.  A provided nil listener removes the rune binding.
	Runes(RuneRegistration)
}

// Keyer is implemented by components who want to take over the user's
// key-input if they are focused.
type Keyer interface {

	// OnKey is provided with every key-press and reported modifiers which
	// were pressed at the same time.
	OnKey(*Env, tcell.Key, tcell.ModMask)
}

// Runer is implemented by components who want to take over the user's
// rune-input if they are focused.
type Runer interface {

	// OnRune is provided with every rune-input.  NOTE modifiers and
	// runes are not really a thing when it comes to terminals.  Hence
	// they are mainly ignored by lines.  But the provided Env
	// instance's Evt-Property provides the original received
	// *tcell.EventKey instance which holds all information which tcell
	// provided about reported rune event.
	OnRune(*Env, rune)
}

func registerKeys(cmp Componenter, cntx *rprContext) {
	if !cmp.hasLayoutWrapper() {
		return
	}
	wrapped := cmp.layoutComponent().wrapped()
	if wrapped.lst == nil {
		wrapped.lst = &listeners{}
	}
	kc, ok := cmp.(KeysRegisterer)
	if ok {
		kc.Keys(func(k tcell.Key, mm tcell.ModMask, l Listener) {
			wrapped.lst.key(k, mm, l)
		})
	}
}

func registerRunes(cmp Componenter, cntx *rprContext) {
	if !cmp.hasLayoutWrapper() {
		return
	}
	wrapped := cmp.layoutComponent().wrapped()
	if wrapped.lst == nil {
		wrapped.lst = &listeners{}
	}
	rc, ok := cmp.(RunesRegisterer)
	if ok {
		rc.Runes(func(r rune, l Listener) {
			wrapped.lst.rune(r, l)
		})
	}
}

func reportKey(cntx *rprContext) (quit bool) {
	evt := cntx.evt.(*tcell.EventKey)
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
	execFeature(cntx, evt)
	if !cntx.scr.root().ff.keyQuits(evt.Key()) {
		return false
	}
	reportQuit(cntx)
	return true
}

func reportKeyListener(
	c layoutComponenter, evt *tcell.EventKey, cntx *rprContext,
) (stopBubbling bool) {
	if c.wrapped().lst == nil {
		return false
	}
	l, ok := c.wrapped().lst.keyListenerOf(evt.Key(), evt.Modifiers())
	if !ok {
		return false
	}
	env := callback(c.userComponent(), cntx, l)
	return env&envStopBubbling == envStopBubbling
}

func keyCurry(
	evt *tcell.EventKey, cb func(*Env, tcell.Key, tcell.ModMask),
) func(*Env) {
	return func(e *Env) {
		cb(e, evt.Key(), evt.Modifiers())
	}
}

func execFeature(cntx *rprContext, evt *tcell.EventKey) {
	usr := cntx.scr.focus.userComponent()
	var f FeatureMask
	if evt.Key() == tcell.KeyRune {
		f = usr.layoutComponent().wrapped().ff.runeFeature(evt.Rune())
	} else {
		f = usr.layoutComponent().wrapped().ff.keyFeature(
			evt.Key(), evt.Modifiers())
	}
	if f == NoFeature {
		return
	}
	usr.enable()
	defer usr.disable()
	execute(cntx, usr, f)
}

func reportOnKey(
	c layoutComponenter, evt *tcell.EventKey, cntx *rprContext,
) (stopBubbling bool) {
	kyr, ok := c.userComponent().(Keyer)
	if !ok {
		return false
	}
	env := callback(c.userComponent(), cntx, keyCurry(evt, kyr.OnKey))
	return env&envStopBubbling == envStopBubbling
}

func reportRune(cntx *rprContext) (quit bool) {
	evt := cntx.evt.(*tcell.EventKey)
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
	execFeature(cntx, evt)
	if !cntx.scr.root().ff.runeQuits(evt.Rune()) {
		return false
	}
	reportQuit(cntx)
	return true
}

func reportRuneListener(
	c layoutComponenter, evt *tcell.EventKey, cntx *rprContext,
) (stopBubbling bool) {
	if c.wrapped().lst == nil {
		return false
	}
	l, ok := c.wrapped().lst.runeListenerOf(evt.Rune())
	if !ok {
		return false
	}
	env := callback(c.userComponent(), cntx, l)
	return env&envStopBubbling == envStopBubbling
}

func runeCurry(
	evt *tcell.EventKey, cb func(*Env, rune),
) func(*Env) {
	return func(e *Env) { cb(e, evt.Rune()) }
}

func reportOnRune(
	c layoutComponenter, evt *tcell.EventKey, cntx *rprContext,
) (stopBubbling bool) {
	rnr, ok := c.userComponent().(Runer)
	if !ok {
		return false
	}
	env := callback(c.userComponent(), cntx, runeCurry(evt, rnr.OnRune))
	return env&envStopBubbling == envStopBubbling
}
