package controller

import (
	"fmt"
	commConsts "github.com/aaronchen2k/deeptest/internal/comm/consts"
	commDomain "github.com/aaronchen2k/deeptest/internal/comm/domain"
	"github.com/aaronchen2k/deeptest/internal/pkg/domain"
	serverDomain "github.com/aaronchen2k/deeptest/internal/server/modules/v1/domain"
	"github.com/aaronchen2k/deeptest/internal/server/modules/v1/model"
	"github.com/aaronchen2k/deeptest/internal/server/modules/v1/service"
	"github.com/kataras/iris/v12"
)

type WorkspaceCtrl struct {
	WorkspaceService *service.WorkspaceService `inject:""`
	BaseCtrl
}

func NewWorkspaceCtrl() *WorkspaceCtrl {
	return &WorkspaceCtrl{}
}

func (c *WorkspaceCtrl) List(ctx iris.Context) {
	currProductId, _ := ctx.URLParamInt("currProductId")

	var req serverDomain.WorkspaceReqPaginate
	if err := ctx.ReadQuery(&req); err != nil {
		ctx.JSON(c.ErrResp(commConsts.CommErr, err.Error()))
		return
	}

	req.ProductId = currProductId
	data, err := c.WorkspaceService.Paginate(req)
	if err != nil {
		ctx.JSON(c.ErrResp(commConsts.CommErr, err.Error()))
		return
	}

	ctx.JSON(c.SuccessResp(data))
}

func (c *WorkspaceCtrl) Get(ctx iris.Context) {
	id, err := ctx.Params().GetInt("id")
	if err != nil {
		ctx.JSON(c.ErrResp(commConsts.CommErr, err.Error()))
		return
	}

	po, err := c.WorkspaceService.Get(uint(id))
	if err != nil {
		ctx.JSON(c.ErrResp(commConsts.CommErr, err.Error()))
		return
	}
	ctx.JSON(c.SuccessResp(po))
}

func (c *WorkspaceCtrl) Create(ctx iris.Context) {
	currSiteId, _ := ctx.URLParamInt("currSiteId")
	currProductId, _ := ctx.URLParamInt("currProductId")

	req := model.Workspace{}
	if err := ctx.ReadJSON(&req); err != nil {
		ctx.JSON(c.ErrResp(commConsts.CommErr, err.Error()))
	}

	req.SiteId = uint(currSiteId)
	req.ProductId = uint(currProductId)
	id, err := c.WorkspaceService.Create(req)
	if err != nil {
		ctx.JSON(c.ErrResp(commConsts.CommErr, err.Error()))
		return
	}

	ctx.JSON(c.SuccessResp(iris.Map{"id": id}))
}

func (c *WorkspaceCtrl) Update(ctx iris.Context) {
	currSiteId, _ := ctx.URLParamInt("currSiteId")
	currProductId, _ := ctx.URLParamInt("currProductId")
	if currProductId <= 0 {
		ctx.JSON(c.ErrResp(commConsts.ParamErr, fmt.Sprintf("参数%s不合法", "currProductId")))
		return
	}

	req := model.Workspace{}
	if err := ctx.ReadJSON(&req); err != nil {
		ctx.JSON(c.ErrResp(commConsts.CommErr, err.Error()))
	}

	req.SiteId = uint(currSiteId)
	req.ProductId = uint(currProductId)
	err := c.WorkspaceService.Update(req)
	if err != nil {
		ctx.JSON(c.ErrResp(commConsts.CommErr, err.Error()))
		return
	}

	ctx.JSON(c.SuccessResp(iris.Map{"id": req.ID}))
}

// Delete 删除
func (c *WorkspaceCtrl) Delete(ctx iris.Context) {
	workspaceId, _ := ctx.Params().GetInt("id")

	if workspaceId <= 0 {
		ctx.JSON(c.ErrResp(commConsts.ParamErr, "id"))
		return
	}
	err := c.WorkspaceService.Delete(uint(workspaceId))
	if err != nil {
		ctx.JSON(c.ErrResp(commConsts.CommErr, err.Error()))
		return
	}

	ctx.JSON(c.SuccessResp(nil))
}

func (c *WorkspaceCtrl) GetByUser(ctx iris.Context) {
	workspacePath := ctx.URLParam("currWorkspace")

	if workspacePath == "" {
		workspaces, _ := c.WorkspaceService.ListWorkspaceByUser()
		data := iris.Map{
			"workspaces":    workspaces,
			"currWorkspace": model.Workspace{},
			"currConfig":    commDomain.WorkspaceConf{},
			"scriptTree":    serverDomain.TestAsset{}}

		ctx.JSON(c.SuccessResp(data))
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

func (c *WorkspaceCtrl) ListByProduct(ctx iris.Context) {
	currSiteId, _ := ctx.URLParamInt("currSiteId")
	currProductId, _ := ctx.URLParamInt("currProductId")
	if currProductId <= 0 {
		ctx.JSON(c.SuccessResp(domain.PageData{}))
		return
	}

	data, err := c.WorkspaceService.ListByProduct(currSiteId, currProductId)
	if err != nil {
		ctx.JSON(c.ErrResp(commConsts.CommErr, err.Error()))
		return
	}

	ctx.JSON(c.SuccessResp(data))
}
