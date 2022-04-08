package action

import (
	commConsts "github.com/aaronchen2k/deeptest/internal/comm/consts"
	commDomain "github.com/aaronchen2k/deeptest/internal/comm/domain"
	configUtils "github.com/aaronchen2k/deeptest/internal/comm/helper/config"
	"github.com/aaronchen2k/deeptest/internal/command"
	stdinUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/stdin"
	stringUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/string"
)

func Checkout(productId, moduleId, suiteId, taskId string, independentFile bool, scriptLang string, actionModule *command.IndexModule) {
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

	config := configUtils.LoadByWorkspacePath(commConsts.WorkDir)

	actionModule.SyncService.SyncFromZentao(settings, config, commConsts.WorkDir)

}
