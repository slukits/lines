// Copyright (c) 2022 Stephan Lukits. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package lines

import (
	"testing"

	. "github.com/slukits/gounit"
)

type mouseEventType uint

const (
	onMove mouseEventType = 1 << iota
	onClick
	onContext
	onDrag
	onDrop
	onMouse
	onEnter
	onExit
)

type mouseEvent struct {
	evtType mouseEventType
	x, y    int
}

type mouseEvents []mouseEvent

func (me *mouseEvents) append(
	evtTyp mouseEventType, x, y int,
) {
	*me = append(*me, mouseEvent{
		evtType: evtTyp, x: x, y: y})
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
	reported     mouseEvents
	stopBubbling mouseEventType
}

func (c *mouseFX) OnMove(e *Env, x, y int) {
	c.reported.append(onMove, x, y)
	if c.stopBubbling&onMove != 0 {
		e.StopBubbling()
	}
}

func (c *mouseFX) OnClick(e *Env, x, y int) {
	c.reported.append(onClick, x, y)
	if c.stopBubbling&onClick != 0 {
		e.StopBubbling()
	}
}

func (c *mouseFX) OnContext(e *Env, x, y int) {
	c.reported.append(onContext, x, y)
	if c.stopBubbling&onContext != 0 {
		e.StopBubbling()
	}
}

func (c *mouseFX) OnDrag(e *Env, b ButtonMask, x, y int) {
	c.reported.append(onDrag, x, y)
	if c.stopBubbling&onDrag != 0 {
		e.StopBubbling()
	}
}

func (c *mouseFX) OnDrop(e *Env, b ButtonMask, x, y int) {
	c.reported.append(onDrop, x, y)
	if c.stopBubbling&onDrop != 0 {
		e.StopBubbling()
	}
}

func (c *mouseFX) OnMouse(e *Env, bm ButtonMask, x, y int) {
	c.reported.append(onMouse, x, y)
	if c.stopBubbling&onMouse != 0 {
		e.StopBubbling()
	}
}

func (c *mouseFX) OnEnter(e *Env) {
	c.reported.append(onEnter, 0, 0)
}

func (c *mouseFX) OnExit(e *Env) {
	c.reported.append(onExit, 0, 0)
}

func (c *mouseFX) HasMove() bool    { return c.reported.has(onMove) }
func (c *mouseFX) HasClick() bool   { return c.reported.has(onClick) }
func (c *mouseFX) HasContext() bool { return c.reported.has(onContext) }
func (c *mouseFX) HasDrag() bool    { return c.reported.has(onDrag) }
func (c *mouseFX) HasDrop() bool    { return c.reported.has(onDrop) }
func (c *mouseFX) HasMouse() bool   { return c.reported.has(onMouse) }
func (c *mouseFX) HasEnter() bool   { return c.reported.has(onEnter) }
func (c *mouseFX) HasExit() bool    { return c.reported.has(onExit) }

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

// nonZero provides the component with x and y not zero having all the mouse
// listeners implemented.  I.e. there are two mouse events reported if
// Click or Context events are reported since the Mouse-event in
// these cases is also reported.
func (c *nonZeroOriginFx) nonZero() *mouseCmpFX {
	return c.CC[1].(*nonZeroOriginChainFx).CC[1].(*mouseCmpFX)
}

func (c *nonZeroOriginFx) OnInit(*Env) {
	c.CC = append(c.CC, &mouseCmpFX{}, &nonZeroOriginChainFx{})
}

type Mouse struct{ Suite }

func (s *Mouse) SetUp(t *T) { t.Parallel() }

func (s *Mouse) Move_is_reported_to_focused_component(t *T) {
	cmp := &mouseCmpFX{}
	fx := fx(t, cmp)
	fx.FireMove(1, 1)
	t.True(cmp.HasMove())
}

func (s *Mouse) Click_is_reported_to_focused_component(t *T) {
	cmp := &mouseCmpFX{}
	fx := fx(t, cmp)
	fx.FireClick(1, 1)
	t.True(cmp.HasClick())
}

func (s *Mouse) Context_is_reported_to_focused_component(t *T) {
	cmp := &mouseCmpFX{}
	fx := fx(t, cmp)
	fx.FireContext(1, 1)
	t.True(cmp.HasContext())
}

func (s *Mouse) DragNDrop_is_reported_to_focused_component(t *T) {
	cmp := &mouseCmpFX{}
	fx := fx(t, cmp)
	fx.FireDragNDrop(1, 1, Primary, ZeroModifier)
	t.True(cmp.HasDrag())
	t.True(cmp.HasDrop())
}

func (s *Mouse) Event_is_reported_to_focused_component(t *T) {
	cmp := &mouseCmpFX{}
	fx := fx(t, cmp)
	fx.FireMouse(1, 1, Middle, 0)
	t.True(cmp.HasMouse())
}

func (s *Mouse) Click_moves_focus_and_reports_to_focusable(t *T) {
	cmp := &nonZeroOriginFx{}
	fx := fx(t, cmp)
	var fxX, fxY int
	t.FatalOn(fx.Lines.Update(cmp.nonZero(), nil, func(e *Env) {
		cmp.nonZero().FF.Set(Focusable)
		t.Not.True(e.Focused() == cmp.nonZero())
		// need an event callback to access component features
		fxX, fxY = cmp.nonZero().Dim().X(), cmp.nonZero().Dim().Y()
	}))
	fx.FireClick(fxX+1, fxY+1)
	t.FatalOn(fx.Lines.Update(cmp.nonZero(), nil, func(e *Env) {
		t.True(e.Focused() == cmp.nonZero())
	}))
	t.True(cmp.nonZero().HasClick())
}

