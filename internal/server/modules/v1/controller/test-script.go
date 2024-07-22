package controller

import (
	"fmt"
	"strings"

	serverConfig "github.com/easysoft/zentaoatf/internal/server/config"

	commConsts "github.com/easysoft/zentaoatf/internal/pkg/consts"
	commDomain "github.com/easysoft/zentaoatf/internal/pkg/domain"
	configHelper "github.com/easysoft/zentaoatf/internal/pkg/helper/config"
	scriptHelper "github.com/easysoft/zentaoatf/internal/pkg/helper/script"
	zentaoHelper "github.com/easysoft/zentaoatf/internal/pkg/helper/zentao"
	serverDomain "github.com/easysoft/zentaoatf/internal/server/modules/v1/domain"
	"github.com/easysoft/zentaoatf/internal/server/modules/v1/service"
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

	displayBy := ctx.URLParam("displayBy")

	testScripts, _ := c.TestScriptService.LoadTestScriptsBySiteProduct(uint(currSiteId), uint(currProductId),
		displayBy, filerType, filerValue)

	ctx.JSON(c.SuccessResp(testScripts))
}

// LoadCodeChildren 子节点
func (c *TestScriptCtrl) LoadCodeChildren(ctx iris.Context) {
	dir := ctx.URLParam("dir")
	workspaceId, _ := ctx.URLParamInt("workspaceId")

	testScripts, _ := c.TestScriptService.LoadCodeChildren(dir, workspaceId)

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

func (c *TestScriptCtrl) Create(ctx iris.Context) {
	req := serverDomain.CreateScriptReq{}
	err := ctx.ReadJSON(&req)
	if err != nil {
		ctx.JSON(c.ErrResp(commConsts.CommErr, err.Error()))
	}

	pth, bizErr := c.TestScriptService.CreateNode(req)
	if bizErr != nil {
		ctx.JSON(c.BizErrResp(bizErr, ""))
		return
	}

	ctx.JSON(c.SuccessResp(pth))
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

func (c *TestScriptCtrl) UpdateName(ctx iris.Context) {
	req := serverDomain.TestScript{}

	err := ctx.ReadJSON(&req)
	if err != nil {
		ctx.JSON(c.ErrResp(commConsts.CommErr, err.Error()))
	}

	err = c.TestScriptService.UpdateName(req)

	if err != nil {
		ctx.JSON(c.ErrResp(commConsts.CommErr, err.Error()))
		return
	}
	ctx.JSON(c.SuccessResp(nil))
}

func (c *TestScriptCtrl) Paste(ctx iris.Context) {
	req := serverDomain.PasteScriptReq{}

	err := ctx.ReadJSON(&req)
	if err != nil {
		ctx.JSON(c.ErrResp(commConsts.CommErr, err.Error()))
	}

	err = c.TestScriptService.Paste(req)
	if err != nil {
		ctx.JSON(c.ErrResp(commConsts.CommErr, err.Error()))
		return
	}

	ctx.JSON(c.SuccessResp(nil))
}

func (c *TestScriptCtrl) Move(ctx iris.Context) {
	req := serverDomain.MoveScriptReq{}

	err := ctx.ReadJSON(&req)
	if err != nil {
		ctx.JSON(c.ErrResp(commConsts.CommErr, err.Error()))
	}

	err = c.TestScriptService.Move(req)
	if err != nil {
		ctx.JSON(c.ErrResp(commConsts.CommErr, err.Error()))
		return
	}

	ctx.JSON(c.SuccessResp(nil))
}

func (c *TestScriptCtrl) Delete(ctx iris.Context) {
	req := serverDomain.TestScript{}

	err := ctx.ReadJSON(&req)
	if err != nil {
		ctx.JSON(c.ErrResp(commConsts.CommErr, err.Error()))
		return
	}

	bizErr := c.TestScriptService.Delete(req.Path)
	if bizErr != nil {
		ctx.JSON(c.BizErrResp(bizErr, ""))
		return
	}

	ctx.JSON(c.SuccessResp(nil))
}

func (c *TestScriptCtrl) Rename(ctx iris.Context) {
	req := serverDomain.TestScript{}

	err := ctx.ReadJSON(&req)
	if err != nil {
		ctx.JSON(c.ErrResp(commConsts.CommErr, err.Error()))
		return
	}

	bizErr := c.TestScriptService.Rename(req.Path, req.Name)
	if bizErr != nil {
		ctx.JSON(c.BizErrResp(bizErr, ""))
		return
	}

	ctx.JSON(c.SuccessResp(nil))
}

// Extract 详情
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
	pths, err := zentaoHelper.Checkout(syncSettings, config, workspace.Path)
	if err != nil {
		ctx.JSON(c.ErrResp(commConsts.CommErr, err.Error()))
		return
	}

	ctx.JSON(c.SuccessResp(pths))
}

func (c *TestScriptCtrl) SyncToZentao(ctx iris.Context) {
	currSiteId, _ := ctx.URLParamInt("currSiteId")
	currProductId := ctx.URLParam("currProductId")

	if len(currProductId) == 0 {
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

	totalNum := 0
	successNum := 0
	for _, set := range sets {
		totalNum += len(set.Cases)

		count, _ := zentaoHelper.CheckIn(currProductId, set.Cases, config, true, true)

		successNum += count
	}

	ctx.JSON(c.SuccessResp(iris.Map{
		"total":   totalNum,
		"success": successNum,
	}))
}

func (c *TestScriptCtrl) SyncDirToZentao(ctx iris.Context) {
	dir := ctx.URLParam("dir")
	productId := ctx.URLParam("productId")

	cases := scriptHelper.GetCaseByDirAndFile([]string{dir})

	config := commDomain.WorkspaceConf{
		Url: serverConfig.CONFIG.Server,
	}

	// TODO 产品ID 新增
	zentaoHelper.CheckIn(productId, cases, config, true, true)

	return
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
