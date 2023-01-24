// Copyright (c) 2022 Stephan Lukits. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package lines

import (
	"sort"
	"strings"
)

// LineFlags control the behavior and layout of a displayed line.
type LineFlags uint64

const (
	dirty LineFlags = 1 << iota

	// NotFocusable flagged line can not be focused
	NotFocusable

	// Highlighted flagged line has additive the global highlight style applied
	Highlighted

	// TrimmedHighlighted flagged line has additive the global highlight
	// style applied except on leading and trailing spaces.
	TrimmedHighlighted

	// ZeroLineFlag is the LineFlags zero type
	ZeroLineFlag LineFlags = 0
)

// A Line structure stores a Line's flags, its content,
// style-information and line filler and provides methods to modify its
// flags and styles.  A [Component]'s Line instance is obtained through
// its LL property (see [ComponentLines.By]).
type Line struct {
	// ff a lines state and format flags like *dirty or *Highlighted*.
	// ff is reset if reset() is called.
	ff LineFlags
	// rr is a lines content and reset by a call of reset().
	rr []rune
	// ss styles the runes of a line and reset by a call of reset().
	ss styleRanges
	// fillAt indicates line filling runes and reset by a call of
	// reset().
	fillAt []int
	// start is the first rune index of the displayed line content rr.
	// start is reset if reset() or resetLineFocus() is called.
	start int
	// ofRight stores a line's overflow at the right side since it was
	// asked for overflowing (see isOverflowing) the last time.  ofRight
	// is reset if reset() or resetLineFocus() is called.
	ofRight bool
	// ofRight stores a line's overflow at the left side since it was
	// asked for overflowing (see isOverflowing) the last time.  ofLeft
	// is reset if reset() or resetLineFocus() is called.
	ofLeft bool
}

func (l *Line) Len() int {
	return len(l.rr)
}

func (l *Line) isOverflowing(width int) (left, right, changed bool) {
	left = l.start > 0
	right = len(l.rr)-l.start > width
	if left == l.ofLeft && right == l.ofRight {
		return left, right, false
	}
	l.ofLeft, l.ofRight = left, right
	return left, right, true
}

func (l *Line) reset(ff LineFlags, s *Style) *Line {
	l.ff = ff
	l.rr = nil
	l.start = 0
	l.ofRight, l.ofLeft = false, false
	if s != nil {
		l.ss = newStyleRanges(*s)
	} else {
		l.ss = nil
	}
	l.fillAt = nil
	l.setDirty()
	return l
}

func (l *Line) setFlags(ff LineFlags) {
	if l.ff == ff {
		return
	}
	l.ff = ff | dirty
}

// Switch turns given flag(s) on if they are not all set otherwise these
// flags are removed.
func (l *Line) Switch(ff LineFlags) {
	if l.ff&ff == ff {
		l.setFlags(l.ff &^ ff)
		return
	}
	ff = cleanForFlagging(ff)
	switch {
	case l.ff&Highlighted != 0 && ff&TrimmedHighlighted != 0:
		l.ff &^= Highlighted
	case l.ff&TrimmedHighlighted != 0 && ff&Highlighted != 0:
		l.ff &^= TrimmedHighlighted
	}
	l.setFlags(l.ff | ff)
}

func (l *Line) Flag(ff LineFlags) {
	if l.ff&ff == ff {
		return
	}
	ff = cleanForFlagging(ff)
	l.setFlags(l.ff | ff)
}

// IsFlagged returns true if given line l has given flags ff set; false
// otherwise.
func (l *Line) IsFlagged(ff LineFlags) bool {
	return l.ff&ff == ff
}

// Unflag removes given flags ff from given line l.
func (l *Line) Unflag(ff LineFlags) { l.setFlags(l.ff &^ ff) }

