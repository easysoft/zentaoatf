package index

import (
	"github.com/aaronchen2k/deeptest/internal/server/core/module"
	"github.com/aaronchen2k/deeptest/internal/server/middleware"
	"github.com/aaronchen2k/deeptest/internal/server/modules/v1/controller"
	"github.com/kataras/iris/v12"
)

type TestFilterModule struct {
	TestFilterCtrl *controller.TestFilterCtrl `inject:""`
}

func NewTestFilterModule() *TestFilterModule {
	return &TestFilterModule{}
}

// Party 脚本
func (m *TestFilterModule) Party() module.WebModule {
	handler := func(index iris.Party) {
		index.Use(middleware.InitCheck())

		index.Get("/listItems", m.TestFilterCtrl.ListItems).Name = "获取脚本过滤器的内容列表"
	}
	return module.NewModule("/filters", handler)
}
