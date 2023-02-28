// Copyright (c) 2022 Stephan Lukits. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

/*
editable_test drives the implementation for a components editable
feature.  A component that is editable must be focusable as well as its
lines and cells must be focusable.  Since focusable line comes also in a
highlighted variant editable must as well. Setting the editable feature
lets a user modify a components content of the screen-area.  Note a
component becomes automatically editable if a component source's Liner
implementation implements the EditLiner interface.
*/

package lines

import (
	"fmt"
	"strings"
	"testing"
	"time"

	. "github.com/slukits/gounit"
)

type _editable struct{ Suite }

func (s *_editable) SetUp(t *T) { t.Parallel() }

func (s *_editable) Component_has_editable_feature_set(t *T) {
	cmp := &cmpFX{
		onInit: func(cf *cmpFX, e *Env) { cf.FF.Set(Editable) },
	}
	fx := fx(t, cmp)

	t.FatalOn(fx.Lines.Update(cmp, nil, func(e *Env) {
		t.True(cmp.FF.Has(Editable))
	}))
}

func (s *_editable) Is_triggered_by_editable_source_liner(t *T) {
	cmp := &srcFX{liner: &editableLinerFX{}}
	fx := fx(t, cmp)

	t.FatalOn(fx.Lines.Update(cmp, nil, func(e *Env) {
		cmp.FF.Has(Editable)
		cmp.FF.Has(LinesFocusable)
		cmp.FF.Has(HighlightEnabled)
	}))
}

func (s *_editable) Component_is_focusable(t *T) {
	fx, cmp := fxCmpFF(t, Editable)
	t.FatalOn(fx.Lines.Update(cmp, nil, func(e *Env) {
		t.True(cmp.FF.Has(Focusable))
	}))
}

func (s *_editable) Component_has_focusable_lines(t *T) {
	fx, cmp := fxCmpFF(t, Editable)
	t.FatalOn(fx.Lines.Update(cmp, nil, func(e *Env) {
		t.True(cmp.FF.Has(LinesFocusable))
	}))
}

func (s *_editable) Component_s_focused_lines_are_unfocusable(t *T) {
	fx, cmp := fxCmpFF(t, Editable)
	t.FatalOn(fx.Lines.Update(cmp, nil, func(e *Env) {
		t.True(cmp.FF.Has(LineUnfocusable))
	}))
}

func (s *_editable) Component_has_focusable_cells(t *T) {
	fx, cmp := fxCmpFF(t, Editable)
	t.FatalOn(fx.Lines.Update(cmp, nil, func(e *Env) {
		t.True(cmp.FF.Has(CellFocusable))
	}))
}

func (s *_editable) Component_has_non_nil_edit_property(t *T) {
	fx_, cmp := fxCmpFF(t, Editable)
	t.FatalOn(fx_.Lines.Update(cmp, nil, func(e *Env) {
		t.True(cmp.Edit != nil)
	}))

	src := &srcFX{liner: (&editableLinerFX{}).initLines(0)}
	fx_ = fx(t, src)
	t.FatalOn(fx_.Lines.Update(src, nil, func(e *Env) {
		src.FF.Has(Editable)
		t.True(src.Edit != nil)
	}))
}

func (s *_editable) Component_s_editor_is_inactive_by_default(t *T) {
	fx_, cmp := fxCmpFF(t, Editable)
	t.FatalOn(fx_.Lines.Update(cmp, nil, func(e *Env) {
		t.True(cmp.Edit != nil)
		t.Not.True(cmp.Edit.IsActive())
	}))

	src := &srcFX{liner: (&editableLinerFX{}).initLines(0)}
	fx_ = fx(t, src)
	t.FatalOn(fx_.Lines.Update(src, nil, func(e *Env) {
		t.True(src.Edit != nil)
		t.Not.True(src.Edit.IsActive())
	}))
}

func (s *_editable) Edit_gaps_are_initially_negative(t *T) {
	fx, cmp := fxCmpFF(t, Editable)
	t.FatalOn(fx.Lines.Update(cmp, nil, func(e *Env) {
		t.True(cmp.Edit.LeftGap < 0 && cmp.Edit.RightGap < 0)
	}))
}

func (s *_editable) Editor_calculates_gaps_on_first_activation(t *T) {
	cmp := &cmpFX{onInit: func(c *cmpFX, e *Env) {
		c.FF.Set(Editable)
		fmt.Fprint(c.Gaps(0).Vertical, "")
		fmt.Fprint(c.Gaps(1).Right, "")
	}, onLayout: func(c *cmpFX, e *Env) { c.Edit.Resume() }}

	fx := fx(t, cmp)
	t.FatalOn(fx.Lines.Update(cmp, nil, func(e *Env) {
		t.Eq(1, cmp.Edit.LeftGap)
		t.Eq(2, cmp.Edit.RightGap)
	}))
}

