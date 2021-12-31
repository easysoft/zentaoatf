package index

import (
	"github.com/aaronchen2k/deeptest/internal/server/core/module"
	"github.com/aaronchen2k/deeptest/internal/server/middleware"
	"github.com/aaronchen2k/deeptest/internal/server/modules/v1/controller"
	"github.com/kataras/iris/v12"
)

type TestSetModule struct {
	TestSetCtrl *controller.TestSetCtrl `inject:""`
}

func NewTestSetModule() *TestSetModule {
	return &TestSetModule{}
}

// Party 测试集
func (m *TestSetModule) Party() module.WebModule {
	handler := func(index iris.Party) {
		index.Use(middleware.InitCheck(), middleware.JwtHandler(), middleware.OperationRecord(), middleware.Casbin())
		index.Get("/", m.TestSetCtrl.List).Name = "测试集查询"
		index.Get("/{id:uint}", m.TestSetCtrl.Get).Name = "测试集详情"
		index.Post("/", m.TestSetCtrl.Create).Name = "创建测试集"
		index.Post("/{id:uint}", m.TestSetCtrl.Update).Name = "编辑测试集"
		index.Delete("/{id:uint}", m.TestSetCtrl.Delete).Name = "删除测试集"
	}
	return module.NewModule("/testSets", handler)
}
