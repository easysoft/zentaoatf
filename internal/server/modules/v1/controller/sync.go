package controller

import (
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
	//workspacePath := ctx.URLParam("currWorkspace")
	//
	//req := commDomain.SyncSettings{}
	//err := ctx.ReadJSON(&req)
	//if err != nil {
	//	c.ErrResp(commConsts.ParamErr, err.Error())
	//	return
	//}
	//
	//err = c.SyncService.SyncFromZentao(req, workspacePath)
	//if err != nil {
	//	ctx.JSON(c.ErrResp(commConsts.CommErr, err.Error()))
	//	return
	//}

	ctx.JSON(c.SuccessResp(nil))
}

func (c *SyncCtrl) SyncToZentao(ctx iris.Context) {
	//workspacePath := ctx.URLParam("currWorkspace")
	//commitProductIdStr := ctx.URLParam("commitProductId")
	//commitProductId, _ := strconv.Atoi(commitProductIdStr)
	//
	//if commitProductId == 0 {
	//	ctx.JSON(c.ErrResp(commConsts.ParamErr, ""))
	//	return
	//}
	//
	//err := c.SyncService.SyncToZentao(workspacePath, commitProductId)
	//if err != nil {
	//	ctx.JSON(c.ErrResp(commConsts.CommErr, err.Error()))
	//	return
	//}

	ctx.JSON(c.SuccessResp(nil))
}
