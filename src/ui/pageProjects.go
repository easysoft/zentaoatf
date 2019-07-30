package ui

import (
	"github.com/jroimartin/gocui"
)

var CurrProjectsButton string
var projectsButtons []string

func InitProjectsPage(g *gocui.Gui) error {
	importLabel := NewLabelWidgetAutoWidth(g, "switch", 0, 2, "Switch Work dir")
	ViewMap["projects"] = append(ViewMap["projects"], importLabel.Name())

	switchLabel := NewLabelWidgetAutoWidth(g, "import", 0, 3, "Import from Zentao")
	ViewMap["projects"] = append(ViewMap["projects"], switchLabel.Name())

	keybindingProjectsButton(g)

	return nil
}

func keybindingProjectsButton(g *gocui.Gui) error {
	if err := g.SetKeybinding("import", gocui.MouseLeft, gocui.ModNone, toggleProjectsButton); err != nil {
		return err
	}
	if err := g.SetKeybinding("switch", gocui.MouseLeft, gocui.ModNone, toggleProjectsButton); err != nil {
		return err
	}

	return nil
}

func toggleProjectsButton(g *gocui.Gui, v *gocui.View) error {
	CurrProjectsButton = v.Name()

	SelectProjectsButton(g)

	if v.Name() == "import" {
		InitImportPage(g)
	} else if v.Name() == "switch" {
		InitSwitchPage(g)
	}

	return nil
}

func SelectProjectsButton(g *gocui.Gui) {
	for _, name := range projectsButtons {
		v, err := g.View(name)

		if err == nil {
			if v.Name() == CurrProjectsButton {
				v.Highlight = true
				v.SelBgColor = gocui.ColorWhite
				v.SelFgColor = gocui.ColorBlack
			} else {
				v.Highlight = false
				v.SelBgColor = gocui.ColorBlack
				v.SelFgColor = gocui.ColorDefault
			}
		}
	}
}

func init() {
	CurrProjectsButton = "import"
	projectsButtons = append(projectsButtons, "import", "switch")
}

func DestoryProjectsPage(g *gocui.Gui) {
	for _, v := range ViewMap["projects"] {
		g.DeleteView(v)
		g.DeleteKeybindings(v)
	}
}
