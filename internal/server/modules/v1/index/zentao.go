package index

import (
	"github.com/aaronchen2k/deeptest/internal/server/core/module"
	"github.com/aaronchen2k/deeptest/internal/server/middleware"
	"github.com/aaronchen2k/deeptest/internal/server/modules/v1/controller"
	"github.com/kataras/iris/v12"
)

type ZentaoModule struct {
	ZentaoCtrl *controller.ZentaoCtrl `inject:""`
}

func NewZentaoModule() *ZentaoModule {
	return &ZentaoModule{}
}

// Party 产品
func (m *ZentaoModule) Party() module.WebModule {
	handler := func(index iris.Party) {
		index.Use(middleware.InitCheck())

		index.Get("/listLang", m.ZentaoCtrl.ListLang).Name = "脚本语言列表"
		index.Get("/listProduct", m.ZentaoCtrl.ListProduct).Name = "产品列表"
		index.Get("/listModule", m.ZentaoCtrl.ListModule).Name = "模块列表"
		index.Get("/listSuite", m.ZentaoCtrl.ListSuite).Name = "套件列表"
		index.Get("/listTask", m.ZentaoCtrl.ListTask).Name = "任务列表"
	}
	return module.NewModule("/zentao", handler)
}
