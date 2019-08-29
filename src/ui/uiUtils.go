package ui

import (
	"github.com/easysoft/zentaoatf/src/utils/log"
	"github.com/easysoft/zentaoatf/src/utils/vari"
	"github.com/jroimartin/gocui"
	"regexp"
	"strings"
)

const (
	Space = 2
)

func BindEventForInputWidgets(arr []string) {
	for _, v := range arr {
		if isInput(v) {
			vari.Cui.SetKeybinding(v, gocui.MouseLeft, gocui.ModNone, SetCurrView(v))
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
	v, err := vari.Cui.View(name)
	if err != nil {
		logUtils.PrintToCmd(err.Error() + ": " + name)
		return nil
	}

	v.Wrap = true

	if err := vari.Cui.SetKeybinding(name, gocui.MouseLeft, gocui.ModNone, SetCurrView(name)); err != nil {
		return err
	}
	if err := vari.Cui.SetKeybinding(name, gocui.KeyArrowUp, gocui.ModNone, scrollEvent(-1)); err != nil {
		return err
	}
	if err := vari.Cui.SetKeybinding(name, gocui.KeyArrowDown, gocui.ModNone, scrollEvent(1)); err != nil {
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
		atBottom := scroll(v, dy)
		if atBottom {
			return true
		}
	} else {
		v.SetCursor(cx, newCy)
	}

	return false
}
func scroll(v *gocui.View, dy int) bool {
	_, h := v.Size()

	ox, oy := v.Origin()
	newOy := oy + dy

	// If we're at the bottom...
	if newOy+h >= strings.Count(v.ViewBuffer(), "\n") {
		// Set autoscroll to normal again.
		v.Autoscroll = true

		return true
	} else {
		// Set autoscroll to false and scroll.
		v.Autoscroll = false
		v.SetOrigin(ox, newOy)

		return false
	}
}

func SupportRowHighlight(name string) error {
	v, _ := vari.Cui.View(name)

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
		v, _ := vari.Cui.View(name)

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
