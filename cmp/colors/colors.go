// Copyright (c) 2022 Stephan Lukits. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

/*
Package colors provides helper functions for dealing with colors in
context of lines.  I.e. it its mainly about coloring text.  This effort
cumulates in a color-scheme picker.
*/
package colors

import "github.com/slukits/lines"

type Mono int32

const (
	Black Mono = Mono(lines.Black)
	White Mono = Mono(lines.White)
)

var MonoColors = []Mono{Black, White}

func MonoForeground(bg Mono) []lines.Style {
	if bg == Black {
		return []lines.Style{lines.NewStyle(
			lines.ZeroStyle, lines.White, lines.Black)}
	}
	return []lines.Style{lines.NewStyle(
		lines.ZeroStyle, lines.Black, lines.White)}
}

func MonoBackground(fg Mono) []lines.Style {
	if fg == Black {
		return []lines.Style{lines.NewStyle(
			lines.ZeroStyle, lines.Black, lines.White)}
	}
	return []lines.Style{lines.NewStyle(
		lines.ZeroStyle, lines.White, lines.Black)}
}

type Eight int32

const ()
