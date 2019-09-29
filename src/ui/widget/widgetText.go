package widget

import (
	"fmt"
	"github.com/awesome-gocui/gocui"
	"github.com/easysoft/zentaoatf/src/utils/vari"
)

const (
	TextWidthFull = 69

	TextHeight = 2
)

type TextWidget struct {
	name string
	x, y int
	w    int
	h    int
	text string
}

func NewTextWidget(name string, x, y, w int, text string) *gocui.View {
	widget := TextWidget{name: name, x: x, y: y, w: w, text: text}
	v, _ := widget.Layout()
	return v
}

func NewTextareaWidget(name string, x, y, w, h int, text string) *gocui.View {
	widget := TextWidget{name: name, x: x, y: y, w: w, h: h, text: text}
	v, _ := widget.Layout()

	return v
}

func (w *TextWidget) Layout() (*gocui.View, error) {
	var h int
	if w.h <= 0 {
		h = TextHeight
		w.h = h
	} else {
		h = w.h
	}

	v, err := vari.Cui.SetView(w.name, w.x, w.y, w.x+w.w, w.y+h, 0)
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
