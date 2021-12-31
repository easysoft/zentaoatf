package index

import (
	"github.com/aaronchen2k/deeptest/internal/server/core/module"
	"github.com/aaronchen2k/deeptest/internal/server/middleware"
	"github.com/aaronchen2k/deeptest/internal/server/modules/v1/controller"
	"github.com/kataras/iris/v12"
)

type RoleModule struct {
	RoleCtrl *controller.RoleCtrl `inject:""`
}

func NewRoleModule() *RoleModule {
	return &RoleModule{}
}

// Party 角色模块
func (m *RoleModule) Party() module.WebModule {
	handler := func(index iris.Party) {
		index.Use(middleware.InitCheck(), middleware.JwtHandler(), middleware.OperationRecord(), middleware.Casbin())
		index.Get("/", m.RoleCtrl.GetAllRoles).Name = "角色列表"
		index.Get("/{id:uint}", m.RoleCtrl.GetRole).Name = "角色详情"
		index.Post("/", m.RoleCtrl.CreateRole).Name = "创建角色"
		index.Post("/{id:uint}", m.RoleCtrl.UpdateRole).Name = "编辑角色"
		index.Delete("/{id:uint}", m.RoleCtrl.DeleteRole).Name = "删除角色"
	}
	return module.NewModule("/roles", handler)
}
