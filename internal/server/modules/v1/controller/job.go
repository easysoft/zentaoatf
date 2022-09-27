package controller

import (
	commConsts "github.com/easysoft/zentaoatf/internal/pkg/consts"
	serverDomain "github.com/easysoft/zentaoatf/internal/server/modules/v1/domain"
	"github.com/easysoft/zentaoatf/internal/server/modules/v1/service"
	"github.com/kataras/iris/v12"
)

type JobCtrl struct {
	BaseCtrl
	JobService  *service.JobService  `inject:""`
	ExecService *service.ExecService `inject:""`
}

func NewJobCtrl() *JobCtrl {
	return &JobCtrl{}
}

func (c *JobCtrl) Add(ctx iris.Context) {
	req := serverDomain.JobReq{}
	if err := ctx.ReadQuery(&req); err != nil {
		ctx.JSON(c.ErrResp(commConsts.CommErr, err.Error()))
		return
	}

	err := c.JobService.Add(req)

	if err != nil {
		ctx.JSON(c.ErrResp(commConsts.CommErr, err.Error()))
		return
	}

	ctx.JSON(c.SuccessResp(nil))
}

func (c *JobCtrl) Remove(ctx iris.Context) {
	req := serverDomain.JobReq{}
	if err := ctx.ReadQuery(&req); err != nil {
		ctx.JSON(c.ErrResp(commConsts.CommErr, err.Error()))
		return
	}

	err := c.JobService.Remove(req)
	if err != nil {
		ctx.JSON(c.ErrResp(commConsts.CommErr, err.Error()))
		return
	}

	ctx.JSON(c.SuccessResp(nil))
}

func (c *JobCtrl) Stop(ctx iris.Context) {
	req := serverDomain.JobReq{}
	if err := ctx.ReadQuery(&req); err != nil {
		ctx.JSON(c.ErrResp(commConsts.CommErr, err.Error()))
		return
	}

	err := c.JobService.Stop()
	if err != nil {
		ctx.JSON(c.ErrResp(commConsts.CommErr, err.Error()))
		return
	}

	ctx.JSON(c.SuccessResp(nil))
}
