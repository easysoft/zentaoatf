package index

import (
	"github.com/aaronchen2k/deeptest/internal/server/core/module"
	"github.com/aaronchen2k/deeptest/internal/server/middleware"
	"github.com/aaronchen2k/deeptest/internal/server/modules/v1/controller"
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

		index.Post("/", m.WorkspaceCtrl.Create).Name = "创建项目"
		index.Delete("/", m.WorkspaceCtrl.Delete).Name = "删除项目"

		index.Get("/getByUser", m.WorkspaceCtrl.GetByUser).Name = "获取用户参与的项目"
	}
	return module.NewModule("/workspaces", handler)
}
