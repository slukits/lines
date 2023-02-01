// Copyright (c) 2022 Stephan Lukits. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

/*
lines' features functionality is defined in three files:
cmpfeatures.go: the public/client API to the features functionality.
features.go: this file providing the underlying internal api.
report_features.go: where features are executed and consequences of this
executions are reported.
*/

package lines

// features provides information about keys/runes/buttons which are
// registered for features provided by the lines-package.  It also
// allows to change these in a consistent and convenient way.  The zero
// value is ready to use.
type features struct {
	keys    map[ModifierMask]map[Key]FeatureMask
	runes   map[ModifierMask]map[rune]FeatureMask
	buttons map[ModifierMask]map[ButtonMask]FeatureMask
	have    FeatureMask
}

// copy creates a new Features instance initialized with the features of
// receiving Features instance.
func (ff *features) copy() *features {
	if ff == nil {
		return nil
	}

	cpy := features{
		keys:    map[ModifierMask]map[Key]FeatureMask{},
		runes:   map[ModifierMask]map[rune]FeatureMask{},
		buttons: map[ModifierMask]map[ButtonMask]FeatureMask{},
		have:    ff.have,
	}

	for m, kk := range ff.keys {
		cpy.keys[m] = map[Key]FeatureMask{}
		for k, f := range kk {
			cpy.keys[m][k] = f
		}
	}

	for m, rr := range ff.runes {
		cpy.runes[m] = map[rune]FeatureMask{}
		for r, f := range rr {
			cpy.runes[m][r] = f
		}
	}

	for m, bb := range ff.buttons {
		cpy.buttons[m] = map[ButtonMask]FeatureMask{}
		for b, f := range bb {
			cpy.buttons[m][b] = f
		}
	}

	return &cpy
}

func (ff *features) has(f FeatureMask) bool {
	if ff == nil {
		return false
	}
	return ff.have&f == f
}

// all returns the set of currently set features.
func (ff *features) all() FeatureMask {
	if ff == nil {
		return NoFeature
	}
	return ff.have
}

// set adds the default bindings of given feature to receiving component
// features.
func (ff *features) set(f FeatureMask) {

	if f == NoFeature || ff == nil {
		return
	}

	_ff := []FeatureMask{}
	for _, _f := range allFeatures {
		if _f&f == NoFeature {
			continue
		}
		_ff = append(_ff, _f)
	}
	if len(_ff) == 0 {
		return
	}
	ff.ensureConsistency(f)
	for _, f := range _ff {
		df := defaultBindings[f]
		if df == nil {
			continue
		}
		ff.ensureInitForDefaults(df)
		for _, k := range df.kk {
			if ff.keys[k.Mod] == nil {
				ff.keys[k.Mod] = map[Key]FeatureMask{}
			}
			ff.keys[k.Mod][k.Key] = f
		}
		for _, b := range df.bb {
			if ff.buttons[b.Mod] == nil {
				ff.buttons[b.Mod] = map[ButtonMask]FeatureMask{}
			}
			ff.buttons[b.Mod][b.Button] = f
		}
		for _, r := range df.rr {
			if ff.runes[r.Mod] == nil {
				ff.runes[r.Mod] = map[rune]FeatureMask{}
			}
			ff.runes[r.Mod][r.Rune] = f
		}
	}

	if ff.have&f != f {
		ff.have |= f
	}
}

func (ff *features) ensureConsistency(f FeatureMask) {
	if f&HighlightEnabled != 0 && ff.have&TrimmedHighlightEnabled != 0 {
		ff.have &^= TrimmedHighlightEnabled
	}
	if f&TrimmedHighlightEnabled != 0 && ff.have&HighlightEnabled != 0 {
		ff.have &^= HighlightEnabled
	}
}

func (ff *features) ensureInitForDefaults(bb *bindings) {
	if bb.kk != nil && ff.keys == nil {
		ff.keys = map[ModifierMask]map[Key]FeatureMask{}
	}
	if bb.rr != nil && ff.runes == nil {
		ff.runes = map[ModifierMask]map[rune]FeatureMask{}
	}
	if bb.bb != nil && ff.buttons == nil {
		ff.buttons = map[ModifierMask]map[ButtonMask]FeatureMask{}
	}
}

