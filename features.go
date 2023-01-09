// Copyright (c) 2022 Stephan Lukits. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

/*
Features are a convenient way to extend a components default behavior.
A typical feature is comprised of elementary features defined as
FeatureMask constants, e.g.:

    const (
        Quitable FeatureMask = 1 << iota
        // ...
        PreviousLineFocusable
        NextLineFocusable
        LineUnfocusable
        // ...
        LinesFocusable = PreviousLineFocusable | NextLineFocusable |
            LineUnfocusable
        HighlightedFocusable = LinesFocusable | highlightedFocusable
        // ...
    )

Most elementary features are public to allow modifying their key, rune
or button bindings.  Elementary features which are not associated with
any key, rune or button like highlightedFocusable are kept private.
They represent usually a mere variation of existing features.  All
elementary features need to be in "allFeatures"

Default key, rune and button bindings for each elementary feature is
defined in defaultBinding, e.g.:

    var defaultBindings = map[FeatureMask]*bindings{
        // ...
        previousLineSelectable: {
            kk: FeatureKeys{{Key: KeyUp, Mod: ZeroModifier}},
            rr: FeatureRunes{'k'},
        },
        // ...
    }

Elementary features not associated with any key, rune or button are
bound to the zero-rune.  A component's associated Features-instance is
used to add (comprised) features, e.g.:

    func (c *Cmp) OnInit(_ lines.Env) {
        c.FF.Add(Scrollable|LinesFocusable)
    }

Finally the "execute" function is used by key, rune or button-reporter
to dispatch the according operation implementing the feature and
reporting associated events if any, e.g.:

    func execute(cntx *rprContext, usr Componenter, f FeatureMask) {
        switch f {
        // ...
        case NextLineFocusable:
            current := usr.Focus.Current()
            if usr.Focus.Next() != current {
                reportOnLineFocus(cntx, usr)
            }
        // ...
        }
    }
*/

package lines

// Features provides access and fine grained control over the behavior
// of a component provided by lines.  Its methods will panic used
// outside an event reporting listener-callback.  Typically you will use
// a component's FF-property to manipulate a component's supported
// features, e.g.
//
//	type Cmp { lines.Component }
//
//	func (c *Cmp) OnInit(_ *lines.Env) {
//		c.FF.Add(lines.LinesFocusable)
//	}
//
// adds the feature "selectable lines" to a component. I.e. if the
// component has the focus up/down keys highlight selectable lines of
// the component while an enter-key-press reports an OnLineSelection of
// the currently highlighted line and an esc-key-press removes the
// line highlighting.
type Features struct{ c *Component }

func (ff *Features) ensureInitialized() *features {
	ff.c.ensureFeatures()
	return ff.c.ff
}

// Add adds the default key, rune and button bindings of given
// feature(s) for associated component.
func (ff *Features) Add(f FeatureMask) {
	if f&editable == editable && !ff.c.isNesting() && ff.c.Edit == nil {
		_, _, hasCursor := ff.c.CursorPosition()
		ff.c.Edit = &Editor{c: ff.c, suspended: !hasCursor}
		if hasCursor {
			ff.c.LL.Focus.eolAfterLastRune = true
		}
	}
	ff.ensureInitialized().add(f, false)
}

// AddRecursive sets the default key, rune and button bindings of given
// feature(s) for associated component.  Whereas the feature(s) are
// flagged recursive, i.e. they apply as well for nested components.
func (ff *Features) AddRecursive(f FeatureMask) {
	ff.ensureInitialized().add(f, true)
}

// OfKey returns the feature bound to given key k with given modifiers
// mm.
func (ff *Features) OfKey(k Key, mm ModifierMask) FeatureMask {
	return ff.ensureInitialized().keyFeature(k, mm)
}

// OfRune returns the feature bound to given rune r with given modifiers
// mm.
func (ff *Features) OfRune(r rune, mm ModifierMask) FeatureMask {
	return ff.ensureInitialized().runeFeature(r, mm)
}

// OfButton returns the feature bound to given buttons b with given modifiers
// mm.
func (ff *Features) OfButton(
	bb ButtonMask, mm ModifierMask,
) FeatureMask {
	ff.ensureInitialized()
	return ff.c.ff.buttonFeature(bb, mm)
}

// Has returns true if receiving component features have key, rune or
// button bindings for given feature(s).
func (ff *Features) Has(f FeatureMask) bool {
	return ff.ensureInitialized().has(f)
}

