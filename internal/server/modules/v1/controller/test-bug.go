package controller

import (
	commDomain "github.com/aaronchen2k/deeptest/internal/comm/domain"
	"github.com/aaronchen2k/deeptest/internal/pkg/domain"
	logUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/log"
	"github.com/aaronchen2k/deeptest/internal/server/modules/v1/service"
	"github.com/kataras/iris/v12"
)

type TestBugCtrl struct {
	TestBugService *service.TestBugService `inject:""`
	BaseCtrl
}

func NewTestBugCtrl() *TestBugCtrl {
	return &TestBugCtrl{}
}

// Gen 生成缺陷
func (c *TestBugCtrl) Gen(ctx iris.Context) {
	req := commDomain.ZtfBug{}
	if err := ctx.ReadJSON(&req); err != nil {
		logUtils.Errorf("参数验证失败，错误%s", err.Error())
		ctx.JSON(domain.Response{Code: domain.SystemErr.Code, Data: nil, Msg: err.Error()})
		return
	}

	err := c.TestBugService.Submit(req)
	if err != nil {
		ctx.JSON(domain.Response{
			Code: c.ErrCode(err),
			Data: nil,
		})
		return
	}

	ctx.JSON(domain.Response{Code: domain.NoErr.Code, Data: nil, Msg: domain.NoErr.Msg})
}

// Submit 提交
func (c *TestBugCtrl) Submit(ctx iris.Context) {
	req := commDomain.ZtfBug{}
	if err := ctx.ReadJSON(&req); err != nil {
		logUtils.Errorf("参数验证失败，错误%s", err.Error())
		ctx.JSON(domain.Response{Code: domain.SystemErr.Code, Data: nil, Msg: err.Error()})
		return
	}

	err := c.TestBugService.Submit(req)
	if err != nil {
		ctx.JSON(domain.Response{
			Code: c.ErrCode(err),
			Data: nil,
		})
		return
	}

	ctx.JSON(domain.Response{Code: domain.NoErr.Code, Data: nil, Msg: domain.NoErr.Msg})
}
