// Copyright (c) 2022 Stephan Lukits. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package lines

// A LineFocus instance is associated with each initialized [Component]
// through its LL property (see [ComponentLines]) and provides the api
// for focusing component lines.
type LineFocus struct {
	c *Component

	// current is the currently as focused flagged line
	current int

	// cursor keeps track of the cursor position during scrolling
	cursor int

	// hlType indicates if highlighted at all and in case of the later
	// if the whole line should be highlighted or if highlighting should
	// be trimmed to non-blank content.
	hlType LineFlags

	// eolAfterLastRune indicates that the last position of the cursor
	// in a given line is after the last rune which is needed for
	// editable components to append to a line.
	eolAfterLastRune bool
}

// Current returns the content line-index of the currently focused line.
func (s *LineFocus) Current() int { return s.current }

// Screen returns the screen line-index of the currently focused line.
func (s *LineFocus) Screen() int { return s.idx() }

// Trimmed switches the highlight-style to "trimmed-highlighted", i.e.
// only a line's space-trimmed content is highlighted and not the whole
// line.  Note the next call of Trimmed switches back to full line
// highlighting.
func (s *LineFocus) Trimmed() {
	if s.hlType != Highlighted {
		s.hlType = Highlighted
		return
	}
	s.hlType = TrimmedHighlighted
}

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
	slIdx, scIdx, hasCursor := f.c.wrapped().cursorPosition()
	if !hasCursor {
		return false
	}
	return f.isEol(f.line(slIdx), scIdx)
}

func (f *LineFocus) isEol(line *Line, columnIdx int) bool {
	availableContent := len(line.rr) - line.start
	if !f.eolAfterLastRune {
		return columnIdx+1 == availableContent
	}
	return columnIdx == availableContent
}

// Next focuses the next focusable line at the currently focused line's
// cursor position if possible and returns its index as well as the cell
// index which defaults to -1 if CellsFocusable feature is not set.
// If highlighted is true the highlight of the current line is removed
// while the next is highlighted.
func (s *LineFocus) Next(highlighted bool) (ln int, cl int) {
	if s.current >= 0 {
		s.line(s.idx()).resetLineFocus()
	}
	ln = s.findNextLine()
	if ln == s.current {
		s.Reset()
		s.c.Scroll.ToBottom()
		return s.current, -1
	}

	_, column, _ := s.c.CursorPosition()
	s.focus(ln, highlighted)

	if cl = s.adjustLineEndCursor(column, NextCellFocusable); cl >= 0 {
		s.c.SetCursor(s.idx(), cl)
	}

	return ln, cl
}

func (s *LineFocus) findNextLine() int {
	if s.c.Src != nil {
		if fl, ok := s.c.Src.Liner.(FocusableLiner); ok {
			return s.nextFromSource(fl)
		}
	}
	if s.current+1 >= s.c.Len() {
		return s.current
	}
	for idx, l := range (*s.c.ll)[s.current+1:] {
		if l.ff&NotFocusable == NotFocusable {
			continue
		}
		return s.current + 1 + idx
	}
	return s.current
}

