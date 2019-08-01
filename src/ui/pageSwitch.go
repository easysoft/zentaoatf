package ui

import (
	"fmt"
	"github.com/easysoft/zentaoatf/src/action"
	"github.com/easysoft/zentaoatf/src/utils"
	"github.com/jroimartin/gocui"
	"strings"
	"time"
)

func InitSwitchPage(g *gocui.Gui) error {
	DestoryRightPages(g)

	maxX, _ := g.Size()
	slideView, _ := g.View("side")
	slideX, _ := slideView.Size()

	left := slideX + 2
	right := left + LabelWidth
	workDirLabel := NewLabelWidget(g, "workDirLabel", left, 1, "WorkDir")
	ViewMap["switch"] = append(ViewMap["switch"], workDirLabel.Name())

	left = right + Space
	right = left + TextWidthFull
	workDirInput := NewTextWidget(g, "workDirInput", left, 1, TextWidthFull, utils.Prefer.WorkDir)
	ViewMap["switch"] = append(ViewMap["switch"], workDirInput.Name())
	if _, err := g.SetCurrentView("workDirInput"); err != nil {
		return err
	}

	buttonX := (maxX-utils.LeftWidth)/2 + utils.LeftWidth - ButtonWidth
	submitInput := NewButtonWidgetAutoWidth(g, "submitInput", buttonX, 4, "Switch", SwitchWorkDir)
	ViewMap["switch"] = append(ViewMap["switch"], submitInput.Name())

	keyBindsInput(ViewMap["switch"])

	return nil
}

func SwitchWorkDir(g *gocui.Gui, v *gocui.View) error {
	workDirView, _ := g.View("workDirInput")

	workDir := strings.TrimSpace(workDirView.Buffer())

	utils.PrintToCmd(g, fmt.Sprintf("#atf switch -d %s", workDir))

	err := action.SwitchWorkDir(workDir)
	if err == nil {
		workDirView.Clear()
		workDirView.Write([]byte(utils.Prefer.WorkDir))

		utils.PrintToCmd(g, fmt.Sprintf("success to switch project to %s at %s",
			workDir, utils.DateTimeStr(time.Now())))
	} else {
		utils.PrintToCmd(g, err.Error())
	}

	return nil
}

func keyBindsSwitch(g *gocui.Gui) {
	for _, v := range ViewMap["switch"] {
		if strings.Index(v, "Input") > -1 {
			setInputEvent(g, v)
		}
	}
}

func DestorySwitchPage(g *gocui.Gui) {
	for _, v := range ViewMap["switch"] {
		g.DeleteView(v)
		g.DeleteKeybindings(v)
	}

	g.DeleteKeybinding("", gocui.KeyTab, gocui.ModNone)
}
