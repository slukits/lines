// Copyright (c) 2022 Stephan Lukits. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package lines

import (
	"fmt"
	"strings"
	"testing"

	. "github.com/slukits/gounit"
)

type Layer struct{ Suite }

func (s *Layer) SetUp(t *T) { t.Parallel() }

func (s *Layer) Overwrites_base_layer(t *T) {
	base := "1st\n2nd\n3rd"
	fx := fx(t, &cmpFX{
		onInit: func(_ *cmpFX, e *Env) { fmt.Fprint(e, base) },
	})
	fx.FireResize(3, 3)
	t.Eq(base, fx.Screen())

	fx.Lines.Update(fx.Root(), nil, func(e *Env) {
		fx.Root().(*cmpFX).Layered(
			e,
			&cmpFX{onInit: func(_ *cmpFX, e *Env) { fmt.Fprint(e, 0) }},
			&LayerPos{},
		)
	})

	t.Eq("1st\n20d\n3rd", fx.Screen())
}

func (s *Layer) Is_removed_from_screen_if_removed_from_layout(t *T) {
	cmp := &cmpFX{onInit: func(_ *cmpFX, e *Env) {
		fmt.Fprint(e, "1st\n2nd\n3rd")
	}}
	fx := fx(t, cmp)
	fx.FireResize(3, 3)
	t.Eq("1st\n2nd\n3rd", fx.Screen())

	lyr := &cmpFX{onInit: func(_ *cmpFX, e *Env) { fmt.Fprint(e, 0) }}
	fx.Lines.Update(cmp, nil, func(e *Env) {
		cmp.Layered(e, lyr, nil)
	})
	t.Eq("1st\n20d\n3rd", fx.Screen())

	fx.Lines.Update(cmp, nil, func(e *Env) {
		cmp.RemoveLayer(e)
	})
	t.Eq("1st\n2nd\n3rd", fx.Screen())
}

type layerFXZZZ struct {
	Component
	init  func(*layerFXZZZ, *Env)
	lyt   func(*layerFXZZZ, *Env)
	mouse func(*layerFXZZZ, *Env)
	key   func(*layerFXZZZ, *Env)
}

func (c *layerFXZZZ) OnInit(e *Env) {
	if c.init == nil {
		return
	}
	c.init(c, e)
}

func (c *layerFXZZZ) OnLayout(e *Env) bool {
	if c.lyt == nil {
		return false
	}
	c.lyt(c, e)
	return false
}

func (c *layerFXZZZ) OnMouse(e *Env, b ButtonMask, x, y int) {
	if c.mouse == nil {
		return
	}
	c.mouse(c, e)
}

func (c *layerFXZZZ) OnKey(e *Env, k Key, mm ModifierMask) {
	if c.key == nil {
		return
	}
	c.key(c, e)
}

func (s *Layer) Gets_on_init_reported(t *T) {
	cmp, lyr := &cmpFX{}, &cmpFX{}
	fx := fx(t, cmp)
	fx.Lines.Update(cmp, nil, func(e *Env) {
		cmp.Layered(e, lyr, nil)
	})

	t.Eq(1, lyr.N(onInit))
}

func (s *Layer) Gets_on_layout_reported(t *T) {
	cmp, lyr := &cmpFX{}, &cmpFX{}
	fx := fx(t, cmp)
	fx.Lines.Update(cmp, nil, func(e *Env) {
		cmp.Layered(e, lyr, nil)
	})

	t.Eq(1, lyr.N(onLayout))
}

func (s *Layer) Is_synced_to_the_screen(t *T) {
	cmp, lyr := &cmpFX{}, &cmpFX{
		onInit: func(cf *cmpFX, e *Env) { fmt.Fprint(e, "overlay") },
	}
	fx := fx(t, cmp)
	fx.Lines.Update(cmp, nil, func(e *Env) {
		cmp.Layered(e, lyr, nil)
	})

	t.Eq("overlay", fx.Screen().Trimmed())
}

func (s *Layer) Position_defaults_to_screen_centered(t *T) {
	cmp := &cmpFX{onInit: func(cf *cmpFX, e *Env) {
		fmt.Fprint(e, "1st\n2nd\n3rd")
	}}
	lyr := &cmpFX{
		onInit: func(cf *cmpFX, e *Env) { fmt.Fprint(e, "X") },
	}
	fx := fx(t, cmp)
	fx.FireResize(3, 3)
	fx.Lines.Update(cmp, nil, func(e *Env) {
		cmp.Layered(e, lyr, nil)
	})

	t.Eq("1st\n2Xd\n3rd", fx.Screen())
}

