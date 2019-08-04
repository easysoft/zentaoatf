package ui

import (
	"github.com/easysoft/zentaoatf/src/utils"
)

func DestoryLeftPages() {
	utils.ClearSide()

	DestoryTestPage()
	DestoryProjectsPage()
	DestorySettingsPage()

	ViewMap["testing"] = make([]string, 0)
	ViewMap["projects"] = make([]string, 0)
	ViewMap["settings"] = make([]string, 0)
}

func DestoryRightPages() {
	mainView, err := utils.Cui.View("main")
	if err == nil {
		mainView.Clear()
	}

	DestoryImportPage()
	DestorySwitchPage()

	ViewMap["import"] = make([]string, 0)
	ViewMap["switch"] = make([]string, 0)
}
