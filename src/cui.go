package main

import (
	"fmt"
	"github.com/easysoft/zentaoatf/src/action"
	httpClient "github.com/easysoft/zentaoatf/src/http"
	"github.com/easysoft/zentaoatf/src/mock"
	"github.com/easysoft/zentaoatf/src/ui"
	"github.com/easysoft/zentaoatf/src/utils"
	"github.com/jroimartin/gocui"
	"log"
	"net/http/httptest"
	"strings"
	"time"
)

const (
	leftWidth = 32
)

var server *httptest.Server
var viewMap map[string][]string

func main() {
	server = mock.Server("case-from-prodoct.json")
	defer server.Close()

	viewMap = map[string][]string{"root": {}, "import": {}}

	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()
	g.Cursor = true
	g.Mouse = true

	layout(g)

	if err := keyBindings(g); err != nil {
		log.Panicln(err)
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}

func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()

	quickBarView := ui.NewPanelWidget(g, "quickBarView", 0, 0, leftWidth, 2, "")
	viewMap["root"] = append(viewMap["root"], quickBarView.Name())

	importView := ui.NewLabelWidget(g, "import", 3, 0, "Import")
	viewMap["root"] = append(viewMap["root"], importView.Name())

	switchView := ui.NewLabelWidget(g, "switch", 19, 0, "Switch")
	viewMap["root"] = append(viewMap["root"], switchView.Name())

	sideView := ui.NewPanelWidget(g, "side", 0, 2, leftWidth, maxY-3, "")
	viewMap["root"] = append(viewMap["root"], sideView.Name())

	mainView := ui.NewPanelWidget(g, "main", leftWidth, 0, maxX-1-leftWidth, maxY-10, "")
	viewMap["root"] = append(viewMap["root"], mainView.Name())

	cmdView := ui.NewPanelWidget(g, "cmd", leftWidth, maxY-10, maxX-1-leftWidth, 9, "")
	viewMap["root"] = append(viewMap["root"], cmdView.Name())

	cmdView.Editable = true
	cmdView.Wrap = true
	cmdView.Autoscroll = true

	ui.NewHelpWidget(g)

	return nil
}

func keyBindings(g *gocui.Gui) error {
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, ui.Quit); err != nil {
		return err
	}
	if err := g.SetKeybinding("", gocui.KeyCtrlH, gocui.ModNone, ui.ShowHelp); err != nil {
		log.Panicln(err)
	}

	if err := g.SetKeybinding("import", gocui.MouseLeft, gocui.ModNone, importProjectUi); err != nil {
		return err
	}
	if err := g.SetKeybinding("switch", gocui.MouseLeft, gocui.ModNone, switchProjectUi); err != nil {
		return err
	}

	if err := g.SetKeybinding("cmd", gocui.MouseLeft, gocui.ModNone, setEdit); err != nil {
		return err
	}

	return nil
}

func setEdit(g *gocui.Gui, v *gocui.View) error {
	if _, err := g.SetCurrentView("cmd"); err != nil {
		return err
	}

	v.Autoscroll = true
	v.Clear()

	return nil
}

func importProjectUi(g *gocui.Gui, v *gocui.View) error {
	ui.HideHelp(g)

	maxX, _ := g.Size()

	slideView, _ := g.View("side")
	slideX, _ := slideView.Size()

	left := slideX + 2
	right := left + ui.LabelWidth
	urlLabel := ui.NewLabelWidget(g, "urlLabel", left, 1, "ZentaoUrl")
	viewMap["import"] = append(viewMap["import"], urlLabel.Name())

	left = right + ui.Space
	right = left + ui.TextWidthFull
	urlInput := ui.NewTextWidget(g, "urlInput", left, 1, ui.TextWidthFull, server.URL)
	viewMap["import"] = append(viewMap["import"], urlInput.Name())
	if _, err := g.SetCurrentView("urlInput"); err != nil {
		return err
	}

	left = slideX + 2
	right = left + ui.LabelWidth
	productLabel := ui.NewLabelWidget(g, "productLabel", left, 4, "ProductId")
	viewMap["import"] = append(viewMap["import"], productLabel.Name())

	left = right + ui.Space
	right = left + ui.TextWidthHalf
	productInput := ui.NewTextWidget(g, "productInput", left, 4, ui.TextWidthHalf, "1")
	viewMap["import"] = append(viewMap["import"], productInput.Name())

	left = right + ui.Space
	right = left + ui.LabelWidth
	taskLabel := ui.NewLabelWidget(g, "taskLabel", left, 4, "or TaskId")
	viewMap["import"] = append(viewMap["import"], taskLabel.Name())

	left = right + ui.Space
	right = left + ui.TextWidthHalf
	taskInput := ui.NewTextWidget(g, "taskInput", left, 4, ui.TextWidthHalf, "1")
	viewMap["import"] = append(viewMap["import"], taskInput.Name())

	left = slideX + 2
	right = left + ui.LabelWidth
	languageLabel := ui.NewLabelWidget(g, "languageLabel", left, 7, "Language")
	viewMap["import"] = append(viewMap["import"], languageLabel.Name())

	left = right + ui.Space
	right = left + ui.TextWidthHalf
	languageInput := ui.NewTextWidget(g, "languageInput", left, 7, ui.TextWidthHalf, "python")
	viewMap["import"] = append(viewMap["import"], languageInput.Name())

	left = right + ui.Space
	right = left + ui.LabelWidth
	singleFileLabel := ui.NewLabelWidget(g, "singleFileLabel", left, 7, "SingleFile")
	viewMap["import"] = append(viewMap["import"], singleFileLabel.Name())

	left = right + ui.Space
	right = left + ui.TextWidthHalf
	singleFileInput := ui.NewLabelWidget(g, "singleFileInput", left, 7, "[*]")
	viewMap["import"] = append(viewMap["import"], singleFileInput.Name())
	if err := g.SetKeybinding("singleFileInput", gocui.KeySpace, gocui.ModNone, changeSingleFile); err != nil {
		return err
	}

	buttonX := (maxX-leftWidth)/2 + leftWidth - ui.ButtonWidth
	submitInput := ui.NewButtonWidgetAutoWidth(g, "submitInput", buttonX, 10, "Submit", importProjectRequest)
	viewMap["import"] = append(viewMap["import"], submitInput.Name())

	if err := g.SetKeybinding("", gocui.KeyTab, gocui.ModNone, ui.ToggleInput(viewMap["import"])); err != nil {
		log.Panicln(err)
	}

	return nil
}

func switchProjectUi(g *gocui.Gui, v *gocui.View) error {
	return nil
}

func importProjectRequest(g *gocui.Gui, v *gocui.View) error {
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
	singleFile := ui.ParseRadioVal(singleFileStr)

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

	err := action.Generate(json, language, singleFile)
	if err == nil {
		fmt.Fprintln(cmdView, fmt.Sprintf("success to generate test scripts in '%s' at %s",
			utils.GenDir, utils.DateTimeStr(time.Now())))
	} else {
		fmt.Fprintln(cmdView, err.Error())
	}

	return nil
}

func changeSingleFile(g *gocui.Gui, v *gocui.View) error {
	val := strings.TrimSpace(v.Buffer())

	v.Clear()
	if val == "[*]" {
		fmt.Fprint(v, "[ ]")
	} else {
		fmt.Fprint(v, "[*]")
	}

	return nil
}
