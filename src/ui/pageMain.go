package ui

import (
	"github.com/easysoft/zentaoatf/src/utils"
	"github.com/jroimartin/gocui"
	"log"
)

func InitMainPage() error {
	maxX, maxY := utils.Cui.Size()
	//if maxX < utils.MinWidth {
	//	maxX = utils.MinWidth
	//}
	//if maxY < utils.MinHeight {
	//	maxY = utils.MinHeight
	//}
	utils.MainViewHeight = maxY - utils.CmdViewHeight - 1

	quickBarView := NewPanelWidget("quickBarView", 0, 0, utils.LeftWidth, 2, "")
	ViewMap["root"] = append(ViewMap["root"], quickBarView.Name())

	sideView := NewPanelWidget("side", 0, 2, utils.LeftWidth, maxY-3, "")
	ViewMap["root"] = append(ViewMap["root"], sideView.Name())
	sideView.Wrap = true
	sideView.Highlight = true
	sideView.SelBgColor = gocui.ColorWhite
	sideView.SelFgColor = gocui.ColorBlack

	x := 2
	for _, name := range Tabs {
		tabView := NewTabWidget(name, x, 0, utils.Ucfirst(name))
		ViewMap["root"] = append(ViewMap["root"], tabView.Name())
		x += 10
	}

	mainView := NewPanelWidget("main", utils.LeftWidth, 0, maxX-utils.LeftWidth-1, utils.MainViewHeight, "")
	ViewMap["root"] = append(ViewMap["root"], mainView.Name())
	mainView.Wrap = true

	cmdView := NewPanelWidget("cmd", utils.LeftWidth, utils.MainViewHeight, maxX-1-utils.LeftWidth, utils.CmdViewHeight, "")
	ViewMap["root"] = append(ViewMap["root"], cmdView.Name())
	mainView.Wrap = true

	utils.PrintPreferenceToView()

	NewHelpWidget()
	MainPageKeyBindings()

	setCurrView("side")
	return nil
}

func MainPageKeyBindings() error {
	if err := utils.Cui.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, Quit); err != nil {
		return err
	}
	if err := utils.Cui.SetKeybinding("", gocui.KeyCtrlH, gocui.ModNone, ShowHelp); err != nil {
		log.Panicln(err)
	}

	setViewScroll("main")
	setViewScroll("cmd")

	return nil
}
