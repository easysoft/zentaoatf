package page

import (
	"fmt"
	"github.com/easysoft/zentaoatf/src/action"
	"github.com/easysoft/zentaoatf/src/model"
	"github.com/easysoft/zentaoatf/src/ui"
	"github.com/easysoft/zentaoatf/src/ui/widget"
	config2 "github.com/easysoft/zentaoatf/src/utils/config"
	"github.com/easysoft/zentaoatf/src/utils/date"
	print2 "github.com/easysoft/zentaoatf/src/utils/print"
	string2 "github.com/easysoft/zentaoatf/src/utils/string"
	"github.com/easysoft/zentaoatf/src/utils/vari"
	"github.com/jroimartin/gocui"
	"time"
)

var CurrProjectId string
var projectHistories []model.WorkHistory

func InitProjectsPage() error {
	his := vari.Prefer.WorkHistories[0]
	id, _, _ := getProjectInfo(his)
	CurrProjectId = id

	y := 2
	for _, his := range vari.Prefer.WorkHistories {
		id, label, _ := getProjectInfo(his)

		hisView := widget.NewLabelWidgetAutoWidth(id, 0, y, label)
		ui.ViewMap["projects"] = append(ui.ViewMap["projects"], hisView.Name())

		y += 1
	}
	keybindingProjectsButton()

	return nil
}

func keybindingProjectsButton() error {
	for _, his := range vari.Prefer.WorkHistories {
		id, _, _ := getProjectInfo(his)
		if err := vari.Cui.SetKeybinding(id, gocui.MouseLeft, gocui.ModNone, toggleProjectsButton); err != nil {
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
	for _, his := range vari.Prefer.WorkHistories {
		id, _, _ := getProjectInfo(his)

		v, err := vari.Cui.View(id)
		if err == nil {
			if id == CurrProjectId {
				v.Highlight = true
				v.SelBgColor = gocui.ColorWhite
				v.SelFgColor = gocui.ColorBlack

				printForSwitch(his) // 显示项目信息
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
	maxX, _ := vari.Cui.Size()

	switchButton := widget.NewButtonWidgetAutoWidth("switchButton", maxX-15, 1, "Switch To", switchProject)
	ui.ViewMap["projects"] = append(ui.ViewMap["projects"], switchButton.Name())

	return nil
}
func switchProject(g *gocui.Gui, v *gocui.View) error {
	for _, his := range vari.Prefer.WorkHistories {
		id, label, path := getProjectInfo(his)
		if id == CurrProjectId {
			action.SwitchWorkDir(path)
			print2.PrintToCmd(fmt.Sprintf("success to switch to project %s: %s at %s",
				label, path, dateUtils.DateTimeStr(time.Now())))
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
		label = string2.PathSimple(his.ProjectPath)
	}

	return id, label, path
}

func printForSwitch(his model.WorkHistory) {
	config := config2.ReadProjectConfig(his.ProjectPath)
	name := config.ProjectName
	if name == "" {
		name = "No Name"
	}

	str := "%s\n Work dir: %s\n Zentao project: %s\n Import type: %s\n Product code: %s\n Language: %s\n " +
		"Independent ExpectResult file: %t"
	str = fmt.Sprintf(str, name, his.ProjectPath, config.Url, config.EntityType, config.EntityVal,
		config.LangType, !config.SingleFile)

	print2.PrintToMainNoScroll(str)
}

func init() {

}

func DestoryProjectsPage() {
	for _, v := range ui.ViewMap["projects"] {
		vari.Cui.DeleteView(v)
		vari.Cui.DeleteKeybindings(v)
	}
}
