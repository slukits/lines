// Copyright (c) 2022 Stephan Lukits. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package colors

import (
	"testing"

	. "github.com/slukits/gounit"
	"github.com/slukits/lines"
)

type Monochrome struct{ Suite }

func (s *Monochrome) SetUp(t *T) { t.Parallel() }

func (s *Monochrome) Has_two_colors(t *T) { t.Eq(2, len(MonoColors)) }

func (s *Monochrome) Offers_one_foreground_to_choose_from(t *T) {
	t.Eq(1, len(MonoForeground(Black)))
	t.Eq(1, len(MonoForeground(White)))
}

func (s *Monochrome) Offers_white_foreground_to_black_background(t *T) {
	t.Eq(lines.White, MonoForeground(Black)[0].FG())
}

func (s *Monochrome) Offers_black_foreground_to_white_background(t *T) {
	t.Eq(lines.Black, MonoForeground(White)[0].FG())
}

func (s *Monochrome) Offers_white_background_to_black_foreground(t *T) {
	t.Eq(lines.White, MonoBackground(Black)[0].BG())
}

func (s *Monochrome) Offers_black_background_to_white_foreground(t *T) {
	t.Eq(lines.Black, MonoBackground(White)[0].BG())
}

func TestMonochrome(t *testing.T) {
	t.Parallel()
	Run(&Monochrome{}, t)
}
