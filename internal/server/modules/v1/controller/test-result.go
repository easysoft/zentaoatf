package controller

import (
	commConsts "github.com/aaronchen2k/deeptest/internal/comm/consts"
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
		c.ErrResp(commConsts.ParamErr, err.Error())
		return
	}

	err = c.TestResultService.Submit(req, projectPath)
	if err != nil {
		ctx.JSON(c.ErrResp(commConsts.CommErr, err.Error()))
		return
	}

	ctx.JSON(c.SuccessResp(nil))
}
