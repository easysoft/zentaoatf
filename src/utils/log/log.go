package logUtils

import (
	"fmt"
	fileUtils "github.com/easysoft/zentaoatf/src/utils/file"
	"github.com/easysoft/zentaoatf/src/utils/i118"
	"github.com/easysoft/zentaoatf/src/utils/vari"
	"github.com/fatih/color"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"strings"
	"unicode/utf8"
)

var Logger *logrus.Logger

func PrintWholeLine(msg string, char string, attr color.Attribute) {
	prefixLen := (vari.ScreenWidth - utf8.RuneCountInString(msg)) / 2
	if prefixLen <= 0 { // no width in debug mode
		prefixLen = 6
	}
	postfixLen := vari.ScreenWidth - utf8.RuneCountInString(msg) - prefixLen
	if postfixLen <= 0 { // no width in debug mode
		postfixLen = 6
	}

	preFixStr := strings.Repeat(char, prefixLen)
	postFixStr := strings.Repeat(char, postfixLen)

	clr := color.New(attr)
	clr.Fprintf(color.Output, fmt.Sprintf("%s%s%s\n", preFixStr, msg, postFixStr))
}

func PrintAndLog(str string) {
	var output io.Writer
	if vari.RunFromCui {
		output, _ = vari.Cui.View("cmd")
	} else {
		output = color.Output
	}

	fmt.Fprintf(output, str+"\n")
}

func PrintAndLogColorLn(str string, attr color.Attribute) {
	var output io.Writer
	if vari.RunFromCui {
		output, _ = vari.Cui.View("cmd")
	} else {
		output = color.Output
	}

	clr := color.New(attr)
	clr.Fprintf(output, str+"\n")
}

func PrintTo(str string) {
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

type MyFormatter struct {
	logrus.TextFormatter
}

func (f *MyFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	return []byte(entry.Message + "\n"), nil
}

func NewLogger(dir string) *logrus.Logger {
	dir = fileUtils.UpdateDir(dir)

	if Logger != nil {
		return Logger
	}

	Logger = logrus.New()
	Logger.Out = ioutil.Discard

	pathMap := lfshook.PathMap{
		logrus.WarnLevel:  dir + "trace.log",
		logrus.ErrorLevel: dir + "result.log",
	}

	Logger.Hooks.Add(lfshook.NewHook(
		pathMap,
		&MyFormatter{},
	))

	Logger.SetFormatter(&MyFormatter{})

	return Logger
}

func Screen(msg string) {
	PrintTo(msg)
}
func Trace(msg string) {
	Logger.Warnln(msg)
}
func Result(msg string) {
	Logger.Errorln(msg)
}

func InitLog(dir string) {
	Logger = NewLogger(dir)
}
