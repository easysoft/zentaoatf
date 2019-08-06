package utils

import (
	"flag"
	"fmt"
	"github.com/fatih/color"
	"reflect"
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

func PrintToSide(msg string) {
	slideView, _ := Cui.View("side")
	slideView.Clear()

	fmt.Fprintln(slideView, msg)
}

func PrintToMainNoScroll(msg string) {
	mainView, _ := Cui.View("main")
	mainView.Clear()

	fmt.Fprintln(mainView, msg)
}

func PrintToCmd(msg string) {
	cmdView, _ := Cui.View("cmd")
	_, _ = fmt.Fprintln(cmdView, msg)

	AdjustOrigin("cmd")
}

func ClearSide() {
	slideView, _ := Cui.View("side")
	slideView.Clear()
}

func PrintStruct(obj interface{}) {
	val := reflect.ValueOf(obj)
	typeOfS := val.Type()
	for i := 0; i < reflect.ValueOf(obj).NumField(); i++ {
		val := val.Field(i)
		fmt.Printf("  %s: %v\n", typeOfS.Field(i).Name, val.Interface())
	}
}

func PrintMap(obj map[string]interface{}) {
	for key, val := range obj {
		fmt.Printf("  %s: %v\n", key, val)
	}
}
