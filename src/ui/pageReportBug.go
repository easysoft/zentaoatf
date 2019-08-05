package ui

import (
	"encoding/json"
	"fmt"
	"github.com/easysoft/zentaoatf/src/action"
	httpClient "github.com/easysoft/zentaoatf/src/http"
	"github.com/easysoft/zentaoatf/src/utils"
	"github.com/jroimartin/gocui"
	"strings"
	"time"
)

func InitReportBugPage() error {
	DestoryReportBugPage()

	maxX, maxY := utils.Cui.Size()
	x := maxX/2 - 50
	y := maxY/2 - 14

	reportBugPanel := NewPanelWidget("reportBugPanel", x, y, 100, 25, "")
	ViewMap["reportBug"] = append(ViewMap["reportBug"], reportBugPanel.Name())

	y += 1
	reportBugTitle := NewLabelWidgetAutoWidth("reportBugTitle", x+2+LabelWidth+Space, y, "Report Bug")
	ViewMap["reportBug"] = append(ViewMap["reportBug"], reportBugTitle.Name())

	// title
	y += 2
	left := x + 2
	right := left + LabelWidth
	titleLabel := NewLabelWidget("titleLabel", left, y+1, "Title")
	ViewMap["reportBug"] = append(ViewMap["reportBug"], titleLabel.Name())

	left = right + Space
	right = left + TextWidthFull
	titleInput := NewTextWidget("titleInput", left, y+1, TextWidthFull, "")
	ViewMap["reportBug"] = append(ViewMap["reportBug"], titleInput.Name())

	// module
	left = x + 2
	right = left + LabelWidth
	moduleLabel := NewLabelWidget("moduleLabel", left, y+4, "Module")
	ViewMap["reportBug"] = append(ViewMap["reportBug"], moduleLabel.Name())

	left = right + Space
	right = left + TextWidthHalf
	moduleInput := NewTextWidget("moduleInput", left, y+4, TextWidthHalf, "")
	ViewMap["reportBug"] = append(ViewMap["reportBug"], moduleInput.Name())

	// category
	left = right + Space
	right = left + LabelWidth
	categoryLabel := NewLabelWidget("categoryLabel", left, y+4, "Category")
	ViewMap["reportBug"] = append(ViewMap["reportBug"], categoryLabel.Name())

	left = right + Space
	right = left + TextWidthHalf
	categoryInput := NewTextWidget("categoryInput", left, y+4, TextWidthHalf, "")
	ViewMap["reportBug"] = append(ViewMap["reportBug"], categoryInput.Name())

	// version
	left = x + 2
	right = left + LabelWidth
	versionLabel := NewLabelWidget("versionLabel", left, y+7, "Version")
	ViewMap["reportBug"] = append(ViewMap["reportBug"], versionLabel.Name())

	left = right + Space
	right = left + TextWidthHalf
	versionInput := NewTextWidget("versionInput", left, y+7, TextWidthHalf, "")
	ViewMap["reportBug"] = append(ViewMap["reportBug"], versionInput.Name())

	// priority
	left = right + Space
	right = left + LabelWidth
	priorityLabel := NewLabelWidget("priorityLabel", left, y+7, "Priority")
	ViewMap["reportBug"] = append(ViewMap["reportBug"], priorityLabel.Name())

	left = right + Space
	right = left + TextWidthHalf
	priorityInput := NewTextWidget("priorityInput", left, y+7, TextWidthHalf, "")
	ViewMap["reportBug"] = append(ViewMap["reportBug"], priorityInput.Name())

	// buttons
	buttonX := maxX/2 - 50 + 2 + LabelWidth + Space
	submitInput := NewButtonWidgetAutoWidth("submitInput", buttonX, y+10, "Submit", reportBug)
	ViewMap["reportBug"] = append(ViewMap["reportBug"], submitInput.Name())

	cancelReportBugInput := NewButtonWidgetAutoWidth("cancelReportBugInput",
		buttonX+12, y+10, "Cancel", cancelReportBug)
	ViewMap["reportBug"] = append(ViewMap["reportBug"], cancelReportBugInput.Name())

	keyBindsInput(ViewMap["reportBug"])

	return nil
}

func reportBug(g *gocui.Gui, v *gocui.View) error {
	urlView, _ := g.View("urlInput")
	productView, _ := g.View("productInput")
	taskView, _ := g.View("taskInput")
	languageView, _ := g.View("languageInput")
	singleFileView, _ := g.View("singleFileInput")

	url := strings.TrimSpace(urlView.ViewBuffer())

	productCode := strings.TrimSpace(productView.Buffer())
	taskId := strings.TrimSpace(taskView.Buffer())
	language := strings.TrimSpace(languageView.Buffer())
	singleFileStr := strings.TrimSpace(singleFileView.Buffer())
	singleFile := ParseRadioVal(singleFileStr)

	params := make(map[string]string)
	if productCode != "" {
		params["entityType"] = "product"
		params["entityVal"] = productCode
	} else {
		params["entityType"] = "task"
		params["entityVal"] = taskId
	}

	jsonStr, _ := json.Marshal(params)
	url = utils.UpdateUrl(url)
	utils.PrintToCmd(fmt.Sprintf("#atf gen -u %s -t %s -v %s -l %s -s %t",
		url, params["entityType"], params["entityVal"], language, singleFile))

	json, e := httpClient.Post(url+utils.UrlImportProject, string(jsonStr))
	if e != nil {
		utils.PrintToCmd(e.Error())
		return nil
	}

	count, err := action.Generate(json, url, params["entityType"], params["entityVal"], language, singleFile)
	if err == nil {
		utils.PrintToCmd(fmt.Sprintf("success to generate %d test scripts in '%s' at %s",
			count, utils.ScriptDir, utils.DateTimeStr(time.Now())))
	} else {
		utils.PrintToCmd(err.Error())
	}

	return nil
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
