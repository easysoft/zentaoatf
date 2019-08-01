package ui

import (
	"fmt"
	"github.com/easysoft/zentaoatf/src/action"
	"github.com/easysoft/zentaoatf/src/model"
	"github.com/easysoft/zentaoatf/src/utils"
	"github.com/jroimartin/gocui"
	"time"
)

var CurrProjectId string
var projectHistories []model.WorkHistory

func InitProjectsPage(g *gocui.Gui) error {
	his := utils.Prefer.WorkHistories[0]
	id, _, _ := getProjectInfo(his)
	CurrProjectId = id

	y := 2
	for _, his := range utils.Prefer.WorkHistories {
		id, label, _ := getProjectInfo(his)

		hisView := NewLabelWidgetAutoWidth(g, id, 0, y, label)
		ViewMap["projects"] = append(ViewMap["projects"], hisView.Name())

		y += 1
	}
	keybindingProjectsButton(g)

	return nil
}

func keybindingProjectsButton(g *gocui.Gui) error {
	for _, his := range utils.Prefer.WorkHistories {
		name, _, _ := getProjectInfo(his)
		if err := g.SetKeybinding(name, gocui.MouseLeft, gocui.ModNone, toggleProjectsButton); err != nil {
			return err
		}
	}

	return nil
}

func toggleProjectsButton(g *gocui.Gui, v *gocui.View) error {
	CurrProjectId = v.Name()
	SelectProjectsButton(g)

	return nil
}

func SelectProjectsButton(g *gocui.Gui) {
	for _, his := range utils.Prefer.WorkHistories {
		name, _, _ := getProjectInfo(his)

		v, err := g.View(name)
		if err == nil {
			if v.Name() == CurrProjectId {
				v.Highlight = true
				v.SelBgColor = gocui.ColorWhite
				v.SelFgColor = gocui.ColorBlack

				printForSwitch(g, his) // 显示项目信息
				showWitchButton(g)
			} else {
				v.Highlight = false
				v.SelBgColor = gocui.ColorBlack
				v.SelFgColor = gocui.ColorDefault
			}
		}
	}
}

func showWitchButton(g *gocui.Gui) error {
	maxX, _ := g.Size()

	switchButton := NewButtonWidgetAutoWidth(g, "switchButton", maxX-15, 1, "Switch To", switchProject)
	ViewMap["projects"] = append(ViewMap["projects"], switchButton.Name())

	return nil
}
func switchProject(g *gocui.Gui, v *gocui.View) error {
	for _, his := range utils.Prefer.WorkHistories {
		id, label, path := getProjectInfo(his)
		if id == CurrProjectId {
			action.SwitchWorkDir(path)
			utils.PrintToCmd(g, fmt.Sprintf("success to switch to project %s: %s at %s",
				label, path, utils.DateTimeStr(time.Now())))
			break
		}
	}
	return nil
}

func getProjectInfo(his model.WorkHistory) (string, string, string) {
	var id string
	var label string
	var path string

	id = his.Id
	path = his.ProjectPath
	if his.EntityType != "" {
		label = his.ProjectName
	} else {
		label = utils.PathSimple(his.ProjectPath)
	}

	return id, label, path
}

func printForSwitch(g *gocui.Gui, his model.WorkHistory) {
	config := utils.ReadConfig()
	name := config.ProjectName
	if name == "" {
		name = "No Name"
	}

	str := "%s\n Work dir: %s\n Zentao project: %s\n Import type: %s\n Product code: %s\n Language: %s\n " +
		"Independent ExpectResult file: %t"
	str = fmt.Sprintf(str, name, his.ProjectPath, config.Url, config.EntityType, config.EntityVal,
		config.LangType, !config.SingleFile)

	utils.PrintToMainNoScroll(g, str)
}

func init() {

}

func DestoryProjectsPage(g *gocui.Gui) {
	for _, v := range ViewMap["projects"] {
		g.DeleteView(v)
		g.DeleteKeybindings(v)
	}
}
