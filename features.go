// Copyright (c) 2022 Stephan Lukits. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package lines

import (
	"github.com/gdamore/tcell/v2"
)

// Features provides access and fine grained control over a components
// (end-user) features provided by lines.  Its methods will panic used
// outside an event reporting listener-callback.
type Features struct{ c *Component }

func (ff *Features) ensureInitialized() *features {
	ff.c.ensureFeatures()
	return ff.c.ff
}

// Add adds the default key, rune and button bindings of given
// feature(s) for associated component.
func (ff *Features) Add(f FeatureMask) {
	ff.ensureInitialized().add(f, false)
}

// AddRecursive sets the default key, rune and button bindings of given
// feature(s) for associated component.  Whereas the feature(s) are
// flagged recursive, i.e. they apply as well for nested components.
func (ff *Features) AddRecursive(f FeatureMask) {
	ff.ensureInitialized().add(f, true)
}

// Has returns true if receiving component features have key, rune or
// button bindings for given feature(s)
func (ff *Features) Has(f FeatureMask) bool {
	return ff.ensureInitialized().has(f)
}

// All returns all features for which currently key, rune or button
// bindings are registered. (note Has is faster to determine if a
// particular feature is set.)
func (ff *Features) All() FeatureMask {
	return ff.ensureInitialized().all()
}

// KeysOf returns the keys with their modifiers bound to given feature
// of associated component.
func (ff *Features) KeysOf(f FeatureMask) FeatureKeys {
	return ff.ensureInitialized().keysOf(f)
}

// SetKeysOf deletes all set keys for given feature (except for Quitable
// defaults) and binds given keys to it instead.  If recursive is true
// the feature becomes applicable for nested components.  The call is
// ignored if given feature is not a power of two i.e. a single feature.
// Providing no keys simply removes all key-bindings for given feature.
func (ff *Features) SetKeysOf(
	f FeatureMask, recursive bool, kk ...FeatureKey,
) {
	ff.ensureInitialized().setKeysOf(f, recursive, kk...)
}

// ButtonsOf returns the buttons with their modifiers bound to given
// feature for associated component.
func (ff *Features) ButtonsOf(f FeatureMask) FeatureButtons {
	return ff.ensureInitialized().buttonsOf(f)
}

// SetButtonsOf deletes all set buttons for given feature and binds
// given buttons to it instead.  If recursive is true the feature
// becomes applicable for nested components.  The call is ignored if
// given feature is not a power of two i.e. a single feature.  Providing
// no buttons simply removes all button-bindings for given feature.
func (ff *Features) SetButtonsOf(
	f FeatureMask, recursive bool, bb ...FeatureButton,
) {
	ff.ensureInitialized().setButtonsOf(f, recursive, bb...)
}

// RunesOf returns the runes bound to given feature for associated
// component.
func (ff *Features) RunesOf(f FeatureMask) FeatureRunes {
	return ff.ensureInitialized().runesOf(f)
}

// SetRunesOf deletes all set runes for given feature and binds given
// runes to it instead.  If recursive is true the feature becomes
// applicable for nested components.  The call is ignored if given
// feature is not a power of two i.e. a single feature.  Providing no
// runes simply removes all runes-bindings for given feature.
func (ff *Features) SetRunesOf(
	f FeatureMask, recursive bool, rr ...rune,
) {
	ff.ensureInitialized().setRunesOf(f, recursive, rr...)
}

// Delete removes all runes, key or button bindings of given feature(s)
// except for Quitable.  The two default Quitable bindings ctrl-c and
// ctrl-d remain.  NOTE you can prevent the processing of the default
// quit bindings by adding to your root component listeners for these
// keys which call StopBubbling on their environment:
//
//	type Root struct { lines.Component }
//
//	func (c *Root) OnInit(e *lines.Env) { fmt.Fprint(e, "hello world") }
//
//	func (c *Root) Keys(register lines.KeyRegistration) {
//	    register(tcell.KeyCtrlC, tcell.ModNone, func(e *Env) {
//	        e.StopBubbling()
//	    })
//	    register(tcell.KeyCtrlD, tcell.ModNone, func(e *Env) {
//	        e.StopBubbling()
//	    })
//	}
//
//	lines.New(&Root{}).Listen()
//
// gives you an application which can't be quit by its users.
func (ff *Features) Delete(f FeatureMask) {
	ff.ensureInitialized().delete(f)
}

