package action

import (
	commConsts "github.com/easysoft/zentaoatf/internal/pkg/consts"
	commDomain "github.com/easysoft/zentaoatf/internal/pkg/domain"
	configHelper "github.com/easysoft/zentaoatf/internal/pkg/helper/config"
	scriptHelper "github.com/easysoft/zentaoatf/internal/pkg/helper/script"
	zentaoHelper "github.com/easysoft/zentaoatf/internal/pkg/helper/zentao"
	stdinUtils "github.com/easysoft/zentaoatf/pkg/lib/stdin"
	stringUtils "github.com/easysoft/zentaoatf/pkg/lib/string"
)

func CheckIn(files []string, noNeedConfirm, withCode bool) {
	cases := scriptHelper.GetCaseByDirAndFile(files)

	config := configHelper.LoadByWorkspacePath(commConsts.ZtfDir)

	zentaoHelper.CheckIn(cases, config, noNeedConfirm, withCode)
}

func Checkout(productId, moduleId, suiteId, taskId string, independentFile bool, scriptLang string) {
	if (productId != "" || moduleId != "" || suiteId != "" || taskId != "") && scriptLang != "" {
		//isReady = true
	} else {
		stdinUtils.InputForCheckout(&productId, &moduleId, &suiteId, &taskId,
			&independentFile, &scriptLang)
	}

	settings := commDomain.SyncSettings{
		ProductId:       stringUtils.ParseInt(productId),
		ModuleId:        stringUtils.ParseInt(moduleId),
		SuiteId:         stringUtils.ParseInt(suiteId),
		TaskId:          stringUtils.ParseInt(taskId),
		IndependentFile: independentFile,
		Lang:            scriptLang,
	}

	config := configHelper.LoadByWorkspacePath(commConsts.ZtfDir)

	zentaoHelper.Checkout(settings, config, commConsts.WorkDir)
}
