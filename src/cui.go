package main

import (
	"github.com/easysoft/zentaoatf/src/mock"
	"github.com/easysoft/zentaoatf/src/ui"
	"github.com/easysoft/zentaoatf/src/utils"
	"github.com/jroimartin/gocui"
	"log"
)

func main() {
	mock.Server = mock.CreateServer("case-from-prodoct.json")
	defer mock.Server.Close()

	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()
	g.Cursor = true
	g.Mouse = true

	ui.InitMainPage(g)

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}

func init() {
	utils.InitPreference()
}
