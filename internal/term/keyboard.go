// Copyright (c) 2022 Stephan Lukits. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package term

import (
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/slukits/lines/internal/api"
)

var apiToTcellKeys = map[api.Key]tcell.Key{
	api.NUL:       tcell.KeyNUL,
	api.SOH:       tcell.KeySOH,
	api.STX:       tcell.KeySTX,
	api.ETX:       tcell.KeyETX,
	api.EOT:       tcell.KeyEOT,
	api.ENQ:       tcell.KeyENQ,
	api.ACK:       tcell.KeyACK,
	api.BEL:       tcell.KeyBEL,
	api.BS:        tcell.KeyBS,
	api.TAB:       tcell.KeyTAB,
	api.LF:        tcell.KeyLF,
	api.VT:        tcell.KeyVT,
	api.FF:        tcell.KeyFF,
	api.CR:        tcell.KeyCR,
	api.SO:        tcell.KeySO,
	api.SI:        tcell.KeySI,
	api.DLE:       tcell.KeyDLE,
	api.DC1:       tcell.KeyDC1,
	api.DC2:       tcell.KeyDC2,
	api.DC3:       tcell.KeyDC3,
	api.DC4:       tcell.KeyDC4,
	api.NAK:       tcell.KeyNAK,
	api.SYN:       tcell.KeySYN,
	api.ETB:       tcell.KeyETB,
	api.CAN:       tcell.KeyCAN,
	api.EM:        tcell.KeyEM,
	api.SUB:       tcell.KeySUB,
	api.ESC:       tcell.KeyESC,
	api.FS:        tcell.KeyFS,
	api.GS:        tcell.KeyGS,
	api.RS:        tcell.KeyRS,
	api.US:        tcell.KeyUS,
	api.DEL:       tcell.KeyDEL,
	api.Up:        tcell.KeyUp,
	api.Down:      tcell.KeyDown,
	api.Right:     tcell.KeyRight,
	api.Left:      tcell.KeyLeft,
	api.UpLeft:    tcell.KeyUpLeft,
	api.UpRight:   tcell.KeyUpRight,
	api.DownLeft:  tcell.KeyDownLeft,
	api.DownRight: tcell.KeyDownRight,
	api.Center:    tcell.KeyCenter,
	api.PgUp:      tcell.KeyPgUp,
	api.PgDn:      tcell.KeyPgDn,
	api.Home:      tcell.KeyHome,
	api.End:       tcell.KeyEnd,
	api.Insert:    tcell.KeyInsert,
	api.Delete:    tcell.KeyDelete,
	api.Help:      tcell.KeyHelp,
	api.Exit:      tcell.KeyExit,
	api.Clear:     tcell.KeyClear,
	api.Cancel:    tcell.KeyCancel,
	api.Print:     tcell.KeyPrint,
	api.Pause:     tcell.KeyPause,
	api.Backtab:   tcell.KeyBacktab,
	api.F1:        tcell.KeyF1,
	api.F2:        tcell.KeyF2,
	api.F3:        tcell.KeyF3,
	api.F4:        tcell.KeyF4,
	api.F5:        tcell.KeyF5,
	api.F6:        tcell.KeyF6,
	api.F7:        tcell.KeyF7,
	api.F8:        tcell.KeyF8,
	api.F9:        tcell.KeyF9,
	api.F10:       tcell.KeyF10,
	api.F11:       tcell.KeyF11,
	api.F12:       tcell.KeyF12,
	api.F13:       tcell.KeyF13,
	api.F14:       tcell.KeyF14,
	api.F15:       tcell.KeyF15,
	api.F16:       tcell.KeyF16,
	api.F17:       tcell.KeyF17,
	api.F18:       tcell.KeyF18,
	api.F19:       tcell.KeyF19,
	api.F20:       tcell.KeyF20,
	api.F21:       tcell.KeyF21,
	api.F22:       tcell.KeyF22,
	api.F23:       tcell.KeyF23,
	api.F24:       tcell.KeyF24,
	api.F25:       tcell.KeyF25,
	api.F26:       tcell.KeyF26,
	api.F27:       tcell.KeyF27,
	api.F28:       tcell.KeyF28,
	api.F29:       tcell.KeyF29,
	api.F30:       tcell.KeyF30,
	api.F31:       tcell.KeyF31,
	api.F32:       tcell.KeyF32,
	api.F33:       tcell.KeyF33,
	api.F34:       tcell.KeyF34,
	api.F35:       tcell.KeyF35,
	api.F36:       tcell.KeyF36,
	api.F37:       tcell.KeyF37,
	api.F38:       tcell.KeyF38,
	api.F39:       tcell.KeyF39,
	api.F40:       tcell.KeyF40,
	api.F41:       tcell.KeyF41,
	api.F42:       tcell.KeyF42,
	api.F43:       tcell.KeyF43,
	api.F44:       tcell.KeyF44,
	api.F45:       tcell.KeyF45,
	api.F46:       tcell.KeyF46,
	api.F47:       tcell.KeyF47,
	api.F48:       tcell.KeyF48,
	api.F49:       tcell.KeyF49,
	api.F50:       tcell.KeyF50,
	api.F51:       tcell.KeyF51,
	api.F52:       tcell.KeyF52,
	api.F53:       tcell.KeyF53,
	api.F54:       tcell.KeyF54,
	api.F55:       tcell.KeyF55,
	api.F56:       tcell.KeyF56,
	api.F57:       tcell.KeyF57,
	api.F58:       tcell.KeyF58,
	api.F59:       tcell.KeyF59,
	api.F60:       tcell.KeyF60,
	api.F61:       tcell.KeyF61,
	api.F62:       tcell.KeyF62,
	api.F63:       tcell.KeyF63,
	api.F64:       tcell.KeyF64,
}

