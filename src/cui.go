package main

import (
	"fmt"
	"github.com/jroimartin/gocui"
	"log"
)

func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()

	if v, err := g.SetView("side", 0, 0, int(0.2*float32(maxX)), maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		v.Frame = true
		v.Highlight = true
		v.SelBgColor = gocui.ColorGreen
		v.SelFgColor = gocui.ColorBlack
		fmt.Fprintln(v, "New From Zentao")
		fmt.Fprintln(v, "Switch Project")
	}

	if v, err := g.SetView("main", int(0.2*float32(maxX)), 0, maxX-1, maxY-5); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Autoscroll = true
	}

	if v, err := g.SetView("cmdline", int(0.2*float32(maxX)), maxY-5, maxX-1, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		v.Editable = true
		v.Wrap = true
	}
	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func main() {
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	g.Cursor = true
	g.Mouse = true
	g.SetManagerFunc(layout)

	if err := keybindings(g); err != nil {
		log.Panicln(err)
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}

func keybindings(g *gocui.Gui) error {
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		return err
	}
	if err := g.SetKeybinding("side", gocui.MouseLeft, gocui.ModNone, newProjectFormZentao); err != nil {
		return err
	}
	if err := g.SetKeybinding("cmdline", gocui.MouseLeft, gocui.ModNone, setEdit); err != nil {
		return err
	}
	if err := g.SetKeybinding("msg", gocui.MouseLeft, gocui.ModNone, delMsg); err != nil {
		return err
	}
	return nil
}

func showMsg(g *gocui.Gui, v *gocui.View) error {
	var l string
	var err error

	if _, err := g.SetCurrentView(v.Name()); err != nil {
		return err
	}

	_, cy := v.Cursor()
	if l, err = v.Line(cy); err != nil {
		l = ""
	}

	mainView, err := g.View("main")
	fmt.Fprintln(mainView, l)

	return nil
}

func delMsg(g *gocui.Gui, v *gocui.View) error {
	if err := g.DeleteView("msg"); err != nil {
		return err
	}

	return nil
}

func setEdit(g *gocui.Gui, v *gocui.View) error {
	if _, err := g.SetCurrentView("cmdline"); err != nil {
		return err
	}

	//v.SetOrigin(0, 0)
	//v.SetCursor(0, 0)

	v.SetCursor(0, 0)
	v.Clear()

	return nil
}

func newProjectFormZentao(g *gocui.Gui, slide *gocui.View) error {
	slideView, _ := g.View("side")
	slideX, _ := slideView.Size()

	//mainView, _ := g.View("main")
	//mainX, mainY := mainView.Size()

	labelWidth := 10
	inputWidth := 30
	space := 2
	left := slideX + 3
	right := left + labelWidth

	if v, err := g.SetView("productLabel", left, 1, right, 3); err != nil {
		v.Frame = false
		fmt.Fprintln(v, "ProjectId")
	}

	left = right + space
	right = left + inputWidth
	if v, err := g.SetView("productInput", left, 1, right, 3); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		v.Editable = true
		v.Wrap = true

		if _, err := g.SetCurrentView("productInput"); err != nil {
			return err
		}
	}

	left = right + space
	right = left + 3
	if v, err := g.SetView("or", left, 1, right, 3); err != nil {
		v.Frame = false
		fmt.Fprintln(v, "or")
	}

	left = right + space
	right = left + (labelWidth - 3)
	if v, err := g.SetView("planLabel", left, 1, right, 3); err != nil {
		v.Frame = false
		fmt.Fprintln(v, "PlanId")
	}

	left = right + space
	right = left + inputWidth
	if v, err := g.SetView("planInput", left, 1, right, 3); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		v.Editable = true
		v.Wrap = true
	}

	return nil
}
