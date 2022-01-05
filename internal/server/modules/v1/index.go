package v1

import (
	serverConfig "github.com/aaronchen2k/deeptest/internal/server/config"
	"github.com/aaronchen2k/deeptest/internal/server/core/module"
	"github.com/aaronchen2k/deeptest/internal/server/modules/v1/index"
	"time"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/rate"
)

type IndexModule struct {
	TestModule *index.TestModule `inject:""`
	FileModule *index.FileModule `inject:""`

	ProductModule *index.ProductModule `inject:""`
	ProjectModule *index.ProjectModule `inject:""`

	TestScriptModule    *index.TestScriptModule    `inject:""`
	TestExecutionModule *index.TestExecutionModule `inject:""`
	TestSuiteModule     *index.TestSuiteModule     `inject:""`
	TestSetModule       *index.TestSetModule       `inject:""`
}

func NewIndexModule() *IndexModule {
	return &IndexModule{}
}

// Party v1 模块
func (m *IndexModule) Party() module.WebModule {
	handler := func(v1 iris.Party) {
		if !serverConfig.CONFIG.Limit.Disable {
			limitV1 := rate.Limit(
				serverConfig.CONFIG.Limit.Limit,
				serverConfig.CONFIG.Limit.Burst,
				rate.PurgeEvery(time.Minute, 5*time.Minute))
			v1.Use(limitV1)
		}
	}
	modules := []module.WebModule{
		m.TestModule.Party(),
		m.FileModule.Party(),

		m.ProjectModule.Party(),
		m.TestScriptModule.Party(),
		m.TestExecutionModule.Party(),
	}
	return module.NewModule(serverConfig.ApiPath, handler, modules...)
}
