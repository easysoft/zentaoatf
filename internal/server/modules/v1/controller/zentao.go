package controller

import (
	commConsts "github.com/aaronchen2k/deeptest/internal/comm/consts"
	commDomain "github.com/aaronchen2k/deeptest/internal/comm/domain"
	zentaoHelper "github.com/aaronchen2k/deeptest/internal/comm/helper/zentao"
	serverDomain "github.com/aaronchen2k/deeptest/internal/server/modules/v1/domain"
	"github.com/kataras/iris/v12"
)

type ZentaoCtrl struct {
	BaseCtrl
}

func NewZentaoCtrl() *ZentaoCtrl {
	return &ZentaoCtrl{}
}

func (c *ZentaoCtrl) GetProfile(ctx iris.Context) {
	projectPath := ctx.URLParam("currProject")
	if projectPath == "" {
		ctx.JSON(c.SuccessResp(make([]serverDomain.ZentaoProduct, 0)))
		return
	}

	data, err := zentaoHelper.GetProfile(projectPath)
	if err != nil {
		ctx.JSON(c.ErrResp(commConsts.BizErrProjectConfig, err.Error()))
		return
	}

	ctx.JSON(c.SuccessResp(data))
}

func (c *ZentaoCtrl) ListProduct(ctx iris.Context) {
	projectPath := ctx.URLParam("currProject")
	if projectPath == "" {
		ctx.JSON(c.SuccessResp(make([]serverDomain.ZentaoProduct, 0)))
		return
	}

	data, err := zentaoHelper.ListProduct(projectPath)
	if err != nil {
		ctx.JSON(c.ErrResp(commConsts.BizErrProjectConfig, err.Error()))
		return
	}

	ctx.JSON(c.SuccessResp(data))
}

func (c *ZentaoCtrl) ListModule(ctx iris.Context) {
	projectPath := ctx.URLParam("currProject")
	productId, err := ctx.URLParamInt("productId")
	if err != nil {
		ctx.JSON(c.ErrResp(commConsts.ParamErr, err.Error()))
		return
	}

	data, err := zentaoHelper.ListModuleForCase(productId, projectPath)
	if err != nil {
		ctx.JSON(c.ErrResp(commConsts.CommErr, err.Error()))
		return
	}

	ctx.JSON(c.SuccessResp(data))
}

func (c *ZentaoCtrl) ListSuite(ctx iris.Context) {
	projectPath := ctx.URLParam("currProject")
	productId, err := ctx.URLParamInt("productId")
	if err != nil {
		ctx.JSON(c.ErrResp(commConsts.ParamErr, err.Error()))
		return
	}

	data, err := zentaoHelper.ListSuiteByProduct(productId, projectPath)
	if err != nil {
		ctx.JSON(c.ErrResp(commConsts.CommErr, err.Error()))
		return
	}

	ctx.JSON(c.SuccessResp(data))
}

func (c *ZentaoCtrl) ListTask(ctx iris.Context) {
	projectPath := ctx.URLParam("currProject")
	productId, err := ctx.URLParamInt("productId")
	if err != nil {
		ctx.JSON(c.ErrResp(commConsts.ParamErr, err.Error()))
		return
	}

	data, err := zentaoHelper.ListTaskByProduct(productId, projectPath)
	if err != nil {
		ctx.JSON(c.ErrResp(commConsts.CommErr, err.Error()))
		return
	}

	ctx.JSON(c.SuccessResp(data))
}

func (c *ZentaoCtrl) GetDataForBugSubmition(ctx iris.Context) {
	projectPath := ctx.URLParam("currProject")

	req := commDomain.FuncResult{}
	if err := ctx.ReadJSON(&req); err != nil {
		ctx.JSON(c.ErrResp(commConsts.ParamErr, err.Error()))
		return
	}

	fields, err := zentaoHelper.GetBugFiledOptions(req, projectPath)
	if err != nil {
		ctx.JSON(c.ErrResp(commConsts.CommErr, err.Error()))
		return
	}

	data := iris.Map{"fields": fields}

	ctx.JSON(c.SuccessResp(data))
}

func (c *ZentaoCtrl) ListLang(ctx iris.Context) {
	data, err := zentaoHelper.ListLang()
	if err != nil {
		ctx.JSON(c.ErrResp(commConsts.CommErr, err.Error()))
		return
	}

	ctx.JSON(c.SuccessResp(data))
}
