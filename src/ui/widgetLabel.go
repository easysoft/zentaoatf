package ui

import (
	"fmt"
	"github.com/jroimartin/gocui"
)

const (
	LabelWidth  = 15
	LabelHeight = 2
)

type LabelWidget struct {
	name  string
	x, y  int
	w     int
	label string
}

func NewLabelWidget(g *gocui.Gui, name string, x, y int, label string) *gocui.View {
	widget := LabelWidget{name: name, x: x, y: y, w: LabelWidth, label: label}
	v, _ := widget.Layout(g)
	v.Frame = false
	return v
}

func NewLabelWidgetAutoWidth(g *gocui.Gui, name string, x, y int, label string) *gocui.View {
	widget := LabelWidget{name: name, x: x, y: y, w: len(label) + 1, label: label}
	v, _ := widget.Layout(g)
	v.Frame = false
	return v
}

func (w *LabelWidget) Layout(g *gocui.Gui) (*gocui.View, error) {
	v, err := g.SetView(w.name, w.x, w.y, w.x+w.w, w.y+LabelHeight)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return nil, err
		}

		fmt.Fprint(v, w.label)
	}
	return v, nil
}
