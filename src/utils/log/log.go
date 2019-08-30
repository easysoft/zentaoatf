package logUtils

import (
	"fmt"
	constant "github.com/easysoft/zentaoatf/src/utils/const"
	"github.com/easysoft/zentaoatf/src/utils/i118"
	"github.com/easysoft/zentaoatf/src/utils/vari"
	"github.com/fatih/color"
	"io"
	"strings"
	"unicode/utf8"
)

func PrintWholeLine(msg string, char string, attr color.Attribute) {
	prefixLen := 6
	var postfixLen int
	if vari.RunFromCui {
		maxX, _ := vari.Cui.Size()
		postfixLen = maxX - constant.LeftWidth - utf8.RuneCountInString(msg) - 9
	} else {
		postfixLen = vari.ScreenWidth - utf8.RuneCountInString(msg) - 6
		if postfixLen < 0 { // no width in debug mode
			postfixLen = 6
		}
	}

	preFixStr := strings.Repeat(char, prefixLen)
	postFixStr := strings.Repeat(char, postfixLen)

	var output io.Writer
	if vari.RunFromCui {
		output, _ = vari.Cui.View("cmd")
	} else {
		output = color.Output
	}

	clr := color.New(attr)
	clr.Fprintf(output, fmt.Sprintf("%s%s%s\n", preFixStr, msg, postFixStr))
}

func PrintAndLog(logs *[]string, str string) {
	*logs = append(*logs, str)

	var output io.Writer
	if vari.RunFromCui {
		output, _ = vari.Cui.View("cmd")
	} else {
		output = color.Output
	}

	fmt.Fprintf(output, str+"\n")
}

func PrintAndLogColorLn(logs *[]string, str string, attr color.Attribute) {
	*logs = append(*logs, str)

	var output io.Writer
	if vari.RunFromCui {
		output, _ = vari.Cui.View("cmd")
	} else {
		output = color.Output
	}

	clr := color.New(attr)
	clr.Fprintf(output, str+"\n")
}

func Printt(str string) {
	var output io.Writer
	if vari.RunFromCui {
		output, _ = vari.Cui.View("cmd")
	} else {
		output = color.Output
	}

	fmt.Fprintf(output, str)
}

func ColoredStatus(status string) string {
	temp := strings.ToLower(status)

	switch temp {
	case "pass":
		return color.GreenString(i118Utils.I118Prt.Sprintf(temp))
	case "fail":
		return color.RedString(i118Utils.I118Prt.Sprintf(temp))
	case "skip":
		return color.YellowString(i118Utils.I118Prt.Sprintf(temp))
	}

	return status
}
