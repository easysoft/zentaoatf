package ui

import (
	"github.com/easysoft/zentaoatf/src/utils/config"
	constant "github.com/easysoft/zentaoatf/src/utils/const"
	string2 "github.com/easysoft/zentaoatf/src/utils/string"
	"github.com/easysoft/zentaoatf/src/utils/vari"
	"github.com/jroimartin/gocui"
	"log"
)

func InitMainPage() error {
	maxX, maxY := vari.Cui.Size()
	if maxX < constant.MinWidth {
		maxX = constant.MinWidth
	}
	if maxY < constant.MinHeight {
		maxY = constant.MinHeight
	}
	vari.MainViewHeight = maxY - constant.CmdViewHeight - 1

	quickBarView := NewPanelWidget("quickBarView", 0, 0, constant.LeftWidth, 2, "")
	ViewMap["root"] = append(ViewMap["root"], quickBarView.Name())

	sideView := NewPanelWidget("side", 0, 2, constant.LeftWidth, maxY-3, "")
	ViewMap["root"] = append(ViewMap["root"], sideView.Name())
	sideView.Wrap = true
	sideView.Highlight = true
	sideView.SelBgColor = gocui.ColorWhite
	sideView.SelFgColor = gocui.ColorBlack

	x := 2
	for _, name := range Tabs {
		tabView := NewTabWidget(name, x, 0, string2.Ucfirst(name))
		ViewMap["root"] = append(ViewMap["root"], tabView.Name())
		x += 10
	}

	mainView := NewPanelWidget("main", constant.LeftWidth, 0, maxX-constant.LeftWidth-1, vari.MainViewHeight, "")
	ViewMap["root"] = append(ViewMap["root"], mainView.Name())
	mainView.Wrap = true

	cmdView := NewPanelWidget("cmd", constant.LeftWidth, vari.MainViewHeight, maxX-1-constant.LeftWidth, constant.CmdViewHeight, "")
	ViewMap["root"] = append(ViewMap["root"], cmdView.Name())
	mainView.Wrap = true

	configUtils.PrintPreferenceToView()

	NewHelpWidget()
	MainPageKeyBindings()

	setCurrView("side")
	return nil
}

func MainPageKeyBindings() error {
	if err := vari.Cui.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, Quit); err != nil {
		return err
	}
	if err := vari.Cui.SetKeybinding("", gocui.KeyCtrlH, gocui.ModNone, ShowHelp); err != nil {
		log.Panicln(err)
	}

	setViewScroll("main")
	setViewScroll("cmd")

	return nil
}
