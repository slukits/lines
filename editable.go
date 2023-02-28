// Copyright (c) 2022 Stephan Lukits. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package lines

import "fmt"

type EditType int

const (
	NoEdit EditType = iota
	Resume
	Suspend
	Ins
	Replace
	Del
	JoinNext
	JoinPrev
)

var OverflowLeft = '…'
var OverflowRight = '…'

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
	// true the following application of the edit request is omitted.
	// Provided Edit instance holds the information about the requested
	// edit.
	OnEdit(*Env, *Edit) (suppressEdit bool)
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

	// LeftGap is calculated by Resume on its first call using the first
	// left gap-index which is not used yet.
	LeftGap int

	// RightGap is calculated by Resume on its first call using the first
	// right gap-index which is not used yet.
	RightGap int
}

// IsActive returns false if given Editor e is nil or suspended or if no
// cursor position is set; otherwise true is returned.
func (e *Editor) IsActive() bool {
	if e == nil || e.suspended {
		return false
	}
	return true
}

// Suspend deactivates a components (non nil) Editor e.
func (e *Editor) Suspend() {
	if e == nil {
		return
	}
	e.suspended = true
}

// Resume (re)activates a components (non nil) Editor e.
// NOTE the first Resume call takes no effect if the component has no
// layout; i.e. OnLayout or OnFocus are good places to call Resume.
// NOTE call Resume for the first time after needed gaps have been
// initialized to have Resume skip this gaps for its edit-gaps
// initialization.
func (e *Editor) Resume() {
	if e == nil {
		return
	}
	if e.LeftGap < 0 {
		_, e.RightGap, _, e.LeftGap = e.c.GapsLen()
		fmt.Fprint(e.c.Gaps(e.RightGap).Right, "")
		fmt.Fprint(e.c.Gaps(e.LeftGap).Left, "")
		e.c.LL.Focus.EolAfterLastRune()
	}
	if e.RightGap+e.LeftGap+2 > e.c.dim.Width() {
		return
	}
	e.suspended = false
	_, _, hasCursor := e.c.cursorPosition()
	if !hasCursor {
		(*e.c.ll).padded(0) // ensure at least one line
		ln, cl := e.c.LL.Focus.Next()
		if ln != 0 || cl != 0 {
			panic("lines: editor: cell at (0,0) must be focusable")
		}
	}
}

func (e *Editor) Replacing() {
	e.mode = Replace
}

func (e *Editor) IsReplacing() bool {
	return e.mode == Replace
}

func (e *Editor) newEdit(t EditType, ln, cl int, r rune) *Edit {
	return &Edit{Line: ln, Cell: cl, Type: t, Rune: r}
}

func (e *Editor) newKeyEdit(evt KeyEventer) *Edit {
	ln, cl, _ := e.c.cursorPosition()
	edt := &Edit{Line: ln, Cell: cl}
	switch evt.Key() {
	case Insert:
		edt.Type = Resume
	case Esc:
		edt.Type = Suspend
	case Left, Right, Up, Down, Home, End: // handled by key-features
		return nil
	default:
		return nil
	}
	return edt
}

func (e *Editor) newRuneEdit(evt RuneEventer) *Edit {
	ln, cl, _ := e.c.cursorPosition()
	return &Edit{Line: ln, Cell: cl, Type: Ins, Rune: evt.Rune()}
}

func (e *Editor) exec(edt *Edit) {
	switch edt.Type {
	case Resume:
		e.Resume()
	case Suspend:
		e.Suspend()
	case Ins:
		if e.c.LL.Focus.Eol() {
			e.c.LL.Focus.Line().appendRune(edt.Rune)
			e.c.LL.Focus.NextCell()
		}
	}
}

func (e *Editor) lineOverflow(left, right bool) {
	ol, or := ' ', ' '
	if left {
		ol = OverflowLeft
	}
	if right {
		or = OverflowRight
	}
	ln, _, _ := e.c.cursorPosition()
	Print(e.c.Gaps(e.LeftGap).Left.At(ln), ol)
	Print(e.c.Gaps(e.RightGap).Right.At(ln), or)
}

// delEdit translates a Backspace or Delete key press into a Del-Edit.
// Note if the first cell is backspaced its preceeding "line-break" is
// considered removed respectively if the last insert-cell is deleted
// the following "line-break" is considered removed, i.e. respective
// line joining Edit-instances are returned.  Is the cursor at the first
// content cell and received key is Backspace respectively at the last
// content cell and the received key is Delete nil is returned.
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
