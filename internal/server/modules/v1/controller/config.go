package controller

import (
	"fmt"

	"github.com/kataras/iris/v12"

	commConsts "github.com/easysoft/zentaoatf/internal/pkg/consts"
)

type ConfigCtrl struct {
	BaseCtrl
}

func NewConfigCtrl() *ConfigCtrl {
	return &ConfigCtrl{}
}

func (c *ConfigCtrl) SetVerbose(ctx iris.Context) {
	val, _ := ctx.URLParamBool("val")

	commConsts.Verbose = val

	ctx.JSON(c.SuccessResp(fmt.Sprintf("set verbose to %t", val)))
}
