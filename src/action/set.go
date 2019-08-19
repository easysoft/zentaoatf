package action

import (
	"github.com/easysoft/zentaoatf/src/utils/config"
	constant "github.com/easysoft/zentaoatf/src/utils/const"
	"github.com/fatih/color"
	"log"
)

func SetLanguage(lang string, dumb bool) {
	if lang == constant.LanguageEN || lang == constant.LanguageZH {
		configUtils.SetLanguage(lang, dumb)
	} else {
		log.Println(color.RedString("only %s or %s language is acceptable", constant.LanguageEN, constant.LanguageZH))
	}
}

func SetWorkDir(dir string, dumb bool) {
	configUtils.SetWorkDir(dir, dumb)
}
