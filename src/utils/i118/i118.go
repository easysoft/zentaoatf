package i118Utils

import (
	"encoding/json"
	"github.com/easysoft/zentaoatf/res"
	"github.com/easysoft/zentaoatf/src/utils/common"
	constant "github.com/easysoft/zentaoatf/src/utils/const"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"io/ioutil"
	"sync"
)

var I118Prt *message.Printer

func InitI118(lang string) {
	var once sync.Once
	once.Do(func() {
		isRelease := commonUtils.IsRelease()

		if isRelease {
			data, _ := res.Asset(constant.EnRes)
			InitResFromAsset(data)
		} else {
			InitRes(constant.EnRes)
		}

		if lang == "zh" {
			if isRelease {
				data, _ := res.Asset(constant.ZhRes)
				InitResFromAsset(data)
			} else {
				InitRes(constant.ZhRes)
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

func InitResFromAsset(bytes []byte) {
	var i18n I18n
	json.Unmarshal(bytes, &i18n)

	msgArr := i18n.Messages
	tag := language.MustParse(i18n.Language)

	for _, e := range msgArr {
		message.SetString(tag, e.Id, e.Translation)
	}
}
