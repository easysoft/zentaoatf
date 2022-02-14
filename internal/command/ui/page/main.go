package page

import (
	commConsts "github.com/aaronchen2k/deeptest/internal/comm/consts"
	"github.com/aaronchen2k/deeptest/internal/command/ui"
	"github.com/aaronchen2k/deeptest/internal/command/ui/widget"
	"github.com/aaronchen2k/deeptest/internal/pkg/consts"
	"github.com/awesome-gocui/gocui"
	"log"
)

func InitMainPage() error {
	maxX, maxY := commConsts.Cui.Size()
	if maxX < consts.MinWidth {
		maxX = consts.MinWidth
	}
	if maxY < consts.MinHeight {
		maxY = consts.MinHeight
	}
	commConsts.MainViewHeight = maxY - consts.CmdViewHeight

	mainView := widget.NewPanelWidget("main", 0, 0, maxX-2, commConsts.MainViewHeight, "")

	ui.ViewMap["root"] = append(ui.ViewMap["root"], mainView.Name())

	cmdView := widget.NewPanelWidget("cmd", 0, commConsts.MainViewHeight, maxX-2, consts.CmdViewHeight-1, "")
	ui.ViewMap["root"] = append(ui.ViewMap["root"], cmdView.Name())

	widget.NewHelpWidget()
	MainPageKeyBindings()

	return nil
}

func MainPageKeyBindings() error {
	if err := commConsts.Cui.SetKeybinding("", gocui.KeyCtrlH, gocui.ModNone, widget.ShowHelpFromView); err != nil {
		log.Panicln(err)
	}
	if err := commConsts.Cui.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, ui.Quit); err != nil {
		return err
	}

	ui.SupportScroll("cmd")

	v, _ := commConsts.Cui.View("cmd")
	v.Autoscroll = true

	return nil
}