// cleanForFlagging removes inconsistent flags preferring usually
// the more specific one.
func cleanForFlagging(ff LineFlags) LineFlags {
	highlight := Highlighted | TrimmedHighlighted
	switch {
	case ff&highlight == highlight:
		ff &^= Highlighted
	}
	return ff
}

// incrementStart moves overflowing content one rune to the left for the
// use case that the cursor is in a component's last content column and
// the user goes to the right.  incrementStart is a no-op if there are
// no overflowing runes.
func (l *Line) incrementStart(width int) {
	if len(l.rr) == 0 || len(l.rr[l.start:]) <= width {
		return
	}
	l.start++
	l.setDirty()
}

func (l *Line) decrementStart() {
	if l.start == 0 {
		return
	}
	l.start--
	l.setDirty()
}

func (l *Line) moveStartToEnd(width int) {
	if len(l.rr) == 0 || len(l.rr[l.start:]) < width {
		return
	}
	l.start = len(l.rr) - width
	l.setDirty()
}

func (l *Line) resetLineFocus() {
	l.ofLeft, l.ofRight = false, false
	if l.start == 0 {
		return
	}
	l.start = 0
	l.setDirty()
}

// setDirty sets the dirty flag if not set.
func (l *Line) setDirty() {
	if l.isDirty() {
		return
	}
	l.ff |= dirty
}

func (l *Line) setClean() {
	if !l.isDirty() {
		return
	}
	l.ff &^= dirty
}

// isDirty returns true if the dirty flag is set.
func (l *Line) isDirty() bool {
	return l.ff&dirty != 0
}

// setDefaultStyle sets given line l's zeroRange style of its style
// ranges.
func (l *Line) setDefaultStyle(s Style) {
	l.ss = newStyleRanges(s)
	l.setDirty()
}

// ensureStyleRanges returns given line l's style ranges and initializes
// them first if necessary.
func (l *Line) ensureStyleRanges() styleRanges {
	if l.ss == nil {
		l.ss = newStyleRanges(DefaultStyle)
	}
	return l.ss
}

// withAA sets of given line l its default style's style attributes.
// NOTE if no default style has been set (see line.setDefaultStyle) the
// DefaultStyle is used and modified with given attributes.
func (l *Line) withAA(aa StyleAttributeMask) {
	l.ensureStyleRanges().withAA(aa)
	l.setDirty()
}

// withFG sets of given line l its default style's foreground color.
// NOTE if no default style has been set (see line.setDefaultStyle) the
// DefaultStyle is used and modified with given foreground color.
func (l *Line) withFG(c Color) {
	l.ensureStyleRanges().withFG(c)
	l.setDirty()
}

// withBG sets of given line l its default style's background color.
// NOTE if no default style has been set (see line.setDefaultStyle) the
// DefaultStyle is used and modified with given background color.
func (l *Line) withBG(c Color) {
	l.ensureStyleRanges().withBG(c)
	l.setDirty()
}

// set sets given line l's content converting given string s to a rune
// slice.
func (l *Line) set(s string) {
	l.rr = []rune(s)
	if len(l.fillAt) > 0 {
		l.fillAt = []int{}
	}
	if len(l.ss) > 0 {
		l.ss = nil
	}
	l.setDirty()
	for i, r := range l.rr {
		if r != filler {
			continue
		}
		l.fillAt = append(l.fillAt, i)
		l.rr[i] = ' '
	}
}

// setStyled sets given line l's content to the conversion of given
// string s to a rune slice and adds a style range for this slice with
// given style.  Note the main difference to set is that if the line has
// a default style d which differs from provided style and the display
// is wider than the l's content then the remaining blanks have the
// style d.
func (l *Line) setStyled(s string, sty Style) {
	l.set(s)
	l.ss = styleRanges{Range{0, len(l.rr)}: sty}
	l.setDirty()
}

