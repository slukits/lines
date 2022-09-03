// Copyright (c) 2022 Stephan Lukits. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package lines

import (
	"fmt"
	"testing"

	"github.com/gdamore/tcell/v2"
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
	tcell   *tcell.EventKey
	evtType kbEventType
}

type kbEvents []kbEvent

func (ke *kbEvents) append(evt *tcell.EventKey, evtTyp kbEventType) {
	*ke = append(*ke, kbEvent{tcell: evt, evtType: evtTyp})
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

func (ke kbEvents) HasUpdatedKeyListener() bool {
	return ke.get(updateKeyListeners) != nil
}
func (ke kbEvents) HasUpdatedRuneListener() bool {
	return ke.get(updateRuneListeners) != nil
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

// Keys should be called by lines at least once during initialization.
// Doesn't affect the event-countdown.  Initially the F1-key will be
// registered.  Events.UpdateKeys should trigger an other call to Keys
// which will remove the F1-listener and add an F2-listener and so on
// until F64, than we panic.
func (fx *kbFX) Keys(register KeyRegistration) {
	delta := fx.len(updateKeyListeners)
	if delta > 63 {
		panic("can't register bigger than F64")
	}
	if delta > 0 {
		// remove previously registered listener
		register(tcell.KeyF1+(tcell.Key(delta)-1), tcell.ModNone, nil)
	}
	// register next F-listener
	register(tcell.KeyF1+tcell.Key(delta), tcell.ModNone,
		func(e *Env) {
			fx.append(e.Evt.(*tcell.EventKey), keyListener)
			if fx.stopBubblingKeys {
				e.StopBubbling()
			}
		},
	)
	// increase len(updateKeyListeners)
	fx.append(nil, updateKeyListeners)
}

// Runes should be called by lines at least once during initialization.
// Doesn't affect the event-countdown.  Initially the A-key will be
// registered.  Events.UpdateRunes should trigger an other call to Runes
// which will remove the A-listener and add an B-listener and so on
// until z, than we panic.
func (fx *kbFX) Runes(register RuneRegistration) {
	delta := fx.len(updateRuneListeners)
	if delta >= 25 {
		delta += 6
	}
	if delta > 63 {
		panic("can't register bigger than Z")
	}
	if delta > 0 {
		// remove previously registered rune listener
		register('A'+(rune(delta)-1), nil)
	}
	// register next rune-listener
	register('A'+rune(delta), func(e *Env) {
		fx.append(e.Evt.(*tcell.EventKey), runeListener)
		if fx.stopBubblingRunes {
			e.StopBubbling()
		}
	})
	// increase len(updateRuneListeners)
	fx.append(nil, updateRuneListeners)
}

func (fx *kbFX) OnKey(e *Env, k tcell.Key, mm tcell.ModMask) {
	fx.append(e.Evt.(*tcell.EventKey), onKey)
	if fx.stopBubblingOnKey {
		e.StopBubbling()
	}
}

func (fx *kbFX) OnRune(e *Env, r rune) {
	fx.append(e.Evt.(*tcell.EventKey), onRune)
	if fx.stopBubblingOnRune {
		e.StopBubbling()
	}
}

type kbCmpFX struct {
	Component
	kbFX
}

type KB struct{ Suite }

func (s *KB) Key_listeners_are_registered(t *T) {
	fx := &kbCmpFX{}
	ee, _ := Test(t.GoT(), fx)
	t.False(fx.HasUpdatedKeyListener())
	ee.Listen()
	ee.QuitListening()
	t.True(fx.HasUpdatedKeyListener())
}

func (s *KB) Rune_listeners_are_registered(t *T) {
	fx := &kbCmpFX{}
	ee, _ := Test(t.GoT(), fx)
	t.False(fx.HasUpdatedRuneListener())
	ee.Listen()
	ee.QuitListening()
	t.True(fx.HasUpdatedRuneListener())
}

func (s *KB) Key_listeners_are_updated(t *T) {
	fx := &kbCmpFX{}
	ee, tt := Test(t.GoT(), fx, 2)
	defer ee.QuitListening()
	ee.UpdateKeys(fx) // deletes F1, registers F2
	tt.FireKey(tcell.KeyF1, 0)
	t.True(fx.len(runeListener) == 0)
	t.True(fx.len(updateKeyListeners) == 2)
}

func (s *KB) Rune_listeners_are_updated(t *T) {
	fx := &kbCmpFX{}
	ee, tt := Test(t.GoT(), fx, 2)
	defer ee.QuitListening()
	ee.UpdateRunes(fx) // deletes 'a', registers 'b'
	tt.FireRune('a')
	t.True(fx.len(keyListener) == 0)
	t.True(fx.len(updateRuneListeners) == 2)
}

func (s *KB) Key_listeners_are_called(t *T) {
	fx := &kbCmpFX{}
	ee, tt := Test(t.GoT(), fx, 4)
	tt.FireKey(tcell.KeyF1, 0)
	ee.UpdateKeys(fx) // deletes F1, registers F2
	tt.FireKey(tcell.KeyF2, 0)
	t.True(fx.len(keyListener) == 2)
	t.False(ee.IsListening())
}

func (s *KB) Rune_listeners_are_called(t *T) {
	fx := &kbCmpFX{}
	ee, tt := Test(t.GoT(), fx, 4)
	tt.FireRune('A')
	ee.UpdateRunes(fx) // deletes A, registers B
	tt.FireRune('B')
	t.True(fx.len(runeListener) == 2)
	t.False(ee.IsListening())
}

func (s *KB) Reports_key(t *T) {
	fx := &kbCmpFX{}
	ee, tt := Test(t.GoT(), fx, 2)
	tt.FireKey(tcell.KeyF1, 0)
	t.False(ee.IsListening())
	t.True(fx.HasKey())
}

func (s *KB) Reports_rune(t *T) {
	fx := &kbCmpFX{}
	ee, tt := Test(t.GoT(), fx, 2)
	tt.FireRune('A')
	t.False(ee.IsListening())
	t.True(fx.HasRune())
}

type bbbKBCmpFX struct {
	Component
	Stacking
	kbCmpFX
}

func (c *bbbKBCmpFX) OnInit(e *Env) {
	c.CC = append(c.CC, &kbCmpFX{})
	e.EE.MoveFocus(c.CC[0])
}

func (c *bbbKBCmpFX) inner() *kbCmpFX {
	return c.CC[0].(*kbCmpFX)
}

func (s *KB) Bubbles_keys(t *T) {
	fx := &bbbKBCmpFX{}
	// 2 x keyListener 2 x OnKey
	ee, tt := Test(t.GoT(), fx, 4)
	tt.FireKey(tcell.KeyF1, 0)
	t.False(ee.IsListening())
	t.True(fx.inner().HasKey())
	t.True(fx.inner().len(keyListener) == 1)
	t.True(fx.HasKey())
	t.True(fx.len(keyListener) == 1)
}

func (s *KB) Bubbles_runes(t *T) {
	fx := &bbbKBCmpFX{}
	// 2 x runeListener 2 x OnRune
	ee, tt := Test(t.GoT(), fx, 4)
	tt.FireRune('A')
	t.False(ee.IsListening())
	t.True(fx.inner().HasRune())
	t.True(fx.inner().len(runeListener) == 1)
	t.True(fx.HasRune())
	t.True(fx.len(runeListener) == 1)
}

func (s *KB) Event_bubbling_may_be_stopped(t *T) {
	fx := &bbbKBCmpFX{}
	// OnInit runeListener keyListener OnRune OnKey
	ee, tt := Test(t.GoT(), fx, 5)
	ee.Listen()
	fx.inner().stopBubblingKeys = true
	fx.inner().stopBubblingOnKey = true
	fx.inner().stopBubblingRunes = true
	fx.inner().stopBubblingOnRune = true
	tt.FireKey(tcell.KeyF1, 0) // reports keyListener only
	tt.FireRune('A')           // reports runeListener only
	tt.FireKey(tcell.KeyBS, 0) // report OnKey only
	tt.FireRune('a')           // report OnRune only
	t.False(ee.IsListening())
	t.True(fx.inner().HasKey())
	t.True(fx.inner().HasRune())
	t.True(fx.inner().HasKeyListener())
	t.True(fx.inner().HasRuneListener())
	t.False(fx.HasKey())
	t.False(fx.HasRune())
	t.False(fx.HasKeyListener())
	t.False(fx.HasRuneListener())
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

func (s *KB) Key_feature_is_executed(t *T) {
	ee, tt := Test(t.GoT(), &icmpFX{init: func(c *icmpFX, e *Env) {
		c.FF.Add(Scrollable)
		c.Dim().SetHeight(2)
		fmt.Fprint(e, "first\nsecond\nthird\nforth")
	}}, 0)
	up := defaultBindings[UpScrollable].kk[0]
	down := defaultBindings[DownScrollable].kk[0]
	defer ee.QuitListening()

	tt.FireKey(down.Key, down.Mod)
	t.Eq(tt.String(), "second\nthird")
	tt.FireKey(down.Key, down.Mod)
	t.Eq(tt.String(), "third\nforth")

	tt.FireKey(up.Key, up.Mod)
	tt.FireKey(up.Key, up.Mod)
	t.Eq(tt.String(), "first\nsecond")
}

func TestKB(t *testing.T) {
	t.Parallel()
	Run(&KB{}, t)
}
