package main

import (
	"github.com/slukits/lines"
)

var term *lines.Lines

func main() {
	term = lines.Term(&cmp{})
	term.WaitForQuit()
}

type cmp struct{ lines.Component }

// func (c *cmp) OnInit(e *lines.Env) {
// 	red, green, blue := []rune("red"), []rune("green"), []rune("blue")
// 	c.Dim().SetWidth(40).SetHeight(12)
// 	lines.Print(e.LL(0).At(0).FG(lines.White).BG(lines.Red), red)
// 	lines.Print(e.LL(0).At(3), []rune("  "))
// 	lines.Print(e.LL(0).At(5).FG(lines.Black).BG(lines.Green), green)
// 	lines.Print(e.LL(0).At(10), []rune("  "))
// 	lines.Print(e.LL(0).At(12).FG(lines.Yellow).BG(lines.Blue), blue)
// 	cm := colorMap(term.CurrentColors())
// 	for i := 0; i < 8; i++ {
// 		lines.Print(e.LL(i+1).At(0).BG(cm(i)).FG(cm(15)),
// 			[]rune(fmt.Sprint(" ", i+1, " ")))
// 		switch i {
// 		case 0:
// 			lines.Print(e.LL(i+1).At(12).BG(cm(i+8)).FG(cm(15)),
// 				[]rune(fmt.Sprint(" ", i+9, "  ")))
// 		default:
// 			lines.Print(e.LL(i+1).At(12).BG(cm(i+8)).FG(cm(15)),
// 				[]rune(fmt.Sprint(" ", i+9, " ")))
// 		}
// 	}
// 	tcll := lines.DBGTcell(term)
// 	lines.Print(e.LL(9).At(0), []rune(lines.Filler+fmt.Sprintf("Colors: %d", tcll.Colors())+lines.Filler))
// 	lines.Print(e.BG(cm(0)).LL(10).BG(cm(0)).FG(cm(15)).At(0), []rune("bg/fg"))
// 	lines.Print(e.LL(10).BG(cm(0)).At(6).FG(cm(4)), []rune("cmm/str"))
// 	lines.Print(e.LL(10).BG(cm(0)).At(15).FG(cm(2)).AA(lines.Bold), []rune("ident"))
// 	lines.Print(e.LL(10).BG(cm(0)).At(22).FG(cm(15)), []rune("op"))
// 	lines.Print(e.LL(10).BG(cm(0)).At(25).FG(cm(1)), []rune("kw"))
// 	fmt.Fprint(e.BG(cm(0)).FG(cm(15)).LL(10), "bg/fg")
// 	lines.Print(e.LL(10).At(6).BG(cm(0)).FG(cm(4)), []rune("cmm/str"))
// 	lines.Print(e.LL(10).At(15).BG(cm(0)).FG(cm(2)).AA(lines.Bold), []rune("ident"))
// 	lines.Print(e.LL(10).At(22).BG(cm(0)).FG(cm(15)), []rune("op"))
// 	lines.Print(e.LL(10).At(25).BG(cm(0)).FG(cm(1)), []rune("kw"))
// 	fmt.Fprint(e.BG(cm(15)).FG(cm(0)).LL(11), "bg/fg")
// 	lines.Print(e.LL(11).At(6).BG(cm(15)).FG(cm(3)), []rune("cmm/str"))
// 	lines.Print(e.LL(11).At(15).BG(cm(15)).FG(cm(4)), []rune("ident"))
// 	lines.Print(e.LL(11).At(22).BG(cm(15)).FG(cm(0)), []rune("op"))
// 	lines.Print(e.LL(11).At(25).BG(cm(15)).FG(cm(1)), []rune("kw"))
// }
