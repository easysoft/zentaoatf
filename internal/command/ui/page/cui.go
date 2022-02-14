package page

import (
	commConsts "github.com/aaronchen2k/deeptest/internal/comm/consts"
	"github.com/aaronchen2k/deeptest/internal/command"
	"github.com/awesome-gocui/gocui"
	"log"
)

func CuiReportBug(dir string, id string, actionModule *command.IndexModule) error {
	g, err := gocui.NewGui(gocui.OutputNormal, true)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()
	//if commonUtils.IsWin() {
	//	g.ASCII = true
	//}

	g.Cursor = true
	g.Mouse = true

	commConsts.Cui = g
	//
	InitMainPage()
	InitReportBugPage(dir, id, actionModule)

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}

	return nil
}
