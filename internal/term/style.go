// Copyright (c) 2022 Stephan Lukits. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package term

import (
	"github.com/gdamore/tcell/v2"
	"github.com/slukits/lines/internal/api"
)

func apiToTcellStyle(s api.Style) tcell.Style {
	return tcell.StyleDefault.Background(tcell.Color(s.BG)).
		Foreground(tcell.Color(s.FG)).
		Attributes(tcell.AttrMask(s.AA))
}

// tcellToApiStyleClosure keeps the last tcell-style conversion and
// returns it until provided tcell style changes.
func tcellToApiStyleClosure() func(tcell.Style) api.Style {
	sty := tcell.StyleDefault
	fg, bg, aa := sty.Decompose()
	apiSty := api.Style{
		AA: api.StyleAttribute(aa),
		FG: api.Color(fg),
		BG: api.Color(bg),
	}
	return func(s tcell.Style) api.Style {
		_fg, _bg, _aa := s.Decompose()
		if _fg == fg && _bg == bg && _aa == aa {
			return apiSty
		}
		fg, bg, aa = _fg, _bg, _aa
		apiSty = api.Style{
			AA: api.StyleAttribute(aa),
			FG: api.Color(fg),
			BG: api.Color(bg),
		}
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
