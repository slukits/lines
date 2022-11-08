// Copyright (c) 2022 Stephan Lukits. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package lines

// A Listener is the most common type of event listener: a callback
// provided with an environment.
type Listener func(*Env)

// Listeners provide the API to register key and rune event listeners
// for a component.  Use a [Component]'s Register property to obtain its
// Listeners instance.  Listeners methods will panic if they are not
// used within an event listener callback.
type Listeners struct {
	c *Component
}

func (ll *Listeners) listeners() *listeners {
	if ll.c.layoutCmp == nil {
		return nil
	}
	ll.c.ensureListeners()
	return ll.c.lst
}

// Key adds to given listeners ll a new listener l which is notified on
// a key event with given key k and modifiers m.  Is l nil the key
// binding is removed.  An already registered listener for k and m is
// overwritten.
func (ll *Listeners) Key(k Key, m ModifierMask, l Listener) {
	cll := ll.listeners()
	if cll == nil {
		return
	}
	cll.key(k, m, l)
}

// Rune adds to given listeners ll a new listener l which is notified on
// a rune event with given rune r and modifiers m.  Is l nil the key
// binding is removed.  An already registered listener for k and m is
// overwritten.
func (ll *Listeners) Rune(r rune, m ModifierMask, l Listener) {
	cll := ll.listeners()
	if cll == nil {
		return
	}
	cll.rune(r, m, l)
}

// listeners hold a components event-listers for particular key or rune
// events.
type listeners struct {
	kk map[ModifierMask]map[Key]Listener
	rr map[ModifierMask]map[rune]Listener
}

// key registers provided listener for given key/mode combination
// respectively removes the registration for given key/mode if the
// listener is nil.  key fails if already a listener is registered for
// given key/mode or if the zero key is given or if given key is
// associated with the quit-feature.
func (ll *listeners) key(k Key, m ModifierMask, l Listener) {
	if ll.kk == nil {
		ll.kk = map[ModifierMask]map[Key]Listener{}
	}
	if l == nil {
		if ll.kk[m] != nil {
			delete(ll.kk[m], k)
		}
		return
	}
	if k == NUL {
		return
	}
	if ll.kk[m] == nil {
		ll.kk[m] = map[Key]Listener{k: l}
		return
	}
	ll.kk[m][k] = l
}

// keyListenerOf returns the listener registered for given key/mode
// combination.  The second return value is false if no listener is
// registered for given key.
func (ll *listeners) keyListenerOf(k Key, m ModifierMask) (Listener, bool) {

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
func (ll *listeners) rune(r rune, m ModifierMask, l Listener) {

	if ll.rr == nil {
		ll.rr = map[ModifierMask]map[rune]Listener{}
	}

	if l == nil {
		if ll.rr[m] != nil {
			delete(ll.rr[m], r)
		}
		return
	}

	if r == rune(0) {
		return
	}

	if ll.rr[m] == nil {
		ll.rr[m] = map[rune]Listener{r: l}
		return
	}

	ll.rr[m][r] = l
}

// runeListenerOf returns the listener registered for given rune.  The
// second return value is false if no listener is registered for given
// rune.
func (ll *listeners) runeListenerOf(r rune, m ModifierMask) (Listener, bool) {

	if ll.rr == nil {
		return nil, false
	}

	if _, ok := ll.rr[m]; !ok {
		return nil, false
	}

	l, ok := ll.rr[m][r]
	return l, ok
}
