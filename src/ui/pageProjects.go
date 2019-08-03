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

func InitProjectsPage() error {
	his := utils.Prefer.WorkHistories[0]
	id, _, _ := getProjectInfo(his)
	CurrProjectId = id

	y := 2
	for _, his := range utils.Prefer.WorkHistories {
		id, label, _ := getProjectInfo(his)

		hisView := NewLabelWidgetAutoWidth(utils.Cui, id, 0, y, label)
		ViewMap["projects"] = append(ViewMap["projects"], hisView.Name())

		y += 1
	}
	keybindingProjectsButton(utils.Cui)

	return nil
}

func keybindingProjectsButton() error {
	for _, his := range utils.Prefer.WorkHistories {
		id, _, _ := getProjectInfo(his)
		if err := utils.Cui.SetKeybinding(id, gocui.MouseLeft, gocui.ModNone, toggleProjectsButton); err != nil {
			return err
		}
	}

	return nil
}

func toggleProjectsButton(g *gocui.Gui, v *gocui.View) error {
	CurrProjectId = v.Name()
	SelectProjectsButton()

	return nil
}

func SelectProjectsButton() {
	for _, his := range utils.Prefer.WorkHistories {
		id, _, _ := getProjectInfo(his)

		v, err := utils.Cui.View(id)
		if err == nil {
			if id == CurrProjectId {
				v.Highlight = true
				v.SelBgColor = gocui.ColorWhite
				v.SelFgColor = gocui.ColorBlack

				printForSwitch(utils.Cui, his) // 显示项目信息
				showWitchButton()
			} else {
				v.Highlight = false
				v.SelBgColor = gocui.ColorBlack
				v.SelFgColor = gocui.ColorDefault
			}
		}
	}
}

func showWitchButton() error {
	maxX, _ := utils.Cui.Size()

	switchButton := NewButtonWidgetAutoWidth("switchButton", maxX-15, 1, "Switch To", switchProject)
	ViewMap["projects"] = append(ViewMap["projects"], switchButton.Name())

	return nil
}
func switchProject(g *gocui.Gui, v *gocui.View) error {
	for _, his := range utils.Prefer.WorkHistories {
		id, label, path := getProjectInfo(his)
		if id == CurrProjectId {
			action.SwitchWorkDir(path)
			utils.PrintToCmd(fmt.Sprintf("success to switch to project %s: %s at %s",
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

func printForSwitch(his model.WorkHistory) {
	config := utils.ReadProjectConfig(his.ProjectPath)
	name := config.ProjectName
	if name == "" {
		name = "No Name"
	}

	str := "%s\n Work dir: %s\n Zentao project: %s\n Import type: %s\n Product code: %s\n Language: %s\n " +
		"Independent ExpectResult file: %t"
	str = fmt.Sprintf(str, name, his.ProjectPath, config.Url, config.EntityType, config.EntityVal,
		config.LangType, !config.SingleFile)

	utils.PrintToMainNoScroll(utils.Cui, str)
}

func init() {

}

func DestoryProjectsPage() {
	for _, v := range ViewMap["projects"] {
		utils.Cui.DeleteView(v)
		utils.Cui.DeleteKeybindings(v)
	}
}
