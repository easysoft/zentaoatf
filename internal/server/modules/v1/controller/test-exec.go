package controller

import (
	"github.com/aaronchen2k/deeptest/internal/pkg/domain"
	logUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/log"
	"github.com/aaronchen2k/deeptest/internal/server/modules/v1/service"
	"github.com/kataras/iris/v12"
)

type TestExecCtrl struct {
	TestExecService *service.TestExecService `inject:""`
	BaseCtrl
}

func NewTestExecCtrl() *TestExecCtrl {
	return &TestExecCtrl{}
}

// List 分页列表
func (c *TestExecCtrl) List(ctx iris.Context) {
	projectPath := ctx.URLParam("currProject")

	data, err := c.TestExecService.List(projectPath)
	if err != nil {
		ctx.JSON(domain.Response{Code: domain.SystemErr.Code, Data: nil, Msg: err.Error()})
		return
	}

	ctx.JSON(domain.Response{Code: domain.NoErr.Code, Data: data, Msg: domain.NoErr.Msg})
}

// Get 详情
func (c *TestExecCtrl) Get(ctx iris.Context) {
	resultPath := ctx.URLParam("resultPath")
	if resultPath == "" {
		logUtils.Errorf("参数解析失败")
		ctx.JSON(domain.Response{Code: domain.ParamErr.Code, Data: nil, Msg: domain.ParamErr.Msg})
		return
	}

	exec, err := c.TestExecService.Get(resultPath)
	if err != nil {
		ctx.JSON(domain.Response{Code: domain.SystemErr.Code, Data: nil, Msg: domain.SystemErr.Msg})
		return
	}
	ctx.JSON(domain.Response{Code: domain.NoErr.Code, Data: exec, Msg: domain.NoErr.Msg})
}

// Delete 删除
func (c *TestExecCtrl) Delete(ctx iris.Context) {
	projectPath := ctx.URLParam("currProject")
	name := ctx.URLParam("name")

	if projectPath == "" || name == "" {
		logUtils.Errorf("参数解析失败")
		ctx.JSON(domain.Response{Code: domain.ParamErr.Code, Data: nil, Msg: domain.ParamErr.Msg})
		return
	}

	err := c.TestExecService.Delete(projectPath, name)
	if err != nil {
		ctx.JSON(domain.Response{Code: domain.SystemErr.Code, Data: nil, Msg: err.Error()})
		return
	}

	ctx.JSON(domain.Response{Code: domain.NoErr.Code, Data: nil, Msg: domain.NoErr.Msg})
}
