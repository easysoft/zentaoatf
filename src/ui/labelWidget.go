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

func NewLabelWidget(g *gocui.Gui, name string, x, y, w int, label string) *gocui.View {
	widget := LabelWidget{name: name, x: x, y: y, w: w, label: label}
	v, _ := widget.Layout(g)
	return v
}

func NewLabelWidgetAutoWidth(g *gocui.Gui, name string, x, y int, label string) *gocui.View {
	widget := &LabelWidget{name: name, x: x, y: y, w: len(label) + 1, label: label}
	v, _ := widget.Layout(g)
	return v
}

func (w *LabelWidget) Layout(g *gocui.Gui) (*gocui.View, error) {
	v, err := g.SetView(w.name, w.x, w.y, w.x+w.w, w.y+2)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return nil, err
		}

		fmt.Fprint(v, w.label)
	}
	return v, nil
}
