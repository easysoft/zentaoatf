package action

import (
	"github.com/easysoft/zentaoatf/src/service/script"
	"github.com/easysoft/zentaoatf/src/service/zentao"
	"github.com/easysoft/zentaoatf/src/utils/const"
	"github.com/easysoft/zentaoatf/src/utils/i118"
	"github.com/easysoft/zentaoatf/src/utils/lang"
	"github.com/easysoft/zentaoatf/src/utils/log"
	stdinUtils "github.com/easysoft/zentaoatf/src/utils/stdin"
	"github.com/fatih/color"
)

func Generate(productId string, moduleId string, suiteId string, taskId string,
	independentFile bool, scriptLang string) {

	CheckRequestConfig()

	if (productId != "") || (moduleId != "" && productId != "") || suiteId != "" || taskId != "" {

	} else {
		stdinUtils.InputForCheckout(&productId, &moduleId, &suiteId, &taskId,
			&independentFile, &scriptLang)
	}

	ok := langUtils.CheckSupportLangages(scriptLang)
	if !ok {
		return
	}

	cases := zentaoService.LoadTestCases(productId, moduleId, suiteId, taskId)

	if cases != nil && len(cases) > 0 {

		productId = cases[0].Product
		zentaoService.GetCaseModules(productId)

		count, err := scriptService.Generate(cases, scriptLang, independentFile)
		if err == nil {
			logUtils.PrintToCmd(i118Utils.I118Prt.Sprintf("success_to_generate", count, constant.ScriptDir)+"\n", -1)
		} else {
			logUtils.PrintToCmd(err.Error(), color.FgRed)
		}
	} else {
		logUtils.PrintToCmd(i118Utils.I118Prt.Sprintf("no_cases"), color.FgRed)
	}
}
