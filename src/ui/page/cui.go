package page

import (
	"github.com/easysoft/zentaoatf/src/utils/common"
	"github.com/easysoft/zentaoatf/src/utils/vari"
	"github.com/jroimartin/gocui"
	"log"
)

func Cui(page string) {
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()
	if commonUtils.IsWin() {
		g.ASCII = true
	}

	g.Cursor = true
	g.Mouse = true

	vari.Cui = g

	if page == "bug" {
		InitReportBugPage()
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}

func init() {

}
