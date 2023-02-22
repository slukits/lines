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
	System8Colors
	System8LinuxColors
	System16Colors
	ANSIColors
	TrueTypeColor
)

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

// System8 types the eight system colors.
type System8 int32

const (
	Black8  System8 = System8(lines.Black)
	Maroon8 System8 = System8(lines.Maroon)
	Green8  System8 = System8(lines.Green)
	Olive8  System8 = System8(lines.Olive)
	Navy8   System8 = System8(lines.Navy)
	Purple8 System8 = System8(lines.Purple)
	Teal8   System8 = System8(lines.Teal)
	Silver8 System8 = System8(lines.Silver)
)

var system8Colors = []System8{
	Black8, Maroon8, Green8, Olive8, Navy8, Purple8, Teal8, Silver8}

// System8Foregrounds provides the possible foreground combinations with
// given background color bg.  A possible foreground combination is any
// System8 color which is not bg.
func System8Foregrounds(bg System8) (ss []lines.Style) {
	for _, c := range system8Colors {
		if c == bg {
			continue
		}
		sty := lines.NewStyle(
			lines.ZeroStyle, lines.Color(c), lines.Color(bg))
		ss = append(ss, sty)
	}
	return ss
}

// System8Backgrounds provides the possible background combinations with
// given foreground color fg.  A possible background combination is any
// System8 color which is not fg.
func System8Backgrounds(fg System8) (ss []lines.Style) {
	for _, c := range system8Colors {
		if c == fg {
			continue
		}
		sty := lines.NewStyle(
			lines.ZeroStyle, lines.Color(fg), lines.Color(c))
		ss = append(ss, sty)
	}
	return ss
}

// LinuxBG types the linux-background colors which are the eight system
// colors.
type LinuxBG int32

const (
	BlackBG  LinuxBG = LinuxBG(lines.Black)
	MaroonBG LinuxBG = LinuxBG(lines.Maroon)
	GreenBG  LinuxBG = LinuxBG(lines.Green)
	OliveBG  LinuxBG = LinuxBG(lines.Olive)
	NavyBG   LinuxBG = LinuxBG(lines.Navy)
	PurpleBG LinuxBG = LinuxBG(lines.Purple)
	TealBG   LinuxBG = LinuxBG(lines.Teal)
	SilverBG LinuxBG = LinuxBG(lines.Silver)
)

var linuxBGs = []LinuxBG{
	BlackBG, MaroonBG, GreenBG, OliveBG, NavyBG, PurpleBG, TealBG,
	SilverBG,
}

// LinuxBackgrounds provides the possible background combinations with
// given foreground color fg.  A possible background combination is any
// LinuxBG color which is not fg.  Note since we have more foreground
// colors than background colors we may get seven or eight combination
// depending if fg is a background color or not.
func LinuxBackgrounds(fg LinuxFG) (ss []lines.Style) {
	for _, c := range linuxBGs {
		if LinuxFG(c) == fg {
			continue
		}
		var sty lines.Style
		if bfg, ok := linuxFGsBold[fg]; ok {
			sty = lines.NewStyle(
				lines.Bold, lines.Color(bfg), lines.Color(c))
		} else {
			sty = lines.NewStyle(
				lines.ZeroStyle, lines.Color(fg), lines.Color(c))
		}
		ss = append(ss, sty)
	}
	return ss
}

// LinuxFG types the linux-foreground color which are the sixteen system
// colors whereas the "upper" eight system-colors are the eight system
// colors with the style-attribute bold.
type LinuxFG int32