// FeatureKey represents a key with its modifier and key value
// which is typically bound to a feature using [Features.SetKeysOf].
type FeatureKey struct {
	Mod ModifierMask
	Key Key
}

// FeatureKeys are reported by [Features.KeysOf] providing the keys
// bound to a given feature.  FeatureKeys may be also used as variadic
// argument for [Features.SetKeysOf] to bind several [FeatureKey] to the
// same feature.
type FeatureKeys []FeatureKey

// Equals returns true iff given feature keys fk and given other feature
// keys contain the same FeatureKey values.
func (fk FeatureKeys) Equals(other FeatureKeys) bool {
	if len(fk) != len(other) {
		return false
	}
	for _, k := range fk {
		has := false
		for _, o := range other {
			if o.Key != k.Key || o.Mod != k.Mod {
				continue
			}
			has = true
		}
		if !has {
			return false
		}
	}
	return true
}

// keysOf returns the keys with their modifiers for given feature.
func (ff *features) keysOf(f FeatureMask) FeatureKeys {
	kk := []FeatureKey{}
	if ff == nil {
		return kk
	}
	for m, _kk := range ff.keys {
		for k, _f := range _kk {
			if f&_f == NoFeature {
				continue
			}
			kk = append(kk, FeatureKey{Mod: m, Key: k})
		}
	}
	return kk
}

// setKeysOf deletes all set keys for given feature and binds given keys
// to it instead.  The call is ignored if given feature is not a power
// of two i.e. a single feature.  NOTE providing no keys simply removes
// all key-bindings for given feature.
func (ff *features) setKeysOf(f FeatureMask, kk ...FeatureKey) {
	if ff == nil || f == 0 || f&(f-1) != 0 { // f is not a power of two
		return
	}
	ff.deleteKeysOfButDefaults(f)
	if ff.keys == nil {
		ff.keys = map[ModifierMask]map[Key]FeatureMask{}
	}
	for _, k := range kk {
		if ff.keys[k.Mod] == nil {
			ff.keys[k.Mod] = map[Key]FeatureMask{}
		}
		ff.keys[k.Mod][k.Key] = f
	}

	if ff.have&f == NoFeature {
		ff.have |= f
	}
}

// FeatureButton represents a button with its modifier and button value
// which is typically bound to a feature using [Features.SetButtonsOf].
type FeatureButton struct {
	Mod    ModifierMask
	Button ButtonMask
}

// FeatureButtons are reported by [Features.ButtonsOf] providing the
// mouse buttons bound to a given feature.  FeatureButtons may be also
// used as variadic argument of [Features.SetButtonsOf] to bind several
// [FeatureButton] to the same feature.
type FeatureButtons []FeatureButton

// Equals returns true if given feature buttons fb and given other
// feature buttons contain the same feature buttons.
func (fb FeatureButtons) Equals(other FeatureButtons) bool {
	if len(fb) != len(other) {
		return false
	}
	for _, b := range fb {
		has := false
		for _, o := range other {
			if o.Button != b.Button || o.Mod != b.Mod {
				continue
			}
			has = true
		}
		if !has {
			return false
		}
	}
	return true
}

// buttonsOf returns the buttons with their modifiers for given feature.
func (ff *features) buttonsOf(f FeatureMask) FeatureButtons {
	bb := []FeatureButton{}
	if ff == nil {
		return bb
	}
	for m, _bb := range ff.buttons {
		for b, _f := range _bb {
			if f&_f == NoFeature {
				continue
			}
			bb = append(bb, FeatureButton{Mod: m, Button: b})
		}
	}
	return bb
}

