// Copyright (c) 2022 Stephan Lukits. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package lyt

import (
	"errors"
	"fmt"
)

// Dimer must be implemented by each component of the layout to provide
// the information for layout calculations which accordingly updates
// this information.
type Dimer interface {

	// Dim provides the layout dimensions of a component which are used
	// and adapted during the layout process.
	Dim() *Dim
}

// Stacker is implemented by components who consist of stacked Dimers
// which are arranged vertically.
type Stacker interface {
	Dimer

	// ForStacked provides the stacked Dimer in the order they are
	// stacked.
	ForStacked(func(Dimer) (stop bool))
}

// Chainer is implemented by components who consist of chained Dimers
// which are arranged horizontally.
type Chainer interface {
	Dimer

	// ForChained provides the chained Dimer in the order they are
	// chained.
	ForChained(func(Dimer) (stop bool))
}

// ErrLyt is the basic error type which is wrapped by all layout errors.
var ErrLyt = errors.New("lty: ")

// A Manager is used to calculates the layout of set Root Dimer i.e.
// its provided Dimers origin, size, margins and clipping.  Is Root not
// set or either its width or height is not positive a Manger's
// operations fail.  Is set Root implementing either the Stacker or
// Chainer interface the layout of provided Dimers by this
// implementation is calculated as well.  If one of these provided
// Dimers implements either of those interfaces its provided Dimers'
// layout is calculated also and so on.  Provided Dimers must not
// implement both interfaces.  In the later case the Stacker supersedes
// the Chainer; no error is reported.  Dimers overflowing their
// available area are clipped, i.e. have either a partial area of their
// wanted area available or are flagged as off-screen (see
// Dim.IsOffScreen).  Dimers which underflow their assigned area receive
// a margin.
type Manager struct {
	Width, Height int
	width, height int
	Root          Dimer
	base          *Manager
	Layers        *Layers
}

func (m *Manager) SetRoot(d Dimer) {
	m.Root, m.width, m.height = d, 0, 0
}

func (m *Manager) validate() error {
	if m.Root == nil {
		return fmt.Errorf("%w%s", ErrLyt, "root must not be nil")
	}
	if m.Height <= 0 && m.Root.Dim().height <= 0 {
		return fmt.Errorf("%w%s", ErrLyt,
			"manager's or root's height must be positive")
	}
	if m.Width <= 0 && m.Root.Dim().width <= 0 {
		return fmt.Errorf("%w%s", ErrLyt,
			"manager's or root's width must be positive")
	}
	if m.Height == 0 {
		m.Height = m.Root.Dim().height
	}
	if m.Width == 0 {
		m.Width = m.Root.Dim().width
	}

	return nil
}

// IsDirty returns true iff one of the components in the layout has been
// flagged as dirty, the set of layers has changed or one of the layers
// is dirty.
func (m *Manager) IsDirty() (dirty bool) {
	if m.height != m.Height || m.width != m.Width {
		return true
	}
	ll := m.ForDimer(nil, func(d Dimer) (stop bool) {
		if !d.Dim().IsDirty() && !d.Dim().IsUpdated() {
			return false
		}
		dirty = true
		return true
	})
	all := ll.all()

	if m.base == nil {
		return dirty || m.Layers.hasDelta(all) || all.isDirty()
	}
	return dirty
}

// HasConsistentLayout returns true iff the following invariants hold
// true.  The layed out width of a Stacker's Dimer must be the layed out
// width of the Stacker and the layed out heights of a Stacker's Dimers
// must sum up the Stacker's layed out height.  The layed out heights of a
// Chainer's Dimers must be the layed out height of the Chainer and all
// layed out widths of a Chainer's Dimers must sum up to the Chainer's
// layed out width.  Whereas the layed out width/height is the
// width/height reduced by its clipping or increased by its relevant
// margins if there is/are any clipping or margins.  NOTE
// HasConsistentLayout returns also false if a Manager is not properly
// initialized.
func (m *Manager) HasConsistentLayout() bool {
	if err := m.validate(); err != nil {
		return false // TODO: coverage
	}
	consistent := true
	m.forContainer(m.Root,
		func(s Stacker) (stop bool) {
			if !isConsistentStacker(s) {
				consistent = false
				return true
			}
			return false
		},
		func(c Chainer) (stop bool) {
			if !isConsistentChainer(c) {
				consistent = false
				return true
			}
			return false
		},
	)
	return consistent
}

