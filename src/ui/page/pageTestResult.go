package page

import (
	"fmt"
	"github.com/easysoft/zentaoatf/src/service/script"
	testingService "github.com/easysoft/zentaoatf/src/service/testing"
	zentaoService "github.com/easysoft/zentaoatf/src/service/zentao"
	"github.com/easysoft/zentaoatf/src/ui"
	"github.com/easysoft/zentaoatf/src/ui/widget"
	constant "github.com/easysoft/zentaoatf/src/utils/const"
	logUtils "github.com/easysoft/zentaoatf/src/utils/log"
	"github.com/easysoft/zentaoatf/src/utils/vari"
	"github.com/jroimartin/gocui"
	"strconv"
	"strings"
)

var runViews []string

func showRun(g *gocui.Gui, v *gocui.View) error {
	DestoryContentPanel()
	ui.HighlightTab(v.Name(), tabs)

	h := vari.MainViewHeight / 2
	logUtils.PrintToCmd(fmt.Sprintf("==%d==", h))
	maxX, _ := g.Size()

	panelResultList := widget.NewPanelWidget("panelResultList", constant.LeftWidth, 2, 50, h-2, "")
	ui.ViewMap["testing"] = append(ui.ViewMap["testing"], panelResultList.Name())
	runViews = append(runViews, panelResultList.Name())

	panelCaseList := widget.NewPanelWidget("panelCaseList", constant.LeftWidth, h, 50, vari.MainViewHeight-h, "")
	ui.ViewMap["testing"] = append(ui.ViewMap["testing"], panelCaseList.Name())
	runViews = append(runViews, panelCaseList.Name())

	panelCaseResult := widget.NewPanelWidget("panelCaseResult", constant.LeftWidth+50, 2,
		maxX-constant.LeftWidth-51, vari.MainViewHeight-2, "")
	ui.ViewMap["testing"] = append(ui.ViewMap["testing"], panelCaseResult.Name())
	runViews = append(runViews, panelCaseResult.Name())

	ui.SupportScroll("panelResultList")
	ui.SupportScroll("panelCaseList")
	ui.SupportScroll("panelCaseResult")

	ui.SupportRowHighlight("panelResultList")
	ui.SupportRowHighlight("panelCaseList")

	ui.AddLineSelectedEvent("panelResultList", selectResultEvent)
	ui.AddLineSelectedEvent("panelCaseList", selectCaseEvent)

	results := scriptService.LoadTestResults(vari.CurrScriptFile)
	fmt.Fprintln(panelResultList, strings.Join(results, "\n"))

	return nil
}

func init() {

}

func selectResultEvent(g *gocui.Gui, v *gocui.View) error {
	clearPanelCaseResult()

	v.Highlight = true

	line := ui.GetSelectedRowVal(v)
	vari.CurrResultDate = line

	content := make([]string, 0)
	report := testingService.GetTestTestReportForSubmit(vari.CurrScriptFile, vari.CurrResultDate)
	for _, cs := range report.Cases {
		id := cs.Id
		title := cs.Title
		result := cs.Status

		str := fmt.Sprintf("%d-%s: %s", id, title, result)
		content = append(content, str)
	}

	panelCaseList, _ := g.View("panelCaseList")
	panelCaseList.Clear()
	fmt.Fprintln(panelCaseList, strings.Join(content, "\n"))

	maxX, _ := g.Size()
	uploadButton := widget.NewButtonWidgetAutoWidth("uploadButton", maxX-35, 0, "[Upload Result]", toUploadResult)
	uploadButton.Frame = false
	runViews = append(runViews, uploadButton.Name())

	return nil
}

func selectCaseEvent(g *gocui.Gui, v *gocui.View) error {
	v.Highlight = true

	caseLine := ui.GetSelectedRowVal(v)
	caseIdStr := strings.Split(caseLine, "-")[0]
	caseId, _ := strconv.Atoi(caseIdStr)
	vari.CurrCaseId = caseId

	content := make([]string, 0)
	report := testingService.GetTestTestReportForSubmit(vari.CurrScriptFile, vari.CurrResultDate)
	for _, cs := range report.Cases {
		if cs.Id == caseId {
			for _, step := range cs.Steps {
				content = append(content, testingService.GetStepText(step))
				content = append(content, "")
			}
		}
	}

	panelCaseResult, _ := g.View("panelCaseResult")
	panelCaseResult.Clear()
	fmt.Fprintln(panelCaseResult, strings.Join(content, "\n"))

	// show submit bug button
	maxX, _ := g.Size()
	bugButton := widget.NewButtonWidgetAutoWidth("bugButton", maxX-18, 0, "[Report Bug]", toReportBug)
	bugButton.Frame = false
	runViews = append(runViews, bugButton.Name())

	return nil
}

func clearPanelCaseResult() {
	panelCaseResult, _ := vari.Cui.View("panelCaseResult")
	if panelCaseResult != nil {
		panelCaseResult.Clear()
	}
	vari.Cui.DeleteView("bugButton")
}

func toUploadResult(g *gocui.Gui, v *gocui.View) error {
	zentaoService.SubmitResult()

	return nil
}

func toReportBug(g *gocui.Gui, v *gocui.View) error {
	InitReportBugPage()

	return nil
}

func DestoryRunPanel() {
	for _, v := range runViews {
		vari.Cui.DeleteView(v)
		vari.Cui.DeleteKeybindings(v)
	}
}
