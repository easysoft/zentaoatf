package ui

import (
	"fmt"
	"github.com/jroimartin/gocui"
)

const (
	labelWidth = 15
)

type LabelWidget struct {
	name  string
	x, y  int
	w     int
	label string
}

func NewLabelWidget(g *gocui.Gui, name string, x, y, w int, label string) {
	widget := LabelWidget{name: name, x: x, y: y, w: w, label: label}
	widget.Layout(g)
}

func NewLabelWidgetAutoWidth(g *gocui.Gui, name string, x, y int, label string) {
	widget := &LabelWidget{name: name, x: x, y: y, w: len(label) + 1, label: label}
	widget.Layout(g)
}

func (w *LabelWidget) Layout(g *gocui.Gui) error {
	v, err := g.SetView(w.name, w.x, w.y, w.x+w.w, w.y+2)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		fmt.Fprint(v, w.label)
	}
	return nil
}
