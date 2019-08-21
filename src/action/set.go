package action

import (
	"github.com/easysoft/zentaoatf/src/utils/config"
	constant "github.com/easysoft/zentaoatf/src/utils/const"
	i118Utils "github.com/easysoft/zentaoatf/src/utils/i118"
	"github.com/fatih/color"
	"log"
)

func SetLanguage(lang string, dumb bool) {
	if lang == constant.LanguageEN || lang == constant.LanguageZH {
		configUtils.SetLanguage(lang, dumb)
	} else {
		log.Println(color.RedString(i118Utils.I118Prt.Sprintf("support_language", i118Utils.I118Prt.Sprintf(constant.LanguageEN),
			i118Utils.I118Prt.Sprintf(constant.LanguageZH))))
	}
}

func SetWorkDir(dir string, dumb bool) {
	configUtils.SetWorkDir(dir, dumb)
}
