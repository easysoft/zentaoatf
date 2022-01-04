package index

import (
	"github.com/aaronchen2k/deeptest/internal/server/core/module"
	"github.com/aaronchen2k/deeptest/internal/server/middleware"
	"github.com/aaronchen2k/deeptest/internal/server/modules/v1/controller"
	"github.com/kataras/iris/v12"
)

type TestExecutionModule struct {
	TestExecutionCtrl *controller.TestExecutionCtrl `inject:""`
}

func NewTestExecutionModule() *TestExecutionModule {
	return &TestExecutionModule{}
}

// Party 执行
func (m *TestExecutionModule) Party() module.WebModule {
	handler := func(index iris.Party) {
		index.Use(middleware.InitCheck())

		index.Get("/", m.TestExecutionCtrl.List).Name = "执行列表"
		index.Get("/{id:uint}", m.TestExecutionCtrl.Get).Name = "执行详情"
		index.Post("/", m.TestExecutionCtrl.Create).Name = "创建执行"
		index.Put("/{id:uint}", m.TestExecutionCtrl.Update).Name = "更新执行"
		index.Delete("/{id:uint}", m.TestExecutionCtrl.Delete).Name = "删除执行"
	}
	return module.NewModule("/executions", handler)
}
