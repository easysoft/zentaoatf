package index

import (
	"github.com/easysoft/zentaoatf/internal/server/core/module"
	"github.com/easysoft/zentaoatf/internal/server/middleware"
	"github.com/easysoft/zentaoatf/internal/server/modules/v1/controller"
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
		index.Get("/loadChildren", m.TestScriptCtrl.LoadCodeChildren).Name = "加载子项"
		index.Get("/get", m.TestScriptCtrl.Get).Name = "脚本详情"
		index.Post("/", m.TestScriptCtrl.Create).Name = "创建脚本"
		index.Delete("/", m.TestScriptCtrl.Delete).Name = "删除脚本"

		index.Put("/updateCode", m.TestScriptCtrl.UpdateCode).Name = "更新脚本"
		index.Put("/updateName", m.TestScriptCtrl.UpdateName).Name = "重命令脚本"
		index.Put("/paste", m.TestScriptCtrl.Paste).Name = "粘贴脚本"
		index.Put("/move", m.TestScriptCtrl.Move).Name = "移动脚本"
		index.Put("/rename", m.TestScriptCtrl.Rename).Name = "重命名脚本"
		index.Get("/extract", m.TestScriptCtrl.Extract).Name = "抽取脚本"

		index.Post("/syncFromZentao", m.TestScriptCtrl.SyncFromZentao).Name = "从禅道导出脚本"
		index.Post("/syncToZentao", m.TestScriptCtrl.SyncToZentao).Name = "更新脚本到禅道"
		index.Get("/getCaseIdsFromReport", m.TestScriptCtrl.GetCaseIdsFromReport).Name = "获取报告中的用例编号列表"
	}
	return module.NewModule("/scripts", handler)
}
