package controller

import (
	commConsts "github.com/aaronchen2k/deeptest/internal/comm/consts"
	"github.com/aaronchen2k/deeptest/internal/server/modules/v1/service"
	"github.com/kataras/iris/v12"
)

type TestFilterCtrl struct {
	TestFilterService *service.TestFilterService `inject:""`
	BaseCtrl
}

func NewTestFilterCtrl() *TestFilterCtrl {
	return &TestFilterCtrl{}
}

// ListItems 获取脚本过滤器的内容列表
func (c *TestFilterCtrl) ListItems(ctx iris.Context) {
	currSiteId, _ := ctx.URLParamInt("currSiteId")
	currProductId, _ := ctx.URLParamInt("currProductId")
	filerType := ctx.URLParam("filerType")

	ret, _ := c.TestFilterService.ListFilterItems(commConsts.ScriptFilterType(filerType), uint(currSiteId), uint(currProductId))
	ctx.JSON(c.SuccessResp(ret))
}