// All returns all features for which currently key, rune or button
// bindings are registered. (note [Features.Has] is faster to determine if a
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
	f FeatureMask, recursive bool, rr ...FeatureRune,
) {
	ff.ensureInitialized().setRunesOf(f, recursive, rr...)
}

// Delete removes all runes, key or button bindings of given feature(s)
// except for Quitable.  The two default Quitable bindings ctrl-c and
// ctrl-d remain.  NOTE use a *Kiosk constructor like [TermKiosk] for a
// [Lines]-instance to avoid having Quitable set by default.
// func (ff *Features) Delete(f FeatureMask) {
// 	ff.ensureInitialized().delete(f)
// }

// FeatureMask classifies keys/runes/buttons for a components default
// behavior like focusable, scrollable etc.
type FeatureMask uint64

const (

	// Quitable makes the application quitable for the user.
	Quitable FeatureMask = 1 << iota

	// Focusable enables a component to be focused by a user's mouse
	// input (default first(left)-button, second(right)-button and
	// third(middle)-button, ModMask == ModeNone).
	Focusable

	// focusBubblable lets a user remove the focus from a component (or
	// a component's line) up through nesting components by pressing the
	// bubble focus key (default esc).
	focusBubblable // TODO: implement

	// previousSelectable components can be selected (i.e. receive the
	// focus through key-board input) by the user. (default shift-tab)
	previousSelectable // TODO: implement

	// nextSelectable components can be selected (i.e. receive the
	// focus through key-board input) by the user. (default tab-key)
	nextSelectable // TODO: implement

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

	// nextLineFocusable a component's previous focusable line can
	// receive the focus. (default down-key)
	PreviousLineFocusable

	// NextLineFocusable a component's next focusable line can receive
	// the focus. (default up-key)
	NextLineFocusable

	// LineUnfocusable a component's set line-focus can be removed
	// (default esc)
	LineUnfocusable

	// linesHighlightedFocusable highlights a component's focused line.
	linesHighlightedFocusable

	// LineSelectable a component's focused line can be reported as
	// selected (default enter)
	LineSelectable

	// PreviousCellFocusable shows the cursor in a components focused
	// line while the left arrow key moves the cursor to ("focuses") the
	// previous rune.
	PreviousCellFocusable

	LastCellFocusable

	// NextCellFocusable shows the cursor in a components focused
	// line while the right arrow key moves the cursor to ("focuses") the
	// previous rune.
	NextCellFocusable

	FirstCellFocusable

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

	// selectable makes a component focusable through keyboard input by
	// combining next- and previous-selectable.
	selectable = previousSelectable | nextSelectable // TODO: implement

	// Scrollable makes a component's content vertically Scrollable by
	// combining up- and down-Scrollable.
	Scrollable = UpScrollable | DownScrollable

	// horizontalScrollable makes a component horizontally scrollable by
	// combining left- and right-scrollable.
	horizontalScrollable = leftScrollable | rightScrollable // TODO: implement

	// lineScrollable makes individual lines of a component horizontally
	// scrollable by combining line-left- and line-right-scrollable.
	lineScrollable = lineLeftScrollable | lineRightScrollable // TODO: implement

	// LinesFocusable makes lines focusable, i.e. a line receiving the
	// focus is reported; see OnLineFocus.
	LinesFocusable = NextLineFocusable | PreviousLineFocusable |
		LineUnfocusable

	// LinesHighlightedFocusable makes lines focusable whereas the focused
	// line is highlighted.
	LinesHighlightedFocusable = LinesFocusable | linesHighlightedFocusable

	// LinesSelectable makes a component's lines selectable by combining
	// HighlightedFocusable and LineSelectable.
	LinesSelectable = LinesHighlightedFocusable | LineSelectable

	// CellFocusable turns on LinesFocusable for a component c and shows
	// the cursor whose positioning indicates the "focused cell".
	CellFocusable = PreviousCellFocusable | NextCellFocusable |
		LinesFocusable | LastCellFocusable | FirstCellFocusable

	// RuneFocusable turns on LinesHighlightedFocusable for a component
	// c and shows the cursor whose positioning indicates the "focused
	// rune".
	CellHighlightedFocusable = CellFocusable | LinesHighlightedFocusable

	Editable = Focusable | CellFocusable | Scrollable | editable

	HighlightedEditable = Focusable | CellHighlightedFocusable |
		Scrollable | editable
)

