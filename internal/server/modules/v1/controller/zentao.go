package controller

import (
	commConsts "github.com/easysoft/zentaoatf/internal/comm/consts"
	configHelper "github.com/easysoft/zentaoatf/internal/comm/helper/config"
	zentaoHelper "github.com/easysoft/zentaoatf/internal/comm/helper/zentao"
	serverDomain "github.com/easysoft/zentaoatf/internal/server/modules/v1/domain"
	"github.com/easysoft/zentaoatf/internal/server/modules/v1/service"
	"github.com/kataras/iris/v12"
)

type ZentaoCtrl struct {
	SiteService    *service.SiteService    `inject:""`
	TestBugService *service.TestBugService `inject:""`
	BaseCtrl
}

func NewZentaoCtrl() *ZentaoCtrl {
	return &ZentaoCtrl{}
}

func (c *ZentaoCtrl) GetProfile(ctx iris.Context) {
	currSiteId, _ := ctx.URLParamInt("currSiteId")
	if currSiteId <= 0 {
		ctx.JSON(c.SuccessResp(iris.Map{}))
		return
	}

	site, _ := c.SiteService.Get(uint(currSiteId))
	if site.Url == "" {
		ctx.JSON(c.SuccessResp(iris.Map{}))
		return
	}

	config := configHelper.LoadBySite(site)
	data, err := zentaoHelper.GetProfile(config)
	if err != nil {
		ctx.JSON(c.ErrResp(commConsts.ErrZentaoRequest, err.Error()))
		return
	}

	ctx.JSON(c.SuccessResp(data))
}

func (c *ZentaoCtrl) ListSiteAndProduct(ctx iris.Context) {
	currSiteId, _ := ctx.URLParamInt("currSiteId")
	currProductId, _ := ctx.URLParamInt("currProductId")

	sites, currSite, _ := c.SiteService.LoadSites(currSiteId)
	products, currProduct, err := zentaoHelper.LoadSiteProduct(currSite, currProductId)

	for idx, _ := range sites {
		sites[idx].Url = ""
	}

	zentaoErr := false
	if err != nil {
		currSite = sites[len(sites)-1]
		zentaoErr = true
	}

	data := iris.Map{
		"zentaoErr": zentaoErr,
		"sites":     sites, "products": products,
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
		ctx.JSON(c.ErrResp(commConsts.ErrZentaoRequest, err.Error()))
		return
	}

	ctx.JSON(c.SuccessResp(data))
}

func (c *ZentaoCtrl) ListModule(ctx iris.Context) {
	siteId, err := ctx.URLParamInt("currSiteId")
	productId, err := ctx.URLParamInt("currProductId")

	if err != nil {
		ctx.JSON(c.ErrResp(commConsts.ParamErr, err.Error()))
		return
	}

	site, _ := c.SiteService.Get(uint(siteId))

	data, err := zentaoHelper.ListCaseModule(uint(productId), site)
	if err != nil {
		ctx.JSON(c.ErrResp(commConsts.ErrZentaoRequest, err.Error()))
		return
	}

	ctx.JSON(c.SuccessResp(data))
}

func (c *ZentaoCtrl) ListSuite(ctx iris.Context) {
	siteId, err := ctx.URLParamInt("currSiteId")
	productId, err := ctx.URLParamInt("currProductId")

	if err != nil {
		ctx.JSON(c.ErrResp(commConsts.ParamErr, err.Error()))
		return
	}

	site, _ := c.SiteService.Get(uint(siteId))
	data, err := zentaoHelper.ListSuite(uint(productId), site)
	if err != nil {
		ctx.JSON(c.ErrResp(commConsts.ErrZentaoRequest, err.Error()))
		return
	}

	ctx.JSON(c.SuccessResp(data))
}

func (c *ZentaoCtrl) ListTask(ctx iris.Context) {
	siteId, err := ctx.URLParamInt("currSiteId")
	productId, err := ctx.URLParamInt("currProductId")

	if err != nil {
		ctx.JSON(c.ErrResp(commConsts.ParamErr, err.Error()))
		return
	}

	site, _ := c.SiteService.Get(uint(siteId))
	data, err := zentaoHelper.ListTask(uint(productId), site)
	if err != nil {
		ctx.JSON(c.ErrResp(commConsts.ErrZentaoRequest, err.Error()))
		return
	}

	ctx.JSON(c.SuccessResp(data))
}

func (c *ZentaoCtrl) ListBugFields(ctx iris.Context) {
	siteId, _ := ctx.URLParamInt("currSiteId")
	productId, _ := ctx.URLParamInt("currProductId")

	if siteId == 0 || productId == 0 {
		ctx.JSON(c.ErrResp(commConsts.ParamErr, "siteId or productId"))
		return
	}

	data, _ := c.TestBugService.GetBugFields(siteId, productId)

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
