package action

import (
	"github.com/easysoft/zentaoatf/src/service/script"
	"github.com/easysoft/zentaoatf/src/service/zentao"
	"github.com/easysoft/zentaoatf/src/utils/config"
	"github.com/easysoft/zentaoatf/src/utils/const"
	"github.com/easysoft/zentaoatf/src/utils/i118"
	"github.com/easysoft/zentaoatf/src/utils/langUtils"
	"github.com/easysoft/zentaoatf/src/utils/log"
)

func GenerateScript(productId string, moduleId string, suiteId string, taskId string,
	independentFile bool, scriptLang string) {

	config := configUtils.ReadCurrConfig()
	if config.Url == "" || config.Account == "" || config.Password == "" {
		configUtils.ConfigForSet()
	}

	if (productId != "") || (moduleId != "" && productId != "") || suiteId != "" || taskId != "" {

	} else {
		configUtils.ConfigForCheckout(&productId, &moduleId, &suiteId, &taskId,
			&independentFile, &scriptLang)
	}

	ok := langUtils.CheckSupportLangages(scriptLang)
	if !ok {
		return
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
