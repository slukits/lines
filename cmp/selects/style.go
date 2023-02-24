// Copyright (c) 2022 Stephan Lukits. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

/*
select implements two features a (style) property-selection and a style
selector.  The property-selection  allows to select the foreground color
the background color, style attributes and to reverse fore- and
background color.  The style selector allows to set one property at the
time and displays its label previewing the selected style.
*/

package selects

import (
	"strings"

	"github.com/slukits/lines"
)

const (
	Properties uint8 = iota
	ForegroundProperty
	BackgroundProperty
	ReverseFgBgProperty
	StyleAttributeProperty
	ResetProperties
)

var monoPP = []uint8{
	ReverseFgBgProperty, StyleAttributeProperty, ResetProperties}

var pp = []uint8{
	ForegroundProperty, BackgroundProperty,
	ReverseFgBgProperty, StyleAttributeProperty, ResetProperties,
}

var PropertyNames = map[uint8]string{
	Properties: "properties", ForegroundProperty: "foreground",
	BackgroundProperty: "background", ReverseFgBgProperty: "reverse",
	StyleAttributeProperty: "attributes", ResetProperties: "reset",
}

type PropertyType uint8

const (
	Foreground PropertyType = iota
	Background
	Attribute
)

type dropdown = DropDown

// StyleProperty is a drop-down-box for selecting the style aspect which
// should be set in associated *Styles* drop-down-box.
//
// NOTE this implementation of StyleProperty and Styles assumes that
// StyleProperty is initialized by lines before an associated Styles
// instanced is initialized.  I.e. it must come before in a wrapping
// Stacker or Chainer.
type StyleProperty struct {
	dropdown

	// Styles references the Styles drop-down-box a StyleProperty
	// instance is operating on.
	Styles *Styles

	value Value
}

// OnInit sets up the property items and passes initialization through.
func (p *StyleProperty) OnInit(e *lines.Env) {
	p.listener = p
	if p.Styles == nil {
		p.items.OnInit(e)
		return
	}
	if p.Styles.Colors == Monochrome {
		p.DefaultItem = NoDefault
		p.ZeroLabel = PropertyNames[Properties]
		for _, pn := range monoPP {
			p.Items = append(p.Items, PropertyNames[pn])
		}
		p.dropdown.OnInit(e)
		return
	}
	for _, pn := range pp {
		p.Items = append(p.Items, PropertyNames[pn])
	}
	p.dropdown.OnInit(e)
}

// OnUpdate controls an associated Styles-instance by either inverting
// its fore- and background colors, resetting it to its initial value or
// by setting the style aspect whose values are offered by the Styles
// instance's selection list.
func (p *StyleProperty) OnUpdate(e *lines.Env, data interface{}) {
	if p.Styles == nil {
		p.value = data.(Value)
		p.items.OnUpdate(e, data)
		return
	}
	switch data.(Value) {
	case 0:
		if p.Styles.Colors == Monochrome {
			p.invertMonochrome(e)
			break
		}
		p.Styles.setForegroundSelection()
		p.items.OnUpdate(e, data)
		p.value = data.(Value)
	case 1:
		if p.Styles.Colors == Monochrome {
			p.Styles.setAttributeSelection()
			p.items.OnUpdate(e, data)
			break
		}
		p.Styles.setBackgroundSelection()
		p.items.OnUpdate(e, data)
		p.value = data.(Value)
	case 2:
		if p.Styles.Colors == Monochrome {
			p.resetMonochrome(e)
			break
		}
		p.invertColored(e)
	case 3:
		p.Styles.setAttributeSelection()
		p.items.OnUpdate(e, data)
		p.value = data.(Value)
	case 4:
		p.resetColored(e)
	default: // zero selection
		if p.Styles.Colors == Monochrome && p.Styles.attrSelection {
			p.Styles.removeAttributeSelection()
			p.items.OnUpdate(e, data)
		}
	}
}

func (p *StyleProperty) invertMonochrome(e *lines.Env) {
	p.Styles.UpdateValue(e.Lines, p.Styles.Value().Invert())
	p.items.OnUpdate(e, Value(NoDefault))
}

