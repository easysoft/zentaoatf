package controller

import (
	"fmt"
	commConsts "github.com/aaronchen2k/deeptest/internal/comm/consts"
	commDomain "github.com/aaronchen2k/deeptest/internal/comm/domain"
	configUtils "github.com/aaronchen2k/deeptest/internal/comm/helper/config"
	scriptUtils "github.com/aaronchen2k/deeptest/internal/comm/helper/script"
	serverDomain "github.com/aaronchen2k/deeptest/internal/server/modules/v1/domain"
	"github.com/aaronchen2k/deeptest/internal/server/modules/v1/service"
	"github.com/kataras/iris/v12"
)

type TestScriptCtrl struct {
	TestScriptService *service.TestScriptService `inject:""`
	SyncService       *service.SyncService       `inject:""`
	WorkspaceService  *service.WorkspaceService  `inject:""`
	SiteService       *service.SiteService       `inject:""`
	BaseCtrl
}

func NewTestScriptCtrl() *TestScriptCtrl {
	return &TestScriptCtrl{}
}

// List 列表
func (c *TestScriptCtrl) List(ctx iris.Context) {
	currSiteId, _ := ctx.URLParamInt("currSiteId")
	currProductId, _ := ctx.URLParamInt("currProductId")

	filerType := ctx.URLParam("filerType")
	filerValue, _ := ctx.URLParamInt("filerValue")

	testScripts, _ := c.TestScriptService.LoadTestScriptsBySiteProduct(currSiteId, currProductId, filerType, filerValue)

	ctx.JSON(c.SuccessResp(testScripts))
}

// Get 详情
func (c *TestScriptCtrl) Get(ctx iris.Context) {
	scriptPath := ctx.URLParam("path")
	workspaceId, _ := ctx.URLParamInt("workspaceId")

	if scriptPath == "" {
		ctx.JSON(c.ErrResp(commConsts.ParamErr, fmt.Sprintf("参数%s不合法", "path")))
		return
	}

	script, err := scriptUtils.GetScriptContent(scriptPath, workspaceId)
	if err != nil {
		ctx.JSON(c.ErrResp(commConsts.CommErr, err.Error()))
		return
	}
	ctx.JSON(c.SuccessResp(script))
}

// LoadCodeChildren 子节点
func (c *TestScriptCtrl) LoadCodeChildren(ctx iris.Context) {
	dir := ctx.URLParam("dir")
	workspaceId, _ := ctx.URLParamInt("workspaceId")

	testScripts, _ := c.TestScriptService.LoadCodeChildren(dir, workspaceId)

	ctx.JSON(c.SuccessResp(testScripts))
}

func (c *TestScriptCtrl) UpdateCode(ctx iris.Context) {
	req := serverDomain.TestScript{}

	err := ctx.ReadJSON(&req)
	if err != nil {
		ctx.JSON(c.ErrResp(commConsts.CommErr, err.Error()))
	}

	err = c.TestScriptService.UpdateCode(req)

	if err != nil {
		ctx.JSON(c.ErrResp(commConsts.CommErr, err.Error()))
		return
	}
	ctx.JSON(c.SuccessResp(nil))
}

// Get 详情
func (c *TestScriptCtrl) Extract(ctx iris.Context) {
	scriptPath := ctx.URLParam("path")
	workspaceId, _ := ctx.URLParamInt("workspaceId")

	if scriptPath == "" {
		ctx.JSON(c.ErrResp(commConsts.ParamErr, fmt.Sprintf("参数%s不合法", "path")))
		return
	}

	scriptUtils.Extract([]string{scriptPath})

	script, err := scriptUtils.GetScriptContent(scriptPath, workspaceId)
	if err != nil {
		ctx.JSON(c.ErrResp(commConsts.CommErr, err.Error()))
		return
	}
	ctx.JSON(c.SuccessResp(script))
}

func (c *TestScriptCtrl) SyncFromZentao(ctx iris.Context) {
	currSiteId, _ := ctx.URLParamInt("currSiteId")
	currProductId, _ := ctx.URLParamInt("currProductId")

	syncSettings := commDomain.SyncSettings{}
	err := ctx.ReadJSON(&syncSettings)
	if err != nil {
		c.ErrResp(commConsts.ParamErr, err.Error())
		return
	}

	site, _ := c.SiteService.Get(uint(currSiteId))
	config := configUtils.LoadBySite(site)
	workspace, _ := c.WorkspaceService.Get(uint(syncSettings.WorkspaceId))

	syncSettings.ProductId = currProductId
	err = c.SyncService.SyncFromZentao(syncSettings, config, workspace.Path)
	if err != nil {
		ctx.JSON(c.ErrResp(commConsts.CommErr, err.Error()))
		return
	}

	ctx.JSON(c.SuccessResp(nil))
}

func (c *TestScriptCtrl) SyncToZentao(ctx iris.Context) {
	currSiteId, _ := ctx.URLParamInt("currSiteId")
	currProductId, _ := ctx.URLParamInt("currProductId")

	if currProductId == 0 {
		ctx.JSON(c.ErrResp(commConsts.ParamErr, ""))
		return
	}

	sets := make([]serverDomain.TestSet, 0)
	err := ctx.ReadJSON(&sets)
	if err != nil {
		c.ErrResp(commConsts.ParamErr, err.Error())
		return
	}

	site, _ := c.SiteService.Get(uint(currSiteId))
	config := configUtils.LoadBySite(site)

	for _, set := range sets {
		workspaceId := set.WorkspaceId
		workspace, _ := c.WorkspaceService.Get(uint(workspaceId))

		err := c.SyncService.SyncToZentao(set.Cases, workspace.Path, currProductId, config)
		if err != nil {
			ctx.JSON(c.ErrResp(commConsts.CommErr, err.Error()))
			return
		}
	}

	ctx.JSON(c.SuccessResp(nil))
}

// Get 根据报告获取用例编号的列表
func (c *TestScriptCtrl) GetCaseIdsFromReport(ctx iris.Context) {
	workspaceId, _ := ctx.URLParamInt("workspaceId")
	seq := ctx.URLParam("seq")
	scope := ctx.URLParam("scope")

	if workspaceId == 0 || seq == "" || scope == "" {
		c.ErrResp(commConsts.ParamErr, "workspaceId, seq and scope")
		return
	}

	caseIds, err := c.TestScriptService.GetCaseIdsFromReport(workspaceId, seq, scope)
	if err != nil {
		ctx.JSON(c.ErrResp(commConsts.CommErr, err.Error()))
		return
	}
	ctx.JSON(c.SuccessResp(caseIds))
}
