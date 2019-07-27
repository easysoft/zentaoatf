package widget

import (
	"fmt"
	"github.com/jroimartin/gocui"
)

const (
	ButtonWidth  = 15
	ButtonHeight = 2
)

type ButtonWidget struct {
	name    string
	x, y    int
	w       int
	label   string
	handler func(g *gocui.Gui, v *gocui.View) error
}

func NewButtonWidget(g *gocui.Gui, name string, x, y, w int, label string,
	handler func(g *gocui.Gui, v *gocui.View) error) *gocui.View {
	widget := ButtonWidget{name: name, x: x, y: y, w: w, label: label, handler: handler}
	v, _ := widget.Layout(g, handler)
	return v
}

func NewButtonWidgetAutoWidth(g *gocui.Gui, name string, x, y int, label string,
	handler func(g *gocui.Gui, v *gocui.View) error) *gocui.View {
	widget := NewButtonWidget(g, name, x, y, len(label)+3, " "+label+" ", handler)
	return widget
}

func (w *ButtonWidget) Layout(g *gocui.Gui, handler func(g *gocui.Gui, v *gocui.View) error) (*gocui.View, error) {
	v, err := g.SetView(w.name, w.x, w.y, w.x+w.w, w.y+ButtonHeight)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return nil, err
		}

		if err := g.SetKeybinding(w.name, gocui.KeyEnter, gocui.ModNone, w.handler); err != nil {
			return nil, err
		}
		if err := g.SetKeybinding(w.name, gocui.MouseLeft, gocui.ModNone, w.handler); err != nil {
			return nil, err
		}

		fmt.Fprint(v, w.label)
	}
	return v, nil
}
