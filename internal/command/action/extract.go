package action

import (
	commConsts "github.com/easysoft/zentaoatf/internal/pkg/consts"
	scriptHelper "github.com/easysoft/zentaoatf/internal/pkg/helper/script"
	serverConfig "github.com/easysoft/zentaoatf/internal/server/config"
	i118Utils "github.com/easysoft/zentaoatf/pkg/lib/i118"
	logUtils "github.com/easysoft/zentaoatf/pkg/lib/log"
)

func Extract(files []string) (err error) {
	serverConfig.InitExecLog(commConsts.ExecLogDir)

	_, err = scriptHelper.Extract(files)
	logUtils.Info(i118Utils.Sprintf("success_to_extract_step"))

	return
}
