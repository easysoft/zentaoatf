package action

import (
	stdinUtils "github.com/aaronchen2k/deeptest/internal/comm/helper/stdin"
	commandConfig "github.com/aaronchen2k/deeptest/internal/command/config"
)

func Set() {
	stdinUtils.InputForSet()
	commandConfig.PrintCurrConfig()
}
