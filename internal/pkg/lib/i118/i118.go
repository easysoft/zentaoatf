package i118Utils

import (
	"encoding/json"
	"fmt"
	"github.com/aaronchen2k/deeptest/internal/pkg/lib/res"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"io/ioutil"
	"path/filepath"
)

var I118Prt *message.Printer

func Init(lang string, app string) {

	langRes := filepath.Join("res", app, lang, "messages.json")

	bytes, _ := resUtils.ReadRes(langRes)
	InitResFromAsset(bytes)

	if lang == "zh" {
		I118Prt = message.NewPrinter(language.SimplifiedChinese)
	} else {
		I118Prt = message.NewPrinter(language.AmericanEnglish)
	}
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

func Sprintf(key message.Reference, a ...interface{}) string {
	if I118Prt == nil {
		return fmt.Sprintf("%s, %#v", key.(string), a)
	} else {
		return I118Prt.Sprintf(key, a...)
	}
}
