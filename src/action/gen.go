package action

import (
	"github.com/easysoft/zentaoatf/src/service/script"
	zentaoService "github.com/easysoft/zentaoatf/src/service/zentao"
	"github.com/easysoft/zentaoatf/src/utils/common"
	constant "github.com/easysoft/zentaoatf/src/utils/const"
	i118Utils "github.com/easysoft/zentaoatf/src/utils/i118"
	logUtils "github.com/easysoft/zentaoatf/src/utils/log"
	"github.com/fatih/color"
)

func GenerateScript(url string, account string, password string,
	productId string, moduleId string, suiteId string, taskId string,
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

	url = commonUtils.UpdateUrl(url)
	cases := zentaoService.LoadTestCases(url, account, password, productId, moduleId, suiteId, taskId)

	if cases != nil {
		count, err := scriptService.Generate(cases, scriptLang, independentFile)
		if err == nil {
			//configUtils.SaveConfig("", url, entityType, entityVal,
			//	productIdInt, projectId, scriptLang, independentFile,
			//	name, account, password)

			//configUtils.UpdateWorkDirHistoryForGenerate()

			logUtils.PrintToCmd(i118Utils.I118Prt.Sprintf("success_to_generate", count, constant.ScriptDir) + "\n")
		} else {
			logUtils.PrintToCmd(err.Error())
		}
	}
}
