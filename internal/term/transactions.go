// Copyright (c) 2022 Stephan Lukits. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

/*
transactions provides the functionality for "transactional" posting of
event.  I.e. an event-post p does not return before all other
event-posts during the processing of p have been processed.
*/

package term

import (
	"fmt"
	"sync/atomic"
	"time"

	"github.com/slukits/lines/internal/api"
)

type transactional struct {
	ui      *UI
	timeout time.Duration
	count   atomic.Int32
	waiting chan bool
}

func (t *transactional) Post(evt api.Eventer) error {
	if t.count.CompareAndSwap(0, 1) {
		return t.postAndWait(evt)
	}
	t.count.Add(1)
	return t.ui.lib.PostEvent(evt)
}

func (t *transactional) postAndWait(evt api.Eventer) error {
	err := t.ui.lib.PostEvent(evt)
	for {
		select {
		case done := <-t.waiting:
			if done {
				return err
			}
		case <-time.After(t.timeout):
			return fmt.Errorf("post transactional %T: timeout", evt)
		}
	}
}

func (t *transactional) polled() {
	t.waiting <- t.count.Add(-1) == 0
}
