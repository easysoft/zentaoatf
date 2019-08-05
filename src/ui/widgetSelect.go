package ui

import (
	"fmt"
	"github.com/easysoft/zentaoatf/src/model"
	"github.com/easysoft/zentaoatf/src/utils"
	"github.com/jroimartin/gocui"
	"strings"
)

type SelectWidget struct {
	name    string
	x, y    int
	w       int
	h       int
	title   string
	options []model.Option
}

func NewSelectWidget(name string, x, y, w, h int, title string, options []model.Option) *gocui.View {
	widget := SelectWidget{name: name, x: x, y: y, w: w, h: h, title: title, options: options}
	v, _ := widget.Layout()

	return v
}

func (w *SelectWidget) Layout() (*gocui.View, error) {
	if w.h < 1 {
		w.h = 3
	}

	v, _ := utils.Cui.SetView(w.name, w.x, w.y, w.x+w.w, w.y+w.h)
	setViewScroll(w.name)
	setViewLineHighlight(w.name)

	v.Title = w.title

	labels := make([]string, 0)
	for _, opt := range w.options {
		labels = append(labels, opt.Name)
	}

	fmt.Fprint(v, strings.Join(labels, "\n"))
	return v, nil
}
