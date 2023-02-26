package api

// Color represents an rgb color.  Predefined colors are expressed in
// the typical hex-format 0xRRGGBB whereas R, G and B are hex-digits,
// i.e. red is 0xFF0000.
type Color int32

type CCC struct {
	C1st  Color
	C2nd  Color
	C3rd  Color
	C4th  Color
	C5th  Color
	C6th  Color
	C7th  Color
	C8th  Color
	C9th  Color
	C10th Color
	C11th Color
	C12th Color
	C13th Color
	C14th Color
	C15th Color
	C16th Color
}

const (
	Black             Color = 0x000000
	Maroon            Color = 0x800000
	Green             Color = 0x008000
	Olive             Color = 0x808000
	Navy              Color = 0x000080
	Purple            Color = 0x800080
	Teal              Color = 0x008080
	Silver            Color = 0xc0c0c0
	Grey              Color = 0x808080
	Red               Color = 0xff0000
	Lime              Color = 0x00ff00
	Yellow            Color = 0xffff00
	Blue              Color = 0x0000ff
	Fuchsia           Color = 0xff00ff
	Aqua              Color = 0x00ffff
	White             Color = 0xffffff
	Grey0             Color = 0x000000
	NavyBlue          Color = 0x00005f
	DarkBlue          Color = 0x000087
	Blue3             Color = 0x0000af
	Blue3_2           Color = 0x0000d7
	Blue1             Color = 0x0000ff
	DarkGreen         Color = 0x005f00
	DeepSkyBlue4      Color = 0x005f5f
	DeepSkyBlue4_2    Color = 0x005f87
	DeepSkyBlue4_3    Color = 0x005faf
	DodgerBlue3       Color = 0x005fd7
	DodgerBlue2       Color = 0x005fff
	Green4            Color = 0x008700
	SpringGreen4      Color = 0x00875f
	Turquoise4        Color = 0x008787
	DeepSkyBlue3      Color = 0x0087af
	DeepSkyBlue3_2    Color = 0x0087d7
	DodgerBlue1       Color = 0x0087ff
	Green3            Color = 0x00af00
	SpringGreen3      Color = 0x00af5f
	DarkCyan          Color = 0x00af87
	LightSeaGreen     Color = 0x00afaf
	DeepSkyBlue2      Color = 0x00afd7
	DeepSkyBlue1      Color = 0x00afff
	Green3_2          Color = 0x00d700
	SpringGreen3_2    Color = 0x00d75f
	SpringGreen2      Color = 0x00d787
	Cyan3             Color = 0x00d7af
	DarkTurquoise     Color = 0x00d7d7
	Turquoise2        Color = 0x00d7ff
	Green1            Color = 0x00ff00
	SpringGreen2_2    Color = 0x00ff5f
	SpringGreen1      Color = 0x00ff87
	MediumSpringGreen Color = 0x00ffaf
	Cyan2             Color = 0x00ffd7
	Cyan1             Color = 0x00ffff
	DarkRed           Color = 0x5f0000
	DeepPink4_3       Color = 0x5f005f
	Purple4           Color = 0x5f0087
	Purple4_2         Color = 0x5f00af
	Purple3           Color = 0x5f00d7
	BlueViolet        Color = 0x5f00ff
	Orange4           Color = 0x5f5f00
	Grey37            Color = 0x5f5f5f
	MediumPurple4     Color = 0x5f5f87
	SlateBlue3        Color = 0x5f5faf
	SlateBlue3_2      Color = 0x5f5fd7
	RoyalBlue1        Color = 0x5f5fff
	Chartreuse4       Color = 0x5f8700
	DarkSeaGreen4     Color = 0x5f875f
	PaleTurquoise4    Color = 0x5f8787
	SteelBlue         Color = 0x5f87af
	SteelBlue3        Color = 0x5f87d7
	CornflowerBlue    Color = 0x5f87ff
	Chartreuse3       Color = 0x5faf00
	DarkSeaGreen4_2   Color = 0x5faf5f
	CadetBlue         Color = 0x5faf87
	CadetBlue_2       Color = 0x5fafaf
	SkyBlue3          Color = 0x5fafd7
	SteelBlue1        Color = 0x5fafff
	Chartreuse3_2     Color = 0x5fd700
	PaleGreen3        Color = 0x5fd75f
	SeaGreen3         Color = 0x5fd787
	Aquamarine3       Color = 0x5fd7af
	MediumTurquoise   Color = 0x5fd7d7
	SteelBlue1_2      Color = 0x5fd7ff
	Chartreuse2       Color = 0x5fff00
	SeaGreen2         Color = 0x5fff5f
	SeaGreen1         Color = 0x5fff87
	SeaGreen1_2       Color = 0x5fffaf
	Aquamarine1       Color = 0x5fffd7
	DarkSlateGray2    Color = 0x5fffff
	DarkRed_2         Color = 0x870000
	DeepPink4_2       Color = 0x87005f
	DarkMagenta       Color = 0x870087
	DarkMagenta_2     Color = 0x8700af
	DarkViolet        Color = 0x8700d7
	Purple_1          Color = 0x8700ff
	Orange4_2         Color = 0x875f00
	LightPink4        Color = 0x875f5f
	Plum4             Color = 0x875f87
	MediumPurple3     Color = 0x875faf
	MediumPurple3_2   Color = 0x875fd7
	SlateBlue1        Color = 0x875fff
	Yellow4           Color = 0x878700
	Wheat4            Color = 0x87875f
	Grey53            Color = 0x878787
	LightSlateGrey    Color = 0x8787af
	MediumPurple      Color = 0x8787d7
	LightSlateBlue    Color = 0x8787ff
	Yellow4_2         Color = 0x87af00
	DarkOliveGreen3   Color = 0x87af5f
	DarkSeaGreen      Color = 0x87af87
	LightSkyBlue3     Color = 0x87afaf
	LightSkyBlue3_2   Color = 0x87afd7
	SkyBlue2          Color = 0x87afff
	Chartreuse2_2     Color = 0x87d700
	DarkOliveGreen3_2 Color = 0x87d75f
	PaleGreen3_2      Color = 0x87d787
	DarkSeaGreen3     Color = 0x87d7af
	DarkSlateGray3    Color = 0x87d7d7
	SkyBlue1          Color = 0x87d7ff
	Chartreuse1       Color = 0x87ff00
	LightGreen        Color = 0x87ff5f
	LightGreen_2      Color = 0x87ff87
	PaleGreen1        Color = 0x87ffaf
	Aquamarine1_2     Color = 0x87ffd7
	DarkSlateGray1    Color = 0x87ffff
	Red3              Color = 0xaf0000
	DeepPink4         Color = 0xaf005f
	MediumVioletRed   Color = 0xaf0087
	Magenta3          Color = 0xaf00af
	DarkViolet_2      Color = 0xaf00d7
	Purple_2          Color = 0xaf00ff
	DarkOrange3       Color = 0xaf5f00
	IndianRed         Color = 0xaf5f5f
	HotPink3          Color = 0xaf5f87
	MediumOrchid3     Color = 0xaf5faf
	MediumOrchid      Color = 0xaf5fd7
	MediumPurple2     Color = 0xaf5fff
	DarkGoldenrod     Color = 0xaf8700
	LightSalmon3      Color = 0xaf875f
	RosyBrown         Color = 0xaf8787
	Grey63            Color = 0xaf87af
	MediumPurple2_2   Color = 0xaf87d7
	MediumPurple1     Color = 0xaf87ff
	Gold3             Color = 0xafaf00
	DarkKhaki         Color = 0xafaf5f
	NavajoWhite3      Color = 0xafaf87
	Grey69            Color = 0xafafaf
	LightSteelBlue3   Color = 0xafafd7
	LightSteelBlue    Color = 0xafafff
	Yellow3           Color = 0xafd700
	DarkOliveGreen3_3 Color = 0xafd75f
	DarkSeaGreen3_2   Color = 0xafd787
	DarkSeaGreen2     Color = 0xafd7af
	LightCyan3        Color = 0xafd7d7
	LightSkyBlue1     Color = 0xafd7ff
	GreenYellow       Color = 0xafff00
	DarkOliveGreen2   Color = 0xafff5f
	PaleGreen1_2      Color = 0xafff87
	DarkSeaGreen2_2   Color = 0xafffaf
	DarkSeaGreen1     Color = 0xafffd7
	PaleTurquoise1    Color = 0xafffff
	Red3_2            Color = 0xd70000
	DeepPink3         Color = 0xd7005f
	DeepPink3_2       Color = 0xd70087
	Magenta3_2        Color = 0xd700af
	Magenta3_3        Color = 0xd700d7
	Magenta2          Color = 0xd700ff
	DarkOrange3_2     Color = 0xd75f00
	IndianRed_2       Color = 0xd75f5f
	HotPink3_2        Color = 0xd75f87
	HotPink2          Color = 0xd75faf
	Orchid            Color = 0xd75fd7
	MediumOrchid1     Color = 0xd75fff
	Orange3           Color = 0xd78700
	LightSalmon3_2    Color = 0xd7875f
	LightPink3        Color = 0xd78787
	Pink3             Color = 0xd787af
	Plum3             Color = 0xd787d7
	Violet            Color = 0xd787ff
	Gold3_2           Color = 0xd7af00
	LightGoldenrod3   Color = 0xd7af5f
	Tan               Color = 0xd7af87
	MistyRose3        Color = 0xd7afaf
	Thistle3          Color = 0xd7afd7
	Plum2             Color = 0xd7afff
	Yellow3_2         Color = 0xd7d700
	Khaki3            Color = 0xd7d75f
	LightGoldenrod2   Color = 0xd7d787
	LightYellow3      Color = 0xd7d7af
	Grey84            Color = 0xd7d7d7
	LightSteelBlue1   Color = 0xd7d7ff
	Yellow2           Color = 0xd7ff00
	DarkOliveGreen1   Color = 0xd7ff5f
	DarkOliveGreen1_2 Color = 0xd7ff87
	DarkSeaGreen1_2   Color = 0xd7ffaf
	Honeydew2         Color = 0xd7ffd7
	LightCyan1        Color = 0xd7ffff
	Red1              Color = 0xff0000
	DeepPink2         Color = 0xff005f
	DeepPink1         Color = 0xff0087
	DeepPink1_2       Color = 0xff00af
	Magenta2_2        Color = 0xff00d7
	Magenta1          Color = 0xff00ff
	OrangeRed1        Color = 0xff5f00
	IndianRed1        Color = 0xff5f5f
	IndianRed1_2      Color = 0xff5f87
	HotPink           Color = 0xff5faf
	HotPink_2         Color = 0xff5fd7
	MediumOrchid1_2   Color = 0xff5fff
	DarkOrange        Color = 0xff8700
	Salmon1           Color = 0xff875f
	LightCoral        Color = 0xff8787
	PaleVioletRed1    Color = 0xff87af
	Orchid2           Color = 0xff87d7
	Orchid1           Color = 0xff87ff
	Orange1           Color = 0xffaf00
	SandyBrown        Color = 0xffaf5f
	LightSalmon1      Color = 0xffaf87
	LightPink1        Color = 0xffafaf
	Pink1             Color = 0xffafd7
	Plum1             Color = 0xffafff
	Gold1             Color = 0xffd700
	LightGoldenrod2_2 Color = 0xffd75f
	LightGoldenrod2_3 Color = 0xffd787
	NavajoWhite1      Color = 0xffd7af
	MistyRose1        Color = 0xffd7d7
	Thistle1          Color = 0xffd7ff
	Yellow1           Color = 0xffff00
	LightGoldenrod1   Color = 0xffff5f
	Khaki1            Color = 0xffff87
	Wheat1            Color = 0xffffaf
	Cornsilk1         Color = 0xffffd7
	Grey100           Color = 0xffffff
	Grey3             Color = 0x080808
	Grey7             Color = 0x121212
	Grey11            Color = 0x1c1c1c
	Grey15            Color = 0x262626
	Grey19            Color = 0x303030
	Grey23            Color = 0x3a3a3a
	Grey27            Color = 0x444444
	Grey30            Color = 0x4e4e4e
	Grey35            Color = 0x585858
	Grey39            Color = 0x626262
	Grey42            Color = 0x6c6c6c
	Grey46            Color = 0x767676
	Grey50            Color = 0x808080
	Grey54            Color = 0x8a8a8a
	Grey58            Color = 0x949494
	Grey62            Color = 0x9e9e9e
	Grey66            Color = 0xa8a8a8
	Grey70            Color = 0xb2b2b2
	Grey74            Color = 0xbcbcbc
	Grey78            Color = 0xc6c6c6
	Grey82            Color = 0xd0d0d0
	Grey85            Color = 0xdadada
	Grey89            Color = 0xe4e4e4
	Grey93            Color = 0xeeeeee

	DefaultColor Color = -1
)
