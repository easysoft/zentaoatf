package index

import (
	"github.com/aaronchen2k/deeptest/internal/server/core/module"
	"github.com/aaronchen2k/deeptest/internal/server/middleware"
	"github.com/aaronchen2k/deeptest/internal/server/modules/v1/controller"
	"github.com/kataras/iris/v12"
)

type SiteModule struct {
	SiteCtrl *controller.SiteCtrl `inject:""`
}

func NewSiteModule() *SiteModule {
	return &SiteModule{}
}

// Party 执行
func (m *SiteModule) Party() module.WebModule {
	handler := func(index iris.Party) {
		index.Use(middleware.InitCheck())

		index.Get("/", m.SiteCtrl.List).Name = "列表"
		index.Get("/{id:int}", m.SiteCtrl.Get).Name = "详情"
		index.Post("/{id:int}", m.SiteCtrl.Create).Name = "保存"
		index.Put("/{id:int}", m.SiteCtrl.Update).Name = "保存"
		index.Delete("/{id:int}", m.SiteCtrl.Delete).Name = "删除"
	}
	return module.NewModule("/sites", handler)
}
