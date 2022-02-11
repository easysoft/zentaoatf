package action

import (
	commConsts "github.com/aaronchen2k/deeptest/internal/comm/consts"
	"github.com/aaronchen2k/deeptest/internal/command"
	stringUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/string"
)

func CommitCases(files []string, actionModule *command.IndexModule) {
	var resultDir string
	if len(files) > 0 {
		resultDir = files[0]
	}

	actionModule.SyncService.SyncToZentao(resultDir, stringUtils.ParseInt(commConsts.ProductId))
}
