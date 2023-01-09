// Copyright (c) 2022 Stephan Lukits. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

/*
editable_test drives the implementation for a components editable
feature.  A component that is editable must be focusable as well as its
lines and cells must be focusable.  Since focusable line comes also in a
highlighted variant editable must as well. Setting the editable feature
lets a user modify a components content of the screen-area.  The cursor
keys move the insert cursor as provided by the selectable cell feature.
As usual a rune-key press triggers an OnRune event, i.e. if an editable
component implements this listener it can stop the bubbling of the rune
before it is processed by the editable feature.  The same holds true for
key-events.  Is the component implementation not suppressing runes they
are inserted at (before) the current cursor position unless the "insert"
key was pressed.  In the later case insert and replace are toggled, i.e.
the rune under the cursor is replaced.  Registered rune events are *not*
triggered if a component is editable.  But once the editable feature was
set deleting and resetting it becomes very cheap.  Backspace deletes the
rune left from the cursor while delete removes the rune under the
cursor.  The cursor can be moved one cell behind the last rune in order
to append runes at the end of the line.  After a cell content was
changed but not updated yet by lines (!) an OnEdit event is triggered.
The OnEdit event provides all information about the requested cell
change.  If OnEdit is not implemented lines just updates the cell as
requested.  If OnEdit is implemented lines ony updates the cell as
requested if the OnEdit implementation returns true otherwise it will do
nothing.  Has a component a source set the OnEdit event is reported to
its Liner implementation.  Note a component becomes automatically
editable if a component source's Liner implementation implements the
EditLiner interface.
*/

package lines

import (
	"fmt"
	"testing"
	"time"

	. "github.com/slukits/gounit"
)

type _editable struct{ Suite }

func (s *_editable) SetUp(t *T) { t.Parallel() }

func (s *_editable) Component_has_editable_feature_set(t *T) {
	cmp := &cmpFX{
		onInit: func(cf *cmpFX, e *Env) { cf.FF.Add(Editable) },
	}
	fx := fx(t, cmp)

	t.FatalOn(fx.Lines.Update(cmp, nil, func(e *Env) {
		t.True(cmp.FF.Has(Editable))
	}))
}

func (s *_editable) Component_has_non_nil_edit_property(t *T) {
	cmp := &cmpFX{}
	fx := fx(t, cmp)
	t.FatalOn(fx.Lines.Update(cmp, nil, func(e *Env) {
		t.Not.True(cmp.Edit != nil)
		cmp.FF.Add(Editable)
	}))
	t.FatalOn(fx.Lines.Update(cmp, nil, func(e *Env) {
		t.True(cmp.Edit != nil)
	}))
}

func (s *_editable) Component_s_editor_is_inactive_by_default(t *T) {
	cmp := &cmpFX{}
	fx := fx(t, cmp)
	t.FatalOn(fx.Lines.Update(cmp, nil, func(e *Env) {
		cmp.FF.Add(Editable)
	}))
	t.FatalOn(fx.Lines.Update(cmp, nil, func(e *Env) {
		t.Not.True(cmp.Edit.IsActive())
	}))
}

func (s *_editable) Component_s_editor_is_activated_on_insert(t *T) {
	fx, cmp := fxFF(t, Editable)
	t.FatalOn(fx.Lines.Update(cmp, nil, func(e *Env) {
		t.Not.True(cmp.Edit.IsActive())
	}))
	fx.FireKey(Insert)
	t.FatalOn(fx.Lines.Update(cmp, nil, func(e *Env) {
		t.True(cmp.Edit.IsActive())
	}))
}

type dbg struct{ Suite }

func (s *dbg) Dbg(t *T) {
	fx, cmp := fxFF(t, Editable, 20*time.Minute)
	t.FatalOn(fx.Lines.Update(cmp, nil, func(e *Env) {
		t.Not.True(cmp.Edit.IsActive())
	}))
	fx.FireKey(Insert)
	t.FatalOn(fx.Lines.Update(cmp, nil, func(e *Env) {
		t.True(cmp.Edit.IsActive())
	}))
}

func TestDBG(t *testing.T) { Run(&dbg{}, t) }

func (s *_editable) Component_is_focusable(t *T) {
	fx, cmp := fxFF(t, Editable)
	t.FatalOn(fx.Lines.Update(cmp, nil, func(e *Env) {
		t.True(cmp.FF.Has(Focusable))
	}))
}

