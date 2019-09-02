package page

import (
	commonUtils "github.com/easysoft/zentaoatf/src/utils/common"
	"github.com/easysoft/zentaoatf/src/utils/vari"
	"github.com/jroimartin/gocui"
	"log"
)

func CuiReportBug(dir string, id string) error {
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
	vari.RunFromCui = true

	//InitMainPage()
	InitReportBugPage(dir, id)

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}

	return nil
}

func init() {

}
