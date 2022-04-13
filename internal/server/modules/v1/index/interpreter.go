package index

import (
	"github.com/easysoft/zentaoatf/internal/server/core/module"
	"github.com/easysoft/zentaoatf/internal/server/middleware"
	"github.com/easysoft/zentaoatf/internal/server/modules/v1/controller"
	"github.com/kataras/iris/v12"
)

type InterpreterModule struct {
	InterpreterCtrl *controller.InterpreterCtrl `inject:""`
}

func NewInterpreterModule() *InterpreterModule {
	return &InterpreterModule{}
}

// Party 执行
func (m *InterpreterModule) Party() module.WebModule {
	handler := func(index iris.Party) {
		index.Use(middleware.InitCheck())

		index.Get("/", m.InterpreterCtrl.List).Name = "列表"
		index.Get("/{id:int}", m.InterpreterCtrl.Get).Name = "详情"
		index.Post("/", m.InterpreterCtrl.Create).Name = "新建"
		index.Put("/{id:int}", m.InterpreterCtrl.Update).Name = "更新"
		index.Delete("/{id:int}", m.InterpreterCtrl.Delete).Name = "删除"

		index.Get("/getLangSettings", m.InterpreterCtrl.GetLangSettings).Name = "获取语言配置项"
		index.Get("/getLangInterpreter", m.InterpreterCtrl.GetLangInterpreter).Name = "获取语言运行环境"

	}
	return module.NewModule("/interpreters", handler)
}
