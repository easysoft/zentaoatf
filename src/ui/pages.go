package ui

import (
	print2 "github.com/easysoft/zentaoatf/src/utils/print"
	"github.com/easysoft/zentaoatf/src/utils/vari"
)

func DestoryLeftPages() {
	print2.ClearSide()

	DestoryTestPage()
	DestoryProjectsPage()
	DestorySettingsPage()

	ViewMap["testing"] = make([]string, 0)
	ViewMap["projects"] = make([]string, 0)
	ViewMap["settings"] = make([]string, 0)
}

func DestoryRightPages() {
	mainView, err := vari.Cui.View("main")
	if err == nil {
		mainView.Clear()
	}

	DestoryImportPage()
	DestorySwitchPage()

	ViewMap["import"] = make([]string, 0)
	ViewMap["switch"] = make([]string, 0)
}
