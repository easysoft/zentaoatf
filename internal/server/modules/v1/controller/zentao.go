package controller

import (
	commDomain "github.com/aaronchen2k/deeptest/internal/comm/domain"
	"github.com/aaronchen2k/deeptest/internal/pkg/domain"
	logUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/log"
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
		ctx.JSON(domain.Response{Code: domain.SystemErr.Code, Data: nil, Msg: "获取禅道产品失败"})
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

	data, err := zentaoUtils.ListModuleForCase(productId, projectPath)
	if err != nil {
		ctx.JSON(domain.Response{Code: domain.SystemErr.Code, Data: nil, Msg: "获取禅道模块失败"})
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
		ctx.JSON(domain.Response{Code: domain.SystemErr.Code, Data: nil, Msg: "获取禅道套件失败"})
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
		ctx.JSON(domain.Response{Code: domain.SystemErr.Code, Data: nil, Msg: "获取禅道任务失败"})
		return
	}

	ctx.JSON(domain.Response{Code: domain.NoErr.Code, Data: data, Msg: domain.NoErr.Msg})
}

func (c *ZentaoCtrl) GetDataForBugSubmition(ctx iris.Context) {
	projectPath := ctx.URLParam("currProject")

	req := commDomain.FuncResult{}
	if err := ctx.ReadJSON(&req); err != nil {
		logUtils.Errorf("参数验证失败 %s", err.Error())
		ctx.JSON(domain.Response{Code: domain.SystemErr.Code, Data: nil, Msg: err.Error()})
		return
	}

	steps, ids, fields, err := zentaoUtils.GetBugFiledOptions(req, projectPath)
	if err != nil {
		ctx.JSON(domain.Response{Code: domain.SystemErr.Code, Data: nil, Msg: "获取禅道缺陷属性失败"})
		return
	}

	data := iris.Map{"steps": steps, "ids": ids, "fields": fields}

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
