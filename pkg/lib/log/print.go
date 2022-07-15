package logUtils

import (
	"encoding/json"
	"fmt"
	"github.com/easysoft/zentaoatf/pkg/consts"
	"strings"
	"unicode/utf8"

	"github.com/fatih/color"
	"go.uber.org/zap"
)

var LoggerStandard *zap.Logger
var LoggerExecConsole *zap.Logger

var LoggerExecFile *zap.Logger
var LoggerExecResult *zap.Logger

func Info(str string) {
	LoggerStandard.Info(str)
}
func Infof(str string, args ...interface{}) {
	msg := fmt.Sprintf(str, args...)
	LoggerStandard.Info(msg)
}
func Warn(str string) {
	LoggerStandard.Warn(str)
}
func Warnf(str string, args ...interface{}) {
	msg := fmt.Sprintf(str, args...)
	LoggerStandard.Warn(msg)
}
func Error(str string) {
	LoggerStandard.Error(str)
}
func Errorf(str string, args ...interface{}) {
	msg := fmt.Sprintf(str, args...)
	LoggerStandard.Error(msg)
}

func ExecConsole(clr color.Attribute, str string) {
	msg := str
	if clr > 0 {
		msg = color.New(clr).Sprint(msg)
	}

	LoggerExecConsole.Info(msg)
}
func ExecConsolef(clr color.Attribute, str string, args ...interface{}) {
	msg := fmt.Sprintf(str, args...)
	if clr > 0 {
		msg = color.New(clr).Sprint(msg)
	}

	LoggerExecConsole.Info(msg)
}

func ExecFile(str string) {
	LoggerExecFile.Info(str)
}
func ExecFilef(str string, args ...interface{}) {
	msg := fmt.Sprintf(str, args...)
	LoggerExecFile.Info(msg)
}

func ExecResult(str string) {
	LoggerExecResult.Info(str)
}
func ExecResultf(str string, args ...interface{}) {
	msg := fmt.Sprintf(str, args...)
	LoggerExecResult.Info(msg)
}

func PrintUnicode(str []byte) {
	msg := ConvertUnicode(str)
	LoggerStandard.Info(msg)
}

func ConvertUnicode(str []byte) string {
	var a interface{}

	temp := strings.Replace(string(str), "\\\\", "\\", -1)

	err := json.Unmarshal([]byte(temp), &a)

	var msg string
	if err == nil {
		bytes, _ := json.Marshal(a)
		msg = string(bytes)
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

func PrintToCmd(msg string, attr color.Attribute) {
	output := color.Output

	if attr == -1 {
		fmt.Fprint(output, msg+"\n")
	} else {
		clr := color.New(attr)
		clr.Fprint(output, msg+"\n")
	}
}

func PrintVersion(appVersion, buildTime, goVersion, gitHash string) {
	fmt.Printf("%s \n", appVersion)
	fmt.Printf("Build TimeStamp: %s \n", buildTime)
	fmt.Printf("GoLang Version: %s \n", goVersion)
	fmt.Printf("Git Commit Hash: %s \n", gitHash)
}
