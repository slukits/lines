// Copyright (c) 2022 Stephan Lukits. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package lines

import (
	"bytes"

	"github.com/gdamore/tcell/v2"
	"github.com/slukits/lines/internal/lyt"
)

// Component enables a user implemented UI-component to be processed by
// lines, i.e. all user UI-components which are provided to lines must
// embed this type.  NOTE accessing features provided by an embedded
// Component instance is only supported during an event reporting
// callback to its embedding component; otherwise the embedded Component
// instance will panic.
type Component struct {

	// FF provides access and fine grained control over a components
	// end-user features (see FeatureMask).
	FF *Features

	// Scroll provides a component's API for scrolling.
	Scroll *Scroller

	// component provides properties/features of an Component.  A
	// Component can't do it directly if it should panic if it is used
	// outside an event reporting callback.  layoutCmp wraps the actual
	// component instance hence it is available for internal use and if
	// this Component is enabled its *components points to the
	// layoutCmp wrapped component.  If it is disabled *component is set
	// to nil.
	*component

	// layoutCmp wraps a *component-instance in a way that it can be
	// used by the layout-manager and that it provides user defined
	// components.  I.e. the problem here is that an embedded Component
	// instance might be enabled or not (*component is nil or not).
	// Thus we can't just pass on a user defined component to the layout
	// manager because its access to the components dimensions which are
	// provided by the wrapped *component might panic.  On the other
	// hand  we need to provide the user provided components
	// (ForStacked, ForChained) to the layout manager.  A
	// layoutComponenter combines a Component's wrapped (non nil)
	// *component with the ForStacked/ForChained method of the user
	// defined associated component.  Emitting other layoutComponenter
	// to the layout manager.
	layoutCmp layoutComponenter
}

type ComponentMode uint

const (

	// Overwriting a components content by an write operation.
	Overwriting ComponentMode = 1 << iota

	// Appending to a components content by an write operation.
	Appending

	// Tailing is appending and displaying the contents "tail"
	// especially if the display area cannot show all the content.
	Tailing
)

func (c *Component) initialize(
	userComponent Componenter,
) layoutComponenter {

	if c.layoutCmp != nil {
		return c.layoutCmp
	}

	inner := &component{
		dim:     lyt.DimFilling(1, 1),
		ll:      &lines{},
		sty:     tcell.StyleDefault,
		userCmp: userComponent,
		mod:     Overwriting,
	}
	c.FF = &Features{c: c}
	switch userComponent.(type) {
	case Stacker:
		c.layoutCmp = &stackingWrapper{component: inner}
	case Chainer:
		c.layoutCmp = &chainingWrapper{component: inner}
	default:
		c.layoutCmp = inner
	}
	return c.layoutCmp
}

// enable component for client usage.
func (c *Component) enable() { c.component = c.layoutCmp.wrapped() }

// disable component for client usage.
func (c *Component) disable() { c.component = nil }

// isInitialized returns true if embedded component-instance is wrapped
// in a layout component and has been initialized.
func (c *Component) isInitialized() bool {
	if c.layoutCmp == nil {
		return false
	}
	return c.layoutCmp.wrapped().initialized
}

func (c *Component) hasLayoutWrapper() bool {
	return c.layoutCmp != nil
}

func (c *Component) layoutComponent() layoutComponenter {
	return c.layoutCmp
}

// component is the actual implementation of a lines-Component.
type component struct {
	userCmp     Componenter
	dim         *lyt.Dim
	mod         ComponentMode
	initialized bool
	ll          *lines
	sty         tcell.Style
	lst         *listeners
	ff          *features
}

// Mod sets how given components content is maintained.
func (c *component) Mod(cm ComponentMode) {
	switch cm {
	case Appending:
		c.mod &^= Overwriting | Tailing
		c.mod |= Appending
	case Overwriting:
		c.mod &^= Appending | Tailing
		c.mod |= Overwriting
	case Tailing:
		c.mod &^= Appending | Overwriting
		c.mod |= Tailing
	}
}

// Len returns the number of lines currently stored in a component.
// Note the line number is independent of a component's associated
// screen area.
func (c *component) Len() int {
	return len(*c.ll)
}

