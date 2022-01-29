package controller

import (
	"github.com/aaronchen2k/deeptest/internal/pkg/domain"
	logUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/log"
	serverDomain "github.com/aaronchen2k/deeptest/internal/server/modules/v1/domain"
	"github.com/aaronchen2k/deeptest/internal/server/modules/v1/service"
	"github.com/kataras/iris/v12"
)

type TestResultCtrl struct {
	TestResultService *service.TestResultService `inject:""`
	BaseCtrl
}

func NewTestResultCtrl() *TestResultCtrl {
	return &TestResultCtrl{}
}

// Submit 提交
func (c *TestResultCtrl) Submit(ctx iris.Context) {
	projectPath := ctx.URLParam("currProject")
	req := serverDomain.ZentaoResultSubmitReq{}
	err := ctx.ReadJSON(&req)
	if err != nil {
		logUtils.Errorf("参数验证失败，错误%s", err.Error())
		ctx.JSON(domain.Response{Code: domain.ParamErr.Code, Data: nil, Msg: err.Error()})
		return
	}

	err = c.TestResultService.Submit(req, projectPath)
	if err != nil {
		ctx.JSON(domain.Response{Code: domain.CommonErr.Code, Msg: err.Error()})
		return
	}

	ctx.JSON(domain.Response{Code: domain.NoErr.Code, Data: nil, Msg: domain.NoErr.Msg})
}
