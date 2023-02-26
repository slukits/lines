// Copyright (c) 2022 Stephan Lukits. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

/*
Package fx provides test fixtures for the cmp-packages.
*/
package fx

import (
	"time"

	"github.com/slukits/gounit"
	"github.com/slukits/lines"
)

var Factory = lines.TermFixture

// New creates new lines terminal test fixture.
func New(
	t *gounit.T, cmp lines.Componenter, timeout ...time.Duration,
) *lines.Fixture {
	d := time.Duration(0)
	if len(timeout) > 0 {
		d = timeout[0]
	}
	if cmp == nil {
		cmp = &Cmp{}
	}
	return Factory(t.GoT(), d, cmp)
}

func Sized(
	t *gounit.T,
	width, height int, cmp lines.Componenter,
	timeout ...time.Duration,
) *lines.Fixture {
	d := time.Duration(0)
	if len(timeout) > 0 {
		d = timeout[0]
	}
	if width <= 0 {
		width = 80
	}
	if height <= 0 {
		height = 25
	}
	if cmp == nil {
		cmp = &Cmp{}
	}
	fx := Factory(t.GoT(), d, &Cmp{})
	fx.FireResize(width, height)
	fx.Lines.SetRoot(cmp)
	return fx
}