func isConsistentStacker(s Stacker) bool {
	lWidth, lHeightSum, ok := s.Dim().layoutWidth(), 0, true
	s.ForStacked(func(d Dimer) (stop bool) {
		if lWidth != d.Dim().layoutWidth() {
			ok = false
			return true
		}
		lHeightSum += d.Dim().layoutHeight()
		return false
	})
	return lHeightSum == s.Dim().layoutHeight() && ok
}

func isConsistentChainer(c Chainer) bool {
	lHeight, lWidthSum, ok := c.Dim().layoutHeight(), 0, true
	c.ForChained(func(d Dimer) (stop bool) {
		if lHeight != d.Dim().layoutHeight() {
			ok = false
			return true
		}
		lWidthSum += d.Dim().layoutWidth()
		return false
	})
	return lWidthSum == c.Dim().layoutWidth() && ok
}

// Reflow calculates the layout of all Dimer provided by Root and
// subsequent stacking, chaining and layered components respecting
// potentially updated widths or heights of layed out Dimer.  Given
// function will receive all dirty layout components i.e. components
// whose area on the screen has changed in origin or size.  Note
// provided layer and their components are always reported as dirty.
func (m *Manager) Reflow(dirty func(Dimer)) (err error) {
	if err := m.validate(); err != nil {
		return err
	}

	ddBefore, ll, err := m.reflowLayer()
	if err != nil {
		return err
	}
	m.Layers = ll.all()

	// if there is no callback we can clean everything up.
	if dirty == nil {
		if err := ll.reflow(nil); err != nil {
			return err
		}
		m.ForDimer(nil, func(d Dimer) (stop bool) {
			d.Dim().cleanUp()
			return false
		})
		return
	}

	m.reportDirty(ddBefore, dirty)

	if err := ll.reflow(dirty); err != nil {
		return err
	}

	return nil
}

func (m *Manager) reflowLayer() (
	before map[Dimer]*dim, ll *Layers, err error,
) {
	// save the printable areas of all Dimers which are layed out to
	// decide after the reflow if it has changed, i.e. if they are
	// considered dirty.
	ddBefore := map[Dimer]*dim{}
	ll = m.ForDimer(nil, func(d Dimer) (stop bool) {
		ddBefore[d] = d.Dim().prepareLayout()
		return false
	})

	// layout all containers.
	m.forContainer(
		m.layedOutRoot(),
		func(s Stacker) (stop bool) {
			if err = layoutStacker(s); err != nil {
				return true
			}
			return false
		},
		func(c Chainer) (stop bool) {
			if err = layoutChainer(c); err != nil {
				return true
			}
			return false
		},
	)
	if err != nil {
		return nil, nil, err
	}
	return ddBefore, ll, nil
}

func (m *Manager) reportDirty(before map[Dimer]*dim, rpr func(Dimer)) {
	m.ForDimer(nil, func(d Dimer) (stop bool) {
		d.Dim().finalizeLayout(before[d])
		if d.Dim().IsDirty() {
			rpr(d)
		}
		d.Dim().cleanUp()
		return false
	})
}

