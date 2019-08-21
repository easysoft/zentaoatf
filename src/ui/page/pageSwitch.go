package page

import (
	"fmt"
	"github.com/easysoft/zentaoatf/src/action"
	"github.com/easysoft/zentaoatf/src/ui"
	"github.com/easysoft/zentaoatf/src/ui/widget"
	constant "github.com/easysoft/zentaoatf/src/utils/const"
	i118Utils "github.com/easysoft/zentaoatf/src/utils/i118"
	"github.com/easysoft/zentaoatf/src/utils/log"
	"github.com/easysoft/zentaoatf/src/utils/vari"
	"github.com/jroimartin/gocui"
	"strings"
)

func InitSwitchPage() error {
	DestoryRightPages()

	maxX, _ := vari.Cui.Size()
	slideView, _ := vari.Cui.View("side")
	slideX, _ := slideView.Size()

	left := slideX + 2
	right := left + widget.LabelWidth
	workDirLabel := widget.NewLabelWidget("workDirLabel", left, 1,
		i118Utils.I118Prt.Sprintf("Workdir"))
	ui.ViewMap["switch"] = append(ui.ViewMap["switch"], workDirLabel.Name())

	left = right + ui.Space
	right = left + widget.TextWidthFull
	workDirInput := widget.NewTextWidget("workDirInput", left, 1, widget.TextWidthFull, vari.Prefer.WorkDir)
	ui.ViewMap["switch"] = append(ui.ViewMap["switch"], workDirInput.Name())
	if _, err := vari.Cui.SetCurrentView("workDirInput"); err != nil {
		return err
	}

	buttonX := (maxX-constant.LeftWidth)/2 + constant.LeftWidth - widget.ButtonWidth
	submitInput := widget.NewButtonWidgetAutoWidth("submitInput", buttonX, 4,
		i118Utils.I118Prt.Sprintf("switch"), SwitchWorkDir)
	ui.ViewMap["switch"] = append(ui.ViewMap["switch"], submitInput.Name())

	ui.AddEventForInputWidgets(ui.ViewMap["switch"])

	return nil
}

func SwitchWorkDir(g *gocui.Gui, v *gocui.View) error {
	workDirView, _ := g.View("workDirInput")

	workDir := strings.TrimSpace(workDirView.Buffer())

	logUtils.PrintToCmd(fmt.Sprintf("#atf switch -d %s", workDir))

	err := action.SwitchWorkDir(workDir)
	if err == nil {
		workDirView.Clear()
		workDirView.Write([]byte(vari.Prefer.WorkDir))
	} else {
		logUtils.PrintToCmd(err.Error())
	}

	return nil
}

func DestorySwitchPage() {
	for _, v := range ui.ViewMap["switch"] {
		vari.Cui.DeleteView(v)
		vari.Cui.DeleteKeybindings(v)
	}

	vari.Cui.DeleteKeybinding("", gocui.KeyTab, gocui.ModNone)
}
