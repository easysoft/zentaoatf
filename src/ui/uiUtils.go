package ui

import (
	"github.com/easysoft/zentaoatf/src/utils"
	"github.com/jroimartin/gocui"
	"strings"
)

const (
	Space = 2
)

//func ToggleInput(views []string) func(g *gocui.Gui, v *gocui.View) error {
//	return func(g *gocui.Gui, v *gocui.View) error {
//		var viewName string
//		if v == nil {
//			viewName = views[0]
//		} else {
//			viewName = v.Name()
//		}
//
//		nextView := GetNextView(viewName, views)
//		utils.PrintToCmd(g, nextView)
//
//		if nextView != "" {
//			_, err := g.SetCurrentView(nextView)
//			return err
//		}
//
//		return nil
//	}
//}
//
//func GetNextView(name string, views []string) string {
//	i := 0
//	found := false
//	for true {
//		if name == views[i] {
//			found = true
//			i++
//			i = i % len(views)
//			continue
//		}
//
//		if found {
//			if IsInput(views[i]) {
//				return views[i]
//			}
//		}
//
//		i++
//		i = i % len(views)
//	}
//
//	return ""
//}

func keyBindsInput(arr []string) {
	for _, v := range arr {
		if IsInput(v) {
			setInputEvent(utils.Cui, v)
		}
	}
}

func IsInput(v string) bool {
	return strings.Index(v, "Input") > -1
}

func Quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func scrollEvent(dy int) func(g *gocui.Gui, v *gocui.View) error {
	return func(g *gocui.Gui, v *gocui.View) error {
		return scrollAction(g, v, dy)
	}
}

func scrollAction(g *gocui.Gui, v *gocui.View, dy int) error {
	v.Autoscroll = false

	if dy > 0 {
		_, oy := v.Origin()
		cx, cy := v.Cursor()

		pos := oy + dy
		_, height := v.Size()

		if err := v.SetCursor(cx, cy+1); err != nil {
			ox, oy := v.Origin()

			if pos < len(v.BufferLines())-height-1 {
				if err := v.SetOrigin(ox, oy+1); err != nil {
					return err
				}
			}
		}
	} else if dy < 0 {
		ox, oy := v.Origin()
		cx, cy := v.Cursor()

		if err := v.SetCursor(cx, cy-1); err != nil && oy > 0 {
			if err := v.SetOrigin(ox, oy-1); err != nil {
				return err
			}
		}
	}

	return nil
}

func setViewScroll(g *gocui.Gui, name string) error {
	if err := g.SetKeybinding(name, gocui.MouseLeft, gocui.ModNone, setCurrView(name)); err != nil {
		return err
	}
	if err := g.SetKeybinding(name, gocui.KeyArrowUp, gocui.ModNone, scrollEvent(-1)); err != nil {
		return err
	}
	if err := g.SetKeybinding(name, gocui.KeyArrowDown, gocui.ModNone, scrollEvent(1)); err != nil {
		return err
	}

	return nil
}

func setHighlight(g *gocui.Gui, name string) error {
	v, _ := g.View(name)

	v.Wrap = true
	v.Highlight = true
	v.SelBgColor = gocui.ColorWhite
	v.SelFgColor = gocui.ColorBlack

	return nil
}

func setViewLineSelected(g *gocui.Gui, name string, selectLine func(g *gocui.Gui, v *gocui.View) error) error {
	if err := g.SetKeybinding(name, gocui.KeyEnter, gocui.ModNone, selectLine); err != nil {
		return err
	}
	if err := g.SetKeybinding(name, gocui.MouseLeft, gocui.ModNone, selectLine); err != nil {
		return err
	}

	return nil
}

func setCurrView(name string) func(g *gocui.Gui, v *gocui.View) error {
	return func(g *gocui.Gui, v *gocui.View) error {
		g.SetCurrentView(name)
		return nil
	}
}

func setInputEvent(g *gocui.Gui, name string) error {
	if err := g.SetKeybinding(name, gocui.MouseLeft, gocui.ModNone, setCurrView(name)); err != nil {
		return err
	}
	return nil
}

func HighlightTab(view string, views []string) {
	for _, name := range views {
		v, _ := utils.Cui.View(name)

		if v.Name() == view {
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
