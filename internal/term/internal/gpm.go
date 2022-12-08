package internal

import (
	"os"

	"github.com/gdamore/tcell/v2"
	"github.com/jackdoe/go-gpmctl"
)

// scr enables to embedded tcell.Screen privately
type scr tcell.Screen

type gpmReader interface {
	Read() (gpmctl.Event, error)
}

// WrapGPMSupport either returns given screen scr and false if we are
// not in a linux console or no gpm support is found.  Otherwise scr is
// wrapped by GPMScreen reporting gpm mouse events and providing a mouse
// cursor.
func WarpGPMSupport(scr tcell.Screen, tt ...string) (
	_ tcell.Screen, haveGPM bool,
) {
	tt = append([]string{"linux"}, tt...)
	properTerminal, envterm := false, os.Getenv("TERM")
	for _, t := range tt {
		if t != envterm {
			continue
		}
		properTerminal = true
		break
	}
	if !properTerminal {
		return scr, false
	}
	gpm, err := gpmctl.NewGPM(gpmctl.GPMConnect{
		EventMask:   gpmctl.ANY,
		DefaultMask: ^gpmctl.HARD,
		MinMod:      0,
		MaxMod:      ^uint16(0),
	})
	if err != nil {
		return scr, false
	}
	scr = &GPMScreen{scr: scr, gpm: gpm}
	go gpmReporter(scr.(*GPMScreen))
	return scr, true
}

type GPMScreen struct {
	gpm gpmReader
	scr
	x, y int
}

func (gpm *GPMScreen) PollEvent() tcell.Event {
	evt := gpm.scr.PollEvent()
	switch evt := evt.(type) {
	case *tcell.EventMouse:
		// update the mouse cursor
		gpm.update(evt)
	case *tcell.EventKey:
		// remove mouse cursor on keyboard input
		gpm.reset()
	}
	return evt
}

func (gpm *GPMScreen) GetContent(x, y int) (
	primary rune, combining []rune, style tcell.Style, width int,
) {
	if x != gpm.x || y != gpm.y {
		return gpm.scr.GetContent(x, y)
	}

	primary, combining, style, width = gpm.scr.GetContent(x, y)
	return primary, combining, switchReverseAttribute(style), width
}

func (gpm *GPMScreen) SetContent(
	x, y int, primary rune, combining []rune, style tcell.Style,
) {
	if x != gpm.x || y != gpm.y {
		gpm.scr.SetContent(x, y, primary, combining, style)
		return
	}

	gpm.scr.SetContent(
		x, y, primary, combining, switchReverseAttribute(style))
}

func (gpm *GPMScreen) update(evt *tcell.EventMouse) {
	gpm.switchCursor(false)
	gpm.x, gpm.y = evt.Position()
	gpm.switchCursor(true)
}

func (gpm *GPMScreen) reset() {
	gpm.switchCursor(true)
	gpm.x, gpm.y = -1, -1
}

func (gpm *GPMScreen) switchCursor(show bool) {
	r, _, sty, _ := gpm.scr.GetContent(gpm.x, gpm.y)
	if r == 0 {
		return
	}
	sty = switchReverseAttribute(sty)
	gpm.scr.SetContent(gpm.x, gpm.y, r, nil, sty)
	if show {
		gpm.Show()
	}
}

func switchReverseAttribute(sty tcell.Style) tcell.Style {
	_, _, aa := sty.Decompose()
	if aa&tcell.AttrReverse != 0 {
		aa &^= tcell.AttrReverse
	} else {
		aa |= tcell.AttrReverse
	}
	return sty.Attributes(aa)
}

func gpmReporter(scr *GPMScreen) {
	baseTypes := (gpmctl.MOVE | gpmctl.DOWN | gpmctl.UP | gpmctl.DRAG)
	for {
		evt, err := scr.gpm.Read()
		if err != nil {
			continue
		}
		x, y := int(evt.X-1), int(evt.Y-1)
		switch evt.Type & baseTypes {
		case gpmctl.MOVE:
			scr.PostEvent(tcell.NewEventMouse(
				x, y,
				gpmButtonsToTcell(evt.Buttons),
				gpmModifiersToTcell(int(evt.Modifiers)),
			))
		case gpmctl.DOWN:
			scr.PostEvent(tcell.NewEventMouse(
				x, y,
				gpmButtonsToTcell(evt.Buttons),
				gpmModifiersToTcell(int(evt.Modifiers)),
			))
		case gpmctl.UP:
			scr.PostEvent(tcell.NewEventMouse(
				x, y,
				tcell.ButtonNone,
				gpmModifiersToTcell(int(evt.Modifiers)),
			))
		case gpmctl.DRAG:
			scr.PostEvent(tcell.NewEventMouse(
				x, y,
				gpmButtonsToTcell(evt.Buttons),
				gpmModifiersToTcell(int(evt.Modifiers)),
			))
		}
	}
}

var gpmTcellModifiers = map[int]tcell.ModMask{
	1: tcell.ModShift,
	4: tcell.ModCtrl,
	8: tcell.ModAlt,
}

func gpmModifiersToTcell(m int) (tm tcell.ModMask) {
	for _, _m := range []int{1, 4, 8} {
		if m&_m == 0 {
			continue
		}
		tm = tm | gpmTcellModifiers[_m]
	}
	return tm
}

var gpmTcellButtons = map[gpmctl.Buttons]tcell.ButtonMask{
	gpmctl.B_LEFT:   tcell.Button1,
	gpmctl.B_RIGHT:  tcell.Button2,
	gpmctl.B_MIDDLE: tcell.Button3,
	gpmctl.B_FOURTH: tcell.Button4,
	gpmctl.B_UP:     tcell.WheelUp,
	gpmctl.B_DOWN:   tcell.WheelDown,
}

var gpmButtons = []gpmctl.Buttons{gpmctl.B_LEFT, gpmctl.B_RIGHT, gpmctl.B_MIDDLE,
	gpmctl.B_FOURTH, gpmctl.B_UP, gpmctl.B_DOWN}

func gpmButtonsToTcell(b gpmctl.Buttons) (tb tcell.ButtonMask) {
	for _, _b := range gpmButtons {
		if b&_b != _b {
			continue
		}
		tb |= gpmTcellButtons[_b]
	}
	return
}
