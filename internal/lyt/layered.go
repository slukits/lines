// Copyright (c) 2022 Stephan Lukits. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package lyt

import "sort"

// Layered implementation associates implementing [Dimer] with an
// provided layer positioned absolute on the screen.
type Layered interface {
	// Layer provides a Dimer which overlays according to provided
	// layer definition and its associated Dimer the calculated layout.
	Layer() *Layer
}

// Layereds is a sequence of overlaying components of a given layout
// (manager) in order of their appearance.
type Layereds []Layered

// LayerPos lets a user define the positioning, i.e. the dimensions and
// alignment, of an layer component.  While a layout wrapper of a
// component implementing the Layered interface will adapt this
// positioning to what's possible according to the screen layout.
type LayerPos struct {
	x, y, width, height int
	z                   int
	movedFrom           *Rect
	isDirty             bool
}

func NewLayerPos(x, y, width, height int) *LayerPos {
	return &LayerPos{x: x, y: y, width: width, height: height}
}

func (p *LayerPos) X() int      { return p.x }
func (p *LayerPos) Y() int      { return p.y }
func (p *LayerPos) Width() int  { return p.width }
func (p *LayerPos) Height() int { return p.height }

func (p *LayerPos) hasMoved() bool {
	if p == nil {
		return false
	}
	return p.movedFrom != nil
}

// SetZ sets a layers z-level influencing the order in which dirty
// layers are reported.  The higher a layers z-value the later it is
// reported dirty after an update.
func (p *LayerPos) SetZ(z int) *LayerPos {
	if z < 0 {
		return p
	}
	p.z = z
	if !p.isDirty {
		p.isDirty = true
	}
	return p
}

func (p *LayerPos) isZero() bool {
	if p == nil {
		return true
	}
	return p.x == 0 && p.y == 0 && p.width == 0 && p.height == 0
}

// MoveTo moves a layers origin to given position (x,y).
func (p *LayerPos) MoveTo(x, y int) {
	if x == p.x && y == p.y {
		return
	}
	if p.movedFrom == nil {
		p.movedFrom = NewRect(p.x, p.y, p.width, p.height)
	}
	p.x, p.y = x, y
	if !p.isDirty {
		p.isDirty = true
	}
}

// Layers struct maintains a set of layers associated with a layout.
// Note each layer may provide new layers.  Layers are positioned
// absolute on the base-layer, i.e. the screen layout without any
// layers, in the order they are retrieved from the component structure.
type Layers struct {
	oo   Layereds
	ll   map[Layered]*Layer
	base *Manager
}

func newLayers(m *Manager, oo Layereds) *Layers {
	if len(oo) == 0 || m == nil {
		return nil
	}
	ll := Layers{
		oo:   oo,
		ll:   map[Layered]*Layer{},
		base: m,
	}
	if m.base != nil {
		ll.base = m.base
	}

	for _, o := range ll.oo {
		l := o.Layer()
		l.Manager = &Manager{Root: l.Dimer, base: ll.base}
		ll.ll[o] = l
	}

	return &ll
}

func (ll *Layers) hasDelta(other *Layers) bool {
	if ll == nil && other == nil {
		return false
	}
	if ll == nil && other != nil || ll != nil && other == nil {
		return true
	}
	if len(ll.ll) != len(other.ll) {
		return true
	}
	for k, v := range ll.ll {
		if _, ok := other.ll[k]; !ok {
			return true
		}
		if other.ll[k].Root == v.Root {
			continue
		}
		return true
	}
	return false
}

// Len returns the number of layers.
func (ll *Layers) Len() int {
	if ll == nil {
		return 0
	}
	return len(ll.oo)
}

func (ll *Layers) For(cb func(*Layer) (stop bool)) {
	if ll == nil || cb == nil {
		return
	}
	for _, o := range ll.oo {
		if cb(ll.ll[o]) {
			return
		}
	}
}

func (ll *Layers) sort() *Layers {
	sort.Slice(ll.oo, func(i, j int) bool {
		return ll.ll[ll.oo[i]].Def.z < ll.ll[ll.oo[j]].Def.z
	})
	return ll
}

func (ll *Layers) ForReversed(cb func(*Layer) (stop bool)) {
	if ll == nil || cb == nil {
		return
	}
	for i := len(ll.oo) - 1; i >= 0; i-- {
		cb(ll.ll[ll.oo[i]])
	}
}

// Encloses returns the top-most layer containing given position
// with the coordinates x and y.  Encloses returns nil if there is
// not layer containing (x,y).
func (ll *Layers) Encloses(x, y int) *Layer {
	if ll == nil {
		return nil
	}
	var lyr *Layer
	ll.ForReversed(func(l *Layer) (stop bool) {
		if !l.Root.Dim().Contains(x, y) {
			return
		}
		lyr = l
		return true
	})
	return lyr
}

func (ll *Layers) Containing(d Dimer) *Layer {
	if ll == nil {
		return nil
	}
	var lyr *Layer
	ll.For(func(l *Layer) (stop bool) {
		l.ForDimer(nil, func(_d Dimer) (stop bool) {
			if d != _d {
				return false
			}
			lyr = l
			return true
		})
		return false
	})
	return lyr
}