func (s *_editable) Component_has_focusable_lines(t *T) {
	fx, cmp := fxFF(t, Editable)
	t.FatalOn(fx.Lines.Update(cmp, nil, func(e *Env) {
		t.True(cmp.FF.Has(LinesFocusable))
	}))
}

func (s *_editable) Component_s_focused_lines_are_unfocusable(t *T) {
	fx, cmp := fxFF(t, Editable)
	t.FatalOn(fx.Lines.Update(cmp, nil, func(e *Env) {
		t.True(cmp.FF.Has(LineUnfocusable))
	}))
}

func (s *_editable) Component_has_focusable_cells(t *T) {
	fx, cmp := fxFF(t, Editable)
	t.FatalOn(fx.Lines.Update(cmp, nil, func(e *Env) {
		t.True(cmp.FF.Has(CellFocusable))
	}))
}

type editLinerFX struct {
	focusLinerFX
}

func (l *editLinerFX) OnEdit(w *EnvLineWriter, e *Edit) bool {
	return true
}

func (s *_editable) Triggered_by_on_edit_source_liner(t *T) {
	el := &editLinerFX{}
	el.cc = []string{"1st", "2nd", "3rd", "4th"}
	fx, cmp := fxFF(t, Editable)
	t.FatalOn(fx.Lines.Update(cmp, nil, func(e *Env) {
		cmp.Src = &ContentSource{Liner: el}
	}))

	t.FatalOn(fx.Lines.Update(cmp, nil, func(e *Env) {
		cmp.FF.Has(Editable)
		cmp.FF.Has(HighlightedEditable)
	}))
}

func (s *_editable) Suppresses_rune_events_having_active_editor(t *T) {
	aRuneReceived := 0
	cmp := &cmpFX{
		onInit: func(c *cmpFX, e *Env) {
			c.Register.Rune('a', ZeroModifier, func(e *Env) {
				aRuneReceived++
			})
			// NOTE we need content to receive a cursor position
			fmt.Fprintf(e, "1st\n2nd\n3rd")
		},
	}
	fx := fx(t, cmp)
	fx.FireRune('a', ZeroModifier)
	t.Eq(1, aRuneReceived)
	fx.Lines.Update(cmp, nil, func(e *Env) { cmp.FF.Add(Editable) })
	fx.FireKey(Down) // set the cursor position to (0,0) which is ...
	// ... a prerequisite for having an active editor suppressing
	// reporting of rune events
	t.True(cmp.Edit.IsActive())
	fx.FireRune('a', ZeroModifier)
	t.Eq(1, aRuneReceived)
}

func (s *_editable) Doesnt_block_on_rune_reporting(t *T) {
	rr := []rune{}
	cmp := &cmpFX{
		onInit: func(cf *cmpFX, e *Env) {
			fmt.Fprintf(e, "1st\n2nd\n3rd")
		},
		onRune: func(c *cmpFX, e *Env, r rune, mm ModifierMask) {
			rr = append(rr, r)
		},
	}
	fx := fx(t, cmp)
	t.Eq(0, cmp.cc[onRune])
	fx.FireRune('a', ZeroModifier)
	t.Eq(1, cmp.cc[onRune])
	fx.Lines.Update(cmp, nil, func(e *Env) { cmp.FF.Add(Editable) })
	fx.FireKey(Down)
	fx.FireRune('b', ZeroModifier)
	t.Eq(2, cmp.cc[onRune])
	t.Eq("ab", string(rr))
}

func (s *_editable) Reports_insert(t *T) {
	reportedInsert := false
	cmp := &cmpFX{
		onInit: func(c *cmpFX, e *Env) {
			c.FF.Add(Editable)
			fmt.Fprintf(e, "1st\n2nd\n3rd")
		},
		onEdit: func(c *cmpFX, e *Env, edt *Edit) bool {
			t.Eq(0, edt.Line)
			t.Eq(0, edt.Cell)
			t.Eq('a', edt.Rune)
			if t.Eq(Ins, edt.Type) {
				reportedInsert = true
			}
			return true
		},
	}
	fx := fx(t, cmp)
	fx.FireKey(Down)
	fx.FireRune('a', ZeroModifier)
	t.True(reportedInsert)
}

