package ui

import (
	"github.com/jroimartin/gocui"
)

func DestoryRightPages(g *gocui.Gui) {
	for key, _ := range ViewMap {
		if key == "import" {
			DestoryImportPage(g)
		} else if key == "switch" {
			DestorySwitchPage(g)
		}
	}
	ViewMap["import"] = make([]string, 0)
	ViewMap["switch"] = make([]string, 0)
}

func DestoryLeftPages(g *gocui.Gui) {
	for key, _ := range ViewMap {
		if key == "testings" {
			//DestoryImportPage(g)

		} else if key == "projects" {
			//DestorySwitchPage(g)
		} else if key == "settings" {
			DestorySettingsPage(g)
		}
	}

	ViewMap["testings"] = make([]string, 0)
	ViewMap["projects"] = make([]string, 0)
	ViewMap["settings"] = make([]string, 0)
}