// Dim provides a components layout dimensions and features to adapt
// them.
func (c *component) Dim() *lyt.Dim { return c.dim }

// Reset blanks out the content of the line with given index the next
// time it is printed to the screen.
func (c *component) Reset(idx int) {
	if idx < 0 || idx >= len(*c.ll) {
		return
	}
	(*c.ll)[idx].set("")
}

func (c *component) setInitialized() {
	c.initialized = true
}

// component gets the component out of a layoutComponenter without using
// a type-switch.
func (c *component) wrapped() *component { return c }

func (c *component) userComponent() Componenter {
	return c.userCmp
}

func (c *component) ensureFeatures() {
	if c.ff != nil {
		return
	}
	c.ff = defaultFeatures.copy()
}

// hardSync clears the screen area of receiving component before its
// content is written to the screen.
func (c *component) hardSync(rw runeWriter) {
	c.clear(rw)
	c.sync(rw)
}

// sync writes receiving components lines to the screen.
func (c *component) sync(rw runeWriter) {
	sx, sy, sw, sh := c.dim.Area()
	if c.mod&Tailing == Tailing {
		if len(*c.ll) >= sh {
			c.syncTailed(rw, sx, sy, sw, sh)
			return
		}
	}
	c.ll.For(func(i int, l *line) (stop bool) {
		if i >= sh {
			return true
		}
		l.sync(sx, sy+i, sw, rw, c.sty)
		return false
	})
}

func (c *component) syncTailed(rw runeWriter, sx, sy, sw, sh int) {
	y := sy + sh - 1
	c.ll.ForInverse(func(_ int, l *line) (stop bool) {
		if y < sy {
			return true
		}
		l.sync(sx, y, sw, rw, c.sty)
		y--
		return false
	})
}

// clear fills the receiving component's printable area with spaces.
func (c *component) clear(rw runeWriter) {
	sx, sy, sw, sh := c.dim.Area()
	for y := sy; y < sh; y++ {
		for x := sx; x < sw; x++ {
			rw.SetContent(x, y, ' ', nil, c.sty)
		}
	}
}

func (c *component) write(bb []byte, at int) (int, error) {
	if at > -1 {
		return c.writeAt(bb, at)
	}
	switch {
	case c.mod&Overwriting == Overwriting:
		c.ll.replace(bytes.Split(bb, []byte("\n"))...)
	case c.mod&(Appending|Tailing) != 0:
		c.ll.append(bytes.Split(bb, []byte("\n"))...)
	}
	return len(bb), nil
}

func (c *component) writeAt(bb []byte, at int) (int, error) {
	c.ll.replaceAt(at, bytes.Split(bb, []byte("\n"))...)
	return len(bb), nil
}

// Features provides access and fine grained control over a components
// (end-user) features provided by lines.  Its methods will panic used
// outside an event reporting listener-callback.
type Features struct{ c *Component }

func (ff *Features) ensureInitialized() *features {
	ff.c.ensureFeatures()
	return ff.c.ff
}

// Add adds the default key, rune and button bindings of given
// feature(s) for associated component.
func (ff *Features) Add(f FeatureMask) {
	ff.ensureInitialized().add(f, false)
}

// AddRecursive sets the default key, rune and button bindings of given
// feature(s) for associated component.  Whereas the feature(s) are
// flagged recursive, i.e. they apply as well for nested components.
func (ff *Features) AddRecursive(f FeatureMask) {
	ff.ensureInitialized().add(f, true)
}

// Has returns true if receiving component features have key, rune or
// button bindings for given feature(s)
func (ff *Features) Has(f FeatureMask) bool {
	return ff.ensureInitialized().has(f)
}

// All returns all features for which currently key, rune or button
// bindings are registered. (note Has is faster to determine if a
// particular feature is set.)
func (ff *Features) All() FeatureMask {
	return ff.ensureInitialized().all()
}

// KeysOf returns the keys with their modifiers bound to given feature
// of associated component.
func (ff *Features) KeysOf(f FeatureMask) FeatureKeys {
	return ff.ensureInitialized().keysOf(f)
}