// FeatureMask classifies keys/runes/buttons for usability features.
// I.e. features enable certain default UI-behavior for components
// having this feature to be used by a user of the final terminal
// application.  E.g. scrolling, editing ...
type FeatureMask uint64

const (

	// Quitable makes the application quitable for the user.
	Quitable FeatureMask = 1 << iota

	// Focusable enables a component to be focused by a user's mouse
	// input (default first(left)-button, second(right)-button and
	// third(middle)-button, ModMask == ModeNone).
	Focusable // TODO: implement

	// focusBubblable lets a user remove the focus from a component (or
	// a component's line) up through nesting components by pressing the
	// bubble focus key (default esc).
	focusBubblable // TODO: implement

	// PreviousSelectable components can be selected (i.e. receive the
	// focus through key-board input) by the user. (default shift-tab)
	PreviousSelectable // TODO: implement

	// NextSelectable components can be selected (i.e. receive the
	// focus through key-board input) by the user. (default tab-key)
	NextSelectable // TODO: implement

	// UpScrollable makes a component's content up-scrollable by the
	// user (default page-up-key).
	UpScrollable

	// DownScrollable makes a component's content down-scrollable by the
	// user (default page-down-key).
	DownScrollable

	// leftScrollable enables a component to be scrolled to the left
	// by the user (default left-key).
	leftScrollable // TODO: implement

	// rightScrollable enables a component to be scrolled to the right
	// by the user (default right-key).
	rightScrollable // TODO: implement

	// lineLeftScrollable makes individual lines of a component
	// scrollable to the left (default left-key)
	lineLeftScrollable // TODO: implement

	// lineRightScrollable makes individual lines of a component
	// scrollable to the right (default right-key)
	lineRightScrollable // TODO: implement

	// previousLineSelectable lets the user select a component's lines
	// in descending line-index direction. (default up-key)
	previousLineSelectable // TODO: implement

	// nextLineSelectable lets the user select a component's lines
	// in ascending line-index direction. (default down-key)
	nextLineSelectable // TODO: implement

	// maximizable lets the user maximize a component, i.e. all siblings
	// which are collapsed to either one line if parent is stacking or
	// to one column if parent is chaining. (default shift-primary-button)
	maximizable // TODO: implement

	// minimizable lets the user minimize a component with this feature,
	// i.e. the component is collapsed to one line in a stacking parent
	// or to one column in a chaining parent.
	minimizable // TODO: implement

	// editable makes a component's content editable by the user.
	editable // TODO: implement

	// _recursive flags component FeatureMask-settings as applicable for
	// all nested components as well.
	_recursive

	// NoFeature classifies keys/runes/buttons not registered for any
	// feature.
	NoFeature FeatureMask = 0

	// Selectable makes a component focusable through keyboard input by
	// combining next- and previous-selectable.
	Selectable = PreviousSelectable | NextSelectable // TODO: implement

	// Scrollable makes a component's content vertically Scrollable by
	// combining up- and down-Scrollable.
	Scrollable = UpScrollable | DownScrollable

	// horizontalScrollable makes a component horizontally scrollable by
	// combining left- and right-scrollable.
	horizontalScrollable = leftScrollable | rightScrollable // TODO: implement

	// lineScrollable makes individual lines of a component horizontally
	// scrollable by combining line-left- and line-right-scrollable.
	lineScrollable = lineLeftScrollable | lineRightScrollable // TODO: implement

	// linesSelectable makes a component's lines selectable by combining
	// previous-line- and next-line-selectable.
	linesSelectable = previousLineSelectable | nextLineSelectable // TODO: implement
)

// features provides information about keys/runes/buttons which are
// registered for features provided by the lines-package.  It also
// allows to change these in a consistent and convenient way.  The zero
// value is not ready to use.  Make a copy of DefaultFeatures to create
// a new features-instance.  Note  A *Register* instance is always with
// a copy of the *DefaultFeatures* features-instance initialized which
// holds the quit-feature only.
type features struct {
	keys    map[tcell.ModMask]map[tcell.Key]FeatureMask
	runes   map[rune]FeatureMask
	buttons map[tcell.ModMask]map[tcell.ButtonMask]FeatureMask
	have    FeatureMask
}

