package logUtils

import (
	"encoding/json"
	"fmt"
	constant "github.com/easysoft/zentaoatf/src/utils/const"
	stringUtils "github.com/easysoft/zentaoatf/src/utils/string"
	"github.com/easysoft/zentaoatf/src/utils/vari"
	"github.com/fatih/color"
	"strings"
)

func PrintUsage() {
	fmt.Println(color.CyanString("\nUsage: "))
	fmt.Fprintf(color.Output, "%s\n", constant.Usage)

	fmt.Println(color.CyanString("Example: "))
	fmt.Fprintf(color.Output, "%s\n", constant.Example)
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
}

func PrintStructToCmd(obj interface{}) {
	str := stringUtils.StructToStr(obj)
	PrintToCmd(str)
}

func ClearSide() {
	slideView, _ := vari.Cui.View("side")
	slideView.Clear()
}

func PrintUnicode(str []byte) {
	var a interface{}

	temp := strings.Replace(string(str), "\\\\", "\\", -1)

	err := json.Unmarshal([]byte(temp), &a)

	var msg string
	if err == nil {
		msg = fmt.Sprint(a)
	} else {
		msg = temp
	}

	if !vari.RunFromCui {
		fmt.Println(msg)
	} else {
		PrintToCmd(msg)
	}
}
