// Copyright (c) 2022 Stephan Lukits. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package term

import (
	"github.com/gdamore/tcell/v2"
	"github.com/slukits/lines/internal/api"
)

func apiToTcellStyle(s api.Style) tcell.Style {
	fg, bg := tcell.ColorDefault, tcell.ColorDefault
	if s.FG() != api.DefaultColor {
		fg = tcell.NewHexColor(int32(s.FG()))
	}
	if s.BG() != api.DefaultColor {
		bg = tcell.NewHexColor(int32(s.BG()))
	}
	return tcell.StyleDefault.
		Foreground(fg).Background(bg).
		Attributes(tcell.AttrMask(s.AA()))
}

func tcellToApiStyle(s tcell.Style) api.Style {
	fg, bg := api.DefaultColor, api.DefaultColor
	tfg, tbg, taa := s.Decompose()
	if tfg != tcell.ColorDefault {
		fg = api.Color(tfg.Hex())
	}
	if tbg != tcell.ColorDefault {
		bg = api.Color(tbg.Hex())
	}
	return api.NewStyle(api.StyleAttribute(taa), fg, bg)
}

// tcellToApiStyleClosure keeps the last tcell-style conversion and
// returns it until provided tcell style changes.
func tcellToApiStyleClosure() func(tcell.Style) api.Style {
	sty := tcell.StyleDefault
	tfg, tbg, aa := sty.Decompose()
	apiSty := tcellToApiStyle(sty)
	return func(s tcell.Style) api.Style {
		_fg, _bg, _aa := s.Decompose()
		if _fg == tfg && _bg == tbg && _aa == aa {
			return apiSty
		}
		tfg, tbg, aa = _fg, _bg, _aa
		apiSty = tcellToApiStyle(s)
		return apiSty
	}
}

// apiToTcellStyleClosure keeps the last api-style conversion and
// returns it until provided api-style changes.
func apiToTcellStyleClosure() func(api.Style) tcell.Style {
	apiSty := api.Style{}
	tclSty := apiToTcellStyle(apiSty)
	return func(s api.Style) tcell.Style {
		if apiSty.Equals(s) {
			return tclSty
		}
		apiSty = s
		tclSty = apiToTcellStyle(s)
		return tclSty
	}
}
