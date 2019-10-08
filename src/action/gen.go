package action

import (
	"github.com/easysoft/zentaoatf/src/service/script"
	"github.com/easysoft/zentaoatf/src/service/zentao"
	configUtils "github.com/easysoft/zentaoatf/src/utils/config"
	fileUtils "github.com/easysoft/zentaoatf/src/utils/file"
	"github.com/easysoft/zentaoatf/src/utils/i118"
	"github.com/easysoft/zentaoatf/src/utils/lang"
	"github.com/easysoft/zentaoatf/src/utils/log"
	stdinUtils "github.com/easysoft/zentaoatf/src/utils/stdin"
	"github.com/fatih/color"
	"os"
)

func Generate(productId string, moduleId string, suiteId string, taskId string, independentFile bool, scriptLang string) {
	configUtils.CheckRequestConfig()

	if ((productId != "") || (moduleId != "" && productId != "") || suiteId != "" || taskId != "") && scriptLang != "" {
		// ready
	} else {
		stdinUtils.InputForCheckout(&productId, &moduleId, &suiteId, &taskId,
			&independentFile, &scriptLang)
	}

	ok := langUtils.CheckSupportLanguages(scriptLang)
	if !ok {
		return
	}

	cases := zentaoService.LoadTestCases(productId, moduleId, suiteId, taskId)

	if cases != nil && len(cases) > 0 {
		productId = cases[0].Product
		//zentaoService.GetCaseModules(productId)

		// target dir
		targetDirDft := "product" + productId + string(os.PathSeparator) // fileUtils.GetCurrDir()
		targetDir := stdinUtils.GetInput("", targetDirDft, "where_to_store_script", targetDirDft)
		targetDir = fileUtils.AbosutePath(targetDir)

		// organize by module
		var byModule bool
		stdinUtils.InputForBool(&byModule, true, "co_organize_by_module")

		// prefix
		prefix := stdinUtils.GetInput("[-_a-z0-9]*", "", "co_script_prefix", "")

		count, err := scriptUtils.Generate(cases, scriptLang, independentFile, targetDir, byModule, prefix)
		if err == nil {
			logUtils.PrintToCmd(i118Utils.I118Prt.Sprintf("success_to_generate", count, targetDir)+"\n", -1)
		} else {
			logUtils.PrintToCmd(err.Error(), color.FgRed)
		}
	} else {
		logUtils.PrintToCmd(i118Utils.I118Prt.Sprintf("no_cases"), color.FgRed)
	}
}
