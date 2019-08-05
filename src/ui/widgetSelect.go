package ui

import (
	"fmt"
	"github.com/easysoft/zentaoatf/src/model"
	"github.com/easysoft/zentaoatf/src/utils"
	"github.com/jroimartin/gocui"
	"strings"
)

type SelectWidget struct {
	name  string
	x, y  int
	w     int
	h     int
	title string

	options []model.Option
	handler func(g *gocui.Gui, v *gocui.View) error
}

func NewSelectWidget(name string, x, y, w, h int, title string, options []model.Option,
	handler func(g *gocui.Gui, v *gocui.View) error) *gocui.View {
	widget := SelectWidget{name: name, x: x, y: y, w: w, h: h, title: title, options: options, handler: handler}
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
	setViewLineSelected(w.name, selectResultEvent)

	v.Title = w.title

	labels := make([]string, 0)
	for _, opt := range w.options {
		labels = append(labels, opt.Name)
	}

	fmt.Fprint(v, strings.Join(labels, "\n"))

	if err := utils.Cui.SetKeybinding(w.name, gocui.KeyEnter, gocui.ModNone, w.handler); err != nil {
		return nil, err
	}
	if err := utils.Cui.SetKeybinding(w.name, gocui.MouseLeft, gocui.ModNone, w.handler); err != nil {
		return nil, err
	}
	if err := utils.Cui.SetKeybinding(w.name, gocui.KeyArrowUp, gocui.ModNone, selectScrollEvent(-1, w)); err != nil {
		return nil, err
	}
	if err := utils.Cui.SetKeybinding(w.name, gocui.KeyArrowDown, gocui.ModNone, selectScrollEvent(1, w)); err != nil {
		return nil, err
	}

	return v, nil
}

func selectScrollEvent(dy int, w *SelectWidget) func(g *gocui.Gui, v *gocui.View) error {
	return func(g *gocui.Gui, v *gocui.View) error {
		scrollAction(v, dy)

		w.handler(filedValMap)

		return nil
	}
}
