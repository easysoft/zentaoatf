package page

import (
	"fmt"
	"github.com/easysoft/zentaoatf/src/model"
	zentaoService "github.com/easysoft/zentaoatf/src/service/zentao"
	"github.com/easysoft/zentaoatf/src/ui"
	"github.com/easysoft/zentaoatf/src/ui/widget"
	"github.com/easysoft/zentaoatf/src/utils/log"
	"github.com/easysoft/zentaoatf/src/utils/vari"
	"github.com/fatih/color"
	"github.com/jroimartin/gocui"
	"strings"
)

var (
	bug      model.Bug
	idInTask string
	stepIds  string
)

var filedValMap map[string]int

func InitReportBugPage() error {
	DestoryReportBugPage()

	zentaoService.GetBugFiledOptions()
	logUtils.PrintToCmd(fmt.Sprintf("%v", vari.ZentaoBugFileds.Modules))

	bug, idInTask, stepIds = zentaoService.GenBug()

	maxX, maxY := vari.Cui.Size()
	pageWidth := 120
	pageHeight := 27
	x := maxX/2 - 60
	y := maxY/2 - 14

	var bugVersion string
	for _, val := range bug.OpenedBuild { // 取字符串值显示
		bugVersion = val
	}

	// panel
	reportBugPanel := widget.NewPanelWidget("reportBugPanel", x, y, pageWidth, pageHeight, "")
	ui.ViewMap["reportBug"] = append(ui.ViewMap["reportBug"], reportBugPanel.Name())

	// title
	y += 1
	left := x + 2
	right := left + widget.TextWidthFull - 5
	titleInput := widget.NewTextWidget("titleInput", left, y, widget.TextWidthFull-5, bug.Title)
	ui.ViewMap["reportBug"] = append(ui.ViewMap["reportBug"], titleInput.Name())

	// steps
	left = right + ui.Space
	stepsWidth := pageWidth - left - ui.Space + x
	stepsInput := widget.NewTextareaWidget("stepsInput", left, y, stepsWidth, pageHeight-2, bug.Steps)
	stepsInput.Title = "Steps"
	ui.ViewMap["reportBug"] = append(ui.ViewMap["reportBug"], stepsInput.Name())

	// module
	y += 3
	left = x + 2
	right = left + widget.SelectWidth
	moduleInput := widget.NewSelectWidgetWithDefault("module", left, y, widget.SelectWidth, 6, "Module",
		vari.ZentaoBugFileds.Modules, zentaoService.GetNameById(bug.Module, vari.ZentaoBugFileds.Modules),
		bugSelectFieldCheckEvent())
	ui.ViewMap["reportBug"] = append(ui.ViewMap["reportBug"], moduleInput.Name())

	// type
	left = right + ui.Space
	right = left + widget.SelectWidth
	typeInput := widget.NewSelectWidgetWithDefault("type", left, y, widget.SelectWidth, 6, "Category",
		vari.ZentaoBugFileds.Categories, zentaoService.GetNameById(bug.Type, vari.ZentaoBugFileds.Categories),
		bugSelectFieldCheckEvent())
	ui.ViewMap["reportBug"] = append(ui.ViewMap["reportBug"], typeInput.Name())

	// version
	left = right + ui.Space
	right = left + widget.SelectWidth
	versionInput := widget.NewSelectWidgetWithDefault("version", left, y, widget.SelectWidth, 6, "Version",
		vari.ZentaoBugFileds.Versions, zentaoService.GetNameById(bugVersion, vari.ZentaoBugFileds.Versions),
		bugSelectFieldCheckEvent())
	ui.ViewMap["reportBug"] = append(ui.ViewMap["reportBug"], versionInput.Name())

	// severity
	y += 7
	left = x + 2
	right = left + widget.SelectWidth
	severityInput := widget.NewSelectWidgetWithDefault("severity", left, y, widget.SelectWidth, 6, "Severity",
		vari.ZentaoBugFileds.Severities, zentaoService.GetNameById(bug.Severity, vari.ZentaoBugFileds.Severities),
		bugSelectFieldCheckEvent())
	ui.ViewMap["reportBug"] = append(ui.ViewMap["reportBug"], severityInput.Name())

	// priority
	left = right + ui.Space
	right = left + widget.SelectWidth
	priorityInput := widget.NewSelectWidgetWithDefault("priority", left, y, widget.SelectWidth, 6, "Priority",
		vari.ZentaoBugFileds.Priorities, zentaoService.GetNameById(bug.Pri, vari.ZentaoBugFileds.Priorities),
		bugSelectFieldCheckEvent())
	ui.ViewMap["reportBug"] = append(ui.ViewMap["reportBug"], priorityInput.Name())

	// msg
	y += 7
	left = x + 2
	reportBugMsg := widget.NewPanelWidget("reportBugMsg", left, y, widget.TextWidthFull-5, 2, "")
	reportBugMsg.Frame = false
	ui.ViewMap["reportBug"] = append(ui.ViewMap["reportBug"], reportBugMsg.Name())

	// buttons
	y += 6
	buttonX := maxX/2 - 49 + widget.LabelWidthSmall + ui.Space
	submitInput := widget.NewButtonWidgetAutoWidth("submitInput", buttonX, y, "Submit", reportBug)
	ui.ViewMap["reportBug"] = append(ui.ViewMap["reportBug"], submitInput.Name())

	cancelReportBugInput := widget.NewButtonWidgetAutoWidth("cancelReportBugInput",
		buttonX+11, y, "Cancel", cancelReportBug)
	ui.ViewMap["reportBug"] = append(ui.ViewMap["reportBug"], cancelReportBugInput.Name())

	ui.AddEventForInputWidgets(ui.ViewMap["reportBug"])

	return nil
}

func reportBug(g *gocui.Gui, v *gocui.View) error {
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
		color.New(color.FgMagenta).Fprintf(v, "Title cannot be empty")
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

	logUtils.PrintStructToCmd(bug)
	zentaoService.SubmitBug(bug, idInTask, stepIds)

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
	DestoryReportBugPage()
	return nil
}

func DestoryReportBugPage() {
	for _, v := range ui.ViewMap["reportBug"] {
		vari.Cui.DeleteView(v)
		vari.Cui.DeleteKeybindings(v)
	}
}