// setAt sets given rune slice rr at given position p in given line l's
// content overwriting what has been at and after this position.  If
// needed rr is padded with spaces to p.
func (l *Line) setAt(p int, rr []rune) {
	if p < -1 {
		return
	}
	l.setDirty()
	if p == -1 {
		l.truncateAt(0)
		p = 0
	} else {
		l.truncateAt(p)
		l.padTo(p)
		l.rr = l.rr[:p]
	}
	l.rr = l.rr[:p]
	for i, r := range rr {
		if r == filler {
			r = ' '
			l.fillAt = append(l.fillAt, p+i)
		}
		l.rr = append(l.rr, r)
	}
}

// truncateAt truncates fillers and styles at and after given position
// p.
func (l *Line) truncateAt(p int) {
	truncate := -1
	for i, f := range l.fillAt {
		if f < p {
			continue
		}
		truncate = i
		break
	}
	if truncate >= 0 {
		l.fillAt = l.fillAt[:truncate]
	}
	for r := range l.ss {
		if r.Start() < p || r == zeroRange {
			continue
		}
		delete(l.ss, r)
	}
}

// setStyledAt sets given rune slice rr at given position p in given
// line l's content overwriting what has been at and after this position
// and adds a corresponding style range.  If needed l's content is
// padded with spaces until p.
func (l *Line) setStyledAt(at int, rr []rune, sty Style) {
	l.setAt(at, rr)
	if at == -1 {
		l.setDefaultStyle(sty)
		return
	}
	if l.ss == nil {
		l.ss = styleRanges{}
	}
	l.ss[Range{at, at + len(rr)}] = sty
}

// setAtFilling sets given rune r at given position p of given line l's
// content truncating all possibly following content.
func (l *Line) setAtFilling(p int, r rune) {
	l.padTo(p)
	l.truncateAt(p)
	l.rr = append(l.rr[:p], r)
	l.fillAt = append(l.fillAt, p)
	l.setDirty()
}

// setStyledAtFilling sets given rune r at given position p of given
// line l's content truncating all possibly following content and adds
// a corresponding style range.
func (l *Line) setStyledAtFilling(at int, r rune, sty Style) {
	l.setAtFilling(at, r)
	if l.ss == nil {
		l.ss = styleRanges{}
	}
	l.ss[Range{at, at + 1}] = sty
	l.setDirty()
}

// AddStyleRange adds given style ranges sr and rr to given line l's
// style ranges iff they don't overlap with already existing style ranges.
func (l *Line) AddStyleRange(sr SR, rr ...SR) {
	if l.ss == nil {
		l.ss = styleRanges{}
	}
	l.ss.add(sr.Range, sr.Style)
	for _, r := range rr {
		l.ss.add(r.Range, r.Style)
	}
	l.setDirty()
}

type runeWriter interface {
	Display(x, y int, r rune, s Style)
}

// sync writes given line l's expanded and styled runes at coordinates x
// and y to the screen with given rune writer rw.  Is l wider than given
// width w it is truncated at w.
func (l *Line) sync(x, y, width int, rw runeWriter, gg *globals) {
	l.setClean()
	rr, ss := l.display(width, gg)
	for i, r := range rr {
		if i == width {
			break
		}
		rw.Display(x+i, y, r, ss.of(i))
	}
}

func (l *Line) vsync(x, y, height int, rw runeWriter, gg *globals) {
	l.setClean()
	rr, ss := l.display(height, gg)
	for i, r := range rr {
		if i == height {
			break
		}
		rw.Display(x, y+i, r, ss.of(i))
	}
}

// display returns a line's calculated content depending on given width
// and set filler as well as corresponding style ranges ready to print
// to the screen.
func (l *Line) display(width int, gg *globals) ([]rune, styleRanges) {
	ss := l.ss.copyWithDefault(gg.Style(Default))
	if len(l.rr) == 0 {
		return l.displayEmpty(width, gg, ss)
	}
	rr := append([]rune{}, l.rr...)
	rr, ss, tc := l.expandLeadingTabs(rr, ss, gg.tabWidth)
	if len(rr) >= width {
		return l.displayOverflowing(width, gg, rr, ss)
	}
	if len(l.fillAt) > 0 {
		rr, ss = l.expandFillerAt(rr, width, ss, tc, gg)
	}
	if len(rr) < width {
		rr = l.pad(rr, width)
	}
	if l.ff&(Highlighted|TrimmedHighlighted) != 0 {
		ss = l.highlighted(rr, ss, gg)
	}
	return rr, ss
}

