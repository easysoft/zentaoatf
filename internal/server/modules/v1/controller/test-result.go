package controller

import (
	"fmt"
	"strings"

	commConsts "github.com/easysoft/zentaoatf/internal/pkg/consts"
	serverDomain "github.com/easysoft/zentaoatf/internal/server/modules/v1/domain"
	"github.com/easysoft/zentaoatf/internal/server/modules/v1/service"
	"github.com/kataras/iris/v12"
)

type TestResultCtrl struct {
	TestResultService *service.TestResultService `inject:""`
	StatisticService  *service.StatisticService  `inject:""`
	BaseCtrl
}

func NewTestResultCtrl() *TestResultCtrl {
	return &TestResultCtrl{}
}

// List 分页列表
func (c *TestResultCtrl) List(ctx iris.Context) {
	currSiteId, _ := ctx.URLParamInt("currSiteId")
	currProductId, _ := ctx.URLParamInt("currProductId")

	var req serverDomain.ReqPaginate
	if err := ctx.ReadQuery(&req); err != nil {
		ctx.JSON(c.ErrResp(commConsts.CommErr, err.Error()))
		return
	}

	data, err := c.TestResultService.Paginate(uint(currSiteId), uint(currProductId), req)
	if err != nil {
		ctx.JSON(c.ErrResp(commConsts.CommErr, err.Error()))
		return
	}

	ctx.JSON(c.SuccessResp(data))
}

func (c *TestResultCtrl) GetLatest(ctx iris.Context) {
	currSiteId, _ := ctx.URLParamInt("currSiteId")
	currProductId, _ := ctx.URLParamInt("currProductId")

	report, err := c.TestResultService.GetLatest(uint(currSiteId), uint(currProductId))
	if err != nil {
		ctx.JSON(c.ErrResp(commConsts.CommErr, err.Error()))
		return
	}

	ctx.JSON(c.SuccessResp(report))
}

// Get 详情
func (c *TestResultCtrl) Get(ctx iris.Context) {
	workspaceId, _ := ctx.Params().GetInt("workspaceId")
	seq := ctx.Params().Get("seq")

	if workspaceId == 0 || seq == "" {
		c.ErrResp(commConsts.ParamErr, "workspaceId and seq")
		return
	}

	exec, err := c.TestResultService.Get(workspaceId, seq)
	if err != nil {
		ctx.JSON(c.ErrResp(commConsts.CommErr, err.Error()))
		return
	}
	ctx.JSON(c.SuccessResp(exec))
}

// Delete 删除
func (c *TestResultCtrl) Delete(ctx iris.Context) {
	workspaceId, _ := ctx.URLParamInt("workspaceId")
	seq := ctx.URLParam("seq")

	if workspaceId == 0 || seq == "" {
		c.ErrResp(commConsts.ParamErr, "workspaceId and seq")
		return
	}

	err := c.TestResultService.Delete(workspaceId, seq)
	if err != nil {
		ctx.JSON(c.ErrResp(commConsts.CommErr, err.Error()))
		return
	}

	ctx.JSON(c.SuccessResp(nil))
}

// Submit 提交
func (c *TestResultCtrl) Submit(ctx iris.Context) {
	siteId, _ := ctx.URLParamInt("currSiteId")
	productId, _ := ctx.URLParamInt("currProductId")

	req := serverDomain.ZentaoResultSubmitReq{}
	err := ctx.ReadJSON(&req)
	if err != nil {
		c.ErrResp(commConsts.ParamErr, err.Error())
		return
	}

	req.ProductId = productId
	err = c.TestResultService.Submit(req, siteId, productId)
	if err != nil {
		ctx.JSON(c.ErrResp(commConsts.CommErr, err.Error()))
		return
	}

	ctx.JSON(c.SuccessResp(nil))
}

func (c *TestResultCtrl) DownloadLog(ctx iris.Context) {
	fileName := ctx.URLParamDefault("file", "")
	zipPath, err := c.TestResultService.ZipLog(fileName)
	if err != nil {
		ctx.JSON(c.ErrResp(commConsts.CommErr, err.Error()))
		return
	}

	ctx.ServeFile(zipPath)
}

func (c *TestResultCtrl) MvLog(ctx iris.Context) {
	type mvlog struct {
		File        string `json:"file"`
		WorkspaceId int    `json:"workspaceId"`
	}
	params := mvlog{}
	err := ctx.ReadJSON(&params)
	if err != nil {
		c.ErrResp(commConsts.ParamErr, err.Error())
		return
	}
	fileName := params.File
	workspaceId := params.WorkspaceId
	logPath, err := c.TestResultService.DownloadFromProxy(fileName, workspaceId)

	if err != nil {
		ctx.JSON(c.ErrResp(commConsts.CommErr, err.Error()))
		return
	}

	ctx.JSON(c.SuccessResp(logPath))
}

// Get 统计信息
func (c *TestResultCtrl) GetStatistic(ctx iris.Context) {
	scriptPath := ctx.URLParam("path")

	if scriptPath == "" {
		ctx.JSON(c.ErrResp(commConsts.ParamErr, fmt.Sprintf("参数%s不合法", "path")))
		return
	}

	if strings.Index(scriptPath, "zentao") == 0 {
		ctx.JSON(c.SuccessResp(""))
		return
	}

	statistic, err := c.StatisticService.GetByPath(scriptPath)
	if err != nil {
		ctx.JSON(c.ErrResp(commConsts.CommErr, err.Error()))
		return
	}
	ctx.JSON(c.SuccessResp(statistic))
}

// Put 统计信息
func (c *TestResultCtrl) UpdateStatistic(ctx iris.Context) {
	resultPath := ctx.URLParam("path")

	if resultPath == "" {
		ctx.JSON(c.ErrResp(commConsts.ParamErr, fmt.Sprintf("参数%s不合法", "path")))
		return
	}

	if strings.Index(resultPath, "zentao") == 0 {
		ctx.JSON(c.SuccessResp(""))
		return
	}

	scriptPaths, err := c.StatisticService.UpdateStatistic(resultPath)
	if err != nil {
		ctx.JSON(c.ErrResp(commConsts.CommErr, err.Error()))
		return
	}
	ctx.JSON(c.SuccessResp(scriptPaths))
}

// Get 失败用例信息
func (c *TestResultCtrl) GetFailure(ctx iris.Context) {
	scriptPath := ctx.URLParam("path")

	if scriptPath == "" {
		ctx.JSON(c.ErrResp(commConsts.ParamErr, fmt.Sprintf("参数%s不合法", "path")))
		return
	}

	if strings.Index(scriptPath, "zentao") == 0 {
		ctx.JSON(c.SuccessResp(""))
		return
	}

	logs, err := c.StatisticService.GetFailureLogs(scriptPath)

	if err != nil {
		ctx.JSON(c.ErrResp(commConsts.CommErr, err.Error()))
		return
	}
	ctx.JSON(c.SuccessResp(logs))
}
