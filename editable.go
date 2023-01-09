// Copyright (c) 2022 Stephan Lukits. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package lines

type EditType int

const (
	Ins EditType = iota
	Rpl
	Del
	JoinNext
	JoinPrev
)

type Edit struct {
	Line int
	Cell int
	Type EditType
	Rune rune
}

// Editer implementations are informed about user edits of a components
// content having the option to suppress the edit.
type Editer interface {

	// OnEdit is called right before a user requested edit of a
	// component's cell content is applied.  In case OnEdit returns
	// false the following application of the edit request is omitted.
	// Provided Edit instance holds the information about the requested
	// edit.
	OnEdit(*Env, *Edit) bool
}

// An Editor provides the client API to control a components editing
// behavior.  The zero-type is NOT ready to use; an Editor instance for
// a component c's Edit-property is automatically created iff c gets the
// Editable (or HighlightedEditable) feature set.  The later may happen
// in two ways: Either using c's FF-property to Add the feature
// explicitly or by setting c's Src property with a ContentSource having
// a Liner implementing the EditLiner interface.
type Editor struct {
	c         *Component
	suspended bool
	mode      EditType
}

// IsActive returns false if given Editor e is nil or suspended;
// otherwise true is returned.
func (e *Editor) IsActive() bool {
	if e == nil || e.suspended {
		return false
	}
	_, _, haveCursor := e.c.layoutCmp.wrapped().cursorPosition()
	return haveCursor
}

// Suspend deactivates a components (non nil) Editor e.
func (e *Editor) Suspend() {
	if e == nil {
		return
	}
	e.suspended = true
}

// Resume reactivates a components (non nil) Editor e.
func (e *Editor) Resume() {
	if e == nil {
		return
	}
	e.suspended = false
}

func (e *Editor) Replacing() {
	e.mode = Rpl
}

func (e *Editor) IsReplacing() bool {
	return e.mode == Rpl
}

func (e *Editor) MapEvent(evt KeyEventer) *Edit {
	switch evt.Key() {
	case Backspace, Delete:
		return e.delEdit(evt)
	}
	return nil
}

// delEdit translates a Backspace or Delete key press into a Del-Edit.
// Note if the first cell is backspaced its preceeding "line-break" is
// considered removed respectively if the last insert-cell is deleted
// the following "line-break" is considered removed, i.e. respective
// line joining Edit-instances are returned.  Is the cursor at the first
// content cell and received key is Backspace respectively at the last
// content cell and teh received key is Delete nil is returned.
func (e *Editor) delEdit(evt KeyEventer) *Edit {
	cmp := e.c.layoutCmp.wrapped()
	ln, cl, haveCursor := cmp.cursorPosition()
	if !haveCursor {
		panic("lines: report: on-edit: cursor position missing")
	}
	key := evt.Key()
	if e.isNothingToDelete(cmp, ln, cl, key) {
		return nil
	}
	edt := &Edit{
		Line: ln,
		Cell: cl,
		Rune: rune(0),
	}
	if cl == 0 && key == Backspace {
		edt.Type = JoinPrev
	}
	if (*cmp.ll)[ln].Len() <= cl && key == Delete {
		edt.Type = JoinNext
	}
	if edt.Type == 0 {
		if key == Backspace {
			edt.Cell--
		}
		edt.Type = Del
	}
	return edt
}

func (e *Editor) isNothingToDelete(
	cmp *component, line, cell int, key Key,
) bool {
	if line == 0 && cell == 0 && key == Backspace {
		return true
	}
	if line+1 < len(*cmp.ll) {
		return false
	}
	if cell >= (*cmp.ll)[line].Len() && key == Delete {
		return true
	}
	return false
}
