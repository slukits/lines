// Copyright (c) 2022 Stephan Lukits. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package term

import (
	"github.com/gdamore/tcell/v2"
	"github.com/slukits/lines/internal/api"
)

var apiToTcellCursorStyles = map[api.CursorStyle]tcell.CursorStyle{
	api.DefaultCursor:           tcell.CursorStyleDefault,
	api.BlockCursorBlinking:     tcell.CursorStyleBlinkingBlock,
	api.BlockCursorSteady:       tcell.CursorStyleSteadyBlock,
	api.UnderlineCursorBlinking: tcell.CursorStyleBlinkingUnderline,
	api.UnderlineCursorSteady:   tcell.CursorStyleSteadyUnderline,
	api.BarCursorBlinking:       tcell.CursorStyleBlinkingBar,
	api.BarCursorSteady:         tcell.CursorStyleSteadyBar,
}

var tcellToApiCursorStyles = map[tcell.CursorStyle]api.CursorStyle{
	tcell.CursorStyleDefault:           api.DefaultCursor,
	tcell.CursorStyleBlinkingBlock:     api.BlockCursorBlinking,
	tcell.CursorStyleSteadyBlock:       api.BarCursorSteady,
	tcell.CursorStyleBlinkingUnderline: api.UnderlineCursorBlinking,
	tcell.CursorStyleSteadyUnderline:   api.UnderlineCursorSteady,
	tcell.CursorStyleBlinkingBar:       api.BarCursorBlinking,
	tcell.CursorStyleSteadyBar:         api.BarCursorSteady,
}
