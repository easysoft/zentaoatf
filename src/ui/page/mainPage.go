package page

import (
	"github.com/easysoft/zentaoatf/src/ui"
	"github.com/easysoft/zentaoatf/src/ui/widget"
	"github.com/easysoft/zentaoatf/src/utils"
	"github.com/jroimartin/gocui"
	"log"
)

func InitMainPage(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if maxX < 130 {
		maxX = 130
	}
	if maxY < 36 {
		maxY = 36
	}

	quickBarView := widget.NewPanelWidget(g, "quickBarView", 0, 0, ui.LeftWidth, 2, "")
	ui.ViewMap["root"] = append(ui.ViewMap["root"], quickBarView.Name())

	sideView := widget.NewPanelWidget(g, "side", 0, 2, ui.LeftWidth, maxY-3, "")
	ui.ViewMap["root"] = append(ui.ViewMap["root"], sideView.Name())

	x := 2
	for _, name := range ui.Tabs {
		tabView := widget.NewTabWidget(g, name, x, 0, utils.Ucfirst(name))
		ui.ViewMap["root"] = append(ui.ViewMap["root"], tabView.Name())
		x += 10
	}

	mainView := widget.NewPanelWidget(g, "main", ui.LeftWidth, 0, maxX-1-ui.LeftWidth, maxY-10, "")
	ui.ViewMap["root"] = append(ui.ViewMap["root"], mainView.Name())

	cmdView := widget.NewPanelWidget(g, "cmd", ui.LeftWidth, maxY-10, maxX-1-ui.LeftWidth, 9, "")
	ui.ViewMap["root"] = append(ui.ViewMap["root"], cmdView.Name())
	cmdView.Autoscroll = true

	utils.PrintPreferenceToView(cmdView)

	widget.NewHelpWidget(g)

	if err := MainPageKeyBindings(g); err != nil {
		log.Panicln(err)
	}

	return nil
}

func MainPageKeyBindings(g *gocui.Gui) error {
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, ui.Quit); err != nil {
		return err
	}
	if err := g.SetKeybinding("", gocui.KeyCtrlH, gocui.ModNone, widget.ShowHelp); err != nil {
		log.Panicln(err)
	}

	//if err := g.SetKeybinding("import", gocui.MouseLeft, gocui.ModNone, InitImportPage); err != nil {
	//	return err
	//}
	//if err := g.SetKeybinding("switch", gocui.MouseLeft, gocui.ModNone, InitSwitchPage); err != nil {
	//	return err
	//}

	return nil
}
