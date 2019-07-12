package misc

import (
	"encoding/json"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"io/ioutil"
	"sync"
)

var printer *message.Printer
var once sync.Once

func GetInstance() *message.Printer {
	once.Do(func() {
		InitConfig("src/res/messages_zh.json")
		InitConfig("src/res/messages_en.json")
		printer = message.NewPrinter(language.SimplifiedChinese)
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

func InitConfig(jsonPath string) {
	var i18n I18n
	str := ReadI18nJson(jsonPath)
	json.Unmarshal([]byte(str), &i18n)
	//fmt.Println(i18n.Language)

	msaArry := i18n.Messages
	tag := language.MustParse(i18n.Language)
	// 以上代码和以下代码都是硬编码方式
	for _, e := range msaArry {
		//fmt.Println(e.Id+"\t"+e.Translation)
		message.SetString(tag, e.Id, e.Translation)

	}
}
