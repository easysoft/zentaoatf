package index

import (
	"github.com/easysoft/zentaoatf/internal/server/core/module"
	"github.com/easysoft/zentaoatf/internal/server/middleware"
	"github.com/easysoft/zentaoatf/internal/server/modules/v1/controller"
	"github.com/kataras/iris/v12"
)

type SettingsModule struct {
	SettingsCtrl *controller.SettingsCtrl `inject:""`
}

func NewSettingsModule() *SettingsModule {
	return &SettingsModule{}
}

// Party 执行
func (m *SettingsModule) Party() module.WebModule {
	handler := func(index iris.Party) {
		index.Use(middleware.InitCheck())

		index.Get("/setLang", m.SettingsCtrl.SetLang).Name = "修改语言"
	}
	return module.NewModule("/settings", handler)
}
