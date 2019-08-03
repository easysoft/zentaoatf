package ui

import (
	"github.com/easysoft/zentaoatf/src/utils"
	"github.com/jroimartin/gocui"
)

var CurrSettingsButton string
var settingsButtons []string

func InitSettingsPage() error {
	importLabel := NewLabelWidgetAutoWidth(utils.Cui, "switch", 0, 2, "Switch Work dir")
	ViewMap["settings"] = append(ViewMap["settings"], importLabel.Name())

	switchLabel := NewLabelWidgetAutoWidth(utils.Cui, "import", 0, 3, "Import from Zentao")
	ViewMap["settings"] = append(ViewMap["settings"], switchLabel.Name())

	keybindingSettingsButton(utils.Cui)

	return nil
}

func keybindingSettingsButton() error {
	if err := utils.Cui.SetKeybinding("import", gocui.MouseLeft, gocui.ModNone, toggleSettingsButton); err != nil {
		return err
	}
	if err := utils.Cui.SetKeybinding("switch", gocui.MouseLeft, gocui.ModNone, toggleSettingsButton); err != nil {
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
		v, err := utils.Cui.View(name)

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
		utils.Cui.DeleteView(v)
		utils.Cui.DeleteKeybindings(v)
	}
}
