package controller

import (
	commConsts "github.com/easysoft/zentaoatf/internal/pkg/consts"
	serverDomain "github.com/easysoft/zentaoatf/internal/server/modules/v1/domain"
	"github.com/easysoft/zentaoatf/internal/server/modules/v1/model"
	"github.com/easysoft/zentaoatf/internal/server/modules/v1/service"
	ztfConsts "github.com/easysoft/zentaoatf/pkg/consts"
	"github.com/easysoft/zentaoatf/pkg/domain"
	"github.com/kataras/iris/v12"
)

type WorkspaceCtrl struct {
	WorkspaceService *service.WorkspaceService `inject:""`
	ProxyService     *service.ProxyService     `inject:""`
	BaseCtrl
}

func NewWorkspaceCtrl() *WorkspaceCtrl {
	return &WorkspaceCtrl{}
}

func (c *WorkspaceCtrl) List(ctx iris.Context) {
	currSiteId, _ := ctx.URLParamInt("currSiteId")
	currProductId, _ := ctx.URLParamInt("currProductId")

	var req serverDomain.WorkspaceReqPaginate
	if err := ctx.ReadQuery(&req); err != nil {
		ctx.JSON(c.ErrResp(commConsts.CommErr, err.Error()))
		return
	}

	req.ProductId = currProductId
	req.SiteId = currSiteId
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

// delete by path
func (c *WorkspaceCtrl) DeleteByPath(ctx iris.Context) {
	path := ctx.URLParam("path")
	currProductId, err := ctx.URLParamInt64("currProductId")
	if err != nil {
		ctx.JSON(c.ErrResp(commConsts.ParamErr, "currProductId"))
		return
	}

	err = c.WorkspaceService.DeleteByPath(path, uint(currProductId))
	if err != nil {
		ctx.JSON(c.ErrResp(commConsts.CommErr, err.Error()))
		return
	}

	ctx.JSON(c.SuccessResp(nil))
}

func (c *WorkspaceCtrl) ListByProduct(ctx iris.Context) {
	currSiteId, _ := ctx.URLParamInt("currSiteId")
	currProductId, _ := ctx.URLParamInt("currProductId")
	if currProductId <= 0 {
		ctx.JSON(c.SuccessResp(domain.PageData{}))
		return
	}

	data, err := c.WorkspaceService.ListByProduct(uint(currSiteId), uint(currProductId))
	if err != nil {
		ctx.JSON(c.ErrResp(commConsts.CommErr, err.Error()))
		return
	}

	ctx.JSON(c.SuccessResp(data))
}

func (c *WorkspaceCtrl) UploadScriptsToProxy(ctx iris.Context) {
	testSets := []serverDomain.TestSet{}
	var err error
	if err := ctx.ReadJSON(&testSets); err != nil {
		ctx.JSON(c.ErrResp(commConsts.CommErr, err.Error()))
		return
	}
	if err != nil {
		ctx.JSON(c.ErrResp(commConsts.CommErr, err.Error()))
		return
	}
	pathMap, err := c.WorkspaceService.UploadScriptsToProxy(testSets)
	if err != nil {
		ctx.JSON(c.ErrResp(commConsts.CommErr, err.Error()))
		return
	}
	ctx.JSON(c.SuccessResp(pathMap))
}

// UploadFile 上传文件
func (c *WorkspaceCtrl) UploadScripts(ctx iris.Context) {
	f, fh, err := ctx.FormFile("file")
	if err != nil {
		ctx.JSON(c.ErrResp(commConsts.CommErr, err.Error()))
		return
	}
	defer f.Close()

	err = c.WorkspaceService.UploadScripts(fh, ctx)
	if err != nil {
		ctx.JSON(c.ErrResp(commConsts.CommErr, err.Error()))
		return
	}
	ctx.JSON(c.SuccessResp(iris.Map{"workPath": commConsts.WorkDir, "sep": ztfConsts.FilePthSep}))
}
