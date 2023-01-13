// Copyright (c) 2022 Stephan Lukits. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package lines

import (
	"testing"

	. "github.com/slukits/gounit"
)

var allFMFX = func() FeatureMask {
	ff := NoFeature
	for _, f := range allFeatures {
		ff |= f
	}
	return ff
}()

type cmpFFX struct {
	Component
	test func(*Features)
}

// OnInit runs the test since we can't access features outside an
// event callback.
func (c *cmpFFX) OnInit(*Env) { c.test(c.FF) }

type _Features struct{ Suite }

func (s *_Features) SetUp(t *T) { t.Parallel() }

func (s *_Features) Panic_outside_event_listener_callback(t *T) {
	fx := &cmpFFX{test: func(*Features) {}}
	TermFixture(t.GoT(), 0, fx)
	t.Panics(func() { fx.FF.Has(Focusable) })
}

func (s *_Features) tt(t *T, test func(ff *Features)) (
	*Fixture, *cmpFFX,
) {
	fx := &cmpFFX{test: func(*Features) {}}
	tt := TermFixture(t.GoT(), 0, fx)
	return tt, fx
}

func (s *_Features) Has_by_default_no_features_registered(t *T) {
	s.tt(t, func(ff *Features) {
		_ff := NoFeature
		for _, f := range allFeatures {
			if ff.Has(f) {
				_ff |= f
			}
		}
		t.Eq(ff.All(), _ff)
		t.Eq(NoFeature, _ff)
	})
}

// func (s *_Features) Ignores_deletion_of_default_quitable_keys(t *T) {
// 	s.tt(t, func(ff *Features) {
// 		ff.ensureInitialized()
// 		exp := quitableFeatures.keysOf(Quitable)
// 		got := ff.KeysOf(Quitable)
// 		t.True(exp.Equals(got))
// 		ff.SetKeysOf(Quitable, false)
// 		got = ff.KeysOf(Quitable)
// 		t.True(exp.Equals(got))
// 		ff.Delete(Quitable)
// 		got = ff.KeysOf(Quitable)
// 		t.True(exp.Equals(got))
// 	})
// }
//
// func (s *_Features) Delete_user_added_quitable_keys(t *T) {
// 	s.tt(t, func(ff *Features) {
// 		ff.ensureInitialized()
// 		exp := quitableFeatures.keysOf(Quitable)
// 		fx := FeatureKey{Key: CtrlX, Mod: 0}
// 		ff.SetKeysOf(Quitable, false, fx)
// 		t.True(ff.KeysOf(Quitable).Equals(append(exp, fx)))
// 		ff.Delete(Quitable)
// 		t.True(exp.Equals(ff.KeysOf(Quitable)))
// 	})
// }

func (s *_Features) Are_inherited_from_nesting_features(t *T) {
	cmp, chn, inner := &stackingFX{}, &chainingFX{}, &cmpFX{}
	cmp.CC = append(cmp.CC, chn)
	chn.CC = append(chn.CC, inner)
	cmp.onInit = func(c *cmpFX, e *Env) {
		c.FF.Add(Focusable)
	}
	fx := fx(t, cmp)
	fx.Lines.Update(inner, nil, func(e *Env) {
		t.True(inner.FF.Has(Focusable))
	})
}

func (s *_Features) Have_set_features(t *T) {
	s.tt(t, func(ff *Features) {
		ff.Add(Focusable)
		t.True(ff.Has(Focusable))
		ff.Add(Scrollable)
		t.True(ff.Has(Scrollable))
		t.True(ff.Has(UpScrollable))
		t.True(ff.Has(DownScrollable))
	})
}

func (s *_Features) Sets_defaults_bindings_of_feature(t *T) {
	s.tt(t, func(ff *Features) {
		ff.Add(Focusable)
		ff.Add(Scrollable)
		t.True(defaultBindings[Focusable].bb.Equals(
			ff.ButtonsOf(Focusable)))
		t.True(defaultBindings[UpScrollable].kk.Equals(
			ff.KeysOf(UpScrollable)))
		t.True(defaultBindings[DownScrollable].kk.Equals(
			ff.KeysOf(DownScrollable)))
	})
}

func (s *_Features) Have_set_runes(t *T) {
	s.tt(t, func(ff *Features) {
		ff.SetRunesOf(Focusable,
			FeatureRunes{{Rune: 'n'}, {Rune: 'm', Mod: Alt}}...)
		t.True(ff.RunesOf(Focusable).Equals(
			FeatureRunes{{Rune: 'n'}, {Rune: 'm', Mod: Alt}}))
	})
}

func (s *_Features) Replaces_rune_bindings_with_set_runes(t *T) {
	s.tt(t, func(ff *Features) {
		ff.SetRunesOf(Focusable,
			FeatureRunes{{Rune: 'n'}, {Rune: 'm', Mod: Alt}}...)
		t.True(ff.RunesOf(Focusable).Equals(FeatureRunes{
			{Rune: 'n'}, {Rune: 'm'}}))
		ff.SetRunesOf(Focusable, FeatureRune{Rune: 'x'})
		t.True(ff.RunesOf(Focusable).Equals(
			FeatureRunes{{Rune: 'x'}}))
	})
}

