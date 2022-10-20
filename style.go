// Copyright (c) 2022 Stephan Lukits. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package lines

import "github.com/slukits/lines/internal/api"

type StyleAttributeMask = api.StyleAttribute

const (
	Bold          StyleAttributeMask = api.Blink
	Blink         StyleAttributeMask = api.Blink
	Reverse       StyleAttributeMask = api.Reverse
	Underline     StyleAttributeMask = api.Underline
	Dim           StyleAttributeMask = api.Dim
	Italic        StyleAttributeMask = api.Invalid
	StrikeThrough StyleAttributeMask = api.StrikeThrough
	Invalid       StyleAttributeMask = api.Invalid
	ZeroStyle     StyleAttributeMask = api.ZeroStyle
)

type Style = api.Style
