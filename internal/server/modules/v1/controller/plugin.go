package controller

import (
	commConsts "github.com/easysoft/zentaoatf/internal/pkg/consts"
	commDomain "github.com/easysoft/zentaoatf/internal/pkg/domain"
	"github.com/easysoft/zentaoatf/internal/server/modules/v1/service"
	"github.com/kataras/iris/v12"
)

type PluginCtrl struct {
	PluginService *service.PluginService `inject:""`
	BaseCtrl
}

func NewPluginCtrl() *PluginCtrl {
	return &PluginCtrl{}
}

func (c *PluginCtrl) Exec(ctx iris.Context) {
	err := c.PluginService.Exec()
	if err != nil {
		ctx.JSON(c.ErrResp(commConsts.CommErr, err.Error()))
		return
	}

	ctx.JSON(c.SuccessResp(nil))
}

func (c *PluginCtrl) Cancel(ctx iris.Context) {
	err := c.PluginService.Cancel()
	if err != nil {
		ctx.JSON(c.ErrResp(commConsts.CommErr, err.Error()))
		return
	}

	ctx.JSON(c.SuccessResp(nil))
}

func (c *PluginCtrl) Start(ctx iris.Context) {
	err := c.PluginService.Start()
	if err != nil {
		ctx.JSON(c.ErrResp(commConsts.CommErr, err.Error()))
		return
	}

	ctx.JSON(c.SuccessResp(nil))
}

func (c *PluginCtrl) Stop(ctx iris.Context) {
	err := c.PluginService.Stop()
	if err != nil {
		ctx.JSON(c.ErrResp(commConsts.CommErr, err.Error()))
		return
	}

	ctx.JSON(c.SuccessResp(nil))
}

func (c *PluginCtrl) Install(ctx iris.Context) {
	req := commDomain.PluginInstallReq{}
	err := ctx.ReadJSON(&req)
	if err != nil {
		ctx.JSON(c.ErrResp(commConsts.CommErr, err.Error()))
	}

	err = c.PluginService.Install(req)
	if err != nil {
		ctx.JSON(c.ErrResp(commConsts.CommErr, err.Error()))
		return
	}

	ctx.JSON(c.SuccessResp(nil))
}

func (c *PluginCtrl) Uninstall(ctx iris.Context) {
	err := c.PluginService.Uninstall()
	if err != nil {
		ctx.JSON(c.ErrResp(commConsts.CommErr, err.Error()))
		return
	}

	ctx.JSON(c.SuccessResp(nil))
}
