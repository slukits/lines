// Copyright (c) 2022 Stephan Lukits. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package lines

import "github.com/gdamore/tcell/v2"

// FmtMask types a control's respectively a control's line formattings.
type FmtMask uint

const (

	// centered center the next written lines.  TODO: implement
	centered FmtMask = 1 << iota

	// filled indicates to fill the whole line with a set background
	// color instead only the non-blank cells.
	filled
)

// llFmt represents a lines (ll) formatting (fmt) properties for a
// component, a line or a partial line.
type llFmt struct {
	sty  tcell.Style
	mask FmtMask
}

// FmtWriter instances provide an API for styling and formatting the
// writing to a component's line(s).
type FmtWriter struct {
	cmp cmpWriter
	fmt *llFmt
}

// FG sets the next write's foreground color.
func (w *FmtWriter) FG(color tcell.Color) *FmtWriter {
	w.fmt.sty = w.fmt.sty.Foreground(color)
	return w
}

// BG sets the next write's foreground color.
func (w *FmtWriter) BG(color tcell.Color) *BGWriter {
	w.fmt.sty = w.fmt.sty.Background(color)
	return &BGWriter{cmp: w.cmp, fmt: w.fmt}
}

// Attr sets the next write's style attributes like bold.
func (w *FmtWriter) Attr(aa tcell.AttrMask) *FmtWriter {
	w.fmt.sty = w.fmt.sty.Attributes(aa)
	return w
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
	return &locWriter{at: idx, ff: _ff, cmp: w.cmp, fmt: w.fmt}
}

// Write to a components screen-portion made available by an Env
// instance provided to a listener implementation.
func (w *FmtWriter) Write(bb []byte) (int, error) {
	return w.cmp.write(bb, 0, 0, w.fmt)
}

// BGWriter instances provide an API for styling and formatting the
// writing to a component's line(s) like FmtWriter.  But it provides one
// additional method [BGWriter.Filled] indicating that the whole line(s)
// should have set background color and not only the part were content
// is written to.
type BGWriter struct {
	cmp cmpWriter
	fmt *llFmt
}

// FG sets the next write's foreground color.
func (w *BGWriter) FG(color tcell.Color) *BGWriter {
	w.fmt.sty = w.fmt.sty.Foreground(color)
	return w
}

// BG sets the next write's background color.
func (w *BGWriter) BG(color tcell.Color) *BGWriter {
	w.fmt.sty = w.fmt.sty.Background(color)
	return w
}

// Attr sets the next write's style attributes like bold.
func (w *BGWriter) Attr(aa tcell.AttrMask) *BGWriter {
	w.fmt.sty = w.fmt.sty.Attributes(aa)
	return w
}

// Fmt sets the next write's formattings like centered.
// func (w *BGWriter) Fmt(ff FmtMask) *BGWriter {
// 	internal := w.fmt.mask & filled // | other future internals
// 	w.fmt.mask = ff | internal
// 	return w
// }

// Filled ensures that the whole line has set background color not only
// the cells which is written to.  E.g.:
//
//	fmt.Fprint(e.BG(tcell.ColorRed), "with red background")
//
// only "with red background" will have a red background while the rest
// of the line has the components background color.
//
//	fmt.Fprint(e.BG(tcell.ColorRed).Filled(), "with red background")
//
// the whole line will have a red background color.
func (w *BGWriter) Filled() *BGWriter {
	w.fmt.mask |= filled
	return w
}

// LL returns a writer which writes to the line and its following lines
// at given index.
func (w *BGWriter) LL(idx int, ff ...LineFlags) *locWriter {
	_ff := LineFlags(0)
	for _, f := range ff {
		_ff |= f
	}
	return &locWriter{at: idx, ff: _ff, cmp: w.cmp, fmt: w.fmt}
}

func (w *BGWriter) Write(bb []byte) (int, error) {
	return w.cmp.write(bb, 0, 0, w.fmt)
}

// locWriter represents a location writer implementing the writer
// interface to write a sequence of lines at a specific location of a
// component which either starts at a given line, at a given line's cell
// or appends at the end.
type locWriter struct {
	fmt *llFmt
	at  int
	ff  LineFlags
	cmp cmpWriter
}

// Write to a specific line an onward.
func (w *locWriter) Write(bb []byte) (int, error) {
	return w.cmp.write(bb, w.at, w.ff, w.fmt)
}
