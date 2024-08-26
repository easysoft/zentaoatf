package controller

import (
	"strings"

	"github.com/kataras/iris/v12"

	commConsts "github.com/easysoft/zentaoatf/internal/pkg/consts"
	i118Utils "github.com/easysoft/zentaoatf/pkg/lib/i118"
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
