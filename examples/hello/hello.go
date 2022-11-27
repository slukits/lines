package main

import (
	"fmt"

	"github.com/slukits/lines"
)

type Cmp struct{ lines.Component }

func (c *Cmp) OnInit(e *lines.Env) {
	hello := "hello world"
	fmt.Fprint(e, hello)
	c.Dim().SetWidth(len([]rune(hello))).SetHeight(1)
}

func main() { lines.Term(&Cmp{}).WaitForQuit() }
