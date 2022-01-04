package index

import (
	"github.com/aaronchen2k/deeptest/internal/server/core/module"
	"github.com/aaronchen2k/deeptest/internal/server/middleware"
	"github.com/aaronchen2k/deeptest/internal/server/modules/v1/controller"
	"github.com/kataras/iris/v12"
)

type TestSuiteModule struct {
	TestSuiteCtrl *controller.TestSuiteCtrl `inject:""`
}

func NewTestSuiteModule() *TestSuiteModule {
	return &TestSuiteModule{}
}

// Party 套件
func (m *TestSuiteModule) Party() module.WebModule {
	handler := func(index iris.Party) {
		index.Use(middleware.InitCheck())
		index.Get("/", m.TestSuiteCtrl.List).Name = "套件查询"
		index.Get("/{id:uint}", m.TestSuiteCtrl.Get).Name = "套件详情"
		index.Post("/", m.TestSuiteCtrl.Create).Name = "创建套件"
		index.Post("/{id:uint}", m.TestSuiteCtrl.Update).Name = "编辑套件"
		index.Delete("/{id:uint}", m.TestSuiteCtrl.Delete).Name = "删除套件"
	}
	return module.NewModule("/suites", handler)
}
