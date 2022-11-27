// Copyright (c) 2022 Stephan Lukits. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package lines

import (
	"github.com/slukits/lines/internal/lyt"
)

// LayerPos lets a user define the dimensions and alignment of an
// overlaying component.
type LayerPos = lyt.LayerPos

// NewLayerPos creates new layer positioning instance from given origin,
// width and height and returns a pointer to it which can be passed to
// [Component.Layered].
var NewLayerPos = lyt.NewLayerPos

// Modaler must be implemented by a layer-component which wants to be
// dealt with by the user before the user does anything else.  A layer
// component is a component provided to [Component.Layered].  The "last"
// modal layer gets the focus and stays focused until the layer is
// removed.  Mouse clicks at a position outside a modal layer are
// reported to OnOutsideMouse.
type Modaler interface {

	// OnOutOfBoundClick gets mouse clicks reported whose position is
	// outside of a modal focused layer.  Provided environment e may be
	// used to retrieve the mouse event e.Evt.(MouseEventer).  If
	// OnOutOfBoundClick returns true the default reporting of the mouse
	// event is executed; otherwise it is not.  The former typically may
	// happen if an outside click removes the modal layer.
	OnOutOfBoundClick(e *Env) (continueReporting bool)
}

// OutOfBoundMover is implemented by a modal layers which want to be
// informed about mouse movement outside their layed out boundaries.
type OutOfBoundMover interface {

	// OnOutOfBoundMove implementation of a modal layer gets mouse
	// movement reported happening outside layed out boundaries.
	// Returned bool value indicates if the mouse movement should be
	// continued to be reported normally or not.
	OnOutOfBoundMove(*Env) bool
}

type layeredComponent struct {
	*component
	layer *lyt.Layer
}

func (cl *layeredComponent) Layer() *lyt.Layer { return cl.layer }

func (cl *layeredComponent) Layered() layoutComponenter {
	return cl.userCmp.layoutComponent()
}

// Layered associates given component c  with given layering component l
// while l will be layed out according to given layer positioning pos.
func (c *Component) Layered(e *Env, l Componenter, pos *LayerPos) {
	if !l.hasLayoutWrapper() {
		l.initialize(l, c.backend(), c.globals().clone())
	}
	if pos == nil {
		pos = &LayerPos{}
	}
	// provoke a panic in case c is disabled
	if lc, ok := c.userCmp.layoutComponent().(*layeredComponent); ok {
		lc.layer = &lyt.Layer{Dimer: l.layoutComponent(), Def: pos}
		return
	}
	lc := &layeredComponent{
		component: c.layoutComponent().wrapped(),
		layer:     &lyt.Layer{Dimer: l.layoutComponent(), Def: pos},
	}
	if c.layoutCmp == e.Lines.scr.lyt.Root {
		e.Lines.scr.lyt.Root = lc
	}
	c.layoutCmp = lc
}

// RemoveLayer removes given component c's association with a layering
// component which is removed from the layout.
func (c *Component) RemoveLayer(e *Env) {
	// provoke a panic in case c is disabled
	lc, ok := c.userCmp.layoutComponent().(*layeredComponent)
	if !ok || lc == nil {
		return
	}
	_, ok = lc.layer.Root.(layoutComponenter).userComponent().(Modaler)
	if ok {
		e.Lines.Focus(lc.userComponent())
	}
	if c.layoutCmp == e.Lines.scr.lyt.Root {
		e.Lines.scr.lyt.Root = lc.component
	}
	c.layoutCmp = lc.component
}
