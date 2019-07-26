package ui

import (
	"fmt"
	"github.com/jroimartin/gocui"
)

const (
	radioWidth = 10
)

type RadioWidget struct {
	name string
	x, y int
	w    int
	body string

	checked bool
	handler func(g *gocui.Gui, v *gocui.View) error
}

func NewRadioWidget(name string, x, y int, checked bool) *RadioWidget {
	var body string
	if checked {
		body = "[*]"
	} else {
		body = "[ ]"
	}

	return &RadioWidget{name: name, x: x, y: y, w: len(body) + 1, body: body, checked: checked, handler: handler()}
}

func (w *RadioWidget) Layout(g *gocui.Gui) error {
	v, err := g.SetView(w.name, w.x, w.y, w.x+w.w, w.y+2)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		if err := g.SetKeybinding(w.name, gocui.KeySpace, gocui.ModNone, w.handler); err != nil {
			return err
		}
		v.Clear()

		fmt.Fprint(v, w.body)
	}
	return nil
}

func handler() func(g *gocui.Gui, v *gocui.View) error {
	return func(g *gocui.Gui, v *gocui.View) error {
		return statusSet(v)
	}
}

func statusSet(v *gocui.View) error {
	body := v.Buffer()
	if body == "[*]" {
		body = "[ ]"
	} else {
		body = "[*]"
	}

	return nil
}

func ParseRadioVal(val string) bool {
	if val == "[*]" {
		return true
	}

	return false
}
