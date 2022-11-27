// Copyright (c) 2022 Stephan Lukits. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"

	"github.com/slukits/lines"
)

type Runes struct {
	lines.Component
}

const (
	waiting  = "waiting for rune"
	received = "received rune: %c"
)

func (c *Runes) OnInit(e *lines.Env) {
	c.Dim().SetHeight(1).SetWidth(len(waiting))
	fmt.Fprint(e, waiting)
}

func (c *Runes) OnRune(e *lines.Env, r rune, _ lines.ModifierMask) {
	fmt.Fprintf(e, received, r)
}

func main() {
	lines.Term(&Runes{}).WaitForQuit()
}
