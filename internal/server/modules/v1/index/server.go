package index

import (
	"github.com/easysoft/zentaoatf/internal/server/core/module"
	"github.com/easysoft/zentaoatf/internal/server/middleware"
	"github.com/easysoft/zentaoatf/internal/server/modules/v1/controller"
	"github.com/kataras/iris/v12"
)

type ServerModule struct {
	ServerCtrl *controller.ServerCtrl `inject:""`
}

func NewServerModule() *ServerModule {
	return &ServerModule{}
}

// Party 执行
func (m *ServerModule) Party() module.WebModule {
	handler := func(index iris.Party) {
		index.Use(middleware.InitCheck())

		index.Get("/", m.ServerCtrl.List).Name = "列表"
		index.Get("/{id:int}", m.ServerCtrl.Get).Name = "详情"
		index.Post("/", m.ServerCtrl.Create).Name = "新建"
		index.Put("/{id:int}", m.ServerCtrl.Update).Name = "更新"
		index.Delete("/{id:int}", m.ServerCtrl.Delete).Name = "删除"
	}
	return module.NewModule("/servers", handler)
}
