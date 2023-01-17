package main

import (
	"fmt"

	"github.com/slukits/lines"
)

type Cmp struct{ lines.Component }

func (c *Cmp) OnInit(e *lines.Env) {
	hello := "hello world on %d x %d cells"
	w, h := e.ScreenSize()
	fmt.Fprintf(e, hello, w, h)
	c.Dim().SetWidth(len([]rune(hello))).SetHeight(1)
}

func main() { lines.Term(&Cmp{}).WaitForQuit() }
