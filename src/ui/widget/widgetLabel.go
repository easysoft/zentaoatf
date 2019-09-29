package widget

import (
	"fmt"
	"github.com/awesome-gocui/gocui"
	"github.com/easysoft/zentaoatf/src/utils/vari"
)

const (
	LabelWidth  = 15
	LabelHeight = 2

	LabelWidthSmall = 10
	SelectWidth     = 20
)

type LabelWidget struct {
	name  string
	x, y  int
	w     int
	label string
}

func NewLabelWidget(name string, x, y int, label string) *gocui.View {
	widget := LabelWidget{name: name, x: x, y: y, w: LabelWidth, label: label}
	v, _ := widget.Layout()
	v.Frame = false
	return v
}

func NewLabelWidgetAutoWidth(name string, x, y int, label string) *gocui.View {
	widget := LabelWidget{name: name, x: x, y: y, w: len(label) + 1, label: label}
	v, _ := widget.Layout()
	v.Frame = false
	return v
}

func (w *LabelWidget) Layout() (*gocui.View, error) {
	v, err := vari.Cui.SetView(w.name, w.x, w.y, w.x+w.w, w.y+LabelHeight, 0)
	if err != nil {
		if !gocui.IsUnknownView(err) {
			return nil, err
		}

		fmt.Fprint(v, w.label)
	}
	return v, nil
}
