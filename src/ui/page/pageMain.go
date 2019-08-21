package page

import (
	"github.com/easysoft/zentaoatf/src/ui"
	"github.com/easysoft/zentaoatf/src/ui/widget"
	"github.com/easysoft/zentaoatf/src/utils/config"
	constant "github.com/easysoft/zentaoatf/src/utils/const"
	i118Utils "github.com/easysoft/zentaoatf/src/utils/i118"
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
	vari.MainViewHeight = maxY - constant.CmdViewHeight

	quickBarView := widget.NewPanelWidget("quickBarView", 0, 0, constant.LeftWidth, 2, "")
	ui.ViewMap["root"] = append(ui.ViewMap["root"], quickBarView.Name())

	sideView := widget.NewPanelWidget("side", 0, 2, constant.LeftWidth, maxY-3, "")
	ui.ViewMap["root"] = append(ui.ViewMap["root"], sideView.Name())
	sideView.Wrap = true
	sideView.Highlight = true
	sideView.SelBgColor = gocui.ColorWhite
	sideView.SelFgColor = gocui.ColorBlack

	x := 2
	for _, name := range ui.ModuleTabs {
		tabView := NewTabWidget(name, x, 0, i118Utils.I118Prt.Sprintf(name))
		ui.ViewMap["root"] = append(ui.ViewMap["root"], tabView.Name())
		x += 10
	}

	mainView := widget.NewPanelWidget("main", constant.LeftWidth, 0, maxX-constant.LeftWidth-1, vari.MainViewHeight, "")
	ui.ViewMap["root"] = append(ui.ViewMap["root"], mainView.Name())

	cmdView := widget.NewPanelWidget("cmd", constant.LeftWidth, vari.MainViewHeight, maxX-1-constant.LeftWidth, constant.CmdViewHeight-1, "")
	ui.ViewMap["root"] = append(ui.ViewMap["root"], cmdView.Name())

	configUtils.PrintPreferenceToView()

	widget.NewHelpWidget()
	MainPageKeyBindings()

	ui.SetCurrView("side")
	return nil
}

func MainPageKeyBindings() error {
	if err := vari.Cui.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, ui.Quit); err != nil {
		return err
	}
	if err := vari.Cui.SetKeybinding("", gocui.KeyCtrlH, gocui.ModNone, widget.ShowHelpFromView); err != nil {
		log.Panicln(err)
	}

	ui.SupportScroll("main")
	ui.SupportScroll("cmd")
	ui.SupportScroll("side")

	v, _ := vari.Cui.View("cmd")
	v.Autoscroll = true

	return nil
}
