package ui

import (
	"fmt"
	"github.com/jroimartin/gocui"
)

type PanelWidget struct {
	name string
	x, y int
	w    int
	h    int
	body string
}

func NewPanelWidget(g *gocui.Gui, name string, x, y, w, h int, body string) *gocui.View {
	widget := PanelWidget{name: name, x: x, y: y, w: w, h: h, body: body}
	v, _ := widget.Layout(g)

	return v
}

func (w *PanelWidget) Layout(g *gocui.Gui) (*gocui.View, error) {
	if w.h < 1 {
		w.h = 3
	}

	v, err := g.SetView(w.name, w.x, w.y, w.x+w.w, w.y+w.h)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return nil, err
		}

		fmt.Fprint(v, w.body)
	}
	return v, nil
}
