package controller

import (
	commConsts "github.com/easysoft/zentaoatf/internal/comm/consts"
	serverDomain "github.com/easysoft/zentaoatf/internal/server/modules/v1/domain"
	"github.com/easysoft/zentaoatf/internal/server/modules/v1/model"
	"github.com/easysoft/zentaoatf/internal/server/modules/v1/service"
	"github.com/kataras/iris/v12"
)

type SiteCtrl struct {
	SiteService *service.SiteService `inject:""`
	BaseCtrl
}

func NewSiteCtrl() *SiteCtrl {
	return &SiteCtrl{}
}

func (c *SiteCtrl) List(ctx iris.Context) {
	var req serverDomain.ReqPaginate
	if err := ctx.ReadQuery(&req); err != nil {
		ctx.JSON(c.ErrResp(commConsts.CommErr, err.Error()))
		return
	}

	data, err := c.SiteService.Paginate(req)
	if err != nil {
		ctx.JSON(c.ErrResp(commConsts.CommErr, err.Error()))
		return
	}

	ctx.JSON(c.SuccessResp(data))
}

func (c *SiteCtrl) Get(ctx iris.Context) {
	id, err := ctx.Params().GetInt("id")
	if err != nil {
		ctx.JSON(c.ErrResp(commConsts.CommErr, err.Error()))
		return
	}

	po, err := c.SiteService.Get(uint(id))
	if err != nil {
		ctx.JSON(c.ErrResp(commConsts.CommErr, err.Error()))
		return
	}
	ctx.JSON(c.SuccessResp(po))
}

func (c *SiteCtrl) Create(ctx iris.Context) {
	req := model.Site{}
	if err := ctx.ReadJSON(&req); err != nil {
		ctx.JSON(c.ErrResp(commConsts.ParamErr, err.Error()))
		return
	}

	id, isDuplicate, err := c.SiteService.Create(req)
	if isDuplicate {
		ctx.JSON(c.ErrResp(commConsts.ErrRecordWithSameNameExist, ""))
	}
	if err != nil {
		ctx.JSON(c.ErrResp(commConsts.ErrZentaoConfig, err.Error()))
		return
	}

	ctx.JSON(c.SuccessResp(iris.Map{"id": id}))
}

func (c *SiteCtrl) Update(ctx iris.Context) {
	req := model.Site{}
	if err := ctx.ReadJSON(&req); err != nil {
		ctx.JSON(c.ErrResp(commConsts.ParamErr, err.Error()))
	}

	isDuplicate, err := c.SiteService.Update(req)
	if isDuplicate {
		ctx.JSON(c.ErrResp(commConsts.ErrRecordWithSameNameExist, ""))
		return
	}
	if err != nil {
		ctx.JSON(c.ErrResp(commConsts.ErrZentaoConfig, err.Error()))
		return
	}

	ctx.JSON(c.SuccessResp(iris.Map{"id": req.ID}))
}

func (c *SiteCtrl) Delete(ctx iris.Context) {
	id, err := ctx.Params().GetInt("id")
	if err != nil {
		ctx.JSON(c.ErrResp(commConsts.CommErr, err.Error()))
		return
	}

	err = c.SiteService.Delete(uint(id))
	if err != nil {
		ctx.JSON(c.ErrResp(commConsts.CommErr, err.Error()))
		return
	}

	ctx.JSON(c.SuccessResp(nil))
}
