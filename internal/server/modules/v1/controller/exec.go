package controller

import (
	commConsts "github.com/easysoft/zentaoatf/internal/pkg/consts"
	serverDomain "github.com/easysoft/zentaoatf/internal/server/modules/v1/domain"
	"github.com/easysoft/zentaoatf/internal/server/modules/v1/service"
	"github.com/kataras/iris/v12"
)

type ExecCtrl struct {
	BaseCtrl
	TestExecService *service.TestExecService `inject:""`
}

func NewExecCtrl() *ExecCtrl {
	return &ExecCtrl{}
}

func (c *ExecCtrl) Start(ctx iris.Context) {
	req := serverDomain.ExecReq{}
	if err := ctx.ReadQuery(&req); err != nil {
		ctx.JSON(c.ErrResp(commConsts.CommErr, err.Error()))
		return
	}

	err := c.TestExecService.Start(req, nil)

	if err != nil {
		ctx.JSON(c.ErrResp(commConsts.CommErr, err.Error()))
		return
	}

	ctx.JSON(c.SuccessResp(nil))
}

func (c *ExecCtrl) Stop(ctx iris.Context) {
	req := serverDomain.ExecReq{}
	if err := ctx.ReadQuery(&req); err != nil {
		ctx.JSON(c.ErrResp(commConsts.CommErr, err.Error()))
		return
	}

	err := c.TestExecService.Stop(req, nil)
	if err != nil {
		ctx.JSON(c.ErrResp(commConsts.CommErr, err.Error()))
		return
	}

	ctx.JSON(c.SuccessResp(nil))
}
