// Copyright (c) 2022 Stephan Lukits. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package lines

import "github.com/slukits/lines/internal/api"

type CursorStyle = api.CursorStyle

const (
	ZeroCursor              CursorStyle = api.ZeroCursor
	DefaultCursor                       = api.DefaultCursor
	BlockCursorBlinking                 = api.BlockCursorBlinking
	BlockCursorSteady                   = api.BlockCursorSteady
	UnderlineCursorBlinking             = api.UnderlineCursorBlinking
	UnderlineCursorSteady               = api.UnderlineCursorSteady
	BarCursorBlinking                   = api.BarCursorBlinking
	BarCursorSteady                     = api.BarCursorSteady
)