const (
	BlackFG   LinuxFG = LinuxFG(lines.Black)
	MaroonFG  LinuxFG = LinuxFG(lines.Maroon)
	GreenFG   LinuxFG = LinuxFG(lines.Green)
	OliveFG   LinuxFG = LinuxFG(lines.Olive)
	NavyFG    LinuxFG = LinuxFG(lines.Navy)
	PurpleFG  LinuxFG = LinuxFG(lines.Purple)
	TealFG    LinuxFG = LinuxFG(lines.Teal)
	SilverFG  LinuxFG = LinuxFG(lines.Silver)
	GreyFG    LinuxFG = LinuxFG(lines.Grey)
	RedFG     LinuxFG = LinuxFG(lines.Red)
	LimeFG    LinuxFG = LinuxFG(lines.Lime)
	YellowFG  LinuxFG = LinuxFG(lines.Yellow)
	BlueFG    LinuxFG = LinuxFG(lines.Blue)
	FuchsiaFG LinuxFG = LinuxFG(lines.Fuchsia)
	AquaFG    LinuxFG = LinuxFG(lines.Aqua)
	WhiteFG   LinuxFG = LinuxFG(lines.White)
)

var linuxFGs = []LinuxFG{
	BlackFG, MaroonFG, GreenFG, OliveFG, NavyFG, PurpleFG, TealFG,
	SilverFG, GreyFG, RedFG, LimeFG, YellowFG, BlueFG, FuchsiaFG, AquaFG,
	WhiteFG,
}
var linuxFGsBold = map[LinuxFG]LinuxBG{
	GreyFG:    BlackBG,
	RedFG:     MaroonBG,
	LimeFG:    GreenBG,
	YellowFG:  OliveBG,
	BlueFG:    NavyBG,
	FuchsiaFG: PurpleBG,
	AquaFG:    TealBG,
	WhiteFG:   SilverBG,
}

