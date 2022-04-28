package controller

import (
	"fmt"
	commConsts "github.com/easysoft/zentaoatf/internal/comm/consts"
	commDomain "github.com/easysoft/zentaoatf/internal/comm/domain"
	configHelper "github.com/easysoft/zentaoatf/internal/comm/helper/config"
	scriptHelper "github.com/easysoft/zentaoatf/internal/comm/helper/script"
	zentaoHelper "github.com/easysoft/zentaoatf/internal/comm/helper/zentao"
	serverDomain "github.com/easysoft/zentaoatf/internal/server/modules/v1/domain"
	"github.com/easysoft/zentaoatf/internal/server/modules/v1/service"
	"github.com/kataras/iris/v12"
	"strings"
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

	testScripts, _ := c.TestScriptService.LoadTestScriptsBySiteProduct(uint(currSiteId), uint(currProductId),
		filerType, filerValue)

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

	if strings.Index(scriptPath, "zentao") == 0 {
		ctx.JSON(c.SuccessResp(""))
		return
	}

	script, err := scriptHelper.GetScriptContent(scriptPath, workspaceId)
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

	done, _ := scriptHelper.Extract([]string{scriptPath})

	script, err := scriptHelper.GetScriptContent(scriptPath, workspaceId)
	if err != nil {
		ctx.JSON(c.ErrResp(commConsts.CommErr, err.Error()))
		return
	}

	ret := iris.Map{"script": script, "done": done}

	ctx.JSON(c.SuccessResp(ret))
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
	config := configHelper.LoadBySite(site)
	workspace, _ := c.WorkspaceService.Get(uint(syncSettings.WorkspaceId))

	syncSettings.ProductId = currProductId
	if syncSettings.Lang == "" {
		syncSettings.Lang = workspace.Lang
	}
	pths, err := zentaoHelper.SyncFromZentao(syncSettings, config, workspace.Path)
	if err != nil {
		ctx.JSON(c.ErrResp(commConsts.CommErr, err.Error()))
		return
	}

	ctx.JSON(c.SuccessResp(pths))
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
	config := configHelper.LoadBySite(site)

	for _, set := range sets {
		//	workspaceId := set.WorkspaceId
		//	workspace, _ := c.WorkspaceService.Get(uint(workspaceId))

		err := zentaoHelper.SyncToZentao(set.Cases, config)
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
