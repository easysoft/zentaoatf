package action

import (
	commConsts "github.com/easysoft/zentaoatf/internal/comm/consts"
	scriptUtils "github.com/easysoft/zentaoatf/internal/comm/helper/script"
	i118Utils "github.com/easysoft/zentaoatf/internal/pkg/lib/i118"
	logUtils "github.com/easysoft/zentaoatf/internal/pkg/lib/log"
	serverConfig "github.com/easysoft/zentaoatf/internal/server/config"
)

func Extract(files []string) (err error) {
	serverConfig.InitExecLog(commConsts.ExecLogDir)

	err = scriptUtils.Extract(files)
	logUtils.Info(i118Utils.Sprintf("success_to_extract_step"))

	return
}
