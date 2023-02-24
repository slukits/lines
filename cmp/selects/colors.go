// Copyright (c) 2022 Stephan Lukits. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

/*
Package colors provides helper functions for dealing with colors in
context of lines.  I.e. it its mainly about coloring text.  This effort
cumulates in a color-scheme picker.
*/
package selects

import "github.com/slukits/lines"

type ColorRange uint8

const (
	Monochrome ColorRange = iota
	System8
	System8Linux
	System16Colors
	ANSIColors
	// TrueColor
)

var RangeNames = map[ColorRange]string{
	Monochrome:     "Monochrome",
	System8:        "System8",
	System8Linux:   "System8-Linux",
	System16Colors: "System16",
	ANSIColors:     "ANSI",
	// TrueColor:          "true-color",
}

// Mono types the colors of a monochrome display.
type Mono int32

const (
	BlackM Mono = Mono(lines.Black)
	WhiteM Mono = Mono(lines.White)
)

var MonoColors = []Mono{BlackM, WhiteM}

var monoColors = map[lines.Color]bool{
	lines.Black: true,
	lines.White: true,
}

// MonoForeground provides the possible foreground color to given
// background color bg.
func MonoForeground(bg Mono) []lines.Style {
	if bg == BlackM {
		return []lines.Style{lines.NewStyle(
			lines.ZeroStyle, lines.White, lines.Black)}
	}
	return []lines.Style{lines.NewStyle(
		lines.ZeroStyle, lines.Black, lines.White)}
}

// MonoBackground provides the possible background color to given
// foreground color fg.
func MonoBackground(fg Mono) []lines.Style {
	if fg == BlackM {
		return []lines.Style{lines.NewStyle(
			lines.ZeroStyle, lines.Black, lines.White)}
	}
	return []lines.Style{lines.NewStyle(
		lines.ZeroStyle, lines.White, lines.Black)}
}

// System8Colors types the eight system colors.
type Colors []lines.Color

func (cc Colors) Have(c lines.Color) bool {
	for _, c_ := range cc {
		if c_ != c {
			continue
		}
		return true
	}
	return false
}

var system8Colors = Colors{
	lines.Black, lines.Maroon, lines.Green, lines.Olive, lines.Navy,
	lines.Purple, lines.Teal, lines.Silver,
}

// System8Foregrounds provides the possible foreground combinations with
// given background color bg.  A possible foreground combination is any
// System8 color which is not bg.
func System8Foregrounds(bg lines.Color) (ss []lines.Style) {
	if !system8Colors.Have(bg) {
		return ss
	}
	for _, c := range system8Colors {
		if c == bg {
			continue
		}
		sty := lines.NewStyle(lines.ZeroStyle, c, bg)
		ss = append(ss, sty)
	}
	return ss
}

// System8Backgrounds provides the possible background combinations with
// given foreground color fg.  A possible background combination is any
// System8 color which is not fg.
func System8Backgrounds(fg lines.Color) (ss []lines.Style) {
	if !system8Colors.Have(fg) {
		return ss
	}
	for _, c := range system8Colors {
		if c == fg {
			continue
		}
		sty := lines.NewStyle(lines.ZeroStyle, fg, c)
		ss = append(ss, sty)
	}
	return ss
}

var linuxBGColors = Colors{
	lines.Black, lines.Maroon, lines.Green, lines.Olive, lines.Navy,
	lines.Purple, lines.Teal, lines.Silver,
}

var linuxFGColors = Colors{
	lines.Black, lines.Maroon, lines.Green, lines.Olive, lines.Navy,
	lines.Purple, lines.Teal, lines.Silver, lines.Grey, lines.Red,
	lines.Lime, lines.Yellow, lines.Blue, lines.Fuchsia, lines.Aqua,
	lines.White,
}

var linuxFGBold = map[lines.Color]lines.Color{
	lines.Grey:    lines.Black,
	lines.Red:     lines.Maroon,
	lines.Lime:    lines.Green,
	lines.Yellow:  lines.Olive,
	lines.Blue:    lines.Navy,
	lines.Fuchsia: lines.Purple,
	lines.Aqua:    lines.Teal,
	lines.White:   lines.Silver,
}

