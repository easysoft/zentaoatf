package ui

import (
	"github.com/easysoft/zentaoatf/src/utils/vari"
	"github.com/jroimartin/gocui"
)

var CurrSettingsButton string
var settingsButtons []string

func InitSettingsPage() error {
	importLabel := NewLabelWidgetAutoWidth("switch", 0, 2, "Switch Work dir")
	ViewMap["settings"] = append(ViewMap["settings"], importLabel.Name())

	switchLabel := NewLabelWidgetAutoWidth("import", 0, 3, "Import from Zentao")
	ViewMap["settings"] = append(ViewMap["settings"], switchLabel.Name())

	keybindingSettingsButton()

	return nil
}

func keybindingSettingsButton() error {
	if err := vari.Cui.SetKeybinding("import", gocui.MouseLeft, gocui.ModNone, toggleSettingsButton); err != nil {
		return err
	}
	if err := vari.Cui.SetKeybinding("switch", gocui.MouseLeft, gocui.ModNone, toggleSettingsButton); err != nil {
		return err
	}

	return nil
}

func toggleSettingsButton(g *gocui.Gui, v *gocui.View) error {
	CurrSettingsButton = v.Name()

	SelectSettingsButton()

	if v.Name() == "import" {
		InitImportPage()
	} else if v.Name() == "switch" {
		InitSwitchPage()
	}

	return nil
}

func SelectSettingsButton() {
	for _, name := range settingsButtons {
		v, err := vari.Cui.View(name)

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

func DestorySettingsPage() {
	for _, v := range ViewMap["settings"] {
		vari.Cui.DeleteView(v)
		vari.Cui.DeleteKeybindings(v)
	}
}