// modifiable returns false for the default features.
func (ff *features) modifiable() bool {
	_, ok := ff.runes[0]
	return ok
}

// keyQuits returns true if given key is associated with the
// quit-feature.
func (ff *features) keyQuits(k tcell.Key) bool {
	return ff.keys[tcell.ModNone][k]&Quitable != NoFeature
}

// runeQuits return true if given rune is associated with the
// quit-feature.
func (ff *features) runeQuits(r rune) bool {
	return ff.runes[r]&Quitable != NoFeature
}

// copy creates a new Features instance initialized with the features of
// receiving Features instance.
func (ff *features) copy() *features {

	cpy := features{
		keys:    map[tcell.ModMask]map[tcell.Key]FeatureMask{},
		runes:   map[rune]FeatureMask{},
		buttons: map[tcell.ModMask]map[tcell.ButtonMask]FeatureMask{},
		have:    ff.have,
	}

	for m, kk := range ff.keys {
		cpy.keys[m] = map[tcell.Key]FeatureMask{}
		for k, f := range kk {
			cpy.keys[m][k] = f
		}
	}

	for r, f := range ff.runes {
		cpy.runes[r] = f
	}

	for m, bb := range ff.buttons {
		cpy.buttons[m] = map[tcell.ButtonMask]FeatureMask{}
		for b, f := range bb {
			cpy.buttons[m][b] = f
		}
	}

	return &cpy
}

func (ff *features) has(f FeatureMask) bool {
	if f&_recursive == NoFeature {
		return ff.have&f != NoFeature
	}

	have := false

	cb := func(_f FeatureMask) (stop bool) {
		if f&_f != f {
			return false
		}
		have = true
		return true
	}

	ff.forRuneFeatures(cb)
	if have {
		return true
	}

	ff.forKeyFeatures(cb)
	if have {
		return true
	}

	ff.forButtonFeatures(cb)
	return have
}

// Registered returns the set of features currently all.
func (ff *features) all() FeatureMask {

	_ff := NoFeature

	cb := func(f FeatureMask) (stoop bool) {
		if f == NoFeature {
			return false
		}
		_ff |= f &^ _recursive
		return false
	}

	ff.forRuneFeatures(cb)
	ff.forKeyFeatures(cb)
	ff.forButtonFeatures(cb)

	return _ff
}

func (ff *features) forRuneFeatures(cb func(FeatureMask) (stoop bool)) {
	for _, f := range ff.runes {
		if cb(f) {
			return
		}
	}
}

func (ff *features) forKeyFeatures(cb func(FeatureMask) (stoop bool)) {
	for _, kk := range ff.keys {
		for _, f := range kk {
			if cb(f) {
				return
			}
		}
	}
}

func (ff *features) forButtonFeatures(cb func(FeatureMask) (stoop bool)) {
	for _, bb := range ff.buttons {
		for _, f := range bb {
			if cb(f) {
				return
			}
		}
	}
}

// add adds the default bindings of given feature to receiving component
// features.
func (ff *features) add(f FeatureMask, recursive bool) {

	_ff := []FeatureMask{}
	for _, _f := range allFeatures {
		if _f&f == NoFeature {
			continue
		}
		_ff = append(_ff, _f)
	}
	if !ff.modifiable() || len(_ff) == 0 {
		return
	}

	for _, f := range _ff {
		df := defaultBindings[f]
		if recursive {
			f |= _recursive
		}
		for _, k := range df.kk {
			if ff.keys[k.Mod] == nil {
				ff.keys[k.Mod] = map[tcell.Key]FeatureMask{}
			}
			ff.keys[k.Mod][k.Key] = f
		}

		for _, b := range df.bb {
			if ff.buttons[b.Mod] == nil {
				ff.buttons[b.Mod] = map[tcell.ButtonMask]FeatureMask{}
			}
			ff.buttons[b.Mod][b.Button] = f
		}

		for _, r := range df.rr {
			ff.runes[r] = f
		}
	}

	if ff.have&f == NoFeature {
		ff.have |= f
	}
}

// FeatureKey represents a key bound to a feature with its tcell.ModMask
// and tcell.Key value.  FeatureKey instances must be also provided to
// SetKeysOf.
type FeatureKey struct {
	Mod tcell.ModMask
	Key tcell.Key
}

