package action

import (
	commConsts "github.com/easysoft/zentaoatf/internal/pkg/consts"
	configHelper "github.com/easysoft/zentaoatf/internal/pkg/helper/config"
	scriptHelper "github.com/easysoft/zentaoatf/internal/pkg/helper/script"
	zentaoHelper "github.com/easysoft/zentaoatf/internal/pkg/helper/zentao"
)

func CommitCases(files []string) {
	cases := scriptHelper.GetCaseByDirAndFile(files)

	config := configHelper.LoadByWorkspacePath(commConsts.ZtfDir)

	zentaoHelper.SyncToZentao(cases, config)
}
