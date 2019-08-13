package ui

import (
	"fmt"
	"github.com/easysoft/zentaoatf/src/utils/vari"
	"github.com/jroimartin/gocui"
	"strings"
)

const (
	HelpGlobal = `KEYBINDINGS
			Mouse: Menu operation on left side
			Tab: Move between form widgets
			Space: Toggle radio box
			Enter: Click button
			^H: Show/Hide help window
			^C: Exit`
)

type HelpWidget struct {
	name string
	x, y int
	w, h int
	body string
}

func NewHelpWidget() {
	maxX, _ := vari.Cui.Size()

	lines := strings.Split(HelpGlobal, "\n")

	w := 0
	for _, l := range lines {
		if len(l) > w {
			w = len(l)
		}
	}
	h := len(lines) + 1
	w = w + 2

	help := HelpWidget{name: "help", x: maxX - w - 3, y: 1, w: w, h: h + 1, body: HelpGlobal}
	help.Layout()
}

func (w *HelpWidget) Layout() error {
	v, err := vari.Cui.SetView(w.name, w.x, w.y, w.x+w.w, w.y+w.h)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		fmt.Fprint(v, w.body)
	}
	return nil
}

func ShowHelp(g *gocui.Gui, v *gocui.View) error {
	help, _ := g.View("help")

	if help != nil {
		HideHelp()
	} else {
		NewHelpWidget()
	}

	return nil
}

func HideHelp() error {
	help, _ := vari.Cui.View("help")

	if help != nil {
		if err := vari.Cui.DeleteView("help"); err != nil {
			return err
		}
	}

	return nil
}
