package biz

import (
	"fmt"
	"github.com/easysoft/zentaoatf/src/utils"
	"github.com/fatih/color"
	"io"
	"strings"
	"unicode/utf8"
)

func PrintWholeLine(msg string, char string, attr color.Attribute) {
	prefixLen := 6
	var postfixLen int
	if utils.RunFromCui {
		maxX, _ := utils.Cui.Size()
		postfixLen = maxX - utils.LeftWidth - utf8.RuneCountInString(msg) - 8
	} else {
		postfixLen = utils.Prefer.Width - utf8.RuneCountInString(msg) - 6
		if postfixLen < 0 { // no width in debug mode
			postfixLen = 6
		}
	}

	preFixStr := strings.Repeat(char, prefixLen)
	postFixStr := strings.Repeat(char, postfixLen)

	var output io.Writer
	if utils.RunFromCui {
		output, _ = utils.Cui.View("main")
	} else {
		output = color.Output
	}

	clr := color.New(attr)
	clr.Fprintf(output, fmt.Sprintf("%s%s%s\n", preFixStr, msg, postFixStr))
}

func PrintAndLog(logs *[]string, str string) {
	*logs = append(*logs, str)

	var output io.Writer
	if utils.RunFromCui {
		output, _ = utils.Cui.View("main")
	} else {
		output = color.Output
	}

	fmt.Fprintf(output, str+"\n")
}

func PrintAndLogColorLn(logs *[]string, str string, attr color.Attribute) {
	*logs = append(*logs, str)

	var output io.Writer
	if utils.RunFromCui {
		output, _ = utils.Cui.View("main")
	} else {
		output = color.Output
	}

	clr := color.New(attr)
	clr.Fprintf(output, str+"\n")
}

func Printt(str string) {
	var output io.Writer
	if utils.RunFromCui {
		output, _ = utils.Cui.View("main")
	} else {
		output = color.Output
	}

	fmt.Fprintf(output, str+"\n")
}

func coloredStatus(status string) string {
	temp := strings.ToLower(status)

	switch temp {
	case "pass":
		return color.GreenString(utils.I118Prt.Sprintf(temp))
	case "fail":
		return color.RedString(utils.I118Prt.Sprintf(temp))
	case "skip":
		return color.YellowString(utils.I118Prt.Sprintf(temp))
	}

	return status
}
