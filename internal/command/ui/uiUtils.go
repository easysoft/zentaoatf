package ui

import (
	commConsts "github.com/aaronchen2k/deeptest/internal/comm/consts"
	"github.com/awesome-gocui/gocui"
	"regexp"
	"strings"
)

const (
	Space = 2
)

func BindEventForInputWidgets(arr []string) {
	for _, v := range arr {
		if isInput(v) {
			commConsts.Cui.SetKeybinding(v, gocui.MouseLeft, gocui.ModNone, SetCurrView(v))
		}
	}
}
func isInput(v string) bool {
	return strings.Index(v, "Input") > -1
}

func Quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func SupportScroll(name string) error {
	v, err := commConsts.Cui.View(name)
	if err != nil {
		return nil
	}

	v.Wrap = true

	if err := commConsts.Cui.SetKeybinding(name, gocui.MouseLeft, gocui.ModNone, SetCurrView(name)); err != nil {
		return err
	}
	if err := commConsts.Cui.SetKeybinding(name, gocui.KeyArrowUp, gocui.ModNone, scrollEvent(-1)); err != nil {
		return err
	}
	if err := commConsts.Cui.SetKeybinding(name, gocui.KeyArrowDown, gocui.ModNone, scrollEvent(1)); err != nil {
		return err
	}

	return nil
}

func scrollEvent(dy int) func(g *gocui.Gui, v *gocui.View) error {
	return func(g *gocui.Gui, v *gocui.View) error {
		ScrollAction(v, dy)
		return nil
	}
}

func ScrollAction(v *gocui.View, dy int) bool {
	// Get the size and position of the view.
	cx, cy := v.Cursor()
	_, h := v.Size()

	newCy := cy + dy
	//logUtils.PrintToCmd(fmt.Sprintf("%d - %d", cy, dy))
	if (cy == 0 && dy < 0) || // top
		(newCy == h && dy > 0) { // bottom

		atBottom := scroll(v, dy) // A. scroll
		if atBottom {
			return true
		}
	} else {
		v.SetCursor(cx, newCy) // B. move
	}

	return false
}
func scroll(v *gocui.View, dy int) bool {
	_, h := v.Size()

	ox, oy := v.Origin()
	newOy := oy + dy

	// If we're at the bottom...
	if newOy+h > strings.Count(v.ViewBuffer(), "\n") {
		//logUtils.PrintToCmd(fmt.Sprintf("=1= %d", time.Now().Unix()), -1)

		// Set autoscroll to normal again.
		v.Autoscroll = true

		return true
	} else {
		//logUtils.PrintToCmd(fmt.Sprintf("=2= %d", time.Now().Unix()), -1)

		// Set autoscroll to false and scroll.
		v.Autoscroll = false
		v.SetOrigin(ox, newOy)

		return false
	}
}

func SupportRowHighlight(name string) error {
	v, _ := commConsts.Cui.View(name)

	v.Wrap = true
	v.SelBgColor = gocui.ColorWhite
	v.SelFgColor = gocui.ColorBlack

	return nil
}

func SetCurrView(name string) func(g *gocui.Gui, v *gocui.View) error {
	return func(g *gocui.Gui, v *gocui.View) error {
		g.SetCurrentView(name)
		return nil
	}
}

func HighlightTab(view string, views []string) {
	for _, name := range views {
		v, _ := commConsts.Cui.View(name)

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

func GetSelectedRowVal(v *gocui.View) string {
	line, _ := getSelectedRow(v, ".*")

	return line
}
func getSelectedRow(v *gocui.View, reg string) (string, error) {
	var line string
	var err error

	_, cy := v.Cursor()
	if line, err = v.Line(cy); err != nil {
		return "", nil
	}
	line = strings.TrimSpace(line)

	pass, _ := regexp.MatchString(reg, line)

	if !pass {
		return "", nil
	}

	return line, nil
}
