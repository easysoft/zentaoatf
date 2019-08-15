package ui

import (
	"github.com/easysoft/zentaoatf/src/model"
	zentaoService "github.com/easysoft/zentaoatf/src/service/zentao"
	"github.com/easysoft/zentaoatf/src/utils/vari"
	zentaoUtils "github.com/easysoft/zentaoatf/src/utils/zentao"
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

	reportBugPanel := NewPanelWidget("reportBugPanel", x, y, 100, 27, "")
	ViewMap["reportBug"] = append(ViewMap["reportBug"], reportBugPanel.Name())

	y += 1
	reportBugTitle := NewLabelWidgetAutoWidth("reportBugTitle", x+2+LabelWidthSmall+Space, y, "Report Bug")
	ViewMap["reportBug"] = append(ViewMap["reportBug"], reportBugTitle.Name())

	// title
	y += 2
	left := x + 2
	right := left + LabelWidthSmall
	titleLabel := NewLabelWidget("titleLabel", left, y+1, "Desc")
	ViewMap["reportBug"] = append(ViewMap["reportBug"], titleLabel.Name())

	left = right + Space
	right = left + TextWidthFull
	titleInput := NewTextWidget("titleInput", left, y+1, TextWidthFull, bug.Title)
	ViewMap["reportBug"] = append(ViewMap["reportBug"], titleInput.Name())

	// module
	left = x + 2 + LabelWidthSmall + Space
	right = left + SelectWidth
	moduleInput := NewSelectWidgetWithDefault("module", left, y+4, SelectWidth, 6, "Module",
		vari.ZendaoSettings.Modules, zentaoService.GetNameById(bug.Module, vari.ZendaoSettings.Modules),
		bugSelectFieldCheckEvent(filedValMap))
	ViewMap["reportBug"] = append(ViewMap["reportBug"], moduleInput.Name())

	// category
	left = right + Space
	right = left + SelectWidth
	categoryInput := NewSelectWidgetWithDefault("category", left, y+4, SelectWidth, 6, "Category",
		vari.ZendaoSettings.Modules, bug.Type,
		bugSelectFieldCheckEvent(filedValMap))
	ViewMap["reportBug"] = append(ViewMap["reportBug"], categoryInput.Name())

	// version
	left = right + Space
	right = left + SelectWidth
	versionInput := NewSelectWidgetWithDefault("version", left, y+4, SelectWidth, 6, "Version",
		vari.ZendaoSettings.Modules, zentaoService.GetNameById(bugVersion, vari.ZendaoSettings.Versions),
		bugSelectFieldCheckEvent(filedValMap))
	ViewMap["reportBug"] = append(ViewMap["reportBug"], versionInput.Name())

	// severity
	left = x + 2 + LabelWidthSmall + Space
	severityInput := NewSelectWidgetWithDefault("severity", left, y+11, SelectWidth, 6, "Severity",
		vari.ZendaoSettings.Modules, zentaoService.GetNameById(bug.Severity, vari.ZendaoSettings.Severities),
		bugSelectFieldCheckEvent(filedValMap))
	ViewMap["reportBug"] = append(ViewMap["reportBug"], severityInput.Name())

	// priority
	left = right + Space
	right = left + SelectWidth
	priorityInput := NewSelectWidgetWithDefault("priority", left, y+11, SelectWidth, 6, "Priority",
		vari.ZendaoSettings.Modules, zentaoService.GetNameById(bug.Pri, vari.ZendaoSettings.Priorities),
		bugSelectFieldCheckEvent(filedValMap))
	ViewMap["reportBug"] = append(ViewMap["reportBug"], priorityInput.Name())

	// msg
	reportBugMsg := NewPanelWidget("reportBugMsg", x+2+LabelWidthSmall+Space, y+18, TextWidthFull, 2, "")
	reportBugMsg.Frame = false
	ViewMap["reportBug"] = append(ViewMap["reportBug"], reportBugMsg.Name())

	// buttons
	buttonX := maxX/2 - 50 + 2 + LabelWidthSmall + Space
	submitInput := NewButtonWidgetAutoWidth("submitInput", buttonX, y+21, "Submit", reportBug)
	ViewMap["reportBug"] = append(ViewMap["reportBug"], submitInput.Name())

	cancelReportBugInput := NewButtonWidgetAutoWidth("cancelReportBugInput",
		buttonX+12, y+21, "Cancel", cancelReportBug)
	ViewMap["reportBug"] = append(ViewMap["reportBug"], cancelReportBugInput.Name())

	keyBindsInput(ViewMap["reportBug"])

	return nil
}

func reportBug(g *gocui.Gui, v *gocui.View) error {
	titleView, _ := g.View("titleInput")
	moduleView, _ := g.View("module")
	categoryView, _ := g.View("category")
	versionView, _ := g.View("version")
	severityView, _ := g.View("severity")
	priorityView, _ := g.View("priority")

	title := strings.TrimSpace(titleView.Buffer())
	moduleStr := strings.TrimSpace(GetSelectedLineVal(moduleView))
	typeStr := strings.TrimSpace(GetSelectedLineVal(categoryView))
	versionStr := strings.TrimSpace(GetSelectedLineVal(versionView))
	severityStr := strings.TrimSpace(GetSelectedLineVal(severityView))
	priorityStr := strings.TrimSpace(GetSelectedLineVal(priorityView))

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
	for _, v := range ViewMap["reportBug"] {
		vari.Cui.DeleteView(v)
		vari.Cui.DeleteKeybindings(v)
	}
}
