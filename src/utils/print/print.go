package printUtils

import (
	"flag"
	"fmt"
	logUtils "github.com/easysoft/zentaoatf/src/utils/log"
	"github.com/easysoft/zentaoatf/src/utils/vari"
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
	if !vari.RunFromCui {
		fmt.Println(msg)
		return
	}
	slideView, _ := vari.Cui.View("side")
	slideView.Clear()

	fmt.Fprintln(slideView, msg)
}

func PrintToMainNoScroll(msg string) {
	if !vari.RunFromCui {
		fmt.Println(msg)
		return
	}
	mainView, _ := vari.Cui.View("main")
	mainView.Clear()

	fmt.Fprintln(mainView, msg)
}

func PrintToCmd(msg string) {
	if !vari.RunFromCui {
		fmt.Println(msg)
		return
	}
	cmdView, _ := vari.Cui.View("cmd")
	_, _ = fmt.Fprintln(cmdView, msg)

	logUtils.AdjustOrigin("cmd")
}

func ClearSide() {
	slideView, _ := vari.Cui.View("side")
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
