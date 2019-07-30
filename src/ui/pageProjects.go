package ui

import (
	"github.com/easysoft/zentaoatf/src/model"
	"github.com/easysoft/zentaoatf/src/utils"
	"github.com/jroimartin/gocui"
)

var CurrProjectName string
var projectHistoriess []model.WorkHistory

func InitProjectsPage(g *gocui.Gui) error {
	his := utils.Prefer.WorkHistories[0]
	name, _ := getProjectInfo(his)
	CurrProjectName = name

	y := 2
	for _, his := range utils.Prefer.WorkHistories {
		name, label := getProjectInfo(his)

		hisView := NewLabelWidgetAutoWidth(g, name, 0, y, label)
		ViewMap["projects"] = append(ViewMap["projects"], hisView.Name())

		y += 1
	}
	keybindingProjectsButton(g)

	return nil
}

func keybindingProjectsButton(g *gocui.Gui) error {
	for _, his := range utils.Prefer.WorkHistories {
		name, _ := getProjectInfo(his)
		if err := g.SetKeybinding(name, gocui.MouseLeft, gocui.ModNone, toggleProjectsButton); err != nil {
			return err
		}
	}

	return nil
}

func toggleProjectsButton(g *gocui.Gui, v *gocui.View) error {
	CurrProjectName = v.Name()
	SelectProjectsButton(g)

	return nil
}

func SelectProjectsButton(g *gocui.Gui) {
	for _, his := range utils.Prefer.WorkHistories {
		name, _ := getProjectInfo(his)

		v, err := g.View(name)

		if err == nil {
			if v.Name() == CurrProjectName {
				v.Highlight = true
				v.SelBgColor = gocui.ColorWhite
				v.SelFgColor = gocui.ColorBlack
			} else {
				v.Highlight = false
				v.SelBgColor = gocui.ColorBlack
				v.SelFgColor = gocui.ColorDefault
			}
		}
	}
}

func getProjectInfo(his model.WorkHistory) (string, string) {
	var name string
	var label string

	if his.EntityType != "" {
		name = his.EntityType + "-" + his.EntityVal
		label = his.ProjectName
	} else {
		name = his.ProjectPath
		label = utils.PathSomple(his.ProjectPath)
	}

	return name, label
}

func init() {

}

func DestoryProjectsPage(g *gocui.Gui) {
	for _, v := range ViewMap["projects"] {
		g.DeleteView(v)
		g.DeleteKeybindings(v)
	}
}
