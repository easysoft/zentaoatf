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
		index.Use(middleware.InitCheck())
	}
	return module.NewModule("/sets", handler)
}
