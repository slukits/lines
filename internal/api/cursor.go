// Copyright (c) 2022 Stephan Lukits. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package api

// CursorStyle represents a given cursor style, which can include the shape and
// whether the cursor blinks or is solid.  Backends need to map their
// provided cursor styles to these cursor styles and set a requested
// cursor style to default if not supported.
type CursorStyle int

const (
	ZeroCursor CursorStyle = iota
	DefaultCursor
	BlockCursorBlinking
	BlockCursorSteady
	UnderlineCursorBlinking
	UnderlineCursorSteady
	BarCursorBlinking
	BarCursorSteady
)