// SetKeysOf deletes all set keys for given feature (except for Quitable
// defaults) and binds given keys to it instead.  If recursive is true
// the feature becomes applicable for nested components.  The call is
// ignored if given feature is not a power of two i.e. a single feature.
// Providing no keys simply removes all key-bindings for given feature.
func (ff *Features) SetKeysOf(
	f FeatureMask, recursive bool, kk ...FeatureKey,
) {
	ff.ensureInitialized().setKeysOf(f, recursive, kk...)
}

// ButtonsOf returns the buttons with their modifiers bound to given
// feature for associated component.
func (ff *Features) ButtonsOf(f FeatureMask) FeatureButtons {
	return ff.ensureInitialized().buttonsOf(f)
}

// SetButtonsOf deletes all set buttons for given feature and binds
// given buttons to it instead.  If recursive is true the feature
// becomes applicable for nested components.  The call is ignored if
// given feature is not a power of two i.e. a single feature.  Providing
// no buttons simply removes all button-bindings for given feature.
func (ff *Features) SetButtonsOf(
	f FeatureMask, recursive bool, bb ...FeatureButton,
) {
	ff.ensureInitialized().setButtonsOf(f, recursive, bb...)
}

// RunesOf returns the runes bound to given feature for associated
// component.
func (ff *Features) RunesOf(f FeatureMask) FeatureRunes {
	return ff.ensureInitialized().runesOf(f)
}

// SetRunesOf deletes all set runes for given feature and binds given
// runes to it instead.  If recursive is true the feature becomes
// applicable for nested components.  The call is ignored if given
// feature is not a power of two i.e. a single feature.  Providing no
// runes simply removes all runes-bindings for given feature.
func (ff *Features) SetRunesOf(
	f FeatureMask, recursive bool, rr ...rune,
) {
	ff.ensureInitialized().setRunesOf(f, recursive, rr...)
}

// Delete removes all runes, key or button bindings of given feature(s)
// except for Quitable.  The two default Quitable bindings ctrl-c and
// ctrl-d remain.  NOTE you can prevent the processing of the default
// quit bindings by adding to your root component listeners for these
// keys which call StopBubbling on their environment:
//
//	type Root struct { lines.Component }
//
//	func (c *Root) OnInit(e *lines.Env) { fmt.Fprint(e, "hello world") }
//
//	func (c *Root) Keys(register lines.KeyRegistration) {
//	    register(tcell.KeyCtrlC, tcell.ModNone, func(e *Env) {
//	        e.StopBubbling()
//	    })
//	    register(tcell.KeyCtrlD, tcell.ModNone, func(e *Env) {
//	        e.StopBubbling()
//	    })
//	}
//
//	lines.New(&Root{}).Listen()
//
// gives you an application which can't be quit by its users.
func (ff *Features) Delete(f FeatureMask) {
	ff.ensureInitialized().delete(f)
}

// layoutComponenter combines the user-provided component with its
// internally created component to have all information to layout the
// component together: Dimer-implementation is provide by the internally
// created component while the user-component potentially provides nested
// user-components.
type layoutComponenter interface {
	lyt.Dimer
	wrapped() *component
	userComponent() Componenter
}

// stackingWrapper wraps a stacking user-component for the
// layout-manager.
type stackingWrapper struct {
	*component
}

func (sw *stackingWrapper) ForStacked(cb func(lyt.Dimer) bool) {
	sw.userCmp.(Stacker).ForStacked(func(cmp Componenter) bool {
		if !cmp.hasLayoutWrapper() {
			cmp.initialize(cmp)
		}
		return cb(cmp.layoutComponent())
	})
}

// chainingWrapper wraps a chaining user-component for the
// layout-manager.
type chainingWrapper struct {
	*component
}

func (cw *chainingWrapper) ForChained(cb func(lyt.Dimer) bool) {
	cw.userCmp.(Chainer).ForChained(func(cmp Componenter) bool {
		if !cmp.hasLayoutWrapper() {
			cmp.initialize(cmp)
		}
		return cb(cmp.layoutComponent())
	})
}
