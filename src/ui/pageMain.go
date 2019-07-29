package ui

import (
	"github.com/easysoft/zentaoatf/src/utils"
	"github.com/jroimartin/gocui"
	"log"
)

func InitMainPage(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if maxX < MinWidth {
		maxX = MinWidth
	}
	if maxY < MinHeight {
		maxY = MinHeight
	}

	quickBarView := NewPanelWidget(g, "quickBarView", 0, 0, LeftWidth, 2, "")
	ViewMap["root"] = append(ViewMap["root"], quickBarView.Name())

	sideView := NewPanelWidget(g, "side", 0, 2, LeftWidth, maxY-3, "")
	ViewMap["root"] = append(ViewMap["root"], sideView.Name())

	x := 2
	for _, name := range Tabs {
		tabView := NewTabWidget(g, name, x, 0, utils.Ucfirst(name))
		ViewMap["root"] = append(ViewMap["root"], tabView.Name())
		x += 10
	}

	mainView := NewPanelWidget(g, "main", LeftWidth, 0, maxX-1-LeftWidth, maxY-10, "")
	ViewMap["root"] = append(ViewMap["root"], mainView.Name())

	cmdView := NewPanelWidget(g, "cmd", LeftWidth, maxY-10, maxX-1-LeftWidth, 9, "")
	ViewMap["root"] = append(ViewMap["root"], cmdView.Name())
	cmdView.Autoscroll = true

	utils.PrintPreferenceToView(cmdView)

	NewHelpWidget(g)

	if err := MainPageKeyBindings(g); err != nil {
		log.Panicln(err)
	}

	return nil
}

func MainPageKeyBindings(g *gocui.Gui) error {
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, Quit); err != nil {
		return err
	}
	if err := g.SetKeybinding("", gocui.KeyCtrlH, gocui.ModNone, ShowHelp); err != nil {
		log.Panicln(err)
	}

	return nil
}
