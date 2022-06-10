package controller

import (
	commConsts "github.com/easysoft/zentaoatf/internal/comm/consts"
	i118Utils "github.com/easysoft/zentaoatf/internal/pkg/lib/i118"
	"github.com/kataras/iris/v12"
	"strings"
)

type SettingsCtrl struct {
	BaseCtrl
}

func NewSettingsCtrl() *SettingsCtrl {
	return &SettingsCtrl{}
}

func (c *SettingsCtrl) SetLang(ctx iris.Context) {
	lang := strings.ToLower(ctx.URLParam("lang"))

	if strings.Index(lang, "en") > -1 {
		lang = "en"
	} else {
		lang = "zh"
	}

	if lang != "" && lang != commConsts.Language {
		commConsts.Language = lang
		i118Utils.Init(lang, commConsts.AppServer)
	}

	ctx.JSON(c.SuccessResp(nil))
}
