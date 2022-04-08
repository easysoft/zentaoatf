package action

import (
	commConsts "github.com/aaronchen2k/deeptest/internal/comm/consts"
	configUtils "github.com/aaronchen2k/deeptest/internal/comm/helper/config"
	zentaoHelper "github.com/aaronchen2k/deeptest/internal/comm/helper/zentao"
	stringUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/string"
)

func CommitCases(files []string) {
	var workspacePath string
	if len(files) > 0 {
		workspacePath = files[0]
	}

	config := configUtils.LoadByWorkspacePath(workspacePath)

	zentaoHelper.SyncToZentao(nil, workspacePath, stringUtils.ParseInt(commConsts.ProductId), config)
}
