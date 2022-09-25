// Copyright (c) 2022 Stephan Lukits. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package lines

// A LineFocus instance is associated with each initialized
// lines-component and provides the api for focusing component lines.
type LineFocus struct {
	c       *Component
	current int
}

// Current returns the line-index of the currently focused line.
func (s *LineFocus) Current() int {
	return s.current
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
			(*s.c.ll)[s.current].SwitchHighlighted()
		}
		if highlighted {
			l.SwitchHighlighted()
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
			(*s.c.ll)[s.current].SwitchHighlighted()
		}
		if highlighted {
			(*s.c.ll)[i].SwitchHighlighted()
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
	if (*s.c.ll)[s.current].IsHighlighted() {
		(*s.c.ll)[s.current].SwitchHighlighted()
	}
	s.current = -1
	return s.current
}
