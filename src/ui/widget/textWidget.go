package widget

import (
	"fmt"
	"github.com/jroimartin/gocui"
)

const (
	TextWidthFull = 69
	TextWidthHalf = 25

	TextHeight = 2
)

type TextWidget struct {
	name string
	x, y int
	w    int
	text string
}

func NewTextWidget(g *gocui.Gui, name string, x, y, w int, text string) *gocui.View {
	widget := TextWidget{name: name, x: x, y: y, w: w, text: text}
	v, _ := widget.Layout(g)
	return v
}

func (w *TextWidget) Layout(g *gocui.Gui) (*gocui.View, error) {
	v, err := g.SetView(w.name, w.x, w.y, w.x+w.w, w.y+TextHeight)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return nil, err
		}

		v.Editable = true
		v.Wrap = true

		v.Clear()
		fmt.Fprint(v, w.text)
	}
	return v, nil
}
