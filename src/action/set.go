package action

import (
	"fmt"
	"github.com/easysoft/zentaoatf/src/config"
	"github.com/fatih/color"
)

func Set(param string, val string) {
	if param == "lang" {
		if val == config.LanguageEN || val == config.LanguageZH {
			config.Set(param, val)
		} else {
			fmt.Println(color.RedString("only %s or %s language is acceptable", config.LanguageEN, config.LanguageZH))
		}
	}
}

func Reset() {
	Set("ZENTAO_LANG", config.LanguageDefault)
}
