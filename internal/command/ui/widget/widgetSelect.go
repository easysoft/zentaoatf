package widget

import (
	"fmt"
	commConsts "github.com/aaronchen2k/deeptest/internal/comm/consts"
	commDomain "github.com/aaronchen2k/deeptest/internal/comm/domain"
	"github.com/aaronchen2k/deeptest/internal/command/ui"
	logUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/log"
	"github.com/awesome-gocui/gocui"
	"strings"
)

type SelectWidget struct {
	name     string
	x, y     int
	w        int
	h        int
	title    string
	options  []commDomain.BugOption
	defaultt string

	checkHandler func(g *gocui.Gui, v *gocui.View) error
}

func NewSelectWidget(name string, x, y, w, h int, title string, options []commDomain.BugOption,
	checkHandler func(g *gocui.Gui, v *gocui.View) error) *gocui.View {
	widget := SelectWidget{name: name, x: x, y: y, w: w, h: h, title: title, options: options,
		checkHandler: checkHandler}
	v, _ := widget.Layout()

	return v
}

func NewSelectWidgetWithDefault(name string, x, y, w, h int, title string, options []commDomain.BugOption, defaultt string,
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

	v, _ := commConsts.Cui.SetView(w.name, w.x, w.y, w.x+w.w, w.y+w.h, 0)
	v.Highlight = true
	ui.SupportScroll(w.name)
	ui.SupportRowHighlight(w.name)

	v.Title = w.title
	logUtils.PrintToCmd(fmt.Sprintf("%s: defalut=%s", v.Name(), w.defaultt), -1)

	labels := make([]string, 0)
	defaultValIndex := -1
	for idx, opt := range w.options {
		labels = append(labels, opt.Name)

		if w.defaultt == opt.Name {
			defaultValIndex = idx
		}
	}

	if defaultValIndex == -1 {
		if len(labels) > 0 {
			w.defaultt = labels[0]
		} else {
			w.defaultt = ""
		}
	}

	fmt.Fprint(v, strings.Join(labels, "\n"))

	for true {
		line := ui.GetSelectedRowVal(v)

		if w.defaultt != "" {
			if line == w.defaultt {
				break
			}
		}

		atBottom := ui.ScrollAction(v, 1)
		if atBottom {
			break
		}
	}

	if err := commConsts.Cui.SetKeybinding(w.name, gocui.KeyEnter, gocui.ModNone, w.checkHandler); err != nil {
		return nil, err
	}
	if err := commConsts.Cui.SetKeybinding(w.name, gocui.MouseLeft, gocui.ModNone, w.checkHandler); err != nil {
		return nil, err
	}

	return v, nil
}
