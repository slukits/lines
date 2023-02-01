// Copyright (c) 2022 Stephan Lukits. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package selects

import (
	"fmt"

	"github.com/slukits/lines"
)

type label struct {
	component
	lbl string
}

func (c *label) OnInit(e *lines.Env) { fmt.Fprint(e, c.lbl) }
func (c *label) width() int          { return len(c.lbl) + 1 }
