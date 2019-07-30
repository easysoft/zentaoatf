package ui

import (
	"github.com/jroimartin/gocui"
)

var CurrSettingsButton string
var settingsButtons []string

func InitSettingsPage(g *gocui.Gui) error {
	importLabel := NewLabelWidgetAutoWidth(g, "switch", 0, 2, "Switch Work dir")
	ViewMap["settings"] = append(ViewMap["settings"], importLabel.Name())

	switchLabel := NewLabelWidgetAutoWidth(g, "import", 0, 3, "Import from Zentao")
	ViewMap["settings"] = append(ViewMap["settings"], switchLabel.Name())

	keybindingSettingsButton(g)

	return nil
}

func keybindingSettingsButton(g *gocui.Gui) error {
	if err := g.SetKeybinding("import", gocui.MouseLeft, gocui.ModNone, toggleSettingsButton); err != nil {
		return err
	}
	if err := g.SetKeybinding("switch", gocui.MouseLeft, gocui.ModNone, toggleSettingsButton); err != nil {
		return err
	}

	return nil
}

func toggleSettingsButton(g *gocui.Gui, v *gocui.View) error {
	CurrSettingsButton = v.Name()

	SelectSettingsButton(g)

	if v.Name() == "import" {
		InitImportPage(g)
	} else if v.Name() == "switch" {
		InitSwitchPage(g)
	}

	return nil
}

func SelectSettingsButton(g *gocui.Gui) {
	for _, name := range settingsButtons {
		v, err := g.View(name)

		if err == nil {
			if v.Name() == CurrSettingsButton {
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
	CurrSettingsButton = "import"
	settingsButtons = append(settingsButtons, "import", "switch")
}

func DestorySettingsPage(g *gocui.Gui) {
	for _, v := range ViewMap["settings"] {
		g.DeleteView(v)
		g.DeleteKeybindings(v)
	}
}
