package logUtils

import (
	"encoding/json"
	"fmt"
	"github.com/easysoft/zentaoatf/src/utils/common"
	"github.com/easysoft/zentaoatf/src/utils/const"
	"github.com/easysoft/zentaoatf/src/utils/file"
	"github.com/fatih/color"
	"os"
	"regexp"
	"strings"
)

var (
	usageFile  = fmt.Sprintf("res%sdoc%susage.txt", string(os.PathSeparator), string(os.PathSeparator))
	sampleFile = fmt.Sprintf("res%sdoc%ssample.txt", string(os.PathSeparator), string(os.PathSeparator))
)

func PrintUsage() {
	PrintToWithColor("Usage: ", color.FgCyan)

	usage := fileUtils.ReadResData(usageFile)
	exeFile := constant.AppName
	if commonUtils.IsWin() {
		exeFile += ".exe"
	}
	usage = fmt.Sprintf(usage, exeFile)
	fmt.Printf("%s\n", usage)

	PrintToWithColor("\nExample: ", color.FgCyan)
	sample := fileUtils.ReadResData(sampleFile)
	if !commonUtils.IsWin() {
		regx, _ := regexp.Compile(`\\`)
		sample = regx.ReplaceAllString(sample, "/")

		regx, _ = regexp.Compile(constant.AppName + `.exe`)
		sample = regx.ReplaceAllString(sample, constant.AppName)

		regx, _ = regexp.Compile(`/bat/`)
		sample = regx.ReplaceAllString(sample, "/shell/")

		regx, _ = regexp.Compile(`\.bat\s{4}`)
		sample = regx.ReplaceAllString(sample, ".shell")
	}
	fmt.Printf("%s\n", sample)
}

func PrintTo(str string) {
	output := color.Output
	fmt.Fprint(output, str+"\n")
}
func PrintTof(format string, params ...interface{}) {
	output := color.Output
	fmt.Fprintf(output, format+"\n", params...)
}

func PrintToWithColor(msg string, attr color.Attribute) {
	output := color.Output

	if attr == -1 {
		fmt.Fprint(output, msg+"\n")
	} else {
		color.New(attr).Fprintf(output, msg+"\n")
	}
}

func PrintToCmd(msg string, attr color.Attribute) {
	output := color.Output

	if attr == -1 {
		fmt.Fprint(output, msg+"\n")
	} else {
		clr := color.New(attr)
		clr.Fprint(output, msg+"\n")
	}
}

func ConvertUnicode(str []byte) string {
	var a interface{}

	temp := strings.Replace(string(str), "\\\\", "\\", -1)

	err := json.Unmarshal([]byte(temp), &a)

	var msg string
	if err == nil {
		msg = fmt.Sprint(a)
	} else {
		msg = temp
	}

	return msg
}
