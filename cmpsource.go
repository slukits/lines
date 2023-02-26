// Copyright (c) 2022 Stephan Lukits. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package lines

// A Liner implementations is a different way to provide a component's
// content in case many lines should be associated with a component.
// This approach allows a component to control its features like
// scrolling and line focusing without bothering the client while only
// the number of screen lines of lines is stored in the component.
type Liner interface {

	// Print prints the line with given index idx to given line writer w
	// and returns true if there are lines with a greater index than
	// idx.
	Print(idx int, w *EnvLineWriter) bool
}

// ScrollableLiner implementations are [Liner] implementations
// triggering the scrolling feature of associated component, i.e. a
// content source with a scrollable liner makes its associated component
// automatically scrolling.
type ScrollableLiner interface {
	Liner

	// Len returns the total number of content lines a liner
	// implementation can provide to its associated component.
	Len() int
}

// FocusableLiner implementations are [ScrollableLiner] implementations
// triggering the lines-focusable feature of associated component, i.e.
// a content source with a focusable liner makes its associated
// component automatically lines-focusable.
type FocusableLiner interface {
	ScrollableLiner

	// IsFocusable returns true iff the line with given index idx is
	// focusable.
	IsFocusable(idx int) bool
}

// EditLiner implementations are [FocusableLiner] implementations
// turning the editable feature of associated component on.
type EditLiner interface {
	FocusableLiner

	// OnEdit implementation gets edit requests of a component's screen
	// cell reported and returns true iff the edit request should be
	// carried out.  Given line writer allows to print to edited line
	// while given Edit-instance provides the information about the
	// edit.
	OnEdit(w *EnvLineWriter, e *Edit) bool
}

// Highlighter provides a highlighter which may be set to a components
// globals.
type Highlighter interface {

	// Highlight implementation of a *Liner is set a Highlighter in the
	// globals of the component whose source has this liner.
	Highlight(Style) Style
}

// A ContentSource instance may be assigned to a [Component]'s Src
// property whose [Liner] is then used by the [Component] to print its
// content. E.g. if MyLiner is a Liner implementation and c a Component:
//
//	c.Src = &lines.ContentSource{Liner: &MyLiner{}}
//
// now c uses provided Liner instance to print its content.  NOTE
// according to a ContentSource's Liner implementation the setting of
// corresponding features is triggered.  E.g. is a Liner implementation
// a ScrollableLiner the component has the feature Scrollable set.
type ContentSource struct {

	// Liner provides a components content
	Liner

	// clean was chosen over dirty since the zero value of a content
	// source should be initially dirty, i.e. not clean.
	clean bool

	// init indicates if initializations have to be done which are
	// derived evaluation of given liner implementation.
	init bool

	first int
}

// IsDirty returns true if an initial content write from set Liner
// to associated component has happened; false otherwise.
func (cs *ContentSource) IsDirty() bool {
	if cs == nil {
		return false
	}
	return !cs.clean
}

func (cs *ContentSource) cleanup(c *component) {
	cs.clean = true
	if cs.Liner == nil {
		return
	}

	if !cs.init {
		cs.init = true
		cs.initialize(c)
	}

	n := c.ContentScreenLines()
	if n <= 0 {
		return
	}
	cs.sync(n, c)
}

func (cs *ContentSource) initialize(c *component) {
	c.ensureFeatures()
	if _, ok := cs.Liner.(EditLiner); ok {
		if !c.ff.has(Editable) {
			if c.userCmp.isEnabled() {
				c.userCmp.embedded().FF.Set(Editable)
			} else {
				c.userCmp.enable()
				c.userCmp.embedded().FF.Set(Editable)
				c.userCmp.disable()
			}
		}
		return
	}
	if sl, ok := cs.Liner.(ScrollableLiner); ok {
		if c.ContentScreenLines() < sl.Len() {
			if !c.ff.has(Scrollable) {
				c.ff.set(Scrollable)
			}
		}
	}
	if _, ok := cs.Liner.(FocusableLiner); ok {
		if !c.ff.has(LinesFocusable) {
			c.ff.set(LinesFocusable)
		}
	}
	if hl, ok := cs.Liner.(Highlighter); ok {
		c.globals().SetHighlighter(hl.Highlight)
	}
}

func (cs *ContentSource) sync(n int, c *component) {
	if cs == nil {
		return
	}
	idx := cs.first
	lw := &EnvLineWriter{
		inner: true, cmp: c.userCmp, line: idx - cs.first}
	for idx-cs.first < n && cs.Print(idx, lw) {
		idx++
		lw = &EnvLineWriter{
			inner: true, cmp: c.userCmp, line: idx - cs.first}
	}
}

func (cs *ContentSource) setFirst(idx int) {
	sl, ok := cs.Liner.(ScrollableLiner)
	if !ok || idx >= sl.Len() || idx < 0 {
		return
	}
	cs.first = idx
	cs.clean = false
}
