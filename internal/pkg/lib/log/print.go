package logUtils

import (
	"encoding/json"
	"fmt"
	"github.com/aaronchen2k/deeptest/internal/pkg/consts"
	"go.uber.org/zap"
	"log"
	"strings"
	"unicode/utf8"
)

var Logger *zap.Logger

func Info(str string) {
	Logger.Info(str)
	log.Println(str)
}
func Infof(str string, args ...interface{}) {
	msg := fmt.Sprintf(str, args...)
	Logger.Info(msg)
	log.Println(msg)
}

func Warn(str string) {
	Logger.Warn(str)
	log.Println(str)
}
func Warnf(str string, args ...interface{}) {
	msg := fmt.Sprintf(str, args...)
	Logger.Warn(msg)
	log.Println(msg)
}

func Error(str string) {
	Logger.Error(str)
	log.Println(str)
}
func Errorf(str string, args ...interface{}) {
	msg := fmt.Sprintf(str, args...)
	Logger.Error(msg)
	log.Printf(msg + "\n")
}

func PrintUnicode(str []byte) {
	msg := ConvertUnicode(str)

	Logger.Info(msg)
	log.Print(msg)
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
