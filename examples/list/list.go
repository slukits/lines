package main

import (
	"github.com/slukits/lines"
	"github.com/slukits/lines/cmp/selects"
	"github.com/slukits/lines/examples/demo"
)

func main() {
	lines.Term(&listDemo{}).WaitForQuit()
}

type listDemo struct {
	lines.Component
	lines.Stacking
	titled *demo.Titled
}

var title = []rune("list-demo")

func (c *listDemo) OnInit(e *lines.Env) {
	c.titled = (&demo.Titled{Gapper: c, Title: title}).Single(e)
	c.Dim().SetWidth(80).SetHeight(24)
	c.CC = append(c.CC, newRow(&zeroList{}))
}

func newRow(c lines.Componenter, cc ...lines.Componenter) *row {
	r := &row{}
	r.CC = append(r.CC, c)
	r.CC = append(r.CC, cc...)
	return r
}

type row struct {
	lines.Component
	lines.Chaining
}

type zeroList struct {
	selects.ModalList
	titled *demo.Titled
	reflow bool
}

var listTitle = []rune("zero-list")

func (c *zeroList) OnInit(e *lines.Env) {
	c.titled = (&demo.Titled{Gapper: c, Title: listTitle}).Single(e)
	e.Lines.Update(&c.ModalList, nil, func(e *lines.Env) {
		c.ModalList.OnInit(e)
	})
}
