package controller

import (
	commConsts "github.com/aaronchen2k/deeptest/internal/comm/consts"
	commDomain "github.com/aaronchen2k/deeptest/internal/comm/domain"
	"github.com/aaronchen2k/deeptest/internal/server/modules/v1/service"
	"github.com/kataras/iris/v12"
	"strconv"
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
		c.ErrResp(commConsts.ParamErr, err.Error())
		return
	}

	err = c.SyncService.SyncFromZentao(req, projectPath)
	if err != nil {
		ctx.JSON(c.ErrResp(commConsts.Failure, err.Error()))
		return
	}

	ctx.JSON(c.SuccessResp(nil))
}

func (c *SyncCtrl) SyncToZentao(ctx iris.Context) {
	projectPath := ctx.URLParam("currProject")
	commitProductIdStr := ctx.URLParam("commitProductId")
	commitProductId, _ := strconv.Atoi(commitProductIdStr)

	if commitProductId == 0 {
		ctx.JSON(c.ErrResp(commConsts.ParamErr, ""))
		return
	}

	err := c.SyncService.SyncToZentao(projectPath, commitProductId)
	if err != nil {
		ctx.JSON(c.ErrResp(commConsts.Failure, err.Error()))
		return
	}

	ctx.JSON(c.SuccessResp(nil))
}
