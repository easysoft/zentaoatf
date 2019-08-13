package ui

import (
	"fmt"
	"github.com/easysoft/zentaoatf/src/utils/vari"
	"github.com/jroimartin/gocui"
)

type PanelWidget struct {
	name string
	x, y int
	w    int
	h    int
	body string
}

func NewPanelWidget(name string, x, y, w, h int, body string) *gocui.View {
	widget := PanelWidget{name: name, x: x, y: y, w: w, h: h, body: body}
	v, _ := widget.Layout()

	return v
}

func (w *PanelWidget) Layout() (*gocui.View, error) {
	if w.h < 1 {
		w.h = 3
	}

	v, err := vari.Cui.SetView(w.name, w.x, w.y, w.x+w.w, w.y+w.h)
	v.Highlight = false
	if err != nil {
		if err != gocui.ErrUnknownView {
			return nil, err
		}

		fmt.Fprint(v, w.body)
	}
	return v, nil
}
