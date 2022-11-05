// Copyright (c) 2022 Stephan Lukits. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package lines

// EnvWriter instances provide an API for styling and formatting the
// writing to a component's line(s) starting at its first line.
type EnvWriter struct {
	cmp cmpWriter
	sty *Style
}

// FG sets the next write's foreground color.
func (w *EnvWriter) FG(color Color) *EnvWriter {
	var sty Style
	if w.sty == nil {
		sty = w.cmp.globals().Style(Default).WithFG(color)
	} else {
		sty = w.sty.WithFG(color)
	}
	w.sty = &sty
	return w
}

// BG sets the next write's foreground color.
func (w *EnvWriter) BG(color Color) *EnvWriter {
	var sty Style
	if w.sty == nil {
		sty = w.cmp.globals().Style(Default).WithBG(color)
	} else {
		sty = w.sty.WithBG(color)
	}
	w.sty = &sty
	return w
}

// AA sets the next write's style attributes like bold.
func (w *EnvWriter) AA(aa StyleAttributeMask) *EnvWriter {
	var sty Style
	if w.sty == nil {
		sty = w.cmp.globals().Style(Default).WithAA(aa)
	} else {
		sty = w.sty.WithAA(aa)
	}
	w.sty = &sty
	return w
}

// LL returns a writer which writes to the line and its following lines
// at given index idx.
func (w *EnvWriter) LL(idx int) *EnvLineWriter {
	return &EnvLineWriter{line: idx, cmp: w.cmp, sty: w.sty}
}

// Write to a components screen-portion made available by an Env
// instance provided to a listener implementation.
func (w *EnvWriter) Write(bb []byte) (int, error) {
	return w.cmp.write(bb, 0, -1, w.sty)
}

// EnvLineWriter instances provide an API for styling and formatting the
// writing to a component's n-th line(s).
type EnvLineWriter struct {
	sty  *Style
	line int
	cmp  cmpWriter
}

// FG sets the next write's foreground color.
func (w *EnvLineWriter) FG(color Color) *EnvLineWriter {
	var sty Style
	if w.sty == nil {
		sty = w.cmp.globals().Style(Default).WithFG(color)
	} else {
		sty = w.sty.WithFG(color)
	}
	w.sty = &sty
	return w
}

// BG sets the next write's foreground color.
func (w *EnvLineWriter) BG(color Color) *EnvLineWriter {
	var sty Style
	if w.sty == nil {
		sty = w.cmp.globals().Style(Default).WithBG(color)
	} else {
		sty = w.sty.WithBG(color)
	}
	w.sty = &sty
	return w
}

// AA sets the next write's style attributes like bold.
func (w *EnvLineWriter) AA(aa StyleAttributeMask) *EnvLineWriter {
	var sty Style
	if w.sty == nil {
		sty = w.cmp.globals().Style(Default).WithAA(aa)
	} else {
		sty = w.sty.WithAA(aa)
	}
	w.sty = &sty
	return w
}

// Write to a specific line an onward.
func (w *EnvLineWriter) Write(bb []byte) (int, error) {
	return w.cmp.write(bb, w.line, -1, w.sty)
}

// At returns a writer which writes at given line writer w's line at
// given cell.  Note you need to use the [lines.Print]-function to write
// to an at-writer and can only provide a rune or a rune-slice.  Styles
// of an at-writer are only applied for the printed range of runes.
func (w *EnvLineWriter) At(cell int) *EnvAtWriter {
	return &EnvAtWriter{line: w.line, cell: cell, cmp: w.cmp, sty: w.sty}
}

type EnvAtWriter struct {
	sty        *Style
	line, cell int
	cmp        cmpWriter
}

// FG sets the next write's foreground color.
func (w *EnvAtWriter) FG(color Color) *EnvAtWriter {
	var sty Style
	if w.sty == nil {
		sty = w.cmp.globals().Style(Default).WithFG(color)
	} else {
		sty = w.sty.WithFG(color)
	}
	w.sty = &sty
	return w
}

// BG sets the next write's foreground color.
func (w *EnvAtWriter) BG(color Color) *EnvAtWriter {
	var sty Style
	if w.sty == nil {
		sty = w.cmp.globals().Style(Default).WithBG(color)
	} else {
		sty = w.sty.WithBG(color)
	}
	w.sty = &sty
	return w
}

// AA sets the next write's style attributes like bold.
func (w *EnvAtWriter) AA(aa StyleAttributeMask) *EnvAtWriter {
	var sty Style
	if w.sty == nil {
		sty = w.cmp.globals().Style(Default).WithAA(aa)
	} else {
		sty = w.sty.WithAA(aa)
	}
	w.sty = &sty
	return w
}

// Filling returns a filling writer which adds first printed rune to
// provided line-index at provided cell as a filling rune.
func (w *EnvAtWriter) Filling() *envAtFillingWriter {
	return &envAtFillingWriter{
		sty:  w.sty,
		line: w.line,
		cell: w.cell,
		cmp:  w.cmp,
	}
}

// WriteAt writes given runes rr to provided line and cell indices with
// set style information.  If there is style information it will be only
// applied for given rune sequence rr.
func (w *EnvAtWriter) WriteAt(rr []rune) {
	if len(rr) == 0 {
		return
	}
	w.cmp.writeAt(rr, w.line, w.cell, w.sty)
}

type envAtFillingWriter struct {
	sty        *Style
	line, cell int
	cmp        cmpWriter
}

func (w *envAtFillingWriter) WriteAt(rr []rune) {
	if len(rr) == 0 {
		return
	}
	w.cmp.writeAtFilling(rr[0], w.line, w.cell, w.sty)
}
