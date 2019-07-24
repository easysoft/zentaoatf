package main

import (
	"fmt"
	"github.com/jroimartin/gocui"
	"log"
)

const (
	leftWidth      = 32
	labelWidth     = 10
	inputNumbWidth = 25
	buttonWidth    = 10
	space          = 2
)

func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()

	if v, err := g.SetView("qickbar", 0, 0, leftWidth, 2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Highlight = true
		v.SelBgColor = gocui.ColorGreen
		v.SelFgColor = gocui.ColorBlack
	}
	if v, err := g.SetView("import", 3, 0, 14, 2); err != nil {
		v.Frame = false
		fmt.Fprintln(v, "  Import   ")
	}
	if v, err := g.SetView("switch", 19, 0, 31, 2); err != nil {
		v.Frame = false
		fmt.Fprintln(v, "  Switch   ")
	}

	if v, err := g.SetView("side", 0, 2, leftWidth, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Highlight = true
		v.SelBgColor = gocui.ColorGreen
		v.SelFgColor = gocui.ColorBlack
	}

	if v, err := g.SetView("main", leftWidth, 0, maxX-1, maxY-5); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Autoscroll = true
	}

	if v, err := g.SetView("cmdline", leftWidth, maxY-5, maxX-1, maxY-1); err != nil {
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

	if err := g.SetKeybinding("import", gocui.MouseLeft, gocui.ModNone, importProject); err != nil {
		return err
	}
	if err := g.SetKeybinding("switch", gocui.MouseLeft, gocui.ModNone, switchProject); err != nil {
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

func importProject(g *gocui.Gui, slide *gocui.View) error {
	maxX, _ := g.Size()

	slideView, _ := g.View("side")
	slideX, _ := slideView.Size()

	mainView, _ := g.View("main")
	_, mainY := mainView.Size()

	left := slideX + 2
	right := left + labelWidth

	if v, err := g.SetView("productLabel", left, 1, right, 3); err != nil {
		v.Frame = false
		fmt.Fprintln(v, "ProdoctId")
	}

	left = right + space
	right = left + inputNumbWidth
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
	right = left + inputNumbWidth
	if v, err := g.SetView("planInput", left, 1, right, 3); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		v.Editable = true
		v.Wrap = true
	}

	buttonX := (maxX-leftWidth)/2 + leftWidth - buttonWidth
	buttonY := mainY - 2
	if v, err := g.SetView("submit", buttonX, buttonY, buttonX+buttonWidth, buttonY+2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		fmt.Fprintln(v, "  Submit  ")
	}

	return nil
}

func switchProject(g *gocui.Gui, slide *gocui.View) error {
	return nil
}
