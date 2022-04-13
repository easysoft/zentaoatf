package action

import (
	commConsts "github.com/easysoft/zentaoatf/internal/comm/consts"
	commDomain "github.com/easysoft/zentaoatf/internal/comm/domain"
	configHelper "github.com/easysoft/zentaoatf/internal/comm/helper/config"
	zentaoHelper "github.com/easysoft/zentaoatf/internal/comm/helper/zentao"
	stdinUtils "github.com/easysoft/zentaoatf/internal/pkg/lib/stdin"
	stringUtils "github.com/easysoft/zentaoatf/internal/pkg/lib/string"
)

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

	config := configHelper.LoadByWorkspacePath(commConsts.WorkDir)

	zentaoHelper.SyncFromZentao(settings, config, commConsts.WorkDir)
}
