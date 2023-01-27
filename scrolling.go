// Copyright (c) 2022 Stephan Lukits. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package lines

type ScrollBarDef struct {
	AtLeft   bool
	GapIndex int
	Style    Style
	Position Style
}

func DefaultScrollbarDef() ScrollBarDef {
	return ScrollBarDef{Style: DefaultStyle.Reverse(), Position: DefaultStyle}
}

// Scroller provides a component's scrolling API.
type Scroller struct {
	c   *Component
	Bar bool
}

// IsAtTop returns true if the first screen line is the first component
// line.
func (s Scroller) IsAtTop() bool { return s.c.First() == 0 }

// IsAtBottom is true if a component's printable area contains the
// component's last line.
func (s Scroller) IsAtBottom() bool {
	return s.c.First()+s.c.ContentScreenLines() >= s.c.Len()
}

// Up scrolls one page up.  Whereas "one page" is in case of a component
// height of 1 is one line.  For a height h with 1 < h < 20 "one page"
// is h - 1.  For h >= 20 "one page" is h - h/10.
func (s Scroller) Up() {
	height, scroll := s.c.ContentScreenLines(), 0
	if height <= 0 || s.c.First() == 0 {
		return
	}
	switch {
	case height == 1:
		scroll = 1
	case height < 20:
		scroll = height - 1
	default:
		scroll = height - (height / 10)
	}
	if scroll >= s.c.First() {
		scroll = s.c.First()
	}
	s.c.LL.Focus.switchScrollingSourcedHighlight(-scroll)
	s.c.setFirst(s.c.First() - scroll)
}

// ToTop scrolls a component's content to its first line, i.e. the first
// screen line displays the first component line.
func (s Scroller) ToTop() { s.c.setFirst(0) }

// ToBottom scrolls to the index that the last screen line displays the
// last component line.
func (s Scroller) ToBottom() {
	height := s.c.ContentScreenLines()
	if height <= 0 {
		return
	}
	s.c.setFirst(s.c.Len() - height)
}

// Down scrolls one page down.  Whereas "one page" is in case of a
// component height of 1 is one line.  For a height h with 1 < h < 20
// "one page" is h - 1.  For h >= 20 "one page" is h - h/10.
func (s Scroller) Down() {
	height, scroll := s.c.ContentScreenLines(), 0
	if height <= 0 || height >= s.c.Len() {
		return
	}
	switch {
	case height == 1:
		scroll = 1
	case height < 20:
		scroll = height - 1
	default:
		scroll = height - (height / 10)
	}
	if s.c.Len()-(s.c.First()+scroll) < height {
		scroll = (s.c.Len() - height) - s.c.First()
	}
	s.c.LL.Focus.switchScrollingSourcedHighlight(scroll)
	s.c.setFirst(s.c.First() + scroll)
}

// To scrolls to the index that the line with given index is displayed.
func (s Scroller) To(idx int) {
	height := s.c.ContentScreenLines()
	if height <= 0 {
		return
	}
	if s.c.First() <= idx && idx < s.c.First()+height {
		return
	}
	if idx <= 0 {
		s.ToTop()
		return
	}
	if idx >= s.c.Len() {
		s.ToBottom()
		return
	}

	for s.c.First() > idx {
		s.Up()
	}
	for idx >= s.c.First()+height {
		s.Down()
	}
}

func (s Scroller) scrollBarGap(c *component) (
	w *GapWriter, sbd ScrollBarDef, ww *GapsWriter,
) {
	sbd = c.globals().ScrollBarDef()
	if c.gaps == nil {
		c.gaps = newGaps(c.gg.Style(Default))
	}
	ww = newGapsWriter(sbd.GapIndex, c.gaps)
	if sbd.AtLeft {
		return ww.Left, sbd, ww
	}
	return ww.Right, sbd, ww
}

func (s Scroller) setScrollBar(pos int, c *component) {
	gw, sbd, _ := s.scrollBarGap(c)
	if pos == -1 {
		Print(gw.Sty(sbd.Style).At(0).Filling(), ' ')
		return
	}
	if pos == 0 {
		Print(gw.At(0).Sty(sbd.Position), ' ')
		Print(gw.At(1).Sty(sbd.Style).Filling(), ' ')
		return
	}
	if pos+1 >= c.ContentScreenLines() {
		Print(gw.At(0).Sty(sbd.Style).Filling(), ' ')
		Print(gw.At(c.ContentScreenLines()-1).Sty(sbd.Position), ' ')
		return
	}
	Print(gw.At(0).Sty(sbd.Style).Filling(), ' ')
	Print(gw.At(pos).Sty(sbd.Position), ' ')
	Print(gw.At(pos+1).Sty(sbd.Style).Filling(), ' ')
}

func (s Scroller) updateBar() {
	if !s.Bar {
		return
	}
	c := s.c.layoutCmp.wrapped()
	if c.ContentScreenLines() >= c.Len() {
		s.setScrollBar(-1, c)
		return
	}
	if c.First() == 0 {
		s.setScrollBar(0, c)
		return
	}
	if c.First()+c.ContentScreenLines() >= c.Len() {
		s.setScrollBar(c.ContentScreenLines()-1, c)
		return
	}
	s.setScrollBar(
		(c.First()+c.ContentScreenLines()/2)/
			(c.Len()/c.ContentScreenLines()), c)
}

func (s Scroller) BarContains(x, y int) bool {
	if s.c == nil || s.c.gaps == nil {
		return false
	}
	top, _, bottom, _ := s.c.GapsLen()
	if y < s.c.dim.Y()+top {
		return false
	}
	if y >= s.c.dim.Y()+s.c.dim.Height()-bottom {
		return false
	}

	sbd := s.c.globals().ScrollBarDef()
	if !sbd.AtLeft {
		return x == s.c.dim.X()+s.c.dim.Width()-(sbd.GapIndex+1)
	}
	return x == s.c.dim.X()+sbd.GapIndex
}
