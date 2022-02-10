package index

import (
	"github.com/aaronchen2k/deeptest/internal/server/core/module"
	"github.com/aaronchen2k/deeptest/internal/server/middleware"
	"github.com/aaronchen2k/deeptest/internal/server/modules/v1/controller"
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

		index.Post("/", m.TestResultCtrl.Submit).Name = "提交测试结果"
	}
	return module.NewModule("/result", handler)
}