// FeatureKeys are provided by KeysOf of an Features instance reporting
// the keys bound to a given feature.  FeaturesKeys may be also used as
// variadic argument for an Features instance's SetKeysOf.
type FeatureKeys []FeatureKey

// Equals returns true if both slices contain the same FeatureKey
// values.
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

// setKeysOf deletes all set keys for given feature (except for Quitable
// defaults) an binds given keys to it instead.  If recursive is true
// the feature becomes applicable for nested components.  The call is
// ignored if given feature is not a power of two i.e. a single feature.
// NOTE providing no keys simply removes all key-bindings for given
// feature.
func (ff *features) setKeysOf(
	f FeatureMask, recursive bool, kk ...FeatureKey,
) {
	if f == 0 || f&(f-1) != 0 { // f is not a power of two
		return
	}

	if f != Quitable {
		ff.deleteKeysOf(f)
	} else {
		ff.deleteKeysOfButDefaults(f)
	}

	if recursive {
		f |= _recursive
	}
	for _, k := range kk {
		if ff.keys[k.Mod] == nil {
			ff.keys[k.Mod] = map[tcell.Key]FeatureMask{}
		}
		ff.keys[k.Mod][k.Key] = f
	}

	ff.have = ff.all()
}

// FeatureButton represents a button (mask) bound to a feature with its
// tcell.ModMask and tcell.ButtonMask value.  FeatureButton instances
// must be also provided to SetButtonsOf.
type FeatureButton struct {
	Mod    tcell.ModMask
	Button tcell.ButtonMask
}

// FeatureButtons are provided by ButtonsOf of an Features instance
// reporting the mouse buttons bound to a given feature.
// FeatureButtons may be also used as variadic argument for an Features
// instance's SetButtonsOf.
type FeatureButtons []FeatureButton

// Equals returns true if receiving and given FeatureButtons contain the
// same buttons.
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
func (ff *features) setButtonsOf(
	f FeatureMask, recursive bool, bb ...FeatureButton,
) {

	if f == 0 || f&(f-1) != 0 { // f is not a power of two
		return
	}

	ff.deleteButtonsOf(f)

	if recursive {
		f |= _recursive
	}
	for _, b := range bb {
		if ff.buttons[b.Mod] == nil {
			ff.buttons[b.Mod] = map[tcell.ButtonMask]FeatureMask{}
		}
		ff.buttons[b.Mod][b.Button] = f
	}

	ff.have = ff.all()
}

// FeatureRunes are provided by RunesOf of an Features instance
// reporting the runes bound to a given feature.  FeaturesRunes may be
// also used as variadic argument for an Features instance's SetRunesOf.
type FeatureRunes []rune

