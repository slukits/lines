// Copyright (c) 2022 Stephan Lukits. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package lines

import (
	"fmt"
	"testing"

	. "github.com/slukits/gounit"
)

type kbEventType uint

const (
	updateKeyListeners kbEventType = iota
	updateRuneListeners
	keyListener
	runeListener
	onKey
	onRune
)

type kbEvent struct {
	Eventer
	evtType kbEventType
}

type kbEvents []kbEvent

func (ke *kbEvents) append(evt Eventer, evtTyp kbEventType) {
	*ke = append(*ke, kbEvent{Eventer: evt, evtType: evtTyp})
}

func (ke kbEvents) has(evtType kbEventType) bool {
	for _, r := range ke {
		if r.evtType != evtType {
			continue
		}
		return true
	}
	return false
}

func (ke kbEvents) len(evtType kbEventType) int {
	n := 0
	for _, r := range ke {
		if r.evtType != evtType {
			continue
		}
		n++
	}
	return n
}

func (ke kbEvents) forEvtOf(
	et kbEventType, cb func(kbEvent) (stop bool),
) {
	for _, r := range ke {
		if r.evtType != et {
			continue
		}
		if cb(r) {
			return
		}
	}
}

func (ke kbEvents) get(evtType kbEventType) *kbEvent {
	for _, r := range ke {
		if r.evtType != evtType {
			continue
		}
		return &r
	}
	return nil
}

func (ke kbEvents) HasKey() bool  { return ke.get(onKey) != nil }
func (ke kbEvents) HasRune() bool { return ke.get(onRune) != nil }
func (ke kbEvents) HasKeyListener() bool {
	return ke.get(keyListener) != nil
}
func (ke kbEvents) HasRuneListener() bool {
	return ke.get(runeListener) != nil
}

type kbFX struct {
	kbEvents
	stopBubblingKeys, stopBubblingRunes   bool
	stopBubblingOnKey, stopBubblingOnRune bool
}

func (fx *kbFX) OnKey(e *Env, k Key, mm Modifier) {
	fx.append(e.Evt, onKey)
	if fx.stopBubblingOnKey {
		e.StopBubbling()
	}
}

func (fx *kbFX) OnRune(e *Env, r rune, mm Modifier) {
	fx.append(e.Evt, onRune)
	if fx.stopBubblingOnRune {
		e.StopBubbling()
	}
}

// registerKeys registers with each call the "next" F*-key and removes the
// registration of the "previous" one; we start at F1 and panic at F64.
func (fx *kbFX) registerKey(register *Listeners) {
	delta := fx.len(updateKeyListeners)
	if delta > 63 {
		panic("can't register bigger than F64")
	}
	if delta > 0 {
		// remove previously registered listener
		register.Key(F1+(Key(delta)-1), ZeroModifier, nil)
	}
	// register next F-listener
	register.Key(F1+Key(delta), ZeroModifier,
		func(e *Env) {
			fx.append(e.Evt, keyListener)
			if fx.stopBubblingKeys {
				e.StopBubbling()
			}
		},
	)
	// increase len(updateKeyListeners)
	fx.append(nil, updateKeyListeners)
}

// registerRunes registers with each call the "next" rune and removes
// the registration of the "previous" one; we start at 'A' and can go
// til 'Z'.
func (fx *kbFX) registerRune(register *Listeners) {
	delta := fx.len(updateRuneListeners)
	if delta >= 25 {
		delta += 6
	}
	if delta > 63 {
		panic("can't register bigger than Z")
	}
	if delta > 0 {
		// remove previously registered rune listener
		register.Rune('A'+(rune(delta)-1), ZeroModifier, nil)
	}
	// register next rune-listener
	register.Rune('A'+rune(delta), ZeroModifier, func(e *Env) {
		fx.append(e.Evt, runeListener)
		if fx.stopBubblingRunes {
			e.StopBubbling()
		}
	})
	// increase len(updateRuneListeners)
	fx.append(nil, updateRuneListeners)
}

type kbCmpFX struct {
	Component
	kbFX
}

func (fx *kbCmpFX) OnInit(e *Env) {
	fx.registerKey(fx.Register)
	fx.registerRune(fx.Register)
}

type KB struct{ Suite }

func (s *KB) tt(t *T, cmp Componenter) *Testing {
	return TermFixture(t.GoT(), 0, cmp)
}

