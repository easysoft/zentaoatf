package widget

import (
	"fmt"
	commConsts "github.com/aaronchen2k/deeptest/internal/comm/consts"
	"github.com/awesome-gocui/gocui"
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

	v, err := commConsts.Cui.SetView(w.name, w.x, w.y, w.x+w.w, w.y+w.h, 0)
	v.Highlight = false
	if err != nil {
		if !gocui.IsUnknownView(err) {
			return nil, err
		}

		fmt.Fprint(v, w.body)
	}

	return v, nil
}