func (s *_Features) Have_set_keys(t *T) {
	kk := FeatureKeys{{Key: Backspace}, {Key: TAB, Mod: Alt}}
	s.tt(t, func(ff *Features) {
		ff.SetKeysOf(Focusable, kk...)
		t.True(ff.KeysOf(Focusable).Equals(kk))
	})
}

func (s *_Features) Replaces_key_bindings_with_set_keys(t *T) {
	kk := FeatureKeys{{Key: Backspace}, {Key: TAB, Mod: Alt}}
	exp := FeatureKeys{{Key: Backtab, Mod: Alt}}
	s.tt(t, func(ff *Features) {
		ff.SetKeysOf(Focusable, kk...)
		t.True(ff.KeysOf(Focusable).Equals(kk))
		ff.SetKeysOf(Focusable, exp...)
		t.True(ff.KeysOf(Focusable).Equals(exp))
	})
}

func (s *_Features) Have_set_buttons(t *T) {
	bb := FeatureButtons{
		{Button: Primary}, {Button: Middle, Mod: Shift}}
	s.tt(t, func(ff *Features) {
		ff.SetButtonsOf(Focusable, bb...)
		t.True(ff.ButtonsOf(Focusable).Equals(bb))
	})
}

func (s *_Features) Replaces_button_bindings_with_set_buttons(t *T) {
	bb := FeatureButtons{
		{Button: Primary}, {Button: Middle, Mod: Shift}}
	exp := FeatureButtons{{Button: Secondary, Mod: Alt}}
	s.tt(t, func(ff *Features) {
		ff.SetButtonsOf(Focusable, bb...)
		t.True(ff.ButtonsOf(Focusable).Equals(bb))
		ff.SetButtonsOf(Focusable, exp...)
		t.True(ff.ButtonsOf(Focusable).Equals(exp))
	})
}

// func (s *_Features) Have_recursively_defined_features(t *T) {
// 	bttFX := FeatureButtons{{Button: Secondary, Mod: Alt}}
// 	s.tt(t, func(ff *Features) {
// 		ff.SetButtonsOf(Focusable, true, bttFX...)
// 		t.True(ff.Has(Focusable | _recursive))
// 	})
// }
//
// func (s *_Features) Reports_rune_bindings_of_recursive_feature(t *T) {
// 	exp := FeatureRunes{{Rune: 'f'}}
// 	s.tt(t, func(ff *Features) {
// 		ff.SetRunesOf(Focusable, true, exp...)
// 		t.True(ff.RunesOf(Focusable | _recursive).Equals(exp))
// 		t.True(ff.RunesOf(Focusable).Equals(exp))
// 	})
// }
//
// func (s *_Features) Reports_key_bindings_of_recursive_feature(t *T) {
// 	exp := FeatureKeys{{Key: Backtab, Mod: Alt}}
// 	s.tt(t, func(ff *Features) {
// 		ff.SetKeysOf(Focusable, true, exp...)
// 		t.True(ff.KeysOf(Focusable | _recursive).Equals(exp))
// 		t.True(ff.KeysOf(Focusable).Equals(exp))
// 	})
// }
//
// func (s *_Features) Reports_button_bindings_of_recursive_feature(t *T) {
// 	exp := FeatureButtons{{Button: Secondary, Mod: Alt}}
// 	s.tt(t, func(ff *Features) {
// 		ff.SetButtonsOf(Focusable, true, exp...)
// 		t.True(ff.ButtonsOf(Focusable | _recursive).Equals(exp))
// 		t.True(ff.ButtonsOf(Focusable).Equals(exp))
// 	})
// }

func (s *_Features) Ignores_setting_runes_for_comprised_features(t *T) {
	s.tt(t, func(ff *Features) {
		ff.SetRunesOf(Scrollable, FeatureRunes{{Rune: 's'}}...)
		t.Not.True(ff.Has(Scrollable))
	})
}

func (s *_Features) Ignores_setting_keys_for_comprised_features(t *T) {
	keyFX := FeatureKeys{{Key: Backtab, Mod: Alt}}
	s.tt(t, func(ff *Features) {
		ff.SetKeysOf(Scrollable, keyFX...)
		t.Not.True(ff.Has(Scrollable))
	})
}

func (s *_Features) Ignores_setting_buttons_for_comprised_features(t *T) {
	bttFX := FeatureButtons{{Button: Secondary, Mod: Alt}}
	s.tt(t, func(ff *Features) {
		ff.SetButtonsOf(Scrollable, bttFX...)
		t.Not.True(ff.Has(Scrollable))
	})
}

func TestFeatures(t *testing.T) {
	t.Parallel()
	Run(&_Features{}, t)
}
