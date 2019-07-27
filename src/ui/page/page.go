package page

import (
	"github.com/easysoft/zentaoatf/src/ui"
	"github.com/jroimartin/gocui"
)

func DestoryPages(g *gocui.Gui, v *gocui.View) {
	for key, _ := range ui.ViewMap {
		if key == "import" {
			DestoryImportPage(g, v)
		} else if key == "switch" {
			DestorySwitchPage(g, v)
		}
	}
}