func (s *_editable) Activation_is_noop_if_to_small_for_edit_gaps(t *T) {
	cmp := &cmpFX{onInit: func(c *cmpFX, e *Env) {
		fmt.Fprint(c.Gaps(0).Vertical, "")
		c.Dim().SetWidth(3)
	}}
	fx := fxFF(t, Editable, cmp)
	t.FatalOn(fx.Lines.Update(cmp, nil, func(e *Env) {
		cmp.Edit.Resume()
	}))
	t.FatalOn(fx.Lines.Update(cmp, nil, func(e *Env) {
		t.Not.True(cmp.Edit.IsActive())
	}))
}

func (s *_editable) Component_s_editor_is_activated_on_resume(t *T) {
	fx, cmp := fxCmpFF(t, Editable)
	t.FatalOn(fx.Lines.Update(cmp, nil, func(e *Env) {
		t.Not.True(cmp.Edit.IsActive())
		cmp.Edit.Resume()
	}))
	t.FatalOn(fx.Lines.Update(cmp, nil, func(e *Env) {
		t.True(cmp.Edit.IsActive())
	}))
}

func (s *_editable) Component_s_editor_is_activated_on_insert_key(t *T) {
	fx, cmp := fxCmpFF(t, Editable)
	t.FatalOn(fx.Lines.Update(cmp, nil, func(e *Env) {
		t.Not.True(cmp.Edit.IsActive())
	}))
	fx.FireKey(Insert)
	t.FatalOn(fx.Lines.Update(cmp, nil, func(e *Env) {
		t.True(cmp.Edit.IsActive())
	}))
}

func (s *_editable) Sourced_component_activates_edit_on_insert(t *T) {
	cmp := &srcFX{liner: &editableLinerFX{}}
	cmp.liner.(*editableLinerFX).cc = []string{""}
	fx := fx(t, cmp)
	t.FatalOn(fx.Lines.Update(cmp, nil, func(e *Env) {
		t.Not.True(cmp.Edit.IsActive())
	}))
	fx.FireKey(Insert)
	t.FatalOn(fx.Lines.Update(cmp, nil, func(e *Env) {
		t.True(cmp.Edit.IsActive())
	}))
	t.Eq("", fx.Screen().Trimmed())
}

func (s *_editable) Has_cursor_set_on_editor_activation(t *T) {
	fx, cmp := fxCmpFF(t, Editable)
	fx.FireKey(Insert)
	hasCursor := false
	t.FatalOn(fx.Lines.Update(cmp, nil, func(e *Env) {
		_, _, hasCursor = cmp.CursorPosition()
	}))
	t.True(hasCursor)
}

func (s *_editable) Component_s_editor_is_deactivated_on_esc(t *T) {
	fx, cmp := fxCmpFF(t, Editable)
	fx.FireKey(Insert)
	fx.FireKey(Esc)
	t.FatalOn(fx.Lines.Update(cmp, nil, func(e *Env) {
		t.Not.True(cmp.Edit.IsActive())
	}))
}

func (s *_editable) Component_loses_line_focus_on_dbl_esc(t *T) {
	fx, cmp := fxCmpFF(t, Editable)
	fx.FireKey(Insert)
	t.FatalOn(fx.Lines.Update(cmp, nil, func(e *Env) {
		t.True(cmp.LL.Focus.Current() != -1)
	}))
	fx.FireKey(Esc)
	t.FatalOn(fx.Lines.Update(cmp, nil, func(e *Env) {
		t.Not.True(cmp.Edit.IsActive())
		t.True(cmp.LL.Focus.Current() != -1)
	}))
	fx.FireKey(Esc)
	t.FatalOn(fx.Lines.Update(cmp, nil, func(e *Env) {
		t.True(cmp.LL.Focus.Current() == -1)
		_, _, hasCursor := cmp.CursorPosition()
		t.Not.True(hasCursor)
	}))
}

func (s *_editable) Editor_is_deactivated_on_loosing_line_focus(t *T) {
	fx, cmp := fxCmpFF(t, Editable)
	fx.FireKey(Insert)
	t.FatalOn(fx.Lines.Update(cmp, nil, func(e *Env) {
		t.True(cmp.Edit.IsActive())
		t.True(cmp.LL.Focus.IsActive())
	}))
	fx.FireKey(Up)
	t.FatalOn(fx.Lines.Update(cmp, nil, func(e *Env) {
		t.Not.True(cmp.Edit.IsActive())
		t.Not.True(cmp.LL.Focus.IsActive())
	}))
}

