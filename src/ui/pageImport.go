package ui

import (
	"fmt"
	"github.com/easysoft/zentaoatf/src/action"
	"github.com/easysoft/zentaoatf/src/http"
	"github.com/easysoft/zentaoatf/src/mock"
	"github.com/easysoft/zentaoatf/src/utils"
	"github.com/jroimartin/gocui"
	"strings"
	"time"
)

func InitImportPage() error {
	DestoryRightPages()

	maxX, _ := utils.Cui.Size()
	slideView, _ := utils.Cui.View("side")
	slideX, _ := slideView.Size()

	left := slideX + 2
	right := left + LabelWidth
	urlLabel := NewLabelWidget("urlLabel", left, 1, "ZentaoUrl")
	ViewMap["import"] = append(ViewMap["import"], urlLabel.Name())

	left = right + Space
	right = left + TextWidthFull
	urlInput := NewTextWidget("urlInput", left, 1, TextWidthFull, mock.GetUrl("importProject"))
	ViewMap["import"] = append(ViewMap["import"], urlInput.Name())

	left = slideX + 2
	right = left + LabelWidth
	productLabel := NewLabelWidget("productLabel", left, 4, "ProductId")
	ViewMap["import"] = append(ViewMap["import"], productLabel.Name())

	left = right + Space
	right = left + TextWidthHalf
	productInput := NewTextWidget("productInput", left, 4, TextWidthHalf, "1")
	ViewMap["import"] = append(ViewMap["import"], productInput.Name())

	left = right + Space
	right = left + LabelWidth
	taskLabel := NewLabelWidget("taskLabel", left, 4, "or TaskId")
	ViewMap["import"] = append(ViewMap["import"], taskLabel.Name())

	left = right + Space
	right = left + TextWidthHalf
	taskInput := NewTextWidget("taskInput", left, 4, TextWidthHalf, "1")
	ViewMap["import"] = append(ViewMap["import"], taskInput.Name())

	left = slideX + 2
	right = left + LabelWidth
	languageLabel := NewLabelWidget("languageLabel", left, 7, "Language")
	ViewMap["import"] = append(ViewMap["import"], languageLabel.Name())

	left = right + Space
	right = left + TextWidthHalf
	languageInput := NewTextWidget("languageInput", left, 7, TextWidthHalf, "python")
	ViewMap["import"] = append(ViewMap["import"], languageInput.Name())

	left = right + Space
	right = left + LabelWidth
	singleFileLabel := NewLabelWidget("singleFileLabel", left, 7, "SingleFile")
	ViewMap["import"] = append(ViewMap["import"], singleFileLabel.Name())

	left = right + Space
	right = left + TextWidthHalf
	singleFileInput := NewRadioWidget("singleFileInput", left, 7, true)
	ViewMap["import"] = append(ViewMap["import"], singleFileInput.Name())

	buttonX := (maxX-utils.LeftWidth)/2 + utils.LeftWidth - ButtonWidth
	submitInput := NewButtonWidgetAutoWidth("submitInput", buttonX, 10, "Submit", ImportRequest)
	ViewMap["import"] = append(ViewMap["import"], submitInput.Name())

	keyBindsInput(ViewMap["import"])

	return nil
}

func ImportRequest(g *gocui.Gui, v *gocui.View) error {
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

	utils.PrintToCmd(fmt.Sprintf("#atf gen -u %s -t %s -v %s -l %s -s %t",
		url, params["entityType"], params["entityVal"], language, singleFile))

	json, e := httpClient.Get(url, params)
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

func DestoryImportPage() {
	for _, v := range ViewMap["import"] {
		utils.Cui.DeleteView(v)
		utils.Cui.DeleteKeybindings(v)
	}

	utils.Cui.DeleteKeybinding("", gocui.KeyTab, gocui.ModNone)
}
