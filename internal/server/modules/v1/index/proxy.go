package index

import (
	"github.com/easysoft/zentaoatf/internal/server/core/module"
	"github.com/easysoft/zentaoatf/internal/server/middleware"
	"github.com/easysoft/zentaoatf/internal/server/modules/v1/controller"
	"github.com/kataras/iris/v12"
)

type ProxyModule struct {
	ProxyCtrl *controller.ProxyCtrl `inject:""`
}

func NewProxyModule() *ProxyModule {
	return &ProxyModule{}
}

// Party 执行
func (m *ProxyModule) Party() module.WebModule {
	handler := func(index iris.Party) {
		index.Use(middleware.InitCheck())

		index.Get("/", m.ProxyCtrl.List).Name = "列表"
		index.Get("/{id:int}/check", m.ProxyCtrl.Check).Name = "测试节点状态"
		index.Get("/{id:int}", m.ProxyCtrl.Get).Name = "详情"
		index.Post("/", m.ProxyCtrl.Create).Name = "新建"
		index.Put("/{id:int}", m.ProxyCtrl.Update).Name = "更新"
		index.Delete("/{id:int}", m.ProxyCtrl.Delete).Name = "删除"
	}
	return module.NewModule("/proxies", handler)
}
