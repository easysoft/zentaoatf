package utils

import (
	"fmt"
	"github.com/fatih/color"
	"io"
	"strings"
	"unicode/utf8"
)

func PrintWholeLine(msg string, char string, attr color.Attribute) {
	prefixLen := 6
	var postfixLen int
	if RunFromCui {
		maxX, _ := Cui.Size()
		postfixLen = maxX - LeftWidth - utf8.RuneCountInString(msg) - 9
	} else {
		postfixLen = Prefer.Width - utf8.RuneCountInString(msg) - 6
		if postfixLen < 0 { // no width in debug mode
			postfixLen = 6
		}
	}

	preFixStr := strings.Repeat(char, prefixLen)
	postFixStr := strings.Repeat(char, postfixLen)

	var output io.Writer
	if RunFromCui {
		output, _ = Cui.View(CuiRunOutputView)
	} else {
		output = color.Output
	}

	clr := color.New(attr)
	clr.Fprintf(output, fmt.Sprintf("%s%s%s\n", preFixStr, msg, postFixStr))
	AdjustOrigin(CuiRunOutputView)
}

func PrintAndLog(logs *[]string, str string) {
	*logs = append(*logs, str)

	var output io.Writer
	if RunFromCui {
		output, _ = Cui.View(CuiRunOutputView)
	} else {
		output = color.Output
	}

	fmt.Fprintf(output, str+"\n")
	AdjustOrigin(CuiRunOutputView)
}

func PrintAndLogColorLn(logs *[]string, str string, attr color.Attribute) {
	*logs = append(*logs, str)

	var output io.Writer
	if RunFromCui {
		output, _ = Cui.View(CuiRunOutputView)
	} else {
		output = color.Output
	}

	clr := color.New(attr)
	clr.Fprintf(output, str+"\n")
	AdjustOrigin(CuiRunOutputView)
}

func Printt(str string) {
	var output io.Writer
	if RunFromCui {
		output, _ = Cui.View(CuiRunOutputView)
	} else {
		output = color.Output
	}

	fmt.Fprintf(output, str)
	AdjustOrigin(CuiRunOutputView)
}

func ColoredStatus(status string) string {
	temp := strings.ToLower(status)

	switch temp {
	case "pass":
		return color.GreenString(I118Prt.Sprintf(temp))
	case "fail":
		return color.RedString(I118Prt.Sprintf(temp))
	case "skip":
		return color.YellowString(I118Prt.Sprintf(temp))
	}

	return status
}
