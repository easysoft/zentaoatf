package action

import (
	"github.com/easysoft/zentaoatf/src/service/script"
	"github.com/easysoft/zentaoatf/src/service/zentao"
	"github.com/easysoft/zentaoatf/src/utils/config"
	"github.com/easysoft/zentaoatf/src/utils/const"
	"github.com/easysoft/zentaoatf/src/utils/i118"
	"github.com/easysoft/zentaoatf/src/utils/log"
	"github.com/fatih/color"
)

func GenerateScript(productId string, moduleId string, suiteId string, taskId string,
	independentFile bool, scriptLang string) {

	LangMap := scriptService.GetSupportedScriptLang()
	langs := ""
	if LangMap[scriptLang] == nil {
		i := 0
		for lang, _ := range LangMap {
			if i > 0 {
				langs += ", "
			}
			langs += lang
			i++
		}
		logUtils.PrintToCmd(color.RedString(i118Utils.I118Prt.Sprintf("only_support_script_language", langs)) + "\n")
		return
	}

	config := configUtils.ReadCurrConfig()
	if config.Url == "" {
		configUtils.ConfigForCheckout()
	}

	cases := zentaoService.LoadTestCases(productId, moduleId, suiteId, taskId)
	productId = cases[0].Product
	zentaoService.GetCaseModules(productId)

	if cases != nil {
		count, err := scriptService.Generate(cases, scriptLang, independentFile)
		if err == nil {
			logUtils.PrintToCmd(i118Utils.I118Prt.Sprintf("success_to_generate", count, constant.ScriptDir) + "\n")
		} else {
			logUtils.PrintToCmd(err.Error())
		}
	}
}
