package command

import (
	commConsts "github.com/aaronchen2k/deeptest/internal/comm/consts"
	commandConfig "github.com/aaronchen2k/deeptest/internal/command/config"
	i118Utils "github.com/aaronchen2k/deeptest/internal/pkg/lib/i118"
	langUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/lang"
	serverConfig "github.com/aaronchen2k/deeptest/internal/server/config"
	serverLog "github.com/aaronchen2k/deeptest/internal/server/core/log"
)

func InitConfig() {
	serverConfig.Init()
	serverLog.Init()

	commandConfig.CheckConfigPermission()

	// screen size
	commandConfig.InitScreenSize()

	// internationalization
	i118Utils.Init(commConsts.Language, commConsts.AppServer)

	langUtils.GetExtToNameMap()
	return
}
