// Copyright (c) 2022 Stephan Lukits. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package lines

import (
	"fmt"
	"testing"

	. "github.com/slukits/gounit"
)

type Line struct{ Suite }

func (s *Line) Is_filled_at_line_fillers(t *T) {
	fx := &icmpFX{init: func(c *icmpFX, e *Env) {
		c.Dim().SetHeight(1).SetWidth(8)
		fmt.Fprintf(e, "a%sb", LineFiller)
	}}
	ee, tt := Test(t.GoT(), fx, 3)
	ee.Listen()

	t.Eq("a      b", tt.Screen().String())

	ee.Update(fx, nil, func(e *Env) {
		fmt.Fprintf(e, "a%sb%[1]sc", LineFiller)
	})

	t.Eq("a   b  c", tt.Screen().String())

	ee.Update(fx, nil, func(e *Env) {
		fmt.Fprintf(e, "ab%scd%[1]sef%[1]sgh", LineFiller)
	})

	t.Eq("ab cd ef", tt.LastScreen.String())
}

func TestLineRun(t *testing.T) {
	Run(&Line{}, t)
}
