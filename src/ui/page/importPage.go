package page

import (
	"fmt"
	"github.com/easysoft/zentaoatf/src/action"
	"github.com/easysoft/zentaoatf/src/http"
	"github.com/easysoft/zentaoatf/src/mock"
	"github.com/easysoft/zentaoatf/src/ui"
	"github.com/easysoft/zentaoatf/src/ui/widget"
	"github.com/easysoft/zentaoatf/src/utils"
	"github.com/jroimartin/gocui"
	"log"
	"strings"
	"time"
)

func InitImportPage(g *gocui.Gui, v *gocui.View) error {
	DestoryPages(g, v)

	maxX, _ := g.Size()
	slideView, _ := g.View("side")
	slideX, _ := slideView.Size()

	left := slideX + 2
	right := left + widget.LabelWidth
	urlLabel := widget.NewLabelWidget(g, "urlLabel", left, 1, "ZentaoUrl")
	ui.ViewMap["import"] = append(ui.ViewMap["import"], urlLabel.Name())

	left = right + ui.Space
	right = left + widget.TextWidthFull
	urlInput := widget.NewTextWidget(g, "urlInput", left, 1, widget.TextWidthFull, mock.Server.URL)
	ui.ViewMap["import"] = append(ui.ViewMap["import"], urlInput.Name())
	if _, err := g.SetCurrentView("urlInput"); err != nil {
		return err
	}

	left = slideX + 2
	right = left + widget.LabelWidth
	productLabel := widget.NewLabelWidget(g, "productLabel", left, 4, "ProductId")
	ui.ViewMap["import"] = append(ui.ViewMap["import"], productLabel.Name())

	left = right + ui.Space
	right = left + widget.TextWidthHalf
	productInput := widget.NewTextWidget(g, "productInput", left, 4, widget.TextWidthHalf, "1")
	ui.ViewMap["import"] = append(ui.ViewMap["import"], productInput.Name())

	left = right + ui.Space
	right = left + widget.LabelWidth
	taskLabel := widget.NewLabelWidget(g, "taskLabel", left, 4, "or TaskId")
	ui.ViewMap["import"] = append(ui.ViewMap["import"], taskLabel.Name())

	left = right + ui.Space
	right = left + widget.TextWidthHalf
	taskInput := widget.NewTextWidget(g, "taskInput", left, 4, widget.TextWidthHalf, "1")
	ui.ViewMap["import"] = append(ui.ViewMap["import"], taskInput.Name())

	left = slideX + 2
	right = left + widget.LabelWidth
	languageLabel := widget.NewLabelWidget(g, "languageLabel", left, 7, "Language")
	ui.ViewMap["import"] = append(ui.ViewMap["import"], languageLabel.Name())

	left = right + ui.Space
	right = left + widget.TextWidthHalf
	languageInput := widget.NewTextWidget(g, "languageInput", left, 7, widget.TextWidthHalf, "python")
	ui.ViewMap["import"] = append(ui.ViewMap["import"], languageInput.Name())

	left = right + ui.Space
	right = left + widget.LabelWidth
	singleFileLabel := widget.NewLabelWidget(g, "singleFileLabel", left, 7, "SingleFile")
	ui.ViewMap["import"] = append(ui.ViewMap["import"], singleFileLabel.Name())

	left = right + ui.Space
	right = left + widget.TextWidthHalf
	singleFileInput := widget.NewRadioWidget(g, "singleFileInput", left, 7, true)
	ui.ViewMap["import"] = append(ui.ViewMap["import"], singleFileInput.Name())

	buttonX := (maxX-ui.LeftWidth)/2 + ui.LeftWidth - widget.ButtonWidth
	submitInput := widget.NewButtonWidgetAutoWidth(g, "submitInput", buttonX, 10, "Submit", ImportRequest)
	ui.ViewMap["import"] = append(ui.ViewMap["import"], submitInput.Name())

	if err := g.SetKeybinding("", gocui.KeyTab, gocui.ModNone, ui.ToggleInput(ui.ViewMap["import"])); err != nil {
		log.Panicln(err)
	}

	ui.HideHelp(g)

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
	singleFile := widget.ParseRadioVal(singleFileStr)

	params := make(map[string]string)
	if productCode != "" {
		params["entityType"] = "product"
		params["entityVal"] = productCode
	} else {
		params["entityType"] = "task"
		params["entityVal"] = taskId
	}

	cmdView, _ := g.View("cmd")
	_, _ = fmt.Fprintln(cmdView, fmt.Sprintf("#atf gen -u %s -t %s -v %s -l %s -s %t",
		url, params["entityType"], params["entityVal"], language, singleFile))

	json, e := httpClient.Get(url, params)
	if e != nil {
		fmt.Fprintln(cmdView, e.Error())
		return nil
	}

	count, err := action.Generate(json, language, singleFile)
	if err == nil {
		fmt.Fprintln(cmdView, fmt.Sprintf("success to generate %d test scripts in '%s' at %s",
			count, utils.GenDir, utils.DateTimeStr(time.Now())))
	} else {
		fmt.Fprintln(cmdView, err.Error())
	}

	return nil
}

func DestoryImportPage(g *gocui.Gui, v *gocui.View) {
	for _, v := range ui.ViewMap["import"] {
		g.DeleteView(v)
	}
	ui.ViewMap["import"] = make([]string, 0)

	g.DeleteKeybinding("", gocui.KeyTab, gocui.ModNone)
	g.DeleteKeybindings("singleFileInput")
	g.DeleteKeybindings("submitInput")
}
