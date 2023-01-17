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

const (

	// Focusable enables a component to be focused by a user's mouse
	// input (default first(left)-button, second(right)-button and
	// third(middle)-button, ModMask == ModeNone).
	Focusable FeatureMask = 1 << iota

	// FocusMoveable enables to move focus between FocusMovable
	// components using the default Tab-key.
	FocusMovable

	// UpScrollable makes a component's content up-scrollable by the
	// user (default page-up-key).
	UpScrollable

	// DownScrollable makes a component's content down-scrollable by the
	// user (default page-down-key).
	DownScrollable

	// nextLineFocusable a component's previous focusable line can
	// receive the focus. (default down-key)
	PreviousLineFocusable

	// NextLineFocusable a component's next focusable line can receive
	// the focus. (default up-key)
	NextLineFocusable

	// LineUnfocusable a component's set line-focus can be removed
	// (default Esc)
	LineUnfocusable

	// HighlightEnabled if set will highlight a component c's
	// focused screen line using c.Globals().Style(Highlight).  (NOTE
	// use c.Globals() to modify the Highlight style)
	HighlightEnabled

	// TrimmedHighlightEnabled if set will highlight a component
	// c's focused screen line whereas the highlight is only applied to
	// content without prefixing or suffixing whitespace using
	// c.Globals().Style(Highlight).  (NOTE use c.Globals() to modify
	// the Highlight style)
	TrimmedHighlightEnabled

	// LineSelectable a component's focused line can be reported as
	// selected (default Enter)
	LineSelectable

	// FirstCellFocusable shows/moves the cursor at the beginning of a
	// focused screen line's content in the line's first cell (default
	// Home).
	FirstCellFocusable

	// PreviousCellFocusable shows the cursor in a components focused
	// line while the left arrow key moves the cursor to ("focuses") the
	// previous rune.
	PreviousCellFocusable

	// LastCellFocusable shows/moves the cursor at the end of a screen
	// line's content in the line's last cell (default End).
	LastCellFocusable

	// NextCellFocusable shows the cursor in a components focused
	// line while the right arrow key moves the cursor to ("focuses") the
	// next cell.
	NextCellFocusable

	// editable makes a component's content editable by the user.
	editable

	// NoFeature classifies keys/runes/buttons not registered for any
	// feature.
	NoFeature FeatureMask = 0

	// Scrollable makes a component's content vertically Scrollable by
	// combining up- and down-Scrollable.
	Scrollable = UpScrollable | DownScrollable

	// LinesFocusable makes lines focusable, i.e. a line receiving the
	// focus is reported; see OnLineFocus.
	LinesFocusable = NextLineFocusable | PreviousLineFocusable |
		LineUnfocusable

	// LinesSelectable makes a component's lines selectable by combining
	// LinesFocusable and LineSelectable.
	LinesSelectable = LinesFocusable | LineSelectable

	// CellFocusable turns on LinesFocusable for a component c and shows
	// the cursor whose positioning indicates the "focused cell".
	CellFocusable = PreviousCellFocusable | NextCellFocusable |
		LinesFocusable | LastCellFocusable | FirstCellFocusable

	Editable = Focusable | CellFocusable | Scrollable | editable
)

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

// Set adds the default key, rune and button bindings of given
// feature(s) for associated component.
func (ff *Features) Set(f FeatureMask) {
	ff.c.ensureFeatures()
	if f&editable == editable && !ff.c.isNesting() && ff.c.Edit == nil {
		_, _, hasCursor := ff.c.CursorPosition()
		ff.c.Edit = &Editor{c: ff.c, suspended: !hasCursor}
		if hasCursor {
			ff.c.LL.Focus.EolAfterLastRune()
		}
	}
	ff.c.ff.set(f)
}

// OfKey returns the feature bound to given key k with given modifiers
// mm.
func (ff *Features) OfKey(k Key, mm ModifierMask) FeatureMask {
	return ff.c.ff.keyFeature(k, mm)
}

// OfRune returns the feature bound to given rune r with given modifiers
// mm.
func (ff *Features) OfRune(r rune, mm ModifierMask) FeatureMask {
	return ff.c.ff.runeFeature(r, mm)
}

// OfButton returns the feature bound to given buttons b with given modifiers
// mm.
func (ff *Features) OfButton(
	bb ButtonMask, mm ModifierMask,
) FeatureMask {
	return ff.c.ff.buttonFeature(bb, mm)
}

// Has returns true if receiving component features have key, rune or
// button bindings for given feature(s).
func (ff *Features) Has(f FeatureMask) bool {
	return ff.c.ff.has(f)
}

// All returns all features for which currently key, rune or button
// bindings are registered. (note [Features.Has] is faster to determine if a
// particular feature is set.)
func (ff *Features) All() FeatureMask {
	return ff.c.ff.all()
}

// KeysOf returns the keys with their modifiers bound to given feature
// of associated component.
func (ff *Features) KeysOf(f FeatureMask) FeatureKeys {
	return ff.c.ff.keysOf(f)
}

// SetKeysOf deletes all set keys for given feature and binds given keys
// to it instead.  Set*Of may be used instead of Add to bind features
// initially to other Runes/Keys/Buttons than default.  The call is
// ignored if given feature is not a power of two i.e. a single feature.
// Providing no keys simply removes all key-bindings for given feature.
func (ff *Features) SetKeysOf(f FeatureMask, kk ...FeatureKey) {
	ff.c.ensureFeatures().setKeysOf(f, kk...)
}

// ButtonsOf returns the buttons with their modifiers bound to given
// feature for associated component.
func (ff *Features) ButtonsOf(f FeatureMask) FeatureButtons {
	return ff.c.ff.buttonsOf(f)
}

// SetButtonsOf deletes all set buttons for given feature and binds
// given buttons to it instead.  Set*Of may be used instead of Add to
// bind features initially to other Runes/Keys/Buttons than default.
// The call is ignored if given feature is not a power of two i.e. a
// single feature.  Providing no button simply removes all
// button-bindings for given feature.
func (ff *Features) SetButtonsOf(f FeatureMask, bb ...FeatureButton) {
	ff.c.ensureFeatures().setButtonsOf(f, bb...)
}

// RunesOf returns the runes bound to given feature for associated
// component.
func (ff *Features) RunesOf(f FeatureMask) FeatureRunes {
	return ff.c.ff.runesOf(f)
}

// SetRunesOf deletes all set runes for given feature and binds given
// runes to it instead.  Set*Of may be used instead of Add to bind
// features initially to other Runes/Keys/Buttons than default.  The
// call is ignored if given feature is not a power of two i.e. a single
// feature.  Providing no runes simply removes all runes-bindings for
// given feature.
func (ff *Features) SetRunesOf(f FeatureMask, rr ...FeatureRune) {
	ff.c.ensureFeatures().setRunesOf(f, rr...)
}

// Delete removes all runes, key or button bindings of given feature(s)
// except for Quitable.  The two default Quitable bindings ctrl-c and
// ctrl-d remain.  NOTE use a *Kiosk constructor like [TermKiosk] for a
// [Lines]-instance to avoid having Quitable set by default.
func (ff *Features) Delete(f FeatureMask) {
	ff.c.ff.delete(f)
}

// FeatureMask classifies keys/runes/buttons for a components default
// behavior like focusable, scrollable etc.
type FeatureMask uint64
