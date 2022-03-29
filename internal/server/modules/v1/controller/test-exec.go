package controller

import (
	commConsts "github.com/aaronchen2k/deeptest/internal/comm/consts"
	serverDomain "github.com/aaronchen2k/deeptest/internal/server/modules/v1/domain"
	"github.com/aaronchen2k/deeptest/internal/server/modules/v1/service"
	"github.com/kataras/iris/v12"
)

type TestExecCtrl struct {
	TestExecService *service.TestExecService `inject:""`
	BaseCtrl
}

func NewTestExecCtrl() *TestExecCtrl {
	return &TestExecCtrl{}
}

// List 分页列表
func (c *TestExecCtrl) List(ctx iris.Context) {
	currSiteId, _ := ctx.URLParamInt("currSiteId")
	currProductId, _ := ctx.URLParamInt("currProductId")

	var req serverDomain.ReqPaginate
	if err := ctx.ReadQuery(&req); err != nil {
		ctx.JSON(c.ErrResp(commConsts.CommErr, err.Error()))
		return
	}

	data, err := c.TestExecService.Paginate(currSiteId, currProductId, req)
	if err != nil {
		ctx.JSON(c.ErrResp(commConsts.CommErr, err.Error()))
		return
	}

	ctx.JSON(c.SuccessResp(data))
}

// Get 详情
func (c *TestExecCtrl) Get(ctx iris.Context) {
	workspaceId, _ := ctx.Params().GetInt("workspaceId")
	seq := ctx.Params().Get("seq")

	if workspaceId == 0 || seq == "" {
		c.ErrResp(commConsts.ParamErr, "workspaceId and seq")
		return
	}

	exec, err := c.TestExecService.Get(workspaceId, seq)
	if err != nil {
		ctx.JSON(c.ErrResp(commConsts.CommErr, err.Error()))
		return
	}
	ctx.JSON(c.SuccessResp(exec))
}

// Delete 删除
func (c *TestExecCtrl) Delete(ctx iris.Context) {
	workspaceId, _ := ctx.URLParamInt("workspaceId")
	seq := ctx.URLParam("seq")

	if workspaceId == 0 || seq == "" {
		c.ErrResp(commConsts.ParamErr, "workspaceId and seq")
		return
	}

	err := c.TestExecService.Delete(workspaceId, seq)
	if err != nil {
		ctx.JSON(c.ErrResp(commConsts.CommErr, err.Error()))
		return
	}

	ctx.JSON(c.SuccessResp(nil))
}
