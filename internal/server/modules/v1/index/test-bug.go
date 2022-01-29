package index

import (
	"github.com/aaronchen2k/deeptest/internal/server/core/module"
	"github.com/aaronchen2k/deeptest/internal/server/middleware"
	"github.com/aaronchen2k/deeptest/internal/server/modules/v1/controller"
	"github.com/kataras/iris/v12"
)

type TestBugModule struct {
	TestBugCtrl *controller.TestBugCtrl `inject:""`
}

func NewTestBugModule() *TestBugModule {
	return &TestBugModule{}
}

// Party 测试结果
func (m *TestBugModule) Party() module.WebModule {
	handler := func(index iris.Party) {
		index.Use(middleware.InitCheck())
		index.Post("/{seq:string}", m.TestBugCtrl.Submit).Name = "提交缺陷"
	}
	return module.NewModule("/bug", handler)
}
