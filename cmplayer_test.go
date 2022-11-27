// Copyright (c) 2022 Stephan Lukits. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package lines

import (
	"fmt"
	"strings"
	"testing"
	"time"

	. "github.com/slukits/gounit"
)

type Layer struct{ Suite }

func (s *Layer) SetUp(t *T) { t.Parallel() }

func fx(t *T, cmp Componenter, timeout ...time.Duration) *Fixture {
	d := time.Duration(0)
	if len(timeout) > 0 {
		d = timeout[0]
	}
	if cmp == nil {
		cmp = &cmpFX{}
	}
	return TermFixture(t.GoT(), d, cmp)
}

func (s *Layer) Overwrites_base_layer(t *T) {
	base := "1st\n2nd\n3rd"
	fx := fx(t, &icmpFX{
		init: func(_ *icmpFX, e *Env) { fmt.Fprint(e, base) },
	})
	fx.FireResize(3, 3)
	t.Eq(base, fx.Screen())

	fx.Lines.Update(fx.Root(), nil, func(e *Env) {
		fx.Root().(*icmpFX).Layered(
			e,
			&icmpFX{init: func(_ *icmpFX, e *Env) { fmt.Fprint(e, 0) }},
			&LayerPos{},
		)
	})

	t.Eq("1st\n20d\n3rd", fx.Screen())
}

type layeredFX struct {
	Component
	layer Componenter
	init  func(*layeredFX, *Env)
	def   *LayerPos
	onAdd func(*layeredFX, *Env)
}

func (c *layeredFX) OnInit(e *Env) {
	if c.init != nil {
		c.init(c, e)
		return
	}
	fmt.Fprint(e, "1st\n2nd\n3rd")
}

func (c *layeredFX) addLayer(fx *Fixture) {
	fx.Lines.Update(c, nil, func(e *Env) {
		if c.layer == nil {
			c.layer = &icmpFX{init: func(_ *icmpFX, e *Env) {
				fmt.Fprint(e, 0)
			}}
		}
		c.Layered(e, c.layer, c.def)
		if c.onAdd != nil {
			c.onAdd(c, e)
		}
	})
}

func (c *layeredFX) rmLayer(fx *Fixture) {
	fx.Lines.Update(c, nil, func(e *Env) {
		fx.Root().(*layeredFX).RemoveLayer(e)
	})
}

func (s *Layer) Is_removed_from_screen_if_removed_from_layout(t *T) {
	fx := fx(t, &layeredFX{})
	fx.FireResize(3, 3)
	fx.Root().(*layeredFX).addLayer(fx)
	t.Eq("1st\n20d\n3rd", fx.Screen())

	fx.Lines.Update(fx.Root(), nil, func(e *Env) {
		fx.Root().(*layeredFX).RemoveLayer(e)
	})
	t.Eq("1st\n2nd\n3rd", fx.Screen())
}

type layerFX struct {
	Component
	init  func(*layerFX, *Env)
	lyt   func(*layerFX, *Env)
	mouse func(*layerFX, *Env)
	key   func(*layerFX, *Env)
}

func (c *layerFX) OnInit(e *Env) {
	if c.init == nil {
		return
	}
	c.init(c, e)
}

func (c *layerFX) OnLayout(e *Env) bool {
	if c.lyt == nil {
		return false
	}
	c.lyt(c, e)
	return false
}

func (c *layerFX) OnMouse(e *Env, b ButtonMask, x, y int) {
	if c.mouse == nil {
		return
	}
	c.mouse(c, e)
}

func (c *layerFX) OnKey(e *Env, k Key, mm ModifierMask) {
	if c.key == nil {
		return
	}
	c.key(c, e)
}

func (s *Layer) Gets_on_init_reported(t *T) {
	olOnInitReported := false
	fx := fx(t, &layeredFX{
		layer: &layerFX{init: func(_ *layerFX, _ *Env) {
			olOnInitReported = true
		}},
	})
	fx.Root().(*layeredFX).addLayer(fx)

	t.True(olOnInitReported)
}