// setButtonsOf deletes all set buttons for given feature an binds given
// buttons to it instead.  If recursive is true the feature becomes
// applicable for nested components.  The call is ignored if given
// feature is not a power of two i.e. a single feature.  NOTE providing
// no buttons simply removes all button-bindings for given feature.
func (ff *features) setButtonsOf(f FeatureMask, bb ...FeatureButton) {

	if ff == nil || f == 0 || f&(f-1) != 0 { // f is not a power of two
		return
	}

	ff.deleteButtonsOf(f)
	if ff.buttons == nil {
		ff.buttons = map[ModifierMask]map[ButtonMask]FeatureMask{}
	}
	for _, b := range bb {
		if ff.buttons[b.Mod] == nil {
			ff.buttons[b.Mod] = map[ButtonMask]FeatureMask{}
		}
		ff.buttons[b.Mod][b.Button] = f
	}

	if ff.have&f == NoFeature {
		ff.have |= f
	}
}

// FeatureRune represents a rune with its modifier and rune value
// which is typically bound to a feature using [Features.SetRunesOf].
type FeatureRune struct {
	Rune rune
	Mod  ModifierMask
}

// FeatureRunes are reported by [Features.RunesOf] providing the runes
// bound to a given feature.  FeatureRunes may be also used as variadic
// argument of [Features.SetRunesOf] to bind several [FeatureRune]
// to the same feature.
type FeatureRunes []FeatureRune

// Equals returns true if given feature runes fr and given other
// feature runes contain the same feature runes.
func (fr FeatureRunes) Equals(other FeatureRunes) bool {
	if len(fr) != len(other) {
		return false
	}
	for _, r := range fr {
		has := false
		for _, o := range other {
			if r.Rune != o.Rune {
				continue
			}
			has = true
		}
		if !has {
			return false
		}
	}
	return true
}

// runesOf returns the runes for given lines-feature.
func (ff *features) runesOf(f FeatureMask) FeatureRunes {
	fr := FeatureRunes{}
	if ff == nil {
		return fr
	}
	for m, rr := range ff.runes {
		for r, _f := range rr {
			if f&_f == NoFeature {
				continue
			}
			fr = append(fr, FeatureRune{Rune: r, Mod: m})
		}
	}
	return fr
}

// setRunesOf deletes all set runes for given feature an binds given
// runes to it instead.  The call is ignored if given feature is not a
// power of two i.e. a single feature.  NOTE providing no runes simply
// removes all runes-bindings for given feature.
func (ff *features) setRunesOf(f FeatureMask, rr ...FeatureRune) {
	if ff == nil || f == 0 || f&(f-1) != 0 { // f is not a power of two
		return
	}

	ff.deleteRunesOf(f)
	if ff.runes == nil {
		ff.runes = map[ModifierMask]map[rune]FeatureMask{}
	}
	for _, r := range rr {
		if ff.runes[r.Mod] == nil {
			ff.runes[r.Mod] = map[rune]FeatureMask{}
		}
		ff.runes[r.Mod][r.Rune] = f
	}

	if ff.have&f == NoFeature {
		ff.have |= f
	}
}

// delete removes all runes, key or button bindings of given feature(s).
func (ff *features) delete(f FeatureMask) {
	if ff == nil || ff.have&f == NoFeature {
		return
	}
	_ff := []FeatureMask{}
	for _, _f := range allFeatures {
		if _f&f == NoFeature {
			continue
		}
		_ff = append(_ff, _f)
	}
	if len(_ff) == 0 {
		return
	}

	for _, f := range _ff {
		ff.deleteKeysOf(f)
		ff.deleteButtonsOf(f)
		ff.deleteRunesOf(f)
		ff.have &^= f
	}
	ff.have &^= f
}

func (ff *features) deleteKeysOf(f FeatureMask) {
	for m, kk := range ff.keys {
		for k, _f := range kk {
			if f&_f == NoFeature {
				continue
			}
			delete(ff.keys[m], k)
		}
	}
}

func (ff *features) deleteKeysOfButDefaults(f FeatureMask) {
	if ff == nil {
		return
	}
	for m, kk := range ff.keys {
		for k, _f := range kk {
			if f&_f == NoFeature {
				continue
			}
			delete(ff.keys[m], k)
		}
	}
}

func (ff *features) deleteButtonsOf(f FeatureMask) {
	if ff == nil {
		return
	}
	for m, bb := range ff.buttons {
		for b, _f := range bb {
			if f&_f == NoFeature {
				continue
			}
			delete(ff.buttons[m], b)
		}
	}
}

