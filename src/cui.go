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
	leftWidth          = 32
	labelWidth         = 15
	inputFullLineWidth = 69
	inputNumbWidth     = 25
	buttonWidth        = 10
	space              = 2
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

	if err := keybindings(g); err != nil {
		log.Panicln(err)
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}

func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()

	qickbarView := ui.NewLabelWidget(g, "qickbar", 0, 0, leftWidth, "")
	viewMap["root"] = append(viewMap["root"], qickbarView.Name())

	importView := ui.NewLabelWidget(g, "import", 3, 0, 9, "Import")
	viewMap["root"] = append(viewMap["root"], importView.Name())
	importView.Frame = false

	switchView := ui.NewLabelWidget(g, "switch", 19, 0, 13, "Switch")
	viewMap["root"] = append(viewMap["root"], switchView.Name())
	switchView.Frame = false

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

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}
func showHelp(g *gocui.Gui, v *gocui.View) error {
	help, _ := g.View("help")

	if help != nil {
		hideHelp(g)
	} else {
		ui.NewHelpWidget(g)
	}

	return nil
}
func hideHelp(g *gocui.Gui) error {
	help, _ := g.View("help")

	if help != nil {
		if err := g.DeleteView("help"); err != nil {
			return err
		}
	}

	return nil
}

func keybindings(g *gocui.Gui) error {
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		return err
	}
	if err := g.SetKeybinding("", gocui.KeyCtrlH, gocui.ModNone, showHelp); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("", gocui.KeyTab, gocui.ModNone, toggleInput); err != nil {
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
	hideHelp(g)

	maxX, _ := g.Size()

	slideView, _ := g.View("side")
	slideX, _ := slideView.Size()

	left := slideX + 2
	right := left + labelWidth
	if v, err := g.SetView("urlLabel", left, 1, right, 3); err != nil {
		v.Frame = false
		fmt.Fprint(v, "ZentaoUrl")

		viewMap["import"] = append(viewMap["import"], v.Name())
	}

	left = right + space
	right = left + inputFullLineWidth
	if v, err := g.SetView("urlInput", left, 1, right, 3); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Clear()
		fmt.Fprint(v, server.URL)

		v.Editable = true
		v.Wrap = true

		if _, err := g.SetCurrentView("urlInput"); err != nil {
			return err
		}

		viewMap["import"] = append(viewMap["import"], v.Name())
	}

	left = slideX + 2
	right = left + labelWidth
	if v, err := g.SetView("productLabel", left, 4, right, 6); err != nil {
		v.Frame = false
		fmt.Fprint(v, "ProdoctId")
		viewMap["import"] = append(viewMap["import"], v.Name())
	}

	left = right + space
	right = left + inputNumbWidth
	if v, err := g.SetView("productInput", left, 4, right, 6); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		v.Editable = true
		v.Wrap = true

		fmt.Fprint(v, "1")
		viewMap["import"] = append(viewMap["import"], v.Name())
	}

	left = right + space
	right = left + labelWidth
	if v, err := g.SetView("taskLabel", left, 4, right, 6); err != nil {
		v.Frame = false
		fmt.Fprint(v, "or TaskId")
		viewMap["import"] = append(viewMap["import"], v.Name())
	}

	left = right + space
	right = left + inputNumbWidth
	if v, err := g.SetView("taskInput", left, 4, right, 6); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		v.Editable = true
		v.Wrap = true

		fmt.Fprint(v, "1")
		viewMap["import"] = append(viewMap["import"], v.Name())
	}

	left = slideX + 2
	right = left + labelWidth
	if v, err := g.SetView("languageLabel", left, 7, right, 9); err != nil {
		v.Frame = false
		fmt.Fprint(v, "Language")
		viewMap["import"] = append(viewMap["import"], v.Name())
	}

	left = right + space
	right = left + inputNumbWidth
	if v, err := g.SetView("languageInput", left, 7, right, 9); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		v.Editable = true
		v.Wrap = true

		fmt.Fprint(v, "python")
		viewMap["import"] = append(viewMap["import"], v.Name())
	}

	left = right + space
	right = left + labelWidth
	if v, err := g.SetView("singleFileLabel", left, 7, right, 9); err != nil {
		v.Frame = false
		fmt.Fprint(v, "SingleFile")
		viewMap["import"] = append(viewMap["import"], v.Name())
	}

	left = right + space
	right = left + inputNumbWidth
	if v, err := g.SetView("singleFileInput", left, 7, right, 9); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Frame = false

		if err := g.SetKeybinding("singleFileInput", gocui.KeySpace, gocui.ModNone, changeSingleFile); err != nil {
			return err
		}

		fmt.Fprint(v, "[*]")
		viewMap["import"] = append(viewMap["import"], v.Name())
	}

	buttonX := (maxX-leftWidth)/2 + leftWidth - buttonWidth
	if v, err := g.SetView("submitInput", buttonX, 10, buttonX+buttonWidth, 12); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		v.Highlight = true
		v.BgColor = gocui.ColorGreen
		v.FgColor = gocui.ColorBlack

		fmt.Fprint(v, "  Submit  ")
		if err := g.SetKeybinding("submitInput", gocui.MouseLeft, gocui.ModNone, importProjectRequest); err != nil {
			return err
		}
		if err := g.SetKeybinding("submitInput", gocui.KeyEnter, gocui.ModNone, importProjectRequest); err != nil {
			return err
		}
		viewMap["import"] = append(viewMap["import"], v.Name())
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

func toggleInput(g *gocui.Gui, v *gocui.View) error {
	nextView := ui.GetNextView(v.Name(), viewMap["import"])

	if nextView != "" {
		_, err := g.SetCurrentView(nextView)
		return err
	}

	return nil
}
