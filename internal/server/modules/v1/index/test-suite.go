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
	}
	return module.NewModule("/suites", handler)
}
