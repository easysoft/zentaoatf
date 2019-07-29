package page

import (
	"fmt"
	"github.com/easysoft/zentaoatf/src/action"
	"github.com/easysoft/zentaoatf/src/ui"
	"github.com/easysoft/zentaoatf/src/ui/widget"
	"github.com/easysoft/zentaoatf/src/utils"
	"github.com/jroimartin/gocui"
	"log"
	"strings"
	"time"
)

func InitSwitchPage(g *gocui.Gui, v *gocui.View) error {
	DestoryPages(g, v)

	maxX, _ := g.Size()
	slideView, _ := g.View("side")
	slideX, _ := slideView.Size()

	left := slideX + 2
	right := left + widget.LabelWidth
	workDirLabel := widget.NewLabelWidget(g, "workDirLabel", left, 1, "WorkDir")
	ui.ViewMap["switch"] = append(ui.ViewMap["switch"], workDirLabel.Name())

	left = right + ui.Space
	right = left + widget.TextWidthFull
	workDirInput := widget.NewTextWidget(g, "workDirInput", left, 1, widget.TextWidthFull, utils.Prefer.WorkDir)
	ui.ViewMap["switch"] = append(ui.ViewMap["switch"], workDirInput.Name())
	if _, err := g.SetCurrentView("workDirInput"); err != nil {
		return err
	}

	buttonX := (maxX-ui.LeftWidth)/2 + ui.LeftWidth - widget.ButtonWidth
	submitInput := widget.NewButtonWidgetAutoWidth(g, "submitInput", buttonX, 4, "Submit", SwitchWorkDir)
	ui.ViewMap["switch"] = append(ui.ViewMap["switch"], submitInput.Name())

	if err := g.SetKeybinding("", gocui.KeyTab, gocui.ModNone, ui.ToggleInput(ui.ViewMap["switch"])); err != nil {
		log.Panicln(err)
	}

	ui.HideHelp(g)

	return nil
}

func SwitchWorkDir(g *gocui.Gui, v *gocui.View) error {
	workDirView, _ := g.View("workDirInput")

	workDir := strings.TrimSpace(workDirView.ViewBuffer())

	cmdView, _ := g.View("cmd")
	_, _ = fmt.Fprintln(cmdView, fmt.Sprintf("#atf switch -d %s", workDir))

	err := action.SwitchWorkDir(workDir)
	if err == nil {
		workDirView.Clear()
		workDirView.Write([]byte(utils.Prefer.WorkDir))

		fmt.Fprintln(cmdView, fmt.Sprintf("success to switch project to %s at %s",
			workDir, utils.DateTimeStr(time.Now())))
	} else {
		fmt.Fprintln(cmdView, err.Error())
	}

	return nil
}

func DestorySwitchPage(g *gocui.Gui, v *gocui.View) {
	for _, v := range ui.ViewMap["switch"] {
		g.DeleteView(v)
	}

	g.DeleteKeybinding("", gocui.KeyTab, gocui.ModNone)
	g.DeleteKeybindings("submitInput")

	ui.ViewMap["import"] = make([]string, 0)
}
