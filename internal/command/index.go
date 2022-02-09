package command

import (
	commConsts "github.com/aaronchen2k/deeptest/internal/comm/consts"
	commandConfig "github.com/aaronchen2k/deeptest/internal/command/config"
	i118Utils "github.com/aaronchen2k/deeptest/internal/pkg/lib/i118"
	langUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/lang"
	serverConfig "github.com/aaronchen2k/deeptest/internal/server/config"
	serverLog "github.com/aaronchen2k/deeptest/internal/server/core/log"
	"github.com/aaronchen2k/deeptest/internal/server/modules/v1/service"
	"github.com/facebookgo/inject"
	"github.com/sirupsen/logrus"
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

type IndexModule struct {
	ZtfCaseService *service.ZtfCaseService `inject:""`
}

func NewIndexModule() *IndexModule {
	return &IndexModule{}
}

func InjectModule() (indexModule *IndexModule) {
	var g inject.Graph
	indexModule = NewIndexModule()

	// inject objects
	if err := g.Provide(
		&inject.Object{Value: indexModule},
	); err != nil {
		logrus.Fatalf("provide usecase objects to the Graph: %v", err)
	}
	err := g.Populate()
	if err != nil {
		logrus.Fatalf("populate the incomplete Objects: %v", err)
	}
	return
}
