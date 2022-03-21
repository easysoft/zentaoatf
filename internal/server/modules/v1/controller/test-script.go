package controller

import (
	"fmt"
	commConsts "github.com/aaronchen2k/deeptest/internal/comm/consts"
	scriptUtils "github.com/aaronchen2k/deeptest/internal/comm/helper/script"
	"github.com/aaronchen2k/deeptest/internal/server/modules/v1/service"
	"github.com/kataras/iris/v12"
)

type TestScriptCtrl struct {
	TestScriptService *service.TestScriptService `inject:""`
	BaseCtrl
}

func NewTestScriptCtrl() *TestScriptCtrl {
	return &TestScriptCtrl{}
}

// List 列表
func (c *TestScriptCtrl) List(ctx iris.Context) {
	currSiteId, _ := ctx.URLParamInt("currSiteId")
	currProductId, _ := ctx.URLParamInt("currProductId")
	workspaceId, _ := ctx.URLParamInt("workspaceId")

	testScripts, _ := c.TestScriptService.LoadTestScriptsBySiteProduct(currSiteId, currProductId, workspaceId)

	ctx.JSON(c.SuccessResp(testScripts))
}

// Get 详情
func (c *TestScriptCtrl) Get(ctx iris.Context) {
	scriptPath := ctx.URLParam("path")
	if scriptPath == "" {
		ctx.JSON(c.ErrResp(commConsts.ParamErr, fmt.Sprintf("参数%s不合法", "path")))
		return
	}

	script, err := scriptUtils.GetScriptContent(scriptPath)
	if err != nil {
		ctx.JSON(c.ErrResp(commConsts.CommErr, err.Error()))
		return
	}
	ctx.JSON(c.SuccessResp(script))
}

// Get 详情
func (c *TestScriptCtrl) Extract(ctx iris.Context) {
	scriptPath := ctx.URLParam("path")

	if scriptPath == "" {
		ctx.JSON(c.ErrResp(commConsts.ParamErr, fmt.Sprintf("参数%s不合法", "path")))
		return
	}

	scriptUtils.Extract([]string{scriptPath})

	script, err := scriptUtils.GetScriptContent(scriptPath)
	if err != nil {
		ctx.JSON(c.ErrResp(commConsts.CommErr, err.Error()))
		return
	}
	ctx.JSON(c.SuccessResp(script))
}
