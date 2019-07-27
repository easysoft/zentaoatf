package ui

import (
	"fmt"
	"github.com/easysoft/zentaoatf/src/ui/widget"
	"github.com/jroimartin/gocui"
	"strings"
)

const (
	Space = 2
)

func ShowHelp(g *gocui.Gui, v *gocui.View) error {
	help, _ := g.View("help")

	if help != nil {
		HideHelp(g)
	} else {
		widget.NewHelpWidget(g)
	}

	return nil
}
func HideHelp(g *gocui.Gui) error {
	help, _ := g.View("help")

	if help != nil {
		if err := g.DeleteView("help"); err != nil {
			return err
		}
	}

	return nil
}

func ToggleInput(views []string) func(g *gocui.Gui, v *gocui.View) error {
	return func(g *gocui.Gui, v *gocui.View) error {
		nextView := GetNextView(v.Name(), views)

		if nextView != "" {
			_, err := g.SetCurrentView(nextView)
			return err
		}

		return nil
	}
}

func GetNextView(name string, views []string) string {
	i := 0
	found := false
	for true {
		if name == views[i] {
			found = true
			i++
			i = i % len(views)
			continue
		}

		if found {
			if strings.Index(views[i], "Input") > -1 {
				return views[i]
			}
		}

		i++
		i = i % len(views)
	}

	return ""
}

func Quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func PrintToView(msg string, v *gocui.View) {
	fmt.Fprintln(v, msg)
}
