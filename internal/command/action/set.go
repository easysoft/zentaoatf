package action

import (
	stdinUtils "github.com/easysoft/zentaoatf/internal/comm/helper/stdin"
	commandConfig "github.com/easysoft/zentaoatf/internal/command/config"
)

func Set() {
	stdinUtils.InputForSet()
	commandConfig.PrintCurrConfig()
}