func (s *KB) Reports_to_keyer_implementation(t *T) {
	fx := &kbCmpFX{}
	tt := s.tt(t, fx)
	t.Not.True(fx.HasKey())
	tt.FireKey(Enter, Shift)
	t.True(fx.HasKey())
	fx.forEvtOf(onKey, func(ke kbEvent) (stop bool) {
		evt := ke.Eventer.(KeyEventer)
		t.Eq(Enter, evt.Key())
		t.Eq(Shift, evt.Mod())
		return true
	})
}

func (s *KB) Reports_to_runer_implementation(t *T) {
	fx := &kbCmpFX{}
	tt := s.tt(t, fx)
	t.Not.True(fx.HasRune())
	tt.FireRune('r', Alt)
	t.True(fx.HasRune())
	fx.forEvtOf(onRune, func(ke kbEvent) (stop bool) {
		evt := ke.Eventer.(RuneEventer)
		t.Eq('r', evt.Rune())
		t.Eq(Alt, evt.Mod())
		return true
	})
}

type bbbKBCmpFX struct {
	Component
	Stacking
	kbFX
}

func (c *bbbKBCmpFX) OnInit(e *Env) {
	c.CC = append(c.CC, &kbCmpFX{})
	c.registerKey(c.Register)
	c.registerRune(c.Register)
}

func (c *bbbKBCmpFX) inner() *kbCmpFX {
	return c.CC[0].(*kbCmpFX)
}

func (s *KB) Bubbles_keys(t *T) {
	fx := &bbbKBCmpFX{}
	tt := s.tt(t, fx)
	t.FatalOn(fx.CC[0].(*kbCmpFX).Focus())
	tt.FireKey(F1)
	t.True(fx.inner().HasKey())
	t.True(fx.inner().len(keyListener) == 1)
	t.True(fx.HasKey())
	t.True(fx.len(keyListener) == 1)
}

func (s *KB) Bubbles_runes(t *T) {
	fx := &bbbKBCmpFX{}
	tt := s.tt(t, fx)
	t.FatalOn(fx.CC[0].(*kbCmpFX).Focus())
	tt.FireRune('A')
	t.True(fx.inner().HasRune())
	t.True(fx.inner().len(runeListener) == 1)
	t.True(fx.HasRune())
	t.True(fx.len(runeListener) == 1)
}

func (s *KB) Event_bubbling_may_be_stopped(t *T) {
	fx := &bbbKBCmpFX{}
	tt := s.tt(t, fx)
	t.FatalOn(fx.CC[0].(*kbCmpFX).Focus())
	fx.inner().stopBubblingKeys = true
	fx.inner().stopBubblingOnKey = true
	fx.inner().stopBubblingRunes = true
	fx.inner().stopBubblingOnRune = true
	tt.FireKey(F1)   // reports keyListener only
	tt.FireRune('A') // reports runeListener only
	tt.FireKey(BS)   // report OnKey only
	tt.FireRune('a') // report OnRune only
	t.True(fx.inner().HasKey())
	t.True(fx.inner().HasRune())
	t.True(fx.inner().HasKeyListener())
	t.True(fx.inner().HasRuneListener())
	t.Not.True(fx.HasKey())
	t.Not.True(fx.HasRune())
	t.Not.True(fx.HasKeyListener())
	t.Not.True(fx.HasRuneListener())
}

type icmpFX struct {
	Component
	init func(*icmpFX, *Env)
}

func (c *icmpFX) OnInit(e *Env) {
	if c.init == nil {
		return
	}
	c.init(c, e)
}

func (s *KB) Executes_key_feature(t *T) {
	tt := s.tt(t, &icmpFX{init: func(c *icmpFX, e *Env) {
		c.FF.Add(Scrollable)
		c.Dim().SetHeight(2)
		fmt.Fprint(e, "first\nsecond\nthird\nforth")
	}})
	up := defaultBindings[UpScrollable].kk[0]
	down := defaultBindings[DownScrollable].kk[0]

	tt.FireKey(down.Key, down.Mod)
	t.Eq("second\nthird ", tt.ScreenOf(tt.Root()).Trimmed().String())
	tt.FireKey(down.Key, down.Mod)
	t.Eq("third\nforth", tt.Screen().Trimmed().String())
	tt.FireKey(up.Key, up.Mod)
	tt.FireKey(up.Key, up.Mod)
	t.Eq("first \nsecond", tt.Screen().Trimmed().String())
}

func TestKB(t *testing.T) {
	t.Parallel()
	Run(&KB{}, t)
}
