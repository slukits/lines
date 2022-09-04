package lines

// Scroller provides a component's scrolling API.
type Scroller struct{ c *Component }

func (s Scroller) IsAtTop() bool { return s.c.first == 0 }

func (s Scroller) IsAtBottom() bool {
	_, _, _, height := s.c.dim.Area()
	return s.c.first+height >= s.c.Len()
}

// Up scrolls one page up or to the last content line is the last
// displayed line.  Whereas "one page" is in case of a component
// height of 1 is one line.  For a height h with 1 < h < 20 "one page"
// is h - 1.  For h >= 20 "one page" is h - h/10.
func (s Scroller) Up() {
	if s.c.first == 0 {
		return
	}
	_, _, _, height := s.c.dim.Area()
	if height == 1 {
		s.c.setFirst(s.c.first - 1)
		return
	}
	scroll := height - 1
	if height >= 20 {
		scroll = height - (height / 10)
	}
	if scroll >= s.c.first {
		s.c.setFirst(0)
		return
	}
	s.c.setFirst(s.c.first - scroll)
}

// ToTop scrolls component content to its first line, i.e. the first
// displayed line is the first content line.
func (s Scroller) ToTop() {
	s.c.setFirst(0)
}

func (s Scroller) ToBottom() {
	_, _, _, height := s.c.dim.Area()
	s.c.setFirst(s.c.Len() - height)
}

// Down scrolls one page down or to the first content line is the first
// displayed line.  Whereas "one page" is in case of a component height
// of 1 is one line.  For a height h with 1 < h < 20 "one page" is h -
// 1.  For h >= 20 "one page" is h - h/10.
func (s Scroller) Down() {
	_, _, _, height := s.c.dim.Area()
	if height >= len((*s.c.ll)[s.c.first:]) {
		return
	}
	if height == 1 {
		s.c.setFirst(s.c.first + 1)
		return
	}
	scroll := height - 1
	if height >= 20 {
		scroll = height - (height / 10)
	}
	if s.c.Len()-(s.c.first+scroll) < height {
		s.c.setFirst(s.c.Len() - height)
		return
	}
	s.c.setFirst(s.c.first + scroll)
}