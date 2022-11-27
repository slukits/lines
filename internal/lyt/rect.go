package lyt

type Rect struct {
	x, y, w, h int
}

func NewRect(x, y, width, height int) *Rect {
	return &Rect{x: x, y: y, w: width, h: height}
}
