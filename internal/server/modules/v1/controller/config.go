package controller

import (
	commConsts "github.com/aaronchen2k/deeptest/internal/comm/consts"
	commDomain "github.com/aaronchen2k/deeptest/internal/comm/domain"
	configUtils "github.com/aaronchen2k/deeptest/internal/comm/helper/config"
	"github.com/aaronchen2k/deeptest/internal/server/modules/v1/service"
	"github.com/kataras/iris/v12"
)

type ConfigCtrl struct {
	WorkspaceService *service.WorkspaceService `inject:""`
	BaseCtrl
}

func NewConfigCtrl() *ConfigCtrl {
	return &ConfigCtrl{}
}

func (c *ConfigCtrl) SaveConfig(ctx iris.Context) {
	workspacePath := ctx.URLParam("currWorkspace")

	req := commDomain.WorkspaceConf{}
	err := ctx.ReadJSON(&req)
	if err != nil {
		ctx.JSON(c.ErrResp(commConsts.ParamErr, err.Error()))
		return
	}

	err = configUtils.SaveConfig(req, workspacePath)
	if err != nil {
		ctx.JSON(c.ErrResp(commConsts.CommErr, err.Error()))
		return
	}

	workspaces, currWorkspace, currWorkspaceConfig, scriptTree, err := c.WorkspaceService.GetByUser(workspacePath)
	if err != nil {
		ctx.JSON(c.ErrResp(commConsts.CommErr, err.Error()))
		return
	}

	data := iris.Map{"workspaces": workspaces, "currWorkspace": currWorkspace,
		"currConfig": currWorkspaceConfig, "scriptTree": scriptTree}

	ctx.JSON(c.SuccessResp(data))
}
