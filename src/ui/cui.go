package ui

import (
	"github.com/easysoft/zentaoatf/src/utils"
	"github.com/jroimartin/gocui"
	"log"
)

func Cui() {
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()
	if utils.IsWin() {
		g.ASCII = true
	}

	g.Cursor = true
	g.Mouse = true

	utils.Cui = g

	InitMainPage()

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}

func init() {

}
