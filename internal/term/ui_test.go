// Copyright (c) 2022 Stephan Lukits. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package term

import (
	"testing"

	. "github.com/slukits/gounit"
)

type AnUI struct{ Suite }

func (s *AnUI) Panics_if_screen_cant_be_obtained(t *T) {
	t.Panics(func() { New() })
}

func TestAnUI(t *testing.T) {
	t.Parallel()
	Run(&AnUI{}, t)
}