func (p *StyleProperty) invertColored(e *lines.Env) {
	if p.Styles.Colors == System8Linux {
		v := p.Styles.Value().WithRemoved(lines.Bold)
		if _, ok := linuxFGBold[p.Styles.Value().FG()]; !ok {
			p.Styles.UpdateValue(e.Lines, v.Invert())
			p.items.OnUpdate(e, p.value)
			return
		}
		fgBg := linuxFGBold[v.FG()]
		if fgBg == v.BG() {
			p.items.OnUpdate(e, p.value)
			return
		}
		p.Styles.UpdateValue(e.Lines, v.WithBG(fgBg).WithFG(v.BG()))
		p.items.OnUpdate(e, p.value)
		return
	}
	p.Styles.UpdateValue(e.Lines, p.Styles.Value().Invert())
	p.items.OnUpdate(e, p.value)
}

func (p *StyleProperty) resetMonochrome(e *lines.Env) {
	p.Styles.UpdateValue(e.Lines, p.Styles.initial)
	p.items.OnUpdate(e, Value(NoDefault))
}

func (p *StyleProperty) resetColored(e *lines.Env) {
	p.Styles.UpdateValue(e.Lines, p.Styles.initial)
	p.items.OnUpdate(e, p.value)
}

// Styles is a style selector allowing at a given time to select one of
// the three style property: foreground color, background color or style
// attribute like bold or underline.  A Styles instance is usually used
// together with a StyleProperty instance which allows to set which of
// the three style should be provided for selection.  By default a
// Styles instance selects the foreground color.
type Styles struct {
	dropdown

	// Colors defines the number of back- and foreground colors which
	// may be selected in a  Styles instance.
	Colors ColorRange

	// SetProperty indicates which of the three style properties
	// foreground color, background color or style-attribute should be
	// set.  SetProperty may be interactively changed by an associated
	// StyleProperty instance.
	SetProperty PropertyType

	// initial style defaulting to a white foreground, a black
	// background and a zero style property.
	initial lines.Style

	value lines.Style

	ss []lines.Style

	initialized bool

	attrSelection bool // offer style-attributes for (un)setting

	fgSelection bool // offer foreground colors for selection

	bgSelection bool // offer background colors for selection
}

// SetInitial sets the initially set style of given Styles selection ss
// which in combination with a StyleProperty component will be the style
// *reset* will revert to.  SetInitial is a noop if given style is not
// in  the styles of ss.
func (ss *Styles) SetInitial(lines.Style) {
	ss.initialized = true
}

// OnInit calculates and sets the available fore- and background
// combinations or style-attributes according to SetProperty's value.
func (ss *Styles) OnInit(e *lines.Env) {
	ss.Styler = ss.itemStyle
	ss.Highlighter = ss.itemHighlight
	ss.value = ss.initialValue()
	if !ss.initialized {
		ss.initial = ss.value
	}
	ss.MinWidth = ss.calculateMinWidth()
	ss.ZeroLabel = RangeNames[ss.Colors]
	ss.setForegroundSelection()
	ss.dropdown.OnInit(e)
}

func (ss *Styles) itemStyle(idx int) lines.Style {
	if idx == LabelStyle {
		return ss.value
	}
	if ss.Colors == Monochrome {
		return ss.value.WithAA(lines.ZeroStyle)
	}
	if ss.fgSelection || ss.bgSelection {
		return ss.ss[idx]
	}
	return ss.value.Invert().WithAA(lines.ZeroStyle)
}

func (ss *Styles) itemHighlight(sty lines.Style) lines.Style {
	return sty.Invert()
}

// Value returns currently selected style.
func (ss *Styles) Value() lines.Style {
	return ss.value
}

// SelectingForeground returns true iff given Styles ss offers
// foreground colors for selection.
func (ss *Styles) SelectingForeground() bool { return ss.fgSelection }

// SelectingBackground returns true iff given Styles ss offers
// background colors for selection.
func (ss *Styles) SelectingBackground() bool { return ss.bgSelection }

// SelectingStyleAttributes return true iff given Styles ss offers style
// attributes for selection.
func (ss *Styles) SelectingStyleAttributes() bool {
	return ss.attrSelection
}

// OnUpdate modifies given styles ss value according to the selected
// item of the current selection which may be the foreground color,
// background color or style-attribute selection.
func (ss *Styles) OnUpdate(e *lines.Env, data interface{}) {
	if ss.attrSelection && data.(Value) != -1 {
		ss.updateAttribute(e, int(data.(Value)))
		return
	}
	if ss.fgSelection && data.(Value) != -1 {
		ss.updateFG(e, int(data.(Value)))
		return
	}
	if ss.bgSelection && data.(Value) != -1 {
		ss.updateBG(e, int(data.(Value)))
		return
	}
	ss.items.OnUpdate(e, data)
}