func (ll *Layers) Have(d Dimer) (found bool) {
	if ll == nil {
		return found
	}
	ll.For(func(l *Layer) (stop bool) {
		if !l.Has(d, nil) {
			return false
		}
		found = true
		return true
	})
	return found
}

func (ll *Layers) isDirty() bool {
	if ll == nil {
		return false
	}
	dirty := false
	ll.For(func(l *Layer) (stop bool) {
		if !l.IsDirty() && !l.Def.isDirty {
			return
		}
		dirty = true
		return true
	})
	return dirty
}

func (ll *Layers) reflow(dirty func(Dimer)) (err error) {
	if ll == nil {
		return nil
	}
	before := map[Dimer]*dim{}
	ll.For(func(l *Layer) (stop bool) {
		l.layoutRoot()
		bfr, _, e := l.reflowLayer()
		if e != nil {
			err = e
			return true
		}
		for k, v := range bfr {
			before[k] = v
		}
		return false
	})
	if err != nil {
		return err
	}
	if dirty == nil {
		ll.For(func(l *Layer) (stop bool) {
			l.ForDimer(nil, func(d Dimer) (stop bool) {
				d.Dim().cleanUp()
				return
			})
			return
		})
		return nil
	}
	dll := []*Layer{}
	ll.For(func(l *Layer) (stop bool) {
		dirty := false
		l.reportDirty(before, func(_ Dimer) {
			if dirty {
				return
			}
			dirty = true
			dll = append(dll, l)
		})
		if l.Def.isDirty {
			l.Def.isDirty = false
			if !dirty {
				dll = append(dll, l)
			}
		}
		return false
	})
	sort.Slice(dll, func(i, j int) bool {
		return dll[i].Def.z < dll[j].Def.z
	})
	ll.sort()
	for _, l := range dll {
		dirty(l.Root)
	}
	return nil
}

func (ll *Layers) copy() *Layers {
	if ll == nil {
		return nil
	}
	cpy := Layers{}
	cpy.oo = append(cpy.oo, ll.oo...)
	cpy.ll = map[Layered]*Layer{}
	for ld, l := range ll.ll {
		cpy.ll[ld] = l
	}
	cpy.base = ll.base
	return &cpy
}

func (ll *Layers) append(other *Layers) {
	if other == nil {
		return
	}
	for _, o := range other.oo {
		ll.oo = append(ll.oo, o)
		ll.ll[o] = other.ll[o]
	}
}

// all adds to given layers ll all further layers which are on top and
// returns them.
func (ll *Layers) all() *Layers {
	if ll == nil {
		return nil
	}
	all := ll.copy()
	ll.For(func(l *Layer) (stop bool) {
		// l.ForDimer returns all layers of a Manager (embedded in Layer)
		all.append(l.ForDimer(nil, func(d Dimer) (stop bool) {
			return
		}).all())
		return false
	})
	return all.sort()
}

func (ll *Layers) BaseDimerLayeredByReMovedLayers(
	base *Manager, other *Layers, cb func(Dimer),
) {
	if other == nil || len(other.oo) == 0 {
		return
	}

	var rm []*Layer // ReMoved layers
	if ll == nil {
		other.For(func(l *Layer) (stop bool) {
			rm = append(rm, l)
			return false
		})
	} else {
		for _, o := range other.oo {
			if ll.ll[o] == other.ll[o] && !other.ll[o].Def.hasMoved() {
				continue
			}
			rm = append(rm, other.ll[o])
		}
	}

	if len(rm) == 0 {
		return
	}

	// max rectangle of ReMoved layers
	x, y, width, height := rectangleOfReMoved(rm[0])
	for _, l := range rm[1:] {
		if l.Dim().x < x {
			x = l.Dim().x
		}
		if l.Dim().y < y {
			y = l.Dim().y
		}
		if l.Dim().width > width {
			width = l.Dim().width
		}
		if l.Dim().height > height {
			height = l.Dim().height
		}
	}

	// find dimer in base layout which were layered by removed layers
	base.ForDimer(nil, func(d Dimer) (stop bool) {
		dx, dy, dw, dh := d.Dim().Screen()
		layered := (x < dx && x+width > dx || x > dx && x < dx+dw) &&
			(y < dy && y+height > dy || y > dy && y < dy+dh)
		if !layered {
			return false
		}
		cb(d)
		return false
	})
}

func rectangleOfReMoved(l *Layer) (x, y, w, h int) {
	if l.Def.hasMoved() {
		from := l.Def.movedFrom
		x, y, w, h = from.x, from.y, from.w, from.h
		l.Def.movedFrom = nil
		return x, y, w, h
	}
	return l.Root.Dim().Screen()
}

type Layer struct {
	*Manager
	Dimer
	Def *LayerPos
}

func (l *Layer) layoutRoot() {
	if l.Def.isZero() {
		width, height := 30, 8
		if l.base.width < width+2 {
			width = l.base.width - 2
		}
		if l.base.height < height+2 {
			height = l.base.height - 2
		}
		l.Root.Dim().setOrigin(
			(l.base.Width-width)/2,
			(l.base.Height-height)/2,
		)
		l.Root.Dim().SetWidth(width).SetHeight(height)
		return
	}
	l.Root.Dim().setOrigin(l.Def.x, l.Def.y)
	l.Root.Dim().SetWidth(l.Def.width).SetHeight(l.Def.height)
}
