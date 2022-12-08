package lyt

type Rect struct {
	x, y, w, h int
}

func NewRect(x, y, width, height int) *Rect {
	return &Rect{x: x, y: y, w: width, h: height}
}

func (r *Rect) Has(x, y int) bool {
	if x < r.x || x >= r.x+r.w || y < r.y || y >= r.y+r.h {
		return false
	}
	return true
}
