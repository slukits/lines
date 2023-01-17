// Copyright (c) 2022 Stephan Lukits. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package lines

import (
	"testing"

	. "github.com/slukits/gounit"
)

type _Features struct{ Suite }

func (s *_Features) SetUp(t *T) { t.Parallel() }

func (s *_Features) Panic_outside_event_listener_callback(t *T) {
	fx := &cmpFX{}
	TermFixture(t.GoT(), 0, fx)
	t.Panics(func() { fx.FF.Has(Focusable) })
}

func (s *_Features) Have_initially_no_features_registered(t *T) {
	fx, cmp := fxCmp(t)
	fx.Lines.Update(cmp, nil, func(e *Env) {
		t.Eq(NoFeature, cmp.FF.All())
	})
}

func (s *_Features) Are_inherited_from_parent_component(t *T) {
	cmp, chn, inner := &stackingFX{}, &chainingFX{}, &cmpFX{}
	cmp.CC = append(cmp.CC, chn)
	chn.CC = append(chn.CC, inner)
	cmp.onInit = func(c *cmpFX, e *Env) {
		c.FF.Set(Focusable)
	}
	fx := fx(t, cmp)
	fx.Lines.Update(inner, nil, func(e *Env) {
		t.True(inner.FF.Has(Focusable))
	})
}

func (s *_Features) Have_set_features(t *T) {
	fx, cmp := fxCmp(t)
	fx.Lines.Update(cmp, nil, func(e *Env) {
		cmp.FF.Set(Focusable)
		t.True(cmp.FF.Has(Focusable))
		cmp.FF.Set(Scrollable)
		t.True(cmp.FF.Has(Scrollable))
		t.True(cmp.FF.Has(UpScrollable))
		t.True(cmp.FF.Has(DownScrollable))
	})
}

func (s *_Features) Adding_sets_defaults_bindings_of_features(t *T) {
	fx, cmp := fxCmp(t)
	fx.Lines.Update(cmp, nil, func(e *Env) {
		cmp.FF.Set(Focusable)
		cmp.FF.Set(Scrollable)
		t.True(defaultBindings[Focusable].bb.Equals(
			cmp.FF.ButtonsOf(Focusable)))
		t.True(defaultBindings[UpScrollable].kk.Equals(
			cmp.FF.KeysOf(UpScrollable)))
		t.True(defaultBindings[DownScrollable].kk.Equals(
			cmp.FF.KeysOf(DownScrollable)))
	})
}

func (s *_Features) Have_set_runes(t *T) {
	fx, cmp := fxCmp(t)
	fx.Lines.Update(cmp, nil, func(e *Env) {
		cmp.FF.SetRunesOf(Focusable,
			FeatureRunes{{Rune: 'n'}, {Rune: 'm', Mod: Alt}}...)
		t.True(cmp.FF.RunesOf(Focusable).Equals(
			FeatureRunes{{Rune: 'n'}, {Rune: 'm', Mod: Alt}}))
		t.Eq(cmp.FF.OfRune('m', Alt), Focusable)
	})
}

func (s *_Features) Replace_default_bindings_with_set_runes(t *T) {
	fx, cmp := fxCmp(t)
	fx.Lines.Update(cmp, nil, func(e *Env) {
		cmp.FF.SetRunesOf(Focusable,
			FeatureRunes{{Rune: 'n'}, {Rune: 'm', Mod: Alt}}...)
		t.True(cmp.FF.RunesOf(Focusable).Equals(FeatureRunes{
			{Rune: 'n'}, {Rune: 'm'}}))
		cmp.FF.SetRunesOf(Focusable, FeatureRune{Rune: 'x'})
		t.True(cmp.FF.RunesOf(Focusable).Equals(
			FeatureRunes{{Rune: 'x'}}))
	})
}

func (s *_Features) Have_set_keys(t *T) {
	kk := FeatureKeys{{Key: Backspace}, {Key: TAB, Mod: Alt}}
	fx, cmp := fxCmp(t)
	fx.Lines.Update(cmp, nil, func(e *Env) {
		cmp.FF.SetKeysOf(Focusable, kk...)
		t.True(cmp.FF.KeysOf(Focusable).Equals(kk))
		t.Eq(cmp.FF.OfKey(TAB, Alt), Focusable)
	})
}

func (s *_Features) Replace_default_bindings_with_set_keys(t *T) {
	kk := FeatureKeys{{Key: Backspace}, {Key: TAB, Mod: Alt}}
	exp := FeatureKeys{{Key: Backtab, Mod: Alt}}
	fx, cmp := fxCmp(t)
	fx.Lines.Update(cmp, nil, func(e *Env) {
		cmp.FF.SetKeysOf(Focusable, kk...)
		t.True(cmp.FF.KeysOf(Focusable).Equals(kk))
		cmp.FF.SetKeysOf(Focusable, exp...)
		t.True(cmp.FF.KeysOf(Focusable).Equals(exp))
	})
}

