package logUtils

import (
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"log"
	"strings"
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
	var a interface{}

	temp := strings.Replace(string(str), "\\\\", "\\", -1)

	err := json.Unmarshal([]byte(temp), &a)

	var msg string
	if err == nil {
		msg = fmt.Sprint(a)
	} else {
		msg = temp
	}

	Logger.Info(msg)
	log.Print(msg)
}
