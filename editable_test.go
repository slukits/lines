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
	"testing"

	. "github.com/slukits/gounit"
)

type _editable struct{ Suite }

func (s *_editable) Activated_by_on_edit_source_liner(t *T) {
	t.TODO()
}

func (s *_editable) Component_is_focusable(t *T) {
	t.TODO()
}

func (s *_editable) Component_has_focusable_lines(t *T) {
	t.TODO()
}

func (s *_editable) Component_has_focusable_cells(t *T) {
	t.TODO()
}

func (s *_editable) Blocks_reporting_to_registered_rune_events(t *T) {
	t.TODO()
}

func (s *_editable) Doesnt_block_on_rune_reporting(t *T) {
	t.TODO()
}

func (s *_editable) Reports_insert(t *T) {
	t.TODO()
}

func (s *_editable) Reports_replacement(t *T) {
	t.TODO()
}

func (s *_editable) Reports_deletion(t *T) {
	t.TODO()
}

func (s *_editable) Updates_content_after_insert(t *T) {
	t.TODO()
}

func (s *_editable) Updates_content_after_replacement(t *T) {
	t.TODO()
}

func (s *_editable) Updates_content_after_deletion(t *T) {
	t.TODO()
}

func (s *_editable) Ignores_edit_if_canceled_by_component(t *T) {
	t.TODO()
}

func (s *_editable) Reports_insert_to_source_liner(t *T) {
	t.TODO()
}

func (s *_editable) Reports_replacement_to_source_liner(t *T) {
	t.TODO()
}

func (s *_editable) Reports_deletion_to_source_liner(t *T) {
	t.TODO()
}

func (s *_editable) Toggles_insert_replace_on_insert_key(t *T) {
	t.TODO()
}

func (s *_editable) Stops_input_processing_on_deactivation(t *T) {
	t.TODO()
}

func TestEditable(t *testing.T) {
	Run(&_editable{}, t)
}
