// Copyright (c) 2022 Stephan Lukits. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package lines

import (
	"github.com/gdamore/tcell/v2"
)

// Listener is the most common type of event listener: a callback
// provided with an environment.
type Listener func(*Env)

// KeyRegistration is the callback provided to a Keyer implementation of
// a component to register specific keys to specific key-listeners.
type KeyRegistration func(tcell.Key, tcell.ModMask, Listener)

// RuneRegistration is the callback provided to a Runer implementation
// of a component to register specific runes to specific rune-listeners
type RuneRegistration func(rune, Listener)

// listeners hold a components event-listers for particular key or rune
// events.
type listeners struct {
	kk map[tcell.ModMask]map[tcell.Key]Listener
	rr map[rune]Listener
}

// key registers provided listener for given key/mode combination
// respectively removes the registration for given key/mode if the
// listener is nil.  key fails if already a listener is registered for
// given key/mode or if the zero key is given or if given key is
// associated with the quit-feature.  NOTE use *Quit* at an
// Register-instance to receive the Quit-event.
func (ll *listeners) key(k tcell.Key, m tcell.ModMask, l Listener) {
	if ll.kk == nil {
		ll.kk = map[tcell.ModMask]map[tcell.Key]Listener{}
	}
	if l == nil {
		if ll.kk[m] != nil {
			delete(ll.kk[m], k)
		}
		return
	}
	if k == tcell.KeyNUL {
		return
	}
	if ll.kk[m] == nil {
		ll.kk[m] = map[tcell.Key]Listener{k: l}
		return
	}
	ll.kk[m][k] = l
}

// keyListenerOf returns the listener registered for given key/mode
// combination.  The second return value is false if no listener is
// registered for given key.
func (ll *listeners) keyListenerOf(
	k tcell.Key, m tcell.ModMask,
) (Listener, bool) {

	if ll.kk == nil {
		return nil, false
	}

	if _, ok := ll.kk[m]; !ok {
		return nil, false
	}

	l, ok := ll.kk[m][k]
	return l, ok
}

// rune registers provided listener for given rune respectively removes
// the registration for given rune if the listener is nil.
func (ll *listeners) rune(r rune, l Listener) error {

	if ll.rr == nil {
		ll.rr = map[rune]Listener{}
	}

	if l == nil {
		delete(ll.rr, r)
		return nil
	}

	if r == rune(0) {
		return nil
	}

	ll.rr[r] = l
	return nil
}

// runeListenerOf returns the listener registered for given rune.  The
// second return value is false if no listener is registered for given
// rune.
func (kk *listeners) runeListenerOf(r rune) (Listener, bool) {

	if kk.rr == nil {
		return nil, false
	}

	l, ok := kk.rr[r]
	return l, ok
}
