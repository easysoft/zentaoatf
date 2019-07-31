package ui

import (
	"github.com/easysoft/zentaoatf/src/utils"
	"github.com/jroimartin/gocui"
	"log"
)

func InitMainPage(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if maxX < utils.MinWidth {
		maxX = utils.MinWidth
	}
	if maxY < utils.MinHeight {
		maxY = utils.MinHeight
	}
	utils.MainViewHeight = maxY - 10

	quickBarView := NewPanelWidget(g, "quickBarView", 0, 0, utils.LeftWidth, 2, "")
	ViewMap["root"] = append(ViewMap["root"], quickBarView.Name())

	sideView := NewPanelWidget(g, "side", 0, 2, utils.LeftWidth, maxY-3, "")
	ViewMap["root"] = append(ViewMap["root"], sideView.Name())

	x := 2
	for _, name := range Tabs {
		tabView := NewTabWidget(g, name, x, 0, utils.Ucfirst(name))
		ViewMap["root"] = append(ViewMap["root"], tabView.Name())
		x += 10
	}

	mainView := NewPanelWidget(g, "main", utils.LeftWidth, 0, maxX-1-utils.LeftWidth, utils.MainViewHeight, "")
	ViewMap["root"] = append(ViewMap["root"], mainView.Name())
	//mainView.Editable = true
	mainView.Wrap = true

	cmdView := NewPanelWidget(g, "cmd", utils.LeftWidth, utils.MainViewHeight, maxX-1-utils.LeftWidth, 9, "")
	ViewMap["root"] = append(ViewMap["root"], cmdView.Name())
	cmdView.Autoscroll = true

	utils.PrintPreferenceToView(cmdView)

	NewHelpWidget(g)
	MainPageKeyBindings(g)

	InitTestingPage(g)

	return nil
}

func MainPageKeyBindings(g *gocui.Gui) error {
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, Quit); err != nil {
		return err
	}
	if err := g.SetKeybinding("", gocui.KeyCtrlH, gocui.ModNone, ShowHelp); err != nil {
		log.Panicln(err)
	}

	if err := g.SetKeybinding("main", gocui.KeyArrowUp, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			scrollView(v, -1)
			return nil
		}); err != nil {
		return err
	}
	if err := g.SetKeybinding("main", gocui.KeyArrowDown, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			scrollView(v, 1)
			return nil
		}); err != nil {
		return err
	}

	return nil
}

func scrollView(v *gocui.View, dy int) error {
	if v != nil {
		v.Autoscroll = false
		ox, oy := v.Origin()
		if err := v.SetOrigin(ox, oy+dy); err != nil {
			return err
		}
	}
	return nil
}
