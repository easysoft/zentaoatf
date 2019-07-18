package action

import (
	"fmt"
	"github.com/easysoft/zentaoatf/src/utils"
	"github.com/fatih/color"
)

func Set(param string, val string) {
	if param == "lang" {
		if val == utils.LanguageEN || val == utils.LanguageZH {
			utils.Set(param, val)
		} else {
			fmt.Println(color.RedString("only %s or %s language is acceptable", utils.LanguageEN, utils.LanguageZH))
		}
	}
}

func Reset() {
	Set("ZENTAO_LANG", utils.LanguageDefault)
}
