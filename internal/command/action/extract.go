package action

import (
	commConsts "github.com/aaronchen2k/deeptest/internal/comm/consts"
	scriptUtils "github.com/aaronchen2k/deeptest/internal/comm/helper/script"
	serverConfig "github.com/aaronchen2k/deeptest/internal/server/config"
)

func Extract(files []string) error {
	serverConfig.InitExecLog(commConsts.ExecLogDir)
	return scriptUtils.Extract(files)
}
