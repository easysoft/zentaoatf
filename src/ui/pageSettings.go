package ui

import (
	"github.com/jroimartin/gocui"
)

var CurrButton string
var buttons []string

func InitSettingsPage(g *gocui.Gui) error {
	DestoryLeftPages(g)
	DestoryRightPages(g)

	importLabel := NewLabelWidgetAutoWidth(g, "switch", 0, 2, "Switch Project")
	ViewMap["settings"] = append(ViewMap["settings"], importLabel.Name())

	switchLabel := NewLabelWidgetAutoWidth(g, "import", 0, 3, "Import from Zentao")
	ViewMap["settings"] = append(ViewMap["settings"], switchLabel.Name())

	keybindings(g)

	return nil
}

func keybindings(g *gocui.Gui) error {
	if err := g.SetKeybinding("import", gocui.MouseLeft, gocui.ModNone, toggleButton); err != nil {
		return err
	}
	if err := g.SetKeybinding("switch", gocui.MouseLeft, gocui.ModNone, toggleButton); err != nil {
		return err
	}

	return nil
}

func toggleButton(g *gocui.Gui, v *gocui.View) error {
	CurrButton = v.Name()

	SelectButton(g)

	if v.Name() == "import" {
		InitImportPage(g)
	} else if v.Name() == "switch" {
		InitSwitchPage(g)
	}

	return nil
}

func SelectButton(g *gocui.Gui) {
	for _, name := range buttons {
		v, err := g.View(name)

		if err == nil {
			if v.Name() == CurrButton {
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
	CurrButton = "import"
	buttons = append(buttons, "import", "switch")
}

func DestorySettingsPage(g *gocui.Gui) {
	for _, v := range ViewMap["settings"] {
		g.DeleteView(v)
		g.DeleteKeybindings(v)
	}
}
