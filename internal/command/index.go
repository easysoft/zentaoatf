package command

import (
	commConsts "github.com/aaronchen2k/deeptest/internal/comm/consts"
	commandConfig "github.com/aaronchen2k/deeptest/internal/command/config"
	i118Utils "github.com/aaronchen2k/deeptest/internal/pkg/lib/i118"
	langUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/lang"
	serverConfig "github.com/aaronchen2k/deeptest/internal/server/config"
	"github.com/aaronchen2k/deeptest/internal/server/modules/v1/service"
)

func InitConfig() {
	serverConfig.Init()
	serverConfig.InitLog()

	commandConfig.CheckConfigPermission()

	// screen size
	commandConfig.InitScreenSize()

	// internationalization
	i118Utils.Init(commConsts.Language, commConsts.AppServer)

	langUtils.GetExtToNameMap()

	commConsts.ComeFrom = "cmd"
	return
}

type IndexModule struct {
	ProjectService    *service.ProjectService    `inject:""`
	SyncService       *service.SyncService       `inject:""`
	TestResultService *service.TestResultService `inject:""`
}

func NewIndexModule() *IndexModule {
	return &IndexModule{}
}
