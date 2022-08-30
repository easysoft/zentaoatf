package controller

import (
	"strconv"

	commConsts "github.com/easysoft/zentaoatf/internal/pkg/consts"
	commDomain "github.com/easysoft/zentaoatf/internal/pkg/domain"
	zentaoHelper "github.com/easysoft/zentaoatf/internal/pkg/helper/zentao"
	"github.com/easysoft/zentaoatf/internal/server/modules/v1/service"
	"github.com/kataras/iris/v12"
)

type TestBugCtrl struct {
	TestBugService   *service.TestBugService   `inject:""`
	WorkspaceService *service.WorkspaceService `inject:""`
	SiteService      *service.SiteService      `inject:""`
	BaseCtrl
}

func NewTestBugCtrl() *TestBugCtrl {
	return &TestBugCtrl{}
}

func (c *TestBugCtrl) PrepareBugData(ctx iris.Context) {
	req := commDomain.FuncResult{}
	err := ctx.ReadJSON(&req)
	if err != nil {
		ctx.JSON(c.ErrResp(commConsts.ParamErr, err.Error()))
		return
	}

	workspace, _ := c.WorkspaceService.Get(uint(req.WorkspaceId))
	bug := zentaoHelper.PrepareBug(workspace.Path, req.Seq, strconv.Itoa(req.Id), 0)

	ctx.JSON(c.SuccessResp(bug))
}

// Submit 提交
func (c *TestBugCtrl) Submit(ctx iris.Context) {
	siteId, _ := ctx.URLParamInt("currSiteId")
	productId, _ := ctx.URLParamInt("currProductId")

	req := commDomain.ZtfBug{}
	if err := ctx.ReadJSON(&req); err != nil {
		ctx.JSON(c.ErrResp(commConsts.CommErr, err.Error()))
		return
	}

	err := c.TestBugService.Submit(req, siteId, productId)
	if err != nil {
		ctx.JSON(c.ErrResp(commConsts.CommErr, err.Error()))
		return
	}

	ctx.JSON(c.SuccessResp(nil))
}

//查询产品下所有bug
func (c *TestBugCtrl) LoadBugs(ctx iris.Context) {
	siteId, _ := ctx.URLParamInt("currSiteId")
	productId, _ := ctx.URLParamInt("currProductId")

	bugs, _ := c.TestBugService.LoadBugs(siteId, productId)

	ctx.JSON(c.SuccessResp(bugs))
}
