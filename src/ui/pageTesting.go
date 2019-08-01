package ui

import (
	"fmt"
	"github.com/easysoft/zentaoatf/src/action"
	"github.com/easysoft/zentaoatf/src/script"
	"github.com/easysoft/zentaoatf/src/utils"
	"github.com/jroimartin/gocui"
	"strings"
)

var CurrAsset string

func InitTestingPage(g *gocui.Gui) error {
	caseFiles, suitesFiles := loadTestAssets()
	dir := utils.Prefer.WorkDir + utils.GenDir

	// left asserts
	y := 2
	suiteLabel := NewLabelWidget(g, "suiteLabel", 0, y, "Test Suite")
	ViewMap["testing"] = append(ViewMap["testing"], suiteLabel.Name())

	y += 1
	for _, suitePath := range suitesFiles {
		suiteName := strings.Replace(suitePath, dir, "", -1)

		suiteView := NewButtonWidgetNoBorderAutoWidth(g, suitePath, 0, y, suiteName, selectTestingItem)
		ViewMap["testing"] = append(ViewMap["testing"], suiteView.Name())

		y += 1
	}

	y += 1
	caseLabel := NewLabelWidget(g, "caseLabel", 0, y, "Test Script")
	ViewMap["testing"] = append(ViewMap["testing"], caseLabel.Name())

	y += 1
	for _, casePath := range caseFiles {
		caseName := strings.Replace(casePath, dir, "", -1)

		caseView := NewButtonWidgetNoBorderAutoWidth(g, casePath, 0, y, caseName, selectTestingItem)
		ViewMap["testing"] = append(ViewMap["testing"], caseView.Name())

		y += 1
	}

	return nil
}

func selectTestingItem(g *gocui.Gui, view *gocui.View) error {
	HideHelp(g)
	CurrAsset = view.Name()

	for _, name := range ViewMap["testing"] {
		v, err := g.View(name)
		if err != nil {
			return err
		}

		if v.Name() == CurrAsset {
			v.Highlight = true
			v.SelBgColor = gocui.ColorWhite
			v.SelFgColor = gocui.ColorBlack
		} else {
			v.Highlight = false
			v.SelBgColor = gocui.ColorBlack
			v.SelFgColor = gocui.ColorDefault
		}
	}

	showRunButton(g)
	content := utils.ReadFile(CurrAsset)
	utils.PrintToMainNoScroll(g, content)

	return nil
}

func showRunButton(g *gocui.Gui) error {
	maxX, _ := g.Size()

	runButton := NewButtonWidgetAutoWidth(g, "runButton", maxX-10, 1, "Run", run)
	ViewMap["testing"] = append(ViewMap["testing"], runButton.Name())

	return nil
}

func run(g *gocui.Gui, v *gocui.View) error {
	if _, err := g.SetCurrentView("main"); err != nil {
		return err
	}

	utils.PrintToCmd(g, fmt.Sprintf("#atf run -d %s -f %s", utils.Prefer.WorkDir, CurrAsset))
	utils.PrintToMain(g, "")
	action.Run(utils.Prefer.WorkDir, []string{CurrAsset}, "")

	return nil
}

func loadTestAssets() ([]string, []string) {
	config := utils.ReadConfig()
	ext := script.GetLangMap()[config.LangType]["extName"]

	caseFiles, _ := utils.GetAllFiles(utils.Prefer.WorkDir+utils.GenDir, ext)
	suitesFiles, _ := utils.GetAllFiles(utils.Prefer.WorkDir+utils.GenDir, "suite")

	return caseFiles, suitesFiles
}

func init() {

}

func DestoryTestingPage(g *gocui.Gui) {
	for _, v := range ViewMap["testing"] {
		g.DeleteView(v)
		g.DeleteKeybindings(v)
	}
}
