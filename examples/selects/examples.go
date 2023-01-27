// Copyright (c) 2022 Stephan Lukits. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

import (
	"github.com/slukits/lines"
	"github.com/slukits/lines/cmp/selects"
)

type empty struct{ selects.List }

func (c *empty) OnInit(e lines.Env) { c.SetDirty() }

func emptyList() lines.Componenter {
	return &empty{}
}
