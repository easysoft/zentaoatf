package utils

import (
	"flag"
	"fmt"
	"github.com/fatih/color"
	"github.com/jroimartin/gocui"
)

func PrintUsage(flagSet flag.FlagSet) {
	PrintUsageWithSpaceLine(flagSet, true)
}
func PrintUsageWithSpaceLine(flagSet flag.FlagSet, spaceLine bool) {
	prefix := ""
	if spaceLine {
		prefix = "\n"
	}

	fmt.Printf("%s %s \n", prefix, color.CyanString(flagSet.Name()))
	flagSet.PrintDefaults()
}

func PrintSample() {
	fmt.Printf("\nSample to use: \n")
	fmt.Printf("TODO... \n")
}

func PrintToSide(g *gocui.Gui, msg string) {
	slideView, _ := g.View("side")
	slideView.Clear()

	fmt.Fprintln(slideView, msg)
}
func PrintToMain(g *gocui.Gui, msg string) {
	PrintToMainNoScroll(g, msg)
	AdjustOrigin("main")
}
func PrintToMainNoScroll(g *gocui.Gui, msg string) {
	mainView, _ := g.View("main")
	mainView.Clear()

	fmt.Fprintln(mainView, msg)
}

func PrintToCmd(g *gocui.Gui, msg string) {
	cmdView, _ := g.View("cmd")
	_, _ = fmt.Fprintln(cmdView, msg)

	AdjustOrigin("cmd")
}

func ClearSide(g *gocui.Gui) {
	slideView, _ := g.View("side")
	slideView.Clear()
}
