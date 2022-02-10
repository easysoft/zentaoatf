package controller

import (
	commDomain "github.com/aaronchen2k/deeptest/internal/comm/domain"
	zentaoUtils "github.com/aaronchen2k/deeptest/internal/comm/helper/zentao"
	"github.com/aaronchen2k/deeptest/internal/pkg/domain"
	logUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/log"
	"github.com/aaronchen2k/deeptest/internal/server/modules/v1/service"
	"github.com/kataras/iris/v12"
	"strconv"
)

type TestBugCtrl struct {
	TestBugService *service.TestBugService `inject:""`
	BaseCtrl
}

func NewTestBugCtrl() *TestBugCtrl {
	return &TestBugCtrl{}
}

func (c *TestBugCtrl) GetBugData(ctx iris.Context) {
	projectPath := ctx.URLParam("currProject")

	req := commDomain.FuncResult{}
	err := ctx.ReadJSON(&req)
	if err != nil {
		logUtils.Errorf("参数验证失败 %s", err.Error())
		ctx.JSON(domain.Response{Code: domain.RequestErr.Code, Data: nil, Msg: err.Error()})
		return
	}

	bug := zentaoUtils.PrepareBug(projectPath, req.Seq, strconv.Itoa(req.Id))

	if err != nil {
		ctx.JSON(domain.Response{Code: domain.RequestErr.Code, Data: nil, Msg: "获取禅道缺陷步骤失败"})
		return
	}

	ctx.JSON(domain.Response{Code: domain.NoErr.Code, Data: bug, Msg: domain.NoErr.Msg})
}

// Submit 提交
func (c *TestBugCtrl) Submit(ctx iris.Context) {
	projectPath := ctx.URLParam("currProject")
	req := commDomain.ZtfBug{}
	if err := ctx.ReadJSON(&req); err != nil {
		logUtils.Errorf("参数验证失败，错误%s", err.Error())
		ctx.JSON(domain.Response{Code: domain.RequestErr.Code, Data: nil, Msg: err.Error()})
		return
	}

	err := c.TestBugService.Submit(req, projectPath)
	if err != nil {
		ctx.JSON(domain.Response{Code: domain.RequestErr.Code, Data: nil, Msg: err.Error()})
		return
	}

	ctx.JSON(domain.Response{Code: domain.NoErr.Code, Data: nil, Msg: domain.NoErr.Msg})
}
