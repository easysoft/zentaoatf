package index

import (
	"github.com/aaronchen2k/deeptest/internal/server/core/module"
	"github.com/aaronchen2k/deeptest/internal/server/middleware"
	"github.com/aaronchen2k/deeptest/internal/server/modules/v1/controller"
	"github.com/kataras/iris/v12"
)

type TestScriptModule struct {
	TestScriptCtrl *controller.TestScriptCtrl `inject:""`
}

func NewTestScriptModule() *TestScriptModule {
	return &TestScriptModule{}
}

// Party 脚本
func (m *TestScriptModule) Party() module.WebModule {
	handler := func(index iris.Party) {
		index.Use(middleware.InitCheck())

		index.Get("/list", m.TestScriptCtrl.List).Name = "脚本列表"
		index.Get("/get", m.TestScriptCtrl.Get).Name = "脚本详情"
		index.Get("/loadChildren", m.TestScriptCtrl.LoadCodeChildren).Name = "脚本详情"

		index.Put("/updateCode", m.TestScriptCtrl.UpdateCode).Name = "抽取脚本"
		index.Get("/extract", m.TestScriptCtrl.Extract).Name = "抽取脚本"

		index.Post("/syncFromZentao", m.TestScriptCtrl.SyncFromZentao).Name = "从禅道导出脚本"
		index.Post("/syncToZentao", m.TestScriptCtrl.SyncToZentao).Name = "更新脚本到禅道"
		index.Get("/getCaseIdsFromReport", m.TestScriptCtrl.GetCaseIdsFromReport).Name = "获取报告中的用例编号列表"
	}
	return module.NewModule("/scripts", handler)
}
