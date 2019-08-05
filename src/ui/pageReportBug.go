package ui

import (
	"github.com/easysoft/zentaoatf/src/biz"
	"github.com/easysoft/zentaoatf/src/utils"
	"github.com/jroimartin/gocui"
)

func InitReportBugPage() error {
	DestoryReportBugPage()

	biz.GetZentaoSettings()

	maxX, maxY := utils.Cui.Size()
	x := maxX/2 - 50
	y := maxY/2 - 14

	reportBugPanel := NewPanelWidget("reportBugPanel", x, y, 100, 26, "")
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
	moduleInput := NewSelectWidget("moduleInput", left, y+4, SelectWidth, 6,
		"Module", utils.ZendaoSettings.Modules)
	ViewMap["reportBug"] = append(ViewMap["reportBug"], moduleInput.Name())

	// category
	left = right + Space
	right = left + SelectWidth
	categoryInput := NewSelectWidget("categoryInput", left, y+4, SelectWidth, 6,
		"Category", utils.ZendaoSettings.Modules)
	ViewMap["reportBug"] = append(ViewMap["reportBug"], categoryInput.Name())

	// version
	left = right + Space
	right = left + SelectWidth
	versionInput := NewSelectWidget("versionInput", left, y+4, SelectWidth, 6,
		"Version", utils.ZendaoSettings.Modules)
	ViewMap["reportBug"] = append(ViewMap["reportBug"], versionInput.Name())

	// priority
	left = x + 2 + LabelWidthSmall + Space
	priorityInput := NewSelectWidget("priorityInput", left, y+11, SelectWidth, 6,
		"Priority", utils.ZendaoSettings.Modules)
	ViewMap["reportBug"] = append(ViewMap["reportBug"], priorityInput.Name())

	// buttons
	buttonX := maxX/2 - 50 + 2 + LabelWidthSmall + Space
	submitInput := NewButtonWidgetAutoWidth("submitInput", buttonX, y+19, "Submit", reportBug)
	ViewMap["reportBug"] = append(ViewMap["reportBug"], submitInput.Name())

	cancelReportBugInput := NewButtonWidgetAutoWidth("cancelReportBugInput",
		buttonX+12, y+19, "Cancel", cancelReportBug)
	ViewMap["reportBug"] = append(ViewMap["reportBug"], cancelReportBugInput.Name())

	keyBindsInput(ViewMap["reportBug"])

	return nil
}

func reportBug(g *gocui.Gui, v *gocui.View) error {
	//titleView, _ := g.View("titleInput")
	//moduleView, _ := g.View("moduleInput")
	//categoryView, _ := g.View("categoryInput")
	//versionView, _ := g.View("versionInput")
	//priorityView, _ := g.View("priorityInput")
	//
	//title := strings.TrimSpace(titleView.Buffer())
	//moduleStr := strings.TrimSpace(moduleView.Buffer())
	//categoryStr := strings.TrimSpace(categoryView.Buffer())
	//versionStr := strings.TrimSpace(versionView.Buffer())
	//priorityStr := strings.TrimSpace(priorityView.Buffer())
	//
	//params := make(map[string]string)
	//if productCode != "" {
	//	params["entityType"] = "product"
	//	params["entityVal"] = productCode
	//} else {
	//	params["entityType"] = "task"
	//	params["entityVal"] = taskId
	//}
	//
	//jsonStr, _ := json.Marshal(params)
	//url := utils.UpdateUrl(mock.BaseUrl)
	//
	//json, e := httpClient.Post(url+utils.UrlImportProject, string(jsonStr))
	//if e != nil {
	//	utils.PrintToCmd(e.Error())
	//	return nil
	//} else {
	//	if json.Code == 1 {
	//		utils.PrintToCmd(fmt.Sprintf("success to report bug at %s", utils.DateTimeStr(time.Now())))
	//	}
	//}

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