// features provides information about keys/runes/buttons which are
// registered for features provided by the lines-package.  It also
// allows to change these in a consistent and convenient way.  The zero
// value is not ready to use.  Make a copy of DefaultFeatures to create
// a new features-instance.  Note  A *Register* instance is always with
// a copy of the *DefaultFeatures* features-instance initialized which
// holds the quit-feature only.
type features struct {
	keys    map[ModifierMask]map[Key]FeatureMask
	runes   map[ModifierMask]map[rune]FeatureMask
	buttons map[ModifierMask]map[ButtonMask]FeatureMask
	have    FeatureMask
}

// modifiable returns false for the default features.
func (ff *features) modifiable() bool {
	_, ok := ff.runes[0]
	return ok
}

// keyQuits returns true if given key is associated with the
// quit-feature.
func (ff *features) keyQuits(k Key) bool {
	return ff.keys[ZeroModifier][k]&Quitable != NoFeature
}

// runeQuits return true if given rune is associated with the
// quit-feature.
func (ff *features) runeQuits(r rune) bool {
	return ff.runes[ZeroModifier][r]&Quitable != NoFeature
}

// copy creates a new Features instance initialized with the features of
// receiving Features instance.
func (ff *features) copy() *features {

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
	if f&_recursive == NoFeature {
		return ff.have&f == f
	}

	// since we can't know which of ff.have features is combined with
	// _recursive we need to find exactly f.  Note that recursive
	// combined features (e.g. LinesFocusable|_recursive) cannot be
	// found as of now.

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

// all returns the set of currently set features.
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
	for _, rr := range ff.runes {
		for _, f := range rr {
			if cb(f) {
				return
			}
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
		if df == nil {
			continue
		}
		if recursive {
			f |= _recursive
		}
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

	if ff.have&f == NoFeature {
		ff.have |= f
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
			ff.keys[k.Mod] = map[Key]FeatureMask{}
		}
		ff.keys[k.Mod][k.Key] = f
	}

	ff.have = ff.all()
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
			ff.buttons[b.Mod] = map[ButtonMask]FeatureMask{}
		}
		ff.buttons[b.Mod][b.Button] = f
	}

	ff.have = ff.all()
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
// runes to it instead.  If recursive is true the feature becomes
// applicable for nested components.  The call is ignored if given
// feature is not a power of two i.e. a single feature.  NOTE providing
// no runes simply removes all runes-bindings for given feature.
func (ff *features) setRunesOf(
	f FeatureMask, recursive bool, rr ...FeatureRune,
) {

	if f == 0 || f&(f-1) != 0 { // f is not a power of two
		return
	}

	ff.deleteRunesOf(f)

	if recursive {
		f |= _recursive
	}
	for _, r := range rr {
		if ff.runes[r.Mod] == nil {
			ff.runes[r.Mod] = map[rune]FeatureMask{}
		}
		ff.runes[r.Mod][r.Rune] = f
	}

	ff.have = ff.all()
}

// delete removes all runes, key or button bindings of given feature(s)
// except for Quitable.  The two default Quitable bindings ctrl-c and
// ctrl-d remain.
// func (ff *features) delete(f FeatureMask) {
// 	_ff := []FeatureMask{}
// 	for _, _f := range allFeatures {
// 		if _f&f == NoFeature {
// 			continue
// 		}
// 		_ff = append(_ff, _f)
// 	}
// 	if !ff.modifiable() || len(_ff) == 0 {
// 		return
// 	}
//
// 	for _, f := range _ff {
// 		ff.deleteKeysOf(f)
// 		ff.deleteButtonsOf(f)
// 		ff.deleteRunesOf(f)
// 		ff.have &^= f
// 	}
// }

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
			if quitableFeatures.keyFeature(k, m)&f != NoFeature {
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
	Quitable, Focusable,
	UpScrollable, DownScrollable,
	leftScrollable, rightScrollable,
	lineLeftScrollable, lineRightScrollable,
	previousSelectable, nextSelectable,
	PreviousLineFocusable, NextLineFocusable, PreviousCellFocusable,
	NextCellFocusable, FirstCellFocusable, LastCellFocusable,
	LineSelectable, LineUnfocusable, linesHighlightedFocusable,
	maximizable, editable,
}

// quitableFeatures are the default runes and keys of a clients root
// component.  NOTE quitableFeatures cannot be modified, a copy of them
// can!
var quitableFeatures = &features{
	keys: map[ModifierMask]map[Key]FeatureMask{
		ZeroModifier: {
			CtrlC: Quitable,
			CtrlD: Quitable,
		},
	},
	runes: map[ModifierMask]map[rune]FeatureMask{
		ZeroModifier: {
			0:   NoFeature, // indicates the immutable default features
			'q': Quitable,
		},
	},
	buttons: map[ModifierMask]map[ButtonMask]FeatureMask{},
	have:    Quitable,
}

