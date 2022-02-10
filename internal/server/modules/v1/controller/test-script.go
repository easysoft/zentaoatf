package controller

import (
	scriptUtils "github.com/aaronchen2k/deeptest/internal/comm/helper/script"
	"github.com/aaronchen2k/deeptest/internal/pkg/domain"
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

// Get 详情
func (c *TestScriptCtrl) Get(ctx iris.Context) {
	scriptPath := ctx.URLParam("path")
	if scriptPath == "" {
		ctx.JSON(domain.Response{Code: domain.ParamErr.Code, Data: nil, Msg: "参数解析失败"})
		return
	}

	script, err := scriptUtils.GetScriptContent(scriptPath)
	if err != nil {
		ctx.JSON(domain.Response{Code: domain.SystemErr.Code, Data: nil, Msg: domain.SystemErr.Msg})
		return
	}
	ctx.JSON(domain.Response{Code: domain.NoErr.Code, Data: script, Msg: domain.NoErr.Msg})
}

// Get 详情
func (c *TestScriptCtrl) Extract(ctx iris.Context) {
	scriptPath := ctx.URLParam("path")
	if scriptPath == "" {
		ctx.JSON(domain.Response{Code: domain.ParamErr.Code, Data: nil, Msg: "参数解析失败"})
		return
	}

	scriptUtils.Extract([]string{scriptPath})

	script, err := scriptUtils.GetScriptContent(scriptPath)
	if err != nil {
		ctx.JSON(domain.Response{Code: domain.SystemErr.Code, Data: nil, Msg: domain.SystemErr.Msg})
		return
	}
	ctx.JSON(domain.Response{Code: domain.NoErr.Code, Data: script, Msg: domain.NoErr.Msg})
}