func (m *Manager) layedOutRoot() Dimer {
	if m.width == m.Width && m.height == m.Height {
		return m.Root
	}
	m.width, m.height = m.Width, m.Height
	m.Root.Dim().resetClippingsAndMargins()
	if m.Root.Dim().width == 0 {
		m.Root.Dim().fillsWidth = 1
	}
	if m.Root.Dim().height == 0 {
		m.Root.Dim().fillsHeight = 1
	}
	if m.Root.Dim().fillsWidth > 0 {
		if m.Root.Dim().fillsWidth > m.Width {
			m.Root.Dim().clipWidth = m.Root.Dim().fillsWidth - m.Width
		} else {
			m.Root.Dim().width = m.Width
		}
	} else {
		if m.Root.Dim().width > m.Width {
			m.Root.Dim().clipWidth = m.Root.Dim().width - m.Width
		} else {
			ml := (m.Width - m.Root.Dim().width) / 2
			m.Root.Dim().mrgLeft = ml
			m.Root.Dim().mrgRight = (m.Width - m.Root.Dim().width) - ml
		}
	}
	if m.Root.Dim().fillsHeight > 0 {
		if m.Root.Dim().fillsHeight > m.Height {
			m.Root.Dim().clipHeight = m.Root.Dim().fillsHeight - m.Height
		} else {
			m.Root.Dim().height = m.Height
		}
	} else {
		if m.Root.Dim().height > m.Height {
			m.Root.Dim().clipHeight = m.Root.Dim().height - m.Height
		} else {
			mt := (m.Height - m.Root.Dim().height) / 2
			m.Root.Dim().mrgTop = mt
			m.Root.Dim().mrgBottom = (m.Height - m.Root.Dim().height) - mt
		}
	}
	return m.Root
}

// Locate returns a path of Stacker and Chainer whose last Stacker or
// Chainer provides given Dimer and each Stacker/Chainer in it is
// provided by its previous Stacker/Chainer (or is root).
func (m *Manager) Locate(dr Dimer) (path []Dimer, err error) {
	if err := m.validate(); err != nil {
		return nil, err
	}
	if l := m.Layers.Containing(dr); l != nil {
		return l.Locate(dr)
	}
	d, path := m.Root, []Dimer{}
	dd, forDD := []Dimer{d}, (func(func(d Dimer) bool))(nil)
	for len(dd) > 0 {
		d = dd[len(dd)-1]
		// if a Stacker/Chainer is visited the first time to evaluate
		// its provided dimmers it is left in the queue as indicator
		// that all its descendants have been processed and added to the
		// path.  The next time it is visited it must also be at the end
		// of the path and it is removed.
		if len(path) != 0 && path[len(path)-1] == d {
			dd = dd[:len(dd)-1]
			path = path[:len(path)-1]
			continue
		}
		if d == dr {
			return path, nil // TODO: coverage
		}
		switch d := d.(type) {
		case Stacker:
			path = append(path, d)
			forDD = d.ForStacked
		case Chainer: // TODO: coverage
			path = append(path, d)
			forDD = d.ForChained
		default: // d == Root implementing neither Stacker nor Chainer
			return nil, nil // TODO: coverage
		}
		found := false
		forDD(func(d Dimer) (stop bool) {
			if d == dr {
				found = true
				return true
			}
			if _, ok := d.(Stacker); ok {
				dd = append(dd, d)
				return false // Stacker supersedes Chainer
			}
			if _, ok := d.(Chainer); ok {
				dd = append(dd, d) // TODO: coverage
			}
			return false
		})
		if found {
			return path, nil
		}
	}
	return nil, nil
}