func (s *_editable) Reports_replacement(t *T) {
	cmp := &cmpFX{
		onInit: func(c *cmpFX, e *Env) {
			c.FF.Add(Editable)
			c.Edit.Replacing()
			fmt.Fprintf(e, "1st\n2nd\n3rd")
		},
		onEdit: func(c *cmpFX, e *Env, edt *Edit) bool {
			t.Eq(0, edt.Line)
			t.Eq(0, edt.Cell)
			t.Eq('a', edt.Rune)
			t.Eq(Rpl, edt.Type)
			return true
		},
	}
	fx := fx(t, cmp)
	fx.FireKey(Down)
	fx.FireRune('a', ZeroModifier)
	t.Eq(1, cmp.cc[onEdit])
}

func (s *_editable) Reports_deletion(t *T) {
	cmp := &cmpFX{
		onInit: func(c *cmpFX, e *Env) {
			c.FF.Add(Editable)
			fmt.Fprintf(e, "1st\n2nd\n3rd")
		},
		onEdit: func(c *cmpFX, e *Env, edt *Edit) (omitEdit bool) {
			switch c.cc[onEdit] {
			case 1:
				t.Eq(0, edt.Line)
				t.Eq(0, edt.Cell)
			case 2:
				t.Eq(1, edt.Line)
				t.Eq(1, edt.Cell)
			}
			t.Eq(rune(0), edt.Rune)
			t.Eq(Del, edt.Type)
			return true // don't apply reported edit
		},
	}
	fx := fx(t, cmp)
	fx.FireKeys(Down, Right, Backspace)
	t.Eq(1, cmp.cc[onEdit])
	fx.FireKeys(Down, Delete)
	t.Eq(2, cmp.cc[onEdit])
}

func (s *_editable) Reports_join_deleting_preceeding_line_break(t *T) {
	cmp := &cmpFX{
		onInit: func(c *cmpFX, e *Env) {
			c.FF.Add(Editable)
			fmt.Fprintf(e, "1st\n2nd\n3rd")
		},
		onEdit: func(c *cmpFX, e *Env, edt *Edit) (omitEdit bool) {
			t.Eq(1, edt.Line)
			t.Eq(0, edt.Cell)
			t.Eq(rune(0), edt.Rune)
			t.Eq(JoinPrev, edt.Type)
			return true // don't apply reported edit
		},
	}
	fx := fx(t, cmp)
	fx.FireKeys(Down, Down, Backspace)
	t.Eq(1, cmp.cc[onEdit])
}

func (s *_editable) Reports_join_deleting_following_line_break(t *T) {
	cmp := &cmpFX{
		onInit: func(c *cmpFX, e *Env) {
			c.FF.Add(Editable)
			fmt.Fprintf(e, "1st\n2nd\n3rd")
		},
		onEdit: func(c *cmpFX, e *Env, edt *Edit) (omitEdit bool) {
			t.Eq(0, edt.Line)
			t.Eq(3, edt.Cell)
			t.Eq(rune(0), edt.Rune)
			t.Eq(JoinNext, edt.Type)
			return true // don't apply reported edit
		},
	}
	fx := fx(t, cmp)
	fx.FireKeys(Down, Right, Right, Right, Delete)
	t.Eq(1, cmp.cc[onEdit])
}

// func (s *_editable) Updates_content_after_insert(t *T) {
// 	t.TODO()
// }

// func (s *_editable) Updates_content_after_replacement(t *T) {
// 	t.TODO()
// }
//
// func (s *_editable) Updates_content_after_deletion(t *T) {
// 	t.TODO()
// }
//
// func (s *_editable) Ignores_edit_if_canceled_by_component(t *T) {
// 	t.TODO()
// }
//
// func (s *_editable) Reports_insert_to_source_liner(t *T) {
// 	t.TODO()
// }
//
// func (s *_editable) Reports_replacement_to_source_liner(t *T) {
// 	t.TODO()
// }
//
// func (s *_editable) Reports_deletion_to_source_liner(t *T) {
// 	t.TODO()
// }
//
// func (s *_editable) Toggles_insert_replace_on_insert_key(t *T) {
// 	t.TODO()
// }
//
// func (s *_editable) Stops_input_processing_on_deactivation(t *T) {
// 	t.TODO()
// }

func TestEditable(t *testing.T) {
	Run(&_editable{}, t)
}
