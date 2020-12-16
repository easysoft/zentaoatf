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

	isReady := false
	if (productId != "" || moduleId != "" || suiteId != "" || taskId != "") && scriptLang != "" {
		isReady = true
	} else {
		stdinUtils.InputForCheckout(&productId, &moduleId, &suiteId, &taskId,
			&independentFile, &scriptLang)
	}

	ok := langUtils.CheckSupportLanguages(scriptLang)
	if !ok {
		return
	}

	cases, loginFail := zentaoService.LoadTestCases(productId, moduleId, suiteId, taskId)

	if cases != nil && len(cases) > 0 {
		productId = cases[0].Product

		// if isReady, no need to set below values

		// 1. target dir
		targetDir := "product" + productId + string(os.PathSeparator)
		if !isReady {
			targetDir = stdinUtils.GetInput("", targetDir, "where_to_store_script", targetDir)
		}
		targetDir = fileUtils.AbosutePath(targetDir)

		// 2. organize by module
		byModule := true
		if !isReady {
			stdinUtils.InputForBool(&byModule, byModule, "co_organize_by_module")
		}

		// 3. prefix
		prefix := ""
		if !isReady {
			prefix = stdinUtils.GetInput("[-_a-z0-9]*", prefix, "co_script_prefix", prefix)
		}

		count, err := scriptUtils.Generate(cases, scriptLang, independentFile, targetDir, byModule, prefix)
		if err == nil {
			logUtils.PrintTo(i118Utils.I118Prt.Sprintf("success_to_generate", count, targetDir) + "\n")
		} else {
			logUtils.PrintToWithColor(err.Error(), color.FgRed)
		}
	} else {
		if !loginFail {
			logUtils.PrintToWithColor(i118Utils.I118Prt.Sprintf("no_cases"), color.FgRed)
		}
	}
}