func (s *Layer) Gets_on_layout_reported(t *T) {
	layerOnLayoutReported := false
	fx := fx(t, &layeredFX{
		layer: &layerFX{
			lyt: func(_ *layerFX, _ *Env) {
				layerOnLayoutReported = true
			}}},
	)
	fx.Root().(*layeredFX).addLayer(fx)

	t.True(layerOnLayoutReported)
}

func (s *Layer) Is_synced_to_the_screen(t *T) {
	fx := fx(t, &layeredFX{
		init: func(lf *layeredFX, e *Env) {},
		layer: &layerFX{init: func(_ *layerFX, e *Env) {
			fmt.Fprint(e, "overlay")
		}},
	})
	fx.Root().(*layeredFX).addLayer(fx)

	t.Eq("overlay", fx.Screen().Trimmed())
}

func (s *Layer) Position_defaults_to_screen_centered(t *T) {
	fx := fx(t, &layeredFX{
		layer: &layerFX{init: func(lf *layerFX, e *Env) {
			fmt.Fprint(e, "X")
		}},
	})
	fx.FireResize(3, 3)
	fx.Root().(*layeredFX).addLayer(fx)

	t.Eq("1st\n2Xd\n3rd", fx.Screen())
}

func (s *Layer) Cant_be_overwritten_by_layered_component(t *T) {
	fx := fx(t, &layeredFX{
		init: func(of *layeredFX, e *Env) { fmt.Fprint(e, "1st\n2nd\n3rd") },
		layer: &layerFX{init: func(lf *layerFX, e *Env) {
			fmt.Fprint(e, "X")
		}},
	})
	fx.FireResize(3, 3)
	fx.Root().(*layeredFX).addLayer(fx)
	t.Eq("1st\n2Xd\n3rd", fx.Screen())

	fx.Lines.Update(fx.Root(), nil, func(e *Env) {
		fmt.Fprintf(e.LL(1), "2nd")
	})
	t.Eq("1st\n2Xd\n3rd", fx.Screen())
}

const layered = `
111111222222
1    12    2
1  llllll  2
111l    l222
333l    l444
3  llllll  4
3    34    4
333333444444
`

const notLayered = `
111111222222
1    12    2
1    12    2
111111222222
333333444444
3    34    4
3    34    4
333333444444
`

type stcFX struct {
	Component
	Stacking
}

func (c *stcFX) OnInit(e *Env) {
	r1, r2 := &chnFX{}, &chnFX{}
	r1.CC = []Componenter{
		&layeredFX{
			init: func(c *layeredFX, e *Env) {
				Print(c.Gaps(0).Filling(), '1')
				fmt.Fprint(c.Gaps(0).Corners, "1")
			},
			layer: &frmFX{filler: 'l'},
			def:   NewLayerPos(3, 2, 6, 4),
		},
		&frmFX{filler: '2'},
	}
	r2.CC = []Componenter{
		&frmFX{filler: '3'}, &frmFX{filler: '4'}}
	c.CC = []Componenter{r1, r2}
}

func (c *stcFX) removeLayer(fx *Fixture) {
	ld := c.CC[0].(*chnFX).CC[0].(*layeredFX)
	fx.Lines.Update(ld, nil, func(e *Env) {
		ld.RemoveLayer(e)
	})
}

func (c *stcFX) addLayer(fx *Fixture) {
	ld := c.CC[0].(*chnFX).CC[0].(*layeredFX)
	fx.Lines.Update(ld, nil, func(e *Env) {
		ld.Layered(e, ld.layer, ld.def)
	})
}

type chnFX struct {
	Component
	Chaining
}

type frmFX struct {
	Component
	filler rune
}

