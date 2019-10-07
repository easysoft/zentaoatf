package logUtils

import (
	"encoding/json"
	"fmt"
	commonUtils "github.com/easysoft/zentaoatf/src/utils/common"
	fileUtils "github.com/easysoft/zentaoatf/src/utils/file"
	"github.com/easysoft/zentaoatf/src/utils/vari"
	"github.com/fatih/color"
	"io"
	"os"
	"regexp"
	"strings"
)

var (
	usageFile  = fmt.Sprintf("res%sdoc%susage.txt", string(os.PathSeparator), string(os.PathSeparator))
	sampleFile = fmt.Sprintf("res%sdoc%ssample.txt", string(os.PathSeparator), string(os.PathSeparator))
)

func PrintUsage() {
	PrintToStdOut("Usage: ", color.FgCyan)

	content := fileUtils.ReadResData(usageFile)
	fmt.Printf("%s\n", content)

	PrintToStdOut("\nExample: ", color.FgCyan)

	content = fileUtils.ReadResData(sampleFile)
	if !commonUtils.IsWin() {
		regx, _ := regexp.Compile(`\\`)
		content = regx.ReplaceAllString(content, "/")

		regx, _ = regexp.Compile(`ztf.exe`)
		content = regx.ReplaceAllString(content, "ztf")

		regx, _ = regexp.Compile(`/bat/`)
		content = regx.ReplaceAllString(content, "/shell/")

		regx, _ = regexp.Compile(`\.bat\s{4}`)
		content = regx.ReplaceAllString(content, ".shell")
	}
	fmt.Printf("%s\n", content)
}

func PrintTo(str string) {
	var output io.Writer
	if vari.RunFromCui {
		output, _ = vari.Cui.View("cmd")
	} else {
		output = color.Output
	}

	fmt.Fprint(output, str+"\n")
}

func PrintToStdOut(msg string, attr color.Attribute) {
	output := color.Output

	if attr == -1 {
		fmt.Fprint(output, msg+"\n")
	} else {
		color.New(attr).Fprintf(output, msg+"\n")
	}
}

func PrintToCmd(msg string, attr color.Attribute) {
	var output io.Writer
	if vari.RunFromCui {
		output, _ = vari.Cui.View("cmd")
	} else {
		output = color.Output
	}

	if attr == -1 {
		fmt.Fprint(output, msg+"\n")
	} else {
		clr := color.New(attr)
		clr.Fprint(output, msg+"\n")
	}
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

	PrintToCmd(msg, -1)
}