func (s *_editable) Reports_activation_to_component_listener(t *T) {
	activated := false
	cmp := &cmpFX{onEdit: func(c *cmpFX, e *Env, edt *Edit) bool {
		if edt.Type == Resume {
			activated = true
		}
		return false
	}}
	fx := fxFF(t, Editable, cmp)
	fx.FireKey(Insert)
	t.True(activated)
}

func (s *_editable) Suppresses_activation_on_listener_request(t *T) {
	cmp := &cmpFX{onEdit: func(c *cmpFX, e *Env, edt *Edit) bool {
		return edt.Type == Resume
	}}
	fx := fxFF(t, Editable, cmp)
	fx.FireKey(Insert)
	fx.Lines.Update(cmp, nil, func(e *Env) {
		t.Not.True(cmp.Edit.IsActive())
	})
}

func (s *_editable) Reports_activation_to_component_source(t *T) {
	cmp := &srcFX{liner: &editableLinerFX{}}
	cmp.liner.(*editableLinerFX).cc = []string{""}
	fx := fx(t, cmp)
	fx.FireKey(Insert)
	fx.Lines.Update(cmp, nil, func(e *Env) {
		t.True(cmp.Src.Liner.(*editableLinerFX).HasReported(Resume))
	})
}

func (s *_editable) Reports_deactivation_to_component_listener(t *T) {
	deactivated := false
	cmp := &cmpFX{onEdit: func(c *cmpFX, e *Env, edt *Edit) bool {
		if edt.Type == Suspend {
			deactivated = true
		}
		return false
	}}
	fx := fxFF(t, Editable, cmp)
	fx.FireKey(Insert)
	fx.FireKey(Esc)
	t.True(deactivated)
}

func (s *_editable) Reports_deactivation_to_component_source(t *T) {
	cmp := &srcFX{liner: &editableLinerFX{}}
	cmp.liner.(*editableLinerFX).cc = []string{""}
	fx := fx(t, cmp)
	fx.FireKey(Insert)
	fx.FireKey(Esc)
	t.FatalOn(fx.Lines.Update(cmp, nil, func(e *Env) {
		t.True(cmp.Src.Liner.(*editableLinerFX).HasReported(Suspend))
	}))
}

func (s *_editable) Reports_key_to_source_and_to_component(t *T) {
	cmpEdt := NoEdit
	cmp := &srcFX{liner: &editableLinerFX{}}
	cmp.onEdit = func(c *cmpFX, e *Env, edt *Edit) bool {
		cmpEdt = edt.Type
		return false
	}
	fx := fx(t, cmp)
	fx.FireKey(Insert)
	t.Eq(Resume, cmpEdt)
	t.FatalOn(fx.Lines.Update(cmp, nil, func(e *Env) {
		t.True(cmp.Src.Liner.(*editableLinerFX).HasReported(Resume))
	}))
	fx.FireKey(Esc)
	t.Eq(Suspend, cmpEdt)
	t.FatalOn(fx.Lines.Update(cmp, nil, func(e *Env) {
		t.True(cmp.Src.Liner.(*editableLinerFX).HasReported(Suspend))
	}))
}

func (s *_editable) Reports_rune_insert_to_component(t *T) {
	ins, exp := false, 'a'
	cmp := &cmpFX{onEdit: func(c *cmpFX, e *Env, edt *Edit) bool {
		if edt.Type == Ins {
			ins = true
			t.Eq(exp, edt.Rune)
		}
		return false
	}}
	fx := fxFF(t, Editable, cmp)
	fx.FireKey(Insert)
	fx.FireRune(exp)
	t.True(ins)
}

func (s *_editable) Reports_rune_insert_to_component_source(t *T) {
	exp := 'a'
	cmp := &srcFX{liner: &editableLinerFX{}}
	cmp.liner.(*editableLinerFX).cc = []string{""}
	fx := fx(t, cmp)
	fx.FireKey(Insert)
	fx.FireRune(exp)
	t.FatalOn(fx.Lines.Update(cmp, nil, func(e *Env) {
		t.True(cmp.Src.Liner.(*editableLinerFX).HasReported(Ins))
		t.True(cmp.Src.Liner.(*editableLinerFX).HasReportedRune(exp))
	}))
}

