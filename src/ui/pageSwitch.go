package ui

import (
	"fmt"
	"github.com/easysoft/zentaoatf/src/action"
	"github.com/easysoft/zentaoatf/src/utils"
	"github.com/jroimartin/gocui"
	"log"
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
	submitInput := NewButtonWidgetAutoWidth(g, "submitInput", buttonX, 4, "Submit", SwitchWorkDir)
	ViewMap["switch"] = append(ViewMap["switch"], submitInput.Name())

	if err := g.SetKeybinding("", gocui.KeyTab, gocui.ModNone, ToggleInput(ViewMap["switch"])); err != nil {
		log.Panicln(err)
	}

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

func DestorySwitchPage(g *gocui.Gui) {
	for _, v := range ViewMap["switch"] {
		g.DeleteView(v)
	}

	g.DeleteKeybinding("", gocui.KeyTab, gocui.ModNone)
	g.DeleteKeybindings("submitInput")
}
