// Copyright (c) 2022 Stephan Lukits. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package lines

import (
	"testing"

	"github.com/gdamore/tcell/v2"
	. "github.com/slukits/gounit"
)

type mouseEventType uint

const (
	onClick mouseEventType = iota
	onContext
	onMouse
)

type mouseEvent struct {
	tcell   *tcell.EventMouse
	evtType mouseEventType
	x, y    int
}

type mouseEvents []mouseEvent

func (me *mouseEvents) append(
	evt *tcell.EventMouse, evtTyp mouseEventType, x, y int,
) {
	*me = append(*me, mouseEvent{
		tcell: evt, evtType: evtTyp, x: x, y: y})
}

func (me mouseEvents) has(evtType mouseEventType) bool {
	for _, r := range me {
		if r.evtType != evtType {
			continue
		}
		return true
	}
	return false
}

func (me mouseEvents) len(evtType mouseEventType) int {
	n := 0
	for _, r := range me {
		if r.evtType != evtType {
			continue
		}
		n++
	}
	return n
}

func (me mouseEvents) get(evtType mouseEventType) *mouseEvent {
	for _, r := range me {
		if r.evtType != evtType {
			continue
		}
		return &r
	}
	return nil
}

func (me mouseEvents) getXY(evtType mouseEventType) (x, y int) {
	r := me.get(evtType)
	if r == nil {
		return -1, -1
	}
	return r.x, r.y
}

type mouseFX struct {
	reported mouseEvents
}

func (c *mouseFX) OnClick(e *Env, x, y int) {
	c.reported.append(e.Evt.(*tcell.EventMouse), onClick, x, y)
}

func (c *mouseFX) OnContext(e *Env, x, y int) {
	c.reported.append(e.Evt.(*tcell.EventMouse), onContext, x, y)
}

func (c *mouseFX) OnMouse(e *Env, x, y int) {
	c.reported.append(e.Evt.(*tcell.EventMouse), onMouse, x, y)
}

func (c *mouseFX) HasClick() bool { return c.reported.has(onClick) }

func (c *mouseFX) HasContext() bool {
	return c.reported.has(onContext)
}

func (c *mouseFX) HasMouse() bool { return c.reported.has(onMouse) }

func (c *mouseFX) LenMouse() int { return c.reported.len(onMouse) }

// ClickXY returns the reported x/y-coordinates of the first reported
// click event
func (c *mouseFX) ClickXY() (x, y int) {
	return c.reported.getXY(onClick)
}

// ContextXY returns the reported x/y-coordinates of the first reported
// context event
func (c *mouseFX) ContextXY() (x, y int) {
	return c.reported.getXY(onContext)
}

// MouseXY returns the reported x/y-coordinates of the first reported
// Mouse event
func (c *mouseFX) MouseXY() (x, y int) {
	return c.reported.getXY(onMouse)
}

type mouseCmpFX struct {
	Component
	mouseFX
}

// nonZeroOriginChainFx ensures a component with a non-zero x-value
type nonZeroOriginChainFx struct {
	Component
	Chaining
	mouseFX
}

func (c *nonZeroOriginChainFx) OnInit(*Env) {
	c.CC = append(c.CC, &mouseCmpFX{}, &mouseCmpFX{})
}

// nonZeroOriginFx provides a test layout of components implementing all
// mouse-event listeners with a component whose x, y values are both
// non-zero.  NOTE a Click or Context event is hence reported twice once
// to the OnClick/OnContext listener and once to the OnMouse listener.
type nonZeroOriginFx struct {
	Component
	Stacking
}

// cmp provides the component with x and y not zero having all the mouse
// listeners implemented.  I.e. there are two mouse events reported if
// Click or Context events are reported since the Mouse-event in
// these cases is also reported.
func (c *nonZeroOriginFx) cmp() *mouseCmpFX {
	return c.CC[1].(*nonZeroOriginChainFx).CC[1].(*mouseCmpFX)
}

func (c *nonZeroOriginFx) OnInit(*Env) {
	c.CC = append(c.CC, &mouseCmpFX{}, &nonZeroOriginChainFx{})
}

type Mouse struct{ Suite }

func (s *Mouse) SetUp(t *T) { t.Parallel() }

func (s *Mouse) Click_is_reported_to_focused_component(t *T) {
	fx := &mouseCmpFX{}
	// reports OnClick and OnMouse bubbling hence Max == 4
	ee, tt := Test(t.GoT(), fx, 2)
	tt.FireClick(1, 1)
	t.False(ee.IsListening())
	t.True(fx.HasClick())
}

func (s *Mouse) Context_is_reported_to_focused_component(t *T) {
	fx := &mouseCmpFX{}
	ee, tt := Test(t.GoT(), fx, 2)
	tt.FireContext(1, 1)
	t.False(ee.IsListening())
	t.True(fx.HasContext())
}

func (s *Mouse) Event_is_reported_to_focused_component(t *T) {
	fx := &mouseCmpFX{}
	ee, tt := Test(t.GoT(), fx, 1)
	tt.FireMouse(1, 1, tcell.ButtonMiddle, 0)
	t.False(ee.IsListening())
	t.True(fx.HasMouse())
}

