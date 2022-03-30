package controller

import (
	commConsts "github.com/aaronchen2k/deeptest/internal/comm/consts"
	serverDomain "github.com/aaronchen2k/deeptest/internal/server/modules/v1/domain"
	"github.com/aaronchen2k/deeptest/internal/server/modules/v1/service"
	"github.com/kataras/iris/v12"
)

type TestHistoryCtrl struct {
	TestHistoryService *service.TestHistoryService `inject:""`
	BaseCtrl
}

func NewTestHistoryCtrl() *TestHistoryCtrl {
	return &TestHistoryCtrl{}
}

// List 分页列表
func (c *TestHistoryCtrl) List(ctx iris.Context) {
	currSiteId, _ := ctx.URLParamInt("currSiteId")
	currProductId, _ := ctx.URLParamInt("currProductId")

	var req serverDomain.ReqPaginate
	if err := ctx.ReadQuery(&req); err != nil {
		ctx.JSON(c.ErrResp(commConsts.CommErr, err.Error()))
		return
	}

	data, err := c.TestHistoryService.Paginate(currSiteId, currProductId, req)
	if err != nil {
		ctx.JSON(c.ErrResp(commConsts.CommErr, err.Error()))
		return
	}

	ctx.JSON(c.SuccessResp(data))
}

// Get 详情
func (c *TestHistoryCtrl) Get(ctx iris.Context) {
	workspaceId, _ := ctx.Params().GetInt("workspaceId")
	seq := ctx.Params().Get("seq")

	if workspaceId == 0 || seq == "" {
		c.ErrResp(commConsts.ParamErr, "workspaceId and seq")
		return
	}

	exec, err := c.TestHistoryService.Get(workspaceId, seq)
	if err != nil {
		ctx.JSON(c.ErrResp(commConsts.CommErr, err.Error()))
		return
	}
	ctx.JSON(c.SuccessResp(exec))
}

// Delete 删除
func (c *TestHistoryCtrl) Delete(ctx iris.Context) {
	workspaceId, _ := ctx.URLParamInt("workspaceId")
	seq := ctx.URLParam("seq")

	if workspaceId == 0 || seq == "" {
		c.ErrResp(commConsts.ParamErr, "workspaceId and seq")
		return
	}

	err := c.TestHistoryService.Delete(workspaceId, seq)
	if err != nil {
		ctx.JSON(c.ErrResp(commConsts.CommErr, err.Error()))
		return
	}

	ctx.JSON(c.SuccessResp(nil))
}
