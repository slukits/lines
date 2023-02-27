// Copyright (c) 2022 Stephan Lukits. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package lines

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
	case editable:
		editorInsert(cntx, usr)
	}
}

func executeLineFocus(
	cntx *rprContext, usr Componenter, f func() (int, int),
) {
	clIdx, slIdx := usr.embedded().LL.Focus.Current(),
		usr.embedded().LL.Focus.Screen()
	_, column, _ := usr.embedded().CursorPosition()
	ln, cl := f()
	if clIdx == ln {
		if cl != column {
			reportCursorChange(cntx, usr)
		}
		return
	}
	reportLineFocus(cntx, usr, clIdx, slIdx)
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
		usr.embedded().cursorMoved = true
		reportCursorChange(cntx, usr)
	}
	reportLineFocus(cntx, usr, cIdx, sIdx)
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

func editorInsert(cntx *rprContext, usr Componenter) {
	if !usr.embedded().Edit.IsActive() {
		editor := usr.embedded().Edit
		edt := editor.newKeyEdit(
			cntx.ll.newKeyEvent(Insert, ZeroModifier))
		if usr.layoutComponent().wrapped().Src != nil {
			if sourcedEdit(cntx, usr, editor, edt) {
				return
			}
		}
		ls, ok := usr.(Editer)
		if ok {
			suppress := false
			callback(usr, cntx, func(e *Env) {
				suppress = ls.OnEdit(e, edt)
			})
			if suppress {
				return
			}
		}
		editor.edit(edt)
	}
}

func sourcedEdit(
	cntx *rprContext, usr Componenter, editor *Editor, edt *Edit,
) (reported bool) {
	ls, ok := usr.layoutComponent().wrapped().Src.Liner.(EditLiner)
	if !ok {
		return false
	}
	suppress := false
	callback(usr, cntx, func(e *Env) {
		suppress = ls.OnEdit(edt)
	})
	return suppress
}