func (s *Mouse) Context_moves_focus_and_reports_to_focusable(t *T) {
	cmp := &nonZeroOriginFx{}
	fx := fx(t, cmp)
	var fxX, fxY int
	t.FatalOn(fx.Lines.Update(cmp.nonZero(), nil, func(e *Env) {
		cmp.nonZero().FF.Set(Focusable)
		t.Not.True(e.Focused() == cmp.nonZero())
		// need an event callback to access component features
		fxX, fxY = cmp.nonZero().Dim().X(), cmp.nonZero().Dim().Y()
	}))
	fx.FireContext(fxX+1, fxY+1)
	t.FatalOn(fx.Lines.Update(cmp.nonZero(), nil, func(e *Env) {
		t.True(e.Focused() == cmp.nonZero())
	}))
	t.True(cmp.nonZero().HasContext())
}

func (s *Mouse) Is_reported_along_with_other_mouse_listener(t *T) {
	cmp := &mouseCmpFX{}
	fx := fx(t, cmp)
	fx.FireClick(1, 1)   // two mouse events
	fx.FireContext(1, 1) // ditto
	t.FatalIfNot(cmp.HasClick())
	t.FatalIfNot(cmp.HasContext())
	t.FatalIfNot(cmp.HasMouse())
	t.True(cmp.LenMouse() == 4)
}

func (s *Mouse) Click_coordinates_are_translated_into_component(t *T) {
	cmp := &nonZeroOriginFx{}
	fx := fx(t, cmp)
	var fxX, fxY int
	t.FatalOn(fx.Lines.Update(cmp.nonZero(), nil, func(e *Env) {
		cmp.nonZero().FF.Set(Focusable)
		fxX, fxY = cmp.nonZero().Dim().X(), cmp.nonZero().Dim().Y()
	}))
	fx.FireClick(fxX+1, fxY+1)
	fx.FireContext(fxX+1, fxY+1)
	x, y := cmp.nonZero().ClickXY()
	t.True(x == 1 && y == 1)
	x, y = cmp.nonZero().ContextXY()
	t.True(x == 1 && y == 1)
	chainer := cmp.CC[1].(*nonZeroOriginChainFx)
	fxWidth := chainer.CC[0].layoutComponent().wrapped().dim.Width()
	x, y = chainer.ClickXY()
	t.True(x == fxWidth+1 && y == 1)
	x, y = chainer.ContextXY()
	t.True(x == fxWidth+1 && y == 1)
}

func (s *Mouse) Events_are_bubbling(t *T) {
	cmp := &nonZeroOriginFx{}
	fx := fx(t, cmp)
	var fxX, fxY int
	t.FatalOn(fx.Lines.Update(cmp.nonZero(), nil, func(e *Env) {
		cmp.nonZero().FF.Set(Focusable)
		fxX, fxY = cmp.nonZero().Dim().X(), cmp.nonZero().Dim().Y()
	}))
	fx.FireClick(fxX+1, fxY+1)
	fx.FireContext(fxX+1, fxY+1)
	fx.FireMouse(fxX+1, fxY+1, Middle, ZeroModifier)
	chainer := cmp.CC[1].(*nonZeroOriginChainFx)
	t.True(chainer.reported.len(onMouse) == 5)
	t.True(chainer.reported.len(onClick) == 1)
	t.True(chainer.reported.len(onContext) == 1)
}

func (s *Mouse) Event_bubbling_may_be_stopped(t *T) {
	cmp := &nonZeroOriginFx{}
	fx := fx(t, cmp)
	cmp.nonZero().stopBubbling = onClick | onContext | onMouse
	var fxX, fxY int
	t.FatalOn(fx.Lines.Update(cmp.nonZero(), nil, func(e *Env) {
		cmp.nonZero().FF.Set(Focusable)
		fxX, fxY = cmp.nonZero().Dim().X(), cmp.nonZero().Dim().Y()
	}))
	fx.FireClick(fxX+1, fxY+1)
	fx.FireContext(fxX+1, fxY+1)
	fx.FireMouse(fxX+1, fxY+1, Middle, ZeroModifier)
	chainer := cmp.CC[1].(*nonZeroOriginChainFx)
	t.Not.True(chainer.HasClick())
	t.Not.True(chainer.HasContext())
	t.Not.True(chainer.HasMouse())
	t.True(cmp.nonZero().LenMouse() == 5)
}

func (s *Mouse) Enter_is_reported_to_moved_over_component(t *T) {
	cmp := &nonZeroOriginFx{}
	fx, fxX, fxY := fx(t, cmp), 0, 0
	t.FatalOn(fx.Lines.Update(cmp.nonZero(), nil, func(e *Env) {
		fxX, fxY = cmp.nonZero().Dim().X(), cmp.nonZero().Dim().Y()
	}))
	t.Not.True(cmp.nonZero().HasEnter())
	fx.FireMove(fxX+1, fxY+1)
	t.True(cmp.nonZero().HasEnter())
}

func (s *Mouse) Exit_is_reported_to_moved_over_component(t *T) {
	cmp := &nonZeroOriginFx{}
	fx, fxX, fxY := fx(t, cmp), 0, 0
	t.FatalOn(fx.Lines.Update(cmp.nonZero(), nil, func(e *Env) {
		fxX, fxY = cmp.nonZero().Dim().X(), cmp.nonZero().Dim().Y()
	}))
	fx.FireMove(fxX+1, fxY+1)
	t.True(cmp.nonZero().HasEnter())
	t.Not.True(cmp.nonZero().HasExit())
	fx.FireMove(0, 0)
	t.True(cmp.nonZero().HasExit())
}

func TestMouse(t *testing.T) {
	t.Parallel()
	Run(&Mouse{}, t)
}