// displayEmpty returns width many space runes and adjust styles in case
// given line l is highlighted.
func (l *Line) displayEmpty(width int, g *globals, ss styleRanges) (
	[]rune, styleRanges,
) {
	rr := []rune(strings.Repeat(" ", width))
	if l.ff&(Highlighted|TrimmedHighlighted) != 0 {
		ss = l.highlighted(rr, ss, g)
	}
	return rr, ss
}

// displayOverflowing trims given runes to given width and adjusts given
// style ranges ss in case given line l is highlighted.
func (l *Line) displayOverflowing(
	width int, g *globals, rr []rune, ss styleRanges,
) ([]rune, styleRanges) {
	if l.ff&(Highlighted|TrimmedHighlighted) != 0 {
		ss = l.highlighted(rr, ss, g)
	}
	ll := len(rr[l.start:])
	if ll > width {
		return rr[l.start : l.start+width], ss
	}
	if ll < width { // shouldn't happen
		return l.pad(rr[l.start:], width), ss
	}
	return rr[l.start:], ss
}

// expandFillerAt expands runes marked as fillers by an equal amount in
// given runes rr to fit given width exactly adjusting given style
// ranges accordingly.  Note a previous tab-expansion shifting the
// filler positions is taken into account by evaluating given tag count
// tc and given globals providing the tab-width.
func (l *Line) expandFillerAt(
	rr []rune, width int, ss styleRanges, tc int, gg *globals,
) ([]rune, styleRanges) {
	fillAt := append([]int{}, l.fillAt...)
	if tc > 0 { // adjust to tab expansion
		tc *= (gg.TabWidth() - 1)
		for i, f := range fillAt {
			fillAt[i] = f + tc
		}
	}
	f := (width - (len(rr) - len(fillAt))) / len(fillAt)
	mf := (width - (len(rr) - len(fillAt))) % len(fillAt)
	ff := map[int][]rune{}
	for i := 0; i < len(fillAt)-mf; i++ {
		ff[fillAt[i]] = []rune(
			strings.Repeat(string(rr[fillAt[i]]), f))
	}
	for i := len(fillAt) - mf; i < len(fillAt); i++ {
		ff[fillAt[i]] = []rune(
			strings.Repeat(string(rr[fillAt[i]]), f+1))
	}
	_rr, last := []rune{}, 0
	for _, at := range fillAt {
		_rr = append(_rr, rr[last:at]...)
		_rr = append(_rr, ff[at]...)
		last = at + 1
	}
	if last < len(rr) {
		_rr = append(_rr, rr[last:]...)
	}
	return _rr, l.adjustFillerExpansionStyles(ff, ss)
}

func (l *Line) adjustFillerExpansionStyles(
	filler map[int][]rune, ss styleRanges,
) styleRanges {
	ff := make([]int, 0, len(filler))
	for k := range filler {
		ff = append(ff, k)
	}
	sort.Ints(ff)
	fills := 0 // move filler positions in ff along with expansions
	for i, f := range ff {
		ss.expand(f+fills, len(filler[f])-1)
		fills += len(filler[ff[i]]) - 1
	}
	return ss
}

