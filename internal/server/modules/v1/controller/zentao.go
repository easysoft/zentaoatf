package controller

import (
	"github.com/aaronchen2k/deeptest/internal/pkg/domain"
	zentaoUtils "github.com/aaronchen2k/deeptest/internal/server/modules/v1/utils/zentao"
	"github.com/kataras/iris/v12"
)

type ZentaoCtrl struct {
	BaseCtrl
}

func NewZentaoCtrl() *ZentaoCtrl {
	return &ZentaoCtrl{}
}

func (c *ZentaoCtrl) ListProduct(ctx iris.Context) {
	projectPath := ctx.URLParam("currProject")

	data, err := zentaoUtils.ListProduct(projectPath)
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

	data, err := zentaoUtils.ListModuleByProduct(productId, projectPath)
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

	data, err := zentaoUtils.ListSuiteByProduct(productId, projectPath)
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

	data, err := zentaoUtils.ListTaskByProduct(productId, projectPath)
	if err != nil {
		ctx.JSON(domain.Response{Code: domain.SystemErr.Code, Data: nil, Msg: err.Error()})
		return
	}

	ctx.JSON(domain.Response{Code: domain.NoErr.Code, Data: data, Msg: domain.NoErr.Msg})
}

func (c *ZentaoCtrl) ListLang(ctx iris.Context) {
	data, err := zentaoUtils.ListLang()
	if err != nil {
		ctx.JSON(domain.Response{Code: domain.SystemErr.Code, Data: nil, Msg: err.Error()})
		return
	}

	ctx.JSON(domain.Response{Code: domain.NoErr.Code, Data: data, Msg: domain.NoErr.Msg})
}
