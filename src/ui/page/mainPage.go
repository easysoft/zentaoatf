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

	quickBarView := widget.NewPanelWidget(g, "quickBarView", 0, 0, ui.LeftWidth, 2, "")
	ui.ViewMap["root"] = append(ui.ViewMap["root"], quickBarView.Name())

	importView := widget.NewLabelWidget(g, "import", 3, 0, "Import")
	ui.ViewMap["root"] = append(ui.ViewMap["root"], importView.Name())

	switchView := widget.NewLabelWidget(g, "switch", 19, 0, "Switch")
	ui.ViewMap["root"] = append(ui.ViewMap["root"], switchView.Name())

	sideView := widget.NewPanelWidget(g, "side", 0, 2, ui.LeftWidth, maxY-3, "")
	ui.ViewMap["root"] = append(ui.ViewMap["root"], sideView.Name())

	mainView := widget.NewPanelWidget(g, "main", ui.LeftWidth, 0, maxX-1-ui.LeftWidth, maxY-10, "")
	ui.ViewMap["root"] = append(ui.ViewMap["root"], mainView.Name())

	cmdView := widget.NewPanelWidget(g, "cmd", ui.LeftWidth, maxY-10, maxX-1-ui.LeftWidth, 9, "")
	ui.ViewMap["root"] = append(ui.ViewMap["root"], cmdView.Name())

	utils.PrintConfigToView(cmdView)

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
	if err := g.SetKeybinding("", gocui.KeyCtrlH, gocui.ModNone, ui.ShowHelp); err != nil {
		log.Panicln(err)
	}

	if err := g.SetKeybinding("import", gocui.MouseLeft, gocui.ModNone, InitImportPage); err != nil {
		return err
	}
	if err := g.SetKeybinding("switch", gocui.MouseLeft, gocui.ModNone, SwitchProjectUi); err != nil {
		return err
	}

	return nil
}

func SwitchProjectUi(g *gocui.Gui, v *gocui.View) error {
	DestoryImportPage(g, v)
	return nil
}