func (l *Line) expandLeadingTabs(
	rr []rune, ss styleRanges, tabWidth int,
) ([]rune, styleRanges, int) {

	if len(rr) == 0 || rr[0] != '\t' {
		return rr, ss, 0
	}
	tc := 0
	for _, r := range rr {
		if r != '\t' {
			break
		}
		tc++
	}
	for i := 0; i < tc; i++ {
		ss.expand(i*tabWidth, tabWidth-1)
		rr = append(rr[:i*tabWidth], append(
			[]rune(strings.Repeat(" ", tabWidth)),
			rr[i*tabWidth+1:]...,
		)...)
	}
	return rr, ss, tc
}

// highlighted highlights the whole line; i.e. the global highlight
// style is applied by switching its style attributes and setting
// non-default colors for each style range.
func (l *Line) highlighted(
	rr []rune, ss styleRanges, gg *globals,
) styleRanges {

	if l.ff&TrimmedHighlighted != 0 {
		return l.highlightTrimmed(rr, ss, gg)
	}
	for r, s := range ss {
		ss[r] = gg.Highlight(s)
	}
	return ss
}

// highlightTrimmed highlights a line from its first non-space rune to
// its last including.  Style ranges overlapping trim-points are split
// accordingly before the sub-range inside the trim-range is highlighted
// by the global highlight style, i.e. switching its style attributes
// and setting non-default colors.
func (l *Line) highlightTrimmed(
	rr []rune, ss styleRanges, gg *globals,
) styleRanges {
	tl, tr := trim(rr)
	if tl == len(rr) {
		tl, tr = tr, tl
	}
	kk := []Range{}
	for r := range ss {
		if r == zeroRange {
			continue
		}
		kk = append(kk, r)
	}
	for _, r := range kk {
		switch {
		case r.Start() < tl && r.End() > tl && r.End() <= tr:
			ss[Range{r.Start(), tl}] = ss[r]
			ss[Range{tl, r.End()}] = gg.Highlight(ss[r])
			delete(ss, r)
		case r.Start() >= tl && r.End() <= tr:
			ss[r] = gg.Highlight(ss[r])
		case r.Start() >= tl && r.Start() < tr && r.End() > tr:
			ss[Range{tr, r.End()}] = ss[r]
			ss[Range{r.Start(), tr}] = gg.Highlight(ss[r])
			delete(ss, r)
		case r.Start() < tl && r.End() > tr:
			ss[Range{r.Start(), tl}] = ss[r]
			ss[Range{tr, r.End()}] = ss[r]
			ss[Range{tl, tr}] = gg.Highlight(ss[r])
			delete(ss, r)
		}
	}
	urr := ss.unstyled(tl, tr)
	if len(urr) == 0 {
		return ss
	}
	for _, r := range urr {
		ss[r] = gg.Highlight(ss[zeroRange])
	}
	return ss
}

// trim returns the index of the first non-space rune and the index
// after the last non-space rune.
func trim(rr []rune) (int, int) {
	tl := trimLeft(rr)
	if tl == len(rr) {
		return len(rr), 0
	}
	return tl, trimRight(rr)
}

// trimLeft returns the index of the first non-space rune.
func trimLeft(rr []rune) int {
	for i, r := range rr {
		if r == ' ' {
			continue
		}
		return i
	}
	return len(rr)
}

// trimRight returns the index of the last non-space rune.
func trimRight(rr []rune) int {
	for i := len(rr) - 1; i >= 0; i-- {
		if rr[i] == ' ' {
			continue
		}
		return i + 1
	}
	return 0
}

// padTo pads a line's l content rr to given position p with spaces.
func (l *Line) padTo(p int) {
	if len(l.rr) >= p {
		return
	}
	l.rr = append(l.rr, []rune(strings.Repeat(" ", p-len(l.rr)))...)
}

// pad returns a rune slice whose len is given width by padding given
// rune slice rr with spaces accordingly.  pad will panic if len(rr) >
// width.
func (l *Line) pad(rr []rune, width int) []rune {
	c := width - len(rr)
	if c < 0 {
		return rr
	}
	return append(rr, []rune(strings.Repeat(" ", c))...)
}
