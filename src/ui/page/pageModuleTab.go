package page

import (
	"fmt"
	"github.com/easysoft/zentaoatf/src/ui"
	"github.com/easysoft/zentaoatf/src/ui/widget"
	"github.com/easysoft/zentaoatf/src/utils/vari"
	"github.com/jroimartin/gocui"
)

var CurrTab string

type TabWidget struct {
	name  string
	x, y  int
	w     int
	label string
}

func NewTabWidget(name string, x, y int, label string) *gocui.View {
	widget := TabWidget{name: name, x: x, y: y, w: len(label) + 1, label: label}

	v, _ := widget.Layout()

	v.Frame = false
	return v
}

func (w *TabWidget) Layout() (*gocui.View, error) {
	v, err := vari.Cui.SetView(w.name, w.x, w.y, w.x+w.w, w.y+widget.LabelHeight)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return nil, err
		}

		if err := vari.Cui.SetKeybinding(w.name, gocui.MouseLeft, gocui.ModNone, ToggleTab); err != nil {
			return nil, err
		}

		ShowTab()
		fmt.Fprint(v, w.label)
	}

	return v, nil
}

func ToggleTab(g *gocui.Gui, v *gocui.View) error {
	CurrTab = v.Name()

	ShowTab()

	return nil
}

func ShowTab() {
	for _, name := range ui.ModuleTabs {
		v, err := vari.Cui.View(name)

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

	DestoryLeftPages()
	DestoryRightPages()
	widget.HideHelp()

	if CurrTab == "testing" {
		InitTestPage()
	} else if CurrTab == "projects" {
		InitProjectsPage()
	} else if CurrTab == "settings" {
		InitSettingsPage()
	}
}

func init() {
	CurrTab = "testing"
}
