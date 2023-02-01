// Copyright (c) 2022 Stephan Lukits. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package lines

import "github.com/slukits/lines/internal/api"

// Color represents an rgb color which is usually given in the typical
// hex-format 0xRRGGBB whereas R, G and B are hex-digits, i.e. red is
// 0xFF0000.
type Color = api.Color

const (
	Black             Color = api.Black
	Maroon            Color = api.Maroon
	Green             Color = api.Green
	Olive             Color = api.Olive
	Navy              Color = api.Navy
	Purple            Color = api.Purple
	Teal              Color = api.Teal
	Silver            Color = api.Silver
	Grey              Color = api.Grey
	Red               Color = api.Red
	Lime              Color = api.Lime
	Yellow            Color = api.Yellow
	Blue              Color = api.Blue
	Fuchsia           Color = api.Fuchsia
	Aqua              Color = api.Aqua
	White             Color = api.White
	Grey0             Color = api.Grey0
	NavyBlue          Color = api.NavyBlue
	DarkBlue          Color = api.DarkBlue
	Blue3             Color = api.Blue3
	Blue3_2           Color = api.Blue3_2
	Blue1             Color = api.Blue1
	DarkGreen         Color = api.DarkGreen
	DeepSkyBlue4      Color = api.DeepSkyBlue4
	DeepSkyBlue4_2    Color = api.DeepSkyBlue4_2
	DeepSkyBlue4_3    Color = api.DeepSkyBlue4_3
	DodgerBlue3       Color = api.DodgerBlue3
	DodgerBlue2       Color = api.DodgerBlue2
	Green4            Color = api.Green4
	SpringGreen4      Color = api.SpringGreen4
	Turquoise4        Color = api.Turquoise4
	DeepSkyBlue3      Color = api.DeepSkyBlue3
	DeepSkyBlue3_2    Color = api.DeepSkyBlue3_2
	DodgerBlue1       Color = api.DodgerBlue1
	Green3            Color = api.Green3
	SpringGreen3      Color = api.SpringGreen3
	DarkCyan          Color = api.DarkCyan
	LightSeaGreen     Color = api.LightSeaGreen
	DeepSkyBlue2      Color = api.DeepSkyBlue2
	DeepSkyBlue1      Color = api.DeepSkyBlue1
	Green3_2          Color = api.Green3_2
	SpringGreen3_2    Color = api.SpringGreen3_2
	SpringGreen2      Color = api.SpringGreen2
	Cyan3             Color = api.Cyan3
	DarkTurquoise     Color = api.DarkTurquoise
	Turquoise2        Color = api.Turquoise2
	Green1            Color = api.Green1
	SpringGreen2_2    Color = api.SpringGreen2_2
	SpringGreen1      Color = api.SpringGreen1
	MediumSpringGreen Color = api.MediumSpringGreen
	Cyan2             Color = api.Cyan2
	Cyan1             Color = api.Cyan1
	DarkRed           Color = api.DarkRed
	DeepPink4_3       Color = api.DeepPink4_3
	Purple4           Color = api.Purple4
	Purple4_2         Color = api.Purple4_2
	Purple3           Color = api.Purple3
	BlueViolet        Color = api.BlueViolet
	Orange4           Color = api.Orange4
	Grey37            Color = api.Grey37
	MediumPurple4     Color = api.MediumPurple4
	SlateBlue3        Color = api.SlateBlue3
	SlateBlue3_2      Color = api.SlateBlue3_2
	RoyalBlue1        Color = api.RoyalBlue1
	Chartreuse4       Color = api.Chartreuse4
	DarkSeaGreen4     Color = api.DarkSeaGreen4
	PaleTurquoise4    Color = api.PaleTurquoise4
	SteelBlue         Color = api.SteelBlue
	SteelBlue3        Color = api.SteelBlue3
	CornflowerBlue    Color = api.CornflowerBlue
	Chartreuse3       Color = api.Chartreuse3
	DarkSeaGreen4_2   Color = api.DarkSeaGreen4_2
	CadetBlue         Color = api.CadetBlue
	CadetBlue_2       Color = api.CadetBlue_2
	SkyBlue3          Color = api.SkyBlue3
	SteelBlue1        Color = api.SteelBlue1
	Chartreuse3_2     Color = api.Chartreuse3_2
	PaleGreen3        Color = api.PaleGreen3
	SeaGreen3         Color = api.SeaGreen3
	Aquamarine3       Color = api.Aquamarine3
	MediumTurquoise   Color = api.MediumTurquoise
	SteelBlue1_2      Color = api.SteelBlue1_2
	Chartreuse2       Color = api.Chartreuse2
	SeaGreen2         Color = api.SeaGreen2
	SeaGreen1         Color = api.SeaGreen1
	SeaGreen1_2       Color = api.SeaGreen1_2
	Aquamarine1       Color = api.Aquamarine1
	DarkSlateGray2    Color = api.DarkSlateGray2
	DarkRed_2         Color = api.DarkRed_2
	DeepPink4_2       Color = api.DeepPink4_2
	DarkMagenta       Color = api.DarkMagenta
	DarkMagenta_2     Color = api.DarkMagenta_2
	DarkViolet        Color = api.DarkViolet
	Purple_1          Color = api.Purple_1
	Orange4_2         Color = api.Orange4_2
	LightPink4        Color = api.LightPink4
	Plum4             Color = api.Plum4
	MediumPurple3     Color = api.MediumPurple3
	MediumPurple3_2   Color = api.MediumPurple3_2
	SlateBlue1        Color = api.SlateBlue1
	Yellow4           Color = api.Yellow4
	Wheat4            Color = api.Wheat4
	Grey53            Color = api.Grey53
	LightSlateGrey    Color = api.LightSlateGrey
	MediumPurple      Color = api.MediumPurple
	LightSlateBlue    Color = api.LightSlateBlue
	Yellow4_2         Color = api.Yellow4_2
	DarkOliveGreen3   Color = api.DarkOliveGreen3
	DarkSeaGreen      Color = api.DarkSeaGreen
	LightSkyBlue3     Color = api.LightSkyBlue3
	LightSkyBlue3_2   Color = api.LightSkyBlue3_2
	SkyBlue2          Color = api.SkyBlue2
	Chartreuse2_2     Color = api.Chartreuse2_2
	DarkOliveGreen3_2 Color = api.DarkOliveGreen3_2
	PaleGreen3_2      Color = api.PaleGreen3_2
	DarkSeaGreen3     Color = api.DarkSeaGreen3
	DarkSlateGray3    Color = api.DarkSlateGray3
	SkyBlue1          Color = api.SkyBlue1
	Chartreuse1       Color = api.Chartreuse1
	LightGreen        Color = api.LightGreen
	LightGreen_2      Color = api.LightGreen_2
	PaleGreen1        Color = api.PaleGreen1
	Aquamarine1_2     Color = api.Aquamarine1_2
	DarkSlateGray1    Color = api.DarkSlateGray1
	Red3              Color = api.Red3
	DeepPink4         Color = api.DeepPink4
	MediumVioletRed   Color = api.MediumVioletRed
	Magenta3          Color = api.Magenta3
	DarkViolet_2      Color = api.DarkViolet_2
	Purple_2          Color = api.Purple_2
	DarkOrange3       Color = api.DarkOrange3
	IndianRed         Color = api.IndianRed
	HotPink3          Color = api.HotPink3
	MediumOrchid3     Color = api.MediumOrchid3
	MediumOrchid      Color = api.MediumOrchid
	MediumPurple2     Color = api.MediumPurple2
	DarkGoldenrod     Color = api.DarkGoldenrod
	LightSalmon3      Color = api.LightSalmon3
	RosyBrown         Color = api.RosyBrown
	Grey63            Color = api.Grey63
	MediumPurple2_2   Color = api.MediumPurple2_2
	MediumPurple1     Color = api.MediumPurple1
	Gold3             Color = api.Gold3
	DarkKhaki         Color = api.DarkKhaki
	NavajoWhite3      Color = api.NavajoWhite3
	Grey69            Color = api.Grey69
	LightSteelBlue3   Color = api.LightSteelBlue3
	LightSteelBlue    Color = api.LightSteelBlue
	Yellow3           Color = api.Yellow3
	DarkOliveGreen3_3 Color = api.DarkOliveGreen3_3
	DarkSeaGreen3_2   Color = api.DarkSeaGreen3_2
	DarkSeaGreen2     Color = api.DarkSeaGreen2
	LightCyan3        Color = api.LightCyan3
	LightSkyBlue1     Color = api.LightSkyBlue1
	GreenYellow       Color = api.GreenYellow
	DarkOliveGreen2   Color = api.DarkOliveGreen2
	PaleGreen1_2      Color = api.PaleGreen1_2
	DarkSeaGreen2_2   Color = api.DarkSeaGreen2_2
	DarkSeaGreen1     Color = api.DarkSeaGreen1
	PaleTurquoise1    Color = api.PaleTurquoise1
	Red3_2            Color = api.Red3_2
	DeepPink3         Color = api.DeepPink3
	DeepPink3_2       Color = api.DeepPink3_2
	Magenta3_2        Color = api.Magenta3_2
	Magenta3_3        Color = api.Magenta3_3
	Magenta2          Color = api.Magenta2
	DarkOrange3_2     Color = api.DarkOrange3_2
	IndianRed_2       Color = api.IndianRed_2
	HotPink3_2        Color = api.HotPink3_2
	HotPink2          Color = api.HotPink2
	Orchid            Color = api.Orchid
	MediumOrchid1     Color = api.MediumOrchid1
	Orange3           Color = api.Orange3
	LightSalmon3_2    Color = api.LightSalmon3_2
	LightPink3        Color = api.LightPink3
	Pink3             Color = api.Pink3
	Plum3             Color = api.Plum3
	Violet            Color = api.Violet
	Gold3_2           Color = api.Gold3_2
	LightGoldenrod3   Color = api.LightGoldenrod3
	Tan               Color = api.Tan
	MistyRose3        Color = api.MistyRose3
	Thistle3          Color = api.Thistle3
	Plum2             Color = api.Plum2
	Yellow3_2         Color = api.Yellow3_2
	Khaki3            Color = api.Khaki3
	LightGoldenrod2   Color = api.LightGoldenrod2
	LightYellow3      Color = api.LightYellow3
	Grey84            Color = api.Grey84
	LightSteelBlue1   Color = api.LightSteelBlue1
	Yellow2           Color = api.Yellow2
	DarkOliveGreen1   Color = api.DarkOliveGreen1
	DarkOliveGreen1_2 Color = api.DarkOliveGreen1_2
	DarkSeaGreen1_2   Color = api.DarkSeaGreen1_2
	Honeydew2         Color = api.Honeydew2
	LightCyan1        Color = api.LightCyan1
	Red1              Color = api.Red1
	DeepPink2         Color = api.DeepPink2
	DeepPink1         Color = api.DeepPink1
	DeepPink1_2       Color = api.DeepPink1_2
	Magenta2_2        Color = api.Magenta2_2
	Magenta1          Color = api.Magenta1
	OrangeRed1        Color = api.OrangeRed1
	IndianRed1        Color = api.IndianRed1
	IndianRed1_2      Color = api.IndianRed1_2
	HotPink           Color = api.HotPink
	HotPink_2         Color = api.HotPink_2
	MediumOrchid1_2   Color = api.MediumOrchid1_2
	DarkOrange        Color = api.DarkOrange
	Salmon1           Color = api.Salmon1
	LightCoral        Color = api.LightCoral
	PaleVioletRed1    Color = api.PaleVioletRed1
	Orchid2           Color = api.Orchid2
	Orchid1           Color = api.Orchid1
	Orange1           Color = api.Orange1
	SandyBrown        Color = api.SandyBrown
	LightSalmon1      Color = api.LightSalmon1
	LightPink1        Color = api.LightPink1
	Pink1             Color = api.Pink1
	Plum1             Color = api.Plum1
	Gold1             Color = api.Gold1
	LightGoldenrod2_2 Color = api.LightGoldenrod2_2
	LightGoldenrod2_3 Color = api.LightGoldenrod2_3
	NavajoWhite1      Color = api.NavajoWhite1
	MistyRose1        Color = api.MistyRose1
	Thistle1          Color = api.Thistle1
	Yellow1           Color = api.Yellow1
	LightGoldenrod1   Color = api.LightGoldenrod1
	Khaki1            Color = api.Khaki1
	Wheat1            Color = api.Wheat1
	Cornsilk1         Color = api.Cornsilk1
	Grey100           Color = api.Grey100
	Grey3             Color = api.Grey3
	Grey7             Color = api.Grey7
	Grey11            Color = api.Grey11
	Grey15            Color = api.Grey15
	Grey19            Color = api.Grey19
	Grey23            Color = api.Grey23
	Grey27            Color = api.Grey27
	Grey30            Color = api.Grey30
	Grey35            Color = api.Grey35
	Grey39            Color = api.Grey39
	Grey42            Color = api.Grey42
	Grey46            Color = api.Grey46
	Grey50            Color = api.Grey50
	Grey54            Color = api.Grey54
	Grey58            Color = api.Grey58
	Grey62            Color = api.Grey62
	Grey66            Color = api.Grey66
	Grey70            Color = api.Grey70
	Grey74            Color = api.Grey74
	Grey78            Color = api.Grey78
	Grey82            Color = api.Grey82
	Grey85            Color = api.Grey85
	Grey89            Color = api.Grey89
	Grey93            Color = api.Grey93

	DefaultColor Color = api.DefaultColor
)