func (c *frmFX) OnInit(e *Env) {
	Print(c.Gaps(0).Filling(), c.filler)
	fmt.Fprint(c.Gaps(0).Corners, string(c.filler))
}

func (s *Layer) Removal_triggers_reprint_of_layered(t *T) {
	fx := fx(t, &stcFX{})
	fx.FireResize(12, 8)
	fx.Root().(*stcFX).addLayer(fx)
	t.Eq(strings.TrimSpace(layered), fx.Screen())

	fx.Root().(*stcFX).removeLayer(fx)
	t.Eq(strings.TrimSpace(notLayered), fx.Screen())
}

func (s *Layer) Gets_mouse_click_reported(t *T) {
	layerOnMouseReported := false
	fx := fx(t, &layeredFX{
		layer: &layerFX{
			mouse: func(_ *layerFX, _ *Env) {
				layerOnMouseReported = true
			}}},
	)
	fx.FireResize(3, 3)
	fx.Root().(*layeredFX).addLayer(fx)
	fx.FireMouse(1, 1, Primary, ZeroModifier)

	t.True(layerOnMouseReported)
}

func (s *Layer) Gets_key_reported(t *T) {
	layerOnKeyReported := false
	fx := fx(t, &layeredFX{
		layer: &layerFX{
			key: func(_ *layerFX, _ *Env) {
				layerOnKeyReported = true
			}},
		onAdd: func(lf *layeredFX, e *Env) {
			e.Lines.Focus(lf.layer)
		}})
	fx.FireResize(3, 3)
	fx.Root().(*layeredFX).addLayer(fx)
	fx.FireKey(Esc)

	t.True(layerOnKeyReported)
}

type stkLayerFX struct {
	layerFX
	Stacking
}

func (s *Layer) Gets_mouse_click_reported_bubbling(t *T) {
	lyr, innerReported, outerReported := &stkLayerFX{}, false, false
	lyr.init = func(lf *layerFX, e *Env) {
		lyr.CC = append(lyr.CC, &layerFX{
			mouse: func(lf *layerFX, e *Env) { innerReported = true },
		})
	}
	lyr.mouse = func(lf *layerFX, e *Env) { outerReported = true }
	fx := fx(t, &layeredFX{layer: lyr})
	fx.FireResize(3, 3)
	fx.Root().(*layeredFX).addLayer(fx)
	fx.FireMouse(1, 1, Primary, ZeroModifier)
	t.True(innerReported)
	t.True(outerReported)
}

func (s *Layer) Gets_key_reported_bubbling(t *T) {
	lyr, innerReported, outerReported := &stkLayerFX{}, false, false
	lyr.init = func(lf *layerFX, e *Env) {
		lyr.CC = append(lyr.CC, &layerFX{
			key: func(lf *layerFX, e *Env) { innerReported = true },
		})
		e.Lines.Focus(lyr.CC[0])
	}
	lyr.key = func(lf *layerFX, e *Env) { outerReported = true }
	fx := fx(t, &layeredFX{layer: lyr})
	fx.FireResize(3, 3)
	fx.Root().(*layeredFX).addLayer(fx)
	fx.FireKey(Esc)
	t.True(innerReported)
	t.True(outerReported)
}

type mdlLayerFX struct {
	layerFX
	onOutOfBound func(*mdlLayerFX, *Env)
}

func (l *mdlLayerFX) OnOutOfBoundClick(e *Env) bool {
	if l.onOutOfBound != nil {
		l.onOutOfBound(l, e)
	}
	return false
}

func (l *mdlLayerFX) OnOutOfBoundMove(e *Env) bool {
	if l.onOutOfBound != nil {
		l.onOutOfBound(l, e)
	}
	return false
}

type chnFcsLstFX struct {
	chnFX
	focusLost bool
}

func (c *chnFcsLstFX) OnFocusLost(_ *Env) { c.focusLost = true }

