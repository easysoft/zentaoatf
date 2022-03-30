package index

import (
	"github.com/aaronchen2k/deeptest/internal/server/core/module"
	"github.com/aaronchen2k/deeptest/internal/server/middleware"
	"github.com/aaronchen2k/deeptest/internal/server/modules/v1/controller"
	"github.com/kataras/iris/v12"
)

type TestHistoryModule struct {
	TestHistoryCtrl *controller.TestHistoryCtrl `inject:""`
}

func NewTestHistoryModule() *TestHistoryModule {
	return &TestHistoryModule{}
}

// Party 执行
func (m *TestHistoryModule) Party() module.WebModule {
	handler := func(index iris.Party) {
		index.Use(middleware.InitCheck())

		index.Get("/", m.TestHistoryCtrl.List).Name = "执行列表"
		index.Get("/{workspaceId:int}/{seq:string}", m.TestHistoryCtrl.Get).Name = "执行详情"
		index.Delete("/", m.TestHistoryCtrl.Delete).Name = "删除执行"
	}
	return module.NewModule("/results", handler)
}
