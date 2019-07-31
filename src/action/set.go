package action

import (
	"github.com/easysoft/zentaoatf/src/utils"
	"github.com/fatih/color"
	"log"
)

func Set(param string, val string, dumb bool) {
	if param == "lang" {
		if val == utils.LanguageEN || val == utils.LanguageZH {
			utils.SetPreference(param, val, dumb)
		} else {
			log.Println(color.RedString("only %s or %s language is acceptable", utils.LanguageEN, utils.LanguageZH))
		}
	} else if param == "workDir" {
		utils.MkDirIfNeeded(val)
		utils.SetPreference("workDir", val, dumb)
	}
}

func Reset() {
	Set("ZENTAO_LANG", utils.LanguageDefault, true)
}