// Equals returns true if receiving and given FeatureRunes contain the
// same runes.
func (fr FeatureRunes) Equals(other FeatureRunes) bool {
	if len(fr) != len(other) {
		return false
	}
	for _, r := range fr {
		has := false
		for _, o := range other {
			if r != o {
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
	rr := []rune{}
	for r, _f := range ff.runes {
		if f&_f == NoFeature {
			continue
		}
		rr = append(rr, r)
	}
	return rr
}

// setRunesOf deletes all set runes for given feature an binds given
// runes to it instead.  If recursive is true the feature becomes
// applicable for nested components.  The call is ignored if given
// feature is not a power of two i.e. a single feature.  NOTE providing
// no runes simply removes all runes-bindings for given feature.
func (ff *features) setRunesOf(
	f FeatureMask, recursive bool, rr ...rune,
) {

	if f == 0 || f&(f-1) != 0 { // f is not a power of two
		return
	}

	ff.deleteRunesOf(f)

	if recursive {
		f |= _recursive
	}
	for _, r := range rr {
		ff.runes[r] = f
	}

	ff.have = ff.all()
}

// delete removes all runes, key or button bindings of given feature(s)
// except for Quitable.  The two default Quitable bindings ctrl-c and
// ctrl-d remain.
func (ff *features) delete(f FeatureMask) {
	_ff := []FeatureMask{}
	for _, _f := range allFeatures {
		if _f&f == NoFeature {
			continue
		}
		_ff = append(_ff, _f)
	}
	if !ff.modifiable() || len(_ff) == 0 {
		return
	}

	for _, f := range _ff {
		ff.deleteKeysOf(f)
		ff.deleteButtonsOf(f)
		ff.deleteRunesOf(f)
		ff.have &^= f
	}
}

func (ff *features) deleteKeysOf(f FeatureMask) {
	if f == Quitable {
		ff.deleteKeysOfButDefaults(f)
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

func (ff *features) deleteKeysOfButDefaults(f FeatureMask) {
	for m, kk := range ff.keys {
		for k, _f := range kk {
			if f&_f == NoFeature {
				continue
			}
			if defaultFeatures.keyFeature(k, m)&f != NoFeature {
				continue
			}
			delete(ff.keys[m], k)
		}
	}
}

func (ff *features) deleteButtonsOf(f FeatureMask) {
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
	for r, _f := range ff.runes {
		if f&_f == NoFeature {
			continue
		}
		delete(ff.runes, r)
	}
}

var allButtons = []tcell.ButtonMask{
	tcell.Button1, tcell.Button2, tcell.Button3, tcell.Button4,
	tcell.Button5, tcell.Button6, tcell.Button7, tcell.Button8,
	tcell.WheelUp, tcell.WheelDown, tcell.WheelLeft, tcell.WheelRight,
}

// keyFeature maps a key to its associated feature or to NoEvent if not
// registered.
func (ff *features) keyFeature(k tcell.Key, m tcell.ModMask) FeatureMask {
	if ff == nil || ff.keys == nil {
		return NoFeature
	}

	if ff.keys[m] == nil {
		return NoFeature
	}

	return ff.keys[m][k]
}

// keyFeature maps a key to its associated feature or to NoEvent if not
// registered.
func (ff *features) buttonFeature(
	b tcell.ButtonMask, m tcell.ModMask,
) FeatureMask {

	if ff.buttons[m] == nil {
		return NoFeature
	}

	return ff.buttons[m][b]
}

// runeFeature maps a rune to its associated feature or to NoEvent if not
// registered.
func (kk *features) runeFeature(r rune) FeatureMask {
	return kk.runes[r]
}

// allFeatures provides a slice of all the potentially internally
// handled features
var allFeatures = []FeatureMask{
	Quitable, Focusable,
	UpScrollable, DownScrollable,
	leftScrollable, rightScrollable,
	lineLeftScrollable, lineRightScrollable,
	PreviousSelectable, NextSelectable,
	previousLineSelectable, nextLineSelectable,
	maximizable, editable,
}

// defaultFeatures are the default runes and keys which are associated
// with (end-user) features.  NOTE defaultFeatures cannot be
// modified, a copy of them can!
var defaultFeatures = &features{
	keys: map[tcell.ModMask]map[tcell.Key]FeatureMask{
		tcell.ModNone: {
			tcell.KeyCtrlC: Quitable,
			tcell.KeyCtrlD: Quitable,
		},
	},
	runes: map[rune]FeatureMask{
		0:   NoFeature, // indicates the immutable default features
		'q': Quitable,
	},
	buttons: map[tcell.ModMask]map[tcell.ButtonMask]FeatureMask{},
	have:    Quitable,
}

type bindings struct {
	kk FeatureKeys
	rr FeatureRunes
	bb FeatureButtons
}

var defaultBindings = map[FeatureMask]*bindings{
	Focusable: {
		bb: FeatureButtons{{
			Button: tcell.ButtonPrimary,
			Mod:    tcell.ModNone,
		}, {
			Button: tcell.ButtonSecondary,
			Mod:    tcell.ModNone,
		}, {
			Button: tcell.ButtonMiddle,
			Mod:    tcell.ModNone,
		}},
	},
	NextSelectable: {
		kk: FeatureKeys{{
			Key: tcell.KeyTAB,
			Mod: tcell.ModNone,
		}},
	},
	PreviousSelectable: {
		kk: FeatureKeys{{
			Key: tcell.KeyTAB,
			Mod: tcell.ModShift,
		}},
	},
	UpScrollable: {
		kk: FeatureKeys{{
			Key: tcell.KeyPgUp,
			Mod: tcell.ModNone,
		}},
	},
	DownScrollable: {
		kk: FeatureKeys{{
			Key: tcell.KeyPgDn,
			Mod: tcell.ModNone,
		}},
	},
}

func execute(usr Componenter, f FeatureMask) {
	switch f {
	case UpScrollable:
		usr.embedded().Scroll.Up()
	case DownScrollable:
		usr.embedded().Scroll.Down()
	}
}