func (s *Layer) Gets_focus_if_modal(t *T) {
	cmp := &stackedCmpFX{Stacking: Stacking{
		CC: []Componenter{&chnFX{Chaining: Chaining{
			CC: []Componenter{&stackedCmpFX{Stacking: Stacking{
				CC: []Componenter{&chnFcsLstFX{chnFX: chnFX{Chaining: Chaining{
					CC: []Componenter{&layeredFX{layer: &mdlLayerFX{}}},
				}}}}}}}}}}}}
	fx := fx(t, cmp)
	parentOfLayered := cmp.CC[0].(*chnFX).CC[0].(*stackedCmpFX).CC[0].(*chnFcsLstFX)
	fx.Lines.Update(fx.Root(), nil, func(e *Env) {
		e.Lines.Focus(parentOfLayered)
	})
	fx.Lines.Update(fx.Root(), nil, func(e *Env) {
		t.Eq(e.Focused(), parentOfLayered)
	})

	parentOfLayered.CC[0].(*layeredFX).addLayer(fx)
	fx.Lines.Update(fx.Root(), nil, func(e *Env) {
		t.Eq(e.Focused(), parentOfLayered.CC[0].(*layeredFX).layer)
	})
	// t.True(parentOfLayered.focusLost)
}

func (s *Layer) Looses_focus_to_layered_if_removed_and_modal(t *T) {
	fx := fx(t, &layeredFX{layer: &mdlLayerFX{}})
	fx.Root().(*layeredFX).addLayer(fx)
	fx.Lines.Update(fx.Root(), nil, func(e *Env) {
		t.Eq(e.Focused(), fx.Root().(*layeredFX).layer)
	})

	fx.Root().(*layeredFX).rmLayer(fx)
	fx.Lines.Update(fx.Root(), nil, func(e *Env) {
		t.Eq(e.Focused(), fx.Root())
	})
}

func (s *Layer) Gets_out_of_bounds_click_reported_if_modal(t *T) {
	cmp, outOfBoundReported := &mdlLayerFX{}, false
	cmp.init = func(lf *layerFX, e *Env) { fmt.Fprint(e, "0") }
	cmp.onOutOfBound = func(mlf *mdlLayerFX, e *Env) {
		outOfBoundReported = true
	}
	fx := fx(t, &layeredFX{layer: cmp})
	fx.FireResize(3, 3)
	fx.Root().(*layeredFX).addLayer(fx)

	fx.FireClick(0, 0)
	t.True(outOfBoundReported)
}

func (s *Layer) Gets_out_of_bounds_move_reported_if_modal(t *T) {
	cmp, outOfBoundReported := &mdlLayerFX{}, false
	cmp.init = func(lf *layerFX, e *Env) { fmt.Fprint(e, "0") }
	cmp.onOutOfBound = func(mlf *mdlLayerFX, e *Env) {
		outOfBoundReported = true
	}
	fx := fx(t, &layeredFX{layer: cmp})
	fx.FireResize(3, 3)
	fx.Root().(*layeredFX).addLayer(fx)

	// omitting 1, 1 arguments would result in a move from (0,0) to
	// (0,0) which would not be reported.
	fx.FireMove(0, 0, 1, 1)
	t.True(outOfBoundReported)
}

func (s *Layer) Is_moved_to_provided_coordinates(t *T) {
	cmp := &layeredFX{def: NewLayerPos(1, 1, 1, 1)}
	fx := fx(t, cmp)
	fx.FireResize(3, 3)
	t.Eq("1st\n2nd\n3rd", fx.Screen())
	fx.Root().(*layeredFX).addLayer(fx)
	t.Eq("1st\n20d\n3rd", fx.Screen())
	fx.Lines.Update(cmp, nil, func(e *Env) {
		cmp.def.MoveTo(0, 0)
	})
	t.Eq("0st\n2nd\n3rd", fx.Screen())
}

func TestLayer(t *testing.T) {
	t.Parallel()
	Run(&Layer{}, t)
}
