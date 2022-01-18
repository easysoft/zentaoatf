package logUtils

import (
	"encoding/json"
	"fmt"
	"github.com/aaronchen2k/deeptest/internal/pkg/consts"
	"go.uber.org/zap"
	"strings"
	"unicode/utf8"
)

var LoggerConsole *zap.Logger
var LoggerLog *zap.Logger
var LoggerResult *zap.Logger

func Info(str string) {
	LoggerConsole.Info(str)
}
func Infof(str string, args ...interface{}) {
	msg := fmt.Sprintf(str, args...)
	LoggerConsole.Info(msg)
}
func Warn(str string) {
	LoggerConsole.Warn(str)
}
func Warnf(str string, args ...interface{}) {
	msg := fmt.Sprintf(str, args...)
	LoggerConsole.Warn(msg)
}
func Error(str string) {
	LoggerConsole.Error(str)
}
func Errorf(str string, args ...interface{}) {
	msg := fmt.Sprintf(str, args...)
	LoggerConsole.Error(msg)
}

func Log(str string) {
	LoggerLog.Info(str)
}
func Logf(str string, args ...interface{}) {
	msg := fmt.Sprintf(str, args...)
	LoggerLog.Info(msg)
}
func Result(str string) {
	LoggerResult.Info(str)
}
func Resultf(str string, args ...interface{}) {
	msg := fmt.Sprintf(str, args...)
	LoggerResult.Info(msg)
}

func PrintUnicode(str []byte) {
	msg := ConvertUnicode(str)
	LoggerConsole.Info(msg)
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

func GetWholeLine(msg string, char string) string {
	prefixLen := (consts.ScreenWidth - utf8.RuneCountInString(msg) - 2) / 2
	if prefixLen <= 0 { // no width in debug mode
		prefixLen = 6
	}
	postfixLen := consts.ScreenWidth - utf8.RuneCountInString(msg) - 2 - prefixLen - 1
	if postfixLen <= 0 { // no width in debug mode
		postfixLen = 6
	}

	preFixStr := strings.Repeat(char, prefixLen)
	postFixStr := strings.Repeat(char, postfixLen)

	return fmt.Sprintf("%s %s %s", preFixStr, msg, postFixStr)
}
