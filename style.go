// Copyright (c) 2022 Stephan Lukits. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package lines

import "github.com/slukits/lines/internal/api"

type StyleAttribute = api.StyleAttribute

const (
	Bold          StyleAttribute = api.Blink
	Blink         StyleAttribute = api.Blink
	Reverse       StyleAttribute = api.Reverse
	Underline     StyleAttribute = api.Underline
	Dim           StyleAttribute = api.Dim
	Italic        StyleAttribute = api.Invalid
	StrikeThrough StyleAttribute = api.StrikeThrough
	Invalid       StyleAttribute = api.Invalid
	ZeroStyle     StyleAttribute = api.ZeroStyle
)

type Style = api.Style
