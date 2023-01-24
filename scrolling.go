// Copyright (c) 2022 Stephan Lukits. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package lines

// Scroller provides a component's scrolling API.
type Scroller struct{ c *Component }

// IsAtTop returns true if the first screen line is the first component
// line.
func (s Scroller) IsAtTop() bool { return s.c.First() == 0 }

// IsAtBottom is true if a component's printable area contains the
// component's last line.
func (s Scroller) IsAtBottom() bool {
	return s.c.First()+s.c.ContentScreenLines() >= s.c.Len()
}

// CoordinateToIndex maps a y-coordinate relative to the components
// origin to its line index taking potential scrolling offsets into
// account.
func (s Scroller) CoordinateToIndex(y int) (line int) {
	if s.c.First() == 0 {
		return y
	}
	// TODO: add error handling
	if s.c.First()+y >= s.c.Len() {
		return s.c.Len() - 1
	}

	return s.c.First() + y
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
// screen line is the first component line.
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