func (s *_Features) Have_set_buttons(t *T) {
	bb := FeatureButtons{
		{Button: Primary}, {Button: Middle, Mod: Shift}}
	fx, cmp := fxCmp(t)
	fx.Lines.Update(cmp, nil, func(e *Env) {
		cmp.FF.SetButtonsOf(Focusable, bb...)
		t.True(cmp.FF.ButtonsOf(Focusable).Equals(bb))
		t.Eq(cmp.FF.OfButton(Middle, Shift), Focusable)
	})
}

func (s *_Features) Replace_default_bindings_with_set_buttons(t *T) {
	bb := FeatureButtons{
		{Button: Primary}, {Button: Middle, Mod: Shift}}
	exp := FeatureButtons{{Button: Secondary, Mod: Alt}}
	fx, cmp := fxCmp(t)
	fx.Lines.Update(cmp, nil, func(e *Env) {
		cmp.FF.SetButtonsOf(Focusable, bb...)
		t.True(cmp.FF.ButtonsOf(Focusable).Equals(bb))
		cmp.FF.SetButtonsOf(Focusable, exp...)
		t.True(cmp.FF.ButtonsOf(Focusable).Equals(exp))
	})
}

func (s *_Features) Have_feature_of_set_binding(t *T) {
	fx, cmp := fxCmp(t)
	fx.Lines.Update(cmp, nil, func(e *Env) {
		cmp.FF.SetRunesOf(Focusable, FeatureRune{Rune: 'n'})
		t.True(cmp.FF.RunesOf(Focusable).Equals(
			FeatureRunes{{Rune: 'n'}}))
		cmp.FF.SetKeysOf(UpScrollable, FeatureKey{Key: PgUp})
		t.True(cmp.FF.KeysOf(UpScrollable).Equals(
			FeatureKeys{{Key: PgUp}}))
		cmp.FF.SetButtonsOf(DownScrollable, FeatureButton{Button: Button5})
		t.True(cmp.FF.ButtonsOf(DownScrollable).Equals(
			FeatureButtons{{Button: Button5}}))
		cmp.FF.Has(Focusable)
		cmp.FF.Has(Scrollable)
	})
}

func (s *_Features) Ignores_setting_runes_for_comprised_features(t *T) {
	fx, cmp := fxCmp(t)
	fx.Lines.Update(cmp, nil, func(e *Env) {
		cmp.FF.SetRunesOf(Scrollable, FeatureRunes{{Rune: 's'}}...)
		t.Not.True(cmp.FF.Has(Scrollable))
	})
}

func (s *_Features) Ignores_setting_keys_for_comprised_features(t *T) {
	keyFX := FeatureKeys{{Key: Backtab, Mod: Alt}}
	fx, cmp := fxCmp(t)
	fx.Lines.Update(cmp, nil, func(e *Env) {
		cmp.FF.SetKeysOf(Scrollable, keyFX...)
		t.Not.True(cmp.FF.Has(Scrollable))
	})
}

func (s *_Features) Ignores_setting_buttons_for_comprised_features(t *T) {
	bttFX := FeatureButtons{{Button: Secondary, Mod: Alt}}
	fx, cmp := fxCmp(t)
	fx.Lines.Update(cmp, nil, func(e *Env) {
		cmp.FF.SetButtonsOf(Scrollable, bttFX...)
		t.Not.True(cmp.FF.Has(Scrollable))
	})
}

func (s *_Features) Deletes_all_bindings_of_given_feature(t *T) {
	fx, cmp := fxCmp(t)
	fx.Lines.Update(cmp, nil, func(e *Env) {
		cmp.FF.SetRunesOf(Focusable, FeatureRune{Rune: 'n'})
		t.True(cmp.FF.RunesOf(Focusable).Equals(
			FeatureRunes{{Rune: 'n'}}))
		cmp.FF.SetKeysOf(Focusable, FeatureKey{Key: Enter})
		t.True(cmp.FF.KeysOf(Focusable).Equals(
			FeatureKeys{{Key: Enter}}))
		cmp.FF.SetButtonsOf(Focusable, FeatureButton{Button: Button1})
		t.True(cmp.FF.ButtonsOf(Focusable).Equals(
			FeatureButtons{{Button: Button1}}))
		cmp.FF.Delete(Focusable)
		t.True(cmp.FF.RunesOf(Focusable).Equals(FeatureRunes{}))
		t.True(cmp.FF.KeysOf(Focusable).Equals(FeatureKeys{}))
		t.True(cmp.FF.ButtonsOf(Focusable).Equals(FeatureButtons{}))
	})
}

func TestFeatures(t *testing.T) {
	t.Parallel()
	Run(&_Features{}, t)
}
