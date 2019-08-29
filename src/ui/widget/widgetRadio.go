package widget

import (
	"fmt"
	"github.com/easysoft/zentaoatf/src/utils/vari"
	"github.com/jroimartin/gocui"
	"strings"
)

type RadioWidget struct {
	name string
	x, y int
	w    int
	val  string

	checked bool
	handler func(g *gocui.Gui, v *gocui.View) error
}

func NewRadioWidget(name string, x, y int, checked bool) *gocui.View {
	var val string
	if checked {
		val = "[*]"
	} else {
		val = "[ ]"
	}

	widget := RadioWidget{name: name, x: x, y: y, w: len(val) + 1, val: val, checked: checked, handler: handler}
	v, _ := widget.Layout()
	return v
}

func (w *RadioWidget) Layout() (*gocui.View, error) {
	v, err := vari.Cui.SetView(w.name, w.x, w.y, w.x+w.w, w.y+2)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return nil, err
		}
		if err := vari.Cui.SetKeybinding(w.name, gocui.KeySpace, gocui.ModNone, w.handler); err != nil {
			return nil, err
		}
		v.Frame = false

		fmt.Fprint(v, w.val)
	}
	return v, nil
}

func handler(g *gocui.Gui, v *gocui.View) error {
	val := strings.TrimSpace(v.Buffer())
	if val == "[*]" {
		val = "[ ]"
	} else {
		val = "[*]"
	}

	v.Clear()
	fmt.Fprint(v, val)

	return nil
}

func ParseRadioVal(val string) bool {
	if val == "[*]" {
		return true
	}

	return false
}
