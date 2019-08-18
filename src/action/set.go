package action

import (
	"github.com/easysoft/zentaoatf/src/utils/config"
	constant "github.com/easysoft/zentaoatf/src/utils/const"
	"github.com/easysoft/zentaoatf/src/utils/file"
	"github.com/fatih/color"
	"log"
)

func Set(param string, val string, dumb bool) {
	if param == "lang" {
		if val == constant.LanguageEN || val == constant.LanguageZH {
			configUtils.SetPreference(param, val, dumb)
		} else {
			log.Println(color.RedString("only %s or %s language is acceptable", constant.LanguageEN, constant.LanguageZH))
		}
	} else if param == "workDir" {
		fileUtils.MkDirIfNeeded(val)
		configUtils.SetPreference("workDir", val, dumb)
	}
}

//func Reset() {
//	Set("ZENTAO_LANG", constant.LanguageDefault, true)
//}