// LocateAt returns a path of Dimers whose last Dimer is
// the smallest Dimer enclosing given coordinates while all other Dimers
// are Stackers or Chainers which also enclose given coordinates as well
// as they enclose each other in a narrowing way.
func (m *Manager) LocateAt(x, y int) (path []Dimer, err error) {
	if err := m.validate(); err != nil {
		return nil, err // TODO: coverage
	}
	rx, ry, rw, rh := m.Root.Dim().Screen()
	if x < rx || y < ry || x >= rw+rx || y >= rh+ry {
		return nil, nil // TODO: coverage
	}
	if l := m.Layers.Encloses(x, y); l != nil {
		return l.LocateAt(x, y)
	}
	last, path := m.Root, []Dimer{}
	forDD := (func(func(d Dimer) bool))(nil)
	for last != nil {
		d := last
		path, last = append(path, d), nil
		switch d := d.(type) {
		case Stacker:
			forDD = d.ForStacked
		case Chainer:
			forDD = d.ForChained
		default:
			forDD = nil
		}
		if forDD == nil {
			break
		}
		forDD(func(d Dimer) (stop bool) {
			dx, dy, dw, dh := d.Dim().Screen()
			if dy <= y && dx <= x && dw+dx > x && dh+dy > y {
				last = d
				return true
			}
			return false
		})
	}
	return path, nil
}

// forContainer iterates in a breadth-first manner over all stacker and
// chainer found in a layout.
func (m *Manager) forContainer(
	d Dimer, s func(Stacker) (stop bool), c func(Chainer) (stop bool),
) *Layers {
	if d == nil {
		return nil // TODO: coverage
	}
	dd, forDD := []Dimer{d}, (func(func(d Dimer) bool))(nil)
	var oo Layereds
	for len(dd) > 0 {
		d, dd = dd[0], dd[1:]
		switch d := d.(type) {
		case Stacker:
			if s(d) {
				return newLayers(m, oo)
			}
			forDD = d.ForStacked
		case Chainer:
			if c(d) {
				return newLayers(m, oo)
			}
			forDD = d.ForChained
		case Layered:
			return newLayers(m, append(oo, d))
		default: // first d is implementing neither Stacker nor Chainer
			return nil // TODO: coverage
		}
		if d, ok := d.(Layered); ok {
			oo = append(oo, d)
		}
		forDD(func(d Dimer) (stop bool) {
			switch d := d.(type) {
			case Stacker:
				dd = append(dd, d)
				return false // Stacker supersedes Chainer
			case Chainer:
				dd = append(dd, d) // TODO: coverage
			case Layered:
				oo = append(oo, d)
			}
			return false
		})
	}
	return newLayers(m, oo)
}

// ForDimer calls back for given Dimer and all nested Dimer and returns
// all dimmer which provide overlays.  Dimer defaults to the Manager's
// Root Dimer if nil.
func (m *Manager) ForDimer(
	d Dimer, cb func(Dimer) (stop bool),
) *Layers {
	if d == nil {
		d = m.Root
	}
	dd, forDD := []Dimer{d}, (func(func(d Dimer) bool))(nil)
	var oo Layereds
	for len(dd) > 0 {
		d, dd, forDD = dd[0], dd[1:], nil
		switch d := d.(type) {
		case Stacker:
			forDD = d.ForStacked
		case Chainer:
			forDD = d.ForChained
		case Layered:
			oo = append(oo, d)
		}
		if cb(d) {
			return newLayers(m, oo) // TODO: coverage
		}
		if forDD == nil {
			continue
		}
		if d, ok := d.(Layered); ok {
			oo = append(oo, d)
		}
		forDD(func(d Dimer) (stop bool) {
			dd = append(dd, d)
			return false
		})
	}
	return newLayers(m, oo)
}

// Has returns true if given manager m has given dimer dr nested in
// given dimer in otherwise false.  If in is nil it defaults to m's root
// dimer.
func (m *Manager) Has(dr, in Dimer) bool {
	if in == nil {
		in = m.Root
	}
	d, dd, forDD := in, []Dimer{in}, (func(func(d Dimer) bool))(nil)
	for len(dd) > 0 {
		d, dd, forDD = dd[0], dd[1:], nil
		if dr == d {
			return true
		}
		switch d := d.(type) {
		case Stacker:
			forDD = d.ForStacked
		case Chainer:
			forDD = d.ForChained
		}
		if forDD == nil {
			continue
		}
		forDD(func(d Dimer) (stop bool) {
			dd = append(dd, d)
			return false
		})
	}
	return false
}
