package main

import (
	"github.com/slukits/lines"
)

type cmp struct{ lines.Component }

func (c *cmp) OnInit(e *lines.Env) {
	red, green, blue := []rune("red"), []rune("green"), []rune("blue")
	c.Dim().SetWidth(16).SetHeight(1)
	lines.Print(e.LL(0).At(0).FG(lines.White).BG(lines.Red), red)
	lines.Print(e.LL(0).At(3), []rune("  "))
	lines.Print(e.LL(0).At(5).FG(lines.Black).BG(lines.Green), green)
	lines.Print(e.LL(0).At(10), []rune("  "))
	lines.Print(e.LL(0).At(12).FG(lines.Yellow).BG(lines.Blue), blue)
}

func main() {
	lines.Term(&cmp{}).WaitForQuit()
}
