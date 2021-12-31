package index

import (
	"github.com/aaronchen2k/deeptest/internal/server/core/module"
	"github.com/aaronchen2k/deeptest/internal/server/middleware"
	"github.com/aaronchen2k/deeptest/internal/server/modules/v1/controller"
	"github.com/kataras/iris/v12"
)

type ProductModule struct {
	ProductCtrl *controller.ProductCtrl `inject:""`
}

func NewProductModule() *ProductModule {
	return &ProductModule{}
}

// Party 项目
func (m *ProductModule) Party() module.WebModule {
	handler := func(index iris.Party) {
		index.Use(middleware.InitCheck(), middleware.JwtHandler(), middleware.OperationRecord(), middleware.Casbin())
		index.Get("/", m.ProductCtrl.List).Name = "项目列表"
		index.Get("/{id:uint}", m.ProductCtrl.Get).Name = "项目详情"
		index.Post("/", m.ProductCtrl.Create).Name = "创建项目"
		index.Post("/{id:uint}", m.ProductCtrl.Update).Name = "编辑项目"
		index.Delete("/{id:uint}", m.ProductCtrl.Delete).Name = "删除项目"
	}
	return module.NewModule("/products", handler)
}
