package ui

import (
	"github.com/easysoft/zentaoatf/src/utils"
	"github.com/jroimartin/gocui"
)

func DestoryLeftPages(g *gocui.Gui) {
	utils.ClearSide(g)

	DestoryTestingPage(g)
	DestoryProjectsPage(g)
	DestorySettingsPage(g)

	ViewMap["testing"] = make([]string, 0)
	ViewMap["projects"] = make([]string, 0)
	ViewMap["settings"] = make([]string, 0)
}

func DestoryRightPages(g *gocui.Gui) {
	mainView, err := g.View("main")
	if err == nil {
		mainView.Clear()
	}

	DestoryImportPage(g)
	DestorySwitchPage(g)

	ViewMap["import"] = make([]string, 0)
	ViewMap["switch"] = make([]string, 0)
}
