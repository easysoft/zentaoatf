package index

import (
	"github.com/aaronchen2k/deeptest/internal/server/core/module"
	"github.com/aaronchen2k/deeptest/internal/server/middleware"
	"github.com/aaronchen2k/deeptest/internal/server/modules/v1/controller"
	"github.com/kataras/iris/v12"
)

type UserModule struct {
	UserCtrl *controller.UserCtrl `inject:""`
}

func NewUserModule() *UserModule {
	return &UserModule{}
}

// Party 用户
func (m *UserModule) Party() module.WebModule {
	handler := func(index iris.Party) {
		index.Use(middleware.InitCheck(), middleware.JwtHandler(), middleware.OperationRecord(), middleware.Casbin())
		index.Get("/", m.UserCtrl.GetAllUsers).Name = "用户列表"
		index.Get("/{id:uint}", m.UserCtrl.GetUser).Name = "用户详情"
		index.Post("/", m.UserCtrl.CreateUser).Name = "创建用户"
		index.Post("/{id:uint}", m.UserCtrl.UpdateUser).Name = "编辑用户"
		index.Delete("/{id:uint}", m.UserCtrl.DeleteUser).Name = "删除用户"

		index.Get("/profile", m.UserCtrl.Profile).Name = "个人信息"
		index.Get("/message", m.UserCtrl.Message).Name = "消息"
		index.Get("/logout", m.UserCtrl.Logout).Name = "退出"
		index.Get("/clear", m.UserCtrl.Clear).Name = "清空 token"
		index.Post("/change_avatar", m.UserCtrl.ChangeAvatar).Name = "修改头像"
		// index.Get("/expire", controllers.Expire).Title = "刷新 token"
	}
	return module.NewModule("/users", handler)
}
