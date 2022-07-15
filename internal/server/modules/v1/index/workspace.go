package index

import (
	"github.com/easysoft/zentaoatf/internal/server/core/module"
	"github.com/easysoft/zentaoatf/internal/server/middleware"
	"github.com/easysoft/zentaoatf/internal/server/modules/v1/controller"
	"github.com/kataras/iris/v12"
)

type WorkspaceModule struct {
	WorkspaceCtrl *controller.WorkspaceCtrl `inject:""`
}

func NewWorkspaceModule() *WorkspaceModule {
	return &WorkspaceModule{}
}

// Party 项目
func (m *WorkspaceModule) Party() module.WebModule {
	handler := func(index iris.Party) {
		index.Use(middleware.InitCheck())

		index.Get("/", m.WorkspaceCtrl.List).Name = "列表"
		index.Get("/listByProduct", m.WorkspaceCtrl.ListByProduct).Name = "列表"
		index.Get("/{id:int}", m.WorkspaceCtrl.Get).Name = "详情"
		index.Post("/", m.WorkspaceCtrl.Create).Name = "新建"
		index.Put("/{id:int}", m.WorkspaceCtrl.Update).Name = "更新"
		index.Delete("/{id:int}", m.WorkspaceCtrl.Delete).Name = "删除"
		index.Delete("/", m.WorkspaceCtrl.DeleteByPath).Name = "删除工作目录"
	}
	return module.NewModule("/workspaces", handler)
}
