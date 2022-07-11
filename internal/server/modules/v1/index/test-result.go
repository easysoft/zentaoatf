package index

import (
	"github.com/easysoft/zentaoatf/internal/server/core/module"
	"github.com/easysoft/zentaoatf/internal/server/middleware"
	"github.com/easysoft/zentaoatf/internal/server/modules/v1/controller"
	"github.com/kataras/iris/v12"
)

type TestResultModule struct {
	TestResultCtrl *controller.TestResultCtrl `inject:""`
}

func NewTestResultModule() *TestResultModule {
	return &TestResultModule{}
}

// Party 测试结果
func (m *TestResultModule) Party() module.WebModule {
	handler := func(index iris.Party) {
		index.Use(middleware.InitCheck())

		index.Get("/", m.TestResultCtrl.List).Name = "执行列表"
		index.Get("/latest", m.TestResultCtrl.GetLatest).Name = "获取最新一次执行详情"
		index.Get("/{workspaceId:int}/{seq:string}", m.TestResultCtrl.Get).Name = "执行详情"
		index.Get("/downloadLog", m.TestResultCtrl.DownloadLog).Name = "获取对应目录下的log zip"
		index.Get("/mvLog", m.TestResultCtrl.MvLog).Name = "从代理机下载log到本地"
		index.Delete("/", m.TestResultCtrl.Delete).Name = "删除执行"

		index.Post("/", m.TestResultCtrl.Submit).Name = "提交测试结果"
	}
	return module.NewModule("/results", handler)
}
