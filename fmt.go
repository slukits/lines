// Copyright (c) 2022 Stephan Lukits. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package lines

// FmtMask types a control's respectively a control's line formattings.
type FmtMask uint

const (

	// centered center the next written lines.  TODO: implement
	centered FmtMask = 1 << iota

	// filled indicates to fill the whole line with a set background
	// color instead only the non-blank cells.
	filled

	// onetimeFilled like filled but the fmt-flag is removed once
	// executed.  E.g. remove highlight.
	onetimeFilled
)

// llFmt represents a lines (ll) formatting (fmt) properties for a
// component, a line or a partial line.
type llFmt struct {
	sty  Style
	mask FmtMask
}

// FmtWriter instances provide an API for styling and formatting the
// writing to a component's line(s).
type FmtWriter struct {
	cmp cmpWriter
	sty Style
}

// FG sets the next write's foreground color.
func (w *FmtWriter) FG(color Color) *FmtWriter {
	w.sty.FG = color
	return w
}

// BG sets the next write's foreground color.
func (w *FmtWriter) BG(color Color) *FmtWriter {
	w.sty.BG = color
	return w
}

// Attr sets the next write's style attributes like bold.
func (w *FmtWriter) Attr(aa StyleAttribute) *FmtWriter {
	w.sty.AA = aa
	return w
}

func (w *FmtWriter) get(line int) *line {
	return (*w.cmp.(Componenter).embedded().ll)[line]
}

func (w *FmtWriter) has(line int) bool {
	if line < 0 || line >= len(*w.cmp.(Componenter).embedded().ll) {
		return false
	}
	return true
}

// Fmt sets the next write's formattings like centered.
// func (w *FmtWriter) Fmt(ff FmtMask) *FmtWriter {
// 	internal := w.fmt.mask & filled // | other future internals
// 	w.fmt.mask = ff | internal
// 	return w
// }

// LL returns a writer which writes to the line and its following lines
// at given index.
func (w *FmtWriter) LL(idx int, ff ...LineFlags) *locWriter {
	_ff := LineFlags(0)
	for _, f := range ff {
		_ff |= f
	}
	return &locWriter{line: idx, cell: -1, ff: _ff, cmp: w.cmp, sty: w.sty}
}

// At sets the collected style attributes and given flags for provided
// range at given line.
func (w *FmtWriter) At(line, cell int, ff ...LineFlags) *locWriter {
	_ff := LineFlags(0)
	for _, f := range ff {
		_ff |= f
	}
	return &locWriter{line: line, cell: cell, ff: _ff, cmp: w.cmp, sty: w.sty}
}

// Write to a components screen-portion made available by an Env
// instance provided to a listener implementation.
func (w *FmtWriter) Write(bb []byte) (int, error) {
	return w.cmp.write(bb, 0, -1, 0, w.sty)
}

// locWriter represents a location writer implementing the writer
// interface to write a sequence of lines at a specific location of a
// component which either starts at a given line, at a given line's cell
// or appends at the end.
type locWriter struct {
	sty        Style
	line, cell int
	ff         LineFlags
	cmp        cmpWriter
}

// Write to a specific line an onward.
func (w *locWriter) Write(bb []byte) (int, error) {
	return w.cmp.write(bb, w.line, w.cell, w.ff, w.sty)
}
