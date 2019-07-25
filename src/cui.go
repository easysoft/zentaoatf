package main

import (
	"fmt"
	httpClient "github.com/easysoft/zentaoatf/src/http"
	"github.com/easysoft/zentaoatf/src/mock"
	"github.com/jroimartin/gocui"
	"log"
	"net/http/httptest"
	"strings"
)

const (
	leftWidth          = 32
	labelWidth         = 10
	inputFullLineWidth = 66
	inputNumbWidth     = 25
	buttonWidth        = 10
	space              = 2
)

var server *httptest.Server

func main() {
	server = mock.Server("case-from-prodoct.json")
	defer server.Close()

	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	g.Cursor = true
	g.Mouse = true
	g.SetManagerFunc(layout)

	if err := keybindings(g); err != nil {
		log.Panicln(err)
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}

func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()

	if v, err := g.SetView("qickbar", 0, 0, leftWidth, 2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Highlight = true
		v.SelBgColor = gocui.ColorGreen
		v.SelFgColor = gocui.ColorBlack
	}
	if v, err := g.SetView("import", 3, 0, 14, 2); err != nil {
		v.Frame = false
		fmt.Fprintln(v, "  Import   ")
	}
	if v, err := g.SetView("switch", 19, 0, 31, 2); err != nil {
		v.Frame = false
		fmt.Fprintln(v, "  Switch   ")
	}

	if v, err := g.SetView("side", 0, 2, leftWidth, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Highlight = true
		v.SelBgColor = gocui.ColorGreen
		v.SelFgColor = gocui.ColorBlack
	}

	if v, err := g.SetView("main", leftWidth, 0, maxX-1, maxY-5); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Autoscroll = true
	}

	if v, err := g.SetView("cmd", leftWidth, maxY-5, maxX-1, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		v.Editable = true
		v.Wrap = true
	}
	return nil

}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func keybindings(g *gocui.Gui) error {
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		return err
	}

	if err := g.SetKeybinding("import", gocui.MouseLeft, gocui.ModNone, importProjectUi); err != nil {
		return err
	}
	if err := g.SetKeybinding("switch", gocui.MouseLeft, gocui.ModNone, switchProjectUi); err != nil {
		return err
	}

	if err := g.SetKeybinding("cmdline", gocui.MouseLeft, gocui.ModNone, setEdit); err != nil {
		return err
	}

	return nil
}

func setEdit(g *gocui.Gui, v *gocui.View) error {
	if _, err := g.SetCurrentView("cmdline"); err != nil {
		return err
	}

	v.SetCursor(0, 0)
	v.Clear()

	return nil
}

func importProjectUi(g *gocui.Gui, v *gocui.View) error {
	maxX, _ := g.Size()

	slideView, _ := g.View("side")
	slideX, _ := slideView.Size()

	mainView, _ := g.View("main")
	_, mainY := mainView.Size()

	left := slideX + 2
	right := left + labelWidth
	if v, err := g.SetView("urlLabel", left, 1, right, 3); err != nil {
		v.Frame = false
		fmt.Fprintln(v, "ZentaoUrl")
	}

	left = right + space
	right = left + inputFullLineWidth
	if v, err := g.SetView("urlInput", left, 1, right, 3); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		fmt.Fprintln(v, server.URL)

		v.Editable = true
		v.Wrap = true
		if _, err := g.SetCurrentView("urlInput"); err != nil {
			return err
		}
	}

	left = slideX + 2
	right = left + labelWidth
	if v, err := g.SetView("productLabel", left, 4, right, 6); err != nil {
		v.Frame = false
		fmt.Fprintln(v, "ProdoctId")
	}

	left = right + space
	right = left + inputNumbWidth
	if v, err := g.SetView("productInput", left, 4, right, 6); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		v.Editable = true
		v.Wrap = true

		fmt.Fprintln(v, "1")
	}

	left = right + space
	right = left + 3
	if v, err := g.SetView("or", left, 4, right, 6); err != nil {
		v.Frame = false
		fmt.Fprintln(v, "or")
	}

	left = right + space
	right = left + (labelWidth - 3)
	if v, err := g.SetView("taskLabel", left, 4, right, 6); err != nil {
		v.Frame = false
		fmt.Fprintln(v, "TaskId")
	}

	left = right + space
	right = left + inputNumbWidth
	if v, err := g.SetView("taskInput", left, 4, right, 6); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		v.Editable = true
		v.Wrap = true

		fmt.Fprintln(v, "1")
	}

	buttonX := (maxX-leftWidth)/2 + leftWidth - buttonWidth
	buttonY := mainY - 2
	if v, err := g.SetView("submit", buttonX, buttonY, buttonX+buttonWidth, buttonY+2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		fmt.Fprintln(v, "  Submit  ")

		if err := g.SetKeybinding("submit", gocui.MouseLeft, gocui.ModNone, importProjectRequest); err != nil {
			return err
		}
	}

	return nil
}

func switchProjectUi(g *gocui.Gui, v *gocui.View) error {
	return nil
}

func importProjectRequest(g *gocui.Gui, v *gocui.View) error {
	urlView, err := g.View("urlInput")
	if err != nil {
		fmt.Println(err.Error())
	}
	productView, err := g.View("productInput")
	if err != nil {
		fmt.Println(err.Error())
	}
	taskView, err := g.View("taskInput")
	if err != nil {
		fmt.Println(err.Error())
	}

	url := strings.TrimSpace(urlView.ViewBuffer())
	productId := strings.TrimSpace(productView.ViewBuffer())
	taskId := strings.TrimSpace(taskView.ViewBuffer())

	fmt.Println(url)

	params := make(map[string]string)
	params["productId"] = productId
	params["taskId"] = taskId

	jsonStr := httpClient.Get(url, params)

	cmdView, _ := g.View("cmd")
	fmt.Fprintln(cmdView, jsonStr)

	return nil
}
