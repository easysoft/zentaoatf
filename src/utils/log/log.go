package logUtils

import (
	"fmt"
	fileUtils "github.com/easysoft/zentaoatf/src/utils/file"
	"github.com/easysoft/zentaoatf/src/utils/i118"
	"github.com/fatih/color"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"strings"
)

var Logger *logrus.Logger

func GetWholeLine(msg string, char string) string {
	//prefixLen := (vari.ScreenWidth - utf8.RuneCountInString(msg)) / 2
	//if prefixLen <= 0 { // no width in debug mode
	//	prefixLen = 6
	//}
	//postfixLen := vari.ScreenWidth - utf8.RuneCountInString(msg) - prefixLen
	//if postfixLen <= 0 { // no width in debug mode
	//	postfixLen = 6
	//}

	preFixStr := strings.Repeat(char, 6)
	postFixStr := strings.Repeat(char, 0)

	return fmt.Sprintf("%s%s%s", preFixStr, msg, postFixStr)
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
func TraceAndResult(msg string) {
	Logger.Warnln(msg)
	Logger.Errorln(msg)
}

func InitLog(dir string) {
	Logger = NewLogger(dir)
}

type MyFormatter struct {
	logrus.TextFormatter
}

func (f *MyFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	return []byte(entry.Message + "\n"), nil
}
