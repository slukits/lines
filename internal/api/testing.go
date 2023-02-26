package api

import (
	"strings"
)

// Tester implementation augments an UIer implementation with additional
// functionality for testing.
type Tester interface {

	// Size returns the number of available lines (height) and the number of
	// runes per line (width) as reported by an backend.
	Size() (width, height int)

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

	// PostKey emulates a user-key-input with underlying backend.
	PostKey(Key, ModifierMask) error

	// PostKey emulates a user-rune-input with underlying backend.
	PostRune(rune, ModifierMask) error

	// PostKey emulates a user-mouse-input with underlying backend.
	PostMouse(x, y int, _ ButtonMask, _ ModifierMask) error

	// PostKey emulates a resize event of the available display area
	// with underlying backend.
	PostResize(width, height int) error
}

// StringLine is a line of a [StringScreen] providing the sequence of
// runes displayed in a particular screen line.
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
	return indent, len([]rune(strings.TrimSpace(string(l))))
}

// StringScreen is the string representation of the screen lines at a
// particular point in time.  Note use [StringScreen.Trimmed] to
// minimize reported screen content.
type StringScreen []string

// String joins the lines of given screen string representation with
// line breaks and returns resulting string.
func (ss StringScreen) String() string {
	return strings.Join(ss, "\n")
}

// Column returns the content of the column with given index as string.
func (ss StringScreen) Column(idx int) string {
	if idx < 0 || len(ss) == 0 || idx >= len([]rune(ss[0])) {
		return ""
	}
	rr := []rune{}
	for _, s := range ss {
		rr = append(rr, []rune(s)[idx])
	}
	return string(rr)
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
	Rune  rune
	Style Style
}

// CellsLine represents a line of a [CellsScreen] providing of each cell
// in the line its displayed rune and style information for
// test-evaluations.
type CellsLine []TestCell

func (l CellsLine) isValidCell(x int) bool {
	return x >= 0 && x < len(l)
}

// HasBG returns true if line cell at given position x in given line
// cells l has given background color.
func (l CellsLine) HasBG(x int, c Color) bool {
	if !l.isValidCell(x) {
		return false
	}
	return l[x].Style.BG() == c
}

// HasFG returns true if line cell at given position x in given line
// cells l has given foreground color.
func (l CellsLine) HasFG(x int, c Color) bool {
	if !l.isValidCell(x) {
		return false
	}
	return l[x].Style.FG() == c
}

// HasAA returns true if line cell at given position x in given line
// cells l has given foreground color.
func (l CellsLine) HasAA(x int, aa StyleAttributeMask) bool {
	if !l.isValidCell(x) {
		return false
	}
	return l[x].Style.AA()&aa == aa
}

// String returns a string representation of given line cells l.
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

// CellsScreen is a screen representation at a specific point in time of
// [CellsLine] instances.  NOTE use [CellsScreen.Trimmed] to minimize the
// reported screen area.
type CellsScreen []CellsLine

// Column returns a column of test-cells from the test screen or nil iff
// the cells-screen cs is zero or given column i doesnt exist.
func (cs CellsScreen) Column(i int) []TestCell {
	if len(cs) == 0 || len(cs[0]) <= i || i < 0 {
		return nil
	}
	c := make([]TestCell, len(cs))
	for j, l := range cs {
		c[j] = l[i]
	}
	return c
}

// String returns a string representation of given screen cells cs.
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

// HasFG returns true if the screen cell at given coordinates x and y in
// given screen cells cs have given foreground color c.
func (cs CellsScreen) HasFG(x, y int, c Color) bool {
	if !cs.isValidLine(y) {
		return false
	}
	return cs[y].HasFG(x, c)
}

// HasBG returns true if the screen cell at given coordinates x and y in
// given screen cells cs have given background color c.
func (cs CellsScreen) HasBG(x, y int, c Color) bool {
	if !cs.isValidLine(y) {
		return false
	}
	return cs[y].HasBG(x, c)
}

// HasAA returns true if the screen cell at given coordinates x and y in
// given screen cells cs have given style attributes aa.
func (cs CellsScreen) HasAA(x, y int, aa StyleAttributeMask) bool {
	if !cs.isValidLine(y) {
		return false
	}
	return cs[y].HasAA(x, aa)
}

/*
Trimmed reduces given screen-cells matrix to its minimum number of
non-empty cells whereas the cells-lines are trimmed to contain all non
white space cells:

	+--------------------+
	|                    |       +------------+
	|   upper left       |       |upper left  |
	|                    |  =>   |            |
	|          right     |       |       right|
	|      bottom        |       |   bottom   |
	|                    |       +------------+
	+--------------------+
*/
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

func (cs CellsScreen) Equals(other CellsScreen) bool {
	for i, l := range cs {
		for j, c := range l {
			if len(other) <= i {
				return false
			}
			if len(other[i]) <= j {
				return false
			}
			oc := other[i][j]
			if oc.Style != c.Style || oc.Rune != c.Rune {
				return false
			}
		}
	}
	return true
}

type CellScreenDiff struct {
	LinesCount int
	CellsCount int
	Line, Cell int
}

func (cs CellsScreen) FirstDiff(other CellsScreen) CellScreenDiff {
	d := CellScreenDiff{
		LinesCount: -1, CellsCount: -1, Line: -1, Cell: -1}
	for i, l := range cs {
		for j, c := range l {
			if len(other) <= i {
				d.LinesCount = i
				return d
			}
			if len(other[i]) <= j {
				d.LinesCount, d.CellsCount = i, j
				return d
			}
			oc := other[i][j]
			if oc.Style != c.Style || oc.Rune != c.Rune {
				d.Line, d.Cell = i, j
				return d
			}
		}
	}
	return d
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
