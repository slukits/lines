// Copyright (c) 2022 Stephan Lukits. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

/*
ui implements an UIer wrapping tcell for lines terminal ui.
*/

package term

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
)

type UI struct {
	lib tcell.Screen
}

func New() *UI {
	lib, err := tcell.NewScreen()
	if err != nil {
		panic(fmt.Sprintf(
			"lines: term: new: can't obtain screen: %v", err))
	}
	if err := lib.Init(); err != nil {
		panic(fmt.Sprintf(
			"lines: term: new: can't obtain screen: %v", err))
	}
	return &UI{lib: lib}
}
