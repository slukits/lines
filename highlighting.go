// Copyright (c) 2022 Stephan Lukits. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package lines

// A Highlighter instance is associated with each initialized
// lines-component and provides the api for highlighting and selecting
// lines in a component.
type Highlighter struct {
	c       *Component
	current int
}

// Current returns the line-index of the currently highlighted line.
func (s *Highlighter) Current() int {
	return s.current
}

// Next returns the index of the next (relative to Current) selectable
// line l.  l is flagged as highlighted while the line highlighted
// before is not highlighted anymore.
func (s *Highlighter) Next() int {
	if s.current+1 == len(*s.c.ll) {
		s.Reset()
		return s.current
	}
	old := s.current
	for idx, l := range (*s.c.ll)[s.current+1:] {
		if l.ff&NotSelectable == NotSelectable {
			continue
		}
		if s.current >= 0 {
			(*s.c.ll)[s.current].ff &^= Highlighted
		}
		l.ff |= Highlighted
		s.current = s.current + 1 + idx
		break
	}
	if old == s.current {
		s.Reset()
		return s.current
	}
	return s.current
}

func (s *Highlighter) Previous() int {
	if s.current <= 0 {
		s.Reset()
		return s.current
	}
	old := s.current
	for i := s.current - 1; i >= 0; i-- {
		if (*s.c.ll)[i].ff&NotSelectable == NotSelectable {
			continue
		}
		if s.current >= 0 {
			(*s.c.ll)[s.current].ff &^= Highlighted
		}
		(*s.c.ll)[i].ff |= Highlighted
		s.current = i
		break
	}
	if old == s.current {
		s.Reset()
		return s.current
	}
	return s.current
}

func (s *Highlighter) Reset() {
	if s.current == -1 {
		return
	}
	(*s.c.ll)[s.current].ff &^= Highlighted
	s.current = -1
}
