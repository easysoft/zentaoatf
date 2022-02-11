package action

import (
	commandConfig "github.com/aaronchen2k/deeptest/internal/command/config"
	stdinUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/stdin"
)

func Set() {
	stdinUtils.InputForSet()
	commandConfig.PrintCurrConfig()
}
