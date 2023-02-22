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
type StyleProperty struct {
	dropdown

	// Styles references the Styles drop-down-box a StyleProperty
	// instance is operating on.
	Styles *Styles
}

// OnInit sets up the property items and passes initialization through.
func (p *StyleProperty) OnInit(e *lines.Env) {
	p.listener = p
	if p.Styles != nil && p.Styles.Colors == Monochrome {
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

func (p *StyleProperty) OnUpdate(e *lines.Env, data interface{}) {
	if p.Styles == nil {
		p.items.OnUpdate(e, data)
		return
	}
	switch data.(Value) {
	case 0:
		if p.Styles.Colors == Monochrome {
			p.reverseMonochrome(e)
		}
	case 1:
		if p.Styles.Colors == Monochrome {
			p.Styles.setAttributeSelection()
			p.items.OnUpdate(e, data)
		}
	case 2:
		if p.Styles.Colors == Monochrome {
			p.resetMonochrome(e)
		}
	default: // zero selection
		if p.Styles.Colors == Monochrome && p.Styles.attrSelection {
			p.Styles.removeAttributeSelection()
			p.items.OnUpdate(e, data)
		}
	}
}

func (p *StyleProperty) reverseMonochrome(e *lines.Env) {
	sty := p.Styles.value.WithFG(p.Styles.value.BG()).
		WithBG(p.Styles.value.FG())
	p.Styles.Items = []string{lines.ColorNames[sty.FG()]}
	p.Styles.updateValue(e.Lines, sty)
	p.Styles.ss = p.Styles.calculateMonochromeStyles()
	p.items.OnUpdate(e, Value(NoDefault))
}

func (p *StyleProperty) resetMonochrome(e *lines.Env) {
	sty := p.Styles.initial
	p.Styles.Items = []string{lines.ColorNames[sty.FG()]}
	p.Styles.updateValue(e.Lines, sty)
	p.Styles.ss = p.Styles.calculateMonochromeStyles()
	p.items.OnUpdate(e, Value(NoDefault))
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

	attrSelection bool
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
func (p *Styles) OnInit(e *lines.Env) {
	p.Styler = p.itemStyle
	p.Highlighter = p.itemHighlight
	p.value = p.initialValue()
	if !p.initialized {
		p.initial = p.value
	}
	p.Items = p.calculateItems()
	p.ss = p.calculateStyles()
	p.MinWidth = p.calculateMinWidth()
	p.dropdown.OnInit(e)
}

func (p *Styles) itemStyle(idx int) lines.Style {
	if idx == LabelStyle {
		return p.value
	}
	if p.Colors == Monochrome {
		return p.value.WithAA(lines.ZeroStyle)
	}
	return lines.DefaultStyle // TODO: implement item style
}

func (p *Styles) OnUpdate(e *lines.Env, data interface{}) {
	if p.attrSelection && data.(Value) != -1 {
		p.updateAttribute(e, int(data.(Value)))
		return
	}
	p.items.OnUpdate(e, data)
}

func (p *Styles) updateAttribute(e *lines.Env, idx int) {
	p.value = p.value.Switch(aa[idx])
	if !strings.HasSuffix(p.Items[idx], SelectedAttr) {
		p.Items[idx] = p.Items[idx] + lines.Filler + SelectedAttr
	} else {
		p.Items[idx] = lines.StyleAttributeNames[aa[idx]]
	}
	p.items.OnUpdate(e, Value(NoDefault))
}

func (p *Styles) updateValue(ll *lines.Lines, s lines.Style) {
	p.value = s
	ll.Update(&p.items, nil, func(e *lines.Env) {
		p.resetItemsLabel(e)
	})
}

func (p *Styles) itemHighlight(sty lines.Style) lines.Style {
	return sty.WithFG(sty.BG()).WithBG(sty.FG())
}

// Value returns currently selected style.
func (p *Styles) Value() lines.Style {
	return p.value
}

func (p *Styles) initialValue() lines.Style {
	switch p.Colors {
	default:
		return lines.NewStyle(lines.ZeroStyle, lines.White,
			lines.Black)
	}
}

func (p *Styles) calculateItems() []string {
	switch p.Colors {
	default:
		return []string{lines.ColorNames[p.value.FG()]}
	}
}

func (p *Styles) calculateStyles() []lines.Style {
	if p.Colors == Monochrome {
		return p.calculateMonochromeStyles()
	}
	return nil
}

func (p *Styles) calculateMonochromeStyles() []lines.Style {
	return []lines.Style{p.value}
}

type attributes []lines.StyleAttributeMask

func (aa attributes) Labels() (ll []string) {
	for _, a := range aa {
		ll = append(ll, lines.StyleAttributeNames[a])
	}
	return ll
}

const SelectedAttr = "✓"

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

func (p *Styles) calculateMinWidth() int {
	decorator := len([]rune(SelectedAttr)) + 1
	width := decorator
	for _, n := range lines.StyleAttributeNames {
		rr := []rune(n)
		if len(rr)+decorator <= width {
			continue
		}
		width = len(rr) + decorator
	}
	clrWidth := p.calculateColorsWidth()
	if clrWidth+decorator > width {
		width = clrWidth + decorator
	}
	return width
}

func (p *Styles) calculateColorsWidth() int {
	var cc map[lines.Color]bool
	switch p.Colors {
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

func (p *Styles) setAttributeSelection() {
	p.attrSelection = true
	p.ZeroLabel = lines.ColorNames[p.value.FG()]
	p.Items = aa.Labels()
	p.DefaultItem = NoDefault
}

func (p *Styles) removeAttributeSelection() {
	p.attrSelection = false
	p.ZeroLabel, p.DefaultItem = "", 0
	p.Items = p.calculateItems()
}
