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
	t.Eq(1, len(MonoForeground(BlackM)))
	t.Eq(1, len(MonoForeground(WhiteM)))
}

func (s *Monochrome) Offers_white_foreground_to_black_background(t *T) {
	t.Eq(lines.White, MonoForeground(BlackM)[0].FG())
}

func (s *Monochrome) Offers_black_foreground_to_white_background(t *T) {
	t.Eq(lines.Black, MonoForeground(WhiteM)[0].FG())
}

func (s *Monochrome) Offers_white_background_to_black_foreground(t *T) {
	t.Eq(lines.White, MonoBackground(BlackM)[0].BG())
}

func (s *Monochrome) Offers_black_background_to_white_foreground(t *T) {
	t.Eq(lines.Black, MonoBackground(WhiteM)[0].BG())
}

func TestMonochrome(t *testing.T) {
	t.Parallel()
	Run(&Monochrome{}, t)
}

type system8 struct{ Suite }

func (s *system8) Has_eight_colors(t *T) {
	t.Eq(8, len(System8Colors))
}

func (s *system8) Provides_seven_foreground_combinations(t *T) {
	t.Eq(7, len(System8Foregrounds(Black8)))
}

func (s *system8) Background_is_not_in_provided_foregrounds(t *T) {
	ff := System8Foregrounds(Black8)
	for _, s := range ff {
		t.FatalIfNot(s.FG() != lines.Color(Black8))
	}
}

func (s *system8) Provides_seven_background_combinations(t *T) {
	t.Eq(7, len(System8Backgrounds(Black8)))
}

func (s *system8) Foreground_is_not_in_provided_backgrounds(t *T) {
	ff := System8Backgrounds(Black8)
	for _, s := range ff {
		t.FatalIfNot(s.BG() != lines.Color(Black8))
	}
}

func TestSystem8(t *testing.T) {
	t.Parallel()
	Run(&system8{}, t)
}
