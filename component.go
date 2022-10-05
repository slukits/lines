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

	// Focus provides a component's line highlighting and line
	// selection API.
	Focus *LineFocus

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
		global:  &global{tabWidth: 4},
		fmt:     llFmt{sty: tcell.StyleDefault},
		userCmp: userComponent,
		mod:     Overwriting,
		dirty:   true,
	}
	c.FF = &Features{c: c}
	c.Scroll = &Scroller{c: c}
	c.Focus = &LineFocus{c: c, current: -1}
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

// func (c *Component) component() *Component

// component is the actual implementation of a lines-Component.
type component struct {
	userCmp     Componenter
	global      *global
	dim         *lyt.Dim
	mod         ComponentMode
	initialized bool
	ll          *lines
	fmt         llFmt
	lst         *listeners
	ff          *features
	dirty       bool

	// first holds the index of the first displayed line
	first int

	// slctd hold the index of the currently selected line
	slctd int
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

// Sty sets a component style attributes like bold or dimmed.  See tcell
// AttrMask.
func (c *component) Sty(attr tcell.AttrMask) {
	c.fmt.sty = c.fmt.sty.Attributes(attr)
}

// FG sets a component's foreground color.
func (c *component) FG(color tcell.Color) {
	c.dirty = true
	c.fmt.sty = c.fmt.sty.Foreground(color)
}

// BG sets a component's background color.
func (c *component) BG(color tcell.Color) {
	c.dirty = true
	c.fmt.sty = c.fmt.sty.Background(color)
}

// Len returns the number of lines currently stored in a component.
// Note the line number is independent of a component's associated
// screen area.
func (c *component) Len() int {
	return len(*c.ll)
}

func (c *component) LL(idx int) *line {
	return (*c.ll)[idx]
}

// IsDirty is true if this component is flagged dirty or one of its
// lines.
func (c *component) IsDirty() bool {
	if c.Len() == 0 {
		return c.dirty
	}
	return c.ll.IsDirty() || c.dirty
}

// SetDirty flags a component as dirty having the effect that at the
// next syncing the component's screen area is cleared before it is
// written to.
func (c *component) SetDirty() {
	c.dirty = true
}

// Dim provides a components layout dimensions and features to adapt
// them.
func (c *component) Dim() *lyt.Dim { return c.dim }

const All = -1

// Reset blanks out the content of the line with given index the next
// time it is printed to the screen.  Provide line flags if for example
// a reset line should not be focusable.
func (c *component) Reset(idx int, ff ...LineFlags) {
	if idx < -1 || idx >= c.Len() {
		return
	}

	c.setFirst(0)

	_ff := LineFlags(0)
	for _, f := range ff {
		_ff |= f
	}

	if idx == -1 {
		for _, l := range *c.ll {
			l.reset(c.fmt.sty, _ff)
		}
		return
	}

	(*c.ll)[idx].reset(c.fmt.sty, _ff)
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
	if c.dirty {
		c.dirty = false
	}
	c.sync(rw)
}

// sync writes receiving components lines to the screen.
func (c *component) sync(rw runeWriter) {
	sx, sy, sw, sh := c.dim.Area()
	if c.mod&Tailing == Tailing && c.Len() >= sh {
		c.setFirst(c.Len() - sh)
	}
	if c.dirty {
		c.clear(rw)
		c.dirty = false
	}
	c.ll.For(c.first, func(i int, l *line) (stop bool) {
		if i >= sh {
			return true
		}
		l.sync(sx, sy+i, sw, rw)
		return false
	})
}

// clear fills the receiving component's printable area with spaces.
func (c *component) clear(rw runeWriter) {
	sx, sy, sw, sh := c.dim.Rect()
	for y := sy; y < sy+sh; y++ {
		for x := sx; x < sx+sw; x++ {
			rw.SetContent(x, y, ' ', nil, c.fmt.sty)
		}
	}
}

// setFirst sets the first displayed line and in case it changes given
// component becomes also dirty (hence the indirection).
func (c *component) setFirst(f int) {
	if f < 0 || f == c.first || f >= c.Len() {
		return
	}

	c.first = f
	if c.dirty {
		return
	}
	c.dirty = true
}

func (c *component) write(
	bb []byte, line, cell int, ff LineFlags, sty tcell.Style,
) (int, error) {
	switch {
	case c.mod&(Appending|Tailing) != 0:
		c.ll.append(
			c.lineFactory, ff, sty, bytes.Split(bb, []byte("\n"))...)
	default:
		if line == -1 {
			c.Reset(line)
			line = 0
		}
		c.ll.replaceAt(
			c.lineFactory, line, cell, ff, sty,
			bytes.Split(bb, []byte("\n"))...)
	}
	return len(bb), nil
}

func (c *component) lineFactory() *line {
	return &line{
		sty:    c.fmt.sty,
		dirty:  true,
		global: c.global,
	}
}

// global represents settings which apply for all lines of a component.
type global struct {
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
