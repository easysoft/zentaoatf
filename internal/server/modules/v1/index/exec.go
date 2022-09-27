package index

import (
	"github.com/easysoft/zentaoatf/internal/server/core/module"
	"github.com/easysoft/zentaoatf/internal/server/middleware"
	"github.com/easysoft/zentaoatf/internal/server/modules/v1/controller"
	"github.com/kataras/iris/v12"
)

type ExecModule struct {
	ExecCtrl *controller.ExecCtrl `inject:""`
}

func NewExecModule() *ExecModule {
	return &ExecModule{}
}

// Party 执行
func (m *ExecModule) Party() module.WebModule {
	handler := func(index iris.Party) {
		index.Use(middleware.InitCheck())

		index.Post("/start", m.ExecCtrl.Start).Name = "执行测试"
		index.Post("/stop", m.ExecCtrl.Stop).Name = "终止测试"
	}
	return module.NewModule("/exec", handler)
}
