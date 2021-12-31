package controller

import (
	"github.com/kataras/iris/v12"
	"github.com/aaronchen2k/deeptest/internal/pkg/domain"
	logUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/log"
)

type TestCtrl struct {
}

func NewTestCtrl() *TestCtrl {
	return &TestCtrl{}
}

func (c *TestCtrl) Test(ctx iris.Context) {
	logUtils.Warn("测试")

	ctx.JSON(domain.Response{Code: domain.NoErr.Code, Data: nil, Msg: domain.NoErr.Msg})
}
