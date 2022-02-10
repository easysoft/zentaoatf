package controller

import (
	commConsts "github.com/aaronchen2k/deeptest/internal/comm/consts"
	commDomain "github.com/aaronchen2k/deeptest/internal/comm/domain"
	zentaoUtils "github.com/aaronchen2k/deeptest/internal/comm/helper/zentao"
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
		ctx.JSON(c.ErrResp(commConsts.Failure, err.Error()))
		return
	}

	ctx.JSON(c.SuccessResp(data))
}

func (c *ZentaoCtrl) ListModule(ctx iris.Context) {
	projectPath := ctx.URLParam("currProject")
	productId, err := ctx.URLParamInt("productId")
	if err != nil {
		c.ErrResp(commConsts.ParamErr, err.Error())
		return
	}

	data, err := zentaoUtils.ListModuleForCase(productId, projectPath)
	if err != nil {
		ctx.JSON(c.ErrResp(commConsts.Failure, err.Error()))
		return
	}

	ctx.JSON(c.SuccessResp(data))
}

func (c *ZentaoCtrl) ListSuite(ctx iris.Context) {
	projectPath := ctx.URLParam("currProject")
	productId, err := ctx.URLParamInt("productId")
	if err != nil {
		c.ErrResp(commConsts.ParamErr, err.Error())
		return
	}

	data, err := zentaoUtils.ListSuiteByProduct(productId, projectPath)
	if err != nil {
		ctx.JSON(c.ErrResp(commConsts.Failure, err.Error()))
		return
	}

	ctx.JSON(c.SuccessResp(data))
}

func (c *ZentaoCtrl) ListTask(ctx iris.Context) {
	projectPath := ctx.URLParam("currProject")
	productId, err := ctx.URLParamInt("productId")
	if err != nil {
		c.ErrResp(commConsts.ParamErr, err.Error())
		return
	}

	data, err := zentaoUtils.ListTaskByProduct(productId, projectPath)
	if err != nil {
		ctx.JSON(c.ErrResp(commConsts.Failure, err.Error()))
		return
	}

	ctx.JSON(c.SuccessResp(data))
}

func (c *ZentaoCtrl) GetDataForBugSubmition(ctx iris.Context) {
	projectPath := ctx.URLParam("currProject")

	req := commDomain.FuncResult{}
	if err := ctx.ReadJSON(&req); err != nil {
		c.ErrResp(commConsts.ParamErr, err.Error())
		return
	}

	fields, err := zentaoUtils.GetBugFiledOptions(req, projectPath)
	if err != nil {
		ctx.JSON(c.ErrResp(commConsts.Failure, err.Error()))
		return
	}

	data := iris.Map{"fields": fields}

	ctx.JSON(c.SuccessResp(data))
}

func (c *ZentaoCtrl) ListLang(ctx iris.Context) {
	data, err := zentaoUtils.ListLang()
	if err != nil {
		ctx.JSON(c.ErrResp(commConsts.Failure, err.Error()))
		return
	}

	ctx.JSON(c.SuccessResp(data))
}
