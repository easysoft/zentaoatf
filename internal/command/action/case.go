package action

import (
	commConsts "github.com/easysoft/zentaoatf/internal/comm/consts"
	configHelper "github.com/easysoft/zentaoatf/internal/comm/helper/config"
	zentaoHelper "github.com/easysoft/zentaoatf/internal/comm/helper/zentao"
	stringUtils "github.com/easysoft/zentaoatf/internal/pkg/lib/string"
)

func CommitCases(files []string) {
	var workspacePath string
	if len(files) > 0 {
		workspacePath = files[0]
	}

	config := configHelper.LoadByWorkspacePath(workspacePath)

	zentaoHelper.SyncToZentao(nil, workspacePath, stringUtils.ParseInt(commConsts.ProductId), config)
}
