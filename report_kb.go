// Copyright (c) 2022 Stephan Lukits. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package lines

import (
	"github.com/slukits/lines/internal/api"
)

// reportKey branches of to reportKeyEdit in case the currently focused
// component has an active Editor otherwise the key event triggers first
// potential listener calls followed by OnKey-implementations of all
// nested focused components.  Finally if bubbling wasn't stopped
// registered key-features are executed.
func reportKey(cntx *rprContext, evt api.KeyEventer) {
	sb := false
	stopBubbling := func() bool {
		sb = true
		return true
	}
	if cntx.scr.focus.userComponent().embedded().Edit.IsActive() {
		if !reportKeyEdit(cntx.scr.focus, evt, cntx) {
			return
		}
	}
	cntx.scr.forFocused(func(c layoutComponenter) (stop bool) {
		if sb := reportKeyListener(c, evt, cntx); sb {
			return stopBubbling()
		}
		if sb := reportOnKey(c, evt, cntx); sb {
			return stopBubbling()
		}
		return
	})
	if sb {
		return
	}
	execKeyFeature(cntx, evt)
}

// reportKeyEdit reports first to OnKey implementations then the event
// is mapped to an Edit edt.  Is the later not possible the event is
// passed on to be executed by a potential key-feature.  Otherwise edt
// is reported to an potential OnEdit implementation.  Both listeners
// (OnKey/OnEdit) can prevent further processing of the key event
// through their return values.  If not the edit on the focused
// component's content is performed.
func reportKeyEdit(
	lc layoutComponenter, evt KeyEventer, cntx *rprContext,
) (proceed bool) {
	if stopBbl := reportOnKey(lc, evt, cntx); stopBbl {
		return
	}
	usr := lc.userComponent()
	editor := usr.embedded().Edit
	if editor == nil {
		panic("lines: report: on-edit: editor missing")
	}
	usr.enable()
	defer usr.disable()
	edt := editor.newKeyEdit(evt)
	if edt == nil {
		return true
	}
	reportEdit(cntx, usr, edt)
	return false
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
	evt api.KeyEventer, cb func(*Env, api.Key, api.ModifierMask),
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

func reportRune(cntx *rprContext, evt RuneEventer) {
	sb := false
	stopBubbling := func() bool {
		sb = true
		return true
	}
	if cntx.scr.focus.userComponent().embedded().Edit.IsActive() {
		reportRuneEdit(cntx.scr.focus, evt, cntx)
		return
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
		return
	}
	execRuneFeature(cntx, evt)
}

func reportRuneEdit(
	lc layoutComponenter, evt RuneEventer, cntx *rprContext,
) {
	usr := lc.userComponent()
	editor := usr.embedded().Edit
	if editor == nil {
		panic("lines: report: on-edit: editor missing")
	}
	usr.enable()
	defer usr.disable()
	edt := editor.newRuneEdit(evt)
	if edt == nil {
		return
	}
	reportEdit(cntx, usr, edt)
}

func reportRuneListener(
	lc layoutComponenter, evt RuneEventer, cntx *rprContext,
) (stopBubbling bool) {
	c := lc.wrapped()
	if c.lst == nil {
		return false
	}
	l, ok := c.lst.runeListenerOf(evt.Rune(), evt.Mod())
	if !ok {
		return false
	}
	env := callback(lc.userComponent(), cntx, l)
	return env&envStopBubbling == envStopBubbling
}

func runeCurry(
	evt RuneEventer, cb func(*Env, rune, ModifierMask),
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
