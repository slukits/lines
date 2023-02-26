// Copyright (c) 2022 Stephan Lukits. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package fx

import (
	"testing"
	"time"

	. "github.com/slukits/gounit"
)

type FX struct{ Suite }

func (s *FX) SetUp(t *T) { t.Parallel() }

func (s *FX) Timeout_defaults_to_ten_seconds(t *T) {
	t.Eq(10*time.Second, New(t, &Cmp{}).TimeOut())
}

func TestFX(t *testing.T) {
	t.Parallel()
	Run(&FX{}, t)
}
