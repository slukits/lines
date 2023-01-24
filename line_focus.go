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

// Next focuses the next focusable line at the currently focused line's
// cursor position if possible and returns its index as well as the cell
// index which defaults to -1 if CellsFocusable feature is not set.
// If highlighted is true the highlight of the current line is removed
// while the next is highlighted.
func (s *LineFocus) Next() (ln int, cl int) {
	if s.idx() >= 0 {
		s.line(s.idx()).resetLineFocus()
	}
	ln = s.findNextLine()
	if ln == s.current {
		s.Reset()
		s.c.Scroll.ToBottom()
		return s.current, -1
	}

	_, column, _ := s.c.CursorPosition()
	s.focus(ln)

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

// Previous focuses the first focusable line previous to
// [LineFocus.Current] and  returns its index along with the cursor
// position which defaults to -1 for an unset cursor.  If highlighted is
// true the highlight of the current line is removed while the previous
// line is highlighted.
func (s *LineFocus) Previous() (slIdx int, cl int) {
	if s.idx() >= 0 {
		s.line(s.idx()).resetLineFocus()
	}
	slIdx = s.findPrevious()
	if slIdx == s.current {
		s.Reset()
		s.c.Scroll.ToTop()
		return s.current, cl
	}

	_, column, _ := s.c.CursorPosition()
	s.focus(slIdx)
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

func (s *LineFocus) hlLineFlag() LineFlags {
	hlFlag := ZeroLineFlag
	if s.c.ff.has(HighlightEnabled) {
		hlFlag = Highlighted
	}
	if s.c.ff.has(TrimmedHighlightEnabled) {
		hlFlag = TrimmedHighlighted
	}
	return hlFlag
}

func (s *LineFocus) focus(idx int) {
	s.Reset()

	if idx == -1 {
		s.c.Scroll.To(0)
		return
	} else {
		s.c.Scroll.To(idx)
	}
	hlFlag := s.hlLineFlag()
	if idx != -1 && hlFlag != ZeroLineFlag {
		s.line(s.lineIndexOfContent(idx)).Switch(hlFlag)
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
	if idx-f.c.First() < 0 || idx-f.c.First() >= len(*f.c.ll) {
		return -1
	}
	return idx - f.c.First()
}

func (f *LineFocus) switchScrollingSourcedHighlight(scroll int) {
	if f.current == -1 {
		return
	}
	idx, hlFlag := f.current-f.c.First(), f.hlLineFlag()
	if idx >= 0 && idx < f.c.ContentScreenLines() {
		_, clm, crsr := f.c.CursorPosition()
		if crsr {
			f.cursor = clm
			f.c.SetCursor(-1, -1)
		}
		l := (*f.c.ll)[idx]
		if l.IsFlagged(hlFlag) {
			l.Switch(hlFlag)
		}
	}
	start := f.c.First() + scroll
	end := start + f.c.ContentScreenLines()
	if f.current < start || f.current >= end {
		return
	}
	l := (*f.c.ll)[f.current-start]
	if l.IsFlagged(hlFlag) {
		return
	}
	if f.cursor >= 0 {
		f.c.setCursor(f.current-start, f.cursor)
		f.cursor = -1
	}
	l.Switch(hlFlag)
}

// Reset removes a set line-focus switching of a potential highlight
// independent of given argument.
func (s *LineFocus) Reset() {
	if s.current == -1 {
		return
	}
	if s.c.Src == nil || s.onDisplay(s.current) {
		s.line(s.idx()).Unflag(Highlighted | TrimmedHighlighted)
	}
	s.current = -1
	if cc := s.c.gg.scr.cursorComponent(); cc != nil {
		if s.c.component == cc.wrapped() {
			s.c.gg.scr.setCursor(-1, -1, ZeroCursor)
		}
	}
}

func (s *LineFocus) onDisplay(idx int) bool {
	return idx >= s.c.First() &&
		idx < s.c.First()+s.c.ContentScreenLines()
}
