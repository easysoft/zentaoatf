package main

import (
	"github.com/easysoft/zentaoatf/src/mock"
	"github.com/easysoft/zentaoatf/src/ui"
	"github.com/easysoft/zentaoatf/src/ui/page"
	"github.com/easysoft/zentaoatf/src/utils"
	"github.com/jroimartin/gocui"
	"log"
)

func main() {
	mock.Server = mock.CreateServer("case-from-prodoct.json")
	defer mock.Server.Close()

	ui.ViewMap = map[string][]string{"root": {}, "import": {}}
	ui.EventMap = map[string][][]interface{}{"root": make([][]interface{}, 0, 2), "import": make([][]interface{}, 0, 2)}

	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()
	g.Cursor = true
	g.Mouse = true

	page.InitMainPage(g)

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}

func init() {
	utils.InitConfig()
}
