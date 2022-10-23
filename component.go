// Copyright (c) 2022 Stephan Lukits. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package lines

import (
	"github.com/slukits/lines/internal/api"
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

	// Register provides the api to register keys and runes listeners
	// for a component.
	Register *Listeners

	// Scroll provides a component's API for scrolling.
	Scroll *Scroller

	// LL provides an API for ui-aspects of a component's lines.  Use an
	// reported event's Env-instance writers to manipulate their
	// content.
	LL *ComponentLines

	// bcknd to post Update and Focus events
	bcknd api.UIer

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
	userComponent Componenter, backend api.UIer,
) layoutComponenter {

	if c.layoutCmp != nil {
		return c.layoutCmp
	}

	c.bcknd = backend

	inner := &component{
		dim:     lyt.DimFilling(1, 1),
		ll:      &lines{},
		globals: &globals{tabWidth: 4},
		fmt:     llFmt{sty: backend.NewStyle()},
		userCmp: userComponent,
		mod:     Overwriting,
		dirty:   true,
	}
	c.FF = &Features{c: c}
	c.Scroll = &Scroller{c: c}
	c.LL = newComponentLines(c)
	c.Register = &Listeners{c: c}
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

func (c *Component) backend() api.UIer {
	return c.bcknd
}

// enable component for client usage.
func (c *Component) enable() {
	if c.component != nil {
		return
	}
	c.component = c.layoutCmp.wrapped()
}

// disable component for client usage.
func (c *Component) disable() {
	if c.component == nil {
		return
	}
	c.component = nil
}

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

func (c *Component) embedded() *Component { return c }

func (c *Component) Gaps(level int) *gapsWriter {
	if c.gg == nil {
		c.gg = newGaps(c.fmt.sty)
	}
	return newGapsWriter(level, c.gg)
}

// func (c *Component) component() *Component

// component is the actual implementation of a lines-Component.
type component struct {
	userCmp     Componenter
	globals     *globals
	dim         *lyt.Dim
	mod         ComponentMode
	initialized bool
	ll          *lines
	fmt         llFmt
	lst         *listeners
	ff          *features
	gg          *gaps
	dirty       bool

	// first holds the index of the first displayed line
	first int

	// slctd hold the index of the currently selected line
	slctd int
}

// component gets the component out of a layoutComponenter without using
// a type-switch.
func (c *component) wrapped() *component { return c }

func (c *component) userComponent() Componenter {
	return c.userCmp
}

func (c *component) setInitialized() {
	c.initialized = true
}

func (c *component) ensureFeatures() {
	if c.ff != nil {
		return
	}
	c.ff = defaultFeatures.copy()
}

func (c *component) ensureListeners() {
	if c.lst != nil {
		return
	}
	c.lst = &listeners{}
}

// globals represents settings which apply for all lines of a component.
type globals struct {
	tabWidth int
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
			cmp.initialize(cmp, sw.userCmp.backend())
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
			cmp.initialize(cmp, cw.userCmp.backend())
		}
		return cb(cmp.layoutComponent())
	})
}
