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

		index.Get("/", m.SiteCtrl.List).Name = "执行列表"
		index.Get("/{seq:string}", m.SiteCtrl.Get).Name = "执行详情"
		index.Delete("/{seq:string}", m.SiteCtrl.Delete).Name = "删除执行"
	}
	return module.NewModule("/sites", handler)
}