// LinuxForegrounds provides the possible foreground combinations with
// given background color bg.  A possible foreground combination is any
// LinuxFG color which is not bg.  Note eight of the linux foreground
// colors are the corresponding background colors (e.g. GreyFG ->
// BlackBG) having the style attribute set to bold.
func LinuxForegrounds(bg LinuxBG) (ss []lines.Style) {
	bgfg := LinuxFG(bg)
	for _, c := range linuxFGs {
		if c == bgfg {
			continue
		}
		var sty lines.Style
		if bg, ok := linuxFGsBold[c]; ok {
			sty = lines.NewStyle(
				lines.Bold, lines.Color(bg), lines.Color(bgfg))
		} else {
			sty = lines.NewStyle(
				lines.ZeroStyle, lines.Color(c), lines.Color(bgfg))
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

var system16Colors = []System{
	Black16, Maroon16, Green16, Olive16, Navy16, Purple16, Teal16,
	Silver16, Grey16, Red16, Lime16, Yellow16, Blue16, Fuchsia16,
	Aqua16, White16}

// System16Foregrounds provides the possible foreground combinations with
// given background color bg.  A possible foreground combination is any
// System color which is not bg.
func System16Foregrounds(bg System) (ss []lines.Style) {
	for _, c := range system16Colors {
		if c == bg {
			continue
		}
		sty := lines.NewStyle(
			lines.ZeroStyle, lines.Color(c), lines.Color(bg))
		ss = append(ss, sty)
	}
	return ss
}

// System16Backgrounds provides the possible background combinations with
// given foreground color fg.  A possible background combination is any
// System color which is not fg.
func System16Backgrounds(fg System) (ss []lines.Style) {
	for _, c := range system16Colors {
		if c == fg {
			continue
		}
		sty := lines.NewStyle(
			lines.ZeroStyle, lines.Color(fg), lines.Color(c))
		ss = append(ss, sty)
	}
	return ss
}

// Ansi types all 256 ANSI colors.
type ANSI int32

const (
	Black             ANSI = ANSI(lines.Black)
	Maroon            ANSI = ANSI(lines.Maroon)
	Green             ANSI = ANSI(lines.Green)
	Olive             ANSI = ANSI(lines.Olive)
	Navy              ANSI = ANSI(lines.Navy)
	Purple            ANSI = ANSI(lines.Purple)
	Teal              ANSI = ANSI(lines.Teal)
	Silver            ANSI = ANSI(lines.Silver)
	Grey              ANSI = ANSI(lines.Grey)
	Red               ANSI = ANSI(lines.Red)
	Lime              ANSI = ANSI(lines.Lime)
	Yellow            ANSI = ANSI(lines.Yellow)
	Blue              ANSI = ANSI(lines.Blue)
	Fuchsia           ANSI = ANSI(lines.Fuchsia)
	Aqua              ANSI = ANSI(lines.Aqua)
	White             ANSI = ANSI(lines.White)
	Grey0             ANSI = ANSI(lines.Grey0)
	NavyBlue          ANSI = ANSI(lines.NavyBlue)
	DarkBlue          ANSI = ANSI(lines.DarkBlue)
	Blue3             ANSI = ANSI(lines.Blue3)
	Blue3_2           ANSI = ANSI(lines.Blue3_2)
	Blue1             ANSI = ANSI(lines.Blue1)
	DarkGreen         ANSI = ANSI(lines.DarkGreen)
	DeepSkyBlue4      ANSI = ANSI(lines.DeepSkyBlue4)
	DeepSkyBlue4_2    ANSI = ANSI(lines.DeepSkyBlue4_2)
	DeepSkyBlue4_3    ANSI = ANSI(lines.DeepSkyBlue4_3)
	DodgerBlue3       ANSI = ANSI(lines.DodgerBlue3)
	DodgerBlue2       ANSI = ANSI(lines.DodgerBlue2)
	Green4            ANSI = ANSI(lines.Green4)
	SpringGreen4      ANSI = ANSI(lines.SpringGreen4)
	Turquoise4        ANSI = ANSI(lines.Turquoise4)
	DeepSkyBlue3      ANSI = ANSI(lines.DeepSkyBlue3)
	DeepSkyBlue3_2    ANSI = ANSI(lines.DeepSkyBlue3_2)
	DodgerBlue1       ANSI = ANSI(lines.DodgerBlue1)
	Green3            ANSI = ANSI(lines.Green3)
	SpringGreen3      ANSI = ANSI(lines.SpringGreen3)
	DarkCyan          ANSI = ANSI(lines.DarkCyan)
	LightSeaGreen     ANSI = ANSI(lines.LightSeaGreen)
	DeepSkyBlue2      ANSI = ANSI(lines.DeepSkyBlue2)
	DeepSkyBlue1      ANSI = ANSI(lines.DeepSkyBlue1)
	Green3_2          ANSI = ANSI(lines.Green3_2)
	SpringGreen3_2    ANSI = ANSI(lines.SpringGreen3_2)
	SpringGreen2      ANSI = ANSI(lines.SpringGreen2)
	Cyan3             ANSI = ANSI(lines.Cyan3)
	DarkTurquoise     ANSI = ANSI(lines.DarkTurquoise)
	Turquoise2        ANSI = ANSI(lines.Turquoise2)
	Green1            ANSI = ANSI(lines.Green1)
	SpringGreen2_2    ANSI = ANSI(lines.SpringGreen2_2)
	SpringGreen1      ANSI = ANSI(lines.SpringGreen1)
	MediumSpringGreen ANSI = ANSI(lines.MediumSpringGreen)
	Cyan2             ANSI = ANSI(lines.Cyan2)
	Cyan1             ANSI = ANSI(lines.Cyan1)
	DarkRed           ANSI = ANSI(lines.DarkRed)
	DeepPink4_3       ANSI = ANSI(lines.DeepPink4_3)
	Purple4           ANSI = ANSI(lines.Purple4)
	Purple4_2         ANSI = ANSI(lines.Purple4_2)
	Purple3           ANSI = ANSI(lines.Purple3)
	BlueViolet        ANSI = ANSI(lines.BlueViolet)
	Orange4           ANSI = ANSI(lines.Orange4)
	Grey37            ANSI = ANSI(lines.Grey37)
	MediumPurple4     ANSI = ANSI(lines.MediumPurple4)
	SlateBlue3        ANSI = ANSI(lines.SlateBlue3)
	SlateBlue3_2      ANSI = ANSI(lines.SlateBlue3_2)
	RoyalBlue1        ANSI = ANSI(lines.RoyalBlue1)
	Chartreuse4       ANSI = ANSI(lines.Chartreuse4)
	DarkSeaGreen4     ANSI = ANSI(lines.DarkSeaGreen4)
	PaleTurquoise4    ANSI = ANSI(lines.PaleTurquoise4)
	SteelBlue         ANSI = ANSI(lines.SteelBlue)
	SteelBlue3        ANSI = ANSI(lines.SteelBlue3)
	CornflowerBlue    ANSI = ANSI(lines.CornflowerBlue)
	Chartreuse3       ANSI = ANSI(lines.Chartreuse3)
	DarkSeaGreen4_2   ANSI = ANSI(lines.DarkSeaGreen4_2)
	CadetBlue         ANSI = ANSI(lines.CadetBlue)
	CadetBlue_2       ANSI = ANSI(lines.CadetBlue_2)
	SkyBlue3          ANSI = ANSI(lines.SkyBlue3)
	SteelBlue1        ANSI = ANSI(lines.SteelBlue1)
	Chartreuse3_2     ANSI = ANSI(lines.Chartreuse3_2)
	PaleGreen3        ANSI = ANSI(lines.PaleGreen3)
	SeaGreen3         ANSI = ANSI(lines.SeaGreen3)
	Aquamarine3       ANSI = ANSI(lines.Aquamarine3)
	MediumTurquoise   ANSI = ANSI(lines.MediumTurquoise)
	SteelBlue1_2      ANSI = ANSI(lines.SteelBlue1_2)
	Chartreuse2       ANSI = ANSI(lines.Chartreuse2)
	SeaGreen2         ANSI = ANSI(lines.SeaGreen2)
	SeaGreen1         ANSI = ANSI(lines.SeaGreen1)
	SeaGreen1_2       ANSI = ANSI(lines.SeaGreen1_2)
	Aquamarine1       ANSI = ANSI(lines.Aquamarine1)
	DarkSlateGray2    ANSI = ANSI(lines.DarkSlateGray2)
	DarkRed_2         ANSI = ANSI(lines.DarkRed_2)
	DeepPink4_2       ANSI = ANSI(lines.DeepPink4_2)
	DarkMagenta       ANSI = ANSI(lines.DarkMagenta)
	DarkMagenta_2     ANSI = ANSI(lines.DarkMagenta_2)
	DarkViolet        ANSI = ANSI(lines.DarkViolet)
	Purple_1          ANSI = ANSI(lines.Purple_1)
	Orange4_2         ANSI = ANSI(lines.Orange4_2)
	LightPink4        ANSI = ANSI(lines.LightPink4)
	Plum4             ANSI = ANSI(lines.Plum4)
	MediumPurple3     ANSI = ANSI(lines.MediumPurple3)
	MediumPurple3_2   ANSI = ANSI(lines.MediumPurple3_2)
	SlateBlue1        ANSI = ANSI(lines.SlateBlue1)
	Yellow4           ANSI = ANSI(lines.Yellow4)
	Wheat4            ANSI = ANSI(lines.Wheat4)
	Grey53            ANSI = ANSI(lines.Grey53)
	LightSlateGrey    ANSI = ANSI(lines.LightSlateGrey)
	MediumPurple      ANSI = ANSI(lines.MediumPurple)
	LightSlateBlue    ANSI = ANSI(lines.LightSlateBlue)
	Yellow4_2         ANSI = ANSI(lines.Yellow4_2)
	DarkOliveGreen3   ANSI = ANSI(lines.DarkOliveGreen3)
	DarkSeaGreen      ANSI = ANSI(lines.DarkSeaGreen)
	LightSkyBlue3     ANSI = ANSI(lines.LightSkyBlue3)
	LightSkyBlue3_2   ANSI = ANSI(lines.LightSkyBlue3_2)
	SkyBlue2          ANSI = ANSI(lines.SkyBlue2)
	Chartreuse2_2     ANSI = ANSI(lines.Chartreuse2_2)
	DarkOliveGreen3_2 ANSI = ANSI(lines.DarkOliveGreen3_2)
	PaleGreen3_2      ANSI = ANSI(lines.PaleGreen3_2)
	DarkSeaGreen3     ANSI = ANSI(lines.DarkSeaGreen3)
	DarkSlateGray3    ANSI = ANSI(lines.DarkSlateGray3)
	SkyBlue1          ANSI = ANSI(lines.SkyBlue1)
	Chartreuse1       ANSI = ANSI(lines.Chartreuse1)
	LightGreen        ANSI = ANSI(lines.LightGreen)
	LightGreen_2      ANSI = ANSI(lines.LightGreen_2)
	PaleGreen1        ANSI = ANSI(lines.PaleGreen1)
	Aquamarine1_2     ANSI = ANSI(lines.Aquamarine1_2)
	DarkSlateGray1    ANSI = ANSI(lines.DarkSlateGray1)
	Red3              ANSI = ANSI(lines.Red3)
	DeepPink4         ANSI = ANSI(lines.DeepPink4)
	MediumVioletRed   ANSI = ANSI(lines.MediumVioletRed)
	Magenta3          ANSI = ANSI(lines.Magenta3)
	DarkViolet_2      ANSI = ANSI(lines.DarkViolet_2)
	Purple_2          ANSI = ANSI(lines.Purple_2)
	DarkOrange3       ANSI = ANSI(lines.DarkOrange3)
	IndianRed         ANSI = ANSI(lines.IndianRed)
	HotPink3          ANSI = ANSI(lines.HotPink3)
	MediumOrchid3     ANSI = ANSI(lines.MediumOrchid3)
	MediumOrchid      ANSI = ANSI(lines.MediumOrchid)
	MediumPurple2     ANSI = ANSI(lines.MediumPurple2)
	DarkGoldenrod     ANSI = ANSI(lines.DarkGoldenrod)
	LightSalmon3      ANSI = ANSI(lines.LightSalmon3)
	RosyBrown         ANSI = ANSI(lines.RosyBrown)
	Grey63            ANSI = ANSI(lines.Grey63)
	MediumPurple2_2   ANSI = ANSI(lines.MediumPurple2_2)
	MediumPurple1     ANSI = ANSI(lines.MediumPurple1)
	Gold3             ANSI = ANSI(lines.Gold3)
	DarkKhaki         ANSI = ANSI(lines.DarkKhaki)
	NavajoWhite3      ANSI = ANSI(lines.NavajoWhite3)
	Grey69            ANSI = ANSI(lines.Grey69)
	LightSteelBlue3   ANSI = ANSI(lines.LightSteelBlue3)
	LightSteelBlue    ANSI = ANSI(lines.LightSteelBlue)
	Yellow3           ANSI = ANSI(lines.Yellow3)
	DarkOliveGreen3_3 ANSI = ANSI(lines.DarkOliveGreen3_3)
	DarkSeaGreen3_2   ANSI = ANSI(lines.DarkSeaGreen3_2)
	DarkSeaGreen2     ANSI = ANSI(lines.DarkSeaGreen2)
	LightCyan3        ANSI = ANSI(lines.LightCyan3)
	LightSkyBlue1     ANSI = ANSI(lines.LightSkyBlue1)
	GreenYellow       ANSI = ANSI(lines.GreenYellow)
	DarkOliveGreen2   ANSI = ANSI(lines.DarkOliveGreen2)
	PaleGreen1_2      ANSI = ANSI(lines.PaleGreen1_2)
	DarkSeaGreen2_2   ANSI = ANSI(lines.DarkSeaGreen2_2)
	DarkSeaGreen1     ANSI = ANSI(lines.DarkSeaGreen1)
	PaleTurquoise1    ANSI = ANSI(lines.PaleTurquoise1)
	Red3_2            ANSI = ANSI(lines.Red3_2)
	DeepPink3         ANSI = ANSI(lines.DeepPink3)
	DeepPink3_2       ANSI = ANSI(lines.DeepPink3_2)
	Magenta3_2        ANSI = ANSI(lines.Magenta3_2)
	Magenta3_3        ANSI = ANSI(lines.Magenta3_3)
	Magenta2          ANSI = ANSI(lines.Magenta2)
	DarkOrange3_2     ANSI = ANSI(lines.DarkOrange3_2)
	IndianRed_2       ANSI = ANSI(lines.IndianRed_2)
	HotPink3_2        ANSI = ANSI(lines.HotPink3_2)
	HotPink2          ANSI = ANSI(lines.HotPink2)
	Orchid            ANSI = ANSI(lines.Orchid)
	MediumOrchid1     ANSI = ANSI(lines.MediumOrchid1)
	Orange3           ANSI = ANSI(lines.Orange3)
	LightSalmon3_2    ANSI = ANSI(lines.LightSalmon3_2)
	LightPink3        ANSI = ANSI(lines.LightPink3)
	Pink3             ANSI = ANSI(lines.Pink3)
	Plum3             ANSI = ANSI(lines.Plum3)
	Violet            ANSI = ANSI(lines.Violet)
	Gold3_2           ANSI = ANSI(lines.Gold3_2)
	LightGoldenrod3   ANSI = ANSI(lines.LightGoldenrod3)
	Tan               ANSI = ANSI(lines.Tan)
	MistyRose3        ANSI = ANSI(lines.MistyRose3)
	Thistle3          ANSI = ANSI(lines.Thistle3)
	Plum2             ANSI = ANSI(lines.Plum2)
	Yellow3_2         ANSI = ANSI(lines.Yellow3_2)
	Khaki3            ANSI = ANSI(lines.Khaki3)
	LightGoldenrod2   ANSI = ANSI(lines.LightGoldenrod2)
	LightYellow3      ANSI = ANSI(lines.LightYellow3)
	Grey84            ANSI = ANSI(lines.Grey84)
	LightSteelBlue1   ANSI = ANSI(lines.LightSteelBlue1)
	Yellow2           ANSI = ANSI(lines.Yellow2)
	DarkOliveGreen1   ANSI = ANSI(lines.DarkOliveGreen1)
	DarkOliveGreen1_2 ANSI = ANSI(lines.DarkOliveGreen1_2)
	DarkSeaGreen1_2   ANSI = ANSI(lines.DarkSeaGreen1_2)
	Honeydew2         ANSI = ANSI(lines.Honeydew2)
	LightCyan1        ANSI = ANSI(lines.LightCyan1)
	Red1              ANSI = ANSI(lines.Red1)
	DeepPink2         ANSI = ANSI(lines.DeepPink2)
	DeepPink1         ANSI = ANSI(lines.DeepPink1)
	DeepPink1_2       ANSI = ANSI(lines.DeepPink1_2)
	Magenta2_2        ANSI = ANSI(lines.Magenta2_2)
	Magenta1          ANSI = ANSI(lines.Magenta1)
	OrangeRed1        ANSI = ANSI(lines.OrangeRed1)
	IndianRed1        ANSI = ANSI(lines.IndianRed1)
	IndianRed1_2      ANSI = ANSI(lines.IndianRed1_2)
	HotPink           ANSI = ANSI(lines.HotPink)
	HotPink_2         ANSI = ANSI(lines.HotPink_2)
	MediumOrchid1_2   ANSI = ANSI(lines.MediumOrchid1_2)
	DarkOrange        ANSI = ANSI(lines.DarkOrange)
	Salmon1           ANSI = ANSI(lines.Salmon1)
	LightCoral        ANSI = ANSI(lines.LightCoral)
	PaleVioletRed1    ANSI = ANSI(lines.PaleVioletRed1)
	Orchid2           ANSI = ANSI(lines.Orchid2)
	Orchid1           ANSI = ANSI(lines.Orchid1)
	Orange1           ANSI = ANSI(lines.Orange1)
	SandyBrown        ANSI = ANSI(lines.SandyBrown)
	LightSalmon1      ANSI = ANSI(lines.LightSalmon1)
	LightPink1        ANSI = ANSI(lines.LightPink1)
	Pink1             ANSI = ANSI(lines.Pink1)
	Plum1             ANSI = ANSI(lines.Plum1)
	Gold1             ANSI = ANSI(lines.Gold1)
	LightGoldenrod2_2 ANSI = ANSI(lines.LightGoldenrod2_2)
	LightGoldenrod2_3 ANSI = ANSI(lines.LightGoldenrod2_3)
	NavajoWhite1      ANSI = ANSI(lines.NavajoWhite1)
	MistyRose1        ANSI = ANSI(lines.MistyRose1)
	Thistle1          ANSI = ANSI(lines.Thistle1)
	Yellow1           ANSI = ANSI(lines.Yellow1)
	LightGoldenrod1   ANSI = ANSI(lines.LightGoldenrod1)
	Khaki1            ANSI = ANSI(lines.Khaki1)
	Wheat1            ANSI = ANSI(lines.Wheat1)
	Cornsilk1         ANSI = ANSI(lines.Cornsilk1)
	Grey100           ANSI = ANSI(lines.Grey100)
	Grey3             ANSI = ANSI(lines.Grey3)
	Grey7             ANSI = ANSI(lines.Grey7)
	Grey11            ANSI = ANSI(lines.Grey11)
	Grey15            ANSI = ANSI(lines.Grey15)
	Grey19            ANSI = ANSI(lines.Grey19)
	Grey23            ANSI = ANSI(lines.Grey23)
	Grey27            ANSI = ANSI(lines.Grey27)
	Grey30            ANSI = ANSI(lines.Grey30)
	Grey35            ANSI = ANSI(lines.Grey35)
	Grey39            ANSI = ANSI(lines.Grey39)
	Grey42            ANSI = ANSI(lines.Grey42)
	Grey46            ANSI = ANSI(lines.Grey46)
	Grey50            ANSI = ANSI(lines.Grey50)
	Grey54            ANSI = ANSI(lines.Grey54)
	Grey58            ANSI = ANSI(lines.Grey58)
	Grey62            ANSI = ANSI(lines.Grey62)
	Grey66            ANSI = ANSI(lines.Grey66)
	Grey70            ANSI = ANSI(lines.Grey70)
	Grey74            ANSI = ANSI(lines.Grey74)
	Grey78            ANSI = ANSI(lines.Grey78)
	Grey82            ANSI = ANSI(lines.Grey82)
	Grey85            ANSI = ANSI(lines.Grey85)
	Grey89            ANSI = ANSI(lines.Grey89)
	Grey93            ANSI = ANSI(lines.Grey93)
)

var ansiColors = []ANSI{
	Black, Maroon, Green, Olive, Navy, Purple, Teal, Silver, Grey, Red,
	Lime, Yellow, Blue, Fuchsia, Aqua, White, Grey0, NavyBlue, DarkBlue,
	Blue3, Blue3_2, Blue1, DarkGreen, DeepSkyBlue4, DeepSkyBlue4_2,
	DeepSkyBlue4_3, DodgerBlue3, DodgerBlue2, Green4, SpringGreen4,
	Turquoise4, DeepSkyBlue3, DeepSkyBlue3_2, DodgerBlue1, Green3,
	SpringGreen3, DarkCyan, LightSeaGreen, DeepSkyBlue2, DeepSkyBlue1,
	Green3_2, SpringGreen3_2, SpringGreen2, Cyan3, DarkTurquoise,
	Turquoise2, Green1, SpringGreen2_2, SpringGreen1, MediumSpringGreen,
	Cyan2, Cyan1, DarkRed, DeepPink4_3, Purple4, Purple4_2, Purple3,
	BlueViolet, Orange4, Grey37, MediumPurple4, SlateBlue3,
	SlateBlue3_2, RoyalBlue1, Chartreuse4, DarkSeaGreen4,
	PaleTurquoise4, SteelBlue, SteelBlue3, CornflowerBlue, Chartreuse3,
	DarkSeaGreen4_2, CadetBlue, CadetBlue_2, SkyBlue3, SteelBlue1,
	Chartreuse3_2, PaleGreen3, SeaGreen3, Aquamarine3, MediumTurquoise,
	SteelBlue1_2, Chartreuse2, SeaGreen2, SeaGreen1, SeaGreen1_2,
	Aquamarine1, DarkSlateGray2, DarkRed_2, DeepPink4_2, DarkMagenta,
	DarkMagenta_2, DarkViolet, Purple_1, Orange4_2, LightPink4, Plum4,
	MediumPurple3, MediumPurple3_2, SlateBlue1, Yellow4, Wheat4, Grey53,
	LightSlateGrey, MediumPurple, LightSlateBlue, Yellow4_2,
	DarkOliveGreen3, DarkSeaGreen, LightSkyBlue3, LightSkyBlue3_2,
	SkyBlue2, Chartreuse2_2, DarkOliveGreen3_2, PaleGreen3_2,
	DarkSeaGreen3, DarkSlateGray3, SkyBlue1, Chartreuse1, LightGreen,
	LightGreen_2, PaleGreen1, Aquamarine1_2, DarkSlateGray1, Red3,
	DeepPink4, MediumVioletRed, Magenta3, DarkViolet_2, Purple_2,
	DarkOrange3, IndianRed, HotPink3, MediumOrchid3, MediumOrchid,
	MediumPurple2, DarkGoldenrod, LightSalmon3, RosyBrown, Grey63,
	MediumPurple2_2, MediumPurple1, Gold3, DarkKhaki, NavajoWhite3,
	Grey69, LightSteelBlue3, LightSteelBlue, Yellow3, DarkOliveGreen3_3,
	DarkSeaGreen3_2, DarkSeaGreen2, LightCyan3, LightSkyBlue1,
	GreenYellow, DarkOliveGreen2, PaleGreen1_2, DarkSeaGreen2_2,
	DarkSeaGreen1, PaleTurquoise1, Red3_2, DeepPink3, DeepPink3_2,
	Magenta3_2, Magenta3_3, Magenta2, DarkOrange3_2, IndianRed_2,
	HotPink3_2, HotPink2, Orchid, MediumOrchid1, Orange3,
	LightSalmon3_2, LightPink3, Pink3, Plum3, Violet, Gold3_2,
	LightGoldenrod3, Tan, MistyRose3, Thistle3, Plum2, Yellow3_2,
	Khaki3, LightGoldenrod2, LightYellow3, Grey84, LightSteelBlue1,
	Yellow2, DarkOliveGreen1, DarkOliveGreen1_2, DarkSeaGreen1_2,
	Honeydew2, LightCyan1, Red1, DeepPink2, DeepPink1, DeepPink1_2,
	Magenta2_2, Magenta1, OrangeRed1, IndianRed1, IndianRed1_2, HotPink,
	HotPink_2, MediumOrchid1_2, DarkOrange, Salmon1, LightCoral,
	PaleVioletRed1, Orchid2, Orchid1, Orange1, SandyBrown, LightSalmon1,
	LightPink1, Pink1, Plum1, Gold1, LightGoldenrod2_2,
	LightGoldenrod2_3, NavajoWhite1, MistyRose1, Thistle1, Yellow1,
	LightGoldenrod1, Khaki1, Wheat1, Cornsilk1, Grey100, Grey3, Grey7,
	Grey11, Grey15, Grey19, Grey23, Grey27, Grey30, Grey35, Grey39,
	Grey42, Grey46, Grey50, Grey54, Grey58, Grey62, Grey66, Grey70,
	Grey74, Grey78, Grey82, Grey85, Grey89, Grey93,
}

// ANSIForegrounds provides the possible foreground combinations with
// given background color bg.  A possible foreground combination is any
// ANSI color which is not bg.
func ANSIForegrounds(bg ANSI) (ss []lines.Style) {
	for _, c := range ansiColors {
		if c == bg {
			continue
		}
		sty := lines.NewStyle(
			lines.ZeroStyle, lines.Color(c), lines.Color(bg))
		ss = append(ss, sty)
	}
	return ss
}

// ANSIBackgrounds provides the possible background combinations with
// given foreground color fg.  A possible background combination is any
// ANSI color which is not fg.
func ANSIBackgrounds(fg ANSI) (ss []lines.Style) {
	for _, c := range ansiColors {
		if c == fg {
			continue
		}
		sty := lines.NewStyle(
			lines.ZeroStyle, lines.Color(fg), lines.Color(c))
		ss = append(ss, sty)
	}
	return ss
}
