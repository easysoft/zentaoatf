package index

import (
	"github.com/aaronchen2k/deeptest/internal/server/core/module"
	"github.com/aaronchen2k/deeptest/internal/server/middleware"
	"github.com/aaronchen2k/deeptest/internal/server/modules/v1/controller"
	"github.com/kataras/iris/v12"
)

type TestCaseModule struct {
	TestCaseCtrl *controller.TestCaseCtrl `inject:""`
}

func NewTestCaseModule() *TestCaseModule {
	return &TestCaseModule{}
}

// Party 用例
func (m *TestCaseModule) Party() module.WebModule {
	handler := func(index iris.Party) {
		index.Use(middleware.InitCheck())
	}
	return module.NewModule("/testCases", handler)
}
