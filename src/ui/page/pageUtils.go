package page

import (
	"github.com/easysoft/zentaoatf/src/ui"
	print2 "github.com/easysoft/zentaoatf/src/utils/print"
	"github.com/easysoft/zentaoatf/src/utils/vari"
)

func DestoryLeftPages() {
	print2.ClearSide()

	DestoryTestPage()
	DestoryProjectsPage()
	DestorySettingsPage()

	ui.ViewMap["testing"] = make([]string, 0)
	ui.ViewMap["projects"] = make([]string, 0)
	ui.ViewMap["settings"] = make([]string, 0)
}

func DestoryRightPages() {
	mainView, err := vari.Cui.View("main")
	if err == nil {
		mainView.Clear()
	}

	DestoryImportPage()
	DestorySwitchPage()

	ui.ViewMap["import"] = make([]string, 0)
	ui.ViewMap["switch"] = make([]string, 0)
}
