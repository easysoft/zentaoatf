package ui

import (
	"github.com/jroimartin/gocui"
	"strings"
)

const (
	Space = 2
)

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

func scroll(dy int) func(g *gocui.Gui, v *gocui.View) error {
	return func(g *gocui.Gui, v *gocui.View) error {
		return scrollView(g, v, dy)
	}
}

func scrollView(g *gocui.Gui, v *gocui.View, dy int) error {
	v.Autoscroll = false
	ox, oy := v.Origin()
	pos := oy + dy
	_, height := v.Size()

	if pos > len(v.BufferLines())-height {
		pos = len(v.BufferLines()) - height
	}
	if pos < 0 {
		pos = 0
	}

	v.SetOrigin(ox, pos)

	return nil
}

func setCurrView(name string) func(g *gocui.Gui, v *gocui.View) error {
	return func(g *gocui.Gui, v *gocui.View) error {
		g.SetCurrentView(name)
		return nil
	}
}