// UpdateValue sets the value of styles ss instance to given style s by
// sending an update-event request to given Lines instance ll.
func (ss *Styles) UpdateValue(ll *lines.Lines, s lines.Style) {
	ll.Update(&ss.items, nil, func(e *lines.Env) {
		ss.value = s
		ss.calculateStylesAndItems()
		ss.resetItemsLabel(e)
	})
}

func (ss *Styles) updateAttribute(e *lines.Env, idx int) {
	aa := aa
	if ss.Colors == System8Linux {
		aa = aaLinux
	}
	ss.value = ss.value.Switch(aa[idx])
	if !strings.HasSuffix(ss.Items[idx], SelectedMark) {
		ss.Items[idx] = ss.Items[idx] + lines.Filler + SelectedMark
	} else {
		ss.Items[idx] = lines.StyleAttributeNames[aa[idx]]
	}
	ss.items.OnUpdate(e, Value(NoDefault))
}

func (ss *Styles) updateFG(e *lines.Env, idx int) {
	currentIdx := ss.calculateCurrentStyleIndex()
	ss.value = ss.value.WithFG(ss.ss[idx].FG())
	if ss.Colors == System8Linux {
		if ss.ss[idx].AA()&lines.Bold != 0 {
			ss.value = ss.value.WithAdded(lines.Bold)
		} else {
			ss.value = ss.value.WithRemoved(lines.Bold)
		}
	}
	ss.Items[currentIdx] = ss.calculateFGName(currentIdx)
	ss.Items[idx] = ss.calculateFGName(idx) + lines.Filler + SelectedMark
	ss.items.OnUpdate(e, Value(NoDefault))
}

func (ss *Styles) updateBG(e *lines.Env, idx int) {
	currentIdx := ss.calculateCurrentStyleIndex()
	ss.value = ss.value.WithBG(ss.ss[idx].BG())
	ss.Items[currentIdx] = lines.ColorNames[ss.ss[currentIdx].BG()]
	ss.Items[idx] = lines.ColorNames[ss.ss[idx].BG()] +
		lines.Filler + SelectedMark
	ss.items.OnUpdate(e, Value(NoDefault))
}

func (ss *Styles) initialValue() lines.Style {
	switch ss.Colors {
	case System8:
		return lines.NewStyle(lines.ZeroStyle, lines.Silver,
			lines.Black)
	case System8Linux:
		return lines.NewStyle(lines.Bold, lines.Silver,
			lines.Black)
	default:
		return lines.NewStyle(lines.ZeroStyle, lines.White,
			lines.Black)
	}
}

type attributes []lines.StyleAttributeMask

// Labels returns the sequence of style-attribute names held by given
// attributes aa.
func (aa attributes) Labels() (ll []string) {
	for _, a := range aa {
		ll = append(ll, lines.StyleAttributeNames[a])
	}
	return ll
}

const SelectedMark = "✓"

var aa = attributes{
	lines.Bold,
	lines.Italic,
	lines.Underline,
	lines.StrikeThrough,
	lines.Dim,
	lines.Reverse,
	lines.Invalid,
	lines.Blink,
}

var aaLinux = attributes{
	lines.Dim,
	lines.Reverse,
	lines.Blink,
}

func (ss *Styles) setAttributeSelection() {
	ss.attrSelection, ss.bgSelection, ss.fgSelection = true, false, false
	if ss.Colors == System8Linux {
		ss.Items = aaLinux.Labels()
	} else {
		ss.Items = aa.Labels()
	}
	ss.DefaultItem = NoDefault
}

func (ss *Styles) setForegroundSelection() {
	ss.fgSelection, ss.bgSelection, ss.attrSelection = true, false, false
	ss.calculateStylesAndItems()
	ss.DefaultItem = NoDefault
}

func (ss *Styles) setBackgroundSelection() {
	ss.bgSelection, ss.fgSelection, ss.attrSelection = true, false, false
	ss.calculateStylesAndItems()
}

func (ss *Styles) removeAttributeSelection() {
	ss.attrSelection = false
	ss.ZeroLabel, ss.DefaultItem = "", 0
	ss.Items = ss.calculateItems()
}

