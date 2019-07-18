package config

import (
	"encoding/json"
	data2 "github.com/easysoft/zentaoatf"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"io/ioutil"
	"sync"
)

var printer *message.Printer

func GetI118(lang string) *message.Printer {
	en := "src/res/messages_en.json"
	zh := "src/res/messages_zh.json"

	var once sync.Once
	once.Do(func() {
		data, err := data2.Asset(en)
		if err == nil {
			InitResFromAsset(data)
		} else {
			InitRes(en)
		}

		if lang == "zh" {
			data, err := data2.Asset(zh)
			if err == nil {
				InitResFromAsset(data)
			} else {
				InitRes(zh)
			}

			printer = message.NewPrinter(language.SimplifiedChinese)
		} else {
			printer = message.NewPrinter(language.AmericanEnglish)
		}
	})
	return printer
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
