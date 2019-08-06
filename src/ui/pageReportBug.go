package ui

import (
	"encoding/json"
	"fmt"
	"github.com/easysoft/zentaoatf/src/biz"
	httpClient "github.com/easysoft/zentaoatf/src/http"
	"github.com/easysoft/zentaoatf/src/mock"
	"github.com/easysoft/zentaoatf/src/utils"
	"github.com/fatih/color"
	"github.com/jroimartin/gocui"
	"strconv"
	"strings"
	"time"
)

var filedValMap map[string]int

func InitReportBugPage() error {
	DestoryReportBugPage()

	biz.GetZentaoSettings()

	maxX, maxY := utils.Cui.Size()
	x := maxX/2 - 50
	y := maxY/2 - 14

	reportBugPanel := NewPanelWidget("reportBugPanel", x, y, 100, 27, "")
	ViewMap["reportBug"] = append(ViewMap["reportBug"], reportBugPanel.Name())

	y += 1
	reportBugTitle := NewLabelWidgetAutoWidth("reportBugTitle", x+2+LabelWidthSmall+Space, y, "Report Bug")
	ViewMap["reportBug"] = append(ViewMap["reportBug"], reportBugTitle.Name())

	// title
	y += 2
	left := x + 2
	right := left + LabelWidthSmall
	titleLabel := NewLabelWidget("titleLabel", left, y+1, "Title")
	ViewMap["reportBug"] = append(ViewMap["reportBug"], titleLabel.Name())

	left = right + Space
	right = left + TextWidthFull
	titleInput := NewTextWidget("titleInput", left, y+1, TextWidthFull, "")
	ViewMap["reportBug"] = append(ViewMap["reportBug"], titleInput.Name())

	// module
	left = x + 2 + LabelWidthSmall + Space
	right = left + SelectWidth
	moduleInput := NewSelectWidget("module", left, y+4, SelectWidth, 6, "Module", utils.ZendaoSettings.Modules,
		bugSelectFieldCheckEvent(filedValMap))
	ViewMap["reportBug"] = append(ViewMap["reportBug"], moduleInput.Name())

	// category
	left = right + Space
	right = left + SelectWidth
	categoryInput := NewSelectWidget("category", left, y+4, SelectWidth, 6, "Category", utils.ZendaoSettings.Modules,
		bugSelectFieldCheckEvent(filedValMap))
	ViewMap["reportBug"] = append(ViewMap["reportBug"], categoryInput.Name())

	// version
	left = right + Space
	right = left + SelectWidth
	versionInput := NewSelectWidget("version", left, y+4, SelectWidth, 6, "Version", utils.ZendaoSettings.Modules,
		bugSelectFieldCheckEvent(filedValMap))
	ViewMap["reportBug"] = append(ViewMap["reportBug"], versionInput.Name())

	// priority
	left = x + 2 + LabelWidthSmall + Space
	priorityInput := NewSelectWidget("priority", left, y+11, SelectWidth, 6, "Priority", utils.ZendaoSettings.Modules,
		bugSelectFieldCheckEvent(filedValMap))
	ViewMap["reportBug"] = append(ViewMap["reportBug"], priorityInput.Name())

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
	priorityView, _ := g.View("priority")

	title := strings.TrimSpace(titleView.Buffer())
	moduleStr := strings.TrimSpace(GetSelectedLineVal(moduleView))
	categoryStr := strings.TrimSpace(GetSelectedLineVal(categoryView))
	versionStr := strings.TrimSpace(GetSelectedLineVal(versionView))
	priorityStr := strings.TrimSpace(GetSelectedLineVal(priorityView))

	if title == "" {
		v, _ := utils.Cui.View("reportBugMsg")

		color.New(color.FgMagenta).Fprintf(v, "Title cannot be empty")
		return nil
	}

	config := utils.ReadCurrConfig()
	params := make(map[string]interface{})
	params["entityType"] = config.EntityType
	params["entityVal"] = config.EntityVal
	params["projectName"] = config.ProjectName

	params["title"] = title

	modulelId, _ := strconv.Atoi(moduleStr)
	params["moduleId"] = modulelId

	categoryId, _ := strconv.Atoi(categoryStr)
	params["categoryId"] = categoryId

	versionId, _ := strconv.Atoi(versionStr)
	params["versionId"] = versionId

	priorityId, _ := strconv.Atoi(priorityStr)
	params["priorityId"] = priorityId

	jsonStr, _ := json.Marshal(params)
	url := utils.UpdateUrl(mock.BaseUrl)

	json, e := httpClient.Post(url+utils.UrlReportBug, string(jsonStr))
	if e != nil {
		utils.PrintToCmd(e.Error())
		return nil
	} else {
		if json.Code == 1 {
			utils.PrintToCmd(fmt.Sprintf("success to report bug at %s", utils.DateTimeStr(time.Now())))
		}
	}

	return nil
}

func bugSelectFieldCheckEvent(filedValMap map[string]int) func(g *gocui.Gui, v *gocui.View) error {
	return func(g *gocui.Gui, v *gocui.View) error {
		name := v.Name()

		g.SetCurrentView(name)

		line, _ := GetSelectedLine(v, ".*")
		line = strings.TrimSpace(line)

		utils.SetBugField(name, line, filedValMap)

		return nil
	}
}

func bugSelectFieldScrollEvent(dy int, filedValMap map[string]int) func(g *gocui.Gui, v *gocui.View) error {
	return func(g *gocui.Gui, v *gocui.View) error {
		scrollAction(v, dy, true)

		name := v.Name()
		line, _ := GetSelectedLine(v, ".*")
		utils.SetBugField(name, strings.TrimSpace(line), filedValMap)

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
		utils.Cui.DeleteView(v)
		utils.Cui.DeleteKeybindings(v)
	}
}
