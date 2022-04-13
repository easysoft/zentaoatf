package action

import (
	stdinHelper "github.com/easysoft/zentaoatf/internal/comm/helper/stdin"
	commandConfig "github.com/easysoft/zentaoatf/internal/command/config"
)

func Set() {
	stdinHelper.InputForSet()
	commandConfig.PrintCurrConfig()
}
