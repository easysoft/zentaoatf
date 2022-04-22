package action

import (
	commConsts "github.com/easysoft/zentaoatf/internal/comm/consts"
	configHelper "github.com/easysoft/zentaoatf/internal/comm/helper/config"
	scriptHelper "github.com/easysoft/zentaoatf/internal/comm/helper/script"
	zentaoHelper "github.com/easysoft/zentaoatf/internal/comm/helper/zentao"
)

func CommitCases(files []string) {
	cases := scriptHelper.GetCaseByDirAndFile(files)

	config := configHelper.LoadByWorkspacePath(commConsts.WorkDir)

	zentaoHelper.SyncToZentao(cases, config)
}
