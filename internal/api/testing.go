package api

import (
	"strings"
)

// Tester implementation augments an UIer implementation with additional
// functionality for testing.
type Tester interface {

	// String returns a string representation of the screen/window
	// content.
	Screen() StringScreen

	// StringArea returns a string representation of given screen/window
	// area.
	ScreenArea(x, y, width, height int) StringScreen

	// Cells returns the content of a test screen as lines of test
	// cells, i.e. in addition to a screens content also the style
	// information is provided.
	Cells() CellsScreen

	// CellsArea returns the content of a test screen area as lines of
	// test cells, i.e. in addition to the area's content also the style
	// information is provided.
	CellsArea(x, y, width, height int) CellsScreen

	// Display brings given string in given style to the screen.
	Display(string, Style)

	PostKey(Key, Modifier) error

	PostRune(rune, Modifier) error

	PostMouse(x, y int, _ Button, _ Modifier) error

	PostResize(width, height int) error
}

type StringLine string

func (l StringLine) isBlank() bool {
	for _, r := range l {
		if r != ' ' {
			return false
		}
	}
	return true
}

func (l StringLine) indentWidth() (int, int) {
	indent := 0
	for _, r := range l {
		if r == ' ' {
			indent++
			continue
		}
		break
	}
	if indent == len(l) {
		return len(l), 0
	}
	return indent, len(strings.TrimSpace(string(l)))
}

type StringScreen []string

func (ss StringScreen) String() string {
	return strings.Join(ss, "\n")
}

// Trimmed reduces given string to its minimum number of
// non-empty lines whereas the lines are trimmed to contain all non
// white space runes:
//
//	+--------------------+
//	|                    |       +------------+
//	|   upper left       |       |upper left  |
//	|                    |  =>   |            |
//	|          right     |       |       right|
//	|      bottom        |       |   bottom   |
//	|                    |       +------------+
//	+--------------------+
func (ss StringScreen) Trimmed() StringScreen {
	start, end := trimVertical(ss)
	if end == 0 {
		return StringScreen{}
	}
	vTrimmed := ss[start:end]
	start, end = trimHorizontal(vTrimmed)
	hTrimmed := StringScreen{}
	for _, s := range vTrimmed {
		hTrimmed = append(hTrimmed, s[start:end])
	}
	return hTrimmed
}

func (ss StringScreen) len() int { return len(ss) }

func (ss StringScreen) forLine(cb func(l liner) (stop bool)) {
	for _, l := range ss {
		if cb(StringLine(l)) {
			return
		}
	}
}

func (ss StringScreen) forReverse(cb func(l liner) (stop bool)) {
	for i := len(ss) - 1; i >= 0; i-- {
		if cb(StringLine(ss[i])) {
			return
		}
	}
}

type TestCell struct {
	Rune rune
	Sty  Style
}

// CellsLine represents a line of a [lines.TestScreen].
type CellsLine []TestCell

func (l CellsLine) isValidCell(x int) bool {
	return x >= 0 && x < len(l)
}

func (l CellsLine) HasBG(x int, c Color) bool {
	if !l.isValidCell(x) {
		return false
	}
	return l[x].Sty.BG == c
}

func (l CellsLine) HasFG(x int, c Color) bool {
	if !l.isValidCell(x) {
		return false
	}
	return l[x].Sty.FG == c
}

// HasAttr returns true if given style attribute mask is set in the
// style at given cell.
func (l CellsLine) HasAttr(x int, aa StyleAttribute) bool {
	if !l.isValidCell(x) {
		return false
	}
	return l[x].Sty.AA&aa == aa
}

func (l CellsLine) String() string {
	b := strings.Builder{}
	for _, c := range l {
		b.WriteRune(c.Rune)
	}
	return b.String()
}

func (l CellsLine) isBlank() bool {
	for _, c := range l {
		if c.Rune != ' ' {
			return false
		}
	}
	return true
}

func (l CellsLine) indentWidth() (int, int) {
	indent := 0
	for _, c := range l {
		if c.Rune == ' ' {
			indent++
			continue
		}
		break
	}
	if indent == len(l) {
		return len(l), 0
	}
	rightBlanks := 0
	for i := len(l) - 1; i >= 0; i-- {
		if l[i].Rune == ' ' {
			rightBlanks++
			continue
		}
		break
	}
	return indent, len(l) - indent - rightBlanks
}

type CellsScreen []CellsLine

func (cs CellsScreen) String() string {
	ll := []string{}
	for _, l := range cs {
		ll = append(ll, l.String())
	}
	return strings.Join(ll, "\n")
}

// isValidPos returns true if at given coordinates a screen cell may be
// looked up.
func (cs CellsScreen) isValidLine(y int) bool {
	return y >= 0 && y < len(cs)
}

func (cs CellsScreen) HasFG(x, y int, c Color) bool {
	if !cs.isValidLine(y) {
		return false
	}
	return cs[y].HasFG(x, c)
}

func (cs CellsScreen) HasBG(x, y int, c Color) bool {
	if !cs.isValidLine(y) {
		return false
	}
	return cs[y].HasBG(x, c)
}

func (cs CellsScreen) HasAttr(x, y int, aa StyleAttribute) bool {
	if !cs.isValidLine(y) {
		return false
	}
	return cs[y].HasAttr(x, aa)
}

// Trimmed reduces given screen-cells matrix to its minimum number of
// non-empty cells whereas the cells-lines are trimmed to contain all non
// white space cells:
//
//	+--------------------+
//	|                    |       +------------+
//	|   upper left       |       |upper left  |
//	|                    |  =>   |            |
//	|          right     |       |       right|
//	|      bottom        |       |   bottom   |
//	|                    |       +------------+
//	+--------------------+
func (cs CellsScreen) Trimmed() CellsScreen {
	start, end := trimVertical(cs)
	if end == 0 {
		return CellsScreen{}
	}
	vTrimmed := cs[start:end]
	start, end = trimHorizontal(vTrimmed)
	hTrimmed := CellsScreen{}
	for _, s := range vTrimmed {
		hTrimmed = append(hTrimmed, s[start:end])
	}
	return hTrimmed
}

func (cs CellsScreen) len() int { return len(cs) }

func (cs CellsScreen) forLine(cb func(l liner) (stop bool)) {
	for _, l := range cs {
		if cb(l) {
			return
		}
	}
}

func (cs CellsScreen) forReverse(cb func(l liner) (stop bool)) {
	for i := len(cs) - 1; i >= 0; i-- {
		if cb(cs[i]) {
			return
		}
	}
}

type liner interface {
	isBlank() bool
	// indentWidth returns a line's indent and width of non-blanks
	// i.e. for a line "  123  " 2 and 3 is returned.
	indentWidth() (int, int)
}

type screener interface {
	forLine(func(liner) (stop bool))
	forReverse(func(liner) (stop bool))
	len() int
}

func trimVertical(s screener) (int, int) {
	start, end := 0, s.len()
	if end == 0 {
		return 0, 0
	}
	s.forLine(func(l liner) (stop bool) {
		if !l.isBlank() {
			return true
		}
		start++
		return false
	})
	if start == s.len() {
		return 0, 0
	}
	s.forReverse(func(l liner) (stop bool) {
		if !l.isBlank() {
			return true
		}
		end--
		return false
	})
	return start, end
}

func trimHorizontal(s screener) (int, int) {
	x, width := -1, 0
	s.forLine(func(l liner) (stop bool) {
		indent, w := l.indentWidth()
		if indent < x || x == -1 {
			x = indent
		}
		if indent+w > width {
			width = indent + w
		}
		return
	})
	return x, width
}
