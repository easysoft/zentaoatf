package ui

import (
	"fmt"
	"github.com/easysoft/zentaoatf/src/action"
	"github.com/easysoft/zentaoatf/src/utils"
	"github.com/jroimartin/gocui"
	"strings"
	"time"
)

func InitSwitchPage() error {
	DestoryRightPages()

	maxX, _ := utils.Cui.Size()
	slideView, _ := utils.Cui.View("side")
	slideX, _ := slideView.Size()

	left := slideX + 2
	right := left + LabelWidth
	workDirLabel := NewLabelWidget("workDirLabel", left, 1, "WorkDir")
	ViewMap["switch"] = append(ViewMap["switch"], workDirLabel.Name())

	left = right + Space
	right = left + TextWidthFull
	workDirInput := NewTextWidget("workDirInput", left, 1, TextWidthFull, utils.Prefer.WorkDir)
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

	utils.PrintToCmd(fmt.Sprintf("#atf switch -d %s", workDir))

	err := action.SwitchWorkDir(workDir)
	if err == nil {
		workDirView.Clear()
		workDirView.Write([]byte(utils.Prefer.WorkDir))

		utils.PrintToCmd(fmt.Sprintf("success to switch project to %s at %s",
			workDir, utils.DateTimeStr(time.Now())))
	} else {
		utils.PrintToCmd(err.Error())
	}

	return nil
}

func DestorySwitchPage() {
	for _, v := range ViewMap["switch"] {
		utils.Cui.DeleteView(v)
		utils.Cui.DeleteKeybindings(v)
	}

	utils.Cui.DeleteKeybinding("", gocui.KeyTab, gocui.ModNone)
}
