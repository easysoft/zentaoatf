package action

import (
	commConsts "github.com/easysoft/zentaoatf/internal/comm/consts"
	stdinHelper "github.com/easysoft/zentaoatf/internal/comm/helper/stdin"
	commandConfig "github.com/easysoft/zentaoatf/internal/command/config"
)

func Set() {
	stdinHelper.InputForSet(commConsts.ZtfDir)
	commandConfig.PrintCurrConfig()
}
