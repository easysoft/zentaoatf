package action

import (
	commConsts "github.com/aaronchen2k/deeptest/internal/comm/consts"
	scriptUtils "github.com/aaronchen2k/deeptest/internal/comm/helper/script"
	i118Utils "github.com/aaronchen2k/deeptest/internal/pkg/lib/i118"
	logUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/log"
	serverConfig "github.com/aaronchen2k/deeptest/internal/server/config"
)

func Extract(files []string) (err error) {
	serverConfig.InitExecLog(commConsts.ExecLogDir)

	err = scriptUtils.Extract(files)
	logUtils.Info(i118Utils.Sprintf("success_to_extract_step"))

	return
}