func (s *LineFocus) nextFromSource(fl FocusableLiner) int {
	if s.current+1 >= fl.Len() {
		return s.current
	}
	for i := s.current + 1; i < fl.Len(); i++ {
		if !fl.IsFocusable(i) {
			continue
		}
		return i
	}
	return s.current
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
	if slIdx != s.idx() {
		panic("lines: line-focus: last cell: cursor-line is not " +
			"focused line")
	}
	line := s.line(slIdx)
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
	if slIdx != s.idx() {
		panic("lines: line-focus: last cell: cursor-line is not " +
			"focused line")
	}
	line := s.line(slIdx)
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

// Previous focuses the first focusable line previous to
// [LineFocus.Current] and  returns its index along with the cursor
// position which defaults to -1 for an unset cursor.  If highlighted is
// true the highlight of the current line is removed while the previous
// line is highlighted.
func (s *LineFocus) Previous(highlighted bool) (slIdx int, cl int) {
	if s.current >= 0 {
		s.line(s.idx()).resetLineFocus()
	}
	slIdx = s.findPrevious()
	if slIdx == s.current {
		s.Reset()
		s.c.Scroll.ToTop()
		return s.current, cl
	}

	_, column, _ := s.c.CursorPosition()
	s.focus(slIdx, highlighted)
	cl = s.adjustLineEndCursor(column, PreviousCellFocusable)
	if cl >= 0 {
		s.c.SetCursor(s.idx(), cl)
	}
	return slIdx, cl
}

func (s *LineFocus) findPrevious() int {
	if s.c.Src != nil {
		if fl, ok := s.c.Src.Liner.(FocusableLiner); ok {
			return s.previousFromSource(fl)
		}
	}
	initI := s.current - 1
	if s.current == -1 {
		initI = len(*s.c.ll) - 1
	}
	for i := initI; i >= 0; i-- {
		if (*s.c.ll)[i].ff&NotFocusable == NotFocusable {
			continue
		}
		return i
	}
	return s.current
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

	l := s.line(s.idx())
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
	slIdx, cl, _ = s.c.SetCursor(s.idx(), 0).CursorPosition()
	s.line(s.idx()).resetLineFocus()
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
		return s.c.SetCursor(s.idx(), cl-1).CursorPosition()
	}
	s.line(s.idx()).decrementStart()
	return slIdx, cl, false
}

func (s *LineFocus) focus(idx int, highlighted bool) {
	s.Reset()

	if idx == -1 {
		s.c.Scroll.To(0)
		return
	} else {
		s.c.Scroll.To(idx)
	}

	if idx != -1 && highlighted {
		s.line(s.lineIndexOfContent(idx)).Switch(s.hlType)
	}

	s.current = idx
}

// line returns the screen line with given index of the component
// associated with given line focus f.  See f.lineIndexOfContent and
// f.idx to map (current) content indices to screen line indices.
func (f *LineFocus) line(idx int) *Line {
	return (*f.c.ll)[idx]
}

// idx maps current content-line index *f.current* to its screen line
// index.
func (f *LineFocus) idx() int { return f.lineIndexOfContent(f.current) }

// lineIndexOfContent returns the screen line index displaying the
// content-line with given index idx.
func (f *LineFocus) lineIndexOfContent(idx int) int {
	if f.c.Src == nil {
		return idx
	}
	if idx-f.c.first() < 0 {
		return -1
	}
	return idx - f.c.first()
}

func (f *LineFocus) switchScrollingSourcedHighlight(scroll int) {
	if f.current == -1 {
		return
	}
	idx := f.current - f.c.first()
	if idx >= 0 && idx < f.c.contentScreenLines() {
		_, clm, crsr := f.c.CursorPosition()
		if crsr {
			f.cursor = clm
			f.c.SetCursor(-1, -1)
		}
		l := (*f.c.ll)[idx]
		if l.IsFlagged(f.hlType) {
			l.Switch(f.hlType)
		}
	}
	start := f.c.first() + scroll
	end := start + f.c.contentScreenLines()
	if f.current < start || f.current >= end {
		return
	}
	l := (*f.c.ll)[f.current-start]
	if l.IsFlagged(f.hlType) {
		return
	}
	if f.cursor >= 0 {
		f.c.setCursor(f.current-start, f.cursor)
		f.cursor = -1
	}
	l.Switch(f.hlType)
}

// Reset removes a set line-focus switching of a potential highlight
// independent of given argument.
func (s *LineFocus) Reset() {
	if s.current == -1 {
		return
	}
	if s.c.Src == nil || s.onDisplay(s.current) {
		l := s.line(s.idx())
		if l.IsFlagged(s.hlType) {
			l.Switch(s.hlType)
		}
	}
	s.current = -1
	if cc := s.c.gg.scr.cursorComponent(); cc != nil {
		if s.c.component == cc.wrapped() {
			s.c.gg.scr.setCursor(-1, -1, ZeroCursor)
		}
	}
}

func (s *LineFocus) onDisplay(idx int) bool {
	return idx >= s.c.first() &&
		idx < s.c.first()+s.c.contentScreenLines()
}
