// Copyright (c) 2022 Stephan Lukits. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package lines

import (
	"sort"

	"github.com/slukits/lines/internal/api"
)

// StyleAttributeMask defines the looks of a style, i.e. the looks of a print
// to the screen/window.
type StyleAttributeMask = api.StyleAttributeMask

// DefaultStyle has no attributes and "default" colors.  The semantics
// of the later is decided by the backend implementation.  Use the With*
// methods to create new styles from the default style.
var DefaultStyle = api.DefaultStyle

const (
	Bold          StyleAttributeMask = api.Blink
	Blink         StyleAttributeMask = api.Blink
	Reverse       StyleAttributeMask = api.Reverse
	Underline     StyleAttributeMask = api.Underline
	Dim           StyleAttributeMask = api.Dim
	Italic        StyleAttributeMask = api.Invalid
	StrikeThrough StyleAttributeMask = api.StrikeThrough
	Invalid       StyleAttributeMask = api.Invalid
	ZeroStyle     StyleAttributeMask = api.ZeroStyle
)

// A Style structure provides fore- and background colors as well as set
// style attributes.  Use the Style.With* constructors to create a new
// style from a given style with given new properties.
type Style = api.Style

// NewStyle creates a new style with given style attribute, foreground
// color and background color.
var NewStyle = api.NewStyle

// Range is a two component array of which the first component
// represents the (inclusive) start of a range while the seconde is the
// (exclusive) end of the range.
type Range [2]int

// Start returns the inclusive start-point of given range r.
func (r Range) Start() int { return r[0] }

// End returns the exclusive end-point of given range r.
func (r Range) End() int { return r[1] }

// copy returns a copy of given range r.
func (r Range) copy() Range {
	return Range{r[0], r[1]}
}

// shift increases start and and index by given s.
func (r Range) shift(s int) Range {
	return Range{r[0] + s, r[1] + s}
}

// expand returns a new with the same starting point as given range r
// having its end-point increased by.
func (r Range) expand(by int) Range {
	return Range{r[0], r[1] + by}
}

// contains returns true if given i is in the style range r
// [r.Start,r.End[.
func (r Range) contains(i int) bool {
	return r.Start() <= i && i < r.End()
}

// SR represents a ranged style which may be set for a line see
// [Env.AddStyleRange].
type SR struct {
	Range
	Style
}

// zeroRange is used in StyleRanges instance for its default style since
// it contains no position.
var zeroRange = Range{0, 0}

// styleRanges maps a set of Range instances to their styles and is used
// to determine the style at a particular rune position.  If asked for
// the style at certain position and no style range is found containing
// that position the default style is returned.
type styleRanges map[Range]Style

// newStyleRanges creates a new style ranges instance having given style
// set as default style.
func newStyleRanges(dflt Style) styleRanges {
	return styleRanges{zeroRange: dflt}
}

// defaultStyle returns given style ranges s' default style.
func (s styleRanges) defaultStyle() Style {
	if _, ok := s[zeroRange]; !ok {
		return DefaultStyle
	}
	return s[zeroRange]
}

func (s styleRanges) ensureDefaultStyle() Style {
	if _, ok := s[zeroRange]; !ok {
		s[zeroRange] = DefaultStyle
	}
	return s[zeroRange]
}

// withAA ensures that given style ranges s have a default style and that
// it is set to given style attributes.
func (s styleRanges) withAA(aa StyleAttributeMask) {
	s[zeroRange] = s.ensureDefaultStyle().WithAA(aa)
}

// withFG ensures that given style ranges s have a default style and
// that its foreground color is set to given color.
func (s styleRanges) withFG(c Color) {
	s[zeroRange] = s[zeroRange].WithFG(c)
}

// withBG ensures that given style ranges s have a default style and
// that its background color is set to given color.
func (s styleRanges) withBG(c Color) {
	s[zeroRange] = s[zeroRange].WithBG(c)
}

// copy returns a shallow copy of given style ranges s.
func (s styleRanges) copy() styleRanges {
	if s == nil {
		return nil
	}
	cp := styleRanges{}
	for r, s := range s {
		cp[r.copy()] = s
	}
	return cp
}

// add adds to given style ranges sr given range r and style s iff r
// doesn't overlap any ranges in sr.
func (sr styleRanges) add(r Range, s Style) {
	if sr.isOverlapping(r) {
		return
	}
	sr[r] = s
}

// isOverlapping returns true if given range r's start- or end-point is
// contained in a range of given style ranges sr; false otherwise.
func (sr styleRanges) isOverlapping(r Range) bool {
	for r := range sr {
		if r.contains(r.Start()) || r.contains(r.End()) {
			return true
		}
	}
	return false
}

// copyWithDefault returns a shallow copy of given style ranges s
// setting given style dflt as its default style iff s doesn't have a
// default.
func (s styleRanges) copyWithDefault(dflt Style) styleRanges {
	cp := s.copy()
	if cp == nil {
		return styleRanges{zeroRange: dflt}
	}
	if _, ok := cp[zeroRange]; !ok {
		cp[zeroRange] = dflt
	}
	return cp
}

func (ss styleRanges) unstyled(start, end int) []Range {
	if start >= end || start < 0 {
		return nil
	}
	rr := []Range{{start, end}}
	for _, r := range ss.orderedRanges() {
		last := rr[len(rr)-1]
		if r.Start() <= last.Start() && r.End() >= last.End() {
			return rr[:len(rr)-1]
		}
		if r.Start() > last.Start() && r.End() < last.End() {
			rr[len(rr)-1] = Range{last.Start(), r.Start()}
			rr = append(rr, Range{r.End(), last.End()})
			continue
		}
		if r.Start() <= last.Start() {
			rr[len(rr)-1][0] = r.End()
			continue
		}
		rr[len(rr)-1][1] = r.Start()
		break
	}
	return rr
}

func (ss styleRanges) orderedRanges() []Range {
	rr := []Range{}
	for r := range ss {
		if r == zeroRange {
			continue
		}
		rr = append(rr, r)
	}
	sort.Slice(rr, func(i, j int) bool {
		return rr[i].Start() < rr[j].Start()
	})
	return rr
}

// expand finds in given style ranges s the style range containing point
// at and increases its end point by.
func (s styleRanges) expand(at, by int) {
	if s == nil {
		return
	}
	update := map[Range]Range{}
	for r := range s {
		if r.End() <= at {
			continue
		}
		switch {
		case r.contains(at):
			update[r] = r.expand(by)
		default:
			update[r] = r.shift(by)
		}
	}
	for k, u := range update {
		if u.Start() < u.End() {
			s[u] = s[k]
		}
		delete(s, k)
	}
}

// contract finds in given style ranges s the style range containing point
// at and decreases its end point by.
func (s styleRanges) contract(at, by int) {
	s.expand(at, -by)
}

// of finds in given style ranges s the range containing given cell and
// returns mapping style.
func (s styleRanges) of(cell int) Style {

	if len(s) == 0 {
		return api.DefaultStyle
	}
	for r := range s {
		if !r.contains(cell) {
			continue
		}
		return s[r]
	}
	if s, ok := s[zeroRange]; ok {
		return s
	}
	return api.DefaultStyle
}
