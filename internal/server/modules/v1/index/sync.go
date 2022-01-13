package index

import (
	"github.com/aaronchen2k/deeptest/internal/server/core/module"
	"github.com/aaronchen2k/deeptest/internal/server/middleware"
	"github.com/aaronchen2k/deeptest/internal/server/modules/v1/controller"
	"github.com/kataras/iris/v12"
)

type SyncModule struct {
	SyncCtrl *controller.SyncCtrl `inject:""`
}

func NewSyncModule() *SyncModule {
	return &SyncModule{}
}

// Party 项目
func (m *SyncModule) Party() module.WebModule {
	handler := func(index iris.Party) {
		index.Use(middleware.InitCheck())
		index.Post("/syncFromZentao", m.SyncCtrl.SyncFromZentao).Name = "从禅道导出脚本"
		index.Post("/syncToZentao", m.SyncCtrl.SyncToZentao).Name = "更新脚本到禅道"
	}
	return module.NewModule("/sync", handler)
}
