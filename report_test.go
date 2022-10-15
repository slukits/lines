package lines

import (
	"testing"

	. "github.com/slukits/gounit"
)

type Report struct{ Suite }

// func (s *Report) Bracket_paste_to_bracket_paster_implementation(t *T) {
// 	t.TODO()
// }

func TestReport(t *testing.T) {
	t.Parallel()
	Run(&Report{}, t)
}
