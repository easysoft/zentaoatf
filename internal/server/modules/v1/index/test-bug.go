package index

import (
	"github.com/easysoft/zentaoatf/internal/server/core/module"
	"github.com/easysoft/zentaoatf/internal/server/middleware"
	"github.com/easysoft/zentaoatf/internal/server/modules/v1/controller"
	"github.com/kataras/iris/v12"
)

type TestBugModule struct {
	TestBugCtrl *controller.TestBugCtrl `inject:""`
}

func NewTestBugModule() *TestBugModule {
	return &TestBugModule{}
}

// Party 测试结果
func (m *TestBugModule) Party() module.WebModule {
	handler := func(index iris.Party) {
		index.Use(middleware.InitCheck())

		index.Post("/prepareBugData", m.TestBugCtrl.PrepareBugData).Name = "获取缺陷步骤"
		index.Post("/", m.TestBugCtrl.Submit).Name = "提交缺陷"
		index.Get("/", m.TestBugCtrl.LoadBugs).Name = "查询产品下所有bug"
	}
	return module.NewModule("/bugs", handler)
}
