package widget

import (
	"fmt"
	"github.com/easysoft/zentaoatf/src/model"
	"github.com/easysoft/zentaoatf/src/ui"
	"github.com/easysoft/zentaoatf/src/utils/vari"
	"github.com/easysoft/zentaoatf/src/utils/zentao"
	"github.com/jroimartin/gocui"
	"strings"
)

type SelectWidget struct {
	name     string
	x, y     int
	w        int
	h        int
	title    string
	options  []model.Option
	defaultt string

	checkHandler func(g *gocui.Gui, v *gocui.View) error
}

func NewSelectWidget(name string, x, y, w, h int, title string, options []model.Option,
	checkHandler func(g *gocui.Gui, v *gocui.View) error) *gocui.View {
	widget := SelectWidget{name: name, x: x, y: y, w: w, h: h, title: title, options: options,
		checkHandler: checkHandler}
	v, _ := widget.Layout()

	return v
}

func NewSelectWidgetWithDefault(name string, x, y, w, h int, title string, options []model.Option, defaultt string,
	checkHandler func(g *gocui.Gui, v *gocui.View) error) *gocui.View {
	widget := SelectWidget{name: name, x: x, y: y, w: w, h: h, title: title, options: options, defaultt: defaultt,
		checkHandler: checkHandler}
	v, _ := widget.Layout()

	return v
}

func (w *SelectWidget) Layout() (*gocui.View, error) {

	if w.h < 1 {
		w.h = 3
	}

	v, _ := vari.Cui.SetView(w.name, w.x, w.y, w.x+w.w, w.y+w.h)
	v.Highlight = true
	ui.SupportScroll(w.name)
	ui.SupportLineHighlight(w.name)

	v.Title = w.title

	labels := make([]string, 0)
	for _, opt := range w.options {
		labels = append(labels, opt.Name)
	}

	fmt.Fprint(v, strings.Join(labels, "\n"))

	_, height := v.Size()
	for true {
		line, _ := ui.GetSelectedLine(v, ".*")
		if w.defaultt != "" {
			if line == w.defaultt {
				break
			}
		} else {
			if zentaoUtils.IsBugFieldDefault(line, w.options) {
				break
			}
		}

		_, oy := v.Origin()
		cx, cy := v.Cursor()

		pos := oy + 1

		if err := v.SetCursor(cx, cy+1); err != nil {
			ox, oy := v.Origin()

			h := len(v.BufferLines()) - height + 1

			if pos < h {
				if err := v.SetOrigin(ox, oy+1); err != nil {
					break
				}
			}
		}

		_, oy1 := v.Origin() // 1
		_, cy1 := v.Cursor() // 4
		if oy1+cy1 >= len(labels)-1 {
			break
		}
	}

	if err := vari.Cui.SetKeybinding(w.name, gocui.KeyEnter, gocui.ModNone, w.checkHandler); err != nil {
		return nil, err
	}
	if err := vari.Cui.SetKeybinding(w.name, gocui.MouseLeft, gocui.ModNone, w.checkHandler); err != nil {
		return nil, err
	}

	return v, nil
}
