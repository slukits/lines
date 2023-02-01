// Copyright (c) 2022 Stephan Lukits. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package fx

import (
	"fmt"
)

func NStrings(n int) []string {
	if n < 0 {
		panic("fx: NStrings: negative number of strings")
	}
	cc := []string{}
	for i := 0; i < n; i++ {
		cc = append(cc, nthString(i+1))
	}
	return cc
}

func nthString(no int) string {
	switch no {
	case 1:
		return "1st"
	case 2:
		return "2nd"
	case 3:
		return "3rd"
	}
	return fmt.Sprintf("%dth", no)
}
