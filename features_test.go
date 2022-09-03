// Copyright (c) 2022 Stephan Lukits. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package lines

import (
	"testing"

	"github.com/gdamore/tcell/v2"
	. "github.com/slukits/gounit"
)

var allFX = func() FeatureMask {
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

func (c *cmpFFX) OnInit(*Env) { c.test(c.FF) }

type _Features struct{ Suite }

func (s *_Features) SetUp(t *T) { t.Parallel() }

func (s *_Features) Panic_outside_event_listener_callback(t *T) {
	fx := &cmpFFX{test: func(*Features) {}}
	ee, _ := Test(t.GoT(), fx)
	ee.Listen()
	t.Panics(func() { fx.FF.Has(Quitable) })
}

func (s *_Features) Has_by_default_only_quitable_registered(t *T) {
	fx := &cmpFFX{test: func(ff *Features) {
		_ff := NoFeature
		for _, f := range allFeatures {
			if ff.Has(f) {
				_ff |= f
			}
		}
		t.Eq(ff.All(), _ff)
		t.Eq(Quitable, _ff)
	}}
	ee, _ := Test(t.GoT(), fx)
	ee.Listen()
}

func (s *_Features) Ignores_deletion_of_default_quitable_keys(t *T) {
	fx := &cmpFFX{test: func(ff *Features) {
		ff.ensureInitialized()
		exp := defaultFeatures.keysOf(Quitable)
		got := ff.KeysOf(Quitable)
		t.True(exp.Equals(got))
		ff.SetKeysOf(Quitable, false)
		got = ff.KeysOf(Quitable)
		t.True(exp.Equals(got))
		ff.Delete(Quitable)
		got = ff.KeysOf(Quitable)
		t.True(exp.Equals(got))
	}}
	ee, _ := Test(t.GoT(), fx)
	ee.Listen()
}

func (s *_Features) Deletes_user_added_quitable_keys(t *T) {
	fx := &cmpFFX{test: func(ff *Features) {
		ff.ensureInitialized()
		exp := defaultFeatures.keysOf(Quitable)
		fx := FeatureKey{Key: tcell.KeyCtrlX, Mod: 0}
		ff.SetKeysOf(Quitable, false, fx)
		t.True(ff.KeysOf(Quitable).Equals(append(exp, fx)))
		ff.Delete(Quitable)
		t.True(exp.Equals(ff.KeysOf(Quitable)))
	}}
	ee, _ := Test(t.GoT(), fx)
	ee.Listen()
}

func (s *_Features) Deletes_default_quitable_rune(t *T) {
	fx := &cmpFFX{test: func(ff *Features) {
		ff.ensureInitialized()
		ff.SetRunesOf(Quitable, false)
		t.True(len(ff.RunesOf(Quitable)) == 0)
	}}
	ee, _ := Test(t.GoT(), fx)
	ee.Listen()
}

func (s *_Features) Has_set_features(t *T) {
	fx := &cmpFFX{test: func(ff *Features) {
		ff.Add(Focusable)
		ff.Add(Selectable)
		t.True(ff.Has(Focusable))
		t.True(ff.Has(Selectable))
		t.True(ff.Has(PreviousSelectable))
		t.True(ff.Has(NextSelectable))
	}}
	ee, _ := Test(t.GoT(), fx)
	ee.Listen()
}

func (s *_Features) Sets_defaults_bindings_of_feature(t *T) {
	fx := &cmpFFX{test: func(ff *Features) {
		ff.Add(Focusable)
		ff.Add(Selectable)
		t.True(defaultBindings[Focusable].bb.Equals(
			ff.ButtonsOf(Focusable)))
		t.True(defaultBindings[NextSelectable].kk.Equals(
			ff.KeysOf(NextSelectable)))
		t.True(defaultBindings[PreviousSelectable].kk.Equals(
			ff.KeysOf(PreviousSelectable)))
	}}
	ee, _ := Test(t.GoT(), fx)
	ee.Listen()
}

func (s *_Features) Have_set_runes(t *T) {
	fx := &cmpFFX{test: func(ff *Features) {
		ff.SetRunesOf(Focusable, false, 'n', 'm')
		t.True(ff.RunesOf(Focusable).Equals(
			FeatureRunes{'n', 'm'}))
	}}
	ee, _ := Test(t.GoT(), fx)
	ee.Listen()
}

func (s *_Features) Replaces_rune_bindings_with_set_runes(t *T) {
	fx := &cmpFFX{test: func(ff *Features) {
		ff.SetRunesOf(Focusable, false, 'n', 'm')
		t.True(ff.RunesOf(Focusable).Equals(FeatureRunes{'n', 'm'}))
		ff.SetRunesOf(Focusable, false, 'x')
		t.True(ff.RunesOf(Focusable).Equals(FeatureRunes{'x'}))
	}}
	ee, _ := Test(t.GoT(), fx)
	ee.Listen()
}

func (s *_Features) Have_set_keys(t *T) {
	kk := FeatureKeys{{Key: tcell.KeyBackspace, Mod: tcell.ModNone},
		{Key: tcell.KeyTAB, Mod: tcell.ModAlt}}
	fx := &cmpFFX{test: func(ff *Features) {
		ff.SetKeysOf(Focusable, false, kk...)
		t.True(ff.KeysOf(Focusable).Equals(kk))
	}}
	ee, _ := Test(t.GoT(), fx)
	ee.Listen()
}

func (s *_Features) Replaces_key_bindings_with_set_keys(t *T) {
	kk := FeatureKeys{{Key: tcell.KeyBackspace, Mod: tcell.ModNone},
		{Key: tcell.KeyTAB, Mod: tcell.ModAlt}}
	exp := FeatureKeys{{Key: tcell.KeyBacktab, Mod: tcell.ModAlt}}
	fx := &cmpFFX{test: func(ff *Features) {
		ff.SetKeysOf(Focusable, false, kk...)
		t.True(ff.KeysOf(Focusable).Equals(kk))
		ff.SetKeysOf(Focusable, false, exp...)
		t.True(ff.KeysOf(Focusable).Equals(exp))
	}}
	ee, _ := Test(t.GoT(), fx)
	ee.Listen()
}

func (s *_Features) Have_set_buttons(t *T) {
	bb := FeatureButtons{
		{Button: tcell.ButtonPrimary, Mod: tcell.ModNone},
		{Button: tcell.ButtonMiddle, Mod: tcell.ModShift}}
	fx := &cmpFFX{test: func(ff *Features) {
		ff.SetButtonsOf(Focusable, false, bb...)
		t.True(ff.ButtonsOf(Focusable).Equals(bb))
	}}
	ee, _ := Test(t.GoT(), fx)
	ee.Listen()
}

func (s *_Features) Replaces_button_bindings_with_set_buttons(t *T) {
	bb := FeatureButtons{
		{Button: tcell.ButtonPrimary, Mod: tcell.ModNone},
		{Button: tcell.ButtonMiddle, Mod: tcell.ModShift}}
	exp := FeatureButtons{
		{Button: tcell.ButtonSecondary, Mod: tcell.ModAlt}}
	fx := &cmpFFX{test: func(ff *Features) {
		ff.SetButtonsOf(Focusable, false, bb...)
		t.True(ff.ButtonsOf(Focusable).Equals(bb))
		ff.SetButtonsOf(Focusable, false, exp...)
		t.True(ff.ButtonsOf(Focusable).Equals(exp))
	}}
	ee, _ := Test(t.GoT(), fx)
	ee.Listen()
}

func (s *_Features) Have_recursively_defined_features(t *T) {
	bttFX := FeatureButtons{
		{Button: tcell.ButtonSecondary, Mod: tcell.ModAlt}}
	fx := &cmpFFX{test: func(ff *Features) {
		ff.SetButtonsOf(Focusable, true, bttFX...)
		t.True(ff.Has(Focusable | _recursive))
	}}
	ee, _ := Test(t.GoT(), fx)
	ee.Listen()
}

func (s *_Features) Reports_rune_bindings_of_recursive_feature(t *T) {
	exp := FeatureRunes{'f'}
	fx := &cmpFFX{test: func(ff *Features) {
		ff.SetRunesOf(Focusable, true, exp...)
		t.True(ff.RunesOf(Focusable | _recursive).Equals(exp))
		t.True(ff.RunesOf(Focusable).Equals(exp))
	}}
	ee, _ := Test(t.GoT(), fx)
	ee.Listen()
}

func (s *_Features) Reports_key_bindings_of_recursive_feature(t *T) {
	exp := FeatureKeys{{Key: tcell.KeyBacktab, Mod: tcell.ModAlt}}
	fx := &cmpFFX{test: func(ff *Features) {
		ff.SetKeysOf(Focusable, true, exp...)
		t.True(ff.KeysOf(Focusable | _recursive).Equals(exp))
		t.True(ff.KeysOf(Focusable).Equals(exp))
	}}
	ee, _ := Test(t.GoT(), fx)
	ee.Listen()
}

func (s *_Features) Reports_button_bindings_of_recursive_feature(t *T) {
	exp := FeatureButtons{
		{Button: tcell.ButtonSecondary, Mod: tcell.ModAlt}}
	fx := &cmpFFX{test: func(ff *Features) {
		ff.SetButtonsOf(Focusable, true, exp...)
		t.True(ff.ButtonsOf(Focusable | _recursive).Equals(exp))
		t.True(ff.ButtonsOf(Focusable).Equals(exp))
	}}
	ee, _ := Test(t.GoT(), fx)
	ee.Listen()
}

func (s *_Features) Ignore_setting_runes_for_comprised_features(t *T) {
	fx := &cmpFFX{test: func(ff *Features) {
		ff.SetRunesOf(Selectable, false, 's')
		t.False(ff.Has(Selectable))
	}}
	ee, _ := Test(t.GoT(), fx)
	ee.Listen()
}

func (s *_Features) Ignore_setting_keys_for_comprised_features(t *T) {
	keyFX := FeatureKeys{{Key: tcell.KeyBacktab, Mod: tcell.ModAlt}}
	fx := &cmpFFX{test: func(ff *Features) {
		ff.SetKeysOf(Selectable, false, keyFX...)
		t.False(ff.Has(Selectable))
	}}
	ee, _ := Test(t.GoT(), fx)
	ee.Listen()
}

func (s *_Features) Ignore_setting_buttons_for_comprised_features(t *T) {
	bttFX := FeatureButtons{
		{Button: tcell.ButtonSecondary, Mod: tcell.ModAlt}}
	fx := &cmpFFX{test: func(ff *Features) {
		ff.SetButtonsOf(Selectable, false, bttFX...)
		t.False(ff.Has(Selectable))
	}}
	ee, _ := Test(t.GoT(), fx)
	ee.Listen()
}

func TestFeatures(t *testing.T) {
	t.Parallel()
	Run(&_Features{}, t)
}
