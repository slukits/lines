// Copyright (c) 2022 Stephan Lukits. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package selects

import (
	"testing"

	. "github.com/slukits/gounit"
	"github.com/slukits/lines"
)

type monochrome struct{ Suite }

func (s *monochrome) SetUp(t *T) { t.Parallel() }

func (s *monochrome) Has_two_colors(t *T) { t.Eq(2, len(MonoColors)) }

func (s *monochrome) Offers_one_foreground_to_choose_from(t *T) {
	t.Eq(1, len(MonoForeground(BlackM)))
	t.Eq(1, len(MonoForeground(WhiteM)))
}

func (s *monochrome) Offers_white_foreground_to_black_background(t *T) {
	t.Eq(lines.White, MonoForeground(BlackM)[0].FG())
}

func (s *monochrome) Offers_black_foreground_to_white_background(t *T) {
	t.Eq(lines.Black, MonoForeground(WhiteM)[0].FG())
}

func (s *monochrome) Offers_white_background_to_black_foreground(t *T) {
	t.Eq(lines.White, MonoBackground(BlackM)[0].BG())
}

func (s *monochrome) Offers_black_background_to_white_foreground(t *T) {
	t.Eq(lines.Black, MonoBackground(WhiteM)[0].BG())
}

func TestMonochrome(t *testing.T) {
	t.Parallel()
	Run(&monochrome{}, t)
}

type system8 struct{ Suite }

func (s *system8) Has_eight_colors(t *T) {
	t.Eq(8, len(system8Colors))
}

func (s *system8) Provides_seven_foreground_combinations(t *T) {
	t.Eq(7, len(System8Foregrounds(lines.Black)))
}

func (s *system8) Provides_no_fg_combinations_if_bg_not_system8(t *T) {
	t.Eq(0, len(System8Foregrounds(lines.Salmon1)))
}

func (s *system8) Background_is_not_in_provided_foregrounds(t *T) {
	ff := System8Foregrounds(lines.Black)
	for _, s := range ff {
		t.FatalIfNot(s.FG() != lines.Color(lines.Black))
	}
}

func (s *system8) Provides_seven_background_combinations(t *T) {
	t.Eq(7, len(System8Backgrounds(lines.Black)))
}

func (s *system8) Provides_no_bg_combinations_if_fg_not_system8(t *T) {
	t.Eq(0, len(System8Backgrounds(lines.Salmon1)))
}

func (s *system8) Foreground_is_not_in_provided_backgrounds(t *T) {
	ff := System8Backgrounds(lines.Black)
	for _, s := range ff {
		t.FatalIfNot(s.BG() != lines.Color(lines.Black))
	}
}

func TestSystem8(t *testing.T) {
	t.Parallel()
	Run(&system8{}, t)
}

type linux struct{ Suite }

func (s *linux) Has_eight_background_colors(t *T) {
	t.Eq(8, len(linuxBGColors))
}

func (s *linux) Has_sixteen_foreground_colors(t *T) {
	t.Eq(16, len(linuxFGColors))
}

func (s *linux) Provides_seven_or_eight_backgrounds_to_given_fg(t *T) {
	eight := 0
	for _, fg := range linuxFGColors {
		bb := LinuxBackgrounds(fg)
		t.True(len(bb) == 7 || len(bb) == 8)
		if len(bb) == 7 {
			for _, b := range bb {
				t.True(b.BG() != lines.Color(fg))
			}
		}
		if len(bb) == 8 {
			eight++
			for _, b := range bb {
				if b.BG() == lines.Color(fg) {
					t.True(b.AA()|lines.Bold == lines.Bold)
				}
			}
		}
	}
	t.Eq(8, eight)
}

func (s *linux) Provides_fifteen_fg_colors_to_given_bg(t *T) {
	fgBright := 0
	for _, bg := range linuxBGColors {
		ff := LinuxForegrounds(bg)
		t.Eq(15, len(ff))
		for _, f := range ff {
			if f.FG() == lines.Color(bg) {
				fgBright++
				t.True(f.AA()|lines.Bold == lines.Bold)
			}
			t.True(f.BG() == lines.Color(bg))
		}
	}
	t.Eq(8, fgBright)
}

func TestLinux(t *testing.T) {
	t.Parallel()
	Run(&linux{}, t)
}

type system16 struct{ Suite }

func (s *system16) Has_sixteen_colors(t *T) {
	t.Eq(16, len(system16Colors))
}

func (s *system16) Provides_fifteen_foreground_combinations(t *T) {
	t.Eq(15, len(System16Foregrounds(lines.Black)))
}

func (s *system16) Background_is_not_in_provided_foregrounds(t *T) {
	ff, bg := System16Foregrounds(lines.Black), lines.Black
	for _, s := range ff {
		t.FatalIfNot(t.True(s.FG() != bg))
	}
}

func (s *system16) Provides_fifteen_background_combinations(t *T) {
	t.Eq(15, len(System16Backgrounds(lines.White)))
}

func (s *system16) Foreground_is_not_in_provided_backgrounds(t *T) {
	ff, fg := System16Backgrounds(lines.White), lines.White
	for _, s := range ff {
		t.FatalIfNot(t.True(s.BG() != fg))
	}
}

func TestSystem16(t *testing.T) {
	t.Parallel()
	Run(&system16{}, t)
}

type ansi struct{ Suite }

func (s *ansi) Has_256_colors(t *T) {
	t.Eq(256, len(ansiColors))
}

func (s *ansi) Provides_254_foreground_combinations_with_system_colors(
	t *T,
) {
	t.Eq(254, len(ANSIForegrounds(lines.Black)))
}

func (s *ansi) Provides_255_foreground_combinations(t *T) {
	t.Eq(255, len(ANSIForegrounds(lines.Tan)))
}

func (s *ansi) Background_is_not_in_provided_foregrounds(t *T) {
	ff, bg := ANSIForegrounds(lines.Black), lines.Black
	for _, s := range ff {
		t.FatalIfNot(t.True(s.FG() != bg))
	}
}

func (s *ansi) Provides_254_background_combinations_with_system_colors(
	t *T,
) {
	t.Eq(254, len(ANSIBackgrounds(lines.White)))
}

func (s *ansi) Provides_255_background_combinations(t *T) {
	t.Eq(255, len(ANSIForegrounds(lines.LightCoral)))
}

func (s *ansi) Foreground_is_not_in_provided_backgrounds(t *T) {
	ff, fg := ANSIBackgrounds(lines.White), lines.White
	for _, s := range ff {
		t.FatalIfNot(t.True(s.BG() != fg))
	}
}

func TestANSI(t *testing.T) {
	t.Parallel()
	Run(&ansi{}, t)
}
