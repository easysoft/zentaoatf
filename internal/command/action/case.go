package action

import (
	commConsts "github.com/aaronchen2k/deeptest/internal/comm/consts"
	configUtils "github.com/aaronchen2k/deeptest/internal/comm/helper/config"
	"github.com/aaronchen2k/deeptest/internal/command"
	stringUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/string"
)

func CommitCases(files []string, actionModule *command.IndexModule) {
	var workspacePath string
	if len(files) > 0 {
		workspacePath = files[0]
	}

	config := configUtils.LoadByWorkspacePath(workspacePath)

	actionModule.SyncService.SyncToZentao(nil, workspacePath, stringUtils.ParseInt(commConsts.ProductId), config)
}
