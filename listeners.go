// Copyright (c) 2022 Stephan Lukits. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package lines

// Listener is the most common type of event listener: a callback
// provided with an environment.
type Listener func(*Env)

// KeyRegistration is the callback provided to a Keyer implementation of
// a component to register specific keys to specific key-listeners.
type KeyRegistration func(Key, Modifier, Listener)

// RuneRegistration is the callback provided to a Runer implementation
// of a component to register specific runes to specific rune-listeners
type RuneRegistration func(rune, Listener)

// Listeners provide the API to register key and rune event listeners.
// While lines will *not* panic if this API is used outside an
// event-listener callback registering methods are not concurrency save.
type Listeners struct {
	c *Component
}

func (ll *Listeners) listeners() *listeners {
	if ll.c.layoutCmp == nil {
		return nil
	}
	wrapped := ll.c.layoutCmp.wrapped()
	wrapped.ensureListeners()
	return wrapped.lst
}

func (ll *Listeners) Key(k Key, m Modifier, l Listener) {
	cll := ll.listeners()
	if cll == nil {
		return
	}
	cll.key(k, m, l)
}

func (ll *Listeners) Rune(r rune, m Modifier, l Listener) {
	cll := ll.listeners()
	if cll == nil {
		return
	}
	cll.rune(r, m, l)
}

// listeners hold a components event-listers for particular key or rune
// events.
type listeners struct {
	kk map[Modifier]map[Key]Listener
	rr map[Modifier]map[rune]Listener
}

// key registers provided listener for given key/mode combination
// respectively removes the registration for given key/mode if the
// listener is nil.  key fails if already a listener is registered for
// given key/mode or if the zero key is given or if given key is
// associated with the quit-feature.  NOTE use *Quit* at an
// Register-instance to receive the Quit-event.
func (ll *listeners) key(k Key, m Modifier, l Listener) {
	if ll.kk == nil {
		ll.kk = map[Modifier]map[Key]Listener{}
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
func (ll *listeners) keyListenerOf(k Key, m Modifier) (Listener, bool) {

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
func (ll *listeners) rune(r rune, m Modifier, l Listener) {

	if ll.rr == nil {
		ll.rr = map[Modifier]map[rune]Listener{}
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
func (ll *listeners) runeListenerOf(r rune, m Modifier) (Listener, bool) {

	if ll.rr == nil {
		return nil, false
	}

	if _, ok := ll.rr[m]; !ok {
		return nil, false
	}

	l, ok := ll.rr[m][r]
	return l, ok
}
