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
		index.Use(middleware.InitCheck(), middleware.JwtHandler(), middleware.OperationRecord(), middleware.Casbin())

		index.Get("/", m.TestScriptCtrl.List).Name = "脚本列表"
		index.Get("/{id:uint}", m.TestScriptCtrl.Get).Name = "脚本详情"
		index.Post("/", m.TestScriptCtrl.Create).Name = "创建脚本"
		index.Put("/{id:uint}", m.TestScriptCtrl.Update).Name = "更新脚本"
		index.Delete("/{id:uint}", m.TestScriptCtrl.Delete).Name = "删除脚本"
	}
	return module.NewModule("/testScripts", handler)
}
