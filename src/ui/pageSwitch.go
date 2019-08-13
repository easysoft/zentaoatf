package ui

import (
	"fmt"
	"github.com/easysoft/zentaoatf/src/action"
	constant "github.com/easysoft/zentaoatf/src/utils/const"
	"github.com/easysoft/zentaoatf/src/utils/date"
	print2 "github.com/easysoft/zentaoatf/src/utils/print"
	"github.com/easysoft/zentaoatf/src/utils/vari"
	"github.com/jroimartin/gocui"
	"strings"
	"time"
)

func InitSwitchPage() error {
	DestoryRightPages()

	maxX, _ := vari.Cui.Size()
	slideView, _ := vari.Cui.View("side")
	slideX, _ := slideView.Size()

	left := slideX + 2
	right := left + LabelWidth
	workDirLabel := NewLabelWidget("workDirLabel", left, 1, "WorkDir")
	ViewMap["switch"] = append(ViewMap["switch"], workDirLabel.Name())

	left = right + Space
	right = left + TextWidthFull
	workDirInput := NewTextWidget("workDirInput", left, 1, TextWidthFull, vari.Prefer.WorkDir)
	ViewMap["switch"] = append(ViewMap["switch"], workDirInput.Name())
	if _, err := vari.Cui.SetCurrentView("workDirInput"); err != nil {
		return err
	}

	buttonX := (maxX-constant.LeftWidth)/2 + constant.LeftWidth - ButtonWidth
	submitInput := NewButtonWidgetAutoWidth("submitInput", buttonX, 4, "Switch", SwitchWorkDir)
	ViewMap["switch"] = append(ViewMap["switch"], submitInput.Name())

	keyBindsInput(ViewMap["switch"])

	return nil
}

func SwitchWorkDir(g *gocui.Gui, v *gocui.View) error {
	workDirView, _ := g.View("workDirInput")

	workDir := strings.TrimSpace(workDirView.Buffer())

	print2.PrintToCmd(fmt.Sprintf("#atf switch -d %s", workDir))

	err := action.SwitchWorkDir(workDir)
	if err == nil {
		workDirView.Clear()
		workDirView.Write([]byte(vari.Prefer.WorkDir))

		print2.PrintToCmd(fmt.Sprintf("success to switch project to %s at %s",
			workDir, dateUtils.DateTimeStr(time.Now())))
	} else {
		print2.PrintToCmd(err.Error())
	}

	return nil
}

func DestorySwitchPage() {
	for _, v := range ViewMap["switch"] {
		vari.Cui.DeleteView(v)
		vari.Cui.DeleteKeybindings(v)
	}

	vari.Cui.DeleteKeybinding("", gocui.KeyTab, gocui.ModNone)
}
