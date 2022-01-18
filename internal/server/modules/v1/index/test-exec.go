package index

import (
	"github.com/aaronchen2k/deeptest/internal/server/core/module"
	"github.com/aaronchen2k/deeptest/internal/server/middleware"
	"github.com/aaronchen2k/deeptest/internal/server/modules/v1/controller"
	"github.com/kataras/iris/v12"
)

type TestExecModule struct {
	TestExecCtrl *controller.TestExecCtrl `inject:""`
}

func NewTestExecutionModule() *TestExecModule {
	return &TestExecModule{}
}

// Party 执行
func (m *TestExecModule) Party() module.WebModule {
	handler := func(index iris.Party) {
		index.Use(middleware.InitCheck())

		index.Get("/", m.TestExecCtrl.List).Name = "执行列表"
		index.Get("/{id:uint}", m.TestExecCtrl.Get).Name = "执行详情"
		index.Delete("/{id:uint}", m.TestExecCtrl.Delete).Name = "删除执行"
	}
	return module.NewModule("/exec", handler)
}
