// Copyright (c) 2022 Stephan Lukits. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package lines

// A LineFocus instance is associated with each initialized
// lines-component and provides the api for focusing component lines.
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

// Next returns the index of the next (relative to Current) selectable
// line l.  l is flagged as highlighted while the line highlighted
// before is not highlighted anymore.
func (s *LineFocus) Next(highlighted bool) int {
	if s.current+1 == len(*s.c.ll) {
		s.Reset(highlighted)
		return s.current
	}
	old := s.current
	for idx, l := range (*s.c.ll)[s.current+1:] {
		if l.ff&NotFocusable == NotFocusable {
			continue
		}
		if s.current >= 0 && highlighted {
			(*s.c.ll)[s.current].Switch(s.hlType)
		}
		if highlighted {
			l.Switch(s.hlType)
		}
		s.current = s.current + 1 + idx
		break
	}
	if old == s.current {
		s.Reset(highlighted)
		return s.current
	}
	s.c.Scroll.To(s.current)
	return s.current
}

func (s *LineFocus) Previous(highlighted bool) int {
	initI := s.current - 1
	if s.current == -1 {
		initI = len(*s.c.ll) - 1
	}
	old := s.current
	for i := initI; i >= 0; i-- {
		if (*s.c.ll)[i].ff&NotFocusable == NotFocusable {
			continue
		}
		if s.current >= 0 && highlighted {
			(*s.c.ll)[s.current].Switch(s.hlType)
		}
		if highlighted {
			(*s.c.ll)[i].Switch(s.hlType)
		}
		s.current = i
		break
	}
	if old == s.current {
		s.Reset(highlighted)
		return s.current
	}
	s.c.Scroll.To(s.current)
	return s.current
}

// Reset removes a set line-focus switching of a potential highlight
// independent of given argument.
func (s *LineFocus) Reset(_ bool) int {
	if s.current == -1 {
		return s.current
	}
	if (*s.c.ll)[s.current].IsFlagged(s.hlType) {
		(*s.c.ll)[s.current].Switch(s.hlType)
	}
	s.current = -1
	return s.current
}
