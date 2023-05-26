package controller

import (
	commConsts "github.com/easysoft/zentaoatf/internal/pkg/consts"
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
	err := c.PluginService.Install()
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
