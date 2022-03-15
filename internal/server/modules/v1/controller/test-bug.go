package controller

import (
	commConsts "github.com/aaronchen2k/deeptest/internal/comm/consts"
	commDomain "github.com/aaronchen2k/deeptest/internal/comm/domain"
	zentaoUtils "github.com/aaronchen2k/deeptest/internal/comm/helper/zentao"
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
	workspacePath := ctx.URLParam("currWorkspace")

	req := commDomain.FuncResult{}
	err := ctx.ReadJSON(&req)
	if err != nil {
		ctx.JSON(c.ErrResp(commConsts.ParamErr, err.Error()))
		return
	}

	bug := zentaoUtils.PrepareBug(workspacePath, req.Seq, strconv.Itoa(req.Id))

	ctx.JSON(c.SuccessResp(bug))
}

// Submit 提交
func (c *TestBugCtrl) Submit(ctx iris.Context) {
	workspacePath := ctx.URLParam("currWorkspace")
	req := commDomain.ZtfBug{}
	if err := ctx.ReadJSON(&req); err != nil {
		ctx.JSON(c.ErrResp(commConsts.CommErr, err.Error()))
		return
	}

	err := c.TestBugService.Submit(req, workspacePath)
	if err != nil {
		ctx.JSON(c.ErrResp(commConsts.CommErr, err.Error()))
		return
	}

	ctx.JSON(c.SuccessResp(nil))
}
