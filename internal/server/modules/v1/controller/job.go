package controller

import (
	commConsts "github.com/easysoft/zentaoatf/internal/pkg/consts"
	serverDomain "github.com/easysoft/zentaoatf/internal/server/modules/v1/domain"
	"github.com/easysoft/zentaoatf/internal/server/modules/v1/service"
	"github.com/kataras/iris/v12"
)

type JobCtrl struct {
	JobService *service.JobService `inject:""`
	BaseCtrl
}

func NewJobCtrl() *JobCtrl {
	return &JobCtrl{}
}

func (c *JobCtrl) List(ctx iris.Context) {
	status := ctx.URLParam("status")

	jobs, err := c.JobService.List(status)
	if err != nil {
		ctx.JSON(c.ErrResp(commConsts.CommErr, err.Error()))
		return
	}

	ctx.JSON(c.SuccessResp(jobs))
	return
}

// @summary 添加下载任务
// @Accept json
// @Produce json
// @Param DownloadReq body []v1.DownloadReq true "Download Request Object"
// @Success 200 {object} _domain.Response "code = success | fail"
// @Router /api/v1/download/add [post]
func (c *JobCtrl) Add(ctx iris.Context) {
	var req serverDomain.ZentaoExecReq
	err := ctx.ReadJSON(&req)
	if err != nil {
		ctx.JSON(c.ErrResp(commConsts.ParamErr, err.Error()))
		return
	}

	err = c.JobService.Add(req)
	if err != nil {
		ctx.JSON(c.ErrResp(commConsts.ParamErr, err.Error()))
		return
	}

	ctx.JSON(c.SuccessResp(nil))
	return
}

// @summary 强制终止下载任务
// @Accept json
// @Produce json
// @Param DownloadCancelReq body v1.DownloadCancelReq true "CancelDate Download Request Object"
// @Success 200 {object} _domain.Response "code = success | fail"
// @Router /api/v1/download/cancel [post]
func (c *JobCtrl) Cancel(ctx iris.Context) {
	req := serverDomain.ZentaoCancelReq{}
	err := ctx.ReadJSON(&req)
	if err != nil {
		ctx.JSON(c.ErrResp(commConsts.ParamErr, err.Error()))
		return
	}

	c.JobService.Cancel(uint(req.Task))

	ctx.JSON(c.SuccessResp(nil))
	return
}
