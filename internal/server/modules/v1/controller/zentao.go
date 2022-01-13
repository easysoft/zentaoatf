package controller

import (
	"github.com/aaronchen2k/deeptest/internal/pkg/domain"
	"github.com/aaronchen2k/deeptest/internal/server/modules/v1/service"
	"github.com/kataras/iris/v12"
)

type ZentaoCtrl struct {
	ZentaoService *service.ZentaoService `inject:""`
	BaseCtrl
}

func NewZentaoCtrl() *ZentaoCtrl {
	return &ZentaoCtrl{}
}

func (c *ZentaoCtrl) ListProduct(ctx iris.Context) {
	projectPath := ctx.URLParam("currProject")

	data, err := c.ZentaoService.ListProduct(projectPath)
	if err != nil {
		ctx.JSON(domain.Response{Code: domain.SystemErr.Code, Data: nil, Msg: err.Error()})
		return
	}

	ctx.JSON(domain.Response{Code: domain.NoErr.Code, Data: data, Msg: domain.NoErr.Msg})
}

func (c *ZentaoCtrl) ListModule(ctx iris.Context) {
	projectPath := ctx.URLParam("currProject")
	productId, err := ctx.URLParamInt("productId")
	if err != nil {
		ctx.JSON(domain.Response{Code: domain.SystemErr.Code, Data: nil, Msg: err.Error()})
		return
	}

	data, err := c.ZentaoService.ListModuleByProduct(productId, projectPath)
	if err != nil {
		ctx.JSON(domain.Response{Code: domain.SystemErr.Code, Data: nil, Msg: err.Error()})
		return
	}

	ctx.JSON(domain.Response{Code: domain.NoErr.Code, Data: data, Msg: domain.NoErr.Msg})
}

func (c *ZentaoCtrl) ListSuite(ctx iris.Context) {
	projectPath := ctx.URLParam("currProject")
	productId, err := ctx.URLParamInt("productId")
	if err != nil {
		ctx.JSON(domain.Response{Code: domain.SystemErr.Code, Data: nil, Msg: err.Error()})
		return
	}

	data, err := c.ZentaoService.ListSuiteByProduct(productId, projectPath)
	if err != nil {
		ctx.JSON(domain.Response{Code: domain.SystemErr.Code, Data: nil, Msg: err.Error()})
		return
	}

	ctx.JSON(domain.Response{Code: domain.NoErr.Code, Data: data, Msg: domain.NoErr.Msg})
}

func (c *ZentaoCtrl) ListTask(ctx iris.Context) {
	projectPath := ctx.URLParam("currProject")
	productId, err := ctx.URLParamInt("productId")
	if err != nil {
		ctx.JSON(domain.Response{Code: domain.SystemErr.Code, Data: nil, Msg: err.Error()})
		return
	}

	data, err := c.ZentaoService.ListTaskByProduct(productId, projectPath)
	if err != nil {
		ctx.JSON(domain.Response{Code: domain.SystemErr.Code, Data: nil, Msg: err.Error()})
		return
	}

	ctx.JSON(domain.Response{Code: domain.NoErr.Code, Data: data, Msg: domain.NoErr.Msg})
}

func (c *ZentaoCtrl) ListLang(ctx iris.Context) {
	data, err := c.ZentaoService.ListLang()
	if err != nil {
		ctx.JSON(domain.Response{Code: domain.SystemErr.Code, Data: nil, Msg: err.Error()})
		return
	}

	ctx.JSON(domain.Response{Code: domain.NoErr.Code, Data: data, Msg: domain.NoErr.Msg})
}
