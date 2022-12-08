/*
Was a first stand alone hack for implementing linux console gpm mouse
support.
*/
package main

import (
	"fmt"
	"log"

	"github.com/gdamore/tcell/v2"
	gpmctl "github.com/jackdoe/go-gpmctl"
)

var content []string

func main() {
	scr, err := tcell.NewScreen()
	if err != nil {
		panic(fmt.Sprintf("can't obtain screen: %v", err))
	}
	if err := scr.Init(); err != nil {
		panic(fmt.Sprintf("can't initialize screen: %v", err))
	}
	watchGPM(scr)
	for {
		evt := scr.PollEvent()
		if evt == nil {
			return
		}
		switch evt := evt.(type) {
		case *tcell.EventResize:
			_, h := scr.Size()
			content = make([]string, h)
			scr.Sync()
		case *tcell.EventKey:
			// removeMouseCursor()
			if evt.Key() == tcell.KeyCtrlC {
				scr.Fini()
			}
		case *tcell.EventMouse:
			printMouseEvent(scr, evt)
			updateMouseCursor(scr, evt)
			scr.Show()
		}
	}
}

func init() { log.SetFlags(0) }

var lastX, lastY = -1, -1

func printMouseEvent(scr tcell.Screen, evt *tcell.EventMouse) {
	x, y := evt.Position()
	s := fmt.Sprintf("Buttons: %s, Modifiers: %s, Position: (%d,%d)",
		tcellButtonsToString(evt.Buttons()),
		tcellModifiersToString(evt.Modifiers()),
		x, y,
	)
	content = append([]string{s}, content[:len(content)-1]...)
	for y, l := range content {
		for x, r := range l {
			scr.SetContent(x, y, r, nil, tcell.StyleDefault)
		}
	}
}

func removeMouseCursor() {}

func updateMouseCursor(scr tcell.Screen, evt *tcell.EventMouse) {
	if lastX >= 0 {
		r, _, sty, _ := scr.GetContent(lastX, lastY)
		scr.SetContent(lastX, lastY, r, nil,
			sty.Attributes(0))
		scr.Show()
	}
	x, y := evt.Position()
	r, _, sty, _ := scr.GetContent(x, y)
	scr.SetContent(x, y, r, nil,
		sty.Attributes(tcell.AttrReverse))
	lastX, lastY = x, y
}

var gpmTcellModifiers = map[int]tcell.ModMask{
	1: tcell.ModShift,
	4: tcell.ModCtrl,
	8: tcell.ModAlt,
}

var gpmTcellButtons = map[gpmctl.Buttons]tcell.ButtonMask{
	gpmctl.B_LEFT:   tcell.Button1,
	gpmctl.B_RIGHT:  tcell.Button2,
	gpmctl.B_MIDDLE: tcell.Button3,
	gpmctl.B_FOURTH: tcell.Button4,
	gpmctl.B_UP:     tcell.WheelUp,
	gpmctl.B_DOWN:   tcell.WheelDown,
}

func watchGPM(scr tcell.Screen) {
	gpm, err := gpmctl.NewGPM(gpmctl.GPMConnect{
		EventMask:   gpmctl.ANY,
		DefaultMask: ^gpmctl.HARD,
		MinMod:      0,
		MaxMod:      ^uint16(0),
	})
	if err != nil {
		return
	}
	go func(scr tcell.Screen, gpm *gpmctl.GPM) {
		baseTypes := (gpmctl.MOVE | gpmctl.DOWN | gpmctl.UP | gpmctl.DRAG)
		for {
			evt, err := gpm.Read()
			if err != nil {
				continue
			}
			x, y := int(evt.X-1), int(evt.Y-1)
			switch evt.Type & baseTypes {
			case gpmctl.MOVE:
				scr.PostEvent(tcell.NewEventMouse(x, y,
					tcell.ButtonNone,
					gpmTcellModifiers[int(evt.Modifiers)],
				))
			case gpmctl.DOWN:
				scr.PostEvent(tcell.NewEventMouse(x, y,
					gpmTcellButtons[evt.Buttons],
					gpmTcellModifiers[int(evt.Modifiers)],
				))
			case gpmctl.UP:
				scr.PostEvent(tcell.NewEventMouse(x, y,
					tcell.ButtonNone,
					gpmTcellModifiers[int(evt.Modifiers)],
				))
			case gpmctl.DRAG:
				scr.PostEvent(tcell.NewEventMouse(x, y,
					gpmTcellButtons[evt.Buttons],
					gpmTcellModifiers[int(evt.Modifiers)],
				))
			}
		}
	}(scr, gpm)
}

var tcellButtons = map[tcell.ButtonMask]string{
	tcell.ButtonNone:      "none",
	tcell.ButtonPrimary:   "left",
	tcell.ButtonSecondary: "right",
	tcell.ButtonMiddle:    "middle",
	tcell.Button4:         "fourth",
	tcell.WheelUp:         "wheel-up",
	tcell.WheelDown:       "wheel-down",
}

func tcellButtonsToString(bb tcell.ButtonMask) string {
	if bb == tcell.ButtonNone {
		return tcellButtons[bb]
	}
	str := ""
	for b, s := range tcellButtons {
		if b == tcell.ButtonNone || bb&b != b {
			continue
		}
		if len(str) == 0 {
			str = s
			continue
		}
		str += "|" + s
	}
	return str
}

var tcellModifiers = map[tcell.ModMask]string{
	tcell.ModNone:  "none",
	tcell.ModShift: "shift",
	tcell.ModCtrl:  "ctrl",
	tcell.ModAlt:   "alt",
}

func tcellModifiersToString(mm tcell.ModMask) string {
	if mm == tcell.ModNone {
		return tcellModifiers[mm]
	}
	str := ""
	for m, s := range tcellModifiers {
		if m == tcell.ModNone || mm&m != m {
			continue
		}
		if len(str) == 0 {
			str = s
			continue
		}
		str += "|" + s
	}
	return str
}
