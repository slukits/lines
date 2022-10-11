// Copyright (c) 2022 Stephan Lukits. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package lines

import "github.com/slukits/lines/internal/api"

// A KeyEventer is implemented by a reported key-event.
type KeyEventer = api.KeyEventer

// A RuneEventer is implemented by a reported rune-event.
type RuneEventer = api.RuneEventer

// A Key is the pressed key of a key event.
type Key = api.Key

const (
	NUL            Key = api.NUL
	SOH            Key = api.SOH
	STX            Key = api.STX
	ETX            Key = api.ETX
	EOT            Key = api.EOT
	ENQ            Key = api.ENQ
	ACK            Key = api.ACK
	BEL            Key = api.BEL
	BS             Key = api.BS
	TAB            Key = api.TAB
	LF             Key = api.LF
	VT             Key = api.VT
	FF             Key = api.FF
	CR             Key = api.CR
	SO             Key = api.SO
	SI             Key = api.SI
	DLE            Key = api.DLE
	DC1            Key = api.DC1
	DC2            Key = api.DC2
	DC3            Key = api.DC3
	DC4            Key = api.DC4
	NAK            Key = api.NAK
	SYN            Key = api.SYN
	ETB            Key = api.ETB
	CAN            Key = api.CAN
	EM             Key = api.EM
	SUB            Key = api.SUB
	ESC            Key = api.ESC
	FS             Key = api.FS
	GS             Key = api.GS
	RS             Key = api.RS
	US             Key = api.US
	DEL            Key = api.DEL
	Up             Key = api.Up
	Down           Key = api.Down
	Right          Key = api.Right
	Left           Key = api.Left
	UpLeft         Key = api.UpLeft
	UpRight        Key = api.UpRight
	DownLeft       Key = api.DownLeft
	DownRight      Key = api.DownRight
	Center         Key = api.Center
	PgUp           Key = api.PgUp
	PgDn           Key = api.PgDn
	Home           Key = api.Home
	End            Key = api.End
	Insert         Key = api.Insert
	Delete         Key = api.Delete
	Help           Key = api.Help
	Exit           Key = api.Exit
	Clear          Key = api.Clear
	Cancel         Key = api.Cancel
	Print          Key = api.Print
	Pause          Key = api.Pause
	Backtab        Key = api.Backtab
	F1             Key = api.F1
	F2             Key = api.F2
	F3             Key = api.F3
	F4             Key = api.F4
	F5             Key = api.F5
	F6             Key = api.F6
	F7             Key = api.F7
	F8             Key = api.F8
	F9             Key = api.F9
	F10            Key = api.F10
	F11            Key = api.F11
	F12            Key = api.F12
	F13            Key = api.F13
	F14            Key = api.F14
	F15            Key = api.F15
	F16            Key = api.F16
	F17            Key = api.F17
	F18            Key = api.F18
	F19            Key = api.F19
	F20            Key = api.F20
	F21            Key = api.F21
	F22            Key = api.F22
	F23            Key = api.F23
	F24            Key = api.F24
	F25            Key = api.F25
	F26            Key = api.F26
	F27            Key = api.F27
	F28            Key = api.F28
	F29            Key = api.F29
	F30            Key = api.F30
	F31            Key = api.F31
	F32            Key = api.F32
	F33            Key = api.F33
	F34            Key = api.F34
	F35            Key = api.F35
	F36            Key = api.F36
	F37            Key = api.F37
	F38            Key = api.F38
	F39            Key = api.F39
	F40            Key = api.F40
	F41            Key = api.F41
	F42            Key = api.F42
	F43            Key = api.F43
	F44            Key = api.F44
	F45            Key = api.F45
	F46            Key = api.F46
	F47            Key = api.F47
	F48            Key = api.F48
	F49            Key = api.F49
	F50            Key = api.F50
	F51            Key = api.F51
	F52            Key = api.F52
	F53            Key = api.F53
	F54            Key = api.F54
	F55            Key = api.F55
	F56            Key = api.F56
	F57            Key = api.F57
	F58            Key = api.F58
	F59            Key = api.F59
	F60            Key = api.F60
	F61            Key = api.F61
	F62            Key = api.F62
	F63            Key = api.F63
	F64            Key = api.F64
	CtrlSpace      Key = api.CtrlSpace
	CtrlA          Key = api.CtrlA
	CtrlB          Key = api.CtrlB
	CtrlC          Key = api.CtrlC
	CtrlD          Key = api.CtrlD
	CtrlE          Key = api.CtrlE
	CtrlF          Key = api.CtrlF
	CtrlG          Key = api.CtrlG
	CtrlH          Key = api.CtrlH
	CtrlI          Key = api.CtrlI
	CtrlJ          Key = api.CtrlJ
	CtrlK          Key = api.CtrlK
	CtrlL          Key = api.CtrlL
	CtrlM          Key = api.CtrlM
	CtrlN          Key = api.CtrlN
	CtrlO          Key = api.CtrlO
	CtrlP          Key = api.CtrlP
	CtrlQ          Key = api.CtrlQ
	CtrlR          Key = api.CtrlR
	CtrlS          Key = api.CtrlS
	CtrlT          Key = api.CtrlT
	CtrlU          Key = api.CtrlU
	CtrlV          Key = api.CtrlV
	CtrlW          Key = api.CtrlW
	CtrlX          Key = api.CtrlX
	CtrlY          Key = api.CtrlY
	CtrlZ          Key = api.CtrlZ
	CtrlLeftSq     Key = api.CtrlLeftSq
	CtrlBackslash  Key = api.CtrlBackslash
	CtrlRightSq    Key = api.CtrlRightSq
	CtrlCarat      Key = api.CtrlCarat
	CtrlUnderscore Key = api.CtrlUnderscore
)

// A Modifier mask are the pressed modifier keys of a key, rune or mouse
// event.  Note that the shift modifier of a rune event is not reported.
type Modifier = api.Modifier

const (
	Shift        Modifier = api.Shift
	Ctrl         Modifier = api.Ctrl
	Alt          Modifier = api.Alt
	Meta         Modifier = api.Meta
	ZeroModifier Modifier = api.ZeroModifier
)
