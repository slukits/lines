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
	Black                Color = api.Black
	Maroon               Color = api.Maroon
	Green                Color = api.Green
	Olive                Color = api.Olive
	Navy                 Color = api.Navy
	Purple               Color = api.Purple
	Teal                 Color = api.Teal
	Silver               Color = api.Silver
	Gray                 Color = api.Gray
	Red                  Color = api.Red
	Lime                 Color = api.Lime
	Yellow               Color = api.Yellow
	Blue                 Color = api.Blue
	Fuchsia              Color = api.Fuchsia
	Aqua                 Color = api.Aqua
	White                Color = api.White
	AliceBlue            Color = api.AliceBlue
	AntiqueWhite         Color = api.AntiqueWhite
	AquaMarine           Color = api.AquaMarine
	Azure                Color = api.Azure
	Beige                Color = api.Beige
	Bisque               Color = api.Bisque
	BlanchedAlmond       Color = api.BlanchedAlmond
	BlueViolet           Color = api.BlueViolet
	Brown                Color = api.Brown
	BurlyWood            Color = api.BurlyWood
	CadetBlue            Color = api.CadetBlue
	Chartreuse           Color = api.Chartreuse
	Chocolate            Color = api.Chocolate
	Coral                Color = api.Coral
	CornflowerBlue       Color = api.CornflowerBlue
	Cornsilk             Color = api.Cornsilk
	Crimson              Color = api.Crimson
	DarkBlue             Color = api.DarkBlue
	DarkCyan             Color = api.DarkCyan
	DarkGoldenrod        Color = api.DarkGoldenrod
	DarkGray             Color = api.DarkGray
	DarkGreen            Color = api.DarkGreen
	DarkKhaki            Color = api.DarkKhaki
	DarkMagenta          Color = api.DarkMagenta
	DarkOliveGreen       Color = api.DarkOliveGreen
	DarkOrange           Color = api.DarkOrange
	DarkOrchid           Color = api.DarkOrchid
	DarkRed              Color = api.DarkRed
	DarkSalmon           Color = api.DarkSalmon
	DarkSeaGreen         Color = api.DarkSeaGreen
	DarkSlateBlue        Color = api.DarkSlateBlue
	DarkSlateGray        Color = api.DarkSlateGray
	DarkTurquoise        Color = api.DarkTurquoise
	DarkViolet           Color = api.DarkViolet
	DeepPink             Color = api.DeepPink
	DeepSkyBlue          Color = api.DeepSkyBlue
	DimGray              Color = api.DimGray
	DodgerBlue           Color = api.DodgerBlue
	FireBrick            Color = api.FireBrick
	FloralWhite          Color = api.FloralWhite
	ForestGreen          Color = api.ForestGreen
	Gainsboro            Color = api.Gainsboro
	GhostWhite           Color = api.GhostWhite
	Gold                 Color = api.Gold
	Goldenrod            Color = api.Goldenrod
	GreenYellow          Color = api.GreenYellow
	Honeydew             Color = api.Honeydew
	HotPink              Color = api.HotPink
	IndianRed            Color = api.IndianRed
	Indigo               Color = api.Indigo
	Ivory                Color = api.Ivory
	Khaki                Color = api.Khaki
	Lavender             Color = api.Lavender
	LavenderBlush        Color = api.LavenderBlush
	LawnGreen            Color = api.LawnGreen
	LemonChiffon         Color = api.LemonChiffon
	LightBlue            Color = api.LightBlue
	LightCoral           Color = api.LightCoral
	LightCyan            Color = api.LightCyan
	LightGoldenrodYellow Color = api.LightGoldenrodYellow
	LightGray            Color = api.LightGray
	LightGreen           Color = api.LightGreen
	LightPink            Color = api.LightPink
	LightSalmon          Color = api.LightSalmon
	LightSeaGreen        Color = api.LightSeaGreen
	LightSkyBlue         Color = api.LightSkyBlue
	LightSlateGray       Color = api.LightSlateGray
	LightSteelBlue       Color = api.LightSteelBlue
	LightYellow          Color = api.LightYellow
	LimeGreen            Color = api.LimeGreen
	Linen                Color = api.Linen
	MediumAquamarine     Color = api.MediumAquamarine
	MediumBlue           Color = api.MediumBlue
	MediumOrchid         Color = api.MediumOrchid
	MediumPurple         Color = api.MediumPurple
	MediumSeaGreen       Color = api.MediumSeaGreen
	MediumSlateBlue      Color = api.MediumSlateBlue
	MediumSpringGreen    Color = api.MediumSpringGreen
	MediumTurquoise      Color = api.MediumTurquoise
	MediumVioletRed      Color = api.MediumVioletRed
	MidnightBlue         Color = api.MidnightBlue
	MintCream            Color = api.MintCream
	MistyRose            Color = api.MistyRose
	Moccasin             Color = api.Moccasin
	NavajoWhite          Color = api.NavajoWhite
	OldLace              Color = api.OldLace
	OliveDrab            Color = api.OliveDrab
	Orange               Color = api.Orange
	OrangeRed            Color = api.OrangeRed
	Orchid               Color = api.Orchid
	PaleGoldenrod        Color = api.PaleGoldenrod
	PaleGreen            Color = api.PaleGreen
	PaleTurquoise        Color = api.PaleTurquoise
	PaleVioletRed        Color = api.PaleVioletRed
	PapayaWhip           Color = api.PapayaWhip
	PeachPuff            Color = api.PeachPuff
	Peru                 Color = api.Peru
	Pink                 Color = api.Pink
	Plum                 Color = api.Plum
	PowderBlue           Color = api.PowderBlue
	RebeccaPurple        Color = api.RebeccaPurple
	RosyBrown            Color = api.RosyBrown
	RoyalBlue            Color = api.RoyalBlue
	SaddleBrown          Color = api.SaddleBrown
	Salmon               Color = api.Salmon
	SandyBrown           Color = api.SandyBrown
	SeaGreen             Color = api.SeaGreen
	Seashell             Color = api.Seashell
	Sienna               Color = api.Sienna
	Skyblue              Color = api.Skyblue
	SlateBlue            Color = api.SlateBlue
	SlateGray            Color = api.SlateGray
	Snow                 Color = api.Snow
	SpringGreen          Color = api.SpringGreen
	SteelBlue            Color = api.SteelBlue
	Tan                  Color = api.Tan
	Thistle              Color = api.Thistle
	Tomato               Color = api.Tomato
	Turquoise            Color = api.Turquoise
	Violet               Color = api.Violet
	Wheat                Color = api.Wheat
	WhiteSmoke           Color = api.WhiteSmoke
	YellowGreen          Color = api.YellowGreen

	DefaultColor Color = api.DefaultColor
)
