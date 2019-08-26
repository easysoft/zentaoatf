package page

import (
	"github.com/easysoft/zentaoatf/src/ui"
	"github.com/easysoft/zentaoatf/src/utils/vari"
	"github.com/jroimartin/gocui"
)

var CurrProjectId string

//var projectHistories []model.WorkHistory

func InitProjectsPage() error {
	//his := vari.Config.WorkHistories[0]
	//id, _, _ := getProjectInfo(his)
	//CurrProjectId = id
	//
	//y := 2
	//for _, his := range vari.Config.WorkHistories {
	//	id, label, _ := getProjectInfo(his)
	//
	//	hisView := widget.NewLabelWidgetAutoWidth(id, 0, y, label)
	//	ui.ViewMap["projects"] = append(ui.ViewMap["projects"], hisView.Name())
	//
	//	y += 1
	//}
	//keybindingProjectsButton()

	return nil
}

func keybindingProjectsButton() error {
	//for _, his := range vari.Config.WorkHistories {
	//	id, _, _ := getProjectInfo(his)
	//	if err := vari.Cui.SetKeybinding(id, gocui.MouseLeft, gocui.ModNone, toggleProjectsButton); err != nil {
	//		return err
	//	}
	//}

	return nil
}

func toggleProjectsButton(g *gocui.Gui, v *gocui.View) error {
	CurrProjectId = v.Name()
	SelectProjectsButton()

	return nil
}

func SelectProjectsButton() {
	//for _, his := range vari.Config.WorkHistories {
	//	id, _, _ := getProjectInfo(his)
	//
	//	v, err := vari.Cui.View(id)
	//	if err == nil {
	//		if id == CurrProjectId {
	//			v.Highlight = true
	//			v.SelBgColor = gocui.ColorWhite
	//			v.SelFgColor = gocui.ColorBlack
	//
	//			//printForSwitch(his) // 显示项目信息
	//			showWitchButton()
	//		} else {
	//			v.Highlight = false
	//			v.SelBgColor = gocui.ColorBlack
	//			v.SelFgColor = gocui.ColorDefault
	//		}
	//	}
	//}
}

func showWitchButton() error {
	//maxX, _ := vari.Cui.Size()

	//switchButton := widget.NewButtonWidgetAutoWidth("switchButton", maxX-15, 1, "Switch To", switchProject)
	//ui.ViewMap["projects"] = append(ui.ViewMap["projects"], switchButton.Name())

	return nil
}

func init() {

}

func DestoryProjectsPage() {
	for _, v := range ui.ViewMap["projects"] {
		vari.Cui.DeleteView(v)
		vari.Cui.DeleteKeybindings(v)
	}
}
