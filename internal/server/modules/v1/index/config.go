package index

import (
	"github.com/easysoft/zentaoatf/internal/server/core/module"
	"github.com/easysoft/zentaoatf/internal/server/middleware"
	"github.com/easysoft/zentaoatf/internal/server/modules/v1/controller"
	"github.com/kataras/iris/v12"
)

type ConfigModule struct {
	ConfigCtrl *controller.ConfigCtrl `inject:""`
}

func NewConfigModule() *ConfigModule {
	return &ConfigModule{}
}

// Party 执行
func (m *ConfigModule) Party() module.WebModule {
	handler := func(index iris.Party) {
		index.Use(middleware.InitCheck())

		index.Get("/setVerbose", m.ConfigCtrl.SetVerbose).Name = "设置日志级别"
	}
	return module.NewModule("/configs", handler)
}
