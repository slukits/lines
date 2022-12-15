// Copyright (c) 2022 Stephan Lukits. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package lines

type EditType int

const (
	Ins EditType = iota
	Del
	Rpl
)

type Edit struct {
	Line int
	Cell int
	Type EditType
	Rune rune
}
