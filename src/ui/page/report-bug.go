package page

import (
	"fmt"
	zentaoService "github.com/easysoft/zentaoatf/src/service/zentao"
	"github.com/easysoft/zentaoatf/src/ui"
	"github.com/easysoft/zentaoatf/src/ui/widget"
	i118Utils "github.com/easysoft/zentaoatf/src/utils/i118"
	"github.com/easysoft/zentaoatf/src/utils/vari"
	"github.com/fatih/color"
	"github.com/jroimartin/gocui"
	"strings"
)

var filedValMap map[string]int

func InitReportBugPage(resultDir string, caseId string) error {
	DestoryReportBugPage()

	vari.CurrBug, vari.CurrBugStepIds = zentaoService.PrepareBug(resultDir, caseId)
	bug := vari.CurrBug

	//w, h := vari.Cui.Size()
	x := 1
	y := 1

	//var bugVersion string
	//for _, val := range bug.OpenedBuild { // 取字符串值显示
	//	bugVersion = val
	//}

	// title
	left := x
	//right := left + widget.TextWidthFull - 5
	titleInput := widget.NewTextWidget("titleInput", left, y, widget.TextWidthFull-5, bug.Title)
	ui.ViewMap["reportBug"] = append(ui.ViewMap["reportBug"], titleInput.Name())

	//// steps
	//left = right + ui.Space
	//stepsWidth := w - left - 3
	//stepsInput := widget.NewTextareaWidget("stepsInput", left, y, stepsWidth, h-constant.CmdViewHeight-2, bug.Steps)
	//stepsInput.Title = i118Utils.I118Prt.Sprintf("steps")
	//ui.ViewMap["reportBug"] = append(ui.ViewMap["reportBug"], stepsInput.Name())
	//
	//// module
	//y += 3
	//left = x
	//right = left + widget.SelectWidth
	//moduleInput := widget.NewSelectWidgetWithDefault("module", left, y, widget.SelectWidth, 6,
	//	i118Utils.I118Prt.Sprintf("module"),
	//	vari.ZentaoBugFileds.Modules, zentaoService.GetNameById(bug.Module, vari.ZentaoBugFileds.Modules),
	//	bugSelectFieldCheckEvent())
	//ui.ViewMap["reportBug"] = append(ui.ViewMap["reportBug"], moduleInput.Name())
	//
	//// type
	//left = right + ui.Space
	//right = left + widget.SelectWidth
	//typeInput := widget.NewSelectWidgetWithDefault("type", left, y, widget.SelectWidth, 6,
	//	i118Utils.I118Prt.Sprintf("category"),
	//	vari.ZentaoBugFileds.Categories, zentaoService.GetNameById(bug.Type, vari.ZentaoBugFileds.Categories),
	//	bugSelectFieldCheckEvent())
	//ui.ViewMap["reportBug"] = append(ui.ViewMap["reportBug"], typeInput.Name())
	//
	//// version
	//left = right + ui.Space
	//right = left + widget.SelectWidth
	//versionInput := widget.NewSelectWidgetWithDefault("version", left, y, widget.SelectWidth, 6,
	//	i118Utils.I118Prt.Sprintf("version"),
	//	vari.ZentaoBugFileds.Versions, zentaoService.GetNameById(bugVersion, vari.ZentaoBugFileds.Versions),
	//	bugSelectFieldCheckEvent())
	//ui.ViewMap["reportBug"] = append(ui.ViewMap["reportBug"], versionInput.Name())
	//
	//// severity
	//y += 7
	//left = x
	//right = left + widget.SelectWidth
	//severityInput := widget.NewSelectWidgetWithDefault("severity", left, y, widget.SelectWidth, 6,
	//	i118Utils.I118Prt.Sprintf("severity"),
	//	vari.ZentaoBugFileds.Severities, zentaoService.GetNameById(bug.Severity, vari.ZentaoBugFileds.Severities),
	//	bugSelectFieldCheckEvent())
	//ui.ViewMap["reportBug"] = append(ui.ViewMap["reportBug"], severityInput.Name())
	//
	//// priority
	//left = right + ui.Space
	//right = left + widget.SelectWidth
	//priorityInput := widget.NewSelectWidgetWithDefault("priority", left, y, widget.SelectWidth, 6,
	//	i118Utils.I118Prt.Sprintf("priority"),
	//	vari.ZentaoBugFileds.Priorities, zentaoService.GetNameById(bug.Pri, vari.ZentaoBugFileds.Priorities),
	//	bugSelectFieldCheckEvent())
	//ui.ViewMap["reportBug"] = append(ui.ViewMap["reportBug"], priorityInput.Name())
	//
	//// msg
	//y += 7
	//left = x
	//reportBugMsg := widget.NewPanelWidget("reportBugMsg", left, y, widget.TextWidthFull-5, 2, "")
	//reportBugMsg.Frame = false
	//ui.ViewMap["reportBug"] = append(ui.ViewMap["reportBug"], reportBugMsg.Name())
	//
	//// buttons
	//y += 5
	//buttonX := x + widget.SelectWidth + ui.Space
	//submitInput := widget.NewButtonWidgetAutoWidth("submitInput", buttonX, y,
	//	i118Utils.I118Prt.Sprintf("submit"), reportBug)
	//ui.ViewMap["reportBug"] = append(ui.ViewMap["reportBug"], submitInput.Name())
	//
	//cancelReportBugInput := widget.NewButtonWidgetAutoWidth("cancelReportBugInput",
	//	buttonX+11, y, i118Utils.I118Prt.Sprintf("cancel"), cancelReportBug)
	//ui.ViewMap["reportBug"] = append(ui.ViewMap["reportBug"], cancelReportBugInput.Name())
	//
	//ui.BindEventForInputWidgets(ui.ViewMap["reportBug"])

	vari.Cui.SetCurrentView("titleInput")

	return nil
}

