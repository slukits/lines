package main

import (
	"fmt"

	"github.com/slukits/lines"
	"github.com/slukits/lines/internal/api"
)

var term *lines.Lines

func main() {
	term = lines.Term(&cmp{})
	term.WaitForQuit()
}

type cmp struct{ lines.Component }

func (c *cmp) OnInit(e *lines.Env) {
	red, green, blue := []rune("red"), []rune("green"), []rune("blue")
	c.Dim().SetWidth(40).SetHeight(12)
	lines.Print(e.LL(0).At(0).FG(lines.White).BG(lines.Red), red)
	lines.Print(e.LL(0).At(3), []rune("  "))
	lines.Print(e.LL(0).At(5).FG(lines.Black).BG(lines.Green), green)
	lines.Print(e.LL(0).At(10), []rune("  "))
	lines.Print(e.LL(0).At(12).FG(lines.Yellow).BG(lines.Blue), blue)
	cm := colorMap(term.CurrentColors())
	for i := 0; i < 8; i++ {
		lines.Print(e.LL(i+1).At(0).BG(cm(i)).FG(cm(15)),
			[]rune(fmt.Sprint(" ", i+1, " ")))
		switch i {
		case 0:
			lines.Print(e.LL(i+1).At(12).BG(cm(i+8)).FG(cm(15)),
				[]rune(fmt.Sprint(" ", i+9, "  ")))
		default:
			lines.Print(e.LL(i+1).At(12).BG(cm(i+8)).FG(cm(15)),
				[]rune(fmt.Sprint(" ", i+9, " ")))
		}
	}
	tcll := lines.DBGTcell(term)
	lines.Print(e.LL(9).At(0), []rune(lines.Filler+fmt.Sprintf("Colors: %d", tcll.Colors())+lines.Filler))
	lines.Print(e.BG(cm(0)).LL(10).BG(cm(0)).FG(cm(15)).At(0), []rune("bg/fg"))
	lines.Print(e.LL(10).BG(cm(0)).At(6).FG(cm(4)), []rune("cmm/str"))
	lines.Print(e.LL(10).BG(cm(0)).At(15).FG(cm(2)).AA(lines.Bold), []rune("ident"))
	lines.Print(e.LL(10).BG(cm(0)).At(22).FG(cm(15)), []rune("op"))
	lines.Print(e.LL(10).BG(cm(0)).At(25).FG(cm(1)), []rune("kw"))
	fmt.Fprint(e.BG(cm(0)).FG(cm(15)).LL(10), "bg/fg")
	lines.Print(e.LL(10).At(6).BG(cm(0)).FG(cm(4)), []rune("cmm/str"))
	lines.Print(e.LL(10).At(15).BG(cm(0)).FG(cm(2)).AA(lines.Bold), []rune("ident"))
	lines.Print(e.LL(10).At(22).BG(cm(0)).FG(cm(15)), []rune("op"))
	lines.Print(e.LL(10).At(25).BG(cm(0)).FG(cm(1)), []rune("kw"))
	fmt.Fprint(e.BG(cm(15)).FG(cm(0)).LL(11), "bg/fg")
	lines.Print(e.LL(11).At(6).BG(cm(15)).FG(cm(3)), []rune("cmm/str"))
	lines.Print(e.LL(11).At(15).BG(cm(15)).FG(cm(4)), []rune("ident"))
	lines.Print(e.LL(11).At(22).BG(cm(15)).FG(cm(0)), []rune("op"))
	lines.Print(e.LL(11).At(25).BG(cm(15)).FG(cm(1)), []rune("kw"))
}

func colorMap(cc api.CCC) func(int) api.Color {
	return func(i int) api.Color {
		switch i {
		case 0:
			return cc.C1st
		case 1:
			return cc.C2nd
		case 2:
			return cc.C3rd
		case 3:
			return cc.C4th
		case 4:
			return cc.C5th
		case 5:
			return cc.C6th
		case 6:
			return cc.C7th
		case 7:
			return cc.C8th
		case 8:
			return cc.C9th
		case 9:
			return cc.C10th
		case 10:
			return cc.C11th
		case 11:
			return cc.C12th
		case 12:
			return cc.C13th
		case 13:
			return cc.C14th
		case 14:
			return cc.C15th
		case 15:
			return cc.C16th
		}
		return api.Color(0)
	}
}
