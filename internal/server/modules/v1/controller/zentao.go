package controller

import (
	commConsts "github.com/aaronchen2k/deeptest/internal/comm/consts"
	commDomain "github.com/aaronchen2k/deeptest/internal/comm/domain"
	zentaoHelper "github.com/aaronchen2k/deeptest/internal/comm/helper/zentao"
	serverDomain "github.com/aaronchen2k/deeptest/internal/server/modules/v1/domain"
	"github.com/aaronchen2k/deeptest/internal/server/modules/v1/service"
	"github.com/kataras/iris/v12"
)

type ZentaoCtrl struct {
	SiteService       *service.SiteService       `inject:""`
	TestScriptService *service.TestScriptService `inject:""`
	BaseCtrl
}

func NewZentaoCtrl() *ZentaoCtrl {
	return &ZentaoCtrl{}
}

func (c *ZentaoCtrl) GetProfile(ctx iris.Context) {
	workspacePath := ctx.URLParam("currWorkspace")
	if workspacePath == "" {
		ctx.JSON(c.SuccessResp(iris.Map{}))
		return
	}

	data, err := zentaoHelper.GetProfile(workspacePath)
	if err != nil {
		ctx.JSON(c.ErrResp(commConsts.BizErrZentaoRequest, err.Error()))
		return
	}

	ctx.JSON(c.SuccessResp(data))
}

func (c *ZentaoCtrl) ListSiteAndProduct(ctx iris.Context) {
	currSiteId, _ := ctx.URLParamInt("currSiteId")
	currProductId, _ := ctx.URLParamInt("currProductId")

	sites, currSite, _ := c.SiteService.LoadSites(currSiteId)
	products, currProduct, err := zentaoHelper.LoadSiteProduct(currSite, currProductId)

	if err != nil {
		ctx.JSON(c.ErrResp(commConsts.BizErrZentaoRequest, err.Error()))
		return
	}

	data := iris.Map{"sites": sites, "products": products,
		"currSite": currSite, "currProduct": currProduct}

	ctx.JSON(c.SuccessResp(data))
}

func (c *ZentaoCtrl) ListProduct(ctx iris.Context) {
	workspacePath := ctx.URLParam("currWorkspace")
	if workspacePath == "" {
		ctx.JSON(c.SuccessResp(make([]serverDomain.ZentaoProduct, 0)))
		return
	}

	data, err := zentaoHelper.ListProduct(workspacePath)
	if err != nil {
		ctx.JSON(c.ErrResp(commConsts.BizErrZentaoRequest, err.Error()))
		return
	}

	ctx.JSON(c.SuccessResp(data))
}

func (c *ZentaoCtrl) ListModule(ctx iris.Context) {
	workspacePath := ctx.URLParam("currWorkspace")
	productId, err := ctx.URLParamInt("productId")
	if err != nil {
		ctx.JSON(c.ErrResp(commConsts.ParamErr, err.Error()))
		return
	}

	data, err := zentaoHelper.ListModule(productId, workspacePath)
	if err != nil {
		ctx.JSON(c.ErrResp(commConsts.BizErrZentaoRequest, err.Error()))
		return
	}

	ctx.JSON(c.SuccessResp(data))
}

func (c *ZentaoCtrl) ListSuite(ctx iris.Context) {
	workspacePath := ctx.URLParam("currWorkspace")
	productId, err := ctx.URLParamInt("productId")
	if err != nil {
		ctx.JSON(c.ErrResp(commConsts.ParamErr, err.Error()))
		return
	}

	data, err := zentaoHelper.ListSuite(productId, workspacePath)
	if err != nil {
		ctx.JSON(c.ErrResp(commConsts.BizErrZentaoRequest, err.Error()))
		return
	}

	ctx.JSON(c.SuccessResp(data))
}

func (c *ZentaoCtrl) ListTask(ctx iris.Context) {
	workspacePath := ctx.URLParam("currWorkspace")
	productId, err := ctx.URLParamInt("productId")
	if err != nil {
		ctx.JSON(c.ErrResp(commConsts.ParamErr, err.Error()))
		return
	}

	data, err := zentaoHelper.ListTask(productId, workspacePath)
	if err != nil {
		ctx.JSON(c.ErrResp(commConsts.BizErrZentaoRequest, err.Error()))
		return
	}

	ctx.JSON(c.SuccessResp(data))
}

func (c *ZentaoCtrl) GetDataForBugSubmition(ctx iris.Context) {
	workspacePath := ctx.URLParam("currWorkspace")

	req := commDomain.FuncResult{}
	if err := ctx.ReadJSON(&req); err != nil {
		ctx.JSON(c.ErrResp(commConsts.ParamErr, err.Error()))
		return
	}

	fields, err := zentaoHelper.GetBugFiledOptions(req, workspacePath)
	if err != nil {
		ctx.JSON(c.ErrResp(commConsts.BizErrZentaoRequest, err.Error()))
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
