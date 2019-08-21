package action

import (
	"github.com/easysoft/zentaoatf/src/service/script"
	zentaoService "github.com/easysoft/zentaoatf/src/service/zentao"
	"github.com/easysoft/zentaoatf/src/utils/common"
	"github.com/easysoft/zentaoatf/src/utils/config"
	constant "github.com/easysoft/zentaoatf/src/utils/const"
	i118Utils "github.com/easysoft/zentaoatf/src/utils/i118"
	logUtils "github.com/easysoft/zentaoatf/src/utils/log"
	"github.com/fatih/color"
)

func GenerateScript(url string, entityType string, entityVal string, langType string, independentFile bool,
	account string, password string) {

	LangMap := scriptService.GetSupportedScriptLang()
	langs := ""
	if LangMap[langType] == nil {
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
	cases, productIdInt, projectId, name := zentaoService.LoadTestCases(url, account, password, entityType, entityVal)
	if cases != nil {
		count, err := scriptService.Generate(cases, langType, independentFile)
		if err == nil {
			configUtils.SaveConfig("", url, entityType, entityVal,
				productIdInt, projectId, langType, independentFile,
				name, account, password)

			configUtils.UpdateWorkDirHistoryForGenerate()

			logUtils.PrintToCmd(i118Utils.I118Prt.Sprintf("success_to_generate", count, constant.ScriptDir) + "\n")
		} else {
			logUtils.PrintToCmd(err.Error())
		}
	}
}