func (s *_editable) Appends_rune_on_eol_cursor_position(t *T) {
	fx, cmp := fxCmpFF(t, Editable, 20*time.Minute)
	exp := 'a'
	fx.FireKey(Insert)
	ln, cl, leftGap := -1, -1, -1
	t.FatalOn(fx.Lines.Update(cmp, nil, func(e *Env) {
		t.True(cmp.LL.Focus.Eol())
		ln, cl, _ = cmp.CursorPosition()
		_, _, _, leftGap = cmp.GapsLen()
	}))
	fx.FireRune(exp)
	t.Eq([]rune(fx.Screen()[ln])[cl+leftGap], exp)
}

func (s *_editable) Advances_cursor_position_on_rune_append(t *T) {
	fx, cmp := fxCmpFF(t, Editable)
	exp := 'a'
	fx.FireKey(Insert)
	ln, cl := -1, -1
	t.FatalOn(fx.Lines.Update(cmp, nil, func(e *Env) {
		t.True(cmp.LL.Focus.Eol())
		ln, cl, _ = cmp.CursorPosition()
	}))
	fx.FireRune(exp)
	afterLn, afterCl := -1, -1
	t.FatalOn(fx.Lines.Update(cmp, nil, func(e *Env) {
		t.True(cmp.LL.Focus.Eol())
		afterLn, afterCl, _ = cmp.CursorPosition()
	}))
	t.Eq(ln, afterLn)
	t.Eq(cl+1, afterCl)
}

func (s *_editable) Editor_indicates_line_overflow(t *T) {
	cmp := &cmpFX{onInit: func(c *cmpFX, e *Env) {
		c.Dim().SetWidth(5).SetHeight(1)
		fmt.Fprint(e, "aa")
	}}
	fx := fxFF(t, Editable, cmp)
	fx.FireKey(Insert) // activate editor
	t.FatalOn(fx.Lines.Update(cmp, nil, func(e *Env) {
		cmp.LL.Focus.LastCell()
	}))
	fx.FireRune('a')
	t.True(strings.HasPrefix(
		fx.ScreenOf(cmp)[0], string(OverflowLeft)+"aa "))
	fx.FireKey(Home)
	t.Eq(" aaa ", fx.ScreenOf(cmp)[0])
	fx.FireKey(End)
	t.True(strings.HasPrefix(
		fx.ScreenOf(cmp)[0], string(OverflowLeft)+"aa "))
	fx.FireRune('a')
	fx.FireKey(Home)
	t.Eq(" aaa"+string(OverflowRight), fx.ScreenOf(cmp)[0])
}

// func (s *_editable) Suppresses_rune_events_having_active_editor(t *T) {
// 	aRuneReceived := 0
// 	cmp := &cmpFX{
// 		onInit: func(c *cmpFX, e *Env) {
// 			c.Register.Rune('a', ZeroModifier, func(e *Env) {
// 				aRuneReceived++
// 			})
// 			// NOTE we need content to receive a cursor position
// 			fmt.Fprintf(e, "1st\n2nd\n3rd")
// 		},
// 	}
// 	fx := fx(t, cmp)
// 	fx.FireRune('a', ZeroModifier)
// 	t.Eq(1, aRuneReceived)
// 	fx.Lines.Update(cmp, nil, func(e *Env) { cmp.FF.Add(Editable) })
// 	fx.FireKey(Down) // set the cursor position to (0,0) which is ...
// 	// ... a prerequisite for having an active editor suppressing
// 	// reporting of rune events
// 	t.True(cmp.Edit.IsActive())
// 	fx.FireRune('a', ZeroModifier)
// 	t.Eq(1, aRuneReceived)
// }

// func (s *_editable) Doesnt_block_on_rune_reporting(t *T) {
// 	rr := []rune{}
// 	cmp := &cmpFX{
// 		onInit: func(cf *cmpFX, e *Env) {
// 			fmt.Fprintf(e, "1st\n2nd\n3rd")
// 		},
// 		onRune: func(c *cmpFX, e *Env, r rune, mm ModifierMask) {
// 			rr = append(rr, r)
// 		},
// 	}
// 	fx := fx(t, cmp)
// 	t.Eq(0, cmp.cc[onRune])
// 	fx.FireRune('a', ZeroModifier)
// 	t.Eq(1, cmp.cc[onRune])
// 	fx.Lines.Update(cmp, nil, func(e *Env) { cmp.FF.Set(Editable) })
// 	fx.FireKey(Down)
// 	fx.FireRune('b', ZeroModifier)
// 	t.Eq(2, cmp.cc[onRune])
// 	t.Eq("ab", string(rr))
// }

