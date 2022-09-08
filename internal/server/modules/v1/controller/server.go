package controller

import (
	commConsts "github.com/easysoft/zentaoatf/internal/pkg/consts"
	"github.com/easysoft/zentaoatf/internal/server/modules/v1/model"
	"github.com/easysoft/zentaoatf/internal/server/modules/v1/service"
	"github.com/kataras/iris/v12"
)

type ServerCtrl struct {
	ServerService    *service.ServerService    `inject:""`
	WorkspaceService *service.WorkspaceService `inject:""`
	BaseCtrl
}

func NewServerCtrl() *ServerCtrl {
	return &ServerCtrl{}
}

func (c *ServerCtrl) List(ctx iris.Context) {
	data, err := c.ServerService.List()
	if err != nil {
		ctx.JSON(c.ErrResp(commConsts.CommErr, err.Error()))
		return
	}

	ctx.JSON(c.SuccessResp(data))
}

func (c *ServerCtrl) Get(ctx iris.Context) {
	id, err := ctx.Params().GetInt("id")
	if err != nil {
		ctx.JSON(c.ErrResp(commConsts.CommErr, err.Error()))
		return
	}

	po, err := c.ServerService.Get(uint(id))
	if err != nil {
		ctx.JSON(c.ErrResp(commConsts.CommErr, err.Error()))
		return
	}
	ctx.JSON(c.SuccessResp(po))
}

func (c *ServerCtrl) Create(ctx iris.Context) {
	req := model.Server{}
	if err := ctx.ReadJSON(&req); err != nil {
		ctx.JSON(c.ErrResp(commConsts.CommErr, err.Error()))
	}

	id, err := c.ServerService.Create(req)
	if err != nil {
		ctx.JSON(c.ErrResp(commConsts.CommErr, err.Error()))
		return
	}

	c.WorkspaceService.UpdateAllConfig()

	ctx.JSON(c.SuccessResp(iris.Map{"id": id}))
}

func (c *ServerCtrl) Update(ctx iris.Context) {
	req := model.Server{}
	if err := ctx.ReadJSON(&req); err != nil {
		ctx.JSON(c.ErrResp(commConsts.CommErr, err.Error()))
	}

	err := c.ServerService.Update(req)
	if err != nil {
		ctx.JSON(c.ErrResp(commConsts.CommErr, err.Error()))
		return
	}

	c.WorkspaceService.UpdateAllConfig()

	ctx.JSON(c.SuccessResp(iris.Map{"id": req.ID}))
}

func (c *ServerCtrl) Delete(ctx iris.Context) {
	id, err := ctx.Params().GetInt("id")
	if err != nil {
		ctx.JSON(c.ErrResp(commConsts.CommErr, err.Error()))
		return
	}

	err = c.ServerService.Delete(uint(id))
	if err != nil {
		ctx.JSON(c.ErrResp(commConsts.CommErr, err.Error()))
		return
	}

	c.WorkspaceService.UpdateAllConfig()

	ctx.JSON(c.SuccessResp(nil))
}
