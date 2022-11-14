// Copyright (c) 2022 Stephan Lukits. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package lines

// A LineFocus instance is associated with each initialized [Component]
// through its LL property (see [ComponentLines]) and provides the api
// for focusing component lines.
type LineFocus struct {
	c       *Component
	current int
	hlType  LineFlags
}

// Current returns the line-index of the currently focused line.
func (s *LineFocus) Current() int {
	return s.current
}

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

// Next focuses the next focusable line after [LineFocus.Current] and
// returns its index.  If highlighted is true the highlight of the
// current line is removed while the next is highlighted.
func (s *LineFocus) Next(highlighted bool) int {
	next := s.findNext()
	if next == s.current {
		s.Reset(highlighted)
		return s.current
	}

	s.focus(next, highlighted)
	return next
}

func (s *LineFocus) findNext() int {
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

// Previous focuses the first focusable line previous to [LineFocus.Current]
// and  returns its index.  If highlighted is true the highlight of the
// current line is removed while the previous line is highlighted.
func (s *LineFocus) Previous(highlighted bool) int {
	prvs := s.findPrevious()
	if prvs == s.current {
		s.Reset(highlighted)
		return s.current
	}

	s.focus(prvs, highlighted)
	return prvs
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

func (s *LineFocus) focus(idx int, highlighted bool) {
	s.Reset(highlighted)

	if idx == -1 {
		s.c.Scroll.To(0)
	} else {
		s.c.Scroll.To(idx)
	}

	if idx != -1 && highlighted {
		s.line(idx).Switch(s.hlType)
	}

	s.current = idx
}

func (f *LineFocus) line(idx int) *Line {
	if f.c.Src == nil {
		return (*f.c.ll)[idx]
	}
	if idx-f.c.first() < 0 {
		panic("lines: component: focusable: line index out of range")
	}
	return (*f.c.ll)[idx-f.c.first()]
}

func (f *LineFocus) switchScrollingSourcedHighlight(scroll int) {
	if f.current == -1 {
		return
	}
	idx := f.current - f.c.first()
	if idx >= 0 && idx < f.c.contentScreenLines() {
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
	l.Switch(f.hlType)
}

// Reset removes a set line-focus switching of a potential highlight
// independent of given argument.
func (s *LineFocus) Reset(_ bool) int {
	if s.current == -1 {
		return s.current
	}
	if s.c.Src == nil || s.onDisplay(s.current) {
		l := s.line(s.current)
		if l.IsFlagged(s.hlType) {
			l.Switch(s.hlType)
		}
	}
	s.current = -1
	return s.current
}

func (s *LineFocus) onDisplay(idx int) bool {
	return idx >= s.c.first() &&
		idx < s.c.first()+s.c.contentScreenLines()
}
