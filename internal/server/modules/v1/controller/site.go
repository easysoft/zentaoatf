package controller

import (
	commConsts "github.com/aaronchen2k/deeptest/internal/comm/consts"
	serverDomain "github.com/aaronchen2k/deeptest/internal/server/modules/v1/domain"
	"github.com/aaronchen2k/deeptest/internal/server/modules/v1/service"
	"github.com/kataras/iris/v12"
)

type SiteCtrl struct {
	SiteService *service.SiteService `inject:""`
	BaseCtrl
}

func NewSiteCtrl() *SiteCtrl {
	return &SiteCtrl{}
}

// List 分页列表
func (c *SiteCtrl) List(ctx iris.Context) {
	projectPath := ctx.URLParam("currProject")

	if projectPath == "" {
		ctx.JSON(c.SuccessResp(make([]serverDomain.TestReportSummary, 0)))
		return
	}

	data, err := c.SiteService.List(projectPath)
	if err != nil {
		ctx.JSON(c.ErrResp(commConsts.Failure, err.Error()))
		return
	}

	ctx.JSON(c.SuccessResp(data))
}

// Get 详情
func (c *SiteCtrl) Get(ctx iris.Context) {
	projectPath := ctx.URLParam("currProject")

	seq := ctx.Params().Get("seq")
	if seq == "" {
		c.ErrResp(commConsts.ParamErr, "seq")
		return
	}

	exec, err := c.SiteService.Get(projectPath, seq)
	if err != nil {
		ctx.JSON(c.ErrResp(commConsts.Failure, err.Error()))
		return
	}
	ctx.JSON(c.SuccessResp(exec))
}

// Delete 删除
func (c *SiteCtrl) Delete(ctx iris.Context) {
	projectPath := ctx.URLParam("currProject")

	seq := ctx.Params().Get("seq")
	if seq == "" {
		c.ErrResp(commConsts.ParamErr, "seq")
		return
	}

	err := c.SiteService.Delete(projectPath, seq)
	if err != nil {
		ctx.JSON(c.ErrResp(commConsts.Failure, err.Error()))
		return
	}

	ctx.JSON(c.SuccessResp(nil))
}