func (s *Mouse) Click_is_reported_to_focusable_component(t *T) {
	fx := &nonZeroOriginFx{}
	// OnInit, 2xOnUpdate, 2xOnClick, 2xOnMouse (because bubbling)
	ee, tt := Test(t.GoT(), fx, 7)
	ee.Listen()
	var fxX, fxY int
	ee.Update(fx.cmp(), nil, func(e *Env) {
		fx.cmp().FF.Add(Focusable)
		t.False(e.Focused() == fx.cmp())
		// need an event callback to access component features
		fxX, fxY = fx.cmp().Dim().X(), fx.cmp().Dim().Y()
	})
	tt.FireClick(fxX+1, fxY+1)
	ee.Update(fx.cmp(), nil, func(e *Env) {
		t.True(e.Focused() == fx.cmp())
	})
	t.False(ee.IsListening())
	t.True(fx.cmp().HasClick())
}

func (s *Mouse) Context_is_reported_to_focusable_component(t *T) {
	fx := &nonZeroOriginFx{}
	// OnInit, 2xOnUpdate, 2xOnClick, 2xOnMouse (because bubbling)
	ee, tt := Test(t.GoT(), fx, 7)
	ee.Listen()
	var fxX, fxY int
	ee.Update(fx.cmp(), nil, func(e *Env) {
		fx.cmp().FF.Add(Focusable)
		t.False(e.Focused() == fx.cmp())
		// need an event callback to access component features
		fxX, fxY = fx.cmp().Dim().X(), fx.cmp().Dim().Y()
	})
	tt.FireContext(fxX+1, fxY+1)
	ee.Update(fx.cmp(), nil, func(e *Env) {
		t.True(e.Focused() == fx.cmp())
	})
	t.False(ee.IsListening())
	t.True(fx.cmp().HasContext())
}

func (s *Mouse) Event_is_reported_to_focusable_component(t *T) {
	fx := &nonZeroOriginFx{}
	// OnInit, 2xOnUpdate, 2xOnMouse (because bubbling)
	ee, tt := Test(t.GoT(), fx, 5)
	ee.Listen()
	var fxX, fxY int
	ee.Update(fx.cmp(), nil, func(e *Env) {
		fx.cmp().FF.Add(Focusable)
		t.False(e.Focused() == fx.cmp())
		// need an event callback to access component features
		fxX, fxY = fx.cmp().Dim().X(), fx.cmp().Dim().Y()
	})
	tt.FireMouse(fxX+1, fxY+1, tcell.ButtonMiddle, tcell.ModNone)
	ee.Update(fx.cmp(), nil, func(e *Env) {
		t.True(e.Focused() == fx.cmp())
	})
	t.False(ee.IsListening())
	t.True(fx.cmp().HasMouse())
}

func (s *Mouse) Is_not_reported_if_not_focusable_and_unfocused(t *T) {
	fx := &nonZeroOriginFx{}
	// OnInit, 2xOnUpdate
	ee, tt := Test(t.GoT(), fx, 3)
	ee.Listen()
	var fxX, fxY int
	ee.Update(fx.cmp(), nil, func(e *Env) {
		t.False(e.Focused() == fx.cmp())
		// need an event callback to access component features
		fxX, fxY = fx.cmp().Dim().X(), fx.cmp().Dim().Y()
	})
	tt.FireClick(fxX+1, fxY+1)
	tt.FireContext(fxX+1, fxY+1)
	tt.FireMouse(fxX+1, fxY+1, tcell.ButtonMiddle, tcell.ModNone)
	ee.Update(fx.cmp(), nil, func(e *Env) {
		t.False(e.Focused() == fx.cmp())
	})
	t.False(ee.IsListening())
	t.False(fx.cmp().HasClick())
	t.False(fx.cmp().HasContext())
	t.False(fx.cmp().HasMouse())
}

func (s *Mouse) Is_reported_along_with_other_mouse_listener(t *T) {
	fx := &mouseCmpFX{}
	// 1xOnClick, 1xOnContext, 2xOnMouse
	ee, tt := Test(t.GoT(), fx, 4)
	ee.Listen()
	tt.FireClick(1, 1)
	tt.FireContext(1, 1)
	t.False(ee.IsListening())
	t.FatalIfNot(fx.HasClick())
	t.FatalIfNot(fx.HasContext())
	t.FatalIfNot(fx.HasMouse())
	t.True(fx.LenMouse() == 2)
}

func (s *Mouse) X_y_is_reported_relative_to_component_s_origin(t *T) {
	fx := &nonZeroOriginFx{}
	// OnInit, OnUpdate, 2xOnClick, 2xOnContext, 4xOnMouse
	ee, tt := Test(t.GoT(), fx, 10)
	ee.Listen()
	var fxX, fxY int
	ee.Update(fx.cmp(), nil, func(e *Env) {
		fx.cmp().FF.Add(Focusable)
		fxX, fxY = fx.cmp().Dim().X(), fx.cmp().Dim().Y()
	})
	tt.FireClick(fxX+1, fxY+1)
	tt.FireContext(fxX+1, fxY+1)
	t.False(ee.IsListening())
	x, y := fx.cmp().ClickXY()
	t.True(x == 1 && y == 1)
	x, y = fx.cmp().ContextXY()
	t.True(x == 1 && y == 1)
	x, y = fx.cmp().MouseXY()
	t.True(x == 1 && y == 1)
	chainer := fx.CC[1].(*nonZeroOriginChainFx)
	fxWidth := chainer.CC[0].layoutComponent().wrapped().dim.Width()
	x, y = chainer.ClickXY()
	t.True(x == fxWidth+1 && y == 1)
	x, y = chainer.ContextXY()
	t.True(x == fxWidth+1 && y == 1)
	x, y = chainer.MouseXY()
	t.True(x == fxWidth+1 && y == 1)
}

func TestMouse(t *testing.T) {
	t.Parallel()
	Run(&Mouse{}, t)
}
