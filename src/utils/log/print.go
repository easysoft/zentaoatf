package logUtils

import (
	"encoding/json"
	"fmt"
	fileUtils "github.com/easysoft/zentaoatf/src/utils/file"
	"github.com/easysoft/zentaoatf/src/utils/vari"
	"github.com/fatih/color"
	"io"
	"os"
	"strings"
)

var (
	usageFile  = fmt.Sprintf("res%sdoc%susage.txt", string(os.PathSeparator), string(os.PathSeparator))
	sampleFile = fmt.Sprintf("res%sdoc%sample.txt", string(os.PathSeparator), string(os.PathSeparator))
)

func PrintUsage() {
	PrintToStdOut("Usage: ", color.FgCyan)

	content := fileUtils.ReadResData(usageFile)
	fmt.Printf(" %s\n", content)

	PrintToStdOut("\nExample: ", color.FgCyan)

	content = fileUtils.ReadResData(usageFile)
	fmt.Printf(" %s", content)
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
