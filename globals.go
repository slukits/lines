// Copyright (c) 2022 Stephan Lukits. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

/*
Globals describe default display properties like tab expansion width or
default style with which any component of the layout is initialized.
Special about globals is that their updates may be propagated.  E.g.

    ll := Term(nil)
    ll.Globals.SetBG(lines.HighlightStyle, lines.White)

will have the consequence that all components part of the layout and
having not explicitly an other background color for highlighting set
will get the highlight background color white.  Has a component's
globals instance explicitly an other background color for highlighting
set the propagation will be ignored.  Hence we have a very simple style
cascaded.

NOTE the propagation of a globals instance is enabled by setting its
propagation property. As of now a propagated value is not further
propagated if it is propagated to a globals instance having a
propagation property set.  The problem here is to find an implementation
guaranteeing to create no infinite propagation loop.
*/

package lines

// StyleType constants identify a style for a certain function, e.g.
// default fore- and background color and style attributes of a
// component, or fore-/background color and style attributes for a
// component's highlighted line.
type StyleType uint

const (
	Default StyleType = iota
	Highlight
)

type globaler interface{ globals() *globals }

// globals represents setup/behavior for a component's lines.
type globals struct {
	scr         *screen
	tabWidth    int
	ss          map[StyleType]Style
	updated     globalsUpdates
	ssUpdated   map[StyleType]globalStyleUpdates
	onUpdate    func(globalsUpdates, StyleType, globalStyleUpdates)
	propagation func(func(globaler))
	highlighter func(Style) Style
}

func newGlobals(propagation func(func(globaler))) *globals {
	gg := &globals{
		tabWidth: 4,
		ss: map[StyleType]Style{
			Default:   DefaultStyle,
			Highlight: DefaultStyle.WithAA(Reverse),
		},
		propagation: propagation,
	}
	gg.highlighter = defaultHighlighter(gg)
	return gg
}

// clone makes a copy of given globals gg without the propagation,
// updated and ssUpdated properties.
func (gg *globals) clone() *globals {
	cpy := &globals{
		scr:      gg.scr,
		tabWidth: gg.tabWidth,
		ss:       map[StyleType]Style{},
	}
	cpy.highlighter = defaultHighlighter(cpy)
	for k, v := range gg.ss {
		cpy.ss[k] = v
	}
	return cpy
}

// setCursor sets the cursor in a components given line at given column.
// Note setCursor's arguments are passed through to screen.setCursor
// whereas column becomes the x- and line the y-coordinate.
func (gg *globals) setCursor(line, column int, cs ...CursorStyle) {
	gg.scr.setCursor(column, line, cs...)
}

func (gg *globals) SetUpdateListener(
	l func(globalsUpdates, StyleType, globalStyleUpdates),
) *globals {
	gg.onUpdate = l
	return gg
}

func (gg *globals) SetHighlighter(h func(Style) Style) *globals {
	if h == nil {
		return gg.setDefaultHighlighter()
	}
	gg.highlighter = h
	if gg.updated&globalHighlighter == 0 {
		gg.updated |= globalHighlighter
	}
	if gg.onUpdate != nil {
		gg.onUpdate(globalHighlighter, 0, 0)
	}
	if gg.propagation == nil {
		return gg
	}
	gg.propagation(func(g globaler) {
		g.globals().prpHighlighter(h)
	})
	return gg
}

func (gg *globals) setDefaultHighlighter() *globals {
	gg.highlighter = defaultHighlighter(gg)
	if gg.updated&globalHighlighter == 0 {
		gg.updated |= globalHighlighter
	}
	if gg.onUpdate != nil {
		gg.onUpdate(globalHighlighter, 0, 0)
	}
	if gg.propagation == nil {
		return gg
	}
	gg.propagation(func(g globaler) {
		g.globals().prpHighlighter(defaultHighlighter(g.globals()))
	})
	return gg
}

func defaultHighlighter(gg *globals) func(s Style) Style {
	return func(s Style) Style {
		h := gg.Style(Highlight)
		if h.AA() != 0 {
			if s.AA()&h.AA() == 0 {
				s = s.WithAdded(h.AA())
			} else {
				s = s.WithRemoved(h.AA())
			}
		}
		if h.FG() != DefaultColor {
			s = s.WithFG(h.FG())
		}
		if h.BG() != DefaultColor {
			s = s.WithBG(h.BG())
		}
		return s
	}
}

func (gg *globals) prpHighlighter(h func(Style) Style) {
	if gg.updated&globalHighlighter != 0 {
		return
	}
	gg.highlighter = h
	if gg.onUpdate != nil {
		gg.onUpdate(globalHighlighter, 0, 0)
	}
}

func (gg *globals) Highlight(s Style) Style {
	return gg.highlighter(s)
}

