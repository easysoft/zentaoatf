package widget

import (
	"fmt"
	"github.com/jroimartin/gocui"
	"strings"
)

const (
	helpWidth = 15

	HelpGlobal = `KEYBINDINGS
				Tab: Move between widgets
				Space: Toggle radio box
				Enter: Click button
				^H: Show/Hide help
				^C: Exit`
)

type HelpWidget struct {
	name string
	x, y int
	w, h int
	body string
}

func NewHelpWidget(g *gocui.Gui) {
	maxX, _ := g.Size()

	lines := strings.Split(HelpGlobal, "\n")

	w := 0
	for _, l := range lines {
		if len(l) > w {
			w = len(l)
		}
	}
	h := len(lines) + 1
	w = w + 1

	help := HelpWidget{name: "help", x: maxX - w - 3, y: 1, w: w, h: h + 1, body: HelpGlobal}
	help.Layout(g)
}

func (w *HelpWidget) Layout(g *gocui.Gui) error {
	v, err := g.SetView(w.name, w.x, w.y, w.x+w.w, w.y+w.h)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		fmt.Fprint(v, w.body)
	}
	return nil
}
