package ui

import (
	"github.com/jroimartin/gocui"
)

func DestoryRightPages(g *gocui.Gui) {
	DestoryImportPage(g)
	DestorySwitchPage(g)

	ViewMap["import"] = make([]string, 0)
	ViewMap["switch"] = make([]string, 0)
}

func DestoryLeftPages(g *gocui.Gui) {
	DestoryTestingPage(g)
	DestoryProjectsPage(g)
	DestorySettingsPage(g)

	ViewMap["testing"] = make([]string, 0)
	ViewMap["projects"] = make([]string, 0)
	ViewMap["settings"] = make([]string, 0)
}