// TabWidth returns the currently set tab-width in given globals gg.
func (gg *globals) TabWidth() int { return gg.tabWidth }

// SetTabWidth sets given width w as tab-width and propagates the change
// if propagation is set and returns given globals gg.  SetTabWidth is
// an no-op if w not positive.
func (gg *globals) SetTabWidth(w int) *globals {
	if w <= 0 {
		return gg
	}
	gg.tabWidth = w
	if gg.updated&globalTabWidth == 0 {
		gg.updated |= globalTabWidth
	}
	if gg.onUpdate != nil {
		gg.onUpdate(globalTabWidth, 0, 0)
	}
	if gg.propagation == nil {
		return gg
	}
	gg.propagation(func(g globaler) {
		g.globals().prpTabWidth(w)
	})
	return gg
}

func (gg *globals) prpTabWidth(w int) {
	if gg.updated&globalTabWidth != 0 {
		return
	}
	gg.tabWidth = w
	if gg.onUpdate != nil {
		gg.onUpdate(globalTabWidth, 0, 0)
	}
}

// AA returns the style attributes mask of given style type st in given
// globals gg.  If no style for st is found the default style's
// attributes are returned.
func (gg *globals) AA(st StyleType) StyleAttributeMask {
	if gg.ss == nil {
		return DefaultStyle.AA()
	}
	if _, ok := gg.ss[st]; !ok {
		return DefaultStyle.AA()
	}
	return gg.ss[st].AA()
}

// SetAA sets in given globals gg for given style type st given style
// attributes aa.
func (gg *globals) SetAA(st StyleType, aa StyleAttributeMask) *globals {
	if gg.ss == nil {
		gg.ss = map[StyleType]Style{}
	}
	if sty, ok := gg.ss[st]; !ok {
		gg.ss[st] = DefaultStyle.WithAA(aa)
	} else {
		gg.ss[st] = sty.WithAA(aa)
	}
	gg.setUpdated(st, glbStyAttribute)
	if gg.onUpdate != nil {
		gg.onUpdate(0, st, glbStyAttribute)
	}
	if gg.propagation == nil {
		return gg
	}
	sty := gg.ss[st]
	gg.propagation(func(g globaler) { g.globals().prpAA(st, aa, sty) })
	return gg
}

func (gg *globals) prpAA(
	st StyleType, aa StyleAttributeMask, sty Style,
) {
	if gg.ss == nil {
		gg.ss = map[StyleType]Style{}
	}
	if _, ok := gg.ss[st]; !ok {
		gg.ss[st] = sty
		if gg.onUpdate != nil {
			gg.onUpdate(0, st, glbStyAttribute)
		}
		return
	}
	if gg.ssUpdated[st]&glbStyAttribute != 0 {
		return
	}
	gg.ss[st] = gg.ss[st].WithAA(aa)
	if gg.onUpdate != nil {
		gg.onUpdate(0, st, glbStyAttribute)
	}
}

// FG returns the foreground color of given style type st in given
// globals gg.  If no style for st is found the default style's
// foreground color is returned.
func (gg *globals) FG(st StyleType) Color {
	if gg.ss == nil {
		return DefaultStyle.FG()
	}
	if _, ok := gg.ss[st]; !ok {
		return DefaultStyle.FG()
	}
	return gg.ss[st].FG()
}

// SetFG sets in given globals gg for given style type st given color c
// as foreground color.
func (gg *globals) SetFG(st StyleType, c Color) *globals {
	if gg.ss == nil {
		gg.ss = map[StyleType]Style{}
	}
	if sty, ok := gg.ss[st]; !ok {
		gg.ss[st] = DefaultStyle.WithFG(c)
	} else {
		gg.ss[st] = sty.WithFG(c)
	}
	gg.setUpdated(st, glbStyForeground)
	if gg.onUpdate != nil {
		gg.onUpdate(0, st, glbStyForeground)
	}
	if gg.propagation == nil {
		return gg
	}
	sty := gg.ss[st]
	gg.propagation(func(g globaler) { g.globals().prpFG(st, c, sty) })
	return gg
}

func (gg *globals) prpFG(st StyleType, c Color, sty Style) {
	if gg.ss == nil {
		gg.ss = map[StyleType]Style{}
	}
	if _, ok := gg.ss[st]; !ok {
		gg.ss[st] = sty
		if gg.onUpdate != nil {
			gg.onUpdate(0, st, glbStyForeground)
		}
		return
	}
	if gg.ssUpdated[st]&glbStyForeground != 0 {
		return
	}
	gg.ss[st] = gg.ss[st].WithFG(c)
	if gg.onUpdate != nil {
		gg.onUpdate(0, st, glbStyForeground)
	}
}

