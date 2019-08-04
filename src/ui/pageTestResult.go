package ui

import (
	"fmt"
	"github.com/easysoft/zentaoatf/src/biz"
	"github.com/easysoft/zentaoatf/src/script"
	"github.com/easysoft/zentaoatf/src/utils"
	"github.com/jroimartin/gocui"
	"strings"
)

var CurrRun string
var CurrResult string

var runViews []string

func showRun(g *gocui.Gui, v *gocui.View) error {
	DestoryContentPanel()
	HighlightTab(v.Name(), tabs)

	h := utils.MainViewHeight / 2
	maxX, _ := g.Size()

	panelResultList := NewPanelWidget("panelResultList", utils.LeftWidth, 2, 50, h, "")
	ViewMap["testing"] = append(ViewMap["testing"], panelResultList.Name())
	runViews = append(runViews, panelResultList.Name())

	panelCaseList := NewPanelWidget("panelCaseList", utils.LeftWidth, h+2, 50, utils.MainViewHeight-h, "")
	ViewMap["testing"] = append(ViewMap["testing"], panelCaseList.Name())
	runViews = append(runViews, panelCaseList.Name())

	panelCaseResult := NewPanelWidget("panelCaseResult", utils.LeftWidth+50, 2,
		maxX-utils.LeftWidth-51, utils.MainViewHeight, "")
	ViewMap["testing"] = append(ViewMap["testing"], panelCaseResult.Name())
	runViews = append(runViews, panelCaseResult.Name())

	for idx, v := range runViews {
		setViewScroll(v)

		if idx < 2 {
			setViewLineHighlight(v)
		}
	}

	setViewLineSelected("panelResultList", selectResultEvent)
	setViewLineSelected("panelCaseList", selectCaseEvent)

	results := script.LoadTestResults(CurrAsset)
	fmt.Fprintln(panelResultList, strings.Join(results, "\n"))

	return nil
}

func init() {

}

func DestoryRunPanel() {
	for _, v := range runViews {
		utils.Cui.DeleteView(v)
		utils.Cui.DeleteKeybindings(v)
	}
}

func selectAssetEvent(g *gocui.Gui, v *gocui.View) error {
	clearPanelCaseResult()

	var line string
	var err error

	_, cy := v.Cursor()
	if line, err = v.Line(cy); err != nil {
		return nil
	}
	line = strings.TrimSpace(line)
	if strings.Index(line, ".") < 0 {
		utils.PrintToMainNoScroll("")
		return nil
	}
	CurrAsset = utils.ScriptDir + line

	// show
	if len(tabs) == 0 {
		HideHelp()
		showTab()
	}

	defaultTab, _ := g.View("tabContentView")
	showContent(g, defaultTab)

	return nil
}

func selectResultEvent(g *gocui.Gui, v *gocui.View) error {
	clearPanelCaseResult()

	v.Highlight = true

	line, _ := SelectLine(v, ".*")
	CurrResult = line
	content := script.GetTestResult(CurrAsset, line)

	panelCaseList, _ := g.View("panelCaseList")
	panelCaseList.Clear()
	fmt.Fprintln(panelCaseList, strings.Join(content, "\n"))

	maxX, _ := g.Size()
	uploadButton := NewButtonWidgetAutoWidth("uploadButton", maxX-35, 0, "[Upload Result]", uploadResult)
	uploadButton.Frame = false
	runViews = append(runViews, uploadButton.Name())

	return nil
}

func selectCaseEvent(g *gocui.Gui, v *gocui.View) error {
	v.Highlight = true

	caseLine, _ := SelectLine(v, ".*")

	content := script.GetCheckpointsResult(CurrAsset, CurrResult, caseLine)
	panelCaseResult, _ := g.View("panelCaseResult")
	panelCaseResult.Clear()
	fmt.Fprintln(panelCaseResult, content)

	// show submit bug button
	maxX, _ := g.Size()
	bugButton := NewButtonWidgetAutoWidth("bugButton", maxX-18, 0, "[Report Bug]", reportBug)
	bugButton.Frame = false
	runViews = append(runViews, bugButton.Name())

	return nil
}

func clearPanelCaseResult() {
	panelCaseResult, _ := utils.Cui.View("panelCaseResult")
	if panelCaseResult != nil {
		panelCaseResult.Clear()
	}
	utils.Cui.DeleteView("bugButton")
}

func uploadResult(g *gocui.Gui, v *gocui.View) error {
	caseList := script.GetTestResult(CurrAsset, CurrResult)

	biz.SubmitResult(caseList)

	return nil
}

func reportBug(g *gocui.Gui, v *gocui.View) error {

	return nil
}
