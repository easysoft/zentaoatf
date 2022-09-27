package v1

import (
	"time"

	serverConfig "github.com/easysoft/zentaoatf/internal/server/config"
	"github.com/easysoft/zentaoatf/internal/server/core/module"
	"github.com/easysoft/zentaoatf/internal/server/modules/v1/index"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/rate"
)

type IndexModule struct {
	ConfigModule   *index.ConfigModule   `inject:""`
	SettingsModule *index.SettingsModule `inject:""`
	FileModule     *index.FileModule     `inject:""`

	ZentaoModule *index.ZentaoModule `inject:""`
	SiteModule   *index.SiteModule   `inject:""`
	ExecModule   *index.ExecModule   `inject:""`

	InterpreterModule *index.InterpreterModule `inject:""`
	WorkspaceModule   *index.WorkspaceModule   `inject:""`
	ProxyModule       *index.ProxyModule       `inject:""`
	ServerModule      *index.ServerModule      `inject:""`
	HeartBeatModule   *index.HeartBeatModule   `inject:""`

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
		m.ExecModule.Party(),
		m.InterpreterModule.Party(),
		m.WorkspaceModule.Party(),
		m.ProxyModule.Party(),
		m.ServerModule.Party(),
		m.HeartBeatModule.Party(),

		m.TestFilterModule.Party(),
		m.TestScriptModule.Party(),
		m.TestBugModule.Party(),
		m.TestResultModule.Party(),
	}
	return module.NewModule(serverConfig.ApiPath, handler, modules...)
}
