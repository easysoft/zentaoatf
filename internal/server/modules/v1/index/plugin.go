package index

import (
	"github.com/easysoft/zentaoatf/internal/server/core/module"
	"github.com/easysoft/zentaoatf/internal/server/middleware"
	"github.com/easysoft/zentaoatf/internal/server/modules/v1/controller"
	"github.com/kataras/iris/v12"
)

type PluginModule struct {
	PluginCtrl *controller.PluginCtrl `inject:""`
}

func NewPluginModule() *PluginModule {
	return &PluginModule{}
}

// Party 插件管理模块
func (m *PluginModule) Party() module.WebModule {
	handler := func(index iris.Party) {
		index.Use(middleware.InitCheck())

		index.Post("/install", m.PluginCtrl.Install).Name = "安装插件"
		index.Post("/uninstall", m.PluginCtrl.Uninstall).Name = "卸载插件"

		index.Post("/start", m.PluginCtrl.Start).Name = "启动插件"
		index.Post("/stop", m.PluginCtrl.Stop).Name = "停止插件"

		index.Post("/exec", m.PluginCtrl.Exec).Name = "执行扫描"
		index.Post("/cancel", m.PluginCtrl.Cancel).Name = "取消扫描"
	}

	return module.NewModule("/plugins", handler)
}
