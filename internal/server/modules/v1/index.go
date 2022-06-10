package v1

import (
	serverConfig "github.com/easysoft/zentaoatf/internal/server/config"
	"github.com/easysoft/zentaoatf/internal/server/core/module"
	"github.com/easysoft/zentaoatf/internal/server/modules/v1/index"
	"time"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/rate"
)

type IndexModule struct {
	ConfigModule   *index.ConfigModule   `inject:""`
	SettingsModule *index.SettingsModule `inject:""`
	FileModule     *index.FileModule     `inject:""`

	ZentaoModule *index.ZentaoModule `inject:""`
	SiteModule   *index.SiteModule   `inject:""`

	InterpreterModule *index.InterpreterModule `inject:""`
	SyncModule        *index.SyncModule        `inject:""`
	WorkspaceModule   *index.WorkspaceModule   `inject:""`

	TestFilterModule *index.TestFilterModule `inject:""`
	TestScriptModule *index.TestScriptModule `inject:""`
	TestResultModule *index.TestResultModule `inject:""`
	TestBugModule    *index.TestBugModule    `inject:""`
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
		m.ConfigModule.Party(),
		m.FileModule.Party(),
		m.SettingsModule.Party(),

		m.ZentaoModule.Party(),
		m.SiteModule.Party(),
		m.InterpreterModule.Party(),
		m.SyncModule.Party(),
		m.WorkspaceModule.Party(),

		m.TestFilterModule.Party(),
		m.TestScriptModule.Party(),
		m.TestBugModule.Party(),
		m.TestResultModule.Party(),
	}
	return module.NewModule(serverConfig.ApiPath, handler, modules...)
}
