package controller

import (
	commDomain "github.com/aaronchen2k/deeptest/internal/comm/domain"
	configUtils "github.com/aaronchen2k/deeptest/internal/comm/helper/config"
	"github.com/aaronchen2k/deeptest/internal/pkg/domain"
	logUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/log"
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
		logUtils.Errorf(domain.ParamErr.Msg, err.Error())
		ctx.JSON(domain.Response{Code: domain.RequestErr.Code, Data: nil, Msg: ""})
		return
	}

	err = configUtils.SaveConfig(req, projectPath)
	if err != nil {
		ctx.JSON(domain.Response{Code: c.ErrCode(err), Data: nil})
		return
	}

	projects, currProject, currProjectConfig, scriptTree, err := c.ProjectService.GetByUser(projectPath)
	if err != nil {
		ctx.JSON(domain.Response{Code: domain.RequestErr.Code, Data: nil, Msg: err.Error()})
		return
	}

	ret := iris.Map{"projects": projects, "currProject": currProject,
		"currConfig": currProjectConfig, "scriptTree": scriptTree}

	ctx.JSON(domain.Response{Code: domain.NoErr.Code, Data: ret, Msg: domain.NoErr.Msg})
}
