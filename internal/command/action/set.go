package action

import (
	commandConfig "github.com/easysoft/zentaoatf/internal/command/config"
	commConsts "github.com/easysoft/zentaoatf/internal/pkg/consts"
	stdinHelper "github.com/easysoft/zentaoatf/internal/pkg/helper/stdin"
)

func Set() {
	stdinHelper.InputForSet(commConsts.ZtfDir)
	commandConfig.PrintCurrConfig()
}