// func (s *_editable) Reports_insert(t *T) {
// 	reportedInsert := false
// 	cmp := &cmpFX{
// 		onInit: func(c *cmpFX, e *Env) {
// 			c.FF.Add(Editable)
// 			fmt.Fprintf(e, "1st\n2nd\n3rd")
// 		},
// 		onEdit: func(c *cmpFX, e *Env, edt *Edit) bool {
// 			t.Eq(0, edt.Line)
// 			t.Eq(0, edt.Cell)
// 			t.Eq('a', edt.Rune)
// 			if t.Eq(Ins, edt.Type) {
// 				reportedInsert = true
// 			}
// 			return true
// 		},
// 	}
// 	fx := fx(t, cmp)
// 	fx.FireKey(Down)
// 	fx.FireRune('a', ZeroModifier)
// 	t.True(reportedInsert)
// }

// func (s *_editable) Reports_replacement(t *T) {
// 	cmp := &cmpFX{
// 		onInit: func(c *cmpFX, e *Env) {
// 			c.FF.Add(Editable)
// 			c.Edit.Replacing()
// 			fmt.Fprintf(e, "1st\n2nd\n3rd")
// 		},
// 		onEdit: func(c *cmpFX, e *Env, edt *Edit) bool {
// 			t.Eq(0, edt.Line)
// 			t.Eq(0, edt.Cell)
// 			t.Eq('a', edt.Rune)
// 			t.Eq(Rpl, edt.Type)
// 			return true
// 		},
// 	}
// 	fx := fx(t, cmp)
// 	fx.FireKey(Down)
// 	fx.FireRune('a', ZeroModifier)
// 	t.Eq(1, cmp.cc[onEdit])
// }

// func (s *_editable) Reports_deletion(t *T) {
// 	cmp := &cmpFX{
// 		onInit: func(c *cmpFX, e *Env) {
// 			c.FF.Add(Editable)
// 			fmt.Fprintf(e, "1st\n2nd\n3rd")
// 		},
// 		onEdit: func(c *cmpFX, e *Env, edt *Edit) (omitEdit bool) {
// 			switch c.cc[onEdit] {
// 			case 1:
// 				t.Eq(0, edt.Line)
// 				t.Eq(0, edt.Cell)
// 			case 2:
// 				t.Eq(1, edt.Line)
// 				t.Eq(1, edt.Cell)
// 			}
// 			t.Eq(rune(0), edt.Rune)
// 			t.Eq(Del, edt.Type)
// 			return true // don't apply reported edit
// 		},
// 	}
// 	fx := fx(t, cmp)
// 	fx.FireKeys(Down, Right, Backspace)
// 	t.Eq(1, cmp.cc[onEdit])
// 	fx.FireKeys(Down, Delete)
// 	t.Eq(2, cmp.cc[onEdit])
// }

// func (s *_editable) Reports_join_deleting_preceding_line_break(t *T) {
// 	cmp := &cmpFX{
// 		onInit: func(c *cmpFX, e *Env) {
// 			c.FF.Add(Editable)
// 			fmt.Fprintf(e, "1st\n2nd\n3rd")
// 		},
// 		onEdit: func(c *cmpFX, e *Env, edt *Edit) (omitEdit bool) {
// 			t.Eq(1, edt.Line)
// 			t.Eq(0, edt.Cell)
// 			t.Eq(rune(0), edt.Rune)
// 			t.Eq(JoinPrev, edt.Type)
// 			return true // don't apply reported edit
// 		},
// 	}
// 	fx := fx(t, cmp)
// 	fx.FireKeys(Down, Down, Backspace)
// 	t.Eq(1, cmp.cc[onEdit])
// }

// func (s *_editable) Reports_join_deleting_following_line_break(t *T) {
// 	cmp := &cmpFX{
// 		onInit: func(c *cmpFX, e *Env) {
// 			c.FF.Add(Editable)
// 			fmt.Fprintf(e, "1st\n2nd\n3rd")
// 		},
// 		onEdit: func(c *cmpFX, e *Env, edt *Edit) (omitEdit bool) {
// 			t.Eq(0, edt.Line)
// 			t.Eq(3, edt.Cell)
// 			t.Eq(rune(0), edt.Rune)
// 			t.Eq(JoinNext, edt.Type)
// 			return true // don't apply reported edit
// 		},
// 	}
// 	fx := fx(t, cmp)
// 	fx.FireKeys(Down, Right, Right, Right, Delete)
// 	t.Eq(1, cmp.cc[onEdit])
// }

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