// BG returns the background color of given style type st in given
// globals gg.  If no style for st is found the default style's
// background color is returned.
func (gg *globals) BG(st StyleType) Color {
	if gg.ss == nil {
		return DefaultStyle.BG()
	}
	if _, ok := gg.ss[st]; !ok {
		return DefaultStyle.BG()
	}
	return gg.ss[st].BG()
}

// SetBG sets in given globals gg for given style type st given color c
// as background color.
func (gg *globals) SetBG(st StyleType, c Color) *globals {
	if gg.ss == nil {
		gg.ss = map[StyleType]Style{}
	}
	if sty, ok := gg.ss[st]; !ok {
		gg.ss[st] = DefaultStyle.WithBG(c)
	} else {
		gg.ss[st] = sty.WithBG(c)
	}
	gg.setUpdated(st, glbStyBackground)
	if gg.onUpdate != nil {
		gg.onUpdate(0, st, glbStyBackground)
	}
	if gg.propagation == nil {
		return gg
	}
	sty := gg.ss[st]
	gg.propagation(func(g globaler) { g.globals().prpBG(st, c, sty) })
	return gg
}

func (gg *globals) prpBG(st StyleType, c Color, sty Style) {
	if gg.ss == nil {
		gg.ss = map[StyleType]Style{}
	}
	if _, ok := gg.ss[st]; !ok {
		gg.ss[st] = sty
		if gg.onUpdate != nil {
			gg.onUpdate(0, st, glbStyBackground)
		}
		return
	}
	if gg.ssUpdated[st]&glbStyBackground != 0 {
		return
	}
	gg.ss[st] = gg.ss[st].WithBG(c)
	if gg.onUpdate != nil {
		gg.onUpdate(0, st, glbStyBackground)
	}
}

// SetStyle sets in given globals gg for given style type st given style
// sty.
func (gg *globals) SetStyle(st StyleType, sty Style) {
	if gg.ss == nil {
		gg.ss = map[StyleType]Style{}
	}
	gg.ss[st] = sty
	gg.setUpdated(st, glbStyAttribute|glbStyForeground|glbStyBackground)
	if gg.onUpdate != nil {
		gg.onUpdate(0, st, 0)
	}
	if gg.propagation == nil {
		return
	}
	gg.propagation(func(g globaler) { g.globals().prpStyle(st, sty) })
}

func (gg *globals) setUpdated(st StyleType, gsu globalStyleUpdates) {
	if gg.ssUpdated == nil {
		gg.ssUpdated = map[StyleType]globalStyleUpdates{}
	}
	gg.ssUpdated[st] |= gsu
}

// prpStyle merges given style sty for given style type st in to given
// globals gg styles.  I.e. only aspects of sty which haven't been set
// differently for gg are updated.
func (gg *globals) prpStyle(st StyleType, sty Style) {
	if gg.ss == nil {
		gg.ss = map[StyleType]Style{}
	}
	if _, ok := gg.ssUpdated[st]; !ok {
		gg.ss[st] = sty
		if gg.onUpdate != nil {
			gg.onUpdate(0, st, 0)
		}
		return
	}
	ssUpd, ss := gg.ssUpdated[st], gg.ss[st]
	for _, su := range allGlbStyAttribute {
		if ssUpd&su != 0 {
			continue
		}
		switch su {
		case glbStyAttribute:
			ss = ss.WithAA(sty.AA())
		case glbStyForeground:
			ss = ss.WithFG(sty.FG())
		case glbStyBackground:
			ss = ss.WithBG(sty.BG())
		}
	}
	gg.ss[st] = ss
	if gg.onUpdate != nil {
		gg.onUpdate(0, st, 0)
	}
}

// Style returns set style in given globals gg for given style type st.
// If no style is set the DefaultStyle is returned.
func (gg *globals) Style(st StyleType) Style {
	if gg.ss == nil {
		return DefaultStyle
	}
	if _, ok := gg.ss[st]; !ok {
		return DefaultStyle
	}
	return gg.ss[st]
}

// globalsUpdates are flags tracking which global properties of a
// component have been updated.
type globalsUpdates uint64

const (
	globalTabWidth globalsUpdates = 1 << iota
	globalFmt
	globalHighlighter
)

type globalStyleUpdates uint8

const (
	glbStyAttribute globalStyleUpdates = 1 << iota
	glbStyForeground
	glbStyBackground
)

var allGlbStyAttribute = []globalStyleUpdates{
	glbStyAttribute, glbStyForeground, glbStyBackground}

func globalsPropagationClosure(
	scr *screen,
) func(func(globaler)) {

	return func(f func(globaler)) {
		scr.forBaseComponents(func(lc layoutComponenter) {
			f(lc.wrapped())
		})
	}
}
