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
	projectPath := ctx.URLParam("currProject")

	if projectPath == "" {
		ctx.JSON(c.SuccessResp(make([]serverDomain.TestReportSummary, 0)))
		return
	}

	data, err := c.TestExecService.List(projectPath)
	if err != nil {
		ctx.JSON(c.ErrResp(commConsts.Failure, err.Error()))
		return
	}

	ctx.JSON(c.SuccessResp(data))
}

// Get 详情
func (c *TestExecCtrl) Get(ctx iris.Context) {
	projectPath := ctx.URLParam("currProject")

	seq := ctx.Params().Get("seq")
	if seq == "" {
		c.ErrResp(commConsts.ParamErr, "seq")
		return
	}

	exec, err := c.TestExecService.Get(projectPath, seq)
	if err != nil {
		ctx.JSON(c.ErrResp(commConsts.Failure, err.Error()))
		return
	}
	ctx.JSON(c.SuccessResp(exec))
}

// Delete 删除
func (c *TestExecCtrl) Delete(ctx iris.Context) {
	projectPath := ctx.URLParam("currProject")

	seq := ctx.Params().Get("seq")
	if seq == "" {
		c.ErrResp(commConsts.ParamErr, "seq")
		return
	}

	err := c.TestExecService.Delete(projectPath, seq)
	if err != nil {
		ctx.JSON(c.ErrResp(commConsts.Failure, err.Error()))
		return
	}

	ctx.JSON(c.SuccessResp(nil))
}
