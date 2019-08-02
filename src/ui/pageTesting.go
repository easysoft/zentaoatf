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
var tabs []string

func InitTestingPage(g *gocui.Gui) error {
	// left
	caseFiles, suitesFiles := loadTestAssets()
	dir := utils.Prefer.WorkDir + utils.GenDir

	content := "Test Suite:" + "\n"
	for _, suitePath := range suitesFiles {
		suiteName := strings.Replace(suitePath, dir, "", -1)
		content += "  " + suiteName + "\n"
	}

	content += "Test Scripts:" + "\n"
	for _, casePath := range caseFiles {
		caseName := strings.Replace(casePath, dir, "", -1)
		content += "  " + caseName + "\n"
	}
	utils.PrintToSide(g, content)

	// right
	setViewScroll(g, "side")
	setViewLineSelected(g, "side", selectAssetEvent)

	return nil
}

func selectAssetEvent(g *gocui.Gui, v *gocui.View) error {
	var line string
	var err error

	_, cy := v.Cursor()
	if line, err = v.Line(cy); err != nil {
		return nil
	}
	line = strings.TrimSpace(line)
	if strings.Index(line, ".") < 0 {
		utils.PrintToMainNoScroll(g, "")
		return nil
	}
	CurrAsset = utils.Prefer.WorkDir + utils.GenDir + line

	showAsset(g)

	return nil
}
func showAsset(g *gocui.Gui) {
	if len(tabs) == 0 {
		HideHelp(g)
		showTab(g)
	}

	defaultTab, _ := g.View("tabContentView")
	showContent(g, defaultTab)
}

func showTab(g *gocui.Gui) error {
	x := utils.LeftWidth + 1
	tabContentView := NewLabelWidgetAutoWidth(g, "tabContentView", x, 0, "Content")
	ViewMap["testing"] = append(ViewMap["testing"], tabContentView.Name())
	tabs = append(tabs, tabContentView.Name())
	if err := g.SetKeybinding("tabContentView", gocui.MouseLeft, gocui.ModNone, showContent); err != nil {
		return nil
	}

	tabResultView := NewLabelWidgetAutoWidth(g, "tabResultView", x+12, 0, "Results")
	ViewMap["testing"] = append(ViewMap["testing"], tabResultView.Name())
	tabs = append(tabs, tabResultView.Name())
	if err := g.SetKeybinding("tabResultView", gocui.MouseLeft, gocui.ModNone, showRun); err != nil {
		return nil
	}

	return nil
}

func showContent(g *gocui.Gui, v *gocui.View) error {
	DestoryRunPanel(g)
	HighlightTab(v.Name(), tabs)

	panelFileContent, _ := g.View("panelFileContent")
	if panelFileContent != nil {
		panelFileContent.Clear()
	} else {
		maxX, _ := g.Size()
		panelFileContent = NewPanelWidget(g, "panelFileContent", utils.LeftWidth, 2,
			maxX-utils.LeftWidth-1, utils.MainViewHeight, "")
		ViewMap["testing"] = append(ViewMap["testing"], panelFileContent.Name())

		runButton := NewButtonWidgetAutoWidth(g, "runButton", maxX-10, 0, "Run", run)
		runButton.Frame = false
		ViewMap["testing"] = append(ViewMap["testing"], runButton.Name())
	}

	content := utils.ReadFile(CurrAsset)
	fmt.Fprintln(panelFileContent, content)

	return nil
}

func showRun(g *gocui.Gui, v *gocui.View) error {
	DestoryContentPanel(g)
	HighlightTab(v.Name(), tabs)

	h := utils.MainViewHeight / 2
	maxX, _ := g.Size()

	panelResultList := NewPanelWidget(g, "panelResultList", utils.LeftWidth, 2,
		60, h, "panelResultList")
	ViewMap["testing"] = append(ViewMap["testing"], panelResultList.Name())

	panelCaseList := NewPanelWidget(g, "panelCaseList", utils.LeftWidth, h+2,
		60, utils.MainViewHeight-h, "panelCaseList")
	ViewMap["testing"] = append(ViewMap["testing"], panelCaseList.Name())

	panelCaseResult := NewPanelWidget(g, "panelCaseResult", utils.LeftWidth+60, 2,
		maxX-utils.LeftWidth-61, utils.MainViewHeight, "panelCaseResult")
	ViewMap["testing"] = append(ViewMap["testing"], panelCaseResult.Name())

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
	config := utils.ReadCurrConfig()
	ext := script.GetLangMap()[config.LangType]["extName"]

	caseFiles, _ := utils.GetAllFiles(utils.Prefer.WorkDir+utils.GenDir, ext)
	suitesFiles, _ := utils.GetAllFiles(utils.Prefer.WorkDir+utils.GenDir, "suite")

	return caseFiles, suitesFiles
}

func init() {

}

func DestoryTestingPage(g *gocui.Gui) {
	g.DeleteKeybindings("side")
	for _, v := range ViewMap["testing"] {
		g.DeleteView(v)
		g.DeleteKeybindings(v)
	}
}

func DestoryContentPanel(g *gocui.Gui) {
	for _, v := range []string{"panelResultList", "runButton"} {
		g.DeleteView(v)
		g.DeleteKeybindings(v)
	}
}
func DestoryRunPanel(g *gocui.Gui) {
	for _, v := range []string{"panelFileContent", "panelCaseList", "panelCaseResult"} {
		g.DeleteView(v)
		g.DeleteKeybindings(v)
	}
}