func (ss *Styles) calculateStylesAndItems() {
	ss.ss = ss.calculateStyles()
	// NOTE item-calculation depends on preceding styles-calculation
	ss.Items = ss.calculateItems()
}

func (ss *Styles) calculateStyles() []lines.Style {
	if ss.attrSelection {
		return nil
	}
	switch ss.Colors {
	case System8:
		if ss.fgSelection {
			return System8Foregrounds(ss.value.BG())
		}
		return System8Backgrounds(ss.value.FG())
	case System8Linux:
		if ss.fgSelection {
			return LinuxForegrounds(ss.value.BG())
		}
		return LinuxBackgrounds(ss.value.FG())
	default:
		return ss.calculateMonochromeStyles()
	}
}

func (ss *Styles) calculateItems() (ii []string) {
	if ss.attrSelection {
		return aa.Labels()
	}
	if ss.fgSelection {
		return ss.calculateFGItems()
	}
	return ss.calculateBGItems()
}

func (ss *Styles) calculateFGName(styleIdx int) string {
	sty := ss.ss[styleIdx]
	if ss.Colors == System8Linux && sty.AA()&lines.Bold != 0 {
		return lines.ColorNames[linuxBoldFG[sty.FG()]]
	}
	return lines.ColorNames[sty.FG()]
}

func (ss *Styles) calculateCurrentStyleIndex() int {
	if ss.Colors != System8Linux || ss.value.AA()&lines.Bold == 0 ||
		ss.SelectingBackground() {
		for i, s := range ss.ss {
			if ss.SelectingBackground() {
				if s.BG() != ss.value.BG() {
					continue
				}
				return i
			}
			if s.FG() != ss.value.FG() {
				continue
			}
			return i
		}
	}
	for i, s := range ss.ss {
		if s.FG() != ss.value.FG() || s.AA()&lines.Bold == 0 {
			continue
		}
		return i
	}
	return 0
}

func (ss *Styles) calculateFGItems() (ii []string) {
	isLinux := ss.Colors == System8Linux
	valueIsBold := ss.value.AA()&lines.Bold != 0
	hasFg := func(sty lines.Style) bool {
		return (ss.value.FG() == sty.FG() && !isLinux) ||
			(ss.value.FG() == sty.FG() &&
				(valueIsBold == (sty.AA()&lines.Bold != 0)))
	}
	for _, s := range ss.ss {
		if isLinux && s.AA()&lines.Bold != 0 {
			name := lines.ColorNames[linuxBoldFG[s.FG()]]
			if s.FG() != ss.Value().FG() || ss.Value().AA()&lines.Bold == 0 {
				ii = append(ii, name)
				continue
			}
			ii = append(ii, name+lines.Filler+SelectedMark)
			continue
		}
		if hasFg(s) {
			ii = append(ii, lines.ColorNames[s.FG()]+
				lines.Filler+SelectedMark)
			continue
		}
		ii = append(ii, lines.ColorNames[s.FG()])
	}
	return ii
}

func (ss *Styles) calculateBGItems() (ii []string) {
	for _, s := range ss.ss {
		if s.BG() != ss.Value().BG() {
			ii = append(ii, lines.ColorNames[s.BG()])
			continue
		}
		ii = append(ii, lines.ColorNames[s.BG()]+
			lines.Filler+SelectedMark)
	}
	return ii
}

func (p *Styles) calculateMonochromeStyles() []lines.Style {
	return []lines.Style{p.value}
}

func (ss *Styles) calculateMinWidth() int {
	decorator := len([]rune(SelectedMark)) + 1
	width := decorator
	for _, n := range lines.StyleAttributeNames {
		rr := []rune(n)
		if len(rr)+decorator <= width {
			continue
		}
		width = len(rr) + decorator
	}
	clrWidth := ss.calculateColorsWidth()
	if clrWidth+decorator > width {
		width = clrWidth + decorator
	}
	return width
}

func (ss *Styles) calculateColorsWidth() int {
	var cc map[lines.Color]bool
	switch ss.Colors {
	case Monochrome:
		cc = monoColors
	}
	width := 0
	for c := range cc {
		cn := lines.ColorNames[c]
		if len(cn) <= width {
			continue
		}
		width = len(cn)
	}
	return width
}
