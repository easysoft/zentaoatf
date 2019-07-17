package config

import (
	"encoding/json"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"io/ioutil"
	"sync"
)

var printer *message.Printer

func GetI118(lang string) *message.Printer {
	var once sync.Once
	once.Do(func() {
		InitRes("src/res/messages_en.json")
		if lang == "zh" {
			InitRes("src/res/messages_zh.json")
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