// defaultFeatures is an unmodifiable initialized features instance
// which may be used to create a new features instance by copy()-ing it.
var defaultFeatures = &features{
	keys: map[ModifierMask]map[Key]FeatureMask{},
	runes: map[ModifierMask]map[rune]FeatureMask{ZeroModifier: {
		0: NoFeature, // indicates the immutable zero features
	}},
	buttons: map[ModifierMask]map[ButtonMask]FeatureMask{},
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
	nextSelectable: {
		kk: FeatureKeys{{
			Key: TAB,
			Mod: ZeroModifier,
		}},
	},
	previousSelectable: {
		kk: FeatureKeys{{
			Key: TAB,
			Mod: Shift,
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
	linesHighlightedFocusable: {
		rr: FeatureRunes{{Rune: rune(0), Mod: ZeroModifier}},
	},
}

// LineSelecter is implemented by a component who wants to be informed
// when its focused line was selected.
type LineSelecter interface {

	// OnLineSelection is called by Lines if the focused line having
	// index i of implementing component was selected.
	OnLineSelection(_ *Env, cIdx, sIdx int)
}

// LineFocuser is implemented by a component who wants to be informed when
// one of its lines receives the focus.
type LineFocuser interface {

	// OnLineFocus is called by Lines if the line with given content
	// line index cIdx and given screen line index sIdx receives the
	// focus.  Note is a component is associated with a source or there
	// are more content lines than fitting on the screen cLine and sLine
	// may differ.
	OnLineFocus(_ *Env, cIdx, sIdx int)
}

// LineFocusLooser is implemented by a component who wants to be informed when
// a focused line looses its focus.
type LineFocusLooser interface {

	// OnLineFocusLost is called by Lines if the line with given content
	// line index cIdx and given screen line index sIdx of implementing
	// component lost the focus.  Note "on line focus lost" is reported
	// after the focus change has happened, i.e. given screen line sIdx
	// not necessary displays the content line cIdx.  "on line focus
	// lost" is reported before "on line focus" is report (if it is
	// reported).
	OnLineFocusLost(_ *Env, cIdx, sIdx int)
}

// LineOverflower implementations are called back when a line receives
// the focus whose content either overflows to the left or to the right
// or at both sides.
type LineOverflower interface {

	// LineOverflows is called by Lines if a components line receives
	// the focus whose content overflows to the left or to the right or
	// at both sides.
	OnLineOverflowing(_ *Env, left, right bool)
}

// Cursorer is implemented by a component which wants to be notified
// about cursor movement.
type Cursorer interface {

	// OnCursor implemented by a component c is called by Lines if the
	// cursor position has changed; use c.CursorPosition() to retrieve
	// the current cursor position.  Note if the display resizes Lines
	// either removes the cursor iff it is not in the content area of a
	// component; otherwise it keeps the cursor in c's content area
	// trying to keep it relative to the content areas origin at the
	// same position.  If the later can be achieved absOnly is true.
	OnCursor(_ *Env, absOnly bool)
}

// execute given feature f on given user-component usr.
func execute(cntx *rprContext, usr Componenter, f FeatureMask) {
	switch f {
	case UpScrollable:
		usr.embedded().Scroll.Up()
	case DownScrollable:
		usr.embedded().Scroll.Down()
	case NextLineFocusable:
		executeLineFocus(cntx, usr, usr.embedded().LL.Focus.Next)
	case PreviousLineFocusable:
		executeLineFocus(cntx, usr, usr.embedded().LL.Focus.Previous)
	case FirstCellFocusable:
		executeCellFocus(cntx, usr, usr.embedded().LL.Focus.FirstCell)
	case PreviousCellFocusable:
		executeCellFocus(cntx, usr, usr.embedded().LL.Focus.PreviousCell)
	case NextCellFocusable:
		executeCellFocus(cntx, usr, usr.embedded().LL.Focus.NextCell)
	case LastCellFocusable:
		executeCellFocus(cntx, usr, usr.embedded().LL.Focus.LastCell)
	case LineUnfocusable:
		executeResetLineFocus(cntx, usr)
	case LineSelectable:
		reportSelectedLine(cntx, usr)
	}
}

func executeLineFocus(
	cntx *rprContext, usr Componenter, f func(bool) (int, int),
) {
	cIdx, sIdx := usr.embedded().LL.Focus.Current(),
		usr.embedded().LL.Focus.Screen()
	_, column, _ := usr.embedded().CursorPosition()
	// TODO: figure out what that does especially in the context that
	// a set zero-rune feature indicates that a feature set is not
	// modifiable.
	rf := usr.embedded().ff.runeFeature(rune(0), ZeroModifier)
	highlighted := rf&linesHighlightedFocusable == linesHighlightedFocusable
	ln, cl := f(highlighted)
	if cIdx == ln {
		if cl != column {
			reportCursorChange(cntx, usr)
		}
		return
	}
	if cIdx < 0 && usr.embedded().Edit != nil {
		usr.embedded().Edit.Resume()
		if usr.embedded().Edit.IsReplacing() {
			usr.embedded().LL.Focus.EolAtLastRune()
		} else {
			usr.embedded().LL.Focus.EolAfterLastRune()
		}
	}
	reportLineFocus(cntx, usr, cIdx, sIdx)
	if cl == column && column == -1 {
		return
	}
	reportCursorChange(cntx, usr)
}

func executeResetLineFocus(cntx *rprContext, usr Componenter) {
	cIdx, sIdx := usr.embedded().LL.Focus.Current(),
		usr.embedded().LL.Focus.Screen()
	_, _, haveCursor := usr.embedded().CursorPosition()
	usr.embedded().LL.Focus.Reset()
	if haveCursor {
		reportCursorChange(cntx, usr)
	}
	reportLineFocus(cntx, usr, cIdx, sIdx)
	if usr.embedded().Edit != nil {
		usr.embedded().Edit.Suspend()
	}
}

func executeCellFocus(
	cntx *rprContext, usr Componenter, f func() (int, int, bool),
) {
	_, _, movedCursor := f()
	if movedCursor && reportCursorChange(cntx, usr) {
		usr.enable()
	}
	reportLineOverflow(cntx, usr, usr.embedded().LL.Focus.Screen())
}

func reportCursorChange(cntx *rprContext, usr Componenter) bool {
	c, ok := usr.(Cursorer)
	if !ok {
		return false
	}
	callback(usr, cntx, func(c Cursorer) func(e *Env) {
		return func(e *Env) { c.OnCursor(e, false) }
	}(c))
	return true
}

func lfCurry(cb func(*Env, int, int), cLine, sLine int) func(*Env) {
	return func(e *Env) { cb(e, cLine, sLine) }
}

func ofCurry(of LineOverflower, left, right bool) func(*Env) {
	return func(e *Env) { of.OnLineOverflowing(e, left, right) }
}

func reportLineFocus(cntx *rprContext, usr Componenter, cIdx, sIdx int) {
	fl, ok := usr.(LineFocusLooser)
	if cIdx >= 0 && ok {
		callback(usr, cntx, lfCurry(fl.OnLineFocusLost, cIdx, sIdx))
		usr.enable()
	}
	cmp := usr.embedded()
	sIdx = cmp.LL.Focus.Screen()
	cIdx = cmp.LL.Focus.Current()
	if cIdx < 0 {
		return
	}
	lf, ok := usr.(LineFocuser)
	if ok {
		callback(usr, cntx, lfCurry(lf.OnLineFocus, cIdx, sIdx))
		usr.enable()
	}
	reportLineOverflow(cntx, usr, sIdx)
}

func reportLineOverflow(cntx *rprContext, usr Componenter, sIdx int) {
	if sIdx < 0 {
		return
	}
	of, ok := usr.(LineOverflower)
	if !ok {
		return
	}
	cmp := usr.embedded()
	_, _, width, _ := cmp.ContentArea()
	l, r, changed := cmp.LL.By(sIdx).isOverflowing(width)
	if !l && !r || !changed {
		return
	}
	callback(usr, cntx, ofCurry(of, l, r))
}

func lsCurry(ls LineSelecter, cIdx, sIdx int) func(*Env) {
	return func(e *Env) { ls.OnLineSelection(e, cIdx, sIdx) }
}

func reportSelectedLine(cntx *rprContext, usr Componenter) {
	cIdx, sIdx := usr.embedded().LL.Focus.Current(),
		usr.embedded().LL.Focus.Screen()
	if cIdx < 0 {
		return
	}
	ls, ok := usr.(LineSelecter)
	if !ok {
		return
	}
	callback(usr, cntx, lsCurry(ls, cIdx, sIdx))
}
