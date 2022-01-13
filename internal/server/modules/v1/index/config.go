package index

import (
	"github.com/aaronchen2k/deeptest/internal/server/core/module"
	"github.com/aaronchen2k/deeptest/internal/server/middleware"
	"github.com/aaronchen2k/deeptest/internal/server/modules/v1/controller"
	"github.com/kataras/iris/v12"
)

type ConfigModule struct {
	ConfigCtrl *controller.ConfigCtrl `inject:""`
}

func NewConfigModule() *ConfigModule {
	return &ConfigModule{}
}

// Party 项目
func (m *ConfigModule) Party() module.WebModule {
	handler := func(index iris.Party) {
		index.Use(middleware.InitCheck())
		index.Post("/saveConfig", m.ConfigCtrl.SaveConfig).Name = "保存项目配置"
	}
	return module.NewModule("/config", handler)
}
