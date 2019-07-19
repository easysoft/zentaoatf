package utils

import (
	"encoding/json"
	"github.com/easysoft/zentaoatf/res"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"io/ioutil"
	"sync"
)

var I118Prt *message.Printer

func InitI118(lang string) {
	var once sync.Once
	once.Do(func() {
		isDebug := IsDebug()

		if isDebug {
			InitRes(EnRes)
		} else {
			data, _ := res.Asset(EnRes)
			InitResFromAsset(data)
		}

		if lang == "zh" {
			if isDebug {
				InitRes(ZhRes)
			} else {
				data, _ := res.Asset(ZhRes)
				InitResFromAsset(data)
			}

			I118Prt = message.NewPrinter(language.SimplifiedChinese)
		} else {
			I118Prt = message.NewPrinter(language.AmericanEnglish)
		}
	})
}

type I18n struct {
	Language string    `json:"language"`
	Messages []Message `json:"messages"`
}

type Message struct {
	Id          string `json:"id"`
	Message     string `json:"message,omitempty"`
	Translation string `json:"translation,omitempty"`
}

func Check(e error) {
	if e != nil {
		panic(e)
	}
}
func ReadI18nJson(file string) string {
	b, err := ioutil.ReadFile(file)
	Check(err)
	str := string(b)
	return str

}

func InitResFromAsset(bytes []byte) {
	var i18n I18n
	json.Unmarshal(bytes, &i18n)

	msgArr := i18n.Messages
	tag := language.MustParse(i18n.Language)

	for _, e := range msgArr {
		message.SetString(tag, e.Id, e.Translation)
	}
}
func InitRes(jsonPath string) {
	var i18n I18n
	str := ReadI18nJson(jsonPath)
	json.Unmarshal([]byte(str), &i18n)

	msgArr := i18n.Messages
	tag := language.MustParse(i18n.Language)

	for _, e := range msgArr {
		message.SetString(tag, e.Id, e.Translation)
	}
}