func (ff *features) deleteRunesOf(f FeatureMask) {
	if ff == nil {
		return
	}
	for m, rr := range ff.runes {
		for r, _f := range rr {
			if f&_f == NoFeature {
				continue
			}
			delete(ff.runes[m], r)
		}
	}
}

var allButtons = []ButtonMask{
	Button1, Button2, Button3, Button4, Button5, Button6, Button7,
	Button8, WheelUp, WheelDown, WheelLeft, WheelRight,
}

// keyFeature maps a key to its associated feature or to NoEvent if not
// registered.
func (ff *features) keyFeature(k Key, m ModifierMask) FeatureMask {
	if ff == nil || ff.keys == nil || ff.keys[m] == nil {
		return NoFeature
	}

	return ff.keys[m][k]
}

// keyFeature maps a key to its associated feature or to NoEvent if not
// registered.
func (ff *features) buttonFeature(
	b ButtonMask, m ModifierMask,
) FeatureMask {
	if ff == nil || ff.buttons[m] == nil {
		return NoFeature
	}

	return ff.buttons[m][b]
}

// runeFeature maps a rune to its associated feature or to NoEvent if not
// registered.
func (ff *features) runeFeature(r rune, m ModifierMask) FeatureMask {
	if ff == nil || ff.runes == nil || ff.runes[ZeroModifier] == nil {
		return NoFeature
	}
	return ff.runes[m][r]
}

// allFeatures provides a slice of all elementary features
var allFeatures = []FeatureMask{
	Focusable, FocusMovable, UpScrollable, DownScrollable,
	PreviousLineFocusable, NextLineFocusable, PreviousCellFocusable,
	NextCellFocusable, FirstCellFocusable, LastCellFocusable,
	LineSelectable, LineUnfocusable, HighlightEnabled,
	TrimmedHighlightEnabled, editable,
}

type bindings struct {
	kk FeatureKeys
	rr FeatureRunes
	bb FeatureButtons
}

var defaultBindings = map[FeatureMask]*bindings{
	Focusable: {
		bb: FeatureButtons{{
			Button: Primary,
			Mod:    ZeroModifier,
		}, {
			Button: Secondary,
			Mod:    ZeroModifier,
		}, {
			Button: Middle,
			Mod:    ZeroModifier,
		}},
	},
	UpScrollable: {
		kk: FeatureKeys{{
			Key: PgUp,
			Mod: ZeroModifier,
		}},
	},
	DownScrollable: {
		kk: FeatureKeys{{
			Key: PgDn,
			Mod: ZeroModifier,
		}},
	},
	PreviousLineFocusable: {
		kk: FeatureKeys{{Key: Up, Mod: ZeroModifier}},
		// rr: FeatureRunes{{Rune: 'k', Mod: ZeroModifier}},
	},
	NextLineFocusable: {
		kk: FeatureKeys{{Key: Down, Mod: ZeroModifier}},
		// rr: FeatureRunes{{Rune: 'j', Mod: ZeroModifier}},
	},
	FirstCellFocusable: {
		kk: FeatureKeys{{Key: Home, Mod: ZeroModifier}},
	},
	PreviousCellFocusable: {
		kk: FeatureKeys{{Key: Left, Mod: ZeroModifier}},
	},
	NextCellFocusable: {
		kk: FeatureKeys{{Key: Right, Mod: ZeroModifier}},
	},
	LastCellFocusable: {
		kk: FeatureKeys{{Key: End, Mod: ZeroModifier}},
	},
	LineSelectable: {
		kk: FeatureKeys{{Key: Enter, Mod: ZeroModifier}},
	},
	LineUnfocusable: {
		kk: FeatureKeys{{Key: Esc, Mod: ZeroModifier}},
	},
	HighlightEnabled: {
		rr: FeatureRunes{{Rune: rune(0), Mod: ZeroModifier}},
	},
	editable: {
		kk: FeatureKeys{{Key: Insert, Mod: ZeroModifier}},
	},
}