var linuxBoldFG = map[lines.Color]lines.Color{
	lines.Black:  lines.Grey,
	lines.Maroon: lines.Red,
	lines.Green:  lines.Lime,
	lines.Olive:  lines.Yellow,
	lines.Navy:   lines.Blue,
	lines.Purple: lines.Fuchsia,
	lines.Teal:   lines.Aqua,
	lines.Silver: lines.White,
}

// LinuxBackgrounds provides the possible background combinations with
// given foreground color fg.  A possible background combination is any
// LinuxBG color which is not fg.  Note since we have more foreground
// colors than background colors we may get seven or eight combination
// depending if fg is a background color or not.
func LinuxBackgrounds(fg lines.Color) (ss []lines.Style) {
	if !linuxFGColors.Have(fg) {
		return
	}
	attr := lines.ZeroStyle
	fg_ := fg
	if _, ok := linuxFGBold[fg]; ok {
		attr = lines.Bold
		fg_ = linuxFGBold[fg]
	}
	for _, c := range linuxBGColors {
		if c == fg {
			continue
		}
		ss = append(ss, lines.NewStyle(attr, fg_, c))
	}
	return ss
}

// LinuxForegrounds provides the possible foreground combinations with
// given background color bg.  A possible foreground combination is any
// LinuxFG color which is not bg.  Note eight of the linux foreground
// colors are the corresponding background colors (e.g. GreyFG ->
// BlackBG) having the style attribute set to bold.
func LinuxForegrounds(bg lines.Color) (ss []lines.Style) {
	if !linuxBGColors.Have(bg) {
		return
	}
	for _, c := range linuxFGColors {
		if c == bg {
			continue
		}
		var sty lines.Style
		if fg_, ok := linuxFGBold[c]; ok {
			sty = lines.NewStyle(lines.Bold, fg_, bg)
		} else {
			sty = lines.NewStyle(lines.ZeroStyle, c, bg)
		}
		ss = append(ss, sty)
	}
	return ss
}

// System types the sixteen terminal system colors.
type System int32

const (
	Black16   System = System(lines.Black)
	Maroon16  System = System(lines.Maroon)
	Green16   System = System(lines.Green)
	Olive16   System = System(lines.Olive)
	Navy16    System = System(lines.Navy)
	Purple16  System = System(lines.Purple)
	Teal16    System = System(lines.Teal)
	Silver16  System = System(lines.Silver)
	Grey16    System = System(lines.Grey)
	Red16     System = System(lines.Red)
	Lime16    System = System(lines.Lime)
	Yellow16  System = System(lines.Yellow)
	Blue16    System = System(lines.Blue)
	Fuchsia16 System = System(lines.Fuchsia)
	Aqua16    System = System(lines.Aqua)
	White16   System = System(lines.White)
)

var system16Colors = Colors{
	lines.Black, lines.Maroon, lines.Green, lines.Olive, lines.Navy,
	lines.Purple, lines.Teal, lines.Silver, lines.Grey, lines.Red,
	lines.Lime, lines.Yellow, lines.Blue, lines.Fuchsia, lines.Aqua,
	lines.White,
}

// System16Foregrounds provides the possible foreground combinations with
// given background color bg.  A possible foreground combination is any
// System color which is not bg.
func System16Foregrounds(bg lines.Color) (ss []lines.Style) {
	if !system16Colors.Have(bg) {
		return ss
	}
	for _, c := range system16Colors {
		if c == bg {
			continue
		}
		sty := lines.NewStyle(lines.ZeroStyle, c, bg)
		ss = append(ss, sty)
	}
	return ss
}

// System16Backgrounds provides the possible background combinations with
// given foreground color fg.  A possible background combination is any
// System color which is not fg.
func System16Backgrounds(fg lines.Color) (ss []lines.Style) {
	if !system16Colors.Have(fg) {
		return ss
	}
	for _, c := range system16Colors {
		if c == fg {
			continue
		}
		sty := lines.NewStyle(lines.ZeroStyle, fg, c)
		ss = append(ss, sty)
	}
	return ss
}