var ColorNames = map[Color]string{
	Black:             "Black",
	Maroon:            "Maroon",
	Green:             "Green",
	Olive:             "Olive",
	Navy:              "Navy",
	Purple:            "Purple",
	Teal:              "Teal",
	Silver:            "Silver",
	Grey:              "Grey",
	Red:               "Red",
	Lime:              "Lime",
	Yellow:            "Yellow",
	Blue:              "Blue",
	Fuchsia:           "Fuchsia",
	Aqua:              "Aqua",
	White:             "White",
	NavyBlue:          "NavyBlue",
	DarkBlue:          "DarkBlue",
	Blue3:             "Blue3",
	Blue3_2:           "Blue3(2)",
	DarkGreen:         "DarkGreen",
	DeepSkyBlue4:      "DeepSkyBlue4",
	DeepSkyBlue4_2:    "DeepSkyBlue4(2)",
	DeepSkyBlue4_3:    "DeepSkyBlue4(3)",
	DodgerBlue3:       "DodgerBlue3",
	DodgerBlue2:       "DodgerBlue2",
	Green4:            "Green4",
	SpringGreen4:      "SpringGreen4",
	Turquoise4:        "Turquoise4",
	DeepSkyBlue3:      "DeepSkyBlue3",
	DeepSkyBlue3_2:    "DeepSkyBlue3(2)",
	DodgerBlue1:       "DodgerBlue1",
	Green3:            "Green3",
	SpringGreen3:      "SpringGreen3",
	DarkCyan:          "DarkCyan",
	LightSeaGreen:     "LightSeaGreen",
	DeepSkyBlue2:      "DeepSkyBlue2",
	DeepSkyBlue1:      "DeepSkyBlue1",
	Green3_2:          "Green3(2)",
	SpringGreen3_2:    "SpringGreen3(2)",
	SpringGreen2:      "SpringGreen2",
	Cyan3:             "Cyan3",
	DarkTurquoise:     "DarkTurquoise",
	Turquoise2:        "Turquoise2",
	SpringGreen2_2:    "SpringGreen2(2)",
	SpringGreen1:      "SpringGreen1",
	MediumSpringGreen: "MediumSpringGreen",
	Cyan2:             "Cyan2",
	DarkRed:           "DarkRed",
	DeepPink4_3:       "DeepPink4(3)",
	Purple4:           "Purple4",
	Purple4_2:         "Purple4(2)",
	Purple3:           "Purple3",
	BlueViolet:        "BlueViolet",
	Orange4:           "Orange4",
	Grey37:            "Grey37",
	MediumPurple4:     "MediumPurple4",
	SlateBlue3:        "SlateBlue3",
	SlateBlue3_2:      "SlateBlue3(2)",
	RoyalBlue1:        "RoyalBlue1",
	Chartreuse4:       "Chartreuse4",
	DarkSeaGreen4:     "DarkSeaGreen4",
	PaleTurquoise4:    "PaleTurquoise4",
	SteelBlue:         "SteelBlue",
	SteelBlue3:        "SteelBlue3",
	CornflowerBlue:    "CornflowerBlue",
	Chartreuse3:       "Chartreuse3",
	DarkSeaGreen4_2:   "DarkSeaGreen4(2)",
	CadetBlue:         "CadetBlue",
	CadetBlue_2:       "CadetBlue(2)",
	SkyBlue3:          "SkyBlue3",
	SteelBlue1:        "SteelBlue1",
	Chartreuse3_2:     "Chartreuse3(2)",
	PaleGreen3:        "PaleGreen3",
	SeaGreen3:         "SeaGreen3",
	Aquamarine3:       "Aquamarine3",
	MediumTurquoise:   "MediumTurquoise",
	SteelBlue1_2:      "SteelBlue1(2)",
	Chartreuse2:       "Chartreuse2",
	SeaGreen2:         "SeaGreen2",
	SeaGreen1:         "SeaGreen1",
	SeaGreen1_2:       "SeaGreen1(2)",
	Aquamarine1:       "Aquamarine1",
	DarkSlateGray2:    "DarkSlateGray2",
	DarkRed_2:         "DarkRed(2)",
	DeepPink4_2:       "DeepPink4(2)",
	DarkMagenta:       "DarkMagenta",
	DarkMagenta_2:     "DarkMagenta(2)",
	DarkViolet:        "DarkViolet",
	Purple_1:          "Purple(1)",
	Orange4_2:         "Orange4(2)",
	LightPink4:        "LightPink4",
	Plum4:             "Plum4",
	MediumPurple3:     "MediumPurple3",
	MediumPurple3_2:   "MediumPurple3(2)",
	SlateBlue1:        "SlateBlue1",
	Yellow4:           "Yellow4",
	Wheat4:            "Wheat4",
	Grey53:            "Grey53",
	LightSlateGrey:    "LightSlateGrey",
	MediumPurple:      "MediumPurple",
	LightSlateBlue:    "LightSlateBlue",
	Yellow4_2:         "Yellow4(2)",
	DarkOliveGreen3:   "DarkOliveGreen3",
	DarkSeaGreen:      "DarkSeaGreen",
	LightSkyBlue3:     "LightSkyBlue3",
	LightSkyBlue3_2:   "LightSkyBlue3(2)",
	SkyBlue2:          "SkyBlue2",
	Chartreuse2_2:     "Chartreuse2(2)",
	DarkOliveGreen3_2: "DarkOliveGreen3(2)",
	PaleGreen3_2:      "PaleGreen3(2)",
	DarkSeaGreen3:     "DarkSeaGreen3",
	DarkSlateGray3:    "DarkSlateGray3",
	SkyBlue1:          "SkyBlue1",
	Chartreuse1:       "Chartreuse1",
	LightGreen:        "LightGreen",
	LightGreen_2:      "LightGreen(2)",
	PaleGreen1:        "PaleGreen1",
	Aquamarine1_2:     "Aquamarine1(2)",
	DarkSlateGray1:    "DarkSlateGray1",
	Red3:              "Red3",
	DeepPink4:         "DeepPink4",
	MediumVioletRed:   "MediumVioletRed",
	Magenta3:          "Magenta3",
	DarkViolet_2:      "DarkViolet(2)",
	Purple_2:          "Purple(2)",
	DarkOrange3:       "DarkOrange3",
	IndianRed:         "IndianRed",
	HotPink3:          "HotPink3",
	MediumOrchid3:     "MediumOrchid3",
	MediumOrchid:      "MediumOrchid",
	MediumPurple2:     "MediumPurple2",
	DarkGoldenrod:     "DarkGoldenrod",
	LightSalmon3:      "LightSalmon3",
	RosyBrown:         "RosyBrown",
	Grey63:            "Grey63",
	MediumPurple2_2:   "MediumPurple2(2)",
	MediumPurple1:     "MediumPurple1",
	Gold3:             "Gold3",
	DarkKhaki:         "DarkKhaki",
	NavajoWhite3:      "NavajoWhite3",
	Grey69:            "Grey69",
	LightSteelBlue3:   "LightSteelBlue3",
	LightSteelBlue:    "LightSteelBlue",
	Yellow3:           "Yellow3",
	DarkOliveGreen3_3: "DarkOliveGreen3(3)",
	DarkSeaGreen3_2:   "DarkSeaGreen3(2)",
	DarkSeaGreen2:     "DarkSeaGreen2",
	LightCyan3:        "LightCyan3",
	LightSkyBlue1:     "LightSkyBlue1",
	GreenYellow:       "GreenYellow",
	DarkOliveGreen2:   "DarkOliveGreen2",
	PaleGreen1_2:      "PaleGreen1(2)",
	DarkSeaGreen2_2:   "DarkSeaGreen2(2)",
	DarkSeaGreen1:     "DarkSeaGreen1",
	PaleTurquoise1:    "PaleTurquoise1",
	Red3_2:            "Red3(2)",
	DeepPink3:         "DeepPink3",
	DeepPink3_2:       "DeepPink3(2)",
	Magenta3_2:        "Magenta3(2)",
	Magenta3_3:        "Magenta3(3)",
	Magenta2:          "Magenta2",
	DarkOrange3_2:     "DarkOrange3(2)",
	IndianRed_2:       "IndianRed(2)",
	HotPink3_2:        "HotPink3(2)",
	HotPink2:          "HotPink2",
	Orchid:            "Orchid",
	MediumOrchid1:     "MediumOrchid1",
	Orange3:           "Orange3",
	LightSalmon3_2:    "LightSalmon3(2)",
	LightPink3:        "LightPink3",
	Pink3:             "Pink3",
	Plum3:             "Plum3",
	Violet:            "Violet",
	Gold3_2:           "Gold3(2)",
	LightGoldenrod3:   "LightGoldenrod3",
	Tan:               "Tan",
	MistyRose3:        "MistyRose3",
	Thistle3:          "Thistle3",
	Plum2:             "Plum2",
	Yellow3_2:         "Yellow3(2)",
	Khaki3:            "Khaki3",
	LightGoldenrod2:   "LightGoldenrod2",
	LightYellow3:      "LightYellow3",
	Grey84:            "Grey84",
	LightSteelBlue1:   "LightSteelBlue1",
	Yellow2:           "Yellow2",
	DarkOliveGreen1:   "DarkOliveGreen1",
	DarkOliveGreen1_2: "DarkOliveGreen1(2)",
	DarkSeaGreen1_2:   "DarkSeaGreen1(2)",
	Honeydew2:         "Honeydew2",
	LightCyan1:        "LightCyan1",
	DeepPink2:         "DeepPink2",
	DeepPink1:         "DeepPink1",
	DeepPink1_2:       "DeepPink1(2)",
	Magenta2_2:        "Magenta2(2)",
	OrangeRed1:        "OrangeRed1",
	IndianRed1:        "IndianRed1",
	IndianRed1_2:      "IndianRed1(2)",
	HotPink:           "HotPink",
	HotPink_2:         "HotPink(2)",
	MediumOrchid1_2:   "MediumOrchid1(2)",
	DarkOrange:        "DarkOrange",
	Salmon1:           "Salmon1",
	LightCoral:        "LightCoral",
	PaleVioletRed1:    "PaleVioletRed1",
	Orchid2:           "Orchid2",
	Orchid1:           "Orchid1",
	Orange1:           "Orange1",
	SandyBrown:        "SandyBrown",
	LightSalmon1:      "LightSalmon1",
	LightPink1:        "LightPink1",
	Pink1:             "Pink1",
	Plum1:             "Plum1",
	Gold1:             "Gold1",
	LightGoldenrod2_2: "LightGoldenrod2(2)",
	LightGoldenrod2_3: "LightGoldenrod2(3)",
	NavajoWhite1:      "NavajoWhite1",
	MistyRose1:        "MistyRose1",
	Thistle1:          "Thistle1",
	LightGoldenrod1:   "LightGoldenrod1",
	Khaki1:            "Khaki1",
	Wheat1:            "Wheat1",
	Cornsilk1:         "Cornsilk1",
	Grey3:             "Grey3",
	Grey7:             "Grey7",
	Grey11:            "Grey11",
	Grey15:            "Grey15",
	Grey19:            "Grey19",
	Grey23:            "Grey23",
	Grey27:            "Grey27",
	Grey30:            "Grey30",
	Grey35:            "Grey35",
	Grey39:            "Grey39",
	Grey42:            "Grey42",
	Grey46:            "Grey46",
	Grey54:            "Grey54",
	Grey58:            "Grey58",
	Grey62:            "Grey62",
	Grey66:            "Grey66",
	Grey70:            "Grey70",
	Grey74:            "Grey74",
	Grey78:            "Grey78",
	Grey82:            "Grey82",
	Grey85:            "Grey85",
	Grey89:            "Grey89",
	Grey93:            "Grey93",
}
