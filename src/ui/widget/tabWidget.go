package widget

import (
	"fmt"
	"github.com/easysoft/zentaoatf/src/ui"
	"github.com/jroimartin/gocui"
)

var CurrTab string

type TabWidget struct {
	name  string
	x, y  int
	w     int
	label string
}

func NewTabWidget(g *gocui.Gui, name string, x, y int, label string) *gocui.View {
	widget := TabWidget{name: name, x: x, y: y, w: len(label) + 1, label: label}

	v, _ := widget.Layout(g)

	v.Frame = false
	return v
}

func (w *TabWidget) Layout(g *gocui.Gui) (*gocui.View, error) {
	v, err := g.SetView(w.name, w.x, w.y, w.x+w.w, w.y+LabelHeight)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return nil, err
		}

		if err := g.SetKeybinding(w.name, gocui.MouseLeft, gocui.ModNone, ToggleTab); err != nil {
			return nil, err
		}

		Select(g)
		fmt.Fprint(v, w.label)
	}

	return v, nil
}

func ToggleTab(g *gocui.Gui, v *gocui.View) error {
	CurrTab = v.Name()

	Select(g)

	return nil
}

func Select(g *gocui.Gui) {
	for _, name := range ui.Tabs {
		v, err := g.View(name)

		if err == nil {
			if v.Name() == CurrTab {
				v.Highlight = true
				v.SelBgColor = gocui.ColorWhite
				v.SelFgColor = gocui.ColorBlack
			} else {
				v.Highlight = false
				v.SelBgColor = gocui.ColorBlack
				v.SelFgColor = gocui.ColorDefault
			}
		}
	}
}

func init() {
	CurrTab = "testing"
}
