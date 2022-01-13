package controller

import (
	commDomain "github.com/aaronchen2k/deeptest/internal/comm/domain"
	"github.com/aaronchen2k/deeptest/internal/pkg/domain"
	logUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/log"
	"github.com/aaronchen2k/deeptest/internal/server/modules/v1/service"
	"github.com/kataras/iris/v12"
)

type SyncCtrl struct {
	SyncService *service.SyncService `inject:""`
	BaseCtrl
}

func NewSyncCtrl() *SyncCtrl {
	return &SyncCtrl{}
}

func (c *SyncCtrl) SyncFromZentao(ctx iris.Context) {
	projectPath := ctx.URLParam("currProject")

	req := commDomain.SyncSettings{}
	err := ctx.ReadJSON(&req)

	if err != nil {
		logUtils.Errorf("参数验证失败 %s", err.Error())
		ctx.JSON(domain.Response{Code: domain.SystemErr.Code, Data: nil, Msg: err.Error()})
		return
	}

	err = c.SyncService.SyncFromZentao(req, projectPath)
	if err != nil {
		ctx.JSON(domain.Response{Code: c.ErrCode(err), Data: nil})
		return
	}

	ctx.JSON(domain.Response{Code: domain.NoErr.Code, Data: nil, Msg: domain.NoErr.Msg})
}

func (c *SyncCtrl) SyncToZentao(ctx iris.Context) {
	projectPath := ctx.URLParam("currProject")

	err := c.SyncService.SyncToZentao(projectPath)
	if err != nil {
		ctx.JSON(domain.Response{Code: c.ErrCode(err), Data: nil})
		return
	}

	ctx.JSON(domain.Response{Code: domain.NoErr.Code, Data: nil, Msg: domain.NoErr.Msg})
}
