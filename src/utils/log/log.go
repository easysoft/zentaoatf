package logUtils

import (
	fileUtils "github.com/easysoft/zentaoatf/src/utils/file"
	"github.com/easysoft/zentaoatf/src/utils/i118"
	"github.com/easysoft/zentaoatf/src/utils/vari"
	"github.com/fatih/color"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"strings"
)

var Logger *logrus.Logger

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

func InitLogger() *logrus.Logger {
	vari.LogDir = fileUtils.GetLogDir()

	if Logger != nil {
		return Logger
	}

	Logger = logrus.New()
	Logger.Out = ioutil.Discard

	pathMap := lfshook.PathMap{
		logrus.WarnLevel:  vari.LogDir + "log.txt",
		logrus.ErrorLevel: vari.LogDir + "result.txt",
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
func Log(msg string) {
	Logger.Warnln(msg)
}
func Result(msg string) {
	Logger.Errorln(msg)
}

func ScreenAndResult(msg string) {
	Screen(msg)
	Result(msg)
}

type MyFormatter struct {
	logrus.TextFormatter
}

func (f *MyFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	return []byte(entry.Message + "\n"), nil
}
