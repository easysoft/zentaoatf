package index

import (
	"github.com/easysoft/zentaoatf/internal/server/core/module"
	"github.com/easysoft/zentaoatf/internal/server/middleware"
	"github.com/easysoft/zentaoatf/internal/server/modules/v1/controller"
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
		index.Get("/getProfile", m.ZentaoCtrl.GetProfile).Name = "脚本语言列表"

		index.Get("/listSiteAndProduct", m.ZentaoCtrl.ListSiteAndProduct).Name = "获取禅道站点及其下产品"

		index.Get("/listProduct", m.ZentaoCtrl.ListProduct).Name = "产品列表"
		index.Get("/listModule", m.ZentaoCtrl.ListModule).Name = "模块列表"
		index.Get("/listSuite", m.ZentaoCtrl.ListSuite).Name = "套件列表"
		index.Get("/listTask", m.ZentaoCtrl.ListTask).Name = "测试单列表"
		index.Get("/listCase", m.ZentaoCtrl.ListCase).Name = "用例列表"

		index.Get("/listBugFields", m.ZentaoCtrl.ListBugFields).Name = "获取缺陷属性数据"
	}
	return module.NewModule("/zentao", handler)
}
