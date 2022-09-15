package controller

import (
	commConsts "github.com/easysoft/zentaoatf/internal/pkg/consts"
	"github.com/easysoft/zentaoatf/internal/server/modules/v1/model"
	"github.com/easysoft/zentaoatf/internal/server/modules/v1/service"
	"github.com/kataras/iris/v12"
)

type ProxyCtrl struct {
	ProxyService     *service.ProxyService     `inject:""`
	WorkspaceService *service.WorkspaceService `inject:""`
	BaseCtrl
}

func NewProxyCtrl() *ProxyCtrl {
	return &ProxyCtrl{}
}

func (c *ProxyCtrl) List(ctx iris.Context) {
	data, err := c.ProxyService.List()
	if err != nil {
		ctx.JSON(c.ErrResp(commConsts.CommErr, err.Error()))
		return
	}

	ctx.JSON(c.SuccessResp(data))
}

func (c *ProxyCtrl) Get(ctx iris.Context) {
	id, err := ctx.Params().GetInt("id")
	if err != nil {
		ctx.JSON(c.ErrResp(commConsts.CommErr, err.Error()))
		return
	}

	po, err := c.ProxyService.Get(uint(id))
	if err != nil {
		ctx.JSON(c.ErrResp(commConsts.CommErr, err.Error()))
		return
	}
	ctx.JSON(c.SuccessResp(po))
}

func (c *ProxyCtrl) Create(ctx iris.Context) {
	req := model.Proxy{}
	if err := ctx.ReadJSON(&req); err != nil {
		ctx.JSON(c.ErrResp(commConsts.CommErr, err.Error()))
	}

	id, err := c.ProxyService.Create(req)
	if err != nil {
		ctx.JSON(c.ErrResp(commConsts.ErrPostParam, err.Error()))
		return
	}

	c.WorkspaceService.UpdateAllConfig()

	ctx.JSON(c.SuccessResp(iris.Map{"id": id}))
}

func (c *ProxyCtrl) Update(ctx iris.Context) {
	req := model.Proxy{}
	if err := ctx.ReadJSON(&req); err != nil {
		ctx.JSON(c.ErrResp(commConsts.CommErr, err.Error()))
	}

	err := c.ProxyService.Update(req)
	if err != nil {
		ctx.JSON(c.ErrResp(commConsts.ErrPostParam, err.Error()))
		return
	}

	c.WorkspaceService.UpdateAllConfig()

	ctx.JSON(c.SuccessResp(iris.Map{"id": req.ID}))
}

func (c *ProxyCtrl) Delete(ctx iris.Context) {
	id, err := ctx.Params().GetInt("id")
	if err != nil {
		ctx.JSON(c.ErrResp(commConsts.CommErr, err.Error()))
		return
	}

	err = c.ProxyService.Delete(uint(id))
	if err != nil {
		ctx.JSON(c.ErrResp(commConsts.CommErr, err.Error()))
		return
	}

	c.WorkspaceService.UpdateAllConfig()

	ctx.JSON(c.SuccessResp(nil))
}