var tcellToApiKeys = map[tcell.Key]api.Key{
	tcell.KeyNUL:       api.NUL,
	tcell.KeySOH:       api.SOH,
	tcell.KeySTX:       api.STX,
	tcell.KeyETX:       api.ETX,
	tcell.KeyEOT:       api.EOT,
	tcell.KeyENQ:       api.ENQ,
	tcell.KeyACK:       api.ACK,
	tcell.KeyBEL:       api.BEL,
	tcell.KeyBS:        api.BS,
	tcell.KeyTAB:       api.TAB,
	tcell.KeyLF:        api.LF,
	tcell.KeyVT:        api.VT,
	tcell.KeyFF:        api.FF,
	tcell.KeyCR:        api.CR,
	tcell.KeySO:        api.SO,
	tcell.KeySI:        api.SI,
	tcell.KeyDLE:       api.DLE,
	tcell.KeyDC1:       api.DC1,
	tcell.KeyDC2:       api.DC2,
	tcell.KeyDC3:       api.DC3,
	tcell.KeyDC4:       api.DC4,
	tcell.KeyNAK:       api.NAK,
	tcell.KeySYN:       api.SYN,
	tcell.KeyETB:       api.ETB,
	tcell.KeyCAN:       api.CAN,
	tcell.KeyEM:        api.EM,
	tcell.KeySUB:       api.SUB,
	tcell.KeyESC:       api.ESC,
	tcell.KeyFS:        api.FS,
	tcell.KeyGS:        api.GS,
	tcell.KeyRS:        api.RS,
	tcell.KeyUS:        api.US,
	tcell.KeyDEL:       api.DEL,
	tcell.KeyUp:        api.Up,
	tcell.KeyDown:      api.Down,
	tcell.KeyRight:     api.Right,
	tcell.KeyLeft:      api.Left,
	tcell.KeyUpLeft:    api.UpLeft,
	tcell.KeyUpRight:   api.UpRight,
	tcell.KeyDownLeft:  api.DownLeft,
	tcell.KeyDownRight: api.DownRight,
	tcell.KeyCenter:    api.Center,
	tcell.KeyPgUp:      api.PgUp,
	tcell.KeyPgDn:      api.PgDn,
	tcell.KeyHome:      api.Home,
	tcell.KeyEnd:       api.End,
	tcell.KeyInsert:    api.Insert,
	tcell.KeyDelete:    api.Delete,
	tcell.KeyHelp:      api.Help,
	tcell.KeyExit:      api.Exit,
	tcell.KeyClear:     api.Clear,
	tcell.KeyCancel:    api.Cancel,
	tcell.KeyPrint:     api.Print,
	tcell.KeyPause:     api.Pause,
	tcell.KeyBacktab:   api.Backtab,
	tcell.KeyF1:        api.F1,
	tcell.KeyF2:        api.F2,
	tcell.KeyF3:        api.F3,
	tcell.KeyF4:        api.F4,
	tcell.KeyF5:        api.F5,
	tcell.KeyF6:        api.F6,
	tcell.KeyF7:        api.F7,
	tcell.KeyF8:        api.F8,
	tcell.KeyF9:        api.F9,
	tcell.KeyF10:       api.F10,
	tcell.KeyF11:       api.F11,
	tcell.KeyF12:       api.F12,
	tcell.KeyF13:       api.F13,
	tcell.KeyF14:       api.F14,
	tcell.KeyF15:       api.F15,
	tcell.KeyF16:       api.F16,
	tcell.KeyF17:       api.F17,
	tcell.KeyF18:       api.F18,
	tcell.KeyF19:       api.F19,
	tcell.KeyF20:       api.F20,
	tcell.KeyF21:       api.F21,
	tcell.KeyF22:       api.F22,
	tcell.KeyF23:       api.F23,
	tcell.KeyF24:       api.F24,
	tcell.KeyF25:       api.F25,
	tcell.KeyF26:       api.F26,
	tcell.KeyF27:       api.F27,
	tcell.KeyF28:       api.F28,
	tcell.KeyF29:       api.F29,
	tcell.KeyF30:       api.F30,
	tcell.KeyF31:       api.F31,
	tcell.KeyF32:       api.F32,
	tcell.KeyF33:       api.F33,
	tcell.KeyF34:       api.F34,
	tcell.KeyF35:       api.F35,
	tcell.KeyF36:       api.F36,
	tcell.KeyF37:       api.F37,
	tcell.KeyF38:       api.F38,
	tcell.KeyF39:       api.F39,
	tcell.KeyF40:       api.F40,
	tcell.KeyF41:       api.F41,
	tcell.KeyF42:       api.F42,
	tcell.KeyF43:       api.F43,
	tcell.KeyF44:       api.F44,
	tcell.KeyF45:       api.F45,
	tcell.KeyF46:       api.F46,
	tcell.KeyF47:       api.F47,
	tcell.KeyF48:       api.F48,
	tcell.KeyF49:       api.F49,
	tcell.KeyF50:       api.F50,
	tcell.KeyF51:       api.F51,
	tcell.KeyF52:       api.F52,
	tcell.KeyF53:       api.F53,
	tcell.KeyF54:       api.F54,
	tcell.KeyF55:       api.F55,
	tcell.KeyF56:       api.F56,
	tcell.KeyF57:       api.F57,
	tcell.KeyF58:       api.F58,
	tcell.KeyF59:       api.F59,
	tcell.KeyF60:       api.F60,
	tcell.KeyF61:       api.F61,
	tcell.KeyF62:       api.F62,
	tcell.KeyF63:       api.F63,
	tcell.KeyF64:       api.F64,
}

