package controller

import (
	commConsts "github.com/aaronchen2k/deeptest/internal/comm/consts"
	commDomain "github.com/aaronchen2k/deeptest/internal/comm/domain"
	configUtils "github.com/aaronchen2k/deeptest/internal/comm/helper/config"
	"github.com/aaronchen2k/deeptest/internal/server/modules/v1/service"
	"github.com/kataras/iris/v12"
)

type ConfigCtrl struct {
	ProjectService *service.ProjectService `inject:""`
	BaseCtrl
}

func NewConfigCtrl() *ConfigCtrl {
	return &ConfigCtrl{}
}

func (c *ConfigCtrl) SaveConfig(ctx iris.Context) {
	projectPath := ctx.URLParam("currProject")

	req := commDomain.ProjectConf{}
	err := ctx.ReadJSON(&req)
	if err != nil {
		ctx.JSON(c.ErrResp(commConsts.ParamErr, err.Error()))
		return
	}

	err = configUtils.SaveConfig(req, projectPath)
	if err != nil {
		ctx.JSON(c.ErrResp(commConsts.CommErr, err.Error()))
		return
	}

	projects, currProject, currProjectConfig, scriptTree, err := c.ProjectService.GetByUser(projectPath)
	if err != nil {
		ctx.JSON(c.ErrResp(commConsts.CommErr, err.Error()))
		return
	}

	data := iris.Map{"projects": projects, "currProject": currProject,
		"currConfig": currProjectConfig, "scriptTree": scriptTree}

	ctx.JSON(c.SuccessResp(data))
}