func reportBug(g *gocui.Gui, v *gocui.View) error {
	bug := vari.CurrBug

	titleView, _ := g.View("titleInput")
	stepsView, _ := g.View("stepsInput")
	moduleView, _ := g.View("module")
	typeView, _ := g.View("type")
	versionView, _ := g.View("version")
	severityView, _ := g.View("severity")
	priorityView, _ := g.View("priority")

	title := strings.TrimSpace(titleView.Buffer())
	stepsStr := strings.TrimSpace(stepsView.Buffer())

	moduleStr := strings.TrimSpace(ui.GetSelectedRowVal(moduleView))
	typeStr := strings.TrimSpace(ui.GetSelectedRowVal(typeView))
	versionStr := strings.TrimSpace(ui.GetSelectedRowVal(versionView))
	severityStr := strings.TrimSpace(ui.GetSelectedRowVal(severityView))
	priorityStr := strings.TrimSpace(ui.GetSelectedRowVal(priorityView))

	if title == "" {
		v, _ := vari.Cui.View("reportBugMsg")
		color.New(color.FgMagenta).Fprintf(v, i118Utils.ReadI18nJson("title_cannot_be_empty"))
		return nil
	}

	bug.Title = title
	bug.Steps = strings.Replace(stepsStr, "\n", "<br/>", -1)

	bug.Type = zentaoService.GetIdByName(typeStr, vari.ZentaoBugFileds.Categories)

	bug.Module = zentaoService.GetIdByName(moduleStr, vari.ZentaoBugFileds.Modules)

	versionKey := zentaoService.GetIdByName(versionStr, vari.ZentaoBugFileds.Versions)
	build := make(map[string]string)
	if versionKey == "trunk" {
		build["0"] = "trunk"
	} else {
		build[versionKey] = versionStr
	}
	bug.OpenedBuild = build

	bug.Severity = zentaoService.GetIdByName(severityStr, vari.ZentaoBugFileds.Severities)
	bug.Pri = zentaoService.GetIdByName(priorityStr, vari.ZentaoBugFileds.Priorities)

	vari.CurrBug = bug
	ok, msg := zentaoService.CommitBug()

	msgView, _ := vari.Cui.View("reportBugMsg")
	msgView.Clear()

	if ok {
		color.New(color.FgGreen).Fprintf(msgView, msg)

		vari.Cui.DeleteView("submitInput")

		cancelReportBugInput, _ := vari.Cui.View("cancelReportBugInput")
		cancelReportBugInput.Clear()
		fmt.Fprint(cancelReportBugInput, " "+i118Utils.I118Prt.Sprintf("close"))
	} else {
		color.New(color.FgMagenta).Fprintf(msgView, msg)
	}

	return nil
}

func bugSelectFieldCheckEvent() func(g *gocui.Gui, v *gocui.View) error {
	return func(g *gocui.Gui, v *gocui.View) error {
		name := v.Name()

		g.SetCurrentView(name)

		return nil
	}
}

func init() {
	filedValMap = make(map[string]int)
}

func cancelReportBug(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func DestoryReportBugPage() {
	for _, v := range ui.ViewMap["reportBug"] {
		vari.Cui.DeleteView(v)
		vari.Cui.DeleteKeybindings(v)
	}
}