func (s *Layer) Cant_be_overwritten_by_layered_component(t *T) {
	cmp := &cmpFX{onInit: func(cf *cmpFX, e *Env) {
		fmt.Fprint(e, "1st\n2nd\n3rd")
	}}
	lyr := &cmpFX{
		onInit: func(cf *cmpFX, e *Env) { fmt.Fprint(e, "X") },
	}
	fx := fx(t, cmp)
	fx.FireResize(3, 3)
	fx.Lines.Update(cmp, nil, func(e *Env) {
		cmp.Layered(e, lyr, nil)
	})
	t.Eq("1st\n2Xd\n3rd", fx.Screen())

	fx.Lines.Update(cmp, nil, func(e *Env) {
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

func (s *Layer) Removal_triggers_reprint_of_layered(t *T) {
	stk, chn1, chn2 := &stackingFX{}, &chainingFX{}, &chainingFX{}
	lyrd := &cmpFX{onInit: func(c *cmpFX, e *Env) {
		Print(c.Gaps(0).Filling(), '1')
		fmt.Fprint(c.Gaps(0).Corners, "1")
	}}
	chn1.CC = append(chn1.CC, lyrd, &framingFX{filler: '2'})
	chn2.CC = append(chn2.CC, &framingFX{filler: '3'},
		&framingFX{filler: '4'})
	stk.CC = append(stk.CC, chn1, chn2)
	fx := fx(t, stk)
	fx.FireResize(12, 8)
	t.Eq(strings.TrimSpace(notLayered), fx.Screen())

	lyr, lyrPos := &framingFX{filler: 'l'}, NewLayerPos(3, 2, 6, 4)
	fx.Lines.Update(lyrd, nil, func(e *Env) {
		lyrd.Layered(e, lyr, lyrPos)
	})
	t.Eq(strings.TrimSpace(layered), fx.Screen())

	fx.Lines.Update(lyrd, nil, func(e *Env) {
		lyrd.RemoveLayer(e)
	})
	t.Eq(strings.TrimSpace(notLayered), fx.Screen())
}

func (s *Layer) Gets_mouse_click_reported(t *T) {
	cmp, lyr := &cmpFX{}, &cmpFX{}
	fx := fx(t, cmp)
	fx.FireResize(3, 3)
	fx.Lines.Update(cmp, nil, func(e *Env) {
		cmp.Layered(e, lyr, nil)
	})

	fx.FireMouse(1, 1, Primary, ZeroModifier)
	t.Eq(1, lyr.N(onMouseN))
}

func (s *Layer) Gets_key_reported(t *T) {
	cmp, lyr := &cmpFX{}, &cmpFX{}
	fx := fx(t, cmp)
	fx.FireResize(3, 3)
	fx.Lines.Update(cmp, nil, func(e *Env) {
		cmp.Layered(e, lyr, nil)
		e.Lines.Focus(lyr)
	})

	fx.FireKey(Esc)
	t.Eq(1, lyr.N(onKeyN))
}

type stkLayerFX struct {
	layerFXZZZ
	Stacking
}

func (s *Layer) Gets_mouse_click_reported_bubbling(t *T) {
	cmp, lyr, inner := &cmpFX{}, &stackingFX{}, &cmpFX{}
	lyr.CC = append(lyr.CC, inner)
	fx := fx(t, cmp)
	fx.FireResize(3, 3)
	fx.Lines.Update(cmp, nil, func(e *Env) {
		cmp.Layered(e, lyr, nil)
	})
	fx.FireMouse(1, 1, Primary, ZeroModifier)
	t.Eq(1, inner.N(onMouseN))
	t.Eq(1, lyr.N(onMouseN))
}

func (s *Layer) Gets_key_reported_bubbling(t *T) {
	cmp, lyr, inner := &cmpFX{}, &stackingFX{}, &cmpFX{}
	lyr.CC = append(lyr.CC, inner)
	fx := fx(t, cmp)
	fx.FireResize(3, 3)
	fx.Lines.Update(cmp, nil, func(e *Env) {
		cmp.Layered(e, lyr, nil)
		e.Lines.Focus(inner)
	})
	fx.FireKey(Esc)
	t.Eq(1, inner.N(onKeyN))
	t.Eq(1, lyr.N(onKeyN))
}

func (s *Layer) Gets_focus_if_modal(t *T) {
	cmp, chn1, stk, chn2, lyrd := &stackingFX{}, &chainingFX{},
		&stackingFX{}, &chainingFX{}, &cmpFX{}
	cmp.CC = append(cmp.CC, chn1)
	chn1.CC = append(chn1.CC, stk)
	stk.CC = append(stk.CC, chn2)
	chn2.CC = append(chn2.CC, lyrd)
	fx := fx(t, cmp)
	fx.Lines.Update(cmp, nil, func(e *Env) {
		e.Lines.Focus(chn2)
	})
	fx.Lines.Update(cmp, nil, func(e *Env) {
		t.Eq(e.Focused(), chn2)
	})

	mdlLyr := &modalLayerFX{}
	fx.Lines.Update(lyrd, nil, func(e *Env) {
		lyrd.Layered(e, mdlLyr, nil)
	})
	fx.Lines.Update(cmp, nil, func(e *Env) {
		t.Eq(e.Focused(), mdlLyr)
	})
}

func (s *Layer) Looses_focus_to_layered_if_removed_and_modal(t *T) {
	cmp, mdlLyr := &cmpFX{}, &modalLayerFX{}
	fx := fx(t, cmp)
	fx.Lines.Update(cmp, nil, func(e *Env) {
		cmp.Layered(e, mdlLyr, nil)
	})
	fx.Lines.Update(cmp, nil, func(e *Env) {
		t.Eq(e.Focused(), mdlLyr)
	})

	fx.Lines.Update(cmp, nil, func(e *Env) {
		cmp.RemoveLayer(e)
	})
	fx.Lines.Update(cmp, nil, func(e *Env) {
		t.Eq(e.Focused(), cmp)
	})
}

func (s *Layer) Gets_out_of_bounds_click_reported_if_modal(t *T) {
	cmp, mdlLyr := &cmpFX{}, &modalLayerFX{}
	mdlLyr.onInit = func(_ *cmpFX, e *Env) { fmt.Fprint(e, "0") }
	fx := fx(t, cmp)
	fx.FireResize(3, 3)
	fx.Lines.Update(cmp, nil, func(e *Env) {
		cmp.Layered(e, mdlLyr, nil)
	})

	fx.FireClick(0, 0)
	t.Eq(1, mdlLyr.N(onOutOfBoundClick))
}

func (s *Layer) Gets_out_of_bounds_move_reported_if_modal(t *T) {
	cmp, mdlLyr := &cmpFX{}, &modalLayerFX{}
	mdlLyr.onInit = func(_ *cmpFX, e *Env) { fmt.Fprint(e, "0") }
	fx := fx(t, cmp)
	fx.FireResize(3, 3)
	fx.Lines.Update(cmp, nil, func(e *Env) {
		cmp.Layered(e, mdlLyr, nil)
	})

	// omitting 1, 1 arguments would result in a move from (0,0) to
	// (0,0) which would not be reported.
	fx.FireMove(0, 0, 1, 1)
	t.Eq(1, mdlLyr.N(onOutOfBoundMove))
}

func (s *Layer) Is_moved_to_provided_coordinates(t *T) {
	lyrPos := NewLayerPos(1, 1, 1, 1)
	cmp := &cmpFX{onInit: func(cf *cmpFX, e *Env) {
		fmt.Fprint(e, "1st\n2nd\n3rd")
	}}
	fx := fx(t, cmp)
	fx.FireResize(3, 3)
	t.Eq("1st\n2nd\n3rd", fx.Screen())
	fx.Lines.Update(cmp, nil, func(e *Env) {
		cmp.Layered(e, &cmpFX{onInit: func(cf *cmpFX, e *Env) {
			fmt.Fprint(e, "0")
		}}, lyrPos)
	})
	t.Eq("1st\n20d\n3rd", fx.Screen())
	fx.Lines.Update(cmp, nil, func(e *Env) {
		lyrPos.MoveTo(0, 0)
	})
	t.Eq("0st\n2nd\n3rd", fx.Screen())
}

func TestLayer(t *testing.T) {
	t.Parallel()
	Run(&Layer{}, t)
}