var ansiColors = Colors{
	lines.Black, lines.Maroon, lines.Green, lines.Olive, lines.Navy,
	lines.Purple, lines.Teal, lines.Silver, lines.Grey, lines.Red,
	lines.Lime, lines.Yellow, lines.Blue, lines.Fuchsia, lines.Aqua,
	lines.White, lines.Grey0, lines.NavyBlue, lines.DarkBlue,
	lines.Blue3, lines.Blue3_2, lines.Blue1, lines.DarkGreen,
	lines.DeepSkyBlue4, lines.DeepSkyBlue4_2, lines.DeepSkyBlue4_3,
	lines.DodgerBlue3, lines.DodgerBlue2, lines.Green4,
	lines.SpringGreen4, lines.Turquoise4, lines.DeepSkyBlue3,
	lines.DeepSkyBlue3_2, lines.DodgerBlue1, lines.Green3,
	lines.SpringGreen3, lines.DarkCyan, lines.LightSeaGreen,
	lines.DeepSkyBlue2, lines.DeepSkyBlue1, lines.Green3_2,
	lines.SpringGreen3_2, lines.SpringGreen2, lines.Cyan3,
	lines.DarkTurquoise, lines.Turquoise2, lines.Green1,
	lines.SpringGreen2_2, lines.SpringGreen1, lines.MediumSpringGreen,
	lines.Cyan2, lines.Cyan1, lines.DarkRed, lines.DeepPink4_3,
	lines.Purple4, lines.Purple4_2, lines.Purple3, lines.BlueViolet,
	lines.Orange4, lines.Grey37, lines.MediumPurple4, lines.SlateBlue3,
	lines.SlateBlue3_2, lines.RoyalBlue1, lines.Chartreuse4,
	lines.DarkSeaGreen4, lines.PaleTurquoise4, lines.SteelBlue,
	lines.SteelBlue3, lines.CornflowerBlue, lines.Chartreuse3,
	lines.DarkSeaGreen4_2, lines.CadetBlue, lines.CadetBlue_2,
	lines.SkyBlue3, lines.SteelBlue1, lines.Chartreuse3_2,
	lines.PaleGreen3, lines.SeaGreen3, lines.Aquamarine3,
	lines.MediumTurquoise, lines.SteelBlue1_2, lines.Chartreuse2,
	lines.SeaGreen2, lines.SeaGreen1, lines.SeaGreen1_2,
	lines.Aquamarine1, lines.DarkSlateGray2, lines.DarkRed_2,
	lines.DeepPink4_2, lines.DarkMagenta, lines.DarkMagenta_2,
	lines.DarkViolet, lines.Purple_1, lines.Orange4_2, lines.LightPink4,
	lines.Plum4, lines.MediumPurple3, lines.MediumPurple3_2,
	lines.SlateBlue1, lines.Yellow4, lines.Wheat4, lines.Grey53,
	lines.LightSlateGrey, lines.MediumPurple, lines.LightSlateBlue,
	lines.Yellow4_2, lines.DarkOliveGreen3, lines.DarkSeaGreen,
	lines.LightSkyBlue3, lines.LightSkyBlue3_2, lines.SkyBlue2,
	lines.Chartreuse2_2, lines.DarkOliveGreen3_2, lines.PaleGreen3_2,
	lines.DarkSeaGreen3, lines.DarkSlateGray3, lines.SkyBlue1,
	lines.Chartreuse1, lines.LightGreen, lines.LightGreen_2,
	lines.PaleGreen1, lines.Aquamarine1_2, lines.DarkSlateGray1,
	lines.Red3, lines.DeepPink4, lines.MediumVioletRed, lines.Magenta3,
	lines.DarkViolet_2, lines.Purple_2, lines.DarkOrange3,
	lines.IndianRed, lines.HotPink3, lines.MediumOrchid3,
	lines.MediumOrchid, lines.MediumPurple2, lines.DarkGoldenrod,
	lines.LightSalmon3, lines.RosyBrown, lines.Grey63,
	lines.MediumPurple2_2, lines.MediumPurple1, lines.Gold3,
	lines.DarkKhaki, lines.NavajoWhite3, lines.Grey69,
	lines.LightSteelBlue3, lines.LightSteelBlue, lines.Yellow3,
	lines.DarkOliveGreen3_3, lines.DarkSeaGreen3_2, lines.DarkSeaGreen2,
	lines.LightCyan3, lines.LightSkyBlue1, lines.GreenYellow,
	lines.DarkOliveGreen2, lines.PaleGreen1_2, lines.DarkSeaGreen2_2,
	lines.DarkSeaGreen1, lines.PaleTurquoise1, lines.Red3_2,
	lines.DeepPink3, lines.DeepPink3_2, lines.Magenta3_2,
	lines.Magenta3_3, lines.Magenta2, lines.DarkOrange3_2,
	lines.IndianRed_2, lines.HotPink3_2, lines.HotPink2, lines.Orchid,
	lines.MediumOrchid1, lines.Orange3, lines.LightSalmon3_2,
	lines.LightPink3, lines.Pink3, lines.Plum3, lines.Violet,
	lines.Gold3_2, lines.LightGoldenrod3, lines.Tan, lines.MistyRose3,
	lines.Thistle3, lines.Plum2, lines.Yellow3_2, lines.Khaki3,
	lines.LightGoldenrod2, lines.LightYellow3, lines.Grey84,
	lines.LightSteelBlue1, lines.Yellow2, lines.DarkOliveGreen1,
	lines.DarkOliveGreen1_2, lines.DarkSeaGreen1_2, lines.Honeydew2,
	lines.LightCyan1, lines.Red1, lines.DeepPink2, lines.DeepPink1,
	lines.DeepPink1_2, lines.Magenta2_2, lines.Magenta1,
	lines.OrangeRed1, lines.IndianRed1, lines.IndianRed1_2,
	lines.HotPink, lines.HotPink_2, lines.MediumOrchid1_2,
	lines.DarkOrange, lines.Salmon1, lines.LightCoral,
	lines.PaleVioletRed1, lines.Orchid2, lines.Orchid1, lines.Orange1,
	lines.SandyBrown, lines.LightSalmon1, lines.LightPink1, lines.Pink1,
	lines.Plum1, lines.Gold1, lines.LightGoldenrod2_2,
	lines.LightGoldenrod2_3, lines.NavajoWhite1, lines.MistyRose1,
	lines.Thistle1, lines.Yellow1, lines.LightGoldenrod1, lines.Khaki1,
	lines.Wheat1, lines.Cornsilk1, lines.Grey100, lines.Grey3,
	lines.Grey7, lines.Grey11, lines.Grey15, lines.Grey19, lines.Grey23,
	lines.Grey27, lines.Grey30, lines.Grey35, lines.Grey39,
	lines.Grey42, lines.Grey46, lines.Grey50, lines.Grey54,
	lines.Grey58, lines.Grey62, lines.Grey66, lines.Grey70,
	lines.Grey74, lines.Grey78, lines.Grey82, lines.Grey85,
	lines.Grey89, lines.Grey93,
}

// ANSIForegrounds provides the possible foreground combinations with
// given background color bg.  A possible foreground combination is any
// ANSI color which is not bg.
func ANSIForegrounds(bg lines.Color) (ss []lines.Style) {
	if !ansiColors.Have(bg) {
		return
	}
	for _, c := range ansiColors {
		if c == bg {
			continue
		}
		sty := lines.NewStyle(lines.ZeroStyle, c, bg)
		ss = append(ss, sty)
	}
	return ss
}

// ANSIBackgrounds provides the possible background combinations with
// given foreground color fg.  A possible background combination is any
// ANSI color which is not fg.
func ANSIBackgrounds(fg lines.Color) (ss []lines.Style) {
	if !ansiColors.Have(fg) {
		return
	}
	for _, c := range ansiColors {
		if c == fg {
			continue
		}
		sty := lines.NewStyle(lines.ZeroStyle, fg, c)
		ss = append(ss, sty)
	}
	return ss
}
