package page

import (
	"github.com/easysoft/zentaoatf/src/model"
	zentaoService "github.com/easysoft/zentaoatf/src/service/zentao"
	"github.com/easysoft/zentaoatf/src/ui"
	"github.com/easysoft/zentaoatf/src/ui/widget"
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

	zentaoService.GetZentaoSettings()
	bug, idInTask, stepIds = zentaoService.GenBug()

	maxX, maxY := vari.Cui.Size()
	x := maxX/2 - 50
	y := maxY/2 - 14

	var bugVersion string
	for _, val := range bug.OpenedBuild {
		bugVersion = val.(string)
	}

	reportBugPanel := widget.NewPanelWidget("reportBugPanel", x, y, 100, 27, "")
	ui.ViewMap["reportBug"] = append(ui.ViewMap["reportBug"], reportBugPanel.Name())

	y += 1
	reportBugTitle := widget.NewLabelWidgetAutoWidth("reportBugTitle", x+2+widget.LabelWidthSmall+ui.Space, y, "Report Bug")
	ui.ViewMap["reportBug"] = append(ui.ViewMap["reportBug"], reportBugTitle.Name())

	// title
	y += 3
	left := x + 2
	right := left + widget.LabelWidthSmall
	titleLabel := widget.NewLabelWidget("titleLabel", left, y, "Desc")
	ui.ViewMap["reportBug"] = append(ui.ViewMap["reportBug"], titleLabel.Name())

	left = right + ui.Space
	right = left + widget.TextWidthFull
	titleInput := widget.NewTextWidget("titleInput", left, y, widget.TextWidthFull, bug.Title)
	ui.ViewMap["reportBug"] = append(ui.ViewMap["reportBug"], titleInput.Name())

	// module
	y += 3
	left = x + 2 + widget.LabelWidthSmall + ui.Space
	right = left + widget.SelectWidth
	moduleInput := widget.NewSelectWidgetWithDefault("module", left, y, widget.SelectWidth, 6, "Module",
		vari.ZendaoSettings.Modules, zentaoService.GetNameById(bug.Module, vari.ZendaoSettings.Modules),
		bugSelectFieldCheckEvent(filedValMap))
	ui.ViewMap["reportBug"] = append(ui.ViewMap["reportBug"], moduleInput.Name())

	// type
	left = right + ui.Space
	right = left + widget.SelectWidth
	typeInput := widget.NewSelectWidgetWithDefault("type", left, y, widget.SelectWidth, 6, "Category",
		vari.ZendaoSettings.Modules, bug.Type,
		bugSelectFieldCheckEvent(filedValMap))
	ui.ViewMap["reportBug"] = append(ui.ViewMap["reportBug"], typeInput.Name())

	// version
	left = right + ui.Space
	right = left + widget.SelectWidth
	versionInput := widget.NewSelectWidgetWithDefault("version", left, y, widget.SelectWidth, 6, "Version",
		vari.ZendaoSettings.Modules, zentaoService.GetNameById(bugVersion, vari.ZendaoSettings.Versions),
		bugSelectFieldCheckEvent(filedValMap))
	ui.ViewMap["reportBug"] = append(ui.ViewMap["reportBug"], versionInput.Name())

	// severity
	y += 7
	left = x + 2 + widget.LabelWidthSmall + ui.Space
	right = left + widget.SelectWidth
	severityInput := widget.NewSelectWidgetWithDefault("severity", left, y, widget.SelectWidth, 6, "Severity",
		vari.ZendaoSettings.Modules, zentaoService.GetNameById(bug.Severity, vari.ZendaoSettings.Severities),
		bugSelectFieldCheckEvent(filedValMap))
	ui.ViewMap["reportBug"] = append(ui.ViewMap["reportBug"], severityInput.Name())

	// priority
	left = right + ui.Space
	right = left + widget.SelectWidth
	priorityInput := widget.NewSelectWidgetWithDefault("priority", left, y, widget.SelectWidth, 6, "Priority",
		vari.ZendaoSettings.Modules, zentaoService.GetNameById(bug.Pri, vari.ZendaoSettings.Priorities),
		bugSelectFieldCheckEvent(filedValMap))
	ui.ViewMap["reportBug"] = append(ui.ViewMap["reportBug"], priorityInput.Name())

	y += 7
	// msg
	reportBugMsg := widget.NewPanelWidget("reportBugMsg", x+2+widget.LabelWidthSmall+ui.Space, y, widget.TextWidthFull, 2, "")
	reportBugMsg.Frame = false
	ui.ViewMap["reportBug"] = append(ui.ViewMap["reportBug"], reportBugMsg.Name())

	// buttons
	y += 3
	buttonX := maxX/2 - 50 + 2 + widget.LabelWidthSmall + ui.Space
	submitInput := widget.NewButtonWidgetAutoWidth("submitInput", buttonX, y, "Submit", reportBug)
	ui.ViewMap["reportBug"] = append(ui.ViewMap["reportBug"], submitInput.Name())

	cancelReportBugInput := widget.NewButtonWidgetAutoWidth("cancelReportBugInput",
		buttonX+12, y, "Cancel", cancelReportBug)
	ui.ViewMap["reportBug"] = append(ui.ViewMap["reportBug"], cancelReportBugInput.Name())

	ui.KeyBindsInput(ui.ViewMap["reportBug"])

	return nil
}

func reportBug(g *gocui.Gui, v *gocui.View) error {
	titleView, _ := g.View("titleInput")
	moduleView, _ := g.View("module")
	typeView, _ := g.View("type")
	versionView, _ := g.View("version")
	severityView, _ := g.View("severity")
	priorityView, _ := g.View("priority")

	title := strings.TrimSpace(titleView.Buffer())
	moduleStr := strings.TrimSpace(ui.GetSelectedLineVal(moduleView))
	typeStr := strings.TrimSpace(ui.GetSelectedLineVal(typeView))
	versionStr := strings.TrimSpace(ui.GetSelectedLineVal(versionView))
	severityStr := strings.TrimSpace(ui.GetSelectedLineVal(severityView))
	priorityStr := strings.TrimSpace(ui.GetSelectedLineVal(priorityView))

	if title == "" {
		v, _ := vari.Cui.View("reportBugMsg")
		color.New(color.FgMagenta).Fprintf(v, "Desc cannot be empty")
		return nil
	}

	bug.Title = title
	bug.Type = typeStr

	bug.Module = zentaoService.GetIdByName(moduleStr, vari.ZendaoSettings.Modules)
	bug.OpenedBuild = map[string]interface{}{
		zentaoService.GetIdByName(versionStr, vari.ZendaoSettings.Versions): versionStr,
	}
	bug.Severity = zentaoService.GetIdByName(severityStr, vari.ZendaoSettings.Severities)
	bug.Pri = zentaoService.GetIdByName(priorityStr, vari.ZendaoSettings.Priorities)

	zentaoService.SubmitBug(bug, idInTask, stepIds)

	return nil
}

func bugSelectFieldCheckEvent(filedValMap map[string]int) func(g *gocui.Gui, v *gocui.View) error {
	return func(g *gocui.Gui, v *gocui.View) error {
		name := v.Name()

		g.SetCurrentView(name)

		//line, _ := GetSelectedLine(v, ".*")
		//line = strings.TrimSpace(line)
		//
		//zentaoUtils.SetBugField(name, line, filedValMap)

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