var apiToTcellModifiers = map[api.ModifierMask]tcell.ModMask{
	api.Shift:        tcell.ModShift,
	api.Ctrl:         tcell.ModCtrl,
	api.Alt:          tcell.ModAlt,
	api.Meta:         tcell.ModMeta,
	api.ZeroModifier: tcell.ModNone,
}

var apiModifiers = []api.ModifierMask{
	api.Shift, api.Ctrl, api.Alt, api.Meta, api.ZeroModifier,
}

func apiModifiersToTcell(mm api.ModifierMask) (tm tcell.ModMask) {
	for _, m := range apiModifiers {
		if mm&m == 0 {
			continue
		}
		tm |= apiToTcellModifiers[m]
	}
	return tm
}

var tcellToApiModifiers = map[tcell.ModMask]api.ModifierMask{
	tcell.ModShift: api.Shift,
	tcell.ModCtrl:  api.Ctrl,
	tcell.ModAlt:   api.Alt,
	tcell.ModMeta:  api.Meta,
	tcell.ModNone:  api.ZeroModifier,
}

var tcellModifiers = []tcell.ModMask{
	tcell.ModShift, tcell.ModCtrl, tcell.ModAlt, tcell.ModMeta,
	tcell.ModNone,
}

func tcellModifiersToApi(mm tcell.ModMask) (am api.ModifierMask) {
	for _, m := range tcellModifiers {
		if mm&m == 0 {
			continue
		}
		am |= tcellToApiModifiers[m]
	}
	return am
}

type keyEvent struct{ evt *tcell.EventKey }

func newKeyEvent(k api.Key, m api.ModifierMask) api.KeyEventer {
	return &keyEvent{evt: tcell.NewEventKey(
		apiToTcellKeys[k],
		rune(tcell.KeyRune),
		apiModifiersToTcell(m),
	)}
}

func (e *keyEvent) Key() api.Key {
	return tcellToApiKeys[e.evt.Key()]
}

func (e *keyEvent) Mod() api.ModifierMask {
	return tcellModifiersToApi(e.evt.Modifiers())
}

func (e *keyEvent) When() time.Time { return e.evt.When() }

func (e *keyEvent) Source() interface{} { return e.evt }

type runeEvent struct{ evt *tcell.EventKey }

func newRuneEvent(r rune, m api.ModifierMask) api.RuneEventer {
	return &runeEvent{evt: tcell.NewEventKey(
		tcell.KeyNUL, r, apiModifiersToTcell(m),
	)}
}

func (e *runeEvent) Rune() rune { return e.evt.Rune() }

func (e *runeEvent) Mod() api.ModifierMask {
	return tcellModifiersToApi(e.evt.Modifiers())
}

func (e *runeEvent) When() time.Time { return e.evt.When() }

func (e *runeEvent) Source() interface{} { return e.evt }
