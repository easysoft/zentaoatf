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

func PrintToCmd(g *gocui.Gui, msg string) {
	cmdView, _ := g.View("cmd")
	_, _ = fmt.Fprintln(cmdView, msg)

	AdjustOrigin("cmd")
}
func PrintToMain(g *gocui.Gui, msg string) {
	mainView, _ := g.View("main")
	mainView.Clear()

	_, _ = fmt.Fprintln(mainView, msg)
	AdjustOrigin("main")
}
