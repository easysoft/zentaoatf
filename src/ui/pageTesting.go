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
var CurrRun string
var CurrResult string

var tabs []string
var contentViews []string
var runViews []string

func InitTestingPage(g *gocui.Gui) error {
	// left
	caseFiles, suitesFiles := script.LoadTestAssets()
	dir := utils.Prefer.WorkDir + utils.ScriptDir

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
		panelFileContent = NewPanelWidget(g, utils.CuiRunOutputView, utils.LeftWidth, 2,
			maxX-utils.LeftWidth-1, utils.MainViewHeight, "")
		ViewMap["testing"] = append(ViewMap["testing"], panelFileContent.Name())
		contentViews = append(contentViews, panelFileContent.Name())
		setViewScroll(g, panelFileContent.Name())

		runButton := NewButtonWidgetAutoWidth(g, "runButton", maxX-10, 0, "Run", run)
		runButton.Frame = false
		contentViews = append(contentViews, runButton.Name())
	}

	panelFileContent.Clear()
	panelFileContent.SetOrigin(0, 0)
	content := utils.ReadFile(CurrAsset)
	fmt.Fprintln(panelFileContent, content)

	return nil
}

func showRun(g *gocui.Gui, v *gocui.View) error {
	DestoryContentPanel(g)
	HighlightTab(v.Name(), tabs)

	h := utils.MainViewHeight / 2
	maxX, _ := g.Size()

	panelResultList := NewPanelWidget(g, "panelResultList", utils.LeftWidth, 2, 60, h, "")
	ViewMap["testing"] = append(ViewMap["testing"], panelResultList.Name())
	runViews = append(runViews, panelResultList.Name())

	panelCaseList := NewPanelWidget(g, "panelCaseList", utils.LeftWidth, h+2,
		60, utils.MainViewHeight-h, "")
	ViewMap["testing"] = append(ViewMap["testing"], panelCaseList.Name())
	runViews = append(runViews, panelCaseList.Name())

	panelCaseResult := NewPanelWidget(g, "panelCaseResult", utils.LeftWidth+60, 2,
		maxX-utils.LeftWidth-61, utils.MainViewHeight, "panelCaseResult")
	ViewMap["testing"] = append(ViewMap["testing"], panelCaseResult.Name())
	runViews = append(runViews, panelCaseResult.Name())

	for idx, v := range runViews {
		setViewScroll(g, v)

		if idx < 2 {
			setHighlight(g, v)
		}
	}

	setViewLineSelected(g, "panelResultList", selectResultEvent)
	setViewLineSelected(g, "panelCaseList", selectCaseEvent)

	results := script.LoadTestResults(CurrAsset)
	fmt.Fprintln(panelResultList, strings.Join(results, "\n"))

	return nil
}

func run(g *gocui.Gui, v *gocui.View) error {
	if _, err := g.SetCurrentView("main"); err != nil {
		return err
	}

	utils.PrintToCmd(g, fmt.Sprintf("#atf run -d %s -f %s", utils.Prefer.WorkDir, CurrAsset))
	output, _ := g.View(utils.CuiRunOutputView)
	output.Clear()
	action.Run(utils.Prefer.WorkDir, []string{CurrAsset}, "")

	return nil
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
	for _, v := range contentViews {
		g.DeleteView(v)
		g.DeleteKeybindings(v)
	}
}
func DestoryRunPanel(g *gocui.Gui) {
	for _, v := range runViews {
		g.DeleteView(v)
		g.DeleteKeybindings(v)
	}
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
	CurrAsset = utils.Prefer.WorkDir + utils.ScriptDir + line

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

func selectResultEvent(g *gocui.Gui, v *gocui.View) error {
	line, _ := SelectLine(v, ".*")
	content := script.GetTestResult(CurrAsset, line)

	panelCaseList, _ := g.View("panelCaseList")
	panelCaseList.Clear()
	fmt.Fprintln(panelCaseList, content)

	return nil
}

func selectCaseEvent(g *gocui.Gui, v *gocui.View) error {

	return nil
}
