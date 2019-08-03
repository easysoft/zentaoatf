package ui

import (
	"fmt"
	"github.com/easysoft/zentaoatf/src/utils"
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

func NewButtonWidget(name string, x, y, w int, label string,
	handler func(g *gocui.Gui, v *gocui.View) error) *gocui.View {
	widget := ButtonWidget{name: name, x: x, y: y, w: w, label: label, handler: handler}

	v, _ := widget.Layout(handler)
	return v
}
func NewButtonWidgetNoBorder(name string, x, y, w int, label string,
	handler func(g *gocui.Gui, v *gocui.View) error) *gocui.View {
	widget := ButtonWidget{name: name, x: x, y: y, w: w, label: label, handler: handler}

	v, _ := widget.Layout(handler)
	v.Frame = false
	return v
}

func NewButtonWidgetAutoWidth(name string, x, y int, label string,
	handler func(g *gocui.Gui, v *gocui.View) error) *gocui.View {
	widget := NewButtonWidget(name, x, y, len(label)+3, " "+label+" ", handler)

	return widget
}
func NewButtonWidgetNoBorderAutoWidth(name string, x, y int, label string,
	handler func(g *gocui.Gui, v *gocui.View) error) *gocui.View {
	widget := NewButtonWidgetNoBorder(name, x, y, len(label)+3, " "+label+" ", handler)
	return widget
}

func (w *ButtonWidget) Layout(handler func(g *gocui.Gui, v *gocui.View) error) (*gocui.View, error) {
	v, err := utils.Cui.SetView(w.name, w.x, w.y, w.x+w.w, w.y+ButtonHeight)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return nil, err
		}

		if err := utils.Cui.SetKeybinding(w.name, gocui.KeyEnter, gocui.ModNone, w.handler); err != nil {
			return nil, err
		}
		if err := utils.Cui.SetKeybinding(w.name, gocui.MouseLeft, gocui.ModNone, w.handler); err != nil {
			return nil, err
		}

		fmt.Fprint(v, w.label)
	}
	return v, nil
}
