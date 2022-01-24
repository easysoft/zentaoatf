package controller

import (
	commDomain "github.com/aaronchen2k/deeptest/internal/comm/domain"
	"github.com/aaronchen2k/deeptest/internal/pkg/domain"
	logUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/log"
	"github.com/aaronchen2k/deeptest/internal/server/core/web/validate"
	"github.com/aaronchen2k/deeptest/internal/server/modules/v1/service"
	configUtils "github.com/aaronchen2k/deeptest/internal/server/modules/v1/utils/config"
	"strings"

	"github.com/kataras/iris/v12"
	"go.uber.org/zap"
)

type ConfigCtrl struct {
	BaseCtrl
	ProjectService *service.ProjectService `inject:""`
}

func NewConfigCtrl() *ConfigCtrl {
	return &ConfigCtrl{}
}

func (c *ConfigCtrl) SaveConfig(ctx iris.Context) {
	projectPath := ctx.URLParam("currProject")

	req := commDomain.ProjectConf{}
	if err := ctx.ReadJSON(&req); err != nil {
		errs := validate.ValidRequest(err)
		if len(errs) > 0 {
			logUtils.Errorf("参数验证失败", zap.String("错误", strings.Join(errs, ";")))
			ctx.JSON(domain.Response{Code: domain.SystemErr.Code, Data: nil, Msg: strings.Join(errs, ";")})
			return
		}
	}

	err := configUtils.SaveConfig(req, projectPath)
	if err != nil {
		ctx.JSON(domain.Response{Code: c.ErrCode(err), Data: nil})
		return
	}

	projects, currProject, currProjectConfig, scriptTree, err := c.ProjectService.GetByUser(projectPath)
	if err != nil {
		ctx.JSON(domain.Response{Code: domain.SystemErr.Code, Data: nil, Msg: err.Error()})
		return
	}

	ret := iris.Map{"projects": projects, "currProject": currProject,
		"currConfig": currProjectConfig, "scriptTree": scriptTree}

	ctx.JSON(domain.Response{Code: domain.NoErr.Code, Data: ret, Msg: domain.NoErr.Msg})
}
