// Copyright (c) 2022 Stephan Lukits. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package lines

// For the definition of LineFocus see line_focus.go here the
// functionality for operations on cell-level are implemented.

// EolAfterLastRune enables the cursor to go one rune over the last
// content rune of a line.
func (s *LineFocus) EolAfterLastRune() *LineFocus {
	s.eolAfterLastRune = true
	return s
}

// EolAtLastRune restricts a line's most right cursor position to the
// lines last content rune.
func (s *LineFocus) EolAtLastRune() *LineFocus {
	s.eolAfterLastRune = false
	return s
}

// Eol returns true iff the cursor is at the most right position of the
// current line's content.
func (f *LineFocus) Eol() bool {
	// grab screen-line and screen-cell indices
	_, scIdx, hasCursor := f.c.wrapped().cursorPosition()
	if !hasCursor {
		return false
	}
	return f.isEol(f.Line(), scIdx)
}

func (f *LineFocus) isEol(line *Line, columnIdx int) bool {
	availableContent := len(line.rr) - line.start
	if !f.eolAfterLastRune {
		return columnIdx+1 == availableContent
	}
	return columnIdx == availableContent
}

// NextCell moves the cursor to the next cell in the currently focused
// line and returns the later's screen line index with the cell index of
// the cursor position and a boolean value indicating if the cursor was
// moved.
func (s *LineFocus) NextCell() (slIdx, scIdx int, moved bool) {
	slIdx, scIdx, haveCursorPos := s.c.CursorPosition()
	if s.current < 0 || !haveCursorPos {
		return -1, -1, false
	}
	if slIdx != s.Screen() {
		panic("lines: line-focus: last cell: cursor-line is not " +
			"focused line")
	}
	line := s.Line()
	if s.isEol(line, scIdx) {
		return slIdx, scIdx, false
	}
	_, _, screenWidth, _ := s.c.ContentArea()
	if scIdx+1 < screenWidth {
		s.c.setCursor(slIdx, scIdx+1)
		return slIdx, scIdx + 1, true
	}
	if s.eolAfterLastRune {
		screenWidth--
	}
	line.incrementStart(screenWidth)
	return slIdx, scIdx, false
}

// LastCell moves the cursor of currently focused component line to the
// right most non empty screen column and moves the content so far to
// the left in case of an overflowing line that the last content rune is
// in the component's last screen column.  Last cell returns the
// currently focused line, the currently focused cell and if the cursor
// was moved.
func (s *LineFocus) LastCell() (slIdx, scIdx int, moved bool) {
	slIdx, scIdx, _ = s.c.CursorPosition()
	if s.current < 0 || scIdx < 0 {
		return -1, -1, false
	}
	if slIdx != s.Screen() {
		panic("lines: line-focus: last cell: cursor-line is not " +
			"focused line")
	}
	line := s.Line()
	if s.isEol(line, scIdx) {
		return slIdx, scIdx, false
	}
	_, _, screenWidth, _ := s.c.ContentArea()
	if s.eolAfterLastRune {
		screenWidth--
	}
	if line.Len() < screenWidth {
		scIdx = line.Len()
		if !s.eolAfterLastRune {
			scIdx--
		}
		s.c.setCursor(slIdx, scIdx)
		return slIdx, scIdx, true
	}
	scIdx = screenWidth
	if !s.eolAfterLastRune {
		scIdx--
	}
	slIdx, scIdx, _ = s.c.SetCursor(slIdx, scIdx).CursorPosition()
	line.moveStartToEnd(screenWidth)
	return slIdx, scIdx, true
}

// adjustLineEndCursor is a helper for Previous to move the cursor
// onto the right position on the pervious line.
func (s *LineFocus) adjustLineEndCursor(
	lastColumn int, f FeatureMask,
) int {
	if !s.c.ff.has(f) {
		return -1
	}
	if lastColumn == -1 {
		return 0
	}

	l := s.Line()
	if l.Len() <= lastColumn {
		return l.Len() - 1
	}
	return lastColumn
}

// FirstCell moves the cursor of the currently focused component line to
// its first screen column and its content to the right that the first
// content rune is in the first column.
func (s *LineFocus) FirstCell() (slIdx, cl int, moved bool) {
	_, cl, _ = s.c.CursorPosition()
	if s.current < 0 || cl < 0 {
		return -1, -1, false
	}
	if cl != 0 {
		moved = true
	}
	slIdx, cl, _ = s.c.SetCursor(s.Screen(), 0).CursorPosition()
	s.Line().resetLineFocus()
	return slIdx, cl, moved
}

func (s *LineFocus) previousFromSource(fl FocusableLiner) int {
	initI := s.current - 1
	if s.current == -1 {
		initI = fl.Len() - 1
	}
	for i := initI; i >= 0; i-- {
		if !fl.IsFocusable(i) {
			continue
		}
		return i
	}
	return s.current
}

// PreviousCell moves a components cursor to the previous cell in the
// currently focused line if possible and returns the screen line index,
// cell index and a boolean indicating if the cursor was moved.  Last
// cell returns the currently focused line, the currently focused cell
// and if the cursor was moved.
func (s *LineFocus) PreviousCell() (slIdx, cl int, moved bool) {
	slIdx, cl, _ = s.c.CursorPosition()
	if s.current < 0 || cl < 0 {
		return -1, -1, false
	}
	if cl > 0 {
		return s.c.SetCursor(s.Screen(), cl-1).CursorPosition()
	}
	s.Line().decrementStart()
	return slIdx, cl, false
}
