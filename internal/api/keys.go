// Copyright (c) 2022 Stephan Lukits. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package api

import "github.com/gdamore/tcell/v2"

type Key int32

const (
	NUL Key = iota
	SOH
	STX
	ETX
	EOT
	ENQ
	ACK
	BEL
	BS
	TAB
	LF
	VT
	FF
	CR
	SO
	SI
	DLE
	DC1
	DC2
	DC3
	DC4
	NAK
	SYN
	ETB
	CAN
	EM
	SUB
	ESC
	FS
	GS
	RS
	US
	DEL
	Up
	Down
	Right
	Left
	UpLeft
	UpRight
	DownLeft
	DownRight
	Center
	PgUp
	PgDn
	Home
	End
	Insert
	Delete
	Help
	Exit
	Clear
	Cancel
	Print
	Pause
	Backtab
	F1
	F2
	F3
	F4
	F5
	F6
	F7
	F8
	F9
	F10
	F11
	F12
	F13
	F14
	F15
	F16
	F17
	F18
	F19
	F20
	F21
	F22
	F23
	F24
	F25
	F26
	F27
	F28
	F29
	F30
	F31
	F32
	F33
	F34
	F35
	F36
	F37
	F38
	F39
	F40
	F41
	F42
	F43
	F44
	F45
	F46
	F47
	F48
	F49
	F50
	F51
	F52
	F53
	F54
	F55
	F56
	F57
	F58
	F59
	F60
	F61
	F62
	F63
	F64

	CtrlSpace      = NUL
	CtrlA          = SOH
	CtrlB          = STX
	CtrlC          = ETX
	CtrlD          = EOT
	CtrlE          = ENQ
	CtrlF          = ACK
	CtrlG          = BEL
	CtrlH          = BS
	CtrlI          = TAB
	CtrlJ          = LF
	CtrlK          = VT
	CtrlL          = FF
	CtrlM          = CR
	CtrlN          = SO
	CtrlO          = SI
	CtrlP          = DLE
	CtrlQ          = DC1
	CtrlR          = DC2
	CtrlS          = DC3
	CtrlT          = DC4
	CtrlU          = NAK
	CtrlV          = SYN
	CtrlW          = ETB
	CtrlX          = CAN
	CtrlY          = EM
	CtrlZ          = SUB
	CtrlLeftSq     = ESC
	CtrlBackslash  = FS
	CtrlRightSq    = GS
	CtrlCarat      = RS
	CtrlUnderscore = US
	Enter          = CR
	Tab            = TAB
	Backspace      = BS
	Esc            = ESC
)

const k = tcell.KeyBEL

// ModifierMask mask is used to provide pressed modifiers of a
// reported/posted key/rune user input event.
type ModifierMask int32

const (
	Shift ModifierMask = 1 << iota
	Ctrl
	Alt
	Meta
	ZeroModifier ModifierMask = 0
)
